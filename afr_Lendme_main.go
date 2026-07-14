package main

import (
	lendme "afr_Loyalty/Lendme"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kardianos/service"
	"github.com/rs/cors"
)

const ApplicationName = "AFR Lendme Backend"
const ApplicationReleaseNumber = "0.1.0"
const ApplicationReleaseDate = "23/01/2024" //"dd/MM/YYYY"

var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	log.Println("-------------------------------------------------------------------")
	log.Println("Application name: ", ApplicationName)
	log.Println("Application release number: ", ApplicationReleaseNumber)
	log.Println("Application release date: ", ApplicationReleaseDate)
	log.Println("-------------------------------------------------------------------")

	//read and parse configuration
	err := lendme.GetDefaultConfiguration()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Establishing connections...")
	UserControl := lendme.NewUserControl()

	//UserControl.Init_Notification_Profiles()

	//Prometheus functions
	lendme.Init_Prometheus_Metrics()
	lendme.Init_Prometheus_Metrics_Latency()
	go lendme.PortlinkInquiry()
	go lendme.Reset_Prometheus_Metrics()
	go lendme.Reset_Prometheus_Metrics_Latency()

	corsOpts := cors.New(cors.Options{
		AllowedOrigins: lendme.Configuration.OKAPIAllowedOrigins, //you service is available and allowed for this base url
		AllowedMethods: []string{
			http.MethodGet, //http methods for your app
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowedHeaders: []string{
			"*", //or you can your header key values which you are using in your application

		},
	})

	//**Loyalty main section
	if lendme.Configuration.HttpAppLoyaltyServicePort != "" {
		//Initialize loyalty
		UserControl.InitializeMongoxRepositories()
		UserControl.RedisDataLoader()
		UserControl.LoyaltyIndexesMaintenanceProcess()
		// UserControl.InitializeLoyaltyDefaultUAT()
		// To be removed in the next deployment
		UserControl.SeedUATLoyaltyTestPointsExpiry()

		//Loyalty service
		log.Println("Add Loyalty service routers to the web service")
		Loyalty_Service_router := mux.NewRouter().StrictSlash(true)
		UserControl.AddToLoyaltyServiceRouter(Loyalty_Service_router, UserControl)
		HttpLoyaltyServicePort := lendme.Configuration.HttpAppLoyaltyServicePort
		log.Println("Loyalty service WS listen and serve on port: " + HttpLoyaltyServicePort) //auc.Configuration.HttpServicePort
		go http.ListenAndServe(":"+HttpLoyaltyServicePort, corsOpts.Handler(Loyalty_Service_router))

		//Loyalty management
		log.Println("Add Loyalty management routers to the web service")
		Loyalty_management_router := mux.NewRouter().StrictSlash(true)
		UserControl.AddToLoyaltyManagementRouter(Loyalty_management_router, UserControl)
		HttpLoyaltyManagementPort := lendme.Configuration.HttpAppLoyaltyManagementPort
		log.Println("Loyalty management WS listen and serve on port: " + HttpLoyaltyManagementPort) //auc.Configuration.HttpServicePort
		go http.ListenAndServe(":"+HttpLoyaltyManagementPort, corsOpts.Handler(Loyalty_management_router))

		//Loyalty Feed
		log.Println("Add Loyalty Feed routers to the web service")
		Loyalty_Feed_router := mux.NewRouter().StrictSlash(true)
		UserControl.AddToLoyaltyFeedRouter(Loyalty_Feed_router, UserControl)
		HttpLoyaltyFeedPort := lendme.Configuration.HttpAppLoyaltyFeedPort
		log.Println("Loyalty Feed WS listen and serve on port: " + HttpLoyaltyFeedPort) //auc.Configuration.HttpServicePort
		go http.ListenAndServe(":"+HttpLoyaltyFeedPort, corsOpts.Handler(Loyalty_Feed_router))

		go UserControl.PointsExpiry_Process()
		go UserControl.Loyalty_Governance_DailyLog_Process()
		go UserControl.Loyalty_Status_Expiry_Daily_Process()
		go UserControl.Loyalty_Customer_Account_Daily_Snapshot()
		go UserControl.LoyaltyGovernancePools_Metrics_Process()
		go UserControl.Auto_GetLoyaltySubsSummary()

	}

	go UserControl.SubQueueExecution()
	go UserControl.Auto_Import_Subscribers_Dump()

	//**Lendme web services
	log.Println("Add App routers to the web service")
	App_router := mux.NewRouter().StrictSlash(true)
	UserControl.AddToAppRouter(App_router, UserControl)
	// //**OKAPI web Service
	// log.Println("Add Lendme routers to the web service")
	// OKAPI_router := mux.NewRouter().StrictSlash(true)
	// UserControl.AddToOKAPIRouter(OKAPI_router, UserControl)
	// HttpOKAPIServicePort := lendme.Configuration.HttpOKAPIServicePort
	// log.Println("Lendme HTTP listen and serve on port: " + HttpOKAPIServicePort) //auc.Configuration.HttpServicePort
	// log.Fatal(http.ListenAndServe(":"+HttpOKAPIServicePort, corsOpts.Handler(OKAPI_router)))
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        ApplicationName,
		DisplayName: ApplicationName,
		Description: ApplicationName + " service",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
