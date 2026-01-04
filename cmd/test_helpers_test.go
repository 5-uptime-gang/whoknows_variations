package main

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// --- Mock dependencies ---

var (
	mockInsertUserQuery        func(*sql.DB, string, string, string) (int64, error)
	mockGetUserByUsernameQuery func(*sql.DB, string) (int, string, string, string, error)
	mockSearchPagesQuery       func(*sql.DB, string, string, int) ([]SearchResult, error)
)

// Patch the global functions to mocks for testing
func init() {
	InsertUserQuery = func(db *sql.DB, u, e, p string) (int64, error) {
		return mockInsertUserQuery(db, u, e, p)
	}
	GetUserByUsernameQuery = func(db *sql.DB, u string) (int, string, string, string, error) {
		return mockGetUserByUsernameQuery(db, u)
	}
	SearchPagesQuery = func(db *sql.DB, q, lang string, limit int) ([]SearchResult, error) {
		return mockSearchPagesQuery(db, q, lang, limit)
	}
}

// --- Helpers ---

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return newRouter()
}

func decode[T any](t *testing.T, body []byte) T {
	var v T
	err := json.Unmarshal(body, &v)
	assert.NoError(t, err)
	return v
}
