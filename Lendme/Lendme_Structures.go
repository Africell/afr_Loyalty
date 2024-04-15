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
	Key                     string    `bson:"Key" json:"Key"` //MSISDN
	Subscriber_Id           int64     `bson:"Subscriber_Id" json:"Subscriber_Id"`
	Add_Date                time.Time `bson:"Add_Date" json:"Add_Date"`
	Last_ProfileUpdate_date time.Time `bson:"Last_ProfileUpdate_date" json:"Last_ProfileUpdate_date"`

	COS           string    `bson:"COS" json:"COS"`
	FirstUse_date time.Time `bson:"FirstUse_date" json:"FirstUse_date"`
	Last_Credit   time.Time `bson:"Last_Credit" json:"Last_Credit"`

	IN_Loyalty_Status                string    `bson:"IN_Loyalty_Status" json:"IN_Loyalty_Status"`
	IN_Credit_Limit                  float64   `bson:"IN_Credit_Limit" json:"IN_Credit_Limit"`
	ARPU                             float64   `bson:"ARPU" json:"ARPU"`
	Recharge                         float64   `bson:"Recharge" json:"Recharge"`
	Last_Recharge_Date               time.Time `bson:"Last_Recharge_Date" json:"Last_Recharge_Date"`
	Dealer_Bundle_Purchase           float64   `bson:"Dealer_Bundle_Purchase" json:"Dealer_Bundle_Purchase"`
	Last_Dealer_Bundle_Purchase_Date time.Time `bson:"Last_Dealer_Bundle_Purchase_Date" json:"Last_Dealer_Bundle_Purchase_Date"`

	IsLendmeEligible         bool      `bson:"IsLendmeEligible" json:"IsLendmeEligible"`
	Credit_Limit_Scheme      string    `bson:"Credit_Limit_Scheme" json:"Credit_Limit_Scheme"`
	Credit_Limit_Scheme_Date time.Time `bson:"Credit_Limit_Scheme_Date" json:"Credit_Limit_Scheme_Date"`
	NotElligibleReason       string    `bson:"NotElligibleReason" json:"NotElligibleReason"`

	Lendme_Outstanding_Amount float64   `bson:"Lendme_Outstanding_Amount" json:"Lendme_Outstanding_Amount"`
	Lendme_Outstanding_Fee    float64   `bson:"Lendme_Outstanding_Fee" json:"Lendme_Outstanding_Fee"`
	Last_Lend_Date            time.Time `bson:"Last_Lend_Date" json:"Last_Lend_Date"`

	Cumulative_Lent_Amount float64   `bson:"Cumulative_Lent_Amount" json:"Cumulative_Lent_Amount"`
	Cumulative_Lent_Fee    float64   `bson:"Cumulative_Lent_Fee" json:"Cumulative_Lent_Fee"`
	Cumulative_Payback     float64   `bson:"Cumulative_Payback" json:"Cumulative_Payback"`
	Last_Payback_Date      time.Time `bson:"Last_Payback_Date" json:"Last_Payback_Date"`

	//Workflow status
	Status            string    `bson:"Status" json:"Status"`
	StatusDate        time.Time `bson:"StatusDate" json:"StatusDate"`
	StatusUser        string    `bson:"StatusUser" json:"StatusUser"`
	StatusDescription string    `bson:"StatusDescription" json:"StatusDescription"`
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
	Amount_From         float64     `bson:"Amount_From" json:"Amount_From"`
	Amount_Till         float64     `bson:"Amount_Till" json:"Amount_Till"`
	AON_From            float64     `bson:"AON_From" json:"AON_From"` //months
	AON_Till            float64     `bson:"AON_Till" json:"AON_Till"` //months
	Credit_limit_Amount float64     `bson:"Credit_limit_Amount" json:"Credit_limit_Amount"`
	Event_Logs          []Event_Log `bson:"Event_Logs" json:"Event_Logs"`
}

type Credit_Limit_Scheme_Add_Request struct {
	Key                 string  `bson:"Key" json:"Key"` //Scheme name
	Scheme_Id           int64   `bson:"Scheme_Id" json:"Scheme_Id"`
	Amount_From         float64 `bson:"Amount_From" json:"Amount_From"`
	Amount_Till         float64 `bson:"Amount_Till" json:"Amount_Till"`
	AON_From            float64 `bson:"AON_From" json:"AON_From"` //months
	AON_Till            float64 `bson:"AON_Till" json:"AON_Till"` //months
	Credit_limit_Amount float64 `bson:"Credit_limit_Amount" json:"Credit_limit_Amount"`
}

type Credit_Limit_Scheme_Edit_Request struct {
	Key                 string  `bson:"Key" json:"Key"` //Scheme name
	NewKey              string  `bson:"NewKey" json:"NewKey"`
	Scheme_Id           int64   `bson:"Scheme_Id" json:"Scheme_Id"`
	Amount_From         float64 `bson:"Amount_From" json:"Amount_From"`
	Amount_Till         float64 `bson:"Amount_Till" json:"Amount_Till"`
	AON_From            float64 `bson:"AON_From" json:"AON_From"` //months
	AON_Till            float64 `bson:"AON_Till" json:"AON_Till"` //months
	Credit_limit_Amount float64 `bson:"Credit_limit_Amount" json:"Credit_limit_Amount"`
}

type DefaultValues struct {
	Key string `bson:"Key" json:"Key"` //"Default"
}

type Lendme_log struct {
	Source        string    `bson:"Source" json:"Source"`
	MSISDN        string    `bson:"MSISDN" json:"MSISDN"`
	Log_Date      time.Time `bson:"Log_Date" json:"Log_Date"`
	Type          string    `bson:"Type" json:"Type"`
	Lendme_Amount float64   `bson:"Lendme_Amount" json:"Lendme_Amount"`
	Lendme_Fee    float64   `bson:"Lendme_Fee" json:"Lendme_Fee"`

	Status            string `bson:"Status" json:"Status"`
	StatusDescription string `bson:"StatusDescription" json:"StatusDescription"`
}

type LendMe_Request struct {
	Source string  `bson:"Source" json:"Source"`
	MSISDN string  `bson:"MSISDN" json:"MSISDN"`
	Amount float64 `bson:"Amount" json:"Amount"`
}

type Sub_Update_Request struct {
	MSISDN                           string    `bson:"MSISDN" json:"MSISDN"`
	COS                              string    `bson:"COS" json:"COS"`
	First_Used                       time.Time `bson:"First_Used" json:"First_Used"`
	Last_Credit                      time.Time `bson:"Last_Credit" json:"Last_Credit"`
	Loyalty_Status                   string    `bson:"Loyalty_Status" json:"Loyalty_Status"`
	Credit_Limit                     float64   `bson:"Credit_Limit" json:"Credit_Limit"`
	ARPU_Amount                      float64   `bson:"ARPU_Amount" json:"ARPU_Amount"`
	Recharge                         float64   `bson:"Recharge" json:"Recharge"`
	Last_Recharge_Date               time.Time `bson:"Last_Recharge_Date" json:"Last_Recharge_Date"`
	Dealer_Bundle_Purchase           float64   `bson:"Dealer_Bundle_Purchase" json:"Dealer_Bundle_Purchase"`
	Last_Dealer_Bundle_Purchase_Date time.Time `bson:"Last_Dealer_Bundle_Purchase_Date" json:"Last_Dealer_Bundle_Purchase_Date"`
}

type Subscriber_USSD struct {
	MSISDN string `bson:"MSISDN" json:"MSISDN"` //MSISDN

	IsLendmeEligible    bool    `bson:"IsLendmeEligible" json:"IsLendmeEligible"`
	Credit_Limit_Scheme string  `bson:"Credit_Limit_Scheme" json:"Credit_Limit_Scheme"`
	NotElligibleReason  string  `bson:"NotElligibleReason" json:"NotElligibleReason"`
	Credit_limit_Amount float64 `bson:"Credit_limit_Amount" json:"Credit_limit_Amount"`

	Min_Allowed_Amount float64 `bson:"Min_Allowed_Amount" json:"Min_Allowed_Amount"`
	Max_Allowed_Amount float64 `bson:"Max_Allowed_Amount" json:"Max_Allowed_Amount"`

	Lendme_Outstanding_Amount float64 `bson:"Lendme_Outstanding_Amount" json:"Lendme_Outstanding_Amount"`
	Lendme_Outstanding_Fee    float64 `bson:"Lendme_Outstanding_Fee" json:"Lendme_Outstanding_Fee"`
}
