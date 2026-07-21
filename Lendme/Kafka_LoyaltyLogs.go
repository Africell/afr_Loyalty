package Lendme

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Per-type topics (one per log type, matching the reference's topic-per-stream style).
const (
	Topic_Loyalty_Redemption_log          = "Loyalty_Redemption_log"
	Topic_Loyalty_AccountCreditPoints_log = "Loyalty_AccountCreditPoints_log"
	Topic_Loyalty_AccountDebitPoints_log  = "Loyalty_AccountDebitPoints_log"
	Topic_Loyalty_Status_log              = "Loyalty_Status_log"
)

// LoyaltyKafkaTopics is consumed by ensureTopics() (Kafka_Producer.go) when
// Configuration.Kafka_Events.CreateTopicsOnStartup is true.
var LoyaltyKafkaTopics = []string{
	Topic_Loyalty_Redemption_log,
	Topic_Loyalty_AccountCreditPoints_log,
	Topic_Loyalty_AccountDebitPoints_log,
	Topic_Loyalty_Status_log,
}

// Two channels per log type: a buffered data channel + a controller (semaphore) channel.
var chan_Loyalty_Redemption_log = make(chan Loyalty_Redemption_log, 500)
var chan_Loyalty_Redemption_log_controler = make(chan int, 50)

var chan_Loyalty_AccountCreditPoints_log = make(chan Loyalty_AccountCreditPoints_log, 500)
var chan_Loyalty_AccountCreditPoints_log_controler = make(chan int, 50)

var chan_Loyalty_AccountDebitPoints_log = make(chan Loyalty_AccountDebitPoints_log, 500)
var chan_Loyalty_AccountDebitPoints_log_controler = make(chan int, 50)

var chan_Loyalty_Status_log = make(chan Loyalty_Status_log, 500)
var chan_Loyalty_Status_log_controler = make(chan int, 50)

// enqueue_* perform a NON-BLOCKING hand-off to the Kafka channel. If the buffer
// is full (e.g. the broker is down and backing up), the record is dropped with a
// log line rather than blocking the calling request handler. Used by the five
// target functions and the HTTP failure handlers.
func enqueue_Loyalty_Redemption_log(rec Loyalty_Redemption_log) {
	select {
	case chan_Loyalty_Redemption_log <- rec:
	default:
		log.Println("kafka enqueue dropped (buffer full): Loyalty_Redemption_log MSISDN=" + rec.MSISDN)
	}
}

func enqueue_Loyalty_AccountCreditPoints_log(rec Loyalty_AccountCreditPoints_log) {
	select {
	case chan_Loyalty_AccountCreditPoints_log <- rec:
	default:
		log.Println("kafka enqueue dropped (buffer full): Loyalty_AccountCreditPoints_log MSISDN=" + rec.MSISDN)
	}
}

func enqueue_Loyalty_AccountDebitPoints_log(rec Loyalty_AccountDebitPoints_log) {
	select {
	case chan_Loyalty_AccountDebitPoints_log <- rec:
	default:
		log.Println("kafka enqueue dropped (buffer full): Loyalty_AccountDebitPoints_log MSISDN=" + rec.MSISDN)
	}
}

func enqueue_Loyalty_Status_log(rec Loyalty_Status_log) {
	select {
	case chan_Loyalty_Status_log <- rec:
	default:
		log.Println("kafka enqueue dropped (buffer full): Loyalty_Status_log MSISDN=" + rec.MSISDN)
	}
}

// loyaltyKafkaPushOptions matches the tuning used by KafkaPush_EventMetrics in the reference.
func loyaltyKafkaPushOptions() *PushOptions {
	return &PushOptions{
		OverallTimeout: 5 * time.Second,
		AttemptTimeout: 2 * time.Second,
		MaxRetries:     5,
		BaseBackoff:    100 * time.Millisecond,
		MaxBackoff:     2 * time.Second,
		JitterFraction: 0.25,
	}
}

// Loyalty_Kafka_Process drains every loyalty log channel and dispatches a push
// goroutine per message. Launched once at startup (afr_Lendme_main.go).
func (Uc *UserControl) Loyalty_Kafka_Process() {
	log.Println("Loyalty Kafka producer process started")
	for {
		select {
		case msg := <-chan_Loyalty_Redemption_log:
			chan_Loyalty_Redemption_log_controler <- 1
			go Uc.KafkaPush_Loyalty_Redemption_log(msg)
		case msg := <-chan_Loyalty_AccountCreditPoints_log:
			chan_Loyalty_AccountCreditPoints_log_controler <- 1
			go Uc.KafkaPush_Loyalty_AccountCreditPoints_log(msg)
		case msg := <-chan_Loyalty_AccountDebitPoints_log:
			chan_Loyalty_AccountDebitPoints_log_controler <- 1
			go Uc.KafkaPush_Loyalty_AccountDebitPoints_log(msg)
		case msg := <-chan_Loyalty_Status_log:
			chan_Loyalty_Status_log_controler <- 1
			go Uc.KafkaPush_Loyalty_Status_log(msg)
		default:
			<-time.After(10 * time.Millisecond)
		}
	}
}

func (Uc *UserControl) KafkaPush_Loyalty_Redemption_log(payload Loyalty_Redemption_log) {
	defer func() { <-chan_Loyalty_Redemption_log_controler }()

	if Uc.KafkaEventClient == nil {
		log.Println("KafkaPush_Loyalty_Redemption_log: kafka client is nil, skipping")
		return
	}
	valueBytes, err := bson.Marshal(payload)
	if err != nil {
		log.Println("KafkaPush_Loyalty_Redemption_log failed to marshal:", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	if err := Uc.KafkaEventClient.Push(ctx, []byte(payload.MSISDN), valueBytes, Topic_Loyalty_Redemption_log, loyaltyKafkaPushOptions()); err != nil {
		log.Println("KafkaPush_Loyalty_Redemption_log failed to push message:", err)
		return
	}
}

func (Uc *UserControl) KafkaPush_Loyalty_AccountCreditPoints_log(payload Loyalty_AccountCreditPoints_log) {
	defer func() { <-chan_Loyalty_AccountCreditPoints_log_controler }()

	if Uc.KafkaEventClient == nil {
		log.Println("KafkaPush_Loyalty_AccountCreditPoints_log: kafka client is nil, skipping")
		return
	}
	valueBytes, err := bson.Marshal(payload)
	if err != nil {
		log.Println("KafkaPush_Loyalty_AccountCreditPoints_log failed to marshal:", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	if err := Uc.KafkaEventClient.Push(ctx, []byte(payload.MSISDN), valueBytes, Topic_Loyalty_AccountCreditPoints_log, loyaltyKafkaPushOptions()); err != nil {
		log.Println("KafkaPush_Loyalty_AccountCreditPoints_log failed to push message:", err)
		return
	}
}

func (Uc *UserControl) KafkaPush_Loyalty_AccountDebitPoints_log(payload Loyalty_AccountDebitPoints_log) {
	defer func() { <-chan_Loyalty_AccountDebitPoints_log_controler }()

	if Uc.KafkaEventClient == nil {
		log.Println("KafkaPush_Loyalty_AccountDebitPoints_log: kafka client is nil, skipping")
		return
	}
	valueBytes, err := bson.Marshal(payload)
	if err != nil {
		log.Println("KafkaPush_Loyalty_AccountDebitPoints_log failed to marshal:", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	if err := Uc.KafkaEventClient.Push(ctx, []byte(payload.MSISDN), valueBytes, Topic_Loyalty_AccountDebitPoints_log, loyaltyKafkaPushOptions()); err != nil {
		log.Println("KafkaPush_Loyalty_AccountDebitPoints_log failed to push message:", err)
		return
	}
}

func (Uc *UserControl) KafkaPush_Loyalty_Status_log(payload Loyalty_Status_log) {
	defer func() { <-chan_Loyalty_Status_log_controler }()

	if Uc.KafkaEventClient == nil {
		log.Println("KafkaPush_Loyalty_Status_log: kafka client is nil, skipping")
		return
	}
	valueBytes, err := bson.Marshal(payload)
	if err != nil {
		log.Println("KafkaPush_Loyalty_Status_log failed to marshal:", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	if err := Uc.KafkaEventClient.Push(ctx, []byte(payload.MSISDN), valueBytes, Topic_Loyalty_Status_log, loyaltyKafkaPushOptions()); err != nil {
		log.Println("KafkaPush_Loyalty_Status_log failed to push message:", err)
		return
	}
}
