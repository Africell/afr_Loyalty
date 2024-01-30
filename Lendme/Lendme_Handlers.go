package Lendme

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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
