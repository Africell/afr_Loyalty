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

	DailyImportSubsStats = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "DailyImportSubsStats",
			Help: "DailyImportSubsStats",
		},
		[]string{"IsElligble", "Reason", "Scheme"},
	)

	LendMeRequestsCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "LendMeRequestsCount",
			Help: "LendMeRequestsCount",
		},
		[]string{"Status", "Reason", "Scheme"},
	)
	LendMeRequestsAmount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "LendMeRequestsAmount",
			Help: "LendMeRequestsAmount",
		},
		[]string{"Status", "Reason", "Scheme"},
	)

	// TransactionsTotalAmount = prometheus.NewGaugeVec(
	// 	prometheus.GaugeOpts{
	// 		Name: "TransactionsTotalAmount",
	// 		Help: "TransactionsTotalAmount",
	// 	},
	// 	[]string{"Stream", "Type", "Description"},
	// )
)

//DailyImportSubsStats.With(prometheus.Labels{"IsElligble":"", "Reason":"", "Scheme":""}).Inc()
//TriviaQuizChargedAmount.With(prometheus.Labels{"Description": "desc"}).Add()

func Init_Prometheus_Metrics() {
	log.Println("Init Prometheus metrics")
	PrometheusRegistry.Register(PortStatus)
	PrometheusRegistry.Register(DailyImportSubsStats)
	PrometheusRegistry.Register(LendMeRequestsCount)
	PrometheusRegistry.Register(LendMeRequestsAmount)
	//PrometheusRegistry.Register(TransactionsTotalAmount)
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
						DailyImportSubsStats.Reset()
						LendMeRequestsCount.Reset()
						LendMeRequestsAmount.Reset()
						//TransactionsTotalAmount.Reset()
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
