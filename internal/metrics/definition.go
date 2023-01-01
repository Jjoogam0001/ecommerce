package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestDuration = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  namespace,
		Name:       "request_duration",
		Help:       "The HTTP request latencies in milliseconds. Partitioned by status code, HTTP method, endpoint and originator.",
		Objectives: map[float64]float64{0.99: 0.001},
		MaxAge:     time.Minute,
	}, []string{"host", "code", "method", "endpoint"})

	RequestCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "request_counter",
		Help:      "How many requests receive app.",
	}, []string{"host", "code", "method", "endpoint", "originator"})

	DbCall = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  namespace,
		Name:       "db_request_duration",
		Help:       "The database request latencies in milliseconds.",
		Objectives: map[float64]float64{0.99: 0.001},
		MaxAge:     time.Minute,
	}, []string{"host", "function"})

	DbCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "db_request_counter",
		Help:      "How many database requests.",
	}, []string{"host", "function"})

	DbError = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "db_errors_counter",
		Help:      "How many unexpected database errors.",
	}, []string{"host"})
)
