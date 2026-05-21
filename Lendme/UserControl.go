package Lendme

import (
	SpinAndWin_client "afr_SpinAndWin_be/SpinAndWinClient"
	APGW "afr_ao_apgw_v2/APGWClientV2"
	AuthCenterClient "afr_auth_center/AuthCenterClient"
	LendmeClient "afr_lendme/LendmeClient"
	Prop "afr_propylaea/PropylaeaClient"
	INClient "afr_sb_in"
	UCGW_client "afr_unified_charging_gateway/Unified_charging_gateway_Client"
	"daoc"
	"log"
	"time"
)

var CGWHostConfig UCGW_client.UC_GW_Client
var SpinAndWinHostConfig SpinAndWin_client.SpinAndWin_Client
var LendmeHostConfig LendmeClient.Lendme_Client

type UserControl struct {
	MongoDB        *daoc.MongoDB
	LoyaltyMongoDB *daoc.MongoDB
	CacheDir       *daoc.CacheRegistry
	AppAUC         *AuthCenterClient.AUC
	OKAPIAUC       *AuthCenterClient.AUC
	IN             *INClient.IN
	CGW            *UCGW_client.UC_GW
	Propylaea      *Prop.Propylaea
	SpinAndWin     *SpinAndWin_client.SpinAndWin
	APGW           *APGW.APGW
	Lendme         *LendmeClient.LendMe
}

func NewUserControl() *UserControl {
	MongoHostConfig := daoc.InitMongoHost(Configuration.MongoDB.ReplicaSet,
		Configuration.MongoDB.UserName,
		Configuration.MongoDB.Password,
		Configuration.MongoDB.HostIP_1,
		Configuration.MongoDB.HostPort_1,
		Configuration.MongoDB.HostIP_2,
		Configuration.MongoDB.HostPort_2,
		Configuration.MongoDB.HostIP_3,
		Configuration.MongoDB.HostPort_3,
		Configuration.MongoDB.HostIP_4,
		Configuration.MongoDB.HostPort_4,
	)
	log.Println(MongoHostConfig)
	LoyaltyMongoHostConfig := daoc.InitMongoHost(Configuration.LoyaltyMongoDB.ReplicaSet,
		Configuration.LoyaltyMongoDB.UserName,
		Configuration.LoyaltyMongoDB.Password,
		Configuration.LoyaltyMongoDB.HostIP_1,
		Configuration.LoyaltyMongoDB.HostPort_1,
		Configuration.LoyaltyMongoDB.HostIP_2,
		Configuration.LoyaltyMongoDB.HostPort_2,
		Configuration.LoyaltyMongoDB.HostIP_3,
		Configuration.LoyaltyMongoDB.HostPort_3,
		Configuration.LoyaltyMongoDB.HostIP_4,
		Configuration.LoyaltyMongoDB.HostPort_4,
	)
	log.Println(LoyaltyMongoHostConfig)
	App_AUCHostConfig := AuthCenterClient.InitHostConfig(Configuration.App_AUC.Protocol,
		Configuration.App_AUC.Hostname,
		Configuration.App_AUC.Port,
		Configuration.App_AUC.Module,
		Configuration.App_AUC.Version,
		"",
		Configuration.App_AUC.S2S_Username,
		Configuration.App_AUC.S2S_Password,
		Configuration.App_AUC.Timeout_After)
	log.Println(App_AUCHostConfig)
	OKAPI_AUCHostConfig := AuthCenterClient.InitHostConfig(Configuration.OKAPI_AUC.Protocol,
		Configuration.OKAPI_AUC.Hostname,
		Configuration.OKAPI_AUC.Port,
		Configuration.OKAPI_AUC.Module,
		Configuration.OKAPI_AUC.Version,
		"",
		Configuration.OKAPI_AUC.S2S_Username,
		Configuration.OKAPI_AUC.S2S_Password,
		Configuration.OKAPI_AUC.Timeout_After)
	log.Println(OKAPI_AUCHostConfig)
	INHostConfig := INClient.InitHostConfig(Configuration.IN.IP,
		Configuration.IN.Port,
		Configuration.IN.WS_SOAP_Endpoint,
		Configuration.IN.WS_XMLNS_SOAP_Env,
		Configuration.IN.WS_XMLNS_Web,
		Configuration.IN.WS_EVC_SOAP_Endpoint,
		Configuration.IN.WS_EVC_XMLNS_SOAP_Env,
		Configuration.IN.WS_EVC_XMLNS_Web,
		Configuration.IN.Default_OpId,
		Configuration.IN.Default_OpPwd,
		Configuration.IN.Is_OpPwd_Required,
		"",
		"",
		Configuration.IN.Timeout,
		Configuration.IN.PrintLogs)
	log.Println(INHostConfig)
	CGWAUC := AuthCenterClient.InitHostConfig(Configuration.CGW_AUC.Protocol,
		Configuration.CGW_AUC.Hostname,
		Configuration.CGW_AUC.Port,
		Configuration.CGW_AUC.Module,
		Configuration.CGW_AUC.Version,
		"",
		Configuration.CGW_AUC.S2S_Username,
		Configuration.CGW_AUC.S2S_Password,
		Configuration.CGW_AUC.Timeout_After)
	log.Println(CGWAUC)
	CGWHostConfig = UCGW_client.UC_GW_Client{
		Protocol:   Configuration.CGW.Protocol,
		Hostname:   Configuration.CGW.Hostname,
		Port:       Configuration.CGW.Port,
		Module:     Configuration.CGW.Module,
		Version:    Configuration.CGW.Version,
		Timeout:    10 * Configuration.CGW.Timeout,
		AUC_client: AuthCenterClient.NewAUCClient(CGWAUC).AUCClient,
	}
	log.Println(CGWHostConfig)
	propylaea_config := Prop.Propylaea_Client{
		Protocol:        Configuration.Propylaea.Protocol, // or https
		Hostname:        Configuration.Propylaea.Hostname,
		Port:            Configuration.Propylaea.Port,
		Module:          Configuration.Propylaea.Module,
		Version:         Configuration.Propylaea.Version,
		S2S_AccessToken: Configuration.Propylaea.S2S_AccessToken,
		Timeout:         10 * Configuration.Propylaea.Timeout_After, //timeout if no reply after X seconds
		AUC_client:      AuthCenterClient.NewAUCClient(OKAPI_AUCHostConfig).AUCClient,
	}
	log.Println(propylaea_config)
	SpinAndWinAUC := AuthCenterClient.InitHostConfig(Configuration.SpinAndWin_AUC.Protocol,
		Configuration.SpinAndWin_AUC.Hostname,
		Configuration.SpinAndWin_AUC.Port,
		Configuration.SpinAndWin_AUC.Module,
		Configuration.SpinAndWin_AUC.Version,
		"",
		Configuration.SpinAndWin_AUC.S2S_Username,
		Configuration.SpinAndWin_AUC.S2S_Password,
		Configuration.SpinAndWin_AUC.Timeout_After)
	log.Println(SpinAndWinAUC)
	SpinAndWinHostConfig = SpinAndWin_client.SpinAndWin_Client{
		Protocol:   Configuration.SpinAndWin.Protocol,
		Hostname:   Configuration.SpinAndWin.Hostname,
		Port:       Configuration.SpinAndWin.Port,
		Module:     Configuration.SpinAndWin.Module,
		Version:    Configuration.SpinAndWin.Version,
		Timeout:    10 * Configuration.SpinAndWin.Timeout,
		AUC_client: AuthCenterClient.NewAUCClient(SpinAndWinAUC).AUCClient,
	}
	log.Println(SpinAndWinHostConfig)
	APGW_config := APGW.APGW_Client{
		Protocol:        Configuration.APGW.Protocol, // or https
		Hostname:        Configuration.APGW.Hostname,
		Port:            Configuration.APGW.Port,
		S2S_Username:    Configuration.APGW.S2S_Username,
		S2S_Password:    Configuration.APGW.S2S_Password,
		S2S_AccessToken: "",
		Timeout:         10 * time.Second, //timeout if no reply after X seconds
	}
	log.Println(APGW_config)
	LendmeAUC := AuthCenterClient.InitHostConfig(
		Configuration.Lendme_AUC.Protocol,
		Configuration.Lendme_AUC.Hostname,
		Configuration.Lendme_AUC.Port,
		Configuration.Lendme_AUC.Module,
		Configuration.Lendme_AUC.Version,
		"",
		Configuration.Lendme_AUC.S2S_Username,
		Configuration.Lendme_AUC.S2S_Password,
		Configuration.Lendme_AUC.Timeout_After)
	log.Println(LendmeAUC)
	LendmeHostConfig = LendmeClient.Lendme_Client{
		Protocol:        Configuration.Lendme.Protocol,
		Hostname:        Configuration.Lendme.Hostname,
		LendmePort:      Configuration.Lendme.Port,
		LendMeModule:    Configuration.Lendme.Module,
		LendMeVersion:   Configuration.Lendme.Version,
		Timeout:         Configuration.Lendme.Timeout,
		AUC_client:      AuthCenterClient.NewAUCClient(LendmeAUC).AUCClient,
	}
	log.Println(LendmeHostConfig)

	UC := &UserControl{
		MongoDB:        daoc.NewMongoDBClient(MongoHostConfig),
		LoyaltyMongoDB: daoc.NewMongoDBClient(LoyaltyMongoHostConfig),
		AppAUC:         AuthCenterClient.NewAUCClient(App_AUCHostConfig),
		OKAPIAUC:       AuthCenterClient.NewAUCClient(OKAPI_AUCHostConfig),
		CacheDir:       daoc.NewCacheRegistry(),
		IN:             INClient.NewINClient(INHostConfig),
		CGW:            UCGW_client.NewUC_GWClient(CGWHostConfig),
		Propylaea:      Prop.NewPropylaeaClient(propylaea_config),
		SpinAndWin:     SpinAndWin_client.NewSpinAndWinClient(SpinAndWinHostConfig),
		APGW:           APGW.NewAPGWClient(APGW_config),
		Lendme:         LendmeClient.NewLendmeClient(LendmeHostConfig),
	}
	return UC
}