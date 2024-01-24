package Lendme

import (
	//"encoding/xml"
	"time"
)

type API_Standard_response struct {
	//response source detail
	SourceIP        string    `bson:"SourceIP" json:"-"`
	Login           string    `bson:"Login" json:"-"`
	SourceApp       string    `bson:"SourceApp" json:"-"`
	AccessKey       string    `bson:"AccessKey" json:"-"`
	AccessMethod    string    `bson:"AccessMethod" json:"-"`
	HostId          string    `bson:"HostId" json:"-"`
	ReceiveDate     time.Time `bson:"ReceiveDate" json:"-"`
	TransactionType string    `bson:"TransactionType" json:"-"`
	//response detail
	Data interface{}
	//response result
	Status            string    `bson:"Status" json:"Status"` //successful, failed
	StatusCode        int       `bson:"StatusCode" json:"StatusCode"`
	StatusDescription string    `bson:"StatusDescription" json:"StatusDescription"` //error description if there is an error
	ErrorDescription  string    `bson:"ErrorDescription" json:"ErrorDescription"`
	StatusDate        time.Time `bson:"StatusDate" json:"-"`
	Elapsedtime       int64     `bson:"Elapsedtime" json:"-"`
}

type Event_Log struct {
	Event_User         string      `bson:"Event_User" json:"Event_User"`
	Event_Time         time.Time   `bson:"Event_Time" json:"Event_Time"`
	Event_Type         string      `bson:"Event_Type" json:"Event_Type"`
	Event_Description  string      `bson:"Event_Description" json:"Event_Description"`
	Event_Entry_Before interface{} `bson:"Event_Entry_Before" json:"Event_Entry_Before"`
	Event_Entry_After  interface{} `bson:"Event_Entry_After" json:"Event_Entry_After"`
}

type Agent struct {
	Key      string `bson:"Key" json:"Key"` //Login
	Agent_Id int64  `bson:"Agent_Id" json:"Agent_Id"`
	//Personal Info
	FirstName  string    `bson:"FirstName" json:"FirstName"`
	MiddleName string    `bson:"MiddleName" json:"MiddleName"`
	LastName   string    `bson:"LastName" json:"LastName"`
	Gender     string    `bson:"Gender" json:"Gender"`
	BirthDate  time.Time `bson:"BirthDate" json:"BirthDate"`
	Language   string    `bson:"Language" json:"Language"`
	//Address
	Provinces      string `bson:"Provinces" json:"Provinces"`
	Municipality   string `bson:"Municipality" json:"Municipality"`
	Street_Address string `bson:"Street_Address" json:"Street_Address"`
	Address_Note   string `bson:"Address_Note" json:"Address_Note"`

	Nationality                  string `bson:"Nationality" json:"Nationality"`
	Fiscal_identification_number string `bson:"Fiscal_identification_number" json:"Fiscal_identification_number"`
	FaceFrontPhotoId             string `bson:"FaceFrontPhotoId" json:"FaceFrontPhotoId"` //Attachement ID
	FaceSidePhotoId              string `bson:"FaceSidePhotoId" json:"FaceSidePhotoId"`   //Attachement ID
	OfficialDocType              string `bson:"OfficialDocType" json:"OfficialDocType"`
	OfficialDocSerialNumber      string `bson:"OfficialDocSerialNumber" json:"OfficialDocSerialNumber"`
	OfficialDocDocFrontPhotoId   string `bson:"OfficialDocDocFrontPhotoId" json:"OfficialDocDocFrontPhotoId"`
	OfficialDocDocBackPhotoId    string `bson:"OfficialDocDocBackPhotoId" json:"OfficialDocDocBackPhotoId"`
	//contacts
	Email            string `bson:"Email" json:"Email"`
	MobilePhone      string `bson:"MobilePhone" json:"MobilePhone"`
	OtherMobilePhone string `bson:"OtherMobilePhone" json:"OtherMobilePhone"`

	//company info
	Company    string `V:"M" bson:"Company" json:"Company"`
	Department string `V:"" bson:"Department" json:"Department"`
	Unit       string `V:"" bson:"Unit" json:"Unit"`

	Access_Level string `V:"" bson:"Access_Level" json:"Access_Level"` //Agent or Supervisor

	//Workflow status
	Status            string    `bson:"Status" json:"Status"`
	StatusDate        time.Time `bson:"StatusDate" json:"StatusDate"`
	StatusUser        string    `bson:"StatusUser" json:"StatusUser"`
	StatusDescription string    `bson:"StatusDescription" json:"StatusDescription"`
	//Events Log trail
	Event_Logs []Event_Log `bson:"Event_Logs" json:"Event_Logs"`
}

type Agent_Last_Location struct {
	Key            string    `bson:"Key" json:"Key"` //Agent_Key
	Location_Time  time.Time `bson:"Location_Time" json:"Location_Time"`
	GPS_Location   Location  `bson:"GPS_Location" json:"GPS_Location"`
	Detection_Type string    `bson:"Detection_Type" json:"Detection_Type"`
}

type Agent_History_Location struct {
	Key            string    `bson:"Key" json:"Key"` //Agent_Key
	Location_Time  time.Time `bson:"Location_Time" json:"Location_Time"`
	GPS_Location   Location  `bson:"GPS_Location" json:"GPS_Location"`
	Detection_Type string    `bson:"Detection_Type" json:"Detection_Type"`
}

type Agent_Location_Add_Request struct {
	Key            string  `bson:"Key" json:"Key"` //Agent_Key
	Latitude       float64 `bson:"Latitude" json:"Latitude"`
	Longitude      float64 `bson:"Longitude" json:"Longitude"`
	Detection_Type string  `bson:"Detection_Type" json:"Detection_Type"`
}

type Agent_Add_Request struct {
	Key      string `bson:"Key" json:"Key"` //Login
	Agent_Id int64  `bson:"Agent_Id" json:"Agent_Id"`
	//Personal Info
	FirstName  string    `bson:"FirstName" json:"FirstName"`
	MiddleName string    `bson:"MiddleName" json:"MiddleName"`
	LastName   string    `bson:"LastName" json:"LastName"`
	Gender     string    `bson:"Gender" json:"Gender"`
	BirthDate  time.Time `bson:"BirthDate" json:"BirthDate"`
	Language   string    `bson:"Language" json:"Language"`
	//Address
	Provinces      string `bson:"Provinces" json:"Provinces"`
	Municipality   string `bson:"Municipality" json:"Municipality"`
	Street_Address string `bson:"Street_Address" json:"Street_Address"`
	Address_Note   string `bson:"Address_Note" json:"Address_Note"`

	Nationality                  string `bson:"Nationality" json:"Nationality"`
	Fiscal_identification_number string `bson:"Fiscal_identification_number" json:"Fiscal_identification_number"`
	FaceFrontPhoto_b64           string `bson:"FaceFrontPhoto_b64" json:"FaceFrontPhoto_b64"` //Attachement ID
	FaceSidePhoto_b64            string `bson:"FaceSidePhoto_b64" json:"FaceSidePhoto_b64"`   //Attachement ID
	OfficialDocFrontPhoto_b64    string `bson:"OfficialDocFrontPhoto_b64" json:"OfficialDocFrontPhoto_b64"`
	OfficialDocBackPhoto_b64     string `bson:"OfficialDocBackPhoto_b64" json:"OfficialDocBackPhoto_b64"`
	OfficialDocType              string `bson:"OfficialDocType" json:"OfficialDocType"`
	OfficialDocSerialNumber      string `bson:"OfficialDocSerialNumber" json:"OfficialDocSerialNumber"`

	//contacts
	Email            string `bson:"Email" json:"Email"`
	MobilePhone      string `bson:"MobilePhone" json:"MobilePhone"`
	OtherMobilePhone string `bson:"OtherMobilePhone" json:"OtherMobilePhone"`

	//company info
	Company    string `V:"M" bson:"Company" json:"Company"`
	Department string `V:"" bson:"Department" json:"Department"`
	Unit       string `V:"" bson:"Unit" json:"Unit"`

	Access_Level string `V:"" bson:"Access_Level" json:"Access_Level"` //Agent or Supervisor
}

type Agent_Edit_Request struct {
	Key      string `bson:"Key" json:"Key"` //Login
	Agent_Id int64  `bson:"Agent_Id" json:"Agent_Id"`
	//NewKey   string `bson:"NewKey" json:"NewKey"` //New agent key
	//Personal Info
	FirstName  string    `bson:"FirstName" json:"FirstName"`
	MiddleName string    `bson:"MiddleName" json:"MiddleName"`
	LastName   string    `bson:"LastName" json:"LastName"`
	Gender     string    `bson:"Gender" json:"Gender"`
	BirthDate  time.Time `bson:"BirthDate" json:"BirthDate"`
	Language   string    `bson:"Language" json:"Language"`
	//Address
	Provinces      string `bson:"Provinces" json:"Provinces"`
	Municipality   string `bson:"Municipality" json:"Municipality"`
	Street_Address string `bson:"Street_Address" json:"Street_Address"`
	Address_Note   string `bson:"Address_Note" json:"Address_Note"`

	Nationality                  string `bson:"Nationality" json:"Nationality"`
	Fiscal_identification_number string `bson:"Fiscal_identification_number" json:"Fiscal_identification_number"`
	FaceFrontPhoto_b64           string `bson:"FaceFrontPhoto_b64" json:"FaceFrontPhoto_b64"` //Attachement ID
	FaceSidePhoto_b64            string `bson:"FaceSidePhoto_b64" json:"FaceSidePhoto_b64"`   //Attachement ID
	OfficialDocFrontPhoto_b64    string `bson:"OfficialDocFrontPhoto_b64" json:"OfficialDocFrontPhoto_b64"`
	OfficialDocBackPhoto_b64     string `bson:"OfficialDocBackPhoto_b64" json:"OfficialDocBackPhoto_b64"`
	OfficialDocType              string `bson:"OfficialDocType" json:"OfficialDocType"`
	OfficialDocSerialNumber      string `bson:"OfficialDocSerialNumber" json:"OfficialDocSerialNumber"`
	//contacts
	Email            string `bson:"Email" json:"Email"`
	MobilePhone      string `bson:"MobilePhone" json:"MobilePhone"`
	OtherMobilePhone string `bson:"OtherMobilePhone" json:"OtherMobilePhone"`

	//company info
	Company    string `bson:"Company" json:"Company"`
	Department string `bson:"Department" json:"Department"`
	Unit       string `bson:"Unit" json:"Unit"`

	Access_Level string `bson:"Access_Level" json:"Access_Level"` //Agent or Supervisor

}

type Agent_Image struct {
	Image_Id         string `bson:"Image_Id" json:"Image_Id"`
	Image_b64        string `bson:"Image_b64" json:"Image_b64"`
	ImageDescription string `bson:"ImageDescription" json:"ImageDescription"`
	UploadedBy       string `bson:"UploadedBy" json:"UploadedBy"`
	UploadTime       string `bson:"UploadTime" json:"UploadTime"`
	Agent_key        string `bson:"Agent_key" json:"Agent_key"`
}

type Dealer struct {
	Key       string `bson:"Key" json:"Key"` //Login
	Dealer_Id int64  `bson:"Dealer_Id" json:"Dealer_Id"`
	//Personal Info
	FirstName  string    `bson:"FirstName" json:"FirstName"`
	MiddleName string    `bson:"MiddleName" json:"MiddleName"`
	LastName   string    `bson:"LastName" json:"LastName"`
	Gender     string    `bson:"Gender" json:"Gender"`
	BirthDate  time.Time `bson:"BirthDate" json:"BirthDate"`
	Address    string    `bson:"Address" json:"Address"`
	Language   string    `bson:"Language" json:"Language"`

	Nationality                string `bson:"Nationality" json:"Nationality"`
	FaceFrontPhotoId           string `bson:"FaceFrontPhotoId" json:"FaceFrontPhotoId"` //Attachement ID
	FaceSidePhotoId            string `bson:"FaceSidePhotoId" json:"FaceSidePhotoId"`   //Attachement ID
	OfficialDocType            string `bson:"OfficialDocType" json:"OfficialDocType"`
	OfficialDocSerialNumber    string `bson:"OfficialDocSerialNumber" json:"OfficialDocSerialNumber"`
	OfficialDocDocFrontPhotoId string `bson:"OfficialDocDocFrontPhotoId" json:"OfficialDocDocFrontPhotoId"`
	OfficialDocDocBackPhotoId  string `bson:"OfficialDocDocBackPhotoId" json:"OfficialDocDocBackPhotoId"`
	//contacts
	Email            string `bson:"Email" json:"Email"`
	MobilePhone      string `bson:"MobilePhone" json:"MobilePhone"`
	OtherMobilePhone string `bson:"OtherMobilePhone" json:"OtherMobilePhone"`

	//company info
	Company      string `V:"M" bson:"Company" json:"Company"`
	ParentDealer string `V:"M" bson:"ParentDealer" json:"ParentDealer"`

	//Workflow status
	Status            string    `bson:"Status" json:"Status"`
	StatusDate        time.Time `bson:"StatusDate" json:"StatusDate"`
	StatusDescription string    `bson:"StatusDescription" json:"StatusDescription"`
	//Events Log trail
	Event_Logs []Event_Log `bson:"Event_Logs" json:"Event_Logs"`
}

type Outlet_Qrcode struct {
	Key            string    `bson:"Key" json:"Key"` //barcode string
	Qrcode_Id      int64     `bson:"Qrcode_Id" json:"Qrcode_Id"`
	GenerationDate time.Time `bson:"GenerationDate" json:"GenerationDate"`
	GeneratedBy    string    `bson:"GeneratedBy" json:"GeneratedBy"`
	BatchNo        string    `bson:"BatchNo" json:"BatchNo"`

	QrCode_Pict_Id string `bson:"QrCode_Pict_Id" json:"QrCode_Pict_Id"`

	Outlet_Key      string    `bson:"Outlet_Key" json:"Outlet_Key"`
	Outlet_Id       int64     `bson:"Outlet_Id" json:"Outlet_Id"`
	Outlet_LinkDate time.Time `bson:"Outlet_LinkDate" json:"Outlet_LinkDate"`

	//Workflow status
	Status            string    `bson:"Status" json:"Status"` //Available, Used, replaced
	StatusDate        time.Time `bson:"StatusDate" json:"StatusDate"`
	StatusDescription string    `bson:"StatusDescription" json:"StatusDescription"`
	StatusUser        string    `bson:"StatusUser" json:"StatusUser"`

	//Events Log trail
	Event_Logs []Event_Log `bson:"Event_Logs" json:"-"`
}

type Outlet_Qrcode_Add_Request struct {
	Barcodes_count  int    `bson:"Barcodes_count" json:"Barcodes_count"`   //how many to generate
	Barcodes_length int    `bson:"Barcodes_length" json:"Barcodes_length"` //string length
	BarcodeSize     int    `bson:"BarcodeSize" json:"BarcodeSize"`         //image size
	BatchNo         string `bson:"BatchNo" json:"BatchNo"`
}

type Outlet_Qrcode_Add_Response struct {
	Barcodes_count  int    `bson:"Barcodes_count" json:"Barcodes_count"`   //how many to generate
	Barcodes_length int    `bson:"Barcodes_length" json:"Barcodes_length"` //string length
	BarcodeSize     int    `bson:"BarcodeSize" json:"BarcodeSize"`         //image size
	BatchNo         string `bson:"BatchNo" json:"BatchNo"`
	From_Id         int64  `bson:"From_Id" json:"From_Id"`
	Till_Id         int64  `bson:"Till_Id" json:"Till_Id"`
}

type Outlet struct {
	Key             string           `bson:"Key" json:"Key"` //outlet Name
	Outlet_Id       int64            `bson:"Outlet_Id" json:"Outlet_Id"`
	Qrcode          string           `bson:"Qrcode" json:"Qrcode"`
	Business_type   string           `bson:"Business_type" json:"Business_type"`
	Outlet_Contacts []Outlet_Contact `bson:"Outlet_Contacts" json:"Outlet_Contacts"`

	Branding_Type string `bson:"Branding_Type" json:"Branding_Type"` //check list including: round, rectangle, paint...

	CreateDate time.Time `bson:"CreateDate" json:"CreateDate"`
	CreatedBy  string    `bson:"CreatedBy" json:"CreatedBy"` //Agent name

	//Address
	Provinces       string `bson:"Provinces" json:"Provinces"`
	Municipality    string `bson:"Municipality" json:"Municipality"`
	Street_Address  string `bson:"Street_Address" json:"Street_Address"`
	Address_Note    string `bson:"Address_Note" json:"Address_Note"`
	Zone            string `bson:"Zone" json:"Zone"`                       //should appear automatically
	Zone_Supervisor string `bson:"Zone_Supervisor" json:"Zone_Supervisor"` //should appear automatically
	//map coordination
	GPS_Location Location `bson:"GPS_Location" json:"GPS_Location"`

	//
	Agent_Key  string `bson:"Agent_Key" json:"Agent_Key"`
	Dealer_Key string `bson:"Dealer_Key" json:"Dealer_Key"`

	//Africell avilable Stock Type (MSL = minimum stock level)
	SIM_Card          bool   `bson:"SIM_Card" json:"SIM_Card"`
	MSL_SIM_Card      int64  `bson:"MSL_SIM_Card" json:"MSL_SIM_Card"`
	Scratch_Card      bool   `bson:"Scratch_Card" json:"Scratch_Card"`
	MSL_Scratch_Card  int64  `bson:"MSL_Scratch_Card" json:"MSL_Scratch_Card"`
	EVoucher          bool   `bson:"EVoucher" json:"EVoucher"`
	MSL_EVoucher      int64  `bson:"MSL_EVoucher" json:"MSL_EVoucher"`
	EVoucher_Account  string `bson:"EVoucher_Account" json:"EVoucher_Account"` //evoucher msisdn
	AfriMoney         bool   `bson:"AfriMoney" json:"AfriMoney"`
	MSL_AfriMoney     int64  `bson:"MSL_AfriMoney" json:"MSL_AfriMoney"`
	AfriMoney_Account string `bson:"AfriMoney_Account" json:"AfriMoney_Account"` //mobile money msisdn

	Have_Unitel_Stock bool `bson:"Have_Unitel_Stock" json:"Have_Unitel_Stock"`

	Outlet_Album_Ids  []int64 `bson:"Outlet_Album_Ids" json:"Outlet_Album_Ids"`
	Active_Alarms_Ids []int64 `bson:"Active_Alarms_Ids" json:"Active_Alarms_Ids"`

	//Last Visit
	Last_Visit_Date            time.Time     `bson:"Last_Visit_Date" json:"Last_Visit_Date"`
	Last_Visit_Agent           string        `bson:"Last_Visit_Agent" json:"Last_Visit_Agent"`
	Next_Planned_Visit_Date    time.Time     `bson:"Next_Planned_Visit_Date" json:"Next_Planned_Visit_Date"`
	Days_till_next_visit       time.Duration `bson:"Days_till_next_visit" json:"Days_till_next_visit"`
	Bln_DelayedVisit_Alarm     bool          `bson:"Bln_DelayedVisit_Alarm" json:"Bln_DelayedVisit_Alarm"`
	DelayedVisit_AlarmId       int64         `bson:"DelayedVisit_AlarmId" json:"DelayedVisit_AlarmId"`
	DelayedVisit_AlarmPosition int           `bson:"DelayedVisit_AlarmPosition" json:"DelayedVisit_AlarmPosition"`
	DelayedVisit_AlarmTime     time.Time     `bson:"DelayedVisit_AlarmTime" json:"DelayedVisit_AlarmTime"`
	DelayedVisit_AlarmSeverity string        `bson:"DelayedVisit_AlarmSeverity" json:"DelayedVisit_AlarmSeverity"`

	EVoucher_Balance       float64   `bson:"EVoucher_Balance" json:"EVoucher_Balance"`
	EVoucher_Balance_time  time.Time `bson:"EVoucher_Balance_time" json:"EVoucher_Balance_time"`
	EVoucher_Balance_err   string    `bson:"EVoucher_Balance_err" json:"EVoucher_Balance_err"`
	Bln_EVoucher_Alarm     bool      `bson:"Bln_EVoucher_Alarm" json:"Bln_EVoucher_Alarm"`
	EVoucher_AlarmId       int64     `bson:"EVoucher_AlarmId" json:"EVoucher_AlarmId"`
	EVoucher_AlarmPosition int       `bson:"EVoucher_AlarmPosition" json:"EVoucher_AlarmPosition"`
	EVoucher_AlarmTime     time.Time `bson:"EVoucher_AlarmTime" json:"EVoucher_AlarmTime"`
	EVoucher_AlarmSeverity string    `bson:"EVoucher_AlarmSeverity" json:"EVoucher_AlarmSeverity"`

	//Opening Hours
	OpeningHours []OutletOpenHours `bson:"OpeningHours" json:"OpeningHours"`
	Visits_Id    []int64           `bson:"Visits_Id" json:"Visits_Id"`

	Note string `bson:"Note" json:"Note"`

	//Workflow status
	Status            string    `bson:"Status" json:"Status"`
	StatusDate        time.Time `bson:"StatusDate" json:"StatusDate"`
	StatusDescription string    `bson:"StatusDescription" json:"StatusDescription"`
	StatusUser        string    `bson:"StatusUser" json:"StatusUser"`

	//Events Log trail
	Event_Logs []Event_Log `bson:"Event_Logs" json:"-"`
}

type Outlets_Nearby_Distance struct {
	Key                  string   `bson:"Key" json:"Key"` //outlet Name
	GPS_Location         Location `bson:"GPS_Location" json:"GPS_Location"`
	DistanceFromStarting float64  `bson:"DistanceFromStarting" json:"DistanceFromStarting"`
}

type OutletOpenHours struct {
	WeekDay     string `bson:"WeekDay" json:"WeekDay"`
	OpeningTime string `bson:"OpeningTime" json:"OpeningTime"`
	ClosingTime string `bson:"ClosingTime" json:"ClosingTime"`
}

type Outlet_Add_Request struct {
	Key             string           `bson:"Key" json:"Key"` //outlet Name
	Outlet_Id       int64            `bson:"Outlet_Id" json:"Outlet_Id"`
	Qrcode          string           `bson:"Qrcode" json:"Qrcode"`
	Business_type   string           `bson:"Business_type" json:"Business_type"`
	Outlet_Contacts []Outlet_Contact `bson:"Outlet_Contacts" json:"Outlet_Contacts"`
	Branding_Type   string           `bson:"Branding_Type" json:"Branding_Type"` //check list including: round, rectangle, paint...

	//Address
	Provinces       string `bson:"Provinces" json:"Provinces"`
	Municipality    string `bson:"Municipality" json:"Municipality"`
	Street_Address  string `bson:"Street_Address" json:"Street_Address"`
	Address_Note    string `bson:"Address_Note" json:"Address_Note"`
	Zone            string `bson:"Zone" json:"Zone"`                       //should appear automatically
	Zone_Supervisor string `bson:"Zone_Supervisor" json:"Zone_Supervisor"` //should appear automatically
	//map coordination
	Latitude  float64 `bson:"Latitude" json:"Latitude"`
	Longitude float64 `bson:"Longitude" json:"Longitude"`

	Dealer_Key string `bson:"Dealer_Key" json:"Dealer_Key"`

	//Africell avilable Stock Type (MSL = minimum stock level)
	SIM_Card          bool   `bson:"SIM_Card" json:"SIM_Card"`
	MSL_SIM_Card      int64  `bson:"MSL_SIM_Card" json:"MSL_SIM_Card"`
	Scratch_Card      bool   `bson:"Scratch_Card" json:"Scratch_Card"`
	MSL_Scratch_Card  int64  `bson:"MSL_Scratch_Card" json:"MSL_Scratch_Card"`
	EVoucher          bool   `bson:"EVoucher" json:"EVoucher"`
	MSL_EVoucher      int64  `bson:"MSL_EVoucher" json:"MSL_EVoucher"`
	EVoucher_Account  string `bson:"EVoucher_Account" json:"EVoucher_Account"` //evoucher msisdn
	AfriMoney         bool   `bson:"AfriMoney" json:"AfriMoney"`
	MSL_AfriMoney     int64  `bson:"MSL_AfriMoney" json:"MSL_AfriMoney"`
	AfriMoney_Account string `bson:"AfriMoney_Account" json:"AfriMoney_Account"` //mobile money msisdn
	Have_Unitel_Stock bool   `bson:"Have_Unitel_Stock" json:"Have_Unitel_Stock"`

	OpeningHours         []OutletOpenHours `bson:"OpeningHours" json:"OpeningHours"`
	Days_till_next_visit time.Duration     `bson:"Days_till_next_visit" json:"Days_till_next_visit"`

	Note string `bson:"Note" json:"Note"`
}

type Outlet_Edit_Request struct {
	Key             string           `bson:"Key" json:"Key"` //outlet Name
	Outlet_Id       int64            `bson:"Outlet_Id" json:"Outlet_Id"`
	NewKey          string           `bson:"NewKey" json:"NewKey"` //New outlet Name
	Business_type   string           `bson:"Business_type" json:"Business_type"`
	Outlet_Contacts []Outlet_Contact `bson:"Outlet_Contacts" json:"Outlet_Contacts"`
	Branding_Type   string           `bson:"Branding_Type" json:"Branding_Type"` //check list including: round, rectangle, paint...

	//Address
	Provinces       string `bson:"Provinces" json:"Provinces"`
	Municipality    string `bson:"Municipality" json:"Municipality"`
	Street_Address  string `bson:"Street_Address" json:"Street_Address"`
	Address_Note    string `bson:"Address_Note" json:"Address_Note"`
	Zone            string `bson:"Zone" json:"Zone"`                       //should appear automatically
	Zone_Supervisor string `bson:"Zone_Supervisor" json:"Zone_Supervisor"` //should appear automatically
	//map coordination
	Latitude  float64 `bson:"Latitude" json:"Latitude"`
	Longitude float64 `bson:"Longitude" json:"Longitude"`

	Dealer_Key string `bson:"Dealer_Key" json:"Dealer_Key"`

	//Africell avilable Stock Type (MSL = minimum stock level)
	SIM_Card          bool   `bson:"SIM_Card" json:"SIM_Card"`
	MSL_SIM_Card      int64  `bson:"MSL_SIM_Card" json:"MSL_SIM_Card"`
	Scratch_Card      bool   `bson:"Scratch_Card" json:"Scratch_Card"`
	MSL_Scratch_Card  int64  `bson:"MSL_Scratch_Card" json:"MSL_Scratch_Card"`
	EVoucher          bool   `bson:"EVoucher" json:"EVoucher"`
	MSL_EVoucher      int64  `bson:"MSL_EVoucher" json:"MSL_EVoucher"`
	EVoucher_Account  string `bson:"EVoucher_Account" json:"EVoucher_Account"` //evoucher msisdn
	AfriMoney         bool   `bson:"AfriMoney" json:"AfriMoney"`
	MSL_AfriMoney     int64  `bson:"MSL_AfriMoney" json:"MSL_AfriMoney"`
	AfriMoney_Account string `bson:"AfriMoney_Account" json:"AfriMoney_Account"` //mobile money msisdn
	Have_Unitel_Stock bool   `bson:"Have_Unitel_Stock" json:"Have_Unitel_Stock"`

	OpeningHours         []OutletOpenHours `bson:"OpeningHours" json:"OpeningHours"`
	Days_till_next_visit time.Duration     `bson:"Days_till_next_visit" json:"Days_till_next_visit"`

	Note string `bson:"Note" json:"Note"`
}

//	Valid longitude values are between -180 and 180, both inclusive.
//
// Valid latitude values are between -90 and 90, both inclusive.
type Location struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

type Outlet_Contact struct {
	Name               string `bson:"Name" json:"Name"`
	Type               string `bson:"Type" json:"Type"`
	Mobile             string `bson:"Mobile" json:"Mobile"`
	Alternative_Mobile string `bson:"Alternative_Mobile" json:"Alternative_Mobile"`
	Email              string `bson:"Email" json:"Email"`
	Note               string `bson:"Note" json:"Note"`
}

type Outlet_Album struct {
	Key                 string    `bson:"Key" json:"Key"` //Outlet_Key | Album_Id
	Album_Id            int64     `bson:"Album_Id" json:"Album_Id"`
	Album_Name          string    `bson:"Album_Name" json:"Album_Name"` //Album Date
	Outlet_Key          string    `bson:"Outlet_Key" json:"Outlet_Key"` //outlet name
	Album_Date          time.Time `bson:"Album_Date" json:"Album_Date"`
	Agent_Key           string    `bson:"Agent_Key" json:"Agent_Key"`
	GPS_Location        Location  `bson:"GPS_Location" json:"GPS_Location"`
	Album_Description   string    `bson:"Album_Description" json:"Album_Description"`
	Image_Ids           []string  `bson:"Image_Ids" json:"Image_Ids"`
	Image_Thumbnail_Ids []string  `bson:"Image_Thumbnail_Ids" json:"Image_Thumbnail_Ids"`
}

type Outlet_Album_Add_Request struct {
	Album_Id          int64   `bson:"Album_Id" json:"Album_Id"`     //Outlet_Key | Album_Id
	Album_Name        string  `bson:"Album_Name" json:"Album_Name"` //Album Date
	Outlet_Key        string  `bson:"Outlet_Key" json:"Outlet_Key"` //outlet name
	Latitude          float64 `bson:"Latitude" json:"Latitude"`
	Longitude         float64 `bson:"Longitude" json:"Longitude"`
	Album_Description string  `bson:"Album_Description" json:"Album_Description"`
	//Image_Ids         []string  `bson:"Image_Ids" json:"Image_Ids"`
}

type Outlet_Album_Edit_Request struct {
	Album_Id          int64  `bson:"Album_Id" json:"Album_Id"`
	Outlet_Key        string `bson:"Outlet_Key" json:"Outlet_Key"` //outlet name
	Album_Name        string `bson:"Album_Name" json:"Album_Name"` //Album Date
	Album_Description string `bson:"Album_Description" json:"Album_Description"`
	//Image_Ids         []string  `bson:"Image_Ids" json:"Image_Ids"`
}

type Outlet_Image struct {
	Image_Id   string  `bson:"Image_Id" json:"Image_Id"`
	Image_b64  string  `bson:"Image_b64" json:"Image_b64"`
	UploadedBy string  `bson:"UploadedBy" json:"UploadedBy"`
	UploadTime string  `bson:"UploadTime" json:"UploadTime"`
	Latitude   float64 `bson:"Latitude" json:"Latitude"`
	Longitude  float64 `bson:"Longitude" json:"Longitude"`
}

type Outlet_Album_Image_Add_Request struct {
	Image_Id           string  `bson:"Image_Id" json:"Image_Id"`
	Thumbnail_Image_Id string  `bson:"Thumbnail_Image_Id" json:"Thumbnail_Image_Id"`
	Album_Id           int64   `bson:"Album_Id" json:"Album_Id"`
	Outlet_Key         string  `bson:"Outlet_Key" json:"Outlet_Key"`
	Image_b64          string  `bson:"Image_b64" json:"Image_b64"`
	Latitude           float64 `bson:"Latitude" json:"Latitude"`
	Longitude          float64 `bson:"Longitude" json:"Longitude"`
}

type Outlet_Visits struct {
	Key        string `bson:"Key" json:"Key"` //Outlet_Key | Visit_Id
	Visit_Id   int64  `bson:"Visit_Id" json:"Visit_Id"`
	Outlet_Key string `bson:"Outlet_Key" json:"Outlet_Key"` //outlet name

	GPS_Location        Location  `bson:"GPS_Location" json:"GPS_Location"`
	Visit_Date          time.Time `bson:"Visit_Date" json:"Visit_Date"`
	Previous_Visit_Date time.Time `bson:"Previous_Visit_Date" json:"Previous_Visit_Date"`

	Outlet_Contact_Name string `bson:"Outlet_Contact_Name" json:"Outlet_Contact_Name"` //name of the person present in the outlet during the visit
	Agent_Key           string `bson:"Agent_Key" json:"Agent_Key"`                     //the agent doing the visit

	//Stock Type
	SIM_Card                   bool      `bson:"SIM_Card" json:"SIM_Card"`
	SIM_Card_Available_Quatity int64     `bson:"SIM_Card_Available_Quatity" json:"SIM_Card_Available_Quatity"`
	SIM_Card_Sold              int64     `bson:"SIM_Card_Sold" json:"SIM_Card_Sold"`
	SIM_Card_Last_Supply_Date  time.Time `bson:"SIM_Card_Last_Supply_Date" json:"SIM_Card_Last_Supply_Date"`

	Scratch_Card                   bool      `bson:"Scratch_Card" json:"Scratch_Card"`
	Scratch_Card_Available_Quatity int64     `bson:"Scratch_Card_Available_Quatity" json:"Scratch_Card_Available_Quatity"`
	Scratch_Card_Sold              int64     `bson:"Scratch_Card_Sold" json:"Scratch_Card_Sold"`
	Scratch_Card_Last_Supply_Date  time.Time `bson:"Scratch_Card_Last_Supply_Date" json:"Scratch_Card_Last_Supply_Date"`

	Scratch_Card_AQ_200   int64 `bson:"Scratch_Card_AQ_200" json:"Scratch_Card_AQ_200"`
	Scratch_Card_AQ_500   int64 `bson:"Scratch_Card_AQ_500" json:"Scratch_Card_AQ_500"`
	Scratch_Card_AQ_1000  int64 `bson:"Scratch_Card_AQ_1000" json:"Scratch_Card_AQ_1000"`
	Scratch_Card_AQ_2000  int64 `bson:"Scratch_Card_AQ_2000" json:"Scratch_Card_AQ_2000"`
	Scratch_Card_AQ_5000  int64 `bson:"Scratch_Card_AQ_5000" json:"Scratch_Card_AQ_5000"`
	Scratch_Card_AQ_10000 int64 `bson:"Scratch_Card_AQ_10000" json:"Scratch_Card_AQ_10000"`

	EVoucher                   bool      `bson:"EVoucher" json:"EVoucher"`
	EVoucher_Account           string    `bson:"EVoucher_Account" json:"EVoucher_Account"`
	EVoucher_Available_Quatity int64     `bson:"EVoucher_Available_Quatity" json:"EVoucher_Available_Quatity"`
	EVoucher_Sold              int64     `bson:"EVoucher_Sold" json:"EVoucher_Sold"`
	EVoucher_Last_Supply_Date  time.Time `bson:"EVoucher_Last_Supply_Date" json:"EVoucher_Last_Supply_Date"`

	AfriMoney                   bool      `bson:"AfriMoney" json:"AfriMoney"`
	AfriMoney_Account           string    `bson:"AfriMoney_Account" json:"AfriMoney_Account"`
	AfriMoney_Available_Quatity int64     `bson:"AfriMoney_Available_Quatity" json:"AfriMoney_Available_Quatity"`
	AfriMoney_Sold              int64     `bson:"AfriMoney_Sold" json:"AfriMoney_Sold"`
	AfriMoney_Last_Supply_Date  time.Time `bson:"AfriMoney_Last_Supply_Date" json:"AfriMoney_Last_Supply_Date"`

	Note string `bson:"Note" json:"Note"`
}

type Outlet_Visits_Add_Request struct {
	Outlet_Key          string `bson:"Outlet_Key" json:"Outlet_Key"` //outlet name
	Visit_Id            int64  `bson:"Visit_Id" json:"Visit_Id"`
	Outlet_Contact_Name string `bson:"Outlet_Contact_Name" json:"Outlet_Contact_Name"` //name of the person present in the outlet during the visit

	Latitude  float64 `bson:"Latitude" json:"Latitude"`
	Longitude float64 `bson:"Longitude" json:"Longitude"`

	//Stock Type
	SIM_Card                   bool      `bson:"SIM_Card" json:"SIM_Card"`
	SIM_Card_Available_Quatity int64     `bson:"SIM_Card_Available_Quatity" json:"SIM_Card_Available_Quatity"`
	SIM_Card_Sold              int64     `bson:"SIM_Card_Sold" json:"SIM_Card_Sold"`
	SIM_Card_Last_Supply_Date  time.Time `bson:"SIM_Card_Last_Supply_Date" json:"SIM_Card_Last_Supply_Date"`

	Scratch_Card                   bool      `bson:"Scratch_Card" json:"Scratch_Card"`
	Scratch_Card_Available_Quatity int64     `bson:"Scratch_Card_Available_Quatity" json:"Scratch_Card_Available_Quatity"`
	Scratch_Card_Sold              int64     `bson:"Scratch_Card_Sold" json:"Scratch_Card_Sold"`
	Scratch_Card_Last_Supply_Date  time.Time `bson:"Scratch_Card_Last_Supply_Date" json:"Scratch_Card_Last_Supply_Date"`

	Scratch_Card_AQ_200   int64 `bson:"Scratch_Card_AQ_200" json:"Scratch_Card_AQ_200"`
	Scratch_Card_AQ_500   int64 `bson:"Scratch_Card_AQ_500" json:"Scratch_Card_AQ_500"`
	Scratch_Card_AQ_1000  int64 `bson:"Scratch_Card_AQ_1000" json:"Scratch_Card_AQ_1000"`
	Scratch_Card_AQ_2000  int64 `bson:"Scratch_Card_AQ_2000" json:"Scratch_Card_AQ_2000"`
	Scratch_Card_AQ_5000  int64 `bson:"Scratch_Card_AQ_5000" json:"Scratch_Card_AQ_5000"`
	Scratch_Card_AQ_10000 int64 `bson:"Scratch_Card_AQ_10000" json:"Scratch_Card_AQ_10000"`

	EVoucher                   bool      `bson:"EVoucher" json:"EVoucher"`
	EVoucher_Account           string    `bson:"EVoucher_Account" json:"EVoucher_Account"`
	EVoucher_Available_Quatity int64     `bson:"EVoucher_Available_Quatity" json:"EVoucher_Available_Quatity"`
	EVoucher_Sold              int64     `bson:"EVoucher_Sold" json:"EVoucher_Sold"`
	EVoucher_Last_Supply_Date  time.Time `bson:"EVoucher_Last_Supply_Date" json:"EVoucher_Last_Supply_Date"`

	AfriMoney                   bool      `bson:"AfriMoney" json:"AfriMoney"`
	AfriMoney_Account           string    `bson:"AfriMoney_Account" json:"AfriMoney_Account"`
	AfriMoney_Available_Quatity int64     `bson:"AfriMoney_Available_Quatity" json:"AfriMoney_Available_Quatity"`
	AfriMoney_Sold              int64     `bson:"AfriMoney_Sold" json:"AfriMoney_Sold"`
	AfriMoney_Last_Supply_Date  time.Time `bson:"AfriMoney_Last_Supply_Date" json:"AfriMoney_Last_Supply_Date"`

	Note string `bson:"Note" json:"Note"`
}

type Outlet_Visits_Edit_Request struct {
	Visit_Id   int64  `bson:"Visit_Id" json:"Visit_Id"`
	Outlet_Key string `bson:"Outlet_Key" json:"Outlet_Key"` //outlet name

	Outlet_Contact_Name string `bson:"Outlet_Contact_Name" json:"Outlet_Contact_Name"` //name of the person present in the outlet during the visit

	Latitude  float64 `bson:"Latitude" json:"Latitude"`
	Longitude float64 `bson:"Longitude" json:"Longitude"`

	//Stock Type
	SIM_Card                   bool      `bson:"SIM_Card" json:"SIM_Card"`
	SIM_Card_Available_Quatity int64     `bson:"SIM_Card_Available_Quatity" json:"SIM_Card_Available_Quatity"`
	SIM_Card_Sold              int64     `bson:"SIM_Card_Sold" json:"SIM_Card_Sold"`
	SIM_Card_Last_Supply_Date  time.Time `bson:"SIM_Card_Last_Supply_Date" json:"SIM_Card_Last_Supply_Date"`

	Scratch_Card                   bool      `bson:"Scratch_Card" json:"Scratch_Card"`
	Scratch_Card_Available_Quatity int64     `bson:"Scratch_Card_Available_Quatity" json:"Scratch_Card_Available_Quatity"`
	Scratch_Card_Sold              int64     `bson:"Scratch_Card_Sold" json:"Scratch_Card_Sold"`
	Scratch_Card_Last_Supply_Date  time.Time `bson:"Scratch_Card_Last_Supply_Date" json:"Scratch_Card_Last_Supply_Date"`

	Scratch_Card_AQ_200   int64 `bson:"Scratch_Card_AQ_200" json:"Scratch_Card_AQ_200"`
	Scratch_Card_AQ_500   int64 `bson:"Scratch_Card_AQ_500" json:"Scratch_Card_AQ_500"`
	Scratch_Card_AQ_1000  int64 `bson:"Scratch_Card_AQ_1000" json:"Scratch_Card_AQ_1000"`
	Scratch_Card_AQ_2000  int64 `bson:"Scratch_Card_AQ_2000" json:"Scratch_Card_AQ_2000"`
	Scratch_Card_AQ_5000  int64 `bson:"Scratch_Card_AQ_5000" json:"Scratch_Card_AQ_5000"`
	Scratch_Card_AQ_10000 int64 `bson:"Scratch_Card_AQ_10000" json:"Scratch_Card_AQ_10000"`

	EVoucher                   bool      `bson:"EVoucher" json:"EVoucher"`
	EVoucher_Account           string    `bson:"EVoucher_Account" json:"EVoucher_Account"`
	EVoucher_Available_Quatity int64     `bson:"EVoucher_Available_Quatity" json:"EVoucher_Available_Quatity"`
	EVoucher_Sold              int64     `bson:"EVoucher_Sold" json:"EVoucher_Sold"`
	EVoucher_Last_Supply_Date  time.Time `bson:"EVoucher_Last_Supply_Date" json:"EVoucher_Last_Supply_Date"`

	AfriMoney                   bool      `bson:"AfriMoney" json:"AfriMoney"`
	AfriMoney_Account           string    `bson:"AfriMoney_Account" json:"AfriMoney_Account"`
	AfriMoney_Available_Quatity int64     `bson:"AfriMoney_Available_Quatity" json:"AfriMoney_Available_Quatity"`
	AfriMoney_Sold              int64     `bson:"AfriMoney_Sold" json:"AfriMoney_Sold"`
	AfriMoney_Last_Supply_Date  time.Time `bson:"AfriMoney_Last_Supply_Date" json:"AfriMoney_Last_Supply_Date"`

	Note string `bson:"Note" json:"Note"`
}

type Alarm struct {
	Alarm_Id    int64  `bson:"Alarm_Id" json:"Alarm_Id"`
	Outlet_Key  string `bson:"Outlet_Key" json:"Outlet_Key"` //outlet Name
	Visit_Id    int64  `bson:"Visit_Id" json:"Visit_Id"`
	Album_Id    int64  `bson:"Album_Id" json:"Album_Id"`
	Type        string `bson:"Type" json:"Type"` //Low Stock, Not Visited
	Description string `bson:"Description" json:"Description"`
	Message     string `bson:"Message" json:"Message"`
	Severity    string `bson:"Severity" json:"Severity"` //Critical, High, Medium, and Low

	CreateTime time.Time `bson:"CreateTime" json:"CreateTime"`
	CreatedBy  string    `bson:"CreatedBy" json:"CreatedBy"`

	Acknowledgment_Time time.Time `bson:"Acknowledgment_Time" json:"Acknowledgment_Time"`
	Acknowledged_By     string    `bson:"Acknowledged_By" json:"Acknowledged_By"`
	Acknowledgment_Note string    `bson:"Acknowledgment_Note" json:"Acknowledgment_Note"`

	Disabled_Time  time.Time `bson:"Disabled_Time" json:"Disabled_Time"`
	Disabled_By    string    `bson:"Disabled_By" json:"Disabled_By"` //User, Logic
	Disable_Reason string    `bson:"Disable_Reason" json:"Disable_Reason"`

	//Workflow status
	Status            string    `bson:"Status" json:"Status"` //Created, acknowledged, Disabled
	StatusDate        time.Time `bson:"StatusDate" json:"StatusDate"`
	StatusDescription string    `bson:"StatusDescription" json:"StatusDescription"`
	StatusUser        string    `bson:"StatusUser" json:"StatusUser"`

	//Events Log trail
	Event_Logs []Event_Log `bson:"Event_Logs" json:"-"`
}

type Alarm_Add_Request struct {
	Alarm_Id    int64  `bson:"Alarm_Id" json:"Alarm_Id"`
	Outlet_Key  string `bson:"Outlet_Key" json:"Outlet_Key"` //outlet Name
	Visit_Id    int64  `bson:"Visit_Id" json:"Visit_Id"`
	Album_Id    int64  `bson:"Album_Id" json:"Album_Id"`
	Type        string `bson:"Type" json:"Type"` //Low Stock, Not Visited
	Description string `bson:"Description" json:"Description"`
	Message     string `bson:"Message" json:"Message"`
	Severity    string `bson:"Severity" json:"Severity"` //Critical, High, Medium, and Low
}

type Alarm_Disable_Request struct {
	Alarm_Id       int64  `bson:"Alarm_Id" json:"Alarm_Id"`
	Disabled_By    string `bson:"Disabled_By" json:"Disabled_By"` //User, Logic
	Disable_Reason string `bson:"Disable_Reason" json:"Disable_Reason"`
}

type DefaultValues struct {
	Key                         string   `bson:"Key" json:"Key"` //"Default"
	MSL_SIM_Card                int64    `bson:"MSL_SIM_Card" json:"MSL_SIM_Card"`
	MSL_Scratch_Card            int64    `bson:"MSL_Scratch_Card" json:"MSL_Scratch_Card"`
	MSL_EVoucher                int64    `bson:"MSL_EVoucher" json:"MSL_EVoucher"`
	MSL_AfriMoney               int64    `bson:"MSL_AfriMoney" json:"MSL_AfriMoney"`
	Scratch_Cards_Denominations []string `bson:"Scratch_Cards_Denominations" json:"Scratch_Cards_Denominations"`
}

type Agent_Sales struct {
	Key                   string    `bson:"Key" json:"Key"` //Outlet_Key | Visit_Id
	Sales_Id              int64     `bson:"Sales_Id" json:"Sales_Id"`
	Outlet_Key            string    `bson:"Outlet_Key" json:"Outlet_Key"`                   //outlet name
	Outlet_Contact_Name   string    `bson:"Outlet_Contact_Name" json:"Outlet_Contact_Name"` //name of the person present in the outlet during the visit
	GPS_Location          Location  `bson:"GPS_Location" json:"GPS_Location"`
	Sale_Date             time.Time `bson:"Sale_Date" json:"Sale_Date"`
	SIM_Card              int64     `bson:"SIM_Card" json:"SIM_Card"`
	Scratch_Card_AQ_200   int64     `bson:"Scratch_Card_AQ_200" json:"Scratch_Card_AQ_200"`
	Scratch_Card_AQ_500   int64     `bson:"Scratch_Card_AQ_500" json:"Scratch_Card_AQ_500"`
	Scratch_Card_AQ_1000  int64     `bson:"Scratch_Card_AQ_1000" json:"Scratch_Card_AQ_1000"`
	Scratch_Card_AQ_2000  int64     `bson:"Scratch_Card_AQ_2000" json:"Scratch_Card_AQ_2000"`
	Scratch_Card_AQ_5000  int64     `bson:"Scratch_Card_AQ_5000" json:"Scratch_Card_AQ_5000"`
	Scratch_Card_AQ_10000 int64     `bson:"Scratch_Card_AQ_10000" json:"Scratch_Card_AQ_10000"`
	EVoucher              int64     `bson:"EVoucher" json:"EVoucher"`
	AfriMoney             int64     `bson:"AfriMoney" json:"AfriMoney"`
	Note                  string    `bson:"Note" json:"Note"`
}

type Agent_Sales_Add_Request struct {
	Sales_Id            int64   `bson:"Sales_Id" json:"Sales_Id"`
	Outlet_Key          string  `bson:"Outlet_Key" json:"Outlet_Key"`                   //outlet name
	Outlet_Contact_Name string  `bson:"Outlet_Contact_Name" json:"Outlet_Contact_Name"` //name of the person present in the outlet during the visit
	Latitude            float64 `bson:"Latitude" json:"Latitude"`
	Longitude           float64 `bson:"Longitude" json:"Longitude"`
	//Stock Type
	SIM_Card              int64  `bson:"SIM_Card" json:"SIM_Card"`
	Scratch_Card_AQ_200   int64  `bson:"Scratch_Card_AQ_200" json:"Scratch_Card_AQ_200"`
	Scratch_Card_AQ_500   int64  `bson:"Scratch_Card_AQ_500" json:"Scratch_Card_AQ_500"`
	Scratch_Card_AQ_1000  int64  `bson:"Scratch_Card_AQ_1000" json:"Scratch_Card_AQ_1000"`
	Scratch_Card_AQ_2000  int64  `bson:"Scratch_Card_AQ_2000" json:"Scratch_Card_AQ_2000"`
	Scratch_Card_AQ_5000  int64  `bson:"Scratch_Card_AQ_5000" json:"Scratch_Card_AQ_5000"`
	Scratch_Card_AQ_10000 int64  `bson:"Scratch_Card_AQ_10000" json:"Scratch_Card_AQ_10000"`
	EVoucher              int64  `bson:"EVoucher" json:"EVoucher"`
	AfriMoney             int64  `bson:"AfriMoney" json:"AfriMoney"`
	Note                  string `bson:"Note" json:"Note"`
}

type Agent_Sales_Edit_Request struct {
	Sales_Id            int64   `bson:"Sales_Id" json:"Sales_Id"`
	Outlet_Key          string  `bson:"Outlet_Key" json:"Outlet_Key"`                   //outlet name
	Outlet_Contact_Name string  `bson:"Outlet_Contact_Name" json:"Outlet_Contact_Name"` //name of the person present in the outlet during the visit
	Latitude            float64 `bson:"Latitude" json:"Latitude"`
	Longitude           float64 `bson:"Longitude" json:"Longitude"`
	//Stock Type
	SIM_Card              int64  `bson:"SIM_Card" json:"SIM_Card"`
	Scratch_Card_AQ_200   int64  `bson:"Scratch_Card_AQ_200" json:"Scratch_Card_AQ_200"`
	Scratch_Card_AQ_500   int64  `bson:"Scratch_Card_AQ_500" json:"Scratch_Card_AQ_500"`
	Scratch_Card_AQ_1000  int64  `bson:"Scratch_Card_AQ_1000" json:"Scratch_Card_AQ_1000"`
	Scratch_Card_AQ_2000  int64  `bson:"Scratch_Card_AQ_2000" json:"Scratch_Card_AQ_2000"`
	Scratch_Card_AQ_5000  int64  `bson:"Scratch_Card_AQ_5000" json:"Scratch_Card_AQ_5000"`
	Scratch_Card_AQ_10000 int64  `bson:"Scratch_Card_AQ_10000" json:"Scratch_Card_AQ_10000"`
	EVoucher              int64  `bson:"EVoucher" json:"EVoucher"`
	AfriMoney             int64  `bson:"AfriMoney" json:"AfriMoney"`
	Note                  string `bson:"Note" json:"Note"`
}
