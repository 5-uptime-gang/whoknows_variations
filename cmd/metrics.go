package main

import (
	"strings"
	"regexp"
	"time"
	"log"
	"database/sql"

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
		[]string{"browser", "version"},
	)
	searchQueryCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_search_queries_total",
			Help: "Total count of successful search queries by term.",
		},
		[]string{"query"},
	)
	versionRegex = regexp.MustCompile(`([0-9.]+[\-0-9.]*)`)

	userTotalGauge = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "app_users_current_total",
            Help: "Det nuværende antal brugere i databasen.",
        },
    )
)

func init() {
	prometheus.MustRegister(requestCounter, requestDuration, userSignupCounter, browserCounter, searchQueryCounter, userTotalGauge)
}

func metricsHandler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

func BrowserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		userAgent := c.GetHeader("User-Agent")
		browser, version := parseUserAgent(userAgent)

		browserCounter.WithLabelValues(browser, version).Inc()

		c.Next()
	}
}

func parseUserAgent(ua string) (browser string, version string) {
	ua = strings.ToLower(ua)
	version = "N/A"

	getSpecificVersion := func(key string) string {
		startIndex := strings.Index(ua, key)
		if startIndex == -1 {
			return "N/A"
		}
		
		startIndex += len(key)
		
		endIndex := startIndex
		for endIndex < len(ua) && (ua[endIndex] >= '0' && ua[endIndex] <= '9' || ua[endIndex] == '.') {
			endIndex++
		}
		
		v := strings.TrimSpace(ua[startIndex:endIndex])
		
		if match := versionRegex.FindString(v); match != "" {
			return match
		}
		return v
	}
	
	if strings.Contains(ua, "bot") || strings.Contains(ua, "crawler") || strings.Contains(ua, "spider") || strings.Contains(ua, "slurp") {
		return "Bot", "N/A"
	}
	if strings.Contains(ua, "curl") || strings.Contains(ua, "wget") || strings.Contains(ua, "postman") {
		return "Dev Tools", "N/A"
	}

	if strings.Contains(ua, "edg/") {
		return "Edge", getSpecificVersion("edg/")
	}
	if strings.Contains(ua, "opr/") || strings.Contains(ua, "opera") {
		return "Opera", getSpecificVersion("opr/")
	}
	if strings.Contains(ua, "samsungbrowser") {
		return "Samsung Internet", getSpecificVersion("samsungbrowser/")
	}
	if strings.Contains(ua, "vivaldi") {
		return "Vivaldi", getSpecificVersion("vivaldi/")
	}
	
	if strings.Contains(ua, "chrome") {
		return "Chrome", getSpecificVersion("chrome/")
	}
	
	if strings.Contains(ua, "firefox") {
		return "Firefox", getSpecificVersion("firefox/")
	}
	
	if strings.Contains(ua, "safari") && !strings.Contains(ua, "android") {
		if v := getSpecificVersion("version/"); v != "N/A" {
			return "Safari (Desktop)", v
		}
		return "Safari (Desktop)", getSpecificVersion("safari/")
	}
	
	if strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") {
		return "Safari on iOS", getSpecificVersion("version/")
	}

	if strings.Contains(ua, "android") {
		return "Android Browser", getSpecificVersion("android/")
	}
	if strings.Contains(ua, "msie") || strings.Contains(ua, "trident") {
		return "Internet Explorer", getSpecificVersion("msie ")
	}

	return "Other", "N/A"
}

func monitorUserCount(db *sql.DB) {
    ticker := time.NewTicker(15 * time.Second)
    defer ticker.Stop()

    for range ticker.C {

        count, err := GetUserCountQuery(db) 
        if err != nil {
            log.Printf("Fejl ved tælling af brugere: %v", err)
            continue
        }
        userTotalGauge.Set(count)
    }
}