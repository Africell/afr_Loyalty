package Lendme

import (
	"afr_auth_center/AuthCenter"
	"daoc"
	"log"
	"reflect"
)

var Map_DefaultValues daoc.Cache_Synch
var DAO_DefaultValues daoc.DAO

var MapAccessEntry daoc.Cache_Synch
var DAO_AccessEntry daoc.DAO

var Map_AutoIncrement daoc.Cache_Synch
var DAO_AutoIncrement daoc.DAO

var Map_Subscribers daoc.Cache_Synch
var DAO_Subscribers daoc.DAO

func (uc *UserControl) InitializeCache() {
	var access_entry AuthCenter.AccessEntry
	MapAccessEntry.Initialize("AccessEntry", "AccessEntry", reflect.TypeOf(AuthCenter.AccessEntry{}), access_entry, true, &DAO_AccessEntry, uc.CacheDir.List)
	var AutoIncr daoc.AutoIncrement
	Map_AutoIncrement.Initialize("AutoIncrement", "AutoIncrement", reflect.TypeOf(daoc.AutoIncrement{}), AutoIncr, true, &DAO_AutoIncrement, uc.CacheDir.List)
	var defaultValues DefaultValues
	Map_DefaultValues.Initialize("DefaultValues", "DefaultValues", reflect.TypeOf(DefaultValues{}), defaultValues, true, &DAO_DefaultValues, uc.CacheDir.List)
}

func (uc *UserControl) InitializeDAO() {
	DAO_AccessEntry.Initialize("AccessEntry", uc.MongoDB.MongoDBClient, reflect.TypeOf(AuthCenter.AccessEntry{}), Configuration.DB_Name, "Col_AccessEntry", "")
	DAO_AutoIncrement.Initialize("AutoIncrement", uc.MongoDB.MongoDBClient, reflect.TypeOf(daoc.AutoIncrement{}), Configuration.DB_Name, "Col_AutoIncrement", "")
	//DAO_UAT_Pool.Initialize("UAT_Pool", uc.MongoDB.MongoDBClient, reflect.TypeOf(UAT_Pool{}), Configuration.DB_Name, "Col_UAT_Pool", "")
	DAO_DefaultValues.Initialize("DefaultValues", uc.MongoDB.MongoDBClient, reflect.TypeOf(DefaultValues{}), Configuration.DB_Name, "Col_DefaultValues", "")
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

	log.Println("DB index manintenance process finished")
}

func (uc *UserControl) LoadDefaultValues() {
	exist := Map_DefaultValues.Check("Default")
	if !exist {
		defaultValues := DefaultValues{
			Key:              "Default",
			MSL_SIM_Card:     10,
			MSL_Scratch_Card: 10000,
			MSL_EVoucher:     5000,
			MSL_AfriMoney:    5000,
		}
		Map_DefaultValues.Put(defaultValues.Key, defaultValues)
	}
}
