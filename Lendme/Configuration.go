package Lendme

import (
	"time"
)

var Configuration ConfigType

type ConfigType struct {
	//HttpOKAPIServicePort string
	HttpAppServicePort  string
	OKAPIAllowedOrigins []string

	Operation string
	HostId    string
	DB_Name   string
	Version   string
	Module    string

	IsProduction bool

	//Lendme Config
	Min_Allowed_Amnt float64
	Service_FeePerc  float64

	Min_Allowed_AON        float64
	Min_Avg3MRecharge      float64
	Min_LastRechargePeriod float64
	Min_Allowed_Balance    float64
	Max_Allowed_Balance    float64

	ARPU_File_Path string

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

	SMPP struct {
		IP       string
		Port     string
		Login    string
		Password string
	}
}

func GetDefaultConfiguration() (err error) {
	//Configuration = setDefaultConfiguration_DRC_Live()
	Configuration = setDefaultConfiguration_GM_Live()
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

	Configuration.Operation = "DRC"
	Configuration.HostId = "Lendme-01"
	Configuration.DB_Name = "Lendme_DB"

	Configuration.Version = "V1"
	Configuration.Module = "Lendme"

	Configuration.IsProduction = false
	Configuration.Min_Allowed_Amnt = 10
	Configuration.Service_FeePerc = 0.1
	Configuration.Min_Allowed_AON = 3
	Configuration.Min_Avg3MRecharge = 50
	Configuration.Min_LastRechargePeriod = 60
	Configuration.Min_Allowed_Balance = 0
	Configuration.Max_Allowed_Balance = 1000
	Configuration.ARPU_File_Path = "/home/Subs_ARPU/"

	Configuration.App_AUC.Description = "App AUC service"
	Configuration.App_AUC.Protocol = "http"
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
	Configuration.IN.PrintLogs = false

	//http://10.95.64.6:15403/?systemid=lendme&password=lendmeP@ssw0rd&Originator=setest&dest_addr=243900100606&msg_text=test&registered_delivery=0&ston=5&snpi=0&dton=1&dnpi=1&encoding=1

	//SMPP
	Configuration.SMPP.IP = "10.95.64.6" //// floating IS IP: "10.250.8.53", test IP VPN: "10.250.0.52" (or .50, .51)
	Configuration.SMPP.Port = "15403"
	Configuration.SMPP.Login = "lendme"
	Configuration.SMPP.Password = "lendmeP@ssw0rd"

	return
}

func setDefaultConfiguration_GM_Live() (Configuration ConfigType) {
	//Configuration.HttpOKAPIServicePort = "9291"
	Configuration.HttpAppServicePort = "9290"

	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:3000")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:5173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4414")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://okpaihr.africell.ao")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://outlet.africell.ao")

	Configuration.Operation = "Gambia"
	Configuration.HostId = "Lendme-01"
	Configuration.DB_Name = "Lendme_DB"

	Configuration.Version = "V1"
	Configuration.Module = "Lendme"

	Configuration.IsProduction = false
	Configuration.Min_Allowed_Amnt = 5
	Configuration.Service_FeePerc = 0.1
	Configuration.Min_Allowed_AON = 3
	Configuration.Min_Avg3MRecharge = 5
	Configuration.Min_LastRechargePeriod = 60
	Configuration.Min_Allowed_Balance = 0
	Configuration.Max_Allowed_Balance = 677
	Configuration.ARPU_File_Path = "/home/Subs_ARPU/"

	Configuration.App_AUC.Description = "App AUC service"
	Configuration.App_AUC.Protocol = "http"
	Configuration.App_AUC.Hostname = "Lendme_auc"
	Configuration.App_AUC.Port = "9293"
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
	Configuration.MongoDB.ReplicaSet = "reps01"
	Configuration.MongoDB.UserName = "mongo-root"
	Configuration.MongoDB.Password = "Speci@LM0nG0P@ssw0rd_F0r_G@mB!A"
	Configuration.MongoDB.HostIP_1 = "10.64.33.49" //==>Primary
	Configuration.MongoDB.HostPort_1 = "9001"
	Configuration.MongoDB.HostIP_2 = "10.64.33.48" //==>Secondary
	Configuration.MongoDB.HostPort_2 = "9002"
	Configuration.MongoDB.HostIP_3 = "10.64.33.101" //==> Aribter
	Configuration.MongoDB.HostPort_3 = "9003"
	Configuration.MongoDB.HostIP_4 = ""
	Configuration.MongoDB.HostPort_4 = ""

	Configuration.IN.IP = "192.168.0.232"
	Configuration.IN.Port = "8080"
	Configuration.IN.WS_SOAP_Endpoint = "/axis2/services/WebService.WebServiceHttpSoap12Endpoint/"
	Configuration.IN.WS_XMLNS_SOAP_Env = "http://schemas.xmlsoap.org/soap/envelope/"
	Configuration.IN.WS_XMLNS_Web = "http://webservice.CSI.omvia.convergys.com"

	Configuration.IN.WS_EVC_SOAP_Endpoint = "/axis2/services/ERechargeWebService.ERechargeWebServiceHttpSoap11Endpoint/"
	Configuration.IN.WS_EVC_XMLNS_SOAP_Env = "http://schemas.xmlsoap.org/soap/envelope/"
	Configuration.IN.WS_EVC_XMLNS_Web = "http://webservice.CSI.omvia.convergys.com"

	Configuration.IN.Default_OpId = "lendme"
	Configuration.IN.Default_OpPwd = "lenmeP@sw0rd"
	Configuration.IN.Is_OpPwd_Required = true
	Configuration.IN.Timeout = 5
	Configuration.IN.PrintLogs = true

	//http://10.95.64.6:15403/?systemid=lendme&password=lendmeP@ssw0rd&Originator=setest&dest_addr=243900100606&msg_text=test&registered_delivery=0&ston=5&snpi=0&dton=1&dnpi=1&encoding=1

	//SMPP
	Configuration.SMPP.IP = "10.30.8.10"
	Configuration.SMPP.Port = "15403"
	Configuration.SMPP.Login = "lendme"
	Configuration.SMPP.Password = "lendmeP@ssw0rd"

	return
}
