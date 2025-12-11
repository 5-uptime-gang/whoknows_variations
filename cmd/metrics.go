package main

import (
	"strings"

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
	browserCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_browser_usage_total",
			Help: "Count of HTTP requests grouped by browser family",
		},
		[]string{"browser"},
	)
	searchQueryCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_search_queries_total",
			Help: "Total count of successful search queries by term.",
		},
		[]string{"query"},
	)
)

func init() {
	prometheus.MustRegister(requestCounter, requestDuration, userSignupCounter, browserCounter, searchQueryCounter)
}

func metricsHandler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

func BrowserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAgent := c.GetHeader("User-Agent")
		browser := parseUserAgent(userAgent)

		browserCounter.WithLabelValues(browser).Inc()

		c.Next()
	}
}

func parseUserAgent(ua string) string {
	ua = strings.ToLower(ua)

	// 1. Specifikke browsere (Tjekkes først, da de ofte indeholder "Chrome" eller "Safari")
	if strings.Contains(ua, "samsungbrowser") {
		return "Samsung Internet"
	}
	if strings.Contains(ua, "vivaldi") {
		return "Vivaldi"
	}
	if strings.Contains(ua, "duckduckgo") {
		return "DuckDuckGo"
	}
	if strings.Contains(ua, "opr") || strings.Contains(ua, "opera") {
		return "Opera"
	}
	if strings.Contains(ua, "yabrowser") {
		return "Yandex"
	}
	if strings.Contains(ua, "brave") {
		return "Brave"
	}

	if strings.Contains(ua, "edg/") || strings.Contains(ua, "edge") {
		return "Edge"
	}

	if strings.Contains(ua, "chrome") {
		return "Chrome"
	}
	if strings.Contains(ua, "safari") {
		return "Safari"
	}
	if strings.Contains(ua, "firefox") {
		return "Firefox"
	}

	if strings.Contains(ua, "msie") || strings.Contains(ua, "trident") {
		return "Internet Explorer"
	}
	if strings.Contains(ua, "bot") || strings.Contains(ua, "crawler") || strings.Contains(ua, "spider") || strings.Contains(ua, "slurp") {
		return "Bot"
	}
	if strings.Contains(ua, "curl") || strings.Contains(ua, "wget") || strings.Contains(ua, "postman") {
		return "Dev Tools"
	}

	return "Other"
}
