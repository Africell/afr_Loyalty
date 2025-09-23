package Lendme

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"time"
)

var EncryptionKey string = "A4ask%a2l&S&rRo1~~2Fo|003j|`XX%&*h1u)U(*@(NOD!3!PH`OPD7#HAHA))Y:)"

var Configuration ConfigType

type ConfigType struct {
	//HttpOKAPIServicePort string
	HttpAppServicePort           string
	HttpAppLoyaltyServicePort    string
	HttpAppLoyaltyManagementPort string
	HttpAppLoyaltyFeedPort       string
	OKAPIAllowedOrigins          []string

	Operation       string
	HostId          string
	DB_Name         string
	DB_Name_Loyalty string
	Version         string
	Module          string

	IsProduction                   bool
	IsLoyaltyProduction            bool
	ISLoyaltyOptIn                 bool // if true customer has to opt in
	ISLoyaltyOptOutGracePeriodDays int  // if true customer has to opt in
	LoyaltyProgramName             string
	LoyaltyVersion                 string
	LoyaltyModule                  string

	MSISDN_Prefix    string
	MSISDN_Short_len int
	CountryCode      string

	//Lendme Config
	Min_Allowed_Amnt float64
	Service_FeePerc  float64

	Min_Allowed_AON        float64
	Min_Avg3MRecharge      float64
	Min_LastRechargePeriod float64
	Min_Allowed_Balance    float64
	Max_Allowed_Balance    float64

	ARPU_File_Path             string
	ARPU_File_Prefix           string
	ARPU_File_Column_Separator string

	Lendme_EVC_Dealer_MSISDN string
	Lendme_EVC_Dealer_PIN    string

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

	OKAPI_AUC struct {
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
	LoyaltyMongoDB struct {
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
		IP                string
		Port              string
		Login             string
		Password          string
		TimeOut           time.Duration
		PrintLogs         bool
		MSISDN_Short_len  int //length without country code and without 0 prefix
		CountryCodePrefix string
		DefaultSender     string
		Encoding          int
	}

	CGW_AUC struct {
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

	CGW struct {
		Protocol        string //http or https
		Hostname        string //name or IP
		Port            string
		Module          string
		Version         string
		S2S_AccessToken string        //system to system access token
		Timeout         time.Duration //timeout if no reply after X seconds
	}

	Propylaea struct {
		Description     string
		Protocol        string
		Hostname        string
		Port            string
		Module          string
		Version         string
		S2S_Username    string
		S2S_Password    string
		S2S_AccessToken string
		Timeout_After   time.Duration
		ChannelName     string
		ChannelPlan     string
		ChannelVersion  string
	}

	SpinAndWin_AUC struct {
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

	SpinAndWin struct {
		Protocol        string //http or https
		Hostname        string //name or IP
		Port            string
		Module          string
		Version         string
		S2S_AccessToken string        //system to system access token
		Timeout         time.Duration //timeout if no reply after X seconds
	}
	APGW struct {
		Description   string
		Protocol      string
		Hostname      string
		Port          string
		S2S_Username  string
		S2S_Password  string
		Timeout_After time.Duration
	}
}

func GetDefaultConfiguration() (err error) {
	// Configuration = setDefaultConfiguration_Dev()
	// Configuration = setDefaultConfiguration_DRC_Live()
	// Configuration = setDefaultConfiguration_DRC_Loyalty_UAT()
	//	Configuration = setDefaultConfiguration_GM_Live()

	// Configuration = setDefaultConfiguration_SL_Live()
	Configuration = setDefaultConfiguration_SL_Loyalty()
	// Configuration = setDefaultConfiguration_SL_Loyalty_UAT()

	// Configuration = setDefaultConfiguration_AO_Loyalty()
	//Configuration = setDefaultConfiguration_AO_Loyalty_UAT()

	// Configuration = setDefaultConfiguration_GM_Loyalty_Live()
	// Configuration = setDefaultConfiguration_GM_Loyalty_UAT()
	return nil
}

func setDefaultConfiguration_DRC_Live() (Configuration ConfigType) {
	//Configuration.HttpOKAPIServicePort = "9291"
	Configuration.HttpAppServicePort = "9290"           //lendme services
	Configuration.HttpAppLoyaltyServicePort = "9280"    //for USSD and Mobile App
	Configuration.HttpAppLoyaltyManagementPort = "9281" //for OKAPI
	Configuration.HttpAppLoyaltyFeedPort = "9282"       //for IN & MM live feed

	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:3000")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:5173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4414")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://okpaihr.africell.cd")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://outlet.africell.cd")

	Configuration.Operation = "DRC"
	Configuration.HostId = "Lendme-01"
	Configuration.DB_Name = "Lendme_DB"

	Configuration.Version = "V1"
	Configuration.Module = "Lendme"

	Configuration.LoyaltyVersion = "V1"
	Configuration.LoyaltyModule = "Loyalty"

	Configuration.MSISDN_Prefix = ""
	Configuration.MSISDN_Short_len = 9
	Configuration.CountryCode = "243"

	Configuration.IsProduction = true
	Configuration.IsLoyaltyProduction = false
	Configuration.Min_Allowed_Amnt = 10
	Configuration.Service_FeePerc = 0.1
	Configuration.Min_Allowed_AON = 3
	Configuration.Min_Avg3MRecharge = 50
	Configuration.Min_LastRechargePeriod = 60
	Configuration.Min_Allowed_Balance = 0
	Configuration.Max_Allowed_Balance = 1000
	Configuration.ARPU_File_Path = "/home/Subs_ARPU/"
	Configuration.ARPU_File_Prefix = "Rgs_"
	Configuration.ARPU_File_Column_Separator = ","

	Configuration.App_AUC.Description = "App AUC service"
	Configuration.App_AUC.Protocol = "http"
	Configuration.App_AUC.Hostname = "Lendme_auc"
	Configuration.App_AUC.Port = "9001"
	Configuration.App_AUC.Module = "AUC"
	Configuration.App_AUC.Version = "V1"
	Configuration.App_AUC.S2S_Username = "Lendme_Admin"
	Configuration.App_AUC.S2S_Password = "s@l$e$IrSW0$4"
	Configuration.App_AUC.Timeout_After = 5 * time.Second

	Configuration.OKAPI_AUC.Description = "OKAPI AUC service"
	Configuration.OKAPI_AUC.Protocol = "http"
	Configuration.OKAPI_AUC.Hostname = "10.95.72.166"
	Configuration.OKAPI_AUC.Port = "9001"
	Configuration.OKAPI_AUC.Module = "AUC"
	Configuration.OKAPI_AUC.Version = "V1"
	Configuration.OKAPI_AUC.S2S_Username = "SalesMonitoring_OKAPI"
	Configuration.OKAPI_AUC.S2S_Password = "s@le$P@s$W0$3"
	Configuration.OKAPI_AUC.Timeout_After = 5 * time.Second

	//mongoDB
	// Configuration.MongoDB.ReplicaSet = ""
	// Configuration.MongoDB.UserName = "db_root"
	// Configuration.MongoDB.Password = "B3202T@soSo0612w6"
	// Configuration.MongoDB.HostIP_1 = "LendMe_mongodb"
	// Configuration.MongoDB.HostPort_1 = "27017"
	// Configuration.MongoDB.HostIP_2 = ""
	// Configuration.MongoDB.HostPort_2 = ""
	// Configuration.MongoDB.HostIP_3 = ""
	// Configuration.MongoDB.HostPort_3 = ""
	// Configuration.MongoDB.HostIP_4 = ""
	// Configuration.MongoDB.HostPort_4 = ""
	Configuration.MongoDB.ReplicaSet = "reps1"
	Configuration.MongoDB.UserName = "db_root"
	Configuration.MongoDB.Password = "B3202T@soSo0612w6"
	Configuration.MongoDB.HostIP_1 = "10.95.72.177" //==> Primary
	Configuration.MongoDB.HostPort_1 = "9001"
	Configuration.MongoDB.HostIP_2 = "10.95.72.176" //==>secondary
	Configuration.MongoDB.HostPort_2 = "9002"
	Configuration.MongoDB.HostIP_3 = "10.95.32.22" //==> Aribter
	Configuration.MongoDB.HostPort_3 = "9003"
	Configuration.MongoDB.HostIP_4 = ""
	Configuration.MongoDB.HostPort_4 = ""

	Configuration.DB_Name_Loyalty = "Loyalty_DB"
	Configuration.LoyaltyProgramName = "Loyalty"
	Configuration.LoyaltyMongoDB.ReplicaSet = Configuration.MongoDB.ReplicaSet
	Configuration.LoyaltyMongoDB.UserName = Configuration.MongoDB.UserName
	Configuration.LoyaltyMongoDB.Password = Configuration.MongoDB.Password
	Configuration.LoyaltyMongoDB.HostIP_1 = Configuration.MongoDB.HostIP_1
	Configuration.LoyaltyMongoDB.HostPort_1 = Configuration.MongoDB.HostPort_1
	Configuration.LoyaltyMongoDB.HostIP_2 = Configuration.MongoDB.HostIP_2
	Configuration.LoyaltyMongoDB.HostPort_2 = Configuration.MongoDB.HostPort_2
	Configuration.LoyaltyMongoDB.HostIP_3 = Configuration.MongoDB.HostIP_3
	Configuration.LoyaltyMongoDB.HostPort_3 = Configuration.MongoDB.HostPort_3
	Configuration.LoyaltyMongoDB.HostIP_4 = Configuration.MongoDB.HostIP_4
	Configuration.LoyaltyMongoDB.HostPort_4 = Configuration.MongoDB.HostPort_4

	Configuration.IN.IP = "10.95.73.12" //"10.70.1.59"
	Configuration.IN.Port = "8444"

	Configuration.IN.WS_SOAP_Endpoint = "/axis2/services/WebService.WebServiceHttpSoap12Endpoint/"
	Configuration.IN.WS_XMLNS_SOAP_Env = "http://schemas.xmlsoap.org/soap/envelope/"
	Configuration.IN.WS_XMLNS_Web = "http://webservice.CSI.omvia.convergys.com"

	Configuration.IN.WS_EVC_SOAP_Endpoint = "/axis2/services/ERechargeWebService.ERechargeWebServiceHttpSoap11Endpoint/"
	Configuration.IN.WS_EVC_XMLNS_SOAP_Env = "http://schemas.xmlsoap.org/soap/envelope/"
	Configuration.IN.WS_EVC_XMLNS_Web = "http://webservice.CSI.omvia.convergys.com"

	Configuration.IN.Default_OpId = "lendme"
	Configuration.IN.Default_OpPwd = "P@ssw0rd"
	Configuration.IN.Is_OpPwd_Required = true
	Configuration.IN.Timeout = 5
	Configuration.IN.PrintLogs = false

	//http://10.95.64.6:15403/?systemid=lendme&password=lendmeP@ssw0rd&Originator=setest&dest_addr=243900100606&msg_text=test&registered_delivery=0&ston=5&snpi=0&dton=1&dnpi=1&encoding=1

	//SMPP
	Configuration.SMPP.IP = "10.95.64.6" //// floating IS IP: "10.250.8.53", test IP VPN: "10.250.0.52" (or .50, .51)
	Configuration.SMPP.Port = "15403"
	Configuration.SMPP.Login = "lendme"
	Configuration.SMPP.Password = "lendmeP@ssw0rd"

	Configuration.CGW_AUC.Description = "UCGW AUC service"
	Configuration.CGW_AUC.Protocol = "http"
	Configuration.CGW_AUC.Hostname = "10.95.72.95"
	Configuration.CGW_AUC.Port = "9994"
	Configuration.CGW_AUC.Module = "AUC"
	Configuration.CGW_AUC.Version = "V1"
	Configuration.CGW_AUC.S2S_Username = "SAW_UCGW"
	Configuration.CGW_AUC.S2S_Password = "uC@g$ASDKJH66&&&RiS6$2"
	Configuration.CGW_AUC.Timeout_After = 5 * time.Second

	Configuration.CGW.Protocol = "http"
	Configuration.CGW.Hostname = "10.95.72.95"
	Configuration.CGW.Port = "9991"
	Configuration.CGW.Module = "UCGW"
	Configuration.CGW.Version = "V1"
	Configuration.CGW.Timeout = 15 * time.Second

	Configuration.Propylaea.Description = "Product Design Center - Propylaea"
	Configuration.Propylaea.Protocol = "http"
	Configuration.Propylaea.Port = "9900"
	Configuration.Propylaea.Hostname = "10.95.72.95"
	Configuration.Propylaea.Module = "Propylaea"
	Configuration.Propylaea.Version = "V1"
	Configuration.Propylaea.S2S_Username = "Propylaea_Admin"
	Configuration.Propylaea.S2S_Password = "uC@g$W$iRiS6$2@333dd"
	Configuration.Propylaea.Timeout_After = 5 * time.Second
	Configuration.Propylaea.ChannelName = "Spin And Win"
	Configuration.Propylaea.ChannelPlan = "Normal SIM"
	Configuration.Propylaea.ChannelVersion = "1"

	return
}

func setDefaultConfiguration_DRC_Loyalty_UAT() (Configuration ConfigType) {
	//Configuration.HttpOKAPIServicePort = "9291"
	Configuration.HttpAppServicePort = "9290"           //lendme services
	Configuration.HttpAppLoyaltyServicePort = "9280"    //for USSD and Mobile App
	Configuration.HttpAppLoyaltyManagementPort = "9281" //for OKAPI
	Configuration.HttpAppLoyaltyFeedPort = "9282"       //for IN & MM live feed

	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:3000")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:5173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4414")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://okapihruat.africell.cd")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://outlet.africell.cd")

	Configuration.Operation = "DRC"
	Configuration.HostId = "Lendme-01"
	Configuration.DB_Name = "Lendme_DB"

	Configuration.Version = "V1"
	Configuration.Module = "Lendme"

	Configuration.LoyaltyVersion = "V1"
	Configuration.LoyaltyModule = "Loyalty"

	Configuration.MSISDN_Prefix = ""
	Configuration.MSISDN_Short_len = 9
	Configuration.CountryCode = "243"

	Configuration.IsProduction = true
	Configuration.IsLoyaltyProduction = false
	Configuration.Min_Allowed_Amnt = 10
	Configuration.Service_FeePerc = 0.1
	Configuration.Min_Allowed_AON = 3
	Configuration.Min_Avg3MRecharge = 50
	Configuration.Min_LastRechargePeriod = 60
	Configuration.Min_Allowed_Balance = 0
	Configuration.Max_Allowed_Balance = 1000
	Configuration.ARPU_File_Path = "/home/Subs_ARPU/"
	Configuration.ARPU_File_Prefix = "Rgs_"
	Configuration.ARPU_File_Column_Separator = ","

	Configuration.App_AUC.Description = "App AUC service"
	Configuration.App_AUC.Protocol = "http"
	Configuration.App_AUC.Hostname = "Lendme_auc"
	Configuration.App_AUC.Port = "9293"
	Configuration.App_AUC.Module = "AUC"
	Configuration.App_AUC.Version = "V1"
	Configuration.App_AUC.S2S_Username = "Lendme_Admin"
	Configuration.App_AUC.S2S_Password = "s@l$e$IrSW0$4"
	Configuration.App_AUC.Timeout_After = 5 * time.Second

	Configuration.OKAPI_AUC.Description = "OKAPI AUC service"
	Configuration.OKAPI_AUC.Protocol = "http"
	Configuration.OKAPI_AUC.Hostname = "okapi_auc"
	Configuration.OKAPI_AUC.Port = "9001"
	Configuration.OKAPI_AUC.Module = "AUC"
	Configuration.OKAPI_AUC.Version = "V1"
	Configuration.OKAPI_AUC.S2S_Username = "Loyalty_OKAPI"
	Configuration.OKAPI_AUC.S2S_Password = "]W8#x3D1USKUyH@p]s&D_"
	Configuration.OKAPI_AUC.Timeout_After = 5 * time.Second

	//mongoDB
	// Configuration.MongoDB.ReplicaSet = ""
	// Configuration.MongoDB.UserName = "db_root"
	// Configuration.MongoDB.Password = "B3202T@soSo0612w6"
	// Configuration.MongoDB.HostIP_1 = "LendMe_mongodb"
	// Configuration.MongoDB.HostPort_1 = "27017"
	// Configuration.MongoDB.HostIP_2 = ""
	// Configuration.MongoDB.HostPort_2 = ""
	// Configuration.MongoDB.HostIP_3 = ""
	// Configuration.MongoDB.HostPort_3 = ""
	// Configuration.MongoDB.HostIP_4 = ""
	// Configuration.MongoDB.HostPort_4 = ""
	Configuration.MongoDB.UserName = "db_root"
	Configuration.MongoDB.Password = "P@s54D0Brdara_r@75S"
	Configuration.MongoDB.HostIP_1 = "mongodb" //"host.docker.internal"
	Configuration.MongoDB.HostPort_1 = "27017" //"27017"

	Configuration.DB_Name_Loyalty = "Loyalty_DB"
	Configuration.LoyaltyMongoDB.ReplicaSet = Configuration.MongoDB.ReplicaSet
	Configuration.LoyaltyMongoDB.UserName = Configuration.MongoDB.UserName
	Configuration.LoyaltyMongoDB.Password = Configuration.MongoDB.Password
	Configuration.LoyaltyMongoDB.HostIP_1 = Configuration.MongoDB.HostIP_1
	Configuration.LoyaltyMongoDB.HostPort_1 = Configuration.MongoDB.HostPort_1
	Configuration.LoyaltyMongoDB.HostIP_2 = Configuration.MongoDB.HostIP_2
	Configuration.LoyaltyMongoDB.HostPort_2 = Configuration.MongoDB.HostPort_2
	Configuration.LoyaltyMongoDB.HostIP_3 = Configuration.MongoDB.HostIP_3
	Configuration.LoyaltyMongoDB.HostPort_3 = Configuration.MongoDB.HostPort_3
	Configuration.LoyaltyMongoDB.HostIP_4 = Configuration.MongoDB.HostIP_4
	Configuration.LoyaltyMongoDB.HostPort_4 = Configuration.MongoDB.HostPort_4

	Configuration.IN.IP = "10.95.73.12" //"10.70.1.59"
	Configuration.IN.Port = "8444"

	Configuration.IN.WS_SOAP_Endpoint = "/axis2/services/WebService.WebServiceHttpSoap12Endpoint/"
	Configuration.IN.WS_XMLNS_SOAP_Env = "http://schemas.xmlsoap.org/soap/envelope/"
	Configuration.IN.WS_XMLNS_Web = "http://webservice.CSI.omvia.convergys.com"

	Configuration.IN.WS_EVC_SOAP_Endpoint = "/axis2/services/ERechargeWebService.ERechargeWebServiceHttpSoap11Endpoint/"
	Configuration.IN.WS_EVC_XMLNS_SOAP_Env = "http://schemas.xmlsoap.org/soap/envelope/"
	Configuration.IN.WS_EVC_XMLNS_Web = "http://webservice.CSI.omvia.convergys.com"

	Configuration.IN.Default_OpId = "lendme"
	Configuration.IN.Default_OpPwd = "P@ssw0rd"
	Configuration.IN.Is_OpPwd_Required = true
	Configuration.IN.Timeout = 5
	Configuration.IN.PrintLogs = true

	//http://10.95.64.6:15403/?systemid=lendme&password=lendmeP@ssw0rd&Originator=setest&dest_addr=243900100606&msg_text=test&registered_delivery=0&ston=5&snpi=0&dton=1&dnpi=1&encoding=1

	//SMPP
	Configuration.SMPP.IP = "10.95.64.6" //// floating IS IP: "10.250.8.53", test IP VPN: "10.250.0.52" (or .50, .51)
	Configuration.SMPP.Port = "15403"
	//Configuration.SMPP.Login = "lendme"
	//Configuration.SMPP.Password = "lendmeP@ssw0rd"
	Configuration.SMPP.Login = "Loyalty"
	Configuration.SMPP.Password = "Loyalty123"
	Configuration.SMPP.TimeOut = 5 //in seconds
	Configuration.SMPP.PrintLogs = true
	Configuration.SMPP.MSISDN_Short_len = 8
	Configuration.SMPP.CountryCodePrefix = "232"
	Configuration.SMPP.DefaultSender = "Africell" //"Africell"
	Configuration.SMPP.Encoding = 0
	//CGW
	Configuration.CGW.Protocol = "http"
	Configuration.CGW.Hostname = "UCGW"
	Configuration.CGW.Port = "9991"
	Configuration.CGW.Module = "UCGW"
	Configuration.CGW.Version = "V1"
	Configuration.CGW.Timeout = 15 * time.Second
	Configuration.CGW.S2S_AccessToken = ""

	Configuration.CGW_AUC.Description = "UCGW AUC service"
	Configuration.CGW_AUC.Protocol = "http"
	Configuration.CGW_AUC.Hostname = "UCGW_AUC"
	Configuration.CGW_AUC.Port = "9001"
	Configuration.CGW_AUC.Module = "AUC"
	Configuration.CGW_AUC.Version = "V1"
	Configuration.CGW_AUC.S2S_Username = "UCGW_Admin"
	Configuration.CGW_AUC.S2S_Password = "uC@g$W$iRiS6$2"
	Configuration.CGW_AUC.Timeout_After = 5 * time.Second

	Configuration.Propylaea.Description = "Product Design Center - Propylaea"
	Configuration.Propylaea.Protocol = "http"
	Configuration.Propylaea.Port = "9900"
	Configuration.Propylaea.Hostname = "Propylaea_Prod"
	Configuration.Propylaea.Module = "Propylaea"
	Configuration.Propylaea.Version = "V1"
	Configuration.Propylaea.S2S_Username = "Propylaea_Admin"
	Configuration.Propylaea.S2S_Password = "uC@g$W$iRiS6$2@333dd"
	Configuration.Propylaea.Timeout_After = 5 * time.Second

	// Configuration.SpinAndWin_AUC.Description = "SAW AUC"
	// Configuration.SpinAndWin_AUC.Protocol = "http"
	// Configuration.SpinAndWin_AUC.Hostname = "10.95.72.166"
	// Configuration.SpinAndWin_AUC.Port = "9001"
	// Configuration.SpinAndWin_AUC.Module = "AUC"
	// Configuration.SpinAndWin_AUC.Version = "V1"
	// Configuration.SpinAndWin_AUC.S2S_Username = "SAW_Admin"
	// Configuration.SpinAndWin_AUC.S2S_Password = "LQaDUp388UNKhz0Ap"
	// Configuration.SpinAndWin_AUC.Timeout_After = 30 * time.Second
	Configuration.SpinAndWin.Protocol = "http"
	Configuration.SpinAndWin.Hostname = "spinwin_be"
	Configuration.SpinAndWin.Port = "9111"
	Configuration.SpinAndWin.Module = "SpinAndWin"
	Configuration.SpinAndWin.Version = "V1"
	Configuration.SpinAndWin.Timeout = 30 * time.Second

	return
}

func setDefaultConfiguration_GM_Live() (Configuration ConfigType) {
	//Configuration.HttpOKAPIServicePort = "9291"
	Configuration.HttpAppServicePort = "9290"           //lendme services
	Configuration.HttpAppLoyaltyServicePort = "9280"    //for USSD and Mobile App
	Configuration.HttpAppLoyaltyManagementPort = "9281" //for OKAPI
	Configuration.HttpAppLoyaltyFeedPort = "9282"       //for IN & MM live feed

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

	Configuration.LoyaltyVersion = "V1"
	Configuration.LoyaltyModule = "Loyalty"

	Configuration.IsProduction = false
	Configuration.IsLoyaltyProduction = false
	Configuration.Min_Allowed_Amnt = 5
	Configuration.Service_FeePerc = 0.04
	Configuration.Min_Allowed_AON = 3
	Configuration.Min_Avg3MRecharge = 5
	Configuration.Min_LastRechargePeriod = 60
	Configuration.Min_Allowed_Balance = 0
	Configuration.Max_Allowed_Balance = 677
	Configuration.ARPU_File_Path = "/home/Subs_ARPU/"
	Configuration.ARPU_File_Prefix = "Rgs_"
	Configuration.ARPU_File_Column_Separator = ","

	Configuration.App_AUC.Description = "App AUC service"
	Configuration.App_AUC.Protocol = "http"
	Configuration.App_AUC.Hostname = "Lendme_auc"
	Configuration.App_AUC.Port = "9293"
	Configuration.App_AUC.Module = "AUC"
	Configuration.App_AUC.Version = "V1"
	Configuration.App_AUC.S2S_Username = "Lendme_Admin"
	Configuration.App_AUC.S2S_Password = "s@l$e$IrSW0$4"
	Configuration.App_AUC.Timeout_After = 5 * time.Second

	Configuration.OKAPI_AUC.Description = "OKAPI AUC service"
	Configuration.OKAPI_AUC.Protocol = "https"
	Configuration.OKAPI_AUC.Hostname = "auc"
	Configuration.OKAPI_AUC.Port = "9001"
	Configuration.OKAPI_AUC.Module = "AUC"
	Configuration.OKAPI_AUC.Version = "V1"
	Configuration.OKAPI_AUC.S2S_Username = "SalesMonitoring_OKAPI"
	Configuration.OKAPI_AUC.S2S_Password = "s@le$P@s$W0$3"
	Configuration.OKAPI_AUC.Timeout_After = 5 * time.Second

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

	Configuration.DB_Name_Loyalty = "Loyalty_DB"
	Configuration.LoyaltyMongoDB.ReplicaSet = Configuration.MongoDB.ReplicaSet
	Configuration.LoyaltyMongoDB.UserName = Configuration.MongoDB.UserName
	Configuration.LoyaltyMongoDB.Password = Configuration.MongoDB.Password
	Configuration.LoyaltyMongoDB.HostIP_1 = Configuration.MongoDB.HostIP_1
	Configuration.LoyaltyMongoDB.HostPort_1 = Configuration.MongoDB.HostPort_1
	Configuration.LoyaltyMongoDB.HostIP_2 = Configuration.MongoDB.HostIP_2
	Configuration.LoyaltyMongoDB.HostPort_2 = Configuration.MongoDB.HostPort_2
	Configuration.LoyaltyMongoDB.HostIP_3 = Configuration.MongoDB.HostIP_3
	Configuration.LoyaltyMongoDB.HostPort_3 = Configuration.MongoDB.HostPort_3
	Configuration.LoyaltyMongoDB.HostIP_4 = Configuration.MongoDB.HostIP_4
	Configuration.LoyaltyMongoDB.HostPort_4 = Configuration.MongoDB.HostPort_4

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
	Configuration.SMPP.Login = "LendME2"
	Configuration.SMPP.Password = "LendMEP@ssw0rd"

	Configuration.CGW_AUC.Description = "UCGW AUC service"
	Configuration.CGW_AUC.Protocol = "http"
	Configuration.CGW_AUC.Hostname = "?.?.?.?"
	Configuration.CGW_AUC.Port = "9994"
	Configuration.CGW_AUC.Module = "AUC"
	Configuration.CGW_AUC.Version = "V1"
	Configuration.CGW_AUC.S2S_Username = "SAW_UCGW"
	Configuration.CGW_AUC.S2S_Password = "uC@g$ASDKJH66&&&RiS6$2"
	Configuration.CGW_AUC.Timeout_After = 5 * time.Second

	Configuration.CGW.Protocol = "http"
	Configuration.CGW.Hostname = "?.?.?.?"
	Configuration.CGW.Port = "9991"
	Configuration.CGW.Module = "UCGW"
	Configuration.CGW.Version = "V1"
	Configuration.CGW.Timeout = 15 * time.Second

	Configuration.Propylaea.Description = "Product Design Center - Propylaea"
	Configuration.Propylaea.Protocol = "https"
	Configuration.Propylaea.Port = "443"
	Configuration.Propylaea.Hostname = "sapp.africell.cd"
	Configuration.Propylaea.Module = "Propylaea"
	Configuration.Propylaea.Version = "V1"
	Configuration.Propylaea.S2S_Username = "Propylaea_Admin"
	Configuration.Propylaea.S2S_Password = "uC@g$W$iRiS6$2@333dd"
	Configuration.Propylaea.Timeout_After = 5 * time.Second
	Configuration.Propylaea.ChannelName = "Spin And Win"
	Configuration.Propylaea.ChannelPlan = "Normal SIM"
	Configuration.Propylaea.ChannelVersion = "1"

	return
}

func setDefaultConfiguration_GM_Loyalty() (Configuration ConfigType) {
	//Configuration.HttpOKAPIServicePort = "9291"
	Configuration.HttpAppServicePort = "9290"           //lendme services
	Configuration.HttpAppLoyaltyServicePort = "9280"    //for USSD and Mobile App
	Configuration.HttpAppLoyaltyManagementPort = "9281" //for OKAPI
	Configuration.HttpAppLoyaltyFeedPort = "9282"       //for IN & MM live feed

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

	Configuration.LoyaltyVersion = "V1"
	Configuration.LoyaltyModule = "Loyalty"

	Configuration.MSISDN_Prefix = ""
	Configuration.MSISDN_Short_len = 7
	Configuration.CountryCode = "220"

	Configuration.IsProduction = false
	Configuration.IsLoyaltyProduction = false
	Configuration.Min_Allowed_Amnt = 5
	Configuration.Service_FeePerc = 0.04
	Configuration.Min_Allowed_AON = 3
	Configuration.Min_Avg3MRecharge = 5
	Configuration.Min_LastRechargePeriod = 60
	Configuration.Min_Allowed_Balance = 0
	Configuration.Max_Allowed_Balance = 677
	Configuration.ARPU_File_Path = "/home/Subs_ARPU/"
	Configuration.ARPU_File_Prefix = "Rgs_"
	Configuration.ARPU_File_Column_Separator = ","

	Configuration.App_AUC.Description = "App AUC service"
	Configuration.App_AUC.Protocol = "http"
	Configuration.App_AUC.Hostname = "Lendme_auc"
	Configuration.App_AUC.Port = "9293"
	Configuration.App_AUC.Module = "AUC"
	Configuration.App_AUC.Version = "V1"
	Configuration.App_AUC.S2S_Username = "Lendme_Admin"
	Configuration.App_AUC.S2S_Password = "s@l$e$IrSW0$4"
	Configuration.App_AUC.Timeout_After = 5 * time.Second

	Configuration.OKAPI_AUC.Description = "OKAPI AUC service"
	Configuration.OKAPI_AUC.Protocol = "http"
	//Configuration.OKAPI_AUC.Hostname = "auc"
	Configuration.OKAPI_AUC.Hostname = "10.30.0.120"
	Configuration.OKAPI_AUC.Port = "9001"
	Configuration.OKAPI_AUC.Module = "AUC"
	Configuration.OKAPI_AUC.Version = "V1"
	Configuration.OKAPI_AUC.S2S_Username = "SalesMonitoring_OKAPI"
	Configuration.OKAPI_AUC.S2S_Password = "s@le$P@s$W0$3"
	Configuration.OKAPI_AUC.Timeout_After = 5 * time.Second

	// //mongoDB
	// Configuration.MongoDB.ReplicaSet = "reps01"
	// Configuration.MongoDB.UserName = "mongo-root"
	// Configuration.MongoDB.Password = "Speci@LM0nG0P@ssw0rd_F0r_G@mB!A"
	// Configuration.MongoDB.HostIP_1 = "10.64.33.49" //==>Primary
	// Configuration.MongoDB.HostPort_1 = "9001"
	// Configuration.MongoDB.HostIP_2 = "10.64.33.48" //==>Secondary
	// Configuration.MongoDB.HostPort_2 = "9002"
	// Configuration.MongoDB.HostIP_3 = "10.64.33.101" //==> Aribter
	// Configuration.MongoDB.HostPort_3 = "9003"
	// Configuration.MongoDB.HostIP_4 = ""
	// Configuration.MongoDB.HostPort_4 = ""

	Configuration.MongoDB.ReplicaSet = ""
	Configuration.MongoDB.UserName = "db_root"
	Configuration.MongoDB.Password = "P@s54D0Brdara_r@75S"
	Configuration.MongoDB.HostIP_1 = "10.30.0.151" //==>Primary
	Configuration.MongoDB.HostPort_1 = "9510"

	Configuration.DB_Name_Loyalty = "Loyalty_DB"
	Configuration.LoyaltyMongoDB.ReplicaSet = Configuration.MongoDB.ReplicaSet
	Configuration.LoyaltyMongoDB.UserName = Configuration.MongoDB.UserName
	Configuration.LoyaltyMongoDB.Password = Configuration.MongoDB.Password
	Configuration.LoyaltyMongoDB.HostIP_1 = Configuration.MongoDB.HostIP_1
	Configuration.LoyaltyMongoDB.HostPort_1 = Configuration.MongoDB.HostPort_1
	Configuration.LoyaltyMongoDB.HostIP_2 = Configuration.MongoDB.HostIP_2
	Configuration.LoyaltyMongoDB.HostPort_2 = Configuration.MongoDB.HostPort_2
	Configuration.LoyaltyMongoDB.HostIP_3 = Configuration.MongoDB.HostIP_3
	Configuration.LoyaltyMongoDB.HostPort_3 = Configuration.MongoDB.HostPort_3
	Configuration.LoyaltyMongoDB.HostIP_4 = Configuration.MongoDB.HostIP_4
	Configuration.LoyaltyMongoDB.HostPort_4 = Configuration.MongoDB.HostPort_4

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
	Configuration.SMPP.Login = "Loyalty"
	Configuration.SMPP.Password = "Loyalty123"
	Configuration.SMPP.TimeOut = 5 //in seconds
	Configuration.SMPP.PrintLogs = true
	Configuration.SMPP.MSISDN_Short_len = 8
	Configuration.SMPP.CountryCodePrefix = "220"
	Configuration.SMPP.DefaultSender = "Africell" //"Africell"
	Configuration.SMPP.Encoding = 0

	Configuration.CGW_AUC.Description = "UCGW AUC service"
	Configuration.CGW_AUC.Protocol = "http"
	Configuration.CGW_AUC.Hostname = "10.30.0.140"
	Configuration.CGW_AUC.Port = "9994"
	Configuration.CGW_AUC.Module = "AUC"
	Configuration.CGW_AUC.Version = "V1"
	Configuration.CGW_AUC.S2S_Username = "UCGW_Admin"
	Configuration.CGW_AUC.S2S_Password = "uC@g$W$iRiS6$2"
	Configuration.CGW_AUC.Timeout_After = 5 * time.Second

	Configuration.CGW.Protocol = "http"
	Configuration.CGW.Hostname = "10.30.0.140"
	Configuration.CGW.Port = "9991"
	Configuration.CGW.Module = "UCGW"
	Configuration.CGW.Version = "V1"
	Configuration.CGW.Timeout = 15 * time.Second

	Configuration.Propylaea.Description = "Product Design Center - Propylaea"
	Configuration.Propylaea.Protocol = "http"
	Configuration.Propylaea.Port = "9900"
	Configuration.Propylaea.Hostname = "10.30.0.140"
	Configuration.Propylaea.Module = "Propylaea"
	Configuration.Propylaea.Version = "V1"
	Configuration.Propylaea.S2S_Username = "Propylaea_Admin"
	Configuration.Propylaea.S2S_Password = "uC@g$W$iRiS6$2@333dd"
	Configuration.Propylaea.Timeout_After = 5 * time.Second
	Configuration.Propylaea.ChannelName = "Spin And Win"
	Configuration.Propylaea.ChannelPlan = "Normal SIM"
	Configuration.Propylaea.ChannelVersion = "1"

	Configuration.SpinAndWin_AUC.Description = "SAW AUC"
	Configuration.SpinAndWin_AUC.Protocol = "http"
	Configuration.SpinAndWin_AUC.Hostname = "10.30.0.120"
	Configuration.SpinAndWin_AUC.Port = "9102"
	Configuration.SpinAndWin_AUC.Module = "AUC"
	Configuration.SpinAndWin_AUC.Version = "V1"
	Configuration.SpinAndWin_AUC.S2S_Username = "SAW_Admin"
	Configuration.SpinAndWin_AUC.S2S_Password = "LQaDUp388UNKhz0Ap"
	Configuration.SpinAndWin_AUC.Timeout_After = 30 * time.Second

	Configuration.SpinAndWin.Protocol = "http"
	Configuration.SpinAndWin.Hostname = "10.30.0.120"
	Configuration.SpinAndWin.Port = "9112"
	Configuration.SpinAndWin.Module = "SpinAndWin"
	Configuration.SpinAndWin.Version = "V1"
	Configuration.SpinAndWin.Timeout = 30 * time.Second

	return
}

func setDefaultConfiguration_GM_Loyalty_Live() (Configuration ConfigType) {
	//Configuration.HttpOKAPIServicePort = "9291"
	Configuration.HttpAppServicePort = "9290"           //lendme services
	Configuration.HttpAppLoyaltyServicePort = "9280"    //for USSD and Mobile App
	Configuration.HttpAppLoyaltyManagementPort = "9281" //for OKAPI
	Configuration.HttpAppLoyaltyFeedPort = "9282"       //for IN & MM live feed

	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:3000")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:5173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4414")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://okpaihr.africell.gm")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://outlet.africell.gm")

	Configuration.Operation = "Gambia"
	Configuration.HostId = "Lendme-01"
	Configuration.DB_Name = "Lendme_DB"

	Configuration.Version = "V1"
	Configuration.Module = "Lendme"

	Configuration.LoyaltyVersion = "V1"
	Configuration.LoyaltyModule = "Loyalty"

	Configuration.MSISDN_Prefix = ""
	Configuration.MSISDN_Short_len = 7
	Configuration.CountryCode = "220"

	Configuration.IsProduction = true
	Configuration.IsLoyaltyProduction = true
	Configuration.Min_Allowed_Amnt = 5
	Configuration.Service_FeePerc = 0.04
	Configuration.Min_Allowed_AON = 3
	Configuration.Min_Avg3MRecharge = 5
	Configuration.Min_LastRechargePeriod = 60
	Configuration.Min_Allowed_Balance = 0
	Configuration.Max_Allowed_Balance = 677
	Configuration.ARPU_File_Path = "/home/Subs_ARPU/"
	Configuration.ARPU_File_Prefix = "Rgs_"
	Configuration.ARPU_File_Column_Separator = ","

	Configuration.App_AUC.Description = "App AUC service"
	Configuration.App_AUC.Protocol = "http"
	Configuration.App_AUC.Hostname = "Lendme_auc"
	Configuration.App_AUC.Port = "9293"
	Configuration.App_AUC.Module = "AUC"
	Configuration.App_AUC.Version = "V1"
	Configuration.App_AUC.S2S_Username = "Lendme_Admin"
	Configuration.App_AUC.S2S_Password = "s@l$e$IrSW0$4"
	Configuration.App_AUC.Timeout_After = 5 * time.Second

	Configuration.OKAPI_AUC.Description = "OKAPI AUC service"
	Configuration.OKAPI_AUC.Protocol = "http"
	//Configuration.OKAPI_AUC.Hostname = "auc"
	Configuration.OKAPI_AUC.Hostname = "10.30.0.120"
	Configuration.OKAPI_AUC.Port = "9001"
	Configuration.OKAPI_AUC.Module = "AUC"
	Configuration.OKAPI_AUC.Version = "V1"
	Configuration.OKAPI_AUC.S2S_Username = "SalesMonitoring_OKAPI"
	Configuration.OKAPI_AUC.S2S_Password = "s@le$P@s$W0$3"
	Configuration.OKAPI_AUC.Timeout_After = 5 * time.Second

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

	// Configuration.MongoDB.ReplicaSet = ""
	// Configuration.MongoDB.UserName = "db_root"
	// Configuration.MongoDB.Password = "P@s54D0Brdara_r@75S"
	// Configuration.MongoDB.HostIP_1 = "10.30.0.151" //==>Primary
	// Configuration.MongoDB.HostPort_1 = "9510"

	Configuration.DB_Name_Loyalty = "Loyalty_DB"
	Configuration.LoyaltyMongoDB.ReplicaSet = Configuration.MongoDB.ReplicaSet
	Configuration.LoyaltyMongoDB.UserName = Configuration.MongoDB.UserName
	Configuration.LoyaltyMongoDB.Password = Configuration.MongoDB.Password
	Configuration.LoyaltyMongoDB.HostIP_1 = Configuration.MongoDB.HostIP_1
	Configuration.LoyaltyMongoDB.HostPort_1 = Configuration.MongoDB.HostPort_1
	Configuration.LoyaltyMongoDB.HostIP_2 = Configuration.MongoDB.HostIP_2
	Configuration.LoyaltyMongoDB.HostPort_2 = Configuration.MongoDB.HostPort_2
	Configuration.LoyaltyMongoDB.HostIP_3 = Configuration.MongoDB.HostIP_3
	Configuration.LoyaltyMongoDB.HostPort_3 = Configuration.MongoDB.HostPort_3
	Configuration.LoyaltyMongoDB.HostIP_4 = Configuration.MongoDB.HostIP_4
	Configuration.LoyaltyMongoDB.HostPort_4 = Configuration.MongoDB.HostPort_4

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
	Configuration.SMPP.Login = "Loyalty"
	Configuration.SMPP.Password = "Loyalty123"
	Configuration.SMPP.TimeOut = 5 //in seconds
	Configuration.SMPP.PrintLogs = true
	Configuration.SMPP.MSISDN_Short_len = 8
	Configuration.SMPP.CountryCodePrefix = "220"
	Configuration.SMPP.DefaultSender = "Africell" //"Africell"
	Configuration.SMPP.Encoding = 0

	Configuration.CGW_AUC.Description = "UCGW AUC service"
	Configuration.CGW_AUC.Protocol = "http"
	Configuration.CGW_AUC.Hostname = "10.30.0.140"
	Configuration.CGW_AUC.Port = "9994"
	Configuration.CGW_AUC.Module = "AUC"
	Configuration.CGW_AUC.Version = "V1"
	Configuration.CGW_AUC.S2S_Username = "UCGW_Admin"
	Configuration.CGW_AUC.S2S_Password = "uC@g$W$iRiS6$2"
	Configuration.CGW_AUC.Timeout_After = 5 * time.Second

	Configuration.CGW.Protocol = "http"
	Configuration.CGW.Hostname = "10.30.0.140"
	Configuration.CGW.Port = "9991"
	Configuration.CGW.Module = "UCGW"
	Configuration.CGW.Version = "V1"
	Configuration.CGW.Timeout = 15 * time.Second

	Configuration.Propylaea.Description = "Product Design Center - Propylaea"
	Configuration.Propylaea.Protocol = "http"
	Configuration.Propylaea.Port = "9900"
	Configuration.Propylaea.Hostname = "10.30.0.140"
	Configuration.Propylaea.Module = "Propylaea"
	Configuration.Propylaea.Version = "V1"
	Configuration.Propylaea.S2S_Username = "Propylaea_Admin"
	Configuration.Propylaea.S2S_Password = "uC@g$W$iRiS6$2@333dd"
	Configuration.Propylaea.Timeout_After = 5 * time.Second
	Configuration.Propylaea.ChannelName = "Spin And Win"
	Configuration.Propylaea.ChannelPlan = "Normal SIM"
	Configuration.Propylaea.ChannelVersion = "1"

	Configuration.SpinAndWin_AUC.Description = "SAW AUC"
	Configuration.SpinAndWin_AUC.Protocol = "http"
	Configuration.SpinAndWin_AUC.Hostname = "10.30.0.120"
	Configuration.SpinAndWin_AUC.Port = "9102"
	Configuration.SpinAndWin_AUC.Module = "AUC"
	Configuration.SpinAndWin_AUC.Version = "V1"
	Configuration.SpinAndWin_AUC.S2S_Username = "SAW_Admin"
	Configuration.SpinAndWin_AUC.S2S_Password = "LQaDUp388UNKhz0Ap"
	Configuration.SpinAndWin_AUC.Timeout_After = 30 * time.Second

	Configuration.SpinAndWin.Protocol = "http"
	Configuration.SpinAndWin.Hostname = "10.30.0.120"
	Configuration.SpinAndWin.Port = "9112"
	Configuration.SpinAndWin.Module = "SpinAndWin"
	Configuration.SpinAndWin.Version = "V1"
	Configuration.SpinAndWin.Timeout = 30 * time.Second

	return
}

func setDefaultConfiguration_GM_Loyalty_UAT() (Configuration ConfigType) {
	//Configuration.HttpOKAPIServicePort = "9291"
	Configuration.HttpAppServicePort = "9290"           //lendme services
	Configuration.HttpAppLoyaltyServicePort = "9280"    //for USSD and Mobile App
	Configuration.HttpAppLoyaltyManagementPort = "9281" //for OKAPI
	Configuration.HttpAppLoyaltyFeedPort = "9282"       //for IN & MM live feed

	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:3000")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:5173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4414")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://okpaihruat.africell.gm")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://outlet.africell.ao")

	Configuration.Operation = "Gambia"
	Configuration.HostId = "Lendme-01"
	Configuration.DB_Name = "Lendme_DB"

	Configuration.Version = "V1"
	Configuration.Module = "Lendme"

	Configuration.LoyaltyVersion = "V1"
	Configuration.LoyaltyModule = "Loyalty"

	Configuration.MSISDN_Prefix = ""
	Configuration.MSISDN_Short_len = 7
	Configuration.CountryCode = "220"

	Configuration.IsProduction = true
	Configuration.IsLoyaltyProduction = true
	Configuration.Min_Allowed_Amnt = 5
	Configuration.Service_FeePerc = 0.04
	Configuration.Min_Allowed_AON = 3
	Configuration.Min_Avg3MRecharge = 5
	Configuration.Min_LastRechargePeriod = 60
	Configuration.Min_Allowed_Balance = 0
	Configuration.Max_Allowed_Balance = 677
	Configuration.ARPU_File_Path = "/home/Subs_ARPU/"
	Configuration.ARPU_File_Prefix = "Rgs_"
	Configuration.ARPU_File_Column_Separator = ","

	Configuration.App_AUC.Description = "App AUC service"
	Configuration.App_AUC.Protocol = "http"
	Configuration.App_AUC.Hostname = "Lendme_auc"
	Configuration.App_AUC.Port = "9293"
	Configuration.App_AUC.Module = "AUC"
	Configuration.App_AUC.Version = "V1"
	Configuration.App_AUC.S2S_Username = "Lendme_Admin"
	Configuration.App_AUC.S2S_Password = "s@l$e$IrSW0$4"
	Configuration.App_AUC.Timeout_After = 5 * time.Second

	Configuration.OKAPI_AUC.Description = "OKAPI AUC service"
	Configuration.OKAPI_AUC.Protocol = "http"
	//Configuration.OKAPI_AUC.Hostname = "auc"
	Configuration.OKAPI_AUC.Hostname = "10.30.0.119"
	Configuration.OKAPI_AUC.Port = "9001"
	Configuration.OKAPI_AUC.Module = "AUC"
	Configuration.OKAPI_AUC.Version = "V1"
	Configuration.OKAPI_AUC.S2S_Username = "SalesMonitoring_OKAPI"
	Configuration.OKAPI_AUC.S2S_Password = "s@le$P@s$W0$3"
	Configuration.OKAPI_AUC.Timeout_After = 5 * time.Second

	//mongoDB
	// Configuration.MongoDB.ReplicaSet = "reps01"
	// Configuration.MongoDB.UserName = "mongo-root"
	// Configuration.MongoDB.Password = "Speci@LM0nG0P@ssw0rd_F0r_G@mB!A"
	// Configuration.MongoDB.HostIP_1 = "10.64.33.49" //==>Primary
	// Configuration.MongoDB.HostPort_1 = "9001"
	// Configuration.MongoDB.HostIP_2 = "10.64.33.48" //==>Secondary
	// Configuration.MongoDB.HostPort_2 = "9002"
	// Configuration.MongoDB.HostIP_3 = "10.64.33.101" //==> Aribter
	// Configuration.MongoDB.HostPort_3 = "9003"
	// Configuration.MongoDB.HostIP_4 = ""
	// Configuration.MongoDB.HostPort_4 = ""

	Configuration.MongoDB.ReplicaSet = ""
	Configuration.MongoDB.UserName = "db_root"
	Configuration.MongoDB.Password = "P@s54D0Brdara_r@75S"
	Configuration.MongoDB.HostIP_1 = "10.30.0.151" //==>Primary
	Configuration.MongoDB.HostPort_1 = "9510"

	Configuration.DB_Name_Loyalty = "Loyalty_DB"
	Configuration.LoyaltyMongoDB.ReplicaSet = Configuration.MongoDB.ReplicaSet
	Configuration.LoyaltyMongoDB.UserName = Configuration.MongoDB.UserName
	Configuration.LoyaltyMongoDB.Password = Configuration.MongoDB.Password
	Configuration.LoyaltyMongoDB.HostIP_1 = Configuration.MongoDB.HostIP_1
	Configuration.LoyaltyMongoDB.HostPort_1 = Configuration.MongoDB.HostPort_1
	Configuration.LoyaltyMongoDB.HostIP_2 = Configuration.MongoDB.HostIP_2
	Configuration.LoyaltyMongoDB.HostPort_2 = Configuration.MongoDB.HostPort_2
	Configuration.LoyaltyMongoDB.HostIP_3 = Configuration.MongoDB.HostIP_3
	Configuration.LoyaltyMongoDB.HostPort_3 = Configuration.MongoDB.HostPort_3
	Configuration.LoyaltyMongoDB.HostIP_4 = Configuration.MongoDB.HostIP_4
	Configuration.LoyaltyMongoDB.HostPort_4 = Configuration.MongoDB.HostPort_4

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
	Configuration.SMPP.Login = "Loyalty"
	Configuration.SMPP.Password = "Loyalty123"
	Configuration.SMPP.TimeOut = 5 //in seconds
	Configuration.SMPP.PrintLogs = true
	Configuration.SMPP.MSISDN_Short_len = 8
	Configuration.SMPP.CountryCodePrefix = "220"
	Configuration.SMPP.DefaultSender = "Africell" //"Africell"
	Configuration.SMPP.Encoding = 0

	Configuration.CGW_AUC.Description = "UCGW AUC service"
	Configuration.CGW_AUC.Protocol = "http"
	Configuration.CGW_AUC.Hostname = "10.30.0.119"
	Configuration.CGW_AUC.Port = "9994"
	Configuration.CGW_AUC.Module = "AUC"
	Configuration.CGW_AUC.Version = "V1"
	Configuration.CGW_AUC.S2S_Username = "UCGW_Admin"
	Configuration.CGW_AUC.S2S_Password = "uC@g$W$iRiS6$2"
	Configuration.CGW_AUC.Timeout_After = 5 * time.Second

	Configuration.CGW.Protocol = "http"
	Configuration.CGW.Hostname = "10.30.0.119"
	Configuration.CGW.Port = "9991"
	Configuration.CGW.Module = "UCGW"
	Configuration.CGW.Version = "V1"
	Configuration.CGW.Timeout = 15 * time.Second

	Configuration.Propylaea.Description = "Product Design Center - Propylaea"
	Configuration.Propylaea.Protocol = "http"
	Configuration.Propylaea.Port = "9900"
	Configuration.Propylaea.Hostname = "10.30.0.119"
	Configuration.Propylaea.Module = "Propylaea"
	Configuration.Propylaea.Version = "V1"
	Configuration.Propylaea.S2S_Username = "Propylaea_Admin"
	Configuration.Propylaea.S2S_Password = "uC@g$W$iRiS6$2@333dd"
	Configuration.Propylaea.Timeout_After = 5 * time.Second
	Configuration.Propylaea.ChannelName = "Spin And Win"
	Configuration.Propylaea.ChannelPlan = "Normal SIM"
	Configuration.Propylaea.ChannelVersion = "1"

	Configuration.SpinAndWin_AUC.Description = "SAW AUC"
	Configuration.SpinAndWin_AUC.Protocol = "http"
	Configuration.SpinAndWin_AUC.Hostname = "10.30.0.119"
	Configuration.SpinAndWin_AUC.Port = "9102"
	Configuration.SpinAndWin_AUC.Module = "AUC"
	Configuration.SpinAndWin_AUC.Version = "V1"
	Configuration.SpinAndWin_AUC.S2S_Username = "SAW_Admin"
	Configuration.SpinAndWin_AUC.S2S_Password = "LQaDUp388UNKhz0Ap"
	Configuration.SpinAndWin_AUC.Timeout_After = 30 * time.Second

	Configuration.SpinAndWin.Protocol = "http"
	Configuration.SpinAndWin.Hostname = "10.30.0.119"
	Configuration.SpinAndWin.Port = "9112"
	Configuration.SpinAndWin.Module = "SpinAndWin"
	Configuration.SpinAndWin.Version = "V1"
	Configuration.SpinAndWin.Timeout = 30 * time.Second

	return
}

func setDefaultConfiguration_SL_Live() (Configuration ConfigType) {
	//Configuration.HttpOKAPIServicePort = "9291"
	Configuration.HttpAppServicePort = "9290"           //lendme services
	Configuration.HttpAppLoyaltyServicePort = "9280"    //for USSD and Mobile App
	Configuration.HttpAppLoyaltyManagementPort = "9281" //for OKAPI
	Configuration.HttpAppLoyaltyFeedPort = "9282"       //for IN & MM live feed

	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:3000")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:5173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4414")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://okapihr.africell.sl")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://okapihruat.africell.sl")

	Configuration.Operation = "SierraLeone"
	Configuration.HostId = "Lendme-01"
	Configuration.DB_Name = "Lendme_DB"

	Configuration.Version = "V1"
	Configuration.Module = "Lendme"

	Configuration.LoyaltyVersion = "V1"
	Configuration.LoyaltyModule = "Loyalty"

	Configuration.MSISDN_Prefix = "0"
	Configuration.MSISDN_Short_len = 8
	Configuration.CountryCode = "232"

	Configuration.IsProduction = true
	Configuration.IsLoyaltyProduction = false
	Configuration.Min_Allowed_Amnt = 1
	Configuration.Service_FeePerc = 0.15
	Configuration.Min_Allowed_AON = 3
	Configuration.Min_Avg3MRecharge = 0
	Configuration.Min_LastRechargePeriod = 60
	Configuration.Min_Allowed_Balance = 0
	Configuration.Max_Allowed_Balance = 225
	Configuration.ARPU_File_Path = "/home/Subs_ARPU/"
	Configuration.ARPU_File_Prefix = "Rgs_"
	Configuration.ARPU_File_Column_Separator = ","

	Configuration.App_AUC.Description = "App AUC service"
	Configuration.App_AUC.Protocol = "http"
	Configuration.App_AUC.Hostname = "Lendme_auc"
	Configuration.App_AUC.Port = "9293"
	Configuration.App_AUC.Module = "AUC"
	Configuration.App_AUC.Version = "V1"
	Configuration.App_AUC.S2S_Username = "Lendme_Admin"
	Configuration.App_AUC.S2S_Password = "s@l$e$IrSW0$4"
	Configuration.App_AUC.Timeout_After = 5 * time.Second

	Configuration.OKAPI_AUC.Description = "OKAPI AUC service"
	Configuration.OKAPI_AUC.Protocol = "http"
	Configuration.OKAPI_AUC.Hostname = "10.10.234.56"
	Configuration.OKAPI_AUC.Port = "9001"
	Configuration.OKAPI_AUC.Module = "AUC"
	Configuration.OKAPI_AUC.Version = "V1"
	Configuration.OKAPI_AUC.S2S_Username = "Loyalty_OKAPI"
	Configuration.OKAPI_AUC.S2S_Password = "]W8#x3D1USKUyH@p]s&D_"
	Configuration.OKAPI_AUC.Timeout_After = 5 * time.Second

	// //mongoDB
	Configuration.MongoDB.ReplicaSet = "reps0"
	Configuration.MongoDB.UserName = "LendMeApp"
	Configuration.MongoDB.Password = "LendMeApp_@@!!"
	Configuration.MongoDB.HostIP_1 = "10.10.247.21" //==>Primary
	Configuration.MongoDB.HostPort_1 = "9771"
	Configuration.MongoDB.HostIP_2 = "10.10.247.22" //==>Secondary
	Configuration.MongoDB.HostPort_2 = "9772"
	Configuration.MongoDB.HostIP_3 = "10.10.247.20" //==> Aribter
	Configuration.MongoDB.HostPort_3 = "9773"
	Configuration.MongoDB.HostIP_4 = ""
	Configuration.MongoDB.HostPort_4 = ""
	////
	// Configuration.MongoDB.ReplicaSet = ""
	// Configuration.MongoDB.UserName = "mongo-root"
	// Configuration.MongoDB.Password = "Speci@LM0nG0P@ssw0rd_F0r_LeNdM#SL"
	// Configuration.MongoDB.HostIP_1 = "10.10.247.21" //==>Primary
	// Configuration.MongoDB.HostPort_1 = "9001"
	// Configuration.MongoDB.HostIP_2 = "" //==>Secondary
	// Configuration.MongoDB.HostPort_2 = ""
	// Configuration.MongoDB.HostIP_3 = "" //==> Aribter
	// Configuration.MongoDB.HostPort_3 = ""
	// Configuration.MongoDB.HostIP_4 = ""
	// Configuration.MongoDB.HostPort_4 = ""

	//mongoDB
	// Configuration.MongoDB.UserName = "db_root"
	// Configuration.MongoDB.Password = "P@s54D0Brdara_r@75S"
	// Configuration.MongoDB.HostIP_1 = "LendMe_db" //"host.docker.internal"
	// Configuration.MongoDB.HostPort_1 = "27017"   //"27017"
	///////////////////////////////////////MONGO DOCKER ///////////////////////////

	Configuration.DB_Name_Loyalty = "Loyalty_DB"
	Configuration.LoyaltyMongoDB.ReplicaSet = Configuration.MongoDB.ReplicaSet
	Configuration.LoyaltyMongoDB.UserName = Configuration.MongoDB.UserName
	Configuration.LoyaltyMongoDB.Password = Configuration.MongoDB.Password
	Configuration.LoyaltyMongoDB.HostIP_1 = Configuration.MongoDB.HostIP_1
	Configuration.LoyaltyMongoDB.HostPort_1 = Configuration.MongoDB.HostPort_1
	Configuration.LoyaltyMongoDB.HostIP_2 = Configuration.MongoDB.HostIP_2
	Configuration.LoyaltyMongoDB.HostPort_2 = Configuration.MongoDB.HostPort_2
	Configuration.LoyaltyMongoDB.HostIP_3 = Configuration.MongoDB.HostIP_3
	Configuration.LoyaltyMongoDB.HostPort_3 = Configuration.MongoDB.HostPort_3
	Configuration.LoyaltyMongoDB.HostIP_4 = Configuration.MongoDB.HostIP_4
	Configuration.LoyaltyMongoDB.HostPort_4 = Configuration.MongoDB.HostPort_4

	Configuration.IN.IP = "10.10.51.51"
	Configuration.IN.Port = "8080"
	Configuration.IN.WS_SOAP_Endpoint = "/axis2/services/WebService.WebServiceHttpSoap12Endpoint/"
	Configuration.IN.WS_XMLNS_SOAP_Env = "http://schemas.xmlsoap.org/soap/envelope/"
	Configuration.IN.WS_XMLNS_Web = "http://webservice.CSI.omvia.convergys.com"

	Configuration.IN.WS_EVC_SOAP_Endpoint = "/axis2/services/ERechargeWebService.ERechargeWebServiceHttpSoap11Endpoint/"
	Configuration.IN.WS_EVC_XMLNS_SOAP_Env = "http://schemas.xmlsoap.org/soap/envelope/"
	Configuration.IN.WS_EVC_XMLNS_Web = "http://webservice.CSI.omvia.convergys.com"

	Configuration.IN.Default_OpId = "lendme"
	Configuration.IN.Default_OpPwd = "gu1vY6Q$"
	Configuration.IN.Is_OpPwd_Required = true
	Configuration.IN.Timeout = 5
	Configuration.IN.PrintLogs = true

	//http://10.95.64.6:15403/?systemid=lendme&password=lendmeP@ssw0rd&Originator=setest&dest_addr=243900100606&msg_text=test&registered_delivery=0&ston=5&snpi=0&dton=1&dnpi=1&encoding=1

	//SMPP
	Configuration.SMPP.IP = "10.10.215.52"
	Configuration.SMPP.Port = "15403"

	//Configuration.SMPP.Login = "lendme"
	//Configuration.SMPP.Password = "lendmeP@ssw0rd"
	Configuration.SMPP.Login = "LendME2"
	Configuration.SMPP.Password = "LendMEP@ssw0rd"

	//CGW
	Configuration.CGW.Protocol = "http"
	Configuration.CGW.Hostname = "10.10.231.51"
	Configuration.CGW.Port = "9991"
	Configuration.CGW.Module = "UCGW"
	Configuration.CGW.Version = "V1"
	Configuration.CGW.Timeout = 15 * time.Second
	Configuration.CGW.S2S_AccessToken = ""

	Configuration.CGW_AUC.Description = "UCGW AUC service"
	Configuration.CGW_AUC.Protocol = "http"
	Configuration.CGW_AUC.Hostname = "10.10.231.51"
	Configuration.CGW_AUC.Port = "9994"
	Configuration.CGW_AUC.Module = "AUC"
	Configuration.CGW_AUC.Version = "V1"
	Configuration.CGW_AUC.S2S_Username = "UCGW_Admin"
	Configuration.CGW_AUC.S2S_Password = "uC@g$W$iRiS6$2"
	Configuration.CGW_AUC.Timeout_After = 5 * time.Second

	Configuration.Propylaea.Description = "Product Design Center - Propylaea"
	Configuration.Propylaea.Protocol = "http"
	Configuration.Propylaea.Port = "9900"
	Configuration.Propylaea.Hostname = "10.10.231.51"
	Configuration.Propylaea.Module = "Propylaea"
	Configuration.Propylaea.Version = "V1"
	Configuration.Propylaea.S2S_Username = "Propylaea_Admin"
	Configuration.Propylaea.S2S_Password = "uC@g$W$iRiS6$2@333dd"
	Configuration.Propylaea.Timeout_After = 5 * time.Second

	return
}

func setDefaultConfiguration_SL_Loyalty() (Configuration ConfigType) {
	//Configuration.HttpOKAPIServicePort = "9291"
	Configuration.HttpAppServicePort = "9290"           //lendme services
	Configuration.HttpAppLoyaltyServicePort = "9280"    //for USSD and Mobile App
	Configuration.HttpAppLoyaltyManagementPort = "9281" //for OKAPI
	Configuration.HttpAppLoyaltyFeedPort = "9282"       //for IN & MM live feed

	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:3000")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:5173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4414")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://okapihr.africell.sl")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://okapihruat.africell.sl")

	Configuration.Operation = "SierraLeone"
	Configuration.HostId = "Lendme-01"
	Configuration.DB_Name = "Lendme_DB"

	Configuration.Version = "V1"
	Configuration.Module = "Lendme"

	Configuration.LoyaltyVersion = "V1"
	Configuration.LoyaltyModule = "Loyalty"

	Configuration.MSISDN_Prefix = "0"
	Configuration.MSISDN_Short_len = 8
	Configuration.CountryCode = "232"

	Configuration.IsProduction = true
	Configuration.IsLoyaltyProduction = false
	Configuration.Min_Allowed_Amnt = 1
	Configuration.Service_FeePerc = 0.15
	Configuration.Min_Allowed_AON = 3
	Configuration.Min_Avg3MRecharge = 0
	Configuration.Min_LastRechargePeriod = 60
	Configuration.Min_Allowed_Balance = 0
	Configuration.Max_Allowed_Balance = 225
	Configuration.ARPU_File_Path = "/home/Subs_ARPU/"
	Configuration.ARPU_File_Prefix = "Rgs_"
	Configuration.ARPU_File_Column_Separator = ","

	Configuration.App_AUC.Description = "App AUC service"
	Configuration.App_AUC.Protocol = "http"
	Configuration.App_AUC.Hostname = "Lendme_auc"
	Configuration.App_AUC.Port = "9293"
	Configuration.App_AUC.Module = "AUC"
	Configuration.App_AUC.Version = "V1"
	Configuration.App_AUC.S2S_Username = "Lendme_Admin"
	Configuration.App_AUC.S2S_Password = "s@l$e$IrSW0$4"
	Configuration.App_AUC.Timeout_After = 5 * time.Second

	Configuration.OKAPI_AUC.Description = "OKAPI AUC service"
	Configuration.OKAPI_AUC.Protocol = "http"
	Configuration.OKAPI_AUC.Hostname = "10.10.234.56"
	Configuration.OKAPI_AUC.Port = "9001"
	Configuration.OKAPI_AUC.Module = "AUC"
	Configuration.OKAPI_AUC.Version = "V1"
	Configuration.OKAPI_AUC.S2S_Username = "Loyalty_OKAPI"
	Configuration.OKAPI_AUC.S2S_Password = "]W8#x3D1USKUyH@p]s&D_"
	Configuration.OKAPI_AUC.Timeout_After = 5 * time.Second

	// //mongoDB
	Configuration.MongoDB.ReplicaSet = "reps0"
	Configuration.MongoDB.UserName = "LendMeApp"
	Configuration.MongoDB.Password = "LendMeApp_@@!!"
	Configuration.MongoDB.HostIP_1 = "10.10.247.21" //==>Primary
	Configuration.MongoDB.HostPort_1 = "9771"
	Configuration.MongoDB.HostIP_2 = "10.10.247.22" //==>Secondary
	Configuration.MongoDB.HostPort_2 = "9772"
	Configuration.MongoDB.HostIP_3 = "10.10.231.20" //==> Aribter
	Configuration.MongoDB.HostPort_3 = "9773"
	Configuration.MongoDB.HostIP_4 = ""
	Configuration.MongoDB.HostPort_4 = ""
	////
	// Configuration.MongoDB.ReplicaSet = ""
	// Configuration.MongoDB.UserName = "mongo-root"
	// Configuration.MongoDB.Password = "Speci@LM0nG0P@ssw0rd_F0r_LeNdM#SL"
	// Configuration.MongoDB.HostIP_1 = "10.10.247.21" //==>Primary
	// Configuration.MongoDB.HostPort_1 = "9001"
	// Configuration.MongoDB.HostIP_2 = "" //==>Secondary
	// Configuration.MongoDB.HostPort_2 = ""
	// Configuration.MongoDB.HostIP_3 = "" //==> Aribter
	// Configuration.MongoDB.HostPort_3 = ""
	// Configuration.MongoDB.HostIP_4 = ""
	// Configuration.MongoDB.HostPort_4 = ""

	// //mongoDB
	// Configuration.MongoDB.UserName = "db_root"
	// Configuration.MongoDB.Password = "P@s54D0Brdara_r@75S"
	// Configuration.MongoDB.HostIP_1 = "LendMe_db" //"host.docker.internal"
	// Configuration.MongoDB.HostPort_1 = "27017"   //"27017"
	///////////////////////////////////////MONGO DOCKER ///////////////////////////

	Configuration.DB_Name_Loyalty = "Loyalty_DB"
	Configuration.LoyaltyMongoDB.ReplicaSet = Configuration.MongoDB.ReplicaSet
	Configuration.LoyaltyMongoDB.UserName = Configuration.MongoDB.UserName
	Configuration.LoyaltyMongoDB.Password = Configuration.MongoDB.Password
	Configuration.LoyaltyMongoDB.HostIP_1 = Configuration.MongoDB.HostIP_1
	Configuration.LoyaltyMongoDB.HostPort_1 = Configuration.MongoDB.HostPort_1
	Configuration.LoyaltyMongoDB.HostIP_2 = Configuration.MongoDB.HostIP_2
	Configuration.LoyaltyMongoDB.HostPort_2 = Configuration.MongoDB.HostPort_2
	Configuration.LoyaltyMongoDB.HostIP_3 = Configuration.MongoDB.HostIP_3
	Configuration.LoyaltyMongoDB.HostPort_3 = Configuration.MongoDB.HostPort_3
	Configuration.LoyaltyMongoDB.HostIP_4 = Configuration.MongoDB.HostIP_4
	Configuration.LoyaltyMongoDB.HostPort_4 = Configuration.MongoDB.HostPort_4

	Configuration.IN.IP = "10.10.51.51"
	Configuration.IN.Port = "8080"
	Configuration.IN.WS_SOAP_Endpoint = "/axis2/services/WebService.WebServiceHttpSoap12Endpoint/"
	Configuration.IN.WS_XMLNS_SOAP_Env = "http://schemas.xmlsoap.org/soap/envelope/"
	Configuration.IN.WS_XMLNS_Web = "http://webservice.CSI.omvia.convergys.com"

	Configuration.IN.WS_EVC_SOAP_Endpoint = "/axis2/services/ERechargeWebService.ERechargeWebServiceHttpSoap11Endpoint/"
	Configuration.IN.WS_EVC_XMLNS_SOAP_Env = "http://schemas.xmlsoap.org/soap/envelope/"
	Configuration.IN.WS_EVC_XMLNS_Web = "http://webservice.CSI.omvia.convergys.com"

	Configuration.IN.Default_OpId = "lendme"
	Configuration.IN.Default_OpPwd = "gu1vY6Q$"
	Configuration.IN.Is_OpPwd_Required = true
	Configuration.IN.Timeout = 5
	Configuration.IN.PrintLogs = true

	//http://10.95.64.6:15403/?systemid=lendme&password=lendmeP@ssw0rd&Originator=setest&dest_addr=243900100606&msg_text=test&registered_delivery=0&ston=5&snpi=0&dton=1&dnpi=1&encoding=1

	//SMPP
	Configuration.SMPP.IP = "10.10.215.52"
	Configuration.SMPP.Port = "15403"

	//Configuration.SMPP.Login = "lendme"
	//Configuration.SMPP.Password = "lendmeP@ssw0rd"
	Configuration.SMPP.Login = "Loyalty"
	Configuration.SMPP.Password = "Loyalty123"
	Configuration.SMPP.TimeOut = 5 //in seconds
	Configuration.SMPP.PrintLogs = true
	Configuration.SMPP.MSISDN_Short_len = 8
	Configuration.SMPP.CountryCodePrefix = "232"
	Configuration.SMPP.DefaultSender = "Africell" //"Africell"
	Configuration.SMPP.Encoding = 0
	//CGW
	Configuration.CGW.Protocol = "http"
	Configuration.CGW.Hostname = "10.10.231.51"
	Configuration.CGW.Port = "9991"
	Configuration.CGW.Module = "UCGW"
	Configuration.CGW.Version = "V1"
	Configuration.CGW.Timeout = 15 * time.Second
	Configuration.CGW.S2S_AccessToken = ""

	Configuration.CGW_AUC.Description = "UCGW AUC service"
	Configuration.CGW_AUC.Protocol = "http"
	Configuration.CGW_AUC.Hostname = "10.10.231.51"
	Configuration.CGW_AUC.Port = "9994"
	Configuration.CGW_AUC.Module = "AUC"
	Configuration.CGW_AUC.Version = "V1"
	Configuration.CGW_AUC.S2S_Username = "UCGW_Admin"
	Configuration.CGW_AUC.S2S_Password = "uC@g$W$iRiS6$2"
	Configuration.CGW_AUC.Timeout_After = 5 * time.Second

	Configuration.Propylaea.Description = "Product Design Center - Propylaea"
	Configuration.Propylaea.Protocol = "http"
	Configuration.Propylaea.Port = "9900"
	Configuration.Propylaea.Hostname = "10.10.231.51"
	Configuration.Propylaea.Module = "Propylaea"
	Configuration.Propylaea.Version = "V1"
	Configuration.Propylaea.S2S_Username = "Propylaea_Admin"
	Configuration.Propylaea.S2S_Password = "uC@g$W$iRiS6$2@333dd"
	Configuration.Propylaea.Timeout_After = 5 * time.Second

	return
}

func setDefaultConfiguration_SL_Loyalty_UAT() (Configuration ConfigType) {
	// Configuration.HttpOKAPIServicePort = "9291"
	Configuration.HttpAppServicePort = "9290"           //lendme services
	Configuration.HttpAppLoyaltyServicePort = "9280"    //for USSD and Mobile App
	Configuration.HttpAppLoyaltyManagementPort = "9281" //for OKAPI
	Configuration.HttpAppLoyaltyFeedPort = "9282"       //for IN & MM live feed

	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:3000")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:5173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4414")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://okpaihruat.africell.sl")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://okpaihruat.africell.sl")

	Configuration.Operation = "SierraLeone"
	Configuration.HostId = "Lendme-01"
	Configuration.DB_Name = "Lendme_DB"

	Configuration.Version = "V1"
	Configuration.Module = "Lendme"

	Configuration.LoyaltyVersion = "V1"
	Configuration.LoyaltyModule = "Loyalty"

	Configuration.MSISDN_Prefix = "0"
	Configuration.MSISDN_Short_len = 8
	Configuration.CountryCode = "232"

	Configuration.IsProduction = true
	Configuration.IsLoyaltyProduction = false
	Configuration.Min_Allowed_Amnt = 1
	Configuration.Service_FeePerc = 0.15
	Configuration.Min_Allowed_AON = 3
	Configuration.Min_Avg3MRecharge = 0
	Configuration.Min_LastRechargePeriod = 60
	Configuration.Min_Allowed_Balance = 0
	Configuration.Max_Allowed_Balance = 225
	Configuration.ARPU_File_Path = "/home/Subs_ARPU/"
	Configuration.ARPU_File_Prefix = "Rgs_"
	Configuration.ARPU_File_Column_Separator = ","

	Configuration.App_AUC.Description = "App AUC service"
	Configuration.App_AUC.Protocol = "http"
	Configuration.App_AUC.Hostname = "Lendme_auc"
	Configuration.App_AUC.Port = "9293"
	Configuration.App_AUC.Module = "AUC"
	Configuration.App_AUC.Version = "V1"
	Configuration.App_AUC.S2S_Username = "Lendme_Admin"
	Configuration.App_AUC.S2S_Password = "s@l$e$IrSW0$4"
	Configuration.App_AUC.Timeout_After = 5 * time.Second

	Configuration.OKAPI_AUC.Description = "OKAPI AUC service"
	Configuration.OKAPI_AUC.Protocol = "http"
	Configuration.OKAPI_AUC.Hostname = "auc"
	Configuration.OKAPI_AUC.Port = "9001"
	Configuration.OKAPI_AUC.Module = "AUC"
	Configuration.OKAPI_AUC.Version = "V1"
	Configuration.OKAPI_AUC.S2S_Username = "Loyalty_OKAPI"
	Configuration.OKAPI_AUC.S2S_Password = "]W8#x3D1USKUyH@p]s&D_"
	Configuration.OKAPI_AUC.Timeout_After = 5 * time.Second

	// //mongoDB
	// Configuration.MongoDB.ReplicaSet = "reps0"
	// Configuration.MongoDB.UserName = "LendMeApp"
	// Configuration.MongoDB.Password = "LendMeApp_@@!!"
	// Configuration.MongoDB.HostIP_1 = "10.10.247.21" //==>Primary
	// Configuration.MongoDB.HostPort_1 = "9771"
	// Configuration.MongoDB.HostIP_2 = "10.10.247.22" //==>Secondary
	// Configuration.MongoDB.HostPort_2 = "9772"
	// Configuration.MongoDB.HostIP_3 = "10.10.231.20" //==> Aribter
	// Configuration.MongoDB.HostPort_3 = "9773"
	// Configuration.MongoDB.HostIP_4 = ""
	// Configuration.MongoDB.HostPort_4 = ""
	////
	// Configuration.MongoDB.ReplicaSet = ""
	// Configuration.MongoDB.UserName = "mongo-root"
	// Configuration.MongoDB.Password = "Speci@LM0nG0P@ssw0rd_F0r_LeNdM#SL"
	// Configuration.MongoDB.HostIP_1 = "10.10.2" //==>Primary
	// Configuration.MongoDB.HostPort_1 = "9001"
	// Configuration.MongoDB.HostIP_2 = "" //==>Secondary
	// Configuration.MongoDB.HostPort_2 = ""
	// Configuration.MongoDB.HostIP_3 = "" //==> Aribter
	// Configuration.MongoDB.HostPort_3 = ""
	// Configuration.MongoDB.HostIP_4 = ""
	// Configuration.MongoDB.HostPort_4 = ""

	//mongoDB
	Configuration.MongoDB.UserName = "db_root"
	Configuration.MongoDB.Password = "P@s54D0Brdara_r@75S"
	Configuration.MongoDB.HostIP_1 = "mongodb" //"host.docker.internal"
	Configuration.MongoDB.HostPort_1 = "27017" //"27017"
	/////////////////////////////////////MONGO DOCKER ///////////////////////////

	Configuration.DB_Name_Loyalty = "Loyalty_DB"
	Configuration.LoyaltyMongoDB.ReplicaSet = Configuration.MongoDB.ReplicaSet
	Configuration.LoyaltyMongoDB.UserName = Configuration.MongoDB.UserName
	Configuration.LoyaltyMongoDB.Password = Configuration.MongoDB.Password
	Configuration.LoyaltyMongoDB.HostIP_1 = Configuration.MongoDB.HostIP_1
	Configuration.LoyaltyMongoDB.HostPort_1 = Configuration.MongoDB.HostPort_1
	Configuration.LoyaltyMongoDB.HostIP_2 = Configuration.MongoDB.HostIP_2
	Configuration.LoyaltyMongoDB.HostPort_2 = Configuration.MongoDB.HostPort_2
	Configuration.LoyaltyMongoDB.HostIP_3 = Configuration.MongoDB.HostIP_3
	Configuration.LoyaltyMongoDB.HostPort_3 = Configuration.MongoDB.HostPort_3
	Configuration.LoyaltyMongoDB.HostIP_4 = Configuration.MongoDB.HostIP_4
	Configuration.LoyaltyMongoDB.HostPort_4 = Configuration.MongoDB.HostPort_4

	Configuration.IN.IP = "10.10.51.51"
	Configuration.IN.Port = "8080"
	Configuration.IN.WS_SOAP_Endpoint = "/axis2/services/WebService.WebServiceHttpSoap12Endpoint/"
	Configuration.IN.WS_XMLNS_SOAP_Env = "http://schemas.xmlsoap.org/soap/envelope/"
	Configuration.IN.WS_XMLNS_Web = "http://webservice.CSI.omvia.convergys.com"

	Configuration.IN.WS_EVC_SOAP_Endpoint = "/axis2/services/ERechargeWebService.ERechargeWebServiceHttpSoap11Endpoint/"
	Configuration.IN.WS_EVC_XMLNS_SOAP_Env = "http://schemas.xmlsoap.org/soap/envelope/"
	Configuration.IN.WS_EVC_XMLNS_Web = "http://webservice.CSI.omvia.convergys.com"

	Configuration.IN.Default_OpId = "lendme"
	Configuration.IN.Default_OpPwd = "gu1vY6Q$"
	Configuration.IN.Is_OpPwd_Required = true
	Configuration.IN.Timeout = 5
	Configuration.IN.PrintLogs = true

	//http://10.95.64.6:15403/?systemid=lendme&password=lendmeP@ssw0rd&Originator=setest&dest_addr=243900100606&msg_text=test&registered_delivery=0&ston=5&snpi=0&dton=1&dnpi=1&encoding=1

	//SMPP
	Configuration.SMPP.IP = "10.10.215.52"
	Configuration.SMPP.Port = "15403"

	//Configuration.SMPP.Login = "lendme"
	//Configuration.SMPP.Password = "lendmeP@ssw0rd"
	Configuration.SMPP.Login = "Loyalty"
	Configuration.SMPP.Password = "Loyalty123"
	Configuration.SMPP.TimeOut = 5 //in seconds
	Configuration.SMPP.PrintLogs = true
	Configuration.SMPP.MSISDN_Short_len = 8
	Configuration.SMPP.CountryCodePrefix = "232"
	Configuration.SMPP.DefaultSender = "Africell" //"Africell"
	Configuration.SMPP.Encoding = 0
	//CGW
	Configuration.CGW.Protocol = "http"
	Configuration.CGW.Hostname = "UCGW"
	Configuration.CGW.Port = "9991"
	Configuration.CGW.Module = "UCGW"
	Configuration.CGW.Version = "V1"
	Configuration.CGW.Timeout = 15 * time.Second
	Configuration.CGW.S2S_AccessToken = ""

	Configuration.CGW_AUC.Description = "UCGW AUC service"
	Configuration.CGW_AUC.Protocol = "http"
	Configuration.CGW_AUC.Hostname = "UCGW_AUC"
	Configuration.CGW_AUC.Port = "9001"
	Configuration.CGW_AUC.Module = "AUC"
	Configuration.CGW_AUC.Version = "V1"
	Configuration.CGW_AUC.S2S_Username = "UCGW_Admin"
	Configuration.CGW_AUC.S2S_Password = "uC@g$W$iRiS6$2"
	Configuration.CGW_AUC.Timeout_After = 5 * time.Second

	Configuration.Propylaea.Description = "Product Design Center - Propylaea"
	Configuration.Propylaea.Protocol = "http"
	Configuration.Propylaea.Port = "9900"
	Configuration.Propylaea.Hostname = "Propylaea_Prod"
	Configuration.Propylaea.Module = "Propylaea"
	Configuration.Propylaea.Version = "V1"
	Configuration.Propylaea.S2S_Username = "Propylaea_Admin"
	Configuration.Propylaea.S2S_Password = "uC@g$W$iRiS6$2@333dd"
	Configuration.Propylaea.Timeout_After = 5 * time.Second

	return
}

func setDefaultConfiguration_AO_Loyalty() (Configuration ConfigType) {
	//Configuration.HttpOKAPIServicePort = "9291"
	Configuration.HttpAppServicePort = "9290"           //lendme services
	Configuration.HttpAppLoyaltyServicePort = "9280"    //for USSD and Mobile App
	Configuration.HttpAppLoyaltyManagementPort = "9281" //for OKAPI
	Configuration.HttpAppLoyaltyFeedPort = "9282"       //for IN & MM live feed

	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:3000")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:5173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4414")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://okpaihruat.africell.ao")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://okpaihruat.africell.ao")

	Configuration.Operation = "Angola"
	Configuration.HostId = "Lendme-01"
	Configuration.DB_Name = "Lendme_DB"

	Configuration.Version = "V1"
	Configuration.Module = "Lendme"

	Configuration.LoyaltyVersion = "V1"
	Configuration.LoyaltyModule = "Loyalty"

	Configuration.MSISDN_Prefix = ""
	Configuration.MSISDN_Short_len = 8
	Configuration.CountryCode = "244"

	Configuration.Lendme_EVC_Dealer_MSISDN = "244951010532"
	Configuration.Lendme_EVC_Dealer_PIN = "8236"

	Configuration.IsProduction = false
	Configuration.IsLoyaltyProduction = false
	Configuration.Min_Allowed_Amnt = 50
	Configuration.Service_FeePerc = 0.15
	Configuration.Min_Allowed_AON = 3
	Configuration.Min_Avg3MRecharge = 200
	Configuration.Min_LastRechargePeriod = 60
	Configuration.Min_Allowed_Balance = 0
	Configuration.Max_Allowed_Balance = 100
	Configuration.ARPU_File_Path = "/home/Subs_ARPU/"
	Configuration.ARPU_File_Prefix = "lendme_"
	Configuration.ARPU_File_Column_Separator = ";"

	Configuration.App_AUC.Description = "App AUC service"
	Configuration.App_AUC.Protocol = "http"
	Configuration.App_AUC.Hostname = "Lendme_auc"
	Configuration.App_AUC.Port = "9293"
	Configuration.App_AUC.Module = "AUC"
	Configuration.App_AUC.Version = "V1"
	Configuration.App_AUC.S2S_Username = "Lendme_Admin"
	Configuration.App_AUC.S2S_Password = "s@l$e$IrSW0$4"
	Configuration.App_AUC.Timeout_After = 5 * time.Second

	Configuration.OKAPI_AUC.Description = "OKAPI AUC service"
	Configuration.OKAPI_AUC.Protocol = "http"
	Configuration.OKAPI_AUC.Hostname = "10.250.3.148"
	Configuration.OKAPI_AUC.Port = "9001"
	Configuration.OKAPI_AUC.Module = "AUC"
	Configuration.OKAPI_AUC.Version = "V1"
	Configuration.OKAPI_AUC.S2S_Username = "Loyalty_OKAPI"
	Configuration.OKAPI_AUC.S2S_Password = "]W8#x3D1USKUyH@p]s&D_"
	Configuration.OKAPI_AUC.Timeout_After = 5 * time.Second

	// //mongoDB
	Configuration.MongoDB.ReplicaSet = "reps0"
	Configuration.MongoDB.UserName = "mongo-root"
	Configuration.MongoDB.Password = "Speci@LM0nG0P@ssw0rd_F0r_@ng0l@LenMeRepl##$$"
	Configuration.MongoDB.HostIP_1 = "10.250.1.198" //==>Primary
	Configuration.MongoDB.HostPort_1 = "9001"
	Configuration.MongoDB.HostIP_2 = "10.250.1.199" //==>Secondary
	Configuration.MongoDB.HostPort_2 = "9002"
	Configuration.MongoDB.HostIP_3 = "10.250.1.197" //==> Aribter
	Configuration.MongoDB.HostPort_3 = "9003"
	Configuration.MongoDB.HostIP_4 = ""
	Configuration.MongoDB.HostPort_4 = ""
	////

	Configuration.DB_Name_Loyalty = "Loyalty_DB"
	Configuration.LoyaltyMongoDB.ReplicaSet = Configuration.MongoDB.ReplicaSet
	Configuration.LoyaltyMongoDB.UserName = Configuration.MongoDB.UserName
	Configuration.LoyaltyMongoDB.Password = Configuration.MongoDB.Password
	Configuration.LoyaltyMongoDB.HostIP_1 = Configuration.MongoDB.HostIP_1
	Configuration.LoyaltyMongoDB.HostPort_1 = Configuration.MongoDB.HostPort_1
	Configuration.LoyaltyMongoDB.HostIP_2 = Configuration.MongoDB.HostIP_2
	Configuration.LoyaltyMongoDB.HostPort_2 = Configuration.MongoDB.HostPort_2
	Configuration.LoyaltyMongoDB.HostIP_3 = Configuration.MongoDB.HostIP_3
	Configuration.LoyaltyMongoDB.HostPort_3 = Configuration.MongoDB.HostPort_3
	Configuration.LoyaltyMongoDB.HostIP_4 = Configuration.MongoDB.HostIP_4
	Configuration.LoyaltyMongoDB.HostPort_4 = Configuration.MongoDB.HostPort_4

	Configuration.IN.IP = "" // "10.10.51.51"
	Configuration.IN.Port = "8080"
	Configuration.IN.WS_SOAP_Endpoint = "/axis2/services/WebService.WebServiceHttpSoap12Endpoint/"
	Configuration.IN.WS_XMLNS_SOAP_Env = "http://schemas.xmlsoap.org/soap/envelope/"
	Configuration.IN.WS_XMLNS_Web = "http://webservice.CSI.omvia.convergys.com"

	Configuration.IN.WS_EVC_SOAP_Endpoint = "/axis2/services/ERechargeWebService.ERechargeWebServiceHttpSoap11Endpoint/"
	Configuration.IN.WS_EVC_XMLNS_SOAP_Env = "http://schemas.xmlsoap.org/soap/envelope/"
	Configuration.IN.WS_EVC_XMLNS_Web = "http://webservice.CSI.omvia.convergys.com"

	Configuration.IN.Default_OpId = "lendme"
	Configuration.IN.Default_OpPwd = "gu1vY6Q$"
	Configuration.IN.Is_OpPwd_Required = true
	Configuration.IN.Timeout = 5
	Configuration.IN.PrintLogs = true

	//http://10.95.64.6:15403/?systemid=lendme&password=lendmeP@ssw0rd&Originator=setest&dest_addr=243900100606&msg_text=test&registered_delivery=0&ston=5&snpi=0&dton=1&dnpi=1&encoding=1

	//SMPP
	Configuration.SMPP.IP = "10.250.11.111" //"10.10.215.52"
	Configuration.SMPP.Port = "15403"       //"15403"

	//Configuration.SMPP.Login = "lendme"
	//Configuration.SMPP.Password = "lendmeP@ssw0rd"
	Configuration.SMPP.Login = "lendme"
	Configuration.SMPP.Password = "Lm*!op25"
	Configuration.SMPP.TimeOut = 5 //in seconds
	Configuration.SMPP.PrintLogs = true
	Configuration.SMPP.MSISDN_Short_len = 8
	Configuration.SMPP.CountryCodePrefix = "244"
	Configuration.SMPP.DefaultSender = "Africell" //"Africell"
	Configuration.SMPP.Encoding = 0
	//CGW
	Configuration.CGW.Protocol = "http"
	Configuration.CGW.Hostname = "" //"10.10.231.51"
	Configuration.CGW.Port = "9991"
	Configuration.CGW.Module = "UCGW"
	Configuration.CGW.Version = "V1"
	Configuration.CGW.Timeout = 15 * time.Second
	Configuration.CGW.S2S_AccessToken = ""

	Configuration.CGW_AUC.Description = "UCGW AUC service"
	Configuration.CGW_AUC.Protocol = "http"
	Configuration.CGW_AUC.Hostname = "" // "10.10.231.51"
	Configuration.CGW_AUC.Port = "9994"
	Configuration.CGW_AUC.Module = "AUC"
	Configuration.CGW_AUC.Version = "V1"
	Configuration.CGW_AUC.S2S_Username = "UCGW_Admin"
	Configuration.CGW_AUC.S2S_Password = "uC@g$W$iRiS6$2"
	Configuration.CGW_AUC.Timeout_After = 5 * time.Second

	Configuration.Propylaea.Description = "Product Design Center - Propylaea"
	Configuration.Propylaea.Protocol = "http"
	Configuration.Propylaea.Port = "9900"
	Configuration.Propylaea.Hostname = "10.250.1.228"
	Configuration.Propylaea.Module = "Propylaea"
	Configuration.Propylaea.Version = "V1"
	Configuration.Propylaea.S2S_Username = "Propylaea_Admin"
	Configuration.Propylaea.S2S_Password = "uC@g$W$iRiS6$2@333dd"
	Configuration.Propylaea.Timeout_After = 5 * time.Second

	Configuration.APGW.Description = "APGW service"
	Configuration.APGW.Protocol = "http"
	Configuration.APGW.Hostname = "10.250.1.228"
	Configuration.APGW.Port = "9903"
	Configuration.APGW.S2S_Username = "LendmeClient"
	Configuration.APGW.S2S_Password = "L3ndM320#25)sdfj&^M@y"
	Configuration.APGW.Timeout_After = 5 * time.Second
	return
}

func setDefaultConfiguration_AO_Loyalty_UAT() (Configuration ConfigType) {
	//Configuration.HttpOKAPIServicePort = "9291"
	Configuration.HttpAppServicePort = "9290"           //lendme services
	Configuration.HttpAppLoyaltyServicePort = "9280"    //for USSD and Mobile App
	Configuration.HttpAppLoyaltyManagementPort = "9281" //for OKAPI
	Configuration.HttpAppLoyaltyFeedPort = "9282"       //for IN & MM live feed

	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:3000")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:5173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4414")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://okapihruat.africell.ao")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://okapihruat.africell.ao")

	Configuration.Operation = "Angola"
	Configuration.HostId = "Lendme-01"
	Configuration.DB_Name = "Lendme_DB"

	Configuration.Version = "V1"
	Configuration.Module = "Lendme"

	Configuration.LoyaltyVersion = "V1"
	Configuration.LoyaltyModule = "Loyalty"

	Configuration.MSISDN_Prefix = ""
	Configuration.MSISDN_Short_len = 8
	Configuration.CountryCode = "244"

	Configuration.Lendme_EVC_Dealer_MSISDN = "244951010534"
	Configuration.Lendme_EVC_Dealer_PIN = "8236"

	Configuration.IsProduction = false
	Configuration.IsLoyaltyProduction = false
	Configuration.Min_Allowed_Amnt = 50
	Configuration.Service_FeePerc = 0.12
	Configuration.Min_Allowed_AON = 3
	Configuration.Min_Avg3MRecharge = 200
	Configuration.Min_LastRechargePeriod = 60
	Configuration.Min_Allowed_Balance = 0
	Configuration.Max_Allowed_Balance = 50
	Configuration.ARPU_File_Path = "/home/Subs_ARPU/"
	Configuration.ARPU_File_Prefix = "lendme_"
	Configuration.ARPU_File_Column_Separator = ";"

	Configuration.App_AUC.Description = "App AUC service"
	Configuration.App_AUC.Protocol = "http"
	Configuration.App_AUC.Hostname = "Lendme_auc"
	Configuration.App_AUC.Port = "9293"
	Configuration.App_AUC.Module = "AUC"
	Configuration.App_AUC.Version = "V1"
	Configuration.App_AUC.S2S_Username = "Lendme_Admin"
	Configuration.App_AUC.S2S_Password = "s@l$e$IrSW0$4"
	Configuration.App_AUC.Timeout_After = 5 * time.Second

	Configuration.OKAPI_AUC.Description = "OKAPI AUC service"
	Configuration.OKAPI_AUC.Protocol = "http"
	Configuration.OKAPI_AUC.Hostname = "10.250.3.149"
	Configuration.OKAPI_AUC.Port = "9001"
	Configuration.OKAPI_AUC.Module = "AUC"
	Configuration.OKAPI_AUC.Version = "V1"
	Configuration.OKAPI_AUC.S2S_Username = "Loyalty_OKAPI"
	Configuration.OKAPI_AUC.S2S_Password = "]W8#x3D1USKUyH@p]s&D_"
	Configuration.OKAPI_AUC.Timeout_After = 5 * time.Second

	// //mongoDB
	Configuration.MongoDB.ReplicaSet = "reps0"
	Configuration.MongoDB.UserName = "mongo-root"
	Configuration.MongoDB.Password = "Speci@LM0nG0P@ssw0rd_F0r_@ng0l@LenMeRepl##$$"
	Configuration.MongoDB.HostIP_1 = "10.250.1.198" //==>Primary
	Configuration.MongoDB.HostPort_1 = "9001"
	Configuration.MongoDB.HostIP_2 = "10.250.1.199" //==>Secondary
	Configuration.MongoDB.HostPort_2 = "9002"
	Configuration.MongoDB.HostIP_3 = "10.250.1.197" //==> Aribter
	Configuration.MongoDB.HostPort_3 = "9003"
	Configuration.MongoDB.HostIP_4 = ""
	Configuration.MongoDB.HostPort_4 = ""
	////

	Configuration.DB_Name_Loyalty = "Loyalty_DB"
	Configuration.LoyaltyMongoDB.ReplicaSet = Configuration.MongoDB.ReplicaSet
	Configuration.LoyaltyMongoDB.UserName = Configuration.MongoDB.UserName
	Configuration.LoyaltyMongoDB.Password = Configuration.MongoDB.Password
	Configuration.LoyaltyMongoDB.HostIP_1 = Configuration.MongoDB.HostIP_1
	Configuration.LoyaltyMongoDB.HostPort_1 = Configuration.MongoDB.HostPort_1
	Configuration.LoyaltyMongoDB.HostIP_2 = Configuration.MongoDB.HostIP_2
	Configuration.LoyaltyMongoDB.HostPort_2 = Configuration.MongoDB.HostPort_2
	Configuration.LoyaltyMongoDB.HostIP_3 = Configuration.MongoDB.HostIP_3
	Configuration.LoyaltyMongoDB.HostPort_3 = Configuration.MongoDB.HostPort_3
	Configuration.LoyaltyMongoDB.HostIP_4 = Configuration.MongoDB.HostIP_4
	Configuration.LoyaltyMongoDB.HostPort_4 = Configuration.MongoDB.HostPort_4

	Configuration.IN.IP = "" // "10.10.51.51"
	Configuration.IN.Port = "8080"
	Configuration.IN.WS_SOAP_Endpoint = "/axis2/services/WebService.WebServiceHttpSoap12Endpoint/"
	Configuration.IN.WS_XMLNS_SOAP_Env = "http://schemas.xmlsoap.org/soap/envelope/"
	Configuration.IN.WS_XMLNS_Web = "http://webservice.CSI.omvia.convergys.com"

	Configuration.IN.WS_EVC_SOAP_Endpoint = "/axis2/services/ERechargeWebService.ERechargeWebServiceHttpSoap11Endpoint/"
	Configuration.IN.WS_EVC_XMLNS_SOAP_Env = "http://schemas.xmlsoap.org/soap/envelope/"
	Configuration.IN.WS_EVC_XMLNS_Web = "http://webservice.CSI.omvia.convergys.com"

	Configuration.IN.Default_OpId = "lendme"
	Configuration.IN.Default_OpPwd = "gu1vY6Q$"
	Configuration.IN.Is_OpPwd_Required = true
	Configuration.IN.Timeout = 5
	Configuration.IN.PrintLogs = true

	//http://10.95.64.6:15403/?systemid=lendme&password=lendmeP@ssw0rd&Originator=setest&dest_addr=243900100606&msg_text=test&registered_delivery=0&ston=5&snpi=0&dton=1&dnpi=1&encoding=1

	//SMPP
	Configuration.SMPP.IP = "10.250.11.111" //"10.10.215.52"
	Configuration.SMPP.Port = "15403"       //"15403"

	//Configuration.SMPP.Login = "lendme"
	//Configuration.SMPP.Password = "lendmeP@ssw0rd"
	Configuration.SMPP.Login = "lendme"
	Configuration.SMPP.Password = "Lm*!op25"
	Configuration.SMPP.TimeOut = 5 //in seconds
	Configuration.SMPP.PrintLogs = true
	Configuration.SMPP.MSISDN_Short_len = 8
	Configuration.SMPP.CountryCodePrefix = "244"
	Configuration.SMPP.DefaultSender = "Africell" //"Africell"
	Configuration.SMPP.Encoding = 0
	//CGW
	Configuration.CGW.Protocol = "http"
	Configuration.CGW.Hostname = "" //"10.10.231.51"
	Configuration.CGW.Port = "9991"
	Configuration.CGW.Module = "UCGW"
	Configuration.CGW.Version = "V1"
	Configuration.CGW.Timeout = 15 * time.Second
	Configuration.CGW.S2S_AccessToken = ""

	Configuration.CGW_AUC.Description = "UCGW AUC service"
	Configuration.CGW_AUC.Protocol = "http"
	Configuration.CGW_AUC.Hostname = "" // "10.10.231.51"
	Configuration.CGW_AUC.Port = "9994"
	Configuration.CGW_AUC.Module = "AUC"
	Configuration.CGW_AUC.Version = "V1"
	Configuration.CGW_AUC.S2S_Username = "UCGW_Admin"
	Configuration.CGW_AUC.S2S_Password = "uC@g$W$iRiS6$2"
	Configuration.CGW_AUC.Timeout_After = 5 * time.Second

	Configuration.Propylaea.Description = "Product Design Center - Propylaea"
	Configuration.Propylaea.Protocol = "http"
	Configuration.Propylaea.Port = "9900"
	Configuration.Propylaea.Hostname = "10.250.1.228"
	Configuration.Propylaea.Module = "Propylaea"
	Configuration.Propylaea.Version = "V1"
	Configuration.Propylaea.S2S_Username = "Propylaea_Admin"
	Configuration.Propylaea.S2S_Password = "uC@g$W$iRiS6$2@333dd"
	Configuration.Propylaea.Timeout_After = 5 * time.Second

	Configuration.APGW.Description = "APGW service"
	Configuration.APGW.Protocol = "http"
	Configuration.APGW.Hostname = "10.250.1.228"
	Configuration.APGW.Port = "9904"
	Configuration.APGW.S2S_Username = "LendmeClient"
	Configuration.APGW.S2S_Password = "L3ndM320#25)sdfj&^M@y"
	Configuration.APGW.Timeout_After = 5 * time.Second
	return
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Encryption functions
// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func EcryptToHexString(str string) (EncHexString string, err error) {
	enc, err := encrypt([]byte(str), EncryptionKey)
	if err != nil {
		return
	}
	EncHexString = hex.EncodeToString(enc)
	return
}

func DecryptHexString(hexStr string) (DecString string, err error) {
	dec, err := hex.DecodeString(hexStr)
	if err != nil {
		return
	}
	dec_c, err := decrypt([]byte(dec), EncryptionKey)
	if err != nil {
		return
	}
	DecString = string(dec_c)
	return
}

func encrypt(data []byte, passphrase string) (ciphertext []byte, err error) {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}
	ciphertext = gcm.Seal(nonce, nonce, data, nil)
	return
}

func decrypt(data []byte, passphrase string) (plaintext []byte, err error) {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err = gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return
	}
	return
}
func setDefaultConfiguration_Dev() (Configuration ConfigType) {
	//Configuration.HttpOKAPIServicePort = "9291"
	Configuration.HttpAppServicePort = "9290"           //lendme services
	Configuration.HttpAppLoyaltyServicePort = "9280"    //for USSD and Mobile App
	Configuration.HttpAppLoyaltyManagementPort = "9281" //for OKAPI
	Configuration.HttpAppLoyaltyFeedPort = "9282"       //for IN & MM live feed

	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:3000")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:5173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4414")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://okpaihr.africell.sl")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://okpaihruat.africell.sl")

	Configuration.Operation = "GM"
	Configuration.HostId = "Lendme-01"
	Configuration.DB_Name = "Lendme_DB"

	Configuration.Version = "V1"
	Configuration.Module = "Lendme"

	Configuration.LoyaltyVersion = "V1"
	Configuration.LoyaltyModule = "Loyalty"

	Configuration.MSISDN_Prefix = "0"
	Configuration.MSISDN_Short_len = 8
	Configuration.CountryCode = "220"

	Configuration.IsProduction = false
	Configuration.IsLoyaltyProduction = false
	Configuration.Min_Allowed_Amnt = 1
	Configuration.Service_FeePerc = 0.15
	Configuration.Min_Allowed_AON = 3
	Configuration.Min_Avg3MRecharge = 0
	Configuration.Min_LastRechargePeriod = 60
	Configuration.Min_Allowed_Balance = 0
	Configuration.Max_Allowed_Balance = 225
	Configuration.ARPU_File_Path = "/home/Subs_ARPU/"

	Configuration.App_AUC.Description = "App AUC service"
	Configuration.App_AUC.Protocol = "http"
	Configuration.App_AUC.Hostname = "localhost"
	Configuration.App_AUC.Port = "9293"
	Configuration.App_AUC.Module = "AUC"
	Configuration.App_AUC.Version = "V1"
	Configuration.App_AUC.S2S_Username = "Lendme_Admin"
	Configuration.App_AUC.S2S_Password = "s@l$e$IrSW0$4"
	Configuration.App_AUC.Timeout_After = 5 * time.Second

	Configuration.OKAPI_AUC.Description = "OKAPI AUC service"
	Configuration.OKAPI_AUC.Protocol = "http"
	Configuration.OKAPI_AUC.Hostname = "localhost"
	Configuration.OKAPI_AUC.Port = "9001"
	Configuration.OKAPI_AUC.Module = "AUC"
	Configuration.OKAPI_AUC.Version = "V1"
	Configuration.OKAPI_AUC.S2S_Username = "Loyalty_OKAPI"
	Configuration.OKAPI_AUC.S2S_Password = "]W8#x3D1USKUyH@p]s&D_"
	Configuration.OKAPI_AUC.Timeout_After = 5 * time.Second

	// //mongoDB
	// Configuration.MongoDB.ReplicaSet = "reps1"
	// Configuration.MongoDB.UserName = "mongo-root"
	// Configuration.MongoDB.Password = "Speci@LM0nG0P@ssw0rd_F0r_LeNdM#SL"
	// Configuration.MongoDB.HostIP_1 = "10.10.247.21" //==>Primary
	// Configuration.MongoDB.HostPort_1 = "9001"
	// Configuration.MongoDB.HostIP_2 = "10.10.247.22" //==>Secondary
	// Configuration.MongoDB.HostPort_2 = "9002"
	// Configuration.MongoDB.HostIP_3 = "10.10.231.52" //==> Aribter
	// Configuration.MongoDB.HostPort_3 = "9003"
	// Configuration.MongoDB.HostIP_4 = ""
	// Configuration.MongoDB.HostPort_4 = ""
	////
	// Configuration.MongoDB.ReplicaSet = ""
	// Configuration.MongoDB.UserName = "mongo-root"
	// Configuration.MongoDB.Password = "Speci@LM0nG0P@ssw0rd_F0r_LeNdM#SL"
	// Configuration.MongoDB.HostIP_1 = "10.10.247.21" //==>Primary
	// Configuration.MongoDB.HostPort_1 = "9001"
	// Configuration.MongoDB.HostIP_2 = "" //==>Secondary
	// Configuration.MongoDB.HostPort_2 = ""
	// Configuration.MongoDB.HostIP_3 = "" //==> Aribter
	// Configuration.MongoDB.HostPort_3 = ""
	// Configuration.MongoDB.HostIP_4 = ""
	// Configuration.MongoDB.HostPort_4 = ""

	//mongoDB
	Configuration.MongoDB.UserName = ""
	Configuration.MongoDB.Password = ""
	Configuration.MongoDB.HostIP_1 = "localhost" //"host.docker.internal"
	Configuration.MongoDB.HostPort_1 = "27017"   //"27017"
	///////////////////////////////////////MONGO DOCKER ///////////////////////////

	Configuration.DB_Name_Loyalty = "Loyalty_DB"
	Configuration.LoyaltyMongoDB.ReplicaSet = Configuration.MongoDB.ReplicaSet
	Configuration.LoyaltyMongoDB.UserName = Configuration.MongoDB.UserName
	Configuration.LoyaltyMongoDB.Password = Configuration.MongoDB.Password
	Configuration.LoyaltyMongoDB.HostIP_1 = Configuration.MongoDB.HostIP_1
	Configuration.LoyaltyMongoDB.HostPort_1 = Configuration.MongoDB.HostPort_1
	Configuration.LoyaltyMongoDB.HostIP_2 = Configuration.MongoDB.HostIP_2
	Configuration.LoyaltyMongoDB.HostPort_2 = Configuration.MongoDB.HostPort_2
	Configuration.LoyaltyMongoDB.HostIP_3 = Configuration.MongoDB.HostIP_3
	Configuration.LoyaltyMongoDB.HostPort_3 = Configuration.MongoDB.HostPort_3
	Configuration.LoyaltyMongoDB.HostIP_4 = Configuration.MongoDB.HostIP_4
	Configuration.LoyaltyMongoDB.HostPort_4 = Configuration.MongoDB.HostPort_4

	Configuration.IN.IP = "10.95.73.12" //"10.70.1.59"
	Configuration.IN.Port = "8444"

	Configuration.IN.WS_SOAP_Endpoint = "/axis2/services/WebService.WebServiceHttpSoap12Endpoint/"
	Configuration.IN.WS_XMLNS_SOAP_Env = "http://schemas.xmlsoap.org/soap/envelope/"
	Configuration.IN.WS_XMLNS_Web = "http://webservice.CSI.omvia.convergys.com"

	Configuration.IN.WS_EVC_SOAP_Endpoint = "/axis2/services/ERechargeWebService.ERechargeWebServiceHttpSoap11Endpoint/"
	Configuration.IN.WS_EVC_XMLNS_SOAP_Env = "http://schemas.xmlsoap.org/soap/envelope/"
	Configuration.IN.WS_EVC_XMLNS_Web = "http://webservice.CSI.omvia.convergys.com"

	Configuration.IN.Default_OpId = "lendme"
	Configuration.IN.Default_OpPwd = "P@ssw0rd"
	Configuration.IN.Is_OpPwd_Required = true
	Configuration.IN.Timeout = 5
	Configuration.IN.PrintLogs = true

	//http://10.95.64.6:15403/?systemid=lendme&password=lendmeP@ssw0rd&Originator=setest&dest_addr=243900100606&msg_text=test&registered_delivery=0&ston=5&snpi=0&dton=1&dnpi=1&encoding=1

	//SMPP
	Configuration.SMPP.IP = "10.10.215.52"
	Configuration.SMPP.Port = "15403"

	//Configuration.SMPP.Login = "lendme"
	//Configuration.SMPP.Password = "lendmeP@ssw0rd"
	Configuration.SMPP.Login = "LendME2"
	Configuration.SMPP.Password = "LendMEP@ssw0rd"

	//CGW
	Configuration.CGW.Protocol = "http"
	Configuration.CGW.Hostname = "UCGW"
	Configuration.CGW.Port = "9991"
	Configuration.CGW.Module = "UCGW"
	Configuration.CGW.Version = "V1"
	Configuration.CGW.Timeout = 15 * time.Second
	Configuration.CGW.S2S_AccessToken = ""

	Configuration.CGW_AUC.Description = "UCGW AUC service"
	Configuration.CGW_AUC.Protocol = "http"
	Configuration.CGW_AUC.Hostname = "UCGW_AUC"
	Configuration.CGW_AUC.Port = "9001"
	Configuration.CGW_AUC.Module = "AUC"
	Configuration.CGW_AUC.Version = "V1"
	Configuration.CGW_AUC.S2S_Username = "UCGW_Admin"
	Configuration.CGW_AUC.S2S_Password = "uC@g$W$iRiS6$2"
	Configuration.CGW_AUC.Timeout_After = 5 * time.Second

	Configuration.Propylaea.Description = "Product Design Center - Propylaea"
	Configuration.Propylaea.Protocol = "https"
	Configuration.Propylaea.Port = "443"
	Configuration.Propylaea.Hostname = "sapp.africell.cd"
	Configuration.Propylaea.Module = "Propylaea"
	Configuration.Propylaea.Version = "V1"
	Configuration.Propylaea.S2S_Username = "Propylaea_Admin"
	Configuration.Propylaea.S2S_Password = "uC@g$W$iRiS6$2@333dd"
	Configuration.Propylaea.Timeout_After = 5 * time.Second
	Configuration.Propylaea.ChannelName = "Spin And Win"
	Configuration.Propylaea.ChannelPlan = "Normal SIM"
	Configuration.Propylaea.ChannelVersion = "1"

	Configuration.SMPP.IP = "10.95.64.6"
	Configuration.SMPP.Port = "15403"
	Configuration.SMPP.Login = "OKAPI"
	Configuration.SMPP.Password = "OKAPIP@ssw0rd"
	Configuration.SMPP.TimeOut = 5 //in seconds
	Configuration.SMPP.PrintLogs = true
	Configuration.SMPP.MSISDN_Short_len = 9
	Configuration.SMPP.CountryCodePrefix = "243"
	Configuration.SMPP.DefaultSender = "Africell" //"Africell"
	Configuration.SMPP.Encoding = 1

	Configuration.ISLoyaltyOptOutGracePeriodDays = 30
	// Configuration.KafkaBrokerUrls = "kafka1:9092,kafka2:9092,kafka3:9092"
	// Configuration.KafkaBrokerUrls = "kafka3:9092,kafka2:9092,kafka1:9092"
	// Configuration.KafkaClientId = "LoyaltyLiveFeed"

	return
}
