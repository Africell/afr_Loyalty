package main

import (
	lendme "afr_lendme/Lendme"
	sysadmin "daoc/SysAdmin"
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
	UserControl.InitializeDAO()
	UserControl.InitializeCache()
	UserControl.IndexesMaintenanceProcess()
	UserControl.LoadDefaultValues()

	//UserControl.Init_Notification_Profiles()

	//Prometheus functions
	lendme.Init_Prometheus_Metrics()
	lendme.Init_Prometheus_Metrics_Latency()
	go lendme.PortlinkInquiry()
	go lendme.Reset_Prometheus_Metrics()
	go lendme.Reset_Prometheus_Metrics_Latency()

	go sysadmin.SysAdminInit(UserControl.CacheDir)

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

	//Add user routers to the web service
	//**App web services
	log.Println("Add App routers to the web service")
	App_router := mux.NewRouter().StrictSlash(true)
	UserControl.AddToAppRouter(App_router, UserControl)
	HttpAppServicePort := lendme.Configuration.HttpAppServicePort
	log.Println("App HTTP listen and serve on port: " + HttpAppServicePort) //auc.Configuration.HttpServicePort
	//go http.ListenAndServe(":"+HttpAppServicePort, corsOpts.Handler(App_router))
	log.Fatal(http.ListenAndServe(":"+HttpAppServicePort, corsOpts.Handler(App_router)))
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
