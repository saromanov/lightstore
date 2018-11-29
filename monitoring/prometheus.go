package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	rpcDurations = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "rpc_durations_seconds",
			Help:       "RPC latency distributions.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"lightstore"},
	)
	rpcDurationsHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "rpc_durations_histogram_seconds",
		Help:    "RPC latency distributions.",
		Buckets: prometheus.LinearBuckets(*normMean-5**normDomain, .5**normDomain, 20),
	})
	opsInsert = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "lightstore_insert_total",
		Help: "Number of inserts",
	})
	opsGet = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "lightstore_get_total",
		Help: "Number of read operations",
	})
)

// Init provides initialization of the Prometheus
func Init() {
	prometheus.MustRegister(rpcDurations)
	prometheus.MustRegister(rpcDurationsHistogram)
	prometheus.MustRegister(opsInsert)
	prometheus.MustRegister(opsGet)
}
