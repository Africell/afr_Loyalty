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

	HTTPPostToKafka = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "HTTPPostToKafka",
			Help: "HTTPPostToKafka",
		},
		[]string{"Description", "Status", "Error"},
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
	LiveFeedCounters = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "LiveFeedCounters",
			Help: "LiveFeedCounters",
		},
		[]string{"Stream", "Type", "Description"},
	)
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
	PointsCreditedCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "PointsCreditedCount",
			Help: "PointsCreditedCount",
		},
		[]string{"EventSource", "EventType", "Level"},
	)
	PointsCredited = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "PointsCredited",
			Help: "PointsCredited",
		},
		[]string{"EventSource", "EventType", "Level"},
	)
	PointsDebitedCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "PointsDebitedCount",
			Help: "PointsDebitedCount",
		},
		[]string{"EventSource", "DebitType", "Level"},
	)
	PointsDebited = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "PointsDebited",
			Help: "PointsDebited",
		},
		[]string{"EventSource", "DebitType", "Level"},
	)

	BundleRedemptionCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "BundleRedemptionCount",
			Help: "BundleRedemptionCount",
		},
		[]string{"EventSource", "BunldeId", "Level"},
	)
	BundleRedemptionPoints = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "BundleRedemptionPoints",
			Help: "BundleRedemptionPoints",
		},
		[]string{"EventSource", "BunldeId", "Level"},
	)
	AirtimeRedemptionCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "AirtimeRedemptionCount",
			Help: "AirtimeRedemptionCount",
		},
		[]string{"EventSource", "Level"},
	)
	AirtimeRedemptionPoints = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "AirtimeRedemptionPoints",
			Help: "AirtimeRedemptionPoints",
		},
		[]string{"EventSource", "Level"},
	)
	AirtimeRedemptionAmount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "AirtimeRedemptionAmount",
			Help: "AirtimeRedemptionAmount",
		},
		[]string{"EventSource", "Level"},
	)
	MobileMoneyRedemptionCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "MobileMoneyRedemptionCount",
			Help: "MobileMoneyRedemptionCount",
		},
		[]string{"EventSource", "Level"},
	)
	MobileMoneyRedemptionPoints = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "MobileMoneyRedemptionPoints",
			Help: "MobileMoneyRedemptionPoints",
		},
		[]string{"EventSource", "Level"},
	)
	MobileMoneyRedemptionAmount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "MobileMoneyRedemptionAmount",
			Help: "MobileMoneyRedemptionAmount",
		},
		[]string{"EventSource", "Level"},
	)

	LoyaltySubsSummary = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "LoyaltySubsSummary",
			Help: "LoyaltySubsSummary",
		},
		[]string{"Level", "Description"},
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
	PrometheusRegistry.Register(HTTPPostToKafka)
	//*************************
	//Loaylty Metrics
	//*************************
	LoyaltyPrometheusRegistry.Register(PortStatus)
	LoyaltyPrometheusRegistry.Register(LiveFeedCounters)
	LoyaltyPrometheusRegistry.Register(LoyaltyGovernancePools)
	LoyaltyPrometheusRegistry.Register(NewJoiningsCount)
	LoyaltyPrometheusRegistry.Register(PointsCreditedCount)
	LoyaltyPrometheusRegistry.Register(PointsCredited)
	LoyaltyPrometheusRegistry.Register(PointsDebitedCount)
	LoyaltyPrometheusRegistry.Register(PointsDebited)
	LoyaltyPrometheusRegistry.Register(BundleRedemptionCount)
	LoyaltyPrometheusRegistry.Register(BundleRedemptionPoints)
	LoyaltyPrometheusRegistry.Register(AirtimeRedemptionCount)
	LoyaltyPrometheusRegistry.Register(AirtimeRedemptionPoints)
	LoyaltyPrometheusRegistry.Register(AirtimeRedemptionAmount)
	LoyaltyPrometheusRegistry.Register(MobileMoneyRedemptionCount)
	LoyaltyPrometheusRegistry.Register(MobileMoneyRedemptionPoints)
	LoyaltyPrometheusRegistry.Register(MobileMoneyRedemptionAmount)
	LoyaltyPrometheusRegistry.Register(LoyaltySubsSummary)

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
						HTTPPostToKafka.Reset()
						//LendMeOutstandingSummary.Reset()
						//SubsDumpFile.Reset()
						//TransactionsTotalAmount.Reset()
						//*************************
						//Loaylty Metrics
						//*************************
						LiveFeedCounters.Reset()
						NewJoiningsCount.Reset()
						PointsCreditedCount.Reset()
						PointsCredited.Reset()
						PointsDebitedCount.Reset()
						PointsDebited.Reset()
						BundleRedemptionCount.Reset()
						BundleRedemptionPoints.Reset()
						AirtimeRedemptionCount.Reset()
						AirtimeRedemptionPoints.Reset()
						AirtimeRedemptionAmount.Reset()
						MobileMoneyRedemptionCount.Reset()
						MobileMoneyRedemptionPoints.Reset()
						MobileMoneyRedemptionAmount.Reset()
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
