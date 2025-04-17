package Lendme

import (
	"afr_auth_center/AuthCenter"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Route struct {
	Name               string
	Method             string
	Pattern            string
	HandlerFunc        http.HandlerFunc
	AddToAccessEntry   bool
	DisplayName        string //create new
	DisplayOrder       int64
	Module             string //inventory
	ModuleDisplayOrder int64
	Level1             string //transactions
	Level1DisplayOrder int64
	Level2             string //purchase request
	Level2DisplayOrder int64
	Level3             string
	Level3DisplayOrder int64
	AllowedFor_OKAPI   bool
	AllowedFor_App     bool
}

type Routes []Route

func CreateRoutes(UC *UserControl) Routes {
	return Routes{
		// Route{
		// 	"HTTP_ReceiveSSR",
		// 	"POST",
		// 	"/HTTP_ReceiveSSR/",
		// 	UC.HTTP_ReceiveSSR,
		// 	false,
		// 	"Receive SSR", // DisplayName
		// 	1,             // DisplayOrder
		// 	"",            // Module
		// 	0,             //ModuleDisplayOrder
		// 	"",            // Level1
		// 	0,             // Level1DisplayOrder
		// 	"",            // Level2
		// 	0,             // Level2DisplayOrder
		// 	"",            // Level3
		// 	0,             // Level3DisplayOrder
		// },
	}
}

func (Uc *UserControl) AddToAppRouter(router *mux.Router, UC *UserControl) {
	// When StrictSlash is set to true, if the route path is "/path/", accessing "/path" will redirect
	// to the former and vice versa
	routes := CreateRoutes(UC)

	Uc.Add_LendmeRoutes(&routes)

	accessEntries, err := Uc.AppAUC.AUCClient.ReadAccessEntries("")
	if err != nil {
		log.Fatalln("Error Reading Existing Access Entries from AUC !!!")
	}

	log.Println("Read existing routes: #", len(accessEntries.Data))

	for _, acc := range accessEntries.Data {
		MapAccessEntry.Put(acc.Key, acc)
	}

	for _, route := range routes {
		if route.AllowedFor_App {
			//var handler http.Handler
			handler := route.HandlerFunc
			//handler = Logger(handler, route.Name)
			router.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(handler)

			//add to OKAPI AUC
			if route.AddToAccessEntry {
				ok := MapAccessEntry.Check(route.Name + "|" + route.Method)
				if !ok {
					var accessEntry AuthCenter.AccessEntry
					accessEntry.AddUser = "Auto Add"
					accessEntry.AddDate = time.Now().UTC()
					accessEntry.LastModifyUser = "Auto Add"
					accessEntry.LastModifyDate = time.Now().UTC()
					accessEntry.Key = route.Name + "|" + route.Method
					accessEntry.AccessKey = route.Name
					accessEntry.AccessMethod = route.Method
					accessEntry.DisplayName = route.DisplayName
					accessEntry.DisplayOrder = route.DisplayOrder
					accessEntry.Module = route.Module
					accessEntry.ModuleDisplayOrder = route.ModuleDisplayOrder
					accessEntry.Level1 = route.Level1
					accessEntry.Level1DisplayOrder = route.Level1DisplayOrder
					accessEntry.Level2 = route.Level2
					accessEntry.Level2DisplayOrder = route.Level2DisplayOrder
					accessEntry.Level3 = route.Level3
					accessEntry.Level3DisplayOrder = route.Level3DisplayOrder

					AccessKeyDescription := strings.Replace(route.Name, "HTTP_", "", -1)
					AccessKeyDescription = strings.Replace(AccessKeyDescription, "_", " ", -1)
					accessEntry.AccessKeyDescription = AccessKeyDescription
					//Insert into Cache
					_, err := Uc.AppAUC.AUCClient.CreateAccessEntries(accessEntry)
					if err == nil {
						log.Println("Created Access Entry: ", accessEntry.Key)
					} else {
						log.Println("Error creating AccessEntry in AUC !!!")
						json.NewEncoder(os.Stdout).Encode(accessEntry)
					}
				}
			}
		}

	}
	router.Path("/metrics").Handler(CustomPrometheusHandler())
	// router.Path("/metrics_latency").Handler(CustomPrometheusLatencyHandler())
}

func (Uc *UserControl) AddToLoyaltyServiceRouter(router *mux.Router, UC *UserControl) {
	// When StrictSlash is set to true, if the route path is "/path/", accessing "/path" will redirect
	// to the former and vice versa
	routes := CreateRoutes(UC)

	Uc.Add_LoyaltyServiceRoutes(&routes)

	accessEntries, err := Uc.AppAUC.AUCClient.ReadAccessEntries("")
	if err != nil {
		log.Fatalln("Error Reading Existing Access Entries from AUC !!!")
	}

	log.Println("Read existing routes: #", len(accessEntries.Data))

	for _, acc := range accessEntries.Data {
		MapAccessEntry.Put(acc.Key, acc)
	}

	for _, route := range routes {
		if route.AllowedFor_App {
			//var handler http.Handler
			handler := route.HandlerFunc
			//handler = Logger(handler, route.Name)
			router.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(handler)

			//add to OKAPI AUC
			if route.AddToAccessEntry {
				ok := MapAccessEntry.Check(route.Name + "|" + route.Method)
				if !ok {
					var accessEntry AuthCenter.AccessEntry
					accessEntry.AddUser = "Auto Add"
					accessEntry.AddDate = time.Now().UTC()
					accessEntry.LastModifyUser = "Auto Add"
					accessEntry.LastModifyDate = time.Now().UTC()
					accessEntry.Key = route.Name + "|" + route.Method
					accessEntry.AccessKey = route.Name
					accessEntry.AccessMethod = route.Method
					accessEntry.DisplayName = route.DisplayName
					accessEntry.DisplayOrder = route.DisplayOrder
					accessEntry.Module = route.Module
					accessEntry.ModuleDisplayOrder = route.ModuleDisplayOrder
					accessEntry.Level1 = route.Level1
					accessEntry.Level1DisplayOrder = route.Level1DisplayOrder
					accessEntry.Level2 = route.Level2
					accessEntry.Level2DisplayOrder = route.Level2DisplayOrder
					accessEntry.Level3 = route.Level3
					accessEntry.Level3DisplayOrder = route.Level3DisplayOrder

					AccessKeyDescription := strings.Replace(route.Name, "HTTP_", "", -1)
					AccessKeyDescription = strings.Replace(AccessKeyDescription, "_", " ", -1)
					accessEntry.AccessKeyDescription = AccessKeyDescription
					//Insert into Cache
					_, err := Uc.AppAUC.AUCClient.CreateAccessEntries(accessEntry)
					if err == nil {
						log.Println("Created Access Entry: ", accessEntry.Key)
					} else {
						log.Println("Error creating AccessEntry in AUC !!!")
						json.NewEncoder(os.Stdout).Encode(accessEntry)
					}
				}
			}
		}

	}
	router.Path("/Loyalty_metrics").Handler(LoyaltyPrometheusHandler())
	// router.Path("/metrics_latency").Handler(CustomPrometheusLatencyHandler())
	Uc.Create_UCGW_AUCUser()
}

func (Uc *UserControl) AddToLoyaltyManagementRouter(router *mux.Router, UC *UserControl) {
	// When StrictSlash is set to true, if the route path is "/path/", accessing "/path" will redirect
	// to the former and vice versa
	routes := CreateRoutes(UC)

	Uc.Add_LoyaltyManagementRoutes(&routes)

	accessEntries, err := Uc.OKAPIAUC.AUCClient.ReadAccessEntries("")
	if err != nil {
		log.Fatalln("Error Reading Existing Access Entries from AUC !!!")
	}

	log.Println("Read existing routes: #", len(accessEntries.Data))

	for _, acc := range accessEntries.Data {
		MapAccessEntry.Put(acc.Key, acc)
	}

	for _, route := range routes {
		if route.AllowedFor_App {
			//var handler http.Handler
			handler := route.HandlerFunc
			//handler = Logger(handler, route.Name)
			router.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(handler)

			//add to OKAPI AUC
			if route.AddToAccessEntry {
				ok := MapAccessEntry.Check(route.Name + "|" + route.Method)
				if !ok {
					var accessEntry AuthCenter.AccessEntry
					accessEntry.AddUser = "Auto Add"
					accessEntry.AddDate = time.Now().UTC()
					accessEntry.LastModifyUser = "Auto Add"
					accessEntry.LastModifyDate = time.Now().UTC()
					accessEntry.Key = route.Name + "|" + route.Method
					accessEntry.AccessKey = route.Name
					accessEntry.AccessMethod = route.Method
					accessEntry.DisplayName = route.DisplayName
					accessEntry.DisplayOrder = route.DisplayOrder
					accessEntry.Module = route.Module
					accessEntry.ModuleDisplayOrder = route.ModuleDisplayOrder
					accessEntry.Level1 = route.Level1
					accessEntry.Level1DisplayOrder = route.Level1DisplayOrder
					accessEntry.Level2 = route.Level2
					accessEntry.Level2DisplayOrder = route.Level2DisplayOrder
					accessEntry.Level3 = route.Level3
					accessEntry.Level3DisplayOrder = route.Level3DisplayOrder

					AccessKeyDescription := strings.Replace(route.Name, "HTTP_", "", -1)
					AccessKeyDescription = strings.Replace(AccessKeyDescription, "_", " ", -1)
					accessEntry.AccessKeyDescription = AccessKeyDescription
					//Insert into Cache
					_, err := Uc.OKAPIAUC.AUCClient.CreateAccessEntries(accessEntry)
					if err == nil {
						log.Println("Created Access Entry: ", accessEntry.Key)
					} else {
						log.Println("Error creating AccessEntry in AUC !!!")
						json.NewEncoder(os.Stdout).Encode(accessEntry)
					}
				}
			}
		}

	}
	//router.Path("/metrics").Handler(CustomPrometheusHandler())
	// router.Path("/metrics_latency").Handler(CustomPrometheusLatencyHandler())
}

func (Uc *UserControl) AddToLoyaltyFeedRouter(router *mux.Router, UC *UserControl) {
	// When StrictSlash is set to true, if the route path is "/path/", accessing "/path" will redirect
	// to the former and vice versa
	routes := CreateRoutes(UC)
	Uc.Add_LoyaltyFeedRoutes(&routes)

	for _, route := range routes {
		//var handler http.Handler
		handler := route.HandlerFunc
		//handler = Logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	//router.Path("/metrics").Handler(CustomPrometheusHandler())
	// router.Path("/metrics_latency").Handler(CustomPrometheusLatencyHandler())
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////
// authentication functions
// /////////////////////////////////////////////////////////////////////////////////////////////////////
func Use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

func (Uc *UserControl) Create_UCGW_AUCUser() {
	// Create User in AUC
	auc_user := AuthCenter.User{
		Key:                      "UCGW_Loyalty",
		FirstName:                "UCGW_Loyalty",
		MiddleName:               "UCGW_Loyalty",
		LastName:                 "UCGW_Loyalty",
		BirthDate:                time.Now(),
		Company:                  "Africell",
		Email:                    "UCGW_Loyalty",
		Phone:                    "UCGW_Loyalty",
		Login:                    "UCGW_Loyalty",
		Password:                 "Tles$h!s$w0P@$rd!$p0397",
		LastPasswordSetDate:      time.Now(),
		PasswordExpiryDate:       time.Now().AddDate(50, 0, 0),
		EnableMFA:                false,
		OTPBySMS:                 false,
		SkipOTPByMail:            true,
		LoginWithoutCaptcha:      true,
		JWEOverwriteDefault:      true,
		JWEOTPValidationValidity: 6000,
		JWEAccessValidity:        3600,
		JWERefreshValidity:       900000000,
		KeepMeLoggedIn:           true,
		PreferedLanguage:         "en",
		AddUser:                  "LoyaltyBE",
		AddDate:                  time.Now(),
		LastModifyUser:           "LoyaltyBE",
		LastModifyDate:           time.Now(),
	}
	sr, err := Uc.AppAUC.AUCClient.ReadUser(auc_user.Key)
	if err == nil && sr.StatusCode != http.StatusOK && sr.ErrorDescription == "login does not exist" {
		sr, err := Uc.AppAUC.AUCClient.CreateUser(auc_user)
		if err != nil {
			log.Println("error creating UCGW_Loyalty user: " + err.Error())
		} else if sr.StatusCode != http.StatusOK {
			log.Println("error creating UCGW_Loyalty user: " + sr.StatusDescription)
		}
	}
	return
}
