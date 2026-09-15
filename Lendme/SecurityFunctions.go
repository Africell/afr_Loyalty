package Lendme

import (
	"afr_auth_center/AuthCenter"
	"encoding/json"
	"fmt"
	"net"
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
		switch port {
		case Configuration.HttpAppLoyaltyServicePort:
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
			r.Header.Set("TokenType", sra.TokenType)
			h.ServeHTTP(w, r)
		case Configuration.HttpAppLoyaltyManagementPort:
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
			r.Header.Set("TokenType", sra.TokenType)
			h.ServeHTTP(w, r)
		case Configuration.HttpAppLoyaltyFeedPort:
			h.ServeHTTP(w, r)
		default:
			h.ServeHTTP(w, r)
		}
	}
}

func GetDestinationPortFromRequest(r *http.Request) (port string, err error) {
	//Derive the port from the actual listener that accepted the connection
	//(9280/9281/9282, each its own http.ListenAndServe) instead of a client-supplied
	//header, since X-Original-Host/Host can be spoofed by any client that reaches
	//this service directly.
	localAddr, ok := r.Context().Value(http.LocalAddrContextKey).(net.Addr)
	if !ok {
		return "", fmt.Errorf("no valid destination port")
	}
	_, port, err = net.SplitHostPort(localAddr.String())
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("no valid destination port")
	}
	return port, nil
}

func (Uc *UserControl) ValidateAccess_AUC(h http.HandlerFunc) http.HandlerFunc {
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

		// Split the path by "/"
		segments := strings.Split(r.URL.Path, "/")

		// Iterate over segments and find the one that contains "HTTP_"
		var httpSegment string
		for _, segment := range segments {
			if strings.Contains(segment, "HTTP_") {
				httpSegment = segment
				break
			}
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

		accReq := AuthCenter.ValidateAccess_request{
			User:         r.Header.Get("Login"),
			AccessKey:    httpSegment,
			AccessMethod: r.Method,
		}

		switch port {
		case Configuration.HttpAppLoyaltyServicePort:
			h.ServeHTTP(w, r)
		case Configuration.HttpAppLoyaltyManagementPort:
			sr_validate, err := Uc.OKAPIAUC.AUCClient.ValidateAccess(accReq)
			if err != nil {
				// HTTP reply
				w.WriteHeader(sr.StatusCode)
				json.NewEncoder(w).Encode(sr)
				return
			}
			if !sr_validate.Data.IsAccessAllowed {
				// HTTP reply
				sr.Status = "failed"
				sr.StatusCode = 400
				sr.StatusDescription = "Error: Not Allowed For this action, Please contact system Admin!"

				w.WriteHeader(sr.StatusCode)
				json.NewEncoder(w).Encode(sr)
				return
			}
			keys := make([]string, len(sr_validate.Data.AccessLevel))
			for k := range sr_validate.Data.AccessLevel {
				keys = append(keys, k)
			}
			r.Header.Set("AccessLevels", fmt.Sprint(keys))
			h.ServeHTTP(w, r)
		case Configuration.HttpAppLoyaltyFeedPort:
			h.ServeHTTP(w, r)
		default:
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusUnauthorized) + ": Invalid Destination Port"
			sr.ErrorDescription = http.StatusText(http.StatusUnauthorized) + ": Invalid Destination Port"
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
	}
}
