package Lendme

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	PrometheusRegistry = prometheus.NewRegistry()
	//PortStatus.With(prometheus.Labels{"DestinationHost": "HHHH"}).Set(1)
	PortStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "PortStatus",
			Help: "1: port up, 2: port down",
		},
		[]string{"DestinationHost"},
	)

	TransactionsCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "TransactionsCount",
			Help: "TransactionsCount",
		},
		[]string{"Stream", "Type", "Description"},
	)
	TransactionsTotalAmount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "TransactionsTotalAmount",
			Help: "TransactionsTotalAmount",
		},
		[]string{"Stream", "Type", "Description"},
	)

	LiveFeedStream = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "LiveFeedStream",
			Help: "LiveFeedStream",
		},
		[]string{"Type"},
	)

	SubscriptionRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "SubscriptionRequestCount",
			Help: "SubscriptionRequestCount",
		},
		[]string{"RequestType", "Status"},
	)

	NotificationCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "NotificationCount",
			Help: "NotificationCount",
		},
		[]string{"Type", "Status"},
	)
)

//NotificationCount.With(prometheus.Labels{"Type": "", "Status":""}).Inc()
//TriviaQuizChargedAmount.With(prometheus.Labels{"Description": "desc"}).Add()

func Init_Prometheus_Metrics() {
	log.Println("Init Prometheus metrics")
	PrometheusRegistry.Register(PortStatus)
	PrometheusRegistry.Register(TransactionsCount)
	PrometheusRegistry.Register(TransactionsTotalAmount)
	PrometheusRegistry.Register(LiveFeedStream)
	PrometheusRegistry.Register(SubscriptionRequestCount)
	PrometheusRegistry.Register(NotificationCount)
}

func Reset_Prometheus_Metrics() {
	exec := 0
	for range time.Tick(time.Second * 1) {
		_CurrentDateTime := time.Now()
		_hr, _mi, _se := _CurrentDateTime.Clock()
		if _hr == 0 {
			if _mi == 0 {
				if _se < 10 {
					if exec == 0 {
						TransactionsCount.Reset()
						TransactionsTotalAmount.Reset()
						SubscriptionRequestCount.Reset()
						NotificationCount.Reset()
						LiveFeedStream.Reset()
						exec = 1
					}
				} else {
					if exec == 1 {
						exec = 0
					}
				}
			}
		}
	}
}

func CustomPrometheusHandler() http.Handler {
	return promhttp.HandlerFor(
		PrometheusRegistry,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: false,
		},
	)
}
