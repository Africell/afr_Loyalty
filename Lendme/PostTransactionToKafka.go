package Lendme

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type EVC_Recharge_request struct {
	TransID                string  `bson:"TransID" json:"TransID"`
	TransStatus            string  `bson:"TransStatus" json:"TransStatus"` //successful, failed
	TransStatusDescription string  `bson:"TransStatusDescription" json:"TransStatusDescription"`
	DealerMSISDN           string  `bson:"DealerMSISDN" json:"DealerMSISDN"`
	DealerName             string  `bson:"DealerName" json:"DealerName"`
	DealerClosingBalance   float64 `bson:"DealerClosingBalance" json:"DealerClosingBalance"`
	TargetMSISDN           string  `bson:"TargetMSISDN" json:"TargetMSISDN"`
	Amount                 float64 `bson:"Amount" json:"Amount"`
	GSMLocation            string  `bson:"GSMLocation" json:"GSMLocation"`
}

func Post_EVC_Recharge_ToKafka(load EVC_Recharge_request) {
	HTTPPostToKafka.With(prometheus.Labels{"Description": "EVC_Recharge", "Status": "Attempt", "Error": ""}).Inc()
	TargetFullURL := "http://10.250.9.180:9501/HTTP_EVC_Recharge_Producer/"
	method := "POST"
	client := &http.Client{
		Timeout: 10 * time.Second, //if not reachable, request will time out after XX sec
	}
	load_byte, err := json.Marshal(load)
	if err != nil {
		HTTPPostToKafka.With(prometheus.Labels{"Description": "EVC_Recharge", "Status": "Failed", "Error": "Load Marshal problem"}).Inc()
		return
	}
	req, err := http.NewRequest(method, TargetFullURL, bytes.NewBuffer(load_byte))
	if err != nil {
		log.Println("error in Post_Consumption_ToTagetURL: ", err.Error())
		HTTPPostToKafka.With(prometheus.Labels{"Description": "EVC_Recharge", "Status": "Failed", "Error": err.Error()}).Inc()
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("SourceApp", "EMIS GW")
	req.Header.Add("Latitude", "0")
	req.Header.Add("Longitude", "0")
	// now POST it
	resp, err := client.Do(req)
	defer req.Body.Close()
	if err != nil {
		HTTPPostToKafka.With(prometheus.Labels{"Description": "EVC_Recharge", "Status": "Failed", "Error": err.Error()}).Inc()
		return
	}
	if resp.StatusCode == http.StatusOK {
		HTTPPostToKafka.With(prometheus.Labels{"Description": "EVC_Recharge", "Status": "Successful", "Error": ""}).Inc()
	} else {
		HTTPPostToKafka.With(prometheus.Labels{"Description": "EVC_Recharge", "Status": "Failed", "Error": "StatusCode " + strconv.Itoa(resp.StatusCode)}).Inc()
	}
}

type EVC_BundlePurchase_request struct {
	TransID                string  `bson:"TransID" json:"TransID"`
	TransStatus            string  `bson:"TransStatus" json:"TransStatus"` //successful, failed
	TransStatusDescription string  `bson:"TransStatusDescription" json:"TransStatusDescription"`
	DealerMSISDN           string  `bson:"DealerMSISDN" json:"DealerMSISDN"`
	DealerName             string  `bson:"DealerName" json:"DealerName"`
	DealerClosingBalance   float64 `bson:"DealerClosingBalance" json:"DealerClosingBalance"`
	TargetMSISDN           string  `bson:"TargetMSISDN" json:"TargetMSISDN"`
	BundleId               string  `bson:"BundleId" json:"BundleId"`
	BundleName             string  `bson:"BundleName" json:"BundleName"`
	BundleType             string  `bson:"BundleType" json:"BundleType"`
	BundleCategory         string  `bson:"BundleCategory" json:"BundleCategory"`
	BundleCost             float64 `bson:"BundleCost" json:"BundleCost"`
	BundleValidityInDays   float64 `bson:"BundleValidityInDays" json:"BundleValidityInDays"`
	GSMLocation            string  `bson:"GSMLocation" json:"GSMLocation"`
}

func Post_EVC_BundlePurchase_ToKafka(load EVC_BundlePurchase_request) {
	HTTPPostToKafka.With(prometheus.Labels{"Description": "EVC_BundlePurchase", "Status": "Attempt", "Error": ""}).Inc()
	TargetFullURL := "http://10.250.9.180:9501/HTTP_EVC_BundlePurchase_Producer/"
	method := "POST"
	client := &http.Client{
		Timeout: 10 * time.Second, //if not reachable, request will time out after XX sec
	}
	load_byte, err := json.Marshal(load)
	if err != nil {
		HTTPPostToKafka.With(prometheus.Labels{"Description": "EVC_BundlePurchase", "Status": "Failed", "Error": "Load Marshal problem"}).Inc()
		return
	}
	req, err := http.NewRequest(method, TargetFullURL, bytes.NewBuffer(load_byte))
	if err != nil {
		log.Println("error in Post_Consumption_ToTagetURL: ", err.Error())
		HTTPPostToKafka.With(prometheus.Labels{"Description": "EVC_BundlePurchase", "Status": "Failed", "Error": err.Error()}).Inc()
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("SourceApp", "EMIS GW")
	req.Header.Add("Latitude", "0")
	req.Header.Add("Longitude", "0")
	// now POST it
	resp, err := client.Do(req)
	defer req.Body.Close()
	if err != nil {
		HTTPPostToKafka.With(prometheus.Labels{"Description": "EVC_BundlePurchase", "Status": "Failed", "Error": err.Error()}).Inc()
		return
	}
	if resp.StatusCode == http.StatusOK {
		HTTPPostToKafka.With(prometheus.Labels{"Description": "EVC_BundlePurchase", "Status": "Successful", "Error": ""}).Inc()
	} else {
		HTTPPostToKafka.With(prometheus.Labels{"Description": "EVC_BundlePurchase", "Status": "Failed", "Error": "StatusCode " + strconv.Itoa(resp.StatusCode)}).Inc()
	}
}
