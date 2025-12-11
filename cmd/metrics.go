package main

import (
	"regexp"
	"strings"
	"unicode/utf8"

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
			Help: "Count of HTTP requests grouped by browser family and version",
		},
		[]string{"browser", "version"},
	)

	searchQueryCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_search_queries_total",
			Help: "Total number of search queries (sanitized to avoid high-cardinality inputs)",
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
		browser, version := parseUserAgentDetails(userAgent)

		browserCounter.WithLabelValues(browser, version).Inc()

		c.Next()
	}
}

func parseUserAgentDetails(ua string) (string, string) {
	ua = strings.ToLower(ua)

	browser := "Other"
	version := "unknown"

	switch {
	case strings.Contains(ua, "samsungbrowser"):
		browser = "Samsung Internet"
		version = extractVersion(ua, "samsungbrowser")
	case strings.Contains(ua, "vivaldi"):
		browser = "Vivaldi"
		version = extractVersion(ua, "vivaldi")
	case strings.Contains(ua, "duckduckgo"):
		browser = "DuckDuckGo"
		version = extractVersion(ua, "duckduckgo")
	case strings.Contains(ua, "opr") || strings.Contains(ua, "opera"):
		browser = "Opera"
		version = extractVersion(ua, "opr")
	case strings.Contains(ua, "yabrowser"):
		browser = "Yandex"
		version = extractVersion(ua, "yabrowser")
	case strings.Contains(ua, "brave"):
		browser = "Brave"
		version = extractVersion(ua, "brave")
	case strings.Contains(ua, "edg/") || strings.Contains(ua, "edge"):
		browser = "Edge"
		version = extractVersion(ua, "edg")
	case strings.Contains(ua, "chrome"):
		browser = "Chrome"
		version = extractVersion(ua, "chrome")
	case strings.Contains(ua, "safari"):
		browser = "Safari"
		version = extractVersion(ua, "version")
	case strings.Contains(ua, "firefox"):
		browser = "Firefox"
		version = extractVersion(ua, "firefox")
	case strings.Contains(ua, "msie") || strings.Contains(ua, "trident"):
		browser = "Internet Explorer"
		version = extractVersion(ua, "msie")
	case strings.Contains(ua, "bot") || strings.Contains(ua, "crawler") || strings.Contains(ua, "spider") || strings.Contains(ua, "slurp"):
		browser = "Bot"
	case strings.Contains(ua, "curl") || strings.Contains(ua, "wget") || strings.Contains(ua, "postman"):
		browser = "Dev Tools"
	}

	return browser, version
}

func extractVersion(ua, token string) string {
	pattern := regexp.MustCompile(token + "/([0-9]+)")
	matches := pattern.FindStringSubmatch(ua)
	if len(matches) == 2 {
		return matches[1]
	}

	pattern = regexp.MustCompile(token + " ([0-9]+)")
	matches = pattern.FindStringSubmatch(ua)
	if len(matches) == 2 {
		return matches[1]
	}

	return "unknown"
}

func sanitizeLabelValue(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.ReplaceAll(value, "\r", " ")
	if value == "" {
		return "unknown"
	}

	const maxRunes = 64
	if utf8.RuneCountInString(value) > maxRunes {
		runes := []rune(value)
		value = string(runes[:maxRunes])
	}

	return value
}
