package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ==== Users + Auth ====

type User struct {
	ID       string
	Username string
	Email    string
	Password string // bcrypt hash
}

var users = []User{
	{
		ID:       "1",
		Username: "alice",
		Password: "$2a$10$s/f.UkN1UVrdLL6Yk8oRku5UoZRG1aaMxlYBDDdLD/LKDwjmucxD6", // "password123"
	},
	{
		ID:       "2",
		Username: "bob",
		Password: "$2a$10$sFC0YeE48GMmLjhx96lZT.qXcWcyC1suCq8/nIVPcOTQDM4Aaq9Di", // "secret"
	},
}

var nextUserID = 3

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

	var user *User
	for _, u := range users {
		if u.Username == creds.Username {
			user = &u
			break
		}
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"user_id": user.ID,
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
	for _, u := range users {
		if u.Username == form.Username {
			c.JSON(http.StatusBadRequest, gin.H{"error": "the username is already taken"})
			return
		}
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}

	newUser := User{
		ID:       strconv.Itoa(nextUserID),
		Username: form.Username,
		Email:    form.Email,
		Password: string(hash),
	}
	nextUserID++
	users = append(users, newUser)

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"user_id": newUser.ID,
	})
}

// ==== Main entry ====

func main() {
	router := gin.Default()
	fmt.Println("Starting server on http://localhost:8080")

	router.POST("/api/login", apiLogin)
	router.POST("/api/register", apiRegister)

	router.Run("localhost:8080")
}