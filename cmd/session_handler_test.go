package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSessionNotLoggedIn(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/session", nil)

	router.ServeHTTP(w, req)
	resp := decode[AuthResponse](t, w.Body.Bytes())
	assert.Equal(t, 401, *resp.StatusCode)
	assert.Equal(t, "not logged in", *resp.Message)
}

func TestSessionLoggedIn(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/session", nil)
	req.AddCookie(&http.Cookie{Name: "user_id", Value: "1"})

	router.ServeHTTP(w, req)
	resp := decode[AuthResponse](t, w.Body.Bytes())
	assert.Equal(t, 200, *resp.StatusCode)
	assert.Equal(t, "logged in", *resp.Message)
}
