package Lendme

import (
	AuthCenterClient "afr_auth_center/AuthCenterClient"
	"daoc"
)

type UserControl struct {
	MongoDB  *daoc.MongoDB
	CacheDir *daoc.CacheRegistry
	AppAUC   *AuthCenterClient.AUC
	OKAPIAUC *AuthCenterClient.AUC
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

	UC := &UserControl{
		MongoDB:  daoc.NewMongoDBClient(MongoHostConfig),
		AppAUC:   AuthCenterClient.NewAUCClient(App_AUCHostConfig),
		OKAPIAUC: AuthCenterClient.NewAUCClient(OKAPI_AUCHostConfig),
		CacheDir: daoc.NewCacheRegistry(),
	}
	return UC
}
