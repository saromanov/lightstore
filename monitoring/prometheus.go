package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	normMean   = 0.0001
	normDomain = 0.002
)

var (
	rpcErrorsHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "errors_histogram_total",
		Help:    "Distribution of errors after handling of data",
		Buckets: prometheus.LinearBuckets(normMean-5*normDomain, .5*normDomain, 20),
	})
	opsInsert = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "lightstore_insert_total",
		Help: "Number of inserts",
	})
	opsGet = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "lightstore_get_total",
		Help: "Number of read operations",
	})
	gaugeInsert = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name:      "ligthstore_insert_gauge",
			Help:      "Number of inserts",
			Subsystem: "storage",
		},
	)
)

// Init provides initialization of the Prometheus
func Init() {
	prometheus.MustRegister(rpcErrorsHistogram)
	prometheus.MustRegister(opsInsert)
	prometheus.MustRegister(opsGet)
}
