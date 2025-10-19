package main

import (
	"WHOKNOWS_VARIATIONS/util"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

// ==== Users + Auth ====

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"` // "-" skjuler password i JSON output
}

// --- Weather types + tiny cache ---
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
type apiEnvelope struct {
	Data any `json:"data"`
}
type cacheEntry struct {
	payload []byte
	expires time.Time
}

var (
	cacheTimeToLiveMinutes = 15
	cacheTTL               = time.Duration(cacheTimeToLiveMinutes) * time.Minute
	wxCache                = struct {
		mu sync.RWMutex
		m  map[string]cacheEntry
	}{m: make(map[string]cacheEntry)}
)

// ==== Database initializer ====
var db *sql.DB

func init() {
	const dbPath = "/usr/src/app/data/whoknows.db"

	// If the DB file doesn't exist, we'll need to initialize it
	dbExists := true
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		dbExists = false
	}

	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database at %s: %v", dbPath, err)
	}

	if !dbExists {
		log.Println("Database not found — initializing schema and seed data...")
		if err := InitDB(db); err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}
		log.Println("Database initialized successfully")
	} else {
		log.Println("Database already exists — skipping initialization")
	}
}

// ==== API Endpoints ====

func apiWeather(c *gin.Context) {
	lat, lon := "55.6761", "12.5683" // Copenhagen
	days, units := "5", "metric"
	key := lat + "|" + lon + "|" + days + "|" + units

	// 1) cache
	wxCache.mu.RLock()
	ce, ok := wxCache.m[key]
	wxCache.mu.RUnlock()
	if ok && time.Now().Before(ce.expires) {
		c.Header("X-Cache", "HIT")
		c.Data(http.StatusOK, "application/json", ce.payload)
		return
	}

	// 2) upstream call
	url := buildOpenMeteoURL(lat, lon, days, units)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "request build failed"})
		return
	}
	req.Header.Set("User-Agent", "who-knows-weather/1.0 (+contact: example@example.com)")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "upstream unavailable"})
		return
	}
	// make errcheck happy + do the right thing
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("resp body close error: %v", cerr)
		}
	}()

	if resp.StatusCode >= 500 {
		c.JSON(http.StatusBadGateway, gin.H{"error": "upstream unavailable"})
		return
	}
	if resp.StatusCode >= 400 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request to upstream"})
		return
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "read upstream failed"})
		return
	}

	// 3) normalize provider JSON to our stable shape
	out, tz, err := normalizeOpenMeteo(raw)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "normalize failed"})
		return
	}
	out.Source = "open-meteo"
	out.Updated = time.Now().UTC().Format(time.RFC3339)
	out.Timezone = tz

	// 4) wrap as { "data": ... } to match your spec
	payload, err := json.Marshal(apiEnvelope{Data: out})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "marshal failed"})
		return
	}

	// 5) save to cache + return
	wxCache.mu.Lock()
	wxCache.m[key] = cacheEntry{payload: payload, expires: time.Now().Add(cacheTTL)}
	wxCache.mu.Unlock()

	c.Header("X-Cache", "MISS")
	c.Data(http.StatusOK, "application/json", payload)
}

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

func apiLogin(c *gin.Context) {
	var creds struct {
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	}
	if err := c.ShouldBind(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Brug GetUserByUsernameQuery fra queries.go
	id, username, email, hashedPassword, err := GetUserByUsernameQuery(db, creds.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	util.SetAuthCookie(c, id)

	c.JSON(http.StatusOK, gin.H{
		"message":  "login successful",
		"user_id":  id,
		"username": username,
		"email":    email,
	})
}

func apiRegister(c *gin.Context) {
	var form struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Password2 string `json:"password2"`
	}

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Validation
	if form.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you have to enter a username"})
		return
	}
	if form.Email == "" || !regexp.MustCompile(`.+@.+\..+`).MatchString(form.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you have to enter a valid email address"})
		return
	}
	if form.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you have to enter a password"})
		return
	}
	if form.Password != form.Password2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "the two passwords do not match"})
		return
	}

	fmt.Println("username: ", form.Username)
	fmt.Println("Email:", form.Email)
	fmt.Println("password1: ", form.Password)
	fmt.Println("password2: ", form.Password2)

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}

	// Brug InsertUserQuery fra queries.go
	userID, err := InsertUserQuery(db, form.Username, form.Email, string(hash))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or email already taken"})
		return
	}

	util.SetAuthCookie(c, int(userID))

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"user_id": userID,
	})
}

func apiLogout(c *gin.Context) {
	// overwrite cookie with empty value and expired time
	util.RemoveAuthCookie(c)

	c.JSON(http.StatusOK, gin.H{
		"message": "logged out",
		"status":  "ok",
	})
}

func apiSearch(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		// q er obligatorisk ifølge openAPI - derfor skal der bruges q i URL hvis man ønsker at finde noget.
		c.JSON(422, gin.H{
			"statusCode": 422,
			"message":    "Query parameter 'q' is required",
		})
		return
	}

	lang := c.DefaultQuery("language", "en") // Default til engelsk

	results, err := SearchPagesQuery(db, q, lang)
	if err != nil {
		// Hvis search fejler returnerer vi 422
		c.JSON(422, gin.H{
			"statusCode": 422,
			"message":    "Search failed: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data": results,
	})
}

func apiSession(c *gin.Context) {
	_, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"logged_in": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"logged_in": true})
}

func serveLoginRegisterFiles(c *gin.Context, fp string) {
	// Debug: confirm file exists and size
	if info, err := filepath.Abs(fp); err == nil {
		log.Println("Serving:", info)
	}
	// If it doesn't exist, Gin would 404. We'll log size after write below.

	// Important: don't return before writing the body
	c.File(fp) // sets 200 + streams file if found
}

func serveLoginFile(c *gin.Context) {
	serveLoginRegisterFiles(c, "./public/login.html")
}

func serveRegisterFile(c *gin.Context) {
	serveLoginRegisterFiles(c, "./public/register.html")
}

func serverWeatherFile(c *gin.Context) {
	serveLoginRegisterFiles(c, "./public/weather.html")
}

func serverAboutFile(c *gin.Context) {
	serveLoginRegisterFiles(c, "./public/about.html")
}

func serveIndexFile(c *gin.Context) {
	serveLoginRegisterFiles(c, "./public/index.html")
}

// ==== Main entry ====

func main() {
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing DB: %v", err)
		}
	}()
	router := gin.Default()

	const PORT = ":8080"

	api := router.Group("/api")
	{
		api.POST("/login", apiLogin)
		api.POST("/register", apiRegister)
		api.POST("/logout", apiLogout)
		api.GET("/search", apiSearch)
		api.GET("/session", apiSession)
		api.GET("/weather", apiWeather)
	}

	router.GET("/", serveIndexFile)
	router.GET("/login", serveLoginFile)
	router.GET("/register", serveRegisterFile)
	router.GET("/weather", serverWeatherFile)
	router.GET("/about", serverAboutFile)

	// This makes everything in ./public available under /public
	router.Static("/public", "./public")

	if err := router.Run(PORT); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
