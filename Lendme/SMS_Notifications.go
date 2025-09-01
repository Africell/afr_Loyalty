package Lendme

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func Send_SMS(sender string, target string, text string) (_rErr error) {
	//check target
	if len(target) < Configuration.SMPP.MSISDN_Short_len {
		err := errors.New("error sending SMS: target is short ")
		return err
	}
	// target = Configuration.SMPP.CountryCodePrefix + target[(len(target)-Configuration.SMPP.MSISDN_Short_len):]
	//check if UAT and if target is in UAT pool
	if !Configuration.IsLoyaltyProduction {
		existsUAT := Map_Customer_UAT.Check(target)
		if !existsUAT {
			// err := errors.New("sending SMS blocked by UAT pool: Sender (" + Sender + "), Target (" + target + "), text (" + SMSText + ") ")
			// return err
			return
		}
	}

	//check if target is in do not disturb list
	existsDND := Map_Customer_DND.Check(target)
	if existsDND {
		err := errors.New("sending SMS blocked by DND: Sender (" + sender + "), Target (" + target + "), text (" + text + ") ")
		return err
	}
	if Configuration.SMPP.PrintLogs {
		log.Println("Sending SMS: Sender (" + sender + "), Target (" + target + "), text (" + text + ") ")
	}
	url := "http://" + Configuration.SMPP.IP + ":" + Configuration.SMPP.Port + "/?systemid=" + Configuration.SMPP.Login + "&password=" + url.QueryEscape(Configuration.SMPP.Password) + "&Originator=" + sender + "&dest_addr=" + target + "&msg_text=" + url.QueryEscape(text) + "&encoding=" + fmt.Sprint(Configuration.SMPP.Encoding) + "&ston=5&snpi=0&dton=1&registered_delivery=0"
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}
	//client := &http.Client{}
	client := &http.Client{
		Timeout: Configuration.SMPP.TimeOut * time.Second, //is SMSC not reachable request will time out after "TimeOut" sec
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			err := errors.New("failed to read SMS submit response body")
			return err
		}
		response_body := string(body)
		if len(response_body) < 11 {
			err := errors.New("error sending SMS: " + response_body)
			return err
		} else {
			if response_body[0:11] != "errorCode=0" { //errorCode=0&message_id=468008000124702 ==> successful response
				err := errors.New("error sending SMS: " + response_body + url)
				return err
			}
		}
		return nil
	} else {
		err := errors.New("error sending SMS (StatusCode: " + fmt.Sprint(resp.StatusCode) + ")")
		return err
	}
}
