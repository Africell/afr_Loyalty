package Lendme

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	PrometheusRegistry        = prometheus.NewRegistry()
	LoyaltyPrometheusRegistry = prometheus.NewRegistry()
	//PortStatus.With(prometheus.Labels{"DestinationHost": "HHHH"}).Set(1)
	PortStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "PortStatus",
			Help: "1: port up, 2: port down",
		},
		[]string{"DestinationHost"},
	)
	//*************************
	//Lendme Metrics
	//*************************
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

	LendMePayBackCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "LendMePayBackCount",
			Help: "LendMePayBackCount",
		},
		[]string{"Status", "Reason", "Description"},
	)
	LendMePayBackAmount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "LendMePayBackAmount",
			Help: "LendMePayBackAmount",
		},
		[]string{"Status", "Reason", "Description"},
	)

	SubsDumpFile = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "SubsDumpFile",
			Help: "SubsDumpFile",
		},
		[]string{"FileName"},
	)

	SubsDumpFileImportTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "SubsDumpFileImportTime",
			Help: "SubsDumpFileImportTime",
		},
		[]string{"TimeInMinutes"},
	)

	LendMeOutstandingSummary = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "LendMeOutstandingSummary",
			Help: "LendMeOutstandingSummary",
		},
		[]string{"Field"},
	)
	// TransactionsTotalAmount = prometheus.NewGaugeVec(
	// 	prometheus.GaugeOpts{
	// 		Name: "TransactionsTotalAmount",
	// 		Help: "TransactionsTotalAmount",
	// 	},
	// 	[]string{"Stream", "Type", "Description"},
	// )

	//*************************
	//Loaylty Metrics
	//*************************
	LoyaltyGovernancePools = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "LoyaltyGovernancePools",
			Help: "LoyaltyGovernancePools",
		},
		[]string{"Pool"},
	)
	NewJoiningsCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "NewJoiningsCount",
			Help: "NewJoiningsCount",
		},
		[]string{"Source"},
	)
	AwardedTransactions = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "AwardedTransactions",
			Help: "AwardedTransactions",
		},
		[]string{"EventSource", "EventType", "EventDetail"},
	)
	AwardedPoints = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "AwardedPoints",
			Help: "AwardedPoints",
		},
		[]string{"EventSource", "EventType", "EventDetail"},
	)
)

//DailyImportSubsStats.With(prometheus.Labels{"IsElligble":"", "Reason":"", "Scheme":""}).Inc()
//LendMeOutstandingSummary.With(prometheus.Labels{""Field"": "desc"}).Add()

func Init_Prometheus_Metrics() {
	log.Println("Init Prometheus metrics")
	PrometheusRegistry.Register(PortStatus)
	//*************************
	//Lendme Metrics
	//*************************
	PrometheusRegistry.Register(DailyImportSubsStats)
	PrometheusRegistry.Register(LendMeRequestsCount)
	PrometheusRegistry.Register(LendMeRequestsAmount)
	PrometheusRegistry.Register(LendMePayBackCount)
	PrometheusRegistry.Register(LendMePayBackAmount)
	PrometheusRegistry.Register(SubsDumpFile)
	PrometheusRegistry.Register(SubsDumpFileImportTime)
	PrometheusRegistry.Register(LendMeOutstandingSummary)
	//*************************
	//Loaylty Metrics
	//*************************
	LoyaltyPrometheusRegistry.Register(PortStatus)
	LoyaltyPrometheusRegistry.Register(LoyaltyGovernancePools)
	LoyaltyPrometheusRegistry.Register(NewJoiningsCount)
	LoyaltyPrometheusRegistry.Register(AwardedPoints)
	LoyaltyPrometheusRegistry.Register(AwardedTransactions)
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
						//*************************
						//Lendme Metrics
						//*************************
						DailyImportSubsStats.Reset()
						LendMeRequestsCount.Reset()
						LendMeRequestsAmount.Reset()
						LendMePayBackCount.Reset()
						LendMePayBackAmount.Reset()
						//LendMeOutstandingSummary.Reset()
						//SubsDumpFile.Reset()
						//TransactionsTotalAmount.Reset()
						//*************************
						//Loaylty Metrics
						//*************************
						NewJoiningsCount.Reset()
						AwardedPoints.Reset()
						AwardedTransactions.Reset()
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

func LoyaltyPrometheusHandler() http.Handler {
	return promhttp.HandlerFor(
		LoyaltyPrometheusRegistry,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: false,
		},
	)
}
