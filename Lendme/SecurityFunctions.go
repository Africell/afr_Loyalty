package Lendme

import (
	"net/http"
	"strings"
	"time"
)

func (Uc *UserControl) ValidateJWEToken(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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
		Authorization := r.Header.Get("Authorization")
		AuthorizationSplit := strings.Split(Authorization, " ")
		if Authorization == "" || len(AuthorizationSplit) < 2 {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusUnauthorized) + ":" + " Authorization information must be provided"
			sr.ErrorDescription = http.StatusText(http.StatusUnauthorized) + ":" + " Authorization information must be provided"
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		if AuthorizationSplit[0] != "Bearer" {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusUnauthorized) + ":" + " Authorization type must be Bearer"
			sr.ErrorDescription = http.StatusText(http.StatusUnauthorized) + ":" + " Authorization type must be Bearer"
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		if len(AuthorizationSplit[1]) < 3 {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusUnauthorized) + ":" + " token is empty"
			sr.ErrorDescription = http.StatusText(http.StatusUnauthorized) + ":" + " token is empty"
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return

		}

		port, err := GetDestinationPortFromRequest(r)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusUnauthorized) + ":" + err.Error()
			sr.ErrorDescription = http.StatusText(http.StatusUnauthorized) + ":" + err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		if port == Configuration.HttpAppServicePort {
			sra, err := Uc.AppAUC.AUCClient.ValidateToken(AuthorizationSplit[1])
			if err != nil {
				sr.Status = sra.Status
				sr.StatusCode = sra.StatusCode
				sr.ErrorDescription = sra.ErrorDescription
				Uc.HTTP_API_Standard_response(w, r, sr, false)
				return
			}
			r.Header.Set("Login", sra.Login)
			r.Header.Set("DestinationPort", port)
			h.ServeHTTP(w, r)
		} else if port == Configuration.HttpOKAPIServicePort {
			sra, err := Uc.OKAPIAUC.AUCClient.ValidateToken(AuthorizationSplit[1])
			if err != nil {
				sr.Status = sra.Status
				sr.StatusCode = sra.StatusCode
				sr.ErrorDescription = sra.ErrorDescription
				Uc.HTTP_API_Standard_response(w, r, sr, false)
				return
			}
			r.Header.Set("Login", sra.Login)
			r.Header.Set("DestinationPort", port)
			h.ServeHTTP(w, r)
		} else {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusUnauthorized) + ": Invalid Destination Port"
			sr.ErrorDescription = http.StatusText(http.StatusUnauthorized) + ": Invalid Destination Port"
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
	}
}
