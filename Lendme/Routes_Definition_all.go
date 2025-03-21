package Lendme

func (UC *UserControl) Add_LendmeRoutes(R *Routes) {
	var r Route
	var Module, Level1 string
	var ModuleDisplayOrder int64
	Module = Configuration.Module
	ModuleDisplayOrder = 1
	var Level1DisplayOrder int64
	var DisplayOrder int64 = 1

	//*****************************
	// Credit Limit Scheme
	//*****************************
	Level1 = "Credit Limit Scheme"
	Level1DisplayOrder = Level1DisplayOrder + 1
	r = Route{
		"HTTP_Credit_Limit_Scheme",
		"GET",
		"/" + Configuration.Module + "/" + Configuration.Version + "/HTTP_Credit_Limit_Scheme/",
		UC.HTTP_Credit_Limit_Scheme,
		true,
		"Credit Limit Scheme - Read", // DisplayName
		DisplayOrder,                 // DisplayOrder
		Module,                       // Module
		ModuleDisplayOrder,           //ModuleDisplayOrder
		Level1,                       // Level1
		Level1DisplayOrder,           // Level1DisplayOrder
		"",                           // Level2
		0,                            // Level2DisplayOrder
		"",                           // Level3
		0,                            // Level3DisplayOrder
		true,                         //AllowedFor_OKAPI
		true,                         //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Credit_Limit_Scheme",
		"POST",
		"/" + Configuration.Module + "/" + Configuration.Version + "/HTTP_Credit_Limit_Scheme/",
		UC.HTTP_Credit_Limit_Scheme,
		true,
		"Credit Limit Scheme - Add", // DisplayName
		DisplayOrder,                // DisplayOrder
		Module,                      // Module
		ModuleDisplayOrder,          //ModuleDisplayOrder
		Level1,                      // Level1
		Level1DisplayOrder,          // Level1DisplayOrder
		"",                          // Level2
		0,                           // Level2DisplayOrder
		"",                          // Level3
		0,                           // Level3DisplayOrder
		true,                        //AllowedFor_OKAPI
		true,                        //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Credit_Limit_Scheme",
		"PUT",
		"/" + Configuration.Module + "/" + Configuration.Version + "/HTTP_Credit_Limit_Scheme/",
		UC.HTTP_Credit_Limit_Scheme,
		true,
		"Credit Limit Scheme - Edit", // DisplayName
		DisplayOrder,                 // DisplayOrder
		Module,                       // Module
		ModuleDisplayOrder,           //ModuleDisplayOrder
		Level1,                       // Level1
		Level1DisplayOrder,           // Level1DisplayOrder
		"",                           // Level2
		0,                            // Level2DisplayOrder
		"",                           // Level3
		0,                            // Level3DisplayOrder
		true,                         //AllowedFor_OKAPI
		true,                         //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Credit_Limit_Scheme",
		"DELETE",
		"/" + Configuration.Module + "/" + Configuration.Version + "/HTTP_Credit_Limit_Scheme/{KeyDelete}/",
		UC.HTTP_Credit_Limit_Scheme,
		true,
		"Credit Limit Scheme - Delete", // DisplayName
		DisplayOrder,                   // DisplayOrder
		Module,                         // Module
		ModuleDisplayOrder,             //ModuleDisplayOrder
		Level1,                         // Level1
		Level1DisplayOrder,             // Level1DisplayOrder
		"",                             // Level2
		0,                              // Level2DisplayOrder
		"",                             // Level3
		0,                              // Level3DisplayOrder
		true,                           //AllowedFor_OKAPI
		true,                           //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	//*****************************
	// Subscribers
	//*****************************
	Level1 = "Subscriber"
	Level1DisplayOrder = Level1DisplayOrder + 1
	r = Route{
		"HTTP_Subscriber",
		"GET",
		"/" + Configuration.Module + "/" + Configuration.Version + "/HTTP_Subscriber/",
		UC.HTTP_Subscriber,
		true,
		"Subscriber - Read", // DisplayName
		DisplayOrder,        // DisplayOrder
		Module,              // Module
		ModuleDisplayOrder,  //ModuleDisplayOrder
		Level1,              // Level1
		Level1DisplayOrder,  // Level1DisplayOrder
		"",                  // Level2
		0,                   // Level2DisplayOrder
		"",                  // Level3
		0,                   // Level3DisplayOrder
		true,                //AllowedFor_OKAPI
		true,                //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Subscriber_USSD",
		"GET",
		"/" + Configuration.Module + "/" + Configuration.Version + "/HTTP_Subscriber_USSD/",
		UC.HTTP_Subscriber_USSD,
		true,
		"Subscriber USSD - Read", // DisplayName
		DisplayOrder,             // DisplayOrder
		Module,                   // Module
		ModuleDisplayOrder,       //ModuleDisplayOrder
		Level1,                   // Level1
		Level1DisplayOrder,       // Level1DisplayOrder
		"",                       // Level2
		0,                        // Level2DisplayOrder
		"",                       // Level3
		0,                        // Level3DisplayOrder
		true,                     //AllowedFor_OKAPI
		true,                     //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Subscriber",
		"POST",
		"/" + Configuration.Module + "/" + Configuration.Version + "/HTTP_Subscriber/",
		UC.HTTP_Subscriber,
		true,
		"Subscriber - Add", // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Subscriber",
		"PUT",
		"/" + Configuration.Module + "/" + Configuration.Version + "/HTTP_Subscriber/",
		UC.HTTP_Subscriber,
		true,
		"Subscriber - Edit", // DisplayName
		DisplayOrder,        // DisplayOrder
		Module,              // Module
		ModuleDisplayOrder,  //ModuleDisplayOrder
		Level1,              // Level1
		Level1DisplayOrder,  // Level1DisplayOrder
		"",                  // Level2
		0,                   // Level2DisplayOrder
		"",                  // Level3
		0,                   // Level3DisplayOrder
		true,                //AllowedFor_OKAPI
		true,                //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Subscriber",
		"DELETE",
		"/" + Configuration.Module + "/" + Configuration.Version + "/HTTP_Subscriber/{KeyDelete}/",
		UC.HTTP_Subscriber,
		true,
		"Subscriber - Delete", // DisplayName
		DisplayOrder,          // DisplayOrder
		Module,                // Module
		ModuleDisplayOrder,    //ModuleDisplayOrder
		Level1,                // Level1
		Level1DisplayOrder,    // Level1DisplayOrder
		"",                    // Level2
		0,                     // Level2DisplayOrder
		"",                    // Level3
		0,                     // Level3DisplayOrder
		true,                  //AllowedFor_OKAPI
		true,                  //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	//*********************************
	// Subscribers ARPU Import Launch
	//*********************************
	Level1 = "Subscribers ARPU Import"
	Level1DisplayOrder = Level1DisplayOrder + 1
	r = Route{
		"HTTP_Subscribers_ARPU_Import_Launch",
		"GET",
		"/" + Configuration.Module + "/" + Configuration.Version + "/HTTP_Subscribers_ARPU_Import_Launch/",
		UC.HTTP_Subscribers_ARPU_Import_Launch,
		true,
		"Subscribers ARPU Import Launch", // DisplayName
		DisplayOrder,                     // DisplayOrder
		Module,                           // Module
		ModuleDisplayOrder,               //ModuleDisplayOrder
		Level1,                           // Level1
		Level1DisplayOrder,               // Level1DisplayOrder
		"",                               // Level2
		0,                                // Level2DisplayOrder
		"",                               // Level3
		0,                                // Level3DisplayOrder
		true,                             //AllowedFor_OKAPI
		true,                             //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	//*****************************
	// Get IN subscriber detail
	//*****************************
	Level1 = "Subscribers IN functions"
	Level1DisplayOrder = Level1DisplayOrder + 1
	r = Route{
		"HTTP_Subscribers_GetINDetail",
		"GET",
		"/" + Configuration.Module + "/" + Configuration.Version + "/HTTP_Subscribers_GetINDetail/",
		UC.HTTP_Subscribers_GetINDetail,
		true,
		"Subscribers Get IN Detail", // DisplayName
		DisplayOrder,                // DisplayOrder
		Module,                      // Module
		ModuleDisplayOrder,          //ModuleDisplayOrder
		Level1,                      // Level1
		Level1DisplayOrder,          // Level1DisplayOrder
		"",                          // Level2
		0,                           // Level2DisplayOrder
		"",                          // Level3
		0,                           // Level3DisplayOrder
		true,                        //AllowedFor_OKAPI
		true,                        //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	//*****************************
	// Subscriber Lendme Request
	//*****************************
	r = Route{
		"HTTP_Lendme_Request",
		"POST",
		"/" + Configuration.Module + "/" + Configuration.Version + "/HTTP_Lendme_Request/",
		UC.HTTP_Lendme_Request,
		true,
		"Lendme Request - Add", // DisplayName
		DisplayOrder,           // DisplayOrder
		Module,                 // Module
		ModuleDisplayOrder,     //ModuleDisplayOrder
		Level1,                 // Level1
		Level1DisplayOrder,     // Level1DisplayOrder
		"",                     // Level2
		0,                      // Level2DisplayOrder
		"",                     // Level3
		0,                      // Level3DisplayOrder
		true,                   //AllowedFor_OKAPI
		true,                   //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	//*****************************
	// Subscriber PayBack
	//*****************************
	r = Route{
		"HTTP_Lendme_PayBack",
		"POST",
		"/" + Configuration.Module + "/" + Configuration.Version + "/HTTP_Lendme_PayBack/",
		UC.HTTP_Lendme_PayBack,
		true,
		"Lendme PayBack - Add", // DisplayName
		DisplayOrder,           // DisplayOrder
		Module,                 // Module
		ModuleDisplayOrder,     //ModuleDisplayOrder
		Level1,                 // Level1
		Level1DisplayOrder,     // Level1DisplayOrder
		"",                     // Level2
		0,                      // Level2DisplayOrder
		"",                     // Level3
		0,                      // Level3DisplayOrder
		true,                   //AllowedFor_OKAPI
		true,                   //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	//*****************************
	// Subscriber PayBack
	//*****************************
	r = Route{
		"HTTP_Lendme_SendSMS",
		"POST",
		"/" + Configuration.Module + "/" + Configuration.Version + "/HTTP_Lendme_SendSMS/",
		UC.HTTP_Lendme_SendSMS,
		true,
		"Lendme Send SMS",  // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1
}

func (UC *UserControl) Add_LoyaltyRoutes(R *Routes) {
	var r Route
	var Module, Level1 string
	var ModuleDisplayOrder int64
	Module = Configuration.LoyaltyModule
	ModuleDisplayOrder = 1
	var Level1DisplayOrder int64
	var DisplayOrder int64 = 1

	//*****************************
	// Loyalty_Governance
	//*****************************
	Level1 = "Loyalty Governance"
	Level1DisplayOrder = Level1DisplayOrder + 1
	r = Route{
		"HTTP_Loyalty_Governance",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Governance/",
		UC.HTTP_Loyalty_Governance,
		true,
		Level1 + " - Read", // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Loyalty_Governance",
		"POST",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Governance/",
		UC.HTTP_Loyalty_Governance,
		true,
		Level1 + " - Add",  // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Loyalty_Governance",
		"PUT",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Governance/",
		UC.HTTP_Loyalty_Governance,
		true,
		Level1 + " - Edit", // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Loyalty_Governance",
		"DELETE",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Governance/{KeyDelete}/",
		UC.HTTP_Loyalty_Governance,
		true,
		Level1 + " - Delete", // DisplayName
		DisplayOrder,         // DisplayOrder
		Module,               // Module
		ModuleDisplayOrder,   //ModuleDisplayOrder
		Level1,               // Level1
		Level1DisplayOrder,   // Level1DisplayOrder
		"",                   // Level2
		0,                    // Level2DisplayOrder
		"",                   // Level3
		0,                    // Level3DisplayOrder
		true,                 //AllowedFor_OKAPI
		true,                 //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	//*****************************
	// Loyalty_Level
	//*****************************
	Level1 = "Loyalty Level"
	Level1DisplayOrder = Level1DisplayOrder + 1
	r = Route{
		"HTTP_Loyalty_Level",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Level/",
		UC.HTTP_Loyalty_Level,
		true,
		Level1 + " - Read", // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Loyalty_Level",
		"POST",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Level/",
		UC.HTTP_Loyalty_Level,
		true,
		Level1 + " - Add",  // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Loyalty_Level",
		"PUT",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Level/",
		UC.HTTP_Loyalty_Level,
		true,
		Level1 + " - Edit", // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Loyalty_Level",
		"DELETE",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Level/{KeyDelete}/",
		UC.HTTP_Loyalty_Level,
		true,
		Level1 + " - Delete", // DisplayName
		DisplayOrder,         // DisplayOrder
		Module,               // Module
		ModuleDisplayOrder,   //ModuleDisplayOrder
		Level1,               // Level1
		Level1DisplayOrder,   // Level1DisplayOrder
		"",                   // Level2
		0,                    // Level2DisplayOrder
		"",                   // Level3
		0,                    // Level3DisplayOrder
		true,                 //AllowedFor_OKAPI
		true,                 //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	//*****************************
	// Loyalty_Account_Segment
	//*****************************
	Level1 = "Loyalty Account Segment"
	Level1DisplayOrder = Level1DisplayOrder + 1
	r = Route{
		"HTTP_Loyalty_Account_Segment",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Account_Segment/",
		UC.HTTP_Loyalty_Account_Segment,
		true,
		Level1 + " - Read", // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Loyalty_Account_Segment",
		"POST",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Account_Segment/",
		UC.HTTP_Loyalty_Account_Segment,
		true,
		Level1 + " - Add",  // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Loyalty_Account_Segment",
		"PUT",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Account_Segment/",
		UC.HTTP_Loyalty_Account_Segment,
		true,
		Level1 + " - Edit", // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Loyalty_Account_Segment",
		"DELETE",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Account_Segment/{KeyDelete}/",
		UC.HTTP_Loyalty_Account_Segment,
		true,
		Level1 + " - Delete", // DisplayName
		DisplayOrder,         // DisplayOrder
		Module,               // Module
		ModuleDisplayOrder,   //ModuleDisplayOrder
		Level1,               // Level1
		Level1DisplayOrder,   // Level1DisplayOrder
		"",                   // Level2
		0,                    // Level2DisplayOrder
		"",                   // Level3
		0,                    // Level3DisplayOrder
		true,                 //AllowedFor_OKAPI
		true,                 //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	//*****************************
	// Loyalty_Point_Earning_Rules
	//*****************************
	Level1 = "Loyalty Point Earning Rules"
	Level1DisplayOrder = Level1DisplayOrder + 1
	r = Route{
		"HTTP_Loyalty_Point_Earning_Rules",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Point_Earning_Rules/",
		UC.HTTP_Loyalty_Point_Earning_Rules,
		true,
		Level1 + " - Read", // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Loyalty_Point_Earning_Rules",
		"POST",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Point_Earning_Rules/",
		UC.HTTP_Loyalty_Point_Earning_Rules,
		true,
		Level1 + " - Add",  // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Loyalty_Point_Earning_Rules",
		"PUT",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Point_Earning_Rules/",
		UC.HTTP_Loyalty_Point_Earning_Rules,
		true,
		Level1 + " - Edit", // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Loyalty_Point_Earning_Rules",
		"DELETE",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Point_Earning_Rules/{KeyDelete}/",
		UC.HTTP_Loyalty_Point_Earning_Rules,
		true,
		Level1 + " - Delete", // DisplayName
		DisplayOrder,         // DisplayOrder
		Module,               // Module
		ModuleDisplayOrder,   //ModuleDisplayOrder
		Level1,               // Level1
		Level1DisplayOrder,   // Level1DisplayOrder
		"",                   // Level2
		0,                    // Level2DisplayOrder
		"",                   // Level3
		0,                    // Level3DisplayOrder
		true,                 //AllowedFor_OKAPI
		true,                 //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	//*****************************
	// Loyalty_Point_Expiry_Rules
	//*****************************
	Level1 = "Loyalty Point Expiry Rules"
	Level1DisplayOrder = Level1DisplayOrder + 1
	r = Route{
		"HTTP_Loyalty_Point_Expiry_Rules",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Point_Expiry_Rules/",
		UC.HTTP_Loyalty_Point_Expiry_Rules,
		true,
		Level1 + " - Read", // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Loyalty_Point_Expiry_Rules",
		"POST",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Point_Expiry_Rules/",
		UC.HTTP_Loyalty_Point_Expiry_Rules,
		true,
		Level1 + " - Add",  // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Loyalty_Point_Expiry_Rules",
		"PUT",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Point_Expiry_Rules/",
		UC.HTTP_Loyalty_Point_Expiry_Rules,
		true,
		Level1 + " - Edit", // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Loyalty_Point_Expiry_Rules",
		"DELETE",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Point_Expiry_Rules/{KeyDelete}/",
		UC.HTTP_Loyalty_Point_Expiry_Rules,
		true,
		Level1 + " - Delete", // DisplayName
		DisplayOrder,         // DisplayOrder
		Module,               // Module
		ModuleDisplayOrder,   //ModuleDisplayOrder
		Level1,               // Level1
		Level1DisplayOrder,   // Level1DisplayOrder
		"",                   // Level2
		0,                    // Level2DisplayOrder
		"",                   // Level3
		0,                    // Level3DisplayOrder
		true,                 //AllowedFor_OKAPI
		true,                 //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	//*****************************
	// Loyalty_Point_Redemption_Rules
	//*****************************
	Level1 = "Loyalty Point Redemption Rules"
	Level1DisplayOrder = Level1DisplayOrder + 1
	r = Route{
		"HTTP_Loyalty_Point_Redemption_Rules",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Point_Redemption_Rules/",
		UC.HTTP_Loyalty_Point_Redemption_Rules,
		true,
		Level1 + " - Read", // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Loyalty_Point_Redemption_Rules",
		"POST",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Point_Redemption_Rules/",
		UC.HTTP_Loyalty_Point_Redemption_Rules,
		true,
		Level1 + " - Add",  // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Loyalty_Point_Redemption_Rules",
		"PUT",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Point_Redemption_Rules/",
		UC.HTTP_Loyalty_Point_Redemption_Rules,
		true,
		Level1 + " - Edit", // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Loyalty_Point_Redemption_Rules",
		"DELETE",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Point_Redemption_Rules/{KeyDelete}/",
		UC.HTTP_Loyalty_Point_Redemption_Rules,
		true,
		Level1 + " - Delete", // DisplayName
		DisplayOrder,         // DisplayOrder
		Module,               // Module
		ModuleDisplayOrder,   //ModuleDisplayOrder
		Level1,               // Level1
		Level1DisplayOrder,   // Level1DisplayOrder
		"",                   // Level2
		0,                    // Level2DisplayOrder
		"",                   // Level3
		0,                    // Level3DisplayOrder
		true,                 //AllowedFor_OKAPI
		true,                 //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	//*****************************
	// Loyalty_Plan
	//*****************************
	Level1 = "Loyalty Plan"
	Level1DisplayOrder = Level1DisplayOrder + 1
	r = Route{
		"HTTP_Loyalty_Plan",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Plan/",
		UC.HTTP_Loyalty_Plan,
		true,
		Level1 + " - Read", // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Loyalty_Plan",
		"POST",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Plan/",
		UC.HTTP_Loyalty_Plan,
		true,
		Level1 + " - Add",  // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Loyalty_Plan",
		"PUT",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Plan/",
		UC.HTTP_Loyalty_Plan,
		true,
		Level1 + " - Edit", // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Loyalty_Plan",
		"DELETE",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Plan/{KeyDelete}/",
		UC.HTTP_Loyalty_Plan,
		true,
		Level1 + " - Delete", // DisplayName
		DisplayOrder,         // DisplayOrder
		Module,               // Module
		ModuleDisplayOrder,   //ModuleDisplayOrder
		Level1,               // Level1
		Level1DisplayOrder,   // Level1DisplayOrder
		"",                   // Level2
		0,                    // Level2DisplayOrder
		"",                   // Level3
		0,                    // Level3DisplayOrder
		true,                 //AllowedFor_OKAPI
		true,                 //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	//*****************************
	// Loyalty_Plan
	//*****************************
	Level1 = "Customer Loyalty Account"
	Level1DisplayOrder = Level1DisplayOrder + 1
	r = Route{
		"HTTP_Customer_Loyalty_Account",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_Loyalty_Account/",
		UC.HTTP_Customer_Loyalty_Account,
		true,
		Level1 + " - Read", // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Customer_Loyalty_Account",
		"POST",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_Loyalty_Account/",
		UC.HTTP_Customer_Loyalty_Account,
		true,
		Level1 + " - Add",  // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Customer_Loyalty_Account",
		"PUT",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_Loyalty_Account/",
		UC.HTTP_Customer_Loyalty_Account,
		true,
		Level1 + " - Edit", // DisplayName
		DisplayOrder,       // DisplayOrder
		Module,             // Module
		ModuleDisplayOrder, //ModuleDisplayOrder
		Level1,             // Level1
		Level1DisplayOrder, // Level1DisplayOrder
		"",                 // Level2
		0,                  // Level2DisplayOrder
		"",                 // Level3
		0,                  // Level3DisplayOrder
		true,               //AllowedFor_OKAPI
		true,               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Customer_Loyalty_Account",
		"DELETE",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_Loyalty_Account/{KeyDelete}/",
		UC.HTTP_Customer_Loyalty_Account,
		true,
		Level1 + " - Delete", // DisplayName
		DisplayOrder,         // DisplayOrder
		Module,               // Module
		ModuleDisplayOrder,   //ModuleDisplayOrder
		Level1,               // Level1
		Level1DisplayOrder,   // Level1DisplayOrder
		"",                   // Level2
		0,                    // Level2DisplayOrder
		"",                   // Level3
		0,                    // Level3DisplayOrder
		true,                 //AllowedFor_OKAPI
		true,                 //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1
}
