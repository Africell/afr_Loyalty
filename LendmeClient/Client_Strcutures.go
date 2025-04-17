package LendmeClient

import (
	"daoc"
	"time"
)

type Loyalty_AccountDebitPoints_Request struct {
	MSISDN               string  `bson:"MSISDN" json:"MSISDN"` //MSISDN
	Debit_Amount         float64 `bson:"Debit_Amount" json:"Debit_Amount"`
	Debit_Reason         string  `bson:"Debit_Reason" json:"Debit_Reason"`
	Redemption_Type      string  `bson:"Redemption_Type" json:"Redemption_Type"` //Airtime, Bundle, MobileMoney, SpinAndWin
	Redemption_Bunlde_Id string  `bson:"Redemption_Bunlde_Id" json:"Redemption_Bunlde_Id"`
	Redemption_Amount    float64 `bson:"Redemption_Amount" json:"Redemption_Amount"`
}

type Loyalty_AccountDebitPoints_log struct {
	//request Header info
	SourceIP    string        `bson:"SourceIP" json:"-"`
	SourceApp   string        `bson:"SourceApp" json:"-"`
	AppLogin    string        `bson:"AppLogin" json:"-"`
	AppVersion  string        `bson:"AppVersion" json:"-"`
	GPSLocation daoc.Location `bson:"GPSLocation" json:"-"`
	GSMLocation string        `bson:"GSMLocation" json:"-"`

	//request detail
	MSISDN               string  `bson:"MSISDN" json:"MSISDN"` //MSISDN
	Debit_Amount         float64 `bson:"Debit_Amount" json:"Debit_Amount"`
	Debit_Reason         string  `bson:"Debit_Reason" json:"Debit_Reason"`
	Redemption_Type      string  `bson:"Redemption_Type" json:"Redemption_Type"` //Airtime, Bundle, MobileMoney, SpinAndWin
	Redemption_Bunlde_Id string  `bson:"Redemption_Bunlde_Id" json:"Redemption_Bunlde_Id"`
	Redemption_Amount    float64 `bson:"Redemption_Amount" json:"Redemption_Amount"`

	Customer_Id                 int64   `bson:"Customer_Id" json:"Customer_Id"`
	Account_Status              string  `bson:"Account_Status" json:"Account_Status"`
	Loyalty_Level_Key           string  `bson:"Loyalty_Level_Key" json:"Loyalty_Level_Key"`
	Loyalty_Account_Segment_Key string  `bson:"Loyalty_Account_Segment_Key" json:"Loyalty_Account_Segment_Key"`
	Awarded_Points              float64 `bson:"Awarded_Points" json:"Awarded_Points"`
	Opening_Redeemed_Points     float64 `bson:"Opening_Redeemed_Points" json:"Opening_Redeemed_Points"`
	Closure_Redeemed_Points     float64 `bson:"Closure_Redeemed_Points" json:"Closure_Redeemed_Points"`
	Available_Points            float64 `bson:"Available_Points" json:"Available_Points"` //(Awarded_Points + Expired_Points) - Redeemed_Points

	MinRequiredPoints               float64 `bson:"MinRequiredPoints" json:"MinRequiredPoints"`
	Allow_Negative_Balance_ToRedeem bool    `bson:"Allow_Negative_Balance_ToRedeem" json:"Allow_Negative_Balance_ToRedeem"`
	Allow_PendingLendme_ToRedeem    bool    `bson:"Allow_PendingLendme_ToRedeem" json:"Allow_PendingLendme_ToRedeem"`

	Airtime_PurchaseResult     interface{} `bson:"Airtime_PurchaseResult" json:"Airtime_PurchaseResult"`
	Bundle_PurchaseResult      interface{} `bson:"Bundle_PurchaseResult" json:"Bundle_PurchaseResult"`
	MobileMoney_PurchaseResult interface{} `bson:"MobileMoney_PurchaseResult" json:"MobileMoney_PurchaseResult"`
	SpinAndWin_PurchaseResult  interface{} `bson:"SpinAndWin_PurchaseResult" json:"SpinAndWin_PurchaseResult"`

	//response result
	ReceiveDate       time.Time `bson:"ReceiveDate" json:"-"`
	Status            string    `bson:"Status" json:"Status"` //successful, failed
	StatusCode        int       `bson:"StatusCode" json:"StatusCode"`
	StatusDescription string    `bson:"StatusDescription" json:"StatusDescription"`
	ErrorDescription  string    `bson:"ErrorDescription" json:"ErrorDescription"`
	StatusDate        time.Time `bson:"StatusDate" json:"StatusDate"`
	E2E_Elapsedtime   int64     `bson:"E2E_Elapsedtime" json:"E2E_Elapsedtime"` //receive date till return
}
