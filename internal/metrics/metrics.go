package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Технические метрики
	HTTPRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)

	// Бизнесовые метрики
	PVZCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "pvz_created_total",
			Help: "Total number of created PVZs",
		},
	)

	ReceptionsCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "receptions_created_total",
			Help: "Total number of created receptions",
		},
	)

	ProductsAddedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "products_added_total",
			Help: "Total number of added products",
		},
	)
)

func Init() {
	prometheus.MustRegister(
		HTTPRequestTotal,
		HTTPRequestDuration,
		PVZCreatedTotal,
		ReceptionsCreatedTotal,
		ProductsAddedTotal,
	)
}
