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

func (UC *UserControl) Add_LoyaltyServiceRoutes(R *Routes) {
	var r Route
	var Module, Level1 string
	var ModuleDisplayOrder int64
	Module = Configuration.LoyaltyModule
	ModuleDisplayOrder = 1
	var Level1DisplayOrder int64
	var DisplayOrder int64 = 1

	Level1 = "Customer Loyalty Account"
	Level1DisplayOrder = Level1DisplayOrder + 1
	r = Route{
		"HTTP_Customer_Loyalty_Account",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_Loyalty_Account/",
		Use(UC.HTTP_Customer_Loyalty_Account, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
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

	Level1 = "Customer Loyalty Account - Products Catalogue"
	Level1DisplayOrder = Level1DisplayOrder + 1
	r = Route{
		"HTTP_Loyalty_Products_Catalogue",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Products_Catalogue/",
		Use(UC.HTTP_Loyalty_Products_Catalogue, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		Level1 + "",        // DisplayName
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

	Level1 = "Customer Loyalty Account - Debit Points"
	Level1DisplayOrder = Level1DisplayOrder + 1
	r = Route{
		"HTTP_Customer_Loyalty_Account_DebitPoints",
		"PUT",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_Loyalty_Account_DebitPoints/",
		Use(UC.HTTP_Customer_Loyalty_Account_DebitPoints, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		Level1 + "",        // DisplayName
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

	Level1 = "Customer Loyalty - Redeem Request"
	Level1DisplayOrder = Level1DisplayOrder + 1
	r = Route{
		"HTTP_Customer_Loyalty_RedeemRequest",
		"PUT",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_Loyalty_RedeemRequest/",
		Use(UC.HTTP_Customer_Loyalty_RedeemRequest, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		Level1 + "",        // DisplayName
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

	Level1 = "Customer Loyalty - Get Redemption Rules"
	Level1DisplayOrder = Level1DisplayOrder + 1
	r = Route{
		"HTTP_Customer_Loyalty_Account_GetRedemption_Rules",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_Loyalty_Account_GetRedemption_Rules/",
		Use(UC.HTTP_Customer_Loyalty_Account_GetRedemption_Rules, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		Level1 + "",        // DisplayName
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

func (UC *UserControl) Add_LoyaltyManagementRoutes(R *Routes) {
	var r Route
	var Module, Level1, Level2 string
	var ModuleDisplayOrder, Level1DisplayOrder, Level2DisplayOrder int64
	Module = "Value Added Services"
	ModuleDisplayOrder = 13
	Level1 = "Loyalty"
	Level1DisplayOrder = 6
	var DisplayOrder int64 = 1
	//*****************************
	// Loyalty_Governance
	//*****************************
	Level2 = "Loyalty Governance"
	Level2DisplayOrder = 1
	r = Route{
		"HTTP_Loyalty_Governance",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Governance/",
		Use(UC.HTTP_Loyalty_Governance, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Governance - Read", // DisplayName
		1,                           // DisplayOrder
		Module,                      // Module
		ModuleDisplayOrder,          //ModuleDisplayOrder
		Level1,                      // Level1
		Level1DisplayOrder,          // Level1DisplayOrder
		Level2,                      // Level2
		Level2DisplayOrder,          // Level2DisplayOrder
		"",                          // Level3
		0,                           // Level3DisplayOrder
		true,                        //AllowedFor_OKAPI
		true,                        //AllowedFor_App
	}
	*R = append(*R, r)

	r = Route{
		"HTTP_Loyalty_Governance",
		"POST",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Governance/",
		Use(UC.HTTP_Loyalty_Governance, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Governance - Add", // DisplayName
		2,                          // DisplayOrder
		Module,                     // Module
		ModuleDisplayOrder,         //ModuleDisplayOrder
		Level1,                     // Level1
		Level1DisplayOrder,         // Level1DisplayOrder
		Level2,                     // Level2
		Level2DisplayOrder,         // Level2DisplayOrder
		"",                         // Level3
		0,                          // Level3DisplayOrder
		true,                       //AllowedFor_OKAPI
		true,                       //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1
	r = Route{
		"HTTP_Loyalty_Governance",
		"PUT",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Governance/",
		Use(UC.HTTP_Loyalty_Governance, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Governance - Edit", // DisplayName
		3,                           // DisplayOrder
		Module,                      // Module
		ModuleDisplayOrder,          //ModuleDisplayOrder
		Level1,                      // Level1
		Level1DisplayOrder,          // Level1DisplayOrder
		Level2,                      // Level2
		Level2DisplayOrder,          // Level2DisplayOrder
		"",                          // Level3
		0,                           // Level3DisplayOrder
		true,                        //AllowedFor_OKAPI
		true,                        //AllowedFor_App
	}
	*R = append(*R, r)
	r = Route{
		"HTTP_Loyalty_Governance",
		"DELETE",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Governance/{KeyDelete}/",
		Use(UC.HTTP_Loyalty_Governance, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Governance - Delete", // DisplayName
		4,                             // DisplayOrder
		Module,                        // Module
		ModuleDisplayOrder,            //ModuleDisplayOrder
		Level1,                        // Level1
		Level1DisplayOrder,            // Level1DisplayOrder
		Level2,                        // Level2
		Level2DisplayOrder,            // Level2DisplayOrder
		"",                            // Level3
		0,                             // Level3DisplayOrder
		true,                          //AllowedFor_OKAPI
		true,                          //AllowedFor_App
	}
	*R = append(*R, r)
	//*****************************
	// Loyalty_Level
	//*****************************
	Level2 = "Loyalty Level"
	Level2DisplayOrder = 2
	r = Route{
		"HTTP_Loyalty_Level",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Level/",
		Use(UC.HTTP_Loyalty_Level, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Level - Read", // DisplayName
		1,                      // DisplayOrder
		Module,                 // Module
		ModuleDisplayOrder,     //ModuleDisplayOrder
		Level1,                 // Level1
		Level1DisplayOrder,     // Level1DisplayOrder
		Level2,                 // Level2
		Level2DisplayOrder,     // Level2DisplayOrder
		"",                     // Level3
		0,                      // Level3DisplayOrder
		true,                   //AllowedFor_OKAPI
		true,                   //AllowedFor_App
	}
	*R = append(*R, r)

	r = Route{
		"HTTP_Loyalty_Level",
		"POST",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Level/",
		Use(UC.HTTP_Loyalty_Level, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Level - Add", // DisplayName
		2,                     // DisplayOrder
		Module,                // Module
		ModuleDisplayOrder,    //ModuleDisplayOrder
		Level1,                // Level1
		Level1DisplayOrder,    // Level1DisplayOrder
		Level2,                // Level2
		Level2DisplayOrder,    // Level2DisplayOrder
		"",                    // Level3
		0,                     // Level3DisplayOrder
		true,                  //AllowedFor_OKAPI
		true,                  //AllowedFor_App
	}
	*R = append(*R, r)

	r = Route{
		"HTTP_Loyalty_Level",
		"PUT",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Level/",
		Use(UC.HTTP_Loyalty_Level, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Level - Edit", // DisplayName
		3,                      // DisplayOrder
		Module,                 // Module
		ModuleDisplayOrder,     //ModuleDisplayOrder
		Level1,                 // Level1
		Level1DisplayOrder,     // Level1DisplayOrder
		Level2,                 // Level2
		Level2DisplayOrder,     // Level2DisplayOrder
		"",                     // Level3
		0,                      // Level3DisplayOrder
		true,                   //AllowedFor_OKAPI
		true,                   //AllowedFor_App
	}
	*R = append(*R, r)

	r = Route{
		"HTTP_Loyalty_Level",
		"DELETE",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Level/{KeyDelete}/",
		Use(UC.HTTP_Loyalty_Level, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Level - Delete", // DisplayName
		4,                        // DisplayOrder
		Module,                   // Module
		ModuleDisplayOrder,       //ModuleDisplayOrder
		Level1,                   // Level1
		Level1DisplayOrder,       // Level1DisplayOrder
		Level2,                   // Level2
		Level2DisplayOrder,       // Level2DisplayOrder
		"",                       // Level3
		0,                        // Level3DisplayOrder
		true,                     //AllowedFor_OKAPI
		true,                     //AllowedFor_App
	}
	*R = append(*R, r)

	//*****************************
	// Loyalty_Account_Segment
	//*****************************
	Level2 = "Loyalty Account Segment"
	Level2DisplayOrder = 3

	r = Route{
		"HTTP_Loyalty_Account_Segment",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Account_Segment/",
		Use(UC.HTTP_Loyalty_Account_Segment, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Account Segment - Read", // DisplayName
		1,                                // DisplayOrder
		Module,                           // Module
		ModuleDisplayOrder,               //ModuleDisplayOrder
		Level1,                           // Level1
		Level1DisplayOrder,               // Level1DisplayOrder
		Level2,                           // Level2
		Level2DisplayOrder,               // Level2DisplayOrder
		"",                               // Level3
		0,                                // Level3DisplayOrder
		true,                             //AllowedFor_OKAPI
		true,                             //AllowedFor_App
	}
	*R = append(*R, r)

	r = Route{
		"HTTP_Loyalty_Account_Segment",
		"POST",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Account_Segment/",
		Use(UC.HTTP_Loyalty_Account_Segment, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Account Segment - Add", // DisplayName
		2,                               // DisplayOrder
		Module,                          // Module
		ModuleDisplayOrder,              //ModuleDisplayOrder
		Level1,                          // Level1
		Level1DisplayOrder,              // Level1DisplayOrder
		Level2,                          // Level2
		Level2DisplayOrder,              // Level2DisplayOrder
		"",                              // Level3
		0,                               // Level3DisplayOrder
		true,                            //AllowedFor_OKAPI
		true,                            //AllowedFor_App
	}
	*R = append(*R, r)

	r = Route{
		"HTTP_Loyalty_Account_Segment",
		"PUT",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Account_Segment/",
		Use(UC.HTTP_Loyalty_Account_Segment, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Account Segment - Edit", // DisplayName
		3,                                // DisplayOrder
		Module,                           // Module
		ModuleDisplayOrder,               //ModuleDisplayOrder
		Level1,                           // Level1
		Level1DisplayOrder,               // Level1DisplayOrder
		Level2,                           // Level2
		Level2DisplayOrder,               // Level2DisplayOrder
		"",                               // Level3
		0,                                // Level3DisplayOrder
		true,                             //AllowedFor_OKAPI
		true,                             //AllowedFor_App
	}
	*R = append(*R, r)

	r = Route{
		"HTTP_Loyalty_Account_Segment",
		"DELETE",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Account_Segment/{KeyDelete}/",
		Use(UC.HTTP_Loyalty_Account_Segment, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Account Segment - Delete", // DisplayName
		4,                                  // DisplayOrder
		Module,                             // Module
		ModuleDisplayOrder,                 //ModuleDisplayOrder
		Level1,                             // Level1
		Level1DisplayOrder,                 // Level1DisplayOrder
		Level2,                             // Level2
		Level2DisplayOrder,                 // Level2DisplayOrder
		"",                                 // Level3
		0,                                  // Level3DisplayOrder
		true,                               //AllowedFor_OKAPI
		true,                               //AllowedFor_App
	}
	*R = append(*R, r)

	//*****************************
	// Loyalty_Point_Earning_Rules
	//*****************************
	Level2 = "Loyalty Point Earning Rules"
	Level2DisplayOrder = 4
	r = Route{
		"HTTP_Loyalty_Point_Earning_Rules",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Point_Earning_Rules/",
		Use(UC.HTTP_Loyalty_Point_Earning_Rules, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Point Earning Rules - Read", // DisplayName
		1,                                    // DisplayOrder
		Module,                               // Module
		ModuleDisplayOrder,                   //ModuleDisplayOrder
		Level1,                               // Level1
		Level1DisplayOrder,                   // Level1DisplayOrder
		"",                                   // Level2
		0,                                    // Level2DisplayOrder
		"",                                   // Level3
		0,                                    // Level3DisplayOrder
		true,                                 //AllowedFor_OKAPI
		true,                                 //AllowedFor_App
	}
	*R = append(*R, r)

	r = Route{
		"HTTP_Loyalty_Point_Earning_Rules",
		"POST",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Point_Earning_Rules/",
		Use(UC.HTTP_Loyalty_Point_Earning_Rules, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Point Earning Rules - Add", // DisplayName
		2,                                   // DisplayOrder
		Module,                              // Module
		ModuleDisplayOrder,                  //ModuleDisplayOrder
		Level1,                              // Level1
		Level1DisplayOrder,                  // Level1DisplayOrder
		"",                                  // Level2
		0,                                   // Level2DisplayOrder
		"",                                  // Level3
		0,                                   // Level3DisplayOrder
		true,                                //AllowedFor_OKAPI
		true,                                //AllowedFor_App
	}
	*R = append(*R, r)

	r = Route{
		"HTTP_Loyalty_Point_Earning_Rules",
		"PUT",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Point_Earning_Rules/",
		Use(UC.HTTP_Loyalty_Point_Earning_Rules, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Point Earning Rules - Edit", // DisplayName
		3,                                    // DisplayOrder
		Module,                               // Module
		ModuleDisplayOrder,                   //ModuleDisplayOrder
		Level1,                               // Level1
		Level1DisplayOrder,                   // Level1DisplayOrder
		"",                                   // Level2
		0,                                    // Level2DisplayOrder
		"",                                   // Level3
		0,                                    // Level3DisplayOrder
		true,                                 //AllowedFor_OKAPI
		true,                                 //AllowedFor_App
	}
	*R = append(*R, r)

	r = Route{
		"HTTP_Loyalty_Point_Earning_Rules",
		"DELETE",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Point_Earning_Rules/{KeyDelete}/",
		Use(UC.HTTP_Loyalty_Point_Earning_Rules, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Point Earning Rules - Delete", // DisplayName
		4,                                      // DisplayOrder
		Module,                                 // Module
		ModuleDisplayOrder,                     //ModuleDisplayOrder
		Level1,                                 // Level1
		Level1DisplayOrder,                     // Level1DisplayOrder
		"",                                     // Level2
		0,                                      // Level2DisplayOrder
		"",                                     // Level3
		0,                                      // Level3DisplayOrder
		true,                                   //AllowedFor_OKAPI
		true,                                   //AllowedFor_App
	}
	*R = append(*R, r)

	//*****************************
	// Loyalty_Point_Expiry_Rules
	//*****************************
	Level2 = "Loyalty Point Expiry Rules"
	Level2DisplayOrder = 5
	r = Route{
		"HTTP_Loyalty_Point_Expiry_Rules",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Point_Expiry_Rules/",
		Use(UC.HTTP_Loyalty_Point_Expiry_Rules, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Point Expiry Rules - Read", // DisplayName
		1,                                   // DisplayOrder
		Module,                              // Module
		ModuleDisplayOrder,                  //ModuleDisplayOrder
		Level1,                              // Level1
		Level1DisplayOrder,                  // Level1DisplayOrder
		"",                                  // Level2
		0,                                   // Level2DisplayOrder
		"",                                  // Level3
		0,                                   // Level3DisplayOrder
		true,                                //AllowedFor_OKAPI
		true,                                //AllowedFor_App
	}
	*R = append(*R, r)

	r = Route{
		"HTTP_Loyalty_Point_Expiry_Rules",
		"POST",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Point_Expiry_Rules/",
		Use(UC.HTTP_Loyalty_Point_Expiry_Rules, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Point Expiry Rules - Add", // DisplayName
		2,                                  // DisplayOrder
		Module,                             // Module
		ModuleDisplayOrder,                 //ModuleDisplayOrder
		Level1,                             // Level1
		Level1DisplayOrder,                 // Level1DisplayOrder
		"",                                 // Level2
		0,                                  // Level2DisplayOrder
		"",                                 // Level3
		0,                                  // Level3DisplayOrder
		true,                               //AllowedFor_OKAPI
		true,                               //AllowedFor_App
	}
	*R = append(*R, r)

	r = Route{
		"HTTP_Loyalty_Point_Expiry_Rules",
		"PUT",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Point_Expiry_Rules/",
		Use(UC.HTTP_Loyalty_Point_Expiry_Rules, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Point Expiry Rules - Edit", // DisplayName
		3,                                   // DisplayOrder
		Module,                              // Module
		ModuleDisplayOrder,                  //ModuleDisplayOrder
		Level1,                              // Level1
		Level1DisplayOrder,                  // Level1DisplayOrder
		"",                                  // Level2
		0,                                   // Level2DisplayOrder
		"",                                  // Level3
		0,                                   // Level3DisplayOrder
		true,                                //AllowedFor_OKAPI
		true,                                //AllowedFor_App
	}
	*R = append(*R, r)

	r = Route{
		"HTTP_Loyalty_Point_Expiry_Rules",
		"DELETE",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Point_Expiry_Rules/{KeyDelete}/",
		Use(UC.HTTP_Loyalty_Point_Expiry_Rules, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Point Expiry Rules - Delete", // DisplayName
		4,                                     // DisplayOrder
		Module,                                // Module
		ModuleDisplayOrder,                    //ModuleDisplayOrder
		Level1,                                // Level1
		Level1DisplayOrder,                    // Level1DisplayOrder
		"",                                    // Level2
		0,                                     // Level2DisplayOrder
		"",                                    // Level3
		0,                                     // Level3DisplayOrder
		true,                                  //AllowedFor_OKAPI
		true,                                  //AllowedFor_App
	}
	*R = append(*R, r)

	//*****************************
	// Loyalty_Point_Redemption_Rules
	//*****************************
	Level2 = "Loyalty Point Redemption Rules"
	Level2DisplayOrder = 6
	r = Route{
		"HTTP_Loyalty_Point_Redemption_Rules",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Point_Redemption_Rules/",
		Use(UC.HTTP_Loyalty_Point_Redemption_Rules, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Point Redemption Rules - Read", // DisplayName
		1,                                       // DisplayOrder
		Module,                                  // Module
		ModuleDisplayOrder,                      //ModuleDisplayOrder
		Level1,                                  // Level1
		Level1DisplayOrder,                      // Level1DisplayOrder
		"",                                      // Level2
		0,                                       // Level2DisplayOrder
		"",                                      // Level3
		0,                                       // Level3DisplayOrder
		true,                                    //AllowedFor_OKAPI
		true,                                    //AllowedFor_App
	}
	*R = append(*R, r)

	r = Route{
		"HTTP_Loyalty_Point_Redemption_Rules",
		"POST",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Point_Redemption_Rules/",
		Use(UC.HTTP_Loyalty_Point_Redemption_Rules, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Point Redemption Rules - Add", // DisplayName
		2,                                      // DisplayOrder
		Module,                                 // Module
		ModuleDisplayOrder,                     //ModuleDisplayOrder
		Level1,                                 // Level1
		Level1DisplayOrder,                     // Level1DisplayOrder
		"",                                     // Level2
		0,                                      // Level2DisplayOrder
		"",                                     // Level3
		0,                                      // Level3DisplayOrder
		true,                                   //AllowedFor_OKAPI
		true,                                   //AllowedFor_App
	}
	*R = append(*R, r)

	r = Route{
		"HTTP_Loyalty_Point_Redemption_Rules",
		"PUT",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Point_Redemption_Rules/",
		Use(UC.HTTP_Loyalty_Point_Redemption_Rules, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Point Redemption Rules - Edit", // DisplayName
		3,                                       // DisplayOrder
		Module,                                  // Module
		ModuleDisplayOrder,                      //ModuleDisplayOrder
		Level1,                                  // Level1
		Level1DisplayOrder,                      // Level1DisplayOrder
		"",                                      // Level2
		0,                                       // Level2DisplayOrder
		"",                                      // Level3
		0,                                       // Level3DisplayOrder
		true,                                    //AllowedFor_OKAPI
		true,                                    //AllowedFor_App
	}
	*R = append(*R, r)

	r = Route{
		"HTTP_Loyalty_Point_Redemption_Rules",
		"DELETE",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Point_Redemption_Rules/{KeyDelete}/",
		Use(UC.HTTP_Loyalty_Point_Redemption_Rules, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Point Redemption Rules - Delete", // DisplayName
		4,                  // DisplayOrder
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

	//*****************************
	// Loyalty_Plan
	//*****************************
	Level2 = "Loyalty Plan"
	Level2DisplayOrder = 7
	r = Route{
		"HTTP_Loyalty_Plan",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Plan/",
		Use(UC.HTTP_Loyalty_Plan, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Plan - Read", // DisplayName
		1,                     // DisplayOrder
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

	r = Route{
		"HTTP_Loyalty_Plan",
		"POST",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Plan/",
		Use(UC.HTTP_Loyalty_Plan, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Plan - Add", // DisplayName
		2,                    // DisplayOrder
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

	r = Route{
		"HTTP_Loyalty_Plan",
		"PUT",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Plan/",
		Use(UC.HTTP_Loyalty_Plan, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Plan - Edit", // DisplayName
		3,                     // DisplayOrder
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

	r = Route{
		"HTTP_Loyalty_Plan",
		"DELETE",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Plan/{KeyDelete}/",
		Use(UC.HTTP_Loyalty_Plan, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Plan - Delete", // DisplayName
		4,                       // DisplayOrder
		Module,                  // Module
		ModuleDisplayOrder,      //ModuleDisplayOrder
		Level1,                  // Level1
		Level1DisplayOrder,      // Level1DisplayOrder
		"",                      // Level2
		0,                       // Level2DisplayOrder
		"",                      // Level3
		0,                       // Level3DisplayOrder
		true,                    //AllowedFor_OKAPI
		true,                    //AllowedFor_App
	}
	*R = append(*R, r)

	//*****************************
	// Customer_Loyalty_Account
	//*****************************
	Level2 = "Customer Loyalty Account"
	Level2DisplayOrder = 8
	r = Route{
		"HTTP_Customer_Loyalty_Account",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_Loyalty_Account/",
		Use(UC.HTTP_Customer_Loyalty_Account, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Customer Loyalty Account - Read", // DisplayName
		1,                                 // DisplayOrder
		Module,                            // Module
		ModuleDisplayOrder,                //ModuleDisplayOrder
		Level1,                            // Level1
		Level1DisplayOrder,                // Level1DisplayOrder
		"",                                // Level2
		0,                                 // Level2DisplayOrder
		"",                                // Level3
		0,                                 // Level3DisplayOrder
		true,                              //AllowedFor_OKAPI
		true,                              //AllowedFor_App
	}
	*R = append(*R, r)

	r = Route{
		"HTTP_Customer_Loyalty_Account",
		"POST",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_Loyalty_Account/",
		Use(UC.HTTP_Customer_Loyalty_Account, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Customer Loyalty Account - Add", // DisplayName
		2,                                // DisplayOrder
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

	r = Route{
		"HTTP_Customer_Loyalty_Account",
		"PUT",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_Loyalty_Account/",
		Use(UC.HTTP_Customer_Loyalty_Account, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Customer Loyalty Account - Edit", // DisplayName
		3,                                 // DisplayOrder
		Module,                            // Module
		ModuleDisplayOrder,                //ModuleDisplayOrder
		Level1,                            // Level1
		Level1DisplayOrder,                // Level1DisplayOrder
		"",                                // Level2
		0,                                 // Level2DisplayOrder
		"",                                // Level3
		0,                                 // Level3DisplayOrder
		true,                              //AllowedFor_OKAPI
		true,                              //AllowedFor_App
	}
	*R = append(*R, r)

	r = Route{
		"HTTP_Customer_Loyalty_Account",
		"DELETE",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_Loyalty_Account/{KeyDelete}/",
		Use(UC.HTTP_Customer_Loyalty_Account, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Customer Loyalty Account - Delete", // DisplayName
		4,                                   // DisplayOrder
		Module,                              // Module
		ModuleDisplayOrder,                  //ModuleDisplayOrder
		Level1,                              // Level1
		Level1DisplayOrder,                  // Level1DisplayOrder
		"",                                  // Level2
		0,                                   // Level2DisplayOrder
		"",                                  // Level3
		0,                                   // Level3DisplayOrder
		true,                                //AllowedFor_OKAPI
		true,                                //AllowedFor_App
	}
	*R = append(*R, r)

	r = Route{
		"HTTP_Customer_Loyalty_Account_Points_Details",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_Loyalty_Account_Points_Details/",
		Use(UC.HTTP_Customer_Loyalty_Account_Points_Details, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Customer Loyalty Account Points Details - Read", // DisplayName
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
		"HTTP_Customer_Loyalty_Account_DebitPoints",
		"PUT",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_Loyalty_Account_DebitPoints/",
		Use(UC.HTTP_Customer_Loyalty_Account_DebitPoints, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		Level1 + " - Debit Points", // DisplayName
		DisplayOrder,               // DisplayOrder
		Module,                     // Module
		ModuleDisplayOrder,         //ModuleDisplayOrder
		Level1,                     // Level1
		Level1DisplayOrder,         // Level1DisplayOrder
		"",                         // Level2
		0,                          // Level2DisplayOrder
		"",                         // Level3
		0,                          // Level3DisplayOrder
		true,                       //AllowedFor_OKAPI
		true,                       //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Customer_Loyalty_Account_CreditPoints",
		"PUT",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_Loyalty_Account_CreditPoints/",
		Use(UC.HTTP_Customer_Loyalty_Account_CreditPoints, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		Level1 + " - Credit Points", // DisplayName
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

	r = Route{
		"HTTP_Customer_Loyalty_RedeemRequest",
		"PUT",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_Loyalty_RedeemRequest/",
		Use(UC.HTTP_Customer_Loyalty_RedeemRequest, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		Level1 + " - Redeem Request", // DisplayName
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
		"HTTP_Customer_Loyalty_Account_GetRedemption_Rules",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_Loyalty_Account_GetRedemption_Rules/",
		Use(UC.HTTP_Customer_Loyalty_Account_GetRedemption_Rules, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		Level1 + " - Get Redemption Rules", // DisplayName
		DisplayOrder,                       // DisplayOrder
		Module,                             // Module
		ModuleDisplayOrder,                 //ModuleDisplayOrder
		Level1,                             // Level1
		Level1DisplayOrder,                 // Level1DisplayOrder
		"",                                 // Level2
		0,                                  // Level2DisplayOrder
		"",                                 // Level3
		0,                                  // Level3DisplayOrder
		true,                               //AllowedFor_OKAPI
		true,                               //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	//*****************************
	// Customer_UAT
	//*****************************
	Level1 = "Customer UAT"
	r = Route{
		"HTTP_Customer_UAT",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_UAT/",
		Use(UC.HTTP_Customer_UAT, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Customer UAT - Read", // DisplayName
		1,                     // DisplayOrder
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

	r = Route{
		"HTTP_Customer_UAT",
		"POST",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_UAT/",
		Use(UC.HTTP_Customer_UAT, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Customer UAT  - Add", // DisplayName
		2,                     // DisplayOrder
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

	r = Route{
		"HTTP_Customer_UAT",
		"PUT",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_UAT/",
		Use(UC.HTTP_Customer_UAT, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Customer UAT  - Edit", // DisplayName
		3,                      // DisplayOrder
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

	r = Route{
		"HTTP_Customer_UAT",
		"DELETE",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_UAT/{KeyDelete}/",
		Use(UC.HTTP_Customer_UAT, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Customer UAT  - Delete", // DisplayName
		4,                        // DisplayOrder
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

	//*****************************
	// Customer_DND
	//*****************************
	Level1 = "Customer DND"
	Level1DisplayOrder = Level1DisplayOrder + 1
	r = Route{
		"HTTP_Customer_DND",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_DND/",
		Use(UC.HTTP_Customer_DND, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Customer DND - Read", // DisplayName
		1,                     // DisplayOrder
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

	r = Route{
		"HTTP_Customer_DND",
		"POST",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_DND/",
		Use(UC.HTTP_Customer_DND, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Customer DND - Add", // DisplayName
		2,                    // DisplayOrder
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

	r = Route{
		"HTTP_Customer_DND",
		"PUT",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_DND/",
		Use(UC.HTTP_Customer_DND, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Customer DND - Edit", // DisplayName
		3,                     // DisplayOrder
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

	r = Route{
		"HTTP_Customer_DND",
		"DELETE",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_DND/{KeyDelete}/",
		Use(UC.HTTP_Customer_DND, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Customer DND - Delete", // DisplayName
		4,                       // DisplayOrder
		Module,                  // Module
		ModuleDisplayOrder,      //ModuleDisplayOrder
		Level1,                  // Level1
		Level1DisplayOrder,      // Level1DisplayOrder
		"",                      // Level2
		0,                       // Level2DisplayOrder
		"",                      // Level3
		0,                       // Level3DisplayOrder
		true,                    //AllowedFor_OKAPI
		true,                    //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	//*****************************
	// Customer_Exclusion
	//*****************************
	Level1 = "Customer Exclusion"
	Level1DisplayOrder = Level1DisplayOrder + 1
	r = Route{
		"HTTP_Customer_Exclusion",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_Exclusion/",
		Use(UC.HTTP_Customer_Exclusion, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Customer Exclusion - Read", // DisplayName
		1,                           // DisplayOrder
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

	r = Route{
		"HTTP_Customer_Exclusion",
		"POST",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_Exclusion/",
		Use(UC.HTTP_Customer_Exclusion, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Customer Exclusion - Add", // DisplayName
		2,                          // DisplayOrder
		Module,                     // Module
		ModuleDisplayOrder,         //ModuleDisplayOrder
		Level1,                     // Level1
		Level1DisplayOrder,         // Level1DisplayOrder
		"",                         // Level2
		0,                          // Level2DisplayOrder
		"",                         // Level3
		0,                          // Level3DisplayOrder
		true,                       //AllowedFor_OKAPI
		true,                       //AllowedFor_App
	}
	*R = append(*R, r)

	r = Route{
		"HTTP_Customer_Exclusion",
		"PUT",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_Exclusion/",
		Use(UC.HTTP_Customer_Exclusion, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Customer Exclusion - Edit", // DisplayName
		3,                           // DisplayOrder
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

	r = Route{
		"HTTP_Customer_Exclusion",
		"DELETE",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_Exclusion/{KeyDelete}/",
		Use(UC.HTTP_Customer_Exclusion, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Customer Exclusion - Delete", // DisplayName
		4,                             // DisplayOrder
		Module,                        // Module
		ModuleDisplayOrder,            //ModuleDisplayOrder
		Level1,                        // Level1
		Level1DisplayOrder,            // Level1DisplayOrder
		"",                            // Level2
		0,                             // Level2DisplayOrder
		"",                            // Level3
		0,                             // Level3DisplayOrder
		true,                          //AllowedFor_OKAPI
		true,                          //AllowedFor_App
	}
	*R = append(*R, r)

	//*****************************
	// Customer_COS_Exclusion
	//*****************************
	Level1 = "Customer COS Exclusion"
	Level1DisplayOrder = Level1DisplayOrder + 1
	r = Route{
		"HTTP_Customer_COS_Exclusion",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_COS_Exclusion/",
		Use(UC.HTTP_Customer_COS_Exclusion, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Customer COS Exclusion - Read", // DisplayName
		1,                               // DisplayOrder
		Module,                          // Module
		ModuleDisplayOrder,              //ModuleDisplayOrder
		Level1,                          // Level1
		Level1DisplayOrder,              // Level1DisplayOrder
		"",                              // Level2
		0,                               // Level2DisplayOrder
		"",                              // Level3
		0,                               // Level3DisplayOrder
		true,                            //AllowedFor_OKAPI
		true,                            //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Customer_COS_Exclusion",
		"POST",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_COS_Exclusion/",
		Use(UC.HTTP_Customer_COS_Exclusion, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Customer COS Exclusion - Add", // DisplayName
		2,                              // DisplayOrder
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

	r = Route{
		"HTTP_Customer_COS_Exclusion",
		"PUT",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_COS_Exclusion/",
		Use(UC.HTTP_Customer_COS_Exclusion, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Customer COS Exclusion - Edit", // DisplayName
		3,                               // DisplayOrder
		Module,                          // Module
		ModuleDisplayOrder,              //ModuleDisplayOrder
		Level1,                          // Level1
		Level1DisplayOrder,              // Level1DisplayOrder
		"",                              // Level2
		0,                               // Level2DisplayOrder
		"",                              // Level3
		0,                               // Level3DisplayOrder
		true,                            //AllowedFor_OKAPI
		true,                            //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	r = Route{
		"HTTP_Customer_COS_Exclusion",
		"DELETE",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Customer_COS_Exclusion/{KeyDelete}/",
		Use(UC.HTTP_Customer_COS_Exclusion, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Customer COS Exclusion - Delete", // DisplayName
		4,                                 // DisplayOrder
		Module,                            // Module
		ModuleDisplayOrder,                //ModuleDisplayOrder
		Level1,                            // Level1
		Level1DisplayOrder,                // Level1DisplayOrder
		"",                                // Level2
		0,                                 // Level2DisplayOrder
		"",                                // Level3
		0,                                 // Level3DisplayOrder
		true,                              //AllowedFor_OKAPI
		true,                              //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

	Level1 = "Loyalty Products Catalogue"
	Level1DisplayOrder = Level1DisplayOrder + 1
	r = Route{
		"HTTP_Loyalty_Products_Catalogue",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_Products_Catalogue/",
		Use(UC.HTTP_Loyalty_Products_Catalogue, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
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

	Level1 = "Loyalty Account Debit Points log"
	Level1DisplayOrder = Level1DisplayOrder + 1
	r = Route{
		"HTTP_Loyalty_AccountDebitPoints_log",
		"GET",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_Loyalty_AccountDebitPoints_log/",
		Use(UC.HTTP_Loyalty_AccountDebitPoints_log, UC.ValidateAccess_AUC, UC.ValidateJWEToken),
		true,
		"Loyalty Account Debit Points log - Read", // DisplayName
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

func (UC *UserControl) Add_LoyaltyFeedRoutes(R *Routes) {
	var r Route
	var Module, Level1 string
	var ModuleDisplayOrder int64
	Module = Configuration.LoyaltyModule
	ModuleDisplayOrder = 1
	var Level1DisplayOrder int64
	var DisplayOrder int64 = 1

	//*****************************
	// Receive IN Live feed
	//*****************************
	Level1 = "Receive IN Live feed"
	Level1DisplayOrder = Level1DisplayOrder + 1
	r = Route{
		"HTTP_INLiveFeed_NewJoining",
		"POST",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_INLiveFeed_NewJoining/",
		UC.HTTP_INLiveFeed_NewJoining,
		true,
		Level1 + " - NewJoining", // DisplayName
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
		"HTTP_INLiveFeed_Churn",
		"DELETE",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_INLiveFeed_Churn/",
		UC.HTTP_INLiveFeed_Churn,
		true,
		Level1 + " - Churn", // DisplayName
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
		"HTTP_INLiveFeed_Consuption",
		"POST",
		"/" + Configuration.LoyaltyModule + "/" + Configuration.LoyaltyVersion + "/HTTP_INLiveFeed_Consuption/",
		UC.HTTP_INLiveFeed_Consuption,
		true,
		Level1 + " - Consuption", // DisplayName
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
}
