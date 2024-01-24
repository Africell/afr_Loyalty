package Lendme

import (
	"time"
)

var Configuration ConfigType

type ConfigType struct {
	HttpOKAPIServicePort string
	HttpAppServicePort   string
	OKAPIAllowedOrigins  []string

	HostId  string
	DB_Name string
	Version string
	Module  string

	IsProduction            bool
	Thumbnail_Width         int
	AgentMaxAllowedDistance float64

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
}

func GetDefaultConfiguration() (err error) {
	// Configuration = setDefaultConfiguration_Dev()
	//setDefaultConfiguration_local_docker()
	// Configuration = setDefaultConfiguration_AO_SalesMonitoring_UAT_229()
	Configuration = setDefaultConfiguration_AO_Live()
	return nil
}

func setDefaultConfiguration_AO_SalesMonitoring_UAT_229() (Configuration ConfigType) {
	Configuration.HttpOKAPIServicePort = "9291"
	Configuration.HttpAppServicePort = "9290"
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:3000")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:5173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4414")

	Configuration.HostId = "Outlet-01"
	Configuration.DB_Name = "Outlet_DB"

	Configuration.Version = "V1"
	Configuration.Module = "Outlet"

	Configuration.IsProduction = false
	Configuration.Thumbnail_Width = 200
	Configuration.AgentMaxAllowedDistance = 0.1

	Configuration.App_AUC.Description = "App AUC service"
	Configuration.App_AUC.Protocol = "http"
	Configuration.App_AUC.Hostname = "SalesMonitoring_UAT_AUC"
	Configuration.App_AUC.Port = "9001"
	Configuration.App_AUC.Module = "AUC"
	Configuration.App_AUC.Version = "V1"
	Configuration.App_AUC.S2S_Username = "SalesMonitoring_Admin"
	Configuration.App_AUC.S2S_Password = "s@le$PaS$W0$9"
	Configuration.App_AUC.Timeout_After = 5 * time.Second

	Configuration.OKAPI_AUC.Description = "OKAPI AUC service"
	Configuration.OKAPI_AUC.Protocol = "http"
	Configuration.OKAPI_AUC.Hostname = "SalesMonitoring_UAT_OKAPI_AUC"
	Configuration.OKAPI_AUC.Port = "9002"
	Configuration.OKAPI_AUC.Module = "AUC"
	Configuration.OKAPI_AUC.Version = "V1"
	Configuration.OKAPI_AUC.S2S_Username = "SalesMonitoring_OKAPI"
	Configuration.OKAPI_AUC.S2S_Password = "s@le$P@s$W0$3"
	Configuration.OKAPI_AUC.Timeout_After = 5 * time.Second

	//mongoDB
	Configuration.MongoDB.ReplicaSet = ""
	Configuration.MongoDB.UserName = "db_root"
	Configuration.MongoDB.Password = "A2023P@sSo0612"
	Configuration.MongoDB.HostIP_1 = "SalesMonitoring_UAT_mongodb"
	Configuration.MongoDB.HostPort_1 = "27017"
	Configuration.MongoDB.HostIP_2 = ""
	Configuration.MongoDB.HostPort_2 = ""
	Configuration.MongoDB.HostIP_3 = ""
	Configuration.MongoDB.HostPort_3 = ""
	Configuration.MongoDB.HostIP_4 = ""
	Configuration.MongoDB.HostPort_4 = ""

	return
}
func setDefaultConfiguration_AO_Live() (Configuration ConfigType) {
	Configuration.HttpOKAPIServicePort = "9291"
	Configuration.HttpAppServicePort = "9290"
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:3000")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:5173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4414")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://okpaihr.africell.ao")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://outlet.africell.ao")

	Configuration.HostId = "Outlet-01"
	Configuration.DB_Name = "Outlet_DB"

	Configuration.Version = "V1"
	Configuration.Module = "Outlet"

	Configuration.IsProduction = false
	Configuration.Thumbnail_Width = 200
	Configuration.AgentMaxAllowedDistance = 0.1

	Configuration.App_AUC.Description = "App AUC service"
	Configuration.App_AUC.Protocol = "https"
	Configuration.App_AUC.Hostname = "outlet_auc"
	Configuration.App_AUC.Port = "9001"
	Configuration.App_AUC.Module = "AUC"
	Configuration.App_AUC.Version = "V1"
	Configuration.App_AUC.S2S_Username = "SalesMonitoring_Admin"
	Configuration.App_AUC.S2S_Password = "s@le$PaS$W0$9"
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
	Configuration.MongoDB.ReplicaSet = ""
	Configuration.MongoDB.UserName = "db_root"
	Configuration.MongoDB.Password = "P@s54D0Brdara_r@75S"
	Configuration.MongoDB.HostIP_1 = "mongodb"
	Configuration.MongoDB.HostPort_1 = "27017"
	Configuration.MongoDB.HostIP_2 = ""
	Configuration.MongoDB.HostPort_2 = ""
	Configuration.MongoDB.HostIP_3 = ""
	Configuration.MongoDB.HostPort_3 = ""
	Configuration.MongoDB.HostIP_4 = ""
	Configuration.MongoDB.HostPort_4 = ""

	return
}

func setDefaultConfiguration_Dev() (Configuration ConfigType) {
	Configuration.HttpOKAPIServicePort = "9291"
	Configuration.HttpAppServicePort = "9290"
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:3000")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:5173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4414")

	Configuration.HostId = "Outlet-01"
	Configuration.DB_Name = "Outlet_DB"

	Configuration.Version = "V1"
	Configuration.Module = "Outlet"

	Configuration.IsProduction = false

	Configuration.App_AUC.Description = "App AUC service"
	Configuration.App_AUC.Protocol = "http"
	Configuration.App_AUC.Hostname = "localhost"
	Configuration.App_AUC.Port = "9002"
	Configuration.App_AUC.Module = "AUC"
	Configuration.App_AUC.Version = "V1"
	Configuration.App_AUC.S2S_Username = "SalesMonitoring_Admin"
	Configuration.App_AUC.S2S_Password = "s@le$PaS$W0$9"
	Configuration.App_AUC.Timeout_After = 5 * time.Second

	Configuration.OKAPI_AUC.Description = "OKAPI AUC service"
	Configuration.OKAPI_AUC.Protocol = "http"
	Configuration.OKAPI_AUC.Hostname = "localhost"
	Configuration.OKAPI_AUC.Port = "9001"
	Configuration.OKAPI_AUC.Module = "AUC"
	Configuration.OKAPI_AUC.Version = "V1"
	Configuration.OKAPI_AUC.S2S_Username = "SalesMonitoring_OKAPI"
	Configuration.OKAPI_AUC.S2S_Password = "s@le$P@s$W0$3"
	Configuration.OKAPI_AUC.Timeout_After = 5 * time.Second

	//mongoDB Dev
	Configuration.MongoDB.UserName = ""
	Configuration.MongoDB.Password = ""
	Configuration.MongoDB.HostIP_1 = "localhost" //"host.docker.internal"
	Configuration.MongoDB.HostPort_1 = "27017"   //"27017"

	return
}

func setDefaultConfiguration_local_docker() (config ConfigType) {
	Configuration.HttpOKAPIServicePort = "9291"
	Configuration.HttpAppServicePort = "9290"
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "https://localhost")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:3000")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:5173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4173")
	Configuration.OKAPIAllowedOrigins = append(Configuration.OKAPIAllowedOrigins, "http://localhost:4414")

	Configuration.HostId = "Outlet-01"
	Configuration.DB_Name = "Outlet_DB"

	Configuration.Version = "V1"
	Configuration.Module = "Outlet"

	Configuration.IsProduction = false

	Configuration.App_AUC.Description = "App AUC service"
	Configuration.App_AUC.Protocol = "http"
	Configuration.App_AUC.Hostname = "outlet_auc"
	Configuration.App_AUC.Port = "9001"
	Configuration.App_AUC.Module = "AUC"
	Configuration.App_AUC.Version = "V1"
	Configuration.App_AUC.S2S_Username = "SalesMonitoring_Admin"
	Configuration.App_AUC.S2S_Password = "s@le$PaS$W0$9"
	Configuration.App_AUC.Timeout_After = 5 * time.Second

	Configuration.OKAPI_AUC.Description = "OKAPI AUC service"
	Configuration.OKAPI_AUC.Protocol = "http"
	Configuration.OKAPI_AUC.Hostname = "auc"
	Configuration.OKAPI_AUC.Port = "9001"
	Configuration.OKAPI_AUC.Module = "AUC"
	Configuration.OKAPI_AUC.Version = "V1"
	Configuration.OKAPI_AUC.S2S_Username = "SalesMonitoring_OKAPI"
	Configuration.OKAPI_AUC.S2S_Password = "s@le$P@s$W0$3"
	Configuration.OKAPI_AUC.Timeout_After = 5 * time.Second

	//mongoDB Dev
	Configuration.MongoDB.UserName = ""
	Configuration.MongoDB.Password = ""
	Configuration.MongoDB.HostIP_1 = "mongodb" //"host.docker.internal"
	Configuration.MongoDB.HostPort_1 = "27017" //"27017"

	return
}
