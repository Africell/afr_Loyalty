package Lendme

import "time"

type Loyalty_Event_Log struct {
	Event_User         string      `bson:"Event_User" json:"Event_User"`
	Event_Time         time.Time   `bson:"Event_Time" json:"Event_Time"`
	Event_AffectedType string      `bson:"Event_AffectedType" json:"Event_AffectedType"`
	Event_ActionType   string      `bson:"Event_AffectedType" json:"Event_AffectedType"`
	Event_Description  string      `bson:"Event_Description" json:"Event_Description"`
	Event_Entry_Before interface{} `bson:"Event_Entry_Before" json:"Event_Entry_Before"`
	Event_Entry_After  interface{} `bson:"Event_Entry_After" json:"Event_Entry_After"`
}

type Loyalty_Governance struct {
	Key                             string  `bson:"Key" json:"Key"` //Program name
	Governance_Id                   int64   `bson:"Governance_Id" json:"Governance_Id"`
	Available_Points_Pool           float64 `bson:"Available_Points_Pool" json:"Available_Points_Pool"`
	Distributed_Points_Pool         float64 `bson:"Distributed_Points_Pool" json:"Distributed_Points_Pool"`
	Redeemed_Points_Pool            float64 `bson:"Redeemed_Points_Pool" json:"Redeemed_Points_Pool"`
	MaxAllowedPoints_PerTransaction float64 `bson:"MaxAllowedPoints_PerTransaction" json:"MaxAllowedPoints_PerTransaction"`
	MaxSubsAwardedPoints_PerMonth   float64 `bson:"MaxSubsAwardedPoints_PerMonth" json:"MaxSubsAwardedPoints_PerMonth"`
	MaxSubsAwardedPoints            float64 `bson:"MaxSubsAwardedPoints" json:"MaxSubsAwardedPoints"`
}

type Loyalty_Governance_AddRequest struct {
	Key                             string  `bson:"Key" json:"Key"` //Program name
	Governance_Id                   int64   `bson:"Governance_Id" json:"Governance_Id"`
	Available_Points_Pool           float64 `bson:"Available_Points_Pool" json:"Available_Points_Pool"`
	Distributed_Points_Pool         float64 `bson:"Distributed_Points_Pool" json:"Distributed_Points_Pool"`
	Redeemed_Points_Pool            float64 `bson:"Redeemed_Points_Pool" json:"Redeemed_Points_Pool"`
	MaxAllowedPoints_PerTransaction float64 `bson:"MaxAllowedPoints_PerTransaction" json:"MaxAllowedPoints_PerTransaction"`
	MaxSubsAwardedPoints_PerMonth   float64 `bson:"MaxSubsAwardedPoints_PerMonth" json:"MaxSubsAwardedPoints_PerMonth"`
	MaxSubsAwardedPoints            float64 `bson:"MaxSubsAwardedPoints" json:"MaxSubsAwardedPoints"`
}

type Loyalty_Governance_EditRequest struct {
	Key                             string  `bson:"Key" json:"Key"` //Program name
	NewKey                          string  `bson:"NewKey" json:"NewKey"`
	Governance_Id                   int64   `bson:"Governance_Id" json:"Governance_Id"`
	Available_Points_Pool           float64 `bson:"Available_Points_Pool" json:"Available_Points_Pool"`
	Distributed_Points_Pool         float64 `bson:"Distributed_Points_Pool" json:"Distributed_Points_Pool"`
	Redeemed_Points_Pool            float64 `bson:"Redeemed_Points_Pool" json:"Redeemed_Points_Pool"`
	MaxAllowedPoints_PerTransaction float64 `bson:"MaxAllowedPoints_PerTransaction" json:"MaxAllowedPoints_PerTransaction"`
	MaxSubsAwardedPoints_PerMonth   float64 `bson:"MaxSubsAwardedPoints_PerMonth" json:"MaxSubsAwardedPoints_PerMonth"`
	MaxSubsAwardedPoints            float64 `bson:"MaxSubsAwardedPoints" json:"MaxSubsAwardedPoints"`
}

type Loyalty_Level struct {
	Key                    string  `bson:"Key" json:"Key"`
	Level_Id               int64   `bson:"Level_Id" json:"Level_Id"`
	Description            string  `bson:"Description" json:"Description"`
	Min_Accumulated_Points float64 `bson:"Min_Accumulated_Points" json:"Min_Accumulated_Points"` //will be used for upgrade
	Max_Accumulated_Points float64 `bson:"Max_Accumulated_Points" json:"Max_Accumulated_Points"` //will be used for downgrade
}

type Loyalty_Level_AddRequest struct {
	Key                    string  `bson:"Key" json:"Key"`
	Level_Id               int64   `bson:"Level_Id" json:"Level_Id"`
	Description            string  `bson:"Description" json:"Description"`
	Min_Accumulated_Points float64 `bson:"Min_Accumulated_Points" json:"Min_Accumulated_Points"` //will be used for upgrade
	Max_Accumulated_Points float64 `bson:"Max_Accumulated_Points" json:"Max_Accumulated_Points"` //will be used for downgrade
}

type Loyalty_Level_EditRequest struct {
	Key                    string  `bson:"Key" json:"Key"`
	NewKey                 string  `bson:"NewKey" json:"NewKey"`
	Level_Id               int64   `bson:"Level_Id" json:"Level_Id"`
	Description            string  `bson:"Description" json:"Description"`
	Min_Accumulated_Points float64 `bson:"Min_Accumulated_Points" json:"Min_Accumulated_Points"` //will be used for upgrade
	Max_Accumulated_Points float64 `bson:"Max_Accumulated_Points" json:"Max_Accumulated_Points"` //will be used for downgrade
}

type Loyalty_Account_Segment struct {
	Key         string  `bson:"Key" json:"Key"`
	Segment_Id  int64   `bson:"Segment_Id" json:"Segment_Id"`
	Description string  `bson:"Description" json:"Description"`
	Amount_From float64 `bson:"Amount_From" json:"Amount_From"`
	Amount_Till float64 `bson:"Amount_Till" json:"Amount_Till"`
	AON_From    float64 `bson:"AON_From" json:"AON_From"` //months
	AON_Till    float64 `bson:"AON_Till" json:"AON_Till"` //months
}

type Loyalty_Account_Segment_AddRequest struct {
	Key         string  `bson:"Key" json:"Key"`
	Segment_Id  int64   `bson:"Segment_Id" json:"Segment_Id"`
	Description string  `bson:"Description" json:"Description"`
	Amount_From float64 `bson:"Amount_From" json:"Amount_From"`
	Amount_Till float64 `bson:"Amount_Till" json:"Amount_Till"`
	AON_From    float64 `bson:"AON_From" json:"AON_From"` //months
	AON_Till    float64 `bson:"AON_Till" json:"AON_Till"` //months
}

type Loyalty_Account_Segment_EditRequest struct {
	Key         string  `bson:"Key" json:"Key"`
	NewKey      string  `bson:"NewKey" json:"NewKey"`
	Segment_Id  int64   `bson:"Segment_Id" json:"Segment_Id"`
	Description string  `bson:"Description" json:"Description"`
	Amount_From float64 `bson:"Amount_From" json:"Amount_From"`
	Amount_Till float64 `bson:"Amount_Till" json:"Amount_Till"`
	AON_From    float64 `bson:"AON_From" json:"AON_From"` //months
	AON_Till    float64 `bson:"AON_Till" json:"AON_Till"` //months
}

// AmountPerPoint_BundlePurchase          float64 `bson:"AmountPerPoint_BundlePurchase" json:"AmountPerPoint_BundlePurchase"`
// AmountPerPoint_MOC_Onnet               float64 `bson:"AmountPerPoint_MOC_Onnet" json:"AmountPerPoint_MOC_Onnet"`
// AmountPerPoint_MOC_Offnet              float64 `bson:"AmountPerPoint_MOC_Offnet" json:"AmountPerPoint_MOC_Offnet"`
// AmountPerPoint_MOC_International       float64 `bson:"AmountPerPoint_MOC_International" json:"AmountPerPoint_MOC_International"`
// AmountPerPoint_MOC_Roaming             float64 `bson:"AmountPerPoint_MOC_Roaming" json:"AmountPerPoint_MOC_Roaming"`
// AmountPerPoint_MTC_Roaming             float64 `bson:"AmountPerPoint_MTC_Roaming" json:"AmountPerPoint_MTC_Roaming"`
// AmountPerPoint_MOC_Other               float64 `bson:"AmountPerPoint_MOC_Other" json:"AmountPerPoint_MOC_Other"`
// AmountPerPoint_MO_SMS                  float64 `bson:"AmountPerPoint_MO_SMS" json:"AmountPerPoint_MO_SMS"`
// AmountPerPoint_VAS_Subscriptions       float64 `bson:"AmountPerPoint_VAS_Subscriptions" json:"AmountPerPoint_VAS_Subscriptions"`
// PointsPerTransaction_BundlePurchase    float64 `bson:"PointsPerTransaction_BundlePurchase" json:"PointsPerTransaction_BundlePurchase"`
// PointsPerTransaction_MOC_Onnet         float64 `bson:"PointsPerTransaction_MOC_Onnet" json:"PointsPerTransaction_MOC_Onnet"`
// PointsPerTransaction_MOC_Offnet        float64 `bson:"PointsPerTransaction_MOC_Offnet" json:"PointsPerTransaction_MOC_Offnet"`
// PointsPerTransaction_MOC_International float64 `bson:"PointsPerTransaction_MOC_International" json:"PointsPerTransaction_MOC_International"`
// PointsPerTransaction_MOC_Other         float64 `bson:"PointsPerTransaction_MOC_Other" json:"PointsPerTransaction_MOC_Other"`
// PointsPerTransaction_MOC_Roaming       float64 `bson:"PointsPerTransaction_MOC_Roaming" json:"PointsPerTransaction_MOC_Roaming"`
// PointsPerTransaction_MTC_Roaming       float64 `bson:"PointsPerTransaction_MTC_Roaming" json:"PointsPerTransaction_MTC_Roaming"`
// PointsPerTransaction_MO_SMS            float64 `bson:"PointsPerTransaction_MO_SMS" json:"PointsPerTransaction_MO_SMS"`
// PointsPerTransaction_VAS_Subscriptions float64 `bson:"PointsPerTransaction_VAS_Subscriptions" json:"PointsPerTransaction_VAS_Subscriptions"`

// AmountPerPoint_MM_BundlePurchase  float64 `bson:"AmountPerPoint_MM_BundlePurchase" json:"AmountPerPoint_MM_BundlePurchase"`
// AmountPerPoint_MM_AirtimeRecharge float64 `bson:"AmountPerPoint_MM_AirtimeRecharge" json:"AmountPerPoint_MM_AirtimeRecharge"`
// AmountPerPoint_MM_CashIN          float64 `bson:"AmountPerPoint_MM_CashIN" json:"AmountPerPoint_MM_CashIN"`
// AmountPerPoint_MM_CashOut         float64 `bson:"AmountPerPoint_MM_CashOut" json:"AmountPerPoint_MM_CashOut"`
// AmountPerPoint_MM_P2P             float64 `bson:"AmountPerPoint_MM_P2P" json:"AmountPerPoint_MM_P2P"`
// AmountPerPoint_MM_MerchantPay     float64 `bson:"AmountPerPoint_MM_MerchantPay" json:"AmountPerPoint_MM_MerchantPay"`
// AmountPerPoint_MM_BillPay         float64 `bson:"AmountPerPoint_MM_BillPay" json:"AmountPerPoint_MM_BillPay"`

// PointsPerTransaction_MM_BundlePurchase  float64 `bson:"PointsPerTransaction_MM_BundlePurchase" json:"PointsPerTransaction_MM_BundlePurchase"`
// PointsPerTransaction_MM_AirtimeRecharge float64 `bson:"PointsPerTransaction_MM_AirtimeRecharge" json:"PointsPerTransaction_MM_AirtimeRecharge"`
// PointsPerTransaction_MM_CashIN          float64 `bson:"PointsPerTransaction_MM_CashIN" json:"PointsPerTransaction_MM_CashIN"`
// PointsPerTransaction_MM_CashOut         float64 `bson:"PointsPerTransaction_MM_CashOut" json:"PointsPerTransaction_MM_CashOut"`
// PointsPerTransaction_MM_P2P             float64 `bson:"PointsPerTransaction_MM_P2P" json:"PointsPerTransaction_MM_P2P"`
// PointsPerTransaction_MM_MerchantPay     float64 `bson:"PointsPerTransaction_MM_MerchantPay" json:"PointsPerTransaction_MM_MerchantPay"`
// PointsPerTransaction_MM_BillPay         float64 `bson:"PointsPerTransaction_MM_BillPay" json:"PointsPerTransaction_MM_BillPay"`

type Loyalty_Point_Earning_Rules struct {
	Key                                   string  `bson:"Key" json:"Key"` //Program name
	Earning_Rules_Id                      int64   `bson:"Earning_Rules_Id" json:"Earning_Rules_Id"`
	Description                           string  `bson:"Description" json:"Description"`
	Welcome_Points                        float64 `bson:"Welcome_Points" json:"Welcome_Points"`
	MobileAppDaily_Login                  float64 `bson:"MobileAppDaily_Login" json:"MobileAppDaily_Login"`
	MainGSMBalance_AmountConsumedPerPoint float64 `bson:"MainGSMBalance_AmountConsumedPerPoint" json:"MainGSMBalance_AmountConsumedPerPoint"`
	MobileMoney_AmountConsumedPerPoint    float64 `bson:"MobileMoney_AmountConsumedPerPoint" json:"MobileMoney_AmountConsumedPerPoint"`
}

type Loyalty_Point_Earning_Rules_AddRequest struct {
	Key                                   string  `bson:"Key" json:"Key"` //Program name
	Earning_Rules_Id                      int64   `bson:"Earning_Rules_Id" json:"Earning_Rules_Id"`
	Description                           string  `bson:"Description" json:"Description"`
	Welcome_Points                        float64 `bson:"Welcome_Points" json:"Welcome_Points"`
	MobileAppDaily_Login                  float64 `bson:"MobileAppDaily_Login" json:"MobileAppDaily_Login"`
	MainGSMBalance_AmountConsumedPerPoint float64 `bson:"MainGSMBalance_AmountConsumedPerPoint" json:"MainGSMBalance_AmountConsumedPerPoint"`
	MobileMoney_AmountConsumedPerPoint    float64 `bson:"MobileMoney_AmountConsumedPerPoint" json:"MobileMoney_AmountConsumedPerPoint"`
}

type Loyalty_Point_Earning_Rules_EditRequest struct {
	Key                                   string  `bson:"Key" json:"Key"` //Program name
	NewKey                                string  `bson:"NewKey" json:"NewKey"`
	Earning_Rules_Id                      int64   `bson:"Earning_Rules_Id" json:"Earning_Rules_Id"`
	Description                           string  `bson:"Description" json:"Description"`
	Welcome_Points                        float64 `bson:"Welcome_Points" json:"Welcome_Points"`
	MobileAppDaily_Login                  float64 `bson:"MobileAppDaily_Login" json:"MobileAppDaily_Login"`
	MainGSMBalance_AmountConsumedPerPoint float64 `bson:"MainGSMBalance_AmountConsumedPerPoint" json:"MainGSMBalance_AmountConsumedPerPoint"`
	MobileMoney_AmountConsumedPerPoint    float64 `bson:"MobileMoney_AmountConsumedPerPoint" json:"MobileMoney_AmountConsumedPerPoint"`
}

type Loyalty_Point_Expiry_Rules struct {
	Key               string `bson:"Key" json:"Key"`
	Expiry_Rules_Id   int64  `bson:"Expiry_Rules_Id" json:"Expiry_Rules_Id"`
	Description       string `bson:"Description" json:"Description"`
	Validity_Unit     string `bson:"Validity_Unit" json:"Validity_Unit"` //Monthly, yearly
	Validity_Duration int    `bson:"Validity_Duration" json:"Validity_Duration"`
}

type Loyalty_Point_Expiry_Rules_AddRequest struct {
	Key               string `bson:"Key" json:"Key"`
	Expiry_Rules_Id   int64  `bson:"Expiry_Rules_Id" json:"Expiry_Rules_Id"`
	Description       string `bson:"Description" json:"Description"`
	Validity_Unit     string `bson:"Validity_Unit" json:"Validity_Unit"` //Monthly, yearly
	Validity_Duration int    `bson:"Validity_Duration" json:"Validity_Duration"`
}

type Loyalty_Point_Expiry_Rules_EditRequest struct {
	Key               string `bson:"Key" json:"Key"`
	NewKey            string `bson:"NewKey" json:"NewKey"`
	Expiry_Rules_Id   int64  `bson:"Expiry_Rules_Id" json:"Expiry_Rules_Id"`
	Description       string `bson:"Description" json:"Description"`
	Validity_Unit     string `bson:"Validity_Unit" json:"Validity_Unit"` //Monthly, yearly
	Validity_Duration int    `bson:"Validity_Duration" json:"Validity_Duration"`
}

type Loyalty_Point_Redemption_Rules struct {
	Key                             string  `bson:"Key" json:"Key"`
	Redemption_Rules_Id             int64   `bson:"Redemption_Rules_Id" json:"Redemption_Rules_Id"`
	Description                     string  `bson:"Description" json:"Description"`
	Min_Accumulated_Points          float64 `bson:"Min_Accumulated_Points" json:"Min_Accumulated_Points"` //will be used for upgrade
	Allow_Negative_Balance_ToRedeem bool    `bson:"Allow_Negative_Balance_ToRedeem" json:"Allow_Negative_Balance_ToRedeem"`
	Allow_PendingLendme_ToRedeem    bool    `bson:"Allow_PendingLendme_ToRedeem" json:"Allow_PendingLendme_ToRedeem"`
	Product_Catalogue               string  `bson:"Product_Catalogue" json:"Product_Catalogue"`
}

type Loyalty_Point_Redemption_Rules_AddRequest struct {
	Key                             string  `bson:"Key" json:"Key"`
	Redemption_Rules_Id             int64   `bson:"Redemption_Rules_Id" json:"Redemption_Rules_Id"`
	Description                     string  `bson:"Description" json:"Description"`
	Min_Accumulated_Points          float64 `bson:"Min_Accumulated_Points" json:"Min_Accumulated_Points"` //will be used for upgrade
	Allow_Negative_Balance_ToRedeem bool    `bson:"Allow_Negative_Balance_ToRedeem" json:"Allow_Negative_Balance_ToRedeem"`
	Allow_PendingLendme_ToRedeem    bool    `bson:"Allow_PendingLendme_ToRedeem" json:"Allow_PendingLendme_ToRedeem"`
	Product_Catalogue               string  `bson:"Product_Catalogue" json:"Product_Catalogue"`
}

type Loyalty_Point_Redemption_Rules_EditRequest struct {
	Key                             string  `bson:"Key" json:"Key"`
	NewKey                          string  `bson:"NewKey" json:"NewKey"`
	Redemption_Rules_Id             int64   `bson:"Redemption_Rules_Id" json:"Redemption_Rules_Id"`
	Description                     string  `bson:"Description" json:"Description"`
	Min_Accumulated_Points          float64 `bson:"Min_Accumulated_Points" json:"Min_Accumulated_Points"` //will be used for upgrade
	Allow_Negative_Balance_ToRedeem bool    `bson:"Allow_Negative_Balance_ToRedeem" json:"Allow_Negative_Balance_ToRedeem"`
	Allow_PendingLendme_ToRedeem    bool    `bson:"Allow_PendingLendme_ToRedeem" json:"Allow_PendingLendme_ToRedeem"`
	Product_Catalogue               string  `bson:"Product_Catalogue" json:"Product_Catalogue"`
}

type Loyalty_Plan struct {
	Key                         string `bson:"Key" json:"Key"` //Loyalty_Level_Key + "|" + Loyalty_Account_Segment_Key
	Plan_Id                     int64  `bson:"Plan_Id" json:"Plan_Id"`
	Description                 string `bson:"Description" json:"Description"`
	Loyalty_Level_Key           string `bson:"Loyalty_Level_Key" json:"Loyalty_Level_Key"`
	Loyalty_Account_Segment_Key string `bson:"Loyalty_Account_Segment_Key" json:"Loyalty_Account_Segment_Key"`
	Earning_Rules_Key           string `bson:"Earning_Rules_Key" json:"Earning_Rules_Key"`
	Expiry_Rules_Key            string `bson:"Expiry_Rules_Key" json:"Expiry_Rules_Key"`
	Redemption_Rules_Key        string `bson:"Redemption_Rules_Key" json:"Redemption_Rules_Key"`
}

type Loyalty_Plan_AddRequest struct {
	Key                         string `bson:"Key" json:"Key"` //Loyalty_Level_Key + "|" + Loyalty_Account_Segment_Key
	Plan_Id                     int64  `bson:"Plan_Id" json:"Plan_Id"`
	Description                 string `bson:"Description" json:"Description"`
	Loyalty_Level_Key           string `bson:"Loyalty_Level_Key" json:"Loyalty_Level_Key"`
	Loyalty_Account_Segment_Key string `bson:"Loyalty_Account_Segment_Key" json:"Loyalty_Account_Segment_Key"`
	Earning_Rules_Key           string `bson:"Earning_Rules_Key" json:"Earning_Rules_Key"`
	Expiry_Rules_Key            string `bson:"Expiry_Rules_Key" json:"Expiry_Rules_Key"`
	Redemption_Rules_Key        string `bson:"Redemption_Rules_Key" json:"Redemption_Rules_Key"`
}

type Loyalty_Plan_EditRequest struct {
	Key                         string `bson:"Key" json:"Key"` //Loyalty_Level_Key + "|" + Loyalty_Account_Segment_Key
	NewKey                      string `bson:"NewKey" json:"NewKey"`
	Plan_Id                     int64  `bson:"Plan_Id" json:"Plan_Id"`
	Description                 string `bson:"Description" json:"Description"`
	Loyalty_Level_Key           string `bson:"Loyalty_Level_Key" json:"Loyalty_Level_Key"`
	Loyalty_Account_Segment_Key string `bson:"Loyalty_Account_Segment_Key" json:"Loyalty_Account_Segment_Key"`
	Earning_Rules_Key           string `bson:"Earning_Rules_Key" json:"Earning_Rules_Key"`
	Expiry_Rules_Key            string `bson:"Expiry_Rules_Key" json:"Expiry_Rules_Key"`
	Redemption_Rules_Key        string `bson:"Redemption_Rules_Key" json:"Redemption_Rules_Key"`
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
	Available_Points float64   `bson:"Available_Points" json:"Available_Points"` //Awarded_Points - Redeemed_Points
	Last_Award_Date  time.Time `bson:"Last_Award_Date" json:"Last_Award_Date"`
	Last_Redeem_Date time.Time `bson:"Last_Redeem_Date" json:"Last_Redeem_Date"`

	MainGSMBalance_PendingAmount float64 `bson:"MainGSMBalance_PendingAmount" json:"MainGSMBalance_PendingAmount"`
	MobileMoney_PendingAmount    float64 `bson:"MobileMoney_AmountConsumedPerPoint" json:"MobileMoney_AmountConsumedPerPoint"`

	LoyaltyPointsDetail map[string]Loyalty_Points_Detail `bson:"LoyaltyPointsDetail" json:"LoyaltyPointsDetail"`
}

type Customer_Loyalty_Account_AddRequest struct {
	Key          string    `bson:"Key" json:"Key"` //MSISDN
	Customer_Id  int64     `bson:"Customer_Id" json:"Customer_Id"`
	EventSource  string    `bson:"EventSource" json:"EventSource"` //MSISDN
	COS          string    `bson:"COS" json:"COS"`
	ARPU         float64   `bson:"ARPU" json:"ARPU"`
	Joining_Date time.Time `bson:"Joining_Date" json:"Joining_Date"`
}

type Customer_Loyalty_Account_EditRequest struct {
	Key               string    `bson:"Key" json:"Key"` //MSISDN
	NewKey            string    `bson:"NewKey" json:"NewKey"`
	Customer_Id       int64     `bson:"Customer_Id" json:"Customer_Id"`
	Loyalty_Level_Key string    `bson:"Loyalty_Level_Key" json:"Loyalty_Level_Key"`
	COS               string    `bson:"COS" json:"COS"`
	ARPU              float64   `bson:"ARPU" json:"ARPU"`
	Joining_Date      time.Time `bson:"Joining_Date" json:"Joining_Date"`
}

type Loyalty_Points_Detail struct {
	Year_Month       string  `bson:"Year_Month" json:"Year_Month"`
	Awarded_Points   float64 `bson:"Awarded_Points" json:"Awarded_Points"`
	Redeemed_Points  float64 `bson:"Redeemed_Points" json:"Redeemed_Points"`
	Available_Points float64 `bson:"Available_Points" json:"Available_Points"` //Awarded_Points - Redeemed_Points
}

type Customer_Loyalty_Account_AwardRequest struct {
	MSISDN      string  `bson:"MSISDN" json:"MSISDN"`           //MSISDN
	EventSource string  `bson:"EventSource" json:"EventSource"` //MobileApp, MobileMoney, USSD,...
	EventType   string  `bson:"EventType" json:"EventType"`     //BundlePurchase, MOC,...
	EventDetail string  `bson:"EventDetail" json:"EventDetail"` //BundleName, ...
	Amount      float64 `bson:"Amount" json:"Amount"`
}

type Customer_Loyalty_Account_DeleteRequest struct {
	Key         string `bson:"Key" json:"Key"` //MSISDN
	Customer_Id int64  `bson:"Customer_Id" json:"Customer_Id"`
	EventSource string `bson:"EventSource" json:"EventSource"` //MSISDN
}

type Customer_UAT struct {
	Key       string    `bson:"Key" json:"Key"` //Key is the MSISDN
	Id        int64     `bson:"Id" json:"Id"`
	AddTime   time.Time `bson:"AddTime" json:"AddTime"`
	AddReason string    `bson:"AddReason" json:"AddReason"`
}

type Customer_UAT_AddRequest struct {
	Key       string    `bson:"Key" json:"Key"` //Key is the MSISDN
	Id        int64     `bson:"Id" json:"Id"`
	AddTime   time.Time `bson:"AddTime" json:"AddTime"`
	AddReason string    `bson:"AddReason" json:"AddReason"`
}

type Customer_UAT_EditRequest struct {
	Key       string    `bson:"Key" json:"Key"` //Key is the MSISDN
	NewKey    string    `bson:"NewKey" json:"NewKey"`
	Id        int64     `bson:"Id" json:"Id"`
	AddTime   time.Time `bson:"AddTime" json:"AddTime"`
	AddReason string    `bson:"AddReason" json:"AddReason"`
}

type Customer_DND struct {
	Key       string    `bson:"Key" json:"Key"` //Key is the MSISDN
	Id        int64     `bson:"Id" json:"Id"`
	AddTime   time.Time `bson:"AddTime" json:"AddTime"`
	AddReason string    `bson:"AddReason" json:"AddReason"`
}

type Customer_DND_AddRequest struct {
	Key       string    `bson:"Key" json:"Key"` //Key is the MSISDN
	Id        int64     `bson:"Id" json:"Id"`
	AddTime   time.Time `bson:"AddTime" json:"AddTime"`
	AddReason string    `bson:"AddReason" json:"AddReason"`
}

type Customer_DND_EditRequest struct {
	Key       string    `bson:"Key" json:"Key"` //Key is the MSISDN
	NewKey    string    `bson:"NewKey" json:"NewKey"`
	Id        int64     `bson:"Id" json:"Id"`
	AddTime   time.Time `bson:"AddTime" json:"AddTime"`
	AddReason string    `bson:"AddReason" json:"AddReason"`
}

type Customer_Exclusion struct {
	Key       string    `bson:"Key" json:"Key"` //Key is the MSISDN
	Id        int64     `bson:"Id" json:"Id"`
	AddTime   time.Time `bson:"AddTime" json:"AddTime"`
	AddReason string    `bson:"AddReason" json:"AddReason"`
}

type Customer_Exclusion_AddRequest struct {
	Key       string    `bson:"Key" json:"Key"` //Key is the MSISDN
	Id        int64     `bson:"Id" json:"Id"`
	AddTime   time.Time `bson:"AddTime" json:"AddTime"`
	AddReason string    `bson:"AddReason" json:"AddReason"`
}

type Customer_Exclusion_EditRequest struct {
	Key       string    `bson:"Key" json:"Key"` //Key is the MSISDN
	NewKey    string    `bson:"NewKey" json:"NewKey"`
	Id        int64     `bson:"Id" json:"Id"`
	AddTime   time.Time `bson:"AddTime" json:"AddTime"`
	AddReason string    `bson:"AddReason" json:"AddReason"`
}

type Customer_COS_Exclusion struct {
	Key       string    `bson:"Key" json:"Key"` //Key is the COS Id
	Id        int64     `bson:"Id" json:"Id"`
	AddTime   time.Time `bson:"AddTime" json:"AddTime"`
	AddReason string    `bson:"AddReason" json:"AddReason"`
}

type Customer_COS_Exclusion_AddRequest struct {
	Key       string    `bson:"Key" json:"Key"` //Key is the COS Id
	Id        int64     `bson:"Id" json:"Id"`
	AddTime   time.Time `bson:"AddTime" json:"AddTime"`
	AddReason string    `bson:"AddReason" json:"AddReason"`
}

type Customer_COS_Exclusion_EditRequest struct {
	Key       string    `bson:"Key" json:"Key"` //Key is the COS Id
	NewKey    string    `bson:"NewKey" json:"NewKey"`
	Id        int64     `bson:"Id" json:"Id"`
	AddTime   time.Time `bson:"AddTime" json:"AddTime"`
	AddReason string    `bson:"AddReason" json:"AddReason"`
}
