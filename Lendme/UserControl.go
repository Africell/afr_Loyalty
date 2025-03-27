package Lendme

import (
	AuthCenterClient "afr_auth_center/AuthCenterClient"
	INClient "afr_sb_in"
	"daoc"
)

type UserControl struct {
	MongoDB        *daoc.MongoDB
	LoyaltyMongoDB *daoc.MongoDB
	CacheDir       *daoc.CacheRegistry
	AppAUC         *AuthCenterClient.AUC
	OKAPIAUC       *AuthCenterClient.AUC
	IN             *INClient.IN
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
	App_AUCHostConfig := AuthCenterClient.InitHostConfig(Configuration.App_AUC.Protocol,
		Configuration.App_AUC.Hostname,
		Configuration.App_AUC.Port,
		Configuration.App_AUC.Module,
		Configuration.App_AUC.Version,
		"",
		Configuration.App_AUC.S2S_Username,
		Configuration.App_AUC.S2S_Password,
		Configuration.App_AUC.Timeout_After)
	OKAPI_AUCHostConfig := AuthCenterClient.InitHostConfig(Configuration.OKAPI_AUC.Protocol,
		Configuration.OKAPI_AUC.Hostname,
		Configuration.OKAPI_AUC.Port,
		Configuration.OKAPI_AUC.Module,
		Configuration.OKAPI_AUC.Version,
		"",
		Configuration.OKAPI_AUC.S2S_Username,
		Configuration.OKAPI_AUC.S2S_Password,
		Configuration.OKAPI_AUC.Timeout_After)
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

	UC := &UserControl{
		MongoDB:        daoc.NewMongoDBClient(MongoHostConfig),
		LoyaltyMongoDB: daoc.NewMongoDBClient(LoyaltyMongoHostConfig),
		AppAUC:         AuthCenterClient.NewAUCClient(App_AUCHostConfig),
		OKAPIAUC:       AuthCenterClient.NewAUCClient(OKAPI_AUCHostConfig),
		CacheDir:       daoc.NewCacheRegistry(),
		IN:             INClient.NewINClient(INHostConfig),
	}
	return UC
}
