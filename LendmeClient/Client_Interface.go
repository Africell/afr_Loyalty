package LendmeClient

import (
	"afr_lendme/Lendme"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
)

func (cl *Lendme_Client) Get_Catalogue_WithBundleDetail(token string, Channel, Plan, Version, TargetLocation string) (response Lendme.API_Standard_response, err error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	url := cl.Protocol + "://" + cl.Hostname + ":" + cl.Port + "/" + cl.Module + "/" + cl.Version + "/HTTP_Catalogue_WithBundleDetail/"
	if Channel != "" {
		url = url + "?Channel=" + UrlQueryEscape(Channel)
	} else {
		url = url + "?Channel"
	}
	if Plan != "" {
		url = url + "&Plan=" + UrlQueryEscape(Plan)
	} else {
		url = url + "&Plan"
	}
	if Version != "" {
		url = url + "&Version=" + UrlQueryEscape(Version)
	} else {
		url = url + "&Version"
	}
	if TargetLocation != "" {
		url = url + "&TargetLocation=" + UrlQueryEscape(TargetLocation)
	} else {
		url = url + "&TargetLocation"
	}
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return response, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{Timeout: cl.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}
