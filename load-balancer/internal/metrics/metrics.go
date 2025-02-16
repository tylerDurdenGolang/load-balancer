package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    RequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "loadbalancer_requests_total",
            Help: "Total number of processed requests",
        },
        []string{"backend", "status"},
    )

    RequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "loadbalancer_request_duration_seconds",
            Help:    "Request processing duration",
            Buckets: prometheus.DefBuckets,
        },
        []string{"backend"},
    )

    BackendHealth = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "loadbalancer_backend_health",
            Help: "Backend health status (1 - healthy, 0 - unhealthy)",
        },
        []string{"backend"},
    )

    AlgorithmWeights = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "loadbalancer_algorithm_weights",
            Help: "Current backend weights from balancing algorithms",
        },
        []string{"backend"},
    )
)

func RegisterCustomMetrics() {
    prometheus.MustRegister(
        RequestsTotal,
        RequestDuration,
        BackendHealth,
        AlgorithmWeights,
    )
}