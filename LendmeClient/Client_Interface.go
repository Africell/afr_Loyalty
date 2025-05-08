package LendmeClient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

func (GWClient *Lendme_Client) generic_http_call(request Generic_http_call_Request) (response Generic_http_call_Response, err error) {
	if len(request.Load) > 0 {
		request.Req, err = http.NewRequest(request.Method, request.Url, bytes.NewBuffer(request.Load))
	} else {
		request.Req, err = http.NewRequest(request.Method, request.Url, nil)
	}
	if err != nil {
		return
	}
	request.Req.Header.Set("Content-Type", "application/json")
	request.Req.Header.Set("Connection", "close")
	request.Req.Header.Set("Authorization", "Bearer "+request.Token)

	if len(request.QueryParameters) > 0 {
		q := request.Req.URL.Query()
		for qry_key, qry_val := range request.QueryParameters {
			q.Add(qry_key, qry_val)
		}
		request.Req.URL.RawQuery = q.Encode()
	}

	client := &http.Client{Timeout: GWClient.Timeout * time.Second}
	resp, err := client.Do(request.Req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	response.Body, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	response.Header = resp.Request.Header
	response.Statuscode = resp.StatusCode
	return response, nil
}

// **********************************************************************************************
// Lendme functions
// **********************************************************************************************
func (GWClient *Lendme_Client) Lendme_Subscriber_Get(MSISDN string) (response Lendme_Subscriber, err error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	http_req := Generic_http_call_Request{
		Url:             GWClient.Protocol + "://" + GWClient.Hostname + ":" + GWClient.LendmePort + "/" + GWClient.LendMeModule + "/" + GWClient.LendMeVersion + "/HTTP_Subscriber/",
		Method:          "GET",
		Token:           GWClient.S2S_AccessToken,
		QueryParameters: map[string]string{"Key": MSISDN},
	}
	response_generic, err := GWClient.generic_http_call(http_req)
	if err != nil {
		log.Println("generic_http_call failed : ", err)
		return response, err
	} else {
		err = json.Unmarshal(response_generic.Body, &response)
		if err != nil {
			log.Println("generic_http_call boby unmarshal error : ", err)
			return
		}
		if response_generic.Statuscode == http.StatusUnauthorized {
			srl, err := GWClient.AUC_client.Login(GWClient.AUC_client.S2S_Username, GWClient.AUC_client.S2S_Password) // try to get a new token using login if token is unauthenticated
			if err != nil {
				log.Println("AUC init - FAILED :: ", err)
				return response, err
			}
			rt, err := GWClient.AUC_client.RefreshToken(srl.Token)
			if err != nil {
				log.Println("AUC init - FAILED :: ", err)
				return response, err
			} else {
				GWClient.AUC_client.S2S_AccessToken = rt.Token // save token global variable to re-use
				GWClient.S2S_AccessToken = rt.Token            // save token global variable to re-use
				response_generic, err = GWClient.generic_http_call(http_req)
				if err != nil {
					log.Println("generic_http_call error : ", err)
					return response, err
				}
				err = json.Unmarshal(response_generic.Body, &response)
				if err != nil {
					log.Println("generic_http_call boby unmarshal error : ", err)
					return response, err
				}
				return response, err
			}
		}
	}
	return
}

// **********************************************************************************************
// Loyalty functions
// **********************************************************************************************
func (GWClient *Lendme_Client) Loyalty_Account_Get(MSISDN string) (response Customer_Loyalty_Account, err error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	http_req := Generic_http_call_Request{
		Url:             GWClient.Protocol + "://" + GWClient.Hostname + ":" + GWClient.LoyaltyPort + "/" + GWClient.LoyaltyModule + "/" + GWClient.LoyaltyVersion + "/HTTP_Customer_Loyalty_Account/",
		Method:          "GET",
		Token:           GWClient.S2S_AccessToken,
		QueryParameters: map[string]string{"Key": MSISDN},
	}
	response_generic, err := GWClient.generic_http_call(http_req)
	if err != nil {
		log.Println("generic_http_call failed : ", err)
		return response, err
	} else {
		err = json.Unmarshal(response_generic.Body, &response)
		if err != nil {
			log.Println("generic_http_call boby unmarshal error : ", err)
			return
		}
		if response_generic.Statuscode == http.StatusUnauthorized {
			srl, err := GWClient.AUC_client.Login(GWClient.AUC_client.S2S_Username, GWClient.AUC_client.S2S_Password) // try to get a new token using login if token is unauthenticated
			if err != nil {
				log.Println("AUC init - FAILED :: ", err)
				return response, err
			}
			rt, err := GWClient.AUC_client.RefreshToken(srl.Token)
			if err != nil {
				log.Println("AUC init - FAILED :: ", err)
				return response, err
			} else {
				GWClient.AUC_client.S2S_AccessToken = rt.Token // save token global variable to re-use
				GWClient.S2S_AccessToken = rt.Token            // save token global variable to re-use
				response_generic, err = GWClient.generic_http_call(http_req)
				if err != nil {
					log.Println("generic_http_call error : ", err)
					return response, err
				}
				err = json.Unmarshal(response_generic.Body, &response)
				if err != nil {
					log.Println("generic_http_call boby unmarshal error : ", err)
					return response, err
				}
				return response, err
			}
		}
	}
	return
}

func (GWClient *Lendme_Client) Loyalty_Account_DebitPoints(request Loyalty_AccountDebitPoints_Request) (response Loyalty_AccountDebitPoints_log, err error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	request_byte, err := json.Marshal(request)
	if err == nil {
		log.Println("generic_http_call request marshal error : ", err)
		return response, err
	}
	http_req := Generic_http_call_Request{
		Url:    GWClient.Protocol + "://" + GWClient.Hostname + ":" + GWClient.LoyaltyPort + "/" + GWClient.LoyaltyModule + "/" + GWClient.LoyaltyVersion + "/HTTP_Customer_Loyalty_Account_DebitPoints/",
		Method: "PUT",
		Token:  GWClient.S2S_AccessToken,
		Load:   request_byte,
	}
	response_generic, err := GWClient.generic_http_call(http_req)
	if err != nil {
		log.Println("generic_http_call failed : ", err)
		return response, err
	} else {
		err = json.Unmarshal(response_generic.Body, &response)
		if err != nil {
			log.Println("generic_http_call boby unmarshal error : ", err)
			return
		}
		if response_generic.Statuscode == http.StatusUnauthorized {
			srl, err := GWClient.AUC_client.Login(GWClient.AUC_client.S2S_Username, GWClient.AUC_client.S2S_Password) // try to get a new token using login if token is unauthenticated
			if err != nil {
				log.Println("AUC init - FAILED :: ", err)
				return response, err
			}
			rt, err := GWClient.AUC_client.RefreshToken(srl.Token)
			if err != nil {
				log.Println("AUC init - FAILED :: ", err)
				return response, err
			} else {
				GWClient.AUC_client.S2S_AccessToken = rt.Token // save token global variable to re-use
				GWClient.S2S_AccessToken = rt.Token            // save token global variable to re-use
				response_generic, err = GWClient.generic_http_call(http_req)
				if err != nil {
					log.Println("generic_http_call error : ", err)
					return response, err
				}
				err = json.Unmarshal(response_generic.Body, &response)
				if err != nil {
					log.Println("generic_http_call boby unmarshal error : ", err)
					return response, err
				}
				return response, err
			}
		}
	}
	return
}
