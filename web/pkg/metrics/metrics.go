package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	ServiceProcessTotalCount  *prometheus.CounterVec
	ServiceProcessFailedCount *prometheus.CounterVec
	ServiceProcessSuccedCount *prometheus.CounterVec
	ServiceProcessDuration    *prometheus.SummaryVec
)

func Init(ns string) {
	ServiceProcessTotalCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: ns,
		Name:      "servicd_process_process_count",
		Help:      "The total counter of handler process",
	}, []string{"group", "method"})
	ServiceProcessFailedCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: ns,
		Name:      "servicd_process_failed_count",
		Help:      "The failed counter of handler process",
	}, []string{"group", "method"})
	ServiceProcessSuccedCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: ns,
		Name:      "servicd_process_succed_count",
		Help:      "The succed counter of handler process",
	}, []string{"group", "method"})
	ServiceProcessDuration = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  ns,
		Name:       "servicd_process_process_durations",
		Help:       "The time duration of handler process",
		MaxAge:     1 * time.Minute,
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"group", "method"})
	prometheus.MustRegister(ServiceProcessTotalCount)
	prometheus.MustRegister(ServiceProcessFailedCount)
	prometheus.MustRegister(ServiceProcessSuccedCount)
	prometheus.MustRegister(ServiceProcessDuration)
}
