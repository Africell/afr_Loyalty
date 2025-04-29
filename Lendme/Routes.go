package Lendme

import (
	"afr_auth_center/AuthCenter"
	"encoding/json"
	"fmt"
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
	MapAccessEntry.Clear()
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
	MapAccessEntry.Clear()
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
	MapAccessEntry.Clear()

	fmt.Println("Read existing routes: #", len(accessEntries.Data))
	var existingEntries = make(map[string]AuthCenter.AccessEntry)

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
	router.Path("/metrics").Handler(CustomPrometheusHandler())
	router.Path("/metrics_latency").Handler(CustomPrometheusLatencyHandler())

	Uc.Add_Loyalty_ToAccessEntry(existingEntries)

	var Loyalty_Group = AuthCenter.UserAccessGroup{
		GroupName:        "VAS - Loyalty",
		GroupDescription: "Access group for VAS - Loyalty",
		Locked:           false,
		AddUser:          "Auto Create",
		AddDate:          time.Now(),
		LastModifyUser:   "Auto Create",
		LastModifyDate:   time.Now(),
	}

	userAccessGroups, err := Uc.OKAPIAUC.AUCClient.ReadUserAccessGroup("Value Added Services")
	if err != nil {
		fmt.Println("Error Reading VAS - Loyalty Group from AUC, shutting down !!!")
	}
	if len(userAccessGroups.Data) < 1 {

		_, err := Uc.OKAPIAUC.AUCClient.CreateUserAccessGroup((Loyalty_Group))
		if err != nil {
			fmt.Println("Error Creating VAS - Loyalty Group, shutting down !!!")
		}

		AccessEntries_to_add := []AuthCenter.GroupAccessEntries_comprehensive{
			{AccessKey: "Value Added Services", AccessMethod: "Module", Allowed: true},
			{AccessKey: "Loyalty", AccessMethod: "Module Main Menu", Allowed: true},
			{AccessKey: "Loyalty Governance", AccessMethod: "Module Sub Menu L1", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Governance", AccessMethod: "GET", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Governance", AccessMethod: "POST", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Governance", AccessMethod: "PUT", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Governance", AccessMethod: "DELETE", Allowed: true},
			{AccessKey: "Loyalty Level", AccessMethod: "Module Sub Menu L1", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Level", AccessMethod: "GET", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Level", AccessMethod: "POST", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Level", AccessMethod: "PUT", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Level", AccessMethod: "DELETE", Allowed: true},
			{AccessKey: "Loyalty Account Segment", AccessMethod: "Module Sub Menu L1", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Account_Segment", AccessMethod: "GET", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Account_Segment", AccessMethod: "POST", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Account_Segment", AccessMethod: "PUT", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Account_Segment", AccessMethod: "DELETE", Allowed: true},
			{AccessKey: "Loyalty Point Earning Rules", AccessMethod: "Module Sub Menu L1", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Point_Earning_Rules", AccessMethod: "GET", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Point_Earning_Rules", AccessMethod: "POST", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Point_Earning_Rules", AccessMethod: "PUT", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Point_Earning_Rules", AccessMethod: "DELETE", Allowed: true},
			{AccessKey: "Loyalty Point Expiry Rules", AccessMethod: "Module Sub Menu L1", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Point_Expiry_Rules", AccessMethod: "GET", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Point_Expiry_Rules", AccessMethod: "POST", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Point_Expiry_Rules", AccessMethod: "PUT", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Point_Expiry_Rules", AccessMethod: "DELETE", Allowed: true},
			{AccessKey: "Loyalty Point Redemption Rules", AccessMethod: "Module Sub Menu L1", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Point_Redemption_Rules", AccessMethod: "GET", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Point_Redemption_Rules", AccessMethod: "POST", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Point_Redemption_Rules", AccessMethod: "PUT", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Point_Redemption_Rules", AccessMethod: "DELETE", Allowed: true},
			{AccessKey: "Loyalty Plan", AccessMethod: "Module Sub Menu L1", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Plan", AccessMethod: "GET", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Plan", AccessMethod: "POST", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Plan", AccessMethod: "PUT", Allowed: true},
			{AccessKey: "HTTP_Loyalty_Plan", AccessMethod: "DELETE", Allowed: true},
			{AccessKey: "Customer Loyalty Account", AccessMethod: "Module Sub Menu L1", Allowed: true},
			{AccessKey: "HTTP_Customer_Loyalty_Account", AccessMethod: "GET", Allowed: true},
			{AccessKey: "HTTP_Customer_Loyalty_Account", AccessMethod: "POST", Allowed: true},
			{AccessKey: "HTTP_Customer_Loyalty_Account", AccessMethod: "PUT", Allowed: true},
			{AccessKey: "HTTP_Customer_Loyalty_Account", AccessMethod: "DELETE", Allowed: true},
			{AccessKey: "Customer UAT", AccessMethod: "Module Sub Menu L1", Allowed: true},
			{AccessKey: "HTTP_Customer_UAT", AccessMethod: "GET", Allowed: true},
			{AccessKey: "HTTP_Customer_UAT", AccessMethod: "POST", Allowed: true},
			{AccessKey: "HTTP_Customer_UAT", AccessMethod: "PUT", Allowed: true},
			{AccessKey: "HTTP_Customer_UAT", AccessMethod: "DELETE", Allowed: true},
			{AccessKey: "Customer DND", AccessMethod: "Module Sub Menu L1", Allowed: true},
			{AccessKey: "HTTP_Customer_DND", AccessMethod: "GET", Allowed: true},
			{AccessKey: "HTTP_Customer_DND", AccessMethod: "POST", Allowed: true},
			{AccessKey: "HTTP_Customer_DND", AccessMethod: "PUT", Allowed: true},
			{AccessKey: "HTTP_Customer_DND", AccessMethod: "DELETE", Allowed: true},
			{AccessKey: "Customer Exclusion", AccessMethod: "Module Sub Menu L1", Allowed: true},
			{AccessKey: "HTTP_Customer_Exclusion", AccessMethod: "GET", Allowed: true},
			{AccessKey: "HTTP_Customer_Exclusion", AccessMethod: "POST", Allowed: true},
			{AccessKey: "HTTP_Customer_Exclusion", AccessMethod: "PUT", Allowed: true},
			{AccessKey: "HTTP_Customer_Exclusion", AccessMethod: "DELETE", Allowed: true},
			{AccessKey: "Customer COS Exclusion", AccessMethod: "Module Sub Menu L1", Allowed: true},
			{AccessKey: "HTTP_Customer_COS_Exclusion", AccessMethod: "GET", Allowed: true},
			{AccessKey: "HTTP_Customer_COS_Exclusion", AccessMethod: "POST", Allowed: true},
			{AccessKey: "HTTP_Customer_COS_Exclusion", AccessMethod: "PUT", Allowed: true},
			{AccessKey: "HTTP_Customer_COS_Exclusion", AccessMethod: "DELETE", Allowed: true},
			{AccessKey: "HTTP_Customer_Loyalty_Account_GetRedemption_Rules", AccessMethod: "GET", Allowed: true},
			{AccessKey: "HTTP_Customer_Loyalty_RedeemRequest", AccessMethod: "PUT", Allowed: true},
			{AccessKey: "HTTP_Customer_Loyalty_Account_Points_Details", AccessMethod: "GET", Allowed: true},
			{AccessKey: "HTTP_Loyalty_AccountDebitPoints_log", AccessMethod: "GET", Allowed: true},
			{AccessKey: "HTTP_Loyalty_AccountCreditPoints_log", AccessMethod: "GET", Allowed: true},
		}
		_, err = Uc.OKAPIAUC.AUCClient.GroupAccessEntriesForGroup_Comprehensive("VAS - Loyalty", AccessEntries_to_add)

		if err != nil {
			fmt.Println(err)
			fmt.Println("Error Creating VAS - Loyalty, shutting down !!!")
		}
	}

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

func (Uc *UserControl) AddToOKAPIAccessEntry(existing map[string]AuthCenter.AccessEntry, accessEntry AuthCenter.AccessEntry) {
	accessEntry.Key = accessEntry.AccessKey + "|" + accessEntry.AccessMethod
	_, ok := existing[accessEntry.Key]

	if !ok {
		json.NewEncoder(os.Stdout).Encode(accessEntry)
		accessEntry.Id = 0
		accessEntry.Status = ""
		accessEntry.AddUser = "Auto Add"
		accessEntry.AddDate = time.Now().UTC()
		accessEntry.LastModifyUser = "Auto Add"
		accessEntry.LastModifyDate = time.Now().UTC()
		//Insert into Cache
		_, err := Uc.OKAPIAUC.AUCClient.CreateAccessEntries(accessEntry)
		if err == nil {
			MapAccessEntry.Put(accessEntry.AccessKey+"|"+accessEntry.AccessMethod, accessEntry)
			fmt.Println("Created Access Entry: ", accessEntry.Key)
		} else {
			json.NewEncoder(os.Stdout).Encode(accessEntry)
			fmt.Println("Error creating AccessEntry in AUC !!!")
		}
	}
}

func (Uc *UserControl) Add_Loyalty_ToAccessEntry(existing map[string]AuthCenter.AccessEntry) {
	// Access Method: Module, Module Main Menu, Module Sub Menu L1, Module Sub Menu L2
	Module := "Value Added Services"
	var ModuleDisplayOrder int64 = 13

	//Module (Purple circles)
	sd_ae := AuthCenter.AccessEntry{
		AccessKey:            "Value Added Services",
		AccessMethod:         "Module",
		AccessKeyDescription: "Value Added Services Module",
		DisplayName:          "Value Added Services",
		DisplayOrder:         13,
		Module:               Module,
		ModuleDisplayOrder:   ModuleDisplayOrder,
		Level1:               "",
		Level1DisplayOrder:   0,
		Level2:               "",
		Level2DisplayOrder:   0,
		Level3:               "",
		Level3DisplayOrder:   0,
	}
	Uc.AddToOKAPIAccessEntry(existing, sd_ae)

	//Module Main Menu (yellow circles)
	Level1 := "Loyalty"
	var Level1DisplayOrder int64 = 6

	sd_ae = AuthCenter.AccessEntry{
		AccessKey:            "Loyalty",
		AccessMethod:         "Module Main Menu",
		AccessKeyDescription: "Loyalty",
		DisplayName:          "Loyalty",
		DisplayOrder:         13,
		Module:               Module,
		ModuleDisplayOrder:   ModuleDisplayOrder,
		Level1:               Level1,
		Level1DisplayOrder:   Level1DisplayOrder,
		Level2:               "",
		Level2DisplayOrder:   0,
		Level3:               "",
		Level3DisplayOrder:   0,
	}
	Uc.AddToOKAPIAccessEntry(existing, sd_ae)

	sd_ae = AuthCenter.AccessEntry{
		AccessKey:            "Loyalty Governance",
		AccessMethod:         "Module Sub Menu L1",
		AccessKeyDescription: "Loyalty Governance",
		DisplayName:          "Loyalty Governance",
		DisplayOrder:         13,
		Module:               Module,
		ModuleDisplayOrder:   ModuleDisplayOrder,
		Level1:               Level1,
		Level1DisplayOrder:   Level1DisplayOrder,
		Level2:               "Loyalty Governance",
		Level2DisplayOrder:   1,
		Level3:               "",
		Level3DisplayOrder:   0,
	}
	Uc.AddToOKAPIAccessEntry(existing, sd_ae)

	sd_ae = AuthCenter.AccessEntry{
		AccessKey:            "Loyalty Level",
		AccessMethod:         "Module Sub Menu L1",
		AccessKeyDescription: "Loyalty Level",
		DisplayName:          "Loyalty Level",
		DisplayOrder:         13,
		Module:               Module,
		ModuleDisplayOrder:   ModuleDisplayOrder,
		Level1:               Level1,
		Level1DisplayOrder:   Level1DisplayOrder,
		Level2:               "Loyalty Level",
		Level2DisplayOrder:   2,
		Level3:               "",
		Level3DisplayOrder:   0,
	}
	Uc.AddToOKAPIAccessEntry(existing, sd_ae)

	sd_ae = AuthCenter.AccessEntry{
		AccessKey:            "Loyalty Account Segment",
		AccessMethod:         "Module Sub Menu L1",
		AccessKeyDescription: "Loyalty Account Segment",
		DisplayName:          "Loyalty Account Segment",
		DisplayOrder:         13,
		Module:               Module,
		ModuleDisplayOrder:   ModuleDisplayOrder,
		Level1:               Level1,
		Level1DisplayOrder:   Level1DisplayOrder,
		Level2:               "Loyalty Account Segment",
		Level2DisplayOrder:   3,
		Level3:               "",
		Level3DisplayOrder:   0,
	}
	Uc.AddToOKAPIAccessEntry(existing, sd_ae)

	sd_ae = AuthCenter.AccessEntry{
		AccessKey:            "Loyalty Point Earning Rules",
		AccessMethod:         "Module Sub Menu L1",
		AccessKeyDescription: "Loyalty Point Earning Rules",
		DisplayName:          "Loyalty Point Earning Rules",
		DisplayOrder:         13,
		Module:               Module,
		ModuleDisplayOrder:   ModuleDisplayOrder,
		Level1:               Level1,
		Level1DisplayOrder:   Level1DisplayOrder,
		Level2:               "Loyalty Point Earning Rules",
		Level2DisplayOrder:   4,
		Level3:               "",
		Level3DisplayOrder:   0,
	}
	Uc.AddToOKAPIAccessEntry(existing, sd_ae)

	sd_ae = AuthCenter.AccessEntry{
		AccessKey:            "Loyalty Point Expiry Rules",
		AccessMethod:         "Module Sub Menu L1",
		AccessKeyDescription: "Loyalty Point Expiry Rules",
		DisplayName:          "Loyalty Point Expiry Rules",
		DisplayOrder:         13,
		Module:               Module,
		ModuleDisplayOrder:   ModuleDisplayOrder,
		Level1:               Level1,
		Level1DisplayOrder:   Level1DisplayOrder,
		Level2:               "Loyalty Point Expiry Rules",
		Level2DisplayOrder:   5,
		Level3:               "",
		Level3DisplayOrder:   0,
	}
	Uc.AddToOKAPIAccessEntry(existing, sd_ae)

	sd_ae = AuthCenter.AccessEntry{
		AccessKey:            "Loyalty Point Redemption Rules",
		AccessMethod:         "Module Sub Menu L1",
		AccessKeyDescription: "Loyalty Point Redemption Rules",
		DisplayName:          "Loyalty Point Redemption Rules",
		DisplayOrder:         13,
		Module:               Module,
		ModuleDisplayOrder:   ModuleDisplayOrder,
		Level1:               Level1,
		Level1DisplayOrder:   Level1DisplayOrder,
		Level2:               "Loyalty Point Redemption Rules",
		Level2DisplayOrder:   6,
		Level3:               "",
		Level3DisplayOrder:   0,
	}
	Uc.AddToOKAPIAccessEntry(existing, sd_ae)

	sd_ae = AuthCenter.AccessEntry{
		AccessKey:            "Loyalty Plan",
		AccessMethod:         "Module Sub Menu L1",
		AccessKeyDescription: "Loyalty Plan",
		DisplayName:          "Loyalty Plan",
		DisplayOrder:         13,
		Module:               Module,
		ModuleDisplayOrder:   ModuleDisplayOrder,
		Level1:               Level1,
		Level1DisplayOrder:   Level1DisplayOrder,
		Level2:               "Loyalty Plan",
		Level2DisplayOrder:   7,
		Level3:               "",
		Level3DisplayOrder:   0,
	}
	Uc.AddToOKAPIAccessEntry(existing, sd_ae)

	sd_ae = AuthCenter.AccessEntry{
		AccessKey:            "Customer Loyalty Account",
		AccessMethod:         "Module Sub Menu L1",
		AccessKeyDescription: "Customer Loyalty Account",
		DisplayName:          "Customer Loyalty Account",
		DisplayOrder:         13,
		Module:               Module,
		ModuleDisplayOrder:   ModuleDisplayOrder,
		Level1:               Level1,
		Level1DisplayOrder:   Level1DisplayOrder,
		Level2:               "Customer Loyalty Account",
		Level2DisplayOrder:   8,
		Level3:               "",
		Level3DisplayOrder:   0,
		Routes:               []string{"HTTP_Customer_Loyalty_Account_GetRedemption_Rules|GET", "HTTP_Customer_Loyalty_RedeemRequest|PUT", "HTTP_Loyalty_AccountDebitPoints_log|GET", "HTTP_Customer_Loyalty_Account_Points_Details|GET", "HTTP_Loyalty_AccountCreditPoints_log|GET"},
	}
	Uc.AddToOKAPIAccessEntry(existing, sd_ae)

	sd_ae = AuthCenter.AccessEntry{
		AccessKey:            "Customer UAT",
		AccessMethod:         "Module Sub Menu L1",
		AccessKeyDescription: "Customer UAT",
		DisplayName:          "Customer UAT",
		DisplayOrder:         13,
		Module:               Module,
		ModuleDisplayOrder:   ModuleDisplayOrder,
		Level1:               Level1,
		Level1DisplayOrder:   Level1DisplayOrder,
		Level2:               "Customer UAT",
		Level2DisplayOrder:   9,
		Level3:               "",
		Level3DisplayOrder:   0,
	}
	Uc.AddToOKAPIAccessEntry(existing, sd_ae)

	sd_ae = AuthCenter.AccessEntry{
		AccessKey:            "Customer DND",
		AccessMethod:         "Module Sub Menu L1",
		AccessKeyDescription: "Customer DND",
		DisplayName:          "Customer DND",
		DisplayOrder:         13,
		Module:               Module,
		ModuleDisplayOrder:   ModuleDisplayOrder,
		Level1:               Level1,
		Level1DisplayOrder:   Level1DisplayOrder,
		Level2:               "Customer DND",
		Level2DisplayOrder:   10,
		Level3:               "",
		Level3DisplayOrder:   0,
	}
	Uc.AddToOKAPIAccessEntry(existing, sd_ae)

	sd_ae = AuthCenter.AccessEntry{
		AccessKey:            "Customer Exclusion",
		AccessMethod:         "Module Sub Menu L1",
		AccessKeyDescription: "Customer Exclusion",
		DisplayName:          "Customer Exclusion",
		DisplayOrder:         13,
		Module:               Module,
		ModuleDisplayOrder:   ModuleDisplayOrder,
		Level1:               Level1,
		Level1DisplayOrder:   Level1DisplayOrder,
		Level2:               "Customer Exclusion",
		Level2DisplayOrder:   11,
		Level3:               "",
		Level3DisplayOrder:   0,
	}
	Uc.AddToOKAPIAccessEntry(existing, sd_ae)

	sd_ae = AuthCenter.AccessEntry{
		AccessKey:            "Customer COS Exclusion",
		AccessMethod:         "Module Sub Menu L1",
		AccessKeyDescription: "Customer COS Exclusion",
		DisplayName:          "Customer COS Exclusion",
		DisplayOrder:         13,
		Module:               Module,
		ModuleDisplayOrder:   ModuleDisplayOrder,
		Level1:               Level1,
		Level1DisplayOrder:   Level1DisplayOrder,
		Level2:               "Customer ExclCOS Exclusionusion",
		Level2DisplayOrder:   12,
		Level3:               "",
		Level3DisplayOrder:   0,
	}
	Uc.AddToOKAPIAccessEntry(existing, sd_ae)

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
