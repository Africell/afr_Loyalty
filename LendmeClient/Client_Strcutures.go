package LendmeClient

import (
	"daoc"
	"net/http"
	"time"
)

// **********************************************************************************************
// Loyalty structures
// **********************************************************************************************
type Generic_http_call_Request struct {
	Req                *http.Request
	Url                string
	Method             string
	Token              string
	BasicAuth_User     string
	BasicAuth_Password string
	OTP                string
	QueryParameters    map[string]string
	Headers            map[string]string
	Load               []byte
}

type Generic_http_call_Response struct {
	Body       []byte
	Header     http.Header
	Statuscode int
}

type Customer_Loyalty_Account struct {
	Key             string    `bson:"Key" json:"Key"` //MSISDN
	COS             string    `bson:"COS" json:"COS"`
	Joining_Date    time.Time `bson:"Joining_Date" json:"Joining_Date"` //
	ARPU            float64   `bson:"ARPU" json:"ARPU"`
	Customer_Id     int64     `bson:"Customer_Id" json:"Customer_Id"`
	Creation_date   time.Time `bson:"Creation_date" json:"Creation_date"`
	Account_Status  string    `bson:"Account_Status" json:"Account_Status"`
	Account_Profile string    `bson:"Account_Profile" json:"Account_Profile"`

	Multiplier_Percentage float64 `bson:"Multiplier_Percentage" json:"Multiplier_Percentage"`

	Previous_Loyalty_Level_Key  string    `bson:"Previous_Loyalty_Level_Key" json:"Previous_Loyalty_Level_Key"`
	Previous_Loyalty_Level_Date time.Time `bson:"Previous_Loyalty_Level_Date" json:"Previous_Loyalty_Level_Date"`

	Loyalty_Level_Key       string    `bson:"Loyalty_Level_Key" json:"Loyalty_Level_Key"`
	Loyalty_Level_Date      time.Time `bson:"Loyalty_Level_Date" json:"Loyalty_Level_Date"`
	Loyalty_Level_Direction string    `bson:"Loyalty_Level_Direction" json:"Loyalty_Level_Direction"` //last change: Up or Down
	Loyalty_Level_SetBy     string    `bson:"Loyalty_Level_SetBy" json:"Loyalty_Level_SetBy"`         //program or admin, if admin program cannot change anymore

	Loyalty_Account_Segment_Key       string    `bson:"Loyalty_Account_Segment_Key" json:"Loyalty_Account_Segment_Key"`
	Loyalty_Account_Segment_Date      time.Time `bson:"Loyalty_Account_Segment_Date" json:"Loyalty_Account_Segment_Date"`
	Loyalty_Account_Segment_Direction string    `bson:"Loyalty_Account_Segment_Direction" json:"Loyalty_Account_Segment_Direction"` //last change: Up or Down
	Loyalty_Account_Segment_SetBy     string    `bson:"Loyalty_Account_Segment_SetBy" json:"Loyalty_Account_Segment_SetBy"`         //program or admin, if admin program cannot change anymore

	Awarded_Points   float64   `bson:"Awarded_Points" json:"Awarded_Points"`
	Redeemed_Points  float64   `bson:"Redeemed_Points" json:"Redeemed_Points"`
	Available_Points float64   `bson:"Available_Points" json:"Available_Points"` //(Awarded_Points + Expired_Points) - Redeemed_Points
	Last_Award_Date  time.Time `bson:"Last_Award_Date" json:"Last_Award_Date"`
	Last_Redeem_Date time.Time `bson:"Last_Redeem_Date" json:"Last_Redeem_Date"`

	Points_To_Expire   float64   `bson:"Points_To_Expire" json:"Points_To_Expire"`
	Coming_Expiry_Date time.Time `bson:"Coming_Expiry_Date" json:"Coming_Expiry_Date"`

	Opt_Status               string    `bson:"Opt_Status" json:"Opt_Status"` //OptedIn, OptedOut
	Last_Opt_Status_Date     time.Time `bson:"Opt_Status_Date" json:"Opt_Status_Date"`
	First_Opt_In_Status_Date time.Time `bson:"First_Opt_In_Status_Date" json:"First_Opt_In_Status_Date"`

	Expired_Points              float64   `bson:"Expired_Points" json:"Expired_Points"`                   //expired are deducted from Awarded_Points
	Redeemed_Expired_Points     float64   `bson:"Redeemed_Expired_Points" json:"Redeemed_Expired_Points"` //redeemed expired are deducted from Awarded_Points
	Expiry_Date                 time.Time `bson:"Expiry_Date" json:"Expiry_Date"`
	Initial_Date                time.Time `bson:"Initial_Date" json:"Initial_Date"`
	Outstanding_fraction_points float64   `bson:"Outstanding_fraction_points" json:"Outstanding_fraction_points"`

	Points_Detail_Keys []string `bson:"Points_Detail_Keys" json:"Points_Detail_Keys"`
}

type Loyalty_Points_Detail struct {
	Year_Month       string    `bson:"Year_Month" json:"Year_Month"`
	Creation_date    time.Time `bson:"Creation_date" json:"Creation_date"`
	Awarded_Points   float64   `bson:"Awarded_Points" json:"Awarded_Points"`
	Redeemed_Points  float64   `bson:"Redeemed_Points" json:"Redeemed_Points"`
	Available_Points float64   `bson:"Available_Points" json:"Available_Points"` //(Awarded_Points + Expired_Points) - Redeemed_Points
	Expired_Points   float64   `bson:"Expired_Points" json:"Expired_Points"`     //expired are deducted from Awarded_Points
	Expiry_Date      time.Time `bson:"Expiry_Date" json:"Expiry_Date"`
}

type Loyalty_AccountDebitPoints_Request struct {
	MSISDN               string  `bson:"MSISDN" json:"MSISDN"`             //MSISDN
	Debit_Amount         float64 `bson:"Debit_Amount" json:"Debit_Amount"` //loyalty points to be deducted
	Debit_Reason         string  `bson:"Debit_Reason" json:"Debit_Reason"`
	Redemption_Type      string  `bson:"Redemption_Type" json:"Redemption_Type"` //Airtime, Bundle, MobileMoney, SpinAndWin
	Redemption_Bunlde_Id string  `bson:"Redemption_Bunlde_Id" json:"Redemption_Bunlde_Id"`
	Redemption_Amount    float64 `bson:"Redemption_Amount" json:"Redemption_Amount"` //airtime or money amount redeemed to customer
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

	Customer_Id                 int64  `bson:"Customer_Id" json:"Customer_Id"`
	Account_Status              string `bson:"Account_Status" json:"Account_Status"`
	Loyalty_Level_Key           string `bson:"Loyalty_Level_Key" json:"Loyalty_Level_Key"`
	Loyalty_Account_Segment_Key string `bson:"Loyalty_Account_Segment_Key" json:"Loyalty_Account_Segment_Key"`
	// Awarded_Points              float64 `bson:"Awarded_Points" json:"Awarded_Points"`
	// Opening_Redeemed_Points     float64 `bson:"Opening_Redeemed_Points" json:"Opening_Redeemed_Points"`
	// Closure_Redeemed_Points     float64 `bson:"Closure_Redeemed_Points" json:"Closure_Redeemed_Points"`
	// Available_Points            float64 `bson:"Available_Points" json:"Available_Points"` //(Awarded_Points + Expired_Points) - Redeemed_Points

	Opening_Awarded_Points   float64 `bson:"Opening_Awarded_Points" json:"Opening_Awarded_Points"`
	Opening_Redeemed_Points  float64 `bson:"Opening_Redeemed_Points" json:"Opening_Redeemed_Points"`
	Opening_Available_Points float64 `bson:"Opening_Available_Points" json:"Opening_Available_Points"` //(Awarded_Points + Expired_Points) - Redeemed_Points
	Closure_Awarded_Points   float64 `bson:"Closure_Awarded_Points" json:"Closure_Awarded_Points"`
	Closure_Redeemed_Points  float64 `bson:"Closure_Redeemed_Points" json:"Closure_Redeemed_Points"`
	Closure_Available_Points float64 `bson:"Closure_Available_Points" json:"Closure_Available_Points"` //(Awarded_Points + Expired_Points) - Redeemed_Points

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
