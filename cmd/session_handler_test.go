package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newRouterWithSessionOnly() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	api := r.Group("/api")
	api.GET("/session", apiSession)
	return r
}

func TestSession_NotLoggedIn(t *testing.T) {
	r := newRouterWithSessionOnly()
	req := httptest.NewRequest(http.MethodGet, "/api/session", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	if w.Body.String() != `{"logged_in":false}` {
		t.Fatalf("bad body: %s", w.Body.String())
	}
}

func TestSession_LoggedIn(t *testing.T) {
	r := newRouterWithSessionOnly()
	req := httptest.NewRequest(http.MethodGet, "/api/session", nil)
	req.AddCookie(&http.Cookie{Name: "user_id", Value: "123"})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	if w.Body.String() != `{"logged_in":true}` {
		t.Fatalf("bad body: %s", w.Body.String())
	}
}
