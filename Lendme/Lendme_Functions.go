package Lendme

import (
	"afr_auth_center/AuthCenter"
	"daoc"
	"errors"
	"log"
	"reflect"
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
	if !exits {
		err = errors.New("scheme already exist")
		return Id, err
	}

	//Prepare new entry
	var NewEntry Credit_Limit_Scheme
	NewEntry.Scheme_Id = Map_AutoIncrement.GetNextAI("Credit_Limit_Scheme-Id")
	Id = NewEntry.Scheme_Id
	NewEntry.Key = request.Key
	NewEntry.ARPU_From = request.ARPU_From
	NewEntry.ARPU_Till = request.ARPU_Till
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
	scheme.ARPU_From = request.ARPU_From
	scheme.ARPU_Till = request.ARPU_Till
	scheme.AON_From = request.AON_From
	scheme.AON_Till = request.AON_Till
	scheme.Credit_limit_Amount = request.Credit_limit_Amount
	log := Event_Log{
		Event_User:         Login,
		Event_Time:         time.Now(),
		Event_Type:         "edit",
		Event_Description:  "edit entry",
		Event_Entry_Before: scheme,
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
		if len(schemes) > 0 {
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
