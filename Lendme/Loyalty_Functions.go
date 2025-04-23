package Lendme

import (
	"daoc"
	"errors"
	"log"
	"net/http"
	"reflect"
	"time"

	PropC "afr_propylaea/PropylaeaClient"
	"afr_unified_charging_gateway/Unified_charging_gateway_Client"

	"github.com/prometheus/client_golang/prometheus"
)

var Map_Loyalty_AutoIncrement daoc.Cache_Synch
var DAO_Loyalty_AutoIncrement daoc.DAO

var Map_Loyalty_Governance daoc.Cache_Synch
var DAO_Loyalty_Governance daoc.DAO

var Map_Loyalty_Level daoc.Cache_Synch
var DAO_Loyalty_Level daoc.DAO
var DAO_Loyalty_Level_Change_log daoc.DAO

var Map_Loyalty_Account_Segment daoc.Cache_Synch
var DAO_Loyalty_Account_Segment daoc.DAO

var Map_Loyalty_Point_Earning_Rules daoc.Cache_Synch
var DAO_Loyalty_Point_Earning_Rules daoc.DAO

var Map_Loyalty_Point_Expiry_Rules daoc.Cache_Synch
var DAO_Loyalty_Point_Expiry_Rules daoc.DAO

var Map_Loyalty_Point_Redemption_Rules daoc.Cache_Synch
var DAO_Loyalty_Point_Redemption_Rules daoc.DAO

var Map_Loyalty_Plan daoc.Cache_Synch
var DAO_Loyalty_Plan daoc.DAO

var Map_Customer_Loyalty_Account daoc.Cache_Synch
var DAO_Customer_Loyalty_Account daoc.DAO

var Map_Customer_DND daoc.Cache_Synch
var DAO_Customer_DND daoc.DAO

var Map_Customer_Exclusion daoc.Cache_Synch
var DAO_Customer_Exclusion daoc.DAO

var Map_Customer_COS_Exclusion daoc.Cache_Synch
var DAO_Customer_COS_Exclusion daoc.DAO

var Map_Customer_UAT daoc.Cache_Synch
var DAO_Customer_UAT daoc.DAO

var DAO_Loyalty_Event_Log daoc.DAO

var DAO_Loyalty_Expiry_log daoc.DAO
var DAO_Loyalty_Redemption_log daoc.DAO

var DAO_Loyalty_AccountCreditPoints_log daoc.DAO
var DAO_Loyalty_AccountDebitPoints_log daoc.DAO

var chan_LoyaltyGovernance_Controler = make(chan int, 1)

var chan_PointsExpiry_Controler = make(chan int, 50)

var LOYALTY_GOVERNANCE_KEY string = "Loyalty_Governance"

func (uc *UserControl) InitializeLoyaltyCache() {
	var AutoIncr daoc.AutoIncrement
	Map_Loyalty_AutoIncrement.Initialize("Loyalty_AutoIncrement", "AutoIncrement", reflect.TypeOf(daoc.AutoIncrement{}), AutoIncr, true, &DAO_Loyalty_AutoIncrement, uc.CacheDir.List)
	var loyalty_Governance Loyalty_Governance
	Map_Loyalty_Governance.Initialize("Loyalty_Governance", "Loyalty_Governance", reflect.TypeOf(Loyalty_Governance{}), loyalty_Governance, true, &DAO_Loyalty_Governance, uc.CacheDir.List)
	var loyalty_Level Loyalty_Level
	Map_Loyalty_Level.Initialize("Loyalty_Level", "Loyalty_Level", reflect.TypeOf(Loyalty_Level{}), loyalty_Level, true, &DAO_Loyalty_Level, uc.CacheDir.List)
	var loyalty_Account_Segment Loyalty_Account_Segment
	Map_Loyalty_Account_Segment.Initialize("Loyalty_Account_Segment", "Loyalty_Account_Segment", reflect.TypeOf(Loyalty_Account_Segment{}), loyalty_Account_Segment, true, &DAO_Loyalty_Account_Segment, uc.CacheDir.List)
	var loyalty_Point_Earning_Rules Loyalty_Point_Earning_Rules
	Map_Loyalty_Point_Earning_Rules.Initialize("Loyalty_Point_Earning_Rules", "Loyalty_Point_Earning_Rules", reflect.TypeOf(Loyalty_Point_Earning_Rules{}), loyalty_Point_Earning_Rules, true, &DAO_Loyalty_Point_Earning_Rules, uc.CacheDir.List)
	var loyalty_Point_Expiry_Rules Loyalty_Point_Expiry_Rules
	Map_Loyalty_Point_Expiry_Rules.Initialize("Loyalty_Point_Expiry_Rules", "Loyalty_Point_Expiry_Rules", reflect.TypeOf(Loyalty_Point_Expiry_Rules{}), loyalty_Point_Expiry_Rules, true, &DAO_Loyalty_Point_Expiry_Rules, uc.CacheDir.List)
	var loyalty_Point_Redemption_Rules Loyalty_Point_Redemption_Rules
	Map_Loyalty_Point_Redemption_Rules.Initialize("Loyalty_Point_Redemption_Rules", "Loyalty_Point_Redemption_Rules", reflect.TypeOf(Loyalty_Point_Redemption_Rules{}), loyalty_Point_Redemption_Rules, true, &DAO_Loyalty_Point_Redemption_Rules, uc.CacheDir.List)
	var loyalty_Plan Loyalty_Plan
	Map_Loyalty_Plan.Initialize("Loyalty_Plan", "Loyalty_Plan", reflect.TypeOf(Loyalty_Plan{}), loyalty_Plan, true, &DAO_Loyalty_Plan, uc.CacheDir.List)
	var customer_Loyalty_Account Customer_Loyalty_Account
	Map_Customer_Loyalty_Account.Initialize("Customer_Loyalty_Account", "Customer_Loyalty_Account", reflect.TypeOf(Customer_Loyalty_Account{}), customer_Loyalty_Account, true, &DAO_Customer_Loyalty_Account, uc.CacheDir.List)
	var customer_DND Customer_DND
	Map_Customer_DND.Initialize("Customer_DND", "Customer_DND", reflect.TypeOf(Customer_DND{}), customer_DND, true, &DAO_Customer_DND, uc.CacheDir.List)
	var customer_Exclusion Customer_Exclusion
	Map_Customer_Exclusion.Initialize("Customer_Exclusion", "Customer_Exclusion", reflect.TypeOf(Customer_Exclusion{}), customer_Exclusion, true, &DAO_Customer_Exclusion, uc.CacheDir.List)
	var customer_COS_Exclusion Customer_COS_Exclusion
	Map_Customer_COS_Exclusion.Initialize("Customer_COS_Exclusion", "Customer_COS_Exclusion", reflect.TypeOf(Customer_COS_Exclusion{}), customer_COS_Exclusion, true, &DAO_Customer_COS_Exclusion, uc.CacheDir.List)
	var customer_UAT Customer_UAT
	Map_Customer_UAT.Initialize("Customer_UAT", "Customer_UAT", reflect.TypeOf(Customer_UAT{}), customer_UAT, true, &DAO_Customer_UAT, uc.CacheDir.List)

}

func (uc *UserControl) InitializeLoyaltyDAO() {
	DAO_Loyalty_AutoIncrement.Initialize("Loyalty_AutoIncrement", uc.LoyaltyMongoDB.MongoDBClient, reflect.TypeOf(daoc.AutoIncrement{}), Configuration.DB_Name_Loyalty, "Col_Loyalty_AutoIncrement", "")
	DAO_Loyalty_Governance.Initialize("Loyalty_Governance", uc.LoyaltyMongoDB.MongoDBClient, reflect.TypeOf(Loyalty_Governance{}), Configuration.DB_Name_Loyalty, "Col_Loyalty_Governance", "")
	DAO_Loyalty_Level.Initialize("Loyalty_Level", uc.LoyaltyMongoDB.MongoDBClient, reflect.TypeOf(Loyalty_Level{}), Configuration.DB_Name_Loyalty, "Col_Loyalty_Level", "")
	DAO_Loyalty_Level_Change_log.Initialize("Loyalty_Level_Change_log", uc.LoyaltyMongoDB.MongoDBClient, reflect.TypeOf(Loyalty_Level_Change_log{}), Configuration.DB_Name_Loyalty, "Col_Loyalty_Level_Change_log", "")
	DAO_Loyalty_Account_Segment.Initialize("Loyalty_Account_Segment", uc.LoyaltyMongoDB.MongoDBClient, reflect.TypeOf(Loyalty_Account_Segment{}), Configuration.DB_Name_Loyalty, "Col_Loyalty_Account_Segment", "")
	DAO_Loyalty_Point_Earning_Rules.Initialize("Loyalty_Point_Earning_Rules", uc.LoyaltyMongoDB.MongoDBClient, reflect.TypeOf(Loyalty_Point_Earning_Rules{}), Configuration.DB_Name_Loyalty, "Col_Loyalty_Point_Earning_Rules", "")
	DAO_Loyalty_Point_Expiry_Rules.Initialize("Loyalty_Point_Expiry_Rules", uc.LoyaltyMongoDB.MongoDBClient, reflect.TypeOf(Loyalty_Point_Expiry_Rules{}), Configuration.DB_Name_Loyalty, "Col_Loyalty_Point_Expiry_Rules", "")
	DAO_Loyalty_Point_Redemption_Rules.Initialize("Loyalty_Point_Redemption_Rules", uc.LoyaltyMongoDB.MongoDBClient, reflect.TypeOf(Loyalty_Point_Redemption_Rules{}), Configuration.DB_Name_Loyalty, "Col_Loyalty_Point_Redemption_Rules", "")
	DAO_Loyalty_Plan.Initialize("Loyalty_Plan", uc.LoyaltyMongoDB.MongoDBClient, reflect.TypeOf(Loyalty_Plan{}), Configuration.DB_Name_Loyalty, "Col_Loyalty_Plan", "")
	DAO_Customer_Loyalty_Account.Initialize("Customer_Loyalty_Account", uc.LoyaltyMongoDB.MongoDBClient, reflect.TypeOf(Customer_Loyalty_Account{}), Configuration.DB_Name_Loyalty, "Col_Customer_Loyalty_Account", "")
	DAO_Loyalty_Event_Log.Initialize("Loyalty_Event_Log", uc.LoyaltyMongoDB.MongoDBClient, reflect.TypeOf(Loyalty_Event_Log{}), Configuration.DB_Name_Loyalty, "Col_Loyalty_Event_Log", "")
	DAO_Customer_DND.Initialize("Customer_DND", uc.LoyaltyMongoDB.MongoDBClient, reflect.TypeOf(Customer_DND{}), Configuration.DB_Name_Loyalty, "Col_Customer_DND", "")
	DAO_Customer_Exclusion.Initialize("Customer_Exclusion", uc.LoyaltyMongoDB.MongoDBClient, reflect.TypeOf(Customer_Exclusion{}), Configuration.DB_Name_Loyalty, "Col_Customer_Exclusion", "")
	DAO_Customer_COS_Exclusion.Initialize("Customer_COS_Exclusion", uc.LoyaltyMongoDB.MongoDBClient, reflect.TypeOf(Customer_COS_Exclusion{}), Configuration.DB_Name_Loyalty, "Col_Customer_COS_Exclusion", "")
	DAO_Customer_UAT.Initialize("Customer_UAT", uc.LoyaltyMongoDB.MongoDBClient, reflect.TypeOf(Customer_UAT{}), Configuration.DB_Name_Loyalty, "Col_Customer_UAT", "")
	DAO_Loyalty_Expiry_log.Initialize("Loyalty_Expiry_log", uc.LoyaltyMongoDB.MongoDBClient, reflect.TypeOf(Loyalty_Expiry_log{}), Configuration.DB_Name_Loyalty, "Col_Loyalty_Expiry_log", "")
	DAO_Loyalty_Redemption_log.Initialize("Loyalty_Redemption_log", uc.LoyaltyMongoDB.MongoDBClient, reflect.TypeOf(Loyalty_Redemption_log{}), Configuration.DB_Name_Loyalty, "Col_Loyalty_Redemption_log", "")
	DAO_Loyalty_AccountCreditPoints_log.Initialize("Loyalty_AccountCreditPoints_log", uc.LoyaltyMongoDB.MongoDBClient, reflect.TypeOf(Loyalty_AccountCreditPoints_log{}), Configuration.DB_Name_Loyalty, "Col_Loyalty_AccountCreditPoints_log", "")
	DAO_Loyalty_AccountDebitPoints_log.Initialize("Loyalty_AccountDebitPoints_log", uc.LoyaltyMongoDB.MongoDBClient, reflect.TypeOf(Loyalty_AccountDebitPoints_log{}), Configuration.DB_Name_Loyalty, "Col_Loyalty_AccountDebitPoints_log", "")

}

func (uc *UserControl) LoyaltyIndexesMaintenanceProcess() {
	log.Println("Loyalty DB index manintenance process started...")
	exists, err := DAO_Loyalty_AutoIncrement.CheckAndCreateIndex("Idx_LoyaltyAutoIncrement_Key", []string{"Key"}, true)
	if err != nil {
		log.Println("Error checking and creating index Idx_LoyaltyAutoIncrement_Key: ", err)
	} else {
		if !exists {
			log.Println("Index Idx_LoyaltyAutoIncrement_Key created")
		}
	}

	exists, err = DAO_Loyalty_Governance.CheckAndCreateIndex("Idx_Loyalty_Governance_Key", []string{"Key"}, true)
	if err != nil {
		log.Println("Error creating index Idx_Loyalty_Governance_Key: ", err)
	} else {
		if !exists {
			log.Println("Index Idx_Loyalty_Governance_Key created")
		}
	}

	exists, err = DAO_Loyalty_Level.CheckAndCreateIndex("Idx_Loyalty_Level_Key", []string{"Key"}, true)
	if err != nil {
		log.Println("Error creating index Idx_Loyalty_Level_Key: ", err)
	} else {
		if !exists {
			log.Println("Index Idx_Loyalty_Level_Key created")
		}
	}

	exists, err = DAO_Loyalty_Account_Segment.CheckAndCreateIndex("Idx_Loyalty_Account_Segment_Key", []string{"Key"}, true)
	if err != nil {
		log.Println("Error creating index Idx_Loyalty_Account_Segment_Key: ", err)
	} else {
		if !exists {
			log.Println("Index Idx_Loyalty_Account_Segment_Key created")
		}
	}

	exists, err = DAO_Loyalty_Point_Earning_Rules.CheckAndCreateIndex("Idx_Loyalty_Point_Earning_Rules_Key", []string{"Key"}, true)
	if err != nil {
		log.Println("Error creating index Idx_Loyalty_Point_Earning_Rules_Key: ", err)
	} else {
		if !exists {
			log.Println("Index Idx_Loyalty_Point_Earning_Rules_Key created")
		}
	}

	exists, err = DAO_Loyalty_Point_Expiry_Rules.CheckAndCreateIndex("Idx_Loyalty_Point_Expiry_Rules_Key", []string{"Key"}, true)
	if err != nil {
		log.Println("Error creating index Idx_Loyalty_Point_Expiry_Rules_Key: ", err)
	} else {
		if !exists {
			log.Println("Index Idx_Loyalty_Point_Expiry_Rules_Key created")
		}
	}

	exists, err = DAO_Loyalty_Point_Redemption_Rules.CheckAndCreateIndex("Idx_Loyalty_Point_Redemption_Rules_Key", []string{"Key"}, true)
	if err != nil {
		log.Println("Error creating index Idx_Loyalty_Point_Redemption_Rules_Key: ", err)
	} else {
		if !exists {
			log.Println("Index Idx_Loyalty_Point_Redemption_Rules_Key created")
		}
	}

	exists, err = DAO_Loyalty_Plan.CheckAndCreateIndex("Idx_Loyalty_Plan_Key", []string{"Key"}, true)
	if err != nil {
		log.Println("Error creating index Idx_Loyalty_Plan_Key: ", err)
	} else {
		if !exists {
			log.Println("Index Idx_Loyalty_Plan_Key created")
		}
	}

	exists, err = DAO_Customer_Loyalty_Account.CheckAndCreateIndex("Idx_Customer_Loyalty_Account_Key", []string{"Key"}, true)
	if err != nil {
		log.Println("Error creating index Idx_Customer_Loyalty_Account_Key: ", err)
	} else {
		if !exists {
			log.Println("Index Idx_Customer_Loyalty_Account_Key created")
		}
	}

	exists, err = DAO_Customer_DND.CheckAndCreateIndex("Idx_Customer_DND_Key", []string{"Key"}, true)
	if err != nil {
		log.Println("Error creating index Idx_Customer_DND_Key: ", err)
	} else {
		if !exists {
			log.Println("Index Idx_Customer_DND_Key created")
		}
	}

	exists, err = DAO_Customer_Exclusion.CheckAndCreateIndex("Idx_Customer_Exclusion_Key", []string{"Key"}, true)
	if err != nil {
		log.Println("Error creating index Idx_Customer_Exclusion_Key: ", err)
	} else {
		if !exists {
			log.Println("Index Idx_Customer_Exclusion_Key created")
		}
	}

	exists, err = DAO_Customer_COS_Exclusion.CheckAndCreateIndex("Idx_Customer_COS_Exclusion_Key", []string{"Key"}, true)
	if err != nil {
		log.Println("Error creating index Idx_Customer_COS_Exclusion_Key: ", err)
	} else {
		if !exists {
			log.Println("Index Idx_Customer_COS_Exclusion_Key created")
		}
	}

	exists, err = DAO_Customer_UAT.CheckAndCreateIndex("Idx_Customer_UAT_Key", []string{"Key"}, true)
	if err != nil {
		log.Println("Error creating index Idx_Customer_UAT_Key: ", err)
	} else {
		if !exists {
			log.Println("Index Idx_Customer_UAT_Key created")
		}
	}
}

func (Uc *UserControl) Write_Loyalty_Event_Log(record Loyalty_Event_Log) {
	YYYY, MM, _, DD, _, _, _ := GetTimeParts(record.Event_Time)
	Db := DAO_Loyalty_Event_Log.DB + "_" + YYYY + MM
	Col := DAO_Loyalty_Event_Log.Collection + "_" + DD
	_, err := DAO_Loyalty_Event_Log.PutOneLogs(record, Db, Col)
	if err != nil {
		log.Println("Error in Write_Lendme_log:", err, " (", record, ")")
		return
	}
}

func (Uc *UserControl) Write_Loyalty_Level_Change_log(record Loyalty_Level_Change_log) {
	YYYY, MM, _, _, _, _, _ := GetTimeParts(record.Level_Change_Date)
	Db := DAO_Loyalty_Level_Change_log.DB + "_" + YYYY + MM
	Col := DAO_Loyalty_Level_Change_log.Collection //+ "_" + DD
	_, err := DAO_Loyalty_Level_Change_log.PutOneLogs(record, Db, Col)
	if err != nil {
		log.Println("Error in Write_Lendme_log:", err, " (", record, ")")
		return
	}
}

func (Uc *UserControl) Write_Loyalty_AccountCreditPoints_log(record Loyalty_AccountCreditPoints_log) {
	YYYY, MM, _, DD, _, _, _ := GetTimeParts(record.ReceiveDate)
	Db := DAO_Loyalty_AccountCreditPoints_log.DB + "_" + YYYY + MM
	Col := DAO_Loyalty_AccountCreditPoints_log.Collection + "_" + DD
	_, err := DAO_Loyalty_AccountCreditPoints_log.PutOneLogs(record, Db, Col)
	if err != nil {
		log.Println("Error in Write_Loyalty_AccountCreditPoints_log:", err, " (", record, ")")
		return
	}
}

func (Uc *UserControl) Write_Loyalty_Expiry_log(record Loyalty_Expiry_log) {
	YYYY, MM, _, DD, _, _, _ := GetTimeParts(record.ExpiryTime)
	Db := DAO_Loyalty_Expiry_log.DB + "_" + YYYY + MM
	Col := DAO_Loyalty_Expiry_log.Collection + "_" + DD
	_, err := DAO_Loyalty_Expiry_log.PutOneLogs(record, Db, Col)
	if err != nil {
		log.Println("Error in Write_Loyalty_Expiry_log:", err, " (", record, ")")
		return
	}
}

func (Uc *UserControl) Write_Loyalty_Redemption_log(record Loyalty_Redemption_log) {
	YYYY, MM, _, DD, _, _, _ := GetTimeParts(record.ReceiveDate)
	Db := DAO_Loyalty_Redemption_log.DB + "_" + YYYY + MM
	Col := DAO_Loyalty_Redemption_log.Collection + "_" + DD
	_, err := DAO_Loyalty_Redemption_log.PutOneLogs(record, Db, Col)
	if err != nil {
		log.Println("Error in Write_Loyalty_Redemption_log:", err, " (", record, ")")
		return
	}
}

func (Uc *UserControl) Write_Loyalty_AccountDebitPoints_log(record Loyalty_AccountDebitPoints_log) {
	YYYY, MM, _, DD, _, _, _ := GetTimeParts(record.ReceiveDate)
	Db := DAO_Loyalty_AccountDebitPoints_log.DB + "_" + YYYY + MM
	Col := DAO_Loyalty_AccountDebitPoints_log.Collection + "_" + DD
	_, err := DAO_Loyalty_AccountDebitPoints_log.PutOneLogs(record, Db, Col)
	if err != nil {
		log.Println("Error in Write_Loyalty_AccountDebitPoints_log:", err, " (", record, ")")
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
		Key:                                   "Default_Earning_Rules",
		Description:                           "Default_Earning_Rules",
		Welcome_Points:                        5,
		MobileAppDaily_Login:                  1,
		MainGSMBalance_AmountConsumedPerPoint: 10,
		MM_P2P_Award_Type:                     "Transaction",
		MM_P2P:                                1,
		MM_CASHIN_Award_Type:                  "Transaction",
		MM_CASHIN:                             1,
		MM_CASHOUT_Award_Type:                 "Transaction",
		MM_CASHOUT:                            1,
		MM_MERCHPAY_Award_Type:                "Transaction",
		MM_MERCHPAY:                           1,
		MM_BILLPAY_Award_Type:                 "Transaction",
		MM_BILLPAY:                            1,
		MM_RC_Award_Type:                      "Amount",
		MM_RC:                                 0.15,
		MM_CTMMOREQ_Award_Type:                "Amount",
		MM_CTMMOREQ:                           0.15,
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
		Key:                               "Default_Redemption_Rules",
		Description:                       "Default_Redemption_Rules",
		Min_Accumulated_Points:            100,
		Allow_Negative_Balance_ToRedeem:   false,
		Allow_PendingLendme_ToRedeem:      false,
		Airtime_MinPoints:                 100,
		Airtime_AmountPerPoint:            0.5,
		MobileMoney_MinPoints:             100,
		MobileMoney_AmountPerPoint:        0.5,
		Bundles_MinPoints:                 100,
		Bundles_Product_Catalogue_Channel: "Loyalty_Default_Channel",
		Bundles_Product_Catalogue_Plan:    "Loyalty_Default_Plan",
		Bundles_Product_Catalogue_Version: "1",
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
	//check if key already used
	exits := Map_Loyalty_Governance.Check(request.Key)
	if exits {
		err = errors.New("key already exist")
		return Id, err
	}

	//Prepare new entry
	var NewEntry Loyalty_Governance
	NewEntry.Governance_Id = Map_Loyalty_AutoIncrement.GetNextAI("Loyalty_Governance-Id")
	Id = NewEntry.Governance_Id
	NewEntry.Key = request.Key
	NewEntry.Available_Points_Pool = request.Available_Points_Pool
	NewEntry.Distributed_Points_Pool = request.Distributed_Points_Pool
	NewEntry.Redeemed_Points_Pool = request.Redeemed_Points_Pool
	NewEntry.MaxAllowedPoints_PerTransaction = request.MaxAllowedPoints_PerTransaction
	NewEntry.MaxSubsAwardedPoints_PerMonth = request.MaxSubsAwardedPoints_PerMonth
	NewEntry.MaxSubsAwardedPoints = request.MaxSubsAwardedPoints
	//add to cache and DB
	Map_Loyalty_Governance.Put(NewEntry.Key, NewEntry)
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
	entry_na, exits := Map_Loyalty_Governance.CheckThenGet(request.Key)
	if !exits {
		err = errors.New("key is not created")
		return Id, err
	}
	entry, ok := entry_na.(Loyalty_Governance)
	if !ok {
		return Id, errors.New("error in type assertion")
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

	if request.NewKey != "" {
		if request.NewKey != request.Key {
			//delete old
			Map_Loyalty_Governance.Delete(request.Key)
			//update key
			entry.Key = request.NewKey
		}
	}
	//add to cache and DB
	Map_Loyalty_Governance.Put(entry.Key, entry)
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
		entries_na := Map_Loyalty_Governance.ConvertToArray()
		if len(entries_na) > 0 {
			for _, entry_na := range entries_na {
				entry, ok := entry_na.(Loyalty_Governance)
				if !ok {
					err = errors.New("error in type assertion")
					return entries, err
				} else {
					entries = append(entries, entry)
				}
			}
		}
		return entries, nil
	} else {
		entry_na, exits := Map_Loyalty_Governance.CheckThenGet(Key)
		if !exits {
			err = errors.New("key does not exist")
			return entries, err
		}
		entry, ok := entry_na.(Loyalty_Governance)
		if !ok {
			err = errors.New("error in type assertion")
			return entries, err
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

	var findparams daoc.DAOFindParams
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
	var paginationparams daoc.DAOPaginate
	paginationparams.Limit = Limit
	paginationparams.Page = Page
	findResult, err := DAO_Loyalty_Governance.FindPaginate(findparams, paginationparams)
	if err != nil {
		return entries, err
	}
	if len(findResult) > 0 {
		for _, findres := range findResult {
			InterfaceValue := reflect.ValueOf(findres).Elem().Interface().(Loyalty_Governance)
			entries = append(entries, InterfaceValue)
		}
	}
	return entries, nil

}

func (Uc *UserControl) Loyalty_Governance_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry_na, exits := Map_Loyalty_Governance.CheckThenGet(Key)
	if !exits {
		err = errors.New("entry does not exist")
		return err
	}
	entry, ok := entry_na.(Loyalty_Governance)
	if !ok {
		err = errors.New("error in type assertion")
		return err
	}
	Map_Loyalty_Governance.Delete(Key)
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

func (Uc *UserControl) Loyalty_Governance_Available_Points_Debit(points float64) (err error) {
	chan_LoyaltyGovernance_Controler <- 1
	loyalty_governance_na, exist := Map_Loyalty_Governance.CheckThenGet(LOYALTY_GOVERNANCE_KEY)
	if !exist {
		<-chan_LoyaltyGovernance_Controler
		return errors.New("loyalty governance entry not found")
	}
	loyalty_governance, ok := loyalty_governance_na.(Loyalty_Governance)
	if !ok {
		<-chan_LoyaltyGovernance_Controler
		return errors.New("loyalty governance type assertion issue")
	}
	if (loyalty_governance.Available_Points_Pool - loyalty_governance.Distributed_Points_Pool) > points {
		loyalty_governance.Distributed_Points_Pool = loyalty_governance.Distributed_Points_Pool + points
		Map_Loyalty_Governance.Put(loyalty_governance.Key, loyalty_governance)
		<-chan_LoyaltyGovernance_Controler
		return
	} else {
		<-chan_LoyaltyGovernance_Controler
		return errors.New("no enough loyalty points to distribute from the governance available balance")
	}
}

func (Uc *UserControl) Loyalty_Governance_Redeem_Points_Debit(points float64) (err error) {
	chan_LoyaltyGovernance_Controler <- 1
	loyalty_governance_na, exist := Map_Loyalty_Governance.CheckThenGet(LOYALTY_GOVERNANCE_KEY)
	if !exist {
		<-chan_LoyaltyGovernance_Controler
		return errors.New("loyalty governance entry not found")
	}
	loyalty_governance, ok := loyalty_governance_na.(Loyalty_Governance)
	if !ok {
		<-chan_LoyaltyGovernance_Controler
		return errors.New("loyalty governance type assertion issue")
	}
	loyalty_governance.Redeemed_Points_Pool = loyalty_governance.Redeemed_Points_Pool + points
	Map_Loyalty_Governance.Put(loyalty_governance.Key, loyalty_governance)
	<-chan_LoyaltyGovernance_Controler
	return
}

func (Uc *UserControl) Loyalty_Governance_Expiry_Points_Credit(points float64) (err error) {
	chan_LoyaltyGovernance_Controler <- 1
	loyalty_governance_na, exist := Map_Loyalty_Governance.CheckThenGet(LOYALTY_GOVERNANCE_KEY)
	if !exist {
		<-chan_LoyaltyGovernance_Controler
		return errors.New("loyalty governance entry not found")
	}
	loyalty_governance, ok := loyalty_governance_na.(Loyalty_Governance)
	if !ok {
		<-chan_LoyaltyGovernance_Controler
		return errors.New("loyalty governance type assertion issue")
	}
	loyalty_governance.Expired_Points_Pool = loyalty_governance.Expired_Points_Pool + points
	Map_Loyalty_Governance.Put(loyalty_governance.Key, loyalty_governance)
	<-chan_LoyaltyGovernance_Controler
	return
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
	exits := Map_Loyalty_Level.Check(request.Key)
	if exits {
		err = errors.New("key already exist")
		return Id, err
	}

	//Prepare new entry
	var NewEntry Loyalty_Level
	NewEntry.Level_Id = Map_Loyalty_AutoIncrement.GetNextAI("Loyalty_Level-Id")
	Id = NewEntry.Level_Id
	NewEntry.Key = request.Key
	NewEntry.Description = request.Description
	NewEntry.Min_Accumulated_Points = request.Min_Accumulated_Points
	NewEntry.Max_Accumulated_Points = request.Max_Accumulated_Points
	NewEntry.EnableRedeem = request.EnableRedeem
	NewEntry.DowngradeToLevel_Key = request.DowngradeToLevel_Key
	//add to cache and DB
	Map_Loyalty_Level.Put(NewEntry.Key, NewEntry)
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
	entry_na, exits := Map_Loyalty_Level.CheckThenGet(request.Key)
	if !exits {
		err = errors.New("key is not created")
		return Id, err
	}
	entry, ok := entry_na.(Loyalty_Level)
	if !ok {
		return Id, errors.New("error in type assertion")
	}
	if entry.Level_Id != request.Level_Id {
		return Id, errors.New("id is not matching")
	}
	Current_Entry := entry

	//Prepare new entry
	entry.Key = request.Key
	entry.Description = request.Description
	entry.Min_Accumulated_Points = request.Min_Accumulated_Points
	entry.Max_Accumulated_Points = request.Max_Accumulated_Points
	entry.EnableRedeem = request.EnableRedeem
	entry.DowngradeToLevel_Key = request.DowngradeToLevel_Key
	if request.NewKey != "" {
		if request.NewKey != request.Key {
			//delete old
			Map_Loyalty_Level.Delete(request.Key)
			//update key
			entry.Key = request.NewKey
		}
	}
	//add to cache and DB
	Map_Loyalty_Level.Put(entry.Key, entry)
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
		entries_na := Map_Loyalty_Level.ConvertToArray()
		if len(entries_na) > 0 {
			for _, entry_na := range entries_na {
				entry, ok := entry_na.(Loyalty_Level)
				if !ok {
					err = errors.New("error in type assertion")
					return entries, err
				} else {
					entries = append(entries, entry)
				}
			}
		}
		return entries, nil
	} else {
		entry_na, exits := Map_Loyalty_Level.CheckThenGet(Key)
		if !exits {
			err = errors.New("key does not exist")
			return entries, err
		}
		entry, ok := entry_na.(Loyalty_Level)
		if !ok {
			err = errors.New("error in type assertion")
			return entries, err
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

	var findparams daoc.DAOFindParams
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
	var paginationparams daoc.DAOPaginate
	paginationparams.Limit = Limit
	paginationparams.Page = Page
	findResult, err := DAO_Loyalty_Level.FindPaginate(findparams, paginationparams)
	if err != nil {
		return entries, err
	}
	if len(findResult) > 0 {
		for _, findres := range findResult {
			InterfaceValue := reflect.ValueOf(findres).Elem().Interface().(Loyalty_Level)
			entries = append(entries, InterfaceValue)
		}
	}
	return entries, nil

}

func (Uc *UserControl) Loyalty_Level_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry_na, exits := Map_Loyalty_Level.CheckThenGet(Key)
	if !exits {
		err = errors.New("entry does not exist")
		return err
	}
	entry, ok := entry_na.(Loyalty_Level)
	if !ok {
		err = errors.New("error in type assertion")
		return err
	}
	Map_Loyalty_Level.Delete(Key)
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
// Loyalty_Account_Segment functions
// ***********************************************************************
func (Uc *UserControl) Loyalty_Account_Segment_Add(Login string, request Loyalty_Account_Segment_AddRequest) (Id int64, err error) {
	//check key if filled and if already used
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	//check if key already used
	exits := Map_Loyalty_Account_Segment.Check(request.Key)
	if exits {
		err = errors.New("key already exist")
		return Id, err
	}

	//Prepare new entry
	var NewEntry Loyalty_Account_Segment
	NewEntry.Segment_Id = Map_Loyalty_AutoIncrement.GetNextAI("Loyalty_Account_Segment-Id")
	Id = NewEntry.Segment_Id
	NewEntry.Key = request.Key
	NewEntry.Description = request.Description
	NewEntry.Amount_From = request.Amount_From
	NewEntry.Amount_Till = request.Amount_Till
	NewEntry.AON_From = request.AON_From
	NewEntry.AON_Till = request.AON_Till

	//add to cache and DB
	Map_Loyalty_Account_Segment.Put(NewEntry.Key, NewEntry)
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
	entry_na, exits := Map_Loyalty_Account_Segment.CheckThenGet(request.Key)
	if !exits {
		err = errors.New("key is not created")
		return Id, err
	}
	entry, ok := entry_na.(Loyalty_Account_Segment)
	if !ok {
		return Id, errors.New("error in type assertion")
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
			Map_Loyalty_Account_Segment.Delete(request.Key)
			//update key
			entry.Key = request.NewKey
		}
	}
	//add to cache and DB
	Map_Loyalty_Account_Segment.Put(entry.Key, entry)
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
		entries_na := Map_Loyalty_Account_Segment.ConvertToArray()
		if len(entries_na) > 0 {
			for _, entry_na := range entries_na {
				entry, ok := entry_na.(Loyalty_Account_Segment)
				if !ok {
					err = errors.New("error in type assertion")
					return entries, err
				} else {
					entries = append(entries, entry)
				}
			}
		}
		return entries, nil
	} else {
		entry_na, exits := Map_Loyalty_Account_Segment.CheckThenGet(Key)
		if !exits {
			err = errors.New("key does not exist")
			return entries, err
		}
		entry, ok := entry_na.(Loyalty_Account_Segment)
		if !ok {
			err = errors.New("error in type assertion")
			return entries, err
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

	var findparams daoc.DAOFindParams
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
	var paginationparams daoc.DAOPaginate
	paginationparams.Limit = Limit
	paginationparams.Page = Page
	findResult, err := DAO_Loyalty_Account_Segment.FindPaginate(findparams, paginationparams)
	if err != nil {
		return entries, err
	}
	if len(findResult) > 0 {
		for _, findres := range findResult {
			InterfaceValue := reflect.ValueOf(findres).Elem().Interface().(Loyalty_Account_Segment)
			entries = append(entries, InterfaceValue)
		}
	}
	return entries, nil

}

func (Uc *UserControl) Loyalty_Account_Segment_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry_na, exits := Map_Loyalty_Account_Segment.CheckThenGet(Key)
	if !exits {
		err = errors.New("entry does not exist")
		return err
	}
	entry, ok := entry_na.(Loyalty_Account_Segment)
	if !ok {
		err = errors.New("error in type assertion")
		return err
	}
	Map_Loyalty_Account_Segment.Delete(Key)
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
	exits := Map_Loyalty_Point_Earning_Rules.Check(request.Key)
	if exits {
		err = errors.New("key already exist")
		return Id, err
	}

	//Prepare new entry
	var NewEntry Loyalty_Point_Earning_Rules
	NewEntry.Earning_Rules_Id = Map_Loyalty_AutoIncrement.GetNextAI("Loyalty_Point_Earning_Rules-Id")
	Id = NewEntry.Earning_Rules_Id
	NewEntry.Key = request.Key
	NewEntry.Description = request.Description
	NewEntry.Welcome_Points = request.Welcome_Points
	NewEntry.MobileAppDaily_Login = request.MobileAppDaily_Login
	NewEntry.MainGSMBalance_AmountConsumedPerPoint = request.MainGSMBalance_AmountConsumedPerPoint
	NewEntry.MM_P2P_Award_Type = request.MM_P2P_Award_Type
	NewEntry.MM_P2P = request.MM_P2P
	NewEntry.MM_CASHIN_Award_Type = request.MM_CASHIN_Award_Type
	NewEntry.MM_CASHIN = request.MM_CASHIN
	NewEntry.MM_CASHOUT_Award_Type = request.MM_CASHOUT_Award_Type
	NewEntry.MM_CASHOUT = request.MM_CASHOUT
	NewEntry.MM_MERCHPAY_Award_Type = request.MM_MERCHPAY_Award_Type
	NewEntry.MM_MERCHPAY = request.MM_MERCHPAY
	NewEntry.MM_BILLPAY_Award_Type = request.MM_BILLPAY_Award_Type
	NewEntry.MM_BILLPAY = request.MM_BILLPAY
	NewEntry.MM_RC_Award_Type = request.MM_RC_Award_Type
	NewEntry.MM_RC = request.MM_RC
	NewEntry.MM_CTMMOREQ_Award_Type = request.MM_CTMMOREQ_Award_Type
	NewEntry.MM_CTMMOREQ = request.MM_CTMMOREQ
	//add to cache and DB
	Map_Loyalty_Point_Earning_Rules.Put(NewEntry.Key, NewEntry)
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
	entry_na, exits := Map_Loyalty_Point_Earning_Rules.CheckThenGet(request.Key)
	if !exits {
		err = errors.New("key is not created")
		return Id, err
	}
	entry, ok := entry_na.(Loyalty_Point_Earning_Rules)
	if !ok {
		return Id, errors.New("error in type assertion")
	}
	if entry.Earning_Rules_Id != request.Earning_Rules_Id {
		return Id, errors.New("id is not matching")
	}
	Current_Entry := entry
	//Prepare new entry
	entry.Key = request.Key
	entry.Description = request.Description
	entry.Welcome_Points = request.Welcome_Points
	entry.Welcome_Points = request.Welcome_Points
	entry.MobileAppDaily_Login = request.MobileAppDaily_Login
	entry.MainGSMBalance_AmountConsumedPerPoint = request.MainGSMBalance_AmountConsumedPerPoint
	entry.MM_P2P_Award_Type = request.MM_P2P_Award_Type
	entry.MM_P2P = request.MM_P2P
	entry.MM_CASHIN_Award_Type = request.MM_CASHIN_Award_Type
	entry.MM_CASHIN = request.MM_CASHIN
	entry.MM_CASHOUT_Award_Type = request.MM_CASHOUT_Award_Type
	entry.MM_CASHOUT = request.MM_CASHOUT
	entry.MM_MERCHPAY_Award_Type = request.MM_MERCHPAY_Award_Type
	entry.MM_MERCHPAY = request.MM_MERCHPAY
	entry.MM_BILLPAY_Award_Type = request.MM_BILLPAY_Award_Type
	entry.MM_BILLPAY = request.MM_BILLPAY
	entry.MM_RC_Award_Type = request.MM_RC_Award_Type
	entry.MM_RC = request.MM_RC
	entry.MM_CTMMOREQ_Award_Type = request.MM_CTMMOREQ_Award_Type
	entry.MM_CTMMOREQ = request.MM_CTMMOREQ

	if request.NewKey != "" {
		if request.NewKey != request.Key {
			//delete old
			Map_Loyalty_Point_Earning_Rules.Delete(request.Key)
			//update key
			entry.Key = request.NewKey
		}
	}
	//add to cache and DB
	Map_Loyalty_Point_Earning_Rules.Put(entry.Key, entry)
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
		entries_na := Map_Loyalty_Point_Earning_Rules.ConvertToArray()
		if len(entries_na) > 0 {
			for _, entry_na := range entries_na {
				entry, ok := entry_na.(Loyalty_Point_Earning_Rules)
				if !ok {
					err = errors.New("error in type assertion")
					return entries, err
				} else {
					entries = append(entries, entry)
				}
			}
		}
		return entries, nil
	} else {
		entry_na, exits := Map_Loyalty_Point_Earning_Rules.CheckThenGet(Key)
		if !exits {
			err = errors.New("key does not exist")
			return entries, err
		}
		entry, ok := entry_na.(Loyalty_Point_Earning_Rules)
		if !ok {
			err = errors.New("error in type assertion")
			return entries, err
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

	var findparams daoc.DAOFindParams
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
	var paginationparams daoc.DAOPaginate
	paginationparams.Limit = Limit
	paginationparams.Page = Page
	findResult, err := DAO_Loyalty_Point_Earning_Rules.FindPaginate(findparams, paginationparams)
	if err != nil {
		return entries, err
	}
	if len(findResult) > 0 {
		for _, findres := range findResult {
			InterfaceValue := reflect.ValueOf(findres).Elem().Interface().(Loyalty_Point_Earning_Rules)
			entries = append(entries, InterfaceValue)
		}
	}
	return entries, nil

}

func (Uc *UserControl) Loyalty_Point_Earning_Rules_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry_na, exits := Map_Loyalty_Point_Earning_Rules.CheckThenGet(Key)
	if !exits {
		err = errors.New("entry does not exist")
		return err
	}
	entry, ok := entry_na.(Loyalty_Point_Earning_Rules)
	if !ok {
		err = errors.New("error in type assertion")
		return err
	}
	Map_Loyalty_Point_Earning_Rules.Delete(Key)
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
	exits := Map_Loyalty_Point_Expiry_Rules.Check(request.Key)
	if exits {
		err = errors.New("key already exist")
		return Id, err
	}

	//Prepare new entry
	var NewEntry Loyalty_Point_Expiry_Rules
	NewEntry.Expiry_Rules_Id = Map_Loyalty_AutoIncrement.GetNextAI("Loyalty_Point_Expiry_Rules-Id")
	Id = NewEntry.Expiry_Rules_Id
	NewEntry.Key = request.Key
	NewEntry.Description = request.Description
	NewEntry.Rolling_Expiration = request.Rolling_Expiration
	NewEntry.Validity_Unit = request.Validity_Unit
	NewEntry.Validity_Duration = request.Validity_Duration
	NewEntry.Fix_Date_Expiration = request.Fix_Date_Expiration
	NewEntry.Expiration_Trigger_date = request.Expiration_Trigger_date
	NewEntry.Expiration_Point_Before = request.Expiration_Point_Before
	//add to cache and DB
	Map_Loyalty_Point_Expiry_Rules.Put(NewEntry.Key, NewEntry)
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
	entry_na, exits := Map_Loyalty_Point_Expiry_Rules.CheckThenGet(request.Key)
	if !exits {
		err = errors.New("key is not created")
		return Id, err
	}
	entry, ok := entry_na.(Loyalty_Point_Expiry_Rules)
	if !ok {
		return Id, errors.New("error in type assertion")
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
	entry.Fix_Date_Expiration = request.Fix_Date_Expiration
	entry.Expiration_Trigger_date = request.Expiration_Trigger_date
	entry.Expiration_Point_Before = request.Expiration_Point_Before
	if request.NewKey != "" {
		if request.NewKey != request.Key {
			//delete old
			Map_Loyalty_Point_Expiry_Rules.Delete(request.Key)
			//update key
			entry.Key = request.NewKey
		}
	}
	//add to cache and DB
	Map_Loyalty_Point_Expiry_Rules.Put(entry.Key, entry)
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
		entries_na := Map_Loyalty_Point_Expiry_Rules.ConvertToArray()
		if len(entries_na) > 0 {
			for _, entry_na := range entries_na {
				entry, ok := entry_na.(Loyalty_Point_Expiry_Rules)
				if !ok {
					err = errors.New("error in type assertion")
					return entries, err
				} else {
					entries = append(entries, entry)
				}
			}
		}
		return entries, nil
	} else {
		entry_na, exits := Map_Loyalty_Point_Expiry_Rules.CheckThenGet(Key)
		if !exits {
			err = errors.New("key does not exist")
			return entries, err
		}
		entry, ok := entry_na.(Loyalty_Point_Expiry_Rules)
		if !ok {
			err = errors.New("error in type assertion")
			return entries, err
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

	var findparams daoc.DAOFindParams
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
	var paginationparams daoc.DAOPaginate
	paginationparams.Limit = Limit
	paginationparams.Page = Page
	findResult, err := DAO_Loyalty_Point_Expiry_Rules.FindPaginate(findparams, paginationparams)
	if err != nil {
		return entries, err
	}
	if len(findResult) > 0 {
		for _, findres := range findResult {
			InterfaceValue := reflect.ValueOf(findres).Elem().Interface().(Loyalty_Point_Expiry_Rules)
			entries = append(entries, InterfaceValue)
		}
	}
	return entries, nil

}

func (Uc *UserControl) Loyalty_Point_Expiry_Rules_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry_na, exits := Map_Loyalty_Point_Expiry_Rules.CheckThenGet(Key)
	if !exits {
		err = errors.New("entry does not exist")
		return err
	}
	entry, ok := entry_na.(Loyalty_Point_Expiry_Rules)
	if !ok {
		err = errors.New("error in type assertion")
		return err
	}
	Map_Loyalty_Point_Expiry_Rules.Delete(Key)
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
	exits := Map_Loyalty_Point_Redemption_Rules.Check(request.Key)
	if exits {
		err = errors.New("key already exist")
		return Id, err
	}

	//Prepare new entry
	var NewEntry Loyalty_Point_Redemption_Rules
	NewEntry.Redemption_Rules_Id = Map_Loyalty_AutoIncrement.GetNextAI("Loyalty_Point_Redemption_Rules-Id")
	Id = NewEntry.Redemption_Rules_Id
	NewEntry.Key = request.Key
	NewEntry.Description = request.Description
	NewEntry.Min_Accumulated_Points = request.Min_Accumulated_Points
	NewEntry.Allow_Negative_Balance_ToRedeem = request.Allow_Negative_Balance_ToRedeem
	NewEntry.Allow_PendingLendme_ToRedeem = request.Allow_PendingLendme_ToRedeem
	NewEntry.Airtime_MinPoints = request.Airtime_MinPoints
	NewEntry.Airtime_AmountPerPoint = request.Airtime_AmountPerPoint
	NewEntry.Airtime_EVC_Account = request.Airtime_EVC_Account
	NewEntry.Airtime_EVC_PIN = request.Airtime_EVC_PIN
	NewEntry.MobileMoney_MinPoints = request.MobileMoney_MinPoints
	NewEntry.MobileMoney_AmountPerPoint = request.MobileMoney_AmountPerPoint
	NewEntry.Bundles_MinPoints = request.Bundles_MinPoints
	NewEntry.Bundles_EVC_Account = request.Bundles_EVC_Account
	NewEntry.Bundles_EVC_PIN = request.Bundles_EVC_PIN
	NewEntry.Bundles_Product_Catalogue_Channel = request.Bundles_Product_Catalogue_Channel
	NewEntry.Bundles_Product_Catalogue_Plan = request.Bundles_Product_Catalogue_Plan
	NewEntry.Bundles_Product_Catalogue_Version = request.Bundles_Product_Catalogue_Version
	NewEntry.FreeSpinAndWin_MinPoints = request.FreeSpinAndWin_MinPoints
	NewEntry.FreeSpinAndWin_PointsPerSpin = request.FreeSpinAndWin_PointsPerSpin
	//add to cache and DB
	Map_Loyalty_Point_Redemption_Rules.Put(NewEntry.Key, NewEntry)
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
	entry_na, exits := Map_Loyalty_Point_Redemption_Rules.CheckThenGet(request.Key)
	if !exits {
		err = errors.New("key is not created")
		return Id, err
	}
	entry, ok := entry_na.(Loyalty_Point_Redemption_Rules)
	if !ok {
		return Id, errors.New("error in type assertion")
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
	entry.Airtime_MinPoints = request.Airtime_MinPoints
	entry.Airtime_AmountPerPoint = request.Airtime_AmountPerPoint
	entry.Airtime_EVC_Account = request.Airtime_EVC_Account
	entry.Airtime_EVC_PIN = request.Airtime_EVC_PIN
	entry.MobileMoney_MinPoints = request.MobileMoney_MinPoints
	entry.MobileMoney_AmountPerPoint = request.MobileMoney_AmountPerPoint
	entry.Bundles_MinPoints = request.Bundles_MinPoints
	entry.Bundles_EVC_Account = request.Bundles_EVC_Account
	entry.Bundles_EVC_PIN = request.Bundles_EVC_PIN
	entry.Bundles_Product_Catalogue_Channel = request.Bundles_Product_Catalogue_Channel
	entry.Bundles_Product_Catalogue_Plan = request.Bundles_Product_Catalogue_Plan
	entry.Bundles_Product_Catalogue_Version = request.Bundles_Product_Catalogue_Version
	entry.FreeSpinAndWin_MinPoints = request.FreeSpinAndWin_MinPoints
	entry.FreeSpinAndWin_PointsPerSpin = request.FreeSpinAndWin_PointsPerSpin
	if request.NewKey != "" {
		if request.NewKey != request.Key {
			//delete old
			Map_Loyalty_Point_Redemption_Rules.Delete(request.Key)
			//update key
			entry.Key = request.NewKey
		}
	}
	//add to cache and DB
	Map_Loyalty_Point_Redemption_Rules.Put(entry.Key, entry)
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
		entries_na := Map_Loyalty_Point_Redemption_Rules.ConvertToArray()
		if len(entries_na) > 0 {
			for _, entry_na := range entries_na {
				entry, ok := entry_na.(Loyalty_Point_Redemption_Rules)
				if !ok {
					err = errors.New("error in type assertion")
					return entries, err
				} else {
					entries = append(entries, entry)
				}
			}
		}
		return entries, nil
	} else {
		entry_na, exits := Map_Loyalty_Point_Redemption_Rules.CheckThenGet(Key)
		if !exits {
			err = errors.New("key does not exist")
			return entries, err
		}
		entry, ok := entry_na.(Loyalty_Point_Redemption_Rules)
		if !ok {
			err = errors.New("error in type assertion")
			return entries, err
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

	var findparams daoc.DAOFindParams
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
	var paginationparams daoc.DAOPaginate
	paginationparams.Limit = Limit
	paginationparams.Page = Page
	findResult, err := DAO_Loyalty_Point_Redemption_Rules.FindPaginate(findparams, paginationparams)
	if err != nil {
		return entries, err
	}
	if len(findResult) > 0 {
		for _, findres := range findResult {
			InterfaceValue := reflect.ValueOf(findres).Elem().Interface().(Loyalty_Point_Redemption_Rules)
			entries = append(entries, InterfaceValue)
		}
	}
	return entries, nil

}

func (Uc *UserControl) Loyalty_Point_Redemption_Rules_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry_na, exits := Map_Loyalty_Point_Redemption_Rules.CheckThenGet(Key)
	if !exits {
		err = errors.New("entry does not exist")
		return err
	}
	entry, ok := entry_na.(Loyalty_Point_Redemption_Rules)
	if !ok {
		err = errors.New("error in type assertion")
		return err
	}
	Map_Loyalty_Point_Redemption_Rules.Delete(Key)
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
	exits := Map_Loyalty_Plan.Check(request.Key)
	if exits {
		err = errors.New("key already exist")
		return Id, err
	}

	//Prepare new entry
	var NewEntry Loyalty_Plan
	NewEntry.Plan_Id = Map_Loyalty_AutoIncrement.GetNextAI("Loyalty_Plan-Id")
	Id = NewEntry.Plan_Id
	NewEntry.Key = request.Key
	NewEntry.Description = request.Description
	NewEntry.Loyalty_Level_Key = request.Loyalty_Level_Key
	NewEntry.Loyalty_Account_Segment_Key = request.Loyalty_Account_Segment_Key
	NewEntry.Earning_Rules_Key = request.Earning_Rules_Key
	NewEntry.Expiry_Rules_Key = request.Expiry_Rules_Key
	NewEntry.Redemption_Rules_Key = request.Redemption_Rules_Key
	//add to cache and DB
	Map_Loyalty_Plan.Put(NewEntry.Key, NewEntry)
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
	entry_na, exits := Map_Loyalty_Plan.CheckThenGet(request.Key)
	if !exits {
		err = errors.New("key is not created")
		return Id, err
	}
	entry, ok := entry_na.(Loyalty_Plan)
	if !ok {
		return Id, errors.New("error in type assertion")
	}
	if entry.Plan_Id != request.Plan_Id {
		return Id, errors.New("id is not matching")
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
			Map_Loyalty_Plan.Delete(request.Key)
			//update key
			entry.Key = request.NewKey
		}
	}
	//add to cache and DB
	Map_Loyalty_Plan.Put(entry.Key, entry)
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
		entries_na := Map_Loyalty_Plan.ConvertToArray()
		if len(entries_na) > 0 {
			for _, entry_na := range entries_na {
				entry, ok := entry_na.(Loyalty_Plan)
				if !ok {
					err = errors.New("error in type assertion")
					return entries, err
				} else {
					entries = append(entries, entry)
				}
			}
		}
		return entries, nil
	} else {
		entry_na, exits := Map_Loyalty_Plan.CheckThenGet(Key)
		if !exits {
			err = errors.New("key does not exist")
			return entries, err
		}
		entry, ok := entry_na.(Loyalty_Plan)
		if !ok {
			err = errors.New("error in type assertion")
			return entries, err
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

	var findparams daoc.DAOFindParams
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
	var paginationparams daoc.DAOPaginate
	paginationparams.Limit = Limit
	paginationparams.Page = Page
	findResult, err := DAO_Loyalty_Plan.FindPaginate(findparams, paginationparams)
	if err != nil {
		return entries, err
	}
	if len(findResult) > 0 {
		for _, findres := range findResult {
			InterfaceValue := reflect.ValueOf(findres).Elem().Interface().(Loyalty_Plan)
			entries = append(entries, InterfaceValue)
		}
	}
	return entries, nil

}

func (Uc *UserControl) Loyalty_Plan_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry_na, exits := Map_Loyalty_Plan.CheckThenGet(Key)
	if !exits {
		err = errors.New("entry does not exist")
		return err
	}
	entry, ok := entry_na.(Loyalty_Plan)
	if !ok {
		err = errors.New("error in type assertion")
		return err
	}
	Map_Loyalty_Plan.Delete(Key)
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
	exits := Map_Customer_UAT.Check(request.Key)
	if exits {
		err = errors.New("key already exist")
		return Id, err
	}
	//Prepare new entry
	var NewEntry Customer_UAT
	NewEntry.Id = Map_Loyalty_AutoIncrement.GetNextAI("Customer_UAT-Id")
	Id = NewEntry.Id
	NewEntry.Key = request.Key
	NewEntry.AddTime = time.Now()
	NewEntry.AddReason = request.AddReason
	//add to cache and DB
	Map_Customer_UAT.Put(NewEntry.Key, NewEntry)
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
	entry_na, exits := Map_Customer_UAT.CheckThenGet(request.Key)
	if !exits {
		err = errors.New("key is not created")
		return Id, err
	}
	entry, ok := entry_na.(Customer_UAT)
	if !ok {
		return Id, errors.New("error in type assertion")
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
			Map_Customer_UAT.Delete(request.Key)
			//update key
			entry.Key = request.NewKey
		}
	}
	//add to cache and DB
	Map_Customer_UAT.Put(entry.Key, entry)
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
		entries_na := Map_Customer_UAT.ConvertToArray()
		if len(entries_na) > 0 {
			for _, entry_na := range entries_na {
				entry, ok := entry_na.(Customer_UAT)
				if !ok {
					err = errors.New("error in type assertion")
					return entries, err
				} else {
					entries = append(entries, entry)
				}
			}
		}
		return entries, nil
	} else {
		entry_na, exits := Map_Customer_UAT.CheckThenGet(Key)
		if !exits {
			err = errors.New("key does not exist")
			return entries, err
		}
		entry, ok := entry_na.(Customer_UAT)
		if !ok {
			err = errors.New("error in type assertion")
			return entries, err
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

	var findparams daoc.DAOFindParams
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
	var paginationparams daoc.DAOPaginate
	paginationparams.Limit = Limit
	paginationparams.Page = Page
	findResult, err := DAO_Customer_UAT.FindPaginate(findparams, paginationparams)
	if err != nil {
		return entries, err
	}
	if len(findResult) > 0 {
		for _, findres := range findResult {
			InterfaceValue := reflect.ValueOf(findres).Elem().Interface().(Customer_UAT)
			entries = append(entries, InterfaceValue)
		}
	}
	return entries, nil

}

func (Uc *UserControl) Customer_UAT_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry_na, exits := Map_Customer_UAT.CheckThenGet(Key)
	if !exits {
		err = errors.New("entry does not exist")
		return err
	}
	entry, ok := entry_na.(Customer_UAT)
	if !ok {
		err = errors.New("error in type assertion")
		return err
	}
	Map_Customer_UAT.Delete(Key)
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
	exits := Map_Customer_DND.Check(request.Key)
	if exits {
		err = errors.New("key already exist")
		return Id, err
	}
	//Prepare new entry
	var NewEntry Customer_DND
	NewEntry.Id = Map_Loyalty_AutoIncrement.GetNextAI("Customer_DND-Id")
	Id = NewEntry.Id
	NewEntry.Key = request.Key
	NewEntry.AddTime = time.Now()
	NewEntry.AddReason = request.AddReason
	//add to cache and DB
	Map_Customer_DND.Put(NewEntry.Key, NewEntry)
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
	entry_na, exits := Map_Customer_DND.CheckThenGet(request.Key)
	if !exits {
		err = errors.New("key is not created")
		return Id, err
	}
	entry, ok := entry_na.(Customer_DND)
	if !ok {
		return Id, errors.New("error in type assertion")
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
			Map_Customer_DND.Delete(request.Key)
			//update key
			entry.Key = request.NewKey
		}
	}
	//add to cache and DB
	Map_Customer_DND.Put(entry.Key, entry)
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
		entries_na := Map_Customer_DND.ConvertToArray()
		if len(entries_na) > 0 {
			for _, entry_na := range entries_na {
				entry, ok := entry_na.(Customer_DND)
				if !ok {
					err = errors.New("error in type assertion")
					return entries, err
				} else {
					entries = append(entries, entry)
				}
			}
		}
		return entries, nil
	} else {
		entry_na, exits := Map_Customer_DND.CheckThenGet(Key)
		if !exits {
			err = errors.New("key does not exist")
			return entries, err
		}
		entry, ok := entry_na.(Customer_DND)
		if !ok {
			err = errors.New("error in type assertion")
			return entries, err
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

	var findparams daoc.DAOFindParams
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
	var paginationparams daoc.DAOPaginate
	paginationparams.Limit = Limit
	paginationparams.Page = Page
	findResult, err := DAO_Customer_DND.FindPaginate(findparams, paginationparams)
	if err != nil {
		return entries, err
	}
	if len(findResult) > 0 {
		for _, findres := range findResult {
			InterfaceValue := reflect.ValueOf(findres).Elem().Interface().(Customer_DND)
			entries = append(entries, InterfaceValue)
		}
	}
	return entries, nil

}

func (Uc *UserControl) Customer_DND_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry_na, exits := Map_Customer_DND.CheckThenGet(Key)
	if !exits {
		err = errors.New("entry does not exist")
		return err
	}
	entry, ok := entry_na.(Customer_DND)
	if !ok {
		err = errors.New("error in type assertion")
		return err
	}
	Map_Customer_DND.Delete(Key)
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
	exits := Map_Customer_Exclusion.Check(request.Key)
	if exits {
		err = errors.New("key already exist")
		return Id, err
	}
	//Prepare new entry
	var NewEntry Customer_Exclusion
	NewEntry.Id = Map_Loyalty_AutoIncrement.GetNextAI("Customer_Exclusion-Id")
	Id = NewEntry.Id
	NewEntry.Key = request.Key
	NewEntry.AddTime = time.Now()
	NewEntry.AddReason = request.AddReason
	//add to cache and DB
	Map_Customer_Exclusion.Put(NewEntry.Key, NewEntry)
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
	entry_na, exits := Map_Customer_Exclusion.CheckThenGet(request.Key)
	if !exits {
		err = errors.New("key is not created")
		return Id, err
	}
	entry, ok := entry_na.(Customer_Exclusion)
	if !ok {
		return Id, errors.New("error in type assertion")
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
			Map_Customer_Exclusion.Delete(request.Key)
			//update key
			entry.Key = request.NewKey
		}
	}
	//add to cache and DB
	Map_Customer_Exclusion.Put(entry.Key, entry)
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
		entries_na := Map_Customer_Exclusion.ConvertToArray()
		if len(entries_na) > 0 {
			for _, entry_na := range entries_na {
				entry, ok := entry_na.(Customer_Exclusion)
				if !ok {
					err = errors.New("error in type assertion")
					return entries, err
				} else {
					entries = append(entries, entry)
				}
			}
		}
		return entries, nil
	} else {
		entry_na, exits := Map_Customer_Exclusion.CheckThenGet(Key)
		if !exits {
			err = errors.New("key does not exist")
			return entries, err
		}
		entry, ok := entry_na.(Customer_Exclusion)
		if !ok {
			err = errors.New("error in type assertion")
			return entries, err
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

	var findparams daoc.DAOFindParams
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
	var paginationparams daoc.DAOPaginate
	paginationparams.Limit = Limit
	paginationparams.Page = Page
	findResult, err := DAO_Customer_Exclusion.FindPaginate(findparams, paginationparams)
	if err != nil {
		return entries, err
	}
	if len(findResult) > 0 {
		for _, findres := range findResult {
			InterfaceValue := reflect.ValueOf(findres).Elem().Interface().(Customer_Exclusion)
			entries = append(entries, InterfaceValue)
		}
	}
	return entries, nil

}

func (Uc *UserControl) Customer_Exclusion_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry_na, exits := Map_Customer_Exclusion.CheckThenGet(Key)
	if !exits {
		err = errors.New("entry does not exist")
		return err
	}
	entry, ok := entry_na.(Customer_Exclusion)
	if !ok {
		err = errors.New("error in type assertion")
		return err
	}
	Map_Customer_Exclusion.Delete(Key)
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
	exits := Map_Customer_COS_Exclusion.Check(request.Key)
	if exits {
		err = errors.New("key already exist")
		return Id, err
	}
	//Prepare new entry
	var NewEntry Customer_COS_Exclusion
	NewEntry.Id = Map_Loyalty_AutoIncrement.GetNextAI("Customer_COS_Exclusion-Id")
	Id = NewEntry.Id
	NewEntry.Key = request.Key
	NewEntry.AddTime = time.Now()
	NewEntry.AddReason = request.AddReason
	//add to cache and DB
	Map_Customer_COS_Exclusion.Put(NewEntry.Key, NewEntry)
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
	entry_na, exits := Map_Customer_COS_Exclusion.CheckThenGet(request.Key)
	if !exits {
		err = errors.New("key is not created")
		return Id, err
	}
	entry, ok := entry_na.(Customer_COS_Exclusion)
	if !ok {
		return Id, errors.New("error in type assertion")
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
			Map_Customer_COS_Exclusion.Delete(request.Key)
			//update key
			entry.Key = request.NewKey
		}
	}
	//add to cache and DB
	Map_Customer_COS_Exclusion.Put(entry.Key, entry)
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
		entries_na := Map_Customer_COS_Exclusion.ConvertToArray()
		if len(entries_na) > 0 {
			for _, entry_na := range entries_na {
				entry, ok := entry_na.(Customer_COS_Exclusion)
				if !ok {
					err = errors.New("error in type assertion")
					return entries, err
				} else {
					entries = append(entries, entry)
				}
			}
		}
		return entries, nil
	} else {
		entry_na, exits := Map_Customer_COS_Exclusion.CheckThenGet(Key)
		if !exits {
			err = errors.New("key does not exist")
			return entries, err
		}
		entry, ok := entry_na.(Customer_COS_Exclusion)
		if !ok {
			err = errors.New("error in type assertion")
			return entries, err
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

	var findparams daoc.DAOFindParams
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
	var paginationparams daoc.DAOPaginate
	paginationparams.Limit = Limit
	paginationparams.Page = Page
	findResult, err := DAO_Customer_COS_Exclusion.FindPaginate(findparams, paginationparams)
	if err != nil {
		return entries, err
	}
	if len(findResult) > 0 {
		for _, findres := range findResult {
			InterfaceValue := reflect.ValueOf(findres).Elem().Interface().(Customer_COS_Exclusion)
			entries = append(entries, InterfaceValue)
		}
	}
	return entries, nil

}

func (Uc *UserControl) Customer_COS_Exclusion_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry_na, exits := Map_Customer_COS_Exclusion.CheckThenGet(Key)
	if !exits {
		err = errors.New("entry does not exist")
		return err
	}
	entry, ok := entry_na.(Customer_COS_Exclusion)
	if !ok {
		err = errors.New("error in type assertion")
		return err
	}
	Map_Customer_COS_Exclusion.Delete(Key)
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
	exits := Map_Customer_Loyalty_Account.Check(request.Key)
	if exits {
		err = errors.New("key already exist")
		return Id, err
	}
	//check exclusion list
	exists_exclusion := Map_Customer_Exclusion.Check(request.Key)
	if exists_exclusion {
		err = errors.New("customer is included in the exclusion list")
		return Id, err
	}
	//check COS exclusion list
	exists_exclusion_cos := Map_Customer_COS_Exclusion.Check(request.COS)
	if exists_exclusion_cos {
		err = errors.New("customer is included in the cos exclusion list")
		return Id, err
	}
	//Prepare new entry
	var NewEntry Customer_Loyalty_Account
	NewEntry.Customer_Id = Map_Loyalty_AutoIncrement.GetNextAI("Customer_Loyalty_Account-Id")
	Id = NewEntry.Customer_Id
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

	if Login != "DWH_Import" {
		NewEntry.Loyalty_Account_Segment_Key = Loyalty_Account_Segment_Selection(request.ARPU, request.Joining_Date)
		NewEntry.Loyalty_Account_Segment_Date = time.Now()
		NewEntry.Loyalty_Account_Segment_Direction = ""
		NewEntry.Loyalty_Account_Segment_SetBy = Login
	} else {
		subscriber_na, subexist := Map_Subscribers.CheckThenGet(request.Key)
		if subexist {
			subscriber, ok := subscriber_na.(Subscriber)
			if !ok {
				NewEntry.Loyalty_Account_Segment_Key = Loyalty_Account_Segment_Selection(request.ARPU, request.Joining_Date)
				NewEntry.Loyalty_Account_Segment_Date = time.Now()
				NewEntry.Loyalty_Account_Segment_Direction = ""
				NewEntry.Loyalty_Account_Segment_SetBy = Login
			} else {
				NewEntry.Loyalty_Account_Segment_Key = Loyalty_Account_Segment_Selection(subscriber.ARPU, subscriber.FirstUse_date)
				NewEntry.Loyalty_Account_Segment_Date = time.Now()
				NewEntry.Loyalty_Account_Segment_Direction = ""
				NewEntry.Loyalty_Account_Segment_SetBy = Login
			}
		} else {
			NewEntry.Loyalty_Account_Segment_Key = Loyalty_Account_Segment_Selection(request.ARPU, request.Joining_Date)
			NewEntry.Loyalty_Account_Segment_Date = time.Now()
			NewEntry.Loyalty_Account_Segment_Direction = ""
			NewEntry.Loyalty_Account_Segment_SetBy = Login
		}
	}

	NewEntry.LoyaltyPointsDetail = make(map[string]Loyalty_Points_Detail)

	//add to cache and DB
	Map_Customer_Loyalty_Account.Put(NewEntry.Key, NewEntry)
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
	return Id, nil
}

func (Uc *UserControl) Customer_Loyalty_Account_Edit(Login string, request Customer_Loyalty_Account_EditRequest) (Id int64, err error) {
	//check and validate
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	entry_na, exits := Map_Customer_Loyalty_Account.CheckThenGet(request.Key)
	if !exits {
		err = errors.New("key is not created")
		return Id, err
	}
	entry, ok := entry_na.(Customer_Loyalty_Account)
	if !ok {
		return Id, errors.New("error in type assertion")
	}
	if entry.Customer_Id != request.Customer_Id {
		return Id, errors.New("id is not matching")
	}
	Current_Entry := entry
	//check exclusion list
	exists_exclusion := Map_Customer_Exclusion.Check(request.Key)
	if exists_exclusion {
		Uc.Customer_Loyalty_Account_Delete(Login, request.Key)
		err = errors.New("customer is included in the exclusion list")
		return Id, err
	}
	//check COS exclusion list
	exists_exclusion_cos := Map_Customer_COS_Exclusion.Check(request.COS)
	if exists_exclusion_cos {
		Uc.Customer_Loyalty_Account_Delete(Login, request.Key)
		err = errors.New("customer is included in the cos exclusion list")
		return Id, err
	}
	//Prepare new entry
	entry.Key = request.Key
	if entry.Loyalty_Level_Key != request.Loyalty_Level_Key {
		entry.Loyalty_Level_Key = request.Loyalty_Level_Key
		entry.Loyalty_Level_Date = time.Now()
		//entry.Loyalty_Level_Direction =
		entry.Loyalty_Level_SetBy = Login
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
	} else {
		subscriber_na, subexist := Map_Subscribers.CheckThenGet(request.Key)
		if subexist {
			subscriber, ok := subscriber_na.(Subscriber)
			if !ok {
				Loyalty_Account_Segment_Key := Loyalty_Account_Segment_Selection(request.ARPU, request.Joining_Date)
				if Loyalty_Account_Segment_Key != entry.Loyalty_Account_Segment_Key {
					entry.Loyalty_Account_Segment_Key = Loyalty_Account_Segment_Key
					entry.Loyalty_Account_Segment_Date = time.Now()
					entry.Loyalty_Account_Segment_Direction = ""
					entry.Loyalty_Account_Segment_SetBy = Login
				}
			} else {
				Loyalty_Account_Segment_Key := Loyalty_Account_Segment_Selection(subscriber.ARPU, subscriber.FirstUse_date)
				if Loyalty_Account_Segment_Key != entry.Loyalty_Account_Segment_Key {
					entry.Loyalty_Account_Segment_Key = Loyalty_Account_Segment_Key
					entry.Loyalty_Account_Segment_Date = time.Now()
					entry.Loyalty_Account_Segment_Direction = ""
					entry.Loyalty_Account_Segment_SetBy = Login
				}
			}
		} else {
			Loyalty_Account_Segment_Key := Loyalty_Account_Segment_Selection(request.ARPU, request.Joining_Date)
			if Loyalty_Account_Segment_Key != entry.Loyalty_Account_Segment_Key {
				entry.Loyalty_Account_Segment_Key = Loyalty_Account_Segment_Key
				entry.Loyalty_Account_Segment_Date = time.Now()
				entry.Loyalty_Account_Segment_Direction = ""
				entry.Loyalty_Account_Segment_SetBy = Login
			}
		}
	}
	if request.NewKey != "" {
		if request.NewKey != request.Key {
			//delete old
			Map_Customer_Loyalty_Account.Delete(request.Key)
			//update key
			entry.Key = request.NewKey
		}
	}

	//add to cache and DB
	Map_Customer_Loyalty_Account.Put(entry.Key, entry)
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
		entries_na := Map_Customer_Loyalty_Account.ConvertToArray()
		if len(entries_na) > 0 {
			for _, entry_na := range entries_na {
				entry, ok := entry_na.(Customer_Loyalty_Account)
				if !ok {
					err = errors.New("error in type assertion")
					return entries, err
				} else {
					entries = append(entries, entry)
				}
			}
		}
		return entries, nil
	} else {
		entry_na, exits := Map_Customer_Loyalty_Account.CheckThenGet(Key)
		if !exits {
			err = errors.New("key does not exist")
			return entries, err
		}
		entry, ok := entry_na.(Customer_Loyalty_Account)
		if !ok {
			err = errors.New("error in type assertion")
			return entries, err
		}
		entries = append(entries, entry)
		return entries, nil
	}
}

func (Uc *UserControl) Customer_Loyalty_Account_GetPaginated(Page, Limit int64) (entries []Customer_Loyalty_Account, err error) {
	if Page < 1 {
		return entries, errors.New("invalid page")
	}
	if Limit < 1 || Limit > 50000 {
		return entries, errors.New("invalid limit (accept value between 1 and 50000)")
	}

	var findparams daoc.DAOFindParams
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
	var paginationparams daoc.DAOPaginate
	paginationparams.Limit = Limit
	paginationparams.Page = Page
	findResult, err := DAO_Customer_Loyalty_Account.FindPaginate(findparams, paginationparams)
	if err != nil {
		return entries, err
	}
	if len(findResult) > 0 {
		for _, findres := range findResult {
			InterfaceValue := reflect.ValueOf(findres).Elem().Interface().(Customer_Loyalty_Account)
			entries = append(entries, InterfaceValue)
		}
	}
	return entries, nil
}

func (Uc *UserControl) Customer_Loyalty_Account_Delete(Login, Key string) (err error) {
	Key = Normalize_International_MSISDN(Key)
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	entry_na, exits := Map_Customer_Loyalty_Account.CheckThenGet(Key)
	if !exits {
		err = errors.New("entry does not exist")
		return err
	}
	entry, ok := entry_na.(Customer_Loyalty_Account)
	if !ok {
		err = errors.New("error in type assertion")
		return err
	}
	Map_Customer_Loyalty_Account.Delete(Key)
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

func (Uc *UserControl) Customer_Loyalty_Account_GetRedemption_Rules(MSISDN string) (Redemption_Rules Loyalty_Point_Redemption_Rules, err error) {
	//get loyalty account detail
	loyalty_account_na, subexist := Map_Customer_Loyalty_Account.CheckThenGet(MSISDN)
	if !subexist {
		return Redemption_Rules, errors.New("loyalty account does not exist")
	}
	loyalty_account, ok := loyalty_account_na.(Customer_Loyalty_Account)
	if !ok {
		return Redemption_Rules, errors.New("type assertion issue with Customer_Loyalty_Account")
	}
	if loyalty_account.Loyalty_Account_Segment_Key == "" {
		return Redemption_Rules, errors.New("type assertion issue with Customer_Loyalty_Account")
	}
	if loyalty_account.Loyalty_Level_Key == "" {
		return Redemption_Rules, errors.New("loyalty account level is not assigned")
	}
	//get the loyalty plan
	plan_na, planexist := Map_Loyalty_Plan.CheckThenGet(loyalty_account.Loyalty_Level_Key + "|" + loyalty_account.Loyalty_Account_Segment_Key)
	if !planexist {
		return Redemption_Rules, errors.New("loyalty plan does not exist")
	}
	plan, ok := plan_na.(Loyalty_Plan)
	if !ok {
		return Redemption_Rules, errors.New("type assertion issue with Loyalty_Plan")
	}
	//validate earning rules
	if plan.Redemption_Rules_Key == "" {
		return Redemption_Rules, errors.New("redemption rules is not defined")
	}
	redemption_Rules_na, redemptionexist := Map_Loyalty_Point_Redemption_Rules.CheckThenGet(plan.Redemption_Rules_Key)
	if !redemptionexist {
		return Redemption_Rules, errors.New("redemption rules is not defined")
	}
	Redemption_Rules, ok = redemption_Rules_na.(Loyalty_Point_Redemption_Rules)
	if !ok {
		return Redemption_Rules, errors.New("type assertion issue with Loyalty_Point_Redemption_Rules")
	}
	return Redemption_Rules, nil
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

func (Uc *UserControl) Customer_Loyalty_RedeemRequest(request_header *Request_Header, request Loyalty_Redemption_Request, response *Loyalty_Redemption_log) {
	response.ReceiveDate = time.Now()
	request.MSISDN = Normalize_International_MSISDN(request.MSISDN)
	//fill the request header info
	response.SourceIP = request_header.SourceIP
	response.SourceApp = request_header.SourceApp
	response.AppLogin = request_header.AppLogin
	response.AppVersion = request_header.AppVersion
	response.GPSLocation = request_header.GPSLocation
	response.GSMLocation = request_header.GSMLocation
	//fill the request info
	response.MSISDN = request.MSISDN
	response.ReceiveDate = time.Now()
	response.Redemption_Type = request.Redemption_Type //Airtime, Bundle, MobileMoney, SpinAndWin
	response.Redemption_Bunlde_Id = request.Redemption_Bunlde_Id
	response.Redemption_Amount = request.Redemption_Amount
	response.Points_To_Redeem = request.Points_To_Redeem

	//get loyalty account detail
	loyalty_Account_na, exits := Map_Customer_Loyalty_Account.CheckThenGet(request.MSISDN)
	if !exits {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "loyalty account not found"
		response.ErrorDescription = "loyalty account not found"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_Redemption_log(*response)
		return
	}
	loyalty_Account, ok := loyalty_Account_na.(Customer_Loyalty_Account)
	if !ok {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "error in loyalty account type assertion"
		response.ErrorDescription = "error in loyalty account type assertion"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_Redemption_log(*response)
		return
	}
	response.Customer_Id = loyalty_Account.Customer_Id
	response.Account_Status = loyalty_Account.Account_Status
	response.Loyalty_Level_Key = loyalty_Account.Loyalty_Level_Key
	response.Loyalty_Account_Segment_Key = loyalty_Account.Loyalty_Account_Segment_Key
	response.Opening_Awarded_Points = loyalty_Account.Available_Points
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
	response.Allow_Negative_Balance_ToRedeem = redemption_Rules.Allow_Negative_Balance_ToRedeem
	response.Allow_PendingLendme_ToRedeem = redemption_Rules.Allow_PendingLendme_ToRedeem
	//***To do: validate nagtive balance and pending lendme

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
		response.MinRequiredPoints = redemption_Rules.Airtime_MinPoints
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
		if redemption_Rules.Airtime_AmountPerPoint <= 0 {
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
			response.Redemption_Amount = response.Points_To_Redeem * redemption_Rules.Airtime_AmountPerPoint
		} else {
			if response.Redemption_Amount > 0 {
				//calculate Points_To_Redeem
				response.Points_To_Redeem = response.Redemption_Amount / redemption_Rules.Airtime_AmountPerPoint
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
		airtimeTransferReply, err := Uc.CGW.UC_GWClient.AirtimePurchase(Unified_charging_gateway_Client.AirtimePurchase_Request{
			PayerMSISDN:            redemption_Rules.Airtime_EVC_Account,
			PayerPIN:               redemption_Rules.Airtime_EVC_PIN,
			PaymentMethod:          "Loyalty Points",
			TargetMSISDN:           request.MSISDN,
			Amount:                 request.Redemption_Amount,
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
			return
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
		bundlePurchaseReply, err := Uc.CGW.UC_GWClient.BundlePurchase(Unified_charging_gateway_Client.BundlePurchase_Request{
			PayerMSISDN:            redemption_Rules.Bundles_EVC_Account,
			PayerPIN:               redemption_Rules.Bundles_EVC_PIN,
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
			return
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
		response.MinRequiredPoints = redemption_Rules.MobileMoney_MinPoints
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
		if redemption_Rules.MobileMoney_AmountPerPoint <= 0 {
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
			response.Redemption_Amount = response.Points_To_Redeem * redemption_Rules.MobileMoney_AmountPerPoint
		} else {
			if response.Redemption_Amount > 0 {
				//calculate Points_To_Redeem
				response.Points_To_Redeem = response.Redemption_Amount / redemption_Rules.MobileMoney_AmountPerPoint
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

	case "SpinAndWin":
		if request.Redemption_Amount < 0 {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "number of spins is not provided"
			response.ErrorDescription = "number of spins is not provided"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_Redemption_log(*response)
			return
		}
		response.MinRequiredPoints = redemption_Rules.FreeSpinAndWin_MinPoints
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
	Schemes_na := Map_Loyalty_Account_Segment.ConvertToArray()
	if len(Schemes_na) > 0 {
		for _, scheme_na := range Schemes_na {
			scheme, ok := scheme_na.(Loyalty_Account_Segment)
			if !ok {
				log.Println("error Loyalty_Account_Segment in type assertion")
				continue
			}
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
	levels_na := Map_Loyalty_Level.ConvertToArray()
	if len(levels_na) > 0 {
		for _, level_na := range levels_na {
			level, ok := level_na.(Loyalty_Level)
			if !ok {
				log.Println("error Loyalty_Level in type assertion")
				continue
			}
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

func (Uc *UserControl) Loyalty_AccountCreditPoints(request_header *Request_Header, request Loyalty_AccountCreditPoints_Request, response *Loyalty_AccountCreditPoints_log) {
	response.ReceiveDate = time.Now()
	request.MSISDN = Normalize_International_MSISDN(request.MSISDN)
	//fill the request header info
	response.SourceIP = request_header.SourceIP
	response.SourceApp = request_header.SourceApp
	response.AppLogin = request_header.AppLogin
	response.AppVersion = request_header.AppVersion
	response.GPSLocation = request_header.GPSLocation
	response.GSMLocation = request_header.GSMLocation
	//fill the request info
	response.MSISDN = request.MSISDN
	response.EventSource = request.EventSource
	response.EventType = request.EventType
	response.EventDetail = request.EventDetail
	response.EventAmount = request.EventAmount
	response.PointsToCredit = request.PointsToCredit
	response.EventDescription = request.EventDescription

	//validate loyalty account
	loyalty_account_na, subexist := Map_Customer_Loyalty_Account.CheckThenGet(request.MSISDN)
	if !subexist {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "failed to get loyalty account"
		response.ErrorDescription = "loyalty account does not exist"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_AccountCreditPoints_log(*response)
		return

	}
	loyalty_account, ok := loyalty_account_na.(Customer_Loyalty_Account)
	if !ok {
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
	response.Opening_MainGSMBalance_PendingAmount = loyalty_account.MainGSMBalance_PendingAmount
	response.Opening_MobileMoney_PendingAmount = loyalty_account.MobileMoney_PendingAmount

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
	loyalty_Level_na, exits := Map_Loyalty_Level.CheckThenGet(loyalty_account.Loyalty_Level_Key)
	if !exits {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "failed to get loyalty account level"
		response.ErrorDescription = "loyalty account level is not defined"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_AccountCreditPoints_log(*response)
		return
	}
	loyalty_Level, ok := loyalty_Level_na.(Loyalty_Level)
	if !ok {
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
	plan_na, planexist := Map_Loyalty_Plan.CheckThenGet(loyalty_account.Loyalty_Level_Key + "|" + loyalty_account.Loyalty_Account_Segment_Key)
	if !planexist {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "failed to get loyalty plan"
		response.ErrorDescription = "loyalty plan does not exist"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_AccountCreditPoints_log(*response)
		return
	}
	plan, ok := plan_na.(Loyalty_Plan)
	if !ok {
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
	point_earning_rules_na, earningexist := Map_Loyalty_Point_Earning_Rules.CheckThenGet(plan.Earning_Rules_Key)
	if !earningexist {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "failed to get earning rules"
		response.ErrorDescription = "point earning rules are not defined"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_AccountCreditPoints_log(*response)
		return
	}
	point_earning_rules, ok := point_earning_rules_na.(Loyalty_Point_Earning_Rules)
	if !ok {
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
	var points, mainGSM_pending, mobileMoney_pending float64
	if response.EventAmount > 0 {
		points, mainGSM_pending, mobileMoney_pending = Calculate_Loyalty_Points(point_earning_rules, request, loyalty_account.MainGSMBalance_PendingAmount, loyalty_account.MobileMoney_PendingAmount)
	} else {
		points = response.PointsToCredit
		mainGSM_pending = loyalty_account.MainGSMBalance_PendingAmount
		mobileMoney_pending = loyalty_account.MobileMoney_PendingAmount
	}
	if points > 0 {
		//response.OpeningAvailablePoints = (loyalty_account.Awarded_Points + loyalty_account.Expired_Points) - loyalty_account.Redeemed_Points
		response.AwardedPoints = points
		//validate governance rules
		loyalty_governance_na, lg_exist := Map_Loyalty_Governance.CheckThenGet(LOYALTY_GOVERNANCE_KEY)
		if !lg_exist {
			response.Status = "failed"
			response.StatusCode = http.StatusBadRequest
			response.StatusDescription = "failed to get governance entry"
			response.ErrorDescription = "loyalty governance entry not found"
			response.StatusDate = time.Now()
			response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
			Uc.Write_Loyalty_AccountCreditPoints_log(*response)
			return
		}
		loyalty_governance, ok := loyalty_governance_na.(Loyalty_Governance)
		if !ok {
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
		loyalty_account.Awarded_Points = loyalty_account.Awarded_Points + points
		loyalty_account.Available_Points = (loyalty_account.Awarded_Points + loyalty_account.Expired_Points) - loyalty_account.Redeemed_Points
		loyalty_account.Last_Award_Date = time.Now()
		loyalty_account.MainGSMBalance_PendingAmount = mainGSM_pending
		loyalty_account.MobileMoney_PendingAmount = mobileMoney_pending

		YYYY, MM, _, _, _, _, _ := GetTimeParts(time.Now())
		var PointsDetail Loyalty_Points_Detail
		var exist bool
		PointsDetail, exist = loyalty_account.LoyaltyPointsDetail[YYYY+MM]
		if !exist {
			PointsDetail.Year_Month = YYYY + MM
			PointsDetail.Creation_date = time.Now()
			PointsDetail.Awarded_Points = points
			PointsDetail.Available_Points = PointsDetail.Awarded_Points
		} else {
			PointsDetail.Awarded_Points = PointsDetail.Awarded_Points + points
			PointsDetail.Available_Points = PointsDetail.Awarded_Points - PointsDetail.Redeemed_Points
		}
		loyalty_account.LoyaltyPointsDetail[YYYY+MM] = PointsDetail
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
		err := Uc.Loyalty_Governance_Available_Points_Debit(points)
		if err == nil {
			Map_Customer_Loyalty_Account.Put(loyalty_account.Key, loyalty_account)
			new_Loyalty_level_key, errNL := Uc.EvaluateAndUpdate_CustomerLoyaltyLevel(response.AppLogin, loyalty_account.Key)
			if errNL != nil {
				loyalty_account.Loyalty_Level_Key = new_Loyalty_level_key
			}
		}
	} else {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "failed to credit loyalty account"
		response.ErrorDescription = "points to credit must be greater than 0"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_AccountCreditPoints_log(*response)
		return
	}
	//response.ClosureAvailablePoints = (loyalty_account.Awarded_Points + loyalty_account.Expired_Points) - loyalty_account.Redeemed_Points
	response.Closure_Loyalty_Level_Key = loyalty_account.Loyalty_Level_Key
	response.Closure_Loyalty_Account_Segment_Key = loyalty_account.Loyalty_Account_Segment_Key
	response.Closure_Awarded_Points = loyalty_account.Awarded_Points
	response.Closure_Redeemed_Points = loyalty_account.Redeemed_Points
	response.Closure_Available_Points = loyalty_account.Available_Points
	response.Closure_MainGSMBalance_PendingAmount = loyalty_account.MainGSMBalance_PendingAmount
	response.Closure_MobileMoney_PendingAmount = loyalty_account.MobileMoney_PendingAmount

	//successful reply
	response.Status = "successful"
	response.StatusCode = http.StatusOK
	response.StatusDescription = ""
	response.ErrorDescription = ""
	response.StatusDate = time.Now()
	response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
	Uc.Write_Loyalty_AccountCreditPoints_log(*response)

	AwardedTransactions.With(prometheus.Labels{"EventSource": request.EventSource, "EventType": request.EventType, "EventDetail": request.EventDetail}).Inc()
	AwardedPoints.With(prometheus.Labels{"EventSource": request.EventSource, "EventType": request.EventType, "EventDetail": request.EventDetail}).Add(points)
}

func (Uc *UserControl) Loyalty_AccountDebitPoints(request_header *Request_Header, request Loyalty_AccountDebitPoints_Request, response *Loyalty_AccountDebitPoints_log) {
	response.ReceiveDate = time.Now()
	request.MSISDN = Normalize_International_MSISDN(request.MSISDN)
	//fill the request header info
	response.SourceIP = request_header.SourceIP
	response.SourceApp = request_header.SourceApp
	response.AppLogin = request_header.AppLogin
	response.AppVersion = request_header.AppVersion
	response.GPSLocation = request_header.GPSLocation
	response.GSMLocation = request_header.GSMLocation
	//fill the request info
	response.MSISDN = request.MSISDN
	response.Debit_Amount = request.Debit_Amount
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
	loyalty_Account_na, exits := Map_Customer_Loyalty_Account.CheckThenGet(request.MSISDN)
	if !exits {
		response.Status = "failed"
		response.StatusCode = http.StatusBadRequest
		response.StatusDescription = "loyalty account not found"
		response.ErrorDescription = "loyalty account not found"
		response.StatusDate = time.Now()
		response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
		Uc.Write_Loyalty_AccountDebitPoints_log(*response)
		return
	}
	loyalty_Account, ok := loyalty_Account_na.(Customer_Loyalty_Account)
	if !ok {
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
	//debit the account
	loyalty_Account.Redeemed_Points = loyalty_Account.Redeemed_Points + request.Debit_Amount
	loyalty_Account.Available_Points = (loyalty_Account.Awarded_Points + loyalty_Account.Expired_Points) - loyalty_Account.Redeemed_Points //(Awarded_Points + Expired_Points) - Redeemed_Points
	loyalty_Account.Last_Redeem_Date = time.Now()
	//update Monthly points
	var start_date time.Time
	start_date = time.Date(2025, 04, 01, 00, 00, 59, 0, time.UTC)
	end_date := start_date.AddDate(50, 0, 0)
	Amount_to_debit := request.Debit_Amount
	for d := start_date; d.After(end_date) == false; d = d.AddDate(0, 1, 0) {
		if d.Before(time.Now().AddDate(0, 0, 1)) {
			//fmt.Println(d.Format("2006-01-02"))
			YYYY, MM, _, _, _, _, _ := GetTimeParts(d)
			var PointsDetail Loyalty_Points_Detail
			var exist bool
			PointsDetail, exist = loyalty_Account.LoyaltyPointsDetail[YYYY+MM]
			if exist {
				if PointsDetail.Available_Points > 0 {
					if PointsDetail.Available_Points >= Amount_to_debit {
						//full amount available
						PointsDetail.Redeemed_Points = PointsDetail.Redeemed_Points + Amount_to_debit
						PointsDetail.Available_Points = (PointsDetail.Awarded_Points + PointsDetail.Expired_Points) - PointsDetail.Redeemed_Points //(Awarded_Points + Expired_Points) - Redeemed_Points
						PointsDetail.Last_Redeem_Date = time.Now()
						Amount_to_debit = 0
					} else {
						//partial amount available
						partial_debit_amount := PointsDetail.Available_Points
						PointsDetail.Redeemed_Points = PointsDetail.Redeemed_Points + partial_debit_amount
						PointsDetail.Available_Points = 0
						PointsDetail.Last_Redeem_Date = time.Now()
						Amount_to_debit = Amount_to_debit - partial_debit_amount
					}
					loyalty_Account.LoyaltyPointsDetail[YYYY+MM] = PointsDetail
					if Amount_to_debit == 0 {
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
	Map_Customer_Loyalty_Account.Put(loyalty_Account.Key, loyalty_Account)
	response.Closure_Awarded_Points = loyalty_Account.Awarded_Points
	response.Closure_Redeemed_Points = loyalty_Account.Redeemed_Points
	response.Closure_Available_Points = loyalty_Account.Available_Points
	//update goveranance
	Uc.Loyalty_Governance_Redeem_Points_Debit(request.Redemption_Amount)
	//successful reply
	response.Status = "successful"
	response.StatusCode = http.StatusOK
	response.StatusDescription = ""
	response.ErrorDescription = ""
	response.StatusDate = time.Now()
	response.E2E_Elapsedtime = (time.Since(response.ReceiveDate).Nanoseconds()) / 1000000
	Uc.Write_Loyalty_AccountDebitPoints_log(*response)
}

func Calculate_Loyalty_Points(rules Loyalty_Point_Earning_Rules, award_request Loyalty_AccountCreditPoints_Request, mainGSM_CurrentPending, mobileMoney_CurrentPending float64) (points, mainGSM_pending, mobileMoney_pending float64) {
	switch award_request.EventSource {
	case "DWH_Import":
		switch award_request.EventType {
		case "NewJoining":
			return rules.Welcome_Points, mainGSM_CurrentPending, mobileMoney_CurrentPending
		default:
			return 0, mainGSM_CurrentPending, mobileMoney_CurrentPending
		}
	case "IN_feed":
		switch award_request.EventType {
		case "NewJoining":
			return rules.Welcome_Points, mainGSM_CurrentPending, mobileMoney_CurrentPending
		default:
			//award points based amount
			if award_request.EventAmount > 0 {
				if rules.MainGSMBalance_AmountConsumedPerPoint > 0 {
					flt_points := (award_request.EventAmount + mainGSM_CurrentPending) / rules.MainGSMBalance_AmountConsumedPerPoint
					int_points := int(flt_points)
					mainGSM_pending = (award_request.EventAmount + mainGSM_CurrentPending) - (float64(int_points) * rules.MainGSMBalance_AmountConsumedPerPoint)
					return float64(int_points), mainGSM_pending, mobileMoney_CurrentPending
				} else {
					return 0, mainGSM_CurrentPending, mobileMoney_CurrentPending
				}
			} else { //to do: award points based transaction
				return 0, mainGSM_CurrentPending, mobileMoney_CurrentPending
			}
		}
	case "WebPortal":
		switch award_request.EventType {
		case "NewJoining":
			return rules.Welcome_Points, mainGSM_CurrentPending, mobileMoney_CurrentPending
		default:
			//award points based amount
			if award_request.EventAmount > 0 {
				if rules.MainGSMBalance_AmountConsumedPerPoint > 0 {
					flt_points := (award_request.EventAmount + mainGSM_CurrentPending) / rules.MainGSMBalance_AmountConsumedPerPoint
					int_points := int(flt_points)
					mainGSM_pending = (award_request.EventAmount + mainGSM_CurrentPending) - (float64(int_points) * rules.MainGSMBalance_AmountConsumedPerPoint)
					return float64(int_points), mainGSM_pending, mobileMoney_CurrentPending
				} else {
					return 0, mainGSM_CurrentPending, mobileMoney_CurrentPending
				}
			} else { //to do: award points based transaction
				return 0, mainGSM_CurrentPending, mobileMoney_CurrentPending
			}
		}
	case "MobileMoney_feed":
		return 0, mainGSM_CurrentPending, mobileMoney_CurrentPending
	case "MyAfricellApp":
		switch award_request.EventType {
		case "MobileAppDaily_Login":
			return rules.MobileAppDaily_Login, mainGSM_CurrentPending, mobileMoney_CurrentPending
		default:
			return 0, mainGSM_CurrentPending, mobileMoney_CurrentPending
		}
	}
	return
}

func (Uc *UserControl) EvaluateAndUpdate_CustomerLoyaltyLevel(Login string, Account_Key string) (New_Loyalty_Level_Key string, err error) {
	loyalty_account_na, subexist := Map_Customer_Loyalty_Account.CheckThenGet(Account_Key)
	if !subexist {
		return New_Loyalty_Level_Key, errors.New("loyalty account does not exist")
	}
	loyalty_account, ok := loyalty_account_na.(Customer_Loyalty_Account)
	if !ok {
		return New_Loyalty_Level_Key, errors.New("type assertion issue with Customer_Loyalty_Account")
	}
	//evaluate loyalty level
	var New_Loyalty_Level Loyalty_Level
	loyalty_Level_na := Map_Loyalty_Level.ConvertToArray()
	if len(loyalty_Level_na) > 0 {
		for _, loyalty_Level_na := range loyalty_Level_na {
			loyalty_Level, ok := loyalty_Level_na.(Loyalty_Level)
			if !ok {
				return New_Loyalty_Level_Key, errors.New("error in type assertion")
			} else {
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
		}
		//update customer level
		current_level_na, lvlexist := Map_Loyalty_Level.CheckThenGet(loyalty_account.Loyalty_Level_Key)
		if !lvlexist {
			return New_Loyalty_Level_Key, errors.New("current level is invalid")
		}
		current_level, ok := current_level_na.(Loyalty_Level)
		if !ok {
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
			Map_Customer_Loyalty_Account.Put(loyalty_account.Key, loyalty_account)

		} else {
			//downgrade ==> Sof landing pending
			loyalty_account.Previous_Loyalty_Level_Key = loyalty_account.Loyalty_Level_Key
			loyalty_account.Previous_Loyalty_Level_Date = loyalty_account.Loyalty_Level_Date
			loyalty_account.Loyalty_Level_Key = New_Loyalty_Level.Key
			loyalty_account.Loyalty_Level_Date = time.Now()
			loyalty_account.Loyalty_Level_Direction = "Downgrade"
			loyalty_account.Loyalty_Level_SetBy = Login
			Map_Customer_Loyalty_Account.Put(loyalty_account.Key, loyalty_account)
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
						count, err := DAO_Customer_Loyalty_Account.Count(daoc.DAOCountParams{})
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
											chan_PointsExpiry_Controler <- 1
											go Uc.PointsExpiry_ProcessExec(loyalty_Account)
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
	var expiry_log Loyalty_Expiry_log
	expiry_log.ExpiryTime = time.Now()
	expiry_log.MSISDN = account.Key
	expiry_log.Opening_Awarded_Points = account.Awarded_Points
	expiry_log.Opening_Redeemed_Points = account.Redeemed_Points
	expiry_log.Opening_Available_Points = account.Available_Points
	expiry_log.Opening_Expired_Points = account.Expired_Points
	expiry_log.OpeningLoyaltyLevel = account.Loyalty_Level_Key
	expiry_log.EndLoyaltyLevel = account.Loyalty_Level_Key
	//validate the loyalty plan
	plan_na, planexist := Map_Loyalty_Plan.CheckThenGet(account.Loyalty_Level_Key + "|" + account.Loyalty_Account_Segment_Key)
	if !planexist {
		expiry_log.ExpiryStatus = "failed"
		expiry_log.ExpiryStatusDescription = "loyalty plan does not exist"
		Uc.Write_Loyalty_Expiry_log(expiry_log)
		<-chan_PointsExpiry_Controler
		return
	}
	plan, ok := plan_na.(Loyalty_Plan)
	if !ok {
		expiry_log.ExpiryStatus = "failed"
		expiry_log.ExpiryStatusDescription = "type assertion issue with Loyalty_Plan"
		Uc.Write_Loyalty_Expiry_log(expiry_log)
		<-chan_PointsExpiry_Controler
		return
	}
	expiry_log.Expiry_Rules_Key = plan.Expiry_Rules_Key
	//validate expiry rules
	if plan.Expiry_Rules_Key == "" {
		expiry_log.ExpiryStatus = "failed"
		expiry_log.ExpiryStatusDescription = "points expiry rules not defined"
		Uc.Write_Loyalty_Expiry_log(expiry_log)
		<-chan_PointsExpiry_Controler
		return
	}

	point_Expiry_Rules_Na, exist := Map_Loyalty_Point_Expiry_Rules.CheckThenGet(plan.Expiry_Rules_Key)
	if !exist {
		expiry_log.ExpiryStatus = "failed"
		expiry_log.ExpiryStatusDescription = "points expiry rules not found"
		Uc.Write_Loyalty_Expiry_log(expiry_log)
		<-chan_PointsExpiry_Controler
		return
	}
	point_Expiry_Rules, ok := point_Expiry_Rules_Na.(Loyalty_Point_Expiry_Rules)
	if !ok {
		expiry_log.ExpiryStatus = "failed"
		expiry_log.ExpiryStatusDescription = "points expiry rules type assertion issue"
		Uc.Write_Loyalty_Expiry_log(expiry_log)
		<-chan_PointsExpiry_Controler
		return
	}
	if point_Expiry_Rules.Rolling_Expiration {
		var Expiry_Date time.Time
		if point_Expiry_Rules.Validity_Unit == "Month" {
			Expiry_Date = account.Creation_date.AddDate(0, -1*point_Expiry_Rules.Validity_Duration, 0)
		} else if point_Expiry_Rules.Validity_Unit == "Year" {
			Expiry_Date = account.Creation_date.AddDate(-1*point_Expiry_Rules.Validity_Duration, 0, 0)
		} else {
			expiry_log.ExpiryStatus = "failed"
			expiry_log.ExpiryStatusDescription = "points expiry validity unit is not defined"
			Uc.Write_Loyalty_Expiry_log(expiry_log)
			<-chan_PointsExpiry_Controler
			return
		}
		YYYY, MM, _, _, _, _, _ := GetTimeParts(Expiry_Date)
		var PointsDetail Loyalty_Points_Detail
		var lpdexist bool
		PointsDetail, lpdexist = account.LoyaltyPointsDetail[YYYY+MM]
		if !lpdexist {
			<-chan_PointsExpiry_Controler
			return
		} else {
			expired_Points := PointsDetail.Available_Points
			expiry_log.Year_Month = YYYY + MM
			expiry_log.Month_Awarded_Points = PointsDetail.Awarded_Points
			expiry_log.Month_Redeemed_Points = PointsDetail.Redeemed_Points
			expiry_log.Month_Available_Points = PointsDetail.Available_Points
			expiry_log.Month_Expired_Points = PointsDetail.Available_Points
			delete(account.LoyaltyPointsDetail, YYYY+MM)
			account.Expired_Points = account.Expired_Points + PointsDetail.Available_Points
			account.Expiry_Date = time.Now()
			account.Awarded_Points = account.Awarded_Points - PointsDetail.Available_Points
			account.Available_Points = (account.Awarded_Points + account.Expired_Points) - account.Redeemed_Points //(Awarded_Points + Expired_Points) - Redeemed_Points
			Map_Customer_Loyalty_Account.Put(account.Key, account)
			//update governance expiry
			Uc.Loyalty_Governance_Expiry_Points_Credit(expired_Points)
			//update logs
			expiry_log.End_Awarded_Points = account.Awarded_Points
			expiry_log.End_Redeemed_Points = account.Redeemed_Points
			expiry_log.End_Available_Points = account.Available_Points
			expiry_log.End_Expired_Points = account.Expired_Points
			//check level downgrade
			if account.Loyalty_Level_SetBy != "DWH_Import" &&
				account.Loyalty_Level_SetBy != "INLiveFeed" &&
				account.Loyalty_Level_SetBy != "Points_Expiry" {
				new_Loyalty_level_key, errNL := Uc.EvaluateAndUpdate_CustomerLoyaltyLevel("Points_Expiry", account.Key)
				if errNL != nil {
					expiry_log.EndLoyaltyLevel = new_Loyalty_level_key
				}
			}
		}
		expiry_log.ExpiryStatus = "successful"
		expiry_log.ExpiryStatusDescription = ""
		Uc.Write_Loyalty_Expiry_log(expiry_log)
		<-chan_PointsExpiry_Controler
	} else if point_Expiry_Rules.Fix_Date_Expiration {

		<-chan_PointsExpiry_Controler
	} else {
		<-chan_PointsExpiry_Controler
	}
}

func (Uc *UserControl) LoyaltyGovernancePools_Metrics_Process() {
	exec := 0
	for range time.Tick(time.Second * 15) {
		if exec == 0 {
			exec = 1
			loyalty_governance_na, exist := Map_Loyalty_Governance.CheckThenGet(LOYALTY_GOVERNANCE_KEY)
			if !exist {
				log.Println("loyalty governance entry not found")
				exec = 0
				continue
			}
			loyalty_governance, ok := loyalty_governance_na.(Loyalty_Governance)
			if !ok {
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
		if len(MSISDN) == len(Configuration.CountryCode)+Configuration.MSISDN_Short_len {
			return MSISDN
		} else if len(MSISDN) == Configuration.MSISDN_Short_len {
			return Configuration.CountryCode + MSISDN
		} else {
			return MSISDN[len(MSISDN)-Configuration.MSISDN_Short_len:] + Configuration.CountryCode
		}
	}
}
