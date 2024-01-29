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

type Subscriber struct {
	Key           string `bson:"Key" json:"Key"` //MSISDN
	Subscriber_Id int64  `bson:"Subscriber_Id" json:"Subscriber_Id"`

	COS                 string    `bson:"COS" json:"COS"`
	FirstUse_date       time.Time `bson:"FirstUse_date" json:"FirstUse_date"`
	ARPU                float64   `bson:"ARPU" json:"ARPU"`
	ARPU_date           time.Time `bson:"ARPU_date" json:"ARPU_date"`
	IsLendmeEligible    bool      `bson:"IsLendmeEligible" json:"IsLendmeEligible"`
	Credit_Limit_Scheme string    `bson:"Credit_Limit_Scheme" json:"Credit_Limit_Scheme"`

	//Workflow status
	Status            string    `bson:"Status" json:"Status"`
	StatusDate        time.Time `bson:"StatusDate" json:"StatusDate"`
	StatusUser        string    `bson:"StatusUser" json:"StatusUser"`
	StatusDescription string    `bson:"StatusDescription" json:"StatusDescription"`
	//Events Log trail
	Event_Logs []Event_Log `bson:"Event_Logs" json:"Event_Logs"`
}

type Subscriber_Add_Request struct {
	Key           string `bson:"Key" json:"Key"` //MSISDN
	Subscriber_Id int64  `bson:"Subscriber_Id" json:"Subscriber_Id"`
}

type Subscriber_Edit_Request struct {
	Key           string `bson:"Key" json:"Key"` //MSISDN
	NewKey        string `bson:"NewKey" json:"NewKey"`
	Subscriber_Id int64  `bson:"Subscriber_Id" json:"Subscriber_Id"`
}

type Credit_Limit_Scheme struct {
	Key                 string      `bson:"Key" json:"Key"` //Scheme name
	Scheme_Id           int64       `bson:"Scheme_Id" json:"Scheme_Id"`
	ARPU_From           float64     `bson:"ARPU_From" json:"ARPU_From"`
	ARPU_Till           float64     `bson:"ARPU_Till" json:"ARPU_Till"`
	AON_From            int         `bson:"AON_From" json:"AON_From"` //months
	AON_Till            int         `bson:"AON_Till" json:"AON_Till"` //months
	Credit_limit_Amount float64     `bson:"Credit_limit_Amount" json:"Credit_limit_Amount"`
	Event_Logs          []Event_Log `bson:"Event_Logs" json:"Event_Logs"`
}

type Credit_Limit_Scheme_Add_Request struct {
	Key                 string  `bson:"Key" json:"Key"` //Scheme name
	Scheme_Id           int64   `bson:"Scheme_Id" json:"Scheme_Id"`
	ARPU_From           float64 `bson:"ARPU_From" json:"ARPU_From"`
	ARPU_Till           float64 `bson:"ARPU_Till" json:"ARPU_Till"`
	AON_From            int     `bson:"AON_From" json:"AON_From"` //months
	AON_Till            int     `bson:"AON_Till" json:"AON_Till"` //months
	Credit_limit_Amount float64 `bson:"Credit_limit_Amount" json:"Credit_limit_Amount"`
}

type Credit_Limit_Scheme_Edit_Request struct {
	Key                 string  `bson:"Key" json:"Key"` //Scheme name
	NewKey              string  `bson:"NewKey" json:"NewKey"`
	Scheme_Id           int64   `bson:"Scheme_Id" json:"Scheme_Id"`
	ARPU_From           float64 `bson:"ARPU_From" json:"ARPU_From"`
	ARPU_Till           float64 `bson:"ARPU_Till" json:"ARPU_Till"`
	AON_From            int     `bson:"AON_From" json:"AON_From"` //months
	AON_Till            int     `bson:"AON_Till" json:"AON_Till"` //months
	Credit_limit_Amount float64 `bson:"Credit_limit_Amount" json:"Credit_limit_Amount"`
}

type DefaultValues struct {
	Key string `bson:"Key" json:"Key"` //"Default"
}
