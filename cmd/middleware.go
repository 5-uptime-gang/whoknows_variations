package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		status := c.Writer.Status()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		labels := []string{c.Request.Method, path, fmt.Sprintf("%d", status)}
		requestCounter.WithLabelValues(labels...).Inc()
		requestDuration.WithLabelValues(labels...).Observe(duration.Seconds())

		log.Printf("[REQ] %s %s -> %d (%v)", c.Request.Method, path, status, duration)
	}
}
