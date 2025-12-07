package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},	
		[]string{"method", "path", "status"},
	)
	userSignupCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "app_user_signups_total",
            Help: "Total number of user signup attempts (success/failed).",
        },
        []string{"status"},
    )
)

func init() {
	prometheus.MustRegister(requestCounter, requestDuration, userSignupCounter)
}

func metricsHandler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}
