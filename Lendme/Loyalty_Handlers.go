package Lendme

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

func (Uc *UserControl) HTTP_Loyalty_Governance(w http.ResponseWriter, r *http.Request) {
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
		sr.TransactionType = "Loyalty Governance - Read"
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
			schemes, err := Uc.Loyalty_Governance_GetPaginated(Page, Limit)
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
			schemes, err := Uc.Loyalty_Governance_Get(key)
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
		sr.TransactionType = "Loyalty Governance - Add"
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
		var request Loyalty_Governance_AddRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		Id, err := Uc.Loyalty_Governance_Add(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to add request"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Governance_Id = Id
		sr.Data = request

	case "PUT":
		sr.TransactionType = "Loyalty Governance - Edit"
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
		var request Loyalty_Governance_EditRequest
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
		id, err := Uc.Loyalty_Governance_Edit(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to edit"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Governance_Id = id
		sr.Data = request
	case "DELETE":
		sr.TransactionType = "Loyalty Governance - Delete"
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
		err := Uc.Loyalty_Governance_Delete(sr.Login, KeyDelete)
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

func (Uc *UserControl) HTTP_Loyalty_Level(w http.ResponseWriter, r *http.Request) {
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
		sr.TransactionType = "Loyalty Level - Read"
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
			schemes, err := Uc.Loyalty_Level_GetPaginated(Page, Limit)
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
			schemes, err := Uc.Loyalty_Level_Get(key)
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
		sr.TransactionType = "Loyalty Level - Add"
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
		var request Loyalty_Level_AddRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		Id, err := Uc.Loyalty_Level_Add(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to add request"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Level_Id = Id
		sr.Data = request

	case "PUT":
		sr.TransactionType = "Loyalty Level - Edit"
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
		var request Loyalty_Level_EditRequest
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
		id, err := Uc.Loyalty_Level_Edit(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to edit"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Level_Id = id
		sr.Data = request
	case "DELETE":
		sr.TransactionType = "Loyalty Level - Delete"
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
		err := Uc.Loyalty_Level_Delete(sr.Login, KeyDelete)
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

func (Uc *UserControl) HTTP_Loyalty_Account_Segment(w http.ResponseWriter, r *http.Request) {
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
	sr.TransactionType = "Loyalty Level"

	method := r.Method
	switch method {
	case "GET":
		sr.TransactionType = sr.TransactionType + " - Read"
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
			schemes, err := Uc.Loyalty_Account_Segment_GetPaginated(Page, Limit)
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
			schemes, err := Uc.Loyalty_Account_Segment_Get(key)
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
		sr.TransactionType = sr.TransactionType + " - Add"
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
		var request Loyalty_Account_Segment_AddRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		Id, err := Uc.Loyalty_Account_Segment_Add(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to add request"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Segment_Id = Id
		sr.Data = request

	case "PUT":
		sr.TransactionType = sr.TransactionType + " - Edit"
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
		var request Loyalty_Account_Segment_EditRequest
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
		id, err := Uc.Loyalty_Account_Segment_Edit(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to edit"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Segment_Id = id
		sr.Data = request
	case "DELETE":
		sr.TransactionType = sr.TransactionType + " - Delete"
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
		err := Uc.Loyalty_Account_Segment_Delete(sr.Login, KeyDelete)
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

func (Uc *UserControl) HTTP_Loyalty_Point_Earning_Rules(w http.ResponseWriter, r *http.Request) {
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
	sr.TransactionType = "Loyalty Point Earning Rules"

	method := r.Method
	switch method {
	case "GET":
		sr.TransactionType = sr.TransactionType + " - Read"
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
			schemes, err := Uc.Loyalty_Point_Earning_Rules_GetPaginated(Page, Limit)
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
			schemes, err := Uc.Loyalty_Point_Earning_Rules_Get(key)
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
		sr.TransactionType = sr.TransactionType + " - Add"
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
		var request Loyalty_Point_Earning_Rules_AddRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		Id, err := Uc.Loyalty_Point_Earning_Rules_Add(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to add request"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Earning_Rules_Id = Id
		sr.Data = request

	case "PUT":
		sr.TransactionType = sr.TransactionType + " - Edit"
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
		var request Loyalty_Point_Earning_Rules_EditRequest
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
		id, err := Uc.Loyalty_Point_Earning_Rules_Edit(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to edit"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Earning_Rules_Id = id
		sr.Data = request
	case "DELETE":
		sr.TransactionType = sr.TransactionType + " - Delete"
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
		err := Uc.Loyalty_Point_Earning_Rules_Delete(sr.Login, KeyDelete)
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

func (Uc *UserControl) HTTP_Loyalty_Point_Expiry_Rules(w http.ResponseWriter, r *http.Request) {
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
	sr.TransactionType = "Loyalty Point Expiry Rules"

	method := r.Method
	switch method {
	case "GET":
		sr.TransactionType = sr.TransactionType + " - Read"
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
			schemes, err := Uc.Loyalty_Point_Expiry_Rules_GetPaginated(Page, Limit)
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
			schemes, err := Uc.Loyalty_Point_Expiry_Rules_Get(key)
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
		sr.TransactionType = sr.TransactionType + " - Add"
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
		var request Loyalty_Point_Expiry_Rules_AddRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		Id, err := Uc.Loyalty_Point_Expiry_Rules_Add(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to add request"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Expiry_Rules_Id = Id
		sr.Data = request

	case "PUT":
		sr.TransactionType = sr.TransactionType + " - Edit"
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
		var request Loyalty_Point_Expiry_Rules_EditRequest
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
		id, err := Uc.Loyalty_Point_Expiry_Rules_Edit(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to edit"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Expiry_Rules_Id = id
		sr.Data = request
	case "DELETE":
		sr.TransactionType = sr.TransactionType + " - Delete"
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
		err := Uc.Loyalty_Point_Expiry_Rules_Delete(sr.Login, KeyDelete)
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

func (Uc *UserControl) HTTP_Loyalty_Point_Redemption_Rules(w http.ResponseWriter, r *http.Request) {
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
	sr.TransactionType = "Loyalty Point Redemption Rules"

	method := r.Method
	switch method {
	case "GET":
		sr.TransactionType = sr.TransactionType + " - Read"
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
			schemes, err := Uc.Loyalty_Point_Redemption_Rules_GetPaginated(Page, Limit)
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
			schemes, err := Uc.Loyalty_Point_Redemption_Rules_Get(key)
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
		sr.TransactionType = sr.TransactionType + " - Add"
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
		var request Loyalty_Point_Redemption_Rules_AddRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		Id, err := Uc.Loyalty_Point_Redemption_Rules_Add(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to add request"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Redemption_Rules_Id = Id
		sr.Data = request

	case "PUT":
		sr.TransactionType = sr.TransactionType + " - Edit"
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
		var request Loyalty_Point_Redemption_Rules_EditRequest
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
		id, err := Uc.Loyalty_Point_Redemption_Rules_Edit(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to edit"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Redemption_Rules_Id = id
		sr.Data = request
	case "DELETE":
		sr.TransactionType = sr.TransactionType + " - Delete"
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
		err := Uc.Loyalty_Point_Redemption_Rules_Delete(sr.Login, KeyDelete)
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

func (Uc *UserControl) HTTP_Loyalty_Plan(w http.ResponseWriter, r *http.Request) {
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
	sr.TransactionType = "Loyalty Plan"

	method := r.Method
	switch method {
	case "GET":
		sr.TransactionType = sr.TransactionType + " - Read"
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
			schemes, err := Uc.Loyalty_Plan_GetPaginated(Page, Limit)
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
			schemes, err := Uc.Loyalty_Plan_Get(key)
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
		sr.TransactionType = sr.TransactionType + " - Add"
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
		var request Loyalty_Plan_AddRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		Id, err := Uc.Loyalty_Plan_Add(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to add request"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Plan_Id = Id
		sr.Data = request

	case "PUT":
		sr.TransactionType = sr.TransactionType + " - Edit"
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
		var request Loyalty_Plan_EditRequest
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
		id, err := Uc.Loyalty_Plan_Edit(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to edit"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Plan_Id = id
		sr.Data = request
	case "DELETE":
		sr.TransactionType = sr.TransactionType + " - Delete"
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
		err := Uc.Loyalty_Plan_Delete(sr.Login, KeyDelete)
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

func (Uc *UserControl) HTTP_Customer_UAT(w http.ResponseWriter, r *http.Request) {
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
		sr.TransactionType = "Customer UAT - Read"
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
			entries, err := Uc.Customer_UAT_GetPaginated(Page, Limit)
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
			entries, err := Uc.Customer_UAT_Get(key)
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
		sr.TransactionType = "Customer UAT - Add"
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
		var request Customer_UAT_AddRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		Id, err := Uc.Customer_UAT_Add(sr.Login, request)
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
		sr.TransactionType = "Customer UAT - Edit"
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
		var request Customer_UAT_EditRequest
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
		id, err := Uc.Customer_UAT_Edit(sr.Login, request)
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
		sr.TransactionType = "Customer UAT - Delete"
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
		err := Uc.Customer_UAT_Delete(sr.Login, KeyDelete)
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

func (Uc *UserControl) HTTP_Customer_DND(w http.ResponseWriter, r *http.Request) {
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
		sr.TransactionType = "Customer DND - Read"
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
			entries, err := Uc.Customer_DND_GetPaginated(Page, Limit)
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
			entries, err := Uc.Customer_DND_Get(key)
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
		sr.TransactionType = "Customer DND - Add"
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
		var request Customer_DND_AddRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		Id, err := Uc.Customer_DND_Add(sr.Login, request)
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
		sr.TransactionType = "Customer DND - Edit"
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
		var request Customer_DND_EditRequest
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
		id, err := Uc.Customer_DND_Edit(sr.Login, request)
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
		sr.TransactionType = "Customer DND - Delete"
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
		err := Uc.Customer_DND_Delete(sr.Login, KeyDelete)
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

func (Uc *UserControl) HTTP_Customer_Exclusion(w http.ResponseWriter, r *http.Request) {
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
			entries, err := Uc.Customer_Exclusion_GetPaginated(Page, Limit)
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
			entries, err := Uc.Customer_Exclusion_Get(key)
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
		Id, err := Uc.Customer_Exclusion_Add(sr.Login, request)
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
		id, err := Uc.Customer_Exclusion_Edit(sr.Login, request)
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
		err := Uc.Customer_Exclusion_Delete(sr.Login, KeyDelete)
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

func (Uc *UserControl) HTTP_Customer_COS_Exclusion(w http.ResponseWriter, r *http.Request) {
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
			entries, err := Uc.Customer_COS_Exclusion_GetPaginated(Page, Limit)
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
			entries, err := Uc.Customer_COS_Exclusion_Get(key)
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
		Id, err := Uc.Customer_COS_Exclusion_Add(sr.Login, request)
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
		id, err := Uc.Customer_COS_Exclusion_Edit(sr.Login, request)
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
		err := Uc.Customer_COS_Exclusion_Delete(sr.Login, KeyDelete)
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

func (Uc *UserControl) HTTP_Customer_Loyalty_Account(w http.ResponseWriter, r *http.Request) {
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
	sr.TransactionType = "Loyalty Plan"

	method := r.Method
	switch method {
	case "GET":
		sr.TransactionType = sr.TransactionType + " - Read"
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
			schemes, err := Uc.Customer_Loyalty_Account_GetPaginated(Page, Limit)
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
			schemes, err := Uc.Customer_Loyalty_Account_Get(key)
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
		sr.TransactionType = sr.TransactionType + " - Add"
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
		var request Customer_Loyalty_Account_AddRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		Id, err := Uc.Customer_Loyalty_Account_Add(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to add request"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Customer_Id = Id
		sr.Data = request

	case "PUT":
		sr.TransactionType = sr.TransactionType + " - Edit"
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
		var request Customer_Loyalty_Account_EditRequest
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
		id, err := Uc.Customer_Loyalty_Account_Edit(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to edit"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Customer_Id = id
		sr.Data = request
	case "DELETE":
		sr.TransactionType = sr.TransactionType + " - Delete"
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
		err := Uc.Customer_Loyalty_Account_Delete(sr.Login, KeyDelete)
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

func (Uc *UserControl) HTTP_INLiveFeed_NewJoining(w http.ResponseWriter, r *http.Request) {
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
		sr.TransactionType = "INLiveFeed NewJoining"
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
		var request Customer_Loyalty_Account_AddRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		Id, err := Uc.Customer_Loyalty_Account_Add("INLiveFeed", request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to add request"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Customer_Id = Id
		sr.Data = request
		LiveFeedCounters.With(prometheus.Labels{"Stream": request.EventSource, "Type": "New Joining", "Description": "New Joining"}).Inc()
	}
	//successful response
	sr.Status = "successful"
	sr.StatusCode = http.StatusOK
	sr.StatusDescription = ""
	sr.ErrorDescription = ""
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
}

func (Uc *UserControl) HTTP_INLiveFeed_Churn(w http.ResponseWriter, r *http.Request) {
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

	case "DELETE":
		sr.TransactionType = "INLiveFeed Churn"
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
		var request Customer_Loyalty_Account_DeleteRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		err = Uc.Customer_Loyalty_Account_Delete(sr.Login, request.Key)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to delete"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		LiveFeedCounters.With(prometheus.Labels{"Stream": "IN_Feed", "Type": "Chrun", "Description": "Chrun"}).Inc()
	}
	//successful response
	sr.Status = "successful"
	sr.StatusCode = http.StatusOK
	sr.StatusDescription = ""
	sr.ErrorDescription = ""
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
}

func (Uc *UserControl) HTTP_INLiveFeed_Consuption(w http.ResponseWriter, r *http.Request) {
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
		sr.TransactionType = "INLiveFeed Consuption"
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
		var request Customer_Loyalty_Account_AwardRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		_, err = Uc.Customer_Loyalty_Account_AwardPoints(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to add request"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		LiveFeedCounters.With(prometheus.Labels{"Stream": request.EventSource, "Type": request.EventType, "Description": "Consumption"}).Inc()
	}
	//successful response
	sr.Status = "successful"
	sr.StatusCode = http.StatusOK
	sr.StatusDescription = ""
	sr.ErrorDescription = ""
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
}
