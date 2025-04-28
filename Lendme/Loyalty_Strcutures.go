package Lendme

import (
	"daoc"
	"time"
)

type Loyalty_Event_Log struct {
	Event_User         string      `bson:"Event_User" json:"Event_User"`
	Event_Time         time.Time   `bson:"Event_Time" json:"Event_Time"`
	Event_AffectedType string      `bson:"Event_AffectedType" json:"Event_AffectedType"`
	Event_ActionType   string      `bson:"Event_ActionType" json:"Event_ActionType"`
	Event_Description  string      `bson:"Event_Description" json:"Event_Description"`
	Event_Entry_Before interface{} `bson:"Event_Entry_Before" json:"Event_Entry_Before"`
	Event_Entry_After  interface{} `bson:"Event_Entry_After" json:"Event_Entry_After"`
}

type Request_Header struct {
	SourceIP              string        `bson:"SourceIP" json:"SourceIP"`
	SourceApp             string        `bson:"SourceApp" json:"SourceApp"`
	AppLogin              string        `bson:"AppLogin" json:"AppLogin"`
	HostId                string        `bson:"HostId" json:"HostId"`
	AppVersion            string        `bson:"AppVersion" json:"AppVersion"`
	GPSLocation           daoc.Location `bson:"GPSLocation" json:"GPSLocation"`
	GSMLocation           string        `bson:"GSMLocation" json:"GSMLocation"`
	Authorization         string        `bson:"Authorization" json:"Authorization"`
	IsValid               bool          `bson:"IsValid" json:"IsValid"`
	ValidationDescription string        `bson:"ValidationDescription" json:"ValidationDescription"`
}

type Loyalty_Governance struct {
	Key                             string  `bson:"Key" json:"Key"` //Program name
	Governance_Id                   int64   `bson:"Governance_Id" json:"Governance_Id"`
	Available_Points_Pool           float64 `bson:"Available_Points_Pool" json:"Available_Points_Pool"`
	Distributed_Points_Pool         float64 `bson:"Distributed_Points_Pool" json:"Distributed_Points_Pool"`
	Redeemed_Points_Pool            float64 `bson:"Redeemed_Points_Pool" json:"Redeemed_Points_Pool"`
	Expired_Points_Pool             float64 `bson:"Expired_Points_Pool" json:"Expired_Points_Pool"`
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
	Expired_Points_Pool             float64 `bson:"Expired_Points_Pool" json:"Expired_Points_Pool"`
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
	Expired_Points_Pool             float64 `bson:"Expired_Points_Pool" json:"Expired_Points_Pool"`
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
	EnableRedeem           bool    `bson:"EnableRedeem" json:"EnableRedeem"`
	DowngradeToLevel_Key   string  `bson:"DowngradeToLevel_Key" json:"DowngradeToLevel_Key"`
}

type Loyalty_Level_AddRequest struct {
	Key                    string  `bson:"Key" json:"Key"`
	Level_Id               int64   `bson:"Level_Id" json:"Level_Id"`
	Description            string  `bson:"Description" json:"Description"`
	Min_Accumulated_Points float64 `bson:"Min_Accumulated_Points" json:"Min_Accumulated_Points"` //will be used for upgrade
	Max_Accumulated_Points float64 `bson:"Max_Accumulated_Points" json:"Max_Accumulated_Points"` //will be used for downgrade
	EnableRedeem           bool    `bson:"EnableRedeem" json:"EnableRedeem"`
	DowngradeToLevel_Key   string  `bson:"DowngradeToLevel_Key" json:"DowngradeToLevel_Key"`
}

type Loyalty_Level_EditRequest struct {
	Key                    string  `bson:"Key" json:"Key"`
	NewKey                 string  `bson:"NewKey" json:"NewKey"`
	Level_Id               int64   `bson:"Level_Id" json:"Level_Id"`
	Description            string  `bson:"Description" json:"Description"`
	Min_Accumulated_Points float64 `bson:"Min_Accumulated_Points" json:"Min_Accumulated_Points"` //will be used for upgrade
	Max_Accumulated_Points float64 `bson:"Max_Accumulated_Points" json:"Max_Accumulated_Points"` //will be used for downgrade
	EnableRedeem           bool    `bson:"EnableRedeem" json:"EnableRedeem"`
	DowngradeToLevel_Key   string  `bson:"DowngradeToLevel_Key" json:"DowngradeToLevel_Key"`
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
	MM_P2P_Award_Type                     string  `bson:"MM_P2P_Award_Type" json:"MM_P2P_Award_Type"` //"Transaction" or "Amount"
	MM_P2P                                float64 `bson:"MM_P2P" json:"MM_P2P"`
	MM_CASHIN_Award_Type                  string  `bson:"MM_CASHIN_Award_Type" json:"MM_CASHIN_Award_Type"` //"Transaction" or "Amount"
	MM_CASHIN                             float64 `bson:"MM_CASHIN" json:"MM_CASHIN"`
	MM_CASHOUT_Award_Type                 string  `bson:"MM_CASHOUT_Award_Type" json:"MM_CASHOUT_Award_Type"` //"Transaction" or "Amount"
	MM_CASHOUT                            float64 `bson:"MM_CASHOUT" json:"MM_CASHOUT"`
	MM_MERCHPAY_Award_Type                string  `bson:"MM_MERCHPAY_Award_Type" json:"MM_MERCHPAY_Award_Type"` //"Transaction" or "Amount"
	MM_MERCHPAY                           float64 `bson:"MM_MERCHPAY" json:"MM_MERCHPAY"`
	MM_BILLPAY_Award_Type                 string  `bson:"MM_BILLPAY_Award_Type" json:"MM_BILLPAY_Award_Type"` //"Transaction" or "Amount"
	MM_BILLPAY                            float64 `bson:"MM_BILLPAY" json:"MM_BILLPAY"`
	MM_RC_Award_Type                      string  `bson:"MM_RC_Award_Type" json:"MM_RC_Award_Type"` //"Transaction" or "Amount"
	MM_RC                                 float64 `bson:"MM_RC" json:"MM_RC"`
	MM_CTMMOREQ_Award_Type                string  `bson:"MM_CTMMOREQ_Award_Type" json:"MM_CTMMOREQ_Award_Type"` //"Transaction" or "Amount"
	MM_CTMMOREQ                           float64 `bson:"MM_CTMMOREQ" json:"MM_CTMMOREQ"`
	MM_CBWREQ_Award_Type                  string  `bson:"MM_CBWREQ_Award_Type" json:"MM_CBWREQ_Award_Type"` //"Transaction" or "Amount"
	MM_CBWREQ                             float64 `bson:"MM_CBWREQ" json:"MM_CBWREQ"`
}

type Loyalty_Point_Earning_Rules_AddRequest struct {
	Key                                   string  `bson:"Key" json:"Key"` //Program name
	Earning_Rules_Id                      int64   `bson:"Earning_Rules_Id" json:"Earning_Rules_Id"`
	Description                           string  `bson:"Description" json:"Description"`
	Welcome_Points                        float64 `bson:"Welcome_Points" json:"Welcome_Points"`
	MobileAppDaily_Login                  float64 `bson:"MobileAppDaily_Login" json:"MobileAppDaily_Login"`
	MainGSMBalance_AmountConsumedPerPoint float64 `bson:"MainGSMBalance_AmountConsumedPerPoint" json:"MainGSMBalance_AmountConsumedPerPoint"`
	MM_P2P_Award_Type                     string  `bson:"MM_P2P_Award_Type" json:"MM_P2P_Award_Type"` //"Transaction" or "Amount"
	MM_P2P                                float64 `bson:"MM_P2P" json:"MM_P2P"`
	MM_CASHIN_Award_Type                  string  `bson:"MM_CASHIN_Award_Type" json:"MM_CASHIN_Award_Type"` //"Transaction" or "Amount"
	MM_CASHIN                             float64 `bson:"MM_CASHIN" json:"MM_CASHIN"`
	MM_CASHOUT_Award_Type                 string  `bson:"MM_CASHOUT_Award_Type" json:"MM_CASHOUT_Award_Type"` //"Transaction" or "Amount"
	MM_CASHOUT                            float64 `bson:"MM_CASHOUT" json:"MM_CASHOUT"`
	MM_MERCHPAY_Award_Type                string  `bson:"MM_MERCHPAY_Award_Type" json:"MM_MERCHPAY_Award_Type"` //"Transaction" or "Amount"
	MM_MERCHPAY                           float64 `bson:"MM_MERCHPAY" json:"MM_MERCHPAY"`
	MM_BILLPAY_Award_Type                 string  `bson:"MM_BILLPAY_Award_Type" json:"MM_BILLPAY_Award_Type"` //"Transaction" or "Amount"
	MM_BILLPAY                            float64 `bson:"MM_BILLPAY" json:"MM_BILLPAY"`
	MM_RC_Award_Type                      string  `bson:"MM_RC_Award_Type" json:"MM_RC_Award_Type"` //"Transaction" or "Amount"
	MM_RC                                 float64 `bson:"MM_RC" json:"MM_RC"`
	MM_CTMMOREQ_Award_Type                string  `bson:"MM_CTMMOREQ_Award_Type" json:"MM_CTMMOREQ_Award_Type"` //"Transaction" or "Amount"
	MM_CTMMOREQ                           float64 `bson:"MM_CTMMOREQ" json:"MM_CTMMOREQ"`
	MM_CBWREQ_Award_Type                  string  `bson:"MM_CBWREQ_Award_Type" json:"MM_CBWREQ_Award_Type"` //"Transaction" or "Amount"
	MM_CBWREQ                             float64 `bson:"MM_CBWREQ" json:"MM_CBWREQ"`
}

type Loyalty_Point_Earning_Rules_EditRequest struct {
	Key                                   string  `bson:"Key" json:"Key"` //Program name
	NewKey                                string  `bson:"NewKey" json:"NewKey"`
	Earning_Rules_Id                      int64   `bson:"Earning_Rules_Id" json:"Earning_Rules_Id"`
	Description                           string  `bson:"Description" json:"Description"`
	Welcome_Points                        float64 `bson:"Welcome_Points" json:"Welcome_Points"`
	MobileAppDaily_Login                  float64 `bson:"MobileAppDaily_Login" json:"MobileAppDaily_Login"`
	MainGSMBalance_AmountConsumedPerPoint float64 `bson:"MainGSMBalance_AmountConsumedPerPoint" json:"MainGSMBalance_AmountConsumedPerPoint"`
	MM_P2P_Award_Type                     string  `bson:"MM_P2P_Award_Type" json:"MM_P2P_Award_Type"` //"Transaction" or "Amount"
	MM_P2P                                float64 `bson:"MM_P2P" json:"MM_P2P"`
	MM_CASHIN_Award_Type                  string  `bson:"MM_CASHIN_Award_Type" json:"MM_CASHIN_Award_Type"` //"Transaction" or "Amount"
	MM_CASHIN                             float64 `bson:"MM_CASHIN" json:"MM_CASHIN"`
	MM_CASHOUT_Award_Type                 string  `bson:"MM_CASHOUT_Award_Type" json:"MM_CASHOUT_Award_Type"` //"Transaction" or "Amount"
	MM_CASHOUT                            float64 `bson:"MM_CASHOUT" json:"MM_CASHOUT"`
	MM_MERCHPAY_Award_Type                string  `bson:"MM_MERCHPAY_Award_Type" json:"MM_MERCHPAY_Award_Type"` //"Transaction" or "Amount"
	MM_MERCHPAY                           float64 `bson:"MM_MERCHPAY" json:"MM_MERCHPAY"`
	MM_BILLPAY_Award_Type                 string  `bson:"MM_BILLPAY_Award_Type" json:"MM_BILLPAY_Award_Type"` //"Transaction" or "Amount"
	MM_BILLPAY                            float64 `bson:"MM_BILLPAY" json:"MM_BILLPAY"`
	MM_RC_Award_Type                      string  `bson:"MM_RC_Award_Type" json:"MM_RC_Award_Type"` //"Transaction" or "Amount"
	MM_RC                                 float64 `bson:"MM_RC" json:"MM_RC"`
	MM_CTMMOREQ_Award_Type                string  `bson:"MM_CTMMOREQ_Award_Type" json:"MM_CTMMOREQ_Award_Type"` //"Transaction" or "Amount"
	MM_CTMMOREQ                           float64 `bson:"MM_CTMMOREQ" json:"MM_CTMMOREQ"`
	MM_CBWREQ_Award_Type                  string  `bson:"MM_CBWREQ_Award_Type" json:"MM_CBWREQ_Award_Type"` //"Transaction" or "Amount"
	MM_CBWREQ                             float64 `bson:"MM_CBWREQ" json:"MM_CBWREQ"`
}

type Loyalty_Point_Expiry_Rules struct {
	Key                     string    `bson:"Key" json:"Key"`
	Expiry_Rules_Id         int64     `bson:"Expiry_Rules_Id" json:"Expiry_Rules_Id"`
	Description             string    `bson:"Description" json:"Description"`
	Rolling_Expiration      bool      `bson:"Rolling_Expiration" json:"Rolling_Expiration"`
	Validity_Unit           string    `bson:"Validity_Unit" json:"Validity_Unit"`         //Month, Year --> only when Rolling_Expiration is true
	Validity_Duration       int       `bson:"Validity_Duration" json:"Validity_Duration"` //only when Rolling_Expiration is true
	Fix_Date_Expiration     bool      `bson:"Fix_Date_Expiration" json:"Fix_Date_Expiration"`
	Expiration_Trigger_date time.Time `bson:"Expiration_Trigger_date" json:"Expiration_Trigger_date"` //when the expiry process will run
	Expiration_Point_Before time.Time `bson:"Expiration_Point_Before" json:"Expiration_Point_Before"` //expiry all points before this date
}

type Loyalty_Point_Expiry_Rules_AddRequest struct {
	Key                     string    `bson:"Key" json:"Key"`
	Expiry_Rules_Id         int64     `bson:"Expiry_Rules_Id" json:"Expiry_Rules_Id"`
	Description             string    `bson:"Description" json:"Description"`
	Rolling_Expiration      bool      `bson:"Rolling_Expiration" json:"Rolling_Expiration"`
	Validity_Unit           string    `bson:"Validity_Unit" json:"Validity_Unit"`         //Month, Year --> only when Rolling_Expiration is true
	Validity_Duration       int       `bson:"Validity_Duration" json:"Validity_Duration"` //only when Rolling_Expiration is true
	Fix_Date_Expiration     bool      `bson:"Fix_Date_Expiration" json:"Fix_Date_Expiration"`
	Expiration_Trigger_date time.Time `bson:"Expiration_Trigger_date" json:"Expiration_Trigger_date"` //when the expiry process will run
	Expiration_Point_Before time.Time `bson:"Expiration_Point_Before" json:"Expiration_Point_Before"` //expiry all points before this date
}

type Loyalty_Point_Expiry_Rules_EditRequest struct {
	Key                     string    `bson:"Key" json:"Key"`
	NewKey                  string    `bson:"NewKey" json:"NewKey"`
	Expiry_Rules_Id         int64     `bson:"Expiry_Rules_Id" json:"Expiry_Rules_Id"`
	Description             string    `bson:"Description" json:"Description"`
	Rolling_Expiration      bool      `bson:"Rolling_Expiration" json:"Rolling_Expiration"`
	Validity_Unit           string    `bson:"Validity_Unit" json:"Validity_Unit"`         //Month, Year --> only when Rolling_Expiration is true
	Validity_Duration       int       `bson:"Validity_Duration" json:"Validity_Duration"` //only when Rolling_Expiration is true
	Fix_Date_Expiration     bool      `bson:"Fix_Date_Expiration" json:"Fix_Date_Expiration"`
	Expiration_Trigger_date time.Time `bson:"Expiration_Trigger_date" json:"Expiration_Trigger_date"` //when the expiry process will run
	Expiration_Point_Before time.Time `bson:"Expiration_Point_Before" json:"Expiration_Point_Before"` //expiry all points before this date
}

type Loyalty_Point_Redemption_Rules struct {
	Key                               string  `bson:"Key" json:"Key"`
	Redemption_Rules_Id               int64   `bson:"Redemption_Rules_Id" json:"Redemption_Rules_Id"`
	Description                       string  `bson:"Description" json:"Description"`
	Min_Accumulated_Points            float64 `bson:"Min_Accumulated_Points" json:"Min_Accumulated_Points"` //will be used for upgrade
	Allow_Negative_Balance_ToRedeem   bool    `bson:"Allow_Negative_Balance_ToRedeem" json:"Allow_Negative_Balance_ToRedeem"`
	Allow_PendingLendme_ToRedeem      bool    `bson:"Allow_PendingLendme_ToRedeem" json:"Allow_PendingLendme_ToRedeem"`
	Airtime_MinPoints                 float64 `bson:"Airtime_MinPoints" json:"Airtime_MinPoints"`
	Airtime_AmountPerPoint            float64 `bson:"Airtime_AmountPerPoint" json:"Airtime_AmountPerPoint"`
	Airtime_EVC_Account               string  `bson:"Airtime_EVC_Account" json:"Airtime_EVC_Account"`
	Airtime_EVC_PIN                   string  `bson:"Airtime_EVC_PIN" json:"Airtime_EVC_PIN"`
	MobileMoney_MinPoints             float64 `bson:"MobileMoney_MinPoints" json:"MobileMoney_MinPoints"`
	MobileMoney_AmountPerPoint        float64 `bson:"MobileMoney_AmountPerPoint" json:"MobileMoney_AmountPerPoint"`
	MobileMoney_MerchantAccount       string  `bson:"MobileMoney_MerchantAccount" json:"MobileMoney_MerchantAccount"`
	MobileMoney_MerchantPIN           string  `bson:"MobileMoney_MerchantPIN" json:"MobileMoney_MerchantPIN"`
	Bundles_MinPoints                 float64 `bson:"Bundles_MinPoints" json:"Bundles_MinPoints"`
	Bundles_Product_Catalogue_Channel string  `bson:"Bundles_Product_Catalogue_Channel" json:"Bundles_Product_Catalogue_Channel"`
	Bundles_Product_Catalogue_Plan    string  `bson:"Bundles_Product_Catalogue_Plan" json:"Bundles_Product_Catalogue_Plan"`
	Bundles_Product_Catalogue_Version string  `bson:"Bundles_Product_Catalogue_Version" json:"Bundles_Product_Catalogue_Version"`
	Bundles_EVC_Account               string  `bson:"Bundles_EVC_Account" json:"Bundles_EVC_Account"`
	Bundles_EVC_PIN                   string  `bson:"Bundles_EVC_PIN" json:"Bundles_EVC_PIN"`
	FreeSpinAndWin_MinPoints          float64 `bson:"FreeSpinAndWin_MinPoints" json:"FreeSpinAndWin_MinPoints"`
	FreeSpinAndWin_PointsPerSpin      float64 `bson:"FreeSpinAndWin_PointsPerSpin" json:"FreeSpinAndWin_PointsPerSpin"`
}

type Loyalty_Point_Redemption_Rules_AddRequest struct {
	Key                               string  `bson:"Key" json:"Key"`
	Redemption_Rules_Id               int64   `bson:"Redemption_Rules_Id" json:"Redemption_Rules_Id"`
	Description                       string  `bson:"Description" json:"Description"`
	Min_Accumulated_Points            float64 `bson:"Min_Accumulated_Points" json:"Min_Accumulated_Points"` //will be used for upgrade
	Allow_Negative_Balance_ToRedeem   bool    `bson:"Allow_Negative_Balance_ToRedeem" json:"Allow_Negative_Balance_ToRedeem"`
	Allow_PendingLendme_ToRedeem      bool    `bson:"Allow_PendingLendme_ToRedeem" json:"Allow_PendingLendme_ToRedeem"`
	Airtime_MinPoints                 float64 `bson:"Airtime_MinPoints" json:"Airtime_MinPoints"`
	Airtime_AmountPerPoint            float64 `bson:"Airtime_AmountPerPoint" json:"Airtime_AmountPerPoint"`
	Airtime_EVC_Account               string  `bson:"Airtime_EVC_Account" json:"Airtime_EVC_Account"`
	Airtime_EVC_PIN                   string  `bson:"Airtime_EVC_PIN" json:"Airtime_EVC_PIN"`
	MobileMoney_MinPoints             float64 `bson:"MobileMoney_MinPoints" json:"MobileMoney_MinPoints"`
	MobileMoney_AmountPerPoint        float64 `bson:"MobileMoney_AmountPerPoint" json:"MobileMoney_AmountPerPoint"`
	Bundles_MinPoints                 float64 `bson:"Bundles_MinPoints" json:"Bundles_MinPoints"`
	Bundles_EVC_Account               string  `bson:"Bundles_EVC_Account" json:"Bundles_EVC_Account"`
	Bundles_EVC_PIN                   string  `bson:"Bundles_EVC_PIN" json:"Bundles_EVC_PIN"`
	Bundles_Product_Catalogue_Channel string  `bson:"Bundles_Product_Catalogue_Channel" json:"Bundles_Product_Catalogue_Channel"`
	Bundles_Product_Catalogue_Plan    string  `bson:"Bundles_Product_Catalogue_Plan" json:"Bundles_Product_Catalogue_Plan"`
	Bundles_Product_Catalogue_Version string  `bson:"Bundles_Product_Catalogue_Version" json:"Bundles_Product_Catalogue_Version"`
	FreeSpinAndWin_MinPoints          float64 `bson:"FreeSpinAndWin_MinPoints" json:"FreeSpinAndWin_MinPoints"`
	FreeSpinAndWin_PointsPerSpin      float64 `bson:"FreeSpinAndWin_PointsPerSpin" json:"FreeSpinAndWin_PointsPerSpin"`
}

type Loyalty_Point_Redemption_Rules_EditRequest struct {
	Key                               string  `bson:"Key" json:"Key"`
	NewKey                            string  `bson:"NewKey" json:"NewKey"`
	Redemption_Rules_Id               int64   `bson:"Redemption_Rules_Id" json:"Redemption_Rules_Id"`
	Description                       string  `bson:"Description" json:"Description"`
	Min_Accumulated_Points            float64 `bson:"Min_Accumulated_Points" json:"Min_Accumulated_Points"` //will be used for upgrade
	Allow_Negative_Balance_ToRedeem   bool    `bson:"Allow_Negative_Balance_ToRedeem" json:"Allow_Negative_Balance_ToRedeem"`
	Allow_PendingLendme_ToRedeem      bool    `bson:"Allow_PendingLendme_ToRedeem" json:"Allow_PendingLendme_ToRedeem"`
	Airtime_MinPoints                 float64 `bson:"Airtime_MinPoints" json:"Airtime_MinPoints"`
	Airtime_AmountPerPoint            float64 `bson:"Airtime_AmountPerPoint" json:"Airtime_AmountPerPoint"`
	Airtime_EVC_Account               string  `bson:"Airtime_EVC_Account" json:"Airtime_EVC_Account"`
	Airtime_EVC_PIN                   string  `bson:"Airtime_EVC_PIN" json:"Airtime_EVC_PIN"`
	MobileMoney_MinPoints             float64 `bson:"MobileMoney_MinPoints" json:"MobileMoney_MinPoints"`
	MobileMoney_AmountPerPoint        float64 `bson:"MobileMoney_AmountPerPoint" json:"MobileMoney_AmountPerPoint"`
	Bundles_MinPoints                 float64 `bson:"Bundles_MinPoints" json:"Bundles_MinPoints"`
	Bundles_Product_Catalogue_Channel string  `bson:"Bundles_Product_Catalogue_Channel" json:"Bundles_Product_Catalogue_Channel"`
	Bundles_Product_Catalogue_Plan    string  `bson:"Bundles_Product_Catalogue_Plan" json:"Bundles_Product_Catalogue_Plan"`
	Bundles_Product_Catalogue_Version string  `bson:"Bundles_Product_Catalogue_Version" json:"Bundles_Product_Catalogue_Version"`
	Bundles_EVC_Account               string  `bson:"Bundles_EVC_Account" json:"Bundles_EVC_Account"`
	Bundles_EVC_PIN                   string  `bson:"Bundles_EVC_PIN" json:"Bundles_EVC_PIN"`
	FreeSpinAndWin_MinPoints          float64 `bson:"FreeSpinAndWin_MinPoints" json:"FreeSpinAndWin_MinPoints"`
	FreeSpinAndWin_PointsPerSpin      float64 `bson:"FreeSpinAndWin_PointsPerSpin" json:"FreeSpinAndWin_PointsPerSpin"`
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

	Expired_Points float64   `bson:"Expired_Points" json:"Expired_Points"` //expired are deducted from Awarded_Points
	Expiry_Date    time.Time `bson:"Expiry_Date" json:"Expiry_Date"`

	MainGSMBalance_PendingAmount float64 `bson:"MainGSMBalance_PendingAmount" json:"MainGSMBalance_PendingAmount"`
	MobileMoney_PendingAmount    float64 `bson:"MobileMoney_AmountConsumedPerPoint" json:"MobileMoney_AmountConsumedPerPoint"`

	Points_Detail_Keys []string `bson:"Points_Detail_Keys" json:"Points_Detail_Keys"`
}

type Customer_Loyalty_Account_AddRequest struct {
	Key          string    `bson:"Key" json:"Key"` //MSISDN
	Customer_Id  int64     `bson:"Customer_Id" json:"Customer_Id"`
	EventSource  string    `bson:"EventSource" json:"EventSource"`
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

type Customer_Loyalty_Account_Points_Detail struct {
	Key              string    `bson:"Key" json:"Key"` //MSISDN + "|"+Year_Month
	Year_Month       string    `bson:"Year_Month" json:"Year_Month"`
	Creation_date    time.Time `bson:"Creation_date" json:"Creation_date"`
	Awarded_Points   float64   `bson:"Awarded_Points" json:"Awarded_Points"`
	Redeemed_Points  float64   `bson:"Redeemed_Points" json:"Redeemed_Points"`
	Available_Points float64   `bson:"Available_Points" json:"Available_Points"` //(Awarded_Points + Expired_Points) - Redeemed_Points
	Expired_Points   float64   `bson:"Expired_Points" json:"Expired_Points"`     //expired are deducted from Awarded_Points
	Expiry_Date      time.Time `bson:"Expiry_Date" json:"Expiry_Date"`
	Last_Redeem_Date time.Time `bson:"Last_Redeem_Date" json:"Last_Redeem_Date"`
	Last_Credit_Date time.Time `bson:"Last_Credit_Date" json:"Last_Credit_Date"`
}

type Loyalty_AccountCreditPoints_Request struct {
	MSISDN           string  `bson:"MSISDN" json:"MSISDN"`           //MSISDN
	EventSource      string  `bson:"EventSource" json:"EventSource"` //MobileApp, MobileMoney, USSD,...
	EventType        string  `bson:"EventType" json:"EventType"`     //BundlePurchase, MOC,...
	EventDetail      string  `bson:"EventDetail" json:"EventDetail"` //BundleName, ...
	EventAmount      float64 `bson:"EventAmount" json:"EventAmount"`
	PointsToCredit   float64 `bson:"PointsToCredit" json:"PointsToCredit"`
	EventDescription string  `bson:"EventDescription" json:"EventDescription"`
}

type Loyalty_AccountCreditPoints_log struct {
	//request Header info
	SourceIP    string        `bson:"SourceIP" json:"-"`
	SourceApp   string        `bson:"SourceApp" json:"-"`
	AppLogin    string        `bson:"AppLogin" json:"-"`
	AppVersion  string        `bson:"AppVersion" json:"-"`
	GPSLocation daoc.Location `bson:"GPSLocation" json:"-"`
	GSMLocation string        `bson:"GSMLocation" json:"-"`

	//request detail
	MSISDN           string  `bson:"MSISDN" json:"MSISDN"`           //MSISDN
	EventSource      string  `bson:"EventSource" json:"EventSource"` //MobileApp, MobileMoney, USSD,...
	EventType        string  `bson:"EventType" json:"EventType"`     //BundlePurchase, MOC,...
	EventDetail      string  `bson:"EventDetail" json:"EventDetail"` //BundleName, ...
	EventAmount      float64 `bson:"EventAmount" json:"EventAmount"`
	PointsToCredit   float64 `bson:"PointsToCredit" json:"PointsToCredit"`
	EventDescription string  `bson:"EventDescription" json:"EventDescription"`

	Opening_Loyalty_Level_Key           string `bson:"Opening_Loyalty_Level_Key" json:"Opening_Loyalty_Level_Key"`
	Opening_Loyalty_Account_Segment_Key string `bson:"Opening_Loyalty_Account_Segment_Key" json:"Opening_Loyalty_Account_Segment_Key"`
	Closure_Loyalty_Level_Key           string `bson:"Closure_Loyalty_Level_Key" json:"Closure_Loyalty_Level_Key"`
	Closure_Loyalty_Account_Segment_Key string `bson:"Closure_Loyalty_Account_Segment_Key" json:"Closure_Loyalty_Account_Segment_Key"`

	AwardedPoints            float64 `bson:"AwardedPoints" json:"AwardedPoints"`
	Opening_Awarded_Points   float64 `bson:"Opening_Awarded_Points" json:"Opening_Awarded_Points"`
	Opening_Redeemed_Points  float64 `bson:"Opening_Redeemed_Points" json:"Opening_Redeemed_Points"`
	Opening_Available_Points float64 `bson:"Opening_Available_Points" json:"Opening_Available_Points"` //(Awarded_Points + Expired_Points) - Redeemed_Points
	Closure_Awarded_Points   float64 `bson:"Closure_Awarded_Points" json:"Closure_Awarded_Points"`
	Closure_Redeemed_Points  float64 `bson:"Closure_Redeemed_Points" json:"Closure_Redeemed_Points"`
	Closure_Available_Points float64 `bson:"Closure_Available_Points" json:"Closure_Available_Points"` //(Awarded_Points + Expired_Points) - Redeemed_Points

	Opening_MainGSMBalance_PendingAmount float64 `bson:"Opening_MainGSMBalance_PendingAmount" json:"Opening_MainGSMBalance_PendingAmount"`
	Opening_MobileMoney_PendingAmount    float64 `bson:"Opening_MobileMoney_PendingAmount" json:"Opening_MobileMoney_PendingAmount"`
	Closure_MainGSMBalance_PendingAmount float64 `bson:"Closure_MainGSMBalance_PendingAmount" json:"Closure_MainGSMBalance_PendingAmount"`
	Closure_MobileMoney_PendingAmount    float64 `bson:"Closure_MobileMoney_PendingAmount" json:"Closure_MobileMoney_PendingAmount"`

	//response result
	ReceiveDate       time.Time `bson:"ReceiveDate" json:"-"`
	Status            string    `bson:"Status" json:"Status"` //successful, failed
	StatusCode        int       `bson:"StatusCode" json:"StatusCode"`
	StatusDescription string    `bson:"StatusDescription" json:"StatusDescription"`
	ErrorDescription  string    `bson:"ErrorDescription" json:"ErrorDescription"`
	StatusDate        time.Time `bson:"StatusDate" json:"StatusDate"`
	E2E_Elapsedtime   int64     `bson:"E2E_Elapsedtime" json:"E2E_Elapsedtime"` //receive date till return
}

type Loyalty_Expiry_log struct {
	ExpiryTime       time.Time `bson:"ExpiryTime" json:"ExpiryTime"`
	MSISDN           string    `bson:"MSISDN" json:"MSISDN"` //MSISDN
	Expiry_Rules_Key string    `bson:"Expiry_Rules_Key" json:"Expiry_Rules_Key"`

	Year_Month             string  `bson:"Year_Month" json:"Year_Month"`
	Month_Awarded_Points   float64 `bson:"Month_Awarded_Points" json:"Month_Awarded_Points"`
	Month_Redeemed_Points  float64 `bson:"Month_Redeemed_Points" json:"Month_Redeemed_Points"`
	Month_Available_Points float64 `bson:"Month_Available_Points" json:"Month_Available_Points"` //(Awarded_Points + Expired_Points) - Redeemed_Points
	Month_Expired_Points   float64 `bson:"Month_Expired_Points" json:"Month_Expired_Points"`     //expired are deducted from Awarded_Points

	Opening_Awarded_Points   float64 `bson:"Opening_Awarded_Points" json:"Opening_Awarded_Points"`
	Opening_Redeemed_Points  float64 `bson:"Opening_Redeemed_Points" json:"Opening_Redeemed_Points"`
	Opening_Available_Points float64 `bson:"Opening_Available_Points" json:"Opening_Available_Points"` //(Awarded_Points + Expired_Points) - Redeemed_Points
	Opening_Expired_Points   float64 `bson:"Opening_Expired_Points" json:"Opening_Expired_Points"`     //expired are deducted from Awarded_Points

	End_Awarded_Points   float64 `bson:"End_Awarded_Points" json:"End_Awarded_Points"`
	End_Redeemed_Points  float64 `bson:"End_Redeemed_Points" json:"End_Redeemed_Points"`
	End_Available_Points float64 `bson:"End_Available_Points" json:"End_Available_Points"` //(Awarded_Points + Expired_Points) - Redeemed_Points
	End_Expired_Points   float64 `bson:"End_Expired_Points" json:"End_Expired_Points"`     //expired are deducted from Awarded_Points

	OpeningLoyaltyLevel string `bson:"OpeningLoyaltyLevel" json:"OpeningLoyaltyLevel"`
	EndLoyaltyLevel     string `bson:"EndLoyaltyLevel" json:"EndLoyaltyLevel"`

	ExpiryStatus            string `bson:"ExpiryStatus" json:"ExpiryStatus"`
	ExpiryStatusDescription string `bson:"ExpiryStatusDescription" json:"ExpiryStatusDescription"`
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

type Loyalty_Level_Change_log struct {
	Level_Change_Date time.Time `bson:"Level_Change_Date" json:"Level_Change_Date"`
	MSISDN            string    `bson:"MSISDN" json:"MSISDN"` //MSISDN
	COS               string    `bson:"COS" json:"COS"`
	Joining_Date      time.Time `bson:"Joining_Date" json:"Joining_Date"` //
	ARPU              float64   `bson:"ARPU" json:"ARPU"`
	Customer_Id       int64     `bson:"Customer_Id" json:"Customer_Id"`
	Creation_date     time.Time `bson:"Creation_date" json:"Creation_date"`

	Previous_Loyalty_Level_Key  string    `bson:"Previous_Loyalty_Level_Key" json:"Previous_Loyalty_Level_Key"`
	Previous_Loyalty_Level_Date time.Time `bson:"Previous_Loyalty_Level_Date" json:"Previous_Loyalty_Level_Date"`

	New_Loyalty_Level_Key       string    `bson:"New_Loyalty_Level_Key" json:"New_Loyalty_Level_Key"`
	New_Loyalty_Level_Date      time.Time `bson:"New_Loyalty_Level_Date" json:"New_Loyalty_Level_Date"`
	New_Loyalty_Level_Direction string    `bson:"New_Loyalty_Level_Direction" json:"New_Loyalty_Level_Direction"` //last change: Up or Down
	New_Loyalty_Level_SetBy     string    `bson:"New_Loyalty_Level_SetBy" json:"New_Loyalty_Level_SetBy"`         //program or admin, if admin program cannot change anymore

	Loyalty_Account_Segment_Key       string    `bson:"Loyalty_Account_Segment_Key" json:"Loyalty_Account_Segment_Key"`
	Loyalty_Account_Segment_Date      time.Time `bson:"Loyalty_Account_Segment_Date" json:"Loyalty_Account_Segment_Date"`
	Loyalty_Account_Segment_Direction string    `bson:"Loyalty_Account_Segment_Direction" json:"Loyalty_Account_Segment_Direction"` //last change: Up or Down
	Loyalty_Account_Segment_SetBy     string    `bson:"Loyalty_Account_Segment_SetBy" json:"Loyalty_Account_Segment_SetBy"`         //program or admin, if admin program cannot change anymore

	Awarded_Points   float64   `bson:"Awarded_Points" json:"Awarded_Points"`
	Redeemed_Points  float64   `bson:"Redeemed_Points" json:"Redeemed_Points"`
	Available_Points float64   `bson:"Available_Points" json:"Available_Points"` //Awarded_Points - Redeemed_Points
	Last_Award_Date  time.Time `bson:"Last_Award_Date" json:"Last_Award_Date"`
	Last_Redeem_Date time.Time `bson:"Last_Redeem_Date" json:"Last_Redeem_Date"`
}

type Loyalty_Redemption_Request struct {
	MSISDN               string  `bson:"MSISDN" json:"MSISDN"`                   //MSISDN
	Redemption_Type      string  `bson:"Redemption_Type" json:"Redemption_Type"` //Airtime, Bundle, MobileMoney, SpinAndWin
	Redemption_Bunlde_Id string  `bson:"Redemption_Bunlde_Id" json:"Redemption_Bunlde_Id"`
	Redemption_Amount    float64 `bson:"Redemption_Amount" json:"Redemption_Amount"` //airtime or money amount the subscriber wish to redeem
	Points_To_Redeem     float64 `bson:"Points_To_Redeem" json:"Points_To_Redeem"`   //the amount of points the subscriber wish to redeem
}

type Loyalty_Redemption_log struct {
	//request Header info
	SourceIP    string        `bson:"SourceIP" json:"-"`
	SourceApp   string        `bson:"SourceApp" json:"-"`
	AppLogin    string        `bson:"AppLogin" json:"-"`
	AppVersion  string        `bson:"AppVersion" json:"-"`
	GPSLocation daoc.Location `bson:"GPSLocation" json:"-"`
	GSMLocation string        `bson:"GSMLocation" json:"-"`

	//request detail
	MSISDN               string    `bson:"MSISDN" json:"MSISDN"` //MSISDN
	ReceiveDate          time.Time `bson:"ReceiveDate" json:"ReceiveDate"`
	Redemption_Type      string    `bson:"Redemption_Type" json:"Redemption_Type"` //Airtime, Bundle, MobileMoney, SpinAndWin
	Redemption_Bunlde_Id string    `bson:"Redemption_Bunlde_Id" json:"Redemption_Bunlde_Id"`
	Redemption_Amount    float64   `bson:"Redemption_Amount" json:"Redemption_Amount"` //airtime or money amount the subscriber wish to redeem
	Points_To_Redeem     float64   `bson:"Points_To_Redeem" json:"Points_To_Redeem"`   //the amount of points the subscriber wish to redeem

	Customer_Id                 int64   `bson:"Customer_Id" json:"Customer_Id"`
	Account_Status              string  `bson:"Account_Status" json:"Account_Status"`
	Loyalty_Level_Key           string  `bson:"Loyalty_Level_Key" json:"Loyalty_Level_Key"`
	Loyalty_Account_Segment_Key string  `bson:"Loyalty_Account_Segment_Key" json:"Loyalty_Account_Segment_Key"`
	Opening_Awarded_Points      float64 `bson:"Opening_Awarded_Points" json:"Opening_Awarded_Points"`
	Opening_Redeemed_Points     float64 `bson:"Opening_Redeemed_Points" json:"Opening_Redeemed_Points"`
	Opening_Available_Points    float64 `bson:"Opening_Available_Points" json:"Opening_Available_Points"` //(Awarded_Points + Expired_Points) - Redeemed_Points
	Closure_Awarded_Points      float64 `bson:"Closure_Awarded_Points" json:"Closure_Awarded_Points"`
	Closure_Redeemed_Points     float64 `bson:"Closure_Redeemed_Points" json:"Closure_Redeemed_Points"`
	Closure_Available_Points    float64 `bson:"Closure_Available_Points" json:"Closure_Available_Points"` //(Awarded_Points + Expired_Points) - Redeemed_Points

	MinRequiredPoints               float64 `bson:"MinRequiredPoints" json:"MinRequiredPoints"`
	Allow_Negative_Balance_ToRedeem bool    `bson:"Allow_Negative_Balance_ToRedeem" json:"Allow_Negative_Balance_ToRedeem"`
	Allow_PendingLendme_ToRedeem    bool    `bson:"Allow_PendingLendme_ToRedeem" json:"Allow_PendingLendme_ToRedeem"`

	//Bundle Redeem detail
	Price_Loyalty_Points float64 `bson:"Price_Loyalty_Points" json:"Price_Loyalty_Points"`

	Points_Debit_Result        interface{} `bson:"Points_Debit_Result" json:"Points_Debit_Result"`
	Airtime_PurchaseResult     interface{} `bson:"Airtime_PurchaseResult" json:"Airtime_PurchaseResult"`
	Bundle_PurchaseResult      interface{} `bson:"Bundle_PurchaseResult" json:"Bundle_PurchaseResult"`
	MobileMoney_PurchaseResult interface{} `bson:"MobileMoney_PurchaseResult" json:"MobileMoney_PurchaseResult"`
	SpinAndWin_PurchaseResult  interface{} `bson:"SpinAndWin_PurchaseResult" json:"SpinAndWin_PurchaseResult"`

	//response result
	Status            string    `bson:"Status" json:"Status"` //successful, failed
	StatusCode        int       `bson:"StatusCode" json:"StatusCode"`
	StatusDescription string    `bson:"StatusDescription" json:"StatusDescription"`
	ErrorDescription  string    `bson:"ErrorDescription" json:"ErrorDescription"`
	StatusDate        time.Time `bson:"StatusDate" json:"StatusDate"`
	E2E_Elapsedtime   int64     `bson:"E2E_Elapsedtime" json:"E2E_Elapsedtime"` //receive date till return

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

// type Loyalty_Point_Redemption_Rules struct {
// 	Key                               string  `bson:"Key" json:"Key"`
// 	Redemption_Rules_Id               int64   `bson:"Redemption_Rules_Id" json:"Redemption_Rules_Id"`
// 	Description                       string  `bson:"Description" json:"Description"`
// 	Min_Accumulated_Points            float64 `bson:"Min_Accumulated_Points" json:"Min_Accumulated_Points"` //will be used for upgrade
// 	Allow_Negative_Balance_ToRedeem   bool    `bson:"Allow_Negative_Balance_ToRedeem" json:"Allow_Negative_Balance_ToRedeem"`
// 	Allow_PendingLendme_ToRedeem      bool    `bson:"Allow_PendingLendme_ToRedeem" json:"Allow_PendingLendme_ToRedeem"`
// 	Airtime_MinPoints                 float64 `bson:"Airtime_MinPoints" json:"Airtime_MinPoints"`
// 	Airtime_AmountPerPoint            float64 `bson:"Airtime_AmountPerPoint" json:"Airtime_AmountPerPoint"`
// 	MobileMoney_MinPoints             float64 `bson:"MobileMoney_MinPoints" json:"MobileMoney_MinPoints"`
// 	MobileMoney_AmountPerPoint        float64 `bson:"MobileMoney_AmountPerPoint" json:"MobileMoney_AmountPerPoint"`
// 	Bundles_MinPoints                 float64 `bson:"Bundles_MinPoints" json:"Bundles_MinPoints"`
// 	Bundles_Product_Catalogue_Channel string  `bson:"Bundles_Product_Catalogue_Channel" json:"Bundles_Product_Catalogue_Channel"`
// 	Bundles_Product_Catalogue_Plan    string  `bson:"Bundles_Product_Catalogue_Plan" json:"Bundles_Product_Catalogue_Plan"`
// 	Bundles_Product_Catalogue_Version string  `bson:"Bundles_Product_Catalogue_Version" json:"Bundles_Product_Catalogue_Version"`
// 	FreeSpinAndWin_MinPoints          float64 `bson:"FreeSpinAndWin_MinPoints" json:"FreeSpinAndWin_MinPoints"`
// 	FreeSpinAndWin_PointsPerSpin      float64 `bson:"FreeSpinAndWin_PointsPerSpin" json:"FreeSpinAndWin_PointsPerSpin"`
// }

type Loyalty_Campaign struct {
	Key             string    `bson:"Key" json:"Key"` //Campaign Name
	Campaign_Id     int64     `bson:"Campaign_Id" json:"Campaign_Id"`
	Description     string    `bson:"Description" json:"Description"`
	Start_Date      time.Time `bson:"Start_Date" json:"Start_Date"`
	End_Date        time.Time `bson:"End_Date" json:"End_Date"`
	Campaign_Status string    `bson:"Campaign_Status" json:"Campaign_Status"` //launched, paused, resumed, cancelled, or stopped

	//target rules
	Target_All_Subs     bool              `bson:"Target_All_Subs" json:"Target_All_Subs"` //if true, all target features will be ignored
	Target_Level_Key    map[string]bool   `bson:"Target_Level_Key" json:"Target_Level_Key"`
	Target_Segment_Key  map[string]bool   `bson:"Target_Segment_Key" json:"Target_Segment_Key"`
	Target_List         map[string]string `bson:"Target_List" json:"Target_List"` //we need to be able to define target list
	LoyaltyPoints_From  float64           `bson:"LoyaltyPoints_From" json:"LoyaltyPoints_From"`
	LoyaltyPoints_Till  float64           `bson:"LoyaltyPoints_Till" json:"LoyaltyPoints_Till"`
	AON_From            float64           `bson:"AON_From" json:"AON_From"` //months
	AON_Till            float64           `bson:"AON_Till" json:"AON_Till"` //months
	BundlesPurchase     map[string]string `bson:"BundlesPurchase" json:"BundlesPurchase"`
	MobileAppDailyUsage bool              `bson:"MobileAppDailyUsage" json:"MobileAppDailyUsage"`
	//Points earned scheme
	Multiplier       float64 `bson:"Multiplier" json:"Multiplier"`
	Fixed_Points     float64 `bson:"Fixed_Points" json:"Fixed_Points"`
	Action_Frequency string  `bson:"Action_Frequency" json:"Action_Frequency"` //Once, Once_Daily, Each time
}

type Loyalty_Campaign_Target_List struct {
	Key          string    `bson:"Key" json:"Key"` //Campaign_Key|MSISDN
	Campaign_Key int64     `bson:"Campaign_Key" json:"Campaign_Key"`
	MSISDN       string    `bson:"MSISDN" json:"MSISDN"`
	Start_Date   time.Time `bson:"Start_Date" json:"Start_Date"`
	End_Date     time.Time `bson:"End_Date" json:"End_Date"`
}

type Loyalty_Campaign_Action_Frequency_Control struct {
	Key          string    `bson:"Key" json:"Key"` //Target List Name|MSISDN
	Campaign_Key int64     `bson:"Campaign_Key" json:"Campaign_Key"`
	MSISDN       string    `bson:"MSISDN" json:"MSISDN"`
	Action_Date  time.Time `bson:"Action_Date" json:"Action_Date"`
	End_Date     time.Time `bson:"End_Date" json:"End_Date"` //when the control will be removed
}
