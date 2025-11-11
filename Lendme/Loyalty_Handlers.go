package Lendme

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"regexp"
	"sync"

	"net/http"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

var uploadedFiles = make(map[string]struct{}) // keep track of file names
var uploadedFilesMu sync.Mutex

func (Uc *UserControl) Validate_Headers(r *http.Request) (response Request_Header) {
	SourceIp, _ := GetRequestIP(r)
	response.SourceIP = SourceIp
	response.SourceApp = r.Header.Get("SourceApp")
	response.AppLogin = r.Header.Get("Login")
	response.HostId = Configuration.HostId
	response.AppVersion = r.Header.Get("AppVersion")
	response.GSMLocation = r.Header.Get("GSMLocation")
	response.Authorization = r.Header.Get("Authorization")
	// Latitude_str := r.Header.Get("Latitude")
	// Longitude_str := r.Header.Get("Longitude")
	//check token type
	// TokenType := r.Header.Get("TokenType")
	// if TokenType != "Access" {
	// 	response.ValidationDescription = "token is not provided"
	// 	return
	// }
	// if Latitude_str != "" && Longitude_str != "" {
	// 	//valid latitude values are between -90 and 90, both inclusive
	// 	Latitude, err := strconv.ParseFloat(Latitude_str, 64)
	// 	if err != nil {
	// 		response.ValidationDescription = "Invalid Latitude: " + err.Error()
	// 		return
	// 	}
	// 	if Latitude == 0 {
	// 		response.ValidationDescription = "latitude is not provided"
	// 		return
	// 	}
	// 	if Latitude < -90 || Latitude > 90 {
	// 		response.ValidationDescription = "invalid latitude value"
	// 		return
	// 	}
	// 	//valid longitude values are between -180 and 180, both inclusive
	// 	Longitude, err := strconv.ParseFloat(Longitude_str, 64)
	// 	if err != nil {
	// 		response.ValidationDescription = "Invalid Longitude: " + err.Error()
	// 		return
	// 	}
	// 	if Longitude == 0 {
	// 		response.ValidationDescription = "longitude is not provided"
	// 		return
	// 	}
	// 	if Longitude < -180 || Latitude > 180 {
	// 		response.ValidationDescription = "invalid longitude value"
	// 		return
	// 	}
	// 	response.GPSLocation = daoc.NewGeospatialPoint(Longitude, Latitude)
	// }
	response.IsValid = true
	return
}

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

func (Uc *UserControl) HTTP_Loyalty_Seniority_Level(w http.ResponseWriter, r *http.Request) {
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
	sr.TransactionType = "Loyalty Seniority Level"

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
			schemes, err := Uc.Loyalty_Seniority_Level_GetPaginated(Page, Limit)
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
			schemes, err := Uc.Loyalty_Seniority_Level_Get(key)
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
		var request Loyalty_Seniority_Level
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		Id, err := Uc.Loyalty_Seniority_Level_Add(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to add request"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Seniority_Id = Id
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
		var request Loyalty_Seniority_Level
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
		id, err := Uc.Loyalty_Seniority_Level_Edit(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to edit"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Seniority_Id = id
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
		err := Uc.Loyalty_Seniority_Level_Delete(sr.Login, KeyDelete)
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
func (Uc *UserControl) HTTP_Loyalty_Point_Earning_Rules_Overwrite(w http.ResponseWriter, r *http.Request) {
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
	sr.TransactionType = "Loyalty Point Earning Rules Overwrite"

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
			schemes, err := Uc.Loyalty_Point_Earning_Rules_Overwrite_GetPaginated(Page, Limit)
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
			schemes, err := Uc.Loyalty_Point_Earning_Rules_Overwrite_Get(key)
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
		var request Loyalty_Point_Earning_Rules_Overwrite_AddRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		Id, err := Uc.Loyalty_Point_Earning_Rules_Overwrite_Add(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to add request"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Earning_Rules_Overwrite_Id = Id
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
		var request Loyalty_Point_Earning_Rules_Overwrite_EditRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		id, err := Uc.Loyalty_Point_Earning_Rules_Overwrite_Edit(sr.Login, request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to edit"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		request.Earning_Rules_Overwrite_Id = id
		sr.Data = request
	case "DELETE":
		sr.TransactionType = sr.TransactionType + " - Delete"
		KeyDelete := r.URL.Query().Get("Key")
		err := Uc.Loyalty_Point_Earning_Rules_Overwrite_Delete(sr.Login, KeyDelete)
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
		key := r.URL.Query().Get("KeyDelete")
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

func (Uc *UserControl) HTTP_Customer_Loyalty_Account_Points_Details(w http.ResponseWriter, r *http.Request) {
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
	sr.TransactionType = "Customer Loyalty Account Points Details"

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
			schemes, err := Uc.Customer_Loyalty_Account_Points_Details_GetPaginated(Page, Limit)
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
			schemes, err := Uc.Customer_Loyalty_Account_Points_Details_Get(key)
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
	}
	//successful response
	sr.Status = "successful"
	sr.StatusCode = http.StatusOK
	sr.StatusDescription = ""
	sr.ErrorDescription = ""
	Uc.HTTP_API_Standard_response(w, r, sr, true)
}

func (Uc *UserControl) HTTP_Customer_Loyalty_Account_GetRedemption_Rules(w http.ResponseWriter, r *http.Request) {
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
	sr.TransactionType = "Customer Loyalty Get Redemption Rules"

	method := r.Method
	switch method {
	case "GET":
		sr.TransactionType = sr.TransactionType + ""
		MSISDN := r.URL.Query().Get("MSISDN")
		response, err := Uc.Customer_Loyalty_Account_GetRedemption_Rules(MSISDN)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		sr.Data = response
	}
	//successful response
	sr.Status = "successful"
	sr.StatusCode = http.StatusOK
	sr.StatusDescription = ""
	sr.ErrorDescription = ""
	Uc.HTTP_API_Standard_response(w, r, sr, true)
}

func (Uc *UserControl) HTTP_Customer_Loyalty_Account_DebitPoints(w http.ResponseWriter, r *http.Request) {
	var transaction Loyalty_AccountDebitPoints_log
	transaction.ReceiveDate = time.Now()
	validated_Headers := Uc.Validate_Headers(r)
	transaction.SourceIP = validated_Headers.SourceIP
	transaction.SourceApp = validated_Headers.SourceApp
	transaction.AppLogin = validated_Headers.AppLogin
	transaction.AppVersion = validated_Headers.AppVersion
	transaction.GPSLocation = validated_Headers.GPSLocation
	transaction.GSMLocation = validated_Headers.GSMLocation
	if !validated_Headers.IsValid {
		transaction.Status = "failed"
		transaction.StatusDescription = validated_Headers.ValidationDescription
		Uc.HTTP_Customer_Loyalty_Account_DebitPoints_Response(w, r, &transaction, true)
		return
	}
	//parse body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		transaction.Status = "failed"
		transaction.StatusDescription = "failed to read request body"
		transaction.ErrorDescription = err.Error()
		Uc.HTTP_Customer_Loyalty_Account_DebitPoints_Response(w, r, &transaction, true)
		return
	}
	var Request Loyalty_AccountDebitPoints_Request
	err = json.Unmarshal(body, &Request)
	if err != nil {
		transaction.Status = "failed"
		transaction.StatusDescription = "failed to parse request body"
		transaction.ErrorDescription = err.Error()
		Uc.HTTP_Customer_Loyalty_Account_DebitPoints_Response(w, r, &transaction, true)
		return
	}
	//execute the request
	Uc.Loyalty_AccountDebitPoints(&validated_Headers, Request, &transaction)
	Uc.HTTP_Customer_Loyalty_Account_DebitPoints_Response(w, r, &transaction, false)
}

func (Uc *UserControl) HTTP_Customer_Loyalty_Account_DebitPoints_Response(w http.ResponseWriter, r *http.Request, transaction *Loyalty_AccountDebitPoints_log, DB_Write bool) {
	switch transaction.Status {
	case "successful":
		transaction.StatusCode = http.StatusOK
	case "failed":
		transaction.StatusCode = http.StatusBadRequest
	}
	transaction.StatusDate = time.Now()
	transaction.E2E_Elapsedtime = (time.Since(transaction.ReceiveDate).Nanoseconds()) / 1000000
	if DB_Write {
		Uc.Write_Loyalty_AccountDebitPoints_log(*transaction)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(transaction.StatusCode)
	json.NewEncoder(w).Encode(transaction)
}

func (Uc *UserControl) HTTP_Customer_Loyalty_RedeemRequest(w http.ResponseWriter, r *http.Request) {
	var transaction Loyalty_Redemption_log
	transaction.ReceiveDate = time.Now()
	validated_Headers := Uc.Validate_Headers(r)
	transaction.SourceIP = validated_Headers.SourceIP
	transaction.SourceApp = validated_Headers.SourceApp
	transaction.AppLogin = validated_Headers.AppLogin
	transaction.AppVersion = validated_Headers.AppVersion
	transaction.GPSLocation = validated_Headers.GPSLocation
	transaction.GSMLocation = validated_Headers.GSMLocation
	if !validated_Headers.IsValid {
		transaction.Status = "failed"
		transaction.StatusDescription = validated_Headers.ValidationDescription
		Uc.HTTP_Customer_Loyalty_RedeemRequest_Response(w, r, &transaction, true)
		return
	}
	//parse body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		transaction.Status = "failed"
		transaction.StatusDescription = "failed to read request body"
		transaction.ErrorDescription = err.Error()
		Uc.HTTP_Customer_Loyalty_RedeemRequest_Response(w, r, &transaction, true)
		return
	}
	var Request Loyalty_Redemption_Request
	err = json.Unmarshal(body, &Request)
	if err != nil {
		transaction.Status = "failed"
		transaction.StatusDescription = "failed to parse request body"
		transaction.ErrorDescription = err.Error()
		Uc.HTTP_Customer_Loyalty_RedeemRequest_Response(w, r, &transaction, true)
		return
	}
	//execute the request
	Uc.Customer_Loyalty_RedeemRequest(&validated_Headers, Request, &transaction)
	Uc.HTTP_Customer_Loyalty_RedeemRequest_Response(w, r, &transaction, false)
}

func (Uc *UserControl) HTTP_Customer_Loyalty_Account_OptRequest(w http.ResponseWriter, r *http.Request) {
	var transaction Loyalty_Status_log
	transaction.StatusDate = time.Now()
	validated_Headers := Uc.Validate_Headers(r)
	transaction.SourceIP = validated_Headers.SourceIP
	transaction.SourceApp = validated_Headers.SourceApp
	transaction.AppLogin = validated_Headers.AppLogin
	transaction.AppVersion = validated_Headers.AppVersion

	if !validated_Headers.IsValid {
		transaction.Request_Status = "failed"
		transaction.StatusDescription = validated_Headers.ValidationDescription
		Uc.HTTP_Customer_Loyalty_OptRequest_Response(w, r, &transaction, true)
		return
	}
	//parse body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		transaction.Request_Status = "failed"
		transaction.StatusDescription = "failed to read request body"
		transaction.ErrorDescription = err.Error()
		Uc.HTTP_Customer_Loyalty_OptRequest_Response(w, r, &transaction, true)
		return
	}
	var Request Loyalty_Opt_Request
	err = json.Unmarshal(body, &Request)
	if err != nil {
		transaction.Request_Status = "failed"
		transaction.StatusDescription = "failed to parse request body"
		transaction.ErrorDescription = err.Error()
		Uc.HTTP_Customer_Loyalty_OptRequest_Response(w, r, &transaction, true)
		return
	}
	//execute the request
	Uc.Customer_Loyalty_OptRequest(&validated_Headers, Request, &transaction)
	Uc.HTTP_Customer_Loyalty_OptRequest_Response(w, r, &transaction, false)
}

func (Uc *UserControl) HTTP_Customer_Loyalty_RedeemRequest_Response(w http.ResponseWriter, r *http.Request, transaction *Loyalty_Redemption_log, DB_Write bool) {
	switch transaction.Status {
	case "successful":
		transaction.StatusCode = http.StatusOK
	case "failed":
		transaction.StatusCode = http.StatusBadRequest
	}
	transaction.StatusDate = time.Now()
	transaction.E2E_Elapsedtime = (time.Since(transaction.ReceiveDate).Nanoseconds()) / 1000000
	if DB_Write {
		Uc.Write_Loyalty_Redemption_log(*transaction)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(transaction.StatusCode)
	json.NewEncoder(w).Encode(transaction)
}

func (Uc *UserControl) HTTP_Customer_Loyalty_OptRequest_Response(w http.ResponseWriter, r *http.Request, transaction *Loyalty_Status_log, DB_Write bool) {
	switch transaction.Request_Status {
	case "successful":
		transaction.Request_StatusCode = http.StatusOK
	case "failed":
		transaction.Request_StatusCode = http.StatusBadRequest
	}
	transaction.E2E_Elapsedtime = (time.Since(transaction.StatusDate).Nanoseconds()) / 1000000
	if DB_Write {
		Uc.Write_Loyalty_Status_log(*transaction)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(transaction.Request_StatusCode)
	json.NewEncoder(w).Encode(transaction)
}
func (Uc *UserControl) HTTP_Customer_Loyalty_Account_CreditPoints(w http.ResponseWriter, r *http.Request) {
	var transaction Loyalty_AccountCreditPoints_log
	transaction.ReceiveDate = time.Now()
	validated_Headers := Uc.Validate_Headers(r)
	transaction.SourceIP = validated_Headers.SourceIP
	transaction.SourceApp = validated_Headers.SourceApp
	transaction.AppLogin = validated_Headers.AppLogin
	transaction.AppVersion = validated_Headers.AppVersion
	transaction.GPSLocation = validated_Headers.GPSLocation
	transaction.GSMLocation = validated_Headers.GSMLocation
	if !validated_Headers.IsValid {
		transaction.Status = "failed"
		transaction.StatusDescription = validated_Headers.ValidationDescription
		Uc.HTTP_Customer_Loyalty_Account_CreditPoints_Response(w, r, &transaction, true)
		return
	}
	//parse body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		transaction.Status = "failed"
		transaction.StatusDescription = "failed to read request body"
		transaction.ErrorDescription = err.Error()
		Uc.HTTP_Customer_Loyalty_Account_CreditPoints_Response(w, r, &transaction, true)
		return
	}
	var Request Loyalty_AccountCreditPoints_Request
	err = json.Unmarshal(body, &Request)
	if err != nil {
		transaction.Status = "failed"
		transaction.StatusDescription = "failed to parse request body"
		transaction.ErrorDescription = err.Error()
		Uc.HTTP_Customer_Loyalty_Account_CreditPoints_Response(w, r, &transaction, true)
		return
	}
	//execute the request
	Uc.Loyalty_AccountCreditPoints(&validated_Headers, Request, &transaction)
	Uc.HTTP_Customer_Loyalty_Account_CreditPoints_Response(w, r, &transaction, false)
}

func (Uc *UserControl) HTTP_Customer_Loyalty_Account_Get_Awarded_Points(w http.ResponseWriter, r *http.Request) {
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
	sr.TransactionType = "Customer Loyalty Get Redemption Rules"

	method := r.Method
	switch method {
	case "GET":
		sr.TransactionType = sr.TransactionType + ""
		// MSISDN := r.URL.Query().Get("MSISDN")
		startDate := r.URL.Query().Get("startDate")
		endDate := r.URL.Query().Get("endDate")
		startDateDD, _ := time.ParseInLocation("1/2/2006", startDate, time.Local)
		endDateDD, _ := time.ParseInLocation("1/2/2006", endDate, time.Local)
		response, err := Uc.Customer_Loyalty_Account_GetAwardedPoints(startDateDD, endDateDD)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		w.Header().Set("Content-Disposition", "attachment; filename=awarded_points.csv")
		w.Header().Set("Content-Type", "text/csv")

		writer := csv.NewWriter(w)
		defer writer.Flush()

		writer.Write([]string{"MSISDN", "Awarded Points", "Loyalty Level", "Afrimony Points", "Current Points", "Acumalted Points"})

		for msisdn, data := range response {
			writer.Write([]string{msisdn, fmt.Sprintf("%.2f", data.TotalPoints), data.LoyaltyLevel, fmt.Sprintf("%.2f", data.AfrimoneyPoints), fmt.Sprintf("%.2f", data.Current_Points), fmt.Sprintf("%.2f", data.Accumulated_Points)})
		}
		fmt.Println("err", err)
		// json.NewEncoder(w).Encode(err)

		// sr.Data = response
	}
	//successful response
	// sr.Status = "successful"
	// sr.StatusCode = http.StatusOK
	// sr.StatusDescription = ""
	// sr.ErrorDescription = ""
	// Uc.HTTP_API_Standard_response(w, r, sr, true)
}

func (Uc *UserControl) HTTP_Customer_Loyalty_Account_CreditPoints_Response(w http.ResponseWriter, r *http.Request, transaction *Loyalty_AccountCreditPoints_log, DB_Write bool) {
	switch transaction.Status {
	case "successful":
		transaction.StatusCode = http.StatusOK
	case "failed":
		transaction.StatusCode = http.StatusBadRequest
	}
	transaction.StatusDate = time.Now()
	transaction.E2E_Elapsedtime = (time.Since(transaction.ReceiveDate).Nanoseconds()) / 1000000
	if DB_Write {
		Uc.Write_Loyalty_AccountCreditPoints_log(*transaction)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(transaction.StatusCode)
	json.NewEncoder(w).Encode(transaction)
}

func (Uc *UserControl) HTTP_Loyalty_Products_Catalogue(w http.ResponseWriter, r *http.Request) {
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
	sr.TransactionType = "Loyalty Products Catalogue"

	method := r.Method
	switch method {
	case "GET":
		sr.TransactionType = sr.TransactionType + " - Read"
		MSISDN := r.URL.Query().Get("MSISDN")
		log.Println("MSISDN:", MSISDN)
		productCatalogue, err := Uc.Customer_Loyalty_Account_GetRedemptionProductCatalogue(MSISDN)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		sr.Data = productCatalogue
	}
	//successful response
	sr.Status = "successful"
	sr.StatusCode = http.StatusOK
	sr.StatusDescription = ""
	sr.ErrorDescription = ""
	Uc.HTTP_API_Standard_response(w, r, sr, true)
}
func (Uc *UserControl) HTTP_Bulk_Loyalty_Points_Crediting(w http.ResponseWriter, r *http.Request) {
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
	sr.TransactionType = "Bulk Loyalty Points Deduction"

	method := r.Method
	switch method {
	case "POST":
		sr.TransactionType = sr.TransactionType + ""
		var transaction Loyalty_Redemption_log
		transaction.ReceiveDate = time.Now()
		validated_Headers := Uc.Validate_Headers(r)
		transaction.SourceIP = validated_Headers.SourceIP
		transaction.SourceApp = validated_Headers.SourceApp
		transaction.AppLogin = validated_Headers.AppLogin
		transaction.AppVersion = validated_Headers.AppVersion
		transaction.GPSLocation = validated_Headers.GPSLocation
		transaction.GSMLocation = validated_Headers.GSMLocation
		if !validated_Headers.IsValid {
			transaction.Status = "failed"
			transaction.StatusDescription = validated_Headers.ValidationDescription
			Uc.HTTP_Customer_Loyalty_RedeemRequest_Response(w, r, &transaction, true)
			return
		}
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		file, fileHeader, err := r.FormFile("fileUpload")
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}

		uploadedFilesMu.Lock()
		if uploadedFiles == nil {
			uploadedFiles = make(map[string]struct{})
		}
		fileName := fileHeader.Filename
		if _, exists := uploadedFiles[fileName]; exists {
			uploadedFilesMu.Unlock()
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = "Duplicate file upload"
			sr.ErrorDescription = "File " + fileName + " has already been uploaded"
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		uploadedFiles[fileName] = struct{}{} // mark as uploaded
		uploadedFilesMu.Unlock()

		jobID := uuid.New().String()

		// Save job status
		jobsMu.Lock()
		jobs[jobID] = &JobStatus{Status: "running", Result: make(map[string][]string)}
		jobsMu.Unlock()
		// Run in background
		go func(jobID string, file multipart.File) {
			defer file.Close()

			var entries []*CustomersUploadList
			if err := gocsv.UnmarshalMultipartFile(&file, &entries); err != nil {
				jobsMu.Lock()
				jobs[jobID].Status = "failed"
				jobs[jobID].Errors = append(jobs[jobID].Errors, err.Error())
				jobsMu.Unlock()
				return
			}

			jobsMu.Lock()
			jobs[jobID].TotalRows = len(entries)
			jobsMu.Unlock()

			for i := range entries {
				jobsMu.Lock()
				jobs[jobID].ProcessedRows = i + 1
				jobsMu.Unlock()
				msisdn := entries[i].MSISDN

				if entries[i].Points <= 0 {
					jobsMu.Lock()
					jobs[jobID].Result["Failed"] = append(jobs[jobID].Result["Failed"], msisdn+" points can not be negative or zero")
					jobsMu.Unlock()
					continue
				}

				processedMu.Lock()
				if processed == nil {
					processed = make(map[string]struct{})
				}
				if _, exists := processed[msisdn]; exists {
					processedMu.Unlock()
					jobsMu.Lock()
					jobs[jobID].Result["Failed"] = append(
						jobs[jobID].Result["Failed"],
						msisdn+" already credited",
					)
					jobsMu.Unlock()
					continue
				}
				processedMu.Unlock()
				var loyalty_AccountCreditPoints_log Loyalty_AccountCreditPoints_log
				var credit_Request Loyalty_AccountCreditPoints_Request
				credit_Request.MSISDN = entries[i].MSISDN
				credit_Request.EventAmount = 0
				credit_Request.EventDescription = "Bulk Points Crediting"
				credit_Request.PointsToCredit = float64(int(entries[i].Points))
				credit_Request.EventSource = "Bulk Points Crediting"
				Uc.Loyalty_AccountCreditPoints(&validated_Headers, credit_Request, &loyalty_AccountCreditPoints_log)
				if loyalty_AccountCreditPoints_log.Status == "failed" {
					jobsMu.Lock()
					jobs[jobID].Result["Failed"] = append(jobs[jobID].Result["Failed"], loyalty_AccountCreditPoints_log.MSISDN+" "+loyalty_AccountCreditPoints_log.StatusDescription)
					jobsMu.Unlock()
				} else {
					processedMu.Lock()
					processed[msisdn] = struct{}{}
					processedMu.Unlock()
					jobsMu.Lock()
					jobs[jobID].Result["Successful"] = append(jobs[jobID].Result["Successful"], loyalty_AccountCreditPoints_log.MSISDN+" "+loyalty_AccountCreditPoints_log.StatusDescription)
					slice := jobs[jobID].Result["Points Credited"]
					jobsMu.Unlock()
					var successfulVal float64
					if len(slice) > 0 && slice[0] != "" {
						successfulVal, _ = strconv.ParseFloat(slice[0], 64)
					}
					pointsCredited := successfulVal + loyalty_AccountCreditPoints_log.AwardedPoints
					jobsMu.Lock()
					jobs[jobID].Result["Points Credited"] = []string{fmt.Sprintf("%.2f", pointsCredited)}
					jobsMu.Unlock()
				}

			}
			Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
				Event_User:         validated_Headers.AppLogin,
				Event_Time:         time.Now(),
				Event_AffectedType: "Bulk Points Crediting",
				Event_ActionType:   "Add",
				Event_Description:  "",
				Event_Entry_Before: nil,
				Event_Entry_After:  jobs[jobID].Result,
			})
			jobsMu.Lock()
			jobs[jobID].Status = "completed"
			jobsMu.Unlock()
		}(jobID, file)

		sr.Data = map[string]interface{}{
			"jobID": jobID,
		}
	}
	//successful response
	processedMu.Lock()
	processed = nil
	processedMu.Unlock()
	sr.Status = "successful"
	sr.StatusCode = http.StatusOK
	sr.StatusDescription = ""
	sr.ErrorDescription = ""
	Uc.HTTP_API_Standard_response(w, r, sr, true)
}
func (Uc *UserControl) HTTP_Bulk_Loyalty_Points_Crediting_Progress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["jobID"]
	if jobID == "" {
		http.Error(w, "missing jobID", http.StatusBadRequest)
		return
	}
	jobsMu.Lock()
	job, ok := jobs[jobID]
	jobsMu.Unlock()

	if !ok {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func (Uc *UserControl) HTTP_Bulk_Loyalty_Points_Deduction(w http.ResponseWriter, r *http.Request) {
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
	sr.TransactionType = "Bulk Loyalty Points Deduction"

	method := r.Method
	switch method {
	case "POST":
		sr.TransactionType = sr.TransactionType + ""
		var transaction Loyalty_Redemption_log
		transaction.ReceiveDate = time.Now()
		validated_Headers := Uc.Validate_Headers(r)
		transaction.SourceIP = validated_Headers.SourceIP
		transaction.SourceApp = validated_Headers.SourceApp
		transaction.AppLogin = validated_Headers.AppLogin
		transaction.AppVersion = validated_Headers.AppVersion
		transaction.GPSLocation = validated_Headers.GPSLocation
		transaction.GSMLocation = validated_Headers.GSMLocation
		if !validated_Headers.IsValid {
			transaction.Status = "failed"
			transaction.StatusDescription = validated_Headers.ValidationDescription
			Uc.HTTP_Customer_Loyalty_RedeemRequest_Response(w, r, &transaction, true)
			return
		}
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		file, fileHeader, err := r.FormFile("fileUpload")
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		re := regexp.MustCompile(`^Loyalty_deduction_\d{6}_\d{2}\.csv$`)
		if !re.MatchString(fileHeader.Filename) {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = "Wrong File name format"
			sr.ErrorDescription = "File " + fileHeader.Filename + " does not match the desired format"
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}

		uploadedFilesMu.Lock()
		defer uploadedFilesMu.Unlock()
		if uploadedFiles == nil {
			uploadedFiles = make(map[string]struct{})
		}
		fileName := fileHeader.Filename
		if _, exists := uploadedFiles[fileName]; exists {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = "Duplicate file upload"
			sr.ErrorDescription = "File " + fileName + " has already been uploaded"
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		uploadedFiles[fileName] = struct{}{} // mark as uploaded

		jobID := uuid.New().String()

		// Save job status
		jobsMu.Lock()
		jobs[jobID] = &JobStatus{Status: "running", Result: make(map[string][]string)}
		jobsMu.Unlock()
		// Run in background
		go func(jobID string, file multipart.File) {
			defer file.Close()

			var entries []*CustomersUploadList
			if err := gocsv.UnmarshalMultipartFile(&file, &entries); err != nil {
				jobsMu.Lock()
				jobs[jobID].Status = "failed"
				jobs[jobID].Errors = append(jobs[jobID].Errors, err.Error())
				jobsMu.Unlock()
				return
			}

			jobsMu.Lock()
			jobs[jobID].TotalRows = len(entries)
			jobsMu.Unlock()

			for i, _ := range entries {
				jobsMu.Lock()
				jobs[jobID].ProcessedRows = i + 1
				jobsMu.Unlock()
				// process each row
				msisdn := entries[i].MSISDN

				if entries[i].Points <= 0 {
					jobsMu.Lock()
					jobs[jobID].Result["Failed"] = append(jobs[jobID].Result["Failed"], msisdn+" points can not be negative or zero")
					jobsMu.Unlock()
					continue
				}
				processedMu.Lock()
				if processed == nil {
					processed = make(map[string]struct{})
				}
				if _, exists := processed[msisdn]; exists {
					processedMu.Unlock()
					jobsMu.Lock()
					jobs[jobID].Result["Failed"] = append(
						jobs[jobID].Result["Failed"],
						msisdn+" already debited",
					)
					jobsMu.Unlock()
					continue
				}
				processedMu.Unlock()

				var loyalty_AccountDebitPoints_log Loyalty_AccountDebitPoints_log
				var debit_Request Loyalty_AccountDebitPoints_Request
				debit_Request.MSISDN = entries[i].MSISDN
				debit_Request.Debit_Amount = float64(int(entries[i].Points))
				debit_Request.Debit_Reason = "Bulk Points Deduction"
				Uc.Loyalty_AccountDebitPoints(&validated_Headers, debit_Request, &loyalty_AccountDebitPoints_log, true)

				if loyalty_AccountDebitPoints_log.Status == "failed" {
					if loyalty_AccountDebitPoints_log.StatusDescription == "no enough points" {
						entry_na, exits := Map_Customer_Loyalty_Account.CheckThenGet(loyalty_AccountDebitPoints_log.MSISDN)
						if !exits {
							jobsMu.Lock()
							jobs[jobID].Result["Failed"] = append(jobs[jobID].Result["Failed"], loyalty_AccountDebitPoints_log.MSISDN+" key does not exist")
							jobsMu.Unlock()
							continue
						}
						entry, ok := entry_na.(Customer_Loyalty_Account)
						if !ok {
							jobsMu.Lock()
							jobs[jobID].Result["Failed"] = append(jobs[jobID].Result["Failed"], loyalty_AccountDebitPoints_log.MSISDN+" error in type assertion")
							jobsMu.Unlock()
							continue
						}
						if entry.Available_Points > 0 {
							if entry.Outstanding_fraction_points > 0 {
								entry.Outstanding_fraction_points = 0
								Map_Customer_Loyalty_Account.Put(entry.Key, entry)
							}
							debit_Request.Debit_Amount = float64(int(entry.Available_Points))
							Uc.Loyalty_AccountDebitPoints(&validated_Headers, debit_Request, &loyalty_AccountDebitPoints_log, true)
							if loyalty_AccountDebitPoints_log.Status == "failed" {
								jobsMu.Lock()
								jobs[jobID].Result["Failed"] = append(jobs[jobID].Result["Failed"], loyalty_AccountDebitPoints_log.MSISDN+" "+loyalty_AccountDebitPoints_log.StatusDescription)
								jobsMu.Unlock()
							} else {
								processedMu.Lock()
								if processed == nil {
									processed = make(map[string]struct{})
								}
								processed[msisdn] = struct{}{}
								processedMu.Unlock()
								jobsMu.Lock()
								//amal here
								jobs[jobID].Result["Partially Successful"] = append(
									jobs[jobID].Result["Partially Successful"],
									loyalty_AccountDebitPoints_log.MSISDN+
										" had insufficient points. Deducted "+strconv.FormatFloat(debit_Request.Debit_Amount, 'f', -1, 64)+
										" and balance is now 0 (requested "+strconv.FormatFloat(entries[i].Points, 'f', -1, 64)+").",
								)
								slice := jobs[jobID].Result["Points Deducted"]
								jobsMu.Unlock()

								var successfulVal float64
								if len(slice) > 0 && slice[0] != "" {
									successfulVal, _ = strconv.ParseFloat(slice[0], 64)
								}
								pointsDeducted := successfulVal + loyalty_AccountDebitPoints_log.Debit_Amount

								jobsMu.Lock()
								jobs[jobID].Result["Points Deducted"] = []string{fmt.Sprintf("%.2f", pointsDeducted)}
								jobsMu.Unlock()
							}
							continue
						}
					}
					jobsMu.Lock()
					jobs[jobID].Result["Failed"] = append(jobs[jobID].Result["Failed"], loyalty_AccountDebitPoints_log.MSISDN+" "+loyalty_AccountDebitPoints_log.StatusDescription)
					jobsMu.Unlock()

				} else {
					processedMu.Lock()
					processed[msisdn] = struct{}{}
					processedMu.Unlock()
					jobsMu.Lock()
					jobs[jobID].Result["Successful"] = append(jobs[jobID].Result["Successful"], loyalty_AccountDebitPoints_log.MSISDN+" "+loyalty_AccountDebitPoints_log.StatusDescription)
					slice := jobs[jobID].Result["Points Deducted"]
					jobsMu.Unlock()

					var successfulVal float64
					if len(slice) > 0 && slice[0] != "" {
						successfulVal, _ = strconv.ParseFloat(slice[0], 64)
					}
					pointsDeducted := successfulVal + loyalty_AccountDebitPoints_log.Debit_Amount
					jobsMu.Lock()
					jobs[jobID].Result["Points Deducted"] = []string{fmt.Sprintf("%.2f", pointsDeducted)}
					jobsMu.Unlock()
				}

			}

			Uc.Write_Loyalty_Event_Log(Loyalty_Event_Log{
				Event_User:         validated_Headers.AppLogin,
				Event_Time:         time.Now(),
				Event_AffectedType: "Bulk Points Deduction",
				Event_ActionType:   "Add",
				Event_Description:  "",
				Event_Entry_Before: nil,
				Event_Entry_After:  jobs[jobID].Result,
			})
			jobsMu.Lock()
			jobs[jobID].Status = "completed"
			jobsMu.Unlock()
		}(jobID, file)

		sr.Data = map[string]interface{}{
			"jobID": jobID,
		}
	}
	//successful response
	processedMu.Lock()
	processed = nil
	processedMu.Unlock()
	sr.Status = "successful"
	sr.StatusCode = http.StatusOK
	sr.StatusDescription = ""
	sr.ErrorDescription = ""
	Uc.HTTP_API_Standard_response(w, r, sr, true)
}
func (Uc *UserControl) HTTP_Bulk_Loyalty_Points_Deduction_Progress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["jobID"]
	if jobID == "" {
		http.Error(w, "missing jobID", http.StatusBadRequest)
		return
	}
	jobsMu.Lock()
	job, ok := jobs[jobID]
	jobsMu.Unlock()

	if !ok {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

// func (Uc *UserControl) HTTP_Loyalty_Redemption_log(w http.ResponseWriter, r *http.Request) {
// 	var sr API_Standard_response
// 	//**fill response source detail
// 	SourceIp, _ := GetRequestIP(r)
// 	sr.SourceIP = SourceIp
// 	sr.Login = r.Header.Get("Login")
// 	sr.SourceApp = r.Header.Get("SourceApp")
// 	sr.AccessKey = r.URL.Path
// 	sr.AccessMethod = r.Method
// 	sr.HostId = Configuration.HostId
// 	sr.ReceiveDate = time.Now()
// 	sr.TransactionType = "Loyalty Redemption log"

// 	method := r.Method
// 	switch method {
// 	case "GET":
// 		sr.TransactionType = sr.TransactionType + " - Read"
// 		MSISDN := r.URL.Query().Get("MSISDN")
// 		logs, err := Uc.Customer_Loyalty_Account_GetRedemptionlogs(MSISDN)
// 		if err != nil {
// 			sr.Status = "failed"
// 			sr.StatusCode = http.StatusBadRequest
// 			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get data"
// 			sr.ErrorDescription = err.Error()
// 			Uc.HTTP_API_Standard_response(w, r, sr, false)
// 			return
// 		}
// 		sr.Data = logs
// 	}
// 	//successful response
// 	sr.Status = "successful"
// 	sr.StatusCode = http.StatusOK
// 	sr.StatusDescription = ""
// 	sr.ErrorDescription = ""
// 	Uc.HTTP_API_Standard_response(w, r, sr, true)
// }

func (Uc *UserControl) HTTP_Loyalty_AccountDebitPoints_log(w http.ResponseWriter, r *http.Request) {
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
		sr.TransactionType = "Loyalty Redemption log"
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
		data, err := Uc.Customer_Loyalty_Account_GetDebitPoints_log(startDateDD, endDateDD, MSISDN, "")
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

func (Uc *UserControl) HTTP_Loyalty_AccountCreditPoints_log(w http.ResponseWriter, r *http.Request) {
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
		sr.TransactionType = "Loyalty Credit log"
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
		data, err := Uc.Customer_Loyalty_Account_GetCreditPoints_log(startDateDD, endDateDD, MSISDN, "")
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

func (Uc *UserControl) HTTP_Loyalty_AccountRedemptionPoints_log(w http.ResponseWriter, r *http.Request) {
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
		sr.TransactionType = "Loyalty Redemption log"
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
		data, err := Uc.Customer_Loyalty_Account_GetRedemptionPoints_log(startDateDD, endDateDD, MSISDN, "")
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

func (Uc *UserControl) HTTP_Loyalty_AccountExpiryPoints_log(w http.ResponseWriter, r *http.Request) {
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
		sr.TransactionType = "Loyalty Expiry log"
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
		data, err := Uc.Customer_Loyalty_Account_GetExpiryPoints_log(startDateDD, endDateDD, MSISDN, "")
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

func (Uc *UserControl) HTTP_Loyalty_AccountLevelChangePoints_log(w http.ResponseWriter, r *http.Request) {
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
		sr.TransactionType = "Loyalty Level Change log"
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
		data, err := Uc.Customer_Loyalty_Account_GetLevelChange_log(startDateDD, endDateDD, MSISDN, "")
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

func (Uc *UserControl) HTTP_Loyalty_AccountEvents_log(w http.ResponseWriter, r *http.Request) {
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
		sr.TransactionType = "Loyalty Events log"
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
		data, err := Uc.Customer_Loyalty_Account_GetEvents_log(startDateDD, endDateDD, MSISDN, "")
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
func (Uc *UserControl) HTTP_Loyalty_logs(w http.ResponseWriter, r *http.Request) {
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
		sr.TransactionType = "Loyalty logs"
		//Filter := r.URL.Query().Get("filter")
		//ServiceType := r.URL.Query().Get("ServiceType")
		// LimitStr := r.URL.Query().Get("Limit")
		// PageStr := r.URL.Query().Get("Page")

		MSISDN := r.URL.Query().Get("MSISDN")
		Type := r.URL.Query().Get("Type")
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
		data, err := Uc.Customer_Loyalty_Account_Getlogs(Type, startDateDD, endDateDD, MSISDN, "")
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
		key := request.Key
		// Check existence
		_, subExists := Map_Subscribers.CheckThenGet(key)
		loyaltyEntry, loyaltyExists := Map_Customer_Loyalty_Account.CheckThenGet(key)
		if !subExists && !loyaltyExists {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": user does not exist in lendme or loyalty"
			sr.ErrorDescription = " user does not exist in lendme or loyalty"
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}

		// Delete from subscriber system if exists
		if subExists {
			if err := Uc.Subscriber_Delete(key); err != nil {
				fmt.Println("Error deleting subscriber:", err)
			}
			_, exits := Map_Lendme_Customer_Exclusion.CheckThenGet(key)
			if exits {
				Map_Lendme_Customer_Exclusion.Delete(key)
			}
		}
		// Delete from loyalty system if exists
		if loyaltyExists {
			entry, ok := loyaltyEntry.(Customer_Loyalty_Account)
			if !ok {
				sr.Status = "failed"
				sr.StatusCode = http.StatusBadRequest
				sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to delete"
				sr.ErrorDescription = "failed to cast entry to Customer_Loyalty_Account"
				Uc.HTTP_API_Standard_response(w, r, sr, false)
				return
			}
			// Add his points to the available points pool
			Uc.Loyalty_Governance_Redeem_Points_Debit(entry.Available_Points, true)
			//delete loyalty points monthly wallets
			for _, pointDetailKey := range entry.Points_Detail_Keys {
				Map_Customer_Loyalty_Account_Points_Detail.Delete(pointDetailKey)
			}
			err = Uc.Customer_Loyalty_Account_Delete(sr.Login, key)
			if err != nil {
				sr.Status = "failed"
				sr.StatusCode = http.StatusBadRequest
				sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to delete"
				sr.ErrorDescription = err.Error()
				Uc.HTTP_API_Standard_response(w, r, sr, false)
				return
			}
			_, exits := Map_Customer_Exclusion.CheckThenGet(key)
			if exits {
				Map_Customer_Exclusion.Delete(key)
			}
			Uc.Write_Loyalty_Account_Churned_log(entry)
		}
		LiveFeedCounters.With(prometheus.Labels{"Stream": request.EventSource, "Type": "Chrun", "Description": "Chrun"}).Inc()
	}
	//successful response
	sr.Status = "successful"
	sr.StatusCode = http.StatusOK
	sr.StatusDescription = ""
	sr.ErrorDescription = ""
	Uc.HTTP_API_Standard_response(w, r, sr, true)
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
		var request Loyalty_AccountCreditPoints_Request
		err = json.Unmarshal(body, &request)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		validated_Headers := Uc.Validate_Headers(r)
		var loyalty_AccountCreditPoints_log Loyalty_AccountCreditPoints_log
		Uc.Loyalty_AccountCreditPoints(&validated_Headers, request, &loyalty_AccountCreditPoints_log)
		if loyalty_AccountCreditPoints_log.StatusCode == http.StatusBadRequest {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = "failed to add request"
			sr.ErrorDescription = loyalty_AccountCreditPoints_log.ErrorDescription
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
