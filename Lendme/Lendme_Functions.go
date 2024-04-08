package Lendme

import (
	"afr_auth_center/AuthCenter"
	"daoc"
	"errors"
	"log"
	"reflect"
	"strconv"
	"time"
)

var Map_DefaultValues daoc.Cache_Synch
var DAO_DefaultValues daoc.DAO

var MapAccessEntry daoc.Cache_Synch
var DAO_AccessEntry daoc.DAO

var Map_AutoIncrement daoc.Cache_Synch
var DAO_AutoIncrement daoc.DAO

var Map_Subscribers daoc.Cache_Synch
var DAO_Subscribers daoc.DAO

var Map_Credit_Limit_Scheme daoc.Cache_Synch
var DAO_Credit_Limit_Scheme daoc.DAO

var DAO_Lendme_log daoc.DAO

func (uc *UserControl) InitializeCache() {
	var access_entry AuthCenter.AccessEntry
	MapAccessEntry.Initialize("AccessEntry", "AccessEntry", reflect.TypeOf(AuthCenter.AccessEntry{}), access_entry, true, &DAO_AccessEntry, uc.CacheDir.List)
	var AutoIncr daoc.AutoIncrement
	Map_AutoIncrement.Initialize("AutoIncrement", "AutoIncrement", reflect.TypeOf(daoc.AutoIncrement{}), AutoIncr, true, &DAO_AutoIncrement, uc.CacheDir.List)
	var defaultValues DefaultValues
	Map_DefaultValues.Initialize("DefaultValues", "DefaultValues", reflect.TypeOf(DefaultValues{}), defaultValues, true, &DAO_DefaultValues, uc.CacheDir.List)
	var subscriber Subscriber
	Map_Subscribers.Initialize("Subscriber", "Subscriber", reflect.TypeOf(Subscriber{}), subscriber, true, &DAO_Subscribers, uc.CacheDir.List)
	var credit_Limit_Scheme Credit_Limit_Scheme
	Map_Credit_Limit_Scheme.Initialize("Credit_Limit_Scheme", "Credit_Limit_Scheme", reflect.TypeOf(Credit_Limit_Scheme{}), credit_Limit_Scheme, true, &DAO_Credit_Limit_Scheme, uc.CacheDir.List)
}

func (uc *UserControl) InitializeDAO() {
	DAO_AccessEntry.Initialize("AccessEntry", uc.MongoDB.MongoDBClient, reflect.TypeOf(AuthCenter.AccessEntry{}), Configuration.DB_Name, "Col_AccessEntry", "")
	DAO_AutoIncrement.Initialize("AutoIncrement", uc.MongoDB.MongoDBClient, reflect.TypeOf(daoc.AutoIncrement{}), Configuration.DB_Name, "Col_AutoIncrement", "")
	DAO_DefaultValues.Initialize("DefaultValues", uc.MongoDB.MongoDBClient, reflect.TypeOf(DefaultValues{}), Configuration.DB_Name, "Col_DefaultValues", "")
	DAO_Subscribers.Initialize("Subscriber", uc.MongoDB.MongoDBClient, reflect.TypeOf(Subscriber{}), Configuration.DB_Name, "Col_Subscriber", "")
	DAO_Credit_Limit_Scheme.Initialize("Credit_Limit_Scheme", uc.MongoDB.MongoDBClient, reflect.TypeOf(Credit_Limit_Scheme{}), Configuration.DB_Name, "Col_Credit_Limit_Scheme", "")
	DAO_Lendme_log.Initialize("Lendme_log", uc.MongoDB.MongoDBClient, reflect.TypeOf(Lendme_log{}), Configuration.DB_Name, "Col_Lendme_log", "")
}

func (uc *UserControl) IndexesMaintenanceProcess() {
	log.Println("DB index manintenance process started...")
	exists, err := DAO_AutoIncrement.CheckAndCreateIndex("Idx_AutoIncrement_Key", []string{"Key"}, true)
	if err != nil {
		log.Println("Error checking and creating index Idx_AutoIncrement_Key: ", err)
	} else {
		if !exists {
			log.Println("Index Idx_AutoIncrement_Key created")
		}
	}

	exists, err = DAO_DefaultValues.CheckAndCreateIndex("Idx_DefaultValues_Key", []string{"Key"}, true)
	if err != nil {
		log.Println("Error creating index Idx_DefaultValues_Key: ", err)
	} else {
		if !exists {
			log.Println("Index Idx_DefaultValues_Key created")
		}
	}

	exists, err = DAO_Subscribers.CheckAndCreateIndex("Idx_Subscriber_Key", []string{"Key"}, true)
	if err != nil {
		log.Println("Error creating index Idx_Subscriber_Key: ", err)
	} else {
		if !exists {
			log.Println("Index Idx_Subscriber_Key created")
		}
	}

	exists, err = DAO_Credit_Limit_Scheme.CheckAndCreateIndex("Idx_Credit_Limit_Scheme_Key", []string{"Key"}, true)
	if err != nil {
		log.Println("Error creating index Idx_Credit_Limit_Scheme_Key: ", err)
	} else {
		if !exists {
			log.Println("Index Idx_Credit_Limit_Scheme_Key created")
		}
	}

	log.Println("DB index manintenance process finished")
}

func (uc *UserControl) LoadDefaultValues() {
	exist := Map_DefaultValues.Check("Default")
	if !exist {
		defaultValues := DefaultValues{
			Key: "Default",
		}
		Map_DefaultValues.Put(defaultValues.Key, defaultValues)
	}
	uc.Credit_Limit_Scheme_LoadDefaultValues()
}

func (Uc *UserControl) Write_Lendme_log(record Lendme_log) {
	YYYY, MM, _, DD, _, _, _ := GetTimeParts(record.Log_Date)
	Db := "Lendme_log_Archive_" + YYYY + MM
	Col := "Col_Lendme_log_" + DD
	_, err := DAO_Lendme_log.PutOneLogs(record, Db, Col)
	if err != nil {
		log.Println("Error in Write_Lendme_log:", err, " (", record, ")")
		return
	}
}

// ///////////////////////////////////////////////////////
// Outlet Sales
// ///////////////////////////////////////////////////////
func (Uc *UserControl) Credit_Limit_Scheme_Add(Login string, request Credit_Limit_Scheme_Add_Request) (Id int64, err error) {
	//check key if filled and if already used
	if request.Key == "" {
		err = errors.New("scheme name cannot be empty")
		return Id, err
	}
	//check if key already used
	exits := Map_Credit_Limit_Scheme.Check(request.Key)
	if exits {
		err = errors.New("scheme already exist")
		return Id, err
	}

	//Prepare new entry
	var NewEntry Credit_Limit_Scheme
	NewEntry.Scheme_Id = Map_AutoIncrement.GetNextAI("Credit_Limit_Scheme-Id")
	Id = NewEntry.Scheme_Id
	NewEntry.Key = request.Key
	NewEntry.Amount_From = request.Amount_From
	NewEntry.Amount_Till = request.Amount_Till
	NewEntry.AON_From = request.AON_From
	NewEntry.AON_Till = request.AON_Till
	NewEntry.Credit_limit_Amount = request.Credit_limit_Amount
	log := Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_Type:         "add",
		Event_Description:  "create new entry",
		Event_Entry_Before: nil,
		Event_Entry_After:  nil,
	}
	NewEntry.Event_Logs = append(NewEntry.Event_Logs, log)
	//add to cache and DB
	Map_Credit_Limit_Scheme.Put(NewEntry.Key, NewEntry)
	return Id, nil
}

func (Uc *UserControl) Credit_Limit_Scheme_Edit(Login string, request Credit_Limit_Scheme_Edit_Request) (Id int64, err error) {
	//check and validate outlet
	if request.Key == "" {
		err = errors.New("scheme name cannot be empty")
		return Id, err
	}
	scheme_na, exits := Map_Credit_Limit_Scheme.CheckThenGet(request.Key)
	if !exits {
		err = errors.New("Scheme is not created")
		return Id, err
	}
	scheme, ok := scheme_na.(Credit_Limit_Scheme)
	if !ok {
		return Id, errors.New("error in credit limit scheme type assertion")
	}
	if scheme.Scheme_Id != request.Scheme_Id {
		return Id, errors.New("id is not matching")
	}
	//take log for current values without the log trail
	current_scheme := scheme
	current_scheme.Event_Logs = nil

	//Prepare new entry
	scheme.Key = request.Key
	scheme.Amount_From = request.Amount_From
	scheme.Amount_Till = request.Amount_Till
	scheme.AON_From = request.AON_From
	scheme.AON_Till = request.AON_Till
	scheme.Credit_limit_Amount = request.Credit_limit_Amount
	log := Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_Type:         "edit",
		Event_Description:  "edit entry",
		Event_Entry_Before: current_scheme,
		Event_Entry_After:  nil,
	}
	scheme.Event_Logs = append(scheme.Event_Logs, log)
	if request.NewKey != "" {
		if request.NewKey != request.Key {
			//delete old
			Map_Credit_Limit_Scheme.Delete(request.Key)
			//update key
			scheme.Key = request.NewKey
		}
	}
	//add to cache and DB
	Map_Credit_Limit_Scheme.Put(scheme.Key, scheme)
	return Id, nil
}

func (Uc *UserControl) Credit_Limit_Scheme_Get(Key string) (schemes []Credit_Limit_Scheme, err error) {
	if Key == "" {
		schemes_na := Map_Credit_Limit_Scheme.ConvertToArray()
		if len(schemes_na) > 0 {
			for _, entry := range schemes_na {
				scheme, ok := entry.(Credit_Limit_Scheme)
				if !ok {
					err = errors.New("error in credit limit scheme type assertion")
					return schemes, err
				} else {
					schemes = append(schemes, scheme)
				}
			}
		}
		return schemes, nil
	} else {
		scheme_na, exits := Map_Credit_Limit_Scheme.CheckThenGet(Key)
		if !exits {
			err = errors.New("credit limit scheme does not exist")
			return schemes, err
		}
		scheme, ok := scheme_na.(Credit_Limit_Scheme)
		if !ok {
			err = errors.New("error in credit limit scheme type assertion")
			return schemes, err
		}
		schemes = append(schemes, scheme)
		return schemes, nil
	}
}

func (Uc *UserControl) Credit_Limit_Scheme_GetPaginated(Key string, Page, Limit int64) (schemes []Credit_Limit_Scheme, err error) {
	if Page < 1 {
		return schemes, errors.New("invalid page")
	}
	if Limit < 1 || Limit > 50000 {
		return schemes, errors.New("invalid limit (accept value between 1 and 50000)")
	}
	if Key == "" {
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
		findResult, err := DAO_Credit_Limit_Scheme.FindPaginate(findparams, paginationparams)
		if err != nil {
			return schemes, err
		}
		if len(findResult) > 0 {
			for _, findres := range findResult {
				InterfaceValue := reflect.ValueOf(findres).Elem().Interface().(Credit_Limit_Scheme)
				schemes = append(schemes, InterfaceValue)
			}
		}
		return schemes, nil
	} else {
		scheme_na, exits := Map_Credit_Limit_Scheme.CheckThenGet(Key)
		if !exits {
			err = errors.New("credit limit scheme does not exist")
			return schemes, err
		}
		scheme, ok := scheme_na.(Credit_Limit_Scheme)
		if !ok {
			err = errors.New("error in credit limit scheme type assertion")
			return schemes, err
		}
		schemes = append(schemes, scheme)
		return schemes, nil
	}
}

func (Uc *UserControl) Credit_Limit_Scheme_Delete(Key string) (err error) {
	if Key == "" {
		err = errors.New("scheme name cannot be empty")
		return err
	}
	exits := Map_Credit_Limit_Scheme.Check(Key)
	if !exits {
		err = errors.New("scheme does not exist")
		return err
	}
	Map_Credit_Limit_Scheme.Delete(Key)
	return nil
}

func (uc *UserControl) Credit_Limit_Scheme_LoadDefaultValues() {
	log.Println("Loading credit limit scheme default values")
	//
	request := Credit_Limit_Scheme_Add_Request{
		Key:                 "50_100_3_6",
		Scheme_Id:           0,
		Amount_From:         50,
		Amount_Till:         100,
		AON_From:            3,
		AON_Till:            6,
		Credit_limit_Amount: 10,
	}
	_, err := uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "50_100_6_12",
		Scheme_Id:           0,
		Amount_From:         50,
		Amount_Till:         100,
		AON_From:            6,
		AON_Till:            12,
		Credit_limit_Amount: 10,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "50_100_12_24",
		Scheme_Id:           0,
		Amount_From:         50,
		Amount_Till:         100,
		AON_From:            12,
		AON_Till:            24,
		Credit_limit_Amount: 20,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "50_100_24_1200",
		Scheme_Id:           0,
		Amount_From:         50,
		Amount_Till:         100,
		AON_From:            24,
		AON_Till:            1200,
		Credit_limit_Amount: 20,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	//
	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "100_200_3_6",
		Scheme_Id:           0,
		Amount_From:         100,
		Amount_Till:         200,
		AON_From:            3,
		AON_Till:            6,
		Credit_limit_Amount: 20,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "100_200_6_12",
		Scheme_Id:           0,
		Amount_From:         100,
		Amount_Till:         200,
		AON_From:            6,
		AON_Till:            12,
		Credit_limit_Amount: 20,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "100_200_12_24",
		Scheme_Id:           0,
		Amount_From:         100,
		Amount_Till:         200,
		AON_From:            12,
		AON_Till:            24,
		Credit_limit_Amount: 30,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "100_200_24_1200",
		Scheme_Id:           0,
		Amount_From:         100,
		Amount_Till:         200,
		AON_From:            24,
		AON_Till:            1200,
		Credit_limit_Amount: 40,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	//
	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "200_400_3_6",
		Scheme_Id:           0,
		Amount_From:         200,
		Amount_Till:         400,
		AON_From:            3,
		AON_Till:            6,
		Credit_limit_Amount: 30,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "200_400_6_12",
		Scheme_Id:           0,
		Amount_From:         200,
		Amount_Till:         400,
		AON_From:            6,
		AON_Till:            12,
		Credit_limit_Amount: 30,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "200_400_12_24",
		Scheme_Id:           0,
		Amount_From:         200,
		Amount_Till:         400,
		AON_From:            12,
		AON_Till:            24,
		Credit_limit_Amount: 50,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "200_400_24_1200",
		Scheme_Id:           0,
		Amount_From:         200,
		Amount_Till:         400,
		AON_From:            24,
		AON_Till:            1200,
		Credit_limit_Amount: 60,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}
	//
	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "400_800_3_6",
		Scheme_Id:           0,
		Amount_From:         400,
		Amount_Till:         800,
		AON_From:            3,
		AON_Till:            6,
		Credit_limit_Amount: 40,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "400_800_6_12",
		Scheme_Id:           0,
		Amount_From:         400,
		Amount_Till:         800,
		AON_From:            6,
		AON_Till:            12,
		Credit_limit_Amount: 40,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "400_800_12_24",
		Scheme_Id:           0,
		Amount_From:         400,
		Amount_Till:         800,
		AON_From:            12,
		AON_Till:            24,
		Credit_limit_Amount: 60,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "400_800_24_1200",
		Scheme_Id:           0,
		Amount_From:         400,
		Amount_Till:         800,
		AON_From:            24,
		AON_Till:            1200,
		Credit_limit_Amount: 100,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	//
	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "800_1600_3_6",
		Scheme_Id:           0,
		Amount_From:         800,
		Amount_Till:         1600,
		AON_From:            3,
		AON_Till:            6,
		Credit_limit_Amount: 80,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "800_1600_6_12",
		Scheme_Id:           0,
		Amount_From:         800,
		Amount_Till:         1600,
		AON_From:            6,
		AON_Till:            12,
		Credit_limit_Amount: 100,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "800_1600_12_24",
		Scheme_Id:           0,
		Amount_From:         800,
		Amount_Till:         1600,
		AON_From:            12,
		AON_Till:            24,
		Credit_limit_Amount: 150,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "800_1600_24_1200",
		Scheme_Id:           0,
		Amount_From:         800,
		Amount_Till:         1600,
		AON_From:            24,
		AON_Till:            1200,
		Credit_limit_Amount: 200,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	//
	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "1600_3000_3_6",
		Scheme_Id:           0,
		Amount_From:         1600,
		Amount_Till:         3000,
		AON_From:            3,
		AON_Till:            6,
		Credit_limit_Amount: 150,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "1600_3000_6_12",
		Scheme_Id:           0,
		Amount_From:         1600,
		Amount_Till:         3000,
		AON_From:            6,
		AON_Till:            12,
		Credit_limit_Amount: 200,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "1600_3000_12_24",
		Scheme_Id:           0,
		Amount_From:         1600,
		Amount_Till:         3000,
		AON_From:            12,
		AON_Till:            24,
		Credit_limit_Amount: 250,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "1600_3000_24_1200",
		Scheme_Id:           0,
		Amount_From:         1600,
		Amount_Till:         3000,
		AON_From:            24,
		AON_Till:            1200,
		Credit_limit_Amount: 300,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	//
	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "3000_9999999_3_6",
		Scheme_Id:           0,
		Amount_From:         3000,
		Amount_Till:         9999999,
		AON_From:            3,
		AON_Till:            6,
		Credit_limit_Amount: 200,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "3000_9999999_6_12",
		Scheme_Id:           0,
		Amount_From:         3000,
		Amount_Till:         9999999,
		AON_From:            6,
		AON_Till:            12,
		Credit_limit_Amount: 300,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "3000_9999999_12_24",
		Scheme_Id:           0,
		Amount_From:         3000,
		Amount_Till:         9999999,
		AON_From:            12,
		AON_Till:            24,
		Credit_limit_Amount: 400,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "3000_9999999_24_1200",
		Scheme_Id:           0,
		Amount_From:         3000,
		Amount_Till:         9999999,
		AON_From:            24,
		AON_Till:            1200,
		Credit_limit_Amount: 500,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}
}

func (Uc *UserControl) Credit_Limit_Scheme_Selection(Amount float64, FirstUse_date time.Time, LastRecharge_date time.Time) (scheme_name string, NotElligibleReason string) {
	//AON >= 3 months
	//Avg 3M Recharge >= 50
	//last recharge date within past 60 days
	//not Optin to SOS
	//balance >=0
	var Min_AON float64
	var Min_Amount float64
	AON_Hours := time.Now().Sub(FirstUse_date).Hours()
	AON_Months := (AON_Hours / 24) / 30
	if FirstUse_date.IsZero() {
		NotElligibleReason = "Min AON"
		return
	}
	if AON_Months < Configuration.Min_Allowed_AON {
		NotElligibleReason = "Min AON"
		return
	}
	if Amount < Configuration.Min_Avg3MRecharge {
		NotElligibleReason = "Min Avg3MRecharge"
		return
	}

	LastRecharge_Hours := time.Now().Sub(LastRecharge_date).Hours()
	LastRecharge_Months := (LastRecharge_Hours / 24) / 30
	if LastRecharge_date.IsZero() {
		NotElligibleReason = "Min Last Recharge Period"
		return
	}
	if LastRecharge_Months > Configuration.Min_LastRechargePeriod {
		NotElligibleReason = "Min Last Recharge Period"
		return
	}

	Schemes_na := Map_Credit_Limit_Scheme.ConvertToArray()
	if len(Schemes_na) > 0 {
		for _, scheme_na := range Schemes_na {
			scheme, ok := scheme_na.(Credit_Limit_Scheme)
			if !ok {
				log.Println("error in Credit_Limit_Scheme_Selection Credit_Limit_Scheme type assertion")
				continue
			}
			if Min_Amount == 0 {
				Min_Amount = scheme.Amount_From
			} else {
				if scheme.Amount_From < Min_Amount {
					Min_Amount = scheme.Amount_From
				}
			}
			if Min_AON == 0 {
				Min_AON = scheme.AON_From
			} else {
				if scheme.AON_From < Min_AON {
					Min_AON = scheme.AON_From
				}
			}

			if Amount >= scheme.Amount_From && Amount < scheme.Amount_Till {

				if AON_Months >= scheme.AON_From && AON_Months < scheme.AON_Till {
					return scheme.Key, ""
				}
			}
		}
	}
	if scheme_name == "" {
		if AON_Months < Min_AON {
			NotElligibleReason = "AON"
		}
		if Amount < Min_Amount {
			NotElligibleReason = "Amount"
		}
	}
	return
}

// ///////////////////////////////////////////////////////
// Subscriber
// ///////////////////////////////////////////////////////
func (Uc *UserControl) Subscriber_Add(Login string, request Subscriber_Add_Request) (Id int64, err error) {
	//check key if filled and if already used
	if request.Key == "" {
		err = errors.New("msisdn cannot be empty")
		return Id, err
	}
	//check if key already used
	exits := Map_Subscribers.Check(request.Key)
	if exits {
		err = errors.New("msisdn already exist")
		return Id, err
	}

	//check if exists on IN and get details
	IN_Response, err := Uc.IN.INClient.GetAccountDetails("", "", request.Key)
	if err != nil {
		err = errors.New("error getting account detail: " + err.Error())
		return Id, err
	}

	//to do: get ARPU
	var ARPU float64
	ARPU = 200

	//Prepare new entry
	var NewEntry Subscriber
	NewEntry.Subscriber_Id = Map_AutoIncrement.GetNextAI("Subscriber-Id")
	Id = NewEntry.Subscriber_Id
	NewEntry.Key = request.Key
	NewEntry.COS = IN_Response.COS
	if len(IN_Response.FirstUse) > 10 { //"FirstUse": "2016-06-20T17:51:00.000+00:00"
		First_Use := IN_Response.FirstUse[0:10]
		First_Use_Date, err := time.Parse("2006-01-02", First_Use)
		if err != nil {
			err = errors.New("error parsing FirstUse date: " + err.Error())
			return Id, err
		}
		NewEntry.FirstUse_date = First_Use_Date
	} else {
		err = errors.New("error getting account detail: FirstUse date not defined")
		return Id, err
	}
	NewEntry.ARPU = ARPU
	NewEntry.Last_ProfileUpdate_date = time.Now()
	// NewEntry.Credit_Limit_Scheme, _ = Uc.Credit_Limit_Scheme_Selection(NewEntry.ARPU, NewEntry.FirstUse_date)
	// if NewEntry.Credit_Limit_Scheme != "" {
	// 	NewEntry.IsLendmeEligible = true
	// } else {
	// 	NewEntry.IsLendmeEligible = false
	// }
	//add to cache and DB
	Map_Subscribers.Put(NewEntry.Key, NewEntry)
	return Id, nil
}

func (Uc *UserControl) Subscriber_Edit(Login string, request Subscriber_Edit_Request) (Id int64, err error) {
	//check and validate outlet
	if request.Key == "" {
		err = errors.New("msisdn cannot be empty")
		return Id, err
	}
	subscriber_na, exits := Map_Subscribers.CheckThenGet(request.Key)
	if !exits {
		err = errors.New("msisdn is not created")
		return Id, err
	}
	subscriber, ok := subscriber_na.(Subscriber)
	if !ok {
		return Id, errors.New("error in subscriber type assertion")
	}
	if subscriber.Subscriber_Id != request.Subscriber_Id {
		return Id, errors.New("id is not matching")
	}

	//check if exists on IN and get details
	IN_Response, err := Uc.IN.INClient.GetAccountDetails("", "", request.Key)
	if err != nil {
		err = errors.New("error getting account detail: " + err.Error())
		return Id, err
	}

	//to do: get ARPU
	var ARPU float64
	ARPU = 200

	//Prepare new entry
	subscriber.Key = request.Key
	subscriber.COS = IN_Response.COS
	if len(IN_Response.FirstUse) > 10 { //"FirstUse": "2016-06-20T17:51:00.000+00:00"
		First_Use := IN_Response.FirstUse[0:10]
		First_Use_Date, err := time.Parse("2006-01-02", First_Use)
		if err != nil {
			err = errors.New("error parsing FirstUse date: " + err.Error())
			return Id, err
		}
		subscriber.FirstUse_date = First_Use_Date
	} else {
		err = errors.New("error getting account detail: FirstUse date not defined")
		return Id, err
	}
	subscriber.ARPU = ARPU
	subscriber.Last_ProfileUpdate_date = time.Now()
	// subscriber.Credit_Limit_Scheme, _ = Uc.Credit_Limit_Scheme_Selection(subscriber.ARPU, subscriber.FirstUse_date)
	// if subscriber.Credit_Limit_Scheme != "" {
	// 	subscriber.IsLendmeEligible = true
	// } else {
	// 	subscriber.IsLendmeEligible = false
	// }
	if request.NewKey != "" {
		if request.NewKey != request.Key {
			//delete old
			Map_Subscribers.Delete(request.Key)
			//update key
			subscriber.Key = request.NewKey
		}
	}
	//add to cache and DB
	Map_Subscribers.Put(subscriber.Key, subscriber)
	return Id, nil
}

func (Uc *UserControl) Subscriber_Get(Key string) (subscribers []Subscriber, err error) {
	if Key == "" {
		subscriber_na := Map_Subscribers.ConvertToArray()
		if len(subscriber_na) > 0 {
			for _, entry := range subscriber_na {
				subscriber, ok := entry.(Subscriber)
				if !ok {
					err = errors.New("error in subscriber type assertion")
					return subscribers, err
				} else {
					subscribers = append(subscribers, subscriber)
				}
			}
		}
		return subscribers, nil
	} else {
		subscriber_na, exits := Map_Subscribers.CheckThenGet(Key)
		if !exits {
			err = errors.New("subscriber does not exist")
			return subscribers, err
		}
		subscriber, ok := subscriber_na.(Subscriber)
		if !ok {
			err = errors.New("error in subcriber type assertion")
			return subscribers, err
		}
		subscribers = append(subscribers, subscriber)
		return subscribers, nil
	}
}

func (Uc *UserControl) Subscriber_GetPaginated(Key string, Page, Limit int64) (subscribers []Subscriber, err error) {
	if Page < 1 {
		return subscribers, errors.New("invalid page")
	}
	if Limit < 1 || Limit > 50000 {
		return subscribers, errors.New("invalid limit (accept value between 1 and 50000)")
	}
	if Key == "" {
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
		findResult, err := DAO_Subscribers.FindPaginate(findparams, paginationparams)
		if err != nil {
			return subscribers, err
		}
		if len(findResult) > 0 {
			for _, findres := range findResult {
				InterfaceValue := reflect.ValueOf(findres).Elem().Interface().(Subscriber)
				subscribers = append(subscribers, InterfaceValue)
			}
		}
		return subscribers, nil
	} else {
		subscriber_na, exits := Map_Subscribers.CheckThenGet(Key)
		if !exits {
			err = errors.New("subscriber does not exist")
			return subscribers, err
		}
		subscriber, ok := subscriber_na.(Subscriber)
		if !ok {
			err = errors.New("error in subscriber type assertion")
			return subscribers, err
		}
		subscribers = append(subscribers, subscriber)
		return subscribers, nil
	}
}

func (Uc *UserControl) Subscriber_Delete(Key string) (err error) {
	if Key == "" {
		err = errors.New("msisdn cannot be empty")
		return err
	}
	exits := Map_Subscribers.Check(Key)
	if !exits {
		err = errors.New("msisdn does not exist")
		return err
	}
	Map_Subscribers.Delete(Key)
	return nil
}

// ///////////////////////////////////////////////////////
// Lendme
// ///////////////////////////////////////////////////////
func (Uc *UserControl) Lendme_exec_Request(Source, MSISDN string, Amount float64) (err error) {
	var lendLog Lendme_log
	lendLog.Source = Source
	lendLog.MSISDN = MSISDN
	lendLog.Log_Date = time.Now()
	lendLog.Type = "Lendme Request"
	lendLog.Lendme_Amount = Amount
	lendLog.Lendme_Fee = (Amount * Configuration.Service_FeePerc)

	// check if subscriber exist, auto add if not
	exist := Map_Subscribers.Check(MSISDN)
	if !exist {
		var addRequest Subscriber_Add_Request
		addRequest.Key = MSISDN
		_, err := Uc.Subscriber_Add("Lendme Request", addRequest)
		if err != nil {
			lendLog.Status = "failed"
			lendLog.StatusDescription = err.Error()
			go Uc.Write_Lendme_log(lendLog)
			return err
		}
	}
	subscriber_na := Map_Subscribers.Get(MSISDN)
	subscriber, ok := subscriber_na.(Subscriber)
	if !ok {
		err = errors.New("error in subscriber type assertion")
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		return err
	}
	if !subscriber.IsLendmeEligible {
		err = errors.New("subscriber is not eligible")
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		return err
	}
	//check if exists on IN and get details
	IN_Response, err := Uc.IN.INClient.GetAccountDetails("", "", MSISDN)
	if err != nil {
		err = errors.New("error getting account detail: " + err.Error())
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		return err
	}
	//check balance min and max
	if IN_Response.Balance < Configuration.Min_Allowed_Balance {
		err = errors.New("balance must be greater than " + strconv.FormatFloat(Round(Configuration.Min_Allowed_Balance, 1, 0), 'f', -1, 64))
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		return err
	}
	if IN_Response.Balance > Configuration.Max_Allowed_Balance {
		err = errors.New("balance must be less than " + strconv.FormatFloat(Round(Configuration.Max_Allowed_Balance, 1, 0), 'f', -1, 64))
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		return err
	}
	//check if recharged in the past 60 days
	if len(IN_Response.LastCredited) > 10 {
		LastCredited := IN_Response.LastCredited[0:10]
		LastCredited_Date, err := time.Parse("2006-01-02", LastCredited)
		if err != nil {
			err = errors.New("error parsing LastCredited date: " + err.Error())
			lendLog.Status = "failed"
			lendLog.StatusDescription = err.Error()
			go Uc.Write_Lendme_log(lendLog)
			return err
		}
		//check last credited (should be in more than 60 days)
		durationSinceLastCredited := (time.Now().Sub(LastCredited_Date).Hours() / 24)
		if durationSinceLastCredited < 60 {
			err = errors.New("subscriber must have recharged at least once in a 60-day period")
			lendLog.Status = "failed"
			lendLog.StatusDescription = err.Error()
			go Uc.Write_Lendme_log(lendLog)
			return err
		}
	} else {
		err = errors.New("subscriber must have recharged at least once in a 60-day period")
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		return err
	}
	//check if SOS is active
	if IN_Response.LoyaltyStatus != "X" {
		//todo: opt out from service
		_, errL := Uc.IN.INClient.SetLoyaltyOverdraft("", "", MSISDN, "X", 0)
		if errL != nil {
			err = errors.New("failed to disable SOS: " + errL.Error())
			lendLog.Status = "failed"
			lendLog.StatusDescription = errL.Error()
			go Uc.Write_Lendme_log(lendLog)
			return err
		}
	}

	//check the amount
	if Amount < Configuration.Min_Allowed_Amnt {
		err = errors.New("amount must greater than " + strconv.FormatFloat(Round(Configuration.Min_Allowed_Amnt, 1, 0), 'f', -1, 64))
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		return err
	}
	scheme_na, scheme_exist := Map_Credit_Limit_Scheme.CheckThenGet(subscriber.Credit_Limit_Scheme)
	if !scheme_exist {
		err = errors.New("credit limit schema does not exist")
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		return err
	}
	scheme, ok := scheme_na.(Credit_Limit_Scheme)
	if !ok {
		err = errors.New("error in credit limit schema type assertion")
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		return err
	}

	if Amount > scheme.Credit_limit_Amount {
		err = errors.New("amount is not allowed")
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		return err
	}
	if (subscriber.Lendme_Outstanding_Amount + Amount) > scheme.Credit_limit_Amount {
		err = errors.New("amount is exceeding credit limit")
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		return err
	}
	//credit the amount
	credit_response, credit_err := Uc.IN.INClient.SetAccountBalances("", "", MSISDN, Amount, "N", 0, 0, 0, 0, "Lendme")
	if credit_err != nil {
		lendLog.Status = "failed"
		lendLog.StatusDescription = credit_err.Error()
		go Uc.Write_Lendme_log(lendLog)
		return credit_err
	}
	if credit_response.ResultCode != "0" {
		err = errors.New(credit_response.ResultText)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		return err
	}

	// update subscription
	subscriber.Lendme_Outstanding_Amount = subscriber.Lendme_Outstanding_Amount + Amount
	subscriber.Lendme_Outstanding_Fee = subscriber.Lendme_Outstanding_Fee + (Amount * Configuration.Service_FeePerc)
	subscriber.Last_Lend_Date = time.Now()
	subscriber.Cumulative_Lent_Amount = subscriber.Cumulative_Lent_Amount + Amount
	subscriber.Cumulative_Lent_Fee = subscriber.Cumulative_Lent_Fee + (Amount * Configuration.Service_FeePerc)
	Map_Subscribers.Put(subscriber.Key, subscriber)

	//save log into DB
	lendLog.Status = "successful"
	lendLog.StatusDescription = ""
	go Uc.Write_Lendme_log(lendLog)
	return nil
}

func (Uc *UserControl) Lendme_PayBack(MSISDN string, Amount float64) (err error) {

	return
}

// ///////////////////////////////////////////////////////
// Lendme ARPU
// ///////////////////////////////////////////////////////
