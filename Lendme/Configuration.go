package Lendme

import (
	"time"
)

var Configuration ConfigType

type ConfigType struct {
	//HttpOKAPIServicePort string
	HttpAppServicePort  string
	OKAPIAllowedOrigins []string

	HostId  string
	DB_Name string
	Version string
	Module  string

	IsProduction bool

	//Lendme Config
	Min_Allowed_Amnt float64
	Service_FeePerc  float64

	Min_Allowed_Balance float64
	Max_Allowed_Balance float64

	App_AUC struct {
		Description   string
		Protocol      string
		Hostname      string
		Port          string
		Module        string
		Version       string
		S2S_Username  string
		S2S_Password  string
		Timeout_After time.Duration
	}

	// OKAPI_AUC struct {
	// 	Description   string
	// 	Protocol      string
	// 	Hostname      string
	// 	Port          string
	// 	Module        string
	// 	Version       string
	// 	S2S_Username  string
	// 	S2S_Password  string
	// 	Timeout_After time.Duration
	// }

	MongoDB struct {
		ReplicaSet string
		UserName   string
		Password   string
		HostIP_1   string
		HostPort_1 string
		HostIP_2   string
		HostPort_2 string
		HostIP_3   string
		HostPort_3 string
		HostIP_4   string
		HostPort_4 string
	}

	IN struct {
		IP                    string
		Port                  string
		WS_SOAP_Endpoint      string
		WS_XMLNS_SOAP_Env     string
		WS_XMLNS_Web          string
		WS_EVC_SOAP_Endpoint  string
		WS_EVC_XMLNS_SOAP_Env string
		WS_EVC_XMLNS_Web      string
		Default_OpId          string
		Default_OpPwd         string
		Is_OpPwd_Required     bool
		Timeout               time.Duration
		PrintLogs             bool
	}
}

func GetDefaultConfiguration() (err error) {
	Configuration = setDefaultConfiguration_DRC_Live()
	return nil
}

func setDefaultConfiguration_DRC_Live() (Configuration ConfigType) {
	//Configuration.HttpOKAPIServicePort = "9291"
	Configuration.HttpAppServicePort = "9290"

	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:3000")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:5173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4414")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://okpaihr.africell.ao")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://outlet.africell.ao")

	Configuration.HostId = "Lendme-01"
	Configuration.DB_Name = "Lendme_DB"

	Configuration.Version = "V1"
	Configuration.Module = "Lendme"

	Configuration.IsProduction = false
	Configuration.Min_Allowed_Amnt = 10
	Configuration.Service_FeePerc = 0.1
	Configuration.Min_Allowed_Balance = 0
	Configuration.Max_Allowed_Balance = 1000

	Configuration.App_AUC.Description = "App AUC service"
	Configuration.App_AUC.Protocol = "https"
	Configuration.App_AUC.Hostname = "Lendme_auc"
	Configuration.App_AUC.Port = "9001"
	Configuration.App_AUC.Module = "AUC"
	Configuration.App_AUC.Version = "V1"
	Configuration.App_AUC.S2S_Username = "Lendme_Admin"
	Configuration.App_AUC.S2S_Password = "s@l$e$IrSW0$4"
	Configuration.App_AUC.Timeout_After = 5 * time.Second

	// Configuration.OKAPI_AUC.Description = "OKAPI AUC service"
	// Configuration.OKAPI_AUC.Protocol = "https"
	// Configuration.OKAPI_AUC.Hostname = "auc"
	// Configuration.OKAPI_AUC.Port = "9001"
	// Configuration.OKAPI_AUC.Module = "AUC"
	// Configuration.OKAPI_AUC.Version = "V1"
	// Configuration.OKAPI_AUC.S2S_Username = "SalesMonitoring_OKAPI"
	// Configuration.OKAPI_AUC.S2S_Password = "s@le$P@s$W0$3"
	// Configuration.OKAPI_AUC.Timeout_After = 5 * time.Second

	//mongoDB
	Configuration.MongoDB.ReplicaSet = ""
	Configuration.MongoDB.UserName = "db_root"
	Configuration.MongoDB.Password = "B3202T@soSo0612w6"
	Configuration.MongoDB.HostIP_1 = "LendMe_mongodb"
	Configuration.MongoDB.HostPort_1 = "27017"
	Configuration.MongoDB.HostIP_2 = ""
	Configuration.MongoDB.HostPort_2 = ""
	Configuration.MongoDB.HostIP_3 = ""
	Configuration.MongoDB.HostPort_3 = ""
	Configuration.MongoDB.HostIP_4 = ""
	Configuration.MongoDB.HostPort_4 = ""

	Configuration.IN.IP = "10.70.1.38"
	Configuration.IN.Port = "8080"
	Configuration.IN.WS_SOAP_Endpoint = "/axis2/services/WebService.WebServiceHttpSoap12Endpoint/"
	Configuration.IN.WS_XMLNS_SOAP_Env = "http://schemas.xmlsoap.org/soap/envelope/"
	Configuration.IN.WS_XMLNS_Web = "http://webservice.CSI.omvia.convergys.com"

	Configuration.IN.WS_EVC_SOAP_Endpoint = "/axis2/services/ERechargeWebService.ERechargeWebServiceHttpSoap11Endpoint/"
	Configuration.IN.WS_EVC_XMLNS_SOAP_Env = "http://schemas.xmlsoap.org/soap/envelope/"
	Configuration.IN.WS_EVC_XMLNS_Web = "http://webservice.CSI.omvia.convergys.com"

	Configuration.IN.Default_OpId = "lendme"
	Configuration.IN.Default_OpPwd = ""
	Configuration.IN.Is_OpPwd_Required = false
	Configuration.IN.Timeout = 5
	Configuration.IN.PrintLogs = true

	//http://10.95.64.6:15403/?systemid=lendme&password=lendmeP@ssw0rd&Originator=setest&dest_addr=243900100606&msg_text=test&registered_delivery=0&ston=5&snpi=0&dton=1&dnpi=1&encoding=1

	return
}
