package main

import (
	"WHOKNOWS_VARIATIONS/util"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

// ==== Data Structures ====

type wxDay struct {
	Date string  `json:"date"`
	TMin float64 `json:"tMin"`
	TMax float64 `json:"tMax"`
	Code int     `json:"code"`
}
type wxOut struct {
	Source   string  `json:"source"`
	Updated  string  `json:"updated"`
	Timezone string  `json:"timezone"`
	Days     []wxDay `json:"daily"`
}

// API Response Schemas
type StandardResponse struct {
	Data any `json:"data"`
}

type SearchResponse struct {
	Data []Page `json:"data"`
}

type AuthResponse struct {
	StatusCode *int    `json:"statusCode"`
	Message    *string `json:"message"`
}

type ValidationError struct {
	Loc  []any  `json:"loc"`
	Msg  string `json:"msg"`
	Type string `json:"type"`
}

type HTTPValidationError struct {
	Detail []ValidationError `json:"detail"`
}

type RequestValidationError struct {
	StatusCode int     `json:"statusCode" default:"422"`
	Message    *string `json:"message"`
}

type cacheEntry struct {
	payload []byte
	expires time.Time
}

var (
	cacheTTL = 15 * time.Minute
	wxCache  = struct {
		mu sync.RWMutex
		m  map[string]cacheEntry
	}{m: make(map[string]cacheEntry)}
	db *sql.DB
)

// ==== Initialization ====

func init() {
	const dbPath = "/usr/src/app/data/whoknows.db"
	dbExists := true
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		dbExists = false
	}

	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	if !dbExists {
		if err := InitDB(db); err != nil {
			log.Fatalf("DB init failed: %v", err)
		}
	}
}

// ==== API Endpoints ====

func apiWeather(c *gin.Context) {
	lat, lon := "55.6761", "12.5683" // Copenhagen
	days, units := "5", "metric"
	key := lat + "|" + lon + "|" + days + "|" + units

	// Cache lookup
	wxCache.mu.RLock()
	ce, ok := wxCache.m[key]
	wxCache.mu.RUnlock()
	if ok && time.Now().Before(ce.expires) {
		c.Header("X-Cache", "HIT")
		c.Data(http.StatusOK, "application/json", ce.payload)
		return
	}

	url := buildOpenMeteoURL(lat, lon, days, units)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "who-knows-weather/1.0")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		c.JSON(http.StatusBadGateway, AuthResponse{nil, ptr("upstream unavailable")})
		return
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	out, tz, err := normalizeOpenMeteo(raw)
	if err != nil {
		c.JSON(http.StatusBadGateway, AuthResponse{nil, ptr("normalize failed")})
		return
	}
	out.Source = "open-meteo"
	out.Updated = time.Now().UTC().Format(time.RFC3339)
	out.Timezone = tz

	payload, _ := json.Marshal(StandardResponse{Data: out})
	wxCache.mu.Lock()
	wxCache.m[key] = cacheEntry{payload, time.Now().Add(cacheTTL)}
	wxCache.mu.Unlock()

	c.Header("X-Cache", "MISS")
	c.Data(http.StatusOK, "application/json", payload)
}

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
	log.Printf("[SEARCH] Search successful: q=%s, lang=%s", q, lang)
	c.JSON(http.StatusOK, SearchResponse{Data: results})
}

func apiLogin(c *gin.Context) {
	var creds struct {
		Username string `form:"username" json:"username"`
		Password string `form:"password" json:"password"`
	}
	if err := c.ShouldBind(&creds); err != nil {
		log.Printf("[LOGIN] Invalid form data: %v", err)
		c.JSON(422, HTTPValidationError{Detail: []ValidationError{{Loc: []any{"body", 0}, Msg: "invalid form data", Type: "validation_error"}}})
		return
	}

	id, _, _, hashed, err := GetUserByUsernameQuery(db, creds.Username)
	if err != nil {
		log.Printf("[LOGIN] Invalid username: %s", creds.Username)
		c.JSON(422, HTTPValidationError{Detail: []ValidationError{{Loc: []any{"username", 0}, Msg: "invalid username", Type: "auth_error"}}})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(hashed), []byte(creds.Password)) != nil {
		log.Printf("[LOGIN] Invalid password for username: %s", creds.Username)
		c.JSON(422, HTTPValidationError{Detail: []ValidationError{{Loc: []any{"password", 0}, Msg: "invalid password", Type: "auth_error"}}})
		return
	}

	util.SetAuthCookie(c, id)
	code := 200
	msg := "login successful"
	log.Printf("[LOGIN] Login successful for username: %s", creds.Username)
	c.JSON(http.StatusOK, AuthResponse{&code, &msg})
}

func apiRegister(c *gin.Context) {
	var form struct {
		Username  string `form:"username" json:"username"`
		Email     string `form:"email" json:"email"`
		Password  string `form:"password" json:"password"`
		Password2 string `form:"password2" json:"password2"`
	}
	if err := c.ShouldBind(&form); err != nil {
		log.Printf("[REGISTER] Invalid form data: %v", err)
		sendValidationError(c, "body", "invalid form data")
		return
	}

	if form.Username == "" {
		log.Printf("[REGISTER] Missing username")
		sendValidationError(c, "username", "you have to enter a username")
		return
	}
	if form.Email == "" || !regexp.MustCompile(`.+@.+\..+`).MatchString(form.Email) {
		log.Printf("[REGISTER] Invalid email: %q", form.Email)
		sendValidationError(c, "email", "you have to enter a valid email address")
		return
	}
	if form.Password == "" {
		log.Printf("[REGISTER] Missing password for username=%q", form.Username)
		sendValidationError(c, "password", "you have to enter a password")
		return
	}
	if form.Password != form.Password2 {
		log.Printf("[REGISTER] Password mismatch for username=%q", form.Username)
		sendValidationError(c, "password2", "the two passwords do not match")
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if _, err := InsertUserQuery(db, form.Username, form.Email, string(hash)); err != nil {
		log.Printf("[REGISTER] Database error: %v", err)
		c.JSON(422, HTTPValidationError{Detail: []ValidationError{{Loc: []any{"database", 0}, Msg: "username or email taken", Type: "db_error"}}})
		return
	}

	log.Printf("[REGISTER] User registered: %s", form.Username)

	code := 200
	msg := "user registered successfully"
	c.JSON(http.StatusOK, AuthResponse{&code, &msg})
}

func sendValidationError(c *gin.Context, field, msg string) {
	c.JSON(http.StatusUnprocessableEntity, HTTPValidationError{
		Detail: []ValidationError{{
			Loc:  []any{field, 0},
			Msg:  msg,
			Type: "validation_error",
		}},
	})
}

func apiLogout(c *gin.Context) {
	util.RemoveAuthCookie(c)
	code := 200
	msg := "logged out"
	log.Printf("[LOGOUT] Logout request from IP=%s", c.ClientIP())
	c.JSON(http.StatusOK, AuthResponse{&code, &msg})
}

// ==== Session Endpoint ====

func apiSession(c *gin.Context) {
	_, err := c.Cookie("user_id")
	if err != nil {
		log.Printf("[SESSION] Error getting user_id cookie from IP=%s: %v", c.ClientIP(), err)
		code := 401
		msg := "not logged in"
		c.JSON(http.StatusOK, AuthResponse{&code, &msg})
		return
	}

	log.Printf("[SESSION] Valid session detected from IP=%s", c.ClientIP())
	
	code := 200
	msg := "logged in"
	c.JSON(http.StatusOK, AuthResponse{&code, &msg})
}

// ==== Static / HTML Endpoints ====

func serveHTML(c *gin.Context, path string) {
	c.File(path)
}

func serveIndexFile(c *gin.Context)    { serveHTML(c, "./public/index.html") }
func serveLoginFile(c *gin.Context)    { serveHTML(c, "./public/login.html") }
func serveRegisterFile(c *gin.Context) { serveHTML(c, "./public/register.html") }
func serverWeatherFile(c *gin.Context) { serveHTML(c, "./public/weather.html") }

// ==== Helper ====

func buildOpenMeteoURL(lat, lon, days, units string) string {
	tempUnit := "celsius"
	if units == "imperial" {
		tempUnit = "fahrenheit"
	}
	return "https://api.open-meteo.com/v1/forecast?latitude=" + lat +
		"&longitude=" + lon +
		"&timezone=auto&daily=weathercode,temperature_2m_max,temperature_2m_min" +
		"&forecast_days=" + days +
		"&temperature_unit=" + tempUnit
}

type omResp struct {
	Daily struct {
		Time        []string  `json:"time"`
		WeatherCode []int     `json:"weathercode"`
		TempMax     []float64 `json:"temperature_2m_max"`
		TempMin     []float64 `json:"temperature_2m_min"`
	} `json:"daily"`
	Timezone string `json:"timezone"`
}

func normalizeOpenMeteo(raw []byte) (wxOut, string, error) {
	var r omResp
	if err := json.Unmarshal(raw, &r); err != nil {
		return wxOut{}, "", err
	}
	o := wxOut{}
	for i := range r.Daily.Time {
		o.Days = append(o.Days, wxDay{
			Date: r.Daily.Time[i],
			TMin: r.Daily.TempMin[i],
			TMax: r.Daily.TempMax[i],
			Code: r.Daily.WeatherCode[i],
		})
	}
	return o, r.Timezone, nil
}

func ptr[T any](v T) *T { return &v }

// ==== Middleware ====

func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Printf("[REQ] %s %s -> %d (%v)", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(start))
	}
}

// ==== Main ====

func main() {
	defer db.Close()
	router := gin.New()
	router.Use(gin.Recovery(), loggingMiddleware())

	api := router.Group("/api")
	{
		api.GET("/weather", apiWeather)
		api.GET("/search", apiSearch)
		api.POST("/login", apiLogin)
		api.POST("/register", apiRegister)
		api.GET("/logout", apiLogout)
		api.GET("/session", apiSession)
	}

	router.GET("/", serveIndexFile)
	router.GET("/login", serveLoginFile)
	router.GET("/register", serveRegisterFile)
	router.GET("/weather", serverWeatherFile)
	router.Static("/public", "./public")

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
