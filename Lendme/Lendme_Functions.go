package Lendme

import (
	"afr_auth_center/AuthCenter"
	"context"
	"daoc"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	apgw "afr_ao_apgw_v2/afr_apgw"
)

const (
	EVC_ACCOUNTNO_PREFIX = "DLR_"
	RBB_BAL_ELEM         = "6010001"
	EVC_BAL_ELEM         = "6010006"
	DEALER_EVC_BAL_ELEM  = "6010000"
	PINEncryptionKey     = "191026789922345378911284"
)

var Map_DefaultValues daoc.Cache_Synch
var DAO_DefaultValues daoc.DAO

var MapAccessEntry daoc.Cache_Synch
var DAO_AccessEntry daoc.DAO

var Map_AutoIncrement daoc.Cache_Synch
var DAO_AutoIncrement daoc.DAO

var Map_Subscribers daoc.Cache_Synch
var DAO_Subscribers daoc.DAO
var DAO_Subscribers_Chrun daoc.DAO

var Map_Credit_Limit_Scheme daoc.Cache_Synch
var DAO_Credit_Limit_Scheme daoc.DAO

var DAO_Lendme_log daoc.DAO

var Map_Lendme_Customer_Exclusion daoc.Cache_Synch
var DAO_Lendme_Customer_Exclusion daoc.DAO

var Map_Lendme_Customer_COS_Exclusion daoc.Cache_Synch
var DAO_Lendme_Customer_COS_Exclusion daoc.DAO

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
	var customer_Exclusion Customer_Exclusion
	Map_Lendme_Customer_Exclusion.Initialize("Customer_Exclusion", "Customer_Exclusion", reflect.TypeOf(Customer_Exclusion{}), customer_Exclusion, true, &DAO_Lendme_Customer_Exclusion, uc.CacheDir.List)
	var customer_COS_Exclusion Customer_COS_Exclusion
	Map_Lendme_Customer_COS_Exclusion.Initialize("Customer_COS_Exclusion", "Customer_COS_Exclusion", reflect.TypeOf(Customer_COS_Exclusion{}), customer_COS_Exclusion, true, &DAO_Lendme_Customer_COS_Exclusion, uc.CacheDir.List)
}

func (uc *UserControl) InitializeDAO() {
	DAO_AccessEntry.Initialize("AccessEntry", uc.MongoDB.MongoDBClient, reflect.TypeOf(AuthCenter.AccessEntry{}), Configuration.DB_Name, "Col_AccessEntry", "")
	DAO_AutoIncrement.Initialize("AutoIncrement", uc.MongoDB.MongoDBClient, reflect.TypeOf(daoc.AutoIncrement{}), Configuration.DB_Name, "Col_AutoIncrement", "")
	DAO_DefaultValues.Initialize("DefaultValues", uc.MongoDB.MongoDBClient, reflect.TypeOf(DefaultValues{}), Configuration.DB_Name, "Col_DefaultValues", "")
	DAO_Subscribers.Initialize("Subscriber", uc.MongoDB.MongoDBClient, reflect.TypeOf(Subscriber{}), Configuration.DB_Name, "Col_Subscriber", "")
	DAO_Subscribers_Chrun.Initialize("Subscriber_Churn", uc.MongoDB.MongoDBClient, reflect.TypeOf(Subscriber{}), Configuration.DB_Name, "Col_Subscriber_Churn", "")
	DAO_Credit_Limit_Scheme.Initialize("Credit_Limit_Scheme", uc.MongoDB.MongoDBClient, reflect.TypeOf(Credit_Limit_Scheme{}), Configuration.DB_Name, "Col_Credit_Limit_Scheme", "")
	DAO_Lendme_log.Initialize("Lendme_log", uc.MongoDB.MongoDBClient, reflect.TypeOf(Lendme_log{}), Configuration.DB_Name, "Col_Lendme_log", "")
	DAO_Lendme_Customer_Exclusion.Initialize("Customer_Exclusion", uc.MongoDB.MongoDBClient, reflect.TypeOf(Customer_Exclusion{}), Configuration.DB_Name, "Col_Customer_Exclusion", "")
	DAO_Lendme_Customer_COS_Exclusion.Initialize("Customer_COS_Exclusion", uc.MongoDB.MongoDBClient, reflect.TypeOf(Customer_COS_Exclusion{}), Configuration.DB_Name, "Col_Customer_COS_Exclusion", "")
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

	exists, err = DAO_Lendme_Customer_Exclusion.CheckAndCreateIndex("Idx_Customer_Exclusion_Key", []string{"Key"}, true)
	if err != nil {
		log.Println("Error creating index Idx_Customer_Exclusion_Key: ", err)
	} else {
		if !exists {
			log.Println("Index Idx_Customer_Exclusion_Key created")
		}
	}

	exists, err = DAO_Lendme_Customer_COS_Exclusion.CheckAndCreateIndex("Idx_Customer_COS_Exclusion_Key", []string{"Key"}, true)
	if err != nil {
		log.Println("Error creating index Idx_Customer_COS_Exclusion_Key: ", err)
	} else {
		if !exists {
			log.Println("Index Idx_Customer_COS_Exclusion_Key created")
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
	if Configuration.Operation == "DRC" {
		uc.Credit_Limit_Scheme_LoadDefaultValues_DRC()
	} else if Configuration.Operation == "Gambia" {
		uc.Credit_Limit_Scheme_LoadDefaultValues_Gambia()
	} else if Configuration.Operation == "SierraLeone" {
		uc.Credit_Limit_Scheme_LoadDefaultValues_SierraLeone()
	} else if Configuration.Operation == "Angola" {
		uc.Credit_Limit_Scheme_LoadDefaultValues_Angola()
	}

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

func (Uc *UserControl) Write_Subscribers_Chrun_log(record Subscriber) {
	YYYY, MM, _, DD, _, _, _ := GetTimeParts(time.Now())
	Db := DAO_Subscribers_Chrun.DB + "_" + YYYY + MM
	Col := DAO_Subscribers_Chrun.Collection + "_" + DD
	_, err := DAO_Subscribers_Chrun.PutOneLogs(record, Db, Col)
	if err != nil {
		log.Println("Error in Write_Subscribers_Chrun_log:", err, " (", record, ")")
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

func (uc *UserControl) Credit_Limit_Scheme_LoadDefaultValues_DRC() {
	log.Println("Loading credit limit scheme default values for DRC")
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

func (uc *UserControl) Credit_Limit_Scheme_LoadDefaultValues_Gambia() {
	log.Println("Loading credit limit scheme default values for Gambia")
	//
	request := Credit_Limit_Scheme_Add_Request{
		Key:                 "5_135_3_6",
		Scheme_Id:           0,
		Amount_From:         5,
		Amount_Till:         135,
		AON_From:            3,
		AON_Till:            6,
		Credit_limit_Amount: 10,
	}
	_, err := uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "5_135_6_12",
		Scheme_Id:           0,
		Amount_From:         5,
		Amount_Till:         135,
		AON_From:            6,
		AON_Till:            12,
		Credit_limit_Amount: 20,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "5_135_12_24",
		Scheme_Id:           0,
		Amount_From:         5,
		Amount_Till:         135,
		AON_From:            12,
		AON_Till:            24,
		Credit_limit_Amount: 25,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "5_135_24_1200",
		Scheme_Id:           0,
		Amount_From:         5,
		Amount_Till:         135,
		AON_From:            24,
		AON_Till:            1200,
		Credit_limit_Amount: 50,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	//
	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "135_270_3_6",
		Scheme_Id:           0,
		Amount_From:         135,
		Amount_Till:         270,
		AON_From:            3,
		AON_Till:            6,
		Credit_limit_Amount: 20,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "135_270_6_12",
		Scheme_Id:           0,
		Amount_From:         135,
		Amount_Till:         270,
		AON_From:            6,
		AON_Till:            12,
		Credit_limit_Amount: 25,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "135_270_12_24",
		Scheme_Id:           0,
		Amount_From:         135,
		Amount_Till:         270,
		AON_From:            12,
		AON_Till:            24,
		Credit_limit_Amount: 50,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "135_270_24_1200",
		Scheme_Id:           0,
		Amount_From:         135,
		Amount_Till:         270,
		AON_From:            24,
		AON_Till:            1200,
		Credit_limit_Amount: 50,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	//
	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "270_340_3_6",
		Scheme_Id:           0,
		Amount_From:         270,
		Amount_Till:         340,
		AON_From:            3,
		AON_Till:            6,
		Credit_limit_Amount: 25,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "270_340_6_12",
		Scheme_Id:           0,
		Amount_From:         270,
		Amount_Till:         340,
		AON_From:            6,
		AON_Till:            12,
		Credit_limit_Amount: 50,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "270_340_12_24",
		Scheme_Id:           0,
		Amount_From:         270,
		Amount_Till:         340,
		AON_From:            12,
		AON_Till:            24,
		Credit_limit_Amount: 50,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "270_340_24_1200",
		Scheme_Id:           0,
		Amount_From:         270,
		Amount_Till:         340,
		AON_From:            24,
		AON_Till:            1200,
		Credit_limit_Amount: 75,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}
	//
	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "340_680_3_6",
		Scheme_Id:           0,
		Amount_From:         340,
		Amount_Till:         680,
		AON_From:            3,
		AON_Till:            6,
		Credit_limit_Amount: 25,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "340_680_6_12",
		Scheme_Id:           0,
		Amount_From:         340,
		Amount_Till:         680,
		AON_From:            6,
		AON_Till:            12,
		Credit_limit_Amount: 50,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "340_680_12_24",
		Scheme_Id:           0,
		Amount_From:         340,
		Amount_Till:         680,
		AON_From:            12,
		AON_Till:            24,
		Credit_limit_Amount: 75,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "340_680_24_1200",
		Scheme_Id:           0,
		Amount_From:         340,
		Amount_Till:         680,
		AON_From:            24,
		AON_Till:            1200,
		Credit_limit_Amount: 75,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	//
	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "680_1300_3_6",
		Scheme_Id:           0,
		Amount_From:         680,
		Amount_Till:         1300,
		AON_From:            3,
		AON_Till:            6,
		Credit_limit_Amount: 50,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "680_1300_6_12",
		Scheme_Id:           0,
		Amount_From:         680,
		Amount_Till:         1300,
		AON_From:            6,
		AON_Till:            12,
		Credit_limit_Amount: 50,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "680_1300_12_24",
		Scheme_Id:           0,
		Amount_From:         680,
		Amount_Till:         1300,
		AON_From:            12,
		AON_Till:            24,
		Credit_limit_Amount: 75,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "680_1300_24_1200",
		Scheme_Id:           0,
		Amount_From:         680,
		Amount_Till:         1300,
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
		Key:                 "1300_2350_3_6",
		Scheme_Id:           0,
		Amount_From:         1300,
		Amount_Till:         2350,
		AON_From:            3,
		AON_Till:            6,
		Credit_limit_Amount: 75,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "1300_2350_6_12",
		Scheme_Id:           0,
		Amount_From:         1300,
		Amount_Till:         2350,
		AON_From:            6,
		AON_Till:            12,
		Credit_limit_Amount: 100,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "1300_2350_12_24",
		Scheme_Id:           0,
		Amount_From:         1300,
		Amount_Till:         2350,
		AON_From:            12,
		AON_Till:            24,
		Credit_limit_Amount: 200,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "1300_2350_24_1200",
		Scheme_Id:           0,
		Amount_From:         1300,
		Amount_Till:         2350,
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
		Key:                 "2350_5000_3_6",
		Scheme_Id:           0,
		Amount_From:         2350,
		Amount_Till:         5000,
		AON_From:            3,
		AON_Till:            6,
		Credit_limit_Amount: 100,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "2350_5000_6_12",
		Scheme_Id:           0,
		Amount_From:         2350,
		Amount_Till:         5000,
		AON_From:            6,
		AON_Till:            12,
		Credit_limit_Amount: 100,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "2350_5000_12_24",
		Scheme_Id:           0,
		Amount_From:         2350,
		Amount_Till:         5000,
		AON_From:            12,
		AON_Till:            24,
		Credit_limit_Amount: 200,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "2350_5000_24_1200",
		Scheme_Id:           0,
		Amount_From:         2350,
		Amount_Till:         5000,
		AON_From:            24,
		AON_Till:            1200,
		Credit_limit_Amount: 300,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	// added 21 Nov 2024
	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "5000_9999999_3_6",
		Scheme_Id:           0,
		Amount_From:         5000,
		Amount_Till:         9999999,
		AON_From:            3,
		AON_Till:            6,
		Credit_limit_Amount: 125,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "5000_9999999_6_12",
		Scheme_Id:           0,
		Amount_From:         5000,
		Amount_Till:         9999999,
		AON_From:            6,
		AON_Till:            12,
		Credit_limit_Amount: 125,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "5000_9999999_12_24",
		Scheme_Id:           0,
		Amount_From:         5000,
		Amount_Till:         9999999,
		AON_From:            12,
		AON_Till:            24,
		Credit_limit_Amount: 250,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "5000_9999999_24_1200",
		Scheme_Id:           0,
		Amount_From:         5000,
		Amount_Till:         9999999,
		AON_From:            24,
		AON_Till:            1200,
		Credit_limit_Amount: 600,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

}

func (uc *UserControl) Credit_Limit_Scheme_LoadDefaultValues_SierraLeone() {
	log.Println("Loading credit limit scheme default values for SierraLeone")
	//
	request := Credit_Limit_Scheme_Add_Request{
		Key:                 "0_33_3_6",
		Scheme_Id:           0,
		Amount_From:         0,
		Amount_Till:         33,
		AON_From:            3,
		AON_Till:            6,
		Credit_limit_Amount: 3,
	}
	_, err := uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "0_33_6_12",
		Scheme_Id:           0,
		Amount_From:         0,
		Amount_Till:         33,
		AON_From:            6,
		AON_Till:            12,
		Credit_limit_Amount: 3,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "0_33_12_24",
		Scheme_Id:           0,
		Amount_From:         0,
		Amount_Till:         33,
		AON_From:            12,
		AON_Till:            24,
		Credit_limit_Amount: 6,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "0_33_24_1200",
		Scheme_Id:           0,
		Amount_From:         0,
		Amount_Till:         33,
		AON_From:            24,
		AON_Till:            1200,
		Credit_limit_Amount: 6,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	//
	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "33_65_3_6",
		Scheme_Id:           0,
		Amount_From:         33,
		Amount_Till:         65,
		AON_From:            3,
		AON_Till:            6,
		Credit_limit_Amount: 6,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "33_65_6_12",
		Scheme_Id:           0,
		Amount_From:         33,
		Amount_Till:         65,
		AON_From:            6,
		AON_Till:            12,
		Credit_limit_Amount: 6,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "33_65_12_24",
		Scheme_Id:           0,
		Amount_From:         33,
		Amount_Till:         65,
		AON_From:            12,
		AON_Till:            24,
		Credit_limit_Amount: 10,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "33_65_24_1200",
		Scheme_Id:           0,
		Amount_From:         33,
		Amount_Till:         65,
		AON_From:            24,
		AON_Till:            1200,
		Credit_limit_Amount: 13,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	//
	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "65_88_3_6",
		Scheme_Id:           0,
		Amount_From:         65,
		Amount_Till:         88,
		AON_From:            3,
		AON_Till:            6,
		Credit_limit_Amount: 7,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "65_88_6_12",
		Scheme_Id:           0,
		Amount_From:         65,
		Amount_Till:         88,
		AON_From:            6,
		AON_Till:            12,
		Credit_limit_Amount: 7,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "65_88_12_24",
		Scheme_Id:           0,
		Amount_From:         65,
		Amount_Till:         88,
		AON_From:            12,
		AON_Till:            24,
		Credit_limit_Amount: 11,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "65_88_24_1200",
		Scheme_Id:           0,
		Amount_From:         65,
		Amount_Till:         88,
		AON_From:            24,
		AON_Till:            1200,
		Credit_limit_Amount: 15,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	//
	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "88_175_3_6",
		Scheme_Id:           0,
		Amount_From:         88,
		Amount_Till:         175,
		AON_From:            3,
		AON_Till:            6,
		Credit_limit_Amount: 8,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "88_175_6_12",
		Scheme_Id:           0,
		Amount_From:         88,
		Amount_Till:         175,
		AON_From:            6,
		AON_Till:            12,
		Credit_limit_Amount: 8,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "88_175_12_24",
		Scheme_Id:           0,
		Amount_From:         88,
		Amount_Till:         175,
		AON_From:            12,
		AON_Till:            24,
		Credit_limit_Amount: 15,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "88_175_24_1200",
		Scheme_Id:           0,
		Amount_From:         88,
		Amount_Till:         175,
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
		Key:         "175_375_3_6",
		Scheme_Id:   0,
		Amount_From: 175,
		Amount_Till: 375,
		AON_From:    3,
		AON_Till:    6,
		// Credit_limit_Amount: 8,
		Credit_limit_Amount: 18,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:         "175_375_6_12",
		Scheme_Id:   0,
		Amount_From: 175,
		Amount_Till: 375,
		AON_From:    6,
		AON_Till:    12,
		// Credit_limit_Amount: 8,
		Credit_limit_Amount: 20,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:         "175_375_12_24",
		Scheme_Id:   0,
		Amount_From: 175,
		Amount_Till: 375,
		AON_From:    12,
		AON_Till:    24,
		// Credit_limit_Amount: 15,
		Credit_limit_Amount: 30,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:         "175_375_24_1200",
		Scheme_Id:   0,
		Amount_From: 175,
		Amount_Till: 375,
		AON_From:    24,
		AON_Till:    1200,
		// Credit_limit_Amount: 20,
		Credit_limit_Amount: 50,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	//
	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "375_760_3_6",
		Scheme_Id:           0,
		Amount_From:         375,
		Amount_Till:         760,
		AON_From:            3,
		AON_Till:            6,
		Credit_limit_Amount: 35,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "375_760_6_12",
		Scheme_Id:           0,
		Amount_From:         375,
		Amount_Till:         760,
		AON_From:            6,
		AON_Till:            12,
		Credit_limit_Amount: 50,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "375_760_12_24",
		Scheme_Id:           0,
		Amount_From:         375,
		Amount_Till:         760,
		AON_From:            12,
		AON_Till:            24,
		Credit_limit_Amount: 60,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "375_760_24_1200",
		Scheme_Id:           0,
		Amount_From:         375,
		Amount_Till:         760,
		AON_From:            24,
		AON_Till:            1200,
		Credit_limit_Amount: 70,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	//
	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "760_9999999_3_6",
		Scheme_Id:           0,
		Amount_From:         760,
		Amount_Till:         9999999,
		AON_From:            3,
		AON_Till:            6,
		Credit_limit_Amount: 40,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "760_9999999_6_12",
		Scheme_Id:           0,
		Amount_From:         760,
		Amount_Till:         9999999,
		AON_From:            6,
		AON_Till:            12,
		Credit_limit_Amount: 60,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "760_9999999_12_24",
		Scheme_Id:           0,
		Amount_From:         760,
		Amount_Till:         9999999,
		AON_From:            12,
		AON_Till:            24,
		Credit_limit_Amount: 70,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "760_9999999_24_1200",
		Scheme_Id:           0,
		Amount_From:         760,
		Amount_Till:         9999999,
		AON_From:            24,
		AON_Till:            1200,
		Credit_limit_Amount: 80,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}
}

func (uc *UserControl) Credit_Limit_Scheme_LoadDefaultValues_Angola() {
	log.Println("Loading credit limit scheme default values for Angola")
	//
	request := Credit_Limit_Scheme_Add_Request{
		Key:                 "0_200_2000_3_9",
		Scheme_Id:           0,
		Amount_From:         200,
		Amount_Till:         2000,
		AON_From:            3,
		AON_Till:            9,
		Credit_limit_Amount: 250,
	}
	_, err := uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "0_2001_5000_3_9",
		Scheme_Id:           0,
		Amount_From:         2001,
		Amount_Till:         5000,
		AON_From:            3,
		AON_Till:            9,
		Credit_limit_Amount: 450,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "0_5001_999999999_3_9",
		Scheme_Id:           0,
		Amount_From:         5001,
		Amount_Till:         999999999,
		AON_From:            3,
		AON_Till:            9,
		Credit_limit_Amount: 1300,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "0_200_2000_9_24",
		Scheme_Id:           0,
		Amount_From:         200,
		Amount_Till:         2000,
		AON_From:            9,
		AON_Till:            24,
		Credit_limit_Amount: 450,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "0_2001_5000_9_24",
		Scheme_Id:           0,
		Amount_From:         2001,
		Amount_Till:         5000,
		AON_From:            9,
		AON_Till:            24,
		Credit_limit_Amount: 1300,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "0_5001_999999999_9_24",
		Scheme_Id:           0,
		Amount_From:         5001,
		Amount_Till:         999999999,
		AON_From:            9,
		AON_Till:            24,
		Credit_limit_Amount: 2200,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "0_200_2000_24_9999999",
		Scheme_Id:           0,
		Amount_From:         200,
		Amount_Till:         2000,
		AON_From:            24,
		AON_Till:            9999999,
		Credit_limit_Amount: 1300,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "0_2001_5000_24_9999999",
		Scheme_Id:           0,
		Amount_From:         2001,
		Amount_Till:         5000,
		AON_From:            24,
		AON_Till:            9999999,
		Credit_limit_Amount: 1800,
	}
	_, err = uc.Credit_Limit_Scheme_Add("LoadDefaultValues", request)
	if err != nil {
		log.Println("error adding credit limit scheme (" + request.Key + "): " + err.Error())
	}

	request = Credit_Limit_Scheme_Add_Request{
		Key:                 "0_5001_999999999_24_9999999",
		Scheme_Id:           0,
		Amount_From:         5001,
		Amount_Till:         999999999,
		AON_From:            24,
		AON_Till:            9999999,
		Credit_limit_Amount: 2700,
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
		DailyImportSubsStats.With(prometheus.Labels{"IsElligble": "false", "Reason": NotElligibleReason, "Scheme": ""}).Inc()
		return
	}
	if AON_Months < Configuration.Min_Allowed_AON {
		NotElligibleReason = "Min AON"
		DailyImportSubsStats.With(prometheus.Labels{"IsElligble": "false", "Reason": NotElligibleReason, "Scheme": ""}).Inc()
		return
	}
	if Amount < Configuration.Min_Avg3MRecharge {
		NotElligibleReason = "Min Avg3MRecharge"
		DailyImportSubsStats.With(prometheus.Labels{"IsElligble": "false", "Reason": NotElligibleReason, "Scheme": ""}).Inc()
		return
	}

	LastRecharge_Hours := time.Since(LastRecharge_date).Hours()
	LastRecharge_Days := (LastRecharge_Hours / 24)
	if LastRecharge_date.IsZero() {
		NotElligibleReason = "Min Last Recharge Period"
		DailyImportSubsStats.With(prometheus.Labels{"IsElligble": "false", "Reason": NotElligibleReason, "Scheme": ""}).Inc()
		return
	}
	if LastRecharge_Days > Configuration.Min_LastRechargePeriod {
		NotElligibleReason = "Min Last Recharge Period"
		DailyImportSubsStats.With(prometheus.Labels{"IsElligble": "false", "Reason": NotElligibleReason, "Scheme": ""}).Inc()
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
					DailyImportSubsStats.With(prometheus.Labels{"IsElligble": "true", "Reason": "", "Scheme": scheme.Key}).Inc()
					DailyImportSubsStats.With(prometheus.Labels{"IsElligble": "true", "Reason": "", "Scheme": "All"}).Inc()
					return scheme.Key, ""
				}
			}
		}
	}
	if scheme_name == "" {
		if AON_Months < Min_AON {
			NotElligibleReason = "AON"
			DailyImportSubsStats.With(prometheus.Labels{"IsElligble": "false", "Reason": NotElligibleReason, "Scheme": ""}).Inc()
		}
		if Amount < Min_Amount {
			NotElligibleReason = "Amount"
			DailyImportSubsStats.With(prometheus.Labels{"IsElligble": "false", "Reason": NotElligibleReason, "Scheme": ""}).Inc()
		}
	}
	return
}

// ///////////////////////////////////////////////////////
// Subscriber
// ///////////////////////////////////////////////////////
func (Uc *UserControl) Subscriber_Add(Login string, request Subscriber_Add_Request) (Id int64, err error) {
	err = errors.New("add subscriber is restricted to the daily importing process")
	return Id, err
	//check key if filled and if already used
	// 	if request.Key == "" {
	// 		err = errors.New("msisdn cannot be empty")
	// 		return Id, err
	// 	}
	// 	//check if key already used
	// 	exits := Map_Subscribers.Check(request.Key)
	// 	if exits {
	// 		err = errors.New("msisdn already exist")
	// 		return Id, err
	// 	}

	// 	//check if exists on IN and get details
	// 	IN_Response, err := Uc.IN.INClient.GetAccountDetails("", "", request.Key)
	// 	if err != nil {
	// 		err = errors.New("error getting account detail: " + err.Error())
	// 		return Id, err
	// 	}

	// 	//Prepare new entry
	// 	var NewEntry Subscriber
	// 	NewEntry.Subscriber_Id = Map_AutoIncrement.GetNextAI("Subscriber-Id")
	// 	Id = NewEntry.Subscriber_Id
	// 	NewEntry.Key = request.Key
	// 	NewEntry.COS = IN_Response.COS
	// 	if len(IN_Response.FirstUse) > 10 { //"FirstUse": "2016-06-20T17:51:00.000+00:00"
	// 		First_Use := IN_Response.FirstUse[0:10]
	// 		First_Use_Date, err := time.Parse("2006-01-02", First_Use)
	// 		if err != nil {
	// 			err = errors.New("error parsing FirstUse date: " + err.Error())
	// 			return Id, err
	// 		}
	// 		NewEntry.FirstUse_date = First_Use_Date
	// 	} else {
	// 		err = errors.New("error getting account detail: FirstUse date not defined")
	// 		return Id, err
	// 	}
	// 	NewEntry.ARPU = ARPU
	// 	NewEntry.Last_ProfileUpdate_date = time.Now()
	// 	// NewEntry.Credit_Limit_Scheme, _ = Uc.Credit_Limit_Scheme_Selection(NewEntry.ARPU, NewEntry.FirstUse_date)
	// 	// if NewEntry.Credit_Limit_Scheme != "" {
	// 	// 	NewEntry.IsLendmeEligible = true
	// 	// } else {
	// 	// 	NewEntry.IsLendmeEligible = false
	// 	// }
	// 	//add to cache and DB
	// 	Map_Subscribers.Put(NewEntry.Key, NewEntry)
	// 	return Id, nil
}

func (Uc *UserControl) Subscriber_Edit(Login string, request Subscriber_Edit_Request) (Id int64, err error) {
	err = errors.New("edit subscriber is restricted to the daily importing process")
	return Id, err
	// //check and validate outlet
	// if request.Key == "" {
	// 	err = errors.New("msisdn cannot be empty")
	// 	return Id, err
	// }
	// subscriber_na, exits := Map_Subscribers.CheckThenGet(request.Key)
	// if !exits {
	// 	err = errors.New("msisdn is not created")
	// 	return Id, err
	// }
	// subscriber, ok := subscriber_na.(Subscriber)
	// if !ok {
	// 	return Id, errors.New("error in subscriber type assertion")
	// }
	// if subscriber.Subscriber_Id != request.Subscriber_Id {
	// 	return Id, errors.New("id is not matching")
	// }

	// //check if exists on IN and get details
	// IN_Response, err := Uc.IN.INClient.GetAccountDetails("", "", request.Key)
	// if err != nil {
	// 	err = errors.New("error getting account detail: " + err.Error())
	// 	return Id, err
	// }

	// //to do: get ARPU
	// var ARPU float64
	// ARPU = 200

	// //Prepare new entry
	// subscriber.Key = request.Key
	// subscriber.COS = IN_Response.COS
	// if len(IN_Response.FirstUse) > 10 { //"FirstUse": "2016-06-20T17:51:00.000+00:00"
	// 	First_Use := IN_Response.FirstUse[0:10]
	// 	First_Use_Date, err := time.Parse("2006-01-02", First_Use)
	// 	if err != nil {
	// 		err = errors.New("error parsing FirstUse date: " + err.Error())
	// 		return Id, err
	// 	}
	// 	subscriber.FirstUse_date = First_Use_Date
	// } else {
	// 	err = errors.New("error getting account detail: FirstUse date not defined")
	// 	return Id, err
	// }
	// subscriber.ARPU = ARPU
	// subscriber.Last_ProfileUpdate_date = time.Now()
	// // subscriber.Credit_Limit_Scheme, _ = Uc.Credit_Limit_Scheme_Selection(subscriber.ARPU, subscriber.FirstUse_date)
	// // if subscriber.Credit_Limit_Scheme != "" {
	// // 	subscriber.IsLendmeEligible = true
	// // } else {
	// // 	subscriber.IsLendmeEligible = false
	// // }
	// if request.NewKey != "" {
	// 	if request.NewKey != request.Key {
	// 		//delete old
	// 		Map_Subscribers.Delete(request.Key)
	// 		//update key
	// 		subscriber.Key = request.NewKey
	// 	}
	// }
	// //add to cache and DB
	// Map_Subscribers.Put(subscriber.Key, subscriber)
	// return Id, nil
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
	subscriber_na, exits := Map_Subscribers.CheckThenGet(Key)
	if !exits {
		err = errors.New("msisdn does not exist")
		return err
	}
	subscriber, ok := subscriber_na.(Subscriber)
	if !ok {
		err = errors.New("error in subscriber type assertion")
		return err
	}
	Uc.Write_Subscribers_Chrun_log(subscriber)
	Map_Subscribers.Delete(Key)
	return nil
}

// ///////////////////////////////////////////////////////
// Lendme
// ///////////////////////////////////////////////////////
func (Uc *UserControl) Lendme_exec_Request(Source, MSISDN string, Amount float64) (err error) {
	log.Println("Lendme Request for "+MSISDN+": ", Amount)
	var lendLog Lendme_log
	lendLog.Source = Source
	lendLog.MSISDN = MSISDN
	lendLog.Log_Date = time.Now()
	lendLog.Type = "Lendme Request"
	lendLog.Lendme_Amount = Amount
	lendLog.Lendme_Fee = (Amount * Configuration.Service_FeePerc)

	// check if subscriber exist, return error if not existing
	subscriber_na, exist := Map_Subscribers.CheckThenGet(MSISDN)
	if !exist {
		error_msg := "subscriber does not exist in the service pool"
		err = errors.New(error_msg)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)

		LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Inc()
		LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Add(Amount)
		return err
		// var addRequest Subscriber_Add_Request
		// addRequest.Key = MSISDN
		// _, err := Uc.Subscriber_Add("Lendme Request", addRequest)
		// if err != nil {
		// 	lendLog.Status = "failed"
		// 	lendLog.StatusDescription = err.Error()
		// 	go Uc.Write_Lendme_log(lendLog)
		// 	return err
		// }
	}
	subscriber, ok := subscriber_na.(Subscriber)
	if !ok {
		error_msg := "error in subscriber type assertion"
		err = errors.New(error_msg)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Inc()
		LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Add(Amount)
		return err
	}
	lendLog.Lendme_Outstanding_Amount = subscriber.Lendme_Outstanding_Amount
	lendLog.Lendme_Outstanding_Fee = subscriber.Lendme_Outstanding_Fee
	if !subscriber.IsLendmeEligible {
		error_msg := "subscriber is not eligible"
		err = errors.New(error_msg)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Inc()
		LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Add(Amount)
		return err
	}
	IN_MSISDN := MSISDN
	if Configuration.Operation == "Gambia" {
		if len(MSISDN) > 7 {
			IN_MSISDN = IN_MSISDN[len(MSISDN)-7 : len(MSISDN)]
		}
	} else if Configuration.Operation == "SierraLeone" { //077928014
		if len(MSISDN) > 8 {
			IN_MSISDN = "0" + IN_MSISDN[len(MSISDN)-8:len(MSISDN)]
		}
	}
	//check if exists on IN and get details
	IN_Response, err := Uc.IN.INClient.GetAccountDetails("", "", IN_MSISDN)
	if err != nil {
		error_msg := "error getting account detail: " + err.Error()
		err = errors.New(error_msg)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Inc()
		LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Add(Amount)
		return err
	}
	lendLog.Subscriber_OpeningBlance = IN_Response.Balance
	//check balance min and max
	if IN_Response.Balance < Configuration.Min_Allowed_Balance {
		err = errors.New("balance must be greater than " + strconv.FormatFloat(Round(Configuration.Min_Allowed_Balance, 1, 0), 'f', -1, 64))
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		error_msg := "balance must be greater than minimum allowed"
		LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Inc()
		LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Add(Amount)
		return err
	}
	if IN_Response.Balance > Configuration.Max_Allowed_Balance {
		err = errors.New("balance must be less than " + strconv.FormatFloat(Round(Configuration.Max_Allowed_Balance, 1, 0), 'f', -1, 64))
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		error_msg := "balance must be less than maximum allowed"
		LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Inc()
		LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Add(Amount)
		return err
	}

	//check if recharged in the past 60 days --> done in the daily importing process
	// if len(IN_Response.LastCredited) > 10 {
	// 	LastCredited := IN_Response.LastCredited[0:10]
	// 	LastCredited_Date, err := time.Parse("2006-01-02", LastCredited)
	// 	if err != nil {
	// 		err = errors.New("error parsing LastCredited date: " + err.Error())
	// 		lendLog.Status = "failed"
	// 		lendLog.StatusDescription = err.Error()
	// 		go Uc.Write_Lendme_log(lendLog)
	// 		return err
	// 	}
	// 	//check last credited (should be in more than 60 days)
	// 	durationSinceLastCredited := (time.Now().Sub(LastCredited_Date).Hours() / 24)
	// 	if durationSinceLastCredited < 60 {
	// 		err = errors.New("subscriber must have recharged at least once in a 60-day period")
	// 		lendLog.Status = "failed"
	// 		lendLog.StatusDescription = err.Error()
	// 		go Uc.Write_Lendme_log(lendLog)
	// 		return err
	// 	}
	// } else {
	// 	err = errors.New("subscriber must have recharged at least once in a 60-day period")
	// 	lendLog.Status = "failed"
	// 	lendLog.StatusDescription = err.Error()
	// 	go Uc.Write_Lendme_log(lendLog)
	// 	return err
	// }
	//check if SOS is active
	if IN_Response.LoyaltyStatus != "X" {
		IN_MSISDN := MSISDN
		if Configuration.Operation == "Gambia" {
			if len(MSISDN) > 7 {
				IN_MSISDN = IN_MSISDN[len(MSISDN)-7 : len(MSISDN)]
			}
		} else if Configuration.Operation == "SierraLeone" { //077928014
			if len(MSISDN) > 8 {
				IN_MSISDN = "0" + IN_MSISDN[len(MSISDN)-8:len(MSISDN)]
			}
		}
		//todo: opt out from service
		_, errL := Uc.IN.INClient.SetLoyaltyOverdraft("", "", IN_MSISDN, "X", 0)
		if errL != nil {
			error_msg := "failed to disable SOS: " + errL.Error()
			err = errors.New(error_msg)
			lendLog.Status = "failed"
			lendLog.StatusDescription = errL.Error()
			go Uc.Write_Lendme_log(lendLog)
			LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Inc()
			LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Add(Amount)
			return err
		}
	}

	//check the amount
	if subscriber.Lendme_Outstanding_Amount == 0 {
		if Amount < Configuration.Min_Allowed_Amnt {
			err = errors.New("amount must greater than " + strconv.FormatFloat(Round(Configuration.Min_Allowed_Amnt, 1, 0), 'f', -1, 64))
			lendLog.Status = "failed"
			lendLog.StatusDescription = err.Error()
			go Uc.Write_Lendme_log(lendLog)
			error_msg := "requested amount must greater than minimum allowed amount"
			LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Inc()
			LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Add(Amount)
			return err
		}
	}
	scheme_na, scheme_exist := Map_Credit_Limit_Scheme.CheckThenGet(subscriber.Credit_Limit_Scheme)
	if !scheme_exist {
		error_msg := "credit limit schema does not exist"
		err = errors.New(error_msg)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Inc()
		LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Add(Amount)
		return err
	}
	scheme, ok := scheme_na.(Credit_Limit_Scheme)
	if !ok {
		error_msg := "error in credit limit schema type assertion"
		err = errors.New(error_msg)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Inc()
		LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Add(Amount)
		return err
	}

	if Amount > scheme.Credit_limit_Amount {
		error_msg := "requested amount is not allowed"
		err = errors.New(error_msg)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Inc()
		LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Add(Amount)
		return err
	}
	if (subscriber.Lendme_Outstanding_Amount + Amount) > scheme.Credit_limit_Amount {
		error_msg := "requested amount is exceeding credit limit"
		err = errors.New(error_msg)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Inc()
		LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Add(Amount)
		return err
	}
	//credit the amount
	credit_response, credit_err := Uc.IN.INClient.SetAccountBalances("", "", IN_MSISDN, Amount, "N", 0, 0, 0, 0, "Lendme")
	if credit_err != nil {
		lendLog.Status = "failed"
		lendLog.StatusDescription = credit_err.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": credit_err.Error(), "Scheme": ""}).Inc()
		LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": credit_err.Error(), "Scheme": ""}).Add(Amount)
		return credit_err
	}
	if credit_response.ResultCode != "0" {
		err = errors.New(credit_response.ResultText)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": err.Error(), "Scheme": ""}).Inc()
		LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": err.Error(), "Scheme": ""}).Add(Amount)
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
	LendMeRequestsCount.With(prometheus.Labels{"Status": "successful", "Reason": "", "Scheme": subscriber.Credit_Limit_Scheme}).Inc()
	LendMeRequestsAmount.With(prometheus.Labels{"Status": "successful", "Reason": "", "Scheme": subscriber.Credit_Limit_Scheme}).Add(Amount)
	return nil
}

func (Uc *UserControl) LendmeAO_exec_Request(Source, MSISDN string, Amount float64) (err error) {
	log.Println("Lendme Request for "+MSISDN+": ", Amount)
	var lendLog Lendme_log
	lendLog.Source = Source
	lendLog.MSISDN = MSISDN
	lendLog.Log_Date = time.Now()
	lendLog.Type = "Lendme Request"
	lendLog.Lendme_Amount = Amount
	lendLog.Lendme_Fee = (Amount * Configuration.Service_FeePerc)

	// check if subscriber exist, return error if not existing
	subscriber_na, exist := Map_Subscribers.CheckThenGet(MSISDN)
	if !exist {
		error_msg := "subscriber does not exist in the service pool"
		err = errors.New(error_msg)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)

		LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Inc()
		LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Add(Amount)
		return err
		// var addRequest Subscriber_Add_Request
		// addRequest.Key = MSISDN
		// _, err := Uc.Subscriber_Add("Lendme Request", addRequest)
		// if err != nil {
		// 	lendLog.Status = "failed"
		// 	lendLog.StatusDescription = err.Error()
		// 	go Uc.Write_Lendme_log(lendLog)
		// 	return err
		// }
	}
	subscriber, ok := subscriber_na.(Subscriber)
	if !ok {
		error_msg := "error in subscriber type assertion"
		err = errors.New(error_msg)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Inc()
		LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Add(Amount)
		return err
	}
	lendLog.Lendme_Outstanding_Amount = subscriber.Lendme_Outstanding_Amount
	lendLog.Lendme_Outstanding_Fee = subscriber.Lendme_Outstanding_Fee
	if !subscriber.IsLendmeEligible {
		error_msg := "subscriber is not eligible"
		err = errors.New(error_msg)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Inc()
		LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Add(Amount)
		return err
	}
	OpeningBalance_RBB, OpeningBalance_EVC, balOpenErr := Uc.GetCustomerMainBalance(MSISDN)
	if balOpenErr != nil {
		lendLog.Status = "failed"
		lendLog.StatusDescription = balOpenErr.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": balOpenErr.Error(), "Scheme": ""}).Inc()
		LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": balOpenErr.Error(), "Scheme": ""}).Add(Amount)
		return balOpenErr
	}
	lendLog.Subscriber_OpeningBlance = OpeningBalance_RBB + OpeningBalance_EVC

	//check balance min and max
	if lendLog.Subscriber_OpeningBlance < Configuration.Min_Allowed_Balance {
		err = errors.New("balance must be greater than " + strconv.FormatFloat(Round(Configuration.Min_Allowed_Balance, 1, 0), 'f', -1, 64))
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		error_msg := "balance must be greater than minimum allowed"
		LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Inc()
		LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Add(Amount)
		return err
	}
	// if lendLog.Subscriber_OpeningBlance > Configuration.Max_Allowed_Balance {
	// 	err = errors.New("balance must be less than " + strconv.FormatFloat(Round(Configuration.Max_Allowed_Balance, 1, 0), 'f', -1, 64))
	// 	lendLog.Status = "failed"
	// 	lendLog.StatusDescription = err.Error()
	// 	go Uc.Write_Lendme_log(lendLog)
	// 	error_msg := "balance must be less than maximum allowed"
	// 	LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Inc()
	// 	LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Add(Amount)
	// 	return err
	// }

	//check the amount
	if subscriber.Lendme_Outstanding_Amount == 0 {
		if Amount < Configuration.Min_Allowed_Amnt {
			err = errors.New("amount must greater than " + strconv.FormatFloat(Round(Configuration.Min_Allowed_Amnt, 1, 0), 'f', -1, 64))
			lendLog.Status = "failed"
			lendLog.StatusDescription = err.Error()
			go Uc.Write_Lendme_log(lendLog)
			error_msg := "requested amount must greater than minimum allowed amount"
			LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Inc()
			LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Add(Amount)
			return err
		}
	}
	scheme_na, scheme_exist := Map_Credit_Limit_Scheme.CheckThenGet(subscriber.Credit_Limit_Scheme)
	if !scheme_exist {
		error_msg := "credit limit schema does not exist"
		err = errors.New(error_msg)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Inc()
		LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Add(Amount)
		return err
	}
	scheme, ok := scheme_na.(Credit_Limit_Scheme)
	if !ok {
		error_msg := "error in credit limit schema type assertion"
		err = errors.New(error_msg)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Inc()
		LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Add(Amount)
		return err
	}

	if Amount > scheme.Credit_limit_Amount {
		error_msg := "requested amount is not allowed"
		err = errors.New(error_msg)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Inc()
		LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Add(Amount)
		return err
	}
	if (subscriber.Lendme_Outstanding_Amount + Amount) > scheme.Credit_limit_Amount {
		error_msg := "requested amount is exceeding credit limit"
		err = errors.New(error_msg)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Inc()
		LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Scheme": ""}).Add(Amount)
		return err
	}

	//credit the amount
	DBT_Response, credit_err := CS_DealerBalanceTransfer("customer", Configuration.Lendme_EVC_Dealer_MSISDN, Configuration.Lendme_EVC_Dealer_PIN, MSISDN, strconv.Itoa(int(Amount)))
	if credit_err != nil {
		lendLog.Status = "failed"
		lendLog.StatusDescription = credit_err.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": credit_err.Error(), "Scheme": ""}).Inc()
		LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": credit_err.Error(), "Scheme": ""}).Add(Amount)

		// Post to Kafka
		var evc_Recharge_request EVC_Recharge_request
		evc_Recharge_request.TransID = ""
		evc_Recharge_request.TransStatus = "failed"
		evc_Recharge_request.TransStatusDescription = credit_err.Error()
		evc_Recharge_request.DealerMSISDN = Configuration.Lendme_EVC_Dealer_MSISDN
		evc_Recharge_request.DealerName = "Lendme"
		evc_Recharge_request.TargetMSISDN = MSISDN
		evc_Recharge_request.Amount = Amount
		evc_Recharge_request.GSMLocation = ""
		go Post_EVC_Recharge_ToKafka(evc_Recharge_request)

		return credit_err
	} else if DBT_Response.Response.Success != "true" {
		err_ret := errors.New("CS success flag is not true")
		lendLog.Status = "failed"
		lendLog.StatusDescription = "CS success flag is not true. " + DBT_Response.Response.Result.Message
		go Uc.Write_Lendme_log(lendLog)
		LendMeRequestsCount.With(prometheus.Labels{"Status": "failed", "Reason": "CS success flag is not true", "Scheme": ""}).Inc()
		LendMeRequestsAmount.With(prometheus.Labels{"Status": "failed", "Reason": "CS success flag is not true", "Scheme": ""}).Add(Amount)

		var evc_Recharge_request EVC_Recharge_request
		evc_Recharge_request.TransID = ""
		evc_Recharge_request.TransStatus = "failed"
		evc_Recharge_request.TransStatusDescription = "CDS response Success flag is not true"
		evc_Recharge_request.DealerMSISDN = Configuration.Lendme_EVC_Dealer_MSISDN
		evc_Recharge_request.DealerName = "Lendme"
		evc_Recharge_request.TargetMSISDN = MSISDN
		evc_Recharge_request.Amount = Amount
		evc_Recharge_request.GSMLocation = ""
		go Post_EVC_Recharge_ToKafka(evc_Recharge_request)

		return err_ret
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
	LendMeRequestsCount.With(prometheus.Labels{"Status": "successful", "Reason": "", "Scheme": subscriber.Credit_Limit_Scheme}).Inc()
	LendMeRequestsAmount.With(prometheus.Labels{"Status": "successful", "Reason": "", "Scheme": subscriber.Credit_Limit_Scheme}).Add(Amount)

	var evc_Recharge_request EVC_Recharge_request
	evc_Recharge_request.TransID = DBT_Response.Response.Result.Arguments.TransID
	evc_Recharge_request.TransStatus = "success"
	evc_Recharge_request.TransStatusDescription = ""
	evc_Recharge_request.DealerMSISDN = Configuration.Lendme_EVC_Dealer_MSISDN
	evc_Recharge_request.DealerName = "Lendme"
	evc_Recharge_request.TargetMSISDN = MSISDN
	evc_Recharge_request.Amount = Amount
	evc_Recharge_request.GSMLocation = ""
	if len(DBT_Response.Response.Result.Arguments.AfricellFldSourceAccount.BalanceGroups) > 0 {
		if len(DBT_Response.Response.Result.Arguments.AfricellFldSourceAccount.BalanceGroups[0].Balances) > 0 {
			for _, value := range DBT_Response.Response.Result.Arguments.AfricellFldSourceAccount.BalanceGroups[0].Balances {
				if value.Elem == "6010000" {
					bal_float, err := strconv.ParseFloat(strings.TrimSpace(value.CurrentBal), 64)
					if err != nil {
						fmt.Println("Failed to parse dealer endbal on recharge ", err.Error())
						break
					}
					evc_Recharge_request.DealerClosingBalance = bal_float * -1
				}
			}
		}
	}

	go Post_EVC_Recharge_ToKafka(evc_Recharge_request)

	return nil
}

func (Uc *UserControl) GetCustomerMainBalance(MSISDN string) (balance_RBB float64, balance_EVC float64, err error) {
	Balance_Response, err := CS_GetAccountBalance_ByDeviceId(MSISDN)
	if err != nil {
		return balance_RBB, balance_EVC, err
	}
	if Balance_Response.Response.Success != "true" {
		errN := errors.New("success is false")
		return balance_RBB, balance_EVC, errN
	}
	for _, balanceGroupes := range Balance_Response.Response.Result.Arguments.BalanceGroups {
		for _, balances := range balanceGroupes.Balances {
			RBB_BAL_ELEM_int, _ := strconv.Atoi(RBB_BAL_ELEM)
			if balances.Elem == RBB_BAL_ELEM_int {
				balance_RBB = balances.CurrentBal * -1
			}
			EVC_BAL_ELEM_int, _ := strconv.Atoi(EVC_BAL_ELEM)
			if balances.Elem == EVC_BAL_ELEM_int {
				balance_EVC = balances.CurrentBal * -1
			}
		}
	}
	return balance_RBB, balance_EVC, nil
}

func (Uc *UserControl) Lendme_PayBack(Source, MSISDN string, RechargeAmount float64, Opid string) (err error) {
	if Opid == "lendme" {
		LendMePayBackCount.With(prometheus.Labels{"Status": "Ignored", "Reason": "Opid is lendme", "Description": Source}).Inc()
		return
	}
	LendMePayBackCount.With(prometheus.Labels{"Status": "requested", "Reason": "", "Description": Source}).Inc()
	var lendLog Lendme_log
	lendLog.Source = Source
	lendLog.MSISDN = MSISDN
	lendLog.Log_Date = time.Now()
	lendLog.Type = "Lendme PayBack"
	// lendLog.Lendme_Amount = Amount
	// lendLog.Lendme_Fee = (Amount * Configuration.Service_FeePerc)

	if RechargeAmount <= 0 {
		error_msg := "recharge amount must be positive"
		err = errors.New(error_msg)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMePayBackCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description": ""}).Inc()
		//LendMePayBackAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description" : ""}).Add(RechargeAmount)
		return err
	}

	// check if subscriber exist, return error if not existing
	subscriber_na, exist := Map_Subscribers.CheckThenGet(MSISDN)
	if !exist {
		error_msg := "subscriber does not exist in the service pool"
		err = errors.New(error_msg)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMePayBackCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description": ""}).Inc()
		//LendMePayBackAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description" : ""}).Add(RechargeAmount)
		return err
	}
	subscriber, ok := subscriber_na.(Subscriber)
	if !ok {
		error_msg := "error in subscriber type assertion"
		err = errors.New(error_msg)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMePayBackCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description": ""}).Inc()
		//LendMePayBackAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description" : ""}).Add(RechargeAmount)
		return err
	}
	lendLog.Lendme_Outstanding_Amount = subscriber.Lendme_Outstanding_Amount
	lendLog.Lendme_Outstanding_Fee = subscriber.Lendme_Outstanding_Fee
	//check if subscriber have outstanding amount
	Outstanding_Amount := subscriber.Lendme_Outstanding_Amount + subscriber.Lendme_Outstanding_Fee
	if Outstanding_Amount <= 0 {
		return nil
	}
	IN_MSISDN := MSISDN
	if Configuration.Operation == "Gambia" {
		if len(MSISDN) > 7 {
			IN_MSISDN = IN_MSISDN[len(MSISDN)-7 : len(MSISDN)]
		}
	} else if Configuration.Operation == "SierraLeone" { //077928014
		if len(MSISDN) > 8 {
			IN_MSISDN = "0" + IN_MSISDN[len(MSISDN)-8:len(MSISDN)]
		}
	}
	//check if exists on IN and get details
	IN_Response, err := Uc.IN.INClient.GetAccountDetails("", "", IN_MSISDN)
	if err != nil {
		error_msg := "error getting account detail: " + err.Error()
		err = errors.New(error_msg)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMePayBackCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description": ""}).Inc()
		//LendMePayBackAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description" : ""}).Add(RechargeAmount)
		return err
	}
	lendLog.Subscriber_OpeningBlance = IN_Response.Balance
	//check balance min and max
	if IN_Response.Balance <= 0 {
		error_msg := "balance must be positive"
		err = errors.New(error_msg)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMePayBackCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description": ""}).Inc()
		//LendMePayBackAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description" : ""}).Add(RechargeAmount)
		return err
	}
	//calculate debit fee/amount
	var DebitfeeAmount float64
	var DebitAmount float64
	if IN_Response.Balance >= Outstanding_Amount {
		DebitfeeAmount = subscriber.Lendme_Outstanding_Fee
		DebitAmount = subscriber.Lendme_Outstanding_Amount
	} else {
		if IN_Response.Balance > subscriber.Lendme_Outstanding_Fee {
			DebitfeeAmount = subscriber.Lendme_Outstanding_Fee
			DebitAmount = IN_Response.Balance - subscriber.Lendme_Outstanding_Fee
		} else {
			DebitfeeAmount = IN_Response.Balance
		}
	}
	lendLog.Lendme_PayBack_Amount = DebitAmount
	lendLog.Lendme_PayBack_Fee = DebitfeeAmount
	//debit the fee amount
	if DebitfeeAmount > 0 {
		IN_MSISDN := MSISDN
		if Configuration.Operation == "Gambia" {
			if len(MSISDN) > 7 {
				IN_MSISDN = IN_MSISDN[len(MSISDN)-7 : len(MSISDN)]
			}
		} else if Configuration.Operation == "SierraLeone" { //077928014
			if len(MSISDN) > 8 {
				IN_MSISDN = "0" + IN_MSISDN[len(MSISDN)-8:len(MSISDN)]
			}
		}
		credit_response, credit_err := Uc.IN.INClient.SetAccountBalances("", "", IN_MSISDN, -1*DebitfeeAmount, "N", 0, 0, 0, 0, "LendmeFee")
		if credit_err != nil {
			error_msg := credit_err.Error()
			lendLog.Status = "failed"
			lendLog.StatusDescription = error_msg
			go Uc.Write_Lendme_log(lendLog)
			LendMePayBackCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description": "LendmeFee"}).Inc()
			LendMePayBackAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description": "LendmeFee"}).Add(DebitfeeAmount)
			return credit_err
		}
		if credit_response.ResultCode != "0" {
			error_msg := "result code is not 0"
			err = errors.New(error_msg)
			lendLog.Status = "failed"
			lendLog.StatusDescription = error_msg
			go Uc.Write_Lendme_log(lendLog)
			LendMePayBackCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description": "LendmeFee"}).Inc()
			LendMePayBackAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description": "LendmeFee"}).Add(DebitfeeAmount)
			return err
		}
		//subscriber.Cumulative_Payback_Amount
		subscriber.Lendme_Outstanding_Fee = subscriber.Lendme_Outstanding_Fee - DebitfeeAmount
		subscriber.Cumulative_Payback_Fee = subscriber.Cumulative_Payback_Fee + DebitfeeAmount
		subscriber.Last_Payback_Fee_Date = time.Now()
		Map_Subscribers.Put(subscriber.Key, subscriber)
		LendMePayBackCount.With(prometheus.Labels{"Status": "successful", "Reason": "", "Description": "LendmeFee"}).Inc()
		LendMePayBackAmount.With(prometheus.Labels{"Status": "successful", "Reason": "", "Description": "LendmeFee"}).Add(DebitfeeAmount)
	}
	//debit the fee amount
	if DebitAmount > 0 {
		IN_MSISDN := MSISDN
		if Configuration.Operation == "Gambia" {
			if len(MSISDN) > 7 {
				IN_MSISDN = IN_MSISDN[len(MSISDN)-7 : len(MSISDN)]
			}
		} else if Configuration.Operation == "SierraLeone" { //077928014
			if len(MSISDN) > 8 {
				IN_MSISDN = "0" + IN_MSISDN[len(MSISDN)-8:len(MSISDN)]
			}
		}
		credit_response, credit_err := Uc.IN.INClient.SetAccountBalances("", "", IN_MSISDN, -1*DebitAmount, "N", 0, 0, 0, 0, "LendmePayBack")
		if credit_err != nil {
			error_msg := credit_err.Error()
			lendLog.Status = "failed"
			lendLog.StatusDescription = error_msg
			go Uc.Write_Lendme_log(lendLog)
			LendMePayBackCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description": "LendmePayBack"}).Inc()
			LendMePayBackAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description": "LendmePayBack"}).Add(DebitAmount)
			return credit_err
		}
		if credit_response.ResultCode != "0" {
			error_msg := "result code is not 0"
			err = errors.New(error_msg)
			lendLog.Status = "failed"
			lendLog.StatusDescription = error_msg
			go Uc.Write_Lendme_log(lendLog)
			LendMePayBackCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description": "LendmePayBack"}).Inc()
			LendMePayBackAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description": "LendmePayBack"}).Add(DebitAmount)
			return err
		}
		subscriber.Lendme_Outstanding_Amount = subscriber.Lendme_Outstanding_Amount - DebitAmount
		subscriber.Cumulative_Payback_Amount = subscriber.Cumulative_Payback_Amount + DebitAmount
		subscriber.Last_Payback_Date = time.Now()
		Map_Subscribers.Put(subscriber.Key, subscriber)
		LendMePayBackCount.With(prometheus.Labels{"Status": "successful", "Reason": "", "Description": "LendmePayBack"}).Inc()
		LendMePayBackAmount.With(prometheus.Labels{"Status": "successful", "Reason": "", "Description": "LendmePayBack"}).Add(DebitAmount)
	}
	//Send notification SMS
	PaidAmount_round := Round((DebitAmount + DebitfeeAmount), 1, 2)
	PaidAmount_str := strconv.FormatFloat(PaidAmount_round, 'f', -1, 64)

	if Configuration.Operation == "Gambia" {
		SMS_MSISDN := subscriber.Key
		if len(SMS_MSISDN) >= 7 {
			SMS_MSISDN = "220" + SMS_MSISDN[len(SMS_MSISDN)-7:]
		}
		SMSText := "Dear subscriber, thank you for paying outstanding Lebal Ma amount " + PaidAmount_str + " GMD"
		go SendSMS("Africell", SMS_MSISDN, SMSText)
	} else if Configuration.Operation == "DRC" {
		//SMSText := "Cher abonne, merci d'avoir paye le montant emprunte par Lendme de " + PaidAmount_str + "u"
		SMSText := "Cher abonne, merci d'avoir paye le montant emprunte par Le service pretez moi de" + PaidAmount_str + "u"
		// SMSText := "Felicitations ! Votre compte a ete credite de " + fmt.Sprint(Round((DebitAmount), 1, 2)) + " unites et " + fmt.Sprint(Round((DebitfeeAmount), 1, 2)) + " unites des frais. Pour verifier votre solde, tapez *1099#"
		go SendSMS("Africell", subscriber.Key, SMSText)
	} else if Configuration.Operation == "SierraLeone" {
		SMS_MSISDN := subscriber.Key
		if len(SMS_MSISDN) >= 8 {
			SMS_MSISDN = "232" + SMS_MSISDN[len(SMS_MSISDN)-8:]
		}
		SMSText := "Dear subscriber, thank you for paying outstanding TrossMi amount NLE " + PaidAmount_str
		go SendSMS("TrossMi", SMS_MSISDN, SMSText)

	}

	//save log into DB
	lendLog.Status = "successful"
	lendLog.StatusDescription = ""
	go Uc.Write_Lendme_log(lendLog)
	return nil
}

func (Uc *UserControl) LendmeAO_PayBack(Source, MSISDN string, RechargeAmount float64, Opid string) (err error) {
	if Opid == Configuration.Lendme_EVC_Dealer_MSISDN {
		LendMePayBackCount.With(prometheus.Labels{"Status": "Ignored", "Reason": "Opid is lendme", "Description": Source}).Inc()
		return
	}
	LendMePayBackCount.With(prometheus.Labels{"Status": "requested", "Reason": "", "Description": Source}).Inc()
	var lendLog Lendme_log
	lendLog.Source = Source
	lendLog.MSISDN = MSISDN
	lendLog.Log_Date = time.Now()
	lendLog.Type = "Lendme PayBack"
	// lendLog.Lendme_Amount = Amount
	// lendLog.Lendme_Fee = (Amount * Configuration.Service_FeePerc)

	if RechargeAmount <= 0 {
		error_msg := "recharge amount must be positive"
		err = errors.New(error_msg)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMePayBackCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description": ""}).Inc()
		//LendMePayBackAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description" : ""}).Add(RechargeAmount)
		return err
	}

	// check if subscriber exist, return error if not existing
	subscriber_na, exist := Map_Subscribers.CheckThenGet(MSISDN)
	if !exist {
		error_msg := "subscriber does not exist in the service pool"
		err = errors.New(error_msg)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMePayBackCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description": ""}).Inc()
		//LendMePayBackAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description" : ""}).Add(RechargeAmount)
		return err
	}
	subscriber, ok := subscriber_na.(Subscriber)
	if !ok {
		error_msg := "error in subscriber type assertion"
		err = errors.New(error_msg)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMePayBackCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description": ""}).Inc()
		//LendMePayBackAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description" : ""}).Add(RechargeAmount)
		return err
	}
	lendLog.Lendme_Outstanding_Amount = subscriber.Lendme_Outstanding_Amount
	lendLog.Lendme_Outstanding_Fee = subscriber.Lendme_Outstanding_Fee
	//check if subscriber have outstanding amount
	Outstanding_Amount := subscriber.Lendme_Outstanding_Amount + subscriber.Lendme_Outstanding_Fee
	if Outstanding_Amount <= 0 {
		return nil
	}
	//check if exists on IN and get details
	OpeningBalance_RBB, OpeningBalance_EVC, balOpenErr := Uc.GetCustomerMainBalance(MSISDN)
	if balOpenErr != nil {
		lendLog.Status = "failed"
		lendLog.StatusDescription = balOpenErr.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMePayBackCount.With(prometheus.Labels{"Status": "failed", "Reason": balOpenErr.Error(), "Description": ""}).Inc()
		//LendMePayBackAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description" : ""}).Add(RechargeAmount)
		return balOpenErr
	}

	lendLog.Subscriber_OpeningBlance = OpeningBalance_RBB + OpeningBalance_EVC
	//check balance min and max
	if lendLog.Subscriber_OpeningBlance <= 0 {
		error_msg := "balance must be positive"
		err = errors.New(error_msg)
		lendLog.Status = "failed"
		lendLog.StatusDescription = err.Error()
		go Uc.Write_Lendme_log(lendLog)
		LendMePayBackCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description": ""}).Inc()
		//LendMePayBackAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description" : ""}).Add(RechargeAmount)
		return err
	}
	//calculate debit fee/amount
	//change here
	var DebitfeeAmount float64
	var DebitAmount float64
	if lendLog.Subscriber_OpeningBlance >= Outstanding_Amount {
		DebitfeeAmount = subscriber.Lendme_Outstanding_Fee
		DebitAmount = subscriber.Lendme_Outstanding_Amount
	} else {
		if lendLog.Subscriber_OpeningBlance > subscriber.Lendme_Outstanding_Fee {
			DebitfeeAmount = subscriber.Lendme_Outstanding_Fee
			DebitAmount = lendLog.Subscriber_OpeningBlance - subscriber.Lendme_Outstanding_Fee
		} else {
			DebitfeeAmount = lendLog.Subscriber_OpeningBlance
		}
	}
	lendLog.Lendme_PayBack_Amount = DebitAmount
	lendLog.Lendme_PayBack_Fee = DebitfeeAmount
	//debit the fee amount
	if DebitfeeAmount > 0 {
		var debitfee_request apgw.DebitSubscriber_V3_request
		debitfee_request.Program_Name = "LendmeFee"
		debitfee_request.MSISDN = MSISDN
		debitfee_request.Amount = DebitfeeAmount
		debitfee_response, err := Uc.APGW.APGWClient.DebitSubscriber(debitfee_request)
		if err != nil {
			error_msg := err.Error()
			lendLog.Status = "failed"
			lendLog.StatusDescription = error_msg
			go Uc.Write_Lendme_log(lendLog)
			LendMePayBackCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description": "LendmeFee"}).Inc()
			LendMePayBackAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description": "LendmeFee"}).Add(DebitfeeAmount)
			return err
		}
		if debitfee_response.StatusCode != http.StatusOK {
			error_msg := debitfee_response.ErrorDescription
			err = errors.New(error_msg)
			lendLog.Status = "failed"
			lendLog.StatusDescription = error_msg
			go Uc.Write_Lendme_log(lendLog)
			LendMePayBackCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description": "LendmeFee"}).Inc()
			LendMePayBackAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description": "LendmeFee"}).Add(DebitfeeAmount)
			return err
		}
		//subscriber.Cumulative_Payback_Amount
		subscriber.Lendme_Outstanding_Fee = subscriber.Lendme_Outstanding_Fee - DebitfeeAmount
		subscriber.Cumulative_Payback_Fee = subscriber.Cumulative_Payback_Fee + DebitfeeAmount
		subscriber.Last_Payback_Fee_Date = time.Now()
		Map_Subscribers.Put(subscriber.Key, subscriber)
		LendMePayBackCount.With(prometheus.Labels{"Status": "successful", "Reason": "", "Description": "LendmeFee"}).Inc()
		LendMePayBackAmount.With(prometheus.Labels{"Status": "successful", "Reason": "", "Description": "LendmeFee"}).Add(DebitfeeAmount)
	}
	//debit the amount
	if DebitAmount > 0 {
		var debitAmount_request apgw.DebitSubscriber_V3_request
		debitAmount_request.Program_Name = "LendmePayBack"
		debitAmount_request.MSISDN = MSISDN
		debitAmount_request.Amount = DebitAmount
		debitAmount_response, err := Uc.APGW.APGWClient.DebitSubscriber(debitAmount_request)
		if err != nil {
			error_msg := err.Error()
			lendLog.Status = "failed"
			lendLog.StatusDescription = error_msg
			go Uc.Write_Lendme_log(lendLog)
			LendMePayBackCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description": "LendmePayBack"}).Inc()
			LendMePayBackAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description": "LendmePayBack"}).Add(DebitAmount)
			return err
		}
		if debitAmount_response.StatusCode != http.StatusOK {
			error_msg := debitAmount_response.ErrorDescription
			err = errors.New(error_msg)
			lendLog.Status = "failed"
			lendLog.StatusDescription = error_msg
			go Uc.Write_Lendme_log(lendLog)
			LendMePayBackCount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description": "LendmePayBack"}).Inc()
			LendMePayBackAmount.With(prometheus.Labels{"Status": "failed", "Reason": error_msg, "Description": "LendmePayBack"}).Add(DebitAmount)
			return err
		}
		subscriber.Lendme_Outstanding_Amount = subscriber.Lendme_Outstanding_Amount - DebitAmount
		subscriber.Cumulative_Payback_Amount = subscriber.Cumulative_Payback_Amount + DebitAmount
		subscriber.Last_Payback_Date = time.Now()
		Map_Subscribers.Put(subscriber.Key, subscriber)
		LendMePayBackCount.With(prometheus.Labels{"Status": "successful", "Reason": "", "Description": "LendmePayBack"}).Inc()
		LendMePayBackAmount.With(prometheus.Labels{"Status": "successful", "Reason": "", "Description": "LendmePayBack"}).Add(DebitAmount)
	}
	//Send notification SMS
	PaidAmount_round := Round((DebitAmount + DebitfeeAmount), 1, 2)
	PaidAmount_str := strconv.FormatFloat(PaidAmount_round, 'f', -1, 64)

	SMSText := "Dear subscriber, thank you for paying outstanding SERVICE_NAME amount " + PaidAmount_str + " Kz"
	go SendSMS("Africell", subscriber.Key, SMSText)

	//save log into DB
	lendLog.Status = "successful"
	lendLog.StatusDescription = ""
	go Uc.Write_Lendme_log(lendLog)
	return nil
}

// ///////////////////////////////////////////////////////
// get subscriber from USSD
// ///////////////////////////////////////////////////////

func (Uc *UserControl) SubscriberUSSD_Get(Key string) (response Subscriber_USSD, err error) {
	if Key == "" {
		err = errors.New("subscriber MSISDN is not provided")
		return response, err
	} else {
		subscriber_na, exits := Map_Subscribers.CheckThenGet(Key)
		if !exits {
			err = errors.New("subscriber does not exist")
			return response, err
		}
		subscriber, ok := subscriber_na.(Subscriber)
		if !ok {
			err = errors.New("error in subcriber type assertion")
			return response, err
		}

		response.MSISDN = Key
		response.IsLendmeEligible = subscriber.IsLendmeEligible
		response.Credit_Limit_Scheme = subscriber.Credit_Limit_Scheme
		response.NotElligibleReason = subscriber.NotElligibleReason
		response.Lendme_Outstanding_Amount = RoundToNearest(subscriber.Lendme_Outstanding_Amount, 2)
		response.Lendme_Outstanding_Fee = RoundToNearest(subscriber.Lendme_Outstanding_Fee, 2)
		//get schema detail for elligble subscriber
		if subscriber.IsLendmeEligible {
			scheme_na, exits := Map_Credit_Limit_Scheme.CheckThenGet(subscriber.Credit_Limit_Scheme)
			if !exits {
				err = errors.New("credit limit scheme does not exist")
				return response, err
			}
			scheme, ok := scheme_na.(Credit_Limit_Scheme)
			if !ok {
				err = errors.New("error in credit limit scheme type assertion")
				return response, err
			}
			response.Credit_limit_Amount = scheme.Credit_limit_Amount
			remaining_Allowed := scheme.Credit_limit_Amount - subscriber.Lendme_Outstanding_Amount
			if subscriber.Lendme_Outstanding_Amount > 0 {
				if remaining_Allowed > 0 {
					response.Min_Allowed_Amount = 1
					response.Max_Allowed_Amount = RoundDown(remaining_Allowed, 0)
				} else {
					response.Min_Allowed_Amount = 0
					response.Max_Allowed_Amount = 0
				}
			} else {
				if remaining_Allowed >= Configuration.Min_Allowed_Amnt {
					response.Min_Allowed_Amount = Configuration.Min_Allowed_Amnt
					response.Max_Allowed_Amount = RoundDown(remaining_Allowed, 0)
				} else {
					response.Min_Allowed_Amount = 0
					response.Max_Allowed_Amount = 0
				}
			}
		}
		return response, nil
	}
}

// ***********************************************************************
// Customer Exclusion functions
// ***********************************************************************
func (Uc *UserControl) Lendme_Customer_Exclusion_Add(Login string, request Customer_Exclusion_AddRequest) (Id int64, err error) {
	request.Key = Normalize_International_MSISDN(request.Key)
	//check key if filled and if already used
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	//check if key already used
	exits := Map_Lendme_Customer_Exclusion.Check(request.Key)
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
	subscriber_na, exits := Map_Subscribers.CheckThenGet(NewEntry.Key)
	if !exits {
		err = errors.New("msisdn is not created")
		return Id, err
	}
	subscriber, ok := subscriber_na.(Subscriber)
	if !ok {
		return Id, errors.New("error in subscriber type assertion")
	}
	//add to cache and DB
	Map_Lendme_Customer_Exclusion.Put(NewEntry.Key, NewEntry)
	subscriber.IsLendmeEligible = false
	Map_Subscribers.Put(subscriber.Key, subscriber)
	return Id, nil
}

func (Uc *UserControl) Lendme_Customer_Exclusion_Edit(Login string, request Customer_Exclusion_EditRequest) (Id int64, err error) {
	//check and validate
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	entry_na, exits := Map_Lendme_Customer_Exclusion.CheckThenGet(request.Key)
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

	//Prepare new entry
	entry.Key = request.Key
	entry.AddReason = request.AddReason

	if request.NewKey != "" {
		if request.NewKey != request.Key {
			//delete old
			Map_Lendme_Customer_Exclusion.Delete(request.Key)
			//update key
			entry.Key = request.NewKey
		}
	}
	//add to cache and DB
	subscriber_na, exits := Map_Subscribers.CheckThenGet(entry.Key)
	if !exits {
		err = errors.New("msisdn is not created")
		return Id, err
	}
	subscriber, ok := subscriber_na.(Subscriber)
	if !ok {
		return Id, errors.New("error in subscriber type assertion")
	}
	Map_Lendme_Customer_Exclusion.Put(entry.Key, entry)
	subscriber.IsLendmeEligible = false
	Map_Subscribers.Put(subscriber.Key, subscriber)
	return Id, nil
}

func (Uc *UserControl) Lendme_Customer_Exclusion_Get(Key string) (entries []Customer_Exclusion, err error) {
	if Key == "" {
		entries_na := Map_Lendme_Customer_Exclusion.ConvertToArray()
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
		entry_na, exits := Map_Lendme_Customer_Exclusion.CheckThenGet(Key)
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

func (Uc *UserControl) Lendme_Customer_Exclusion_GetPaginated(Page, Limit int64) (entries []Customer_Exclusion, err error) {
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
	findResult, err := DAO_Lendme_Customer_Exclusion.FindPaginate(findparams, paginationparams)
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

func (Uc *UserControl) Lendme_Customer_Exclusion_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	_, exits := Map_Lendme_Customer_Exclusion.CheckThenGet(Key)
	if !exits {
		err = errors.New("entry does not exist")
		return err
	}
	Map_Lendme_Customer_Exclusion.Delete(Key)

	return nil
}

// ***********************************************************************
// Customer Exclusion functions
// ***********************************************************************
func (Uc *UserControl) Lendme_Customer_COS_Exclusion_Add(Login string, request Customer_COS_Exclusion_AddRequest) (Id int64, err error) {
	//check key if filled and if already used
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	//check if key already used
	exits := Map_Lendme_Customer_COS_Exclusion.Check(request.Key)
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
	Map_Lendme_Customer_COS_Exclusion.Put(NewEntry.Key, NewEntry)
	var findparams daoc.DAOFindParams
	var array []daoc.DAOFindCriteria
	//restrict access for records that belong to this user
	var criteria daoc.DAOFindCriteria = daoc.DAOFindCriteria{
		Field:    "COS",
		Value:    NewEntry.Key,
		Operator: "EQUAL",
	}
	array = append(array, criteria)

	findparams.FindCriteria = array
	//
	findResult, err := DAO_Subscribers.Find(findparams)
	if err != nil {
		return Id, err
	}

	if len(findResult) > 0 {
		for _, entry_na := range findResult {
			subscriber := reflect.ValueOf(entry_na).Elem().Interface().(Subscriber)
			subscriber.IsLendmeEligible = false
			Map_Subscribers.Put(subscriber.Key, subscriber)
		}
	}
	return Id, nil
}

func (Uc *UserControl) Lendme_Customer_COS_Exclusion_Edit(Login string, request Customer_COS_Exclusion_EditRequest) (Id int64, err error) {
	//check and validate
	if request.Key == "" {
		err = errors.New("key cannot be empty")
		return Id, err
	}
	entry_na, exits := Map_Lendme_Customer_COS_Exclusion.CheckThenGet(request.Key)
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

	//Prepare new entry
	entry.Key = request.Key
	entry.AddReason = request.AddReason

	if request.NewKey != "" {
		if request.NewKey != request.Key {
			//delete old
			Map_Lendme_Customer_COS_Exclusion.Delete(request.Key)
			//update key
			entry.Key = request.NewKey
		}
	}
	//add to cache and DB
	Map_Lendme_Customer_COS_Exclusion.Put(entry.Key, entry)
	var findparams daoc.DAOFindParams
	var array []daoc.DAOFindCriteria
	//restrict access for records that belong to this user
	var criteria daoc.DAOFindCriteria = daoc.DAOFindCriteria{
		Field:    "COS",
		Value:    entry.Key,
		Operator: "EQUAL",
	}
	array = append(array, criteria)

	findparams.FindCriteria = array
	//
	findResult, err := DAO_Subscribers.Find(findparams)
	if err != nil {
		return Id, err
	}

	if len(findResult) > 0 {
		for _, entry_na := range findResult {
			subscriber := reflect.ValueOf(entry_na).Elem().Interface().(Subscriber)
			subscriber.IsLendmeEligible = false
			Map_Subscribers.Put(subscriber.Key, subscriber)
		}
	}
	return Id, nil
}

func (Uc *UserControl) Lendme_Customer_COS_Exclusion_Get(Key string) (entries []Customer_COS_Exclusion, err error) {
	if Key == "" {
		entries_na := Map_Lendme_Customer_COS_Exclusion.ConvertToArray()
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
		entry_na, exits := Map_Lendme_Customer_COS_Exclusion.CheckThenGet(Key)
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

func (Uc *UserControl) Lendme_Customer_COS_Exclusion_GetPaginated(Page, Limit int64) (entries []Customer_COS_Exclusion, err error) {
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
	findResult, err := DAO_Lendme_Customer_COS_Exclusion.FindPaginate(findparams, paginationparams)
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

func (Uc *UserControl) Lendme_Customer_COS_Exclusion_Delete(Login, Key string) (err error) {
	if Key == "" {
		err = errors.New("key cannot be empty")
		return err
	}
	_, exits := Map_Lendme_Customer_COS_Exclusion.CheckThenGet(Key)
	if !exits {
		err = errors.New("entry does not exist")
		return err
	}

	Map_Lendme_Customer_COS_Exclusion.Delete(Key)
	return nil
}

func (Uc *UserControl) GetOutstandingSummary() (err error) {
	pipe := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "Lendme_Outstanding_Amount", Value: bson.D{
					{Key: "$gt", Value: 0},
				}},
			}},
		},
		bson.D{
			{Key: "$group", Value: bson.D{
				// {Key: "_id", Value: bson.D{
				// 	{Key: "CreateTime", Value: bson.D{
				// 		{Key: "$dateToString", Value: bson.D{
				// 			{Key: "format", Value: "%Y-%m-%d"},
				// 			{Key: "date", Value: "$CreateTime"},
				// 		}},
				// 	}},
				// 	{Key: "Type", Value: "$Type"},
				// }},

				{Key: "_id", Value: nil},

				{Key: "TotalSubsWithOutstanding", Value: bson.D{
					{Key: "$sum", Value: 1},
				}},
				{Key: "TotalOutstandingAmount", Value: bson.D{
					{Key: "$sum", Value: "$Lendme_Outstanding_Amount"},
				}},
				{Key: "TotalOutstandingFee", Value: bson.D{
					{Key: "$sum", Value: "$Lendme_Outstanding_Fee"},
				}},
			}},
		},
		// bson.D{
		// 	{Key: "$project", Value: bson.D{
		// 		{Key: "CreateTime", Value: "$_id.CreateTime"},
		// 		{Key: "Type", Value: "$_id.Type"},
		// 		{Key: "Total", Value: "$total"},
		// 		{Key: "_id", Value: 0},
		// 	}},
		// },
		// bson.D{
		// 	{Key: "$sort", Value: bson.D{
		// 		{Key: "CreateTime", Value: 1},
		// 	}},
		// },
	}

	cur, err := Uc.MongoDB.MongoDBClient.Database(Configuration.DB_Name).Collection(DAO_Subscribers.Collection).Aggregate(context.TODO(), pipe, options.Aggregate())

	if err != nil {
		log.Println("Error in GetOutstandingSummary: ", err)
		return
	}
	defer cur.Close(context.Background())
	//var output []AlarmsDailyByType
	for cur.Next(context.Background()) {
		var entry Lendme_Outstanding_Summary
		err := cur.Decode(&entry)
		if err != nil {
			log.Println("Error in GetOutstandingSummary: ", err)
			return err
		}
		//log.Println(entry)
		LendMeOutstandingSummary.Reset()
		LendMeOutstandingSummary.With(prometheus.Labels{"Field": "TotalOutstandingAmount"}).Add(entry.TotalOutstandingAmount)
		LendMeOutstandingSummary.With(prometheus.Labels{"Field": "TotalOutstandingFee"}).Add(entry.TotalOutstandingFee)
		LendMeOutstandingSummary.With(prometheus.Labels{"Field": "TotalSubsWithOutstanding"}).Add(entry.TotalSubsWithOutstanding)
	}
	return
}

func (Uc *UserControl) Auto_GetOutstandingSummary() {
	Uc.GetOutstandingSummary()
	for range time.Tick(time.Second * 300) {
		Uc.GetOutstandingSummary()
	}
}

func (Uc *UserControl) Lendme_Subscriber_Daily_Snapshot() {
	exec := 0
	LOG_ID := "<<Lendme Subscriber Daily Snapshot>>"
	for range time.Tick(time.Second * 1) {
		_CurrentDateTime := time.Now()
		_hr, _mi, _se := _CurrentDateTime.Clock()
		if _hr == 00 {
			if _mi == 00 {
				if _se < 60 {
					if exec == 0 {
						exec = 1
						log.Println(LOG_ID + " triggered")
						yesterday := time.Now().AddDate(0, 0, -1)
						YYYY, MM, _, DD, _, _, _ := GetTimeParts(yesterday)
						Db := DAO_Lendme_log.DB + "_" + YYYY + MM
						Col := DAO_Subscribers.Collection + "_" + DD
						err := DAO_Subscribers.CollectionSnapshot(Db, Col)
						if err != nil {
							log.Println("error while taking a snapshot from subscriber collection", err)
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

// /////////////////////////////////////////////////////////////////////////////////////////////////////
// /////SEND SMS////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////
func SendSMS(Sender string, target string, SMSText string) (_rErr error) {
	log.Println("Sending SMS: Sender (" + Sender + "), Target (" + target + "), text (" + SMSText + ") ")
	url := "http://" + Configuration.SMPP.IP + ":" + Configuration.SMPP.Port + "/?systemid=" + Configuration.SMPP.Login + "&password=" + url.QueryEscape(Configuration.SMPP.Password) + "&Originator=" + Sender + "&dest_addr=" + target + "&msg_text=" + url.QueryEscape(SMSText) + "&encoding=1&ston=5&snpi=0&dton=1&registered_delivery=0"
	//-------------- Encoding used in DRC and GM Start
	//"&ston=5&snpi=0&dton=1&dnpi=1&encoding=1"
	//-------------- Encoding used in DRC and GM End
	method := "POST"
	req, err := http.NewRequest(method, url, nil)
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
	if resp.StatusCode == 200 {
		return nil
	} else {
		body, _ := io.ReadAll(resp.Body)
		log.Println("Error sending SMS: ", string(body))
		err := errors.New("error sending SMS")
		return err
	}
}
