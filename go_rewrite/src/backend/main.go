package main

import (
    "net/http"
	"fmt"
	"regexp"

    "github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func main() {
    router := gin.Default()
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
    fmt.Println(string(hash))

    router.GET("/albums", getAlbums)
    router.POST("/albums", postAlbums)
	router.POST("/api/login", apiLogin)



    router.Run("localhost:8080")
}

type User struct {
    ID       string
    Username string
    Password string // hashed password
}

var users = []User{
    {
        ID:       "1",
        Username: "alice",
        // bcrypt hash of "password123"
        Password: "$2a$10$s/f.UkN1UVrdLL6Yk8oRku5UoZRG1aaMxlYBDDdLD/LKDwjmucxD6",
    },
    {
        ID:       "2",
        Username: "bob",
        // bcrypt hash of "secret"
        Password: "$2a$10$sFC0YeE48GMmLjhx96lZT.qXcWcyC1suCq8/nIVPcOTQDM4Aaq9Di",
    },
}

// Beginning of rewrite API handlers

func apiLogin(c *gin.Context) {
    var creds struct {
        Username string `json:"username" form:"username"`
        Password string `json:"password" form:"password"`
    }

    if err := c.ShouldBind(&creds); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    // Look up user in our fake "DB"
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

    // Compare password hash
    if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
        return
    }

    // Login successful
    c.JSON(http.StatusOK, gin.H{
        "message": "login successful",
        "user_id": user.ID,
    })
}
















// Learning go testing
type album struct {
    ID     string  `json:"id"`
    Title  string  `json:"title"`
    Artist string  `json:"artist"`
    Price  float64 `json:"price"`
}

var albums = []album{
    {ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
    {ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
    {ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context) {
    var newAlbum album

    // Call BindJSON to bind the received JSON to
    // newAlbum.
    if err := c.BindJSON(&newAlbum); err != nil {
        return
    }

    // Add the new album to the slice.
    albums = append(albums, newAlbum)
    c.IndentedJSON(http.StatusCreated, newAlbum)
}


