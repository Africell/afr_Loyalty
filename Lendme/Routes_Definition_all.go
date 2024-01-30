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
		false,                       //AllowedFor_App
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
		false,                        //AllowedFor_App
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
		false,                          //AllowedFor_App
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
		false,              //AllowedFor_App
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
		false,               //AllowedFor_App
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
		false,                 //AllowedFor_App
	}
	*R = append(*R, r)
	DisplayOrder = DisplayOrder + 1

}
