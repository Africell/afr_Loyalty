package Lendme

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (Uc *UserControl) HTTP_API_Standard_response(w http.ResponseWriter, r *http.Request, transaction API_Standard_response, KeepDataInDB bool) {
	//CollectPrometheusHTTPMetrics(r.URL.Path, transaction.Status, "", transaction.Elapsedtime)
	transaction.StatusDate = time.Now()
	transaction.Elapsedtime = (time.Since(transaction.ReceiveDate).Nanoseconds()) / 1000000
	//Uc.Write_StandardResponse_log(transaction, "", KeepDataInDB)
	w.Header().Set("Content-Type", "application/json")
	if transaction.Status == "successful" {
		w.WriteHeader(transaction.StatusCode)
	} else if transaction.Status == "failed" {
		w.WriteHeader(transaction.StatusCode)
	}
	json.NewEncoder(w).Encode(transaction)
}

func (Uc *UserControl) HTTP_Credit_Limit_Scheme(w http.ResponseWriter, r *http.Request) {
	var sr API_Standard_response
	//**fill response source detail
	SourceIp, _ := GetRequestIP(r)
	sr.SourceIP = SourceIp
	sr.Login = r.Header.Get("Login")
	sr.SourceApp = r.Header.Get("SourceApp")
	sr.AccessKey = r.URL.Path
	sr.AccessMethod = r.Method
	sr.HostId = Configuration.HostId
	sr.ReceiveDate = time.Now()

	method := r.Method
	switch method {
	case "GET":
		sr.TransactionType = "Credit Limit Scheme - Read"
		key := r.URL.Query().Get("Key")
		LimitStr := r.URL.Query().Get("Limit")
		PageStr := r.URL.Query().Get("Page")

		if LimitStr != "" || PageStr != "" {
			Limit, limiterr := strconv.ParseInt(LimitStr, 10, 64)
			if limiterr != nil {
				sr.Status = "failed"
				sr.StatusCode = http.StatusBadRequest
				sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
				sr.ErrorDescription = limiterr.Error()
				Uc.HTTP_API_Standard_response(w, r, sr, false)
				return
			}
			Page, pageerr := strconv.ParseInt(PageStr, 10, 64)
			if pageerr != nil {
				sr.Status = "failed"
				sr.StatusCode = http.StatusBadRequest
				sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
				sr.ErrorDescription = pageerr.Error()
				Uc.HTTP_API_Standard_response(w, r, sr, false)
				return
			}
			schemes, err := Uc.Credit_Limit_Scheme_GetPaginated(key, Page, Limit)
			if err != nil {
				sr.Status = "failed"
				sr.StatusCode = http.StatusBadRequest
				sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
				sr.ErrorDescription = err.Error()
				Uc.HTTP_API_Standard_response(w, r, sr, false)
				return
			}
			sr.Data = schemes
		} else {
			schemes, err := Uc.Credit_Limit_Scheme_Get(key)
			if err != nil {
				sr.Status = "failed"
				sr.StatusCode = http.StatusBadRequest
				sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
				sr.ErrorDescription = err.Error()
				Uc.HTTP_API_Standard_response(w, r, sr, false)
				return
			}
			sr.Data = schemes
		}
	case "POST":
		sr.TransactionType = "Credit Limit Scheme - Add"
		//parse body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to read request body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		var request Credit_Limit_Scheme_Add_Request
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		Id, err := Uc.Credit_Limit_Scheme_Add(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to add request"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Scheme_Id = Id
		sr.Data = request

	case "PUT":
		sr.TransactionType = "Credit Limit Scheme - Edit"
		//parse body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to read request body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		var request Credit_Limit_Scheme_Edit_Request
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		// if request.Key == "" {
		// 	sr.Status = "failed"
		// 	sr.StatusCode = http.StatusBadRequest
		// 	sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": key cannot be empty"
		// 	sr.ErrorDescription = "key cannot be empty"
		// 	Uc.HTTP_API_Standard_response(w, r, sr, false)
		// 	return
		// }
		id, err := Uc.Credit_Limit_Scheme_Edit(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to edit"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Scheme_Id = id
		sr.Data = request
	case "DELETE":
		sr.TransactionType = "Credit Limit Scheme - Delete"
		vars := mux.Vars(r)
		KeyDelete := vars["KeyDelete"]
		if KeyDelete == "" {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": key is empty"
			sr.ErrorDescription = "key is empty"
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		err := Uc.Credit_Limit_Scheme_Delete(KeyDelete)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to delete"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
	}
	//successful response
	sr.Status = "successful"
	sr.StatusCode = http.StatusOK
	sr.StatusDescription = ""
	sr.ErrorDescription = ""
	Uc.HTTP_API_Standard_response(w, r, sr, true)
}

func (Uc *UserControl) HTTP_Subscriber(w http.ResponseWriter, r *http.Request) {
	var sr API_Standard_response
	//**fill response source detail
	SourceIp, _ := GetRequestIP(r)
	sr.SourceIP = SourceIp
	sr.Login = r.Header.Get("Login")
	sr.SourceApp = r.Header.Get("SourceApp")
	sr.AccessKey = r.URL.Path
	sr.AccessMethod = r.Method
	sr.HostId = Configuration.HostId
	sr.ReceiveDate = time.Now()

	method := r.Method
	switch method {
	case "GET":
		sr.TransactionType = "Subscriber - Read"
		key := r.URL.Query().Get("Key")
		LimitStr := r.URL.Query().Get("Limit")
		PageStr := r.URL.Query().Get("Page")

		if LimitStr != "" || PageStr != "" {
			Limit, limiterr := strconv.ParseInt(LimitStr, 10, 64)
			if limiterr != nil {
				sr.Status = "failed"
				sr.StatusCode = http.StatusBadRequest
				sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
				sr.ErrorDescription = limiterr.Error()
				Uc.HTTP_API_Standard_response(w, r, sr, false)
				return
			}
			Page, pageerr := strconv.ParseInt(PageStr, 10, 64)
			if pageerr != nil {
				sr.Status = "failed"
				sr.StatusCode = http.StatusBadRequest
				sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
				sr.ErrorDescription = pageerr.Error()
				Uc.HTTP_API_Standard_response(w, r, sr, false)
				return
			}
			subscriber, err := Uc.Subscriber_GetPaginated(key, Page, Limit)
			if err != nil {
				sr.Status = "failed"
				sr.StatusCode = http.StatusBadRequest
				sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
				sr.ErrorDescription = err.Error()
				Uc.HTTP_API_Standard_response(w, r, sr, false)
				return
			}
			sr.Data = subscriber
		} else {
			subscriber, err := Uc.Subscriber_Get(key)
			if err != nil {
				sr.Status = "failed"
				sr.StatusCode = http.StatusBadRequest
				sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
				sr.ErrorDescription = err.Error()
				Uc.HTTP_API_Standard_response(w, r, sr, false)
				return
			}
			sr.Data = subscriber
		}
	case "POST":
		sr.TransactionType = "Subscriber - Add"
		//parse body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to read request body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		var request Subscriber_Add_Request
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		Id, err := Uc.Subscriber_Add(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to add request"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Subscriber_Id = Id
		sr.Data = request

	case "PUT":
		sr.TransactionType = "Subscriber - Edit"
		//parse body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to read request body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		var request Subscriber_Edit_Request
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		// if request.Key == "" {
		// 	sr.Status = "failed"
		// 	sr.StatusCode = http.StatusBadRequest
		// 	sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": key cannot be empty"
		// 	sr.ErrorDescription = "key cannot be empty"
		// 	Uc.HTTP_API_Standard_response(w, r, sr, false)
		// 	return
		// }
		id, err := Uc.Subscriber_Edit(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to edit"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Subscriber_Id = id
		sr.Data = request
	case "DELETE":
		sr.TransactionType = "Subscriber - Delete"
		vars := mux.Vars(r)
		KeyDelete := vars["KeyDelete"]
		if KeyDelete == "" {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": key is empty"
			sr.ErrorDescription = "key is empty"
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		err := Uc.Subscriber_Delete(KeyDelete)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to delete"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
	}
	//successful response
	sr.Status = "successful"
	sr.StatusCode = http.StatusOK
	sr.StatusDescription = ""
	sr.ErrorDescription = ""
	Uc.HTTP_API_Standard_response(w, r, sr, true)
}

func (Uc *UserControl) HTTP_Subscriber_USSD(w http.ResponseWriter, r *http.Request) {
	var sr API_Standard_response
	//**fill response source detail
	SourceIp, _ := GetRequestIP(r)
	sr.SourceIP = SourceIp
	sr.Login = r.Header.Get("Login")
	sr.SourceApp = r.Header.Get("SourceApp")
	sr.AccessKey = r.URL.Path
	sr.AccessMethod = r.Method
	sr.HostId = Configuration.HostId
	sr.ReceiveDate = time.Now()

	method := r.Method
	switch method {
	case "GET":
		sr.TransactionType = "Subscriber USSD - Read"
		MSISDN := r.URL.Query().Get("MSISDN")

		subscriber, err := Uc.SubscriberUSSD_Get(MSISDN)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		sr.Data = subscriber

	}
	//successful response
	sr.Status = "successful"
	sr.StatusCode = http.StatusOK
	sr.StatusDescription = ""
	sr.ErrorDescription = ""
	Uc.HTTP_API_Standard_response(w, r, sr, true)
}

func (Uc *UserControl) HTTP_Subscribers_ARPU_Import_Launch(w http.ResponseWriter, r *http.Request) {
	var sr API_Standard_response
	//**fill response source detail
	SourceIp, _ := GetRequestIP(r)
	sr.SourceIP = SourceIp
	sr.Login = r.Header.Get("Login")
	sr.SourceApp = r.Header.Get("SourceApp")
	sr.AccessKey = r.URL.Path
	sr.AccessMethod = r.Method
	sr.HostId = Configuration.HostId
	sr.ReceiveDate = time.Now()

	method := r.Method
	switch method {
	case "GET":
		sr.TransactionType = "Subscriber - Read"
		FileName := r.URL.Query().Get("FileName")
		go Uc.Import_Subscribers_Dump_LineByLine(FileName)
	}
	//successful response
	sr.Status = "successful"
	sr.StatusCode = http.StatusOK
	sr.StatusDescription = "importing process launched"
	sr.ErrorDescription = ""
	Uc.HTTP_API_Standard_response(w, r, sr, true)
}

func (Uc *UserControl) HTTP_Subscribers_GetINDetail(w http.ResponseWriter, r *http.Request) {
	var sr API_Standard_response
	//**fill response source detail
	SourceIp, _ := GetRequestIP(r)
	sr.SourceIP = SourceIp
	sr.Login = r.Header.Get("Login")
	sr.SourceApp = r.Header.Get("SourceApp")
	sr.AccessKey = r.URL.Path
	sr.AccessMethod = r.Method
	sr.HostId = Configuration.HostId
	sr.ReceiveDate = time.Now()

	method := r.Method
	switch method {
	case "GET":
		sr.TransactionType = "Subscriber Get IN Detail- Read"
		MSISDN := r.URL.Query().Get("MSISDN")

		IN_MSISDN := MSISDN
		if Configuration.Operation == "Gambia" {
			if len(MSISDN) > 7 {
				IN_MSISDN = IN_MSISDN[len(MSISDN)-7 : len(MSISDN)]
			}
		} else if Configuration.Operation == "SierraLeone" { //077928014
			if len(MSISDN) > 8 {
				IN_MSISDN = "0" + IN_MSISDN[len(MSISDN)-8:len(MSISDN)]
			}
		}
		IN_Response, err := Uc.IN.INClient.GetAccountDetails("", "", IN_MSISDN)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest)
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		sr.Data = IN_Response
	}
	//successful response
	sr.Status = "successful"
	sr.StatusCode = http.StatusOK
	sr.StatusDescription = ""
	sr.ErrorDescription = ""
	Uc.HTTP_API_Standard_response(w, r, sr, true)
}

func (Uc *UserControl) HTTP_Lendme_Request(w http.ResponseWriter, r *http.Request) {
	var sr API_Standard_response
	//**fill response source detail
	SourceIp, _ := GetRequestIP(r)
	sr.SourceIP = SourceIp
	sr.Login = r.Header.Get("Login")
	sr.SourceApp = r.Header.Get("SourceApp")
	sr.AccessKey = r.URL.Path
	sr.AccessMethod = r.Method
	sr.HostId = Configuration.HostId
	sr.ReceiveDate = time.Now()

	method := r.Method
	switch method {
	case "POST":
		sr.TransactionType = "Lendme Request - Add"
		//parse body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to read request body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		var request LendMe_Request
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		err = Uc.Lendme_exec_Request(request.Source, request.MSISDN, request.Amount)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to add request"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		//request.Subscriber_Id = Id
		sr.Data = request
	}
	//successful response
	sr.Status = "successful"
	sr.StatusCode = http.StatusOK
	sr.StatusDescription = ""
	sr.ErrorDescription = ""
	Uc.HTTP_API_Standard_response(w, r, sr, true)
}

func (Uc *UserControl) HTTP_AO_Lendme_Request(w http.ResponseWriter, r *http.Request) {
	var sr API_Standard_response
	//**fill response source detail
	SourceIp, _ := GetRequestIP(r)
	sr.SourceIP = SourceIp
	sr.Login = r.Header.Get("Login")
	sr.SourceApp = r.Header.Get("SourceApp")
	sr.AccessKey = r.URL.Path
	sr.AccessMethod = r.Method
	sr.HostId = Configuration.HostId
	sr.ReceiveDate = time.Now()

	method := r.Method
	switch method {
	case "POST":
		sr.TransactionType = "Lendme Request - Add"
		//parse body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to read request body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		var request LendMe_Request
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		err = Uc.LendmeAO_exec_Request(request.Source, request.MSISDN, request.Amount)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to add request"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		//request.Subscriber_Id = Id
		sr.Data = request
	}
	//successful response
	sr.Status = "successful"
	sr.StatusCode = http.StatusOK
	sr.StatusDescription = ""
	sr.ErrorDescription = ""
	Uc.HTTP_API_Standard_response(w, r, sr, true)
}

func (Uc *UserControl) HTTP_Lendme_PayBack(w http.ResponseWriter, r *http.Request) {
	var sr API_Standard_response
	//**fill response source detail
	SourceIp, _ := GetRequestIP(r)
	sr.SourceIP = SourceIp
	sr.Login = r.Header.Get("Login")
	sr.SourceApp = r.Header.Get("SourceApp")
	sr.AccessKey = r.URL.Path
	sr.AccessMethod = r.Method
	sr.HostId = Configuration.HostId
	sr.ReceiveDate = time.Now()

	method := r.Method
	switch method {

	case "POST":
		sr.TransactionType = "Lendme PayBack - Add"
		//parse body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to read request body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		var request Lendme_PayBack_Request
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		err = Uc.Lendme_PayBack(request.Source, request.MSISDN, request.RechargeAmount, request.Opid)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest)
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		sr.Data = request
	}
	//successful response
	sr.Status = "successful"
	sr.StatusCode = http.StatusOK
	sr.StatusDescription = ""
	sr.ErrorDescription = ""
	Uc.HTTP_API_Standard_response(w, r, sr, true)
}

func (Uc *UserControl) HTTP_AO_Lendme_PayBack(w http.ResponseWriter, r *http.Request) {
	var sr API_Standard_response
	//**fill response source detail
	SourceIp, _ := GetRequestIP(r)
	sr.SourceIP = SourceIp
	sr.Login = r.Header.Get("Login")
	sr.SourceApp = r.Header.Get("SourceApp")
	sr.AccessKey = r.URL.Path
	sr.AccessMethod = r.Method
	sr.HostId = Configuration.HostId
	sr.ReceiveDate = time.Now()

	method := r.Method
	switch method {

	case "POST":
		sr.TransactionType = "Lendme PayBack - Add"
		//parse body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to read request body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		var request Lendme_PayBack_Request
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		if request.MSISDN == "244959560801" {
			fmt.Println("Payback Received ", request.Source, request.MSISDN, request.RechargeAmount, request.Opid)
		}
		err = Uc.LendmeAO_PayBack(request.Source, request.MSISDN, request.RechargeAmount, request.Opid)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest)
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		sr.Data = request
	}
	//successful response
	sr.Status = "successful"
	sr.StatusCode = http.StatusOK
	sr.StatusDescription = ""
	sr.ErrorDescription = ""
	Uc.HTTP_API_Standard_response(w, r, sr, true)
}

func (Uc *UserControl) HTTP_Lendme_SendSMS(w http.ResponseWriter, r *http.Request) {
	var sr API_Standard_response
	//**fill response source detail
	SourceIp, _ := GetRequestIP(r)
	sr.SourceIP = SourceIp
	sr.Login = r.Header.Get("Login")
	sr.SourceApp = r.Header.Get("SourceApp")
	sr.AccessKey = r.URL.Path
	sr.AccessMethod = r.Method
	sr.HostId = Configuration.HostId
	sr.ReceiveDate = time.Now()

	method := r.Method
	switch method {
	case "GET":
		sr.TransactionType = "Lendme Send SMS"
		Sender := r.URL.Query().Get("Sender")
		Target := r.URL.Query().Get("Target")
		Text := r.URL.Query().Get("Text")

		err := SendSMS(Sender, Target, Text)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = "failed to send sms"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
	}
	//successful response
	sr.Status = "successful"
	sr.StatusCode = http.StatusOK
	sr.StatusDescription = ""
	sr.ErrorDescription = ""
	Uc.HTTP_API_Standard_response(w, r, sr, true)
}

func (Uc *UserControl) HTTP_Lendme_Logs(w http.ResponseWriter, r *http.Request) {
	var sr API_Standard_response
	//**fill response source detail
	SourceIp, _ := GetRequestIP(r)
	sr.SourceIP = SourceIp
	sr.Login = r.Header.Get("Login")
	sr.SourceApp = r.Header.Get("SourceApp")
	sr.AccessKey = r.URL.Path
	sr.AccessMethod = r.Method
	sr.HostId = Configuration.HostId
	sr.ReceiveDate = time.Now()
	method := r.Method
	switch method {
	case "GET":
		sr.TransactionType = "Lendme logs"
		//Filter := r.URL.Query().Get("filter")
		//ServiceType := r.URL.Query().Get("ServiceType")
		// LimitStr := r.URL.Query().Get("Limit")
		// PageStr := r.URL.Query().Get("Page")

		MSISDN := r.URL.Query().Get("MSISDN")
		startDate := r.URL.Query().Get("startDate")
		endDate := r.URL.Query().Get("endDate")
		startDateDD, _ := time.ParseInLocation("1/2/2006", startDate, time.Local)
		endDateDD, _ := time.ParseInLocation("1/2/2006", endDate, time.Local)
		// endDateDD = time.Date(endDateDD.Year(), endDateDD.Month(), endDateDD.Day(), 23, 59, 59, 0, endDateDD.Location())
		//
		// if LimitStr != "" || PageStr != "" {
		// 	Limit, limiterr := strconv.ParseInt(LimitStr, 10, 64)
		// 	if limiterr != nil {
		// 		sr.Status = "failed"
		// 		sr.StatusCode = http.StatusBadRequest
		// 		sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
		// 		sr.ErrorDescription = limiterr.Error()
		// 		Uc.HTTP_API_Standard_response(w, r, sr, false)
		// 		return
		// 	}
		// 	Page, pageerr := strconv.ParseInt(PageStr, 10, 64)
		// 	if pageerr != nil {
		// 		sr.Status = "failed"
		// 		sr.StatusCode = http.StatusBadRequest
		// 		sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
		// 		sr.ErrorDescription = pageerr.Error()
		// 		Uc.HTTP_API_Standard_response(w, r, sr, false)
		// 		return
		// 	}
		// 	fCDMRequests, err := Uc.Customer_Loyalty_Account_DebitPoints_log_GetPaginated(startDateDD, endDateDD, MSISDN, Page, Limit)
		// 	if err != nil {
		// 		sr.Status = "failed"
		// 		sr.StatusCode = http.StatusBadRequest
		// 		sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
		// 		sr.ErrorDescription = err.Error()
		// 		Uc.HTTP_API_Standard_response(w, r, sr, false)
		// 		return
		// 	}
		// 	sr.Data = fCDMRequests
		// } else {
		data, err := Uc.Customer_Lendme_Subscriber_Getlogs(startDateDD, endDateDD, MSISDN, "")
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		sr.Data = data
		// }
	}
	//successful response
	sr.Status = "successful"
	sr.StatusCode = http.StatusOK
	sr.StatusDescription = ""
	sr.ErrorDescription = ""
	Uc.HTTP_API_Standard_response(w, r, sr, true)
}

func (Uc *UserControl) Customer_Lendme_Subscriber_Getlogs(startDate, endDate time.Time, MSISDN string, Filter string) (response []Lendme_log, err error) {
	findResult, err := Uc.ReadSubscriberLogsFromMongoDB(startDate, endDate, MSISDN, 1, 0)
	if err != nil {
		return response, err
	}
	response = append(response, findResult...)

	return response, nil
}
func (Uc *UserControl) ReadSubscriberLogsFromMongoDB(startDate, endDate time.Time, MSISDN string, page, limit int64) (response []Lendme_log, err error) {
	// Ensure startDate is before or equal to endDate

	if startDate.After(endDate) {
		return response, fmt.Errorf("start date cannot be after end date")
	}

	if page < 1 {
		page = 1 // Default to first page if invalid page number
	}

	// Iterate over the range of dates
	for currentDate := startDate; !currentDate.After(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		monthStr := strconv.Itoa(int(currentDate.Month()))
		if len(monthStr) < 2 {
			monthStr = "0" + monthStr
		}
		MongoDB_DB_Name := "Lendme_log_Archive_" + strconv.Itoa(currentDate.Year()) + monthStr

		var dayStr = strconv.Itoa(currentDate.Day())
		if len(dayStr) < 2 {
			dayStr = "0" + dayStr
		}

		collName := "Col_Lendme_log_" + dayStr

		// Fetch the collection
		collection := Uc.LoyaltyMongoDB.MongoDBClient.Database(MongoDB_DB_Name).Collection(collName)
		// Build the filter for the date range
		filter := bson.D{
			{Key: "MSISDN", Value: MSISDN},
		}
		skip := (page - 1) * limit

		// Fetch a paginated list of documents using limit and skip
		cursor, err := collection.Find(
			context.Background(),
			filter,
			options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)),
		)
		if err != nil {
			return response, err
		}
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var result Lendme_log

			if err := cursor.Decode(&result); err != nil {
				return nil, fmt.Errorf("failed to decode result: %w", err)
			}
			response = append(response, result)
		}

		if err := cursor.Err(); err != nil {
			return nil, fmt.Errorf("cursor error: %w", err)
		}
	}

	return response, err
}

func (Uc *UserControl) HTTP_Lendme_Customer_Exclusion(w http.ResponseWriter, r *http.Request) {
	var sr API_Standard_response
	//**fill response source detail
	SourceIp, _ := GetRequestIP(r)
	sr.SourceIP = SourceIp
	sr.Login = r.Header.Get("Login")
	sr.SourceApp = r.Header.Get("SourceApp")
	sr.AccessKey = r.URL.Path
	sr.AccessMethod = r.Method
	sr.HostId = Configuration.HostId
	sr.ReceiveDate = time.Now()

	method := r.Method
	switch method {
	case "GET":
		sr.TransactionType = "Customer Exclusion - Read"
		key := r.URL.Query().Get("Key")
		LimitStr := r.URL.Query().Get("Limit")
		PageStr := r.URL.Query().Get("Page")

		if LimitStr != "" || PageStr != "" {
			Limit, limiterr := strconv.ParseInt(LimitStr, 10, 64)
			if limiterr != nil {
				sr.Status = "failed"
				sr.StatusCode = http.StatusBadRequest
				sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
				sr.ErrorDescription = limiterr.Error()
				Uc.HTTP_API_Standard_response(w, r, sr, false)
				return
			}
			Page, pageerr := strconv.ParseInt(PageStr, 10, 64)
			if pageerr != nil {
				sr.Status = "failed"
				sr.StatusCode = http.StatusBadRequest
				sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
				sr.ErrorDescription = pageerr.Error()
				Uc.HTTP_API_Standard_response(w, r, sr, false)
				return
			}
			entries, err := Uc.Lendme_Customer_Exclusion_GetPaginated(Page, Limit)
			if err != nil {
				sr.Status = "failed"
				sr.StatusCode = http.StatusBadRequest
				sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
				sr.ErrorDescription = err.Error()
				Uc.HTTP_API_Standard_response(w, r, sr, false)
				return
			}
			sr.Data = entries
		} else {
			entries, err := Uc.Lendme_Customer_Exclusion_Get(key)
			if err != nil {
				sr.Status = "failed"
				sr.StatusCode = http.StatusBadRequest
				sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
				sr.ErrorDescription = err.Error()
				Uc.HTTP_API_Standard_response(w, r, sr, false)
				return
			}
			sr.Data = entries
		}
	case "POST":
		sr.TransactionType = "Customer Exclusion - Add"
		//parse body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to read request body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		var request Customer_Exclusion_AddRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		Id, err := Uc.Lendme_Customer_Exclusion_Add(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to add request"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Id = Id
		sr.Data = request

	case "PUT":
		sr.TransactionType = "Customer Exclusion - Edit"
		//parse body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to read request body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		var request Customer_Exclusion_EditRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		// if request.Key == "" {
		// 	sr.Status = "failed"
		// 	sr.StatusCode = http.StatusBadRequest
		// 	sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": key cannot be empty"
		// 	sr.ErrorDescription = "key cannot be empty"
		// 	Uc.HTTP_API_Standard_response(w, r, sr, false)
		// 	return
		// }
		id, err := Uc.Lendme_Customer_Exclusion_Edit(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to edit"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Id = id
		sr.Data = request
	case "DELETE":
		sr.TransactionType = "Customer Exclusion - Delete"
		vars := mux.Vars(r)
		KeyDelete := vars["KeyDelete"]
		if KeyDelete == "" {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": key is empty"
			sr.ErrorDescription = "key is empty"
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		err := Uc.Lendme_Customer_Exclusion_Delete(sr.Login, KeyDelete)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to delete"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
	}
	//successful response
	sr.Status = "successful"
	sr.StatusCode = http.StatusOK
	sr.StatusDescription = ""
	sr.ErrorDescription = ""
	Uc.HTTP_API_Standard_response(w, r, sr, true)
}

func (Uc *UserControl) HTTP_Lendme_Customer_COS_Exclusion(w http.ResponseWriter, r *http.Request) {
	var sr API_Standard_response
	//**fill response source detail
	SourceIp, _ := GetRequestIP(r)
	sr.SourceIP = SourceIp
	sr.Login = r.Header.Get("Login")
	sr.SourceApp = r.Header.Get("SourceApp")
	sr.AccessKey = r.URL.Path
	sr.AccessMethod = r.Method
	sr.HostId = Configuration.HostId
	sr.ReceiveDate = time.Now()

	method := r.Method
	switch method {
	case "GET":
		sr.TransactionType = "Customer COS Exclusion - Read"
		key := r.URL.Query().Get("Key")
		LimitStr := r.URL.Query().Get("Limit")
		PageStr := r.URL.Query().Get("Page")

		if LimitStr != "" || PageStr != "" {
			Limit, limiterr := strconv.ParseInt(LimitStr, 10, 64)
			if limiterr != nil {
				sr.Status = "failed"
				sr.StatusCode = http.StatusBadRequest
				sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
				sr.ErrorDescription = limiterr.Error()
				Uc.HTTP_API_Standard_response(w, r, sr, false)
				return
			}
			Page, pageerr := strconv.ParseInt(PageStr, 10, 64)
			if pageerr != nil {
				sr.Status = "failed"
				sr.StatusCode = http.StatusBadRequest
				sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
				sr.ErrorDescription = pageerr.Error()
				Uc.HTTP_API_Standard_response(w, r, sr, false)
				return
			}
			entries, err := Uc.Lendme_Customer_COS_Exclusion_GetPaginated(Page, Limit)
			if err != nil {
				sr.Status = "failed"
				sr.StatusCode = http.StatusBadRequest
				sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
				sr.ErrorDescription = err.Error()
				Uc.HTTP_API_Standard_response(w, r, sr, false)
				return
			}
			sr.Data = entries
		} else {
			entries, err := Uc.Lendme_Customer_COS_Exclusion_Get(key)
			if err != nil {
				sr.Status = "failed"
				sr.StatusCode = http.StatusBadRequest
				sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
				sr.ErrorDescription = err.Error()
				Uc.HTTP_API_Standard_response(w, r, sr, false)
				return
			}
			sr.Data = entries
		}
	case "POST":
		sr.TransactionType = "Customer COS Exclusion - Add"
		//parse body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to read request body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		var request Customer_COS_Exclusion_AddRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		Id, err := Uc.Lendme_Customer_COS_Exclusion_Add(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to add request"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Id = Id
		sr.Data = request

	case "PUT":
		sr.TransactionType = "Customer COS Exclusion - Edit"
		//parse body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to read request body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		var request Customer_COS_Exclusion_EditRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		// if request.Key == "" {
		// 	sr.Status = "failed"
		// 	sr.StatusCode = http.StatusBadRequest
		// 	sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": key cannot be empty"
		// 	sr.ErrorDescription = "key cannot be empty"
		// 	Uc.HTTP_API_Standard_response(w, r, sr, false)
		// 	return
		// }
		id, err := Uc.Lendme_Customer_COS_Exclusion_Edit(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to edit"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Id = id
		sr.Data = request
	case "DELETE":
		sr.TransactionType = "Customer COS Exclusion - Delete"
		vars := mux.Vars(r)
		KeyDelete := vars["KeyDelete"]
		if KeyDelete == "" {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": key is empty"
			sr.ErrorDescription = "key is empty"
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		err := Uc.Lendme_Customer_COS_Exclusion_Delete(sr.Login, KeyDelete)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to delete"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
	}
	//successful response
	sr.Status = "successful"
	sr.StatusCode = http.StatusOK
	sr.StatusDescription = ""
	sr.ErrorDescription = ""
	Uc.HTTP_API_Standard_response(w, r, sr, true)
}
