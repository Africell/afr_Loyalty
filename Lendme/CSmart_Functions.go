package Lendme

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var CDS_TOKEN = "eyJhbGciOiJIUzUxMiJ9.eyJ0ZW5hbnRUeXBlIjoicGFydG5lciIsInRlbmFuY3lJZCI6Ik1WTk9fQUZSSUNFTExDTElFTlRfQU5HX0xBTkRNRSIsImVudmlyb25tZW50IjoicHJvZCIsInRva2VuVHlwZSI6ImVudlRva2VuIiwidmVyc2lvbiI6MSwiaWF0IjoxNzQxOTUxNTY0fQ.FVtxtuPzSqXKLo44-e8kplmfzMJLwPScsJyoTjHGXyNoHYSvKqj4X0o71oVd5JXOIVLrTrHK0OGA6C6lMcwNwg"
var CDS_Endpoint = "http://10.250.8.245:8080"        //"http://102.218.87.225:8080"  //before improvements
var CDS_Endpoint_Etopup = "http://10.250.8.245:8081" //after improvements
var Source = "LandMe"

// ----------------------------------------------------------------------------------------
// CSMART -- Generic error reply structure
// ----------------------------------------------------------------------------------------
type CS_APIReplyError struct {
	Response CS_Err_Response `bson:"response" json:"response"`
}
type CS_Err_Response struct {
	Success string        `bson:"success" json:"success"`
	Result  CS_Err_Result `bson:"result" json:"result"`
}
type CS_Err_Result struct {
	Arguments CS_Err_Arguments `bson:"arguments" json:"arguments"`
}
type CS_Err_Arguments struct {
	StatusCode   string `bson:"statusCode" json:"statusCode"`
	ErrorCode    string `bson:"errorCode" json:"errorCode"`
	ErrorMessage string `bson:"errorMessage" json:"errorMessage"`
}

// ----------------------------------------------------------------------------------------
// CSMART -- Get Account Balance
// ----------------------------------------------------------------------------------------
type CS_Balance_Response struct {
	Response struct {
		Success string `json:"success"`
		Result  struct {
			Arguments struct {
				BalanceGroups []struct {
					Elem        int    `json:"elem"`
					BillinfoObj string `json:"BILLINFO_OBJ"`
					Balances    []struct {
						Elem        int     `json:"elem"`
						ValidTo     int     `json:"VALID_TO"`
						CreditFloor string  `json:"CREDIT_FLOOR"`
						CreditLimit string  `json:"CREDIT_LIMIT"`
						ResourceID  int     `json:"RESOURCE_ID"`
						TypeStr     string  `json:"TYPE_STR"`
						CurrentBal  float64 `json:"CURRENT_BAL"`
						Info        string  `json:"INFO"`
						ValidFrom   int     `json:"VALID_FROM"`
						AmountOrig  int     `json:"AMOUNT_ORIG"`
						Descr       string  `json:"DESCR"`
					} `json:"BALANCES"`
					BalGrpObj string `json:"BAL_GRP_OBJ"`
				} `json:"BALANCE_GROUPS"`
				ErrorNum int    `json:"ERROR_NUM"`
				Poid     string `json:"POID"`
			} `json:"arguments"`
		} `json:"result"`
	} `json:"response"`
}

func CS_GetAccountBalance_ByDeviceId(MSISDN string) (Balance_Response CS_Balance_Response, err error) {
	url := CDS_Endpoint + "/dx-sync/balance?deviceId=" + MSISDN
	//log.Println("url: ", url)
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("srdate", time.Now().String())
	req.Header.Set("operation", "balance")
	req.Header.Set("correlationId", uuid.New().String())
	req.Header.Set("serviceProvider", Source)
	req.Header.Set("source", Source)
	req.Header.Set("destination", "CRM")
	req.Header.Set("token", CDS_TOKEN)
	//client := &http.Client{}
	client := &http.Client{
		Timeout: 4 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error Connecting to Covalense API -- CS_dealerLogin: ", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	// fmt.Println("here ", string(body))
	// fmt.Println("status ", resp.Status)
	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(body, &Balance_Response)
		if err != nil {
			fmt.Println("Error parsing Covalense API reply -- CS_dealerLogin: ", err, string(body))
		}
	} else {
		var error_reply CS_APIReplyError
		err = json.Unmarshal(body, &error_reply)
		if err != nil {
			fmt.Println("Error parsing Covalense API reply -- CS_dealerLogin: ", err, string(body))
		}
		err = errors.New(error_reply.Response.Result.Arguments.ErrorMessage)
	}
	return Balance_Response, err
}

type CS_Dealer_Balance_Response struct {
	Response struct {
		Success string `json:"success"`
		Result  struct {
			Arguments struct {
				BalanceGroups []struct {
					Elem     int `json:"elem"`
					Balances []struct {
						Elem         int     `json:"elem"`
						StartTString string  `json:"startTString"`
						EndTString   string  `json:"endTString"`
						CreditFloor  string  `json:"CREDIT_FLOOR"`
						CreditLimit  string  `json:"CREDIT_LIMIT"`
						ResourceID   int     `json:"RESOURCE_ID"`
						TypeStr      string  `json:"TYPE_STR"`
						CurrentBal   float64 `json:"CURRENT_BAL"`
						Info         string  `json:"INFO"`
						AmountOrig   int     `json:"AMOUNT_ORIG"`
						Descr        string  `json:"DESCR"`
						EndT         int     `json:"END_T"`
						StartT       int     `json:"START_T"`
						ValidTo      string  `json:"VALID_TO"`
						ValidFrom    string  `json:"VALID_FROM"`
					} `json:"BALANCES"`
					BalGrpObj   string `json:"BAL_GRP_OBJ"`
					BillinfoObj string `json:"BILLINFO_OBJ"`
				} `json:"BALANCE_GROUPS"`
				ErrorNum int    `json:"ERROR_NUM"`
				Poid     string `json:"POID"`
			} `json:"arguments"`
		} `json:"result"`
	} `json:"response"`
}

func CS_GetAccountBalance_ByAccountNo(AccountNo string) (Balance_Response CS_Dealer_Balance_Response, err error) {
	url := CDS_Endpoint + "/dx-sync/balance?accountNo=" + AccountNo
	log.Println("url: ", url)
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("srdate", time.Now().String())
	req.Header.Set("operation", "balance")
	req.Header.Set("correlationId", uuid.New().String())
	req.Header.Set("serviceProvider", Source)
	req.Header.Set("source", Source)
	req.Header.Set("destination", "CRM")
	req.Header.Set("token", CDS_TOKEN)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error Connecting to Covalense API -- CS_dealerLogin: ", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	log.Println("here ", string(body))
	log.Println("status ", resp.Status)
	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(body, &Balance_Response)
		if err != nil {
			fmt.Println("Error parsing Covalense API reply -- CS_dealerLogin: ", err, string(body))
		}
	} else {
		var error_reply CS_APIReplyError
		err = json.Unmarshal(body, &error_reply)
		if err != nil {
			fmt.Println("Error parsing Covalense API reply -- CS_dealerLogin: ", err, string(body))
		}
		err = errors.New(error_reply.Response.Result.Arguments.ErrorMessage)
	}
	return Balance_Response, err
}

// ----------------------------------------------------------------------------------------
// CSMART -- Dealer Balance Transfer
// ----------------------------------------------------------------------------------------
// type CS_DealerBalanceTransfer_Response struct {
// 	Response struct {
// 		Success string `json:"success"`
// 		Result  struct {
// 			Success bool   `json:"success"`
// 			Record  string `json:"record"`
// 			Message string `json:"message"`
// 		} `json:"result"`
// 	} `json:"response"`
// }

type CS_DealerBalanceTransfer_Failed_Response struct {
	Response struct {
		Success string `json:"success"`
		Result  struct {
			Arguments struct {
				ErrorCode    string `json:"errorCode"`
				ErrorMessage string `json:"errorMessage"`
				StatusCode   string `json:"statusCode"`
			} `json:"arguments"`
		} `json:"result"`
	} `json:"response"`
}

// type CS_DealerBalanceTransfer_Success_Response struct {
// 	Response struct {
// 		Success string `json:"success"`
// 		Result  struct {
// 			Arguments interface{} `json:"arguments"`
// 			Message   string      `json:"message"`
// 			Success   string      `json:"success"`
// 			Record    string      `json:"record"`
// 		} `json:"result"`
// 	} `json:"response"`
// }

type CS_DealerBalanceTransfer_Success_Response struct {
	Response struct {
		Success string `json:"success"`
		Result  struct {
			Arguments struct {
				ProviderDescr                 string `json:"PROVIDER_DESCR"`
				AfricellFldDestinationAccount struct {
					Amount        string `json:"AMOUNT"`
					BalanceGroups []struct {
						Elem     string `json:"elem"`
						Balances []struct {
							Elem        string `json:"elem"`
							EndT        string `json:"END_T"`
							StartT      string `json:"START_T"`
							CreditFloor string `json:"CREDIT_FLOOR"`
							CreditLimit string `json:"CREDIT_LIMIT"`
							ResourceID  string `json:"RESOURCE_ID"`
							TypeStr     string `json:"TYPE_STR"`
							CurrentBal  string `json:"CURRENT_BAL"`
							Info        string `json:"INFO"`
							AmountOrig  string `json:"AMOUNT_ORIG"`
							Descr       string `json:"DESCR"`
						} `json:"BALANCES"`
						BalGrpObj string `json:"BAL_GRP_OBJ"`
					} `json:"BALANCE_GROUPS"`
				} `json:"AFRICELL_FLD_DESTINATION_ACCOUNT"`
				Amount                   string `json:"AMOUNT"`
				ProgramName              string `json:"PROGRAM_NAME"`
				AfricellFldSourceAccount struct {
					Amount        string `json:"AMOUNT"`
					BalanceGroups []struct {
						Elem     string `json:"elem"`
						Balances []struct {
							Elem        string `json:"elem"`
							EndT        string `json:"END_T"`
							StartT      string `json:"START_T"`
							CreditFloor string `json:"CREDIT_FLOOR"`
							CreditLimit string `json:"CREDIT_LIMIT"`
							ResourceID  string `json:"RESOURCE_ID"`
							TypeStr     string `json:"TYPE_STR"`
							CurrentBal  string `json:"CURRENT_BAL"`
							Info        string `json:"INFO"`
							AmountOrig  string `json:"AMOUNT_ORIG"`
							Descr       string `json:"DESCR"`
						} `json:"BALANCES"`
						BalGrpObj string `json:"BAL_GRP_OBJ"`
					} `json:"BALANCE_GROUPS"`
				} `json:"AFRICELL_FLD_SOURCE_ACCOUNT"`
				DealerCode                string `json:"DEALER_CODE"`
				EventObj                  string `json:"EVENT_OBJ"`
				CsmartFldSourceDealerCode string `json:"CSMART_FLD_SOURCE_DEALER_CODE"`
				TransID                   string `json:"TRANS_ID"`
				Descr                     string `json:"DESCR"`
				ErrorNum                  string `json:"ERROR_NUM"`
				Poid                      string `json:"POID"`
			} `json:"arguments"`
			Message string `json:"message"`
			Success string `json:"success"`
			Record  string `json:"record"`
		} `json:"result"`
	} `json:"response"`
}

func CS_DealerBalanceTransfer(transferType string, DealerMSISDN string, DealerPIN string, TargetMSISDN string, amount string) (DBT_Response CS_DealerBalanceTransfer_Success_Response, err error) {
	//url := CDS_Endpoint + "/dx-sync/dealer/balancetransfer"
	url := CDS_Endpoint_Etopup + "/csmart-brmservices/brm/v1/dealer/balance/transfer"
	log.Println("url: ", url)
	payload := map[string]string{
		"targetMsisdn": TargetMSISDN,
		"senderMsisdn": DealerMSISDN,
		"senderPin":    DealerPIN,
		"transferType": transferType, //"customer", //"account", //"dealer",
		"amount":       amount,
		"appSource":    Source,
	}
	log.Println("payload:", payload)
	_Json, _ := json.Marshal(payload)
	_jsonByte := []byte(_Json)
	method := "POST"
	req, err := http.NewRequest(method, url, bytes.NewBuffer(_jsonByte))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("srdate", time.Now().String())
	req.Header.Set("operation", "balanceTransfer")
	req.Header.Set("correlationId", uuid.New().String())
	req.Header.Set("serviceProvider", Source)
	req.Header.Set("source", Source)
	req.Header.Set("destination", "CRM")
	req.Header.Set("token", CDS_TOKEN)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error Connecting to Covalense API -- CS_DealerBalanceTransfer: ", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	log.Println("Request Body ", string(body))
	log.Println("status ", resp.Status)
	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(body, &DBT_Response)
		if err != nil {
			fmt.Println("Error parsing Covalense API reply -- CS_DealerBalanceTransfer: ", err, string(body))
		}
	} else {
		var error_reply CS_APIReplyError
		err = json.Unmarshal(body, &error_reply)
		if err != nil {
			fmt.Println("Error parsing Covalense API reply -- CS_DealerBalanceTransfer: ", err, string(body))
		}
		err = errors.New(error_reply.Response.Result.Arguments.ErrorMessage + " (ErrorCode: " + error_reply.Response.Result.Arguments.ErrorCode + ")")
	}
	return DBT_Response, err
}

// ----------------------------------------------------------------------------------------
// CSMART -- Validate PIN
// ----------------------------------------------------------------------------------------

type CS_ValidatePIN_Response struct {
	Response struct {
		Success string `json:"success"`
		Result  struct {
			DearlerNo string `json:"dearlerNo"`
			Message   string `json:"message"`
		} `json:"result"`
	} `json:"response"`
}

func CS_ValidatePIN(MSISDN string, PIN string) (GPD_Response CS_ValidatePIN_Response, err error) {
	url := CDS_Endpoint + "/dx-sync/validate/pin"
	log.Println("url: ", url)
	payload := map[string]string{
		"etopupmsisdn": MSISDN,
		"newPIN":       PIN,
	}
	log.Println("payload:", payload)
	_Json, _ := json.Marshal(payload)
	_jsonByte := []byte(_Json)
	method := "POST"
	req, err := http.NewRequest(method, url, bytes.NewBuffer(_jsonByte))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("srdate", time.Now().String())
	req.Header.Set("operation", "validateEtopupPIN")
	req.Header.Set("correlationId", uuid.New().String())
	req.Header.Set("serviceProvider", Source)
	req.Header.Set("source", Source)
	req.Header.Set("destination", "CRM")
	req.Header.Set("token", CDS_TOKEN)
	client := &http.Client{}
	log.Println("Validate PIN before client.do")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error Connecting to Covalense API -- CS_GetPackageDetails: ", err)
		return
	}
	log.Println("Validate PIN after client.do")
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	log.Println("Request Body ", string(body))
	log.Println("status ", resp.Status)
	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(body, &GPD_Response)
		if err != nil {
			fmt.Println("Error parsing Covalense API reply -- CS_GetPackageDetails: ", err, string(body))
		}
	} else {
		var error_reply CS_APIReplyError
		err = json.Unmarshal(body, &error_reply)
		if err != nil {
			fmt.Println("Error parsing Covalense API reply -- CS_GetPackageDetails: ", err, string(body))
		}
		err = errors.New(error_reply.Response.Result.Arguments.ErrorMessage)
	}
	return GPD_Response, err
}

// ----------------------------------------------------------------------------------------
// CSMART -- Get Account Balance
// ----------------------------------------------------------------------------------------
type CS_AccountDetails_Response struct {
	Response struct {
		Success string `json:"success"`
		Result  struct {
			Arguments struct {
				RESULTS struct {
					Elem           string `json:"elem"`
					BABUSINESSTYPE string `json:"BA_BUSINESS_TYPE"`
					BAACCOUNTNO    string `json:"BA_ACCOUNT_NO"`
				} `json:"RESULTS"`
				DATE         string `json:"DATE"`
				MSISDN       string `json:"MSISDN"`
				BUSINESSTYPE string `json:"BUSINESS_TYPE"`
				IMSI         string `json:"IMSI"`
				RATEPLANNAME string `json:"RATE_PLAN_NAME"`
				ACCOUNTNO    string `json:"ACCOUNT_NO"`
				LASTNAME     string `json:"LAST_NAME"`
				FIRSTNAME    string `json:"FIRST_NAME"`
				CREATEDT     string `json:"CREATED_T"`
				POID         string `json:"POID"`
			} `json:"arguments"`
			Message string `json:"message"`
		} `json:"result"`
	} `json:"response"`
}

func CS_GetAccountDetails(MSISDN string) (response CS_AccountDetails_Response, err error) {
	url := CDS_Endpoint + "/dx-sync/brm/v1/getAccountDetails/" + MSISDN
	log.Println("url: ", url)
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("srdate", time.Now().String())
	req.Header.Set("operation", "accountDetails")
	req.Header.Set("correlationId", uuid.New().String())
	req.Header.Set("serviceProvider", "AFRPGW")
	req.Header.Set("source", "AFRPGW")
	req.Header.Set("destination", "CRM")
	req.Header.Set("token", CDS_TOKEN)

	//client := &http.Client{}
	client := &http.Client{
		Timeout: 4 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error Connecting to Covalense API -- CS_dealerLogin: ", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	log.Println("here ", string(body))
	log.Println("status ", resp.Status)
	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Println("Error parsing Covalense API reply -- CS_dealerLogin: ", err, string(body))
		}
	} else {
		var error_reply CS_APIReplyError
		err = json.Unmarshal(body, &error_reply)
		if err != nil {
			log.Println("Error parsing Covalense API reply -- CS_dealerLogin: ", err, string(body))
		}
		err = errors.New(error_reply.Response.Result.Arguments.ErrorMessage)
	}
	return response, err
}
