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

		labels := prometheusLabels{method: c.Request.Method, path: path, status: fmt.Sprintf("%d", status)}
		observeRequest(duration, labels)
		log.Printf("[REQ] %s %s -> %d (%v)", c.Request.Method, path, status, duration)
	}
}
