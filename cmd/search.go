package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func apiSearch(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	if q == "" {
		msg := "Query parameter 'q' is required"
		log.Printf("[SEARCH] Invalid request: %v", msg)
		c.JSON(http.StatusUnprocessableEntity, RequestValidationError{StatusCode: 422, Message: &msg})
		return
	}

	lang := resolveLanguage(q, c.Query("language"))
	limit := parseLimit(c.DefaultQuery("limit", "10"))

	results, err := SearchPagesQuery(db, q, lang, limit)
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
	safeLimit := strings.ReplaceAll(strings.ReplaceAll(strconv.Itoa(limit), "\n", "_"), "\r", "_")

	log.Printf("[SEARCH] Search successful: q=%q, lang=%q, limit=%s", safeQ, safeLang, safeLimit)
	c.JSON(http.StatusOK, SearchResponse{Data: results})
}

func resolveLanguage(query, langParam string) string {
	normalized := strings.ToLower(strings.TrimSpace(langParam))
	switch normalized {
	case "da", "danish":
		return "da"
	case "en", "english":
		return "en"
	}

	lower := strings.ToLower(query)
	if strings.ContainsAny(lower, "\u00e6\u00f8\u00e5") {
		return "da"
	}

	danishHints := map[string]struct{}{
		"og": {}, "ikke": {}, "det": {}, "der": {}, "som": {},
		"jeg": {}, "du": {}, "vi": {}, "jer": {}, "for": {},
		"med": {}, "uden": {}, "hvor": {}, "hvordan": {},
	}
	for _, token := range strings.Fields(lower) {
		if _, ok := danishHints[token]; ok {
			return "da"
		}
	}
	return "en"
}

func parseLimit(raw string) int {
	if raw == "" {
		return 10
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return 10
	}
	if n < 1 {
		return 1
	}
	if n > 50 {
		return 50
	}
	return n
}
