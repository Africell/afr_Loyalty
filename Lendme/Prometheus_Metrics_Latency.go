package Lendme

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	PrometheusLatencyRegistry = prometheus.NewRegistry()
	EndToEndLatency           = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "EndToEndLatency",
			Help: "How long it took to process the request in milliseconds ",
			//Buckets: prometheus.ExponentialBucketsRange(50, 10000, 15),
			//Buckets: prometheus.ExponentialBuckets(100, 2, 10),
			Buckets: []float64{50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 1000, 1250, 1500, 1750, 2000, 3000, 4000, 5000, 7500, 10000, 15000, 20000},
		},
		[]string{"host", "path", "status"},
	)
	//	BRMLatency.With(prometheus.Labels{"host": Configuration.HostId, "path": route}).Observe(float64((time.Since(StartDate).Nanoseconds()) / 1000000))

)

func Init_Prometheus_Metrics_Latency() {
	log.Println("Init Prometheus metrics latency")
	PrometheusLatencyRegistry.Register(EndToEndLatency)
}

func Reset_Prometheus_Metrics_Latency() {
	exec := 0
	for range time.Tick(time.Second * 1) {
		_CurrentDateTime := time.Now()
		_hr, _mi, _se := _CurrentDateTime.Clock()
		if _hr == 0 {
			if _mi == 0 {
				if _se < 10 {
					if exec == 0 {
						EndToEndLatency.Reset()
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

func CustomPrometheusLatencyHandler() http.Handler {
	return promhttp.HandlerFor(
		PrometheusLatencyRegistry,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: false,
		},
	)
}
