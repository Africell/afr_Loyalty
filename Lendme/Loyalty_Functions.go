package Lendme

import (
	"afr_ao_apgw_v2/APGWClientV2"
	apgw "afr_ao_apgw_v2/afr_apgw"
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"mongox"
	"net/http"
	"net/url"
	"redisx"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	SpinAndWin "afr_SpinAndWin_be/SpinAndWinClient"
	PropC "afr_propylaea/PropylaeaClient"
	"afr_propylaea/propylaea"
	"afr_unified_charging_gateway/Unified_charging_gateway_Client"

	MM "afr_sb_mm"

	"github.com/jinzhu/copier"
	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// mongox repositories â€” main MongoDB (for AccessEntry and NotificationLog)
var Mdb_AccessEntry *mongox.Repository

// mongox repositories â€” loyalty MongoDB
var Mdb_Loyalty_AutoIncrement *mongox.Repository
var Mdb_Loyalty_Governance *mongox.Repository
var Mdb_Loyalty_Governance_log *mongox.Repository
var Mdb_Loyalty_Level *mongox.Repository
var Mdb_Loyalty_Level_Change_log *mongox.Repository
var Mdb_Loyalty_Seniority_Level *mongox.Repository
var Mdb_Loyalty_Account_Segment *mongox.Repository
var Mdb_Loyalty_Point_Earning_Rules *mongox.Repository
var Mdb_Loyalty_Point_Earning_Rules_Overwrite *mongox.Repository
var Mdb_Loyalty_Point_Expiry_Rules *mongox.Repository
var Mdb_Loyalty_Point_Redemption_Rules *mongox.Repository
var Mdb_Loyalty_Plan *mongox.Repository
var Mdb_Customer_Loyalty_Account *mongox.Repository
var Mdb_Churned_Customer_Loyalty_Account *mongox.Repository
var Mdb_Customer_Loyalty_Account_Points_Detail *mongox.Repository
var Mdb_Customer_DND *mongox.Repository
var Mdb_Customer_Exclusion *mongox.Repository
var Mdb_Customer_COS_Exclusion *mongox.Repository
var Mdb_Customer_UAT *mongox.Repository
var Mdb_Loyalty_Event_Log *mongox.Repository
var Mdb_Loyalty_Expiry_log *mongox.Repository
var Mdb_Loyalty_Full_Expiry_log *mongox.Repository
var Mdb_Loyalty_Redemption_log *mongox.Repository
var Mdb_Loyalty_Status_log *mongox.Repository
var Mdb_Loyalty_AccountCreditPoints_log *mongox.Repository
var Mdb_Loyalty_AccountDebitPoints_log *mongox.Repository
var Mdb_NotificationLog *mongox.Repository
var Mdb_Loyalty_Campaign *mongox.Repository
var Mdb_Loyalty_Campaign_Target_List *mongox.Repository
var Mdb_Loyalty_Campaign_Account *mongox.Repository

// redisx client (shared)
var RedisClient *redisx.Client

func (e Loyalty_Governance) RedisKey() string        				{ return "Loyalty_Governance:" + e.Key }
func (e Loyalty_Level) RedisKey() string               				{ return "Loyalty_Level:" + e.Key }
func (e Loyalty_Seniority_Level) RedisKey() string     				{ return "Loyalty_Seniority_Level:" + e.Key }
func (e Loyalty_Account_Segment) RedisKey() string    				{ return "Loyalty_Account_Segment:" + e.Key }
func (e Loyalty_Point_Earning_Rules) RedisKey() string 				{ return "Loyalty_Point_Earning_Rules:" + e.Key }
func (e Loyalty_Point_Earning_Rules_Overwrite) RedisKey() string 	{ return "Loyalty_Point_Earning_Rules_Overwrite:" + e.Key }
func (e Loyalty_Point_Expiry_Rules) RedisKey() string 				{ return "Loyalty_Point_Expiry_Rules:" + e.Key }
func (e Loyalty_Point_Redemption_Rules) RedisKey() string 			{ return "Loyalty_Point_Redemption_Rules:" + e.Key }
func (e Loyalty_Plan) RedisKey() string             				{ return "Loyalty_Plan:" + e.Key }
func (e Customer_Loyalty_Account) RedisKey() string 				{ return "Customer_Loyalty_Account:" + e.Key }
func (e Customer_Loyalty_Account_Points_Detail) RedisKey() string 	{ return "Customer_Loyalty_Account_Points_Detail:" + e.Key}
func (e Customer_DND) RedisKey() string           					{ return "Customer_DND:" + e.Key }
func (e Customer_Exclusion) RedisKey() string     					{ return "Customer_Exclusion:" + e.Key }
func (e Customer_COS_Exclusion) RedisKey() string 					{ return "Customer_COS_Exclusion:" + e.Key }
func (e Customer_UAT) RedisKey() string           					{ return "Customer_UAT:" + e.Key }
func (e Loyalty_Campaign) RedisKey() string      					{ return "Loyalty_Campaign:" + e.Key }
func (e Loyalty_Campaign_Target_List) RedisKey() string 			{ return "Loyalty_Campaign_Target_List:" + e.Key }
func (e Loyalty_Campaign_Account) RedisKey() string 				{ return "Loyalty_Campaign_Account:" + e.Key }

var chan_LoyaltyGovernance_Controler = make(chan int, 1)

var chan_PointsExpiry_Controler = make(chan int, 50)

var LOYALTY_GOVERNANCE_KEY string = "Loyalty_Governance"

// resetDailyWeeklyCountersIfNeeded resets daily and weekly counters when the
// calendar day or ISO week has rolled over since the last reset.
func resetDailyWeeklyCountersIfNeeded(acc *Customer_Loyalty_Account) {
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// ISO week starts Monday (weekday 1). Sunday maps to 7.
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	startOfWeek := today.AddDate(0, 0, -(weekday - 1))

	if acc.Daily_Reset_Date.Before(today) {
		acc.Daily_Earned_Points = 0
		acc.Daily_Redeemed_Points = 0
		acc.Daily_Redemption_Attempts = 0
		acc.Daily_Reset_Date = today
	}
	if acc.Weekly_Reset_Date.Before(startOfWeek) {
		acc.Weekly_Earned_Points = 0
		acc.Weekly_Redeemed_Points = 0
		acc.Weekly_Redemption_Attempts = 0
		acc.Weekly_Reset_Date = startOfWeek
	}
}

var processed = make(map[string]struct{})
var processedMu sync.Mutex

var jobs = make(map[string]*JobStatus)
var jobsMu sync.Mutex

func (UC *UserControl) InitializeMongoxRepositories() error {
	// main MongoDB (for AccessEntry and NotificationLog)
	mainDB, err := mongox.NewDB(UC.MongoClient.Mongo, Configuration.DB_Name_Loyalty, 5*time.Second)
	if err != nil {
		return err
	}
	if Mdb_AccessEntry, err = mongox.NewRepository(mainDB, "Col_AccessEntry"); err != nil {
		return err
	}
	if Mdb_NotificationLog, err = mongox.NewRepository(mainDB, "Col_NotificationLog"); err != nil {
		return err
	}

	// loyalty MongoDB
	loyaltyDB, err := mongox.NewDB(UC.LoyaltyMongoClient.Mongo, Configuration.DB_Name_Loyalty, 5*time.Second)
	if err != nil {
		return err
	}
	if Mdb_Loyalty_AutoIncrement, err = mongox.NewRepository(loyaltyDB, "Col_Loyalty_AutoIncrement"); err != nil {
		return err
	}
	if Mdb_Loyalty_Governance, err = mongox.NewRepository(loyaltyDB, "Col_Loyalty_Governance"); err != nil {
		return err
	}
	if Mdb_Loyalty_Governance_log, err = mongox.NewRepository(loyaltyDB, "Loyalty_Governance_log"); err != nil {
		return err
	}
	if Mdb_Loyalty_Level, err = mongox.NewRepository(loyaltyDB, "Col_Loyalty_Level"); err != nil {
		return err
	}
	if Mdb_Loyalty_Level_Change_log, err = mongox.NewRepository(loyaltyDB, "Col_Loyalty_Level_Change_log"); err != nil {
		return err
	}
	if Mdb_Loyalty_Seniority_Level, err = mongox.NewRepository(loyaltyDB, "Col_Loyalty_Seniority_Level"); err != nil {
		return err
	}
	if Mdb_Loyalty_Account_Segment, err = mongox.NewRepository(loyaltyDB, "Col_Loyalty_Account_Segment"); err != nil {
		return err
	}
	if Mdb_Loyalty_Point_Earning_Rules, err = mongox.NewRepository(loyaltyDB, "Col_Loyalty_Point_Earning_Rules"); err != nil {
		return err
	}
	if Mdb_Loyalty_Point_Earning_Rules_Overwrite, err = mongox.NewRepository(loyaltyDB, "Col_Loyalty_Point_Earning_Rules_Overwrite"); err != nil {
		return err
	}
	if Mdb_Loyalty_Point_Expiry_Rules, err = mongox.NewRepository(loyaltyDB, "Col_Loyalty_Point_Expiry_Rules"); err != nil {
		return err
	}
	if Mdb_Loyalty_Point_Redemption_Rules, err = mongox.NewRepository(loyaltyDB, "Col_Loyalty_Point_Redemption_Rules"); err != nil {
		return err
	}
	if Mdb_Loyalty_Plan, err = mongox.NewRepository(loyaltyDB, "Col_Loyalty_Plan"); err != nil {
		return err
	}
	if Mdb_Customer_Loyalty_Account, err = mongox.NewRepository(loyaltyDB, "Col_Customer_Loyalty_Account"); err != nil {
		return err
	}
	if Mdb_Churned_Customer_Loyalty_Account, err = mongox.NewRepository(loyaltyDB, "Col_Churned_Customer_Loyalty_Account"); err != nil {
		return err
	}
	if Mdb_Customer_Loyalty_Account_Points_Detail, err = mongox.NewRepository(loyaltyDB, "Col_Customer_Loyalty_Account_Points_Detail"); err != nil {
		return err
	}
	if Mdb_Customer_DND, err = mongox.NewRepository(loyaltyDB, "Col_Customer_DND"); err != nil {
		return err
	}
	if Mdb_Customer_Exclusion, err = mongox.NewRepository(loyaltyDB, "Col_Customer_Exclusion"); err != nil {
		return err
	}
	if Mdb_Customer_COS_Exclusion, err = mongox.NewRepository(loyaltyDB, "Col_Customer_COS_Exclusion"); err != nil {
		return err
	}
	if Mdb_Customer_UAT, err = mongox.NewRepository(loyaltyDB, "Col_Customer_UAT"); err != nil {
		return err
	}
	if Mdb_Loyalty_Event_Log, err = mongox.NewRepository(loyaltyDB, "Col_Loyalty_Event_Log"); err != nil {
		return err
	}
	if Mdb_Loyalty_Expiry_log, err = mongox.NewRepository(loyaltyDB, "Col_Loyalty_Expiry_log"); err != nil {
		return err
	}
	if Mdb_Loyalty_Full_Expiry_log, err = mongox.NewRepository(loyaltyDB, "Col_Loyalty_Full_Expiry_log"); err != nil {
		return err
	}
	if Mdb_Loyalty_Redemption_log, err = mongox.NewRepository(loyaltyDB, "Col_Loyalty_Redemption_log"); err != nil {
		return err
	}
	if Mdb_Loyalty_Status_log, err = mongox.NewRepository(loyaltyDB, "Col_Loyalty_Status_log"); err != nil {
		return err
	}
	if Mdb_Loyalty_AccountCreditPoints_log, err = mongox.NewRepository(loyaltyDB, "Col_Loyalty_AccountCreditPoints_log"); err != nil {
		return err
	}
	if Mdb_Loyalty_AccountDebitPoints_log, err = mongox.NewRepository(loyaltyDB, "Col_Loyalty_AccountDebitPoints_log"); err != nil {
		return err
	}
	if Mdb_Loyalty_Campaign, err = mongox.NewRepository(loyaltyDB, "Col_Loyalty_Campaign"); err != nil {
		return err
	}
	if Mdb_Loyalty_Campaign_Target_List, err = mongox.NewRepository(loyaltyDB, "Col_Loyalty_Campaign_Target_List"); err != nil {
		return err
	}
	if Mdb_Loyalty_Campaign_Account, err = mongox.NewRepository(loyaltyDB, "Col_Loyalty_Campaign_Account"); err != nil {
		return err
	}

	RedisClient = UC.Redis
	return nil
}

func (uc *UserControl) RedisDataLoader() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if _, _, err := redisx.LoadMongoToRedis[Loyalty_Governance](ctx, RedisClient, Mdb_Loyalty_Governance.Coll, redisx.MongoLoadOptions{
		BatchSize: 2000, FlushBeforeLoad: true, FlushPattern: "Loyalty_Governance:*", UseUnlink: true,
	}); err != nil {
		return fmt.Errorf("load Loyalty_Governance: %w", err)
	}
	if _, _, err := redisx.LoadMongoToRedis[Loyalty_Level](ctx, RedisClient, Mdb_Loyalty_Level.Coll, redisx.MongoLoadOptions{
		BatchSize: 2000, FlushBeforeLoad: true, FlushPattern: "Loyalty_Level:*", UseUnlink: true,
	}); err != nil {
		return fmt.Errorf("load Loyalty_Level: %w", err)
	}
	if _, _, err := redisx.LoadMongoToRedis[Loyalty_Seniority_Level](ctx, RedisClient, Mdb_Loyalty_Seniority_Level.Coll, redisx.MongoLoadOptions{
		BatchSize: 2000, FlushBeforeLoad: true, FlushPattern: "Loyalty_Seniority_Level:*", UseUnlink: true,
	}); err != nil {
		return fmt.Errorf("load Loyalty_Seniority_Level: %w", err)
	}
	if _, _, err := redisx.LoadMongoToRedis[Loyalty_Account_Segment](ctx, RedisClient, Mdb_Loyalty_Account_Segment.Coll, redisx.MongoLoadOptions{
		BatchSize: 2000, FlushBeforeLoad: true, FlushPattern: "Loyalty_Account_Segment:*", UseUnlink: true,
	}); err != nil {
		return fmt.Errorf("load Loyalty_Account_Segment: %w", err)
	}
	if _, _, err := redisx.LoadMongoToRedis[Loyalty_Point_Earning_Rules](ctx, RedisClient, Mdb_Loyalty_Point_Earning_Rules.Coll, redisx.MongoLoadOptions{
		BatchSize: 2000, FlushBeforeLoad: true, FlushPattern: "Loyalty_Point_Earning_Rules:*", UseUnlink: true,
	}); err != nil {
		return fmt.Errorf("load Loyalty_Point_Earning_Rules: %w", err)
	}
	if _, _, err := redisx.LoadMongoToRedis[Loyalty_Point_Earning_Rules_Overwrite](ctx, RedisClient, Mdb_Loyalty_Point_Earning_Rules_Overwrite.Coll, redisx.MongoLoadOptions{
		BatchSize: 2000, FlushBeforeLoad: true, FlushPattern: "Loyalty_Point_Earning_Rules_Overwrite:*", UseUnlink: true,
	}); err != nil {
		return fmt.Errorf("load Loyalty_Point_Earning_Rules_Overwrite: %w", err)
	}
	if _, _, err := redisx.LoadMongoToRedis[Loyalty_Point_Expiry_Rules](ctx, RedisClient, Mdb_Loyalty_Point_Expiry_Rules.Coll, redisx.MongoLoadOptions{
		BatchSize: 2000, FlushBeforeLoad: true, FlushPattern: "Loyalty_Point_Expiry_Rules:*", UseUnlink: true,
	}); err != nil {
		return fmt.Errorf("load Loyalty_Point_Expiry_Rules: %w", err)
	}
	if _, _, err := redisx.LoadMongoToRedis[Loyalty_Point_Redemption_Rules](ctx, RedisClient, Mdb_Loyalty_Point_Redemption_Rules.Coll, redisx.MongoLoadOptions{
		BatchSize: 2000, FlushBeforeLoad: true, FlushPattern: "Loyalty_Point_Redemption_Rules:*", UseUnlink: true,
	}); err != nil {
		return fmt.Errorf("load Loyalty_Point_Redemption_Rules: %w", err)
	}
	if _, _, err := redisx.LoadMongoToRedis[Loyalty_Plan](ctx, RedisClient, Mdb_Loyalty_Plan.Coll, redisx.MongoLoadOptions{
		BatchSize: 2000, FlushBeforeLoad: true, FlushPattern: "Loyalty_Plan:*", UseUnlink: true,
	}); err != nil {
		return fmt.Errorf("load Loyalty_Plan: %w", err)
	}
	if _, _, err := redisx.LoadMongoToRedis[Customer_Loyalty_Account](ctx, RedisClient, Mdb_Customer_Loyalty_Account.Coll, redisx.MongoLoadOptions{
		BatchSize: 2000, FlushBeforeLoad: true, FlushPattern: "Customer_Loyalty_Account:*", UseUnlink: true,
	}); err != nil {
		return fmt.Errorf("load Customer_Loyalty_Account: %w", err)
	}
	if _, _, err := redisx.LoadMongoToRedis[Customer_Loyalty_Account_Points_Detail](ctx, RedisClient, Mdb_Customer_Loyalty_Account_Points_Detail.Coll, redisx.MongoLoadOptions{
		BatchSize: 2000, FlushBeforeLoad: true, FlushPattern: "Customer_Loyalty_Account_Points_Detail:*", UseUnlink: true,
	}); err != nil {
		return fmt.Errorf("load Customer_Loyalty_Account_Points_Detail: %w", err)
	}
	if _, _, err := redisx.LoadMongoToRedis[Customer_DND](ctx, RedisClient, Mdb_Customer_DND.Coll, redisx.MongoLoadOptions{
		BatchSize: 2000, FlushBeforeLoad: true, FlushPattern: "Customer_DND:*", UseUnlink: true,
	}); err != nil {
		return fmt.Errorf("load Customer_DND: %w", err)
	}
	if _, _, err := redisx.LoadMongoToRedis[Customer_Exclusion](ctx, RedisClient, Mdb_Customer_Exclusion.Coll, redisx.MongoLoadOptions{
		BatchSize: 2000, FlushBeforeLoad: true, FlushPattern: "Customer_Exclusion:*", UseUnlink: true,
	}); err != nil {
		return fmt.Errorf("load Customer_Exclusion: %w", err)
	}
	if _, _, err := redisx.LoadMongoToRedis[Customer_COS_Exclusion](ctx, RedisClient, Mdb_Customer_COS_Exclusion.Coll, redisx.MongoLoadOptions{
		BatchSize: 2000, FlushBeforeLoad: true, FlushPattern: "Customer_COS_Exclusion:*", UseUnlink: true,
	}); err != nil {
		return fmt.Errorf("load Customer_COS_Exclusion: %w", err)
	}
	if _, _, err := redisx.LoadMongoToRedis[Customer_UAT](ctx, RedisClient, Mdb_Customer_UAT.Coll, redisx.MongoLoadOptions{
		BatchSize: 2000, FlushBeforeLoad: true, FlushPattern: "Customer_UAT:*", UseUnlink: true,
	}); err != nil {
		return fmt.Errorf("load Customer_UAT: %w", err)
	}
	if _, _, err := redisx.LoadMongoToRedis[Loyalty_Campaign](ctx, RedisClient, Mdb_Loyalty_Campaign.Coll, redisx.MongoLoadOptions{
		BatchSize: 2000, FlushBeforeLoad: true, FlushPattern: "Loyalty_Campaign:*", UseUnlink: true,
	}); err != nil {
		return fmt.Errorf("load Loyalty_Campaign: %w", err)
	}
	if _, _, err := redisx.LoadMongoToRedis[Loyalty_Campaign_Target_List](ctx, RedisClient, Mdb_Loyalty_Campaign_Target_List.Coll, redisx.MongoLoadOptions{
		BatchSize: 2000, FlushBeforeLoad: true, FlushPattern: "Loyalty_Campaign_Target_List:*", UseUnlink: true,
	}); err != nil {
		return fmt.Errorf("load Loyalty_Campaign_Target_List: %w", err)
	}
	if _, _, err := redisx.LoadMongoToRedis[Loyalty_Campaign_Account](ctx, RedisClient, Mdb_Loyalty_Campaign_Account.Coll, redisx.MongoLoadOptions{
		BatchSize: 2000, FlushBeforeLoad: true, FlushPattern: "Loyalty_Campaign_Account:*", UseUnlink: true,
	}); err != nil {
		return fmt.Errorf("load Loyalty_Campaign_Account: %w", err)
	}

	return nil
}

func (uc *UserControl) LoyaltyIndexesMaintenanceProcess() {
	log.Println("Loyalty DB index maintenance process started...")

	ctx := context.Background()
	keyIdx, keyOpt := mongox.UniqueIndex("Key", true)
	msisdnIdx, msisdnOpt := mongox.UniqueIndex("MSISDN", true)

	if _, err := mongox.CreateIndex(ctx, Mdb_Loyalty_AutoIncrement.Coll, keyIdx, keyOpt); err != nil {
		log.Println("CreateIndex Loyalty_AutoIncrement/Key:", err)
	}
	if _, err := mongox.CreateIndex(ctx, Mdb_Loyalty_Governance.Coll, keyIdx, keyOpt); err != nil {
		log.Println("CreateIndex Loyalty_Governance/Key:", err)
	}
	if _, err := mongox.CreateIndex(ctx, Mdb_Loyalty_Level.Coll, keyIdx, keyOpt); err != nil {
		log.Println("CreateIndex Loyalty_Level/Key:", err)
	}
	if _, err := mongox.CreateIndex(ctx, Mdb_Loyalty_Account_Segment.Coll, keyIdx, keyOpt); err != nil {
		log.Println("CreateIndex Loyalty_Account_Segment/Key:", err)
	}
	if _, err := mongox.CreateIndex(ctx, Mdb_Loyalty_Point_Earning_Rules.Coll, keyIdx, keyOpt); err != nil {
		log.Println("CreateIndex Loyalty_Point_Earning_Rules/Key:", err)
	}
	if _, err := mongox.CreateIndex(ctx, Mdb_Loyalty_Point_Expiry_Rules.Coll, keyIdx, keyOpt); err != nil {
		log.Println("CreateIndex Loyalty_Point_Expiry_Rules/Key:", err)
	}
	if _, err := mongox.CreateIndex(ctx, Mdb_Loyalty_Point_Redemption_Rules.Coll, keyIdx, keyOpt); err != nil {
		log.Println("CreateIndex Loyalty_Point_Redemption_Rules/Key:", err)
	}
	if _, err := mongox.CreateIndex(ctx, Mdb_Loyalty_Plan.Coll, keyIdx, keyOpt); err != nil {
		log.Println("CreateIndex Loyalty_Plan/Key:", err)
	}
	if _, err := mongox.CreateIndex(ctx, Mdb_Customer_Loyalty_Account.Coll, keyIdx, keyOpt); err != nil {
		log.Println("CreateIndex Customer_Loyalty_Account/Key:", err)
	}
	if _, err := mongox.CreateIndex(ctx, Mdb_Customer_Loyalty_Account_Points_Detail.Coll, keyIdx, keyOpt); err != nil {
		log.Println("CreateIndex Customer_Loyalty_Account_Points_Detail/Key:", err)
	}
	if _, err := mongox.CreateIndex(ctx, Mdb_Customer_DND.Coll, keyIdx, keyOpt); err != nil {
		log.Println("CreateIndex Customer_DND/Key:", err)
	}
	if _, err := mongox.CreateIndex(ctx, Mdb_Customer_Exclusion.Coll, keyIdx, keyOpt); err != nil {
		log.Println("CreateIndex Customer_Exclusion/Key:", err)
	}
	if _, err := mongox.CreateIndex(ctx, Mdb_Customer_COS_Exclusion.Coll, keyIdx, keyOpt); err != nil {
		log.Println("CreateIndex Customer_COS_Exclusion/Key:", err)
	}
	if _, err := mongox.CreateIndex(ctx, Mdb_Customer_UAT.Coll, keyIdx, keyOpt); err != nil {
		log.Println("CreateIndex Customer_UAT/Key:", err)
	}
	if _, err := mongox.CreateIndex(ctx, Mdb_Loyalty_Campaign.Coll, keyIdx, keyOpt); err != nil {
		log.Println("CreateIndex Loyalty_Campaign/Key:", err)
	}
	if _, err := mongox.CreateIndex(ctx, Mdb_Loyalty_Campaign_Target_List.Coll, keyIdx, keyOpt); err != nil {
		log.Println("CreateIndex Loyalty_Campaign_Target_List/Key:", err)
	}
	if _, err := mongox.CreateIndex(ctx, Mdb_Loyalty_Campaign_Account.Coll, keyIdx, keyOpt); err != nil {
		log.Println("CreateIndex Loyalty_Campaign_Account/Key:", err)
	}
	if _, err := mongox.CreateIndex(ctx, Mdb_Loyalty_AccountCreditPoints_log.Coll, msisdnIdx, msisdnOpt); err != nil {
		log.Println("CreateIndex Loyalty_AccountCreditPoints_log/MSISDN:", err)
	}
	if _, err := mongox.CreateIndex(ctx, Mdb_Loyalty_AccountDebitPoints_log.Coll, msisdnIdx, msisdnOpt); err != nil {
		log.Println("CreateIndex Loyalty_AccountDebitPoints_log/MSISDN:", err)
	}
	if _, err := mongox.CreateIndex(ctx, Mdb_Loyalty_Point_Redemption_Rules.Coll, msisdnIdx, msisdnOpt); err != nil {
		log.Println("CreateIndex Loyalty_Point_Redemption_Rules/MSISDN:", err)
	}
	if _, err := mongox.CreateIndex(ctx, Mdb_Loyalty_Status_log.Coll, msisdnIdx, msisdnOpt); err != nil {
		log.Println("CreateIndex Loyalty_Status_log/MSISDN:", err)
	}

	log.Println("Loyalty DB index maintenance process completed")
}
func (Uc *UserControl) Write_Loyalty_Event_Log(record Loyalty_Event_Log) {
	YYYY, MM, _, DD, _, _, _ := GetTimeParts(record.Event_Time)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := Uc.LoyaltyMongoClient.Mongo.Database(Configuration.DB_Name_Loyalty + "_" + YYYY + MM).Collection("Col_Loyalty_Event_Log_" + DD)
	_, err := col.InsertOne(ctx, record)
	if err != nil {
		log.Println("Error in Write_Loyalty_Event_Log:", err, " (", record, ")")
		return
	}
}

func (Uc *UserControl) Write_Loyalty_Level_Change_log(record Loyalty_Level_Change_log) {
	YYYY, MM, _, _, _, _, _ := GetTimeParts(record.Level_Change_Date)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := Uc.LoyaltyMongoClient.Mongo.Database(Configuration.DB_Name_Loyalty + "_" + YYYY + MM).Collection("Col_Loyalty_Level_Change_log")
	_, err := col.InsertOne(ctx, record)
	if err != nil {
		log.Println("Error in Write_Loyalty_Level_Change_log:", err, " (", record, ")")
		return
	}

	Earningrecord, err := Uc.Customer_Loyalty_Account_GetEarning_Rule(record.MSISDN)
	if err != nil {
		log.Println("failed to get data")
		return
	}
	if Earningrecord.Level_Change_Notification {
		LevelChangeNotiLog := NotificationLog{
			SourceAction:  "LevelChange",
			TransactionId: "",
			Medium:        "SMS",
			SourceAddress: Earningrecord.Level_Change_Notification_Sender,
			Destination:   record.MSISDN,
			Subject:       "LevelChange",
			AddUser:       "SYSTEM",
			AddDate:       time.Now(),
		}
		LevelChange_Noti_Text := ""
		LevelChange_Noti_Text = Earningrecord.Level_Change_Notification_Text
		if LevelChange_Noti_Text != "" {
			LevelChange_Noti_Text = strings.ReplaceAll(LevelChange_Noti_Text, "{{PreviousLevel}}", fmt.Sprint(record.Previous_Loyalty_Level_Key))
			LevelChange_Noti_Text = strings.ReplaceAll(LevelChange_Noti_Text, "{{NewLevel}}", fmt.Sprint(record.New_Loyalty_Level_Key))
			LevelChange_Noti_Text = strings.ReplaceAll(LevelChange_Noti_Text, "{{LevelChangeDirection}}", fmt.Sprint(record.New_Loyalty_Level_Direction))
			LevelChange_Noti_Text = strings.ReplaceAll(LevelChange_Noti_Text, "{{LoyaltyBalance}}", fmt.Sprint(record.Available_Points))
			LevelChangeNotiLog.Payload = LevelChange_Noti_Text
			err := error(nil)
			if Configuration.Operation == "Angola" {
				err = SendSMS(Earningrecord.Level_Change_Notification_Sender, record.MSISDN, LevelChange_Noti_Text)
			} else {
				err = Send_SMS(Earningrecord.Level_Change_Notification_Sender, record.MSISDN, LevelChange_Noti_Text)
			}
			if err != nil {
				LevelChangeNotiLog.Status = "Failed"
				LevelChangeNotiLog.Error = err.Error()
			} else {
				LevelChangeNotiLog.Status = "Successful"
			}
		} else {
			LevelChangeNotiLog.Payload = LevelChange_Noti_Text
			LevelChangeNotiLog.Status = "Failed"
			LevelChangeNotiLog.Error = "Undefined level change notification for transaction"
		}
		YYYY, MM, _, _, _, _, _ = GetTimeParts(record.Level_Change_Date)
		notiCtx, notiCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer notiCancel()
		notiCol := Uc.MongoClient.Mongo.Database(Configuration.DB_Name_Loyalty + "_Logs_" + YYYY + MM).Collection("Col_NotificationLog")
		_, err = notiCol.InsertOne(notiCtx, LevelChangeNotiLog)
		if err != nil {
			log.Println("Error in Write level change Notification Logs:", err, " (", LevelChangeNotiLog, ")")
		}
	}

}

func (Uc *UserControl) Write_Loyalty_AccountCreditPoints_log(record Loyalty_AccountCreditPoints_log) {
	YYYY, MM, _, DD, _, _, _ := GetTimeParts(record.ReceiveDate)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := Uc.LoyaltyMongoClient.Mongo.Database(Configuration.DB_Name_Loyalty + "_" + YYYY + MM).Collection("Col_Loyalty_AccountCreditPoints_log_" + DD)
	_, err := col.InsertOne(ctx, record)
	if err != nil {
		log.Println("Error in Write_Loyalty_AccountCreditPoints_log:", err, " (", record, ")")
		return
	}
}

func (Uc *UserControl) Write_Loyalty_Monthly_Expiry_log(record Loyalty_Monthly_Expiry_log) {
	YYYY, MM, _, DD, _, _, _ := GetTimeParts(record.ExpiryTime)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := Uc.LoyaltyMongoClient.Mongo.Database(Configuration.DB_Name_Loyalty + "_" + YYYY + MM).Collection("Col_Loyalty_Expiry_log_" + DD)
	_, err := col.InsertOne(ctx, record)
	if err != nil {
		log.Println("Error in Write_Loyalty_Monthly_Expiry_log:", err, " (", record, ")")
		return
	}
}

func (Uc *UserControl) Write_Loyalty_Full_Expiry_log(record Loyalty_Full_Expiry_Log) {
	YYYY, MM, _, DD, _, _, _ := GetTimeParts(record.ExpiryTime)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := Uc.LoyaltyMongoClient.Mongo.Database(Configuration.DB_Name_Loyalty + "_" + YYYY + MM).Collection("Col_Loyalty_Full_Expiry_log_" + DD)
	_, err := col.InsertOne(ctx, record)
	if err != nil {
		log.Println("Error in Write_Full_Loyalty_Expiry_log:", err, " (", record, ")")
		return
	}
}

func (Uc *UserControl) Write_Loyalty_Redemption_log(record Loyalty_Redemption_log) {
	YYYY, MM, _, DD, _, _, _ := GetTimeParts(record.ReceiveDate)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := Uc.LoyaltyMongoClient.Mongo.Database(Configuration.DB_Name_Loyalty + "_" + YYYY + MM).Collection("Col_Loyalty_Redemption_log_" + DD)
	_, err := col.InsertOne(ctx, record)
	if err != nil {
		log.Println("Error in Write_Loyalty_Redemption_log:", err, " (", record, ")")
		return
	}
}

func (Uc *UserControl) Write_Loyalty_Account_Churned_log(record Customer_Loyalty_Account) {
	YYYY, MM, _, DD, _, _, _ := GetTimeParts(time.Now())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := Uc.LoyaltyMongoClient.Mongo.Database(Configuration.DB_Name_Loyalty + "_" + YYYY + MM).Collection("Col_Churned_Customer_Loyalty_Account_" + DD)
	_, err := col.InsertOne(ctx, record)
	if err != nil {
		log.Println("Error in Write_Loyalty_Account_Churned_log:", err, " (", record, ")")
		return
	}
}

func (Uc *UserControl) Write_Loyalty_Status_log(record Loyalty_Status_log) {
	YYYY, MM, _, DD, _, _, _ := GetTimeParts(record.StatusDate)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := Uc.LoyaltyMongoClient.Mongo.Database(Configuration.DB_Name_Loyalty + "_" + YYYY + MM).Collection("Col_Loyalty_Status_log_" + DD)
	_, err := col.InsertOne(ctx, record)
	if err != nil {
		log.Println("Error in Write_Loyalty_Status_log:", err, " (", record, ")")
		return
	}
}

func (Uc *UserControl) Write_Loyalty_AccountDebitPoints_log(record Loyalty_AccountDebitPoints_log) {
	YYYY, MM, _, DD, _, _, _ := GetTimeParts(record.ReceiveDate)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := Uc.LoyaltyMongoClient.Mongo.Database(Configuration.DB_Name_Loyalty + "_" + YYYY + MM).Collection("Col_Loyalty_AccountDebitPoints_log_" + DD)
	_, err := col.InsertOne(ctx, record)
	if err != nil {
		log.Println("Error in Write_Loyalty_AccountDebitPoints_log:", err, " (", record, ")")
		return
	}
}

func (Uc *UserControl) Write_Loyalty_Governance_log(record Loyalty_Governance_log) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := Mdb_Loyalty_Governance_log.Coll.InsertOne(ctx, record)
	if err != nil {
		log.Println("Error in Write_Loyalty_Governance_log:", err, " (", record, ")")
		return
	}
}

func (Uc *UserControl) InitializeLoyaltyDefaultUAT() {
	Uc.Loyalty_Governance_Add("Default", Loyalty_Governance_AddRequest{
		Key:                             LOYALTY_GOVERNANCE_KEY,
		Available_Points_Pool:           5000000000,
		Distributed_Points_Pool:         0,
		Redeemed_Points_Pool:            0,
		Expired_Points_Pool:             0,
		MaxAllowedPoints_PerTransaction: 100,
		MaxSubsAwardedPoints_PerMonth:   10000,
		MaxSubsAwardedPoints:            100000,
		EVC_Account_Balance:             0,
		Merchant_Account_Balance:        0,
		DailyEarningLimit:               10000,
		DailyPointsRedemptionLimit:      10000,
		DailyRedemptionAttemptLimit:     10,
		WeeklyEarningLimit:              100000,
		WeeklyPointsRedemptionLimit:     100000,
		WeeklyRedemptionAttemptLimit:    100,
	})
	Uc.Loyalty_Level_Add("Default", Loyalty_Level_AddRequest{
		Key:                    "Member",
		Description:            "Member",
		Min_Accumulated_Points: 0,
		Max_Accumulated_Points: 100,
		EnableRedeem:           true,
	})
	Uc.Loyalty_Level_Add("Default", Loyalty_Level_AddRequest{
		Key:                    "Silver",
		Description:            "Silver",
		Min_Accumulated_Points: 101,
		Max_Accumulated_Points: 500,
		EnableRedeem:           true,
	})
	Uc.Loyalty_Level_Add("Default", Loyalty_Level_AddRequest{
		Key:                    "Gold",
		Description:            "Gold",
		Min_Accumulated_Points: 501,
		Max_Accumulated_Points: 1000,
		EnableRedeem:           true,
	})
	Uc.Loyalty_Level_Add("Default", Loyalty_Level_AddRequest{
		Key:                    "Platinum",
		Description:            "Platinum",
		Min_Accumulated_Points: 1001,
		Max_Accumulated_Points: 999999999,
		EnableRedeem:           true,
	})
	Uc.Loyalty_Account_Segment_Add("Default", Loyalty_Account_Segment_AddRequest{
		Key:         "Main_Segment",
		Description: "Main_Segment",
		Amount_From: 0,
		Amount_Till: 999999999,
		AON_From:    0,
		AON_Till:    999999999,
	})
	Uc.Loyalty_Point_Earning_Rules_Add("Default", Loyalty_Point_Earning_Rules_AddRequest{
		Key:                            "Default_Earning_Rules",
		Description:                    "Default_Earning_Rules",
		Welcome_Points:                 5,
		MobileAppDaily_Login:           1,
		MainGSMBalance_Amount:          10,
		MainGSMBalance_Points:          1,
		GSM_SC_Airtime_Award_Type:      "Transaction",
		GSM_EVC_Airtime_Award_Type:     "Transaction",
		GSM_EVC_Bundle_Award_Type:      "Transaction",
		MM_P2P_Award_Type:              "Transaction",
		MM_P2P_Amount:                  0,
		MM_P2P_Points:                  1,
		MM_CASHIN_Award_Type:           "Transaction",
		MM_CASHIN_Amount:               0,
		MM_CASHIN_Points:               1,
		MM_CASHOUT_Award_Type:          "Transaction",
		MM_CASHOUT_Amount:              0,
		MM_CASHOUT_Points:              1,
		MM_MERCHPAY_Award_Type:         "Transaction",
		MM_MERCHPAY_Amount:             0,
		MM_MERCHPAY_Points:             1,
		MM_BILLPAY_Award_Type:          "Transaction",
		MM_BILLPAY_Amount:              0,
		MM_BILLPAY_Points:              1,
		MM_RC_Bundle_Award_Type:        "Amount",
		MM_RC_Bundle_Amount:            15,
		MM_RC_Bundle_Points:            1,
		MM_RC_Airtime_Award_Type:       "Amount",
		MM_RC_Airtime_Amount:           15,
		MM_RC_Airtime_Points:           1,
		MM_CTMMOREQ_Bundle_Award_Type:  "Amount",
		MM_CTMMOREQ_Bundle_Amount:      15,
		MM_CTMMOREQ_Bundle_Points:      1,
		MM_CTMMOREQ_Airtime_Award_Type: "Amount",
		MM_CTMMOREQ_Airtime_Amount:     15,
		MM_CTMMOREQ_Airtime_Points:     1,
		MM_CBWREQ_Award_Type:           "Transaction",
		MM_CBWREQ_Amount:               0,
		MM_CBWREQ_Points:               1,
	})
	Uc.Loyalty_Point_Expiry_Rules_Add("Default", Loyalty_Point_Expiry_Rules_AddRequest{
		Key:                     "Default_Expiry_Rules",
		Description:             "Default_Expiry_Rules",
		Rolling_Expiration:      false,
		Validity_Unit:           "", //Month, Year --> only when Rolling_Expiration is true
		Validity_Duration:       0,  //only when Rolling_Expiration is true
		Fix_Date_Expiration:     true,
		Expiration_Trigger_date: time.Date(2026, 01, 31, 00, 00, 59, 0, time.UTC), //when the expiry process will run
		Expiration_Point_Before: time.Date(2025, 12, 31, 00, 00, 59, 0, time.UTC), //expiry all points before this date
	})
	Uc.Loyalty_Point_Redemption_Rules_Add("Default", Loyalty_Point_Redemption_Rules_AddRequest{
		Key:                                 "Default_Redemption_Rules",
		Description:                         "Default_Redemption_Rules",
		Min_Accumulated_Points:              100,
		Allow_Negative_Balance_ToRedeem:     false,
		Allow_PendingLendme_ToRedeem:        false,
		Available_MinPoints_for_Airtime:     100,
		Airtime_Amount:                      2,
		Airtime_Points:                      1,
		Available_MinPoints_for_MobileMoney: 100,
		MobileMoney_Amount:                  2,
		MobileMoney_Points:                  1,
		Bundles_MinPoints:                   100,
		Bundles_Product_Catalogue_Channel:   "Loyalty_Default_Channel",
		Bundles_Product_Catalogue_Plan:      "Loyalty_Default_Plan",
		Bundles_Product_Catalogue_Version:   "1",
	})
	Uc.Loyalty_Plan_Add("Default", Loyalty_Plan_AddRequest{
		Key:                         "Member|Main_Segment", //Loyalty_Level_Key + "|" + Loyalty_Account_Segment_Key
		Description:                 "",
		Loyalty_Level_Key:           "Member",
		Loyalty_Account_Segment_Key: "Main_Segment",
		Earning_Rules_Key:           "Default_Earning_Rules",
		Expiry_Rules_Key:            "Default_Expiry_Rules",
		Redemption_Rules_Key:        "Default_Redemption_Rules",
	})
	Uc.Loyalty_Plan_Add("Default", Loyalty_Plan_AddRequest{
		Key:                         "Silver|Main_Segment", //Loyalty_Level_Key + "|" + Loyalty_Account_Segment_Key
		Description:                 "",
		Loyalty_Level_Key:           "Silver",
		Loyalty_Account_Segment_Key: "Main_Segment",
		Earning_Rules_Key:           "Default_Earning_Rules",
		Expiry_Rules_Key:            "Default_Expiry_Rules",
		Redemption_Rules_Key:        "Default_Redemption_Rules",
	})
	Uc.Loyalty_Plan_Add("Default", Loyalty_Plan_AddRequest{
		Key:                         "Gold|Main_Segment", //Loyalty_Level_Key + "|" + Loyalty_Account_Segment_Key
		Description:                 "",
		Loyalty_Level_Key:           "Gold",
		Loyalty_Account_Segment_Key: "Main_Segment",
		Earning_Rules_Key:           "Default_Earning_Rules",
		Expiry_Rules_Key:            "Default_Expiry_Rules",
		Redemption_Rules_Key:        "Default_Redemption_Rules",
	})
	Uc.Loyalty_Plan_Add("Default", Loyalty_Plan_AddRequest{
		Key:                         "Platinum|Main_Segment", //Loyalty_Level_Key + "|" + Loyalty_Account_Segment_Key
		Description:                 "",
		Loyalty_Level_Key:           "Platinum",
		Loyalty_Account_Segment_Key: "Main_Segment",
		Earning_Rules_Key:           "Default_Earning_Rules",
		Expiry_Rules_Key:            "Default_Expiry_Rules",
		Redemption_Rules_Key:        "Default_Redemption_Rules",
	})
}

// ***********************************************************************
// Loyalty_Governance functions
// ***********************************************************************
func (Uc *UserControl) Loyalty_Governance_Add(Login string, request Loyalty_Governance_AddRequest) (Id int64, err error) {
	//check key if filled and if already used
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	if request.DailyEarningLimit > request.WeeklyEarningLimit {
		return Id, errors.New("Daily Limit cannot exceed Weekly Limit")
	}
	if request.DailyPointsRedemptionLimit > request.WeeklyPointsRedemptionLimit {
		return Id, errors.New("Daily Limit cannot exceed Weekly Limit")
	}
	if request.DailyRedemptionAttemptLimit > request.WeeklyRedemptionAttemptLimit {
		return Id, errors.New("Daily Limit cannot exceed Weekly Limit")
	}
	//check if key already used
	{
		chkCtx, chkCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, chkErr := redisx.GetJSON[Loyalty_Governance](chkCtx, RedisClient, Loyalty_Governance{Key: request.Key}.RedisKey())
		chkCancel()
		if chkErr == nil {
			err = errors.New("key already exist")
			return Id, err
		}
	}

	//Prepare new entry
	var NewEntry Loyalty_Governance
	{
		aiCtx, aiCancel := context.WithTimeout(context.Background(), 10*time.Second)
		NewEntry.Governance_Id, err = redisx.NextAutoIncrementID(aiCtx, RedisClient, Mdb_Loyalty_AutoIncrement.Coll, "Loyalty_Governance-Id", redisx.NextIDOptions{
			RedisBase: "AutoIncrement", MongoRetries: 3, RetryBackoff: 500 * time.Millisecond,
		})
		aiCancel()
		if err != nil {
			return Id, err
		}
		Id = NewEntry.Governance_Id
	}
	NewEntry.Key = request.Key
	NewEntry.Available_Points_Pool = request.Available_Points_Pool
	NewEntry.Distributed_Points_Pool = request.Distributed_Points_Pool
	NewEntry.Redeemed_Points_Pool = request.Redeemed_Points_Pool
	NewEntry.MaxAllowedPoints_PerTransaction = request.MaxAllowedPoints_PerTransaction
	NewEntry.MaxSubsAwardedPoints_PerMonth = request.MaxSubsAwardedPoints_PerMonth
	NewEntry.MaxSubsAwardedPoints = request.MaxSubsAwardedPoints
	NewEntry.DailyEarningLimit = request.DailyEarningLimit
	NewEntry.DailyPointsRedemptionLimit = request.DailyPointsRedemptionLimit
	NewEntry.DailyRedemptionAttemptLimit = request.DailyRedemptionAttemptLimit
	NewEntry.WeeklyEarningLimit = request.WeeklyEarningLimit
	NewEntry.WeeklyPointsRedemptionLimit = request.WeeklyPointsRedemptionLimit
	NewEntry.WeeklyRedemptionAttemptLimit = request.WeeklyRedemptionAttemptLimit
	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Loyalty_Governance.Coll.UpdateOne(putCtx, bson.M{"Key": NewEntry.Key}, bson.M{"$set": NewEntry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			putCancel()
			err = putErr
			return Id, err
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, NewEntry.RedisKey(), NewEntry); putSetErr != nil {
			log.Println("redisx.SetJSON Loyalty_Governance error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Governance",
		Event_ActionType:   "Add",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  NewEntry,
	})
	return Id, nil
}

func (Uc *UserControl) Loyalty_Governance_Edit(Login string, request Loyalty_Governance_EditRequest) (Id int64, err error) {
	//check and validate
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	if request.DailyEarningLimit > request.WeeklyEarningLimit {
		return Id, errors.New("Daily Limit cannot exceed Weekly Limit")
	}
	if request.DailyPointsRedemptionLimit > request.WeeklyPointsRedemptionLimit {
		return Id, errors.New("Daily Limit cannot exceed Weekly Limit")
	}
	if request.DailyRedemptionAttemptLimit > request.WeeklyRedemptionAttemptLimit {
		return Id, errors.New("Daily Limit cannot exceed Weekly Limit")
	}
	entry, entryErr := redisx.GetJSON[Loyalty_Governance](context.Background(), RedisClient, Loyalty_Governance{Key: request.Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("key is not created")
		return Id, err
	}
	if entryErr != nil {
		return Id, entryErr
	}
	if entry.Governance_Id != request.Governance_Id {
		return Id, errors.New("id is not matching")
	}
	Current_Entry := entry

	//Prepare new entry
	entry.Key = request.Key
	entry.Available_Points_Pool = request.Available_Points_Pool
	entry.Distributed_Points_Pool = request.Distributed_Points_Pool
	entry.Redeemed_Points_Pool = request.Redeemed_Points_Pool
	entry.MaxAllowedPoints_PerTransaction = request.MaxAllowedPoints_PerTransaction
	entry.MaxSubsAwardedPoints_PerMonth = request.MaxSubsAwardedPoints_PerMonth
	entry.MaxSubsAwardedPoints = request.MaxSubsAwardedPoints
	entry.DailyEarningLimit = request.DailyEarningLimit
	entry.DailyPointsRedemptionLimit = request.DailyPointsRedemptionLimit
	entry.DailyRedemptionAttemptLimit = request.DailyRedemptionAttemptLimit
	entry.WeeklyEarningLimit = request.WeeklyEarningLimit
	entry.WeeklyPointsRedemptionLimit = request.WeeklyPointsRedemptionLimit
	entry.WeeklyRedemptionAttemptLimit = request.WeeklyRedemptionAttemptLimit

	if request.NewKey != "" {
		if request.NewKey != request.Key {
			//delete old
			{
				delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
				if _, delErr := Mdb_Loyalty_Governance.Coll.DeleteOne(delCtx, bson.M{"Key": request.Key}); delErr != nil {
					log.Println("Mdb_Loyalty_Governance DeleteOne error:", delErr)
				}
				if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Loyalty_Governance{Key: request.Key}.RedisKey()); delRedisErr != nil {
					log.Println("redisx.DelJSON Loyalty_Governance error:", delRedisErr)
				}
				delCancel()
			}
			//update key
			entry.Key = request.NewKey
		}
	}
	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Loyalty_Governance.Coll.UpdateOne(putCtx, bson.M{"Key": entry.Key}, bson.M{"$set": entry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			putCancel()
			err = putErr
			return Id, err
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, entry.RedisKey(), entry); putSetErr != nil {
			log.Println("redisx.SetJSON Loyalty_Governance error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Governance",
		Event_ActionType:   "Edit",
		Event_Description:  "",
		Event_Entry_Before: Current_Entry,
		Event_Entry_After:  entry,
	})
	return Id, nil
}

func (Uc *UserControl) Loyalty_Governance_Get(Key string) (entries []Loyalty_Governance, err error) {
	if Key == "" {
		entries, err = redisx.GetAllJSONByPattern[Loyalty_Governance](context.Background(), RedisClient, redisx.ScanJSONOptions{
			Pattern: "Loyalty_Governance:*", ScanCount: 500, PipelineSize: 250,
		})
		if err != nil {
			return entries, err
		}
		if Configuration.Operation != "Angola" {
			redemptions, redemErr := redisx.GetAllJSONByPattern[Loyalty_Point_Redemption_Rules](context.Background(), RedisClient, redisx.ScanJSONOptions{
				Pattern: "Loyalty_Point_Redemption_Rules:*", ScanCount: 500, PipelineSize: 250,
			})
			if redemErr != nil {
				return entries, redemErr
			}
			for i := range entries {
				for _, redemption := range redemptions {
					var Airtime_EVC_PIN = ""
					if redemption.Airtime_EVC_PIN != "" {
						Airtime_EVC_PIN, err = DecryptHexString(redemption.Airtime_EVC_PIN)
						if err != nil {
							fmt.Println("error in decrypting artime evc pin", err.Error())
						}
					}
					res, err := Uc.CGW.UC_GWClient.GetERDealerBalance(redemption.Airtime_EVC_Account, Airtime_EVC_PIN)
					if err == nil {
						entries[i].EVC_Account_Balance = float64(res.Balance)
					}
				}
			}
		}
		return entries, nil
	} else {
		entry, entryErr := redisx.GetJSON[Loyalty_Governance](context.Background(), RedisClient, Loyalty_Governance{Key: Key}.RedisKey())
		if redisx.IsNil(entryErr) {
			err = errors.New("key does not exist")
			return entries, err
		}
		if entryErr != nil {
			return entries, entryErr
		}
		if Configuration.Operation != "Angola" {
			redemptions, redemErr := redisx.GetAllJSONByPattern[Loyalty_Point_Redemption_Rules](context.Background(), RedisClient, redisx.ScanJSONOptions{
				Pattern: "Loyalty_Point_Redemption_Rules:*", ScanCount: 500, PipelineSize: 250,
			})
			if redemErr != nil {
				return entries, redemErr
			}
			for _, redemption := range redemptions {
				var Airtime_EVC_PIN = ""
				if redemption.Airtime_EVC_PIN != "" {
					Airtime_EVC_PIN, err = DecryptHexString(redemption.Airtime_EVC_PIN)
					if err != nil {
						fmt.Println("error in decrypting artime evc pin", err.Error())
					}
				}
				res, err := Uc.CGW.UC_GWClient.GetERDealerBalance(redemption.Airtime_EVC_Account, Airtime_EVC_PIN)
				if err == nil {
					entry.EVC_Account_Balance = float64(res.Balance)
				}
			}
		}

		entries = append(entries, entry)
		return entries, nil
	}
}

func (Uc *UserControl) Loyalty_Governance_GetPaginated(Page, Limit int64) (entries []Loyalty_Governance, err error) {
	if Page < 1 {
		return entries, errors.New("invalid page")
	}
	if Limit < 1 || Limit > 50000 {
		return entries, errors.New("invalid limit (accept value between 1 and 50000)")
	}

	//var array []daoc.DAOFindCriteria
	// if Outlet_Key != "" {
	// 	//restrict access for records that belong to this user
	// 	var criteria daoc.DAOFindCriteria = daoc.DAOFindCriteria{
	// 		Field:    "Outlet_Key",
	// 		Value:    Outlet_Key,
	// 		Operator: "EQUAL",
	// 	}
	// 	array = append(array, criteria)
	// }
	// if Agent_Key != "" {
	// 	//restrict access for records that belong to this user
	// 	var criteria daoc.DAOFindCriteria = daoc.DAOFindCriteria{
	// 		Field:    "Agent_Key",
	// 		Value:    Agent_Key,
	// 		Operator: "EQUAL",
	// 	}
	// 	array = append(array, criteria)
	// }
	// if len(array) > 0 {
	// 	findparams.FindCriteria = array
	// }
	pgCtx, pgCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer pgCancel()
	result, err := mongox.FindPage[Loyalty_Governance](pgCtx, Mdb_Loyalty_Governance.Coll, bson.M{}, mongox.PageRequest{
		Page: Page, PageSize: Limit, MaxPageSize: 50000,
	})
	if err != nil {
		return entries, err
	}
	entries = result.Items
	return entries, nil

}

func (Uc *UserControl) Loyalty_Governance_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry, entryErr := redisx.GetJSON[Loyalty_Governance](context.Background(), RedisClient, Loyalty_Governance{Key: Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("entry does not exist")
		return err
	}
	if entryErr != nil {
		return entryErr
	}
	{
		delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, delErr := Mdb_Loyalty_Governance.Coll.DeleteOne(delCtx, bson.M{"Key": Key}); delErr != nil {
			log.Println("Mdb_Loyalty_Governance DeleteOne error:", delErr)
		}
		if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Loyalty_Governance{Key: Key}.RedisKey()); delRedisErr != nil {
			log.Println("redisx.DelJSON Loyalty_Governance error:", delRedisErr)
		}
		delCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Governance",
		Event_ActionType:   "Delete",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  entry,
	})
	return nil
}

func (Uc *UserControl) Loyalty_Governance_Available_Points_Debit(points float64, refund ...bool) (err error) {
	chan_LoyaltyGovernance_Controler <- 1
	loyalty_governance, loyaltyErr := redisx.GetJSON[Loyalty_Governance](context.Background(), RedisClient, Loyalty_Governance{Key: LOYALTY_GOVERNANCE_KEY}.RedisKey())
	if loyaltyErr != nil {
		<-chan_LoyaltyGovernance_Controler
		return errors.New("loyalty governance entry not found")
	}
	if (loyalty_governance.Available_Points_Pool - (loyalty_governance.Distributed_Points_Pool + loyalty_governance.Expired_Points_Pool + loyalty_governance.Redeemed_Points_Pool)) > points {
		refundValue := false
		if len(refund) > 0 {
			refundValue = refund[0]
		}
		loyalty_governance.Distributed_Points_Pool = loyalty_governance.Distributed_Points_Pool + points
		if refundValue {
			loyalty_governance.Redeemed_Points_Pool = loyalty_governance.Redeemed_Points_Pool - points
		}
		{
			putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
			if _, putErr := Mdb_Loyalty_Governance.Coll.UpdateOne(putCtx, bson.M{"Key": loyalty_governance.Key}, bson.M{"$set": loyalty_governance}, options.UpdateOne().SetUpsert(true)); putErr != nil {
				log.Println("Mdb_Loyalty_Governance upsert error:", putErr)
			}
			if putSetErr := redisx.SetJSON(putCtx, RedisClient, loyalty_governance.RedisKey(), loyalty_governance); putSetErr != nil {
				log.Println("redisx.SetJSON Loyalty_Governance error:", putSetErr)
			}
			putCancel()
		}
		<-chan_LoyaltyGovernance_Controler
		return
	} else {
		<-chan_LoyaltyGovernance_Controler
		return errors.New("no enough loyalty points to distribute from the governance available balance")
	}
}

func (Uc *UserControl) Loyalty_Governance_Redeem_Points_Debit(points float64, notRedemption ...bool) (err error) {
	chan_LoyaltyGovernance_Controler <- 1
	loyalty_governance, loyaltyErr := redisx.GetJSON[Loyalty_Governance](context.Background(), RedisClient, Loyalty_Governance{Key: LOYALTY_GOVERNANCE_KEY}.RedisKey())
	if loyaltyErr != nil {
		<-chan_LoyaltyGovernance_Controler
		return errors.New("loyalty governance entry not found")
	}
	notRedemptionValue := false
	if len(notRedemption) > 0 {
		notRedemptionValue = notRedemption[0]
	}
	if !notRedemptionValue {
		loyalty_governance.Redeemed_Points_Pool = loyalty_governance.Redeemed_Points_Pool + points
	}
	loyalty_governance.Distributed_Points_Pool = loyalty_governance.Distributed_Points_Pool - points
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Loyalty_Governance.Coll.UpdateOne(putCtx, bson.M{"Key": loyalty_governance.Key}, bson.M{"$set": loyalty_governance}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Loyalty_Governance upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, loyalty_governance.RedisKey(), loyalty_governance); putSetErr != nil {
			log.Println("redisx.SetJSON Loyalty_Governance error:", putSetErr)
		}
		putCancel()
	}
	<-chan_LoyaltyGovernance_Controler
	return
}

func (Uc *UserControl) Loyalty_Governance_Expiry_Points_Credit(points float64) (err error) {
	chan_LoyaltyGovernance_Controler <- 1
	loyalty_governance, loyaltyErr := redisx.GetJSON[Loyalty_Governance](context.Background(), RedisClient, Loyalty_Governance{Key: LOYALTY_GOVERNANCE_KEY}.RedisKey())
	if loyaltyErr != nil {
		<-chan_LoyaltyGovernance_Controler
		return errors.New("loyalty governance entry not found")
	}
	loyalty_governance.Expired_Points_Pool = loyalty_governance.Expired_Points_Pool + points
	loyalty_governance.Distributed_Points_Pool = loyalty_governance.Distributed_Points_Pool - points
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Loyalty_Governance.Coll.UpdateOne(putCtx, bson.M{"Key": loyalty_governance.Key}, bson.M{"$set": loyalty_governance}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Loyalty_Governance upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, loyalty_governance.RedisKey(), loyalty_governance); putSetErr != nil {
			log.Println("redisx.SetJSON Loyalty_Governance error:", putSetErr)
		}
		putCancel()
	}
	<-chan_LoyaltyGovernance_Controler
	return
}

func (Uc *UserControl) Loyalty_Governance_Status_Expiry_Points_Credit(points float64) (err error) {
	chan_LoyaltyGovernance_Controler <- 1
	loyalty_governance, loyaltyErr := redisx.GetJSON[Loyalty_Governance](context.Background(), RedisClient, Loyalty_Governance{Key: LOYALTY_GOVERNANCE_KEY}.RedisKey())
	if loyaltyErr != nil {
		<-chan_LoyaltyGovernance_Controler
		return errors.New("loyalty governance entry not found")
	}
	loyalty_governance.Distributed_Points_Pool = loyalty_governance.Distributed_Points_Pool - points
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Loyalty_Governance.Coll.UpdateOne(putCtx, bson.M{"Key": loyalty_governance.Key}, bson.M{"$set": loyalty_governance}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Loyalty_Governance upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, loyalty_governance.RedisKey(), loyalty_governance); putSetErr != nil {
			log.Println("redisx.SetJSON Loyalty_Governance error:", putSetErr)
		}
		putCancel()
	}
	<-chan_LoyaltyGovernance_Controler
	return
}

func (Uc *UserControl) Loyalty_Governance_DailyLog_Process() {
	exec := 0
	LOG_ID := "<<Loyalty Governance Daily Log Process>>"
	for range time.Tick(time.Second * 1) {
		_CurrentDateTime := time.Now()
		_hr, _mi, _se := _CurrentDateTime.Clock()
		if _hr == 00 {
			if _mi == 00 {
				if _se < 60 {
					if exec == 0 {
						exec = 1
						log.Println(LOG_ID + " triggered")
						loyalty_governance, loyaltyErr := redisx.GetJSON[Loyalty_Governance](context.Background(), RedisClient, Loyalty_Governance{Key: LOYALTY_GOVERNANCE_KEY}.RedisKey())
						if loyaltyErr == nil {
							var log_entry Loyalty_Governance_log
							log_entry.Log_Date = time.Now()
							log_entry.Available_Points_Pool = loyalty_governance.Available_Points_Pool
							log_entry.Distributed_Points_Pool = loyalty_governance.Distributed_Points_Pool
							log_entry.Redeemed_Points_Pool = loyalty_governance.Redeemed_Points_Pool
							log_entry.Expired_Points_Pool = loyalty_governance.Expired_Points_Pool
							if Configuration.Operation != "Angola" {
								redemptions, redemErr := redisx.GetAllJSONByPattern[Loyalty_Point_Redemption_Rules](context.Background(), RedisClient, redisx.ScanJSONOptions{
									Pattern: "Loyalty_Point_Redemption_Rules:*", ScanCount: 500, PipelineSize: 250,
								})
								if redemErr != nil {
									fmt.Println("error fetching redemption rules:", redemErr)
								} else {
									for _, redemption := range redemptions {
										var Airtime_EVC_PIN = ""
										var err error
										if redemption.Airtime_EVC_PIN != "" {
											Airtime_EVC_PIN, err = DecryptHexString(redemption.Airtime_EVC_PIN)
											if err != nil {
												fmt.Println("error in decrypting artime evc pin", err.Error())
											}
										}
										res, err := Uc.CGW.UC_GWClient.GetERDealerBalance(redemption.Airtime_EVC_Account, Airtime_EVC_PIN)
										if err == nil {
											log_entry.EVC_Account_Balance = float64(res.Balance)
										}
									}
								}
							}
							Uc.Write_Loyalty_Governance_log(log_entry)
						}
						log.Println(LOG_ID + " finished")
					}
				}
			} else {
				if exec == 1 {
					exec = 0
				}
			}
		}
	}
}

func (Uc *UserControl) Loyalty_Customer_Account_Daily_Snapshot() {
	exec := 0
	LOG_ID := "<<Loyalty Customer Account Daily Snapshot>>"
	for range time.Tick(time.Second * 1) {
		_CurrentDateTime := time.Now()
		_hr, _mi, _se := _CurrentDateTime.Clock()
		if _hr == 00 {
			if _mi == 00 {
				if _se < 60 {
					if exec == 0 {
						exec = 1
						log.Println(LOG_ID + " triggered")
						// TODO: CollectionSnapshot not yet implemented in mongox migration.
						// Original: snapshot Col_Customer_Loyalty_Account to dynamic DB+collection by date.
						log.Println(LOG_ID + " snapshot skipped (pending mongox implementation)")
						log.Println(LOG_ID + " finished")
					}
				}
			} else {
				if exec == 1 {
					exec = 0
				}
			}
		}
	}
}

func (Uc *UserControl) Loyalty_Status_Expiry_Daily_Process() {
	if !Configuration.ISLoyaltyOptIn {
		return
	}
	LOG_ID := "<<Loyalty Status Expiry Daily check>>"

	for {
		now := time.Now()

		// Next run today at 00:10
		nextRun := time.Date(
			now.Year(), now.Month(), now.Day(),
			00, 13, 0, 0,
			now.Location(),
		)

		// If already past 00:10 â†’ schedule for tomorrow
		if now.After(nextRun) {
			nextRun = nextRun.Add(24 * time.Hour)
		}

		// Sleep until next run
		time.Sleep(time.Until(nextRun))
		log.Println(LOG_ID + " triggered")
		graceDays := Configuration.ISLoyaltyOptOutGracePeriodDays
		cutoffDate := time.Now().AddDate(0, 0, -graceDays)
		pipeline := bson.A{bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "Opt_Status", Value: bson.D{{Key: "$eq", Value: "OptedOut"}}},
				{Key: "Opt_Status_Date", Value: bson.D{{Key: "$lte", Value: cutoffDate}}},
			}},
		}}
		MongoDB_DB_Name := "Loyalty_DB"
		collName := "Col_Customer_Loyalty_Account"
		// Fetch the collection
		collection := Uc.LoyaltyMongoClient.Mongo.Database(MongoDB_DB_Name).Collection(collName)

		cursor, err := collection.Aggregate(context.Background(), pipeline)
		if err != nil {
			log.Printf("Aggregation failed for %s: %v", collName, err)
			continue
		}

		for cursor.Next(context.Background()) {
			var doc Customer_Loyalty_Account
			if err := cursor.Decode(&doc); err != nil {
				log.Printf("Failed decoding result for %s: %v", collName, err)
				continue
			}
			var expiry_log Loyalty_Full_Expiry_Log
			expiry_log.ExpiryTime = time.Now()
			expiry_log.MSISDN = doc.Key
			expiry_log.Opening_Awarded_Points = doc.Awarded_Points
			expiry_log.Opening_Redeemed_Points = doc.Redeemed_Points
			expiry_log.Opening_Available_Points = doc.Available_Points
			expiry_log.Opening_OutStanding_Points = doc.Outstanding_fraction_points
			expiry_log.Opening_Expired_Points = doc.Expired_Points
			expiry_log.OpeningLoyaltyLevel = doc.Loyalty_Level_Key
			expiry_log.EndLoyaltyLevel = doc.Loyalty_Level_Key
			expiry_log.Last_OptOut = doc.Last_Opt_Status_Date
			expiry_log.Grace_Period_Given_Days = graceDays
			expiry_log.ExpiryReason = "Opt out grace period reached"
			expiry_log.ExpiryAmount = doc.Available_Points
			var monthly_expiry_log Loyalty_Monthly_Expiry_log
			monthly_expiry_log.OpeningLoyaltyLevel = doc.Loyalty_Level_Key
			monthly_expiry_log.EndLoyaltyLevel = doc.Loyalty_Level_Key
			monthly_expiry_log.Expiry_Rules_Key = "Opt Out Expiry"
			monthly_expiry_log.MSISDN = doc.Key
			monthly_expiry_log.ExpiryTime = time.Now()
			for _, pointKey := range doc.Points_Detail_Keys {
				pointsDetail, err := Uc.Customer_Loyalty_Account_Points_Details_Get(pointKey)
				if err != nil {
					expiry_log.ExpiryStatus = "failed"
					expiry_log.ExpiryStatusDescription = "points expiry rules not found"
					Uc.Write_Loyalty_Full_Expiry_log(expiry_log)
					// <-chan_PointsExpiry_Controler
					return
				}
				monthly_expiry_log.Year_Month = pointsDetail[0].Year_Month
				monthly_expiry_log.Opening_Awarded_Points = pointsDetail[0].Awarded_Points
				monthly_expiry_log.Opening_Redeemed_Points = pointsDetail[0].Redeemed_Points
				monthly_expiry_log.Opening_Available_Points = pointsDetail[0].Available_Points
				monthly_expiry_log.Opening_Expired_Points = pointsDetail[0].Expired_Points

				pointsDetail[0].Expired_Points = pointsDetail[0].Available_Points
				pointsDetail[0].Available_Points = 0
				pointsDetail[0].Expiry_Date = time.Now()
				{
					putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
					if _, putErr := Mdb_Customer_Loyalty_Account_Points_Detail.Coll.UpdateOne(putCtx, bson.M{"Key": pointsDetail[0].Key}, bson.M{"$set": pointsDetail[0]}, options.UpdateOne().SetUpsert(true)); putErr != nil {
						log.Println("Mdb_Customer_Loyalty_Account_Points_Detail upsert error:", putErr)
					}
					if putSetErr := redisx.SetJSON(putCtx, RedisClient, pointsDetail[0].RedisKey(), pointsDetail[0]); putSetErr != nil {
						log.Println("redisx.SetJSON Customer_Loyalty_Account_Points_Detail error:", putSetErr)
					}
					putCancel()
				}

				monthly_expiry_log.End_Expired_Points = pointsDetail[0].Expired_Points
				monthly_expiry_log.ExpiryTime = time.Now()
				monthly_expiry_log.End_Available_Points = 0
				monthly_expiry_log.End_Awarded_Points = pointsDetail[0].Awarded_Points
				monthly_expiry_log.End_Redeemed_Points = pointsDetail[0].Redeemed_Points
				//check level downgrade

				monthly_expiry_log.ExpiryStatus = "successful"
				monthly_expiry_log.ExpiryStatusDescription = ""
				// Uc.Write_Loyalty_Monthly_Expiry_log(monthly_expiry_log)

				// }
				// <-chan_PointsExpiry_Controler
			}
			expired_Points := doc.Available_Points
			expiry_log.End_Expired_Points = 0
			doc.Expired_Points = 0
			doc.Available_Points = 0
			doc.Awarded_Points = 0
			doc.Redeemed_Points = 0
			doc.Outstanding_fraction_points = 0
			doc.Opt_Status = "OptedOutExpired"
			//update governance expiry
			Uc.Loyalty_Governance_Status_Expiry_Points_Credit(expired_Points)
			//update logs
			expiry_log.End_Awarded_Points = 0
			expiry_log.End_Redeemed_Points = 0
			expiry_log.End_Available_Points = 0
			expiry_log.End_Outstanding_Points = 0
			{
				putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
				if _, putErr := Mdb_Customer_Loyalty_Account.Coll.UpdateOne(putCtx, bson.M{"Key": doc.Key}, bson.M{"$set": doc}, options.UpdateOne().SetUpsert(true)); putErr != nil {
					log.Println("Mdb_Customer_Loyalty_Account upsert error:", putErr)
				}
				if putSetErr := redisx.SetJSON(putCtx, RedisClient, doc.RedisKey(), doc); putSetErr != nil {
					log.Println("redisx.SetJSON Customer_Loyalty_Account error:", putSetErr)
				}
				putCancel()
			}
			var New_Loyalty_Level Loyalty_Level
			loyaltyLevels, levelsErr := redisx.GetAllJSONByPattern[Loyalty_Level](context.Background(), RedisClient, redisx.ScanJSONOptions{
				Pattern: "Loyalty_Level:*", ScanCount: 500, PipelineSize: 250,
			})
			if levelsErr != nil {
				fmt.Println("error fetching loyalty levels:", levelsErr)
			} else {
				for _, loyalty_Level := range loyaltyLevels {
					//evaluate
					if doc.Awarded_Points >= loyalty_Level.Min_Accumulated_Points && doc.Awarded_Points < loyalty_Level.Max_Accumulated_Points {
						New_Loyalty_Level = loyalty_Level
						if New_Loyalty_Level.Key == doc.Loyalty_Level_Key {
							//===>> no level change
							break
						} else {
							current_level, lvlErr := redisx.GetJSON[Loyalty_Level](context.Background(), RedisClient, Loyalty_Level{Key: doc.Loyalty_Level_Key}.RedisKey())
							if lvlErr != nil {
								fmt.Println("current level is invalid")
							}
							if New_Loyalty_Level.Min_Accumulated_Points < current_level.Min_Accumulated_Points &&
								New_Loyalty_Level.Max_Accumulated_Points < current_level.Max_Accumulated_Points {
								//Downgrade level
								expiry_log.EndLoyaltyLevel = New_Loyalty_Level.Key
								doc.Previous_Loyalty_Level_Key = doc.Loyalty_Level_Key
								doc.Previous_Loyalty_Level_Date = doc.Loyalty_Level_Date
								doc.Loyalty_Level_Key = New_Loyalty_Level.Key
								doc.Loyalty_Level_Date = time.Now()
								doc.Loyalty_Level_Direction = "Downgrade"
								doc.Loyalty_Level_SetBy = "Opt Expiry Process"
								{
									putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
									if _, putErr := Mdb_Customer_Loyalty_Account.Coll.UpdateOne(putCtx, bson.M{"Key": doc.Key}, bson.M{"$set": doc}, options.UpdateOne().SetUpsert(true)); putErr != nil {
										log.Println("Mdb_Customer_Loyalty_Account upsert error:", putErr)
									}
									if putSetErr := redisx.SetJSON(putCtx, RedisClient, doc.RedisKey(), doc); putSetErr != nil {
										log.Println("redisx.SetJSON Customer_Loyalty_Account error:", putSetErr)
									}
									putCancel()
								}
							}
						}
					}
				}
			}
			expiry_log.ExpiryStatus = "successful"
			expiry_log.ExpiryStatusDescription = "opt expiry"
			Uc.Write_Loyalty_Full_Expiry_log(expiry_log)
		}

		// yesterday := time.Now().AddDate(0, 0, -1)
		// YYYY, MM, _, DD, _, _, _ := GetTimeParts(yesterday)
		// Db := DAO_Loyalty_AccountCreditPoints_log.DB + "_" + YYYY + MM
		// Col := DAO_Customer_Loyalty_Account.Collection + "_" + DD
		// err := DAO_Customer_Loyalty_Account.CollectionSnapshot(Db, Col)
		// if err != nil {
		// 	log.Println("error while taking a snapshot from customer account collection", err)
		// }
		log.Println(LOG_ID + " finished")
	}
}

// ***********************************************************************
// Loyalty Level functions
// ***********************************************************************
func (Uc *UserControl) Loyalty_Level_Add(Login string, request Loyalty_Level_AddRequest) (Id int64, err error) {
	//check key if filled and if already used
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	//check if key already used
	{
		chkCtx, chkCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, chkErr := redisx.GetJSON[Loyalty_Level](chkCtx, RedisClient, Loyalty_Level{Key: request.Key}.RedisKey())
		chkCancel()
		if chkErr == nil {
			err = errors.New("key already exist")
			return Id, err
		}
	}
	loyaltyLevels, levelsErr := redisx.GetAllJSONByPattern[Loyalty_Level](context.Background(), RedisClient, redisx.ScanJSONOptions{
		Pattern: "Loyalty_Level:*", ScanCount: 500, PipelineSize: 250,
	})
	if levelsErr != nil {
		fmt.Println("error fetching loyalty levels:", levelsErr)
	} else {
		for _, entry := range loyaltyLevels {
			if request.Min_Accumulated_Points <= entry.Max_Accumulated_Points && request.Max_Accumulated_Points >= entry.Min_Accumulated_Points {
				err = errors.New("Points intersect with " + entry.Key + " level")
				return Id, err
			}
		}
	}
	//Prepare new entry
	var NewEntry Loyalty_Level
	{
		aiCtx, aiCancel := context.WithTimeout(context.Background(), 10*time.Second)
		NewEntry.Level_Id, err = redisx.NextAutoIncrementID(aiCtx, RedisClient, Mdb_Loyalty_AutoIncrement.Coll, "Loyalty_Level-Id", redisx.NextIDOptions{
			RedisBase: "AutoIncrement", MongoRetries: 3, RetryBackoff: 500 * time.Millisecond,
		})
		aiCancel()
		if err != nil {
			return Id, err
		}
		Id = NewEntry.Level_Id
	}
	NewEntry.Key = request.Key
	NewEntry.Description = request.Description
	NewEntry.Min_Accumulated_Points = request.Min_Accumulated_Points
	NewEntry.Max_Accumulated_Points = request.Max_Accumulated_Points
	NewEntry.EnableRedeem = request.EnableRedeem
	NewEntry.DowngradeToLevel_Key = request.DowngradeToLevel_Key
	NewEntry.Seniority_Levels = request.Seniority_Levels
	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Loyalty_Level.Coll.UpdateOne(putCtx, bson.M{"Key": NewEntry.Key}, bson.M{"$set": NewEntry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			putCancel()
			err = putErr
			return Id, err
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, NewEntry.RedisKey(), NewEntry); putSetErr != nil {
			log.Println("redisx.SetJSON Loyalty_Level error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Level",
		Event_ActionType:   "Add",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  NewEntry,
	})
	return Id, nil
}

func (Uc *UserControl) Loyalty_Level_Edit(Login string, request Loyalty_Level_EditRequest) (Id int64, err error) {
	//check and validate
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	entry, entryErr := redisx.GetJSON[Loyalty_Level](context.Background(), RedisClient, Loyalty_Level{Key: request.Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("key is not created")
		return Id, err
	}
	if entryErr != nil {
		return Id, entryErr
	}
	if entry.Level_Id != request.Level_Id {
		return Id, errors.New("id is not matching")
	}
	Current_Entry := entry
	loyaltyLevels, levelsErr := redisx.GetAllJSONByPattern[Loyalty_Level](context.Background(), RedisClient, redisx.ScanJSONOptions{
		Pattern: "Loyalty_Level:*", ScanCount: 500, PipelineSize: 250,
	})
	if levelsErr != nil {
		fmt.Println("error fetching loyalty levels:", levelsErr)
	} else {
		for _, lvl := range loyaltyLevels {
			if request.Key != lvl.Key && request.Min_Accumulated_Points <= lvl.Max_Accumulated_Points && request.Max_Accumulated_Points >= lvl.Min_Accumulated_Points {
				err = errors.New("Points intersect with " + lvl.Key + " level")
				return Id, err
			}
		}
	}
	//Prepare new entry
	entry.Key = request.Key
	entry.Description = request.Description
	entry.Min_Accumulated_Points = request.Min_Accumulated_Points
	entry.Max_Accumulated_Points = request.Max_Accumulated_Points
	entry.EnableRedeem = request.EnableRedeem
	entry.DowngradeToLevel_Key = request.DowngradeToLevel_Key
	entry.Seniority_Levels = request.Seniority_Levels
	if request.NewKey != "" {
		if request.NewKey != request.Key {
			//delete old
			{
				delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
				if _, delErr := Mdb_Loyalty_Level.Coll.DeleteOne(delCtx, bson.M{"Key": request.Key}); delErr != nil {
					log.Println("Mdb_Loyalty_Level DeleteOne error:", delErr)
				}
				if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Loyalty_Level{Key: request.Key}.RedisKey()); delRedisErr != nil {
					log.Println("redisx.DelJSON Loyalty_Level error:", delRedisErr)
				}
				delCancel()
			}
			//update key
			entry.Key = request.NewKey
		}
	}
	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Loyalty_Level.Coll.UpdateOne(putCtx, bson.M{"Key": entry.Key}, bson.M{"$set": entry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			putCancel()
			err = putErr
			return Id, err
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, entry.RedisKey(), entry); putSetErr != nil {
			log.Println("redisx.SetJSON Loyalty_Level error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Level",
		Event_ActionType:   "Edit",
		Event_Description:  "",
		Event_Entry_Before: Current_Entry,
		Event_Entry_After:  entry,
	})
	return Id, nil
}

func (Uc *UserControl) Loyalty_Level_Get(Key string) (entries []Loyalty_Level, err error) {
	if Key == "" {
		entries, err = redisx.GetAllJSONByPattern[Loyalty_Level](context.Background(), RedisClient, redisx.ScanJSONOptions{
			Pattern: "Loyalty_Level:*", ScanCount: 500, PipelineSize: 250,
		})
		if err != nil {
			return entries, err
		}
		return entries, nil
	} else {
		entry, entryErr := redisx.GetJSON[Loyalty_Level](context.Background(), RedisClient, Loyalty_Level{Key: Key}.RedisKey())
		if redisx.IsNil(entryErr) {
			err = errors.New("key does not exist")
			return entries, err
		}
		if entryErr != nil {
			return entries, entryErr
		}
		entries = append(entries, entry)
		return entries, nil
	}
}

func (Uc *UserControl) Loyalty_Level_GetPaginated(Page, Limit int64) (entries []Loyalty_Level, err error) {
	if Page < 1 {
		return entries, errors.New("invalid page")
	}
	if Limit < 1 || Limit > 50000 {
		return entries, errors.New("invalid limit (accept value between 1 and 50000)")
	}

	pgCtx, pgCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer pgCancel()
	result, err := mongox.FindPage[Loyalty_Level](pgCtx, Mdb_Loyalty_Level.Coll, bson.M{}, mongox.PageRequest{
		Page: Page, PageSize: Limit, MaxPageSize: 50000,
	})
	if err != nil {
		return entries, err
	}
	entries = result.Items
	return entries, nil

}

func (Uc *UserControl) Loyalty_Level_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry, entryErr := redisx.GetJSON[Loyalty_Level](context.Background(), RedisClient, Loyalty_Level{Key: Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("entry does not exist")
		return err
	}
	if entryErr != nil {
		return entryErr
	}
	{
		delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, delErr := Mdb_Loyalty_Level.Coll.DeleteOne(delCtx, bson.M{"Key": Key}); delErr != nil {
			log.Println("Mdb_Loyalty_Level DeleteOne error:", delErr)
		}
		if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Loyalty_Level{Key: Key}.RedisKey()); delRedisErr != nil {
			log.Println("redisx.DelJSON Loyalty_Level error:", delRedisErr)
		}
		delCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Level",
		Event_ActionType:   "Delete",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  entry,
	})
	return nil
}

// ***********************************************************************
// Loyalty Seniiority Level functions
// ***********************************************************************
func (Uc *UserControl) Loyalty_Seniority_Level_Add(Login string, request Loyalty_Seniority_Level) (Id int64, err error) {
	//check key if filled and if already used
	if request.Name == "" {
		err = errors.New("name cannot be empty")
		return Id, err
	}
	//check if key already used
	{
		chkCtx, chkCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, chkErr := redisx.GetJSON[Loyalty_Seniority_Level](chkCtx, RedisClient, Loyalty_Seniority_Level{Key: request.Key}.RedisKey())
		chkCancel()
		if chkErr == nil {
			err = errors.New("key already exist")
			return Id, err
		}
	}
	if request.AON_From > request.AON_Till {
		err = errors.New("invalid AON values")
		return Id, err
	}
	seniorityLevels, seniorityErr := redisx.GetAllJSONByPattern[Loyalty_Seniority_Level](context.Background(), RedisClient, redisx.ScanJSONOptions{
		Pattern: "Loyalty_Seniority_Level:*", ScanCount: 500, PipelineSize: 250,
	})
	if seniorityErr != nil {
		fmt.Println("error fetching seniority levels:", seniorityErr)
	} else {
		for _, entry := range seniorityLevels {
			if request.AON_From <= entry.AON_Till && request.AON_Till >= entry.AON_From {
				err = errors.New("AON intersect with " + entry.Key + " AON")
				return Id, err
			}
		}
	}
	//Prepare new entry
	var NewEntry Loyalty_Seniority_Level
	{
		aiCtx, aiCancel := context.WithTimeout(context.Background(), 10*time.Second)
		NewEntry.Seniority_Id, err = redisx.NextAutoIncrementID(aiCtx, RedisClient, Mdb_Loyalty_AutoIncrement.Coll, "Loyalty_Seniority_Level-Id", redisx.NextIDOptions{
			RedisBase: "AutoIncrement", MongoRetries: 3, RetryBackoff: 500 * time.Millisecond,
		})
		aiCancel()
		if err != nil {
			return Id, err
		}
		Id = NewEntry.Seniority_Id
	}
	convertedId := strconv.FormatInt(Id, 10)
	NewEntry.Key = convertedId
	NewEntry.Name = request.Name
	NewEntry.Description = request.Description
	NewEntry.AON_From = request.AON_From
	NewEntry.AON_Till = request.AON_Till

	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Loyalty_Seniority_Level.Coll.UpdateOne(putCtx, bson.M{"Key": NewEntry.Key}, bson.M{"$set": NewEntry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			putCancel()
			err = putErr
			return Id, err
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, NewEntry.RedisKey(), NewEntry); putSetErr != nil {
			log.Println("redisx.SetJSON Loyalty_Seniority_Level error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Seniority_Level",
		Event_ActionType:   "Add",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  NewEntry,
	})
	return Id, nil
}

func (Uc *UserControl) Loyalty_Seniority_Level_Edit(Login string, request Loyalty_Seniority_Level) (Id int64, err error) {
	//check and validate
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	entry, entryErr := redisx.GetJSON[Loyalty_Seniority_Level](context.Background(), RedisClient, Loyalty_Seniority_Level{Key: request.Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("key is not created")
		return Id, err
	}
	if entryErr != nil {
		return Id, entryErr
	}
	if entry.Seniority_Id != request.Seniority_Id {
		return Id, errors.New("id is not matching")
	}
	if request.AON_From > request.AON_Till {
		err = errors.New("invalid AON values")
		return Id, err
	}
	Current_Entry := entry
	seniorityLevels, seniorityErr := redisx.GetAllJSONByPattern[Loyalty_Seniority_Level](context.Background(), RedisClient, redisx.ScanJSONOptions{
		Pattern: "Loyalty_Seniority_Level:*", ScanCount: 500, PipelineSize: 250,
	})
	if seniorityErr != nil {
		fmt.Println("error fetching seniority levels:", seniorityErr)
	} else {
		for _, lvl := range seniorityLevels {
			if request.Key != lvl.Key && request.AON_From <= lvl.AON_Till && request.AON_Till >= lvl.AON_From {
				err = errors.New("AON intersect with " + lvl.Key + " AON")
				return Id, err
			}
		}
	}
	//Prepare new entry
	entry.Name = request.Name
	entry.Description = request.Description
	entry.AON_From = request.AON_From
	entry.AON_Till = request.AON_Till

	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Loyalty_Seniority_Level.Coll.UpdateOne(putCtx, bson.M{"Key": entry.Key}, bson.M{"$set": entry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			putCancel()
			err = putErr
			return Id, err
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, entry.RedisKey(), entry); putSetErr != nil {
			log.Println("redisx.SetJSON Loyalty_Seniority_Level error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Seniority_Level",
		Event_ActionType:   "Edit",
		Event_Description:  "",
		Event_Entry_Before: Current_Entry,
		Event_Entry_After:  entry,
	})
	return Id, nil
}

func (Uc *UserControl) Loyalty_Seniority_Level_Get(Key string) (entries []Loyalty_Seniority_Level, err error) {
	if Key == "" {
		entries, err = redisx.GetAllJSONByPattern[Loyalty_Seniority_Level](context.Background(), RedisClient, redisx.ScanJSONOptions{
			Pattern: "Loyalty_Seniority_Level:*", ScanCount: 500, PipelineSize: 250,
		})
		if err != nil {
			return entries, err
		}
		return entries, nil
	} else {
		entry, entryErr := redisx.GetJSON[Loyalty_Seniority_Level](context.Background(), RedisClient, Loyalty_Seniority_Level{Key: Key}.RedisKey())
		if redisx.IsNil(entryErr) {
			err = errors.New("key does not exist")
			return entries, err
		}
		if entryErr != nil {
			return entries, entryErr
		}
		entries = append(entries, entry)
		return entries, nil
	}
}

func (Uc *UserControl) Loyalty_Seniority_Level_GetPaginated(Page, Limit int64) (entries []Loyalty_Seniority_Level, err error) {
	if Page < 1 {
		return entries, errors.New("invalid page")
	}
	if Limit < 1 || Limit > 50000 {
		return entries, errors.New("invalid limit (accept value between 1 and 50000)")
	}

	pgCtx, pgCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer pgCancel()
	result, err := mongox.FindPage[Loyalty_Seniority_Level](pgCtx, Mdb_Loyalty_Seniority_Level.Coll, bson.M{}, mongox.PageRequest{
		Page: Page, PageSize: Limit, MaxPageSize: 50000,
	})
	if err != nil {
		return entries, err
	}
	entries = result.Items
	return entries, nil

}

func (Uc *UserControl) Loyalty_Seniority_Level_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry, entryErr := redisx.GetJSON[Loyalty_Seniority_Level](context.Background(), RedisClient, Loyalty_Seniority_Level{Key: Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("entry does not exist")
		return err
	}
	if entryErr != nil {
		return entryErr
	}
	{
		delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, delErr := Mdb_Loyalty_Seniority_Level.Coll.DeleteOne(delCtx, bson.M{"Key": Key}); delErr != nil {
			log.Println("Mdb_Loyalty_Seniority_Level DeleteOne error:", delErr)
		}
		if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Loyalty_Seniority_Level{Key: Key}.RedisKey()); delRedisErr != nil {
			log.Println("redisx.DelJSON Loyalty_Seniority_Level error:", delRedisErr)
		}
		delCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Seniority_Level",
		Event_ActionType:   "Delete",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  entry,
	})
	return nil
}

// ***********************************************************************
// Loyalty_Account_Segment functions
// ***********************************************************************
func (Uc *UserControl) Loyalty_Account_Segment_Add(Login string, request Loyalty_Account_Segment_AddRequest) (Id int64, err error) {
	//check key if filled and if already used
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	//check if key already used
	{
		chkCtx, chkCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, chkErr := redisx.GetJSON[Loyalty_Account_Segment](chkCtx, RedisClient, Loyalty_Account_Segment{Key: request.Key}.RedisKey())
		chkCancel()
		if chkErr == nil {
			err = errors.New("key already exist")
			return Id, err
		}
	}

	//Prepare new entry
	var NewEntry Loyalty_Account_Segment
	{
		aiCtx, aiCancel := context.WithTimeout(context.Background(), 10*time.Second)
		NewEntry.Segment_Id, err = redisx.NextAutoIncrementID(aiCtx, RedisClient, Mdb_Loyalty_AutoIncrement.Coll, "Loyalty_Account_Segment-Id", redisx.NextIDOptions{
			RedisBase: "AutoIncrement", MongoRetries: 3, RetryBackoff: 500 * time.Millisecond,
		})
		aiCancel()
		if err != nil {
			return Id, err
		}
		Id = NewEntry.Segment_Id
	}
	NewEntry.Key = request.Key
	NewEntry.Description = request.Description
	NewEntry.Amount_From = request.Amount_From
	NewEntry.Amount_Till = request.Amount_Till
	NewEntry.AON_From = request.AON_From
	NewEntry.AON_Till = request.AON_Till

	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Loyalty_Account_Segment.Coll.UpdateOne(putCtx, bson.M{"Key": NewEntry.Key}, bson.M{"$set": NewEntry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			putCancel()
			err = putErr
			return Id, err
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, NewEntry.RedisKey(), NewEntry); putSetErr != nil {
			log.Println("redisx.SetJSON Loyalty_Account_Segment error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Account_Segment",
		Event_ActionType:   "Add",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  NewEntry,
	})
	return Id, nil
}

func (Uc *UserControl) Loyalty_Account_Segment_Edit(Login string, request Loyalty_Account_Segment_EditRequest) (Id int64, err error) {
	//check and validate
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	entry, entryErr := redisx.GetJSON[Loyalty_Account_Segment](context.Background(), RedisClient, Loyalty_Account_Segment{Key: request.Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("key is not created")
		return Id, err
	}
	if entryErr != nil {
		return Id, entryErr
	}
	if entry.Segment_Id != request.Segment_Id {
		return Id, errors.New("id is not matching")
	}
	Current_Entry := entry
	//Prepare new entry
	entry.Key = request.Key
	entry.Description = request.Description
	entry.Amount_From = request.Amount_From
	entry.Amount_Till = request.Amount_Till
	entry.AON_From = request.AON_From
	entry.AON_Till = request.AON_Till
	if request.NewKey != "" {
		if request.NewKey != request.Key {
			//delete old
			{
				delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
				if _, delErr := Mdb_Loyalty_Account_Segment.Coll.DeleteOne(delCtx, bson.M{"Key": request.Key}); delErr != nil {
					log.Println("Mdb_Loyalty_Account_Segment DeleteOne error:", delErr)
				}
				if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Loyalty_Account_Segment{Key: request.Key}.RedisKey()); delRedisErr != nil {
					log.Println("redisx.DelJSON Loyalty_Account_Segment error:", delRedisErr)
				}
				delCancel()
			}
			//update key
			entry.Key = request.NewKey
		}
	}
	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Loyalty_Account_Segment.Coll.UpdateOne(putCtx, bson.M{"Key": entry.Key}, bson.M{"$set": entry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			putCancel()
			err = putErr
			return Id, err
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, entry.RedisKey(), entry); putSetErr != nil {
			log.Println("redisx.SetJSON Loyalty_Account_Segment error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Account_Segment",
		Event_ActionType:   "Edit",
		Event_Description:  "",
		Event_Entry_Before: Current_Entry,
		Event_Entry_After:  entry,
	})
	return Id, nil
}

func (Uc *UserControl) Loyalty_Account_Segment_Get(Key string) (entries []Loyalty_Account_Segment, err error) {
	if Key == "" {
		entries, err = redisx.GetAllJSONByPattern[Loyalty_Account_Segment](context.Background(), RedisClient, redisx.ScanJSONOptions{
			Pattern: "Loyalty_Account_Segment:*", ScanCount: 500, PipelineSize: 250,
		})
		if err != nil {
			return entries, err
		}
		return entries, nil
	} else {
		entry, entryErr := redisx.GetJSON[Loyalty_Account_Segment](context.Background(), RedisClient, Loyalty_Account_Segment{Key: Key}.RedisKey())
		if redisx.IsNil(entryErr) {
			err = errors.New("key does not exist")
			return entries, err
		}
		if entryErr != nil {
			return entries, entryErr
		}
		entries = append(entries, entry)
		return entries, nil
	}
}

func (Uc *UserControl) Loyalty_Account_Segment_GetPaginated(Page, Limit int64) (entries []Loyalty_Account_Segment, err error) {
	if Page < 1 {
		return entries, errors.New("invalid page")
	}
	if Limit < 1 || Limit > 50000 {
		return entries, errors.New("invalid limit (accept value between 1 and 50000)")
	}

	pgCtx, pgCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer pgCancel()
	result, err := mongox.FindPage[Loyalty_Account_Segment](pgCtx, Mdb_Loyalty_Account_Segment.Coll, bson.M{}, mongox.PageRequest{
		Page: Page, PageSize: Limit, MaxPageSize: 50000,
	})
	if err != nil {
		return entries, err
	}
	entries = result.Items
	return entries, nil

}

func (Uc *UserControl) Loyalty_Account_Segment_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry, entryErr := redisx.GetJSON[Loyalty_Account_Segment](context.Background(), RedisClient, Loyalty_Account_Segment{Key: Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("entry does not exist")
		return err
	}
	if entryErr != nil {
		return entryErr
	}
	{
		delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, delErr := Mdb_Loyalty_Account_Segment.Coll.DeleteOne(delCtx, bson.M{"Key": Key}); delErr != nil {
			log.Println("Mdb_Loyalty_Account_Segment DeleteOne error:", delErr)
		}
		if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Loyalty_Account_Segment{Key: Key}.RedisKey()); delRedisErr != nil {
			log.Println("redisx.DelJSON Loyalty_Account_Segment error:", delRedisErr)
		}
		delCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Account_Segment",
		Event_ActionType:   "Delete",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  entry,
	})
	return nil
}

// ***********************************************************************
// Loyalty_Point_Earning_Rules functions
// ***********************************************************************
func (Uc *UserControl) Loyalty_Point_Earning_Rules_Add(Login string, request Loyalty_Point_Earning_Rules_AddRequest) (Id int64, err error) {
	//check key if filled and if already used
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	//check if key already used
	{
		chkCtx, chkCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, chkErr := redisx.GetJSON[Loyalty_Point_Earning_Rules](chkCtx, RedisClient, Loyalty_Point_Earning_Rules{Key: request.Key}.RedisKey())
		chkCancel()
		if chkErr == nil {
			err = errors.New("key already exist")
			return Id, err
		}
	}
	//check values
	if request.MainGSMBalance_Amount > 0 && request.MainGSMBalance_Amount < 1 {
		err = errors.New("invalid Main GSM value")
		return Id, err
	}
	if request.GSM_SC_Airtime_Award_Type == "Amount" {
		if request.GSM_SC_Airtime_Amount > 0 && request.GSM_SC_Airtime_Amount < 1 {
			err = errors.New("invalid SC Airtime value")
			return Id, err
		}
	} else if request.GSM_SC_Airtime_Award_Type != "Transaction" {
		err = errors.New("invalid SC Airtime Award Type")
		return Id, err
	} else {
		request.GSM_SC_Airtime_Amount = 0
	}
	if request.GSM_EVC_Airtime_Award_Type == "Amount" {
		if request.GSM_EVC_Airtime_Amount > 0 && request.GSM_EVC_Airtime_Amount < 1 {
			err = errors.New("invalid EVC Airtime value")
			return Id, err
		}
	} else if request.GSM_EVC_Airtime_Award_Type != "Transaction" {
		err = errors.New("invalid EVC Airtime Award Type")
		return Id, err
	} else {
		request.GSM_EVC_Airtime_Amount = 0
	}
	if request.GSM_EVC_Bundle_Award_Type == "Amount" {
		if request.GSM_EVC_Bundle_Amount > 0 && request.GSM_EVC_Bundle_Amount < 1 {
			err = errors.New("invalid EVC Bundle value")
			return Id, err
		}
	} else if request.GSM_EVC_Bundle_Award_Type != "Transaction" {
		err = errors.New("invalid EVC Bundle Award Type")
		return Id, err
	} else {
		request.GSM_EVC_Bundle_Amount = 0
	}
	if request.MM_P2P_Award_Type == "Amount" {
		if request.MM_P2P_Amount > 0 && request.MM_P2P_Amount < 1 {
			err = errors.New("invalid MM P2P value")
			return Id, err
		}
	} else if request.MM_P2P_Award_Type != "Transaction" {
		err = errors.New("invalid MM P2P Award Type")
		return Id, err
	} else {
		request.MM_P2P_Amount = 0
	}
	if request.MM_CASHIN_Award_Type == "Amount" {
		if request.MM_CASHIN_Amount > 0 && request.MM_CASHIN_Amount < 1 {
			err = errors.New("invalid MM CASHIN value")
			return Id, err
		}
	} else if request.MM_CASHIN_Award_Type != "Transaction" {
		err = errors.New("invalid MM CASHIN Award Type")
		return Id, err
	} else {
		request.MM_CASHIN_Amount = 0
	}
	if request.MM_CASHOUT_Award_Type == "Amount" {
		if request.MM_CASHOUT_Amount > 0 && request.MM_CASHOUT_Amount < 1 {
			err = errors.New("invalid MM CASHOUT value")
			return Id, err
		}
	} else if request.MM_CASHOUT_Award_Type != "Transaction" {
		err = errors.New("invalid MM CASHOUT Award Type")
		return Id, err
	} else {
		request.MM_CASHOUT_Amount = 0
	}
	if request.MM_MERCHPAY_Award_Type == "Amount" {
		if request.MM_MERCHPAY_Amount > 0 && request.MM_MERCHPAY_Amount < 1 {
			err = errors.New("invalid MM MERCHPAY value")
			return Id, err
		}
	} else if request.MM_MERCHPAY_Award_Type != "Transaction" {
		err = errors.New("invalid MM MERCHPAY Award Type")
		return Id, err
	} else {
		request.MM_MERCHPAY_Amount = 0
	}
	if request.MM_BILLPAY_Award_Type == "Amount" {
		if request.MM_BILLPAY_Amount > 0 && request.MM_BILLPAY_Amount < 1 {
			err = errors.New("invalid MM BILLPAY value")
			return Id, err
		}
	} else if request.MM_BILLPAY_Award_Type != "Transaction" {
		err = errors.New("invalid MM BILLPAY Award Type")
		return Id, err
	} else {
		request.MM_BILLPAY_Amount = 0
	}
	if request.MM_RC_Bundle_Award_Type == "Amount" {
		if request.MM_RC_Bundle_Amount > 0 && request.MM_RC_Bundle_Amount < 1 {
			err = errors.New("invalid MM recharge for self value")
			return Id, err
		}
	} else if request.MM_RC_Bundle_Award_Type != "Transaction" {
		err = errors.New("invalid MM recharge for self Award Type")
		return Id, err
	} else {
		request.MM_RC_Bundle_Amount = 0
	}
	if request.MM_RC_Airtime_Award_Type == "Amount" {
		if request.MM_RC_Airtime_Amount > 0 && request.MM_RC_Airtime_Amount < 1 {
			err = errors.New("invalid MM recharge for self value")
			return Id, err
		}
	} else if request.MM_RC_Airtime_Award_Type != "Transaction" {
		err = errors.New("invalid MM recharge for self Award Type")
		return Id, err
	} else {
		request.MM_RC_Airtime_Amount = 0
	}
	if request.MM_CTMMOREQ_Bundle_Award_Type == "Amount" {
		if request.MM_CTMMOREQ_Bundle_Amount > 0 && request.MM_CTMMOREQ_Bundle_Amount < 1 {
			err = errors.New("invalid MM recharge for others value")
			return Id, err
		}
	} else if request.MM_CTMMOREQ_Bundle_Award_Type != "Transaction" {
		err = errors.New("invalid MM recharge for others Award Type")
		return Id, err
	} else {
		request.MM_CTMMOREQ_Bundle_Amount = 0
	}
	if request.MM_CTMMOREQ_Airtime_Award_Type == "Amount" {
		if request.MM_CTMMOREQ_Airtime_Amount > 0 && request.MM_CTMMOREQ_Airtime_Amount < 1 {
			err = errors.New("invalid MM recharge for others value")
			return Id, err
		}
	} else if request.MM_CTMMOREQ_Airtime_Award_Type != "Transaction" {
		err = errors.New("invalid MM recharge for others Award Type")
		return Id, err
	} else {
		request.MM_CTMMOREQ_Airtime_Amount = 0
	}
	if request.MM_CBWREQ_Award_Type == "Amount" {
		if request.MM_CBWREQ_Amount > 0 && request.MM_CBWREQ_Amount < 1 {
			err = errors.New("invalid MM bank to wallet value")
			return Id, err
		}
	} else if request.MM_CBWREQ_Award_Type != "Transaction" {
		err = errors.New("invalid MM bank to wallet Award Type")
		return Id, err
	} else {
		request.MM_CBWREQ_Amount = 0
	}
	// if request.MM_Airtime_Award_Type == "Amount" {
	// 	if request.MM_Airtime_Amount > 0 && request.MM_Airtime_Amount < 1 {
	// 		err = errors.New("invalid MM Airtime value")
	// 		return Id, err
	// 	}
	// } else if request.MM_Airtime_Award_Type != "Transaction" {
	// 	err = errors.New("invalid MM Airtime Award Type")
	// 	return Id, err
	// } else {
	// 	request.MM_Airtime_Amount = 0
	// }
	// if request.MM_Bundle_Award_Type == "Amount" {
	// 	if request.MM_Bundle_Amount > 0 && request.MM_Bundle_Amount < 1 {
	// 		err = errors.New("invalid MM Bundle value")
	// 		return Id, err
	// 	}
	// } else if request.MM_Bundle_Award_Type != "Transaction" {
	// 	err = errors.New("invalid MM Bundle Award Type")
	// 	return Id, err
	// } else {
	// 	request.MM_Bundle_Amount = 0
	// }

	governanceEntries, governanceErr := redisx.GetAllJSONByPattern[Loyalty_Governance](context.Background(), RedisClient, redisx.ScanJSONOptions{
		Pattern: "Loyalty_Governance:*", ScanCount: 500, PipelineSize: 250,
	})
	if governanceErr != nil {
		fmt.Println("error fetching loyalty governance:", governanceErr)
	}
	if len(governanceEntries) > 0 {
		if governanceEntries[0].MaxAllowedPoints_PerTransaction < request.MM_BILLPAY_Points || governanceEntries[0].MaxAllowedPoints_PerTransaction < request.MM_CASHIN_Points || governanceEntries[0].MaxAllowedPoints_PerTransaction < request.MM_CASHOUT_Points || governanceEntries[0].MaxAllowedPoints_PerTransaction < request.MM_CBWREQ_Points || governanceEntries[0].MaxAllowedPoints_PerTransaction < request.MM_CTMMOREQ_Bundle_Points || governanceEntries[0].MaxAllowedPoints_PerTransaction < request.MM_CTMMOREQ_Airtime_Points || governanceEntries[0].MaxAllowedPoints_PerTransaction < request.MM_MERCHPAY_Points || governanceEntries[0].MaxAllowedPoints_PerTransaction < request.MM_P2P_Points || governanceEntries[0].MaxAllowedPoints_PerTransaction < request.MM_RC_Bundle_Points || governanceEntries[0].MaxAllowedPoints_PerTransaction < request.MM_RC_Airtime_Points || governanceEntries[0].MaxAllowedPoints_PerTransaction < request.MainGSMBalance_Points {
			err = errors.New("points can not exceed the maximum allowed points per transaction")
			return Id, err
		}
	}
	//Prepare new entry
	var NewEntry Loyalty_Point_Earning_Rules
	{
		aiCtx, aiCancel := context.WithTimeout(context.Background(), 10*time.Second)
		NewEntry.Earning_Rules_Id, err = redisx.NextAutoIncrementID(aiCtx, RedisClient, Mdb_Loyalty_AutoIncrement.Coll, "Loyalty_Point_Earning_Rules-Id", redisx.NextIDOptions{
			RedisBase: "AutoIncrement", MongoRetries: 3, RetryBackoff: 500 * time.Millisecond,
		})
		aiCancel()
		if err != nil {
			return Id, err
		}
		Id = NewEntry.Earning_Rules_Id
	}
	NewEntry.Key = request.Key
	NewEntry.Description = request.Description
	NewEntry.Welcome_Points = request.Welcome_Points
	NewEntry.MobileAppDaily_Notification = request.MobileAppDaily_Notification
	NewEntry.MobileAppDaily_Notification_Sender = request.MobileAppDaily_Notification_Sender
	NewEntry.MobileAppDaily_Notification_Text = request.MobileAppDaily_Notification_Text
	NewEntry.MobileAppDaily_Login = request.MobileAppDaily_Login
	NewEntry.Welcome_Notification = request.Welcome_Notification
	NewEntry.Welcome_Notification_Sender = request.Welcome_Notification_Sender
	NewEntry.Welcome_Notification_Text = request.Welcome_Notification_Text
	NewEntry.Rejoiner_Notification = request.Rejoiner_Notification
	NewEntry.Rejoiner_Notification_Sender = request.Rejoiner_Notification_Sender
	NewEntry.Rejoiner_Notification_Text = request.Rejoiner_Notification_Text
	NewEntry.Level_Change_Notification = request.Level_Change_Notification
	NewEntry.Level_Change_Notification_Sender = request.Level_Change_Notification_Sender
	NewEntry.Level_Change_Notification_Text = request.Level_Change_Notification_Text
	NewEntry.MainGSM_Notification = request.MainGSM_Notification
	NewEntry.MainGSM_Notification_Sender = request.MainGSM_Notification_Sender
	NewEntry.MainGSM_Notification_Text = request.MainGSM_Notification_Text
	NewEntry.MainGSMBalance_Amount = request.MainGSMBalance_Amount
	NewEntry.MainGSMBalance_Points = request.MainGSMBalance_Points
	NewEntry.GSM_SC_Notification = request.GSM_SC_Notification
	NewEntry.GSM_SC_Notification_Sender = request.GSM_SC_Notification_Sender
	NewEntry.GSM_SC_Notification_Text = request.GSM_SC_Notification_Text
	NewEntry.GSM_SC_Airtime_Award_Type = request.GSM_SC_Airtime_Award_Type
	NewEntry.GSM_SC_Airtime_Amount = request.GSM_SC_Airtime_Amount
	NewEntry.GSM_SC_Airtime_Points = request.GSM_SC_Airtime_Points
	NewEntry.GSM_EVC_Airtime_Notification = request.GSM_EVC_Airtime_Notification
	NewEntry.GSM_EVC_Airtime_Notification_Sender = request.GSM_EVC_Airtime_Notification_Sender
	NewEntry.GSM_EVC_Airtime_Notification_Text = request.GSM_EVC_Airtime_Notification_Text
	NewEntry.GSM_EVC_Airtime_Award_Type = request.GSM_EVC_Airtime_Award_Type
	NewEntry.GSM_EVC_Airtime_Amount = request.GSM_EVC_Airtime_Amount
	NewEntry.GSM_EVC_Airtime_Points = request.GSM_EVC_Airtime_Points
	NewEntry.GSM_EVC_Bundle_Notification = request.GSM_EVC_Bundle_Notification
	NewEntry.GSM_EVC_Bundle_Notification_Sender = request.GSM_EVC_Bundle_Notification_Sender
	NewEntry.GSM_EVC_Bundle_Notification_Text = request.GSM_EVC_Bundle_Notification_Text
	NewEntry.GSM_EVC_Bundle_Award_Type = request.GSM_EVC_Bundle_Award_Type
	NewEntry.GSM_EVC_Bundle_Amount = request.GSM_EVC_Bundle_Amount
	NewEntry.GSM_EVC_Bundle_Points = request.GSM_EVC_Bundle_Points
	NewEntry.MM_P2P_Notification = request.MM_P2P_Notification
	NewEntry.MM_P2P_Notification_Sender = request.MM_P2P_Notification_Sender
	NewEntry.MM_P2P_Notification_Text = request.MM_P2P_Notification_Text
	NewEntry.MM_P2P_Award_Type = request.MM_P2P_Award_Type
	NewEntry.MM_P2P_Amount = request.MM_P2P_Amount
	NewEntry.MM_P2P_Points = request.MM_P2P_Points
	NewEntry.MM_CASHIN_Notification = request.MM_CASHIN_Notification
	NewEntry.MM_CASHIN_Notification_Sender = request.MM_CASHIN_Notification_Sender
	NewEntry.MM_CASHIN_Notification_Text = request.MM_CASHIN_Notification_Text
	NewEntry.MM_CASHIN_Award_Type = request.MM_CASHIN_Award_Type
	NewEntry.MM_CASHIN_Amount = request.MM_CASHIN_Amount
	NewEntry.MM_CASHIN_Points = request.MM_CASHIN_Points
	NewEntry.MM_CASHOUT_Notification = request.MM_CASHOUT_Notification
	NewEntry.MM_CASHOUT_Notification_Sender = request.MM_CASHOUT_Notification_Sender
	NewEntry.MM_CASHOUT_Notification_Text = request.MM_CASHOUT_Notification_Text
	NewEntry.MM_CASHOUT_Award_Type = request.MM_CASHOUT_Award_Type
	NewEntry.MM_CASHOUT_Amount = request.MM_CASHOUT_Amount
	NewEntry.MM_CASHOUT_Points = request.MM_CASHOUT_Points
	NewEntry.MM_MERCHPAY_Notification = request.MM_MERCHPAY_Notification
	NewEntry.MM_MERCHPAY_Notification_Sender = request.MM_MERCHPAY_Notification_Sender
	NewEntry.MM_MERCHPAY_Notification_Text = request.MM_MERCHPAY_Notification_Text
	NewEntry.MM_MERCHPAY_Award_Type = request.MM_MERCHPAY_Award_Type
	NewEntry.MM_MERCHPAY_Amount = request.MM_MERCHPAY_Amount
	NewEntry.MM_MERCHPAY_Points = request.MM_MERCHPAY_Points
	NewEntry.MM_BILLPAY_Notification = request.MM_BILLPAY_Notification
	NewEntry.MM_BILLPAY_Notification_Sender = request.MM_BILLPAY_Notification_Sender
	NewEntry.MM_BILLPAY_Notification_Text = request.MM_BILLPAY_Notification_Text
	NewEntry.MM_BILLPAY_Award_Type = request.MM_BILLPAY_Award_Type
	NewEntry.MM_BILLPAY_Amount = request.MM_BILLPAY_Amount
	NewEntry.MM_BILLPAY_Points = request.MM_BILLPAY_Points
	NewEntry.Overwrite_MM_MERCHPAY_Notification = request.Overwrite_MM_MERCHPAY_Notification
	NewEntry.Overwrite_MM_MERCHPAY_Notification_Sender = request.Overwrite_MM_MERCHPAY_Notification_Sender
	NewEntry.Overwrite_MM_MERCHPAY_Notification_Text = request.Overwrite_MM_MERCHPAY_Notification_Text
	NewEntry.Earning_Rules_Overwrite_MM_MERCHPAY_Keys = request.Earning_Rules_Overwrite_MM_MERCHPAY_Keys
	NewEntry.Overwrite_MM_BILLPAY_Notification = request.Overwrite_MM_BILLPAY_Notification
	NewEntry.Overwrite_MM_BILLPAY_Notification_Sender = request.Overwrite_MM_BILLPAY_Notification_Sender
	NewEntry.Overwrite_MM_BILLPAY_Notification_Text = request.Overwrite_MM_BILLPAY_Notification_Text
	NewEntry.Earning_Rules_Overwrite_MM_BILLPAY_Keys = request.Earning_Rules_Overwrite_MM_BILLPAY_Keys
	NewEntry.MM_RC_Bundle_Notification = request.MM_RC_Bundle_Notification
	NewEntry.MM_RC_Bundle_Notification_Sender = request.MM_RC_Bundle_Notification_Sender
	NewEntry.MM_RC_Bundle_Notification_Text = request.MM_RC_Bundle_Notification_Text
	NewEntry.MM_RC_Bundle_Award_Type = request.MM_RC_Bundle_Award_Type
	NewEntry.MM_RC_Bundle_Amount = request.MM_RC_Bundle_Amount
	NewEntry.MM_RC_Bundle_Points = request.MM_RC_Bundle_Points
	NewEntry.MM_RC_Airtime_Notification = request.MM_RC_Airtime_Notification
	NewEntry.MM_RC_Airtime_Notification_Sender = request.MM_RC_Airtime_Notification_Sender
	NewEntry.MM_RC_Airtime_Notification_Text = request.MM_RC_Airtime_Notification_Text
	NewEntry.MM_RC_Airtime_Award_Type = request.MM_RC_Airtime_Award_Type
	NewEntry.MM_RC_Airtime_Amount = request.MM_RC_Airtime_Amount
	NewEntry.MM_RC_Airtime_Points = request.MM_RC_Airtime_Points
	NewEntry.MM_CTMMOREQ_Bundle_Notification = request.MM_CTMMOREQ_Bundle_Notification
	NewEntry.MM_CTMMOREQ_Bundle_Notification_Sender = request.MM_CTMMOREQ_Bundle_Notification_Sender
	NewEntry.MM_CTMMOREQ_Bundle_Notification_Text = request.MM_CTMMOREQ_Bundle_Notification_Text
	NewEntry.MM_CTMMOREQ_Bundle_Award_Type = request.MM_CTMMOREQ_Bundle_Award_Type
	NewEntry.MM_CTMMOREQ_Bundle_Amount = request.MM_CTMMOREQ_Bundle_Amount
	NewEntry.MM_CTMMOREQ_Bundle_Points = request.MM_CTMMOREQ_Bundle_Points
	NewEntry.MM_CTMMOREQ_Airtime_Notification = request.MM_CTMMOREQ_Airtime_Notification
	NewEntry.MM_CTMMOREQ_Airtime_Notification_Sender = request.MM_CTMMOREQ_Airtime_Notification_Sender
	NewEntry.MM_CTMMOREQ_Airtime_Notification_Text = request.MM_CTMMOREQ_Airtime_Notification_Text
	NewEntry.MM_CTMMOREQ_Airtime_Award_Type = request.MM_CTMMOREQ_Airtime_Award_Type
	NewEntry.MM_CTMMOREQ_Airtime_Amount = request.MM_CTMMOREQ_Airtime_Amount
	NewEntry.MM_CTMMOREQ_Airtime_Points = request.MM_CTMMOREQ_Airtime_Points
	NewEntry.MM_CBWREQ_Notification = request.MM_CBWREQ_Notification
	NewEntry.MM_CBWREQ_Notification_Sender = request.MM_CBWREQ_Notification_Sender
	NewEntry.MM_CBWREQ_Notification_Text = request.MM_CBWREQ_Notification_Text
	NewEntry.MM_CBWREQ_Award_Type = request.MM_CBWREQ_Award_Type
	NewEntry.MM_CBWREQ_Amount = request.MM_CBWREQ_Amount
	NewEntry.MM_CBWREQ_Points = request.MM_CBWREQ_Points
	NewEntry.MM_Airtime_Notification = request.MM_Airtime_Notification
	NewEntry.MM_Airtime_Notification_Sender = request.MM_Airtime_Notification_Sender
	NewEntry.MM_Airtime_Notification_Text = request.MM_Airtime_Notification_Text
	NewEntry.MM_Airtime_Award_Type = request.MM_Airtime_Award_Type
	NewEntry.MM_Airtime_Amount = request.MM_Airtime_Amount
	NewEntry.MM_Airtime_Points = request.MM_Airtime_Points
	NewEntry.MM_Bundle_Notification = request.MM_Bundle_Notification
	NewEntry.MM_Bundle_Notification_Sender = request.MM_Bundle_Notification_Sender
	NewEntry.MM_Bundle_Notification_Text = request.MM_Bundle_Notification_Text
	NewEntry.MM_Bundle_Award_Type = request.MM_Bundle_Award_Type
	NewEntry.MM_Bundle_Amount = request.MM_Bundle_Amount
	NewEntry.MM_Bundle_Points = request.MM_Bundle_Points

	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Loyalty_Point_Earning_Rules.Coll.UpdateOne(putCtx, bson.M{"Key": NewEntry.Key}, bson.M{"$set": NewEntry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			putCancel()
			err = putErr
			return Id, err
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, NewEntry.RedisKey(), NewEntry); putSetErr != nil {
			log.Println("redisx.SetJSON Loyalty_Point_Earning_Rules error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Point_Earning_Rules",
		Event_ActionType:   "Add",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  NewEntry,
	})
	return Id, nil
}

func (Uc *UserControl) Loyalty_Point_Earning_Rules_Edit(Login string, request Loyalty_Point_Earning_Rules_EditRequest) (Id int64, err error) {
	//check and validate
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	//check values
	if request.MainGSMBalance_Amount > 0 && request.MainGSMBalance_Amount < 1 {
		err = errors.New("invalid Main GSM value")
		return Id, err
	}
	if request.GSM_SC_Airtime_Award_Type == "Amount" {
		if request.GSM_SC_Airtime_Amount > 0 && request.GSM_SC_Airtime_Amount < 1 {
			err = errors.New("invalid SC Airtime value")
			return Id, err
		}
	} else if request.GSM_SC_Airtime_Award_Type != "Transaction" {
		err = errors.New("invalid SC Airtime Award Type")
		return Id, err
	} else {
		request.GSM_SC_Airtime_Amount = 0
	}
	if request.GSM_EVC_Airtime_Award_Type == "Amount" {
		if request.GSM_EVC_Airtime_Amount > 0 && request.GSM_EVC_Airtime_Amount < 1 {
			err = errors.New("invalid EVC Airtime value")
			return Id, err
		}
	} else if request.GSM_EVC_Airtime_Award_Type != "Transaction" {
		err = errors.New("invalid EVC Airtime Award Type")
		return Id, err
	} else {
		request.GSM_EVC_Airtime_Amount = 0
	}
	if request.GSM_EVC_Bundle_Award_Type == "Amount" {
		if request.GSM_EVC_Bundle_Amount > 0 && request.GSM_EVC_Bundle_Amount < 1 {
			err = errors.New("invalid EVC Bundle value")
			return Id, err
		}
	} else if request.GSM_EVC_Bundle_Award_Type != "Transaction" {
		err = errors.New("invalid EVC Bundle Award Type")
		return Id, err
	} else {
		request.GSM_EVC_Bundle_Amount = 0
	}
	if request.MM_P2P_Award_Type == "Amount" {
		if request.MM_P2P_Amount > 0 && request.MM_P2P_Amount < 1 {
			err = errors.New("invalid MM P2P value")
			return Id, err
		}
	} else if request.MM_P2P_Award_Type != "Transaction" {
		err = errors.New("invalid MM P2P Award Type")
		return Id, err
	} else {
		request.MM_P2P_Amount = 0
	}
	if request.MM_CASHIN_Award_Type == "Amount" {
		if request.MM_CASHIN_Amount > 0 && request.MM_CASHIN_Amount < 1 {
			err = errors.New("invalid MM CASHIN value")
			return Id, err
		}
	} else if request.MM_CASHIN_Award_Type != "Transaction" {
		err = errors.New("invalid MM CASHIN Award Type")
		return Id, err
	} else {
		request.MM_CASHIN_Amount = 0
	}
	if request.MM_CASHOUT_Award_Type == "Amount" {
		if request.MM_CASHOUT_Amount > 0 && request.MM_CASHOUT_Amount < 1 {
			err = errors.New("invalid MM CASHOUT value")
			return Id, err
		}
	} else if request.MM_CASHOUT_Award_Type != "Transaction" {
		err = errors.New("invalid MM CASHOUT Award Type")
		return Id, err
	} else {
		request.MM_CASHOUT_Amount = 0
	}
	if request.MM_MERCHPAY_Award_Type == "Amount" {
		if request.MM_MERCHPAY_Amount > 0 && request.MM_MERCHPAY_Amount < 1 {
			err = errors.New("invalid MM MERCHPAY value")
			return Id, err
		}
	} else if request.MM_MERCHPAY_Award_Type != "Transaction" {
		err = errors.New("invalid MM MERCHPAY Award Type")
		return Id, err
	} else {
		request.MM_MERCHPAY_Amount = 0
	}
	if request.MM_BILLPAY_Award_Type == "Amount" {
		if request.MM_BILLPAY_Amount > 0 && request.MM_BILLPAY_Amount < 1 {
			err = errors.New("invalid MM BILLPAY value")
			return Id, err
		}
	} else if request.MM_BILLPAY_Award_Type != "Transaction" {
		err = errors.New("invalid MM BILLPAY Award Type")
		return Id, err
	} else {
		request.MM_BILLPAY_Amount = 0
	}
	if request.MM_RC_Bundle_Award_Type == "Amount" {
		if request.MM_RC_Bundle_Amount > 0 && request.MM_RC_Bundle_Amount < 1 {
			err = errors.New("invalid MM recharge for self value")
			return Id, err
		}
	} else if request.MM_RC_Bundle_Award_Type != "Transaction" {
		err = errors.New("invalid MM recharge for self Award Type")
		return Id, err
	} else {
		request.MM_RC_Bundle_Amount = 0
	}
	if request.MM_RC_Airtime_Award_Type == "Amount" {
		if request.MM_RC_Airtime_Amount > 0 && request.MM_RC_Airtime_Amount < 1 {
			err = errors.New("invalid MM recharge for self value")
			return Id, err
		}
	} else if request.MM_RC_Airtime_Award_Type != "Transaction" {
		err = errors.New("invalid MM recharge for self Award Type")
		return Id, err
	} else {
		request.MM_RC_Airtime_Amount = 0
	}
	if request.MM_CTMMOREQ_Bundle_Award_Type == "Amount" {
		if request.MM_CTMMOREQ_Bundle_Amount > 0 && request.MM_CTMMOREQ_Bundle_Amount < 1 {
			err = errors.New("invalid MM recharge for others value")
			return Id, err
		}
	} else if request.MM_CTMMOREQ_Bundle_Award_Type != "Transaction" {
		err = errors.New("invalid MM recharge for others Award Type")
		return Id, err
	} else {
		request.MM_CTMMOREQ_Bundle_Amount = 0
	}
	if request.MM_CTMMOREQ_Airtime_Award_Type == "Amount" {
		if request.MM_CTMMOREQ_Airtime_Amount > 0 && request.MM_CTMMOREQ_Airtime_Amount < 1 {
			err = errors.New("invalid MM recharge for others value")
			return Id, err
		}
	} else if request.MM_CTMMOREQ_Airtime_Award_Type != "Transaction" {
		err = errors.New("invalid MM recharge for others Award Type")
		return Id, err
	} else {
		request.MM_CTMMOREQ_Airtime_Amount = 0
	}
	if request.MM_CBWREQ_Award_Type == "Amount" {
		if request.MM_CBWREQ_Amount > 0 && request.MM_CBWREQ_Amount < 1 {
			err = errors.New("invalid MM bank to wallet value")
			return Id, err
		}
	} else if request.MM_CBWREQ_Award_Type != "Transaction" {
		err = errors.New("invalid MM bank to wallet Award Type")
		return Id, err
	} else {
		request.MM_CBWREQ_Amount = 0
	}
	// if request.MM_Airtime_Award_Type == "Amount" {
	// 	if request.MM_Airtime_Amount > 0 && request.MM_Airtime_Amount < 1 {
	// 		err = errors.New("invalid MM Airtime value")
	// 		return Id, err
	// 	}
	// } else if request.MM_Airtime_Award_Type != "Transaction" {
	// 	err = errors.New("invalid MM Airtime Award Type")
	// 	return Id, err
	// } else {
	// 	request.MM_Airtime_Amount = 0
	// }
	// if request.MM_Bundle_Award_Type == "Amount" {
	// 	if request.MM_Bundle_Amount > 0 && request.MM_Bundle_Amount < 1 {
	// 		err = errors.New("invalid MM Bundle value")
	// 		return Id, err
	// 	}
	// } else if request.MM_Bundle_Award_Type != "Transaction" {
	// 	err = errors.New("invalid MM Bundle Award Type")
	// 	return Id, err
	// } else {
	// 	request.MM_Bundle_Amount = 0
	// }
	entry, entryErr := redisx.GetJSON[Loyalty_Point_Earning_Rules](context.Background(), RedisClient, Loyalty_Point_Earning_Rules{Key: request.Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("key is not created")
		return Id, err
	}
	if entryErr != nil {
		return Id, entryErr
	}
	if entry.Earning_Rules_Id != request.Earning_Rules_Id {
		return Id, errors.New("id is not matching")
	}
	governanceEntries, governanceErr := redisx.GetAllJSONByPattern[Loyalty_Governance](context.Background(), RedisClient, redisx.ScanJSONOptions{
		Pattern: "Loyalty_Governance:*", ScanCount: 500, PipelineSize: 250,
	})
	if governanceErr != nil {
		fmt.Println("error fetching loyalty governance:", governanceErr)
	}
	if len(governanceEntries) > 0 {
		if governanceEntries[0].MaxAllowedPoints_PerTransaction < request.MM_BILLPAY_Points || governanceEntries[0].MaxAllowedPoints_PerTransaction < request.MM_CASHIN_Points || governanceEntries[0].MaxAllowedPoints_PerTransaction < request.MM_CASHOUT_Points || governanceEntries[0].MaxAllowedPoints_PerTransaction < request.MM_CBWREQ_Points || governanceEntries[0].MaxAllowedPoints_PerTransaction < request.MM_CTMMOREQ_Bundle_Points || governanceEntries[0].MaxAllowedPoints_PerTransaction < request.MM_CTMMOREQ_Airtime_Points || governanceEntries[0].MaxAllowedPoints_PerTransaction < request.MM_MERCHPAY_Points || governanceEntries[0].MaxAllowedPoints_PerTransaction < request.MM_P2P_Points || governanceEntries[0].MaxAllowedPoints_PerTransaction < request.MM_RC_Bundle_Points || governanceEntries[0].MaxAllowedPoints_PerTransaction < request.MM_RC_Airtime_Points || governanceEntries[0].MaxAllowedPoints_PerTransaction < request.MainGSMBalance_Points {
			err = errors.New("points can not exceed the maximum allowed points per transaction")
			return Id, err
		}
	}
	Current_Entry := entry
	//Prepare new entry
	entry.Key = request.Key
	entry.Description = request.Description
	entry.Welcome_Points = request.Welcome_Points
	entry.Welcome_Points = request.Welcome_Points
	entry.MobileAppDaily_Notification = request.MobileAppDaily_Notification
	entry.MobileAppDaily_Notification_Sender = request.MobileAppDaily_Notification_Sender
	entry.MobileAppDaily_Notification_Text = request.MobileAppDaily_Notification_Text
	entry.MobileAppDaily_Login = request.MobileAppDaily_Login
	entry.Welcome_Notification = request.Welcome_Notification
	entry.Welcome_Notification_Sender = request.Welcome_Notification_Sender
	entry.Welcome_Notification_Text = request.Welcome_Notification_Text
	entry.Rejoiner_Notification = request.Rejoiner_Notification
	entry.Rejoiner_Notification_Sender = request.Rejoiner_Notification_Sender
	entry.Rejoiner_Notification_Text = request.Rejoiner_Notification_Text
	entry.Level_Change_Notification = request.Level_Change_Notification
	entry.Level_Change_Notification_Sender = request.Level_Change_Notification_Sender
	entry.Level_Change_Notification_Text = request.Level_Change_Notification_Text
	entry.MainGSM_Notification = request.MainGSM_Notification
	entry.MainGSM_Notification_Sender = request.MainGSM_Notification_Sender
	entry.MainGSM_Notification_Text = request.MainGSM_Notification_Text
	entry.MainGSMBalance_Amount = request.MainGSMBalance_Amount
	entry.MainGSMBalance_Points = request.MainGSMBalance_Points
	entry.GSM_SC_Notification = request.GSM_SC_Notification
	entry.GSM_SC_Notification_Sender = request.GSM_SC_Notification_Sender
	entry.GSM_SC_Notification_Text = request.GSM_SC_Notification_Text
	entry.GSM_SC_Airtime_Award_Type = request.GSM_SC_Airtime_Award_Type
	entry.GSM_SC_Airtime_Amount = request.GSM_SC_Airtime_Amount
	entry.GSM_SC_Airtime_Points = request.GSM_SC_Airtime_Points
	entry.GSM_EVC_Airtime_Notification = request.GSM_EVC_Airtime_Notification
	entry.GSM_EVC_Airtime_Notification_Sender = request.GSM_EVC_Airtime_Notification_Sender
	entry.GSM_EVC_Airtime_Notification_Text = request.GSM_EVC_Airtime_Notification_Text
	entry.GSM_EVC_Airtime_Award_Type = request.GSM_EVC_Airtime_Award_Type
	entry.GSM_EVC_Airtime_Amount = request.GSM_EVC_Airtime_Amount
	entry.GSM_EVC_Airtime_Points = request.GSM_EVC_Airtime_Points
	entry.GSM_EVC_Bundle_Notification = request.GSM_EVC_Bundle_Notification
	entry.GSM_EVC_Bundle_Notification_Sender = request.GSM_EVC_Bundle_Notification_Sender
	entry.GSM_EVC_Bundle_Notification_Text = request.GSM_EVC_Bundle_Notification_Text
	entry.GSM_EVC_Bundle_Award_Type = request.GSM_EVC_Bundle_Award_Type
	entry.GSM_EVC_Bundle_Amount = request.GSM_EVC_Bundle_Amount
	entry.GSM_EVC_Bundle_Points = request.GSM_EVC_Bundle_Points
	entry.MM_P2P_Notification = request.MM_P2P_Notification
	entry.MM_P2P_Notification_Sender = request.MM_P2P_Notification_Sender
	entry.MM_P2P_Notification_Text = request.MM_P2P_Notification_Text
	entry.MM_P2P_Award_Type = request.MM_P2P_Award_Type
	entry.MM_P2P_Amount = request.MM_P2P_Amount
	entry.MM_P2P_Points = request.MM_P2P_Points
	entry.MM_CASHIN_Notification = request.MM_CASHIN_Notification
	entry.MM_CASHIN_Notification_Sender = request.MM_CASHIN_Notification_Sender
	entry.MM_CASHIN_Notification_Text = request.MM_CASHIN_Notification_Text
	entry.MM_CASHIN_Award_Type = request.MM_CASHIN_Award_Type
	entry.MM_CASHIN_Amount = request.MM_CASHIN_Amount
	entry.MM_CASHIN_Points = request.MM_CASHIN_Points
	entry.MM_CASHOUT_Notification = request.MM_CASHOUT_Notification
	entry.MM_CASHOUT_Notification_Sender = request.MM_CASHOUT_Notification_Sender
	entry.MM_CASHOUT_Notification_Text = request.MM_CASHOUT_Notification_Text
	entry.MM_CASHOUT_Award_Type = request.MM_CASHOUT_Award_Type
	entry.MM_CASHOUT_Amount = request.MM_CASHOUT_Amount
	entry.MM_CASHOUT_Points = request.MM_CASHOUT_Points
	entry.MM_MERCHPAY_Notification = request.MM_MERCHPAY_Notification
	entry.MM_MERCHPAY_Notification_Sender = request.MM_MERCHPAY_Notification_Sender
	entry.MM_MERCHPAY_Notification_Text = request.MM_MERCHPAY_Notification_Text
	entry.MM_MERCHPAY_Award_Type = request.MM_MERCHPAY_Award_Type
	entry.MM_MERCHPAY_Amount = request.MM_MERCHPAY_Amount
	entry.MM_MERCHPAY_Points = request.MM_MERCHPAY_Points
	entry.MM_BILLPAY_Notification = request.MM_BILLPAY_Notification
	entry.MM_BILLPAY_Notification_Sender = request.MM_BILLPAY_Notification_Sender
	entry.MM_BILLPAY_Notification_Text = request.MM_BILLPAY_Notification_Text
	entry.MM_BILLPAY_Award_Type = request.MM_BILLPAY_Award_Type
	entry.MM_BILLPAY_Amount = request.MM_BILLPAY_Amount
	entry.MM_BILLPAY_Points = request.MM_BILLPAY_Points
	entry.Overwrite_MM_MERCHPAY_Notification = request.Overwrite_MM_MERCHPAY_Notification
	entry.Overwrite_MM_MERCHPAY_Notification_Sender = request.Overwrite_MM_MERCHPAY_Notification_Sender
	entry.Overwrite_MM_MERCHPAY_Notification_Text = request.Overwrite_MM_MERCHPAY_Notification_Text
	entry.Earning_Rules_Overwrite_MM_MERCHPAY_Keys = request.Earning_Rules_Overwrite_MM_MERCHPAY_Keys
	entry.Overwrite_MM_BILLPAY_Notification = request.Overwrite_MM_BILLPAY_Notification
	entry.Overwrite_MM_BILLPAY_Notification_Sender = request.Overwrite_MM_BILLPAY_Notification_Sender
	entry.Overwrite_MM_BILLPAY_Notification_Text = request.Overwrite_MM_BILLPAY_Notification_Text
	entry.Earning_Rules_Overwrite_MM_BILLPAY_Keys = request.Earning_Rules_Overwrite_MM_BILLPAY_Keys
	entry.MM_RC_Bundle_Notification = request.MM_RC_Bundle_Notification
	entry.MM_RC_Bundle_Notification_Sender = request.MM_RC_Bundle_Notification_Sender
	entry.MM_RC_Bundle_Notification_Text = request.MM_RC_Bundle_Notification_Text
	entry.MM_RC_Bundle_Award_Type = request.MM_RC_Bundle_Award_Type
	entry.MM_RC_Bundle_Amount = request.MM_RC_Bundle_Amount
	entry.MM_RC_Bundle_Points = request.MM_RC_Bundle_Points
	entry.MM_RC_Airtime_Notification = request.MM_RC_Airtime_Notification
	entry.MM_RC_Airtime_Notification_Sender = request.MM_RC_Airtime_Notification_Sender
	entry.MM_RC_Airtime_Notification_Text = request.MM_RC_Airtime_Notification_Text
	entry.MM_RC_Airtime_Award_Type = request.MM_RC_Airtime_Award_Type
	entry.MM_RC_Airtime_Amount = request.MM_RC_Airtime_Amount
	entry.MM_RC_Airtime_Points = request.MM_RC_Airtime_Points
	entry.MM_CTMMOREQ_Bundle_Notification = request.MM_CTMMOREQ_Bundle_Notification
	entry.MM_CTMMOREQ_Bundle_Notification_Sender = request.MM_CTMMOREQ_Bundle_Notification_Sender
	entry.MM_CTMMOREQ_Bundle_Notification_Text = request.MM_CTMMOREQ_Bundle_Notification_Text
	entry.MM_CTMMOREQ_Bundle_Award_Type = request.MM_CTMMOREQ_Bundle_Award_Type
	entry.MM_CTMMOREQ_Bundle_Amount = request.MM_CTMMOREQ_Bundle_Amount
	entry.MM_CTMMOREQ_Bundle_Points = request.MM_CTMMOREQ_Bundle_Points
	entry.MM_CTMMOREQ_Airtime_Notification = request.MM_CTMMOREQ_Airtime_Notification
	entry.MM_CTMMOREQ_Airtime_Notification_Sender = request.MM_CTMMOREQ_Airtime_Notification_Sender
	entry.MM_CTMMOREQ_Airtime_Notification_Text = request.MM_CTMMOREQ_Airtime_Notification_Text
	entry.MM_CTMMOREQ_Airtime_Award_Type = request.MM_CTMMOREQ_Airtime_Award_Type
	entry.MM_CTMMOREQ_Airtime_Amount = request.MM_CTMMOREQ_Airtime_Amount
	entry.MM_CTMMOREQ_Airtime_Points = request.MM_CTMMOREQ_Airtime_Points
	entry.MM_CBWREQ_Notification = request.MM_CBWREQ_Notification
	entry.MM_CBWREQ_Notification_Sender = request.MM_CBWREQ_Notification_Sender
	entry.MM_CBWREQ_Notification_Text = request.MM_CBWREQ_Notification_Text
	entry.MM_CBWREQ_Award_Type = request.MM_CBWREQ_Award_Type
	entry.MM_CBWREQ_Amount = request.MM_CBWREQ_Amount
	entry.MM_CBWREQ_Points = request.MM_CBWREQ_Points
	entry.MM_Airtime_Notification = request.MM_Airtime_Notification
	entry.MM_Airtime_Notification_Sender = request.MM_Airtime_Notification_Sender
	entry.MM_Airtime_Notification_Text = request.MM_Airtime_Notification_Text
	entry.MM_Airtime_Award_Type = request.MM_Airtime_Award_Type
	entry.MM_Airtime_Amount = request.MM_Airtime_Amount
	entry.MM_Airtime_Points = request.MM_Airtime_Points
	entry.MM_Bundle_Notification = request.MM_Bundle_Notification
	entry.MM_Bundle_Notification_Sender = request.MM_Bundle_Notification_Sender
	entry.MM_Bundle_Notification_Text = request.MM_Bundle_Notification_Text
	entry.MM_Bundle_Award_Type = request.MM_Bundle_Award_Type
	entry.MM_Bundle_Amount = request.MM_Bundle_Amount
	entry.MM_Bundle_Points = request.MM_Bundle_Points
	if request.NewKey != "" {
		if request.NewKey != request.Key {
			//delete old
			{
				delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
				if _, delErr := Mdb_Loyalty_Point_Earning_Rules.Coll.DeleteOne(delCtx, bson.M{"Key": request.Key}); delErr != nil {
					log.Println("Mdb_Loyalty_Point_Earning_Rules DeleteOne error:", delErr)
				}
				if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Loyalty_Point_Earning_Rules{Key: request.Key}.RedisKey()); delRedisErr != nil {
					log.Println("redisx.DelJSON Loyalty_Point_Earning_Rules error:", delRedisErr)
				}
				delCancel()
			}
			//update key
			entry.Key = request.NewKey
		}
	}
	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Loyalty_Point_Earning_Rules.Coll.UpdateOne(putCtx, bson.M{"Key": entry.Key}, bson.M{"$set": entry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			putCancel()
			err = putErr
			return Id, err
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, entry.RedisKey(), entry); putSetErr != nil {
			log.Println("redisx.SetJSON Loyalty_Point_Earning_Rules error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Point_Earning_Rules",
		Event_ActionType:   "Edit",
		Event_Description:  "",
		Event_Entry_Before: Current_Entry,
		Event_Entry_After:  entry,
	})
	return Id, nil
}

func (Uc *UserControl) Loyalty_Point_Earning_Rules_Get(Key string) (entries []Loyalty_Point_Earning_Rules, err error) {
	if Key == "" {
		entries, err = redisx.GetAllJSONByPattern[Loyalty_Point_Earning_Rules](context.Background(), RedisClient, redisx.ScanJSONOptions{
			Pattern: "Loyalty_Point_Earning_Rules:*", ScanCount: 500, PipelineSize: 250,
		})
		if err != nil {
			return entries, err
		}
		return entries, nil
	} else {
		entry, entryErr := redisx.GetJSON[Loyalty_Point_Earning_Rules](context.Background(), RedisClient, Loyalty_Point_Earning_Rules{Key: Key}.RedisKey())
		if redisx.IsNil(entryErr) {
			err = errors.New("key does not exist")
			return entries, err
		}
		if entryErr != nil {
			return entries, entryErr
		}
		entries = append(entries, entry)
		return entries, nil
	}
}

func (Uc *UserControl) Loyalty_Point_Earning_Rules_GetPaginated(Page, Limit int64) (entries []Loyalty_Point_Earning_Rules, err error) {
	if Page < 1 {
		return entries, errors.New("invalid page")
	}
	if Limit < 1 || Limit > 50000 {
		return entries, errors.New("invalid limit (accept value between 1 and 50000)")
	}

	pgCtx, pgCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer pgCancel()
	result, err := mongox.FindPage[Loyalty_Point_Earning_Rules](pgCtx, Mdb_Loyalty_Point_Earning_Rules.Coll, bson.M{}, mongox.PageRequest{
		Page: Page, PageSize: Limit, MaxPageSize: 50000,
	})
	if err != nil {
		return entries, err
	}
	entries = result.Items
	return entries, nil

}

func (Uc *UserControl) Loyalty_Point_Earning_Rules_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry, entryErr := redisx.GetJSON[Loyalty_Point_Earning_Rules](context.Background(), RedisClient, Loyalty_Point_Earning_Rules{Key: Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("entry does not exist")
		return err
	}
	if entryErr != nil {
		return entryErr
	}
	{
		delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, delErr := Mdb_Loyalty_Point_Earning_Rules.Coll.DeleteOne(delCtx, bson.M{"Key": Key}); delErr != nil {
			log.Println("Mdb_Loyalty_Point_Earning_Rules DeleteOne error:", delErr)
		}
		if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Loyalty_Point_Earning_Rules{Key: Key}.RedisKey()); delRedisErr != nil {
			log.Println("redisx.DelJSON Loyalty_Point_Earning_Rules error:", delRedisErr)
		}
		delCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Point_Earning_Rules",
		Event_ActionType:   "Delete",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  entry,
	})
	return nil
}

//Earning rules overwrite

func (Uc *UserControl) Loyalty_Point_Earning_Rules_Overwrite_Add(Login string, request Loyalty_Point_Earning_Rules_Overwrite_AddRequest) (Id int64, err error) {
	//check key if filled and if already used
	if request.AgentCode == "" {
		err = errors.New("Agent Code cannot be empty")
		return Id, err
	}
	if request.MM_Transaction_Type == "" {
		err = errors.New("transaction type can not be empty")
		return Id, err
	}
	if request.Earning_Rule_Key == "" {
		err = errors.New("earning rule key can not be empty")
		return Id, err
	}
	//check if key already used
	{
		chkCtx, chkCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, chkErr := redisx.GetJSON[Loyalty_Point_Earning_Rules_Overwrite](chkCtx, RedisClient, Loyalty_Point_Earning_Rules_Overwrite{Key: request.Earning_Rule_Key + "|" + request.MM_Transaction_Type + "|" + request.AgentCode}.RedisKey())
		chkCancel()
		if chkErr == nil {
			err = errors.New("key already exist")
			return Id, err
		}
	}
	//check values
	if request.Award_Type == "Amount" {
		if request.Amount > 0 && request.Amount < 1 {
			err = errors.New("invalid amount value")
			return Id, err
		}
	} else if request.Award_Type != "Transaction" {
		err = errors.New("invalid Award Type")
		return Id, err
	} else {
		request.Amount = 0
	}

	var entries []Loyalty_Governance
	entries, redisGovErr := redisx.GetAllJSONByPattern[Loyalty_Governance](context.Background(), RedisClient, redisx.ScanJSONOptions{
		Pattern: "Loyalty_Governance:*", ScanCount: 500, PipelineSize: 250,
	})
	if redisGovErr != nil {
		err = redisGovErr
		return Id, err
	}
	if len(entries) > 0 {
		if entries[0].MaxAllowedPoints_PerTransaction < request.Points {
			err = errors.New("points can not exceed the maximum allowed points per transaction")
			return Id, err
		}
	}
	//Prepare new entry
	var NewEntry Loyalty_Point_Earning_Rules_Overwrite
	{
		aiCtx, aiCancel := context.WithTimeout(context.Background(), 10*time.Second)
		NewEntry.Earning_Rules_Overwrite_Id, err = redisx.NextAutoIncrementID(aiCtx, RedisClient, Mdb_Loyalty_AutoIncrement.Coll, "Loyalty_Point_Earning_Rules_Overwrite-Id", redisx.NextIDOptions{
			RedisBase: "AutoIncrement", MongoRetries: 3, RetryBackoff: 500 * time.Millisecond,
		})
		aiCancel()
		if err != nil {
			return Id, err
		}
		Id = NewEntry.Earning_Rules_Overwrite_Id
	}
	NewEntry.Key = request.Earning_Rule_Key + "|" + request.MM_Transaction_Type + "|" + request.AgentCode
	NewEntry.Description = request.Description
	NewEntry.Earning_Rule_Key = request.Earning_Rule_Key
	NewEntry.Amount = request.Amount
	NewEntry.Award_Type = request.Award_Type
	NewEntry.Description = request.Description
	NewEntry.MM_Transaction_Type = request.MM_Transaction_Type
	NewEntry.AgentCode = request.AgentCode
	NewEntry.Points = request.Points

	//update earning rules
	entry, entryEarningErr := redisx.GetJSON[Loyalty_Point_Earning_Rules](context.Background(), RedisClient, Loyalty_Point_Earning_Rules{Key: request.Earning_Rule_Key}.RedisKey())
	if redisx.IsNil(entryEarningErr) {
		err = errors.New("earning rule does not exist")
		return Id, err
	}
	if entryEarningErr != nil {
		return Id, entryEarningErr
	}
	if request.MM_Transaction_Type == "MERCHPAY" {
		entry.Earning_Rules_Overwrite_MM_MERCHPAY_Keys = append(entry.Earning_Rules_Overwrite_MM_MERCHPAY_Keys, NewEntry.Key)
	} else if request.MM_Transaction_Type == "BILLPAY" {
		entry.Earning_Rules_Overwrite_MM_BILLPAY_Keys = append(entry.Earning_Rules_Overwrite_MM_BILLPAY_Keys, NewEntry.Key)
	}
	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Loyalty_Point_Earning_Rules.Coll.UpdateOne(putCtx, bson.M{"Key": entry.Key}, bson.M{"$set": entry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Loyalty_Point_Earning_Rules upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, entry.RedisKey(), entry); putSetErr != nil {
			log.Println("redisx.SetJSON Loyalty_Point_Earning_Rules error:", putSetErr)
		}
		putCancel()
	}
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Loyalty_Point_Earning_Rules_Overwrite.Coll.UpdateOne(putCtx, bson.M{"Key": NewEntry.Key}, bson.M{"$set": NewEntry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Loyalty_Point_Earning_Rules_Overwrite upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, NewEntry.RedisKey(), NewEntry); putSetErr != nil {
			log.Println("redisx.SetJSON Loyalty_Point_Earning_Rules_Overwrite error:", putSetErr)
		}
		putCancel()
	}

	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Point_Earning_Rules_Overwrite",
		Event_ActionType:   "Add",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  NewEntry,
	})
	return Id, nil
}

func (Uc *UserControl) Loyalty_Point_Earning_Rules_Overwrite_Edit(Login string, request Loyalty_Point_Earning_Rules_Overwrite_EditRequest) (Id int64, err error) {
	//check and validate
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	if request.AgentCode == "" {
		err = errors.New("Agent Code cannot be empty")
		return Id, err
	}
	if request.MM_Transaction_Type == "" {
		err = errors.New("transaction type can not be empty")
		return Id, err
	}
	if request.Earning_Rule_Key == "" {
		err = errors.New("earning rule key can not be empty")
		return Id, err
	}
	//check values
	if request.Award_Type == "Amount" {
		if request.Amount > 0 && request.Amount < 1 {
			err = errors.New("invalid value")
			return Id, err
		}
	} else if request.Award_Type != "Transaction" {
		err = errors.New("invalid Award Type")
		return Id, err
	} else {
		request.Amount = 0
	}

	entry, entryErr := redisx.GetJSON[Loyalty_Point_Earning_Rules_Overwrite](context.Background(), RedisClient, Loyalty_Point_Earning_Rules_Overwrite{Key: request.Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("key is not created")
		return Id, err
	}
	if entryErr != nil {
		return Id, entryErr
	}
	if entry.Earning_Rules_Overwrite_Id != request.Earning_Rules_Overwrite_Id {
		return Id, errors.New("id is not matching")
	}
	var entries []Loyalty_Governance
	entries, redisGovErr2 := redisx.GetAllJSONByPattern[Loyalty_Governance](context.Background(), RedisClient, redisx.ScanJSONOptions{
		Pattern: "Loyalty_Governance:*", ScanCount: 500, PipelineSize: 250,
	})
	if redisGovErr2 != nil {
		err = redisGovErr2
		return Id, err
	}
	if len(entries) > 0 {
		if entries[0].MaxAllowedPoints_PerTransaction < request.Points {
			err = errors.New("points can not exceed the maximum allowed points per transaction")
			return Id, err
		}
	}
	Current_Entry := entry
	//Prepare new entry

	//add to cache and DB
	if request.AgentCode != entry.AgentCode {
		//delete old
		{
			delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
			if _, delErr := Mdb_Loyalty_Point_Earning_Rules_Overwrite.Coll.DeleteOne(delCtx, bson.M{"Key": request.Key}); delErr != nil {
				log.Println("Mdb_Loyalty_Point_Earning_Rules_Overwrite DeleteOne error:", delErr)
			}
			if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Loyalty_Point_Earning_Rules_Overwrite{Key: request.Key}.RedisKey()); delRedisErr != nil {
				log.Println("redisx.DelJSON Loyalty_Point_Earning_Rules_Overwrite error:", delRedisErr)
			}
			delCancel()
		}
		//update key
		entry.Key = request.Earning_Rule_Key + "|" + request.MM_Transaction_Type + "|" + request.AgentCode

		//update earning rules
		earning, earningErr := redisx.GetJSON[Loyalty_Point_Earning_Rules](context.Background(), RedisClient, Loyalty_Point_Earning_Rules{Key: request.Earning_Rule_Key}.RedisKey())
		if redisx.IsNil(earningErr) {
			err = errors.New("earning rule does not exist")
			return Id, err
		}
		if earningErr != nil {
			return Id, earningErr
		}
		var newKeys []string

		for _, i := range earning.Earning_Rules_Overwrite_MM_MERCHPAY_Keys {
			if i != request.Key {
				newKeys = append(newKeys, i)
			}
		}
		newKeys = append(newKeys, entry.Key)
		if request.MM_Transaction_Type == "MERCHPAY" {
			earning.Earning_Rules_Overwrite_MM_MERCHPAY_Keys = newKeys
		} else if request.MM_Transaction_Type == "BILLPAY" {
			earning.Earning_Rules_Overwrite_MM_BILLPAY_Keys = newKeys
		}
		{
			putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
			if _, putErr := Mdb_Loyalty_Point_Earning_Rules.Coll.UpdateOne(putCtx, bson.M{"Key": entry.Earning_Rule_Key}, bson.M{"$set": earning}, options.UpdateOne().SetUpsert(true)); putErr != nil {
				log.Println("Mdb_Loyalty_Point_Earning_Rules upsert error:", putErr)
			}
			if putSetErr := redisx.SetJSON(putCtx, RedisClient, earning.RedisKey(), earning); putSetErr != nil {
				log.Println("redisx.SetJSON Loyalty_Point_Earning_Rules error:", putSetErr)
			}
			putCancel()
		}
	}
	entry.Description = request.Description
	entry.Amount = request.Amount
	entry.Award_Type = request.Award_Type
	entry.Earning_Rule_Key = request.Earning_Rule_Key
	entry.Earning_Rules_Overwrite_Id = request.Earning_Rules_Overwrite_Id
	entry.MM_Transaction_Type = request.MM_Transaction_Type
	entry.AgentCode = request.AgentCode
	entry.Points = request.Points

	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Loyalty_Point_Earning_Rules_Overwrite.Coll.UpdateOne(putCtx, bson.M{"Key": entry.Key}, bson.M{"$set": entry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Loyalty_Point_Earning_Rules_Overwrite upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, entry.RedisKey(), entry); putSetErr != nil {
			log.Println("redisx.SetJSON Loyalty_Point_Earning_Rules_Overwrite error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Point_Earning_Rules_Overwrite",
		Event_ActionType:   "Edit",
		Event_Description:  "",
		Event_Entry_Before: Current_Entry,
		Event_Entry_After:  entry,
	})
	return Id, nil
}

func (Uc *UserControl) Loyalty_Point_Earning_Rules_Overwrite_Get(Key string) (entries []Loyalty_Point_Earning_Rules_Overwrite, err error) {
	if Key == "" {
		entries, err = redisx.GetAllJSONByPattern[Loyalty_Point_Earning_Rules_Overwrite](context.Background(), RedisClient, redisx.ScanJSONOptions{
			Pattern: "Loyalty_Point_Earning_Rules_Overwrite:*", ScanCount: 500, PipelineSize: 250,
		})
		if err != nil {
			return entries, err
		}
		return entries, nil
	} else {
		entry, entryErr := redisx.GetJSON[Loyalty_Point_Earning_Rules_Overwrite](context.Background(), RedisClient, Loyalty_Point_Earning_Rules_Overwrite{Key: Key}.RedisKey())
		if redisx.IsNil(entryErr) {
			err = errors.New("key does not exist")
			return entries, err
		}
		if entryErr != nil {
			return entries, entryErr
		}
		entries = append(entries, entry)
		return entries, nil
	}
}

func (Uc *UserControl) Loyalty_Point_Earning_Rules_Overwrite_GetPaginated(Page, Limit int64) (entries []Loyalty_Point_Earning_Rules_Overwrite, err error) {
	if Page < 1 {
		return entries, errors.New("invalid page")
	}
	if Limit < 1 || Limit > 50000 {
		return entries, errors.New("invalid limit (accept value between 1 and 50000)")
	}

	pgCtx, pgCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer pgCancel()
	result, err := mongox.FindPage[Loyalty_Point_Earning_Rules_Overwrite](pgCtx, Mdb_Loyalty_Point_Earning_Rules_Overwrite.Coll, bson.M{}, mongox.PageRequest{
		Page: Page, PageSize: Limit, MaxPageSize: 50000,
	})
	if err != nil {
		return entries, err
	}
	entries = result.Items
	return entries, nil

}

func (Uc *UserControl) Loyalty_Point_Earning_Rules_Overwrite_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	parts := strings.Split(Key, "|")
	if len(parts) != 3 {
		err = errors.New("valid key format")
		return err
	}
	earningRuleKey := parts[0]
	mmTransactionType := parts[1]

	entry, entryErr := redisx.GetJSON[Loyalty_Point_Earning_Rules_Overwrite](context.Background(), RedisClient, Loyalty_Point_Earning_Rules_Overwrite{Key: Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("entry does not exist")
		return err
	}
	if entryErr != nil {
		return entryErr
	}

	earning, earningErr := redisx.GetJSON[Loyalty_Point_Earning_Rules](context.Background(), RedisClient, Loyalty_Point_Earning_Rules{Key: earningRuleKey}.RedisKey())
	if redisx.IsNil(earningErr) {
		err = errors.New("earning rule does not exist")
		return err
	}
	if earningErr != nil {
		return earningErr
	}
	var newKeys []string

	for _, i := range earning.Earning_Rules_Overwrite_MM_MERCHPAY_Keys {
		if i != Key {
			newKeys = append(newKeys, i)
		}
	}
	if mmTransactionType == "MERCHPAY" {
		earning.Earning_Rules_Overwrite_MM_MERCHPAY_Keys = newKeys
	} else if mmTransactionType == "BILLPAY" {
		earning.Earning_Rules_Overwrite_MM_BILLPAY_Keys = newKeys
	}
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Loyalty_Point_Earning_Rules.Coll.UpdateOne(putCtx, bson.M{"Key": earningRuleKey}, bson.M{"$set": earning}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Loyalty_Point_Earning_Rules upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, earning.RedisKey(), earning); putSetErr != nil {
			log.Println("redisx.SetJSON Loyalty_Point_Earning_Rules error:", putSetErr)
		}
		putCancel()
	}
	{
		delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, delErr := Mdb_Loyalty_Point_Earning_Rules_Overwrite.Coll.DeleteOne(delCtx, bson.M{"Key": Key}); delErr != nil {
			log.Println("Mdb_Loyalty_Point_Earning_Rules_Overwrite DeleteOne error:", delErr)
		}
		if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Loyalty_Point_Earning_Rules_Overwrite{Key: Key}.RedisKey()); delRedisErr != nil {
			log.Println("redisx.DelJSON Loyalty_Point_Earning_Rules_Overwrite error:", delRedisErr)
		}
		delCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Point_Earning_Rules_Overwrite",
		Event_ActionType:   "Delete",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  entry,
	})
	return nil
}

// ***********************************************************************
// Loyalty_Point_Expiry_Rules functions
// ***********************************************************************
func (Uc *UserControl) Loyalty_Point_Expiry_Rules_Add(Login string, request Loyalty_Point_Expiry_Rules_AddRequest) (Id int64, err error) {
	//check key if filled and if already used
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	//check if key already used
	{
		chkCtx, chkCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, chkErr := redisx.GetJSON[Loyalty_Point_Expiry_Rules](chkCtx, RedisClient, Loyalty_Point_Expiry_Rules{Key: request.Key}.RedisKey())
		chkCancel()
		if chkErr == nil {
			err = errors.New("key already exist")
			return Id, err
		}
	}

	//Prepare new entry
	var NewEntry Loyalty_Point_Expiry_Rules
	{
		aiCtx, aiCancel := context.WithTimeout(context.Background(), 10*time.Second)
		NewEntry.Expiry_Rules_Id, err = redisx.NextAutoIncrementID(aiCtx, RedisClient, Mdb_Loyalty_AutoIncrement.Coll, "Loyalty_Point_Expiry_Rules-Id", redisx.NextIDOptions{
			RedisBase: "AutoIncrement", MongoRetries: 3, RetryBackoff: 500 * time.Millisecond,
		})
		aiCancel()
		if err != nil {
			return Id, err
		}
		Id = NewEntry.Expiry_Rules_Id
	}
	NewEntry.Key = request.Key
	NewEntry.Description = request.Description
	NewEntry.Rolling_Expiration = request.Rolling_Expiration
	NewEntry.Validity_Unit = request.Validity_Unit
	NewEntry.Validity_Duration = request.Validity_Duration
	NewEntry.Grace_Validity_Unit = request.Grace_Validity_Unit
	NewEntry.Grace_Validity_Duration = request.Grace_Validity_Duration
	NewEntry.Fix_Date_Expiration = request.Fix_Date_Expiration
	NewEntry.Expiration_Trigger_date = request.Expiration_Trigger_date
	NewEntry.Expiration_Point_Before = request.Expiration_Point_Before
	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Loyalty_Point_Expiry_Rules.Coll.UpdateOne(putCtx, bson.M{"Key": NewEntry.Key}, bson.M{"$set": NewEntry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Loyalty_Point_Expiry_Rules upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, NewEntry.RedisKey(), NewEntry); putSetErr != nil {
			log.Println("redisx.SetJSON Loyalty_Point_Expiry_Rules error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Point_Expiry_Rules",
		Event_ActionType:   "Add",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  NewEntry,
	})
	return Id, nil
}

func (Uc *UserControl) Loyalty_Point_Expiry_Rules_Edit(Login string, request Loyalty_Point_Expiry_Rules_EditRequest) (Id int64, err error) {
	//check and validate
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	entry, entryErr := redisx.GetJSON[Loyalty_Point_Expiry_Rules](context.Background(), RedisClient, Loyalty_Point_Expiry_Rules{Key: request.Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("key is not created")
		return Id, err
	}
	if entryErr != nil {
		return Id, entryErr
	}
	if entry.Expiry_Rules_Id != request.Expiry_Rules_Id {
		return Id, errors.New("id is not matching")
	}
	Current_Entry := entry
	//Prepare new entry
	entry.Key = request.Key
	entry.Description = request.Description
	entry.Rolling_Expiration = request.Rolling_Expiration
	entry.Validity_Unit = request.Validity_Unit
	entry.Validity_Duration = request.Validity_Duration
	entry.Grace_Validity_Unit = request.Grace_Validity_Unit
	entry.Grace_Validity_Duration = request.Grace_Validity_Duration
	entry.Fix_Date_Expiration = request.Fix_Date_Expiration
	entry.Expiration_Trigger_date = request.Expiration_Trigger_date
	entry.Expiration_Point_Before = request.Expiration_Point_Before
	if request.NewKey != "" {
		if request.NewKey != request.Key {
			//delete old
			{
				delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
				if _, delErr := Mdb_Loyalty_Point_Expiry_Rules.Coll.DeleteOne(delCtx, bson.M{"Key": request.Key}); delErr != nil {
					log.Println("Mdb_Loyalty_Point_Expiry_Rules DeleteOne error:", delErr)
				}
				if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Loyalty_Point_Expiry_Rules{Key: request.Key}.RedisKey()); delRedisErr != nil {
					log.Println("redisx.DelJSON Loyalty_Point_Expiry_Rules error:", delRedisErr)
				}
				delCancel()
			}
			//update key
			entry.Key = request.NewKey
		}
	}
	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Loyalty_Point_Expiry_Rules.Coll.UpdateOne(putCtx, bson.M{"Key": entry.Key}, bson.M{"$set": entry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Loyalty_Point_Expiry_Rules upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, entry.RedisKey(), entry); putSetErr != nil {
			log.Println("redisx.SetJSON Loyalty_Point_Expiry_Rules error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Point_Expiry_Rules",
		Event_ActionType:   "Edit",
		Event_Description:  "",
		Event_Entry_Before: Current_Entry,
		Event_Entry_After:  entry,
	})
	return Id, nil
}

func (Uc *UserControl) Loyalty_Point_Expiry_Rules_Get(Key string) (entries []Loyalty_Point_Expiry_Rules, err error) {
	if Key == "" {
		entries, err = redisx.GetAllJSONByPattern[Loyalty_Point_Expiry_Rules](context.Background(), RedisClient, redisx.ScanJSONOptions{
			Pattern: "Loyalty_Point_Expiry_Rules:*", ScanCount: 500, PipelineSize: 250,
		})
		if err != nil {
			return entries, err
		}
		return entries, nil
	} else {
		entry, entryErr := redisx.GetJSON[Loyalty_Point_Expiry_Rules](context.Background(), RedisClient, Loyalty_Point_Expiry_Rules{Key: Key}.RedisKey())
		if redisx.IsNil(entryErr) {
			err = errors.New("key does not exist")
			return entries, err
		}
		if entryErr != nil {
			return entries, entryErr
		}
		entries = append(entries, entry)
		return entries, nil
	}
}

func (Uc *UserControl) Loyalty_Point_Expiry_Rules_GetPaginated(Page, Limit int64) (entries []Loyalty_Point_Expiry_Rules, err error) {
	if Page < 1 {
		return entries, errors.New("invalid page")
	}
	if Limit < 1 || Limit > 50000 {
		return entries, errors.New("invalid limit (accept value between 1 and 50000)")
	}

	//var array []daoc.DAOFindCriteria
	// if Outlet_Key != "" {
	// 	//restrict access for records that belong to this user
	// 	var criteria daoc.DAOFindCriteria = daoc.DAOFindCriteria{
	// 		Field:    "Outlet_Key",
	// 		Value:    Outlet_Key,
	// 		Operator: "EQUAL",
	// 	}
	// 	array = append(array, criteria)
	// }
	// if Agent_Key != "" {
	// 	//restrict access for records that belong to this user
	// 	var criteria daoc.DAOFindCriteria = daoc.DAOFindCriteria{
	// 		Field:    "Agent_Key",
	// 		Value:    Agent_Key,
	// 		Operator: "EQUAL",
	// 	}
	// 	array = append(array, criteria)
	// }
	// if len(array) > 0 {
	// 	findparams.FindCriteria = array
	// }
	pgCtx, pgCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer pgCancel()
	result, err := mongox.FindPage[Loyalty_Point_Expiry_Rules](pgCtx, Mdb_Loyalty_Point_Expiry_Rules.Coll, bson.M{}, mongox.PageRequest{
		Page: Page, PageSize: Limit, MaxPageSize: 50000,
	})
	if err != nil {
		return entries, err
	}
	entries = result.Items
	return entries, nil

}

func (Uc *UserControl) Loyalty_Point_Expiry_Rules_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry, entryErr := redisx.GetJSON[Loyalty_Point_Expiry_Rules](context.Background(), RedisClient, Loyalty_Point_Expiry_Rules{Key: Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("entry does not exist")
		return err
	}
	if entryErr != nil {
		return entryErr
	}
	{
		delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, delErr := Mdb_Loyalty_Point_Expiry_Rules.Coll.DeleteOne(delCtx, bson.M{"Key": Key}); delErr != nil {
			log.Println("Mdb_Loyalty_Point_Expiry_Rules DeleteOne error:", delErr)
		}
		if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Loyalty_Point_Expiry_Rules{Key: Key}.RedisKey()); delRedisErr != nil {
			log.Println("redisx.DelJSON Loyalty_Point_Expiry_Rules error:", delRedisErr)
		}
		delCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Point_Expiry_Rules",
		Event_ActionType:   "Delete",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  entry,
	})
	return nil
}

// ***********************************************************************
// Loyalty_Point_Redemption_Rules functions
// ***********************************************************************
func (Uc *UserControl) Loyalty_Point_Redemption_Rules_Add(Login string, request Loyalty_Point_Redemption_Rules_AddRequest) (Id int64, err error) {
	//check key if filled and if already used
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	//check if key already used
	{
		chkCtx, chkCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, chkErr := redisx.GetJSON[Loyalty_Point_Redemption_Rules](chkCtx, RedisClient, Loyalty_Point_Redemption_Rules{Key: request.Key}.RedisKey())
		chkCancel()
		if chkErr == nil {
			err = errors.New("key already exist")
			return Id, err
		}
	}

	//Prepare new entry
	var NewEntry Loyalty_Point_Redemption_Rules
	{
		aiCtx, aiCancel := context.WithTimeout(context.Background(), 10*time.Second)
		NewEntry.Redemption_Rules_Id, err = redisx.NextAutoIncrementID(aiCtx, RedisClient, Mdb_Loyalty_AutoIncrement.Coll, "Loyalty_Point_Redemption_Rules-Id", redisx.NextIDOptions{
			RedisBase: "AutoIncrement", MongoRetries: 3, RetryBackoff: 500 * time.Millisecond,
		})
		aiCancel()
		if err != nil {
			return Id, err
		}
		Id = NewEntry.Redemption_Rules_Id
	}
	NewEntry.Key = request.Key
	NewEntry.Description = request.Description
	NewEntry.Min_Accumulated_Points = request.Min_Accumulated_Points
	NewEntry.Allow_Negative_Balance_ToRedeem = request.Allow_Negative_Balance_ToRedeem
	NewEntry.Allow_PendingLendme_ToRedeem = request.Allow_PendingLendme_ToRedeem
	NewEntry.Airtime_Notification = request.Airtime_Notification
	NewEntry.Airtime_Notification_Sender = request.Airtime_Notification_Sender
	NewEntry.Airtime_Notification_Text = request.Airtime_Notification_Text
	NewEntry.Airtime_MinPoints = request.Airtime_MinPoints
	NewEntry.Available_MinPoints_for_Airtime = request.Available_MinPoints_for_Airtime
	NewEntry.Airtime_Amount = request.Airtime_Amount
	NewEntry.Airtime_Points = request.Airtime_Points
	NewEntry.Airtime_EVC_Account = request.Airtime_EVC_Account
	if request.Airtime_EVC_PIN == "" {
		NewEntry.Airtime_EVC_PIN = ""
	} else {
		enc_string, err := EcryptToHexString(request.Airtime_EVC_PIN)
		if err != nil {
			return Id, errors.New("error encrypting password: " + err.Error())
		}
		NewEntry.Airtime_EVC_PIN = enc_string

	}
	NewEntry.MobileMoney_MinPoints = request.MobileMoney_MinPoints
	NewEntry.Available_MinPoints_for_MobileMoney = request.Available_MinPoints_for_MobileMoney
	NewEntry.MobileMoney_Amount = request.MobileMoney_Amount
	NewEntry.MobileMoney_Points = request.MobileMoney_Points
	NewEntry.MobileMoney_MerchantAccount = request.MobileMoney_MerchantAccount
	if request.MobileMoney_MerchantPIN == "" {
		NewEntry.MobileMoney_MerchantPIN = ""
	} else {
		enc_string, err := EcryptToHexString(request.MobileMoney_MerchantPIN)
		if err != nil {
			return Id, errors.New("error encrypting password: " + err.Error())
		}
		NewEntry.MobileMoney_MerchantPIN = enc_string

	}
	NewEntry.MobileMoney_Notification = request.MobileMoney_Notification
	NewEntry.MobileMoney_Notification_Sender = request.MobileMoney_Notification_Sender
	NewEntry.MobileMoney_Notification_Text = request.MobileMoney_Notification_Text
	NewEntry.Bundles_MinPoints = request.Bundles_MinPoints
	NewEntry.Bundles_EVC_Account = request.Bundles_EVC_Account
	if request.Bundles_EVC_Account == "" {
		NewEntry.Bundles_EVC_Account = ""
	} else {
		enc_string, err := EcryptToHexString(request.Bundles_EVC_Account)
		if err != nil {
			return Id, errors.New("error encrypting password: " + err.Error())
		}
		NewEntry.Bundles_EVC_Account = enc_string

	}
	NewEntry.Bundles_Product_Catalogue_Channel = request.Bundles_Product_Catalogue_Channel
	NewEntry.Bundles_Product_Catalogue_Plan = request.Bundles_Product_Catalogue_Plan
	NewEntry.Bundles_Product_Catalogue_Version = request.Bundles_Product_Catalogue_Version
	NewEntry.Bundles_Notification = request.Bundles_Notification
	NewEntry.Bundles_Notification_Sender = request.Bundles_Notification_Sender
	NewEntry.Bundles_Notification_Text = request.Bundles_Notification_Text
	NewEntry.FreeSpinAndWin_MinPoints = request.FreeSpinAndWin_MinPoints
	NewEntry.Available_MinPoints_for_SpinAndWin = request.Available_MinPoints_for_SpinAndWin
	NewEntry.FreeSpinAndWin_PointsPerSpin = request.FreeSpinAndWin_PointsPerSpin
	NewEntry.FreeSpinAndWin_Notification = request.FreeSpinAndWin_Notification
	NewEntry.FreeSpinAndWin_Notification_Sender = request.FreeSpinAndWin_Notification_Sender
	NewEntry.FreeSpinAndWin_Notification_Text = request.FreeSpinAndWin_Notification_Text
	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Loyalty_Point_Redemption_Rules.Coll.UpdateOne(putCtx, bson.M{"Key": NewEntry.Key}, bson.M{"$set": NewEntry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Loyalty_Point_Redemption_Rules upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, NewEntry.RedisKey(), NewEntry); putSetErr != nil {
			log.Println("redisx.SetJSON Loyalty_Point_Redemption_Rules error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Point_Redemption_Rules",
		Event_ActionType:   "Add",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  NewEntry,
	})
	return Id, nil
}

func (Uc *UserControl) Loyalty_Point_Redemption_Rules_Edit(Login string, request Loyalty_Point_Redemption_Rules_EditRequest) (Id int64, err error) {
	//check and validate
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	entry, entryErr := redisx.GetJSON[Loyalty_Point_Redemption_Rules](context.Background(), RedisClient, Loyalty_Point_Redemption_Rules{Key: request.Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("key is not created")
		return Id, err
	}
	if entryErr != nil {
		return Id, entryErr
	}
	if entry.Redemption_Rules_Id != request.Redemption_Rules_Id {
		return Id, errors.New("id is not matching")
	}
	Current_Entry := entry
	//Prepare new entry
	entry.Key = request.Key
	entry.Description = request.Description
	entry.Min_Accumulated_Points = request.Min_Accumulated_Points
	entry.Allow_Negative_Balance_ToRedeem = request.Allow_Negative_Balance_ToRedeem
	entry.Allow_PendingLendme_ToRedeem = request.Allow_PendingLendme_ToRedeem
	entry.Airtime_Notification = request.Airtime_Notification
	entry.Airtime_Notification_Sender = request.Airtime_Notification_Sender
	entry.Airtime_Notification_Text = request.Airtime_Notification_Text
	entry.Airtime_MinPoints = request.Airtime_MinPoints
	entry.Available_MinPoints_for_Airtime = request.Available_MinPoints_for_Airtime
	entry.Airtime_Amount = request.Airtime_Amount
	entry.Airtime_Points = request.Airtime_Points
	entry.Airtime_EVC_Account = request.Airtime_EVC_Account
	if request.Airtime_EVC_PIN == "" {
		entry.Airtime_EVC_PIN = ""
	} else if request.Airtime_EVC_PIN != "****" {
		enc_string, err := EcryptToHexString(request.Airtime_EVC_PIN)
		if err != nil {
			return Id, errors.New("error encrypting password: " + err.Error())
		}
		entry.Airtime_EVC_PIN = enc_string
	}
	entry.MobileMoney_MinPoints = request.MobileMoney_MinPoints
	entry.Available_MinPoints_for_MobileMoney = request.Available_MinPoints_for_MobileMoney
	entry.MobileMoney_Amount = request.MobileMoney_Amount
	entry.MobileMoney_Points = request.MobileMoney_Points
	entry.MobileMoney_MerchantAccount = request.MobileMoney_MerchantAccount

	if request.MobileMoney_MerchantPIN == "" {
		entry.MobileMoney_MerchantPIN = ""
	} else if request.MobileMoney_MerchantPIN != "****" {
		enc_string, err := EcryptToHexString(request.MobileMoney_MerchantPIN)
		if err != nil {
			return Id, errors.New("error encrypting password: " + err.Error())
		}
		entry.MobileMoney_MerchantPIN = enc_string
	}
	entry.MobileMoney_Notification = request.MobileMoney_Notification
	entry.MobileMoney_Notification_Sender = request.MobileMoney_Notification_Sender
	entry.MobileMoney_Notification_Text = request.MobileMoney_Notification_Text
	entry.Bundles_MinPoints = request.Bundles_MinPoints
	entry.Bundles_EVC_Account = request.Bundles_EVC_Account
	if request.Bundles_EVC_PIN == "" {
		entry.Bundles_EVC_PIN = ""
	} else if request.Bundles_EVC_PIN != "****" {
		enc_string, err := EcryptToHexString(request.Bundles_EVC_PIN)
		if err != nil {
			return Id, errors.New("error encrypting password: " + err.Error())
		}
		entry.Bundles_EVC_PIN = enc_string
	}
	entry.Bundles_Product_Catalogue_Channel = request.Bundles_Product_Catalogue_Channel
	entry.Bundles_Product_Catalogue_Plan = request.Bundles_Product_Catalogue_Plan
	entry.Bundles_Product_Catalogue_Version = request.Bundles_Product_Catalogue_Version
	entry.Bundles_Notification = request.Bundles_Notification
	entry.Bundles_Notification_Sender = request.Bundles_Notification_Sender
	entry.Bundles_Notification_Text = request.Bundles_Notification_Text
	entry.FreeSpinAndWin_MinPoints = request.FreeSpinAndWin_MinPoints
	entry.Available_MinPoints_for_SpinAndWin = request.Available_MinPoints_for_SpinAndWin
	entry.FreeSpinAndWin_PointsPerSpin = request.FreeSpinAndWin_PointsPerSpin
	entry.FreeSpinAndWin_Notification = request.FreeSpinAndWin_Notification
	entry.FreeSpinAndWin_Notification_Sender = request.FreeSpinAndWin_Notification_Sender
	entry.FreeSpinAndWin_Notification_Text = request.FreeSpinAndWin_Notification_Text
	if request.NewKey != "" {
		if request.NewKey != request.Key {
			//delete old
			{
				delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
				if _, delErr := Mdb_Loyalty_Point_Redemption_Rules.Coll.DeleteOne(delCtx, bson.M{"Key": request.Key}); delErr != nil {
					log.Println("Mdb_Loyalty_Point_Redemption_Rules DeleteOne error:", delErr)
				}
				if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Loyalty_Point_Redemption_Rules{Key: request.Key}.RedisKey()); delRedisErr != nil {
					log.Println("redisx.DelJSON Loyalty_Point_Redemption_Rules error:", delRedisErr)
				}
				delCancel()
			}
			//update key
			entry.Key = request.NewKey
		}
	}
	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Loyalty_Point_Redemption_Rules.Coll.UpdateOne(putCtx, bson.M{"Key": entry.Key}, bson.M{"$set": entry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Loyalty_Point_Redemption_Rules upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, entry.RedisKey(), entry); putSetErr != nil {
			log.Println("redisx.SetJSON Loyalty_Point_Redemption_Rules error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Point_Redemption_Rules",
		Event_ActionType:   "Edit",
		Event_Description:  "",
		Event_Entry_Before: Current_Entry,
		Event_Entry_After:  entry,
	})
	return Id, nil
}

func (Uc *UserControl) Loyalty_Point_Redemption_Rules_Get(Key string) (entries []Loyalty_Point_Redemption_Rules, err error) {
	if Key == "" {
		entries, err = redisx.GetAllJSONByPattern[Loyalty_Point_Redemption_Rules](context.Background(), RedisClient, redisx.ScanJSONOptions{
			Pattern: "Loyalty_Point_Redemption_Rules:*", ScanCount: 500, PipelineSize: 250,
		})
		if err != nil {
			return entries, err
		}
		return entries, nil
	} else {
		entry, entryErr := redisx.GetJSON[Loyalty_Point_Redemption_Rules](context.Background(), RedisClient, Loyalty_Point_Redemption_Rules{Key: Key}.RedisKey())
		if redisx.IsNil(entryErr) {
			err = errors.New("key does not exist")
			return entries, err
		}
		if entryErr != nil {
			return entries, entryErr
		}
		entries = append(entries, entry)
		return entries, nil
	}
}

func (Uc *UserControl) Loyalty_Point_Redemption_Rules_GetPaginated(Page, Limit int64) (entries []Loyalty_Point_Redemption_Rules, err error) {
	if Page < 1 {
		return entries, errors.New("invalid page")
	}
	if Limit < 1 || Limit > 50000 {
		return entries, errors.New("invalid limit (accept value between 1 and 50000)")
	}

	//var array []daoc.DAOFindCriteria
	// if Outlet_Key != "" {
	// 	//restrict access for records that belong to this user
	// 	var criteria daoc.DAOFindCriteria = daoc.DAOFindCriteria{
	// 		Field:    "Outlet_Key",
	// 		Value:    Outlet_Key,
	// 		Operator: "EQUAL",
	// 	}
	// 	array = append(array, criteria)
	// }
	// if Agent_Key != "" {
	// 	//restrict access for records that belong to this user
	// 	var criteria daoc.DAOFindCriteria = daoc.DAOFindCriteria{
	// 		Field:    "Agent_Key",
	// 		Value:    Agent_Key,
	// 		Operator: "EQUAL",
	// 	}
	// 	array = append(array, criteria)
	// }
	// if len(array) > 0 {
	// 	findparams.FindCriteria = array
	// }
	pgCtx, pgCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer pgCancel()
	result, err := mongox.FindPage[Loyalty_Point_Redemption_Rules](pgCtx, Mdb_Loyalty_Point_Redemption_Rules.Coll, bson.M{}, mongox.PageRequest{
		Page: Page, PageSize: Limit, MaxPageSize: 50000,
	})
	if err != nil {
		return entries, err
	}
	entries = result.Items
	return entries, nil

}

func (Uc *UserControl) Loyalty_Point_Redemption_Rules_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry, entryErr := redisx.GetJSON[Loyalty_Point_Redemption_Rules](context.Background(), RedisClient, Loyalty_Point_Redemption_Rules{Key: Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("entry does not exist")
		return err
	}
	if entryErr != nil {
		return entryErr
	}
	{
		delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, delErr := Mdb_Loyalty_Point_Redemption_Rules.Coll.DeleteOne(delCtx, bson.M{"Key": Key}); delErr != nil {
			log.Println("Mdb_Loyalty_Point_Redemption_Rules DeleteOne error:", delErr)
		}
		if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Loyalty_Point_Redemption_Rules{Key: Key}.RedisKey()); delRedisErr != nil {
			log.Println("redisx.DelJSON Loyalty_Point_Redemption_Rules error:", delRedisErr)
		}
		delCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Point_Redemption_Rules",
		Event_ActionType:   "Delete",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  entry,
	})
	return nil
}

// ***********************************************************************
// Loyalty_Plan functions
// ***********************************************************************
func (Uc *UserControl) Loyalty_Plan_Add(Login string, request Loyalty_Plan_AddRequest) (Id int64, err error) {
	//check key if filled and if already used
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	//check if key already used
	{
		chkCtx, chkCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, chkErr := redisx.GetJSON[Loyalty_Plan](chkCtx, RedisClient, Loyalty_Plan{Key: request.Key}.RedisKey())
		chkCancel()
		if chkErr == nil {
			err = errors.New("key already exist")
			return Id, err
		}
	}

	planEntries, planEntriesErr := redisx.GetAllJSONByPattern[Loyalty_Plan](context.Background(), RedisClient, redisx.ScanJSONOptions{
		Pattern: "Loyalty_Plan:*", ScanCount: 500, PipelineSize: 250,
	})
	if planEntriesErr != nil {
		err = planEntriesErr
		return Id, err
	}
	for _, planEntry := range planEntries {
		if request.Loyalty_Level_Key == planEntry.Loyalty_Level_Key && request.Loyalty_Account_Segment_Key == planEntry.Loyalty_Account_Segment_Key {
			err = errors.New("a plan with the same Loyalty Level and Account Segment already exists. Please choose a different combination")
			return Id, err
		}
	}
	//Prepare new entry
	var NewEntry Loyalty_Plan
	{
		aiCtx, aiCancel := context.WithTimeout(context.Background(), 10*time.Second)
		NewEntry.Plan_Id, err = redisx.NextAutoIncrementID(aiCtx, RedisClient, Mdb_Loyalty_AutoIncrement.Coll, "Loyalty_Plan-Id", redisx.NextIDOptions{
			RedisBase: "AutoIncrement", MongoRetries: 3, RetryBackoff: 500 * time.Millisecond,
		})
		aiCancel()
		if err != nil {
			return Id, err
		}
		Id = NewEntry.Plan_Id
	}
	NewEntry.Key = request.Key
	NewEntry.Description = request.Description
	NewEntry.Loyalty_Level_Key = request.Loyalty_Level_Key
	NewEntry.Loyalty_Account_Segment_Key = request.Loyalty_Account_Segment_Key
	NewEntry.Earning_Rules_Key = request.Earning_Rules_Key
	NewEntry.Expiry_Rules_Key = request.Expiry_Rules_Key
	NewEntry.Redemption_Rules_Key = request.Redemption_Rules_Key
	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Loyalty_Plan.Coll.UpdateOne(putCtx, bson.M{"Key": NewEntry.Key}, bson.M{"$set": NewEntry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Loyalty_Plan upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, NewEntry.RedisKey(), NewEntry); putSetErr != nil {
			log.Println("redisx.SetJSON Loyalty_Plan error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Plan",
		Event_ActionType:   "Add",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  NewEntry,
	})
	return Id, nil
}

func (Uc *UserControl) Loyalty_Plan_Edit(Login string, request Loyalty_Plan_EditRequest) (Id int64, err error) {
	//check and validate
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	entry, entryErr := redisx.GetJSON[Loyalty_Plan](context.Background(), RedisClient, Loyalty_Plan{Key: request.Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("key is not created")
		return Id, err
	}
	if entryErr != nil {
		return Id, entryErr
	}
	if entry.Plan_Id != request.Plan_Id {
		return Id, errors.New("id is not matching")
	}
	planEntries2, planEntriesErr2 := redisx.GetAllJSONByPattern[Loyalty_Plan](context.Background(), RedisClient, redisx.ScanJSONOptions{
		Pattern: "Loyalty_Plan:*", ScanCount: 500, PipelineSize: 250,
	})
	if planEntriesErr2 != nil {
		err = planEntriesErr2
		return Id, err
	}
	for _, planEntry2 := range planEntries2 {
		if request.Key != planEntry2.Key && request.Loyalty_Level_Key == planEntry2.Loyalty_Level_Key && request.Loyalty_Account_Segment_Key == planEntry2.Loyalty_Account_Segment_Key {
			err = errors.New("a plan with the same Loyalty Level and Account Segment already exists. Please choose a different combination")
			return Id, err
		}
	}
	Current_Entry := entry
	//Prepare new entry
	entry.Key = request.Key
	entry.Description = request.Description
	entry.Loyalty_Level_Key = request.Loyalty_Level_Key
	entry.Loyalty_Account_Segment_Key = request.Loyalty_Account_Segment_Key
	entry.Earning_Rules_Key = request.Earning_Rules_Key
	entry.Expiry_Rules_Key = request.Expiry_Rules_Key
	entry.Redemption_Rules_Key = request.Redemption_Rules_Key

	if request.NewKey != "" {
		if request.NewKey != request.Key {
			//delete old
			{
				delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
				if _, delErr := Mdb_Loyalty_Plan.Coll.DeleteOne(delCtx, bson.M{"Key": request.Key}); delErr != nil {
					log.Println("Mdb_Loyalty_Plan DeleteOne error:", delErr)
				}
				if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Loyalty_Plan{Key: request.Key}.RedisKey()); delRedisErr != nil {
					log.Println("redisx.DelJSON Loyalty_Plan error:", delRedisErr)
				}
				delCancel()
			}
			//update key
			entry.Key = request.NewKey
		}
	}
	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Loyalty_Plan.Coll.UpdateOne(putCtx, bson.M{"Key": entry.Key}, bson.M{"$set": entry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Loyalty_Plan upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, entry.RedisKey(), entry); putSetErr != nil {
			log.Println("redisx.SetJSON Loyalty_Plan error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Plan",
		Event_ActionType:   "Edit",
		Event_Description:  "",
		Event_Entry_Before: Current_Entry,
		Event_Entry_After:  entry,
	})
	return Id, nil
}

func (Uc *UserControl) Loyalty_Plan_Get(Key string) (entries []Loyalty_Plan, err error) {
	if Key == "" {
		entries, err = redisx.GetAllJSONByPattern[Loyalty_Plan](context.Background(), RedisClient, redisx.ScanJSONOptions{
			Pattern: "Loyalty_Plan:*", ScanCount: 500, PipelineSize: 250,
		})
		if err != nil {
			return entries, err
		}
		return entries, nil
	} else {
		entry, entryErr := redisx.GetJSON[Loyalty_Plan](context.Background(), RedisClient, Loyalty_Plan{Key: Key}.RedisKey())
		if redisx.IsNil(entryErr) {
			err = errors.New("key does not exist")
			return entries, err
		}
		if entryErr != nil {
			return entries, entryErr
		}
		entries = append(entries, entry)
		return entries, nil
	}
}

func (Uc *UserControl) Loyalty_Plan_GetPaginated(Page, Limit int64) (entries []Loyalty_Plan, err error) {
	if Page < 1 {
		return entries, errors.New("invalid page")
	}
	if Limit < 1 || Limit > 50000 {
		return entries, errors.New("invalid limit (accept value between 1 and 50000)")
	}

	//var array []daoc.DAOFindCriteria
	// if Outlet_Key != "" {
	// 	//restrict access for records that belong to this user
	// 	var criteria daoc.DAOFindCriteria = daoc.DAOFindCriteria{
	// 		Field:    "Outlet_Key",
	// 		Value:    Outlet_Key,
	// 		Operator: "EQUAL",
	// 	}
	// 	array = append(array, criteria)
	// }
	// if Agent_Key != "" {
	// 	//restrict access for records that belong to this user
	// 	var criteria daoc.DAOFindCriteria = daoc.DAOFindCriteria{
	// 		Field:    "Agent_Key",
	// 		Value:    Agent_Key,
	// 		Operator: "EQUAL",
	// 	}
	// 	array = append(array, criteria)
	// }
	// if len(array) > 0 {
	// 	findparams.FindCriteria = array
	// }
	pgCtx, pgCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer pgCancel()
	result, err := mongox.FindPage[Loyalty_Plan](pgCtx, Mdb_Loyalty_Plan.Coll, bson.M{}, mongox.PageRequest{
		Page: Page, PageSize: Limit, MaxPageSize: 50000,
	})
	if err != nil {
		return entries, err
	}
	entries = result.Items
	return entries, nil

}

func (Uc *UserControl) Loyalty_Plan_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry, entryErr := redisx.GetJSON[Loyalty_Plan](context.Background(), RedisClient, Loyalty_Plan{Key: Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("entry does not exist")
		return err
	}
	if entryErr != nil {
		return entryErr
	}
	{
		delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, delErr := Mdb_Loyalty_Plan.Coll.DeleteOne(delCtx, bson.M{"Key": Key}); delErr != nil {
			log.Println("Mdb_Loyalty_Plan DeleteOne error:", delErr)
		}
		if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Loyalty_Plan{Key: Key}.RedisKey()); delRedisErr != nil {
			log.Println("redisx.DelJSON Loyalty_Plan error:", delRedisErr)
		}
		delCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Plan",
		Event_ActionType:   "Delete",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  entry,
	})
	return nil
}

// ***********************************************************************
// Customer UAT functions
// ***********************************************************************
func (Uc *UserControl) Customer_UAT_Add(Login string, request Customer_UAT_AddRequest) (Id int64, err error) {
	//check key if filled and if already used
	request.Key = Normalize_International_MSISDN(request.Key)
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	//check if key already used
	{
		chkCtx, chkCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, chkErr := redisx.GetJSON[Customer_UAT](chkCtx, RedisClient, Customer_UAT{Key: request.Key}.RedisKey())
		chkCancel()
		if chkErr == nil {
			err = errors.New("key already exist")
			return Id, err
		}
	}
	//Prepare new entry
	var NewEntry Customer_UAT
	{
		aiCtx, aiCancel := context.WithTimeout(context.Background(), 10*time.Second)
		NewEntry.Id, err = redisx.NextAutoIncrementID(aiCtx, RedisClient, Mdb_Loyalty_AutoIncrement.Coll, "Customer_UAT-Id", redisx.NextIDOptions{
			RedisBase: "AutoIncrement", MongoRetries: 3, RetryBackoff: 500 * time.Millisecond,
		})
		aiCancel()
		if err != nil {
			return Id, err
		}
		Id = NewEntry.Id
	}
	NewEntry.Key = request.Key
	NewEntry.AddTime = time.Now()
	NewEntry.AddReason = request.AddReason
	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Customer_UAT.Coll.UpdateOne(putCtx, bson.M{"Key": NewEntry.Key}, bson.M{"$set": NewEntry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Customer_UAT upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, NewEntry.RedisKey(), NewEntry); putSetErr != nil {
			log.Println("redisx.SetJSON Customer_UAT error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Customer_UAT",
		Event_ActionType:   "Add",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  NewEntry,
	})
	return Id, nil
}

func (Uc *UserControl) Customer_UAT_Edit(Login string, request Customer_UAT_EditRequest) (Id int64, err error) {
	//check and validate
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	entry, entryErr := redisx.GetJSON[Customer_UAT](context.Background(), RedisClient, Customer_UAT{Key: request.Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("key is not created")
		return Id, err
	}
	if entryErr != nil {
		return Id, entryErr
	}
	if entry.Id != request.Id {
		return Id, errors.New("id is not matching")
	}
	Current_Entry := entry

	//Prepare new entry
	entry.Key = request.Key
	entry.AddReason = request.AddReason

	if request.NewKey != "" {
		if request.NewKey != request.Key {
			//delete old
			{
				delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
				if _, delErr := Mdb_Customer_UAT.Coll.DeleteOne(delCtx, bson.M{"Key": request.Key}); delErr != nil {
					log.Println("Mdb_Customer_UAT DeleteOne error:", delErr)
				}
				if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Customer_UAT{Key: request.Key}.RedisKey()); delRedisErr != nil {
					log.Println("redisx.DelJSON Customer_UAT error:", delRedisErr)
				}
				delCancel()
			}
			//update key
			entry.Key = request.NewKey
		}
	}
	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Customer_UAT.Coll.UpdateOne(putCtx, bson.M{"Key": entry.Key}, bson.M{"$set": entry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Customer_UAT upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, entry.RedisKey(), entry); putSetErr != nil {
			log.Println("redisx.SetJSON Customer_UAT error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Customer_UAT",
		Event_ActionType:   "Edit",
		Event_Description:  "",
		Event_Entry_Before: Current_Entry,
		Event_Entry_After:  entry,
	})
	return Id, nil
}

func (Uc *UserControl) Customer_UAT_Get(Key string) (entries []Customer_UAT, err error) {
	if Key == "" {
		entries, err = redisx.GetAllJSONByPattern[Customer_UAT](context.Background(), RedisClient, redisx.ScanJSONOptions{
			Pattern: "Customer_UAT:*", ScanCount: 500, PipelineSize: 250,
		})
		if err != nil {
			return entries, err
		}
		return entries, nil
	} else {
		entry, entryErr := redisx.GetJSON[Customer_UAT](context.Background(), RedisClient, Customer_UAT{Key: Key}.RedisKey())
		if redisx.IsNil(entryErr) {
			err = errors.New("key does not exist")
			return entries, err
		}
		if entryErr != nil {
			return entries, entryErr
		}
		entries = append(entries, entry)
		return entries, nil
	}
}

func (Uc *UserControl) Customer_UAT_GetPaginated(Page, Limit int64) (entries []Customer_UAT, err error) {
	if Page < 1 {
		return entries, errors.New("invalid page")
	}
	if Limit < 1 || Limit > 50000 {
		return entries, errors.New("invalid limit (accept value between 1 and 50000)")
	}

	//var array []daoc.DAOFindCriteria
	// if Outlet_Key != "" {
	// 	//restrict access for records that belong to this user
	// 	var criteria daoc.DAOFindCriteria = daoc.DAOFindCriteria{
	// 		Field:    "Outlet_Key",
	// 		Value:    Outlet_Key,
	// 		Operator: "EQUAL",
	// 	}
	// 	array = append(array, criteria)
	// }
	// if Agent_Key != "" {
	// 	//restrict access for records that belong to this user
	// 	var criteria daoc.DAOFindCriteria = daoc.DAOFindCriteria{
	// 		Field:    "Agent_Key",
	// 		Value:    Agent_Key,
	// 		Operator: "EQUAL",
	// 	}
	// 	array = append(array, criteria)
	// }
	// if len(array) > 0 {
	// 	findparams.FindCriteria = array
	// }
	pgCtx, pgCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer pgCancel()
	result, err := mongox.FindPage[Customer_UAT](pgCtx, Mdb_Customer_UAT.Coll, bson.M{}, mongox.PageRequest{
		Page: Page, PageSize: Limit, MaxPageSize: 50000,
	})
	if err != nil {
		return entries, err
	}
	entries = result.Items
	return entries, nil

}

func (Uc *UserControl) Customer_UAT_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry, entryErr := redisx.GetJSON[Customer_UAT](context.Background(), RedisClient, Customer_UAT{Key: Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("entry does not exist")
		return err
	}
	if entryErr != nil {
		return entryErr
	}
	{
		delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, delErr := Mdb_Customer_UAT.Coll.DeleteOne(delCtx, bson.M{"Key": Key}); delErr != nil {
			log.Println("Mdb_Customer_UAT DeleteOne error:", delErr)
		}
		if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Customer_UAT{Key: Key}.RedisKey()); delRedisErr != nil {
			log.Println("redisx.DelJSON Customer_UAT error:", delRedisErr)
		}
		delCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Customer_UAT",
		Event_ActionType:   "Delete",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  entry,
	})
	return nil
}

// ***********************************************************************
// Customer UAT functions
// ***********************************************************************
func (Uc *UserControl) Customer_DND_Add(Login string, request Customer_DND_AddRequest) (Id int64, err error) {
	request.Key = Normalize_International_MSISDN(request.Key)
	//check key if filled and if already used
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	//check if key already used
	{
		chkCtx, chkCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, chkErr := redisx.GetJSON[Customer_DND](chkCtx, RedisClient, Customer_DND{Key: request.Key}.RedisKey())
		chkCancel()
		if chkErr == nil {
			err = errors.New("key already exist")
			return Id, err
		}
	}
	//Prepare new entry
	var NewEntry Customer_DND
	{
		aiCtx, aiCancel := context.WithTimeout(context.Background(), 10*time.Second)
		NewEntry.Id, err = redisx.NextAutoIncrementID(aiCtx, RedisClient, Mdb_Loyalty_AutoIncrement.Coll, "Customer_DND-Id", redisx.NextIDOptions{
			RedisBase: "AutoIncrement", MongoRetries: 3, RetryBackoff: 500 * time.Millisecond,
		})
		aiCancel()
		if err != nil {
			return Id, err
		}
		Id = NewEntry.Id
	}
	NewEntry.Key = request.Key
	NewEntry.AddTime = time.Now()
	NewEntry.AddReason = request.AddReason
	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Customer_DND.Coll.UpdateOne(putCtx, bson.M{"Key": NewEntry.Key}, bson.M{"$set": NewEntry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Customer_DND upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, NewEntry.RedisKey(), NewEntry); putSetErr != nil {
			log.Println("redisx.SetJSON Customer_DND error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Customer_DND",
		Event_ActionType:   "Add",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  NewEntry,
	})
	return Id, nil
}

func (Uc *UserControl) Customer_DND_Edit(Login string, request Customer_DND_EditRequest) (Id int64, err error) {
	//check and validate
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	entry, entryErr := redisx.GetJSON[Customer_DND](context.Background(), RedisClient, Customer_DND{Key: request.Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("key is not created")
		return Id, err
	}
	if entryErr != nil {
		return Id, entryErr
	}
	if entry.Id != request.Id {
		return Id, errors.New("id is not matching")
	}
	Current_Entry := entry

	//Prepare new entry
	entry.Key = request.Key
	entry.AddReason = request.AddReason

	if request.NewKey != "" {
		if request.NewKey != request.Key {
			//delete old
			{
				delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
				if _, delErr := Mdb_Customer_DND.Coll.DeleteOne(delCtx, bson.M{"Key": request.Key}); delErr != nil {
					log.Println("Mdb_Customer_DND DeleteOne error:", delErr)
				}
				if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Customer_DND{Key: request.Key}.RedisKey()); delRedisErr != nil {
					log.Println("redisx.DelJSON Customer_DND error:", delRedisErr)
				}
				delCancel()
			}
			//update key
			entry.Key = request.NewKey
		}
	}
	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Customer_DND.Coll.UpdateOne(putCtx, bson.M{"Key": entry.Key}, bson.M{"$set": entry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Customer_DND upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, entry.RedisKey(), entry); putSetErr != nil {
			log.Println("redisx.SetJSON Customer_DND error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Customer_DND",
		Event_ActionType:   "Edit",
		Event_Description:  "",
		Event_Entry_Before: Current_Entry,
		Event_Entry_After:  entry,
	})
	return Id, nil
}

func (Uc *UserControl) Customer_DND_Get(Key string) (entries []Customer_DND, err error) {
	if Key == "" {
		entries, err = redisx.GetAllJSONByPattern[Customer_DND](context.Background(), RedisClient, redisx.ScanJSONOptions{
			Pattern: "Customer_DND:*", ScanCount: 500, PipelineSize: 250,
		})
		if err != nil {
			return entries, err
		}
		return entries, nil
	} else {
		entry, entryErr := redisx.GetJSON[Customer_DND](context.Background(), RedisClient, Customer_DND{Key: Key}.RedisKey())
		if redisx.IsNil(entryErr) {
			err = errors.New("key does not exist")
			return entries, err
		}
		if entryErr != nil {
			return entries, entryErr
		}
		entries = append(entries, entry)
		return entries, nil
	}
}

func (Uc *UserControl) Customer_DND_GetPaginated(Page, Limit int64) (entries []Customer_DND, err error) {
	if Page < 1 {
		return entries, errors.New("invalid page")
	}
	if Limit < 1 || Limit > 50000 {
		return entries, errors.New("invalid limit (accept value between 1 and 50000)")
	}

	//var array []daoc.DAOFindCriteria
	// if Outlet_Key != "" {
	// 	//restrict access for records that belong to this user
	// 	var criteria daoc.DAOFindCriteria = daoc.DAOFindCriteria{
	// 		Field:    "Outlet_Key",
	// 		Value:    Outlet_Key,
	// 		Operator: "EQUAL",
	// 	}
	// 	array = append(array, criteria)
	// }
	// if Agent_Key != "" {
	// 	//restrict access for records that belong to this user
	// 	var criteria daoc.DAOFindCriteria = daoc.DAOFindCriteria{
	// 		Field:    "Agent_Key",
	// 		Value:    Agent_Key,
	// 		Operator: "EQUAL",
	// 	}
	// 	array = append(array, criteria)
	// }
	// if len(array) > 0 {
	// 	findparams.FindCriteria = array
	// }
	pgCtx, pgCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer pgCancel()
	result, err := mongox.FindPage[Customer_DND](pgCtx, Mdb_Customer_DND.Coll, bson.M{}, mongox.PageRequest{
		Page: Page, PageSize: Limit, MaxPageSize: 50000,
	})
	if err != nil {
		return entries, err
	}
	entries = result.Items
	return entries, nil

}

func (Uc *UserControl) Customer_DND_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry, entryErr := redisx.GetJSON[Customer_DND](context.Background(), RedisClient, Customer_DND{Key: Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("entry does not exist")
		return err
	}
	if entryErr != nil {
		return entryErr
	}
	{
		delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, delErr := Mdb_Customer_DND.Coll.DeleteOne(delCtx, bson.M{"Key": Key}); delErr != nil {
			log.Println("Mdb_Customer_DND DeleteOne error:", delErr)
		}
		if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Customer_DND{Key: Key}.RedisKey()); delRedisErr != nil {
			log.Println("redisx.DelJSON Customer_DND error:", delRedisErr)
		}
		delCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Customer_DND",
		Event_ActionType:   "Delete",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  entry,
	})
	return nil
}

// ***********************************************************************
// Customer Exclusion functions
// ***********************************************************************
func (Uc *UserControl) Customer_Exclusion_Add(Login string, request Customer_Exclusion_AddRequest) (Id int64, err error) {
	request.Key = Normalize_International_MSISDN(request.Key)
	//check key if filled and if already used
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	//check if key already used
	{
		chkCtx, chkCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, chkErr := redisx.GetJSON[Customer_Exclusion](chkCtx, RedisClient, Customer_Exclusion{Key: request.Key}.RedisKey())
		chkCancel()
		if chkErr == nil {
			err = errors.New("key already exist")
			return Id, err
		}
	}
	//Prepare new entry
	var NewEntry Customer_Exclusion
	{
		aiCtx, aiCancel := context.WithTimeout(context.Background(), 10*time.Second)
		NewEntry.Id, err = redisx.NextAutoIncrementID(aiCtx, RedisClient, Mdb_Loyalty_AutoIncrement.Coll, "Customer_Exclusion-Id", redisx.NextIDOptions{
			RedisBase: "AutoIncrement", MongoRetries: 3, RetryBackoff: 500 * time.Millisecond,
		})
		aiCancel()
		if err != nil {
			return Id, err
		}
		Id = NewEntry.Id
	}
	NewEntry.Key = request.Key
	NewEntry.AddTime = time.Now()
	NewEntry.AddReason = request.AddReason
	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Customer_Exclusion.Coll.UpdateOne(putCtx, bson.M{"Key": NewEntry.Key}, bson.M{"$set": NewEntry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Customer_Exclusion upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, NewEntry.RedisKey(), NewEntry); putSetErr != nil {
			log.Println("redisx.SetJSON Customer_Exclusion error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Customer_Exclusion",
		Event_ActionType:   "Add",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  NewEntry,
	})
	return Id, nil
}

func (Uc *UserControl) Customer_Exclusion_Edit(Login string, request Customer_Exclusion_EditRequest) (Id int64, err error) {
	//check and validate
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	entry, entryErr := redisx.GetJSON[Customer_Exclusion](context.Background(), RedisClient, Customer_Exclusion{Key: request.Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("key is not created")
		return Id, err
	}
	if entryErr != nil {
		return Id, entryErr
	}
	if entry.Id != request.Id {
		return Id, errors.New("id is not matching")
	}
	Current_Entry := entry

	//Prepare new entry
	entry.Key = request.Key
	entry.AddReason = request.AddReason

	if request.NewKey != "" {
		if request.NewKey != request.Key {
			//delete old
			{
				delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
				if _, delErr := Mdb_Customer_Exclusion.Coll.DeleteOne(delCtx, bson.M{"Key": request.Key}); delErr != nil {
					log.Println("Mdb_Customer_Exclusion DeleteOne error:", delErr)
				}
				if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Customer_Exclusion{Key: request.Key}.RedisKey()); delRedisErr != nil {
					log.Println("redisx.DelJSON Customer_Exclusion error:", delRedisErr)
				}
				delCancel()
			}
			//update key
			entry.Key = request.NewKey
		}
	}
	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Customer_Exclusion.Coll.UpdateOne(putCtx, bson.M{"Key": entry.Key}, bson.M{"$set": entry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Customer_Exclusion upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, entry.RedisKey(), entry); putSetErr != nil {
			log.Println("redisx.SetJSON Customer_Exclusion error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Customer_Exclusion",
		Event_ActionType:   "Edit",
		Event_Description:  "",
		Event_Entry_Before: Current_Entry,
		Event_Entry_After:  entry,
	})
	return Id, nil
}

func (Uc *UserControl) Customer_Exclusion_Get(Key string) (entries []Customer_Exclusion, err error) {
	if Key == "" {
		entries, err = redisx.GetAllJSONByPattern[Customer_Exclusion](context.Background(), RedisClient, redisx.ScanJSONOptions{
			Pattern: "Customer_Exclusion:*", ScanCount: 500, PipelineSize: 250,
		})
		if err != nil {
			return entries, err
		}
		return entries, nil
	} else {
		entry, entryErr := redisx.GetJSON[Customer_Exclusion](context.Background(), RedisClient, Customer_Exclusion{Key: Key}.RedisKey())
		if redisx.IsNil(entryErr) {
			err = errors.New("key does not exist")
			return entries, err
		}
		if entryErr != nil {
			return entries, entryErr
		}
		entries = append(entries, entry)
		return entries, nil
	}
}

func (Uc *UserControl) Customer_Exclusion_GetPaginated(Page, Limit int64) (entries []Customer_Exclusion, err error) {
	if Page < 1 {
		return entries, errors.New("invalid page")
	}
	if Limit < 1 || Limit > 50000 {
		return entries, errors.New("invalid limit (accept value between 1 and 50000)")
	}

	//var array []daoc.DAOFindCriteria
	// if Outlet_Key != "" {
	// 	//restrict access for records that belong to this user
	// 	var criteria daoc.DAOFindCriteria = daoc.DAOFindCriteria{
	// 		Field:    "Outlet_Key",
	// 		Value:    Outlet_Key,
	// 		Operator: "EQUAL",
	// 	}
	// 	array = append(array, criteria)
	// }
	// if Agent_Key != "" {
	// 	//restrict access for records that belong to this user
	// 	var criteria daoc.DAOFindCriteria = daoc.DAOFindCriteria{
	// 		Field:    "Agent_Key",
	// 		Value:    Agent_Key,
	// 		Operator: "EQUAL",
	// 	}
	// 	array = append(array, criteria)
	// }
	// if len(array) > 0 {
	// 	findparams.FindCriteria = array
	// }
	pgCtx, pgCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer pgCancel()
	result, err := mongox.FindPage[Customer_Exclusion](pgCtx, Mdb_Customer_Exclusion.Coll, bson.M{}, mongox.PageRequest{
		Page: Page, PageSize: Limit, MaxPageSize: 50000,
	})
	if err != nil {
		return entries, err
	}
	entries = result.Items
	return entries, nil

}

func (Uc *UserControl) Customer_Exclusion_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry, entryErr := redisx.GetJSON[Customer_Exclusion](context.Background(), RedisClient, Customer_Exclusion{Key: Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("entry does not exist")
		return err
	}
	if entryErr != nil {
		return entryErr
	}
	{
		delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, delErr := Mdb_Customer_Exclusion.Coll.DeleteOne(delCtx, bson.M{"Key": Key}); delErr != nil {
			log.Println("Mdb_Customer_Exclusion DeleteOne error:", delErr)
		}
		if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Customer_Exclusion{Key: Key}.RedisKey()); delRedisErr != nil {
			log.Println("redisx.DelJSON Customer_Exclusion error:", delRedisErr)
		}
		delCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Customer_Exclusion",
		Event_ActionType:   "Delete",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  entry,
	})
	return nil
}

// ***********************************************************************
// Customer Exclusion functions
// ***********************************************************************
func (Uc *UserControl) Customer_COS_Exclusion_Add(Login string, request Customer_COS_Exclusion_AddRequest) (Id int64, err error) {
	//check key if filled and if already used
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	//check if key already used
	{
		chkCtx, chkCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, chkErr := redisx.GetJSON[Customer_COS_Exclusion](chkCtx, RedisClient, Customer_COS_Exclusion{Key: request.Key}.RedisKey())
		chkCancel()
		if chkErr == nil {
			err = errors.New("key already exist")
			return Id, err
		}
	}
	//Prepare new entry
	var NewEntry Customer_COS_Exclusion
	{
		aiCtx, aiCancel := context.WithTimeout(context.Background(), 10*time.Second)
		NewEntry.Id, err = redisx.NextAutoIncrementID(aiCtx, RedisClient, Mdb_Loyalty_AutoIncrement.Coll, "Customer_COS_Exclusion-Id", redisx.NextIDOptions{
			RedisBase: "AutoIncrement", MongoRetries: 3, RetryBackoff: 500 * time.Millisecond,
		})
		aiCancel()
		if err != nil {
			return Id, err
		}
		Id = NewEntry.Id
	}
	NewEntry.Key = request.Key
	NewEntry.AddTime = time.Now()
	NewEntry.AddReason = request.AddReason
	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Customer_COS_Exclusion.Coll.UpdateOne(putCtx, bson.M{"Key": NewEntry.Key}, bson.M{"$set": NewEntry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Customer_COS_Exclusion upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, NewEntry.RedisKey(), NewEntry); putSetErr != nil {
			log.Println("redisx.SetJSON Customer_COS_Exclusion error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Customer_COS_Exclusion",
		Event_ActionType:   "Add",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  NewEntry,
	})
	return Id, nil
}

func (Uc *UserControl) Customer_COS_Exclusion_Edit(Login string, request Customer_COS_Exclusion_EditRequest) (Id int64, err error) {
	//check and validate
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	entry, entryErr := redisx.GetJSON[Customer_COS_Exclusion](context.Background(), RedisClient, Customer_COS_Exclusion{Key: request.Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("key is not created")
		return Id, err
	}
	if entryErr != nil {
		return Id, entryErr
	}
	if entry.Id != request.Id {
		return Id, errors.New("id is not matching")
	}
	Current_Entry := entry

	//Prepare new entry
	entry.Key = request.Key
	entry.AddReason = request.AddReason

	if request.NewKey != "" {
		if request.NewKey != request.Key {
			//delete old
			{
				delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
				if _, delErr := Mdb_Customer_COS_Exclusion.Coll.DeleteOne(delCtx, bson.M{"Key": request.Key}); delErr != nil {
					log.Println("Mdb_Customer_COS_Exclusion DeleteOne error:", delErr)
				}
				if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Customer_COS_Exclusion{Key: request.Key}.RedisKey()); delRedisErr != nil {
					log.Println("redisx.DelJSON Customer_COS_Exclusion error:", delRedisErr)
				}
				delCancel()
			}
			//update key
			entry.Key = request.NewKey
		}
	}
	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Customer_COS_Exclusion.Coll.UpdateOne(putCtx, bson.M{"Key": entry.Key}, bson.M{"$set": entry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Customer_COS_Exclusion upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, entry.RedisKey(), entry); putSetErr != nil {
			log.Println("redisx.SetJSON Customer_COS_Exclusion error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Customer_COS_Exclusion",
		Event_ActionType:   "Edit",
		Event_Description:  "",
		Event_Entry_Before: Current_Entry,
		Event_Entry_After:  entry,
	})
	return Id, nil
}

func (Uc *UserControl) Customer_COS_Exclusion_Get(Key string) (entries []Customer_COS_Exclusion, err error) {
	if Key == "" {
		entries, err = redisx.GetAllJSONByPattern[Customer_COS_Exclusion](context.Background(), RedisClient, redisx.ScanJSONOptions{
			Pattern: "Customer_COS_Exclusion:*", ScanCount: 500, PipelineSize: 250,
		})
		if err != nil {
			return entries, err
		}
		return entries, nil
	} else {
		entry, entryErr := redisx.GetJSON[Customer_COS_Exclusion](context.Background(), RedisClient, Customer_COS_Exclusion{Key: Key}.RedisKey())
		if redisx.IsNil(entryErr) {
			err = errors.New("key does not exist")
			return entries, err
		}
		if entryErr != nil {
			return entries, entryErr
		}
		entries = append(entries, entry)
		return entries, nil
	}
}

func (Uc *UserControl) Customer_COS_Exclusion_GetPaginated(Page, Limit int64) (entries []Customer_COS_Exclusion, err error) {
	if Page < 1 {
		return entries, errors.New("invalid page")
	}
	if Limit < 1 || Limit > 50000 {
		return entries, errors.New("invalid limit (accept value between 1 and 50000)")
	}

	//var array []daoc.DAOFindCriteria
	// if Outlet_Key != "" {
	// 	//restrict access for records that belong to this user
	// 	var criteria daoc.DAOFindCriteria = daoc.DAOFindCriteria{
	// 		Field:    "Outlet_Key",
	// 		Value:    Outlet_Key,
	// 		Operator: "EQUAL",
	// 	}
	// 	array = append(array, criteria)
	// }
	// if Agent_Key != "" {
	// 	//restrict access for records that belong to this user
	// 	var criteria daoc.DAOFindCriteria = daoc.DAOFindCriteria{
	// 		Field:    "Agent_Key",
	// 		Value:    Agent_Key,
	// 		Operator: "EQUAL",
	// 	}
	// 	array = append(array, criteria)
	// }
	// if len(array) > 0 {
	// 	findparams.FindCriteria = array
	// }
	pgCtx, pgCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer pgCancel()
	result, err := mongox.FindPage[Customer_COS_Exclusion](pgCtx, Mdb_Customer_COS_Exclusion.Coll, bson.M{}, mongox.PageRequest{
		Page: Page, PageSize: Limit, MaxPageSize: 50000,
	})
	if err != nil {
		return entries, err
	}
	entries = result.Items
	return entries, nil

}

func (Uc *UserControl) Customer_COS_Exclusion_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry, entryErr := redisx.GetJSON[Customer_COS_Exclusion](context.Background(), RedisClient, Customer_COS_Exclusion{Key: Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("entry does not exist")
		return err
	}
	if entryErr != nil {
		return entryErr
	}
	{
		delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, delErr := Mdb_Customer_COS_Exclusion.Coll.DeleteOne(delCtx, bson.M{"Key": Key}); delErr != nil {
			log.Println("Mdb_Customer_COS_Exclusion DeleteOne error:", delErr)
		}
		if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Customer_COS_Exclusion{Key: Key}.RedisKey()); delRedisErr != nil {
			log.Println("redisx.DelJSON Customer_COS_Exclusion error:", delRedisErr)
		}
		delCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Customer_COS_Exclusion",
		Event_ActionType:   "Delete",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  entry,
	})
	return nil
}

// ***********************************************************************
// Customer_Loyalty_Account functions
// ***********************************************************************
func (Uc *UserControl) Customer_Loyalty_Account_Add(Login string, request Customer_Loyalty_Account_AddRequest) (Id int64, err error) {
	//check key if filled and if already used
	request.Key = Normalize_International_MSISDN(request.Key)
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	//check if key already used
	{
		chkCtx, chkCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, chkErr := redisx.GetJSON[Customer_Loyalty_Account](chkCtx, RedisClient, Customer_Loyalty_Account{Key: request.Key}.RedisKey())
		chkCancel()
		if chkErr == nil {
			err = errors.New("key already exist")
			return Id, err
		}
	}
	//check exclusion list
	{
		chkCtx, chkCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, chkErr := redisx.GetJSON[Customer_Exclusion](chkCtx, RedisClient, Customer_Exclusion{Key: request.Key}.RedisKey())
		chkCancel()
		if chkErr == nil {
			err = errors.New("customer is included in the exclusion list")
			return Id, err
		}
	}
	//check COS exclusion list
	{
		chkCtx, chkCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, chkErr := redisx.GetJSON[Customer_COS_Exclusion](chkCtx, RedisClient, Customer_COS_Exclusion{Key: request.COS}.RedisKey())
		chkCancel()
		if chkErr == nil {
			err = errors.New("customer is included in the cos exclusion list")
			return Id, err
		}
	}
	//Prepare new entry
	var NewEntry Customer_Loyalty_Account
	{
		aiCtx, aiCancel := context.WithTimeout(context.Background(), 10*time.Second)
		NewEntry.Customer_Id, err = redisx.NextAutoIncrementID(aiCtx, RedisClient, Mdb_Loyalty_AutoIncrement.Coll, "Customer_Loyalty_Account-Id", redisx.NextIDOptions{
			RedisBase: "AutoIncrement", MongoRetries: 3, RetryBackoff: 500 * time.Millisecond,
		})
		aiCancel()
		if err != nil {
			return Id, err
		}
		Id = NewEntry.Customer_Id
	}
	NewEntry.Key = request.Key

	NewEntry.COS = request.COS
	NewEntry.ARPU = request.ARPU
	NewEntry.Joining_Date = request.Joining_Date

	NewEntry.Creation_date = time.Now()
	NewEntry.Account_Status = "active"
	NewEntry.Account_Profile = ""

	NewEntry.Loyalty_Level_Key = Loyalty_Level_Selection(0)
	NewEntry.Loyalty_Level_Date = time.Now()
	NewEntry.Loyalty_Level_Direction = ""
	NewEntry.Loyalty_Level_SetBy = Login //program or admin, if admin program cannot change anymore
	if Configuration.ISLoyaltyOptIn {
		NewEntry.Opt_Status = "OptedOut"
	} else {
		NewEntry.Opt_Status = "OptedIn"
		NewEntry.First_Opt_In_Status_Date = time.Now()
	}
	NewEntry.Last_Opt_Status_Date = time.Now()

	if Login != "DWH_Import" {
		NewEntry.Loyalty_Account_Segment_Key = Loyalty_Account_Segment_Selection(request.ARPU, request.Joining_Date)
		NewEntry.Loyalty_Account_Segment_Date = time.Now()
		NewEntry.Loyalty_Account_Segment_Direction = ""
		NewEntry.Loyalty_Account_Segment_SetBy = Login
	} else {
		subscriber, err := Uc.Lendme.LendmeClient.Lendme_Subscriber_Get(request.Key)
		if err == nil {
			NewEntry.Loyalty_Account_Segment_Key = Loyalty_Account_Segment_Selection(subscriber.ARPU, subscriber.FirstUse_date)
		} else {
			NewEntry.Loyalty_Account_Segment_Key = Loyalty_Account_Segment_Selection(request.ARPU, request.Joining_Date)
		}
		NewEntry.Loyalty_Account_Segment_Date = time.Now()
		NewEntry.Loyalty_Account_Segment_Direction = ""
		NewEntry.Loyalty_Account_Segment_SetBy = Login
	}

	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Customer_Loyalty_Account.Coll.UpdateOne(putCtx, bson.M{"Key": NewEntry.Key}, bson.M{"$set": NewEntry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Customer_Loyalty_Account upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, NewEntry.RedisKey(), NewEntry); putSetErr != nil {
			log.Println("redisx.SetJSON Customer_Loyalty_Account error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	if Login != "DWH_Import" && Login != "INLiveFeed" { //--> off to avoid filling up logs
		Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
			Event_User:         Login,
			Event_Time:         time.Now(),
			Event_AffectedType: "Customer_Loyalty_Account",
			Event_ActionType:   "Add",
			Event_Description:  "",
			Event_Entry_Before: nil,
			Event_Entry_After:  NewEntry,
		})
	}
	NewJoiningsCount.With(prometheus.Labels{"Source": request.EventSource}).Inc()
	if NewEntry.Opt_Status == "OptedIn" {
		var loyalty_AccountCreditPoints_log Loyalty_AccountCreditPoints_log
		var loyalty_AccountCreditPoints_Request Loyalty_AccountCreditPoints_Request
		loyalty_AccountCreditPoints_Request.MSISDN = NewEntry.Key
		loyalty_AccountCreditPoints_Request.EventSource = request.EventSource
		loyalty_AccountCreditPoints_Request.EventType = "NewJoining"
		loyalty_AccountCreditPoints_Request.EventDetail = ""
		loyalty_AccountCreditPoints_Request.EventAmount = 0
		loyalty_AccountCreditPoints_Request.EventDescription = ""
		var request_header Request_Header
		request_header.SourceIP = "127.0.0.1"
		request_header.SourceApp = request.EventSource
		request_header.AppLogin = Login
		Uc.Loyalty_AccountCreditPoints(&request_header, loyalty_AccountCreditPoints_Request, &loyalty_AccountCreditPoints_log)

		Earningrecord, rule_err := Uc.Customer_Loyalty_Account_GetEarning_Rule(request.Key)
		if rule_err != nil {
			log.Println("failed to get data")
			return
		}
		if Earningrecord.Welcome_Notification {
			WelcomeNotiLog := NotificationLog{
				SourceAction:  "Welcome",
				TransactionId: "",
				Medium:        "SMS",
				SourceAddress: Earningrecord.Welcome_Notification_Sender,
				Destination:   request.Key,
				Subject:       "Welcome",
				AddUser:       "SYSTEM",
				AddDate:       time.Now(),
			}
			Welcome_Noti_Text := ""
			Welcome_Noti_Text = Earningrecord.Welcome_Notification_Text
			if Welcome_Noti_Text != "" {
				Welcome_Noti_Text = strings.ReplaceAll(Welcome_Noti_Text, "{{WelcomePoints}}", fmt.Sprint(Earningrecord.Welcome_Points))
				Welcome_Noti_Text = strings.ReplaceAll(Welcome_Noti_Text, "{{LoyaltyBalance}}", fmt.Sprint(Earningrecord.Welcome_Points))
				Welcome_Noti_Text = strings.ReplaceAll(Welcome_Noti_Text, "{{NewLevel}}", fmt.Sprint(NewEntry.Loyalty_Level_Key))
				WelcomeNotiLog.Payload = Welcome_Noti_Text
				err := error(nil)
				if Configuration.Operation == "Angola" {
					err = SendSMS(Earningrecord.Welcome_Notification_Sender, request.Key, Welcome_Noti_Text)
				} else {
					err = Send_SMS(Earningrecord.Welcome_Notification_Sender, request.Key, Welcome_Noti_Text)
				}
				if err != nil {
					WelcomeNotiLog.Status = "Failed"
					WelcomeNotiLog.Error = err.Error()
				} else {
					WelcomeNotiLog.Status = "Successful"
				}
			} else {
				WelcomeNotiLog.Payload = Welcome_Noti_Text
				WelcomeNotiLog.Status = "Failed"
				WelcomeNotiLog.Error = "Undefined welcome notification for transaction"
			}
			YYYY, MM, _, _, _, _, _ := GetTimeParts(time.Now())
			Db := Configuration.DB_Name_Loyalty + "_Logs_" + YYYY + MM
			{
				notiCtx, notiCancel := context.WithTimeout(context.Background(), 10*time.Second)
				_, notiErr := Uc.MongoClient.Mongo.Database(Db).Collection("Col_NotificationLog").InsertOne(notiCtx, WelcomeNotiLog)
				notiCancel()
				if notiErr != nil {
					log.Println("Error in Write welcome Notification Logs:", notiErr, " (", WelcomeNotiLog, ")")
				}
			}
		}
	}
	return Id, nil

}

func (Uc *UserControl) Customer_Loyalty_Account_Edit(Login string, request Customer_Loyalty_Account_EditRequest) (Id int64, err error) {
	//check and validate
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	entry, entryErr := redisx.GetJSON[Customer_Loyalty_Account](context.Background(), RedisClient, Customer_Loyalty_Account{Key: request.Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("key is not created")
		return Id, err
	}
	if entryErr != nil {
		return Id, entryErr
	}
	if entry.Customer_Id != request.Customer_Id {
		return Id, errors.New("id is not matching")
	}
	Current_Entry := entry
	//check opt status
	if entry.Opt_Status != "OptedIn" && Login != "DWH_Import" && Configuration.ISLoyaltyOptIn {
		err = errors.New("customer is opted out")
		return Id, err
	}
	//check exclusion list
	{
		chkCtx, chkCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, chkErr := redisx.GetJSON[Customer_Exclusion](chkCtx, RedisClient, Customer_Exclusion{Key: request.Key}.RedisKey())
		chkCancel()
		if chkErr == nil {
			Uc.Customer_Loyalty_Account_Delete(Login, request.Key)
			err = errors.New("customer is included in the exclusion list")
			return Id, err
		}
	}
	//check COS exclusion list
	{
		chkCtx, chkCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, chkErr := redisx.GetJSON[Customer_COS_Exclusion](chkCtx, RedisClient, Customer_COS_Exclusion{Key: request.COS}.RedisKey())
		chkCancel()
		if chkErr == nil {
			Uc.Customer_Loyalty_Account_Delete(Login, request.Key)
			err = errors.New("customer is included in the cos exclusion list")
			return Id, err
		}
	}
	//Prepare new entry
	entry.Key = request.Key
	if entry.Loyalty_Level_Key != request.Loyalty_Level_Key {

		current_level, currentLvlErr := redisx.GetJSON[Loyalty_Level](context.Background(), RedisClient, Loyalty_Level{Key: entry.Loyalty_Level_Key}.RedisKey())
		if redisx.IsNil(currentLvlErr) {
			err = errors.New("current level is invalid")
			return Id, err
		}
		if currentLvlErr != nil {
			return Id, currentLvlErr
		}
		new_level, newLvlErr := redisx.GetJSON[Loyalty_Level](context.Background(), RedisClient, Loyalty_Level{Key: request.Loyalty_Level_Key}.RedisKey())
		if redisx.IsNil(newLvlErr) {
			err = errors.New("new level is invalid")
			return Id, err
		}
		if newLvlErr != nil {
			return Id, newLvlErr
		}

		if entry.Loyalty_Level_Key != request.Loyalty_Level_Key {
			entry.Previous_Loyalty_Level_Key = entry.Loyalty_Level_Key
			entry.Previous_Loyalty_Level_Date = entry.Loyalty_Level_Date
			entry.Loyalty_Level_Key = new_level.Key
			entry.Loyalty_Level_Date = time.Now()
			entry.Loyalty_Level_SetBy = Login
			if new_level.Min_Accumulated_Points > current_level.Min_Accumulated_Points &&
				new_level.Max_Accumulated_Points > current_level.Max_Accumulated_Points {
				entry.Loyalty_Level_Direction = "Upgrade"
			} else {
				entry.Loyalty_Level_Direction = "Downgrade"
			}

			Uc.Write_Loyalty_Level_Change_log(Loyalty_Level_Change_log{
				Level_Change_Date:                 time.Now(),
				MSISDN:                            entry.Key,
				COS:                               entry.COS,
				Joining_Date:                      entry.Joining_Date,
				ARPU:                              entry.ARPU,
				Customer_Id:                       entry.Customer_Id,
				Creation_date:                     entry.Creation_date,
				Previous_Loyalty_Level_Key:        entry.Previous_Loyalty_Level_Key,
				Previous_Loyalty_Level_Date:       entry.Previous_Loyalty_Level_Date,
				New_Loyalty_Level_Key:             entry.Loyalty_Level_Key,
				New_Loyalty_Level_Date:            entry.Loyalty_Level_Date,
				New_Loyalty_Level_Direction:       entry.Loyalty_Level_Direction,
				New_Loyalty_Level_SetBy:           entry.Loyalty_Level_SetBy,
				Loyalty_Account_Segment_Key:       entry.Loyalty_Account_Segment_Key,
				Loyalty_Account_Segment_Date:      entry.Loyalty_Account_Segment_Date,
				Loyalty_Account_Segment_Direction: entry.Loyalty_Account_Segment_Direction,
				Loyalty_Account_Segment_SetBy:     entry.Loyalty_Account_Segment_SetBy,
				Awarded_Points:                    entry.Awarded_Points,
				Redeemed_Points:                   entry.Redeemed_Points,
				Available_Points:                  entry.Available_Points,
				Last_Award_Date:                   entry.Last_Award_Date,
				Last_Redeem_Date:                  entry.Last_Redeem_Date,
			})
		}
	}
	if Configuration.ISLoyaltyOptIn && entry.First_Opt_In_Status_Date.IsZero() {
		entry.Opt_Status = "OptedOut"
		entry.Last_Opt_Status_Date = time.Now()
	} else if entry.First_Opt_In_Status_Date.IsZero() {
		entry.Opt_Status = "OptedIn"
		entry.First_Opt_In_Status_Date = time.Now()
		entry.Last_Opt_Status_Date = time.Now()
	}
	//evaluate the loyalty Account segment
	if Login == "DWH_Import" {
		entry.ARPU = request.ARPU
		entry.Joining_Date = request.Joining_Date
		Loyalty_Account_Segment_Key := Loyalty_Account_Segment_Selection(request.ARPU, request.Joining_Date)
		if Loyalty_Account_Segment_Key != entry.Loyalty_Account_Segment_Key {
			entry.Loyalty_Account_Segment_Key = Loyalty_Account_Segment_Key
			entry.Loyalty_Account_Segment_Date = time.Now()
			entry.Loyalty_Account_Segment_Direction = ""
			entry.Loyalty_Account_Segment_SetBy = Login
		}

		plan, planErr := redisx.GetJSON[Loyalty_Plan](context.Background(), RedisClient, Loyalty_Plan{Key: entry.Loyalty_Level_Key + "|" + entry.Loyalty_Account_Segment_Key}.RedisKey())
		if redisx.IsNil(planErr) {
			err = errors.New("loyalty plan does not exist")
			return Id, err
		}
		if planErr != nil {
			return Id, planErr
		}
		if plan.Expiry_Rules_Key == "" {
			err = errors.New("expiry rules is not defined")
			return Id, err
		}
		expiry_Rule, expiryErr := redisx.GetJSON[Loyalty_Point_Expiry_Rules](context.Background(), RedisClient, Loyalty_Point_Expiry_Rules{Key: plan.Expiry_Rules_Key}.RedisKey())
		if redisx.IsNil(expiryErr) {
			err = errors.New("expiry rules is not defined")
			return Id, err
		}
		if expiryErr != nil {
			return Id, expiryErr
		}
		var initialDate time.Time
		if !entry.Expiry_Date.IsZero() {
			initialDate = entry.Expiry_Date
		} else if !entry.First_Opt_In_Status_Date.IsZero() {
			initialDate = entry.First_Opt_In_Status_Date
		} else {
			initialDate = entry.Creation_date
		}
		initialexpiryDate := addValidity(initialDate, expiry_Rule.Validity_Unit, expiry_Rule.Validity_Duration)
		finalexpiryDate := addValidity(initialexpiryDate, expiry_Rule.Grace_Validity_Unit, expiry_Rule.Grace_Validity_Duration)
		entry.Coming_Expiry_Date = finalexpiryDate
		entry.Initial_Date = initialexpiryDate
		var expiryPoints float64 = 0
		for _, pointKey := range entry.Points_Detail_Keys {
			pointsDetail, err := Uc.Customer_Loyalty_Account_Points_Details_Get(pointKey)
			if err != nil {
				return Id, err
			}
			year, err := strconv.Atoi(pointsDetail[0].Year_Month[:4])
			if err != nil {
				return Id, err
			}
			month, err := strconv.Atoi(pointsDetail[0].Year_Month[4:])
			if err == nil && year < int(entry.Initial_Date.Year()) || (year == entry.Initial_Date.Year() && month < int(entry.Initial_Date.Month())) {
				expiryPoints += pointsDetail[0].Available_Points
			}
			entry.Points_To_Expire = expiryPoints
		}
	} else {
		subscriber, err := Uc.Lendme.LendmeClient.Lendme_Subscriber_Get(request.Key)
		var segKey string
		if err == nil {
			segKey = Loyalty_Account_Segment_Selection(subscriber.ARPU, subscriber.FirstUse_date)
		} else {
			segKey = Loyalty_Account_Segment_Selection(request.ARPU, request.Joining_Date)
		}
		if segKey != entry.Loyalty_Account_Segment_Key {
			entry.Loyalty_Account_Segment_Key = segKey
			entry.Loyalty_Account_Segment_Date = time.Now()
			entry.Loyalty_Account_Segment_Direction = ""
			entry.Loyalty_Account_Segment_SetBy = Login
		}
	}
	if request.NewKey != "" {
		if request.NewKey != request.Key {
			//delete old
			{
				delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
				if _, delErr := Mdb_Customer_Loyalty_Account.Coll.DeleteOne(delCtx, bson.M{"Key": request.Key}); delErr != nil {
					log.Println("Mdb_Customer_Loyalty_Account DeleteOne error:", delErr)
				}
				if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Customer_Loyalty_Account{Key: request.Key}.RedisKey()); delRedisErr != nil {
					log.Println("redisx.DelJSON Customer_Loyalty_Account error:", delRedisErr)
				}
				delCancel()
			}
			//update key
			entry.Key = request.NewKey
		}
	}

	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Customer_Loyalty_Account.Coll.UpdateOne(putCtx, bson.M{"Key": entry.Key}, bson.M{"$set": entry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Customer_Loyalty_Account upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, entry.RedisKey(), entry); putSetErr != nil {
			log.Println("redisx.SetJSON Customer_Loyalty_Account error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	if Login != "DWH_Import" { // --> off to avoid filling up logs
		Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
			Event_User:         Login,
			Event_Time:         time.Now(),
			Event_AffectedType: "Customer_Loyalty_Account",
			Event_ActionType:   "Edit",
			Event_Description:  "",
			Event_Entry_Before: Current_Entry,
			Event_Entry_After:  entry,
		})
	}

	return Id, nil
}

func (Uc *UserControl) Customer_Loyalty_Account_Get(Key string) (entries []Customer_Loyalty_Account, err error) {
	if Key == "" {
		entries, err = redisx.GetAllJSONByPattern[Customer_Loyalty_Account](context.Background(), RedisClient, redisx.ScanJSONOptions{
			Pattern: "Customer_Loyalty_Account:*", ScanCount: 500, PipelineSize: 250,
		})
		if err != nil {
			return entries, err
		}
		return entries, nil
	} else {
		Key = Normalize_International_MSISDN(Key)
		if Key == "" {
			err = errors.New("key cannot be empty")
			return entries, err
		}
		entry, entryErr := redisx.GetJSON[Customer_Loyalty_Account](context.Background(), RedisClient, Customer_Loyalty_Account{Key: Key}.RedisKey())
		if redisx.IsNil(entryErr) {
			err = errors.New("key does not exist")
			return entries, err
		}
		if entryErr != nil {
			return entries, entryErr
		}

		plan, planErr := redisx.GetJSON[Loyalty_Plan](context.Background(), RedisClient, Loyalty_Plan{Key: entry.Loyalty_Level_Key + "|" + entry.Loyalty_Account_Segment_Key}.RedisKey())
		if redisx.IsNil(planErr) {
			err = errors.New("loyalty plan does not exist")
			return entries, err
		}
		if planErr != nil {
			return entries, planErr
		}
		//validate earning rules
		if plan.Expiry_Rules_Key == "" {
			err = errors.New("expiry rules is not defined")
			return entries, err
		}
		expiry_Rule, expiryErr := redisx.GetJSON[Loyalty_Point_Expiry_Rules](context.Background(), RedisClient, Loyalty_Point_Expiry_Rules{Key: plan.Expiry_Rules_Key}.RedisKey())
		if redisx.IsNil(expiryErr) {
			err = errors.New("expiry rules is not defined")
			return entries, err
		}
		if expiryErr != nil {
			return entries, expiryErr
		}
		var initialDate time.Time
		if !entry.Expiry_Date.IsZero() {
			initialDate = entry.Expiry_Date
		} else if !entry.First_Opt_In_Status_Date.IsZero() {
			initialDate = entry.First_Opt_In_Status_Date
		} else {
			initialDate = entry.Creation_date
		}
		initialexpiryDate := addValidity(initialDate, expiry_Rule.Validity_Unit, expiry_Rule.Validity_Duration)
		finalexpiryDate := addValidity(initialexpiryDate, expiry_Rule.Grace_Validity_Unit, expiry_Rule.Grace_Validity_Duration)
		entry.Coming_Expiry_Date = finalexpiryDate
		entry.Initial_Date = initialexpiryDate
		var expiryPoints float64 = 0
		for _, pointKey := range entry.Points_Detail_Keys {
			pointsDetail, err := Uc.Customer_Loyalty_Account_Points_Details_Get(pointKey)
			if err != nil {
				return entries, err
			}
			year, err := strconv.Atoi(pointsDetail[0].Year_Month[:4])
			if err != nil {
				return entries, err
			}
			month, err := strconv.Atoi(pointsDetail[0].Year_Month[4:])
			if err == nil && year < int(entry.Initial_Date.Year()) || (year == entry.Initial_Date.Year() && month < int(entry.Initial_Date.Month())) {
				expiryPoints += pointsDetail[0].Available_Points
			}
			entry.Points_To_Expire = expiryPoints
		}
		if !entry.Joining_Date.IsZero() && entry.Joining_Date.Year() != 1 {
			now := time.Now()

			years := now.Year() - entry.Joining_Date.Year()
			months := int(now.Month()) - int(entry.Joining_Date.Month())

			totalMonths := years*12 + months

			// Adjust if current day is before the joining day in the month
			if now.Day() < entry.Joining_Date.Day() {
				totalMonths--
			}
			if totalMonths > 0 {
				loyalty_Level, loyaltyLvlErr := redisx.GetJSON[Loyalty_Level](context.Background(), RedisClient, Loyalty_Level{Key: entry.Loyalty_Level_Key}.RedisKey())
				if redisx.IsNil(loyaltyLvlErr) {
					err = errors.New("failed to get loyalty account level")
					return entries, err
				}
				if loyaltyLvlErr != nil {
					return entries, loyaltyLvlErr
				}
				for _, lvl := range loyalty_Level.Seniority_Levels {
					seniority, seniorityErr := redisx.GetJSON[Loyalty_Seniority_Level](context.Background(), RedisClient, Loyalty_Seniority_Level{Key: lvl.Loyalty_Seniority_Level_Key}.RedisKey())
					if seniorityErr == nil {
						if months >= int(seniority.AON_From) && months <= int(seniority.AON_Till) {
							entry.Multiplier_Percentage = lvl.Multiplier_Percentage
						}
					}

				}
			}

		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (Uc *UserControl) Customer_Loyalty_Account_GetPaginated(Page, Limit int64) (entries []Customer_Loyalty_Account, err error) {
	if Page < 1 {
		return entries, errors.New("invalid page")
	}
	if Limit < 1 || Limit > 50000 {
		return entries, errors.New("invalid limit (accept value between 1 and 50000)")
	}

	//var array []daoc.DAOFindCriteria
	// if Outlet_Key != "" {
	// 	//restrict access for records that belong to this user
	// 	var criteria daoc.DAOFindCriteria = daoc.DAOFindCriteria{
	// 		Field:    "Outlet_Key",
	// 		Value:    Outlet_Key,
	// 		Operator: "EQUAL",
	// 	}
	// 	array = append(array, criteria)
	// }
	// if Agent_Key != "" {
	// 	//restrict access for records that belong to this user
	// 	var criteria daoc.DAOFindCriteria = daoc.DAOFindCriteria{
	// 		Field:    "Agent_Key",
	// 		Value:    Agent_Key,
	// 		Operator: "EQUAL",
	// 	}
	// 	array = append(array, criteria)
	// }
	// if len(array) > 0 {
	// 	findparams.FindCriteria = array
	// }
	pgCtx, pgCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer pgCancel()
	result, err := mongox.FindPage[Customer_Loyalty_Account](pgCtx, Mdb_Customer_Loyalty_Account.Coll, bson.M{}, mongox.PageRequest{
		Page: Page, PageSize: Limit, MaxPageSize: 50000,
	})
	if err != nil {
		return entries, err
	}
	entries = result.Items
	return entries, nil
}

func (Uc *UserControl) Customer_Loyalty_Account_Delete(Login, Key string) (err error) {
	Key = Normalize_International_MSISDN(Key)
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry, entryErr := redisx.GetJSON[Customer_Loyalty_Account](context.Background(), RedisClient, Customer_Loyalty_Account{Key: Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("entry does not exist")
		return err
	}
	if entryErr != nil {
		return entryErr
	}
	{
		delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, delErr := Mdb_Customer_Loyalty_Account.Coll.DeleteOne(delCtx, bson.M{"Key": Key}); delErr != nil {
			log.Println("Mdb_Customer_Loyalty_Account DeleteOne error:", delErr)
		}
		if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Customer_Loyalty_Account{Key: Key}.RedisKey()); delRedisErr != nil {
			log.Println("redisx.DelJSON Customer_Loyalty_Account error:", delRedisErr)
		}
		delCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Customer_Loyalty_Account",
		Event_ActionType:   "Delete",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  entry,
	})
	return nil
}

func addValidity(date time.Time, unit string, duration int) time.Time {
	switch unit {
	case "Month":
		return date.AddDate(0, duration, 0)
	case "Year":
		return date.AddDate(duration, 0, 0)
	default:
		// Handle invalid unit
		fmt.Println("Invalid validity unit:", unit)
		return date
	}
}

func (Uc *UserControl) Customer_Loyalty_Account_Points_Details_Get(Key string) (entries []Customer_Loyalty_Account_Points_Detail, err error) {
	if Key == "" {
		entries, err = redisx.GetAllJSONByPattern[Customer_Loyalty_Account_Points_Detail](context.Background(), RedisClient, redisx.ScanJSONOptions{
			Pattern: "Customer_Loyalty_Account_Points_Detail:*", ScanCount: 500, PipelineSize: 250,
		})
		if err != nil {
			return entries, err
		}
		return entries, nil
	} else {
		entry, entryErr := redisx.GetJSON[Customer_Loyalty_Account_Points_Detail](context.Background(), RedisClient, Customer_Loyalty_Account_Points_Detail{Key: Key}.RedisKey())
		if redisx.IsNil(entryErr) {
			err = errors.New("key does not exist")
			return entries, err
		}
		if entryErr != nil {
			return entries, entryErr
		}
		entries = append(entries, entry)
		return entries, nil
	}
}

func (Uc *UserControl) Customer_Loyalty_Account_Points_Details_GetPaginated(Page, Limit int64) (entries []Customer_Loyalty_Account_Points_Detail, err error) {
	if Page < 1 {
		return entries, errors.New("invalid page")
	}
	if Limit < 1 || Limit > 50000 {
		return entries, errors.New("invalid limit (accept value between 1 and 50000)")
	}

	pgCtx, pgCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer pgCancel()
	result, err := mongox.FindPage[Customer_Loyalty_Account_Points_Detail](pgCtx, Mdb_Customer_Loyalty_Account_Points_Detail.Coll, bson.M{}, mongox.PageRequest{
		Page: Page, PageSize: Limit, MaxPageSize: 50000,
	})
	if err != nil {
		return entries, err
	}
	entries = result.Items
	return entries, nil
}
func (Uc *UserControl) Customer_Loyalty_Account_GetRedemption_Rules(MSISDN string) (Redemption_Rules Loyalty_Point_Redemption_Rule, err error) {
	MSISDN = Normalize_International_MSISDN(MSISDN)
	if MSISDN == "" {
		return Redemption_Rules, errors.New("msisdn cannot be empty")
	}
	//get loyalty account detail
	loyalty_account, loyaltyAccErr := redisx.GetJSON[Customer_Loyalty_Account](context.Background(), RedisClient, Customer_Loyalty_Account{Key: MSISDN}.RedisKey())
	if redisx.IsNil(loyaltyAccErr) {
		return Redemption_Rules, errors.New("loyalty account does not exist")
	}
	if loyaltyAccErr != nil {
		return Redemption_Rules, loyaltyAccErr
	}
	if loyalty_account.Loyalty_Account_Segment_Key == "" {
		return Redemption_Rules, errors.New("type assertion issue with Customer_Loyalty_Account")
	}
	if loyalty_account.Loyalty_Level_Key == "" {
		return Redemption_Rules, errors.New("loyalty account level is not assigned")
	}
	//get the loyalty plan
	plan, planErr := redisx.GetJSON[Loyalty_Plan](context.Background(), RedisClient, Loyalty_Plan{Key: loyalty_account.Loyalty_Level_Key + "|" + loyalty_account.Loyalty_Account_Segment_Key}.RedisKey())
	if redisx.IsNil(planErr) {
		return Redemption_Rules, errors.New("loyalty plan does not exist")
	}
	if planErr != nil {
		return Redemption_Rules, planErr
	}
	//validate earning rules
	if plan.Redemption_Rules_Key == "" {
		return Redemption_Rules, errors.New("redemption rules is not defined")
	}
	Redemption_Rule, redemptionErr := redisx.GetJSON[Loyalty_Point_Redemption_Rules](context.Background(), RedisClient, Loyalty_Point_Redemption_Rules{Key: plan.Redemption_Rules_Key}.RedisKey())
	if redisx.IsNil(redemptionErr) {
		return Redemption_Rules, errors.New("redemption rules is not defined")
	}
	if redemptionErr != nil {
		return Redemption_Rules, redemptionErr
	}
	var cleanRedemptionRule Loyalty_Point_Redemption_Rule
	err = copier.Copy(&cleanRedemptionRule, &Redemption_Rule)
	if err != nil {
		return Loyalty_Point_Redemption_Rule{}, err
	}
	return cleanRedemptionRule, nil
}

func (Uc *UserControl) Customer_Loyalty_Account_GetEarning_Rule(MSISDN string) (Earning_Rules Loyalty_Point_Earning_Rules, err error) {
	MSISDN = Normalize_International_MSISDN(MSISDN)
	if MSISDN == "" {
		return Earning_Rules, errors.New("msisdn cannot be empty")
	}
	//get loyalty account detail
	loyalty_account, loyaltyAccErr := redisx.GetJSON[Customer_Loyalty_Account](context.Background(), RedisClient, Customer_Loyalty_Account{Key: MSISDN}.RedisKey())
	if redisx.IsNil(loyaltyAccErr) {
		return Earning_Rules, errors.New("loyalty account does not exist")
	}
	if loyaltyAccErr != nil {
		return Earning_Rules, loyaltyAccErr
	}
	if loyalty_account.Loyalty_Account_Segment_Key == "" {
		return Earning_Rules, errors.New("type assertion issue with Customer_Loyalty_Account")
	}
	if loyalty_account.Loyalty_Level_Key == "" {
		return Earning_Rules, errors.New("loyalty account level is not assigned")
	}
	//get the loyalty plan
	plan, planErr := redisx.GetJSON[Loyalty_Plan](context.Background(), RedisClient, Loyalty_Plan{Key: loyalty_account.Loyalty_Level_Key + "|" + loyalty_account.Loyalty_Account_Segment_Key}.RedisKey())
	if redisx.IsNil(planErr) {
		return Earning_Rules, errors.New("loyalty plan does not exist")
	}
	if planErr != nil {
		return Earning_Rules, planErr
	}
	//validate earning rules
	if plan.Earning_Rules_Key == "" {
		return Earning_Rules, errors.New("earning rule is not defined")
	}
	Earning_Rule, earningErr := redisx.GetJSON[Loyalty_Point_Earning_Rules](context.Background(), RedisClient, Loyalty_Point_Earning_Rules{Key: plan.Earning_Rules_Key}.RedisKey())
	if redisx.IsNil(earningErr) {
		return Earning_Rules, errors.New("earning rule is not defined")
	}
	if earningErr != nil {
		return Earning_Rules, earningErr
	}

	return Earning_Rule, nil
}

func (Uc *UserControl) Customer_Loyalty_Account_GetRedemptionProductCatalogue(MSISDN string) (response PropC.Catalogue_WithBundleDetail_response, err error) {
	//get redemption plan
	redemption_Rules, err := Uc.Customer_Loyalty_Account_GetRedemption_Rules(MSISDN)
	if err != nil {
		return response, err
	}
	if redemption_Rules.Bundles_Product_Catalogue_Channel == "" {
		return response, errors.New("channel for product catalogue is not defined in redemption rules")
	}
	if redemption_Rules.Bundles_Product_Catalogue_Plan == "" {
		return response, errors.New("plan for product catalogue is not defined in redemption rules")
	}
	if redemption_Rules.Bundles_Product_Catalogue_Version == "" {
		return response, errors.New("version for product catalogue is not defined in redemption rules")
	}
	//Get catalogue
	response, err = Uc.Propylaea.PropylaeaClient.Get_Catalogue_WithBundleDetail_WithoutLocationRestriction(
		redemption_Rules.Bundles_Product_Catalogue_Channel,
		redemption_Rules.Bundles_Product_Catalogue_Plan,
		redemption_Rules.Bundles_Product_Catalogue_Version)
	if err != nil {
		return response, err
	}
	// for _, catalog := range bundles.Data {
	// 	for _, catalogEntry := range catalog.Catalogue_Entries {
	// 		for _, bundle := range catalogEntry.Bundles {
	// 			bundlesToReturn = append(bundlesToReturn, bundle)
	// 		}
	// 	}
	// }
	return response, nil
}
func FormatMBs1024(mb float64) string {
	if mb >= 1024*1024 {
		// 1 TB or more
		return fmt.Sprintf("%.2f TB", mb/1024/1024)
	} else if mb >= 1024 {
		// 1 GB or more
		return fmt.Sprintf("%.2f GB", mb/1024)
	} else {
		// Less than 1 GB
		return fmt.Sprintf("%.0f MB", mb)
	}
}
func FormatMBs1000(mb float64) string {
	if mb >= 1000*1000 {
		// 1 TB or more
		return fmt.Sprintf("%.2f TB", mb/1000/1000)
	} else if mb >= 1000 {
		// 1 GB or more
		return fmt.Sprintf("%.2f GB", mb/1000)
	} else {
		// Less than 1 GB
		return fmt.Sprintf("%.0f MB", mb)
	}
}
func (Uc *UserControl) Customer_Loyalty_RedeemRequest(request_header *Request_Header, request Loyalty_Redemption_Request, response *Loyalty_Redemption_log) {
	response.ReceiveDate = time.Now()
	//fill the request header info
	response.SourceIP = request_header.SourceIP
	response.SourceApp = request_header.SourceApp
	response.AppLogin = request_header.AppLogin
	response.AppVersion = request_header.AppVersion
	response.GPSLocation = request_header.GPSLocation
	response.GSMLocation = request_header.GSMLocation
	request.MSISDN = Normalize_International_MSISDN(request.MSISDN)
	if request.MSISDN == "" {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "incorrect msisdn"
		response.ErrorDescription = "incorrect msisdn"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_Redemption_log(*response)
		return
	}
	//fill the request info
	response.MSISDN = request.MSISDN
	response.ReceiveDate = time.Now()
	response.Redemption_Type = request.Redemption_Type //Airtime, Bundle, MobileMoney, SpinAndWin
	response.Redemption_Bunlde_Id = request.Redemption_Bunlde_Id
	response.Redemption_Amount = request.Redemption_Amount
	response.Points_To_Redeem = request.Points_To_Redeem

	//get loyalty account detail
	loyalty_Account, loyaltyAccErr := redisx.GetJSON[Customer_Loyalty_Account](context.Background(), RedisClient, Customer_Loyalty_Account{Key: request.MSISDN}.RedisKey())
	if redisx.IsNil(loyaltyAccErr) {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "loyalty account not found"
		response.ErrorDescription = "loyalty account not found"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_Redemption_log(*response)
		return
	}
	if loyaltyAccErr != nil {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "error in loyalty account type assertion"
		response.ErrorDescription = loyaltyAccErr.Error()
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_Redemption_log(*response)
		return
	}
	response.Customer_Id = loyalty_Account.Customer_Id
	response.Account_Status = loyalty_Account.Account_Status
	response.Loyalty_Level_Key = loyalty_Account.Loyalty_Level_Key
	response.Loyalty_Account_Segment_Key = loyalty_Account.Loyalty_Account_Segment_Key
	response.Opening_Awarded_Points = loyalty_Account.Awarded_Points
	response.Opening_Redeemed_Points = loyalty_Account.Redeemed_Points
	response.Opening_Available_Points = loyalty_Account.Available_Points

	//get redemption rules
	redemption_Rules, err := Uc.Customer_Loyalty_Account_GetRedemption_Rules(request.MSISDN)
	if err != nil {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "failed to get redemption rules"
		response.ErrorDescription = err.Error()
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_Redemption_log(*response)
		return
	}
	if response.Opening_Available_Points < redemption_Rules.Min_Accumulated_Points {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "no enough points"
		response.ErrorDescription = "no enough points"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_Redemption_log(*response)
		return
	}
	response.Allow_Negative_Balance_ToRedeem = redemption_Rules.Allow_Negative_Balance_ToRedeem
	//validate nagtive balance and pending lendme
	if !redemption_Rules.Allow_Negative_Balance_ToRedeem {
		IN_MSISDN := response.MSISDN
		if Configuration.Operation == "Gambia" {
			if len(response.MSISDN) > 7 {
				IN_MSISDN = IN_MSISDN[len(response.MSISDN)-7 : len(response.MSISDN)]
			}
		} else if Configuration.Operation == "SierraLeone" { //077928014
			if len(response.MSISDN) > 8 {
				IN_MSISDN = "0" + IN_MSISDN[len(response.MSISDN)-8:len(response.MSISDN)]
			}
		}
		IN_Response, err := Uc.CGW.UC_GWClient.IN_GetAccountDetails(response.SourceApp, request.MSISDN)
		// IN_Response, err := Uc.IN.INClient.GetAccountDetails("", "", IN_MSISDN)
		if err != nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to get IN account detail"
			response.ErrorDescription = err.Error()
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		if IN_Response.Data.Balance < 0 {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "balance must be positive"
			response.ErrorDescription = "main GSM balance must be positive"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
	}
	response.Allow_PendingLendme_ToRedeem = redemption_Rules.Allow_PendingLendme_ToRedeem
	if !redemption_Rules.Allow_PendingLendme_ToRedeem {
		subscriber, err := Uc.Lendme.LendmeClient.Lendme_Subscriber_Get(response.MSISDN)
		if err == nil && subscriber.Lendme_Outstanding_Amount > 0 {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "pending outstanding amount must be closed"
			response.ErrorDescription = "pending outstanding amount must be closed"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
	}

	//validate and execute redemption request
	switch request.Redemption_Type {
	case "Airtime":
		if request.Redemption_Amount <= 0 && request.Points_To_Redeem <= 0 {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "invalid redemption amount"
			response.ErrorDescription = "invalid redemption airtime amount"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		response.MinAvailableRequiredPoints = redemption_Rules.Available_MinPoints_for_Airtime
		if response.Opening_Available_Points < response.MinAvailableRequiredPoints {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "no enough points"
			response.ErrorDescription = "no enough points"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		if redemption_Rules.Bundles_Product_Catalogue_Channel == "" {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "redemption product catalogue channel is not provided"
			response.ErrorDescription = "redemption product catalogue channel is not provided"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		if redemption_Rules.Airtime_EVC_Account == "" {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "airtime payer account is not provided"
			response.ErrorDescription = "airtime payer account is not provided"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		//debit loyalty points
		if redemption_Rules.Airtime_Points <= 0 || redemption_Rules.Airtime_Amount <= 0 {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "airtime redeem rules are not defined"
			response.ErrorDescription = "airtime redeem rules are not defined"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		//do the calculation
		if response.Points_To_Redeem > 0 {
			//calculate Redemption_Amount
			response.Redemption_Amount = (response.Points_To_Redeem * redemption_Rules.Airtime_Amount) / redemption_Rules.Airtime_Points
		} else {
			if response.Redemption_Amount > 0 {
				//calculate Points_To_Redeem
				response.Points_To_Redeem = (response.Redemption_Amount * redemption_Rules.Airtime_Points) / redemption_Rules.Airtime_Amount
			} else {
				//return error
				response.Status = "failed"
				response.StatusCode = http.StatusBadRequest
				response.StatusDescription = "invalid redemption amount"
				response.ErrorDescription = "invalid redemption airtime amount"
				response.StatusDate = time.Now()
				response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
				Uc.Write_Loyalty_Redemption_log(*response)
				return
			}
		}
		response.MinRequiredPoints = redemption_Rules.Airtime_MinPoints
		if response.Points_To_Redeem < response.MinRequiredPoints {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "requested points are less than the minimum allowed points for redemption"
			response.ErrorDescription = "requested points are less than the minimum allowed points for redemption"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		//check if subscriber have enough points
		if response.Points_To_Redeem > response.Opening_Available_Points {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "no enough points"
			response.ErrorDescription = "no enough points"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}

		var Debit_Request Loyalty_AccountDebitPoints_Request
		Debit_Request.MSISDN = response.MSISDN
		Debit_Request.Debit_Amount = response.Points_To_Redeem
		Debit_Request.Debit_Reason = "Redeem Request"
		Debit_Request.Redemption_Type = "Airtime" //Airtime, Bundle, MobileMoney, SpinAndWin
		Debit_Request.Redemption_Bunlde_Id = ""
		Debit_Request.Redemption_Amount = response.Redemption_Amount
		var debit_Log Loyalty_AccountDebitPoints_log
		Uc.Loyalty_AccountDebitPoints(request_header, Debit_Request, &debit_Log)
		response.Points_Debit_Result = debit_Log
		if debit_Log.StatusCode != http.StatusOK {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = debit_Log.StatusDescription
			response.ErrorDescription = debit_Log.ErrorDescription
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		response.Closure_Awarded_Points = debit_Log.Closure_Awarded_Points
		response.Closure_Redeemed_Points = debit_Log.Closure_Redeemed_Points
		response.Closure_Available_Points = debit_Log.Closure_Available_Points
		//credit airtime
		Airtime_EVC_PIN, err := DecryptHexString(redemption_Rules.Airtime_EVC_PIN)
		if err != nil {
			fmt.Println("error in decrypting artime evc pin", err.Error())
		}
		airtimeTransferReply, err := Uc.CGW.UC_GWClient.AirtimePurchase(Unified_charging_gateway_Client.AirtimePurchase_Request{
			PayerMSISDN:            redemption_Rules.Airtime_EVC_Account,
			PayerPIN:               Airtime_EVC_PIN,
			PaymentMethod:          "Loyalty Points",
			TargetMSISDN:           request.MSISDN,
			Amount:                 response.Redemption_Amount,
			SendPayerNotification:  false,
			SendTargetNotification: true,
			Language:               "EN",
		}, redemption_Rules.Bundles_Product_Catalogue_Channel)
		response.Airtime_PurchaseResult = airtimeTransferReply
		if err != nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to recharge airtime"
			response.ErrorDescription = err.Error()
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			//refund points
			var refund_response Loyalty_AccountCreditPoints_log
			var refundRequest Loyalty_AccountCreditPoints_Request
			refundRequest.MSISDN = response.MSISDN
			refundRequest.EventSource = response.SourceApp
			refundRequest.EventType = "refund"
			refundRequest.EventDetail = request.Redemption_Type + " - refund"
			refundRequest.EventAmount = 0
			refundRequest.PointsToCredit = response.Points_To_Redeem
			refundRequest.EventDescription = ""
			Uc.Loyalty_AccountCreditPoints(request_header, refundRequest, &refund_response, true)
			return
		}
		if airtimeTransferReply.StatusCode != http.StatusOK {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = airtimeTransferReply.StatusDescription
			response.ErrorDescription = airtimeTransferReply.ErrorDescription
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			//refund points
			var refund_response Loyalty_AccountCreditPoints_log
			var refundRequest Loyalty_AccountCreditPoints_Request
			refundRequest.MSISDN = response.MSISDN
			refundRequest.EventSource = response.SourceApp
			refundRequest.EventType = "refund"
			refundRequest.EventDetail = request.Redemption_Type + " - refund"
			refundRequest.EventAmount = 0
			refundRequest.PointsToCredit = response.Points_To_Redeem
			refundRequest.EventDescription = ""
			Uc.Loyalty_AccountCreditPoints(request_header, refundRequest, &refund_response, true)
			return
		}
		AirtimeRedemptionCount.With(prometheus.Labels{"EventSource": response.SourceApp, "Level": loyalty_Account.Loyalty_Level_Key}).Inc()
		AirtimeRedemptionPoints.With(prometheus.Labels{"EventSource": response.SourceApp, "Level": loyalty_Account.Loyalty_Level_Key}).Add(response.Points_To_Redeem)
		AirtimeRedemptionAmount.With(prometheus.Labels{"EventSource": response.SourceApp, "Level": loyalty_Account.Loyalty_Level_Key}).Add(response.Redemption_Amount)

		record, err := Uc.Customer_Loyalty_Account_GetRedemption_Rules(response.MSISDN)
		if err != nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
			response.ErrorDescription = err.Error()
			return
		}
		if record.Airtime_Notification {
			AirtimeNotiLog := NotificationLog{
				SourceAction:  "AirtimeRedemption",
				TransactionId: "",
				Medium:        "SMS",
				SourceAddress: record.Airtime_Notification_Sender,
				Destination:   response.MSISDN,
				Subject:       "Redemption",
				AddUser:       "SYSTEM",
				AddDate:       time.Now(),
			}
			Airtime_Noti_Text := ""
			Airtime_Noti_Text = record.Airtime_Notification_Text

			if Airtime_Noti_Text != "" {
				Airtime_Noti_Text = strings.ReplaceAll(Airtime_Noti_Text, "{{ PointsDeducted }}", fmt.Sprint(response.Points_To_Redeem))
				Airtime_Noti_Text = strings.ReplaceAll(Airtime_Noti_Text, "{{ AirtimeAwarded }}", fmt.Sprint(response.Redemption_Amount))
				Airtime_Noti_Text = strings.ReplaceAll(Airtime_Noti_Text, "{{ LoyaltyBalance }}", fmt.Sprint(response.Closure_Available_Points))
				AirtimeNotiLog.Payload = Airtime_Noti_Text
				err := Send_SMS(record.Airtime_Notification_Sender, response.MSISDN, Airtime_Noti_Text)
				if err != nil {
					AirtimeNotiLog.Status = "Failed"
					AirtimeNotiLog.Error = err.Error()
				} else {
					AirtimeNotiLog.Status = "Successful"
				}
			} else {
				AirtimeNotiLog.Payload = Airtime_Noti_Text
				AirtimeNotiLog.Status = "Failed"
				AirtimeNotiLog.Error = "Undefined Airtime notification for transaction, check bundle definition"
			}
			YYYY, MM, _, _, _, _, _ := GetTimeParts(response.ReceiveDate)
			Db := Configuration.DB_Name_Loyalty + "_Logs_" + YYYY + MM
			{
				notiCtx, notiCancel := context.WithTimeout(context.Background(), 10*time.Second)
				_, notiErr := Uc.MongoClient.Mongo.Database(Db).Collection("Col_NotificationLog").InsertOne(notiCtx, AirtimeNotiLog)
				notiCancel()
				if notiErr != nil {
					log.Println("Error in Write Airtime Notification Logs:", notiErr, " (", AirtimeNotiLog, ")")
				}
			}
		}

	case "Bundle":
		if request.Redemption_Bunlde_Id == "" {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "bundle is not provided"
			response.ErrorDescription = "bundle is not provided"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		response.MinRequiredPoints = redemption_Rules.Bundles_MinPoints
		if response.Opening_Available_Points < response.MinRequiredPoints {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "no enough points"
			response.ErrorDescription = "no enough points"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		//get bundles detail
		bundle_response, err := Uc.Propylaea.PropylaeaClient.Get_Bundle(
			request.Redemption_Bunlde_Id,
			"", "", "")
		if err != nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to get bundle detail"
			response.ErrorDescription = err.Error()
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		if len(bundle_response.Data) < 1 {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "bundle not found"
			response.ErrorDescription = "bundle not found"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}

		//check if subscriber have enough points
		response.Price_Loyalty_Points = bundle_response.Data[0].Price_Loyalty_Points
		response.Points_To_Redeem = response.Price_Loyalty_Points
		//check if subscriber have enough points
		if response.Points_To_Redeem > response.Opening_Available_Points {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "no enough points"
			response.ErrorDescription = "no enough points"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		//debit the loyalty points
		var Debit_Request Loyalty_AccountDebitPoints_Request
		Debit_Request.MSISDN = response.MSISDN
		Debit_Request.Debit_Amount = response.Points_To_Redeem
		Debit_Request.Debit_Reason = "Redeem Request"
		Debit_Request.Redemption_Type = "Bundle" //Airtime, Bundle, MobileMoney, SpinAndWin
		Debit_Request.Redemption_Bunlde_Id = request.Redemption_Bunlde_Id
		Debit_Request.Redemption_Amount = response.Price_Loyalty_Points
		var debit_Log Loyalty_AccountDebitPoints_log
		Uc.Loyalty_AccountDebitPoints(request_header, Debit_Request, &debit_Log)
		response.Points_Debit_Result = debit_Log
		if debit_Log.StatusCode != http.StatusOK {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = debit_Log.StatusDescription
			response.ErrorDescription = debit_Log.ErrorDescription
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		response.Closure_Awarded_Points = debit_Log.Closure_Awarded_Points
		response.Closure_Redeemed_Points = debit_Log.Closure_Redeemed_Points
		response.Closure_Available_Points = debit_Log.Closure_Available_Points
		//recharge the bundle
		Bundles_EVC_PIN, err := DecryptHexString(redemption_Rules.Bundles_EVC_PIN)
		if err != nil {
			fmt.Println("error in decrypting artime evc pin", err.Error())
		}
		bundlePurchaseReply, err := Uc.CGW.UC_GWClient.BundlePurchase(Unified_charging_gateway_Client.BundlePurchase_Request{
			PayerMSISDN:            redemption_Rules.Bundles_EVC_Account,
			PayerPIN:               Bundles_EVC_PIN,
			PaymentMethod:          "Loyalty Points",
			TargetMSISDN:           request.MSISDN,
			BundleKey:              request.Redemption_Bunlde_Id,
			SendPayerNotification:  false,
			SendTargetNotification: true,
			Language:               "EN",
		}, redemption_Rules.Bundles_Product_Catalogue_Channel)
		response.Bundle_PurchaseResult = bundlePurchaseReply
		if err != nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to recharge bundle"
			response.ErrorDescription = err.Error()
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			//refund points
			var refund_response Loyalty_AccountCreditPoints_log
			var refundRequest Loyalty_AccountCreditPoints_Request
			refundRequest.MSISDN = response.MSISDN
			refundRequest.EventSource = response.SourceApp
			refundRequest.EventType = "refund"
			refundRequest.EventDetail = request.Redemption_Type + " - refund"
			refundRequest.EventAmount = 0
			refundRequest.PointsToCredit = response.Points_To_Redeem
			refundRequest.EventDescription = ""
			Uc.Loyalty_AccountCreditPoints(request_header, refundRequest, &refund_response, true)
			return
		}
		if bundlePurchaseReply.StatusCode != http.StatusOK {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = bundlePurchaseReply.StatusDescription
			response.ErrorDescription = bundlePurchaseReply.ErrorDescription
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			//refund points
			var refund_response Loyalty_AccountCreditPoints_log
			var refundRequest Loyalty_AccountCreditPoints_Request
			refundRequest.MSISDN = response.MSISDN
			refundRequest.EventSource = response.SourceApp
			refundRequest.EventType = "refund"
			refundRequest.EventDetail = request.Redemption_Type + " - refund"
			refundRequest.EventAmount = 0
			refundRequest.PointsToCredit = response.Points_To_Redeem
			refundRequest.EventDescription = ""
			Uc.Loyalty_AccountCreditPoints(request_header, refundRequest, &refund_response, true)
			return
		}
		BundleRedemptionCount.With(prometheus.Labels{"EventSource": response.SourceApp, "BunldeId": request.Redemption_Bunlde_Id, "Level": loyalty_Account.Loyalty_Level_Key}).Inc()
		BundleRedemptionPoints.With(prometheus.Labels{"EventSource": response.SourceApp, "BunldeId": request.Redemption_Bunlde_Id, "Level": loyalty_Account.Loyalty_Level_Key}).Add(response.Points_To_Redeem)

		record, err := Uc.Customer_Loyalty_Account_GetRedemption_Rules(response.MSISDN)
		if record.Bundles_Notification {
			if err != nil {
				response.Status = "failed"
				response.StatusCode = http.StatusBadRequest
				response.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
				response.ErrorDescription = err.Error()
				return
			}
			BundleNotiLog := NotificationLog{
				SourceAction:  "BundleRedemption",
				TransactionId: "",
				Medium:        "SMS",
				SourceAddress: record.Bundles_Notification_Sender,
				Destination:   response.MSISDN,
				Subject:       "Redemption",
				AddUser:       "SYSTEM",
				AddDate:       time.Now(),
			}
			Bundle_Noti_Text := ""
			var bundle propylaea.Bundle
			Bundle_Noti_Text = record.Bundles_Notification_Text
			bundleResponse, err := Uc.Propylaea.PropylaeaClient.Get_Bundle(response.Redemption_Bunlde_Id, "", "", "")
			if err != nil {
				fmt.Println("failed to get bundle")
			}
			var bundleName string
			if err == nil && len(bundleResponse.Data) > 0 {
				bundle = bundleResponse.Data[0]
				bundleName = bundle.Name_Lang1
			} else {
				bundleName = ""
			}
			// Variables to replace
			Minutes := 0.0
			MBs := 0
			SMS := 0
			Bonus_Minutes := 0.0
			Bonus_MBs := 0
			Bonus_SMS := 0

			Validity_value := bundle.Validity_Value
			Validity_type := ""
			// Minutes and SMSs from IN tokens
			seconds := 0.0
			bonus_seconds := 0.0
			for _, token := range bundle.IN_Tokens {
				if token.Token_Type == "SMS (whole SMS messages)" {
					if token.IsBonus {
						Bonus_SMS += int(token.Token_Value)
					} else {
						SMS += int(token.Token_Value)
					}
				} else if token.Token_Type == "time (seconds)" {
					if token.IsBonus {
						bonus_seconds += token.Token_Value
					} else {
						seconds += token.Token_Value
					}
				}
			}

			if seconds > 0 {
				Minutes = seconds / 60
			}
			if bonus_seconds > 0 {
				Bonus_Minutes = bonus_seconds / 60
			}
			// MBs from Alepo and Protei
			var DO_bytes float64 = 0
			var bonus_DO_bytes float64 = 0
			for _, DO_A := range bundle.Data_Offers {
				var size float64
				size, err := strconv.ParseFloat(DO_A.Offer_Size_Value, 64)
				if err != nil {
					fmt.Println("Error parsing Alepo Offer size value:", DO_A.Offer_Size_Value)
					continue
				}

				// Use size as float64 here
				fmt.Printf("Parsed size: %.2f\n", size) // e.g., show 2 decimal places

				if DO_A.IsBonus {
					if Configuration.Operation == "DRC" {
						if DO_A.Offer_Size_Unit == "MB" {
							size += (size * 1000 * 1000)
						} else if DO_A.Offer_Size_Unit == "GB" {
							size += (size * 1000 * 1000 * 1000)
						}
						bonus_DO_bytes += size
					} else {
						if DO_A.Offer_Size_Unit == "MB" {
							size += (size * 1024 * 1024)
						} else if DO_A.Offer_Size_Unit == "GB" {
							size += (size * 1024 * 1024 * 1024)
						}
						bonus_DO_bytes += size
					}

				} else {
					if Configuration.Operation == "DRC" {
						if DO_A.Offer_Size_Unit == "MB" {
							size += (size * 1000 * 1000)
						} else if DO_A.Offer_Size_Unit == "GB" {
							size += (size * 1000 * 1000 * 1000)
						}
						DO_bytes += size
					} else {
						if DO_A.Offer_Size_Unit == "MB" {
							size += (size * 1024 * 1024)
						} else if DO_A.Offer_Size_Unit == "GB" {
							size += (size * 1024 * 1024 * 1024)
						}
						DO_bytes += size
					}
				}
			}

			if DO_bytes > 0 {
				if Configuration.Operation == "DRC" {
					MBs = int(DO_bytes / 1000 / 1000)
				} else {
					MBs = int(DO_bytes / 1024 / 1024)
				}
			}
			if bonus_DO_bytes > 0 {
				if Configuration.Operation == "DRC" {
					Bonus_MBs = int(bonus_DO_bytes / 1000 / 1000)
				} else {
					Bonus_MBs = int(bonus_DO_bytes / 1024 / 1024)
				}
			}

			// Expiry Date from bundle Validity
			switch bundle.Validity_Type {
			case "Hourly":
				if bundle.Validity_Value == 1 {
					Validity_type = "hour"
				} else {
					Validity_type = "hours"
				}
			case "Daily":
				if bundle.Validity_Value == 1 {
					Validity_type = "day"
				} else {
					Validity_type = "days"
				}
			case "Weekly":
				if bundle.Validity_Value == 1 {
					Validity_type = "week"
				} else {
					Validity_type = "weeks"
				}
			case "Monthly":
				if bundle.Validity_Value == 1 {
					Validity_type = "month"
				} else {
					Validity_type = "months"
				}
			case "Yearly":
				if bundle.Validity_Value == 1 {
					Validity_type = "year"
				} else {
					Validity_type = "years"
				}
			}
			sizeParts := []string{}

			totalMBs := MBs + Bonus_MBs
			if totalMBs > 0 {
				formatted := ""
				if Configuration.Operation == "DRC" {
					formatted = FormatMBs1000(float64(totalMBs))

				} else {
					formatted = FormatMBs1024(float64(totalMBs))
				}
				sizeParts = append(sizeParts, formatted)
			}

			totalSMS := SMS + Bonus_SMS
			if totalSMS > 0 {
				sizeParts = append(sizeParts, fmt.Sprintf("%d SMS", totalSMS))
			}

			totalMinutes := Minutes + Bonus_Minutes

			if totalMinutes > 0 {
				sizeParts = append(sizeParts, fmt.Sprintf("%.1f Minutes", totalMinutes))
			}

			Size := strings.Join(sizeParts, " + ")

			Validity := fmt.Sprintf("%d %s", Validity_value, Validity_type)
			if Bundle_Noti_Text != "" {
				Bundle_Noti_Text = strings.ReplaceAll(Bundle_Noti_Text, "{{ PointsDeducted }}", fmt.Sprint(response.Points_To_Redeem))
				Bundle_Noti_Text = strings.ReplaceAll(Bundle_Noti_Text, "{{ BundleName }}", fmt.Sprint(bundleName))
				Bundle_Noti_Text = strings.ReplaceAll(Bundle_Noti_Text, "{{ BundleSize }}", fmt.Sprint(Size))
				Bundle_Noti_Text = strings.ReplaceAll(Bundle_Noti_Text, "{{ BundleValidity }}", fmt.Sprint(Validity))
				Bundle_Noti_Text = strings.ReplaceAll(Bundle_Noti_Text, "{{ LoyaltyBalance }}", fmt.Sprint(response.Closure_Available_Points))
				BundleNotiLog.Payload = Bundle_Noti_Text
				err := Send_SMS(record.Bundles_Notification_Sender, response.MSISDN, Bundle_Noti_Text)
				if err != nil {
					BundleNotiLog.Status = "Failed"
					BundleNotiLog.Error = err.Error()
				} else {
					BundleNotiLog.Status = "Successful"
				}
			} else {
				BundleNotiLog.Payload = Bundle_Noti_Text
				BundleNotiLog.Status = "Failed"
				BundleNotiLog.Error = "Undefined Bundle notification for transaction"
			}
			YYYY, MM, _, _, _, _, _ := GetTimeParts(response.ReceiveDate)
			Db := Configuration.DB_Name_Loyalty + "_Logs_" + YYYY + MM
			{
				notiCtx, notiCancel := context.WithTimeout(context.Background(), 10*time.Second)
				_, notiErr := Uc.MongoClient.Mongo.Database(Db).Collection("Col_NotificationLog").InsertOne(notiCtx, BundleNotiLog)
				notiCancel()
				if notiErr != nil {
					log.Println("Error in Write Bundle Notification Logs:", notiErr, " (", BundleNotiLog, ")")
				}
			}
		}
	case "MobileMoney":
		if request.Redemption_Amount <= 0 && request.Points_To_Redeem <= 0 {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "invalid redemption amount"
			response.ErrorDescription = "invalid redemption airtime amount"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		response.MinAvailableRequiredPoints = redemption_Rules.Available_MinPoints_for_MobileMoney
		if response.Opening_Available_Points < response.MinAvailableRequiredPoints {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "no enough points"
			response.ErrorDescription = "no enough points"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		if redemption_Rules.MobileMoney_MerchantAccount == "" {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "mobile money payer account is not provided"
			response.ErrorDescription = "mobile money payer account is not provided"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		//debit loyalty points
		if redemption_Rules.MobileMoney_Points <= 0 || redemption_Rules.MobileMoney_Amount <= 0 {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "mobile money redeem rules are not defined"
			response.ErrorDescription = "mobile money redeem rules are not defined"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		//do the calculation
		if response.Points_To_Redeem > 0 {
			//calculate Redemption_Amount
			response.Redemption_Amount = (response.Points_To_Redeem * redemption_Rules.MobileMoney_Amount) / redemption_Rules.MobileMoney_Points
		} else {
			if response.Redemption_Amount > 0 {
				//calculate Points_To_Redeem
				response.Points_To_Redeem = (response.Redemption_Amount * redemption_Rules.MobileMoney_Points) / redemption_Rules.MobileMoney_Amount
			} else {
				//return error
				response.Status = "failed"
				response.StatusCode = http.StatusBadRequest
				response.StatusDescription = "invalid redemption amount"
				response.ErrorDescription = "invalid redemption mobile money amount"
				response.StatusDate = time.Now()
				response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
				Uc.Write_Loyalty_Redemption_log(*response)
				return
			}
		}
		response.MinRequiredPoints = redemption_Rules.MobileMoney_MinPoints
		if response.Points_To_Redeem < response.MinRequiredPoints {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "requested points are less than the minimum allowed points for redemption"
			response.ErrorDescription = "requested points are less than the minimum allowed points for redemption"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		//check if subscriber have enough points
		if response.Points_To_Redeem > response.Opening_Available_Points {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "no enough points"
			response.ErrorDescription = "no enough points"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		var Debit_Request Loyalty_AccountDebitPoints_Request
		Debit_Request.MSISDN = response.MSISDN
		Debit_Request.Debit_Amount = response.Points_To_Redeem
		Debit_Request.Debit_Reason = "Redeem Request"
		Debit_Request.Redemption_Type = "MobileMoney" //Airtime, Bundle, MobileMoney, SpinAndWin
		Debit_Request.Redemption_Bunlde_Id = ""
		Debit_Request.Redemption_Amount = response.Redemption_Amount
		var debit_Log Loyalty_AccountDebitPoints_log
		Uc.Loyalty_AccountDebitPoints(request_header, Debit_Request, &debit_Log)
		response.Points_Debit_Result = debit_Log
		if debit_Log.StatusCode != http.StatusOK {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = debit_Log.StatusDescription
			response.ErrorDescription = debit_Log.ErrorDescription
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		response.Closure_Awarded_Points = debit_Log.Closure_Awarded_Points
		response.Closure_Redeemed_Points = debit_Log.Closure_Redeemed_Points
		response.Closure_Available_Points = debit_Log.Closure_Available_Points
		//to do: credit mobile money amount --> merchant transfer
		MobileMoney_MerchantPIN, err := DecryptHexString(redemption_Rules.MobileMoney_MerchantPIN)
		if err != nil {
			fmt.Println("error in decrypting artime evc pin", err.Error())
		}

		mm_CashIn_Reply, err := Uc.CGW.UC_GWClient.MM_Agent_CashIN(MM.CashIN_Request{
			SenderMSISDN:   redemption_Rules.MobileMoney_MerchantAccount,
			SenderPIN:      MobileMoney_MerchantPIN,
			ReceiverMSISDN: request.MSISDN,
			Amount:         fmt.Sprintf("%f", response.Redemption_Amount),
			Remark:         "Loyalty Redemption",
			// Currency:       "102",
		})
		response.MobileMoney_PurchaseResult = mm_CashIn_Reply
		if err != nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to redeem mobile money"
			response.ErrorDescription = err.Error()
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			//refund points
			var refund_response Loyalty_AccountCreditPoints_log
			var refundRequest Loyalty_AccountCreditPoints_Request
			refundRequest.MSISDN = response.MSISDN
			refundRequest.EventSource = response.SourceApp
			refundRequest.EventType = "refund"
			refundRequest.EventDetail = request.Redemption_Type + " - refund"
			refundRequest.EventAmount = 0
			refundRequest.PointsToCredit = response.Points_To_Redeem
			refundRequest.EventDescription = ""
			Uc.Loyalty_AccountCreditPoints(request_header, refundRequest, &refund_response, true)
			return
		}
		if mm_CashIn_Reply.Status != "SUCCEEDED" {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to redeem mobile money"
			if len(mm_CashIn_Reply.Errors) > 0 {
				response.ErrorDescription = mm_CashIn_Reply.Errors[0].Message
			}
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			//refund points
			var refund_response Loyalty_AccountCreditPoints_log
			var refundRequest Loyalty_AccountCreditPoints_Request
			refundRequest.MSISDN = response.MSISDN
			refundRequest.EventSource = response.SourceApp
			refundRequest.EventType = "refund"
			refundRequest.EventDetail = request.Redemption_Type + " - refund"
			refundRequest.EventAmount = 0
			refundRequest.PointsToCredit = response.Points_To_Redeem
			refundRequest.EventDescription = ""
			Uc.Loyalty_AccountCreditPoints(request_header, refundRequest, &refund_response, true)
			return
		}
		MobileMoneyRedemptionCount.With(prometheus.Labels{"EventSource": response.SourceApp, "Level": loyalty_Account.Loyalty_Level_Key}).Inc()
		MobileMoneyRedemptionPoints.With(prometheus.Labels{"EventSource": response.SourceApp, "Level": loyalty_Account.Loyalty_Level_Key}).Add(response.Points_To_Redeem)
		MobileMoneyRedemptionAmount.With(prometheus.Labels{"EventSource": response.SourceApp, "Level": loyalty_Account.Loyalty_Level_Key}).Add(response.Redemption_Amount)

		record, err := Uc.Customer_Loyalty_Account_GetRedemption_Rules(response.MSISDN)
		if err != nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
			response.ErrorDescription = err.Error()
			return
		}
		if record.MobileMoney_Notification {
			MobileMoneyNotiLog := NotificationLog{
				SourceAction:  "MobileMoneyRedemption",
				TransactionId: "",
				Medium:        "SMS",
				SourceAddress: record.MobileMoney_Notification_Sender,
				Destination:   response.MSISDN,
				Subject:       "Redemption",
				AddUser:       "SYSTEM",
				AddDate:       time.Now(),
			}
			MobileMoney_Noti_Text := ""
			MobileMoney_Noti_Text = record.MobileMoney_Notification_Text

			if MobileMoney_Noti_Text != "" {
				MobileMoney_Noti_Text = strings.ReplaceAll(MobileMoney_Noti_Text, "{{ PointsDeducted }}", fmt.Sprint(response.Points_To_Redeem))
				MobileMoney_Noti_Text = strings.ReplaceAll(MobileMoney_Noti_Text, "{{ MobileMoneyAwarded }}", fmt.Sprint(response.Redemption_Amount))
				MobileMoney_Noti_Text = strings.ReplaceAll(MobileMoney_Noti_Text, "{{ LoyaltyBalance }}", fmt.Sprint(response.Closure_Available_Points))
				MobileMoneyNotiLog.Payload = MobileMoney_Noti_Text
				err := Send_SMS(record.MobileMoney_Notification_Sender, response.MSISDN, MobileMoney_Noti_Text)
				if err != nil {
					MobileMoneyNotiLog.Status = "Failed"
					MobileMoneyNotiLog.Error = err.Error()
				} else {
					MobileMoneyNotiLog.Status = "Successful"
				}
			} else {
				MobileMoneyNotiLog.Payload = MobileMoney_Noti_Text
				MobileMoneyNotiLog.Status = "Failed"
				MobileMoneyNotiLog.Error = "Undefined Mobile Money notification for transaction"
			}
			YYYY, MM, _, _, _, _, _ := GetTimeParts(response.ReceiveDate)
			Db := Configuration.DB_Name_Loyalty + "_Logs_" + YYYY + MM
			{
				notiCtx, notiCancel := context.WithTimeout(context.Background(), 10*time.Second)
				_, notiErr := Uc.MongoClient.Mongo.Database(Db).Collection("Col_NotificationLog").InsertOne(notiCtx, MobileMoneyNotiLog)
				notiCancel()
				if notiErr != nil {
					log.Println("Error in Write Mobile Money Notification Logs:", notiErr, " (", MobileMoneyNotiLog, ")")
				}
			}
		}
	case "SpinAndWin":
		if request.Redemption_Amount <= 0 {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "invalid redemption amount"
			response.ErrorDescription = "invalid redemption spin and win amount"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		response.MinAvailableRequiredPoints = redemption_Rules.Available_MinPoints_for_SpinAndWin
		if response.Opening_Available_Points < response.MinAvailableRequiredPoints {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "no enough points"
			response.ErrorDescription = "no enough points"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		//debit loyalty points
		if redemption_Rules.FreeSpinAndWin_PointsPerSpin <= 0 {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "airtime redeem rules are not defined"
			response.ErrorDescription = "airtime redeem rules are not defined"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		//do the calculation
		response.Points_To_Redeem = response.Redemption_Amount * redemption_Rules.FreeSpinAndWin_PointsPerSpin
		response.MinRequiredPoints = redemption_Rules.FreeSpinAndWin_MinPoints
		if response.Points_To_Redeem < response.MinRequiredPoints {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "requested points are less than the minimum allowed points for redemption"
			response.ErrorDescription = "requested points are less than the minimum allowed points for redemption"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		//check if subscriber have enough points
		if response.Points_To_Redeem > response.Opening_Available_Points {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "no enough points"
			response.ErrorDescription = "no enough points"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}

		var Debit_Request Loyalty_AccountDebitPoints_Request
		Debit_Request.MSISDN = response.MSISDN
		Debit_Request.Debit_Amount = response.Points_To_Redeem
		Debit_Request.Debit_Reason = "Redeem Request"
		Debit_Request.Redemption_Type = "SpinAndWin" //Airtime, Bundle, MobileMoney, SpinAndWin
		Debit_Request.Redemption_Bunlde_Id = ""
		Debit_Request.Redemption_Amount = response.Redemption_Amount
		var debit_Log Loyalty_AccountDebitPoints_log
		Uc.Loyalty_AccountDebitPoints(request_header, Debit_Request, &debit_Log)
		response.Points_Debit_Result = debit_Log
		if debit_Log.StatusCode != http.StatusOK {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = debit_Log.StatusDescription
			response.ErrorDescription = debit_Log.ErrorDescription
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		response.Closure_Awarded_Points = debit_Log.Closure_Awarded_Points
		response.Closure_Redeemed_Points = debit_Log.Closure_Redeemed_Points
		response.Closure_Available_Points = debit_Log.Closure_Available_Points
		//awar spins chances
		spinAndWin_Reply, err := Uc.SpinAndWin.SpinAndWinClient.EligibleSubs_AddChances(SpinAndWin.EligibleSubs_AddChances_Request{
			Key:            request.MSISDN,
			SpinCountToAdd: int64(response.Redemption_Amount),
		})
		response.SpinAndWin_PurchaseResult = spinAndWin_Reply
		if err != nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to redeem Spin And Win Chances"
			response.ErrorDescription = err.Error()
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			//refund points
			var refund_response Loyalty_AccountCreditPoints_log
			var refundRequest Loyalty_AccountCreditPoints_Request
			refundRequest.MSISDN = response.MSISDN
			refundRequest.EventSource = response.SourceApp
			refundRequest.EventType = "refund"
			refundRequest.EventDetail = request.Redemption_Type + " - refund"
			refundRequest.EventAmount = 0
			refundRequest.PointsToCredit = response.Points_To_Redeem
			refundRequest.EventDescription = ""
			Uc.Loyalty_AccountCreditPoints(request_header, refundRequest, &refund_response, true)
			return
		}
		if spinAndWin_Reply.StatusCode != http.StatusOK {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = spinAndWin_Reply.StatusDescription
			response.ErrorDescription = spinAndWin_Reply.ErrorDescription
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			//refund points
			var refund_response Loyalty_AccountCreditPoints_log
			var refundRequest Loyalty_AccountCreditPoints_Request
			refundRequest.MSISDN = response.MSISDN
			refundRequest.EventSource = response.SourceApp
			refundRequest.EventType = "refund"
			refundRequest.EventDetail = request.Redemption_Type + " - refund"
			refundRequest.EventAmount = 0
			refundRequest.PointsToCredit = response.Points_To_Redeem
			refundRequest.EventDescription = ""
			Uc.Loyalty_AccountCreditPoints(request_header, refundRequest, &refund_response, true)
			return
		}

		record, err := Uc.Customer_Loyalty_Account_GetRedemption_Rules(response.MSISDN)
		if err != nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
			response.ErrorDescription = err.Error()
			return
		}
		if record.FreeSpinAndWin_Notification {
			SpinWinNotiLog := NotificationLog{
				SourceAction:  "SpinAndWinRedemption",
				TransactionId: "",
				Medium:        "SMS",
				SourceAddress: record.FreeSpinAndWin_Notification_Sender,
				Destination:   response.MSISDN,
				Subject:       "Redemption",
				AddUser:       "SYSTEM",
				AddDate:       time.Now(),
			}
			SpinWin_Noti_Text := ""
			SpinWin_Noti_Text = record.FreeSpinAndWin_Notification_Text

			if SpinWin_Noti_Text != "" {
				SpinWin_Noti_Text = strings.ReplaceAll(SpinWin_Noti_Text, "{{ PointsDeducted }}", fmt.Sprint(response.Points_To_Redeem))
				SpinWin_Noti_Text = strings.ReplaceAll(SpinWin_Noti_Text, "{{ SpinsAwarded }}", fmt.Sprint(response.Redemption_Amount))
				SpinWin_Noti_Text = strings.ReplaceAll(SpinWin_Noti_Text, "{{ LoyaltyBalance }}", fmt.Sprint(response.Closure_Available_Points))
				SpinWinNotiLog.Payload = SpinWin_Noti_Text
				err := Send_SMS(record.FreeSpinAndWin_Notification_Sender, response.MSISDN, SpinWin_Noti_Text)
				if err != nil {
					SpinWinNotiLog.Status = "Failed"
					SpinWinNotiLog.Error = err.Error()
				} else {
					SpinWinNotiLog.Status = "Successful"
				}
			} else {
				SpinWinNotiLog.Payload = SpinWin_Noti_Text
				SpinWinNotiLog.Status = "Failed"
				SpinWinNotiLog.Error = "Undefined Mobile Money notification for transaction"
			}
			YYYY, MM, _, _, _, _, _ := GetTimeParts(response.ReceiveDate)
			Db := Configuration.DB_Name_Loyalty + "_Logs_" + YYYY + MM
			{
				notiCtx, notiCancel := context.WithTimeout(context.Background(), 10*time.Second)
				_, notiErr := Uc.MongoClient.Mongo.Database(Db).Collection("Col_NotificationLog").InsertOne(notiCtx, SpinWinNotiLog)
				notiCancel()
				if notiErr != nil {
					log.Println("Error in Write Mobile Money Notification Logs:", notiErr, " (", SpinWinNotiLog, ")")
				}
			}
		}

	default:
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "invalid redemption type"
		response.ErrorDescription = "invalid redemption type (accepted values: Airtime, Bundle, MobileMoney, SpinAndWin)"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_Redemption_log(*response)
		return
	}
	//successful reply
	response.Status = "successful"
	response.StatusCode = http.StatusOK
	response.StatusDescription = ""
	response.ErrorDescription = ""
	response.StatusDate = time.Now()
	response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
	Uc.Write_Loyalty_Redemption_log(*response)
}
func (Uc *UserControl) Customer_Loyalty_RedeemRequest_Angola(request_header *Request_Header, request Loyalty_Redemption_Request, response *Loyalty_Redemption_log) {
	response.ReceiveDate = time.Now()
	//fill the request header info
	response.SourceIP = request_header.SourceIP
	response.SourceApp = request_header.SourceApp
	response.AppLogin = request_header.AppLogin
	response.AppVersion = request_header.AppVersion
	response.GPSLocation = request_header.GPSLocation
	response.GSMLocation = request_header.GSMLocation
	request.MSISDN = Normalize_International_MSISDN(request.MSISDN)
	if request.MSISDN == "" {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "incorrect msisdn"
		response.ErrorDescription = "incorrect msisdn"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_Redemption_log(*response)
		return
	}
	//fill the request info
	response.MSISDN = request.MSISDN
	response.ReceiveDate = time.Now()
	response.Redemption_Type = request.Redemption_Type //Airtime, Bundle, MobileMoney, SpinAndWin
	response.Redemption_Bunlde_Id = request.Redemption_Bunlde_Id
	response.Redemption_Amount = request.Redemption_Amount
	response.Points_To_Redeem = request.Points_To_Redeem

	//get loyalty account detail
	loyalty_Account, loyaltyAccErr := redisx.GetJSON[Customer_Loyalty_Account](context.Background(), RedisClient, Customer_Loyalty_Account{Key: request.MSISDN}.RedisKey())
	if redisx.IsNil(loyaltyAccErr) {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "loyalty account not found"
		response.ErrorDescription = "loyalty account not found"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_Redemption_log(*response)
		return
	}
	if loyaltyAccErr != nil {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "error in loyalty account type assertion"
		response.ErrorDescription = loyaltyAccErr.Error()
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_Redemption_log(*response)
		return
	}
	response.Customer_Id = loyalty_Account.Customer_Id
	response.Account_Status = loyalty_Account.Account_Status
	response.Loyalty_Level_Key = loyalty_Account.Loyalty_Level_Key
	response.Loyalty_Account_Segment_Key = loyalty_Account.Loyalty_Account_Segment_Key
	response.Opening_Awarded_Points = loyalty_Account.Awarded_Points
	response.Opening_Redeemed_Points = loyalty_Account.Redeemed_Points
	response.Opening_Available_Points = loyalty_Account.Available_Points

	//get redemption rules
	redemption_Rules, err := Uc.Customer_Loyalty_Account_GetRedemption_Rules(request.MSISDN)
	if err != nil {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "failed to get redemption rules"
		response.ErrorDescription = err.Error()
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_Redemption_log(*response)
		return
	}
	if response.Opening_Available_Points < redemption_Rules.Min_Accumulated_Points {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "no enough points"
		response.ErrorDescription = "no enough points"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_Redemption_log(*response)
		return
	}
	response.Allow_Negative_Balance_ToRedeem = redemption_Rules.Allow_Negative_Balance_ToRedeem
	//validate nagtive balance and pending lendme
	if !redemption_Rules.Allow_Negative_Balance_ToRedeem {
		IN_MSISDN := response.MSISDN
		// IN_Response, err := Uc.IN.INClient.GetAccountDetails("", "", IN_MSISDN)
		IN_Response, err := Uc.APGW.APGWClient.CS_AccountDetails(IN_MSISDN)
		fmt.Println("IN_Response", IN_Response)
		fmt.Println("err", err)
		success := IN_Response.Data.Response.Success == "true"
		message := IN_Response.Data.Response.Result.Message
		if err != nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to get IN account detail"
			response.ErrorDescription = message
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		if !success {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to get IN account detail" + message
			response.ErrorDescription = message
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
	}
	response.Allow_PendingLendme_ToRedeem = redemption_Rules.Allow_PendingLendme_ToRedeem
	if !redemption_Rules.Allow_PendingLendme_ToRedeem {
		subscriber, err := Uc.Lendme.LendmeClient.Lendme_Subscriber_Get(response.MSISDN)
		if err == nil && subscriber.Lendme_Outstanding_Amount > 0 {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "pending outstanding amount must be closed"
			response.ErrorDescription = "pending outstanding amount must be closed"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
	}

	//validate and execute redemption request
	switch request.Redemption_Type {
	case "Airtime":
		if request.Redemption_Amount <= 0 && request.Points_To_Redeem <= 0 {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "invalid redemption amount"
			response.ErrorDescription = "invalid redemption airtime amount"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		response.MinAvailableRequiredPoints = redemption_Rules.Available_MinPoints_for_Airtime
		if response.Opening_Available_Points < response.MinAvailableRequiredPoints {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "no enough points"
			response.ErrorDescription = "no enough points"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		if redemption_Rules.Bundles_Product_Catalogue_Channel == "" {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "redemption product catalogue channel is not provided"
			response.ErrorDescription = "redemption product catalogue channel is not provided"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		if redemption_Rules.Airtime_EVC_Account == "" {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "airtime payer account is not provided"
			response.ErrorDescription = "airtime payer account is not provided"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		//debit loyalty points
		if redemption_Rules.Airtime_Points <= 0 || redemption_Rules.Airtime_Amount <= 0 {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "airtime redeem rules are not defined"
			response.ErrorDescription = "airtime redeem rules are not defined"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		//do the calculation
		if response.Points_To_Redeem > 0 {
			//calculate Redemption_Amount
			response.Redemption_Amount = (response.Points_To_Redeem * redemption_Rules.Airtime_Amount) / redemption_Rules.Airtime_Points
		} else {
			if response.Redemption_Amount > 0 {
				//calculate Points_To_Redeem
				response.Points_To_Redeem = (response.Redemption_Amount * redemption_Rules.Airtime_Points) / redemption_Rules.Airtime_Amount
			} else {
				//return error
				response.Status = "failed"
				response.StatusCode = http.StatusBadRequest
				response.StatusDescription = "invalid redemption amount"
				response.ErrorDescription = "invalid redemption airtime amount"
				response.StatusDate = time.Now()
				response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
				Uc.Write_Loyalty_Redemption_log(*response)
				return
			}
		}
		response.MinRequiredPoints = redemption_Rules.Airtime_MinPoints
		if response.Points_To_Redeem < response.MinRequiredPoints {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "requested points are less than the minimum allowed points for redemption"
			response.ErrorDescription = "requested points are less than the minimum allowed points for redemption"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		//check if subscriber have enough points
		if response.Points_To_Redeem > response.Opening_Available_Points {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "no enough points"
			response.ErrorDescription = "no enough points"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}

		var Debit_Request Loyalty_AccountDebitPoints_Request
		Debit_Request.MSISDN = response.MSISDN
		Debit_Request.Debit_Amount = response.Points_To_Redeem
		Debit_Request.Debit_Reason = "Redeem Request"
		Debit_Request.Redemption_Type = "Airtime" //Airtime, Bundle, MobileMoney, SpinAndWin
		Debit_Request.Redemption_Bunlde_Id = ""
		Debit_Request.Redemption_Amount = response.Redemption_Amount
		var debit_Log Loyalty_AccountDebitPoints_log
		Uc.Loyalty_AccountDebitPoints(request_header, Debit_Request, &debit_Log)
		response.Points_Debit_Result = debit_Log
		if debit_Log.StatusCode != http.StatusOK {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = debit_Log.StatusDescription
			response.ErrorDescription = debit_Log.ErrorDescription
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		response.Closure_Awarded_Points = debit_Log.Closure_Awarded_Points
		response.Closure_Redeemed_Points = debit_Log.Closure_Redeemed_Points
		response.Closure_Available_Points = debit_Log.Closure_Available_Points
		//credit airtime
		Airtime_EVC_PIN, err := DecryptHexString(redemption_Rules.Airtime_EVC_PIN)
		if err != nil {
			fmt.Println("error in decrypting artime evc pin", err.Error())
		}
		// airtimeTransferReply, err := Uc.CGW.UC_GWClient.AirtimePurchase(Unified_charging_gateway_Client.AirtimePurchase_Request{
		// 	PayerMSISDN:            redemption_Rules.Airtime_EVC_Account,
		// 	PayerPIN:               Airtime_EVC_PIN,
		// 	PaymentMethod:          "Loyalty Points",
		// 	TargetMSISDN:           request.MSISDN,
		// 	Amount:                 response.Redemption_Amount,
		// 	SendPayerNotification:  false,
		// 	SendTargetNotification: true,
		// 	Language:               "EN",
		// }, redemption_Rules.Bundles_Product_Catalogue_Channel)

		airtimeTransferReply, err := Uc.APGW.APGWClient.D2C_AirtimePurchase(apgw.AirtimePurchase_Request{
			TargetMSISDN:           request.MSISDN,
			EVCAccount:             redemption_Rules.Airtime_EVC_Account,
			EVCPIN:                 Airtime_EVC_PIN,
			Amount:                 response.Redemption_Amount,
			SendPayerNotification:  false,
			SendTargetNotification: true,
		}, "Loyalty")
		response.Airtime_PurchaseResult = airtimeTransferReply
		if err != nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to recharge airtime"
			response.ErrorDescription = err.Error()
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			//refund points
			var refund_response Loyalty_AccountCreditPoints_log
			var refundRequest Loyalty_AccountCreditPoints_Request
			refundRequest.MSISDN = response.MSISDN
			refundRequest.EventSource = response.SourceApp
			refundRequest.EventType = "refund"
			refundRequest.EventDetail = request.Redemption_Type + " - refund"
			refundRequest.EventAmount = 0
			refundRequest.PointsToCredit = response.Points_To_Redeem
			refundRequest.EventDescription = ""
			Uc.Loyalty_AccountCreditPoints(request_header, refundRequest, &refund_response, true)
			return
		}
		if airtimeTransferReply.StatusCode != http.StatusOK {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = airtimeTransferReply.AirtimeAllocate_StatusDescription
			response.ErrorDescription = airtimeTransferReply.ErrorCode
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			//refund points
			var refund_response Loyalty_AccountCreditPoints_log
			var refundRequest Loyalty_AccountCreditPoints_Request
			refundRequest.MSISDN = response.MSISDN
			refundRequest.EventSource = response.SourceApp
			refundRequest.EventType = "refund"
			refundRequest.EventDetail = request.Redemption_Type + " - refund"
			refundRequest.EventAmount = 0
			refundRequest.PointsToCredit = response.Points_To_Redeem
			refundRequest.EventDescription = ""
			Uc.Loyalty_AccountCreditPoints(request_header, refundRequest, &refund_response, true)
			return
		}
		AirtimeRedemptionCount.With(prometheus.Labels{"EventSource": response.SourceApp, "Level": loyalty_Account.Loyalty_Level_Key}).Inc()
		AirtimeRedemptionPoints.With(prometheus.Labels{"EventSource": response.SourceApp, "Level": loyalty_Account.Loyalty_Level_Key}).Add(response.Points_To_Redeem)
		AirtimeRedemptionAmount.With(prometheus.Labels{"EventSource": response.SourceApp, "Level": loyalty_Account.Loyalty_Level_Key}).Add(response.Redemption_Amount)

		record, err := Uc.Customer_Loyalty_Account_GetRedemption_Rules(response.MSISDN)
		if err != nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
			response.ErrorDescription = err.Error()
			return
		}
		if record.Airtime_Notification {
			AirtimeNotiLog := NotificationLog{
				SourceAction:  "AirtimeRedemption",
				TransactionId: "",
				Medium:        "SMS",
				SourceAddress: record.Airtime_Notification_Sender,
				Destination:   response.MSISDN,
				Subject:       "Redemption",
				AddUser:       "SYSTEM",
				AddDate:       time.Now(),
			}
			Airtime_Noti_Text := ""
			Airtime_Noti_Text = record.Airtime_Notification_Text

			if Airtime_Noti_Text != "" {
				Airtime_Noti_Text = strings.ReplaceAll(Airtime_Noti_Text, "{{ PointsDeducted }}", fmt.Sprint(response.Points_To_Redeem))
				Airtime_Noti_Text = strings.ReplaceAll(Airtime_Noti_Text, "{{ AirtimeAwarded }}", fmt.Sprint(response.Redemption_Amount))
				Airtime_Noti_Text = strings.ReplaceAll(Airtime_Noti_Text, "{{ LoyaltyBalance }}", fmt.Sprint(response.Closure_Available_Points))
				AirtimeNotiLog.Payload = Airtime_Noti_Text
				err := SendSMS(record.Airtime_Notification_Sender, response.MSISDN, Airtime_Noti_Text)
				if err != nil {
					AirtimeNotiLog.Status = "Failed"
					AirtimeNotiLog.Error = err.Error()
				} else {
					AirtimeNotiLog.Status = "Successful"
				}
			} else {
				AirtimeNotiLog.Payload = Airtime_Noti_Text
				AirtimeNotiLog.Status = "Failed"
				AirtimeNotiLog.Error = "Undefined Airtime notification for transaction, check bundle definition"
			}
			YYYY, MM, _, _, _, _, _ := GetTimeParts(response.ReceiveDate)
			Db := Configuration.DB_Name_Loyalty + "_Logs_" + YYYY + MM
			{
				notiCtx, notiCancel := context.WithTimeout(context.Background(), 10*time.Second)
				_, notiErr := Uc.MongoClient.Mongo.Database(Db).Collection("Col_NotificationLog").InsertOne(notiCtx, AirtimeNotiLog)
				notiCancel()
				if notiErr != nil {
					log.Println("Error in Write Airtime Notification Logs:", notiErr, " (", AirtimeNotiLog, ")")
				}
			}
		}

	case "Bundle":
		if request.Redemption_Bunlde_Id == "" {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "bundle is not provided"
			response.ErrorDescription = "bundle is not provided"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		response.MinRequiredPoints = redemption_Rules.Bundles_MinPoints
		if response.Opening_Available_Points < response.MinRequiredPoints {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "no enough points"
			response.ErrorDescription = "no enough points"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		//get bundles detail
		bundle_response, err := Uc.Propylaea.PropylaeaClient.Get_Bundle(
			request.Redemption_Bunlde_Id,
			"", "", "")
		if err != nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to get bundle detail"
			response.ErrorDescription = err.Error()
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		if len(bundle_response.Data) < 1 {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "bundle not found"
			response.ErrorDescription = "bundle not found"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}

		//check if subscriber have enough points
		response.Price_Loyalty_Points = bundle_response.Data[0].Price_Loyalty_Points
		response.Points_To_Redeem = response.Price_Loyalty_Points
		//check if subscriber have enough points
		if response.Points_To_Redeem > response.Opening_Available_Points {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "no enough points"
			response.ErrorDescription = "no enough points"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		//debit the loyalty points
		var Debit_Request Loyalty_AccountDebitPoints_Request
		Debit_Request.MSISDN = response.MSISDN
		Debit_Request.Debit_Amount = response.Points_To_Redeem
		Debit_Request.Debit_Reason = "Redeem Request"
		Debit_Request.Redemption_Type = "Bundle" //Airtime, Bundle, MobileMoney, SpinAndWin
		Debit_Request.Redemption_Bunlde_Id = request.Redemption_Bunlde_Id
		Debit_Request.Redemption_Amount = response.Price_Loyalty_Points
		var debit_Log Loyalty_AccountDebitPoints_log
		Uc.Loyalty_AccountDebitPoints(request_header, Debit_Request, &debit_Log)
		response.Points_Debit_Result = debit_Log
		if debit_Log.StatusCode != http.StatusOK {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = debit_Log.StatusDescription
			response.ErrorDescription = debit_Log.ErrorDescription
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		response.Closure_Awarded_Points = debit_Log.Closure_Awarded_Points
		response.Closure_Redeemed_Points = debit_Log.Closure_Redeemed_Points
		response.Closure_Available_Points = debit_Log.Closure_Available_Points
		//recharge the bundle
		Bundles_EVC_PIN, err := DecryptHexString(redemption_Rules.Bundles_EVC_PIN)
		if err != nil {
			fmt.Println("error in decrypting artime evc pin", err.Error())
		}

		bundlePurchaseReply, err := Uc.APGW.APGWClient.BundlePurchase(apgw.BundlePurchase_Request{
			PayerMSISDN:            redemption_Rules.Bundles_EVC_Account,
			PaymentMethod:          "Loyalty Points",
			PayerPIN:               Bundles_EVC_PIN,
			TargetMSISDN:           request.MSISDN,
			EVCAccount:             redemption_Rules.Bundles_EVC_Account,
			EVCPIN:                 Bundles_EVC_PIN,
			BundleKey:              request.Redemption_Bunlde_Id,
			SendPayerNotification:  false,
			SendTargetNotification: true,
		})
		response.Bundle_PurchaseResult = bundlePurchaseReply
		if err != nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to recharge bundle"
			response.ErrorDescription = err.Error()
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			//refund points
			var refund_response Loyalty_AccountCreditPoints_log
			var refundRequest Loyalty_AccountCreditPoints_Request
			refundRequest.MSISDN = response.MSISDN
			refundRequest.EventSource = response.SourceApp
			refundRequest.EventType = "refund"
			refundRequest.EventDetail = request.Redemption_Type + " - refund"
			refundRequest.EventAmount = 0
			refundRequest.PointsToCredit = response.Points_To_Redeem
			refundRequest.EventDescription = ""
			Uc.Loyalty_AccountCreditPoints(request_header, refundRequest, &refund_response, true)
			return
		}
		if bundlePurchaseReply.StatusCode != http.StatusOK {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = bundlePurchaseReply.ErrorTechnicalReason
			response.ErrorDescription = bundlePurchaseReply.ErrorCode
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			//refund points
			var refund_response Loyalty_AccountCreditPoints_log
			var refundRequest Loyalty_AccountCreditPoints_Request
			refundRequest.MSISDN = response.MSISDN
			refundRequest.EventSource = response.SourceApp
			refundRequest.EventType = "refund"
			refundRequest.EventDetail = request.Redemption_Type + " - refund"
			refundRequest.EventAmount = 0
			refundRequest.PointsToCredit = response.Points_To_Redeem
			refundRequest.EventDescription = ""
			Uc.Loyalty_AccountCreditPoints(request_header, refundRequest, &refund_response, true)
			return
		}
		BundleRedemptionCount.With(prometheus.Labels{"EventSource": response.SourceApp, "BunldeId": request.Redemption_Bunlde_Id, "Level": loyalty_Account.Loyalty_Level_Key}).Inc()
		BundleRedemptionPoints.With(prometheus.Labels{"EventSource": response.SourceApp, "BunldeId": request.Redemption_Bunlde_Id, "Level": loyalty_Account.Loyalty_Level_Key}).Add(response.Points_To_Redeem)

		record, err := Uc.Customer_Loyalty_Account_GetRedemption_Rules(response.MSISDN)
		if record.Bundles_Notification {
			if err != nil {
				response.Status = "failed"
				response.StatusCode = http.StatusBadRequest
				response.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
				response.ErrorDescription = err.Error()
				return
			}
			BundleNotiLog := NotificationLog{
				SourceAction:  "BundleRedemption",
				TransactionId: "",
				Medium:        "SMS",
				SourceAddress: record.Bundles_Notification_Sender,
				Destination:   response.MSISDN,
				Subject:       "Redemption",
				AddUser:       "SYSTEM",
				AddDate:       time.Now(),
			}
			Bundle_Noti_Text := ""
			var bundle propylaea.Bundle
			Bundle_Noti_Text = record.Bundles_Notification_Text
			bundleResponse, err := Uc.Propylaea.PropylaeaClient.Get_Bundle(response.Redemption_Bunlde_Id, "", "", "")
			if err != nil {
				fmt.Println("failed to get bundle")
			}
			var bundleName string
			if err == nil && len(bundleResponse.Data) > 0 {
				bundle = bundleResponse.Data[0]
				bundleName = bundle.Name_Lang1
			} else {
				bundleName = ""
			}
			// Variables to replace
			Minutes := 0.0
			MBs := 0
			SMS := 0
			Bonus_Minutes := 0.0
			Bonus_MBs := 0
			Bonus_SMS := 0

			Validity_value := bundle.Validity_Value
			Validity_type := ""
			// Minutes and SMSs from IN tokens
			seconds := 0.0
			bonus_seconds := 0.0
			for _, token := range bundle.IN_Tokens {
				if token.Token_Type == "SMS (whole SMS messages)" {
					if token.IsBonus {
						Bonus_SMS += int(token.Token_Value)
					} else {
						SMS += int(token.Token_Value)
					}
				} else if token.Token_Type == "time (seconds)" {
					if token.IsBonus {
						bonus_seconds += token.Token_Value
					} else {
						seconds += token.Token_Value
					}
				}
			}

			if seconds > 0 {
				Minutes = seconds / 60
			}
			if bonus_seconds > 0 {
				Bonus_Minutes = bonus_seconds / 60
			}
			// MBs from Alepo and Protei
			var DO_bytes float64 = 0
			var bonus_DO_bytes float64 = 0
			for _, DO_A := range bundle.Data_Offers {
				var size float64
				size, err := strconv.ParseFloat(DO_A.Offer_Size_Value, 64)
				if err != nil {
					fmt.Println("Error parsing Alepo Offer size value:", DO_A.Offer_Size_Value)
					continue
				}

				// Use size as float64 here
				fmt.Printf("Parsed size: %.2f\n", size) // e.g., show 2 decimal places

				if DO_A.IsBonus {
					if Configuration.Operation == "DRC" {
						if DO_A.Offer_Size_Unit == "MB" {
							size += (size * 1000 * 1000)
						} else if DO_A.Offer_Size_Unit == "GB" {
							size += (size * 1000 * 1000 * 1000)
						}
						bonus_DO_bytes += size
					} else {
						if DO_A.Offer_Size_Unit == "MB" {
							size += (size * 1024 * 1024)
						} else if DO_A.Offer_Size_Unit == "GB" {
							size += (size * 1024 * 1024 * 1024)
						}
						bonus_DO_bytes += size
					}

				} else {
					if Configuration.Operation == "DRC" {
						if DO_A.Offer_Size_Unit == "MB" {
							size += (size * 1000 * 1000)
						} else if DO_A.Offer_Size_Unit == "GB" {
							size += (size * 1000 * 1000 * 1000)
						}
						DO_bytes += size
					} else {
						if DO_A.Offer_Size_Unit == "MB" {
							size += (size * 1024 * 1024)
						} else if DO_A.Offer_Size_Unit == "GB" {
							size += (size * 1024 * 1024 * 1024)
						}
						DO_bytes += size
					}
				}
			}

			if DO_bytes > 0 {
				if Configuration.Operation == "DRC" {
					MBs = int(DO_bytes / 1000 / 1000)
				} else {
					MBs = int(DO_bytes / 1024 / 1024)
				}
			}
			if bonus_DO_bytes > 0 {
				if Configuration.Operation == "DRC" {
					Bonus_MBs = int(bonus_DO_bytes / 1000 / 1000)
				} else {
					Bonus_MBs = int(bonus_DO_bytes / 1024 / 1024)
				}
			}

			// Expiry Date from bundle Validity
			switch bundle.Validity_Type {
			case "Hourly":
				if bundle.Validity_Value == 1 {
					Validity_type = "hour"
				} else {
					Validity_type = "hours"
				}
			case "Daily":
				if bundle.Validity_Value == 1 {
					Validity_type = "day"
				} else {
					Validity_type = "days"
				}
			case "Weekly":
				if bundle.Validity_Value == 1 {
					Validity_type = "week"
				} else {
					Validity_type = "weeks"
				}
			case "Monthly":
				if bundle.Validity_Value == 1 {
					Validity_type = "month"
				} else {
					Validity_type = "months"
				}
			case "Yearly":
				if bundle.Validity_Value == 1 {
					Validity_type = "year"
				} else {
					Validity_type = "years"
				}
			}
			sizeParts := []string{}

			totalMBs := MBs + Bonus_MBs
			if totalMBs > 0 {
				formatted := ""
				if Configuration.Operation == "DRC" {
					formatted = FormatMBs1000(float64(totalMBs))

				} else {
					formatted = FormatMBs1024(float64(totalMBs))
				}
				sizeParts = append(sizeParts, formatted)
			}

			totalSMS := SMS + Bonus_SMS
			if totalSMS > 0 {
				sizeParts = append(sizeParts, fmt.Sprintf("%d SMS", totalSMS))
			}

			totalMinutes := Minutes + Bonus_Minutes

			if totalMinutes > 0 {
				sizeParts = append(sizeParts, fmt.Sprintf("%.1f Minutes", totalMinutes))
			}

			Size := strings.Join(sizeParts, " + ")

			Validity := fmt.Sprintf("%d %s", Validity_value, Validity_type)
			if Bundle_Noti_Text != "" {
				Bundle_Noti_Text = strings.ReplaceAll(Bundle_Noti_Text, "{{ PointsDeducted }}", fmt.Sprint(response.Points_To_Redeem))
				Bundle_Noti_Text = strings.ReplaceAll(Bundle_Noti_Text, "{{ BundleName }}", fmt.Sprint(bundleName))
				Bundle_Noti_Text = strings.ReplaceAll(Bundle_Noti_Text, "{{ BundleSize }}", fmt.Sprint(Size))
				Bundle_Noti_Text = strings.ReplaceAll(Bundle_Noti_Text, "{{ BundleValidity }}", fmt.Sprint(Validity))
				Bundle_Noti_Text = strings.ReplaceAll(Bundle_Noti_Text, "{{ LoyaltyBalance }}", fmt.Sprint(response.Closure_Available_Points))
				BundleNotiLog.Payload = Bundle_Noti_Text

				err := SendSMS(record.Bundles_Notification_Sender, response.MSISDN, Bundle_Noti_Text)
				if err != nil {
					BundleNotiLog.Status = "Failed"
					BundleNotiLog.Error = err.Error()
				} else {
					BundleNotiLog.Status = "Successful"
				}
			} else {
				BundleNotiLog.Payload = Bundle_Noti_Text
				BundleNotiLog.Status = "Failed"
				BundleNotiLog.Error = "Undefined Bundle notification for transaction"
			}
			YYYY, MM, _, _, _, _, _ := GetTimeParts(response.ReceiveDate)
			Db := Configuration.DB_Name_Loyalty + "_Logs_" + YYYY + MM
			{
				notiCtx, notiCancel := context.WithTimeout(context.Background(), 10*time.Second)
				_, notiErr := Uc.MongoClient.Mongo.Database(Db).Collection("Col_NotificationLog").InsertOne(notiCtx, BundleNotiLog)
				notiCancel()
				if notiErr != nil {
					log.Println("Error in Write Bundle Notification Logs:", notiErr, " (", BundleNotiLog, ")")
				}
			}
		}
	case "MobileMoney":
		if request.Redemption_Amount <= 0 && request.Points_To_Redeem <= 0 {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "invalid redemption amount"
			response.ErrorDescription = "invalid redemption airtime amount"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		response.MinAvailableRequiredPoints = redemption_Rules.Available_MinPoints_for_MobileMoney
		if response.Opening_Available_Points < response.MinAvailableRequiredPoints {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "no enough points"
			response.ErrorDescription = "no enough points"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		if redemption_Rules.MobileMoney_MerchantAccount == "" {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "mobile money payer account is not provided"
			response.ErrorDescription = "mobile money payer account is not provided"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		//debit loyalty points
		if redemption_Rules.MobileMoney_Points <= 0 || redemption_Rules.MobileMoney_Amount <= 0 {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "mobile money redeem rules are not defined"
			response.ErrorDescription = "mobile money redeem rules are not defined"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		//do the calculation
		if response.Points_To_Redeem > 0 {
			//calculate Redemption_Amount
			response.Redemption_Amount = (response.Points_To_Redeem * redemption_Rules.MobileMoney_Amount) / redemption_Rules.MobileMoney_Points
		} else {
			if response.Redemption_Amount > 0 {
				//calculate Points_To_Redeem
				response.Points_To_Redeem = (response.Redemption_Amount * redemption_Rules.MobileMoney_Points) / redemption_Rules.MobileMoney_Amount
			} else {
				//return error
				response.Status = "failed"
				response.StatusCode = http.StatusBadRequest
				response.StatusDescription = "invalid redemption amount"
				response.ErrorDescription = "invalid redemption mobile money amount"
				response.StatusDate = time.Now()
				response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
				Uc.Write_Loyalty_Redemption_log(*response)
				return
			}
		}
		response.MinRequiredPoints = redemption_Rules.MobileMoney_MinPoints
		if response.Points_To_Redeem < response.MinRequiredPoints {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "requested points are less than the minimum allowed points for redemption"
			response.ErrorDescription = "requested points are less than the minimum allowed points for redemption"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		//check if subscriber have enough points
		if response.Points_To_Redeem > response.Opening_Available_Points {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "no enough points"
			response.ErrorDescription = "no enough points"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		var Debit_Request Loyalty_AccountDebitPoints_Request
		Debit_Request.MSISDN = response.MSISDN
		Debit_Request.Debit_Amount = response.Points_To_Redeem
		Debit_Request.Debit_Reason = "Redeem Request"
		Debit_Request.Redemption_Type = "MobileMoney" //Airtime, Bundle, MobileMoney, SpinAndWin
		Debit_Request.Redemption_Bunlde_Id = ""
		Debit_Request.Redemption_Amount = response.Redemption_Amount
		var debit_Log Loyalty_AccountDebitPoints_log
		Uc.Loyalty_AccountDebitPoints(request_header, Debit_Request, &debit_Log)
		response.Points_Debit_Result = debit_Log
		if debit_Log.StatusCode != http.StatusOK {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = debit_Log.StatusDescription
			response.ErrorDescription = debit_Log.ErrorDescription
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		response.Closure_Awarded_Points = debit_Log.Closure_Awarded_Points
		response.Closure_Redeemed_Points = debit_Log.Closure_Redeemed_Points
		response.Closure_Available_Points = debit_Log.Closure_Available_Points
		//to do: credit mobile money amount --> merchant transfer
		MobileMoney_MerchantPIN, err := DecryptHexString(redemption_Rules.MobileMoney_MerchantPIN)
		if err != nil {
			fmt.Println("error in decrypting artime evc pin", err.Error())
		}

		// mm_CashIn_Reply, err := Uc.CGW.UC_GWClient.MM_Agent_CashIN_ToBonusWallet(MM.CashIN_Request{
		// 	SenderMSISDN:   redemption_Rules.MobileMoney_MerchantAccount,
		// 	SenderPIN:      MobileMoney_MerchantPIN,
		// 	ReceiverMSISDN: request.MSISDN,
		// 	Amount:         fmt.Sprintf("%f", response.Redemption_Amount),
		// 	Remark:         "Loyalty Redemption",
		// 	// Currency:       "102",
		// })

		// var request MM_CashIn_request
		// mm_CashIn_Reply, err :=Uc.APGW.APGWClient.MM_CashIn(request{
		// TransactionAmount :strconv.FormatFloat(response.Redemption_Amount, 'f', -1, 64),
		// Remarks : "Loyalty Redemption",
		// Transactor.IdValue : redemption_Rules.MobileMoney_MerchantAccount,
		// Transactor.Mpin : MobileMoney_MerchantPIN,
		// Receiver.IdType : "mobileNumber",
		// // Receiver.ProductId : request.ReceiverProductId,
		// Receiver.IdValue  :request.MSISDN,
		// })

		var mm_request APGWClientV2.MM_CashIn_Request
		mm_request.TransactionAmount = strconv.FormatFloat(response.Redemption_Amount, 'f', -1, 64)
		mm_request.Remarks = "Loyalty Redemption"

		mm_request.Transactor.IdType = "mobileNumber"
		mm_request.Transactor.IdValue = redemption_Rules.MobileMoney_MerchantAccount
		mm_request.Transactor.Mpin = MobileMoney_MerchantPIN

		mm_request.Receiver.IdType = "mobileNumber"
		mm_request.Receiver.ProductId = "12"
		mm_request.Receiver.IdValue = request.MSISDN
		mm_CashIn_Reply, err := Uc.APGW.APGWClient.MM_CashIn(mm_request)
		response.MobileMoney_PurchaseResult = mm_CashIn_Reply
		fmt.Println("mm_CashIn_Reply.Status", mm_CashIn_Reply.Status)
		if err != nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to redeem mobile money"
			response.ErrorDescription = err.Error()
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			//refund points
			var refund_response Loyalty_AccountCreditPoints_log
			var refundRequest Loyalty_AccountCreditPoints_Request
			refundRequest.MSISDN = response.MSISDN
			refundRequest.EventSource = response.SourceApp
			refundRequest.EventType = "refund"
			refundRequest.EventDetail = request.Redemption_Type + " - refund"
			refundRequest.EventAmount = 0
			refundRequest.PointsToCredit = response.Points_To_Redeem
			refundRequest.EventDescription = ""
			Uc.Loyalty_AccountCreditPoints(request_header, refundRequest, &refund_response, true)
			return
		}
		if mm_CashIn_Reply.Data.Status != "SUCCEEDED" {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to redeem mobile money"
			response.ErrorDescription = mm_CashIn_Reply.Data.Status
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			//refund points
			var refund_response Loyalty_AccountCreditPoints_log
			var refundRequest Loyalty_AccountCreditPoints_Request
			refundRequest.MSISDN = response.MSISDN
			refundRequest.EventSource = response.SourceApp
			refundRequest.EventType = "refund"
			refundRequest.EventDetail = request.Redemption_Type + " - refund"
			refundRequest.EventAmount = 0
			refundRequest.PointsToCredit = response.Points_To_Redeem
			refundRequest.EventDescription = ""
			Uc.Loyalty_AccountCreditPoints(request_header, refundRequest, &refund_response, true)
			return
		}
		MobileMoneyRedemptionCount.With(prometheus.Labels{"EventSource": response.SourceApp, "Level": loyalty_Account.Loyalty_Level_Key}).Inc()
		MobileMoneyRedemptionPoints.With(prometheus.Labels{"EventSource": response.SourceApp, "Level": loyalty_Account.Loyalty_Level_Key}).Add(response.Points_To_Redeem)
		MobileMoneyRedemptionAmount.With(prometheus.Labels{"EventSource": response.SourceApp, "Level": loyalty_Account.Loyalty_Level_Key}).Add(response.Redemption_Amount)

		record, err := Uc.Customer_Loyalty_Account_GetRedemption_Rules(response.MSISDN)
		if err != nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
			response.ErrorDescription = err.Error()
			return
		}
		if record.MobileMoney_Notification {
			MobileMoneyNotiLog := NotificationLog{
				SourceAction:  "MobileMoneyRedemption",
				TransactionId: "",
				Medium:        "SMS",
				SourceAddress: record.MobileMoney_Notification_Sender,
				Destination:   response.MSISDN,
				Subject:       "Redemption",
				AddUser:       "SYSTEM",
				AddDate:       time.Now(),
			}
			MobileMoney_Noti_Text := ""
			MobileMoney_Noti_Text = record.MobileMoney_Notification_Text

			if MobileMoney_Noti_Text != "" {
				MobileMoney_Noti_Text = strings.ReplaceAll(MobileMoney_Noti_Text, "{{ PointsDeducted }}", fmt.Sprint(response.Points_To_Redeem))
				MobileMoney_Noti_Text = strings.ReplaceAll(MobileMoney_Noti_Text, "{{ MobileMoneyAwarded }}", fmt.Sprint(response.Redemption_Amount))
				MobileMoney_Noti_Text = strings.ReplaceAll(MobileMoney_Noti_Text, "{{ LoyaltyBalance }}", fmt.Sprint(response.Closure_Available_Points))
				MobileMoneyNotiLog.Payload = MobileMoney_Noti_Text
				err := SendSMS(record.MobileMoney_Notification_Sender, response.MSISDN, MobileMoney_Noti_Text)
				if err != nil {
					MobileMoneyNotiLog.Status = "Failed"
					MobileMoneyNotiLog.Error = err.Error()
				} else {
					MobileMoneyNotiLog.Status = "Successful"
				}
			} else {
				MobileMoneyNotiLog.Payload = MobileMoney_Noti_Text
				MobileMoneyNotiLog.Status = "Failed"
				MobileMoneyNotiLog.Error = "Undefined Mobile Money notification for transaction"
			}
			YYYY, MM, _, _, _, _, _ := GetTimeParts(response.ReceiveDate)
			Db := Configuration.DB_Name_Loyalty + "_Logs_" + YYYY + MM
			{
				notiCtx, notiCancel := context.WithTimeout(context.Background(), 10*time.Second)
				_, notiErr := Uc.MongoClient.Mongo.Database(Db).Collection("Col_NotificationLog").InsertOne(notiCtx, MobileMoneyNotiLog)
				notiCancel()
				if notiErr != nil {
					log.Println("Error in Write Mobile Money Notification Logs:", notiErr, " (", MobileMoneyNotiLog, ")")
				}
			}
		}
	case "SpinAndWin":
		if request.Redemption_Amount <= 0 {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "invalid redemption amount"
			response.ErrorDescription = "invalid redemption spin and win amount"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		response.MinAvailableRequiredPoints = redemption_Rules.Available_MinPoints_for_SpinAndWin
		if response.Opening_Available_Points < response.MinAvailableRequiredPoints {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "no enough points"
			response.ErrorDescription = "no enough points"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		//debit loyalty points
		if redemption_Rules.FreeSpinAndWin_PointsPerSpin <= 0 {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "airtime redeem rules are not defined"
			response.ErrorDescription = "airtime redeem rules are not defined"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		//do the calculation
		response.Points_To_Redeem = response.Redemption_Amount * redemption_Rules.FreeSpinAndWin_PointsPerSpin
		response.MinRequiredPoints = redemption_Rules.FreeSpinAndWin_MinPoints
		if response.Points_To_Redeem < response.MinRequiredPoints {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "requested points are less than the minimum allowed points for redemption"
			response.ErrorDescription = "requested points are less than the minimum allowed points for redemption"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		//check if subscriber have enough points
		if response.Points_To_Redeem > response.Opening_Available_Points {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "no enough points"
			response.ErrorDescription = "no enough points"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}

		var Debit_Request Loyalty_AccountDebitPoints_Request
		Debit_Request.MSISDN = response.MSISDN
		Debit_Request.Debit_Amount = response.Points_To_Redeem
		Debit_Request.Debit_Reason = "Redeem Request"
		Debit_Request.Redemption_Type = "SpinAndWin" //Airtime, Bundle, MobileMoney, SpinAndWin
		Debit_Request.Redemption_Bunlde_Id = ""
		Debit_Request.Redemption_Amount = response.Redemption_Amount
		var debit_Log Loyalty_AccountDebitPoints_log
		Uc.Loyalty_AccountDebitPoints(request_header, Debit_Request, &debit_Log)
		response.Points_Debit_Result = debit_Log
		if debit_Log.StatusCode != http.StatusOK {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = debit_Log.StatusDescription
			response.ErrorDescription = debit_Log.ErrorDescription
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		response.Closure_Awarded_Points = debit_Log.Closure_Awarded_Points
		response.Closure_Redeemed_Points = debit_Log.Closure_Redeemed_Points
		response.Closure_Available_Points = debit_Log.Closure_Available_Points
		//awar spins chances
		spinAndWin_Reply, err := Uc.SpinAndWin.SpinAndWinClient.EligibleSubs_AddChances(SpinAndWin.EligibleSubs_AddChances_Request{
			Key:            request.MSISDN,
			SpinCountToAdd: int64(response.Redemption_Amount),
		})
		response.SpinAndWin_PurchaseResult = spinAndWin_Reply
		if err != nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to redeem Spin And Win Chances"
			response.ErrorDescription = err.Error()
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			//refund points
			var refund_response Loyalty_AccountCreditPoints_log
			var refundRequest Loyalty_AccountCreditPoints_Request
			refundRequest.MSISDN = response.MSISDN
			refundRequest.EventSource = response.SourceApp
			refundRequest.EventType = "refund"
			refundRequest.EventDetail = request.Redemption_Type + " - refund"
			refundRequest.EventAmount = 0
			refundRequest.PointsToCredit = response.Points_To_Redeem
			refundRequest.EventDescription = ""
			Uc.Loyalty_AccountCreditPoints(request_header, refundRequest, &refund_response, true)
			return
		}
		if spinAndWin_Reply.StatusCode != http.StatusOK {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = spinAndWin_Reply.StatusDescription
			response.ErrorDescription = spinAndWin_Reply.ErrorDescription
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			//refund points
			var refund_response Loyalty_AccountCreditPoints_log
			var refundRequest Loyalty_AccountCreditPoints_Request
			refundRequest.MSISDN = response.MSISDN
			refundRequest.EventSource = response.SourceApp
			refundRequest.EventType = "refund"
			refundRequest.EventDetail = request.Redemption_Type + " - refund"
			refundRequest.EventAmount = 0
			refundRequest.PointsToCredit = response.Points_To_Redeem
			refundRequest.EventDescription = ""
			Uc.Loyalty_AccountCreditPoints(request_header, refundRequest, &refund_response, true)
			return
		}

		record, err := Uc.Customer_Loyalty_Account_GetRedemption_Rules(response.MSISDN)
		if err != nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
			response.ErrorDescription = err.Error()
			return
		}
		if record.FreeSpinAndWin_Notification {
			SpinWinNotiLog := NotificationLog{
				SourceAction:  "SpinAndWinRedemption",
				TransactionId: "",
				Medium:        "SMS",
				SourceAddress: record.FreeSpinAndWin_Notification_Sender,
				Destination:   response.MSISDN,
				Subject:       "Redemption",
				AddUser:       "SYSTEM",
				AddDate:       time.Now(),
			}
			SpinWin_Noti_Text := ""
			SpinWin_Noti_Text = record.FreeSpinAndWin_Notification_Text

			if SpinWin_Noti_Text != "" {
				SpinWin_Noti_Text = strings.ReplaceAll(SpinWin_Noti_Text, "{{ PointsDeducted }}", fmt.Sprint(response.Points_To_Redeem))
				SpinWin_Noti_Text = strings.ReplaceAll(SpinWin_Noti_Text, "{{ SpinsAwarded }}", fmt.Sprint(response.Redemption_Amount))
				SpinWin_Noti_Text = strings.ReplaceAll(SpinWin_Noti_Text, "{{ LoyaltyBalance }}", fmt.Sprint(response.Closure_Available_Points))
				SpinWinNotiLog.Payload = SpinWin_Noti_Text
				err := SendSMS(record.FreeSpinAndWin_Notification_Sender, response.MSISDN, SpinWin_Noti_Text)
				if err != nil {
					SpinWinNotiLog.Status = "Failed"
					SpinWinNotiLog.Error = err.Error()
				} else {
					SpinWinNotiLog.Status = "Successful"
				}
			} else {
				SpinWinNotiLog.Payload = SpinWin_Noti_Text
				SpinWinNotiLog.Status = "Failed"
				SpinWinNotiLog.Error = "Undefined Mobile Money notification for transaction"
			}
			YYYY, MM, _, _, _, _, _ := GetTimeParts(response.ReceiveDate)
			Db := Configuration.DB_Name_Loyalty + "_Logs_" + YYYY + MM
			{
				notiCtx, notiCancel := context.WithTimeout(context.Background(), 10*time.Second)
				_, notiErr := Uc.MongoClient.Mongo.Database(Db).Collection("Col_NotificationLog").InsertOne(notiCtx, SpinWinNotiLog)
				notiCancel()
				if notiErr != nil {
					log.Println("Error in Write Mobile Money Notification Logs:", notiErr, " (", SpinWinNotiLog, ")")
				}
			}
		}

	default:
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "invalid redemption type"
		response.ErrorDescription = "invalid redemption type (accepted values: Airtime, Bundle, MobileMoney, SpinAndWin)"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_Redemption_log(*response)
		return
	}
	//successful reply
	response.Status = "successful"
	response.StatusCode = http.StatusOK
	response.StatusDescription = ""
	response.ErrorDescription = ""
	response.StatusDate = time.Now()
	response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
	Uc.Write_Loyalty_Redemption_log(*response)
}

func Loyalty_Account_Segment_Selection(Amount float64, FirstUse_date time.Time) (scheme_name string) {
	//AON_Hours := time.Now().Sub(FirstUse_date).Hours()
	AON_Hours := time.Since(FirstUse_date).Hours()
	AON_Months := (AON_Hours / 24) / 30
	Schemes, redisErr := redisx.GetAllJSONByPattern[Loyalty_Account_Segment](context.Background(), RedisClient, redisx.ScanJSONOptions{
		Pattern: "Loyalty_Account_Segment:*", ScanCount: 500, PipelineSize: 250,
	})
	if redisErr != nil {
		log.Println("error getting Loyalty_Account_Segment from Redis:", redisErr)
	}
	if len(Schemes) > 0 {
		for _, scheme := range Schemes {
			if Amount >= scheme.Amount_From && Amount < scheme.Amount_Till {
				if AON_Months >= scheme.AON_From && AON_Months < scheme.AON_Till {
					//DailyImportSubsStats.With(prometheus.Labels{"IsElligble": "true", "Reason": "", "Scheme": scheme.Key}).Inc()
					//DailyImportSubsStats.With(prometheus.Labels{"IsElligble": "true", "Reason": "", "Scheme": "All"}).Inc()
					return scheme.Key
				}
			}
		}
	}
	// if scheme_name == "" {
	// 	DailyImportSubsStats.With(prometheus.Labels{"IsElligble": "false", "Reason": "", "Scheme": ""}).Inc()
	// }
	return
}

func Loyalty_Level_Selection(Accumulated_Points float64) (level_key string) {
	levels, redisErr := redisx.GetAllJSONByPattern[Loyalty_Level](context.Background(), RedisClient, redisx.ScanJSONOptions{
		Pattern: "Loyalty_Level:*", ScanCount: 500, PipelineSize: 250,
	})
	if redisErr != nil {
		log.Println("error getting Loyalty_Level from Redis:", redisErr)
	}
	if len(levels) > 0 {
		for _, level := range levels {
			if Accumulated_Points >= level.Min_Accumulated_Points && Accumulated_Points < level.Max_Accumulated_Points {
				return level.Key
			}
		}
	}
	// if scheme_name == "" {
	// 	DailyImportSubsStats.With(prometheus.Labels{"IsElligble": "false", "Reason": "", "Scheme": ""}).Inc()
	// }
	return
}

func (Uc *UserControl) Loyalty_AccountCreditPoints(request_header *Request_Header, request Loyalty_AccountCreditPoints_Request, response *Loyalty_AccountCreditPoints_log, refund ...bool) {
	response.ReceiveDate = time.Now()
	//fill the request header info
	response.SourceIP = request_header.SourceIP
	response.SourceApp = request.EventSource
	response.AppLogin = request_header.AppLogin
	response.AppVersion = request_header.AppVersion
	response.GPSLocation = request_header.GPSLocation
	response.GSMLocation = request_header.GSMLocation
	response.MSISDN = request.MSISDN
	request.MSISDN = Normalize_International_MSISDN(request.MSISDN)
	if request.MSISDN == "" {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "invalid msisdn"
		response.ErrorDescription = "invalid msisdn"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		//Uc.Write_Loyalty_AccountCreditPoints_log(*response)
		return
	}
	//fill the request info
	response.MSISDN = request.MSISDN
	response.EventSource = request.EventSource
	response.EventType = request.EventType
	response.EventDetail = request.EventDetail
	response.EventAmount = request.EventAmount
	response.PointsToCredit = request.PointsToCredit
	response.EventDescription = request.EventDescription
	response.EventDetailCode = request.EventDetailCode

	//validate loyalty account
	loyalty_account, loyaltyAccErr := redisx.GetJSON[Customer_Loyalty_Account](context.Background(), RedisClient, Customer_Loyalty_Account{Key: request.MSISDN}.RedisKey())
	if redisx.IsNil(loyaltyAccErr) {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "failed to get loyalty account"
		response.ErrorDescription = "loyalty account does not exist"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_AccountCreditPoints_log(*response)
		return
	}
	if loyaltyAccErr != nil {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "failed to get loyalty account"
		response.ErrorDescription = "type assertion issue with Customer_Loyalty_Account"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_AccountCreditPoints_log(*response)
		return
	}
	response.Opening_Loyalty_Level_Key = loyalty_account.Loyalty_Level_Key
	response.Opening_Loyalty_Account_Segment_Key = loyalty_account.Loyalty_Account_Segment_Key
	response.Opening_Awarded_Points = loyalty_account.Awarded_Points
	response.Opening_Redeemed_Points = loyalty_account.Redeemed_Points
	response.Opening_Available_Points = loyalty_account.Available_Points
	response.Opening_Outstanding_fraction_points = loyalty_account.Outstanding_fraction_points

	//check opt status
	if loyalty_account.Opt_Status != "OptedIn" && Configuration.ISLoyaltyOptIn {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "customer is opted out"
		response.ErrorDescription = "customer is opted out"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_AccountCreditPoints_log(*response)
		return
	}
	//check exclusion list
	{
		chkCtx, chkCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, chkErr := redisx.GetJSON[Customer_Exclusion](chkCtx, RedisClient, Customer_Exclusion{Key: loyalty_account.Key}.RedisKey())
		chkCancel()
		if chkErr == nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "customer is included in the exclusion list"
			response.ErrorDescription = "customer is included in the exclusion list"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_AccountCreditPoints_log(*response)
			return
		}
	}
	//check COS exclusion list
	{
		chkCtx, chkCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, chkErr := redisx.GetJSON[Customer_COS_Exclusion](chkCtx, RedisClient, Customer_COS_Exclusion{Key: loyalty_account.COS}.RedisKey())
		chkCancel()
		if chkErr == nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "customer is included in the cos exclusion list"
			response.ErrorDescription = "customer is included in the cos exclusion list"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_AccountCreditPoints_log(*response)
			return
		}
	}
	//exclude loyalty redemption accounts
	if response.EventSource == "MobileMoney_feed" && response.EventType == "CASHIN" && response.EventDetail != "" {
		redemption_Rules, err := Uc.Customer_Loyalty_Account_GetRedemption_Rules(loyalty_account.Key)
		if err != nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to get redemption rules"
			response.ErrorDescription = err.Error()
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_AccountCreditPoints_log(*response)
			return
		}
		response.EventDetail = Normalize_International_MSISDN(response.EventDetail)
		var transactionAgent, loyaltyAgent string
		if len(response.EventDetail) >= Configuration.MSISDN_Short_len {
			transactionAgent = response.EventDetail[len(response.EventDetail)-Configuration.MSISDN_Short_len:]
		}
		if len(redemption_Rules.MobileMoney_MerchantAccount) >= Configuration.MSISDN_Short_len {
			loyaltyAgent = redemption_Rules.MobileMoney_MerchantAccount[len(redemption_Rules.MobileMoney_MerchantAccount)-Configuration.MSISDN_Short_len:]
		}
		if transactionAgent != "" && loyaltyAgent != "" {
			if transactionAgent == loyaltyAgent {
				response.Status = "failed"
				response.StatusCode = http.StatusBadRequest
				response.StatusDescription = "loyalty redemption transaction"
				response.ErrorDescription = "loyalty redemption transaction is not entitled for loyalty redemption"
				response.StatusDate = time.Now()
				response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
				Uc.Write_Loyalty_AccountCreditPoints_log(*response)
				return
			}
		}
	}
	if loyalty_account.Loyalty_Account_Segment_Key == "" {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "loyalty account segment is not assigned"
		response.ErrorDescription = "loyalty account segment is not assigned"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_AccountCreditPoints_log(*response)
		return
	}
	if loyalty_account.Loyalty_Level_Key == "" {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "loyalty account level is not assigned"
		response.ErrorDescription = "loyalty account level is not assigned"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_AccountCreditPoints_log(*response)
		return
	}
	//validate loyalty level
	loyalty_Level, loyaltyLevelErr := redisx.GetJSON[Loyalty_Level](context.Background(), RedisClient, Loyalty_Level{Key: loyalty_account.Loyalty_Level_Key}.RedisKey())
	if redisx.IsNil(loyaltyLevelErr) {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "failed to get loyalty account level"
		response.ErrorDescription = "loyalty account level is not defined"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_AccountCreditPoints_log(*response)
		return
	}
	if loyaltyLevelErr != nil {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "failed to get loyalty account level"
		response.ErrorDescription = "loyalty account level assertion issue"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_AccountCreditPoints_log(*response)
		return
	}
	if !loyalty_Level.EnableRedeem {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "loyalty account level redeemption is disabled"
		response.ErrorDescription = "loyalty account level redeemption is disabled"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_AccountCreditPoints_log(*response)
		return
	}
	//validate the loyalty plan
	plan, planErr := redisx.GetJSON[Loyalty_Plan](context.Background(), RedisClient, Loyalty_Plan{Key: loyalty_account.Loyalty_Level_Key + "|" + loyalty_account.Loyalty_Account_Segment_Key}.RedisKey())
	if redisx.IsNil(planErr) {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "failed to get loyalty plan"
		response.ErrorDescription = "loyalty plan does not exist"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_AccountCreditPoints_log(*response)
		return
	}
	if planErr != nil {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "failed to get loyalty plan"
		response.ErrorDescription = "type assertion issue with Loyalty_Plan"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_AccountCreditPoints_log(*response)
		return
	}
	//validate earning rules
	if plan.Earning_Rules_Key == "" {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "failed to get earning rules"
		response.ErrorDescription = "earning rules are not defined"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_AccountCreditPoints_log(*response)
		return
	}
	point_earning_rules, earningErr := redisx.GetJSON[Loyalty_Point_Earning_Rules](context.Background(), RedisClient, Loyalty_Point_Earning_Rules{Key: plan.Earning_Rules_Key}.RedisKey())
	if redisx.IsNil(earningErr) {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "failed to get earning rules"
		response.ErrorDescription = "point earning rules are not defined"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_AccountCreditPoints_log(*response)
		return
	}
	if earningErr != nil {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "failed to get earning rules"
		response.ErrorDescription = "type assertion issue with Loyalty_Point_Earning_Rules"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_AccountCreditPoints_log(*response)
		return
	}
	//calucate points
	sms := false
	sender := ""
	text := ""
	var points, Outstanding_fraction_points float64
	initial_Outstanding_fraction_points := loyalty_account.Outstanding_fraction_points
	if response.PointsToCredit == 0 {
		points, Outstanding_fraction_points, sms, sender, text = Calculate_Loyalty_Points(point_earning_rules, request, loyalty_account.Outstanding_fraction_points)
		if !loyalty_account.Joining_Date.IsZero() {
			now := time.Now()
			years := now.Year() - loyalty_account.Joining_Date.Year()
			months := int(now.Month()) - int(loyalty_account.Joining_Date.Month())
			totalMonths := years*12 + months
			// Adjust if current day is before the joining day in the month
			if now.Day() < loyalty_account.Joining_Date.Day() {
				totalMonths--
			}
			if totalMonths >= 0 {
				for _, lvl := range loyalty_Level.Seniority_Levels {
					seniority, seniorityErr := redisx.GetJSON[Loyalty_Seniority_Level](context.Background(), RedisClient, Loyalty_Seniority_Level{Key: lvl.Loyalty_Seniority_Level_Key}.RedisKey())
					if seniorityErr == nil {
						if totalMonths >= int(seniority.AON_From) && totalMonths <= int(seniority.AON_Till) {
							changedPointsflt := (points + Outstanding_fraction_points) - initial_Outstanding_fraction_points
							seniorityPoints := changedPointsflt * lvl.Multiplier_Percentage / 100
							endPoints := seniorityPoints + changedPointsflt + initial_Outstanding_fraction_points
							points = math.Floor(endPoints)
							Outstanding_fraction_points = endPoints - points
							break
						}
					}

				}
			}

		}
	} else {
		points = response.PointsToCredit
		Outstanding_fraction_points = loyalty_account.Outstanding_fraction_points
	}
	if points > 0 {
		//response.OpeningAvailablePoints = (loyalty_account.Awarded_Points + loyalty_account.Expired_Points) - loyalty_account.Redeemed_Points
		response.AwardedPoints = points
		//validate governance rules
		loyalty_governance, lgErr := redisx.GetJSON[Loyalty_Governance](context.Background(), RedisClient, Loyalty_Governance{Key: LOYALTY_GOVERNANCE_KEY}.RedisKey())
		if redisx.IsNil(lgErr) {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to get governance entry"
			response.ErrorDescription = "loyalty governance entry not found"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_AccountCreditPoints_log(*response)
			return
		}
		if lgErr != nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to get governance entry"
			response.ErrorDescription = "loyalty governance type assertion issue"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_AccountCreditPoints_log(*response)
			return
		}
		if points > loyalty_governance.MaxAllowedPoints_PerTransaction {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to credit loyalty account"
			response.ErrorDescription = "awarded points per transaction is exceeding loyalty governance rules"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_AccountCreditPoints_log(*response)
			return
		}
		if (loyalty_account.Awarded_Points + points) > loyalty_governance.MaxSubsAwardedPoints {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to credit loyalty account"
			response.ErrorDescription = "awarded points per subscriber is exceeding loyalty governance rules"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_AccountCreditPoints_log(*response)
			return
		}
		//credit loyalty account
		refundValue := false
		if len(refund) > 0 {
			refundValue = refund[0]
		}
		if !refundValue {
			loyalty_account.Awarded_Points = loyalty_account.Awarded_Points + points
		} else {
			loyalty_account.Redeemed_Points = loyalty_account.Redeemed_Points - points
		}
		loyalty_account.Available_Points = loyalty_account.Awarded_Points - loyalty_account.Redeemed_Points
		loyalty_account.Last_Award_Date = time.Now()
		loyalty_account.Outstanding_fraction_points = Outstanding_fraction_points
		YYYY, MM, _, _, _, _, _ := GetTimeParts(time.Now())
		//prepare the monthly points detail
		var PointsDetail Customer_Loyalty_Account_Points_Detail
		PointsDetail, pdErr := redisx.GetJSON[Customer_Loyalty_Account_Points_Detail](context.Background(), RedisClient, Customer_Loyalty_Account_Points_Detail{Key: request.MSISDN + "|" + YYYY + MM}.RedisKey())
		if redisx.IsNil(pdErr) {
			PointsDetail.Key = request.MSISDN + "|" + YYYY + MM
			PointsDetail.Year_Month = YYYY + MM
			PointsDetail.Creation_date = time.Now()
			PointsDetail.Awarded_Points = points
			PointsDetail.Available_Points = PointsDetail.Awarded_Points
			PointsDetail.Last_Credit_Date = time.Now()
			loyalty_account.Points_Detail_Keys = append(loyalty_account.Points_Detail_Keys, request.MSISDN+"|"+YYYY+MM)
		} else if pdErr != nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to credit loyalty account"
			response.ErrorDescription = "issue with Customer_Loyalty_Account_Points_Detail type assertion"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_AccountCreditPoints_log(*response)
			return
		} else {
			if !refundValue {
				PointsDetail.Awarded_Points = PointsDetail.Awarded_Points + points
			} else {
				PointsDetail.Redeemed_Points = PointsDetail.Redeemed_Points - points
			}
			PointsDetail.Available_Points = PointsDetail.Awarded_Points - PointsDetail.Redeemed_Points
			PointsDetail.Last_Credit_Date = time.Now()
		}
		if PointsDetail.Awarded_Points > loyalty_governance.MaxSubsAwardedPoints_PerMonth {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to credit loyalty account"
			response.ErrorDescription = "monthly awarded points per subscriber is exceeding loyalty governance rules"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_AccountCreditPoints_log(*response)
			return
		}
		resetDailyWeeklyCountersIfNeeded(&loyalty_account)
		if loyalty_governance.DailyEarningLimit > 0 && (loyalty_account.Daily_Earned_Points + points) > loyalty_governance.DailyEarningLimit {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to credit loyalty account"
			response.ErrorDescription = "daily earning limit exceeded"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_AccountCreditPoints_log(*response)
			return
		}
		if loyalty_governance.WeeklyEarningLimit > 0 && (loyalty_account.Weekly_Earned_Points + points) > loyalty_governance.WeeklyEarningLimit {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to credit loyalty account"
			response.ErrorDescription = "weekly earning limit exceeded"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_AccountCreditPoints_log(*response)
			return
		}
		err := Uc.Loyalty_Governance_Available_Points_Debit(points, refundValue)
		if err == nil {
			if !refundValue {
				loyalty_account.Daily_Earned_Points += points
				loyalty_account.Weekly_Earned_Points += points
			}
			{
				putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
				if _, putErr := Mdb_Customer_Loyalty_Account.Coll.UpdateOne(putCtx, bson.M{"Key": loyalty_account.Key}, bson.M{"$set": loyalty_account}, options.UpdateOne().SetUpsert(true)); putErr != nil {
					log.Println("Mdb_Customer_Loyalty_Account upsert error:", putErr)
				}
				if putSetErr := redisx.SetJSON(putCtx, RedisClient, loyalty_account.RedisKey(), loyalty_account); putSetErr != nil {
					log.Println("redisx.SetJSON Customer_Loyalty_Account error:", putSetErr)
				}
				putCancel()
			}
			{
				putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
				if _, putErr := Mdb_Customer_Loyalty_Account_Points_Detail.Coll.UpdateOne(putCtx, bson.M{"Key": PointsDetail.Key}, bson.M{"$set": PointsDetail}, options.UpdateOne().SetUpsert(true)); putErr != nil {
					log.Println("Mdb_Customer_Loyalty_Account_Points_Detail upsert error:", putErr)
				}
				if putSetErr := redisx.SetJSON(putCtx, RedisClient, PointsDetail.RedisKey(), PointsDetail); putSetErr != nil {
					log.Println("redisx.SetJSON Customer_Loyalty_Account_Points_Detail error:", putSetErr)
				}
				putCancel()
			}
			new_Loyalty_level_key, errNL := Uc.EvaluateAndUpdate_CustomerLoyaltyLevel(response.AppLogin, loyalty_account.Key)
			if errNL != nil {
				loyalty_account.Loyalty_Level_Key = new_Loyalty_level_key
			}
		}
	} else {
		if Outstanding_fraction_points > 0 {
			loyalty_account.Outstanding_fraction_points = Outstanding_fraction_points
			{
				putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
				if _, putErr := Mdb_Customer_Loyalty_Account.Coll.UpdateOne(putCtx, bson.M{"Key": loyalty_account.Key}, bson.M{"$set": loyalty_account}, options.UpdateOne().SetUpsert(true)); putErr != nil {
					log.Println("Mdb_Customer_Loyalty_Account upsert error:", putErr)
				}
				if putSetErr := redisx.SetJSON(putCtx, RedisClient, loyalty_account.RedisKey(), loyalty_account); putSetErr != nil {
					log.Println("redisx.SetJSON Customer_Loyalty_Account error:", putSetErr)
				}
				putCancel()
			}
		} else {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to credit loyalty account"
			response.ErrorDescription = "points to credit or outstanding fractions point must be greater than 0"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_AccountCreditPoints_log(*response)
			return
		}
	}
	//send sms earning notification

	//response.ClosureAvailablePoints = (loyalty_account.Awarded_Points + loyalty_account.Expired_Points) - loyalty_account.Redeemed_Points
	response.Closure_Loyalty_Level_Key = loyalty_account.Loyalty_Level_Key
	response.Closure_Loyalty_Account_Segment_Key = loyalty_account.Loyalty_Account_Segment_Key
	response.Closure_Awarded_Points = loyalty_account.Awarded_Points
	response.Closure_Redeemed_Points = loyalty_account.Redeemed_Points
	response.Closure_Available_Points = loyalty_account.Available_Points
	response.Closure_Outstanding_fraction_points = loyalty_account.Outstanding_fraction_points
	if sms && points > 0 {
		EarnNotiLog := NotificationLog{
			SourceAction:  "Earning",
			TransactionId: "",
			Medium:        "SMS",
			SourceAddress: sender,
			Destination:   loyalty_account.Key,
			Subject:       "Loyalty Earning",
			AddUser:       "SYSTEM",
			AddDate:       time.Now(),
		}
		if text != "" && sender != "" {
			text = strings.ReplaceAll(text, "{{EarnedPoints}}", fmt.Sprint(points))
			text = strings.ReplaceAll(text, "{{LoyaltyBalance}}", fmt.Sprint(loyalty_account.Available_Points))
			text = strings.ReplaceAll(text, "{{NewLevel}}", fmt.Sprint(loyalty_account.Loyalty_Level_Key))
			EarnNotiLog.Payload = text
			err := error(nil)
			if Configuration.Operation == "Angola" {
				err = SendSMS(sender, loyalty_account.Key, text)
			} else {
				err = Send_SMS(sender, loyalty_account.Key, text)
			}
			if err != nil {
				EarnNotiLog.Status = "Failed"
				EarnNotiLog.Error = err.Error()
			} else {
				EarnNotiLog.Status = "Successful"
			}
		} else {
			EarnNotiLog.Payload = text
			EarnNotiLog.Status = "Failed"
			EarnNotiLog.Error = "Undefined earning notification for transaction"
		}
		YYYY, MM, _, _, _, _, _ := GetTimeParts(time.Now())
		Db := Configuration.DB_Name_Loyalty + "_Logs_" + YYYY + MM
		{
			notiCtx, notiCancel := context.WithTimeout(context.Background(), 10*time.Second)
			_, notiErr := Uc.MongoClient.Mongo.Database(Db).Collection("Col_NotificationLog").InsertOne(notiCtx, EarnNotiLog)
			notiCancel()
			if notiErr != nil {
				log.Println("Error in Write earning Notification Logs:", notiErr, " (", EarnNotiLog, ")")
			}
		}
	}
	//successful reply
	response.Status = "successful"
	response.StatusCode = http.StatusOK
	response.StatusDescription = ""
	response.ErrorDescription = ""
	response.StatusDate = time.Now()
	response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
	Uc.Write_Loyalty_AccountCreditPoints_log(*response)

	PointsCreditedCount.With(prometheus.Labels{"EventSource": request.EventSource, "EventType": request.EventType, "Level": loyalty_account.Loyalty_Level_Key}).Inc()
	PointsCredited.With(prometheus.Labels{"EventSource": request.EventSource, "EventType": request.EventType, "Level": loyalty_account.Loyalty_Level_Key}).Add(points)
}

func (Uc *UserControl) Loyalty_AccountDebitPoints(request_header *Request_Header, request Loyalty_AccountDebitPoints_Request, response *Loyalty_AccountDebitPoints_log, notRedemption ...bool) {
	response.ReceiveDate = time.Now()
	//fill the request header info
	response.SourceIP = request_header.SourceIP
	response.SourceApp = request_header.SourceApp
	response.AppLogin = request_header.AppLogin
	response.AppVersion = request_header.AppVersion
	response.GPSLocation = request_header.GPSLocation
	response.GSMLocation = request_header.GSMLocation
	response.MSISDN = request.MSISDN
	request.MSISDN = Normalize_International_MSISDN(request.MSISDN)
	if request.MSISDN == "" {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "incorrect msisdn"
		response.ErrorDescription = "incorrect msisdn"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_AccountDebitPoints_log(*response)
		return
	}
	//fill the request info
	response.MSISDN = request.MSISDN
	response.Debit_Amount = request.Debit_Amount
	response.Debit_Reason = request.Debit_Reason
	response.ReceiveDate = time.Now()
	response.Redemption_Type = request.Redemption_Type //Airtime, Bundle, MobileMoney, SpinAndWin
	response.Redemption_Bunlde_Id = request.Redemption_Bunlde_Id
	response.Redemption_Amount = request.Redemption_Amount
	if response.Debit_Amount < 0 {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "invalid amount"
		response.ErrorDescription = "amount must be greater than 0"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_AccountDebitPoints_log(*response)
		return
	}
	//get loyalty account detail
	loyalty_Account, loyaltyAccErr := redisx.GetJSON[Customer_Loyalty_Account](context.Background(), RedisClient, Customer_Loyalty_Account{Key: request.MSISDN}.RedisKey())
	if redisx.IsNil(loyaltyAccErr) {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "loyalty account not found"
		response.ErrorDescription = "loyalty account not found"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_AccountDebitPoints_log(*response)
		return
	}
	if loyaltyAccErr != nil {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "error in loyalty account type assertion"
		response.ErrorDescription = "error in loyalty account type assertion"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_AccountDebitPoints_log(*response)
		return
	}
	response.Customer_Id = loyalty_Account.Customer_Id
	response.Account_Status = loyalty_Account.Account_Status
	response.Loyalty_Level_Key = loyalty_Account.Loyalty_Level_Key
	response.Loyalty_Account_Segment_Key = loyalty_Account.Loyalty_Account_Segment_Key
	response.Opening_Awarded_Points = loyalty_Account.Awarded_Points
	response.Opening_Redeemed_Points = loyalty_Account.Redeemed_Points
	response.Opening_Available_Points = loyalty_Account.Available_Points
	//check exclusion list
	if loyalty_Account.Opt_Status != "OptedIn" && Configuration.ISLoyaltyOptIn {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "customer is opted out"
		response.ErrorDescription = "customer is opted out"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_AccountDebitPoints_log(*response)
		return
	}
	{
		chkCtx, chkCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, chkErr := redisx.GetJSON[Customer_Exclusion](chkCtx, RedisClient, Customer_Exclusion{Key: loyalty_Account.Key}.RedisKey())
		chkCancel()
		if chkErr == nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "customer is included in the exclusion list"
			response.ErrorDescription = "customer is included in the exclusion list"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_AccountDebitPoints_log(*response)
			return
		}
	}
	//check COS exclusion list
	{
		chkCtx, chkCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, chkErr := redisx.GetJSON[Customer_COS_Exclusion](chkCtx, RedisClient, Customer_COS_Exclusion{Key: loyalty_Account.COS}.RedisKey())
		chkCancel()
		if chkErr == nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "customer is included in the cos exclusion list"
			response.ErrorDescription = "customer is included in the cos exclusion list"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_AccountDebitPoints_log(*response)
			return
		}
	}
	//check if available balance is enough
	if response.Opening_Available_Points < request.Debit_Amount {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "no enough points"
		response.ErrorDescription = "no enough points"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_AccountDebitPoints_log(*response)
		return
	}
	notRedemptionValue := false
	if len(notRedemption) > 0 {
		notRedemptionValue = notRedemption[0]
	}
	if !notRedemptionValue {
		resetDailyWeeklyCountersIfNeeded(&loyalty_Account)
		loyalty_governance, lgErr := redisx.GetJSON[Loyalty_Governance](context.Background(), RedisClient, Loyalty_Governance{Key: LOYALTY_GOVERNANCE_KEY}.RedisKey())
		if redisx.IsNil(lgErr) {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to debit loyalty account"
			response.ErrorDescription = "loyalty governance entry not found"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_AccountDebitPoints_log(*response)
			return
		}
		if lgErr != nil {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to debit loyalty account"
			response.ErrorDescription = "loyalty governance type assertion issue"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_AccountDebitPoints_log(*response)
			return
		}
		if loyalty_governance.DailyPointsRedemptionLimit > 0 && (loyalty_Account.Daily_Redeemed_Points+request.Debit_Amount) > loyalty_governance.DailyPointsRedemptionLimit {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "daily redemption points limit exceeded"
			response.ErrorDescription = "daily redemption points limit exceeded"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_AccountDebitPoints_log(*response)
			return
		}
		if loyalty_governance.WeeklyPointsRedemptionLimit > 0 && (loyalty_Account.Weekly_Redeemed_Points+request.Debit_Amount) > loyalty_governance.WeeklyPointsRedemptionLimit {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "weekly redemption points limit exceeded"
			response.ErrorDescription = "weekly redemption points limit exceeded"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_AccountDebitPoints_log(*response)
			return
		}
		if loyalty_governance.DailyRedemptionAttemptLimit > 0 && (loyalty_Account.Daily_Redemption_Attempts+1) > loyalty_governance.DailyRedemptionAttemptLimit {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "daily redemption attempt limit exceeded"
			response.ErrorDescription = "daily redemption attempt limit exceeded"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_AccountDebitPoints_log(*response)
			return
		}
		if loyalty_governance.WeeklyRedemptionAttemptLimit > 0 && (loyalty_Account.Weekly_Redemption_Attempts+1) > loyalty_governance.WeeklyRedemptionAttemptLimit {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "weekly redemption attempt limit exceeded"
			response.ErrorDescription = "weekly redemption attempt limit exceeded"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_AccountDebitPoints_log(*response)
			return
		}
	}
	//debit the account
	if notRedemptionValue {
		loyalty_Account.Awarded_Points = loyalty_Account.Awarded_Points - request.Debit_Amount
	} else {
		loyalty_Account.Redeemed_Points = loyalty_Account.Redeemed_Points + request.Debit_Amount
		loyalty_Account.Last_Redeem_Date = time.Now()
	}
	loyalty_Account.Available_Points = loyalty_Account.Awarded_Points - loyalty_Account.Redeemed_Points //(Awarded_Points + Expired_Points) - Redeemed_Points
	//update Monthly points
	start_date := time.Date(2025, 04, 01, 00, 00, 59, 0, time.UTC)
	end_date := start_date.AddDate(50, 0, 0)
	Amount_to_debit := request.Debit_Amount
	for d := start_date; d.After(end_date) == false; d = d.AddDate(0, 1, 0) {
		if d.Before(time.Now().AddDate(0, 0, 1)) {
			//fmt.Println(d.Format("2006-01-02"))
			YYYY, MM, _, _, _, _, _ := GetTimeParts(d)
			//prepare the monthly points detail
			PointsDetail, pdErr := redisx.GetJSON[Customer_Loyalty_Account_Points_Detail](context.Background(), RedisClient, Customer_Loyalty_Account_Points_Detail{Key: request.MSISDN + "|" + YYYY + MM}.RedisKey())
			if pdErr != nil && !redisx.IsNil(pdErr) {
				response.Status = "failed"
				response.StatusCode = http.StatusBadRequest
				response.StatusDescription = "failed to credit loyalty account"
				response.ErrorDescription = "issue with Customer_Loyalty_Account_Points_Detail type assertion"
				response.StatusDate = time.Now()
				response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
				Uc.Write_Loyalty_AccountDebitPoints_log(*response)
				return
			}
			if pdErr == nil {
				if PointsDetail.Available_Points > 0 {
					if PointsDetail.Available_Points >= Amount_to_debit {
						//full amount available
						if notRedemptionValue {
							PointsDetail.Awarded_Points = PointsDetail.Awarded_Points - Amount_to_debit
						} else {
							PointsDetail.Redeemed_Points = PointsDetail.Redeemed_Points + Amount_to_debit
							PointsDetail.Last_Redeem_Date = time.Now()
						}
						PointsDetail.Available_Points = PointsDetail.Awarded_Points - PointsDetail.Redeemed_Points //(Awarded_Points + Expired_Points) - Redeemed_Points
						Amount_to_debit = 0
					} else {
						//partial amount available
						partial_debit_amount := PointsDetail.Available_Points
						if notRedemptionValue {
							PointsDetail.Awarded_Points = PointsDetail.Awarded_Points - partial_debit_amount
						} else {
							PointsDetail.Redeemed_Points = PointsDetail.Redeemed_Points + partial_debit_amount
							PointsDetail.Last_Redeem_Date = time.Now()
						}
						PointsDetail.Available_Points = 0
						Amount_to_debit = Amount_to_debit - partial_debit_amount
					}
					{
						putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
						if _, putErr := Mdb_Customer_Loyalty_Account_Points_Detail.Coll.UpdateOne(putCtx, bson.M{"Key": PointsDetail.Key}, bson.M{"$set": PointsDetail}, options.UpdateOne().SetUpsert(true)); putErr != nil {
							log.Println("Mdb_Customer_Loyalty_Account_Points_Detail upsert error:", putErr)
						}
						if putSetErr := redisx.SetJSON(putCtx, RedisClient, PointsDetail.RedisKey(), PointsDetail); putSetErr != nil {
							log.Println("redisx.SetJSON Customer_Loyalty_Account_Points_Detail error:", putSetErr)
						}
						putCancel()
					}
					if Amount_to_debit == 0 {
						if !notRedemptionValue {
							loyalty_Account.Daily_Redeemed_Points += request.Debit_Amount
							loyalty_Account.Weekly_Redeemed_Points += request.Debit_Amount
							loyalty_Account.Daily_Redemption_Attempts++
							loyalty_Account.Weekly_Redemption_Attempts++
						}
						break
					}
				}
			} else {
				continue
			}
		} else {
			break
		}
	}
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Customer_Loyalty_Account.Coll.UpdateOne(putCtx, bson.M{"Key": loyalty_Account.Key}, bson.M{"$set": loyalty_Account}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Customer_Loyalty_Account upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, loyalty_Account.RedisKey(), loyalty_Account); putSetErr != nil {
			log.Println("redisx.SetJSON Customer_Loyalty_Account error:", putSetErr)
		}
		putCancel()
	}
	response.Closure_Awarded_Points = loyalty_Account.Awarded_Points
	response.Closure_Redeemed_Points = loyalty_Account.Redeemed_Points
	response.Closure_Available_Points = loyalty_Account.Available_Points
	//update goveranance
	Uc.Loyalty_Governance_Redeem_Points_Debit(request.Debit_Amount, notRedemptionValue)
	//successful reply
	response.Status = "successful"
	response.StatusCode = http.StatusOK
	response.StatusDescription = ""
	response.ErrorDescription = ""
	response.StatusDate = time.Now()
	response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
	Uc.Write_Loyalty_AccountDebitPoints_log(*response)
	PointsDebitedCount.With(prometheus.Labels{"EventSource": response.SourceApp, "DebitType": response.Redemption_Type, "Level": loyalty_Account.Loyalty_Level_Key}).Inc()
	PointsDebited.With(prometheus.Labels{"EventSource": response.SourceApp, "DebitType": response.Redemption_Type, "Level": loyalty_Account.Loyalty_Level_Key}).Add(request.Debit_Amount)

}

func (Uc *UserControl) Customer_Loyalty_OptRequest(request_header *Request_Header, request Loyalty_Opt_Request, response *Loyalty_Status_log) {
	request.MSISDN = Normalize_International_MSISDN(request.MSISDN)
	if request.MSISDN == "" {
		response.Request_Status = "failed"
		response.Request_StatusCode = http.StatusBadRequest
		response.StatusDescription = "invalid msisdn"
		response.ErrorDescription = "invalid msisdn"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.StatusDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_Status_log(*response)
		return
	}
	response.StatusDate = time.Now()
	//fill the request header info
	response.SourceIP = request_header.SourceIP
	response.SourceApp = request.EventSource
	response.AppLogin = request_header.AppLogin
	response.AppVersion = request_header.AppVersion
	response.MSISDN = request.MSISDN
	// request.MSISDN = Normalize_International_MSISDN(request.MSISDN)
	if request.Opt_Status == "" {
		response.Request_Status = "failed"
		response.Request_StatusCode = http.StatusBadRequest
		response.StatusDescription = "invalid opt status"
		response.ErrorDescription = "invalid opt status"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.StatusDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_Status_log(*response)
		return
	}
	//fill the request info
	response.MSISDN = request.MSISDN
	response.Opt_Status = request.Opt_Status

	//validate loyalty account
	loyalty_account, loyaltyAccErr := redisx.GetJSON[Customer_Loyalty_Account](context.Background(), RedisClient, Customer_Loyalty_Account{Key: request.MSISDN}.RedisKey())
	if redisx.IsNil(loyaltyAccErr) {
		response.Request_Status = "failed"
		response.Request_StatusCode = http.StatusBadRequest
		response.StatusDescription = "failed to get loyalty account"
		response.ErrorDescription = "loyalty account does not exist"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.StatusDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_Status_log(*response)
		return
	}
	if loyaltyAccErr != nil {
		response.Request_Status = "failed"
		response.Request_StatusCode = http.StatusBadRequest
		response.StatusDescription = "failed to get loyalty account"
		response.ErrorDescription = "type assertion issue with Customer_Loyalty_Account"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.StatusDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_Status_log(*response)
		return
	}

	if loyalty_account.Opt_Status == request.Opt_Status || (loyalty_account.Opt_Status == "OptedOutExpired" && request.Opt_Status == "OptedOut") {
		response.Request_Status = "failed"
		response.Request_StatusCode = http.StatusBadRequest
		status := "opted out"
		if request.Opt_Status == "OptedIn" {
			status = "opted in"
		}
		response.StatusDescription = fmt.Sprint("customer is already ", status)
		response.ErrorDescription = fmt.Sprint("customer is already ", status)
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.StatusDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_Status_log(*response)
		return
	}
	//check exclusion list
	{
		chkCtx, chkCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, chkErr := redisx.GetJSON[Customer_Exclusion](chkCtx, RedisClient, Customer_Exclusion{Key: loyalty_account.Key}.RedisKey())
		chkCancel()
		if chkErr == nil {
			response.Request_Status = "failed"
			response.Request_StatusCode = http.StatusBadRequest
			response.StatusDescription = "customer is included in the exclusion list"
			response.ErrorDescription = "customer is included in the exclusion list"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.StatusDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Status_log(*response)
			return
		}
	}
	//check COS exclusion list
	{
		chkCtx, chkCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, chkErr := redisx.GetJSON[Customer_COS_Exclusion](chkCtx, RedisClient, Customer_COS_Exclusion{Key: loyalty_account.COS}.RedisKey())
		chkCancel()
		if chkErr == nil {
			response.Request_Status = "failed"
			response.Request_StatusCode = http.StatusBadRequest
			response.StatusDescription = "customer is included in the cos exclusion list"
			response.ErrorDescription = "customer is included in the cos exclusion list"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.StatusDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Status_log(*response)
			return
		}
	}

	if loyalty_account.Loyalty_Account_Segment_Key == "" {
		response.Request_Status = "failed"
		response.Request_StatusCode = http.StatusBadRequest
		response.StatusDescription = "loyalty account segment is not assigned"
		response.ErrorDescription = "loyalty account segment is not assigned"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.StatusDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_Status_log(*response)
		return
	}
	if loyalty_account.Loyalty_Level_Key == "" {
		response.Request_Status = "failed"
		response.Request_StatusCode = http.StatusBadRequest
		response.StatusDescription = "loyalty account level is not assigned"
		response.ErrorDescription = "loyalty account level is not assigned"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.StatusDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_Status_log(*response)
		return
	}
	//validate loyalty level
	_, loyaltyLevelErr := redisx.GetJSON[Loyalty_Level](context.Background(), RedisClient, Loyalty_Level{Key: loyalty_account.Loyalty_Level_Key}.RedisKey())
	if redisx.IsNil(loyaltyLevelErr) {
		response.Request_Status = "failed"
		response.Request_StatusCode = http.StatusBadRequest
		response.StatusDescription = "failed to get loyalty account level"
		response.ErrorDescription = "loyalty account level is not defined"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.StatusDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_Status_log(*response)
		return
	}

	//validate the loyalty plan
	plan, planErr := redisx.GetJSON[Loyalty_Plan](context.Background(), RedisClient, Loyalty_Plan{Key: loyalty_account.Loyalty_Level_Key + "|" + loyalty_account.Loyalty_Account_Segment_Key}.RedisKey())
	if redisx.IsNil(planErr) {
		response.Request_Status = "failed"
		response.Request_StatusCode = http.StatusBadRequest
		response.StatusDescription = "failed to get loyalty plan"
		response.ErrorDescription = "loyalty plan does not exist"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.StatusDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_Status_log(*response)
		return
	}
	if planErr != nil {
		response.Request_Status = "failed"
		response.Request_StatusCode = http.StatusBadRequest
		response.StatusDescription = "failed to get loyalty plan"
		response.ErrorDescription = "type assertion issue with Loyalty_Plan"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.StatusDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_Status_log(*response)
		return
	}
	//validate earning rules
	if plan.Earning_Rules_Key == "" {
		response.Request_Status = "failed"
		response.Request_StatusCode = http.StatusBadRequest
		response.StatusDescription = "failed to get earning rules"
		response.ErrorDescription = "earning rules are not defined"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.StatusDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_Status_log(*response)
		return
	}
	loyalty_account.Opt_Status = request.Opt_Status
	loyalty_account.Last_Opt_Status_Date = time.Now()
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Customer_Loyalty_Account.Coll.UpdateOne(putCtx, bson.M{"Key": loyalty_account.Key}, bson.M{"$set": loyalty_account}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Customer_Loyalty_Account upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, loyalty_account.RedisKey(), loyalty_account); putSetErr != nil {
			log.Println("redisx.SetJSON Customer_Loyalty_Account error:", putSetErr)
		}
		putCancel()
	}

	if loyalty_account.First_Opt_In_Status_Date.IsZero() && request.Opt_Status == "OptedIn" {
		loyalty_account.First_Opt_In_Status_Date = time.Now()
		{
			putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
			if _, putErr := Mdb_Customer_Loyalty_Account.Coll.UpdateOne(putCtx, bson.M{"Key": loyalty_account.Key}, bson.M{"$set": loyalty_account}, options.UpdateOne().SetUpsert(true)); putErr != nil {
				log.Println("Mdb_Customer_Loyalty_Account upsert error:", putErr)
			}
			if putSetErr := redisx.SetJSON(putCtx, RedisClient, loyalty_account.RedisKey(), loyalty_account); putSetErr != nil {
				log.Println("redisx.SetJSON Customer_Loyalty_Account error:", putSetErr)
			}
			putCancel()
		}
		var loyalty_AccountCreditPoints_log Loyalty_AccountCreditPoints_log
		var loyalty_AccountCreditPoints_Request Loyalty_AccountCreditPoints_Request
		loyalty_AccountCreditPoints_Request.MSISDN = loyalty_account.Key
		loyalty_AccountCreditPoints_Request.EventSource = "First_Opt_In"
		loyalty_AccountCreditPoints_Request.EventType = "NewJoining"
		loyalty_AccountCreditPoints_Request.EventDetail = ""
		loyalty_AccountCreditPoints_Request.EventAmount = 0
		loyalty_AccountCreditPoints_Request.EventDescription = ""
		var credit_request_header Request_Header
		credit_request_header.SourceIP = "127.0.0.1"
		credit_request_header.SourceApp = request.EventSource
		credit_request_header.AppLogin = request_header.AppLogin
		Uc.Loyalty_AccountCreditPoints(&credit_request_header, loyalty_AccountCreditPoints_Request, &loyalty_AccountCreditPoints_log)

		Earningrecord, err := Uc.Customer_Loyalty_Account_GetEarning_Rule(loyalty_account.Key)
		if err != nil {
			log.Println("failed to get data")
			return
		}
		loyalty_account, loyaltyAccErr2 := redisx.GetJSON[Customer_Loyalty_Account](context.Background(), RedisClient, Customer_Loyalty_Account{Key: request.MSISDN}.RedisKey())
		if redisx.IsNil(loyaltyAccErr2) {
			response.Request_Status = "failed"
			response.Request_StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to get loyalty account"
			response.ErrorDescription = "loyalty account does not exist"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.StatusDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Status_log(*response)
			return
		}
		if loyaltyAccErr2 != nil {
			response.Request_Status = "failed"
			response.Request_StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to get loyalty account"
			response.ErrorDescription = "type assertion issue with Customer_Loyalty_Account"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.StatusDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Status_log(*response)
			return
		}
		if Earningrecord.Welcome_Notification {
			WelcomeNotiLog := NotificationLog{
				SourceAction:  "Welcome",
				TransactionId: "",
				Medium:        "SMS",
				SourceAddress: Earningrecord.Welcome_Notification_Sender,
				Destination:   loyalty_account.Key,
				Subject:       "Welcome",
				AddUser:       "SYSTEM",
				AddDate:       time.Now(),
			}
			Welcome_Noti_Text := ""
			Welcome_Noti_Text = Earningrecord.Welcome_Notification_Text
			if Welcome_Noti_Text != "" {
				Welcome_Noti_Text = strings.ReplaceAll(Welcome_Noti_Text, "{{WelcomePoints}}", fmt.Sprint(loyalty_account.Awarded_Points))
				Welcome_Noti_Text = strings.ReplaceAll(Welcome_Noti_Text, "{{LoyaltyBalance}}", fmt.Sprint(loyalty_account.Available_Points))
				Welcome_Noti_Text = strings.ReplaceAll(Welcome_Noti_Text, "{{NewLevel}}", fmt.Sprint(loyalty_account.Loyalty_Level_Key))
				WelcomeNotiLog.Payload = Welcome_Noti_Text
				err := error(nil)
				if Configuration.Operation == "Angola" {
					err = SendSMS(Earningrecord.Welcome_Notification_Sender, loyalty_account.Key, Welcome_Noti_Text)
				} else {
					err = Send_SMS(Earningrecord.Welcome_Notification_Sender, loyalty_account.Key, Welcome_Noti_Text)
				}
				if err != nil {
					WelcomeNotiLog.Status = "Failed"
					WelcomeNotiLog.Error = err.Error()
				} else {
					WelcomeNotiLog.Status = "Successful"
				}
			} else {
				WelcomeNotiLog.Payload = Welcome_Noti_Text
				WelcomeNotiLog.Status = "Failed"
				WelcomeNotiLog.Error = "Undefined welcome notification for transaction"
			}
			YYYY, MM, _, _, _, _, _ := GetTimeParts(time.Now())
			Db := Configuration.DB_Name_Loyalty + "_Logs_" + YYYY + MM
			{
				notiCtx, notiCancel := context.WithTimeout(context.Background(), 10*time.Second)
				_, notiErr := Uc.MongoClient.Mongo.Database(Db).Collection("Col_NotificationLog").InsertOne(notiCtx, WelcomeNotiLog)
				notiCancel()
				if notiErr != nil {
					log.Println("Error in Write welcome Notification Logs:", notiErr, " (", WelcomeNotiLog, ")")
				}
			}
		}
	} else if !loyalty_account.First_Opt_In_Status_Date.IsZero() && request.Opt_Status == "OptedIn" {
		Earningrecord, earningErr := Uc.Customer_Loyalty_Account_GetEarning_Rule(loyalty_account.Key)
		if earningErr != nil {
			log.Println("failed to get earning rule for rejoiner notification")
		} else if Earningrecord.Rejoiner_Notification {
			RejoinerNotiLog := NotificationLog{
				SourceAction:  "Rejoiner",
				TransactionId: "",
				Medium:        "SMS",
				SourceAddress: Earningrecord.Rejoiner_Notification_Sender,
				Destination:   loyalty_account.Key,
				Subject:       "Rejoiner",
				AddUser:       "SYSTEM",
				AddDate:       time.Now(),
			}
			Rejoiner_Noti_Text := Earningrecord.Rejoiner_Notification_Text
			if Rejoiner_Noti_Text != "" {
				Rejoiner_Noti_Text = strings.ReplaceAll(Rejoiner_Noti_Text, "{{LoyaltyBalance}}", fmt.Sprint(loyalty_account.Available_Points))
				Rejoiner_Noti_Text = strings.ReplaceAll(Rejoiner_Noti_Text, "{{NewLevel}}", fmt.Sprint(loyalty_account.Loyalty_Level_Key))
				Rejoiner_Noti_Text = strings.ReplaceAll(Rejoiner_Noti_Text, "{{PreviousLevel}}", fmt.Sprint(loyalty_account.Previous_Loyalty_Level_Key))
				Rejoiner_Noti_Text = strings.ReplaceAll(Rejoiner_Noti_Text, "{{LevelChangeDirection}}", fmt.Sprint(loyalty_account.Loyalty_Level_Direction))
				RejoinerNotiLog.Payload = Rejoiner_Noti_Text
				err := error(nil)
				if Configuration.Operation == "Angola" {
					err = SendSMS(Earningrecord.Rejoiner_Notification_Sender, loyalty_account.Key, Rejoiner_Noti_Text)
				} else {
					err = Send_SMS(Earningrecord.Rejoiner_Notification_Sender, loyalty_account.Key, Rejoiner_Noti_Text)
				}
				if err != nil {
					RejoinerNotiLog.Status = "Failed"
					RejoinerNotiLog.Error = err.Error()
				} else {
					RejoinerNotiLog.Status = "Successful"
				}
			} else {
				RejoinerNotiLog.Payload = Rejoiner_Noti_Text
				RejoinerNotiLog.Status = "Failed"
				RejoinerNotiLog.Error = "Undefined rejoiner notification for transaction"
			}
			YYYY, MM, _, _, _, _, _ := GetTimeParts(time.Now())
			Db := Configuration.DB_Name_Loyalty + "_Logs_" + YYYY + MM
			{
				notiCtx, notiCancel := context.WithTimeout(context.Background(), 10*time.Second)
				_, notiErr := Uc.MongoClient.Mongo.Database(Db).Collection("Col_NotificationLog").InsertOne(notiCtx, RejoinerNotiLog)
				notiCancel()
				if notiErr != nil {
					log.Println("Error in Write rejoiner Notification Logs:", notiErr, " (", RejoinerNotiLog, ")")
				}
			}
		}
	}
	//response.ClosureAvailablePoints = (loyalty_account.Awarded_Points + loyalty_account.Expired_Points) - loyalty_account.Redeemed_Points

	//successful reply
	response.Request_Status = "successful"
	response.Request_StatusCode = http.StatusOK
	response.StatusDescription = ""
	response.ErrorDescription = ""
	response.StatusDate = time.Now()
	response.E2E_Elapsedtime = (time.Since(response.StatusDate).Nanoseconds()) / 1000000
	Uc.Write_Loyalty_Status_log(*response)
}

func Calculate_Loyalty_Points(rules Loyalty_Point_Earning_Rules, award_request Loyalty_AccountCreditPoints_Request, current_outstanding_points float64) (points, outstanding_points float64, sms_notificatoin bool, notification_sender string, notification_text string) {
	switch award_request.EventSource {
	case "DWH_Import":
		switch award_request.EventType {
		case "NewJoining":
			return rules.Welcome_Points, current_outstanding_points, false, "", ""
		default:
			return 0, current_outstanding_points, false, "", ""
		}
	case "First_Opt_In":
		return rules.Welcome_Points, current_outstanding_points, false, "", ""
	case "IN_feed":
		switch award_request.EventType {
		case "NewJoining":
			return rules.Welcome_Points, current_outstanding_points, false, "", ""
		case "SSR_3": //scratch card recharge
			if rules.GSM_SC_Airtime_Award_Type == "Transaction" {
				if rules.GSM_SC_Airtime_Points > 0 {
					// return rules.GSM_SC_Airtime_Points, current_outstanding_points
					flt_points := rules.GSM_SC_Airtime_Points + current_outstanding_points
					int_points := int(flt_points)
					outstanding_points = flt_points - float64(int_points)
					return float64(int_points), outstanding_points, rules.GSM_SC_Notification, rules.GSM_SC_Notification_Sender, rules.GSM_SC_Notification_Text
				}
			} else if rules.GSM_SC_Airtime_Award_Type == "Amount" {
				if rules.GSM_SC_Airtime_Amount > 0 && award_request.EventAmount > 0 {
					flt_fractions := award_request.EventAmount / rules.GSM_SC_Airtime_Amount
					flt_points := (flt_fractions * rules.GSM_SC_Airtime_Points) + current_outstanding_points
					int_points := int(flt_points)
					outstanding_points = flt_points - float64(int_points)
					return float64(int_points), outstanding_points, rules.GSM_SC_Notification, rules.GSM_SC_Notification_Sender, rules.GSM_SC_Notification_Text
				}
			}
			return 0, current_outstanding_points, false, "", ""
		case "SSR_211": //scratch card recharge
			if rules.GSM_EVC_Bundle_Award_Type == "Transaction" {
				if rules.GSM_EVC_Bundle_Points > 0 {
					// return rules.GSM_EVC_Bundle_Points, current_outstanding_points
					flt_points := rules.GSM_EVC_Bundle_Points + current_outstanding_points
					int_points := int(flt_points)
					outstanding_points = flt_points - float64(int_points)
					return float64(int_points), outstanding_points, rules.GSM_EVC_Bundle_Notification, rules.GSM_EVC_Bundle_Notification_Sender, rules.GSM_EVC_Bundle_Notification_Text
				}
			} else if rules.GSM_EVC_Bundle_Award_Type == "Amount" {
				if rules.GSM_EVC_Bundle_Amount > 0 && award_request.EventAmount > 0 {
					flt_fractions := award_request.EventAmount / rules.GSM_EVC_Bundle_Amount
					flt_points := (flt_fractions * rules.GSM_EVC_Bundle_Points) + current_outstanding_points
					int_points := int(flt_points)
					outstanding_points = flt_points - float64(int_points)
					return float64(int_points), outstanding_points, rules.GSM_EVC_Bundle_Notification, rules.GSM_EVC_Bundle_Notification_Sender, rules.GSM_EVC_Bundle_Notification_Text
				}
			}
			return 0, current_outstanding_points, false, "", ""
		case "SSR_97": //EVC recharge
			if rules.GSM_EVC_Airtime_Award_Type == "Transaction" {
				if rules.GSM_EVC_Airtime_Points > 0 {
					// return rules.GSM_EVC_Airtime_Points, current_outstanding_points
					flt_points := rules.GSM_EVC_Airtime_Points + current_outstanding_points
					int_points := int(flt_points)
					outstanding_points = flt_points - float64(int_points)
					return float64(int_points), outstanding_points, rules.GSM_EVC_Airtime_Notification, rules.GSM_EVC_Airtime_Notification_Sender, rules.GSM_EVC_Airtime_Notification_Text
				}
			} else if rules.GSM_EVC_Airtime_Award_Type == "Amount" {
				if rules.GSM_EVC_Airtime_Amount > 0 && award_request.EventAmount > 0 {
					flt_fractions := award_request.EventAmount / rules.GSM_EVC_Airtime_Amount
					flt_points := (flt_fractions * rules.GSM_EVC_Airtime_Points) + current_outstanding_points
					int_points := int(flt_points)
					outstanding_points = flt_points - float64(int_points)
					return float64(int_points), outstanding_points, rules.GSM_EVC_Airtime_Notification, rules.GSM_EVC_Airtime_Notification_Sender, rules.GSM_EVC_Airtime_Notification_Text
				}
			}
			return 0, current_outstanding_points, false, "", ""
		case "MM_SSR_97": //mobile money airtime recharge
			if rules.MM_Airtime_Award_Type == "Transaction" {
				if rules.MM_Airtime_Points > 0 {
					// return rules.MM_Airtime_Points, current_outstanding_points
					flt_points := rules.MM_Airtime_Points + current_outstanding_points
					int_points := int(flt_points)
					outstanding_points = flt_points - float64(int_points)
					return float64(int_points), outstanding_points, rules.MM_Airtime_Notification, rules.MM_Airtime_Notification_Sender, rules.MM_Airtime_Notification_Text
				}
			} else if rules.MM_Airtime_Award_Type == "Amount" {
				if rules.MM_Airtime_Amount > 0 && award_request.EventAmount > 0 {
					flt_fractions := award_request.EventAmount / rules.MM_Airtime_Amount
					flt_points := (flt_fractions * rules.MM_Airtime_Points) + current_outstanding_points
					int_points := int(flt_points)
					outstanding_points = flt_points - float64(int_points)
					return float64(int_points), outstanding_points, rules.MM_Airtime_Notification, rules.MM_Airtime_Notification_Sender, rules.MM_Airtime_Notification_Text
				}
			}
			return 0, current_outstanding_points, false, "", ""
		case "MM_Bundles_Recharge": //mobile money bundle recharge
			if rules.MM_Bundle_Award_Type == "Transaction" {
				if rules.MM_Bundle_Points > 0 {
					// return rules.MM_Bundle_Points, current_outstanding_points
					flt_points := rules.MM_Bundle_Points + current_outstanding_points
					int_points := int(flt_points)
					outstanding_points = flt_points - float64(int_points)
					return float64(int_points), outstanding_points, rules.MM_Bundle_Notification, rules.MM_Bundle_Notification_Sender, rules.MM_Bundle_Notification_Text
				}
			} else if rules.MM_Bundle_Award_Type == "Amount" {
				if rules.MM_Bundle_Amount > 0 && award_request.EventAmount > 0 {
					flt_fractions := award_request.EventAmount / rules.MM_Bundle_Amount
					flt_points := (flt_fractions * rules.MM_Bundle_Points) + current_outstanding_points
					int_points := int(flt_points)
					outstanding_points = flt_points - float64(int_points)
					return float64(int_points), outstanding_points, rules.MM_Bundle_Notification, rules.MM_Bundle_Notification_Sender, rules.MM_Bundle_Notification_Text
				}
			}
			return 0, current_outstanding_points, false, "", ""
		default:
			//award points on main GSM balance consumption based on amount
			if award_request.EventAmount > 0 {
				if rules.MainGSMBalance_Amount > 0 {
					flt_fractions := award_request.EventAmount / rules.MainGSMBalance_Amount
					flt_points := (flt_fractions * rules.MainGSMBalance_Points) + current_outstanding_points
					int_points := int(flt_points)
					outstanding_points = flt_points - float64(int_points)
					return float64(int_points), outstanding_points, rules.MainGSM_Notification, rules.MainGSM_Notification_Sender, rules.MainGSM_Notification_Text
				} else {
					return 0, current_outstanding_points, false, "", ""
				}
			} else { //to do: award points based transaction
				return 0, current_outstanding_points, false, "", ""
			}
		}
	case "WebPortal":
		switch award_request.EventType {
		case "NewJoining":
			return rules.Welcome_Points, current_outstanding_points, false, "", ""
		default:
			//award points based amount
			if award_request.EventAmount > 0 {
				if rules.MainGSMBalance_Amount > 0 {
					flt_fractions := award_request.EventAmount / rules.MainGSMBalance_Amount
					flt_points := (flt_fractions * rules.MainGSMBalance_Points) + current_outstanding_points
					int_points := int(flt_points)
					outstanding_points = flt_points - float64(int_points)
					return float64(int_points), outstanding_points, rules.MainGSM_Notification, rules.MainGSM_Notification_Sender, rules.MainGSM_Notification_Text
				} else {
					return 0, current_outstanding_points, false, "", ""
				}
			} else { //to do: award points based transaction
				return 0, current_outstanding_points, false, "", ""
			}
		}
	case "MobileMoney_feed":
		switch award_request.EventType {
		case "P2P":
			if rules.MM_P2P_Award_Type == "Transaction" {
				if rules.MM_P2P_Points > 0 {
					// return rules.MM_P2P_Points, current_outstanding_points
					flt_points := rules.MM_P2P_Points + current_outstanding_points
					int_points := int(flt_points)
					outstanding_points = flt_points - float64(int_points)
					return float64(int_points), outstanding_points, rules.MM_P2P_Notification, rules.MM_P2P_Notification_Sender, rules.MM_P2P_Notification_Text
				}
			} else if rules.MM_P2P_Award_Type == "Amount" {
				if rules.MM_P2P_Amount > 0 && award_request.EventAmount > 0 {
					flt_fractions := award_request.EventAmount / rules.MM_P2P_Amount
					flt_points := (flt_fractions * rules.MM_P2P_Points) + current_outstanding_points
					int_points := int(flt_points)
					outstanding_points = flt_points - float64(int_points)
					return float64(int_points), outstanding_points, rules.MM_P2P_Notification, rules.MM_P2P_Notification_Sender, rules.MM_P2P_Notification_Text
				}
			}
			return 0, current_outstanding_points, false, "", ""
		case "EMISP2POUT":
			if rules.MM_P2P_Award_Type == "Transaction" {
				if rules.MM_P2P_Points > 0 {
					// return rules.MM_P2P_Points, current_outstanding_points
					flt_points := rules.MM_P2P_Points + current_outstanding_points
					int_points := int(flt_points)
					outstanding_points = flt_points - float64(int_points)
					return float64(int_points), outstanding_points, rules.MM_P2P_Notification, rules.MM_P2P_Notification_Sender, rules.MM_P2P_Notification_Text
				}
			} else if rules.MM_P2P_Award_Type == "Amount" {
				if rules.MM_P2P_Amount > 0 && award_request.EventAmount > 0 {
					flt_fractions := award_request.EventAmount / rules.MM_P2P_Amount
					flt_points := (flt_fractions * rules.MM_P2P_Points) + current_outstanding_points
					int_points := int(flt_points)
					outstanding_points = flt_points - float64(int_points)
					return float64(int_points), outstanding_points, rules.MM_P2P_Notification, rules.MM_P2P_Notification_Sender, rules.MM_P2P_Notification_Text
				}
			}
			return 0, current_outstanding_points, false, "", ""
		case "EMISINTRA":
			if rules.MM_P2P_Award_Type == "Transaction" {
				if rules.MM_P2P_Points > 0 {
					// return rules.MM_P2P_Points, current_outstanding_points
					flt_points := rules.MM_P2P_Points + current_outstanding_points
					int_points := int(flt_points)
					outstanding_points = flt_points - float64(int_points)
					return float64(int_points), outstanding_points, rules.MM_P2P_Notification, rules.MM_P2P_Notification_Sender, rules.MM_P2P_Notification_Text
				}
			} else if rules.MM_P2P_Award_Type == "Amount" {
				if rules.MM_P2P_Amount > 0 && award_request.EventAmount > 0 {
					flt_fractions := award_request.EventAmount / rules.MM_P2P_Amount
					flt_points := (flt_fractions * rules.MM_P2P_Points) + current_outstanding_points
					int_points := int(flt_points)
					outstanding_points = flt_points - float64(int_points)
					return float64(int_points), outstanding_points, rules.MM_P2P_Notification, rules.MM_P2P_Notification_Sender, rules.MM_P2P_Notification_Text
				}
			}
			return 0, current_outstanding_points, false, "", ""
		case "CASHIN":
			if rules.MM_CASHIN_Award_Type == "Transaction" {
				if rules.MM_CASHIN_Points > 0 {
					// return rules.MM_CASHIN_Points, current_outstanding_points
					flt_points := rules.MM_CASHIN_Points + current_outstanding_points
					int_points := int(flt_points)
					outstanding_points = flt_points - float64(int_points)
					return float64(int_points), outstanding_points, rules.MM_CASHIN_Notification, rules.MM_CASHIN_Notification_Sender, rules.MM_CASHIN_Notification_Text
				}
			} else if rules.MM_CASHIN_Award_Type == "Amount" {
				if rules.MM_CASHIN_Amount > 0 && award_request.EventAmount > 0 {
					flt_fractions := award_request.EventAmount / rules.MM_CASHIN_Amount
					flt_points := (flt_fractions * rules.MM_CASHIN_Points) + current_outstanding_points
					int_points := int(flt_points)
					outstanding_points = flt_points - float64(int_points)
					return float64(int_points), outstanding_points, rules.MM_CASHIN_Notification, rules.MM_CASHIN_Notification_Sender, rules.MM_CASHIN_Notification_Text
				}
			}
			return 0, current_outstanding_points, false, "", ""
		case "CASHOUT":
			if rules.MM_CASHOUT_Award_Type == "Transaction" {
				if rules.MM_CASHOUT_Points > 0 {
					// return rules.MM_CASHOUT_Points, current_outstanding_points
					flt_points := rules.MM_CASHOUT_Points + current_outstanding_points
					int_points := int(flt_points)
					outstanding_points = flt_points - float64(int_points)
					return float64(int_points), outstanding_points, rules.MM_CASHOUT_Notification, rules.MM_CASHOUT_Notification_Sender, rules.MM_CASHOUT_Notification_Text
				}
			} else if rules.MM_CASHOUT_Award_Type == "Amount" {
				if rules.MM_CASHOUT_Amount > 0 && award_request.EventAmount > 0 {
					flt_fractions := award_request.EventAmount / rules.MM_CASHOUT_Amount
					flt_points := (flt_fractions * rules.MM_CASHOUT_Points) + current_outstanding_points
					int_points := int(flt_points)
					outstanding_points = flt_points - float64(int_points)
					return float64(int_points), outstanding_points, rules.MM_CASHOUT_Notification, rules.MM_CASHOUT_Notification_Sender, rules.MM_CASHOUT_Notification_Text
				}
			}
			return 0, current_outstanding_points, false, "", ""
		case "MERCHPAY":
			entry, entryErr := redisx.GetJSON[Loyalty_Point_Earning_Rules_Overwrite](context.Background(), RedisClient, Loyalty_Point_Earning_Rules_Overwrite{Key: rules.Key + "|MERCHPAY|" + award_request.EventDetailCode}.RedisKey())
			if entryErr == nil {
				if entry.Award_Type == "Transaction" {
					if entry.Points > 0 {
						// return entry.MM_MERCHPAY_Points, current_outstanding_points
						flt_points := entry.Points + current_outstanding_points
						int_points := int(flt_points)
						outstanding_points = flt_points - float64(int_points)
						return float64(int_points), outstanding_points, rules.Overwrite_MM_MERCHPAY_Notification, rules.Overwrite_MM_MERCHPAY_Notification_Sender, rules.Overwrite_MM_MERCHPAY_Notification_Text
					}
				} else if entry.Award_Type == "Amount" {
					if entry.Amount > 0 && award_request.EventAmount > 0 {
						flt_fractions := award_request.EventAmount / entry.Amount
						flt_points := (flt_fractions * entry.Points) + current_outstanding_points
						int_points := int(flt_points)
						outstanding_points = flt_points - float64(int_points)
						return float64(int_points), outstanding_points, rules.Overwrite_MM_MERCHPAY_Notification, rules.Overwrite_MM_MERCHPAY_Notification_Sender, rules.Overwrite_MM_MERCHPAY_Notification_Text
					}
				}
			} else {
				if rules.MM_MERCHPAY_Award_Type == "Transaction" {
					if rules.MM_MERCHPAY_Points > 0 {
						// return rules.MM_MERCHPAY_Points, current_outstanding_points
						flt_points := rules.MM_MERCHPAY_Points + current_outstanding_points
						int_points := int(flt_points)
						outstanding_points = flt_points - float64(int_points)
						return float64(int_points), outstanding_points, rules.MM_MERCHPAY_Notification, rules.MM_MERCHPAY_Notification_Sender, rules.MM_MERCHPAY_Notification_Text
					}
				} else if rules.MM_MERCHPAY_Award_Type == "Amount" {
					if rules.MM_MERCHPAY_Amount > 0 && award_request.EventAmount > 0 {
						flt_fractions := award_request.EventAmount / rules.MM_MERCHPAY_Amount
						flt_points := (flt_fractions * rules.MM_MERCHPAY_Points) + current_outstanding_points
						int_points := int(flt_points)
						outstanding_points = flt_points - float64(int_points)
						return float64(int_points), outstanding_points, rules.MM_MERCHPAY_Notification, rules.MM_MERCHPAY_Notification_Sender, rules.MM_MERCHPAY_Notification_Text
					}
				}
			}

			return 0, current_outstanding_points, false, "", ""
		case "BILLPAY":
			entry, entryErr := redisx.GetJSON[Loyalty_Point_Earning_Rules_Overwrite](context.Background(), RedisClient, Loyalty_Point_Earning_Rules_Overwrite{Key: rules.Key + "|BILLPAY|" + award_request.EventDetailCode}.RedisKey())
			if entryErr == nil {
				if entry.Award_Type == "Transaction" {
					if entry.Points > 0 {
						// return entry.MM_BILLPAY_Points, current_outstanding_points
						flt_points := entry.Points + current_outstanding_points
						int_points := int(flt_points)
						outstanding_points = flt_points - float64(int_points)
						return float64(int_points), outstanding_points, rules.Overwrite_MM_BILLPAY_Notification, rules.Overwrite_MM_BILLPAY_Notification_Sender, rules.Overwrite_MM_BILLPAY_Notification_Text
					}
				} else if entry.Award_Type == "Amount" {
					if entry.Amount > 0 && award_request.EventAmount > 0 {
						flt_fractions := award_request.EventAmount / rules.MM_BILLPAY_Amount
						flt_points := (flt_fractions * entry.Points) + current_outstanding_points
						int_points := int(flt_points)
						outstanding_points = flt_points - float64(int_points)
						return float64(int_points), outstanding_points, rules.Overwrite_MM_BILLPAY_Notification, rules.Overwrite_MM_BILLPAY_Notification_Sender, rules.Overwrite_MM_BILLPAY_Notification_Text
					}
				}
			} else {
				if rules.MM_BILLPAY_Award_Type == "Transaction" {
					if rules.MM_BILLPAY_Points > 0 {
						// return rules.MM_BILLPAY_Points, current_outstanding_points
						flt_points := rules.MM_BILLPAY_Points + current_outstanding_points
						int_points := int(flt_points)
						outstanding_points = flt_points - float64(int_points)
						return float64(int_points), outstanding_points, rules.MM_BILLPAY_Notification, rules.MM_BILLPAY_Notification_Sender, rules.MM_BILLPAY_Notification_Text
					}
				} else if rules.MM_BILLPAY_Award_Type == "Amount" {
					if rules.MM_BILLPAY_Amount > 0 && award_request.EventAmount > 0 {
						flt_fractions := award_request.EventAmount / rules.MM_BILLPAY_Amount
						flt_points := (flt_fractions * rules.MM_BILLPAY_Points) + current_outstanding_points
						int_points := int(flt_points)
						outstanding_points = flt_points - float64(int_points)
						return float64(int_points), outstanding_points, rules.MM_BILLPAY_Notification, rules.MM_BILLPAY_Notification_Sender, rules.MM_BILLPAY_Notification_Text
					}
				}
			}
			return 0, current_outstanding_points, false, "", ""
		case "RC": //self recharge
			if CheckMMEventDetailType(award_request.EventDetail) == Configuration.LoyaltyMMBundleCode {
				if rules.MM_RC_Bundle_Award_Type == "Transaction" {
					if rules.MM_RC_Bundle_Points > 0 {
						// return rules.MM_RC_Bundle_Points, current_outstanding_points
						flt_points := rules.MM_RC_Bundle_Points + current_outstanding_points
						int_points := int(flt_points)
						outstanding_points = flt_points - float64(int_points)
						return float64(int_points), outstanding_points, rules.MM_RC_Bundle_Notification, rules.MM_RC_Bundle_Notification_Sender, rules.MM_RC_Bundle_Notification_Text
					}
				} else if rules.MM_RC_Bundle_Award_Type == "Amount" {
					if rules.MM_RC_Bundle_Amount > 0 && award_request.EventAmount > 0 {
						flt_fractions := award_request.EventAmount / rules.MM_RC_Bundle_Amount
						flt_points := (flt_fractions * rules.MM_RC_Bundle_Points) + current_outstanding_points
						int_points := int(flt_points)
						outstanding_points = flt_points - float64(int_points)
						return float64(int_points), outstanding_points, rules.MM_RC_Bundle_Notification, rules.MM_RC_Bundle_Notification_Sender, rules.MM_RC_Bundle_Notification_Text
					}
				}
				return 0, current_outstanding_points, false, "", ""
			} else {
				if rules.MM_RC_Airtime_Award_Type == "Transaction" {
					if rules.MM_RC_Airtime_Points > 0 {
						// return rules.MM_RC_Airtime_Points, current_outstanding_points
						flt_points := rules.MM_RC_Airtime_Points + current_outstanding_points
						int_points := int(flt_points)
						outstanding_points = flt_points - float64(int_points)
						return float64(int_points), outstanding_points, rules.MM_RC_Airtime_Notification, rules.MM_RC_Airtime_Notification_Sender, rules.MM_RC_Airtime_Notification_Text
					}
				} else if rules.MM_RC_Airtime_Award_Type == "Amount" {
					if rules.MM_RC_Airtime_Amount > 0 && award_request.EventAmount > 0 {
						flt_fractions := award_request.EventAmount / rules.MM_RC_Airtime_Amount
						flt_points := (flt_fractions * rules.MM_RC_Airtime_Points) + current_outstanding_points
						int_points := int(flt_points)
						outstanding_points = flt_points - float64(int_points)
						return float64(int_points), outstanding_points, rules.MM_RC_Airtime_Notification, rules.MM_RC_Airtime_Notification_Sender, rules.MM_RC_Airtime_Notification_Text
					}
				}
				return 0, current_outstanding_points, false, "", ""
			}
		case "CTMMOREQ": //recharge for others
			if CheckMMEventDetailType(award_request.EventDetail) == Configuration.LoyaltyMMBundleCode {
				if rules.MM_CTMMOREQ_Bundle_Award_Type == "Transaction" {
					if rules.MM_CTMMOREQ_Bundle_Points > 0 {
						// return rules.MM_CTMMOREQ_Bundle_Points, current_outstanding_points
						flt_points := rules.MM_CTMMOREQ_Bundle_Points + current_outstanding_points
						int_points := int(flt_points)
						outstanding_points = flt_points - float64(int_points)
						return float64(int_points), outstanding_points, rules.MM_CTMMOREQ_Bundle_Notification, rules.MM_CTMMOREQ_Bundle_Notification_Sender, rules.MM_CTMMOREQ_Bundle_Notification_Text
					}
				} else if rules.MM_CTMMOREQ_Bundle_Award_Type == "Amount" {
					if rules.MM_CTMMOREQ_Bundle_Amount > 0 && award_request.EventAmount > 0 {
						flt_fractions := award_request.EventAmount / rules.MM_CTMMOREQ_Bundle_Amount
						flt_points := (flt_fractions * rules.MM_CTMMOREQ_Bundle_Points) + current_outstanding_points
						int_points := int(flt_points)
						outstanding_points = flt_points - float64(int_points)
						return float64(int_points), outstanding_points, rules.MM_CTMMOREQ_Bundle_Notification, rules.MM_CTMMOREQ_Bundle_Notification_Sender, rules.MM_CTMMOREQ_Bundle_Notification_Text
					}
				}
				return 0, current_outstanding_points, false, "", ""
			} else {
				if rules.MM_CTMMOREQ_Airtime_Award_Type == "Transaction" {
					if rules.MM_CTMMOREQ_Airtime_Points > 0 {
						flt_points := rules.MM_CTMMOREQ_Airtime_Points + current_outstanding_points
						int_points := int(flt_points)
						outstanding_points = flt_points - float64(int_points)
						return float64(int_points), outstanding_points, rules.MM_CTMMOREQ_Airtime_Notification, rules.MM_CTMMOREQ_Airtime_Notification_Sender, rules.MM_CTMMOREQ_Airtime_Notification_Text
						// return rules.MM_CTMMOREQ_Airtime_Points, current_outstanding_points
					}
				} else if rules.MM_CTMMOREQ_Airtime_Award_Type == "Amount" {
					if rules.MM_CTMMOREQ_Airtime_Amount > 0 && award_request.EventAmount > 0 {
						flt_fractions := award_request.EventAmount / rules.MM_CTMMOREQ_Airtime_Amount
						flt_points := (flt_fractions * rules.MM_CTMMOREQ_Airtime_Points) + current_outstanding_points
						int_points := int(flt_points)
						outstanding_points = flt_points - float64(int_points)
						return float64(int_points), outstanding_points, rules.MM_CTMMOREQ_Airtime_Notification, rules.MM_CTMMOREQ_Airtime_Notification_Sender, rules.MM_CTMMOREQ_Airtime_Notification_Text
					}
				}
				return 0, current_outstanding_points, false, "", ""
			}
		case "CBWREQ": //recharge for others
			if rules.MM_CBWREQ_Award_Type == "Transaction" {
				if rules.MM_CBWREQ_Points > 0 {
					// return rules.MM_CBWREQ_Points, current_outstanding_points
					flt_points := rules.MM_CBWREQ_Points + current_outstanding_points
					int_points := int(flt_points)
					outstanding_points = flt_points - float64(int_points)
					return float64(int_points), outstanding_points, rules.MM_CBWREQ_Notification, rules.MM_CBWREQ_Notification_Sender, rules.MM_CBWREQ_Notification_Text
				}
			} else if rules.MM_CBWREQ_Award_Type == "Amount" {
				if rules.MM_CBWREQ_Amount > 0 && award_request.EventAmount > 0 {
					flt_fractions := award_request.EventAmount / rules.MM_CBWREQ_Amount
					flt_points := (flt_fractions * rules.MM_CBWREQ_Points) + current_outstanding_points
					int_points := int(flt_points)
					outstanding_points = flt_points - float64(int_points)
					return float64(int_points), outstanding_points, rules.MM_CBWREQ_Notification, rules.MM_CBWREQ_Notification_Sender, rules.MM_CBWREQ_Notification_Text
				}
			}
			return 0, current_outstanding_points, false, "", ""
		default:
			return 0, current_outstanding_points, false, "", ""
		}
	case "MyAfricellApp":
		switch award_request.EventType {
		case "MobileAppDaily_Login":
			return rules.MobileAppDaily_Login, current_outstanding_points, rules.MobileAppDaily_Notification, rules.MobileAppDaily_Notification_Sender, rules.MobileAppDaily_Notification_Text
		default:
			return 0, current_outstanding_points, false, "", ""
		}
	}
	return
}

func CheckMMEventDetailType(EventDetail string) (response string) {
	indexOf := strings.Index(EventDetail, " - ")
	if indexOf > 0 {
		return EventDetail[0:indexOf]
	} else {
		return ""
	}
}

func (Uc *UserControl) Customer_Loyalty_Account_GetDebitPoints_log(startDate, endDate time.Time, MSISDN string, Filter string) (response []Loyalty_AccountDebitPoints_log, err error) {

	findResult, err := Uc.ReadAccountDebitPointsDetailsFromMongoDB(startDate, endDate, MSISDN, 1, 0)
	if err != nil {
		return response, err
	}

	response = append(response, findResult...)

	return response, nil
}

func (Uc *UserControl) ReadAccountDebitPointsDetailsFromMongoDB(startDate, endDate time.Time, MSISDN string, page, limit int64) (response []Loyalty_AccountDebitPoints_log, err error) {
	// Ensure startDate is before or equal to endDate

	if startDate.After(endDate) {
		return response, fmt.Errorf("start date cannot be after end date")
	}

	if page < 1 {
		page = 1 // Default to first page if invalid page number
	}

	// Iterate over the range of dates
	for currentDate := startDate; !currentDate.After(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		monthStr := strconv.Itoa(int(currentDate.Month()))
		if len(monthStr) < 2 {
			monthStr = "0" + monthStr
		}
		MongoDB_DB_Name := "Loyalty_DB_" + strconv.Itoa(currentDate.Year()) + monthStr

		var dayStr = strconv.Itoa(currentDate.Day())
		if len(dayStr) < 2 {
			dayStr = "0" + dayStr
		}

		collName := "Col_Loyalty_AccountDebitPoints_log_" + dayStr

		// Fetch the collection
		collection := Uc.LoyaltyMongoClient.Mongo.Database(MongoDB_DB_Name).Collection(collName)
		// Build the filter for the date range
		filter := bson.D{
			{Key: "MSISDN", Value: MSISDN},
		}
		skip := (page - 1) * limit

		// Fetch a paginated list of documents using limit and skip
		cursor, err := collection.Find(
			context.Background(),
			filter,
			options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)),
		)
		if err != nil {
			return response, err
		}
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var result Loyalty_AccountDebitPoints_log

			if err := cursor.Decode(&result); err != nil {
				return nil, fmt.Errorf("failed to decode result: %w", err)
			}
			response = append(response, result)
		}

		if err := cursor.Err(); err != nil {
			return nil, fmt.Errorf("cursor error: %w", err)
		}
	}

	return response, err
}

func (Uc *UserControl) Customer_Loyalty_Account_GetCreditPoints_log(startDate, endDate time.Time, MSISDN string, Filter string) (response []Loyalty_AccountCreditPoints_log, err error) {

	findResult, err := Uc.ReadAccountCreditPointsDetailsFromMongoDB(startDate, endDate, MSISDN, 1, 0)
	if err != nil {
		return response, err
	}
	response = append(response, findResult...)

	return response, nil
}
func (Uc *UserControl) ReadAccountCreditPointsDetailsFromMongoDB(startDate, endDate time.Time, MSISDN string, page, limit int64) (response []Loyalty_AccountCreditPoints_log, err error) {
	// Ensure startDate is before or equal to endDate

	if startDate.After(endDate) {
		return response, fmt.Errorf("start date cannot be after end date")
	}

	if page < 1 {
		page = 1 // Default to first page if invalid page number
	}

	// Iterate over the range of dates
	for currentDate := startDate; !currentDate.After(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		monthStr := strconv.Itoa(int(currentDate.Month()))
		if len(monthStr) < 2 {
			monthStr = "0" + monthStr
		}
		MongoDB_DB_Name := "Loyalty_DB_" + strconv.Itoa(currentDate.Year()) + monthStr

		var dayStr = strconv.Itoa(currentDate.Day())
		if len(dayStr) < 2 {
			dayStr = "0" + dayStr
		}

		collName := "Col_Loyalty_AccountCreditPoints_log_" + dayStr

		// Fetch the collection
		collection := Uc.LoyaltyMongoClient.Mongo.Database(MongoDB_DB_Name).Collection(collName)
		// Build the filter for the date range
		filter := bson.D{
			{Key: "MSISDN", Value: MSISDN},
		}
		skip := (page - 1) * limit

		// Fetch a paginated list of documents using limit and skip
		cursor, err := collection.Find(
			context.Background(),
			filter,
			options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)),
		)
		if err != nil {
			return response, err
		}
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var result Loyalty_AccountCreditPoints_log

			if err := cursor.Decode(&result); err != nil {
				return nil, fmt.Errorf("failed to decode result: %w", err)
			}
			response = append(response, result)
		}

		if err := cursor.Err(); err != nil {
			return nil, fmt.Errorf("cursor error: %w", err)
		}
	}

	return response, err
}

func (Uc *UserControl) Customer_Loyalty_Account_GetRedemptionPoints_log(startDate, endDate time.Time, MSISDN string, Filter string) (response []Loyalty_Redemption_log, err error) {

	findResult, err := Uc.ReadAccountRedemptionPointsDetailsFromMongoDB(startDate, endDate, MSISDN, 1, 0)
	if err != nil {
		return response, err
	}
	response = append(response, findResult...)

	return response, nil
}

func (Uc *UserControl) ReadAccountRedemptionPointsDetailsFromMongoDB(startDate, endDate time.Time, MSISDN string, page, limit int64) (response []Loyalty_Redemption_log, err error) {
	// Ensure startDate is before or equal to endDate

	if startDate.After(endDate) {
		return response, fmt.Errorf("start date cannot be after end date")
	}

	if page < 1 {
		page = 1 // Default to first page if invalid page number
	}

	// Iterate over the range of dates
	for currentDate := startDate; !currentDate.After(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		monthStr := strconv.Itoa(int(currentDate.Month()))
		if len(monthStr) < 2 {
			monthStr = "0" + monthStr
		}
		MongoDB_DB_Name := "Loyalty_DB_" + strconv.Itoa(currentDate.Year()) + monthStr

		var dayStr = strconv.Itoa(currentDate.Day())
		if len(dayStr) < 2 {
			dayStr = "0" + dayStr
		}

		collName := "Col_Loyalty_Redemption_log_" + dayStr

		// Fetch the collection
		collection := Uc.LoyaltyMongoClient.Mongo.Database(MongoDB_DB_Name).Collection(collName)
		// Build the filter for the date range
		filter := bson.D{
			{Key: "MSISDN", Value: MSISDN},
		}
		skip := (page - 1) * limit

		// Fetch a paginated list of documents using limit and skip
		cursor, err := collection.Find(
			context.Background(),
			filter,
			options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)),
		)
		if err != nil {
			return response, err
		}
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var result Loyalty_Redemption_log

			if err := cursor.Decode(&result); err != nil {
				return nil, fmt.Errorf("failed to decode result: %w", err)
			}
			response = append(response, result)
		}

		if err := cursor.Err(); err != nil {
			return nil, fmt.Errorf("cursor error: %w", err)
		}
	}

	return response, err
}

func (Uc *UserControl) Customer_Loyalty_Account_GetExpiryPoints_log(startDate, endDate time.Time, MSISDN string, Filter string) (response []Loyalty_Monthly_Expiry_log, err error) {

	findResult, err := Uc.ReadAccountExpiryPointsDetailsFromMongoDB(startDate, endDate, MSISDN, 1, 0)
	if err != nil {
		return response, err
	}
	response = append(response, findResult...)

	return response, nil
}

func (Uc *UserControl) ReadAccountExpiryPointsDetailsFromMongoDB(startDate, endDate time.Time, MSISDN string, page, limit int64) (response []Loyalty_Monthly_Expiry_log, err error) {
	// Ensure startDate is before or equal to endDate

	if startDate.After(endDate) {
		return response, fmt.Errorf("start date cannot be after end date")
	}

	if page < 1 {
		page = 1 // Default to first page if invalid page number
	}

	// Iterate over the range of dates
	for currentDate := startDate; !currentDate.After(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		monthStr := strconv.Itoa(int(currentDate.Month()))
		if len(monthStr) < 2 {
			monthStr = "0" + monthStr
		}
		MongoDB_DB_Name := "Loyalty_DB_" + strconv.Itoa(currentDate.Year()) + monthStr

		var dayStr = strconv.Itoa(currentDate.Day())
		if len(dayStr) < 2 {
			dayStr = "0" + dayStr
		}

		collName := "Col_Loyalty_Expiry_log_" + dayStr

		// Fetch the collection
		collection := Uc.LoyaltyMongoClient.Mongo.Database(MongoDB_DB_Name).Collection(collName)
		// Build the filter for the date range
		filter := bson.D{
			{Key: "MSISDN", Value: MSISDN},
		}
		skip := (page - 1) * limit

		// Fetch a paginated list of documents using limit and skip
		cursor, err := collection.Find(
			context.Background(),
			filter,
			options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)),
		)
		if err != nil {
			return response, err
		}
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var result Loyalty_Monthly_Expiry_log

			if err := cursor.Decode(&result); err != nil {
				return nil, fmt.Errorf("failed to decode result: %w", err)
			}
			response = append(response, result)
		}

		if err := cursor.Err(); err != nil {
			return nil, fmt.Errorf("cursor error: %w", err)
		}
	}

	return response, err
}

func (Uc *UserControl) Customer_Loyalty_Account_GetLevelChange_log(startDate, endDate time.Time, MSISDN string, Filter string) (response []Loyalty_Level_Change_log, err error) {

	findResult, err := Uc.ReadAccountLevelChangeDetailsFromMongoDB(startDate, endDate, MSISDN, 1, 0)
	if err != nil {
		return response, err
	}
	response = append(response, findResult...)

	return response, nil
}

func (Uc *UserControl) ReadAccountLevelChangeDetailsFromMongoDB(startDate, endDate time.Time, MSISDN string, page, limit int64) (response []Loyalty_Level_Change_log, err error) {
	// Ensure startDate is before or equal to endDate

	if startDate.After(endDate) {
		return response, fmt.Errorf("start date cannot be after end date")
	}

	if page < 1 {
		page = 1 // Default to first page if invalid page number
	}

	// Iterate over the range of dates
	for currentDate := startDate; !currentDate.After(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		monthStr := strconv.Itoa(int(currentDate.Month()))
		if len(monthStr) < 2 {
			monthStr = "0" + monthStr
		}
		MongoDB_DB_Name := "Loyalty_DB_" + strconv.Itoa(currentDate.Year()) + monthStr

		collName := "Col_Loyalty_Level_Change_log"

		// Fetch the collection
		collection := Uc.LoyaltyMongoClient.Mongo.Database(MongoDB_DB_Name).Collection(collName)
		// Build the filter for the date range
		filter := bson.D{
			{Key: "MSISDN", Value: MSISDN},
			{Key: "Level_Change_Date", Value: bson.D{
				{Key: "$gte", Value: startDate},
				{Key: "$lte", Value: endDate},
			}},
		}
		skip := (page - 1) * limit

		// Fetch a paginated list of documents using limit and skip
		cursor, err := collection.Find(
			context.Background(),
			filter,
			options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)),
		)
		if err != nil {
			return response, err
		}
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var result Loyalty_Level_Change_log

			if err := cursor.Decode(&result); err != nil {
				return nil, fmt.Errorf("failed to decode result: %w", err)
			}
			response = append(response, result)
		}

		if err := cursor.Err(); err != nil {
			return nil, fmt.Errorf("cursor error: %w", err)
		}
	}

	return response, err
}

func (Uc *UserControl) Customer_Loyalty_Account_GetEvents_log(startDate, endDate time.Time, MSISDN string, Filter string) (response []Loyalty_Event_Log, err error) {

	findResult, err := Uc.ReadAccountEventsDetailsFromMongoDB(startDate, endDate, MSISDN, 1, 0)
	if err != nil {
		return response, err
	}
	response = append(response, findResult...)

	return response, nil
}

func (Uc *UserControl) ReadAccountEventsDetailsFromMongoDB(startDate, endDate time.Time, MSISDN string, page, limit int64) (response []Loyalty_Event_Log, err error) {
	// Ensure startDate is before or equal to endDate

	if startDate.After(endDate) {
		return response, fmt.Errorf("start date cannot be after end date")
	}

	if page < 1 {
		page = 1 // Default to first page if invalid page number
	}

	// Iterate over the range of dates
	for currentDate := startDate; !currentDate.After(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		monthStr := strconv.Itoa(int(currentDate.Month()))
		if len(monthStr) < 2 {
			monthStr = "0" + monthStr
		}
		MongoDB_DB_Name := "Loyalty_DB_" + strconv.Itoa(currentDate.Year()) + monthStr

		var dayStr = strconv.Itoa(currentDate.Day())
		if len(dayStr) < 2 {
			dayStr = "0" + dayStr
		}

		collName := "Col_Loyalty_Event_Log_" + dayStr

		// Fetch the collection
		collection := Uc.LoyaltyMongoClient.Mongo.Database(MongoDB_DB_Name).Collection(collName)
		// Build the filter for the date range
		filter := bson.D{
			{Key: "MSISDN", Value: MSISDN},
		}
		skip := (page - 1) * limit

		// Fetch a paginated list of documents using limit and skip
		cursor, err := collection.Find(
			context.Background(),
			filter,
			options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)),
		)
		if err != nil {
			return response, err
		}
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var result Loyalty_Event_Log

			if err := cursor.Decode(&result); err != nil {
				return nil, fmt.Errorf("failed to decode result: %w", err)
			}
			response = append(response, result)
		}

		if err := cursor.Err(); err != nil {
			return nil, fmt.Errorf("cursor error: %w", err)
		}
	}

	return response, err
}

func (Uc *UserControl) Customer_Loyalty_Account_GetStatus_log(startDate, endDate time.Time, MSISDN string, Filter string) (response []Loyalty_Status_log, err error) {

	findResult, err := Uc.ReadAccountStatusDetailsFromMongoDB(startDate, endDate, MSISDN, 1, 0)
	if err != nil {
		return response, err
	}

	response = append(response, findResult...)

	return response, nil
}

func (Uc *UserControl) ReadAccountStatusDetailsFromMongoDB(startDate, endDate time.Time, MSISDN string, page, limit int64) (response []Loyalty_Status_log, err error) {
	// Ensure startDate is before or equal to endDate

	if startDate.After(endDate) {
		return response, fmt.Errorf("start date cannot be after end date")
	}

	if page < 1 {
		page = 1 // Default to first page if invalid page number
	}

	// Iterate over the range of dates
	for currentDate := startDate; !currentDate.After(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		monthStr := strconv.Itoa(int(currentDate.Month()))
		if len(monthStr) < 2 {
			monthStr = "0" + monthStr
		}
		MongoDB_DB_Name := "Loyalty_DB_" + strconv.Itoa(currentDate.Year()) + monthStr

		var dayStr = strconv.Itoa(currentDate.Day())
		if len(dayStr) < 2 {
			dayStr = "0" + dayStr
		}

		collName := "Col_Loyalty_Status_log_" + dayStr

		// Fetch the collection
		collection := Uc.LoyaltyMongoClient.Mongo.Database(MongoDB_DB_Name).Collection(collName)
		// Build the filter for the date range
		filter := bson.D{
			{Key: "MSISDN", Value: MSISDN},
		}
		skip := (page - 1) * limit

		// Fetch a paginated list of documents using limit and skip
		cursor, err := collection.Find(
			context.Background(),
			filter,
			options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)),
		)
		if err != nil {
			return response, err
		}
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var result Loyalty_Status_log

			if err := cursor.Decode(&result); err != nil {
				return nil, fmt.Errorf("failed to decode result: %w", err)
			}
			response = append(response, result)
		}

		if err := cursor.Err(); err != nil {
			return nil, fmt.Errorf("cursor error: %w", err)
		}
	}

	return response, err
}

func (Uc *UserControl) Customer_Loyalty_Account_GetStatusExpiry_log(startDate, endDate time.Time, MSISDN string, Filter string) (response []Loyalty_Full_Expiry_Log, err error) {

	findResult, err := Uc.ReadAccountStatusExpiryDetailsFromMongoDB(startDate, endDate, MSISDN, 1, 0)
	if err != nil {
		return response, err
	}

	response = append(response, findResult...)

	return response, nil
}

func (Uc *UserControl) ReadAccountStatusExpiryDetailsFromMongoDB(startDate, endDate time.Time, MSISDN string, page, limit int64) (response []Loyalty_Full_Expiry_Log, err error) {
	// Ensure startDate is before or equal to endDate

	if startDate.After(endDate) {
		return response, fmt.Errorf("start date cannot be after end date")
	}

	if page < 1 {
		page = 1 // Default to first page if invalid page number
	}

	// Iterate over the range of dates
	for currentDate := startDate; !currentDate.After(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		monthStr := strconv.Itoa(int(currentDate.Month()))
		if len(monthStr) < 2 {
			monthStr = "0" + monthStr
		}
		MongoDB_DB_Name := "Loyalty_DB_" + strconv.Itoa(currentDate.Year()) + monthStr

		var dayStr = strconv.Itoa(currentDate.Day())
		if len(dayStr) < 2 {
			dayStr = "0" + dayStr
		}

		collName := "Col_Loyalty_Full_Expiry_log_" + dayStr

		// Fetch the collection
		collection := Uc.LoyaltyMongoClient.Mongo.Database(MongoDB_DB_Name).Collection(collName)
		// Build the filter for the date range
		filter := bson.D{
			{Key: "MSISDN", Value: MSISDN},
		}
		skip := (page - 1) * limit

		// Fetch a paginated list of documents using limit and skip
		cursor, err := collection.Find(
			context.Background(),
			filter,
			options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)),
		)
		if err != nil {
			return response, err
		}
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var result Loyalty_Full_Expiry_Log

			if err := cursor.Decode(&result); err != nil {
				return nil, fmt.Errorf("failed to decode result: %w", err)
			}
			response = append(response, result)
		}

		if err := cursor.Err(); err != nil {
			return nil, fmt.Errorf("cursor error: %w", err)
		}
	}

	return response, err
}

type LoyaltySummary struct {
	MSISDN             string  `bson:"_id"`
	TotalPoints        float64 `bson:"totalPoints"`
	LoyaltyLevel       string  `bson:"loyaltyLevel"`
	AfrimoneyPoints    float64 `bson:"afrimoneyPoints"`
	Accumulated_Points float64 `bson:"accumulated_Points"`
	Current_Points     float64 `bson:"current_Points"`
}

func (Uc *UserControl) Customer_Loyalty_Account_GetAwardedPoints(startDate, endDate time.Time) (response map[string]LoyaltySummary, err error) {
	results := make(map[string]LoyaltySummary)

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		monthStr := strconv.Itoa(int(d.Month()))
		if len(monthStr) < 2 {
			monthStr = "0" + monthStr
		}
		MongoDB_DB_Name := "Loyalty_DB_" + strconv.Itoa(d.Year()) + monthStr

		var dayStr = strconv.Itoa(d.Day())
		if len(dayStr) < 2 {
			dayStr = "0" + dayStr
		}

		collName := "Col_Loyalty_AccountCreditPoints_log_" + dayStr
		// Fetch the collection
		collection := Uc.LoyaltyMongoClient.Mongo.Database(MongoDB_DB_Name).Collection(collName)
		pipeline := mongo.Pipeline{
			{{Key: "$match", Value: bson.D{
				{Key: "AwardedPoints", Value: bson.D{{Key: "$gt", Value: 0}}},
			}}},
			{{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$MSISDN"},
				{Key: "totalPoints", Value: bson.D{{Key: "$sum", Value: "$AwardedPoints"}}},
				{Key: "afrimoneyPoints", Value: bson.D{
					{Key: "$sum", Value: bson.D{
						{Key: "$cond", Value: bson.A{
							bson.D{{Key: "$eq", Value: bson.A{"$EventSource", "MobileMoney_feed"}}},
							"$AwardedPoints", // if condition is true
							0,                // if false
						}},
					}},
				}},
			}}},
		}

		cursor, err := collection.Aggregate(context.Background(), pipeline)
		if err != nil {
			log.Printf("Aggregation failed for %s: %v", collName, err)
			continue
		}

		for cursor.Next(context.Background()) {
			var doc LoyaltySummary
			if err := cursor.Decode(&doc); err != nil {
				log.Printf("Failed decoding result for %s: %v", collName, err)
				continue
			}
			existing := results[doc.MSISDN]
			cusAccount, caErr := redisx.GetJSON[Customer_Loyalty_Account](context.Background(), RedisClient, Customer_Loyalty_Account{Key: doc.MSISDN}.RedisKey())
			if redisx.IsNil(caErr) {
				fmt.Println("key does not exist")
			} else if caErr != nil {
				fmt.Println("error in type assertion")
			}
			existing.LoyaltyLevel = cusAccount.Loyalty_Level_Key
			existing.Current_Points = cusAccount.Available_Points
			existing.Accumulated_Points = cusAccount.Awarded_Points
			existing.TotalPoints += doc.TotalPoints
			existing.AfrimoneyPoints += doc.AfrimoneyPoints
			results[doc.MSISDN] = existing
		}
	}

	return results, nil
}

func (Uc *UserControl) Customer_Loyalty_Account_Getlogs(Type string, startDate, endDate time.Time, MSISDN string, Filter string) (response []Loyalty_Logs, err error) {
	var allLogs []Loyalty_Logs
	if Type == "All Logs" || Type == "Debit Logs" {
		findResult, err := Uc.Customer_Loyalty_Account_GetDebitPoints_log(startDate, endDate, MSISDN, Filter)
		if err != nil {
			return response, err
		}
		for _, log := range findResult {
			var newLog = Loyalty_Logs{}
			newLog.Logs_Type = "Debit"
			newLog.Debit_Logs = log
			newLog.Status = log.Status
			newLog.Date = log.ReceiveDate
			newLog.PointsToCredit = log.Debit_Amount
			newLog.Opening_Available_Points = log.Opening_Available_Points
			newLog.Closure_Available_Points = log.Closure_Available_Points
			allLogs = append(allLogs, newLog)
		}
	}
	if Type == "All Logs" || Type == "Credit Logs" {
		findResult, err := Uc.Customer_Loyalty_Account_GetCreditPoints_log(startDate, endDate, MSISDN, Filter)
		if err != nil {
			return response, err
		}
		for _, log := range findResult {
			var newLog = Loyalty_Logs{}
			newLog.Logs_Type = "Credit"
			newLog.Credit_Logs = log
			newLog.Date = log.ReceiveDate
			newLog.Status = log.Status
			newLog.PointsToCredit = log.AwardedPoints
			newLog.Opening_Available_Points = log.Opening_Available_Points
			newLog.Closure_Available_Points = log.Closure_Available_Points
			allLogs = append(allLogs, newLog)
		}
	}
	if Type == "All Logs" || Type == "Redemption Logs" {
		findResult, err := Uc.Customer_Loyalty_Account_GetRedemptionPoints_log(startDate, endDate, MSISDN, Filter)
		if err != nil {
			return response, err
		}
		for _, log := range findResult {
			var newLog = Loyalty_Logs{}
			newLog.Logs_Type = "Redemption"
			newLog.Redemption_Logs = log
			newLog.Status = log.Status
			newLog.Date = log.ReceiveDate
			newLog.PointsToCredit = log.Points_To_Redeem
			newLog.Opening_Available_Points = log.Opening_Available_Points
			newLog.Closure_Available_Points = log.Closure_Available_Points
			allLogs = append(allLogs, newLog)
		}
	}
	if Type == "All Logs" || Type == "Expiry Logs" {
		findResult, err := Uc.Customer_Loyalty_Account_GetExpiryPoints_log(startDate, endDate, MSISDN, Filter)
		if err != nil {
			return response, err
		}
		for _, log := range findResult {
			var newLog = Loyalty_Logs{}
			newLog.Logs_Type = "Expiry " + log.Year_Month
			newLog.Expiry_Logs = log
			newLog.Status = log.ExpiryStatus
			newLog.Date = log.ExpiryTime
			newLog.Opening_Available_Points = log.Opening_Available_Points
			newLog.Closure_Available_Points = log.End_Available_Points
			allLogs = append(allLogs, newLog)
		}
	}
	if Type == "All Logs" || Type == "Level Change Logs" {
		findResult, err := Uc.Customer_Loyalty_Account_GetLevelChange_log(startDate, endDate, MSISDN, Filter)
		if err != nil {
			return response, err
		}
		for _, log := range findResult {
			var newLog = Loyalty_Logs{}
			newLog.Logs_Type = "Level Change"
			newLog.Status = "successful"
			newLog.Opening_Available_Points = log.Available_Points
			newLog.Closure_Available_Points = log.Available_Points
			newLog.Level_Change_Logs = log
			newLog.Date = log.Level_Change_Date
			allLogs = append(allLogs, newLog)
		}
	}
	if Type == "All Logs" || Type == "Status Logs" {
		findResult, err := Uc.Customer_Loyalty_Account_GetStatus_log(startDate, endDate, MSISDN, Filter)
		if err != nil {
			return response, err
		}
		for _, log := range findResult {
			var newLog = Loyalty_Logs{}
			newLog.Logs_Type = "Status"
			newLog.Status_Logs = log
			newLog.Status = log.Request_Status
			newLog.Date = log.StatusDate
			newLog.PointsToCredit = 0
			newLog.Opening_Available_Points = 0
			newLog.Closure_Available_Points = 0
			allLogs = append(allLogs, newLog)
		}
	}
	if Type == "All Logs" || Type == "Full Expiry Logs" {
		findResult, err := Uc.Customer_Loyalty_Account_GetStatusExpiry_log(startDate, endDate, MSISDN, Filter)
		if err != nil {
			return response, err
		}
		for _, log := range findResult {
			var newLog = Loyalty_Logs{}
			newLog.Logs_Type = "Full Expiry"
			newLog.Full_Expiry_Logs = log
			newLog.Status = log.ExpiryStatus
			newLog.Date = log.ExpiryTime
			newLog.PointsToCredit = log.ExpiryAmount
			newLog.Opening_Available_Points = log.Opening_Available_Points
			newLog.Closure_Available_Points = log.End_Available_Points
			allLogs = append(allLogs, newLog)
		}
	}
	sort.Slice(allLogs, func(i, j int) bool {
		// Sort by date descending: newest first
		return allLogs[i].Date.After(allLogs[j].Date)
	})
	response = append(response, allLogs...)

	return response, nil
}

func (Uc *UserControl) ReadAccountLogsDetailsFromMongoDB(Type string, startDate, endDate time.Time, MSISDN string, page, limit int64) (response []Loyalty_Event_Log, err error) {
	// Ensure startDate is before or equal to endDate

	if startDate.After(endDate) {
		return response, fmt.Errorf("start date cannot be after end date")
	}

	if page < 1 {
		page = 1 // Default to first page if invalid page number
	}

	// Iterate over the range of dates
	for currentDate := startDate; !currentDate.After(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		monthStr := strconv.Itoa(int(currentDate.Month()))
		if len(monthStr) < 2 {
			monthStr = "0" + monthStr
		}
		MongoDB_DB_Name := "Loyalty_DB_" + strconv.Itoa(currentDate.Year()) + monthStr

		var dayStr = strconv.Itoa(currentDate.Day())
		if len(dayStr) < 2 {
			dayStr = "0" + dayStr
		}

		collName := "Col_Loyalty_Event_Log_" + dayStr

		// Fetch the collection
		collection := Uc.LoyaltyMongoClient.Mongo.Database(MongoDB_DB_Name).Collection(collName)
		// Build the filter for the date range
		filter := bson.D{
			{Key: "MSISDN", Value: MSISDN},
		}
		skip := (page - 1) * limit

		// Fetch a paginated list of documents using limit and skip
		cursor, err := collection.Find(
			context.Background(),
			filter,
			options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)),
		)
		if err != nil {
			return response, err
		}
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var result Loyalty_Event_Log

			if err := cursor.Decode(&result); err != nil {
				return nil, fmt.Errorf("failed to decode result: %w", err)
			}
			response = append(response, result)
		}

		if err := cursor.Err(); err != nil {
			return nil, fmt.Errorf("cursor error: %w", err)
		}
	}

	return response, err
}

func (Uc *UserControl) EvaluateAndUpdate_CustomerLoyaltyLevel(Login string, Account_Key string) (New_Loyalty_Level_Key string, err error) {
	loyalty_account, laErr := redisx.GetJSON[Customer_Loyalty_Account](context.Background(), RedisClient, Customer_Loyalty_Account{Key: Account_Key}.RedisKey())
	if redisx.IsNil(laErr) {
		return New_Loyalty_Level_Key, errors.New("loyalty account does not exist")
	}
	if laErr != nil {
		return New_Loyalty_Level_Key, errors.New("type assertion issue with Customer_Loyalty_Account")
	}
	//evaluate loyalty level
	var New_Loyalty_Level Loyalty_Level
	llScanCtx, llScanCancel := context.WithTimeout(context.Background(), 10*time.Second)
	loyalty_Level_slice, llScanErr := redisx.GetAllJSONByPattern[Loyalty_Level](llScanCtx, RedisClient, redisx.ScanJSONOptions{Pattern: "Loyalty_Level:*"})
	llScanCancel()
	if llScanErr != nil {
		return New_Loyalty_Level_Key, fmt.Errorf("error scanning Loyalty_Level: %w", llScanErr)
	}
	if len(loyalty_Level_slice) > 0 {
		for _, loyalty_Level := range loyalty_Level_slice {
			//evaluate
			if loyalty_account.Awarded_Points >= loyalty_Level.Min_Accumulated_Points && loyalty_account.Awarded_Points < loyalty_Level.Max_Accumulated_Points {
				New_Loyalty_Level = loyalty_Level
				New_Loyalty_Level_Key = New_Loyalty_Level.Key
				if New_Loyalty_Level.Key == loyalty_account.Loyalty_Level_Key {
					//===>> no level change
					return New_Loyalty_Level_Key, nil
				} else {
					break
				}
			}
		}
		//update customer level
		current_level, clErr := redisx.GetJSON[Loyalty_Level](context.Background(), RedisClient, Loyalty_Level{Key: loyalty_account.Loyalty_Level_Key}.RedisKey())
		if redisx.IsNil(clErr) {
			return New_Loyalty_Level_Key, errors.New("current level is invalid")
		}
		if clErr != nil {
			return New_Loyalty_Level_Key, errors.New("error in type assertion")
		}
		if New_Loyalty_Level.Min_Accumulated_Points > current_level.Min_Accumulated_Points &&
			New_Loyalty_Level.Max_Accumulated_Points > current_level.Max_Accumulated_Points {
			//Upgrade level
			loyalty_account.Previous_Loyalty_Level_Key = loyalty_account.Loyalty_Level_Key
			loyalty_account.Previous_Loyalty_Level_Date = loyalty_account.Loyalty_Level_Date
			loyalty_account.Loyalty_Level_Key = New_Loyalty_Level.Key
			loyalty_account.Loyalty_Level_Date = time.Now()
			loyalty_account.Loyalty_Level_Direction = "Upgrade"
			loyalty_account.Loyalty_Level_SetBy = Login
			{
				putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
				if _, putErr := Mdb_Customer_Loyalty_Account.Coll.UpdateOne(putCtx, bson.M{"Key": loyalty_account.Key}, bson.M{"$set": loyalty_account}, options.UpdateOne().SetUpsert(true)); putErr != nil {
					log.Println("Mdb_Customer_Loyalty_Account upsert error:", putErr)
				}
				if putSetErr := redisx.SetJSON(putCtx, RedisClient, loyalty_account.RedisKey(), loyalty_account); putSetErr != nil {
					log.Println("redisx.SetJSON Customer_Loyalty_Account error:", putSetErr)
				}
				putCancel()
			}

		} else {
			//downgrade ==> Sof landing (one downgrade during anniversary year)
			if Login != "Points_Expiry" {
				return loyalty_account.Loyalty_Level_Key, errors.New("downgrade is allowed only on expiry")
			}
			//check if the last action was downgrade and allow 12 months for the next dowgrade
			// if loyalty_account.Loyalty_Level_Direction == "Downgrade" {
			// 	if loyalty_account.Loyalty_Level_Date.AddDate(0, 12, 0).After(time.Now()) {
			// 		return loyalty_account.Loyalty_Level_Key, errors.New("only one downgrade is allowed during one anniversary year")
			// 	}
			// }
			if current_level.DowngradeToLevel_Key != "" {
				New_Loyalty_Level_Key = current_level.DowngradeToLevel_Key
				loyalty_account.Previous_Loyalty_Level_Key = loyalty_account.Loyalty_Level_Key
				loyalty_account.Previous_Loyalty_Level_Date = loyalty_account.Loyalty_Level_Date
				loyalty_account.Loyalty_Level_Key = New_Loyalty_Level_Key
				loyalty_account.Loyalty_Level_Date = time.Now()
				loyalty_account.Loyalty_Level_Direction = "Downgrade"
				loyalty_account.Loyalty_Level_SetBy = Login
				{
					putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
					if _, putErr := Mdb_Customer_Loyalty_Account.Coll.UpdateOne(putCtx, bson.M{"Key": loyalty_account.Key}, bson.M{"$set": loyalty_account}, options.UpdateOne().SetUpsert(true)); putErr != nil {
						log.Println("Mdb_Customer_Loyalty_Account upsert error:", putErr)
					}
					if putSetErr := redisx.SetJSON(putCtx, RedisClient, loyalty_account.RedisKey(), loyalty_account); putSetErr != nil {
						log.Println("redisx.SetJSON Customer_Loyalty_Account error:", putSetErr)
					}
					putCancel()
				}
			} else {
				return loyalty_account.Loyalty_Level_Key, errors.New("downgrade level is not defined")
			}
		}
		//write change log
		Uc.Write_Loyalty_Level_Change_log(Loyalty_Level_Change_log{
			Level_Change_Date:                 time.Now(),
			MSISDN:                            loyalty_account.Key,
			COS:                               loyalty_account.COS,
			Joining_Date:                      loyalty_account.Joining_Date,
			ARPU:                              loyalty_account.ARPU,
			Customer_Id:                       loyalty_account.Customer_Id,
			Creation_date:                     loyalty_account.Creation_date,
			Previous_Loyalty_Level_Key:        loyalty_account.Previous_Loyalty_Level_Key,
			Previous_Loyalty_Level_Date:       loyalty_account.Previous_Loyalty_Level_Date,
			New_Loyalty_Level_Key:             loyalty_account.Loyalty_Level_Key,
			New_Loyalty_Level_Date:            loyalty_account.Loyalty_Level_Date,
			New_Loyalty_Level_Direction:       loyalty_account.Loyalty_Level_Direction,
			New_Loyalty_Level_SetBy:           loyalty_account.Loyalty_Level_SetBy,
			Loyalty_Account_Segment_Key:       loyalty_account.Loyalty_Account_Segment_Key,
			Loyalty_Account_Segment_Date:      loyalty_account.Loyalty_Account_Segment_Date,
			Loyalty_Account_Segment_Direction: loyalty_account.Loyalty_Account_Segment_Direction,
			Loyalty_Account_Segment_SetBy:     loyalty_account.Loyalty_Account_Segment_SetBy,
			Awarded_Points:                    loyalty_account.Awarded_Points,
			Redeemed_Points:                   loyalty_account.Redeemed_Points,
			Available_Points:                  loyalty_account.Available_Points,
			Last_Award_Date:                   loyalty_account.Last_Award_Date,
			Last_Redeem_Date:                  loyalty_account.Last_Redeem_Date,
		})
	}
	return New_Loyalty_Level_Key, nil
}

func (Uc *UserControl) PointsExpiry_Process() {
	exec := 0
	LOG_ID := "<<Points Expiry>>"
	for range time.Tick(time.Second * 1) {
		_CurrentDateTime := time.Now()
		_hr, _mi, _se := _CurrentDateTime.Clock()
		if _hr == 00 {
			if _mi == 00 {
				if _se < 60 {
					if exec == 0 {
						exec = 1
						log.Println(LOG_ID + " triggered")
						countCtx, countCancel := context.WithTimeout(context.Background(), 10*time.Second)
						count, err := Mdb_Customer_Loyalty_Account.Coll.CountDocuments(countCtx, bson.D{})
						countCancel()
						if err != nil {
							log.Println(LOG_ID + " count get error: " + err.Error())
						} else {
							if count > 0 {
								var QueryLimit int64 = 1000
								var QueryPage int64 = 1
								var endReached bool = false
								var QueryIdx int64 = 0
								for !endReached {
									loyalty_Accounts, err := Uc.Customer_Loyalty_Account_GetPaginated(QueryPage, QueryLimit)
									if err == nil {
										if QueryIdx < count {
											QueryPage = QueryPage + 1
											QueryIdx = QueryIdx + QueryLimit
										} else {
											endReached = true
										}
										// do the work here
										for _, loyalty_Account := range loyalty_Accounts {
											finalAccount, err := Uc.Customer_Loyalty_Account_Get(loyalty_Account.Key)
											if err != nil || len(finalAccount) == 0 {
												break
											}
											if finalAccount[0].Coming_Expiry_Date.Before(time.Now()) {
												go Uc.PointsExpiry_ProcessExec(finalAccount[0])
											}
										}
									}

								}

							}
						}
						log.Println(LOG_ID + " finished")
					}
				}
			} else {
				if exec == 1 {
					exec = 0
				}
			}

		}
	}
}

func (Uc *UserControl) PointsExpiry_ProcessExec(account Customer_Loyalty_Account) {
	chan_PointsExpiry_Controler <- 1
	var expiry_log Loyalty_Full_Expiry_Log
	expiry_log.ExpiryTime = time.Now()
	expiry_log.MSISDN = account.Key
	expiry_log.Opening_Awarded_Points = account.Awarded_Points
	expiry_log.Opening_Redeemed_Points = account.Redeemed_Points
	expiry_log.Opening_Available_Points = account.Available_Points
	expiry_log.Opening_OutStanding_Points = account.Outstanding_fraction_points
	expiry_log.Opening_Expired_Points = account.Expired_Points
	expiry_log.Opening_Redeemed_Expired_Points = account.Redeemed_Expired_Points
	expiry_log.OpeningLoyaltyLevel = account.Loyalty_Level_Key
	expiry_log.EndLoyaltyLevel = account.Loyalty_Level_Key
	expiry_log.ExpiryReason = "Cycle expiry reached"
	expiry_log.ExpiryAmount = account.Available_Points
	var monthly_expiry_log Loyalty_Monthly_Expiry_log
	monthly_expiry_log.ExpiryTime = time.Now()
	monthly_expiry_log.MSISDN = account.Key
	monthly_expiry_log.Opening_Awarded_Points = account.Awarded_Points
	monthly_expiry_log.Opening_Redeemed_Points = account.Redeemed_Points
	monthly_expiry_log.Opening_Available_Points = account.Available_Points
	monthly_expiry_log.Opening_Expired_Points = account.Expired_Points
	monthly_expiry_log.OpeningLoyaltyLevel = account.Loyalty_Level_Key
	monthly_expiry_log.EndLoyaltyLevel = account.Loyalty_Level_Key
	//validate the loyalty plan
	plan, planErr := redisx.GetJSON[Loyalty_Plan](context.Background(), RedisClient, Loyalty_Plan{Key: account.Loyalty_Level_Key + "|" + account.Loyalty_Account_Segment_Key}.RedisKey())
	if redisx.IsNil(planErr) {
		monthly_expiry_log.ExpiryStatus = "failed"
		monthly_expiry_log.ExpiryStatusDescription = "loyalty plan does not exist"
		Uc.Write_Loyalty_Monthly_Expiry_log(monthly_expiry_log)
		<-chan_PointsExpiry_Controler
		return
	}
	if planErr != nil {
		monthly_expiry_log.ExpiryStatus = "failed"
		monthly_expiry_log.ExpiryStatusDescription = "type assertion issue with Loyalty_Plan"
		Uc.Write_Loyalty_Monthly_Expiry_log(monthly_expiry_log)
		<-chan_PointsExpiry_Controler
		return
	}
	monthly_expiry_log.Expiry_Rules_Key = plan.Expiry_Rules_Key
	//validate expiry rules
	if plan.Expiry_Rules_Key == "" {
		monthly_expiry_log.ExpiryStatus = "failed"
		monthly_expiry_log.ExpiryStatusDescription = "points expiry rules not defined"
		Uc.Write_Loyalty_Monthly_Expiry_log(monthly_expiry_log)
		<-chan_PointsExpiry_Controler
		return
	}
	for _, pointKey := range account.Points_Detail_Keys {
		pointsDetail, err := Uc.Customer_Loyalty_Account_Points_Details_Get(pointKey)
		if err != nil {
			monthly_expiry_log.ExpiryStatus = "failed"
			monthly_expiry_log.ExpiryStatusDescription = "points expiry rules not found"
			Uc.Write_Loyalty_Monthly_Expiry_log(monthly_expiry_log)
			<-chan_PointsExpiry_Controler
			return
		}
		year, err := strconv.Atoi(pointsDetail[0].Year_Month[:4])
		if err != nil {
			monthly_expiry_log.ExpiryStatus = "failed"
			monthly_expiry_log.ExpiryStatusDescription = "points expiry rules not found"
			Uc.Write_Loyalty_Monthly_Expiry_log(monthly_expiry_log)
			<-chan_PointsExpiry_Controler
			return
		}
		changedPoints := 0.0
		month, err := strconv.Atoi(pointsDetail[0].Year_Month[4:])
		if err == nil && (year < int(account.Initial_Date.Year()) || (year == account.Initial_Date.Year() && month <= int(account.Initial_Date.Month()))) {
			changedPoints += pointsDetail[0].Available_Points
			monthly_expiry_log.Year_Month = pointsDetail[0].Year_Month
			monthly_expiry_log.Opening_Awarded_Points = pointsDetail[0].Awarded_Points
			monthly_expiry_log.Opening_Redeemed_Points = pointsDetail[0].Redeemed_Points
			monthly_expiry_log.Opening_Available_Points = pointsDetail[0].Available_Points
			monthly_expiry_log.Opening_Expired_Points = pointsDetail[0].Expired_Points
			monthly_expiry_log.Opening_Expired_Points = pointsDetail[0].Expired_Points

			expired_Points := pointsDetail[0].Available_Points
			redeemed_Points := pointsDetail[0].Redeemed_Points

			account.Expired_Points = account.Expired_Points + expired_Points
			account.Redeemed_Expired_Points = account.Redeemed_Expired_Points + redeemed_Points
			account.Expiry_Date = account.Coming_Expiry_Date
			account.Awarded_Points = account.Awarded_Points - (expired_Points + redeemed_Points)
			account.Redeemed_Points = account.Redeemed_Points - redeemed_Points
			account.Available_Points = account.Awarded_Points - account.Redeemed_Points
			//remove detail key from account
			account.Points_Detail_Keys = removeStringFromArray(account.Points_Detail_Keys, pointsDetail[0].Key)

			//update governance expiry
			Uc.Loyalty_Governance_Expiry_Points_Credit(expired_Points)
			{
				putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
				if _, putErr := Mdb_Customer_Loyalty_Account.Coll.UpdateOne(putCtx, bson.M{"Key": account.Key}, bson.M{"$set": account}, options.UpdateOne().SetUpsert(true)); putErr != nil {
					log.Println("Mdb_Customer_Loyalty_Account upsert error:", putErr)
				}
				if putSetErr := redisx.SetJSON(putCtx, RedisClient, account.RedisKey(), account); putSetErr != nil {
					log.Println("redisx.SetJSON Customer_Loyalty_Account error:", putSetErr)
				}
				putCancel()
			}

			//update logs
			monthly_expiry_log.End_Awarded_Points = 0
			monthly_expiry_log.End_Redeemed_Points = 0
			monthly_expiry_log.End_Available_Points = 0
			monthly_expiry_log.End_Expired_Points = expired_Points
			//check level downgrade
			new_Loyalty_level_key, errNL := Uc.EvaluateAndUpdate_CustomerLoyaltyLevel("Points_Expiry", account.Key)
			if errNL != nil {
				fmt.Println("error", errNL)
			}
			if errNL == nil {
				monthly_expiry_log.EndLoyaltyLevel = new_Loyalty_level_key
			}

			monthly_expiry_log.ExpiryStatus = "successful"
			monthly_expiry_log.ExpiryStatusDescription = "Cycle expiry"
			Uc.Write_Loyalty_Monthly_Expiry_log(monthly_expiry_log)
			{
				delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
				if _, delErr := Mdb_Customer_Loyalty_Account_Points_Detail.Coll.DeleteOne(delCtx, bson.M{"Key": pointsDetail[0].Key}); delErr != nil {
					log.Println("Mdb_Customer_Loyalty_Account_Points_Detail DeleteOne error:", delErr)
				}
				if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Customer_Loyalty_Account_Points_Detail{Key: pointsDetail[0].Key}.RedisKey()); delRedisErr != nil {
					log.Println("redisx.DelJSON Customer_Loyalty_Account_Points_Detail error:", delRedisErr)
				}
				delCancel()
			}

		}
		entry, entryErr := redisx.GetJSON[Customer_Loyalty_Account](context.Background(), RedisClient, Customer_Loyalty_Account{Key: account.Key}.RedisKey())
		if redisx.IsNil(entryErr) {
			err = errors.New("key does not exist")
		} else if entryErr != nil {
			err = errors.New("error in type assertion")
		}

		expiry_log.EndLoyaltyLevel = entry.Loyalty_Level_Key
		expiry_log.End_Available_Points = entry.Available_Points
		expiry_log.End_Awarded_Points = entry.Awarded_Points
		expiry_log.End_Expired_Points = entry.Expired_Points
		expiry_log.End_Redeemed_Points = entry.Redeemed_Points
		expiry_log.End_Outstanding_Points = entry.Outstanding_fraction_points
		expiry_log.End_Redeemed_Points = entry.Redeemed_Points
		expiry_log.End_Redeemed_Expired_Points = entry.Redeemed_Expired_Points
		expiry_log.ExpiryAmount = changedPoints
		Uc.Write_Loyalty_Full_Expiry_log(expiry_log)

	}
	<-chan_PointsExpiry_Controler
}

func removeStringFromArray(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func (Uc *UserControl) LoyaltyGovernancePools_Metrics_Process() {
	exec := 0
	for range time.Tick(time.Second * 15) {
		if exec == 0 {
			exec = 1
			loyalty_governance, lgErr := redisx.GetJSON[Loyalty_Governance](context.Background(), RedisClient, Loyalty_Governance{Key: LOYALTY_GOVERNANCE_KEY}.RedisKey())
			if redisx.IsNil(lgErr) {
				log.Println("loyalty governance entry not found")
				exec = 0
				continue
			}
			if lgErr != nil {
				log.Println("loyalty governance type assertion issue")
				exec = 0
				continue
			}
			LoyaltyGovernancePools.With(prometheus.Labels{"Pool": "Available"}).Set(loyalty_governance.Available_Points_Pool)
			LoyaltyGovernancePools.With(prometheus.Labels{"Pool": "Distributed"}).Set(loyalty_governance.Distributed_Points_Pool)
			LoyaltyGovernancePools.With(prometheus.Labels{"Pool": "Redeemed"}).Set(loyalty_governance.Redeemed_Points_Pool)
			LoyaltyGovernancePools.With(prometheus.Labels{"Pool": "Expired"}).Set(loyalty_governance.Expired_Points_Pool)
			exec = 0
		}
	}

}

func Normalize_International_MSISDN(MSISDN string) (N_MSISDN string) {
	if len(MSISDN) < Configuration.MSISDN_Short_len {
		return ""
	} else {
		//check if msisdn contain character
		_, err := strconv.Atoi(MSISDN)
		if err != nil {
			return ""
		}
		//normalize
		if len(MSISDN) == len(Configuration.CountryCode)+Configuration.MSISDN_Short_len {
			return MSISDN
		} else if len(MSISDN) == Configuration.MSISDN_Short_len {
			return Configuration.CountryCode + MSISDN
		} else {
			return Configuration.CountryCode + MSISDN[len(MSISDN)-Configuration.MSISDN_Short_len:]
		}
	}
}

func (Uc *UserControl) Auto_GetLoyaltySubsSummary() {
	Uc.GetLoyaltySubsSummary()
	for range time.Tick(time.Second * 300) {
		Uc.GetLoyaltySubsSummary()
	}
}

func (Uc *UserControl) GetLoyaltySubsSummary() (err error) {
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id": bson.M{
					"Loyalty_Level_Key": "$Loyalty_Level_Key",
				},
				"Accounts_Count":         bson.M{"$sum": 1},
				"Total_Awarded_Points":   bson.M{"$sum": "$Awarded_Points"},
				"Total_Redeemed_Points":  bson.M{"$sum": "$Redeemed_Points"},
				"Total_Available_Points": bson.M{"$sum": "$Available_Points"},
			},
		},
	}

	aggCtx, aggCancel := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := Mdb_Customer_Loyalty_Account.Coll.Aggregate(aggCtx, pipeline)
	aggCancel()
	if err != nil {
		log.Println("Error in GetLoyaltySubsSummary: ", err)
		return
	}
	defer cur.Close(context.Background())
	//var output []AlarmsDailyByType
	LoyaltySubsSummary.Reset()
	for cur.Next(context.Background()) {
		var entry Loyalty_Subs_Summary
		err := cur.Decode(&entry)
		if err != nil {
			log.Println("Error in GetLoyaltySubsSummary: ", err)
			return err
		}
		//log.Println(entry)
		LoyaltySubsSummary.With(prometheus.Labels{"Level": entry.ID.Loyalty_Level_Key, "Description": "Accounts_Count"}).Set(entry.Accounts_Count)
		LoyaltySubsSummary.With(prometheus.Labels{"Level": entry.ID.Loyalty_Level_Key, "Description": "Total_Awarded_Points"}).Set(entry.Total_Awarded_Points)
		LoyaltySubsSummary.With(prometheus.Labels{"Level": entry.ID.Loyalty_Level_Key, "Description": "Total_Redeemed_Points"}).Set(entry.Total_Redeemed_Points)
		LoyaltySubsSummary.With(prometheus.Labels{"Level": entry.ID.Loyalty_Level_Key, "Description": "Total_Available_Points"}).Set(entry.Total_Available_Points)
	}
	return
}

// ***********************************************************************
// Loyalty Campaign
// ***********************************************************************
func (Uc *UserControl) Loyalty_Campaign_Add(Login string, request Loyalty_Campaign_AddRequest) (Id int64, err error) {
	//check key if filled and if already used
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	//check if key already used
	{
		chkCtx, chkCancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, chkErr := redisx.GetJSON[Loyalty_Campaign](chkCtx, RedisClient, Loyalty_Campaign{Key: request.Key}.RedisKey())
		chkCancel()
		if chkErr == nil {
			err = errors.New("key already exist")
			return Id, err
		}
	}

	//Prepare new entry
	var NewEntry Loyalty_Campaign
	{
		aiCtx, aiCancel := context.WithTimeout(context.Background(), 10*time.Second)
		NewEntry.Campaign_Id, err = redisx.NextAutoIncrementID(aiCtx, RedisClient, Mdb_Loyalty_AutoIncrement.Coll, "Loyalty_Campaign-Id", redisx.NextIDOptions{
			RedisBase: "AutoIncrement", MongoRetries: 3, RetryBackoff: 500 * time.Millisecond,
		})
		aiCancel()
		if err != nil {
			return Id, err
		}
		Id = NewEntry.Campaign_Id
	}
	NewEntry.Key = request.Key
	NewEntry.Description = request.Description
	NewEntry.Add_Date = time.Now()
	NewEntry.Start_Date = request.Start_Date
	NewEntry.End_Date = request.End_Date
	NewEntry.Campaign_Status = "Created"
	NewEntry.Campaign_Status_Date = time.Now()
	NewEntry.Campaign_Status_User = Login
	NewEntry.Target_All_Subs = request.Target_All_Subs
	NewEntry.Target_Level_Key = request.Target_Level_Key
	NewEntry.Target_Segment_Key = request.Target_Segment_Key
	NewEntry.LoyaltyPoints_From = request.LoyaltyPoints_From
	NewEntry.LoyaltyPoints_Till = request.LoyaltyPoints_Till
	NewEntry.AON_From = request.AON_From
	NewEntry.AON_Till = request.AON_Till
	NewEntry.ARPU_From = request.ARPU_From
	NewEntry.ARPU_Till = request.ARPU_Till
	NewEntry.Welcome_Points_Award_type = request.Welcome_Points_Award_type
	NewEntry.Welcome_Points_Frequency = request.Welcome_Points_Frequency
	NewEntry.Welcome_Points_Max_Award = request.Welcome_Points_Max_Award
	NewEntry.Welcome_Points = request.Welcome_Points
	NewEntry.MobileAppDaily_Login_Award_type = request.MobileAppDaily_Login_Award_type
	NewEntry.MobileAppDaily_Login_Frequency = request.MobileAppDaily_Login_Frequency
	NewEntry.MobileAppDaily_Login_Max_Award = request.MobileAppDaily_Login_Max_Award
	NewEntry.MobileAppDaily_Login = request.MobileAppDaily_Login
	NewEntry.MainGSM_Award_type = request.MainGSM_Award_type
	NewEntry.MainGSM_Frequency = request.MainGSM_Frequency
	NewEntry.MainGSM_Max_Award = request.MainGSM_Max_Award
	NewEntry.MainGSM = request.MainGSM
	NewEntry.MM_P2P_Award_type = request.MM_P2P_Award_type
	NewEntry.MM_P2P_Frequency = request.MM_P2P_Frequency
	NewEntry.MM_P2P_Max_Award = request.MM_P2P_Max_Award
	NewEntry.MM_P2P = request.MM_P2P
	NewEntry.MM_CASHIN_Award_type = request.MM_CASHIN_Award_type
	NewEntry.MM_CASHIN_Frequency = request.MM_CASHIN_Frequency
	NewEntry.MM_CASHIN_Max_Award = request.MM_CASHIN_Max_Award
	NewEntry.MM_CASHIN = request.MM_CASHIN
	NewEntry.MM_CASHOUT_Award_type = request.MM_CASHOUT_Award_type
	NewEntry.MM_CASHOUT_Frequency = request.MM_CASHOUT_Frequency
	NewEntry.MM_CASHOUT_Max_Award = request.MM_CASHOUT_Max_Award
	NewEntry.MM_CASHOUT = request.MM_CASHOUT
	NewEntry.MM_MERCHPAY_Award_type = request.MM_MERCHPAY_Award_type
	NewEntry.MM_MERCHPAY_Frequency = request.MM_MERCHPAY_Frequency
	NewEntry.MM_MERCHPAY_Max_Award = request.MM_MERCHPAY_Max_Award
	NewEntry.MM_MERCHPAY = request.MM_MERCHPAY
	NewEntry.MM_BILLPAY_Award_type = request.MM_BILLPAY_Award_type
	NewEntry.MM_BILLPAY_Frequency = request.MM_BILLPAY_Frequency
	NewEntry.MM_BILLPAY_Max_Award = request.MM_BILLPAY_Max_Award
	NewEntry.MM_BILLPAY = request.MM_BILLPAY
	NewEntry.MM_RC_Award_type = request.MM_RC_Award_type
	NewEntry.MM_RC_Frequency = request.MM_RC_Frequency
	NewEntry.MM_RC_Max_Award = request.MM_RC_Max_Award
	NewEntry.MM_RC = request.MM_RC
	NewEntry.MM_CTMMOREQ_Award_type = request.MM_CTMMOREQ_Award_type
	NewEntry.MM_Frequency = request.MM_Frequency
	NewEntry.MM_Max_Award = request.MM_Max_Award
	NewEntry.MM_CTMMOREQ = request.MM_CTMMOREQ
	NewEntry.MM_CBWREQ_Award_type = request.MM_CBWREQ_Award_type
	NewEntry.MM_CBWREQ_Frequency = request.MM_CBWREQ_Frequency
	NewEntry.MM_CBWREQ_Max_Award = request.MM_CBWREQ_Max_Award
	NewEntry.MM_CBWREQ = request.MM_CBWREQ
	NewEntry.Invitation_SMS_Sender = request.Invitation_SMS_Sender
	NewEntry.Invitation_SMS_Text = request.Invitation_SMS_Text
	NewEntry.PointsAward_SMS_Text = request.PointsAward_SMS_Text

	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Loyalty_Campaign.Coll.UpdateOne(putCtx, bson.M{"Key": NewEntry.Key}, bson.M{"$set": NewEntry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Loyalty_Campaign upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, NewEntry.RedisKey(), NewEntry); putSetErr != nil {
			log.Println("redisx.SetJSON Loyalty_Campaign error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Campaign",
		Event_ActionType:   "Add",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  NewEntry,
	})
	return Id, nil
}

func (Uc *UserControl) Loyalty_Campaign_Edit(Login string, request Loyalty_Campaign_EditRequest) (Id int64, err error) {
	//check and validate
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	entry, entryErr := redisx.GetJSON[Loyalty_Campaign](context.Background(), RedisClient, Loyalty_Campaign{Key: request.Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("key is not created")
		return Id, err
	}
	if entryErr != nil {
		return Id, errors.New("error in type assertion")
	}
	if entry.Campaign_Id != request.Campaign_Id {
		return Id, errors.New("id is not matching")
	}
	Current_Entry := entry

	//Prepare new entry
	entry.Key = request.Key
	entry.Description = request.Description
	entry.Start_Date = request.Start_Date
	entry.End_Date = request.End_Date
	if entry.Campaign_Status != request.Campaign_Status {
		entry.Campaign_Status = request.Campaign_Status
		entry.Campaign_Status_Date = time.Now()
		entry.Campaign_Status_User = Login
	}
	entry.Target_All_Subs = request.Target_All_Subs
	entry.Target_Level_Key = request.Target_Level_Key
	entry.Target_Segment_Key = request.Target_Segment_Key
	entry.LoyaltyPoints_From = request.LoyaltyPoints_From
	entry.LoyaltyPoints_Till = request.LoyaltyPoints_Till
	entry.AON_From = request.AON_From
	entry.AON_Till = request.AON_Till
	entry.ARPU_From = request.ARPU_From
	entry.ARPU_Till = request.ARPU_Till
	entry.Welcome_Points_Award_type = request.Welcome_Points_Award_type
	entry.Welcome_Points_Frequency = request.Welcome_Points_Frequency
	entry.Welcome_Points_Max_Award = request.Welcome_Points_Max_Award
	entry.Welcome_Points = request.Welcome_Points
	entry.MobileAppDaily_Login_Award_type = request.MobileAppDaily_Login_Award_type
	entry.MobileAppDaily_Login_Frequency = request.MobileAppDaily_Login_Frequency
	entry.MobileAppDaily_Login_Max_Award = request.MobileAppDaily_Login_Max_Award
	entry.MobileAppDaily_Login = request.MobileAppDaily_Login
	entry.MainGSM_Award_type = request.MainGSM_Award_type
	entry.MainGSM_Frequency = request.MainGSM_Frequency
	entry.MainGSM_Max_Award = request.MainGSM_Max_Award
	entry.MainGSM = request.MainGSM
	entry.MM_P2P_Award_type = request.MM_P2P_Award_type
	entry.MM_P2P_Frequency = request.MM_P2P_Frequency
	entry.MM_P2P_Max_Award = request.MM_P2P_Max_Award
	entry.MM_P2P = request.MM_P2P
	entry.MM_CASHIN_Award_type = request.MM_CASHIN_Award_type
	entry.MM_CASHIN_Frequency = request.MM_CASHIN_Frequency
	entry.MM_CASHIN_Max_Award = request.MM_CASHIN_Max_Award
	entry.MM_CASHIN = request.MM_CASHIN
	entry.MM_CASHOUT_Award_type = request.MM_CASHOUT_Award_type
	entry.MM_CASHOUT_Frequency = request.MM_CASHOUT_Frequency
	entry.MM_CASHOUT_Max_Award = request.MM_CASHOUT_Max_Award
	entry.MM_CASHOUT = request.MM_CASHOUT
	entry.MM_MERCHPAY_Award_type = request.MM_MERCHPAY_Award_type
	entry.MM_MERCHPAY_Frequency = request.MM_MERCHPAY_Frequency
	entry.MM_MERCHPAY_Max_Award = request.MM_MERCHPAY_Max_Award
	entry.MM_MERCHPAY = request.MM_MERCHPAY
	entry.MM_BILLPAY_Award_type = request.MM_BILLPAY_Award_type
	entry.MM_BILLPAY_Frequency = request.MM_BILLPAY_Frequency
	entry.MM_BILLPAY_Max_Award = request.MM_BILLPAY_Max_Award
	entry.MM_BILLPAY = request.MM_BILLPAY
	entry.MM_RC_Award_type = request.MM_RC_Award_type
	entry.MM_RC_Frequency = request.MM_RC_Frequency
	entry.MM_RC_Max_Award = request.MM_RC_Max_Award
	entry.MM_RC = request.MM_RC
	entry.MM_CTMMOREQ_Award_type = request.MM_CTMMOREQ_Award_type
	entry.MM_Frequency = request.MM_Frequency
	entry.MM_Max_Award = request.MM_Max_Award
	entry.MM_CTMMOREQ = request.MM_CTMMOREQ
	entry.MM_CBWREQ_Award_type = request.MM_CBWREQ_Award_type
	entry.MM_CBWREQ_Frequency = request.MM_CBWREQ_Frequency
	entry.MM_CBWREQ_Max_Award = request.MM_CBWREQ_Max_Award
	entry.MM_CBWREQ = request.MM_CBWREQ
	entry.Invitation_SMS_Sender = request.Invitation_SMS_Sender
	entry.Invitation_SMS_Text = request.Invitation_SMS_Text
	entry.PointsAward_SMS_Text = request.PointsAward_SMS_Text
	if request.NewKey != "" {
		if request.NewKey != request.Key {
			//delete old
			{
				delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
				if _, delErr := Mdb_Loyalty_Campaign.Coll.DeleteOne(delCtx, bson.M{"Key": request.Key}); delErr != nil {
					log.Println("Mdb_Loyalty_Campaign DeleteOne error:", delErr)
				}
				if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Loyalty_Campaign{Key: request.Key}.RedisKey()); delRedisErr != nil {
					log.Println("redisx.DelJSON Loyalty_Campaign error:", delRedisErr)
				}
				delCancel()
			}
			//update key
			entry.Key = request.NewKey
		}
	}
	//add to cache and DB
	{
		putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, putErr := Mdb_Loyalty_Campaign.Coll.UpdateOne(putCtx, bson.M{"Key": entry.Key}, bson.M{"$set": entry}, options.UpdateOne().SetUpsert(true)); putErr != nil {
			log.Println("Mdb_Loyalty_Campaign upsert error:", putErr)
		}
		if putSetErr := redisx.SetJSON(putCtx, RedisClient, entry.RedisKey(), entry); putSetErr != nil {
			log.Println("redisx.SetJSON Loyalty_Campaign error:", putSetErr)
		}
		putCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Campaign",
		Event_ActionType:   "Edit",
		Event_Description:  "",
		Event_Entry_Before: Current_Entry,
		Event_Entry_After:  entry,
	})
	return Id, nil
}

func (Uc *UserControl) Loyalty_Campaign_Get(Key string) (entries []Loyalty_Campaign, err error) {
	if Key == "" {
		{
			scanCtx, scanCancel := context.WithTimeout(context.Background(), 10*time.Second)
			entries, err = redisx.GetAllJSONByPattern[Loyalty_Campaign](scanCtx, RedisClient, redisx.ScanJSONOptions{Pattern: "Loyalty_Campaign:*"})
			scanCancel()
		}
		return entries, err
	} else {
		entry, getErr := redisx.GetJSON[Loyalty_Campaign](context.Background(), RedisClient, Loyalty_Campaign{Key: Key}.RedisKey())
		if redisx.IsNil(getErr) {
			err = errors.New("key does not exist")
			return entries, err
		}
		if getErr != nil {
			err = errors.New("error in type assertion")
			return entries, err
		}
		entries = append(entries, entry)
		return entries, nil
	}
}

func (Uc *UserControl) Loyalty_Campaign_GetPaginated(Page, Limit int64) (entries []Loyalty_Campaign, err error) {
	if Page < 1 {
		return entries, errors.New("invalid page")
	}
	if Limit < 1 || Limit > 50000 {
		return entries, errors.New("invalid limit (accept value between 1 and 50000)")
	}

	pgCtx, pgCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer pgCancel()
	result, err := mongox.FindPage[Loyalty_Campaign](pgCtx, Mdb_Loyalty_Campaign.Coll, bson.M{}, mongox.PageRequest{
		Page: Page, PageSize: Limit, MaxPageSize: 50000,
	})
	if err != nil {
		return entries, err
	}
	entries = result.Items
	return entries, nil

}

func (Uc *UserControl) Loyalty_Campaign_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry, entryErr := redisx.GetJSON[Loyalty_Campaign](context.Background(), RedisClient, Loyalty_Campaign{Key: Key}.RedisKey())
	if redisx.IsNil(entryErr) {
		err = errors.New("entry does not exist")
		return err
	}
	if entryErr != nil {
		err = errors.New("error in type assertion")
		return err
	}
	{
		delCtx, delCancel := context.WithTimeout(context.Background(), 10*time.Second)
		if _, delErr := Mdb_Loyalty_Campaign.Coll.DeleteOne(delCtx, bson.M{"Key": Key}); delErr != nil {
			log.Println("Mdb_Loyalty_Campaign DeleteOne error:", delErr)
		}
		if _, delRedisErr := redisx.DelJSON(delCtx, RedisClient, Loyalty_Campaign{Key: Key}.RedisKey()); delRedisErr != nil {
			log.Println("redisx.DelJSON Loyalty_Campaign error:", delRedisErr)
		}
		delCancel()
	}
	//add logs
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Loyalty_Campaign",
		Event_ActionType:   "Delete",
		Event_Description:  "",
		Event_Entry_Before: nil,
		Event_Entry_After:  entry,
	})
	return nil
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////
// /////SEND SMS////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////
func SendSMS(Sender string, target string, SMSText string) (_rErr error) {
	log.Println("Sending SMS: Sender (" + Sender + "), Target (" + target + "), text (" + SMSText + ") ")
	requrl := "http://" + Configuration.SMPP.IP + ":" + Configuration.SMPP.Port + "/?systemid=" + Configuration.SMPP.Login + "&password=" + url.QueryEscape(Configuration.SMPP.Password) + "&Originator=" + Sender + "&dest_addr=" + target + "&msg_text=" + url.QueryEscape(SMSText) + "&encoding=1&ston=5&snpi=0&dton=1&registered_delivery=0"
	//-------------- Encoding used in DRC and GM Start
	//"&ston=5&snpi=0&dton=1&dnpi=1&encoding=1"
	//-------------- Encoding used in DRC and GM End
	method := "GET"
	if Configuration.Operation == "Angola" {
		requrl = "http://" + Configuration.SMPP.IP + ":" + Configuration.SMPP.Port + "/sendsms?username=" + Configuration.SMPP.Login + "&password=" + Configuration.SMPP.Password + "&from=" + Sender + "&to=" + target + "&text=" + url.QueryEscape(SMSText) + "&coding=2"
	}
	req, err := http.NewRequest(method, requrl, nil)
	if err != nil {
		log.Println("Error sending SMS: ", err)
		return err
	}
	//client := &http.Client{}
	client := &http.Client{
		Timeout: 15 * time.Second, //is SMSC not reachable request will time out after 15 sec
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending SMS: ", err)
		return err
	}
	defer resp.Body.Close()
	if Configuration.Operation == "Angola" || resp.StatusCode == 200 {
		return nil
	} else {
		err := errors.New("error sending SMS: " + strconv.Itoa(resp.StatusCode))
		return err
	}
}
