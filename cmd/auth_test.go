package main

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// --- /api/register ---

func TestRegisterValidationErrors(t *testing.T) {
	router := setupRouter()
	body := `{"email":"x@y.z","password":"abc","password2":"abc"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestRegisterSuccess(t *testing.T) {
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

func TestRegisterUsernameTaken(t *testing.T) {
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

func TestLoginInvalidBody(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBufferString("invalid-json"))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	assert.Equal(t, 422, w.Code)
}

func TestLoginInvalidUsername(t *testing.T) {
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

func TestLoginWrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost) //NOSONAR
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

func TestApiLoginSuccess(t *testing.T) {
	t.Skip("DB tests temporarily disabled during PostgreSQL migration")
	// ensure global db is non-nil to prevent nil deref
	db, _ = sql.Open("sqlite", ":memory:")

	// create bcrypt hash for the expected password
	hash, _ := bcrypt.GenerateFromPassword([]byte("goodpw"), bcrypt.DefaultCost) //NOSONAR

	// mock database query
	mockGetUserByUsernameQuery = func(_ *sql.DB, u string) (int, string, string, string, error) {
		return 1, "u", "e@example.com", string(hash), nil
	}
	GetUserByUsernameQuery = func(db *sql.DB, u string) (int, string, string, string, error) {
		return mockGetUserByUsernameQuery(db, u)
	}

	router := setupRouter()
	body := `{"username":"u","password":"goodpw"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "should return 200 OK")

	resp := decode[AuthResponse](t, w.Body.Bytes())
	if assert.NotNil(t, resp.Message, "response message should not be nil") {
		assert.Equal(t, "login successful", *resp.Message)
	}
	if assert.NotNil(t, resp.StatusCode, "status code should not be nil") {
		assert.Equal(t, 200, *resp.StatusCode)
	}
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
