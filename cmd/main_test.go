package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// --- Mock dependencies ---

var (
	mockInsertUserQuery        func(*sql.DB, string, string, string) (int64, error)
	mockGetUserByUsernameQuery func(*sql.DB, string) (int, string, string, string, error)
	mockSearchPagesQuery       func(*sql.DB, string, string) ([]Page, error)
)

// Patch the global functions to mocks for testing
func init() {
	InsertUserQuery = func(db *sql.DB, u, e, p string) (int64, error) {
		return mockInsertUserQuery(db, u, e, p)
	}
	GetUserByUsernameQuery = func(db *sql.DB, u string) (int, string, string, string, error) {
		return mockGetUserByUsernameQuery(db, u)
	}
	SearchPagesQuery = func(db *sql.DB, q, lang string) ([]Page, error) {
		return mockSearchPagesQuery(db, q, lang)
	}
}

// --- Helpers ---

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api")
	{
		api.GET("/weather", apiWeather)
		api.GET("/search", apiSearch)
		api.POST("/login", apiLogin)
		api.POST("/register", apiRegister)
		api.GET("/logout", apiLogout)
		api.GET("/session", apiSession)
	}
	return router
}

func decode[T any](t *testing.T, body []byte) T {
	var v T
	err := json.Unmarshal(body, &v)
	assert.NoError(t, err)
	return v
}

// --- /api/register ---

func TestRegister_ValidationErrors(t *testing.T) {
	router := setupRouter()
	body := `{"email":"x@y.z","password":"abc","password2":"abc"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestRegister_Success(t *testing.T) {
	mockInsertUserQuery = func(_ *sql.DB, u, e, p string) (int64, error) { return 1, nil }

	router := setupRouter()
	body := `{"username":"testuser","email":"user@example.com","password":"abc","password2":"abc"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	resp := decode[AuthResponse](t, w.Body.Bytes())
	assert.Equal(t, 200, *resp.StatusCode)
	assert.Equal(t, "user registered successfully", *resp.Message)
}

func TestRegister_UsernameTaken(t *testing.T) {
	mockInsertUserQuery = func(_ *sql.DB, u, e, p string) (int64, error) { return 0, errors.New("duplicate") }

	router := setupRouter()
	body := `{"username":"exists","email":"user@example.com","password":"abc","password2":"abc"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	assert.Equal(t, 422, w.Code)

	resp := decode[HTTPValidationError](t, w.Body.Bytes())
	assert.Equal(t, "database", resp.Detail[0].Loc[0])
}

// --- /api/login ---

func TestLogin_InvalidBody(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBufferString("invalid-json"))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	assert.Equal(t, 422, w.Code)
}

func TestLogin_InvalidUsername(t *testing.T) {
	mockGetUserByUsernameQuery = func(_ *sql.DB, u string) (int, string, string, string, error) {
		return 0, "", "", "", errors.New("not found")
	}
	router := setupRouter()
	body := `{"username":"nope","password":"x"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	assert.Equal(t, 422, w.Code)
	resp := decode[HTTPValidationError](t, w.Body.Bytes())
	assert.Equal(t, "username", resp.Detail[0].Loc[0])
}

func TestLogin_WrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("goodpw"), bcrypt.DefaultCost)
	mockGetUserByUsernameQuery = func(_ *sql.DB, u string) (int, string, string, string, error) {
		return 1, "u", "e", string(hash), nil
	}
	router := setupRouter()
	body := `{"username":"u","password":"badpw"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	assert.Equal(t, 422, w.Code)
	resp := decode[HTTPValidationError](t, w.Body.Bytes())
	assert.Equal(t, "password", resp.Detail[0].Loc[0])
}

func TestLogin_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("goodpw"), bcrypt.DefaultCost)
	mockGetUserByUsernameQuery = func(_ *sql.DB, u string) (int, string, string, string, error) {
		return 1, "u", "e", string(hash), nil
	}
	router := setupRouter()
	body := `{"username":"u","password":"goodpw"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	resp := decode[AuthResponse](t, w.Body.Bytes())
	assert.Equal(t, "login successful", *resp.Message)
}

// --- /api/logout ---

func TestLogout(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/logout", nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	resp := decode[AuthResponse](t, w.Body.Bytes())
	assert.Equal(t, "logged out", *resp.Message)
}

// --- /api/search ---

func TestSearch_MissingQuery(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/search", nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, 422, w.Code)
	resp := decode[RequestValidationError](t, w.Body.Bytes())
	assert.Equal(t, 422, resp.StatusCode)
	assert.Contains(t, *resp.Message, "required")
}

func TestSearch_DBError(t *testing.T) {
	mockSearchPagesQuery = func(_ *sql.DB, q, l string) ([]Page, error) {
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

func TestSearch_Success(t *testing.T) {
	mockSearchPagesQuery = func(_ *sql.DB, q, l string) ([]Page, error) {
		return []Page{
			{
				Title:       "hi",
				URL:         "https://example.com",
				Language:    "en",
				LastUpdated: time.Now(),
				Content:     "hello world",
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
