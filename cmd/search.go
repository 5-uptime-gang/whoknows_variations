package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func apiSearch(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		msg := "Query parameter 'q' is required"
		log.Printf("[SEARCH] Invalid request: %v", msg)
		c.JSON(http.StatusUnprocessableEntity, RequestValidationError{StatusCode: 422, Message: &msg})
		return
	}
	lang := c.DefaultQuery("language", "en")
	results, err := SearchPagesQuery(db, q, lang)
	if err != nil {
		msg := "Search failed: " + err.Error()
		log.Printf("[SEARCH] Search failed: %v", msg)
		c.JSON(http.StatusUnprocessableEntity, RequestValidationError{StatusCode: 422, Message: &msg})
		return
	}

	normalizedQuery := strings.ToLower(strings.TrimSpace(q))

	searchQueryCounter.WithLabelValues(normalizedQuery).Inc()

	

	safeQ := strings.ReplaceAll(strings.ReplaceAll(q, "\n", "_"), "\r", "_")
	safeLang := strings.ReplaceAll(strings.ReplaceAll(lang, "\n", "_"), "\r", "_")

	log.Printf("[SEARCH] Search successful: q=%q, lang=%q", safeQ, safeLang)
	c.JSON(http.StatusOK, SearchResponse{Data: results})
}
