package Lendme

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

// ----------------------------------------------------------------------------------------
// CSMART -- EVC Get Dealer Balance
// ----------------------------------------------------------------------------------------
type CS_EVC_GetDealerBalance_Response struct {
	Success bool `json:"success"`
	Data    struct {
		Operation      string `json:"operation"`
		ResultCode     int    `json:"resultCode"`
		ResultText     string `json:"resultText"`
		Status         string `json:"status"`
		Balance        string `json:"balance"`
		Success        bool   `json:"success"`
		ResponseFields struct {
			ResultCode        string `json:"resultCode"`
			ResultText        string `json:"resultText"`
			Balance           string `json:"balance"`
			AvailableServices string `json:"availableServices"`
			VoucherGroupName  string `json:"voucherGroupName"`
		} `json:"responseFields"`
	} `json:"data"`
	Message string `json:"message"`
}

func CS_EVC_GetDealerBalance(DealerName, DealerPIN, DealerNumber string) (response CS_EVC_GetDealerBalance_Response, err error) {
	query := url.Values{}
	query.Set("dealerName", DealerName)
	query.Set("dealerPin", DealerPIN)
	query.Set("dealerNumber", DealerNumber)

	url := CDS_Endpoint + "/dx-sync/v1/dealer/balance?" + query.Encode()
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("srDate", time.Now().String())
	req.Header.Set("operation", "ocsDealerBalance")
	req.Header.Set("correlationId", uuid.New().String())
	req.Header.Set("serviceProvider", "OCS")
	req.Header.Set("source", "USSD")
	req.Header.Set("destination", "OCS")
	req.Header.Set("token", CDS_TOKEN)

	client := &http.Client{
		Timeout: 4 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error Connecting to Covalense API -- CS_EVC_GetDealerBalance: ", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading Covalense API reply -- CS_EVC_GetDealerBalance: ", err)
		return response, errors.New("error reading Covalense API reply: " + err.Error())
	}
	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Println("Error parsing Covalense API reply -- CS_EVC_GetDealerBalance: ", err, string(body))
		}
	} else {
		err = json.Unmarshal(body, &response)
		if err != nil {
			var error_reply CS_APIReplyError
			err = json.Unmarshal(body, &error_reply)
			if err != nil {
				log.Println("Error parsing Covalense API reply -- CS_EVC_GetDealerBalance: ", err, string(body))
				return response, errors.New("invalid Covalense API error JSON with HTTP status " + resp.Status + ": " + err.Error())
			}
			return response, errors.New(error_reply.Response.Result.Arguments.ErrorMessage)
		}
		if response.Message != "" {
			return response, errors.New(response.Message)
		}
		return response, errors.New("Covalense API returned HTTP status " + resp.Status)
	}
	return response, err
}

// ----------------------------------------------------------------------------------------
// CSMART -- EVC Dealer Transfer
// ----------------------------------------------------------------------------------------
type CS_EVC_Dealer_Transfer_Response struct {
	Success bool `json:"success"`
	Data    struct {
		Operation      string `json:"operation"`
		ResultCode     int    `json:"resultCode"`
		ResultText     string `json:"resultText"`
		Status         string `json:"status"`
		Balance        string `json:"balance"`
		Success        bool   `json:"success"`
		ResponseFields struct {
			ResultCode             string `json:"resultCode"`
			ResultText             string `json:"resultText"`
			Balance                string `json:"balance"`
			DealerTransactionID    string `json:"dealerTransactionID"`
			SubDealerTransactionID string `json:"subDealerTransactionID"`
		} `json:"responseFields"`
	} `json:"data"`
	Message string `json:"message"`
}

func CS_EVC_Dealer_Transfer(FromDealerID, DealerName, DealerPIN, ToDealerID, ToDealerName string, Amount float64, SendSMSNotifications bool) (response CS_EVC_Dealer_Transfer_Response, err error) {
	url := CDS_Endpoint + "/dx-sync/v1/dealer/transfer"
	payload := map[string]interface{}{
		"fromDealerId":         FromDealerID,
		"dealerName":           DealerName,
		"dealerPin":            DealerPIN,
		"toDealerId":           ToDealerID,
		"toDealerName":         ToDealerName,
		"amount":               Amount,
		"sendSmsNotifications": SendSMSNotifications,
	}
	_Json, _ := json.Marshal(payload)
	_jsonByte := []byte(_Json)
	method := "POST"
	req, err := http.NewRequest(method, url, bytes.NewBuffer(_jsonByte))
	if err != nil {
		return
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("srDate", time.Now().String())
	req.Header.Set("operation", "ocsDealerTransfer")
	req.Header.Set("correlationId", uuid.New().String())
	req.Header.Set("serviceProvider", "OCS")
	req.Header.Set("source", "USSD")
	req.Header.Set("destination", "OCS")
	req.Header.Set("token", CDS_TOKEN)

	client := &http.Client{
		Timeout: 4 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error Connecting to Covalense API -- CS_EVC_Dealer_Transfer: ", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading Covalense API reply -- CS_EVC_Dealer_Transfer: ", err)
		return response, errors.New("error reading Covalense API reply: " + err.Error())
	}
	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Println("Error parsing Covalense API reply -- CS_EVC_Dealer_Transfer: ", err, string(body))
		}
	} else {
		err = json.Unmarshal(body, &response)
		if err != nil {
			var error_reply CS_APIReplyError
			err = json.Unmarshal(body, &error_reply)
			if err != nil {
				log.Println("Error parsing Covalense API reply -- CS_EVC_Dealer_Transfer: ", err, string(body))
				return response, errors.New("invalid Covalense API error JSON with HTTP status " + resp.Status + ": " + err.Error())
			}
			return response, errors.New(error_reply.Response.Result.Arguments.ErrorMessage)
		}
		if response.Message != "" {
			return response, errors.New(response.Message)
		}
		return response, errors.New("Covalense API returned HTTP status " + resp.Status)
	}
	return response, err
}

// ----------------------------------------------------------------------------------------
// CSMART -- EVC Any Dealer Transfer
// ----------------------------------------------------------------------------------------
type CS_EVC_Any_Dealer_Transfer_Response = CS_EVC_Dealer_Transfer_Response

func CS_EVC_Any_Dealer_Transfer(DealerName, DealerNumber, DealerPIN, TargetDealerName, TargetDealerNumber string, Amount float64, SendSMSNotifications bool) (response CS_EVC_Any_Dealer_Transfer_Response, err error) {
	url := CDS_Endpoint + "/dx-sync/v1/dealer/any-dealer-credit-transfer"
	payload := map[string]interface{}{
		"dealerName":           DealerName,
		"dealerNumber":         DealerNumber,
		"dealerPin":            DealerPIN,
		"targetDealerName":     TargetDealerName,
		"targetDealerNumber":   TargetDealerNumber,
		"amount":               Amount,
		"sendSmsNotifications": SendSMSNotifications,
	}
	_Json, _ := json.Marshal(payload)
	_jsonByte := []byte(_Json)
	method := "POST"
	req, err := http.NewRequest(method, url, bytes.NewBuffer(_jsonByte))
	if err != nil {
		return
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("srDate", time.Now().String())
	req.Header.Set("operation", "ocsAnyDealerCredTransfer")
	req.Header.Set("correlationId", uuid.New().String())
	req.Header.Set("serviceProvider", "OCS")
	req.Header.Set("source", "USSD")
	req.Header.Set("destination", "OCS")
	req.Header.Set("token", CDS_TOKEN)

	client := &http.Client{
		Timeout: 4 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error Connecting to Covalense API -- CS_EVC_Any_Dealer_Transfer: ", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading Covalense API reply -- CS_EVC_Any_Dealer_Transfer: ", err)
		return response, errors.New("error reading Covalense API reply: " + err.Error())
	}
	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Println("Error parsing Covalense API reply -- CS_EVC_Any_Dealer_Transfer: ", err, string(body))
		}
	} else {
		err = json.Unmarshal(body, &response)
		if err != nil {
			var error_reply CS_APIReplyError
			err = json.Unmarshal(body, &error_reply)
			if err != nil {
				log.Println("Error parsing Covalense API reply -- CS_EVC_Any_Dealer_Transfer: ", err, string(body))
				return response, errors.New("invalid Covalense API error JSON with HTTP status " + resp.Status + ": " + err.Error())
			}
			return response, errors.New(error_reply.Response.Result.Arguments.ErrorMessage)
		}
		if response.Message != "" {
			return response, errors.New(response.Message)
		}
		return response, errors.New("Covalense API returned HTTP status " + resp.Status)
	}
	return response, err
}

// ----------------------------------------------------------------------------------------
// CSMART -- EVC Dealer Refund
// ----------------------------------------------------------------------------------------
type CS_EVC_Dealer_Refund_Response struct {
	Success bool `json:"success"`
	Data    struct {
		Operation      string `json:"operation"`
		ResultCode     int    `json:"resultCode"`
		ResultText     string `json:"resultText"`
		Status         string `json:"status"`
		Balance        string `json:"balance"`
		Success        bool   `json:"success"`
		ResponseFields struct {
			ResultCode          string `json:"resultCode"`
			ResultText          string `json:"resultText"`
			Balance             string `json:"balance"`
			DealerTransactionID string `json:"dealerTransactionID"`
			VoucherGroupName    string `json:"voucherGroupName"`
		} `json:"responseFields"`
	} `json:"data"`
	Message string `json:"message"`
}

func CS_EVC_Dealer_Refund(DealerID, DealerName, DealerPIN string, Amount float64, TargetNumber string) (response CS_EVC_Dealer_Refund_Response, err error) {
	url := CDS_Endpoint + "/dx-sync/v1/dealer/refund"
	payload := map[string]interface{}{
		"dealerId":     DealerID,
		"dealerName":   DealerName,
		"dealerPin":    DealerPIN,
		"amount":       Amount,
		"targetNumber": TargetNumber,
	}
	_Json, _ := json.Marshal(payload)
	_jsonByte := []byte(_Json)
	method := "POST"
	req, err := http.NewRequest(method, url, bytes.NewBuffer(_jsonByte))
	if err != nil {
		return
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("srDate", time.Now().String())
	req.Header.Set("operation", "ocsDealerRefund")
	req.Header.Set("correlationId", uuid.New().String())
	req.Header.Set("serviceProvider", "OCS")
	req.Header.Set("source", "USSD")
	req.Header.Set("destination", "OCS")
	req.Header.Set("token", CDS_TOKEN)

	client := &http.Client{
		Timeout: 4 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error Connecting to Covalense API -- CS_EVC_Dealer_Refund: ", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading Covalense API reply -- CS_EVC_Dealer_Refund: ", err)
		return response, errors.New("error reading Covalense API reply: " + err.Error())
	}
	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Println("Error parsing Covalense API reply -- CS_EVC_Dealer_Refund: ", err, string(body))
		}
	} else {
		err = json.Unmarshal(body, &response)
		if err != nil {
			var error_reply CS_APIReplyError
			err = json.Unmarshal(body, &error_reply)
			if err != nil {
				log.Println("Error parsing Covalense API reply -- CS_EVC_Dealer_Refund: ", err, string(body))
				return response, errors.New("invalid Covalense API error JSON with HTTP status " + resp.Status + ": " + err.Error())
			}
			return response, errors.New(error_reply.Response.Result.Arguments.ErrorMessage)
		}
		if response.Message != "" {
			return response, errors.New(response.Message)
		}
		return response, errors.New("Covalense API returned HTTP status " + resp.Status)
	}
	return response, err
}

// ----------------------------------------------------------------------------------------
// CSMART -- EVC Dealer Debit
// ----------------------------------------------------------------------------------------
type CS_EVC_Dealer_Debit_Response struct {
	Success bool `json:"success"`
	Data    struct {
		Operation      string `json:"operation"`
		ResultCode     int    `json:"resultCode"`
		ResultText     string `json:"resultText"`
		Status         string `json:"status"`
		Balance        string `json:"balance"`
		Success        bool   `json:"success"`
		ResponseFields struct {
			ResultCode          string `json:"resultCode"`
			ResultText          string `json:"resultText"`
			Balance             string `json:"balance"`
			DealerTransactionID string `json:"dealerTransactionID"`
			VoucherGroupName    string `json:"voucherGroupName"`
		} `json:"responseFields"`
	} `json:"data"`
	Message string `json:"message"`
}

func CS_EVC_Dealer_Debit(DealerID, DealerName, DealerPIN string, Amount float64, TargetNumber string) (response CS_EVC_Dealer_Debit_Response, err error) {
	url := CDS_Endpoint + "/dx-sync/v1/dealer/debit"
	payload := map[string]interface{}{
		"dealerId":     DealerID,
		"dealerName":   DealerName,
		"dealerPin":    DealerPIN,
		"amount":       Amount,
		"targetNumber": TargetNumber,
	}
	_Json, _ := json.Marshal(payload)
	_jsonByte := []byte(_Json)
	method := "POST"
	req, err := http.NewRequest(method, url, bytes.NewBuffer(_jsonByte))
	if err != nil {
		return
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("srDate", time.Now().String())
	req.Header.Set("operation", "ocsDealerDebit")
	req.Header.Set("correlationId", uuid.New().String())
	req.Header.Set("serviceProvider", "OCS")
	req.Header.Set("source", "USSD")
	req.Header.Set("destination", "OCS")
	req.Header.Set("token", CDS_TOKEN)

	client := &http.Client{
		Timeout: 4 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error Connecting to Covalense API -- CS_EVC_Dealer_Debit: ", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading Covalense API reply -- CS_EVC_Dealer_Debit: ", err)
		return response, errors.New("error reading Covalense API reply: " + err.Error())
	}
	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Println("Error parsing Covalense API reply -- CS_EVC_Dealer_Debit: ", err, string(body))
		}
	} else {
		err = json.Unmarshal(body, &response)
		if err != nil {
			var error_reply CS_APIReplyError
			err = json.Unmarshal(body, &error_reply)
			if err != nil {
				log.Println("Error parsing Covalense API reply -- CS_EVC_Dealer_Debit: ", err, string(body))
				return response, errors.New("invalid Covalense API error JSON with HTTP status " + resp.Status + ": " + err.Error())
			}
			return response, errors.New(error_reply.Response.Result.Arguments.ErrorMessage)
		}
		if response.Message != "" {
			return response, errors.New(response.Message)
		}
		return response, errors.New("Covalense API returned HTTP status " + resp.Status)
	}
	return response, err
}

// ----------------------------------------------------------------------------------------
// CSMART -- EVC Dealer Withdraw
// ----------------------------------------------------------------------------------------
type CS_EVC_Dealer_Withdraw_Response struct {
	Success bool `json:"success"`
	Data    struct {
		Operation      string  `json:"operation"`
		ResultCode     int     `json:"resultCode"`
		ResultText     string  `json:"resultText"`
		Status         string  `json:"status"`
		Balance        *string `json:"balance"`
		Success        bool    `json:"success"`
		ResponseFields struct {
			ResultCode          string `json:"resultCode"`
			ResultText          string `json:"resultText"`
			Balance             string `json:"balance"`
			DealerTransactionID string `json:"dealerTransactionID"`
			VoucherGroupName    string `json:"voucherGroupName"`
		} `json:"responseFields"`
	} `json:"data"`
	Message string `json:"message"`
}

func CS_EVC_Dealer_Withdraw(DealerID, DealerName, DealerPIN, SubDealerName, SubDealerNumber string, Amount float64, SendSMSNotifications bool) (response CS_EVC_Dealer_Withdraw_Response, err error) {
	url := CDS_Endpoint + "/dx-sync/v1/dealer/withdraw"
	payload := map[string]interface{}{
		"dealerId":             DealerID,
		"dealerName":           DealerName,
		"dealerPin":            DealerPIN,
		"subDealerName":        SubDealerName,
		"subDealerNumber":      SubDealerNumber,
		"amount":               Amount,
		"sendSMSNotifications": SendSMSNotifications,
	}
	_Json, _ := json.Marshal(payload)
	_jsonByte := []byte(_Json)
	method := "POST"
	req, err := http.NewRequest(method, url, bytes.NewBuffer(_jsonByte))
	if err != nil {
		return
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("srDate", time.Now().String())
	req.Header.Set("operation", "ocsDealerWithdraw")
	req.Header.Set("correlationId", uuid.New().String())
	req.Header.Set("serviceProvider", "OCS")
	req.Header.Set("source", "USSD")
	req.Header.Set("destination", "OCS")
	req.Header.Set("token", CDS_TOKEN)

	client := &http.Client{
		Timeout: 4 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error Connecting to Covalense API -- CS_EVC_Dealer_Withdraw: ", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading Covalense API reply -- CS_EVC_Dealer_Withdraw: ", err)
		return response, errors.New("error reading Covalense API reply: " + err.Error())
	}
	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Println("Error parsing Covalense API reply -- CS_EVC_Dealer_Withdraw: ", err, string(body))
		}
	} else {
		err = json.Unmarshal(body, &response)
		if err != nil {
			var error_reply CS_APIReplyError
			err = json.Unmarshal(body, &error_reply)
			if err != nil {
				log.Println("Error parsing Covalense API reply -- CS_EVC_Dealer_Withdraw: ", err, string(body))
				return response, errors.New("invalid Covalense API error JSON with HTTP status " + resp.Status + ": " + err.Error())
			}
			return response, errors.New(error_reply.Response.Result.Arguments.ErrorMessage)
		}
		if response.Message != "" {
			return response, errors.New(response.Message)
		}
		return response, errors.New("Covalense API returned HTTP status " + resp.Status)
	}
	return response, err
}

type CS_EVC_DealerWithdraw_Response = CS_EVC_Dealer_Withdraw_Response

func CS_EVC_DealerWithdraw(DealerID, DealerName, DealerPIN, SubDealerName, SubDealerNumber string, Amount float64, SendSMSNotifications bool) (response CS_EVC_DealerWithdraw_Response, err error) {
	return CS_EVC_Dealer_Withdraw(DealerID, DealerName, DealerPIN, SubDealerName, SubDealerNumber, Amount, SendSMSNotifications)
}

// ----------------------------------------------------------------------------------------
// CSMART -- EVC Dealer Change PIN
// ----------------------------------------------------------------------------------------
type CS_EVC_Dealer_Changepin_Response struct {
	Success bool `json:"success"`
	Data    struct {
		Operation      string  `json:"operation"`
		ResultCode     int     `json:"resultCode"`
		ResultText     string  `json:"resultText"`
		Status         string  `json:"status"`
		Balance        *string `json:"balance"`
		Success        bool    `json:"success"`
		ResponseFields struct {
			ResultCode string `json:"resultCode"`
			ResultText string `json:"resultText"`
		} `json:"responseFields"`
	} `json:"data"`
	Message string `json:"message"`
}

func CS_EVC_Dealer_Changepin(DealerName, DealerNumber, DealerPIN, NewPIN string) (response CS_EVC_Dealer_Changepin_Response, err error) {
	url := CDS_Endpoint + "/dx-sync/v1/dealer/change-pin"
	payload := map[string]interface{}{
		"dealerName":   DealerName,
		"dealerNumber": DealerNumber,
		"dealerPin":    DealerPIN,
		"newPin":       NewPIN,
	}
	_Json, _ := json.Marshal(payload)
	_jsonByte := []byte(_Json)
	method := "POST"
	req, err := http.NewRequest(method, url, bytes.NewBuffer(_jsonByte))
	if err != nil {
		return
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("srDate", time.Now().String())
	req.Header.Set("operation", "ocsDealerChangePin")
	req.Header.Set("correlationId", uuid.New().String())
	req.Header.Set("serviceProvider", "OCS")
	req.Header.Set("source", "USSD")
	req.Header.Set("destination", "OCS")
	req.Header.Set("token", CDS_TOKEN)

	client := &http.Client{
		Timeout: 4 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error Connecting to Covalense API -- CS_EVC_Dealer_Changepin: ", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading Covalense API reply -- CS_EVC_Dealer_Changepin: ", err)
		return response, errors.New("error reading Covalense API reply: " + err.Error())
	}
	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Println("Error parsing Covalense API reply -- CS_EVC_Dealer_Changepin: ", err, string(body))
		}
	} else {
		err = json.Unmarshal(body, &response)
		if err != nil {
			var error_reply CS_APIReplyError
			err = json.Unmarshal(body, &error_reply)
			if err != nil {
				log.Println("Error parsing Covalense API reply -- CS_EVC_Dealer_Changepin: ", err, string(body))
				return response, errors.New("invalid Covalense API error JSON with HTTP status " + resp.Status + ": " + err.Error())
			}
			return response, errors.New(error_reply.Response.Result.Arguments.ErrorMessage)
		}
		if response.Message != "" {
			return response, errors.New(response.Message)
		}
		return response, errors.New("Covalense API returned HTTP status " + resp.Status)
	}
	return response, err
}

// ----------------------------------------------------------------------------------------
// CSMART -- EVC Dealer Add Subdealer
// ----------------------------------------------------------------------------------------
type CS_EVC_Dealer_Addsubdealer_Response struct {
	Success bool `json:"success"`
	Data    struct {
		Operation      string  `json:"operation"`
		ResultCode     int     `json:"resultCode"`
		ResultText     string  `json:"resultText"`
		Status         string  `json:"status"`
		Balance        *string `json:"balance"`
		Success        bool    `json:"success"`
		ResponseFields struct {
			ResultCode string `json:"resultCode"`
			ResultText string `json:"resultText"`
		} `json:"responseFields"`
	} `json:"data"`
	Message string `json:"message"`
}

func CS_EVC_Dealer_Addsubdealer(DealerName, DealerNumber, DealerPIN, SubDealerName, SubDealerNumber, SubDealerPIN, WebServicesBitmap string) (response CS_EVC_Dealer_Addsubdealer_Response, err error) {
	url := CDS_Endpoint + "/dx-sync/v1/dealer/add-subdealer"
	payload := map[string]interface{}{
		"dealerName":        DealerName,
		"dealerNumber":      DealerNumber,
		"dealerPin":         DealerPIN,
		"subDealerName":     SubDealerName,
		"subDealerNumber":   SubDealerNumber,
		"subDealerPin":      SubDealerPIN,
		"webServicesBitmap": WebServicesBitmap,
	}
	_Json, _ := json.Marshal(payload)
	_jsonByte := []byte(_Json)
	method := "POST"
	req, err := http.NewRequest(method, url, bytes.NewBuffer(_jsonByte))
	if err != nil {
		return
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("srDate", time.Now().String())
	req.Header.Set("operation", "ocsAddSubDealer")
	req.Header.Set("correlationId", uuid.New().String())
	req.Header.Set("serviceProvider", "OCS")
	req.Header.Set("source", "USSD")
	req.Header.Set("destination", "OCS")
	req.Header.Set("token", CDS_TOKEN)

	client := &http.Client{
		Timeout: 4 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error Connecting to Covalense API -- CS_EVC_Dealer_Addsubdealer: ", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading Covalense API reply -- CS_EVC_Dealer_Addsubdealer: ", err)
		return response, errors.New("error reading Covalense API reply: " + err.Error())
	}
	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Println("Error parsing Covalense API reply -- CS_EVC_Dealer_Addsubdealer: ", err, string(body))
		}
	} else {
		err = json.Unmarshal(body, &response)
		if err != nil {
			var error_reply CS_APIReplyError
			err = json.Unmarshal(body, &error_reply)
			if err != nil {
				log.Println("Error parsing Covalense API reply -- CS_EVC_Dealer_Addsubdealer: ", err, string(body))
				return response, errors.New("invalid Covalense API error JSON with HTTP status " + resp.Status + ": " + err.Error())
			}
			return response, errors.New(error_reply.Response.Result.Arguments.ErrorMessage)
		}
		if response.Message != "" {
			return response, errors.New(response.Message)
		}
		return response, errors.New("Covalense API returned HTTP status " + resp.Status)
	}
	return response, err
}
