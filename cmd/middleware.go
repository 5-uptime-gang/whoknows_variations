package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Printf("[REQ] %s %s -> %d (%v)", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(start))
	}
}
