package main

import (
	"WHOKNOWS_VARIATIONS/util"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"


  ginSwagger "github.com/swaggo/gin-swagger"
  swaggerFiles "github.com/swaggo/files"
)

// @title WHOKNOWS_VARIATIONS - 5-uptime-gang - API
// @version 1.0
// @host 68.221.201.252:8080
// @BasePath /api

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

// @Summary Login existing user
// @Tags API
// @Accept json
// @Produce json
// @Param username formData string true "username"
// @Param password formData string true "password"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
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


type RegisterRequest struct {
    Username  string `json:"username" example:"testuser"`
    Email     string `json:"email" example:"test@example.com"`
    Password  string `json:"password" example:"supersecret"`
    Password2 string `json:"password2" example:"supersecret"`
}
// @Summary Create new user
// @Tags API
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /register [post]
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


type SessionResponse struct {
    LoggedIn bool `json:"logged_in" example:"true"`
}

// @Summary Logout the current user
// @Tags API
// @Produce json
// @Success 200 {object} map[string]string "Logout confirmation"
// @Router /logout [post]
func apiLogout(c *gin.Context) {
	// overwrite cookie with empty value and expired time
	util.RemoveAuthCookie(c)

	c.JSON(http.StatusOK, gin.H{
		"message": "logged out",
		"status":  "ok",
	})
}


// @Summary Search after pages
// @Tags API
// @Produce json
// @Param q query string true "search word"
// @Param language query string false "(default=en)"
// @Success 200 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Router /search [get]
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

// @Summary checking the users session
// @Tags API
// @Produce json
// @Success 200 {object} SessionResponse "Status for login"
// @Router /session [get]
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

// @Summary getting login-page
// @Tags Pages
// @Produce html
// @Success 200 "HTML page"
// @Router /login [get]
func serveLoginFile(c *gin.Context) {
	serveLoginRegisterFiles(c, "./public/login.html")
}

// @Summary getting register-page
// @Tags Pages
// @Produce html
// @Success 200 "HTML page"
// @Router /register [get]
func serveRegisterFile(c *gin.Context) {
	serveLoginRegisterFiles(c, "./public/register.html")
}

// @Summary get weather-page
// @Tags Pages
// @Produce html
// @Success 200 "HTML page"
// @Router /weather [get]
func serverWeatherFile(c *gin.Context) {
	serveLoginRegisterFiles(c, "./public/weather.html")
}

// @Summary get about-page
// @Tags Pages
// @Produce html
// @Success 200 "HTML page"
// @Router /about [get]
func serverAboutFile(c *gin.Context) {
	serveLoginRegisterFiles(c, "./public/about.html")
}

// @Summary get index-page
// @Tags Pages
// @Produce html
// @Success 200 "HTML page"
// @Router / [get]
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
	}

	router.GET("/", serveIndexFile)
	router.GET("/login", serveLoginFile)
	router.GET("/register", serveRegisterFile)
	router.GET("/weather", serverWeatherFile)
	router.GET("/about", serverAboutFile)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))


	// This makes everything in ./public available under /public
	router.Static("/public", "./public")

	if err := router.Run(PORT); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
