package main

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// --- /api/search ---

func TestSearchMissingQuery(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/search", nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, 422, w.Code)
	resp := decode[RequestValidationError](t, w.Body.Bytes())
	assert.Equal(t, 422, resp.StatusCode)
	assert.Contains(t, *resp.Message, "required")
}

func TestSearchDBError(t *testing.T) {
	mockSearchPagesQuery = func(_ *sql.DB, q, l string, limit int) ([]SearchResult, error) {
		return nil, errors.New("boom")
	}
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/search?q=hello", nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, 422, w.Code)
	resp := decode[RequestValidationError](t, w.Body.Bytes())
	assert.Contains(t, *resp.Message, "Search failed")
}

func TestSearchSuccess(t *testing.T) {
	mockSearchPagesQuery = func(_ *sql.DB, q, l string, limit int) ([]SearchResult, error) {
		now := time.Now()
		return []SearchResult{
			{
				Title:       "hi",
				URL:         "https://example.com",
				Language:    "en",
				LastUpdated: &now,
				Snippet:     "hello world",
			},
		}, nil
	}
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/search?q=hi", nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	resp := decode[SearchResponse](t, w.Body.Bytes())
	assert.Len(t, resp.Data, 1)
}
