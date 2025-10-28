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
	EVC_Account_Balance             float64 `bson:"EVC_Account_Balance" json:"EVC_Account_Balance"`
	Merchant_Account_Balance        float64 `bson:"Merchant_Account_Balance" json:"Merchant_Account_Balance"`
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

type Loyalty_Governance_log struct {
	Log_Date                 time.Time `bson:"Log_Date" json:"Log_Date"`
	Available_Points_Pool    float64   `bson:"Available_Points_Pool" json:"Available_Points_Pool"`
	Distributed_Points_Pool  float64   `bson:"Distributed_Points_Pool" json:"Distributed_Points_Pool"`
	Redeemed_Points_Pool     float64   `bson:"Redeemed_Points_Pool" json:"Redeemed_Points_Pool"`
	Expired_Points_Pool      float64   `bson:"Expired_Points_Pool" json:"Expired_Points_Pool"`
	EVC_Account_Balance      float64   `bson:"EVC_Account_Balance" json:"EVC_Account_Balance"`
	Merchant_Account_Balance float64   `bson:"Merchant_Account_Balance" json:"Merchant_Account_Balance"`
}
type Seniority_Level struct {
	Key                         string  `bson:"Key" json:"Key"`
	Loyalty_Seniority_Level_Key string  `bson:"Loyalty_Seniority_Level_Key" json:"Loyalty_Seniority_Level_Key"`
	Multiplier_Percentage       float64 `bson:"Multiplier_Percentage" json:"Multiplier_Percentage"`
}
type Loyalty_Level struct {
	Key                    string            `bson:"Key" json:"Key"`
	Level_Id               int64             `bson:"Level_Id" json:"Level_Id"`
	Description            string            `bson:"Description" json:"Description"`
	Min_Accumulated_Points float64           `bson:"Min_Accumulated_Points" json:"Min_Accumulated_Points"` //will be used for upgrade
	Max_Accumulated_Points float64           `bson:"Max_Accumulated_Points" json:"Max_Accumulated_Points"` //will be used for downgrade
	EnableRedeem           bool              `bson:"EnableRedeem" json:"EnableRedeem"`
	DowngradeToLevel_Key   string            `bson:"DowngradeToLevel_Key" json:"DowngradeToLevel_Key"`
	Seniority_Levels       []Seniority_Level `bson:"Seniority_Levels" json:"Seniority_Levels"`
}

type Loyalty_Level_AddRequest struct {
	Key                    string            `bson:"Key" json:"Key"`
	Level_Id               int64             `bson:"Level_Id" json:"Level_Id"`
	Description            string            `bson:"Description" json:"Description"`
	Min_Accumulated_Points float64           `bson:"Min_Accumulated_Points" json:"Min_Accumulated_Points"` //will be used for upgrade
	Max_Accumulated_Points float64           `bson:"Max_Accumulated_Points" json:"Max_Accumulated_Points"` //will be used for downgrade
	EnableRedeem           bool              `bson:"EnableRedeem" json:"EnableRedeem"`
	DowngradeToLevel_Key   string            `bson:"DowngradeToLevel_Key" json:"DowngradeToLevel_Key"`
	Seniority_Levels       []Seniority_Level `bson:"Seniority_Levels" json:"Seniority_Levels"`
}

type Loyalty_Level_EditRequest struct {
	Key                    string            `bson:"Key" json:"Key"`
	NewKey                 string            `bson:"NewKey" json:"NewKey"`
	Level_Id               int64             `bson:"Level_Id" json:"Level_Id"`
	Description            string            `bson:"Description" json:"Description"`
	Min_Accumulated_Points float64           `bson:"Min_Accumulated_Points" json:"Min_Accumulated_Points"` //will be used for upgrade
	Max_Accumulated_Points float64           `bson:"Max_Accumulated_Points" json:"Max_Accumulated_Points"` //will be used for downgrade
	EnableRedeem           bool              `bson:"EnableRedeem" json:"EnableRedeem"`
	DowngradeToLevel_Key   string            `bson:"DowngradeToLevel_Key" json:"DowngradeToLevel_Key"`
	Seniority_Levels       []Seniority_Level `bson:"Seniority_Levels" json:"Seniority_Levels"`
}

type Loyalty_Seniority_Level struct {
	Key          string  `bson:"Key" json:"Key"`
	Seniority_Id int64   `bson:"Seniority_Id" json:"Seniority_Id"`
	Name         string  `bson:"Name" json:"Name"`
	Description  string  `bson:"Description" json:"Description"`
	AON_From     float64 `bson:"AON_From" json:"AON_From"` //months
	AON_Till     float64 `bson:"AON_Till" json:"AON_Till"` //months
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

type Loyalty_Point_Earning_Rules struct {
	Key                              string  `bson:"Key" json:"Key"` //Program name
	Earning_Rules_Id                 int64   `bson:"Earning_Rules_Id" json:"Earning_Rules_Id"`
	Description                      string  `bson:"Description" json:"Description"`
	Welcome_Points                   float64 `bson:"Welcome_Points" json:"Welcome_Points"`
	MobileAppDaily_Login             float64 `bson:"MobileAppDaily_Login" json:"MobileAppDaily_Login"`
	Welcome_Notification             bool    `bson:"Welcome_Notification" json:"Welcome_Notification"`
	Welcome_Notification_Sender      string  `bson:"Welcome_Notification_Sender" json:"Welcome_Notification_Sender"`
	Welcome_Notification_Text        string  `bson:"Welcome_Notification_Text" json:"Welcome_Notification_Text"`
	Level_Change_Notification        bool    `bson:"Level_Change_Notification" json:"Level_Change_Notification"`
	Level_Change_Notification_Sender string  `bson:"Level_Change_Notification_Sender" json:"Level_Change_Notification_Sender"`
	Level_Change_Notification_Text   string  `bson:"Level_Change_Notification_Text" json:"Level_Change_Notification_Text"`
	//GSM main balance consumption
	MainGSMBalance_Amount float64 `bson:"MainGSMBalance_Amount" json:"MainGSMBalance_Amount"`
	MainGSMBalance_Points float64 `bson:"MainGSMBalance_Points" json:"MainGSMBalance_Points"`
	//GSM Scratch Card Recharge
	GSM_SC_Airtime_Award_Type string  `bson:"GSM_SC_Airtime_Award_Type" json:"GSM_SC_Airtime_Award_Type"` //"Transaction" or "Amount"
	GSM_SC_Airtime_Amount     float64 `bson:"GSM_SC_Airtime_Amount" json:"GSM_SC_Airtime_Amount"`
	GSM_SC_Airtime_Points     float64 `bson:"GSM_SC_Airtime_Points" json:"GSM_SC_Airtime_Points"`
	//GSM EVC airtime recharge
	GSM_EVC_Airtime_Award_Type string  `bson:"GSM_EVC_Airtime_Award_Type" json:"GSM_EVC_Airtime_Award_Type"` //"Transaction" or "Amount"
	GSM_EVC_Airtime_Amount     float64 `bson:"GSM_EVC_Airtime_Amount" json:"GSM_EVC_Airtime_Amount"`
	GSM_EVC_Airtime_Points     float64 `bson:"GSM_EVC_Airtime_Points" json:"GSM_EVC_Airtime_Points"`
	//GSM EVC Bundle purchase
	GSM_EVC_Bundle_Award_Type string  `bson:"GSM_EVC_Bundle_Award_Type" json:"GSM_EVC_Bundle_Award_Type"` //"Transaction" or "Amount"
	GSM_EVC_Bundle_Amount     float64 `bson:"GSM_EVC_Bundle_Amount" json:"GSM_EVC_Bundle_Amount"`
	GSM_EVC_Bundle_Points     float64 `bson:"GSM_EVC_Bundle_Points" json:"GSM_EVC_Bundle_Points"`
	//Mobile Money
	MM_P2P_Award_Type                        string   `bson:"MM_P2P_Award_Type" json:"MM_P2P_Award_Type"` //"Transaction" or "Amount"
	MM_P2P_Amount                            float64  `bson:"MM_P2P_Amount" json:"MM_P2P_Amount"`
	MM_P2P_Points                            float64  `bson:"MM_P2P_Points" json:"MM_P2P_Points"`
	MM_CASHIN_Award_Type                     string   `bson:"MM_CASHIN_Award_Type" json:"MM_CASHIN_Award_Type"` //"Transaction" or "Amount"
	MM_CASHIN_Amount                         float64  `bson:"MM_CASHIN_Amount" json:"MM_CASHIN_Amount"`
	MM_CASHIN_Points                         float64  `bson:"MM_CASHIN_Points" json:"MM_CASHIN_Points"`
	MM_CASHOUT_Award_Type                    string   `bson:"MM_CASHOUT_Award_Type" json:"MM_CASHOUT_Award_Type"` //"Transaction" or "Amount"
	MM_CASHOUT_Amount                        float64  `bson:"MM_CASHOUT_Amount" json:"MM_CASHOUT_Amount"`
	MM_CASHOUT_Points                        float64  `bson:"MM_CASHOUT_Points" json:"MM_CASHOUT_Points"`
	MM_MERCHPAY_Award_Type                   string   `bson:"MM_MERCHPAY_Award_Type" json:"MM_MERCHPAY_Award_Type"` //"Transaction" or "Amount"
	MM_MERCHPAY_Amount                       float64  `bson:"MM_MERCHPAY_Amount" json:"MM_MERCHPAY_Amount"`
	MM_MERCHPAY_Points                       float64  `bson:"MM_MERCHPAY_Points" json:"MM_MERCHPAY_Points"`
	MM_BILLPAY_Award_Type                    string   `bson:"MM_BILLPAY_Award_Type" json:"MM_BILLPAY_Award_Type"` //"Transaction" or "Amount"
	MM_BILLPAY_Amount                        float64  `bson:"MM_BILLPAY_Amount" json:"MM_BILLPAY_Amount"`
	MM_BILLPAY_Points                        float64  `bson:"MM_BILLPAY_Points" json:"MM_BILLPAY_Points"`
	Earning_Rules_Overwrite_MM_MERCHPAY_Keys []string `bson:"Earning_Rules_Overwrite_MM_MERCHPAY_Keys" json:"Earning_Rules_Overwrite_MM_MERCHPAY_Keys"`
	Earning_Rules_Overwrite_MM_BILLPAY_Keys  []string `bson:"Earning_Rules_Overwrite_MM_BILLPAY_Keys" json:"Earning_Rules_Overwrite_MM_BILLPAY_Keys"`
	MM_RC_Bundle_Award_Type                  string   `bson:"MM_RC_Bundle_Award_Type" json:"MM_RC_Bundle_Award_Type"` //"Transaction" or "Amount"
	MM_RC_Bundle_Amount                      float64  `bson:"MM_RC_Bundle_Amount" json:"MM_RC_Bundle_Amount"`
	MM_RC_Bundle_Points                      float64  `bson:"MM_RC_Bundle_Points" json:"MM_RC_Bundle_Points"`
	MM_RC_Airtime_Award_Type                 string   `bson:"MM_RC_Airtime_Award_Type" json:"MM_RC_Airtime_Award_Type"` //"Transaction" or "Amount"
	MM_RC_Airtime_Amount                     float64  `bson:"MM_RC_Airtime_Amount" json:"MM_RC_Airtime_Amount"`
	MM_RC_Airtime_Points                     float64  `bson:"MM_RC_Airtime_Points" json:"MM_RC_Airtime_Points"`
	MM_CTMMOREQ_Bundle_Award_Type            string   `bson:"MM_CTMMOREQ_Bundle_Award_Type" json:"MM_CTMMOREQ_Bundle_Award_Type"` //"Transaction" or "Amount"
	MM_CTMMOREQ_Bundle_Amount                float64  `bson:"MM_CTMMOREQ_Bundle_Amount" json:"MM_CTMMOREQ_Bundle_Amount"`
	MM_CTMMOREQ_Bundle_Points                float64  `bson:"MM_CTMMOREQ_Bundle_Points" json:"MM_CTMMOREQ_Bundle_Points"`
	MM_CTMMOREQ_Airtime_Award_Type           string   `bson:"MM_CTMMOREQ_Airtime_Award_Type" json:"MM_CTMMOREQ_Airtime_Award_Type"` //"Transaction" or "Amount"
	MM_CTMMOREQ_Airtime_Amount               float64  `bson:"MM_CTMMOREQ_Airtime_Amount" json:"MM_CTMMOREQ_Airtime_Amount"`
	MM_CTMMOREQ_Airtime_Points               float64  `bson:"MM_CTMMOREQ_Airtime_Points" json:"MM_CTMMOREQ_Airtime_Points"`
	MM_CBWREQ_Award_Type                     string   `bson:"MM_CBWREQ_Award_Type" json:"MM_CBWREQ_Award_Type"` //"Transaction" or "Amount"
	MM_CBWREQ_Amount                         float64  `bson:"MM_CBWREQ_Amount" json:"MM_CBWREQ_Amount"`
	MM_CBWREQ_Points                         float64  `bson:"MM_CBWREQ_Points" json:"MM_CBWREQ_Points"`
	//Mobile money airtime --> detected from IN live feed
	MM_Airtime_Award_Type string  `bson:"MM_Airtime_Award_Type" json:"MM_Airtime_Award_Type"` //"Transaction" or "Amount"
	MM_Airtime_Amount     float64 `bson:"MM_Airtime_Amount" json:"MM_Airtime_Amount"`
	MM_Airtime_Points     float64 `bson:"MM_Airtime_Points" json:"MM_Airtime_Points"`
	//Mobile money data purchase --> detected from IN live feed
	MM_Bundle_Award_Type string  `bson:"MM_Bundle_Award_Type" json:"MM_Bundle_Award_Type"` //"Transaction" or "Amount"
	MM_Bundle_Amount     float64 `bson:"MM_Bundle_Amount" json:"MM_Bundle_Amount"`
	MM_Bundle_Points     float64 `bson:"MM_Bundle_Points" json:"MM_Bundle_Points"`
}

type Loyalty_Point_Earning_Rules_AddRequest struct {
	Key                              string  `bson:"Key" json:"Key"` //Program name
	Earning_Rules_Id                 int64   `bson:"Earning_Rules_Id" json:"Earning_Rules_Id"`
	Description                      string  `bson:"Description" json:"Description"`
	Welcome_Points                   float64 `bson:"Welcome_Points" json:"Welcome_Points"`
	MobileAppDaily_Login             float64 `bson:"MobileAppDaily_Login" json:"MobileAppDaily_Login"`
	Welcome_Notification             bool    `bson:"Welcome_Notification" json:"Welcome_Notification"`
	Welcome_Notification_Sender      string  `bson:"Welcome_Notification_Sender" json:"Welcome_Notification_Sender"`
	Welcome_Notification_Text        string  `bson:"Welcome_Notification_Text" json:"Welcome_Notification_Text"`
	Level_Change_Notification        bool    `bson:"Level_Change_Notification" json:"Level_Change_Notification"`
	Level_Change_Notification_Sender string  `bson:"Level_Change_Notification_Sender" json:"Level_Change_Notification_Sender"`
	Level_Change_Notification_Text   string  `bson:"Level_Change_Notification_Text" json:"Level_Change_Notification_Text"`
	//GSM main balance consumption
	MainGSMBalance_Amount float64 `bson:"MainGSMBalance_Amount" json:"MainGSMBalance_Amount"`
	MainGSMBalance_Points float64 `bson:"MainGSMBalance_Points" json:"MainGSMBalance_Points"`
	//GSM Scratch Card Recharge
	GSM_SC_Airtime_Award_Type string  `bson:"GSM_SC_Airtime_Award_Type" json:"GSM_SC_Airtime_Award_Type"` //"Transaction" or "Amount"
	GSM_SC_Airtime_Amount     float64 `bson:"GSM_SC_Airtime_Amount" json:"GSM_SC_Airtime_Amount"`
	GSM_SC_Airtime_Points     float64 `bson:"GSM_SC_Airtime_Points" json:"GSM_SC_Airtime_Points"`
	//GSM EVC airtime recharge
	GSM_EVC_Airtime_Award_Type string  `bson:"GSM_EVC_Airtime_Award_Type" json:"GSM_EVC_Airtime_Award_Type"` //"Transaction" or "Amount"
	GSM_EVC_Airtime_Amount     float64 `bson:"GSM_EVC_Airtime_Amount" json:"GSM_EVC_Airtime_Amount"`
	GSM_EVC_Airtime_Points     float64 `bson:"GSM_EVC_Airtime_Points" json:"GSM_EVC_Airtime_Points"`
	//GSM EVC Bundle purchase
	GSM_EVC_Bundle_Award_Type string  `bson:"GSM_EVC_Bundle_Award_Type" json:"GSM_EVC_Bundle_Award_Type"` //"Transaction" or "Amount"
	GSM_EVC_Bundle_Amount     float64 `bson:"GSM_EVC_Bundle_Amount" json:"GSM_EVC_Bundle_Amount"`
	GSM_EVC_Bundle_Points     float64 `bson:"GSM_EVC_Bundle_Points" json:"GSM_EVC_Bundle_Points"`
	//Mobile Money
	MM_P2P_Award_Type                        string   `bson:"MM_P2P_Award_Type" json:"MM_P2P_Award_Type"` //"Transaction" or "Amount"
	MM_P2P_Amount                            float64  `bson:"MM_P2P_Amount" json:"MM_P2P_Amount"`
	MM_P2P_Points                            float64  `bson:"MM_P2P_Points" json:"MM_P2P_Points"`
	MM_CASHIN_Award_Type                     string   `bson:"MM_CASHIN_Award_Type" json:"MM_CASHIN_Award_Type"` //"Transaction" or "Amount"
	MM_CASHIN_Amount                         float64  `bson:"MM_CASHIN_Amount" json:"MM_CASHIN_Amount"`
	MM_CASHIN_Points                         float64  `bson:"MM_CASHIN_Points" json:"MM_CASHIN_Points"`
	MM_CASHOUT_Award_Type                    string   `bson:"MM_CASHOUT_Award_Type" json:"MM_CASHOUT_Award_Type"` //"Transaction" or "Amount"
	MM_CASHOUT_Amount                        float64  `bson:"MM_CASHOUT_Amount" json:"MM_CASHOUT_Amount"`
	MM_CASHOUT_Points                        float64  `bson:"MM_CASHOUT_Points" json:"MM_CASHOUT_Points"`
	MM_MERCHPAY_Award_Type                   string   `bson:"MM_MERCHPAY_Award_Type" json:"MM_MERCHPAY_Award_Type"` //"Transaction" or "Amount"
	MM_MERCHPAY_Amount                       float64  `bson:"MM_MERCHPAY_Amount" json:"MM_MERCHPAY_Amount"`
	MM_MERCHPAY_Points                       float64  `bson:"MM_MERCHPAY_Points" json:"MM_MERCHPAY_Points"`
	MM_BILLPAY_Award_Type                    string   `bson:"MM_BILLPAY_Award_Type" json:"MM_BILLPAY_Award_Type"` //"Transaction" or "Amount"
	MM_BILLPAY_Amount                        float64  `bson:"MM_BILLPAY_Amount" json:"MM_BILLPAY_Amount"`
	MM_BILLPAY_Points                        float64  `bson:"MM_BILLPAY_Points" json:"MM_BILLPAY_Points"`
	Earning_Rules_Overwrite_MM_MERCHPAY_Keys []string `bson:"Earning_Rules_Overwrite_MM_MERCHPAY_Keys" json:"Earning_Rules_Overwrite_MM_MERCHPAY_Keys"`
	Earning_Rules_Overwrite_MM_BILLPAY_Keys  []string `bson:"Earning_Rules_Overwrite_MM_BILLPAY_Keys" json:"Earning_Rules_Overwrite_MM_BILLPAY_Keys"`
	MM_RC_Bundle_Award_Type                  string   `bson:"MM_RC_Bundle_Award_Type" json:"MM_RC_Bundle_Award_Type"` //"Transaction" or "Amount"
	MM_RC_Bundle_Amount                      float64  `bson:"MM_RC_Bundle_Amount" json:"MM_RC_Bundle_Amount"`
	MM_RC_Bundle_Points                      float64  `bson:"MM_RC_Bundle_Points" json:"MM_RC_Bundle_Points"`
	MM_RC_Airtime_Award_Type                 string   `bson:"MM_RC_Airtime_Award_Type" json:"MM_RC_Airtime_Award_Type"` //"Transaction" or "Amount"
	MM_RC_Airtime_Amount                     float64  `bson:"MM_RC_Airtime_Amount" json:"MM_RC_Airtime_Amount"`
	MM_RC_Airtime_Points                     float64  `bson:"MM_RC_Airtime_Points" json:"MM_RC_Airtime_Points"`
	MM_CTMMOREQ_Bundle_Award_Type            string   `bson:"MM_CTMMOREQ_Bundle_Award_Type" json:"MM_CTMMOREQ_Bundle_Award_Type"` //"Transaction" or "Amount"
	MM_CTMMOREQ_Bundle_Amount                float64  `bson:"MM_CTMMOREQ_Bundle_Amount" json:"MM_CTMMOREQ_Bundle_Amount"`
	MM_CTMMOREQ_Bundle_Points                float64  `bson:"MM_CTMMOREQ_Bundle_Points" json:"MM_CTMMOREQ_Bundle_Points"`
	MM_CTMMOREQ_Airtime_Award_Type           string   `bson:"MM_CTMMOREQ_Airtime_Award_Type" json:"MM_CTMMOREQ_Airtime_Award_Type"` //"Transaction" or "Amount"
	MM_CTMMOREQ_Airtime_Amount               float64  `bson:"MM_CTMMOREQ_Airtime_Amount" json:"MM_CTMMOREQ_Airtime_Amount"`
	MM_CTMMOREQ_Airtime_Points               float64  `bson:"MM_CTMMOREQ_Airtime_Points" json:"MM_CTMMOREQ_Airtime_Points"`
	MM_CBWREQ_Award_Type                     string   `bson:"MM_CBWREQ_Award_Type" json:"MM_CBWREQ_Award_Type"` //"Transaction" or "Amount"
	MM_CBWREQ_Amount                         float64  `bson:"MM_CBWREQ_Amount" json:"MM_CBWREQ_Amount"`
	MM_CBWREQ_Points                         float64  `bson:"MM_CBWREQ_Points" json:"MM_CBWREQ_Points"`
	//Mobile money airtime --> detected from IN live feed
	MM_Airtime_Award_Type string  `bson:"MM_Airtime_Award_Type" json:"MM_Airtime_Award_Type"` //"Transaction" or "Amount"
	MM_Airtime_Amount     float64 `bson:"MM_Airtime_Amount" json:"MM_Airtime_Amount"`
	MM_Airtime_Points     float64 `bson:"MM_Airtime_Points" json:"MM_Airtime_Points"`
	//Mobile money data purchase --> detected from IN live feed
	MM_Bundle_Award_Type string  `bson:"MM_Bundle_Award_Type" json:"MM_Bundle_Award_Type"` //"Transaction" or "Amount"
	MM_Bundle_Amount     float64 `bson:"MM_Bundle_Amount" json:"MM_Bundle_Amount"`
	MM_Bundle_Points     float64 `bson:"MM_Bundle_Points" json:"MM_Bundle_Points"`
}

type Loyalty_Point_Earning_Rules_EditRequest struct {
	Key                              string  `bson:"Key" json:"Key"` //Program name
	NewKey                           string  `bson:"NewKey" json:"NewKey"`
	Earning_Rules_Id                 int64   `bson:"Earning_Rules_Id" json:"Earning_Rules_Id"`
	Description                      string  `bson:"Description" json:"Description"`
	Welcome_Points                   float64 `bson:"Welcome_Points" json:"Welcome_Points"`
	MobileAppDaily_Login             float64 `bson:"MobileAppDaily_Login" json:"MobileAppDaily_Login"`
	Welcome_Notification             bool    `bson:"Welcome_Notification" json:"Welcome_Notification"`
	Welcome_Notification_Sender      string  `bson:"Welcome_Notification_Sender" json:"Welcome_Notification_Sender"`
	Welcome_Notification_Text        string  `bson:"Welcome_Notification_Text" json:"Welcome_Notification_Text"`
	Level_Change_Notification        bool    `bson:"Level_Change_Notification" json:"Level_Change_Notification"`
	Level_Change_Notification_Sender string  `bson:"Level_Change_Notification_Sender" json:"Level_Change_Notification_Sender"`
	Level_Change_Notification_Text   string  `bson:"Level_Change_Notification_Text" json:"Level_Change_Notification_Text"`
	//GSM main balance consumption
	MainGSMBalance_Amount float64 `bson:"MainGSMBalance_Amount" json:"MainGSMBalance_Amount"`
	MainGSMBalance_Points float64 `bson:"MainGSMBalance_Points" json:"MainGSMBalance_Points"`
	//GSM Scratch Card Recharge
	GSM_SC_Airtime_Award_Type string  `bson:"GSM_SC_Airtime_Award_Type" json:"GSM_SC_Airtime_Award_Type"` //"Transaction" or "Amount"
	GSM_SC_Airtime_Amount     float64 `bson:"GSM_SC_Airtime_Amount" json:"GSM_SC_Airtime_Amount"`
	GSM_SC_Airtime_Points     float64 `bson:"GSM_SC_Airtime_Points" json:"GSM_SC_Airtime_Points"`
	//GSM EVC airtime recharge
	GSM_EVC_Airtime_Award_Type string  `bson:"GSM_EVC_Airtime_Award_Type" json:"GSM_EVC_Airtime_Award_Type"` //"Transaction" or "Amount"
	GSM_EVC_Airtime_Amount     float64 `bson:"GSM_EVC_Airtime_Amount" json:"GSM_EVC_Airtime_Amount"`
	GSM_EVC_Airtime_Points     float64 `bson:"GSM_EVC_Airtime_Points" json:"GSM_EVC_Airtime_Points"`
	//GSM EVC Bundle purchase
	GSM_EVC_Bundle_Award_Type string  `bson:"GSM_EVC_Bundle_Award_Type" json:"GSM_EVC_Bundle_Award_Type"` //"Transaction" or "Amount"
	GSM_EVC_Bundle_Amount     float64 `bson:"GSM_EVC_Bundle_Amount" json:"GSM_EVC_Bundle_Amount"`
	GSM_EVC_Bundle_Points     float64 `bson:"GSM_EVC_Bundle_Points" json:"GSM_EVC_Bundle_Points"`
	//Mobile Money
	MM_P2P_Award_Type                        string   `bson:"MM_P2P_Award_Type" json:"MM_P2P_Award_Type"` //"Transaction" or "Amount"
	MM_P2P_Amount                            float64  `bson:"MM_P2P_Amount" json:"MM_P2P_Amount"`
	MM_P2P_Points                            float64  `bson:"MM_P2P_Points" json:"MM_P2P_Points"`
	MM_CASHIN_Award_Type                     string   `bson:"MM_CASHIN_Award_Type" json:"MM_CASHIN_Award_Type"` //"Transaction" or "Amount"
	MM_CASHIN_Amount                         float64  `bson:"MM_CASHIN_Amount" json:"MM_CASHIN_Amount"`
	MM_CASHIN_Points                         float64  `bson:"MM_CASHIN_Points" json:"MM_CASHIN_Points"`
	MM_CASHOUT_Award_Type                    string   `bson:"MM_CASHOUT_Award_Type" json:"MM_CASHOUT_Award_Type"` //"Transaction" or "Amount"
	MM_CASHOUT_Amount                        float64  `bson:"MM_CASHOUT_Amount" json:"MM_CASHOUT_Amount"`
	MM_CASHOUT_Points                        float64  `bson:"MM_CASHOUT_Points" json:"MM_CASHOUT_Points"`
	MM_MERCHPAY_Award_Type                   string   `bson:"MM_MERCHPAY_Award_Type" json:"MM_MERCHPAY_Award_Type"` //"Transaction" or "Amount"
	MM_MERCHPAY_Amount                       float64  `bson:"MM_MERCHPAY_Amount" json:"MM_MERCHPAY_Amount"`
	MM_MERCHPAY_Points                       float64  `bson:"MM_MERCHPAY_Points" json:"MM_MERCHPAY_Points"`
	MM_BILLPAY_Award_Type                    string   `bson:"MM_BILLPAY_Award_Type" json:"MM_BILLPAY_Award_Type"` //"Transaction" or "Amount"
	MM_BILLPAY_Amount                        float64  `bson:"MM_BILLPAY_Amount" json:"MM_BILLPAY_Amount"`
	MM_BILLPAY_Points                        float64  `bson:"MM_BILLPAY_Points" json:"MM_BILLPAY_Points"`
	Earning_Rules_Overwrite_MM_MERCHPAY_Keys []string `bson:"Earning_Rules_Overwrite_MM_MERCHPAY_Keys" json:"Earning_Rules_Overwrite_MM_MERCHPAY_Keys"`
	Earning_Rules_Overwrite_MM_BILLPAY_Keys  []string `bson:"Earning_Rules_Overwrite_MM_BILLPAY_Keys" json:"Earning_Rules_Overwrite_MM_BILLPAY_Keys"`
	MM_RC_Bundle_Award_Type                  string   `bson:"MM_RC_Bundle_Award_Type" json:"MM_RC_Bundle_Award_Type"` //"Transaction" or "Amount"
	MM_RC_Bundle_Amount                      float64  `bson:"MM_RC_Bundle_Amount" json:"MM_RC_Bundle_Amount"`
	MM_RC_Bundle_Points                      float64  `bson:"MM_RC_Bundle_Points" json:"MM_RC_Bundle_Points"`
	MM_RC_Airtime_Award_Type                 string   `bson:"MM_RC_Airtime_Award_Type" json:"MM_RC_Airtime_Award_Type"` //"Transaction" or "Amount"
	MM_RC_Airtime_Amount                     float64  `bson:"MM_RC_Airtime_Amount" json:"MM_RC_Airtime_Amount"`
	MM_RC_Airtime_Points                     float64  `bson:"MM_RC_Airtime_Points" json:"MM_RC_Airtime_Points"`
	MM_CTMMOREQ_Bundle_Award_Type            string   `bson:"MM_CTMMOREQ_Bundle_Award_Type" json:"MM_CTMMOREQ_Bundle_Award_Type"` //"Transaction" or "Amount"
	MM_CTMMOREQ_Bundle_Amount                float64  `bson:"MM_CTMMOREQ_Bundle_Amount" json:"MM_CTMMOREQ_Bundle_Amount"`
	MM_CTMMOREQ_Bundle_Points                float64  `bson:"MM_CTMMOREQ_Bundle_Points" json:"MM_CTMMOREQ_Bundle_Points"`
	MM_CTMMOREQ_Airtime_Award_Type           string   `bson:"MM_CTMMOREQ_Airtime_Award_Type" json:"MM_CTMMOREQ_Airtime_Award_Type"` //"Transaction" or "Amount"
	MM_CTMMOREQ_Airtime_Amount               float64  `bson:"MM_CTMMOREQ_Airtime_Amount" json:"MM_CTMMOREQ_Airtime_Amount"`
	MM_CTMMOREQ_Airtime_Points               float64  `bson:"MM_CTMMOREQ_Airtime_Points" json:"MM_CTMMOREQ_Airtime_Points"`
	MM_CBWREQ_Award_Type                     string   `bson:"MM_CBWREQ_Award_Type" json:"MM_CBWREQ_Award_Type"` //"Transaction" or "Amount"
	MM_CBWREQ_Amount                         float64  `bson:"MM_CBWREQ_Amount" json:"MM_CBWREQ_Amount"`
	MM_CBWREQ_Points                         float64  `bson:"MM_CBWREQ_Points" json:"MM_CBWREQ_Points"`
	//Mobile money airtime --> detected from IN live feed
	MM_Airtime_Award_Type string  `bson:"MM_Airtime_Award_Type" json:"MM_Airtime_Award_Type"` //"Transaction" or "Amount"
	MM_Airtime_Amount     float64 `bson:"MM_Airtime_Amount" json:"MM_Airtime_Amount"`
	MM_Airtime_Points     float64 `bson:"MM_Airtime_Points" json:"MM_Airtime_Points"`
	//Mobile money data purchase --> detected from IN live feed
	MM_Bundle_Award_Type string  `bson:"MM_Bundle_Award_Type" json:"MM_Bundle_Award_Type"` //"Transaction" or "Amount"
	MM_Bundle_Amount     float64 `bson:"MM_Bundle_Amount" json:"MM_Bundle_Amount"`
	MM_Bundle_Points     float64 `bson:"MM_Bundle_Points" json:"MM_Bundle_Points"`
}

type Loyalty_Point_Earning_Rules_Overwrite struct {
	Key                        string  `bson:"Key" json:"Key"` // + Earning_Rule_Key"|" + MM_Transaction_Type + "|" + MSISDN
	Earning_Rules_Overwrite_Id int64   `bson:"Earning_Rules_Overwrite_Id" json:"Earning_Rules_Overwrite_Id"`
	Earning_Rule_Key           string  `bson:"Earning_Rule_Key" json:"Earning_Rule_Key"`
	MM_Transaction_Type        string  `bson:"MM_Transaction_Type" json:"MM_Transaction_Type"`
	AgentCode                  string  `bson:"AgentCode" json:"AgentCode"`
	Description                string  `bson:"Description" json:"Description"`
	Award_Type                 string  `bson:"Award_Type" json:"Award_Type"` //"Transaction" or "Amount"
	Amount                     float64 `bson:"Amount" json:"Amount"`
	Points                     float64 `bson:"Points" json:"Points"`
}

type Loyalty_Point_Earning_Rules_Overwrite_AddRequest struct {
	Key                        string  `bson:"Key" json:"Key"` //MSISDN + "|" + MM_Transaction_Type
	Earning_Rules_Overwrite_Id int64   `bson:"Earning_Rules_Overwrite_Id" json:"Earning_Rules_Overwrite_Id"`
	Earning_Rule_Key           string  `bson:"Earning_Rule_Key" json:"Earning_Rule_Key"`
	MM_Transaction_Type        string  `bson:"MM_Transaction_Type" json:"MM_Transaction_Type"`
	AgentCode                  string  `bson:"AgentCode" json:"AgentCode"`
	Description                string  `bson:"Description" json:"Description"`
	Award_Type                 string  `bson:"Award_Type" json:"Award_Type"` //"Transaction" or "Amount"
	Amount                     float64 `bson:"Amount" json:"Amount"`
	Points                     float64 `bson:"Points" json:"Points"`
}

type Loyalty_Point_Earning_Rules_Overwrite_EditRequest struct {
	Key                        string  `bson:"Key" json:"Key"` //MSISDN + "|" + MM_Transaction_Type
	Earning_Rules_Overwrite_Id int64   `bson:"Earning_Rules_Overwrite_Id" json:"Earning_Rules_Overwrite_Id"`
	Earning_Rule_Key           string  `bson:"Earning_Rule_Key" json:"Earning_Rule_Key"`
	MM_Transaction_Type        string  `bson:"MM_Transaction_Type" json:"MM_Transaction_Type"`
	AgentCode                  string  `bson:"AgentCode" json:"AgentCode"`
	Description                string  `bson:"Description" json:"Description"`
	Award_Type                 string  `bson:"Award_Type" json:"Award_Type"` //"Transaction" or "Amount"
	Amount                     float64 `bson:"Amount" json:"Amount"`
	Points                     float64 `bson:"Points" json:"Points"`
}

type Loyalty_Point_Expiry_Rules struct {
	Key                     string    `bson:"Key" json:"Key"`
	Expiry_Rules_Id         int64     `bson:"Expiry_Rules_Id" json:"Expiry_Rules_Id"`
	Description             string    `bson:"Description" json:"Description"`
	Rolling_Expiration      bool      `bson:"Rolling_Expiration" json:"Rolling_Expiration"`
	Validity_Unit           string    `bson:"Validity_Unit" json:"Validity_Unit"`                     //Month, Year --> only when Rolling_Expiration is true
	Validity_Duration       int       `bson:"Validity_Duration" json:"Validity_Duration"`             //only when Rolling_Expiration is true
	Grace_Validity_Unit     string    `bson:"Grace_Validity_Unit" json:"Grace_Validity_Unit"`         //actual expiry unit
	Grace_Validity_Duration int       `bson:"Grace_Validity_Duration" json:"Grace_Validity_Duration"` //actual expiry duration
	Fix_Date_Expiration     bool      `bson:"Fix_Date_Expiration" json:"Fix_Date_Expiration"`
	Expiration_Trigger_date time.Time `bson:"Expiration_Trigger_date" json:"Expiration_Trigger_date"` //when the expiry process will run
	Expiration_Point_Before time.Time `bson:"Expiration_Point_Before" json:"Expiration_Point_Before"` //expiry all points before this date
}

type Loyalty_Point_Expiry_Rules_AddRequest struct {
	Key                     string    `bson:"Key" json:"Key"`
	Expiry_Rules_Id         int64     `bson:"Expiry_Rules_Id" json:"Expiry_Rules_Id"`
	Description             string    `bson:"Description" json:"Description"`
	Rolling_Expiration      bool      `bson:"Rolling_Expiration" json:"Rolling_Expiration"`
	Validity_Unit           string    `bson:"Validity_Unit" json:"Validity_Unit"`                     //Month, Year --> only when Rolling_Expiration is true
	Validity_Duration       int       `bson:"Validity_Duration" json:"Validity_Duration"`             //only when Rolling_Expiration is true
	Grace_Validity_Unit     string    `bson:"Grace_Validity_Unit" json:"Grace_Validity_Unit"`         //actual expiry unit
	Grace_Validity_Duration int       `bson:"Grace_Validity_Duration" json:"Grace_Validity_Duration"` //actual expiry duration
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
	Validity_Unit           string    `bson:"Validity_Unit" json:"Validity_Unit"`                     //Month, Year --> only when Rolling_Expiration is true
	Validity_Duration       int       `bson:"Validity_Duration" json:"Validity_Duration"`             //only when Rolling_Expiration is true
	Grace_Validity_Unit     string    `bson:"Grace_Validity_Unit" json:"Grace_Validity_Unit"`         //actual expiry unit
	Grace_Validity_Duration int       `bson:"Grace_Validity_Duration" json:"Grace_Validity_Duration"` //actual expiry duration
	Fix_Date_Expiration     bool      `bson:"Fix_Date_Expiration" json:"Fix_Date_Expiration"`
	Expiration_Trigger_date time.Time `bson:"Expiration_Trigger_date" json:"Expiration_Trigger_date"` //when the expiry process will run
	Expiration_Point_Before time.Time `bson:"Expiration_Point_Before" json:"Expiration_Point_Before"` //expiry all points before this date
}

type Loyalty_Point_Redemption_Rules struct {
	Key                                 string  `bson:"Key" json:"Key"`
	Redemption_Rules_Id                 int64   `bson:"Redemption_Rules_Id" json:"Redemption_Rules_Id"`
	Description                         string  `bson:"Description" json:"Description"`
	Min_Accumulated_Points              float64 `bson:"Min_Accumulated_Points" json:"Min_Accumulated_Points"` //will be used for upgrade
	Allow_Negative_Balance_ToRedeem     bool    `bson:"Allow_Negative_Balance_ToRedeem" json:"Allow_Negative_Balance_ToRedeem"`
	Allow_PendingLendme_ToRedeem        bool    `bson:"Allow_PendingLendme_ToRedeem" json:"Allow_PendingLendme_ToRedeem"`
	Airtime_MinPoints                   float64 `bson:"Airtime_MinPoints" json:"Airtime_MinPoints"`
	Available_MinPoints_for_Airtime     float64 `bson:"Available_MinPoints_for_Airtime" json:"Available_MinPoints_for_Airtime"`
	Airtime_Amount                      float64 `bson:"Airtime_Amount" json:"Airtime_Amount"`
	Airtime_Points                      float64 `bson:"Airtime_Points" json:"Airtime_Points"`
	Airtime_EVC_Account                 string  `bson:"Airtime_EVC_Account" json:"Airtime_EVC_Account"`
	Airtime_EVC_PIN                     string  `bson:"Airtime_EVC_PIN" json:"Airtime_EVC_PIN"`
	Airtime_Notification                bool    `bson:"Airtime_Notification" json:"Airtime_Notification"`
	Airtime_Notification_Sender         string  `bson:"Airtime_Notification_Sender" json:"Airtime_Notification_Sender"`
	Airtime_Notification_Text           string  `bson:"Airtime_Notification_Text" json:"Airtime_Notification_Text"`
	MobileMoney_MinPoints               float64 `bson:"MobileMoney_MinPoints" json:"MobileMoney_MinPoints"`
	Available_MinPoints_for_MobileMoney float64 `bson:"Available_MinPoints_for_MobileMoney" json:"Available_MinPoints_for_MobileMoney"`
	MobileMoney_Amount                  float64 `bson:"MobileMoney_Amount" json:"MobileMoney_Amount"`
	MobileMoney_Points                  float64 `bson:"MobileMoney_Points" json:"MobileMoney_Points"`
	MobileMoney_MerchantAccount         string  `bson:"MobileMoney_MerchantAccount" json:"MobileMoney_MerchantAccount"`
	MobileMoney_MerchantPIN             string  `bson:"MobileMoney_MerchantPIN" json:"MobileMoney_MerchantPIN"`
	MobileMoney_Notification            bool    `bson:"MobileMoney_Notification" json:"MobileMoney_Notification"`
	MobileMoney_Notification_Sender     string  `bson:"MobileMoney_Notification_Sender" json:"MobileMoney_Notification_Sender"`
	MobileMoney_Notification_Text       string  `bson:"MobileMoney_Notification_Text" json:"MobileMoney_Notification_Text"`
	Bundles_MinPoints                   float64 `bson:"Bundles_MinPoints" json:"Bundles_MinPoints"`
	Bundles_Product_Catalogue_Channel   string  `bson:"Bundles_Product_Catalogue_Channel" json:"Bundles_Product_Catalogue_Channel"`
	Bundles_Product_Catalogue_Plan      string  `bson:"Bundles_Product_Catalogue_Plan" json:"Bundles_Product_Catalogue_Plan"`
	Bundles_Product_Catalogue_Version   string  `bson:"Bundles_Product_Catalogue_Version" json:"Bundles_Product_Catalogue_Version"`
	Bundles_EVC_Account                 string  `bson:"Bundles_EVC_Account" json:"Bundles_EVC_Account"`
	Bundles_EVC_PIN                     string  `bson:"Bundles_EVC_PIN" json:"Bundles_EVC_PIN"`
	Bundles_Notification                bool    `bson:"Bundles_Notification" json:"Bundles_Notification"`
	Bundles_Notification_Sender         string  `bson:"Bundles_Notification_Sender" json:"Bundles_Notification_Sender"`
	Bundles_Notification_Text           string  `bson:"Bundles_Notification_Text" json:"Bundles_Notification_Text"`
	FreeSpinAndWin_MinPoints            float64 `bson:"FreeSpinAndWin_MinPoints" json:"FreeSpinAndWin_MinPoints"`
	Available_MinPoints_for_SpinAndWin  float64 `bson:"Available_MinPoints_for_SpinAndWin" json:"Available_MinPoints_for_SpinAndWin"`
	FreeSpinAndWin_PointsPerSpin        float64 `bson:"FreeSpinAndWin_PointsPerSpin" json:"FreeSpinAndWin_PointsPerSpin"`
	FreeSpinAndWin_Notification         bool    `bson:"FreeSpinAndWin_Notification" json:"FreeSpinAndWin_Notification"`
	FreeSpinAndWin_Notification_Sender  string  `bson:"FreeSpinAndWin_Notification_Sender" json:"FreeSpinAndWin_Notification_Sender"`
	FreeSpinAndWin_Notification_Text    string  `bson:"FreeSpinAndWin_Notification_Text" json:"FreeSpinAndWin_Notification_Text"`
}

type Loyalty_Point_Redemption_Rule struct {
	Key                                 string  `bson:"Key" json:"Key"`
	Redemption_Rules_Id                 int64   `bson:"Redemption_Rules_Id" json:"Redemption_Rules_Id"`
	Description                         string  `bson:"Description" json:"Description"`
	Min_Accumulated_Points              float64 `bson:"Min_Accumulated_Points" json:"Min_Accumulated_Points"` //will be used for upgrade
	Allow_Negative_Balance_ToRedeem     bool    `bson:"Allow_Negative_Balance_ToRedeem" json:"Allow_Negative_Balance_ToRedeem"`
	Allow_PendingLendme_ToRedeem        bool    `bson:"Allow_PendingLendme_ToRedeem" json:"Allow_PendingLendme_ToRedeem"`
	Airtime_MinPoints                   float64 `bson:"Airtime_MinPoints" json:"Airtime_MinPoints"`
	Available_MinPoints_for_Airtime     float64 `bson:"Available_MinPoints_for_Airtime" json:"Available_MinPoints_for_Airtime"`
	Airtime_Amount                      float64 `bson:"Airtime_Amount" json:"Airtime_Amount"`
	Airtime_Points                      float64 `bson:"Airtime_Points" json:"Airtime_Points"`
	Airtime_EVC_Account                 string  `bson:"Airtime_EVC_Account" json:"-"`
	Airtime_EVC_PIN                     string  `bson:"Airtime_EVC_PIN" json:"-"`
	Airtime_Notification                bool    `bson:"Airtime_Notification" json:"Airtime_Notification"`
	Airtime_Notification_Sender         string  `bson:"Airtime_Notification_Sender" json:"Airtime_Notification_Sender"`
	Airtime_Notification_Text           string  `bson:"Airtime_Notification_Text" json:"Airtime_Notification_Text"`
	MobileMoney_MinPoints               float64 `bson:"MobileMoney_MinPoints" json:"MobileMoney_MinPoints"`
	Available_MinPoints_for_MobileMoney float64 `bson:"Available_MinPoints_for_MobileMoney" json:"Available_MinPoints_for_MobileMoney"`
	MobileMoney_Amount                  float64 `bson:"MobileMoney_Amount" json:"MobileMoney_Amount"`
	MobileMoney_Points                  float64 `bson:"MobileMoney_Points" json:"MobileMoney_Points"`
	MobileMoney_MerchantAccount         string  `bson:"MobileMoney_MerchantAccount" json:"-"`
	MobileMoney_MerchantPIN             string  `bson:"MobileMoney_MerchantPIN" json:"-"`
	MobileMoney_Notification            bool    `bson:"MobileMoney_Notification" json:"MobileMoney_Notification"`
	MobileMoney_Notification_Sender     string  `bson:"MobileMoney_Notification_Sender" json:"MobileMoney_Notification_Sender"`
	MobileMoney_Notification_Text       string  `bson:"MobileMoney_Notification_Text" json:"MobileMoney_Notification_Text"`
	Bundles_MinPoints                   float64 `bson:"Bundles_MinPoints" json:"Bundles_MinPoints"`
	Bundles_Product_Catalogue_Channel   string  `bson:"Bundles_Product_Catalogue_Channel" json:"Bundles_Product_Catalogue_Channel"`
	Bundles_Product_Catalogue_Plan      string  `bson:"Bundles_Product_Catalogue_Plan" json:"Bundles_Product_Catalogue_Plan"`
	Bundles_Product_Catalogue_Version   string  `bson:"Bundles_Product_Catalogue_Version" json:"Bundles_Product_Catalogue_Version"`
	Bundles_EVC_Account                 string  `bson:"Bundles_EVC_Account" json:"-"`
	Bundles_EVC_PIN                     string  `bson:"Bundles_EVC_PIN" json:"-"`
	Bundles_Notification                bool    `bson:"Bundles_Notification" json:"Bundles_Notification"`
	Bundles_Notification_Sender         string  `bson:"Bundles_Notification_Sender" json:"Bundles_Notification_Sender"`
	Bundles_Notification_Text           string  `bson:"Bundles_Notification_Text" json:"Bundles_Notification_Text"`
	FreeSpinAndWin_MinPoints            float64 `bson:"FreeSpinAndWin_MinPoints" json:"FreeSpinAndWin_MinPoints"`
	Available_MinPoints_for_SpinAndWin  float64 `bson:"Available_MinPoints_for_SpinAndWin" json:"Available_MinPoints_for_SpinAndWin"`
	FreeSpinAndWin_PointsPerSpin        float64 `bson:"FreeSpinAndWin_PointsPerSpin" json:"FreeSpinAndWin_PointsPerSpin"`
	FreeSpinAndWin_Notification         bool    `bson:"FreeSpinAndWin_Notification" json:"FreeSpinAndWin_Notification"`
	FreeSpinAndWin_Notification_Sender  string  `bson:"FreeSpinAndWin_Notification_Sender" json:"FreeSpinAndWin_Notification_Sender"`
	FreeSpinAndWin_Notification_Text    string  `bson:"FreeSpinAndWin_Notification_Text" json:"FreeSpinAndWin_Notification_Text"`
}

type Loyalty_Point_Redemption_Rules_AddRequest struct {
	Key                                 string  `bson:"Key" json:"Key"`
	Redemption_Rules_Id                 int64   `bson:"Redemption_Rules_Id" json:"Redemption_Rules_Id"`
	Description                         string  `bson:"Description" json:"Description"`
	Min_Accumulated_Points              float64 `bson:"Min_Accumulated_Points" json:"Min_Accumulated_Points"` //will be used for upgrade
	Allow_Negative_Balance_ToRedeem     bool    `bson:"Allow_Negative_Balance_ToRedeem" json:"Allow_Negative_Balance_ToRedeem"`
	Allow_PendingLendme_ToRedeem        bool    `bson:"Allow_PendingLendme_ToRedeem" json:"Allow_PendingLendme_ToRedeem"`
	Airtime_MinPoints                   float64 `bson:"Airtime_MinPoints" json:"Airtime_MinPoints"`
	Available_MinPoints_for_Airtime     float64 `bson:"Available_MinPoints_for_Airtime" json:"Available_MinPoints_for_Airtime"`
	Airtime_Amount                      float64 `bson:"Airtime_Amount" json:"Airtime_Amount"`
	Airtime_Points                      float64 `bson:"Airtime_Points" json:"Airtime_Points"`
	Airtime_EVC_Account                 string  `bson:"Airtime_EVC_Account" json:"Airtime_EVC_Account"`
	Airtime_EVC_PIN                     string  `bson:"Airtime_EVC_PIN" json:"Airtime_EVC_PIN"`
	Airtime_Notification                bool    `bson:"Airtime_Notification" json:"Airtime_Notification"`
	Airtime_Notification_Sender         string  `bson:"Airtime_Notification_Sender" json:"Airtime_Notification_Sender"`
	Airtime_Notification_Text           string  `bson:"Airtime_Notification_Text" json:"Airtime_Notification_Text"`
	MobileMoney_MinPoints               float64 `bson:"MobileMoney_MinPoints" json:"MobileMoney_MinPoints"`
	Available_MinPoints_for_MobileMoney float64 `bson:"Available_MinPoints_for_MobileMoney" json:"Available_MinPoints_for_MobileMoney"`
	MobileMoney_Amount                  float64 `bson:"MobileMoney_Amount" json:"MobileMoney_Amount"`
	MobileMoney_Points                  float64 `bson:"MobileMoney_Points" json:"MobileMoney_Points"`
	MobileMoney_MerchantAccount         string  `bson:"MobileMoney_MerchantAccount" json:"MobileMoney_MerchantAccount"`
	MobileMoney_MerchantPIN             string  `bson:"MobileMoney_MerchantPIN" json:"MobileMoney_MerchantPIN"`
	MobileMoney_Notification            bool    `bson:"MobileMoney_Notification" json:"MobileMoney_Notification"`
	MobileMoney_Notification_Sender     string  `bson:"MobileMoney_Notification_Sender" json:"MobileMoney_Notification_Sender"`
	MobileMoney_Notification_Text       string  `bson:"MobileMoney_Notification_Text" json:"MobileMoney_Notification_Text"`
	Bundles_MinPoints                   float64 `bson:"Bundles_MinPoints" json:"Bundles_MinPoints"`
	Bundles_EVC_Account                 string  `bson:"Bundles_EVC_Account" json:"Bundles_EVC_Account"`
	Bundles_EVC_PIN                     string  `bson:"Bundles_EVC_PIN" json:"Bundles_EVC_PIN"`
	Bundles_Notification                bool    `bson:"Bundles_Notification" json:"Bundles_Notification"`
	Bundles_Notification_Sender         string  `bson:"Bundles_Notification_Sender" json:"Bundles_Notification_Sender"`
	Bundles_Notification_Text           string  `bson:"Bundles_Notification_Text" json:"Bundles_Notification_Text"`
	Bundles_Product_Catalogue_Channel   string  `bson:"Bundles_Product_Catalogue_Channel" json:"Bundles_Product_Catalogue_Channel"`
	Bundles_Product_Catalogue_Plan      string  `bson:"Bundles_Product_Catalogue_Plan" json:"Bundles_Product_Catalogue_Plan"`
	Bundles_Product_Catalogue_Version   string  `bson:"Bundles_Product_Catalogue_Version" json:"Bundles_Product_Catalogue_Version"`
	FreeSpinAndWin_MinPoints            float64 `bson:"FreeSpinAndWin_MinPoints" json:"FreeSpinAndWin_MinPoints"`
	Available_MinPoints_for_SpinAndWin  float64 `bson:"Available_MinPoints_for_SpinAndWin" json:"Available_MinPoints_for_SpinAndWin"`
	FreeSpinAndWin_PointsPerSpin        float64 `bson:"FreeSpinAndWin_PointsPerSpin" json:"FreeSpinAndWin_PointsPerSpin"`
	FreeSpinAndWin_Notification         bool    `bson:"FreeSpinAndWin_Notification" json:"FreeSpinAndWin_Notification"`
	FreeSpinAndWin_Notification_Sender  string  `bson:"FreeSpinAndWin_Notification_Sender" json:"FreeSpinAndWin_Notification_Sender"`
	FreeSpinAndWin_Notification_Text    string  `bson:"FreeSpinAndWin_Notification_Text" json:"FreeSpinAndWin_Notification_Text"`
}

type Loyalty_Point_Redemption_Rules_EditRequest struct {
	Key                                 string  `bson:"Key" json:"Key"`
	NewKey                              string  `bson:"NewKey" json:"NewKey"`
	Redemption_Rules_Id                 int64   `bson:"Redemption_Rules_Id" json:"Redemption_Rules_Id"`
	Description                         string  `bson:"Description" json:"Description"`
	Min_Accumulated_Points              float64 `bson:"Min_Accumulated_Points" json:"Min_Accumulated_Points"` //will be used for upgrade
	Allow_Negative_Balance_ToRedeem     bool    `bson:"Allow_Negative_Balance_ToRedeem" json:"Allow_Negative_Balance_ToRedeem"`
	Allow_PendingLendme_ToRedeem        bool    `bson:"Allow_PendingLendme_ToRedeem" json:"Allow_PendingLendme_ToRedeem"`
	Airtime_MinPoints                   float64 `bson:"Airtime_MinPoints" json:"Airtime_MinPoints"`
	Available_MinPoints_for_Airtime     float64 `bson:"Available_MinPoints_for_Airtime" json:"Available_MinPoints_for_Airtime"`
	Airtime_Amount                      float64 `bson:"Airtime_Amount" json:"Airtime_Amount"`
	Airtime_Points                      float64 `bson:"Airtime_Points" json:"Airtime_Points"`
	Airtime_EVC_Account                 string  `bson:"Airtime_EVC_Account" json:"Airtime_EVC_Account"`
	Airtime_EVC_PIN                     string  `bson:"Airtime_EVC_PIN" json:"Airtime_EVC_PIN"`
	Airtime_Notification                bool    `bson:"Airtime_Notification" json:"Airtime_Notification"`
	Airtime_Notification_Sender         string  `bson:"Airtime_Notification_Sender" json:"Airtime_Notification_Sender"`
	Airtime_Notification_Text           string  `bson:"Airtime_Notification_Text" json:"Airtime_Notification_Text"`
	MobileMoney_MinPoints               float64 `bson:"MobileMoney_MinPoints" json:"MobileMoney_MinPoints"`
	Available_MinPoints_for_MobileMoney float64 `bson:"Available_MinPoints_for_MobileMoney" json:"Available_MinPoints_for_MobileMoney"`
	MobileMoney_Amount                  float64 `bson:"MobileMoney_Amount" json:"MobileMoney_Amount"`
	MobileMoney_Points                  float64 `bson:"MobileMoney_Points" json:"MobileMoney_Points"`
	MobileMoney_MerchantAccount         string  `bson:"MobileMoney_MerchantAccount" json:"MobileMoney_MerchantAccount"`
	MobileMoney_MerchantPIN             string  `bson:"MobileMoney_MerchantPIN" json:"MobileMoney_MerchantPIN"`
	MobileMoney_Notification            bool    `bson:"MobileMoney_Notification" json:"MobileMoney_Notification"`
	MobileMoney_Notification_Sender     string  `bson:"MobileMoney_Notification_Sender" json:"MobileMoney_Notification_Sender"`
	MobileMoney_Notification_Text       string  `bson:"MobileMoney_Notification_Text" json:"MobileMoney_Notification_Text"`
	Bundles_MinPoints                   float64 `bson:"Bundles_MinPoints" json:"Bundles_MinPoints"`
	Bundles_Product_Catalogue_Channel   string  `bson:"Bundles_Product_Catalogue_Channel" json:"Bundles_Product_Catalogue_Channel"`
	Bundles_Product_Catalogue_Plan      string  `bson:"Bundles_Product_Catalogue_Plan" json:"Bundles_Product_Catalogue_Plan"`
	Bundles_Product_Catalogue_Version   string  `bson:"Bundles_Product_Catalogue_Version" json:"Bundles_Product_Catalogue_Version"`
	Bundles_EVC_Account                 string  `bson:"Bundles_EVC_Account" json:"Bundles_EVC_Account"`
	Bundles_EVC_PIN                     string  `bson:"Bundles_EVC_PIN" json:"Bundles_EVC_PIN"`
	Bundles_Notification                bool    `bson:"Bundles_Notification" json:"Bundles_Notification"`
	Bundles_Notification_Sender         string  `bson:"Bundles_Notification_Sender" json:"Bundles_Notification_Sender"`
	Bundles_Notification_Text           string  `bson:"Bundles_Notification_Text" json:"Bundles_Notification_Text"`
	FreeSpinAndWin_MinPoints            float64 `bson:"FreeSpinAndWin_MinPoints" json:"FreeSpinAndWin_MinPoints"`
	Available_MinPoints_for_SpinAndWin  float64 `bson:"Available_MinPoints_for_SpinAndWin" json:"Available_MinPoints_for_SpinAndWin"`
	FreeSpinAndWin_PointsPerSpin        float64 `bson:"FreeSpinAndWin_PointsPerSpin" json:"FreeSpinAndWin_PointsPerSpin"`
	FreeSpinAndWin_Notification         bool    `bson:"FreeSpinAndWin_Notification" json:"FreeSpinAndWin_Notification"`
	FreeSpinAndWin_Notification_Sender  string  `bson:"FreeSpinAndWin_Notification_Sender" json:"FreeSpinAndWin_Notification_Sender"`
	FreeSpinAndWin_Notification_Text    string  `bson:"FreeSpinAndWin_Notification_Text" json:"FreeSpinAndWin_Notification_Text"`
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
type Loyalty_Opt_Request struct {
	EventSource string `bson:"EventSource" json:"EventSource"`
	MSISDN      string `bson:"MSISDN" json:"MSISDN"`         //MSISDN
	Opt_Status  string `bson:"Opt_Status" json:"Opt_Status"` //Optn
}
type Loyalty_Status_log struct {
	//request Header info
	SourceIP   string `bson:"SourceIP" json:"-"`
	SourceApp  string `bson:"SourceApp" json:"-"`
	AppLogin   string `bson:"AppLogin" json:"-"`
	AppVersion string `bson:"AppVersion" json:"-"`

	Opt_Status string `bson:"Opt_Status" json:"Opt_Status"`
	MSISDN     string `bson:"MSISDN" json:"MSISDN"`

	Request_Status     string    `bson:"Request_Status" json:"Request_Status"`
	Request_StatusCode int       `bson:"Request_StatusCode" json:"Request_StatusCode"`
	StatusDescription  string    `bson:"StatusDescription" json:"StatusDescription"`
	ErrorDescription   string    `bson:"ErrorDescription" json:"ErrorDescription"`
	StatusDate         time.Time `bson:"StatusDate" json:"StatusDate"`
	E2E_Elapsedtime    int64     `bson:"E2E_Elapsedtime" json:"E2E_Elapsedtime"` //receive date till return
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
	EventDetailCode  string  `bson:"EventDetailCode" json:"EventDetailCode"`
	EventAmount      float64 `bson:"EventAmount" json:"EventAmount"`
	PointsToCredit   float64 `bson:"PointsToCredit" json:"PointsToCredit"`
	EventDescription string  `bson:"EventDescription" json:"EventDescription"`
}

type Loyalty_Logs struct {
	Logs_Type                string                          `bson:"Logs_Type" json:"Logs_Type"`
	Date                     time.Time                       `bson:"Date" json:"Date"`
	Status                   string                          `bson:"Status" json:"Status"` //successful, failed
	PointsToCredit           float64                         `bson:"PointsToCredit" json:"PointsToCredit"`
	Opening_Available_Points float64                         `bson:"Opening_Available_Points" json:"Opening_Available_Points"` //(Awarded_Points + Expired_Points) - Redeemed_Points
	Closure_Available_Points float64                         `bson:"Closure_Available_Points" json:"Closure_Available_Points"` //(Awarded_Points + Expired_Points) - Redeemed_Points
	Debit_Logs               Loyalty_AccountDebitPoints_log  `bson:"Debit_Logs" json:"Debit_Logs"`
	Credit_Logs              Loyalty_AccountCreditPoints_log `bson:"Credit_Logs" json:"Credit_Logs"`
	Redemption_Logs          Loyalty_Redemption_log          `bson:"Redemption_Logs" json:"Redemption_Logs"`
	Expiry_Logs              Loyalty_Monthly_Expiry_log      `bson:"Expiry_Logs" json:"Expiry_Logs"`
	Level_Change_Logs        Loyalty_Level_Change_log        `bson:"Level_Change_Logs" json:"Level_Change_Logs"`
	Status_Logs              Loyalty_Status_log              `bson:"Status_Logs" json:"Status_Logs"`
	Full_Expiry_Logs         Loyalty_Full_Expiry_Log         `bson:"Full_Expiry_Logs" json:"Full_Expiry_Logs"`
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
	MSISDN           string  `bson:"MSISDN" json:"MSISDN"`                   //MSISDN
	EventSource      string  `bson:"EventSource" json:"EventSource"`         //MobileApp, MobileMoney, USSD,...
	EventType        string  `bson:"EventType" json:"EventType"`             //BundlePurchase, MOC,...
	EventDetail      string  `bson:"EventDetail" json:"EventDetail"`         //BundleName, ...
	EventDetailCode  string  `bson:"EventDetailCode" json:"EventDetailCode"` //agent code
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

	Opening_Outstanding_fraction_points float64 `bson:"Opening_Outstanding_fraction_points" json:"Opening_Outstanding_fraction_points"`
	Closure_Outstanding_fraction_points float64 `bson:"Closure_Outstanding_fraction_points" json:"Closure_Outstanding_fraction_points"`

	//response result
	ReceiveDate       time.Time `bson:"ReceiveDate" json:"-"`
	Status            string    `bson:"Status" json:"Status"` //successful, failed
	StatusCode        int       `bson:"StatusCode" json:"StatusCode"`
	StatusDescription string    `bson:"StatusDescription" json:"StatusDescription"`
	ErrorDescription  string    `bson:"ErrorDescription" json:"ErrorDescription"`
	StatusDate        time.Time `bson:"StatusDate" json:"StatusDate"`
	E2E_Elapsedtime   int64     `bson:"E2E_Elapsedtime" json:"E2E_Elapsedtime"` //receive date till return
}

type Loyalty_Monthly_Expiry_log struct {
	ExpiryTime       time.Time `bson:"ExpiryTime" json:"ExpiryTime"`
	MSISDN           string    `bson:"MSISDN" json:"MSISDN"` //MSISDN
	Expiry_Rules_Key string    `bson:"Expiry_Rules_Key" json:"Expiry_Rules_Key"`

	Year_Month string `bson:"Year_Month" json:"Year_Month"`

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

type Loyalty_Full_Expiry_Log struct {
	MSISDN string `bson:"MSISDN" json:"MSISDN"` //MSISDN

	Opening_Awarded_Points     float64 `bson:"Opening_Awarded_Points" json:"Opening_Awarded_Points"`
	Opening_Redeemed_Points    float64 `bson:"Opening_Redeemed_Points" json:"Opening_Redeemed_Points"`
	Opening_Available_Points   float64 `bson:"Opening_Available_Points" json:"Opening_Available_Points"`     //(Awarded_Points + Expired_Points) - Redeemed_Points
	Opening_Expired_Points     float64 `bson:"Opening_Expired_Points" json:"Opening_Expired_Points"`         //expired are deducted from Awarded_Points
	Opening_OutStanding_Points float64 `bson:"Opening_OutStanding_Points" json:"Opening_OutStanding_Points"` //

	End_Awarded_Points     float64 `bson:"End_Awarded_Points" json:"End_Awarded_Points"`
	End_Redeemed_Points    float64 `bson:"End_Redeemed_Points" json:"End_Redeemed_Points"`
	End_Available_Points   float64 `bson:"End_Available_Points" json:"End_Available_Points"`     //(Awarded_Points + Expired_Points) - Redeemed_Points
	End_Expired_Points     float64 `bson:"End_Expired_Points" json:"End_Expired_Points"`         //expired are deducted from Awarded_Points
	End_Outstanding_Points float64 `bson:"End_Outstanding_Points" json:"End_Outstanding_Points"` //

	OpeningLoyaltyLevel string `bson:"OpeningLoyaltyLevel" json:"OpeningLoyaltyLevel"`
	EndLoyaltyLevel     string `bson:"EndLoyaltyLevel" json:"EndLoyaltyLevel"`

	Grace_Period_Given_Days int       `bson:"Grace_Period_Given_Days" json:"Grace_Period_Given_Days"`
	Last_OptOut             time.Time `bson:"Last_OptOut" json:"Last_OptOut"`

	ExpiryAmount            float64   `bson:"ExpiryAmount" json:"ExpiryAmount"`
	ExpiryTime              time.Time `bson:"ExpiryTime" json:"ExpiryTime"`
	ExpiryReason            string    `bson:"ExpiryReason" json:"ExpiryReason"`
	ExpiryStatus            string    `bson:"ExpiryStatus" json:"ExpiryStatus"`
	ExpiryStatusDescription string    `bson:"ExpiryStatusDescription" json:"ExpiryStatusDescription"`
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
	MinAvailableRequiredPoints      float64 `bson:"MinAvailableRequiredPoints" json:"MinAvailableRequiredPoints"`
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

type Loyalty_Subs_Summary struct {
	ID struct {
		Loyalty_Level_Key string
	} `bson:"_id"`
	Accounts_Count         float64 `bson:"Accounts_Count"`
	Total_Awarded_Points   float64 `bson:"Total_Awarded_Points"`
	Total_Redeemed_Points  float64 `bson:"Total_Redeemed_Points"`
	Total_Available_Points float64 `bson:"Total_Available_Points"`
}
type CustomersUploadList struct {
	MSISDN string  `csv:"MSISDN"`
	Points float64 `csv:"POINTS"`
}
type CSVMSISDN struct {
	MSISDN string
	Exists bool
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
	Key                  string    `bson:"Key" json:"Key"` //Campaign Name
	Campaign_Id          int64     `bson:"Campaign_Id" json:"Campaign_Id"`
	Description          string    `bson:"Description" json:"Description"`
	Add_Date             time.Time `bson:"Add_Date" json:"Add_Date"`
	Start_Date           time.Time `bson:"Start_Date" json:"Start_Date"`
	End_Date             time.Time `bson:"End_Date" json:"End_Date"`
	Campaign_Status      string    `bson:"Campaign_Status" json:"Campaign_Status"` //Created, Launched, Paused, Resumed, Cancelled, or Stopped
	Campaign_Status_Date time.Time `bson:"Campaign_Status_Date" json:"Campaign_Status_Date"`
	Campaign_Status_User string    `bson:"Campaign_Status_User" json:"Campaign_Status_User"`

	//target rules
	Target_All_Subs    bool            `bson:"Target_All_Subs" json:"Target_All_Subs"` //if true, all target features will be ignored
	Target_Level_Key   map[string]bool `bson:"Target_Level_Key" json:"Target_Level_Key"`
	Target_Segment_Key map[string]bool `bson:"Target_Segment_Key" json:"Target_Segment_Key"`
	LoyaltyPoints_From float64         `bson:"LoyaltyPoints_From" json:"LoyaltyPoints_From"`
	LoyaltyPoints_Till float64         `bson:"LoyaltyPoints_Till" json:"LoyaltyPoints_Till"`
	AON_From           float64         `bson:"AON_From" json:"AON_From"` //months
	AON_Till           float64         `bson:"AON_Till" json:"AON_Till"` //months
	ARPU_From          float64         `bson:"ARPU_From" json:"ARPU_From"`
	ARPU_Till          float64         `bson:"ARPU_Till" json:"ARPU_Till"`

	//campaign award rules
	Welcome_Points_Award_type string  `bson:"Welcome_Points_Award_type" json:"Welcome_Points_Award_type"` //Multiplier or Fixed Amount
	Welcome_Points_Frequency  string  `bson:"Welcome_Points_Frequency" json:"Welcome_Points_Frequency"`   //Once, Once Daily, Each time
	Welcome_Points_Max_Award  float64 `bson:"Welcome_Points_Max_Award" json:"Welcome_Points_Max_Award"`
	Welcome_Points            float64 `bson:"Welcome_Points" json:"Welcome_Points"`

	MobileAppDaily_Login_Award_type string  `bson:"MobileAppDaily_Login_Award_type" json:"MobileAppDaily_Login_Award_type"` //Multiplier or Fixed Amount
	MobileAppDaily_Login_Frequency  string  `bson:"MobileAppDaily_Login_Frequency" json:"MobileAppDaily_Login_Frequency"`   //Once, Once Daily, Each time
	MobileAppDaily_Login_Max_Award  float64 `bson:"MobileAppDaily_Login_Max_Award" json:"MobileAppDaily_Login_Max_Award"`
	MobileAppDaily_Login            float64 `bson:"MobileAppDaily_Login" json:"MobileAppDaily_Login"`

	MainGSM_Award_type string  `bson:"MainGSM_Award_type" json:"MainGSM_Award_type"` //Multiplier or Fixed Amount
	MainGSM_Frequency  string  `bson:"MainGSM_Frequency" json:"MainGSM_Frequency"`   //Once, Once Daily, Each time
	MainGSM_Max_Award  float64 `bson:"MainGSM_Max_Award" json:"MainGSM_Max_Award"`
	MainGSM            float64 `bson:"MainGSM" json:"MainGSM"`

	MM_P2P_Award_type string  `bson:"MM_P2P_Award_type" json:"MM_P2P_Award_type"` //Multiplier or Fixed Amount
	MM_P2P_Frequency  string  `bson:"MM_P2P_Frequency" json:"MM_P2P_Frequency"`   //Once, Once Daily, Each time
	MM_P2P_Max_Award  float64 `bson:"MM_P2P_Max_Award" json:"MM_P2P_Max_Award"`
	MM_P2P            float64 `bson:"MM_P2P" json:"MM_P2P"`

	MM_CASHIN_Award_type string  `bson:"MM_CASHIN_Award_type" json:"MM_CASHIN_Award_type"` //Multiplier or Fixed Amount
	MM_CASHIN_Frequency  string  `bson:"MM_CASHIN_Frequency" json:"MM_CASHIN_Frequency"`   //Once, Once Daily, Each time
	MM_CASHIN_Max_Award  float64 `bson:"MM_CASHIN_Max_Award" json:"MM_CASHIN_Max_Award"`
	MM_CASHIN            float64 `bson:"MM_CASHIN" json:"MM_CASHIN"`

	MM_CASHOUT_Award_type string  `bson:"MM_CASHOUT_Award_type" json:"MM_CASHOUT_Award_type"` //Multiplier or Fixed Amount
	MM_CASHOUT_Frequency  string  `bson:"MM_CASHOUT_Frequency" json:"MM_CASHOUT_Frequency"`   //Once, Once Daily, Each time
	MM_CASHOUT_Max_Award  float64 `bson:"MM_CASHOUT_Max_Award" json:"MM_CASHOUT_Max_Award"`
	MM_CASHOUT            float64 `bson:"MM_CASHOUT" json:"MM_CASHOUT"`

	MM_MERCHPAY_Award_type string  `bson:"MM_MERCHPAY_Award_type" json:"MM_MERCHPAY_Award_type"` //Multiplier or Fixed Amount
	MM_MERCHPAY_Frequency  string  `bson:"MM_MERCHPAY_Frequency" json:"MM_MERCHPAY_Frequency"`   //Once, Once Daily, Each time
	MM_MERCHPAY_Max_Award  float64 `bson:"MM_MERCHPAY_Max_Award" json:"MM_MERCHPAY_Max_Award"`
	MM_MERCHPAY            float64 `bson:"MM_MERCHPAY" json:"MM_MERCHPAY"`

	MM_BILLPAY_Award_type string  `bson:"MM_BILLPAY_Award_type" json:"MM_BILLPAY_Award_type"` //Multiplier or Fixed Amount
	MM_BILLPAY_Frequency  string  `bson:"MM_BILLPAY_Frequency" json:"MM_BILLPAY_Frequency"`   //Once, Once Daily, Each time
	MM_BILLPAY_Max_Award  float64 `bson:"MM_BILLPAY_Max_Award" json:"MM_BILLPAY_Max_Award"`
	MM_BILLPAY            float64 `bson:"MM_BILLPAY" json:"MM_BILLPAY"`

	MM_RC_Award_type string  `bson:"MM_RC_Award_type" json:"MM_RC_Award_type"` //Multiplier or Fixed Amount
	MM_RC_Frequency  string  `bson:"MM_RC_Frequency" json:"MM_RC_Frequency"`   //Once, Once Daily, Each time
	MM_RC_Max_Award  float64 `bson:"MM_RC_Max_Award" json:"MM_RC_Max_Award"`
	MM_RC            float64 `bson:"MM_RC" json:"MM_RC"`

	MM_CTMMOREQ_Award_type string  `bson:"MM_CTMMOREQ_Award_type" json:"MM_CTMMOREQ_Award_type"` //Multiplier or Fixed Amount
	MM_Frequency           string  `bson:"MM_Frequency" json:"MM_Frequency"`                     //Once, Once Daily, Each time
	MM_Max_Award           float64 `bson:"MM_Max_Award" json:"MM_Max_Award"`
	MM_CTMMOREQ            float64 `bson:"MM_CTMMOREQ" json:"MM_CTMMOREQ"`

	MM_CBWREQ_Award_type string  `bson:"MM_CBWREQ_Award_type" json:"MM_CBWREQ_Award_type"` //Multiplier or Fixed Amount
	MM_CBWREQ_Frequency  string  `bson:"MM_CBWREQ_Frequency" json:"MM_CBWREQ_Frequency"`   //Once, Once Daily, Each time
	MM_CBWREQ_Max_Award  float64 `bson:"MM_CBWREQ_Max_Award" json:"MM_CBWREQ_Max_Award"`
	MM_CBWREQ            float64 `bson:"MM_CBWREQ" json:"MM_CBWREQ"`

	//campaign notification
	Invitation_SMS_Sender string `bson:"Invitation_SMS_Sender" json:"Invitation_SMS_Sender"`
	Invitation_SMS_Text   string `bson:"Invitation_SMS_Text" json:"Invitation_SMS_Text"`
	PointsAward_SMS_Text  string `bson:"PointsAward_SMS_Text" json:"PointsAward_SMS_Text"`
}

type Loyalty_Campaign_AddRequest struct {
	Key         string    `bson:"Key" json:"Key"` //Campaign Name
	Campaign_Id int64     `bson:"Campaign_Id" json:"Campaign_Id"`
	Description string    `bson:"Description" json:"Description"`
	Start_Date  time.Time `bson:"Start_Date" json:"Start_Date"`
	End_Date    time.Time `bson:"End_Date" json:"End_Date"`
	//target rules
	Target_All_Subs    bool            `bson:"Target_All_Subs" json:"Target_All_Subs"` //if true, all target features will be ignored
	Target_Level_Key   map[string]bool `bson:"Target_Level_Key" json:"Target_Level_Key"`
	Target_Segment_Key map[string]bool `bson:"Target_Segment_Key" json:"Target_Segment_Key"`
	LoyaltyPoints_From float64         `bson:"LoyaltyPoints_From" json:"LoyaltyPoints_From"`
	LoyaltyPoints_Till float64         `bson:"LoyaltyPoints_Till" json:"LoyaltyPoints_Till"`
	AON_From           float64         `bson:"AON_From" json:"AON_From"` //months
	AON_Till           float64         `bson:"AON_Till" json:"AON_Till"` //months
	ARPU_From          float64         `bson:"ARPU_From" json:"ARPU_From"`
	ARPU_Till          float64         `bson:"ARPU_Till" json:"ARPU_Till"`

	//campaign award rules
	Welcome_Points_Award_type string  `bson:"Welcome_Points_Award_type" json:"Welcome_Points_Award_type"` //Multiplier or Fixed Amount
	Welcome_Points_Frequency  string  `bson:"Welcome_Points_Frequency" json:"Welcome_Points_Frequency"`   //Once, Once Daily, Each time
	Welcome_Points_Max_Award  float64 `bson:"Welcome_Points_Max_Award" json:"Welcome_Points_Max_Award"`
	Welcome_Points            float64 `bson:"Welcome_Points" json:"Welcome_Points"`

	MobileAppDaily_Login_Award_type string  `bson:"MobileAppDaily_Login_Award_type" json:"MobileAppDaily_Login_Award_type"` //Multiplier or Fixed Amount
	MobileAppDaily_Login_Frequency  string  `bson:"MobileAppDaily_Login_Frequency" json:"MobileAppDaily_Login_Frequency"`   //Once, Once Daily, Each time
	MobileAppDaily_Login_Max_Award  float64 `bson:"MobileAppDaily_Login_Max_Award" json:"MobileAppDaily_Login_Max_Award"`
	MobileAppDaily_Login            float64 `bson:"MobileAppDaily_Login" json:"MobileAppDaily_Login"`

	MainGSM_Award_type string  `bson:"MainGSM_Award_type" json:"MainGSM_Award_type"` //Multiplier or Fixed Amount
	MainGSM_Frequency  string  `bson:"MainGSM_Frequency" json:"MainGSM_Frequency"`   //Once, Once Daily, Each time
	MainGSM_Max_Award  float64 `bson:"MainGSM_Max_Award" json:"MainGSM_Max_Award"`
	MainGSM            float64 `bson:"MainGSM" json:"MainGSM"`

	MM_P2P_Award_type string  `bson:"MM_P2P_Award_type" json:"MM_P2P_Award_type"` //Multiplier or Fixed Amount
	MM_P2P_Frequency  string  `bson:"MM_P2P_Frequency" json:"MM_P2P_Frequency"`   //Once, Once Daily, Each time
	MM_P2P_Max_Award  float64 `bson:"MM_P2P_Max_Award" json:"MM_P2P_Max_Award"`
	MM_P2P            float64 `bson:"MM_P2P" json:"MM_P2P"`

	MM_CASHIN_Award_type string  `bson:"MM_CASHIN_Award_type" json:"MM_CASHIN_Award_type"` //Multiplier or Fixed Amount
	MM_CASHIN_Frequency  string  `bson:"MM_CASHIN_Frequency" json:"MM_CASHIN_Frequency"`   //Once, Once Daily, Each time
	MM_CASHIN_Max_Award  float64 `bson:"MM_CASHIN_Max_Award" json:"MM_CASHIN_Max_Award"`
	MM_CASHIN            float64 `bson:"MM_CASHIN" json:"MM_CASHIN"`

	MM_CASHOUT_Award_type string  `bson:"MM_CASHOUT_Award_type" json:"MM_CASHOUT_Award_type"` //Multiplier or Fixed Amount
	MM_CASHOUT_Frequency  string  `bson:"MM_CASHOUT_Frequency" json:"MM_CASHOUT_Frequency"`   //Once, Once Daily, Each time
	MM_CASHOUT_Max_Award  float64 `bson:"MM_CASHOUT_Max_Award" json:"MM_CASHOUT_Max_Award"`
	MM_CASHOUT            float64 `bson:"MM_CASHOUT" json:"MM_CASHOUT"`

	MM_MERCHPAY_Award_type string  `bson:"MM_MERCHPAY_Award_type" json:"MM_MERCHPAY_Award_type"` //Multiplier or Fixed Amount
	MM_MERCHPAY_Frequency  string  `bson:"MM_MERCHPAY_Frequency" json:"MM_MERCHPAY_Frequency"`   //Once, Once Daily, Each time
	MM_MERCHPAY_Max_Award  float64 `bson:"MM_MERCHPAY_Max_Award" json:"MM_MERCHPAY_Max_Award"`
	MM_MERCHPAY            float64 `bson:"MM_MERCHPAY" json:"MM_MERCHPAY"`

	MM_BILLPAY_Award_type string  `bson:"MM_BILLPAY_Award_type" json:"MM_BILLPAY_Award_type"` //Multiplier or Fixed Amount
	MM_BILLPAY_Frequency  string  `bson:"MM_BILLPAY_Frequency" json:"MM_BILLPAY_Frequency"`   //Once, Once Daily, Each time
	MM_BILLPAY_Max_Award  float64 `bson:"MM_BILLPAY_Max_Award" json:"MM_BILLPAY_Max_Award"`
	MM_BILLPAY            float64 `bson:"MM_BILLPAY" json:"MM_BILLPAY"`

	MM_RC_Award_type string  `bson:"MM_RC_Award_type" json:"MM_RC_Award_type"` //Multiplier or Fixed Amount
	MM_RC_Frequency  string  `bson:"MM_RC_Frequency" json:"MM_RC_Frequency"`   //Once, Once Daily, Each time
	MM_RC_Max_Award  float64 `bson:"MM_RC_Max_Award" json:"MM_RC_Max_Award"`
	MM_RC            float64 `bson:"MM_RC" json:"MM_RC"`

	MM_CTMMOREQ_Award_type string  `bson:"MM_CTMMOREQ_Award_type" json:"MM_CTMMOREQ_Award_type"` //Multiplier or Fixed Amount
	MM_Frequency           string  `bson:"MM_Frequency" json:"MM_Frequency"`                     //Once, Once Daily, Each time
	MM_Max_Award           float64 `bson:"MM_Max_Award" json:"MM_Max_Award"`
	MM_CTMMOREQ            float64 `bson:"MM_CTMMOREQ" json:"MM_CTMMOREQ"`

	MM_CBWREQ_Award_type string  `bson:"MM_CBWREQ_Award_type" json:"MM_CBWREQ_Award_type"` //Multiplier or Fixed Amount
	MM_CBWREQ_Frequency  string  `bson:"MM_CBWREQ_Frequency" json:"MM_CBWREQ_Frequency"`   //Once, Once Daily, Each time
	MM_CBWREQ_Max_Award  float64 `bson:"MM_CBWREQ_Max_Award" json:"MM_CBWREQ_Max_Award"`
	MM_CBWREQ            float64 `bson:"MM_CBWREQ" json:"MM_CBWREQ"`

	//campaign notification
	Invitation_SMS_Sender string `bson:"Invitation_SMS_Sender" json:"Invitation_SMS_Sender"`
	Invitation_SMS_Text   string `bson:"Invitation_SMS_Text" json:"Invitation_SMS_Text"`
	PointsAward_SMS_Text  string `bson:"PointsAward_SMS_Text" json:"PointsAward_SMS_Text"`
}

type Loyalty_Campaign_EditRequest struct {
	Key                  string    `bson:"Key" json:"Key"` //Campaign Name
	NewKey               string    `bson:"NewKey" json:"NewKey"`
	Campaign_Id          int64     `bson:"Campaign_Id" json:"Campaign_Id"`
	Description          string    `bson:"Description" json:"Description"`
	Start_Date           time.Time `bson:"Start_Date" json:"Start_Date"`
	End_Date             time.Time `bson:"End_Date" json:"End_Date"`
	Campaign_Status      string    `bson:"Campaign_Status" json:"Campaign_Status"` //launched, paused, resumed, cancelled, or stopped
	Campaign_Status_Date time.Time `bson:"Campaign_Status_Date" json:"Campaign_Status_Date"`
	Campaign_Status_User string    `bson:"Campaign_Status_User" json:"Campaign_Status_User"`

	//target rules
	Target_All_Subs    bool            `bson:"Target_All_Subs" json:"Target_All_Subs"` //if true, all target features will be ignored
	Target_Level_Key   map[string]bool `bson:"Target_Level_Key" json:"Target_Level_Key"`
	Target_Segment_Key map[string]bool `bson:"Target_Segment_Key" json:"Target_Segment_Key"`
	LoyaltyPoints_From float64         `bson:"LoyaltyPoints_From" json:"LoyaltyPoints_From"`
	LoyaltyPoints_Till float64         `bson:"LoyaltyPoints_Till" json:"LoyaltyPoints_Till"`
	AON_From           float64         `bson:"AON_From" json:"AON_From"` //months
	AON_Till           float64         `bson:"AON_Till" json:"AON_Till"` //months
	ARPU_From          float64         `bson:"ARPU_From" json:"ARPU_From"`
	ARPU_Till          float64         `bson:"ARPU_Till" json:"ARPU_Till"`

	//campaign award rules
	Welcome_Points_Award_type string  `bson:"Welcome_Points_Award_type" json:"Welcome_Points_Award_type"` //Multiplier or Fixed Amount
	Welcome_Points_Frequency  string  `bson:"Welcome_Points_Frequency" json:"Welcome_Points_Frequency"`   //Once, Once Daily, Each time
	Welcome_Points_Max_Award  float64 `bson:"Welcome_Points_Max_Award" json:"Welcome_Points_Max_Award"`
	Welcome_Points            float64 `bson:"Welcome_Points" json:"Welcome_Points"`

	MobileAppDaily_Login_Award_type string  `bson:"MobileAppDaily_Login_Award_type" json:"MobileAppDaily_Login_Award_type"` //Multiplier or Fixed Amount
	MobileAppDaily_Login_Frequency  string  `bson:"MobileAppDaily_Login_Frequency" json:"MobileAppDaily_Login_Frequency"`   //Once, Once Daily, Each time
	MobileAppDaily_Login_Max_Award  float64 `bson:"MobileAppDaily_Login_Max_Award" json:"MobileAppDaily_Login_Max_Award"`
	MobileAppDaily_Login            float64 `bson:"MobileAppDaily_Login" json:"MobileAppDaily_Login"`

	MainGSM_Award_type string  `bson:"MainGSM_Award_type" json:"MainGSM_Award_type"` //Multiplier or Fixed Amount
	MainGSM_Frequency  string  `bson:"MainGSM_Frequency" json:"MainGSM_Frequency"`   //Once, Once Daily, Each time
	MainGSM_Max_Award  float64 `bson:"MainGSM_Max_Award" json:"MainGSM_Max_Award"`
	MainGSM            float64 `bson:"MainGSM" json:"MainGSM"`

	MM_P2P_Award_type string  `bson:"MM_P2P_Award_type" json:"MM_P2P_Award_type"` //Multiplier or Fixed Amount
	MM_P2P_Frequency  string  `bson:"MM_P2P_Frequency" json:"MM_P2P_Frequency"`   //Once, Once Daily, Each time
	MM_P2P_Max_Award  float64 `bson:"MM_P2P_Max_Award" json:"MM_P2P_Max_Award"`
	MM_P2P            float64 `bson:"MM_P2P" json:"MM_P2P"`

	MM_CASHIN_Award_type string  `bson:"MM_CASHIN_Award_type" json:"MM_CASHIN_Award_type"` //Multiplier or Fixed Amount
	MM_CASHIN_Frequency  string  `bson:"MM_CASHIN_Frequency" json:"MM_CASHIN_Frequency"`   //Once, Once Daily, Each time
	MM_CASHIN_Max_Award  float64 `bson:"MM_CASHIN_Max_Award" json:"MM_CASHIN_Max_Award"`
	MM_CASHIN            float64 `bson:"MM_CASHIN" json:"MM_CASHIN"`

	MM_CASHOUT_Award_type string  `bson:"MM_CASHOUT_Award_type" json:"MM_CASHOUT_Award_type"` //Multiplier or Fixed Amount
	MM_CASHOUT_Frequency  string  `bson:"MM_CASHOUT_Frequency" json:"MM_CASHOUT_Frequency"`   //Once, Once Daily, Each time
	MM_CASHOUT_Max_Award  float64 `bson:"MM_CASHOUT_Max_Award" json:"MM_CASHOUT_Max_Award"`
	MM_CASHOUT            float64 `bson:"MM_CASHOUT" json:"MM_CASHOUT"`

	MM_MERCHPAY_Award_type string  `bson:"MM_MERCHPAY_Award_type" json:"MM_MERCHPAY_Award_type"` //Multiplier or Fixed Amount
	MM_MERCHPAY_Frequency  string  `bson:"MM_MERCHPAY_Frequency" json:"MM_MERCHPAY_Frequency"`   //Once, Once Daily, Each time
	MM_MERCHPAY_Max_Award  float64 `bson:"MM_MERCHPAY_Max_Award" json:"MM_MERCHPAY_Max_Award"`
	MM_MERCHPAY            float64 `bson:"MM_MERCHPAY" json:"MM_MERCHPAY"`

	MM_BILLPAY_Award_type string  `bson:"MM_BILLPAY_Award_type" json:"MM_BILLPAY_Award_type"` //Multiplier or Fixed Amount
	MM_BILLPAY_Frequency  string  `bson:"MM_BILLPAY_Frequency" json:"MM_BILLPAY_Frequency"`   //Once, Once Daily, Each time
	MM_BILLPAY_Max_Award  float64 `bson:"MM_BILLPAY_Max_Award" json:"MM_BILLPAY_Max_Award"`
	MM_BILLPAY            float64 `bson:"MM_BILLPAY" json:"MM_BILLPAY"`

	MM_RC_Award_type string  `bson:"MM_RC_Award_type" json:"MM_RC_Award_type"` //Multiplier or Fixed Amount
	MM_RC_Frequency  string  `bson:"MM_RC_Frequency" json:"MM_RC_Frequency"`   //Once, Once Daily, Each time
	MM_RC_Max_Award  float64 `bson:"MM_RC_Max_Award" json:"MM_RC_Max_Award"`
	MM_RC            float64 `bson:"MM_RC" json:"MM_RC"`

	MM_CTMMOREQ_Award_type string  `bson:"MM_CTMMOREQ_Award_type" json:"MM_CTMMOREQ_Award_type"` //Multiplier or Fixed Amount
	MM_Frequency           string  `bson:"MM_Frequency" json:"MM_Frequency"`                     //Once, Once Daily, Each time
	MM_Max_Award           float64 `bson:"MM_Max_Award" json:"MM_Max_Award"`
	MM_CTMMOREQ            float64 `bson:"MM_CTMMOREQ" json:"MM_CTMMOREQ"`

	MM_CBWREQ_Award_type string  `bson:"MM_CBWREQ_Award_type" json:"MM_CBWREQ_Award_type"` //Multiplier or Fixed Amount
	MM_CBWREQ_Frequency  string  `bson:"MM_CBWREQ_Frequency" json:"MM_CBWREQ_Frequency"`   //Once, Once Daily, Each time
	MM_CBWREQ_Max_Award  float64 `bson:"MM_CBWREQ_Max_Award" json:"MM_CBWREQ_Max_Award"`
	MM_CBWREQ            float64 `bson:"MM_CBWREQ" json:"MM_CBWREQ"`

	//campaign notification
	Invitation_SMS_Sender string `bson:"Invitation_SMS_Sender" json:"Invitation_SMS_Sender"`
	Invitation_SMS_Text   string `bson:"Invitation_SMS_Text" json:"Invitation_SMS_Text"`
	PointsAward_SMS_Text  string `bson:"PointsAward_SMS_Text" json:"PointsAward_SMS_Text"`
}

type Loyalty_Campaign_Target_List struct {
	Key          string    `bson:"Key" json:"Key"` //Campaign_Key|MSISDN
	Campaign_Key int64     `bson:"Campaign_Key" json:"Campaign_Key"`
	MSISDN       string    `bson:"MSISDN" json:"MSISDN"`
	Start_Date   time.Time `bson:"Start_Date" json:"Start_Date"`
	End_Date     time.Time `bson:"End_Date" json:"End_Date"`
}

type Loyalty_Campaign_Account struct {
	Key                      string    `bson:"Key" json:"Key"` //Campaign_Key|MSISDN
	Campaign_Key             int64     `bson:"Campaign_Key" json:"Campaign_Key"`
	MSISDN                   string    `bson:"MSISDN" json:"MSISDN"`
	Last_Award_Date          time.Time `bson:"Last_Award_Date" json:"Last_Award_Date"` //the last date when subscriber get benefit from the campaign
	Cumulative_Points_Earned float64   `bson:"Cumulative_Points_Earned" json:"Cumulative_Points_Earned"`
}

type NotificationLog struct {
	SourceAction  string      `bson:"SourceAction" json:"SourceAction"`
	TransactionId string      `bson:"TransactionId" json:"TransactionId"`
	Medium        string      `bson:"Medium" json:"Medium"`
	SourceAddress string      `bson:"SourceAddress" json:"SourceAddress"`
	Destination   string      `bson:"Destination" json:"Destination"`
	Payload       interface{} `bson:"Payload" json:"Payload"`
	Subject       string      `bson:"Subject" json:"Subject"`
	Status        string      `bson:"Status" json:"Status"`
	Error         interface{} `bson:"Error" json:"Error"`
	AddUser       string      `bson:"AddUser" json:"AddUser"`
	AddDate       time.Time   `bson:"AddDate" json:"AddDate"`
}

type JobStatus struct {
	TotalRows     int
	ProcessedRows int
	Status        string // "running", "completed", "failed"
	Errors        []string
	Result        map[string][]string
}
