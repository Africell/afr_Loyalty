package Lendme

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var chan_SubQueueExecution_controler = make(chan int, 10)
var chan_SubDump_EXECQueue = make(chan Sub_Update_Request, 1000)

func (Uc *UserControl) Import_Subscribers_Dump_LineByLine(FileName string) (err error) {
	log.Println("***subscribers dump importing process started ***")
	//Section 1: open file for reading
	if FileName == "" {
		return errors.New("subscribers dump file name cannot be empty")
	}
	FileName = Configuration.ARPU_File_Path + FileName
	file, err := os.Open(FileName)
	if err != nil {
		log.Println("error opening subscribers dump file " + FileName + ": " + err.Error())
		return err
	}
	defer file.Close()

	// Section 2: read lines
	r := bufio.NewReader(file)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				log.Println("***subscribers dump importing process finished ***")
				return nil
			} else {
				log.Println("error reading line from subscribers ARPU file: " + err.Error())
				continue
			}
		}
		if len(line) > 0 {
			//log.Println(string(line))
			result := strings.Split(string(line), ",")
			if len(result) == 11 {
				//file fields: MSISDN, Cos ID, Enrol Date, Last Credited, Loyalty Status, Credit Limit, Revenue,
				//	Recharge, Last Recharge date, Dealer bundle recharge, Last Dealer bundle recharge date

				Credit_Limit_str := result[5]
				if Credit_Limit_str == "" {
					Credit_Limit_str = "0"
				}
				Credit_Limit, err := strconv.ParseFloat(Credit_Limit_str, 64)
				if err != nil {
					log.Println("error converting Credit_Limit: ", err)
					continue
				}
				ARPU_Amount_str := result[6]
				if ARPU_Amount_str == "" {
					ARPU_Amount_str = "0"
				}
				ARPU_Amount, err := strconv.ParseFloat(ARPU_Amount_str, 64)
				if err != nil {
					log.Println("error converting ARPU_Amount: ", err)
					continue
				}
				var firstUsed time.Time
				if len(result[2]) == 10 {
					firstUsed, _ = time.Parse("2006-01-02", result[2])
				}
				var lastCredit time.Time
				if len(result[3]) == 10 {
					lastCredit, _ = time.Parse("2006-01-02", result[3])
				}
				Recharge_amnt_str := result[7]
				if Recharge_amnt_str == "" {
					Recharge_amnt_str = "0"
				}
				Recharge_amnt, err := strconv.ParseFloat(Recharge_amnt_str, 64)
				if err != nil {
					log.Println("error converting Recharge_amnt: ", err)
					continue
				}
				var lastRecharge time.Time
				if len(result[8]) == 10 {
					lastRecharge, _ = time.Parse("2006-01-02", result[8])
				}
				dealerPurchase_amnt_str := result[9]
				if dealerPurchase_amnt_str == "" {
					dealerPurchase_amnt_str = "0"
				}
				dealerPurchase_amnt, err := strconv.ParseFloat(dealerPurchase_amnt_str, 64)
				if err != nil {
					log.Println("error converting dealerPurchase_amnt: ", err)
					continue
				}
				var lastdealerPurchase time.Time
				if len(result[10]) == 10 {
					lastdealerPurchase, _ = time.Parse("2006-01-02", result[10])
				}
				request := Sub_Update_Request{
					MSISDN:                           result[0],
					COS:                              result[1],
					First_Used:                       firstUsed,
					Last_Credit:                      lastCredit,
					Loyalty_Status:                   result[4],
					Credit_Limit:                     Credit_Limit,
					ARPU_Amount:                      ARPU_Amount,
					Recharge:                         Recharge_amnt,
					Last_Recharge_Date:               lastRecharge,
					Dealer_Bundle_Purchase:           dealerPurchase_amnt,
					Last_Dealer_Bundle_Purchase_Date: lastdealerPurchase,
				}
				//log.Println("fill request:", request)
				chan_SubDump_EXECQueue <- request
			}
		}
	}
}

func (Uc *UserControl) SubQueueExecution() {
	for {
		select {
		case msg_sub := <-chan_SubDump_EXECQueue:
			chan_SubQueueExecution_controler <- 1
			go Uc.Subscriber_Update(msg_sub)
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (Uc *UserControl) Subscriber_Update(request Sub_Update_Request) {
	//log.Println("request:", request)
	//check key if filled and if already used
	if request.MSISDN == "" {
		<-chan_SubQueueExecution_controler
		return
	}
	//check if key already used
	var subscriber Subscriber
	subscriber_na, exits := Map_Subscribers.CheckThenGet(request.MSISDN)
	if exits {
		var ok bool
		subscriber, ok = subscriber_na.(Subscriber)
		if !ok {
			<-chan_SubQueueExecution_controler
			return
		}
	} else {
		subscriber.Subscriber_Id = Map_AutoIncrement.GetNextAI("Subscriber-Id")
		subscriber.Key = request.MSISDN
		subscriber.Add_Date = time.Now()
	}
	subscriber.Last_ProfileUpdate_date = time.Now()

	subscriber.COS = request.COS
	subscriber.FirstUse_date = request.First_Used
	subscriber.Last_Credit = request.Last_Credit
	subscriber.IN_Loyalty_Status = request.Loyalty_Status
	subscriber.IN_Credit_Limit = request.Credit_Limit
	subscriber.ARPU = request.ARPU_Amount
	subscriber.Recharge = request.Recharge
	subscriber.Last_Recharge_Date = request.Last_Recharge_Date
	subscriber.Dealer_Bundle_Purchase = request.Dealer_Bundle_Purchase
	subscriber.Last_Dealer_Bundle_Purchase_Date = request.Last_Dealer_Bundle_Purchase_Date

	LastRechargeDate := request.Last_Recharge_Date
	if LastRechargeDate.Before(request.Last_Dealer_Bundle_Purchase_Date) {
		LastRechargeDate = request.Last_Dealer_Bundle_Purchase_Date
	}
	subscriber.Credit_Limit_Scheme, subscriber.NotElligibleReason = Uc.Credit_Limit_Scheme_Selection((request.Recharge + request.Dealer_Bundle_Purchase), request.First_Used, LastRechargeDate)
	if subscriber.Credit_Limit_Scheme != "" {
		subscriber.IsLendmeEligible = true
	} else {
		subscriber.IsLendmeEligible = false
	}
	//add to cache and DB
	Map_Subscribers.Put(subscriber.Key, subscriber)
	<-chan_SubQueueExecution_controler
}
