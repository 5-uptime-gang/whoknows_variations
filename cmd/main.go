package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

// ==== Weather types + tiny cache ====
type wxDaily struct {
	Date string  `json:"date"`
	TMin float64 `json:"tMin"`
	TMax float64 `json:"tMax"`
	Code int     `json:"code"`
}
type wxOut struct {
	Source   string    `json:"source"`
	Updated  string    `json:"updated"`
	Timezone string    `json:"timezone"`
	Days     []wxDaily `json:"daily"`
}

type cacheEntry struct {
	payload []byte
	expires time.Time
}

var (
	wxCache = struct {
		mu sync.RWMutex
		m  map[string]cacheEntry
	}{m: make(map[string]cacheEntry)}
	cacheTTL = 15 * time.Minute
)

// ==== Users + Auth ====

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"` // "-" skjuler password i JSON output
}

// ==== Database initializer ====
var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite", "whoknows.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Initialize database schema
	if err := InitDB(db); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database initialized successfully")
}

// ==== API Endpoints ====

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

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"user_id": userID,
	})
}

func apiLogout(c *gin.Context) {
	// ingen logik/tilstand — bare et simpelt svar
	c.JSON(http.StatusOK, gin.H{
		"message": "logged out",
		"status":  "ok",
	})
}

func apiSearch(c *gin.Context) {
	q := c.Query("q")
	lang := c.DefaultQuery("lang", "en") // Default til engelsk hvis ikke specificeret

	// Brug SearchPagesQuery fra queries.go
	results, err := SearchPagesQuery(db, q, lang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query":   q,
		"count":   len(results),
		"results": results,
	})
}

func apiWeather(c *gin.Context) {
	lat := c.Query("lat")
	lon := c.Query("lon")
	if lat == "" || lon == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing lat/lon"})
		return
	}
	days := c.DefaultQuery("days", "5")
	units := c.DefaultQuery("units", "metric")
	key := lat + "|" + lon + "|" + days + "|" + units

	// 1) cache check
	wxCache.mu.RLock()
	ce, ok := wxCache.m[key]
	wxCache.mu.RUnlock()
	if ok && time.Now().Before(ce.expires) {
		c.Header("X-Cache", "HIT")
		c.Data(http.StatusOK, "application/json", ce.payload)
		return
	}

	// 2) call Open-Meteo
	url := buildOpenMeteoURL(lat, lon, days, units)
	req, _ := http.NewRequest("GET", url, nil)
	// A polite UA (and useful if you later add MET Norway which requires it)
	req.Header.Set("User-Agent", "who-knows-weather/1.0 (+contact: you@example.com)")

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil || resp.StatusCode >= 500 {
		c.JSON(http.StatusBadGateway, gin.H{"error": "upstream unavailable"})
		return
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)

	// 3) normalize provider JSON to our tiny shape
	out, tz, err := normalizeOpenMeteo(raw)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "normalize failed"})
		return
	}
	out.Source = "open-meteo"
	out.Updated = time.Now().UTC().Format(time.RFC3339)
	out.Timezone = tz

	payload, _ := json.Marshal(out)

	// 4) cache + return
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
	// daily fields: weather code + min/max temps; timezone auto
	return fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&timezone=auto&daily=weathercode,temperature_2m_max,temperature_2m_min&forecast_days=%s&temperature_unit=%s",
		lat, lon, days, tempUnit)
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
	out := wxOut{}
	for i := range r.Daily.Time {
		out.Days = append(out.Days, wxDaily{
			Date: r.Daily.Time[i],
			TMin: r.Daily.TempMin[i],
			TMax: r.Daily.TempMax[i],
			Code: r.Daily.WeatherCode[i],
		})
	}
	return out, r.Timezone, nil
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
	serveLoginRegisterFiles(c, "../public/login.html")
}

func serveRegisterFile(c *gin.Context) {
	serveLoginRegisterFiles(c, "../public/register.html")
}
func serveIndexFile(c *gin.Context) {
	serveLoginRegisterFiles(c, "../public/index.html")
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
		api.GET("/logout", apiLogout)
		api.POST("/logout", apiLogout)
		api.GET("/search", apiSearch)
		api.GET("/weather", apiWeather)
	}

	router.GET("/", serveIndexFile)
	router.GET("/login", serveLoginFile)
	router.GET("/register", serveRegisterFile)

	// Maps /css, /js, /images to ./public/css, ./public/js, ./public/images
	// So it can be used in HTML like <link href="/css/styles.css">
	router.Static("/css", "../public/css")
	router.Static("/js", "../public/js")
	router.Static("/images", "../public/images") // or /img if you use that

	if err := router.Run(PORT); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
