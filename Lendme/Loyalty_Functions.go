package Lendme

import (
	"daoc"
	"errors"
	"log"
	"reflect"
	"time"
)

var Map_Loyalty_AutoIncrement daoc.Cache_Synch
var DAO_Loyalty_AutoIncrement daoc.DAO

var Map_Loyalty_Governance daoc.Cache_Synch
var DAO_Loyalty_Governance daoc.DAO

var Map_Loyalty_Level daoc.Cache_Synch
var DAO_Loyalty_Level daoc.DAO

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

var chan_LoyaltyGovernanceAvailable_Debit_Controler = make(chan int, 1)

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
	Db := DAO_Lendme_log.DB + "_" + YYYY + MM
	Col := DAO_Lendme_log.Collection + "_" + DD
	_, err := DAO_Lendme_log.PutOneLogs(record, Db, Col)
	if err != nil {
		log.Println("Error in Write_Lendme_log:", err, " (", record, ")")
		return
	}
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
	chan_LoyaltyGovernanceAvailable_Debit_Controler <- 1
	loyalty_governance_na, exist := Map_Loyalty_Governance.CheckThenGet(LOYALTY_GOVERNANCE_KEY)
	if !exist {
		<-chan_LoyaltyGovernanceAvailable_Debit_Controler
		return errors.New("loyalty governance entry not found")
	}
	loyalty_governance, ok := loyalty_governance_na.(Loyalty_Governance)
	if !ok {
		<-chan_LoyaltyGovernanceAvailable_Debit_Controler
		return errors.New("loyalty governance type assertion issue")
	}
	if (loyalty_governance.Available_Points_Pool - loyalty_governance.Distributed_Points_Pool) > points {
		loyalty_governance.Distributed_Points_Pool = loyalty_governance.Distributed_Points_Pool - points
		Map_Loyalty_Governance.Put(loyalty_governance.Key, loyalty_governance)
		<-chan_LoyaltyGovernanceAvailable_Debit_Controler
		return
	} else {
		<-chan_LoyaltyGovernanceAvailable_Debit_Controler
		return errors.New("no enough loyalty points to distribute from the governance available balance")
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
	NewEntry.MobileMoney_AmountConsumedPerPoint = request.MobileMoney_AmountConsumedPerPoint

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
	entry.MobileMoney_AmountConsumedPerPoint = request.MobileMoney_AmountConsumedPerPoint

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
	NewEntry.Validity_Unit = request.Validity_Unit
	NewEntry.Validity_Duration = request.Validity_Duration

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
	entry.Validity_Unit = request.Validity_Unit
	entry.Validity_Duration = request.Validity_Duration
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
	NewEntry.Product_Catalogue = request.Product_Catalogue
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
	entry.Product_Catalogue = request.Product_Catalogue
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
	//Prepare new entry
	var NewEntry Customer_Loyalty_Account
	NewEntry.Customer_Id = Map_Loyalty_AutoIncrement.GetNextAI("Customer_Loyalty_Account-Id")
	Id = NewEntry.Customer_Id
	NewEntry.Key = request.Key

	NewEntry.Creation_date = time.Now()
	NewEntry.Account_Status = "active"
	NewEntry.Account_Profile = ""

	NewEntry.Loyalty_Level_Key = Loyalty_Level_Selection(0)
	NewEntry.Loyalty_Level_Date = time.Now()
	NewEntry.Loyalty_Level_Direction = ""
	NewEntry.Loyalty_Level_SetBy = "program" //program or admin, if admin program cannot change anymore

	subscriber_na, subexist := Map_Subscribers.CheckThenGet(request.Key)
	if subexist {
		subscriber, ok := subscriber_na.(Subscriber)
		if !ok {
			NewEntry.Loyalty_Account_Segment_Key = Loyalty_Account_Segment_Selection(0, time.Now())
			NewEntry.Loyalty_Account_Segment_Date = time.Now()
			NewEntry.Loyalty_Account_Segment_Direction = ""
			NewEntry.Loyalty_Account_Segment_SetBy = "program"
		} else {
			NewEntry.Loyalty_Account_Segment_Key = Loyalty_Account_Segment_Selection(subscriber.ARPU, subscriber.FirstUse_date)
			NewEntry.Loyalty_Account_Segment_Date = time.Now()
			NewEntry.Loyalty_Account_Segment_Direction = ""
			NewEntry.Loyalty_Account_Segment_SetBy = "program"
		}
	} else {
		NewEntry.Loyalty_Account_Segment_Key = Loyalty_Account_Segment_Selection(0, time.Now())
		NewEntry.Loyalty_Account_Segment_Date = time.Now()
		NewEntry.Loyalty_Account_Segment_Direction = ""
		NewEntry.Loyalty_Account_Segment_SetBy = "program"
	}
	NewEntry.LoyaltyPointsDetail = make(map[string]Loyalty_Points_Detail)

	//add to cache and DB
	Map_Customer_Loyalty_Account.Put(NewEntry.Key, NewEntry)
	//add logs --> off to avoid filling up logs
	// Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
	// 	Event_User:         Login,
	// 	Event_Time:         time.Now(),
	// 	Event_AffectedType: "Customer_Loyalty_Account",
	// 	Event_ActionType:   "Add",
	// 	Event_Description:  "",
	// 	Event_Entry_Before: nil,
	// 	Event_Entry_After:  NewEntry,
	// })
	Uc.Customer_Loyalty_Account_AwardPoints(Login, Customer_Loyalty_Account_AwardRequest{
		MSISDN:      NewEntry.Key,
		EventSource: request.EventSource,
		EventType:   "NewJoining",
		Amount:      0,
	})
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
	//Prepare new entry
	entry.Key = request.Key
	if entry.Loyalty_Level_Key != request.Loyalty_Level_Key {
		entry.Loyalty_Level_Key = request.Loyalty_Level_Key
		entry.Loyalty_Level_Date = time.Now()
		//entry.Loyalty_Level_Direction =
		entry.Loyalty_Level_SetBy = Login
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
	Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_AffectedType: "Customer_Loyalty_Account",
		Event_ActionType:   "Edit",
		Event_Description:  "",
		Event_Entry_Before: Current_Entry,
		Event_Entry_After:  entry,
	})
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

func (Uc *UserControl) Customer_Loyalty_Account_AwardPoints(Login string, request Customer_Loyalty_Account_AwardRequest) (err error) {
	//validate loyalty account
	loyalty_account_na, subexist := Map_Customer_Loyalty_Account.CheckThenGet(request.MSISDN)
	if !subexist {
		return errors.New("loyalty account does not exist")
	}
	loyalty_account, ok := loyalty_account_na.(Customer_Loyalty_Account)
	if !ok {
		return errors.New("type assertion issue with Customer_Loyalty_Account")
	}
	if loyalty_account.Loyalty_Account_Segment_Key == "" {
		return errors.New("loyalty account segment is not assigned")
	}
	if loyalty_account.Loyalty_Level_Key == "" {
		return errors.New("loyalty account level is not assigned")
	}
	//validate the loyalty plan
	plan_na, planexist := Map_Loyalty_Plan.CheckThenGet(loyalty_account.Loyalty_Level_Key + "|" + loyalty_account.Loyalty_Account_Segment_Key)
	if !planexist {
		return errors.New("loyalty plan does not exist")
	}
	plan, ok := plan_na.(Loyalty_Plan)
	if !ok {
		return errors.New("type assertion issue with Loyalty_Plan")
	}
	//validate earning rules
	if plan.Earning_Rules_Key == "" {
		return errors.New("earning rules is not defined")
	}
	point_earning_rules_na, earningexist := Map_Loyalty_Point_Earning_Rules.CheckThenGet(plan.Earning_Rules_Key)
	if !earningexist {
		return errors.New("point earning rules is not defined")
	}
	point_earning_rules, ok := point_earning_rules_na.(Loyalty_Point_Earning_Rules)
	if !ok {
		return errors.New("type assertion issue with Loyalty_Point_Earning_Rules")
	}
	//calucate points
	points, mainGSM_pending, mobileMoney_pending := Calculate_Loyalty_Points(point_earning_rules, request, loyalty_account.MainGSMBalance_PendingAmount, loyalty_account.MobileMoney_PendingAmount)

	if points > 0 {
		//validate governance rules
		loyalty_governance_na, lg_exist := Map_Loyalty_Governance.CheckThenGet(LOYALTY_GOVERNANCE_KEY)
		if !lg_exist {
			return errors.New("loyalty governance entry not found")
		}
		loyalty_governance, ok := loyalty_governance_na.(Loyalty_Governance)
		if !ok {
			return errors.New("loyalty governance type assertion issue")
		}
		if points > loyalty_governance.MaxAllowedPoints_PerTransaction {
			return errors.New("awarded points per transaction is exceeding loyalty governance rules")
		}
		if (loyalty_account.Awarded_Points + points) > loyalty_governance.MaxSubsAwardedPoints {
			return errors.New("awarded points per subscriber is exceeding loyalty governance rules")
		}
		//credit loyalty account
		loyalty_account.Awarded_Points = loyalty_account.Awarded_Points + points
		loyalty_account.Available_Points = loyalty_account.Awarded_Points - loyalty_account.Redeemed_Points
		loyalty_account.Last_Award_Date = time.Now()
		loyalty_account.MainGSMBalance_PendingAmount = mainGSM_pending
		loyalty_account.MobileMoney_PendingAmount = mobileMoney_pending
		YYYY, MM, _, _, _, _, _ := GetTimeParts(time.Now())
		var PointsDetail Loyalty_Points_Detail
		var exist bool
		PointsDetail = loyalty_account.LoyaltyPointsDetail[YYYY+MM]
		if !exist {
			PointsDetail.Year_Month = YYYY + MM
			PointsDetail.Awarded_Points = points
			PointsDetail.Available_Points = PointsDetail.Awarded_Points
		} else {
			PointsDetail.Awarded_Points = PointsDetail.Awarded_Points + points
			PointsDetail.Available_Points = PointsDetail.Awarded_Points - PointsDetail.Redeemed_Points
		}
		loyalty_account.LoyaltyPointsDetail[YYYY+MM] = PointsDetail
		if PointsDetail.Awarded_Points > loyalty_governance.MaxSubsAwardedPoints_PerMonth {
			return errors.New("monthly awarded points per subscriber is exceeding loyalty governance rules")
		}
		err := Uc.Loyalty_Governance_Available_Points_Debit(points)
		if err == nil {
			Map_Customer_Loyalty_Account.Put(loyalty_account.Key, loyalty_account)
		}
	}
	return
}

func Calculate_Loyalty_Points(rules Loyalty_Point_Earning_Rules, award_request Customer_Loyalty_Account_AwardRequest, mainGSM_CurrentPending, mobileMoney_CurrentPending float64) (points, mainGSM_pending, mobileMoney_pending float64) {
	switch award_request.EventSource {
	case "IN_feed":
		switch award_request.EventType {
		case "NewJoining":
			return rules.Welcome_Points, mainGSM_CurrentPending, mobileMoney_CurrentPending
		default:
			//award points based amount
			if award_request.Amount > 0 {
				if rules.MainGSMBalance_AmountConsumedPerPoint > 0 {
					flt_points := (award_request.Amount + mainGSM_CurrentPending) / rules.MainGSMBalance_AmountConsumedPerPoint
					int_points := int(flt_points)
					mainGSM_pending = (award_request.Amount + mainGSM_CurrentPending) - (float64(int_points) * rules.MainGSMBalance_AmountConsumedPerPoint)
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
