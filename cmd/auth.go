package main

import (
	"log"
	"net/http"
	"regexp"

	"WHOKNOWS_VARIATIONS/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

type RegisterRequest struct {
	Username  string `form:"username" json:"username"`
	Email     string `form:"email" json:"email"`
	Password  string `form:"password" json:"password"`
	Password2 string `form:"password2" json:"password2"`
}

// apiLogin godoc
// @Summary Log in and set auth cookie
// @Tags Auth
// @Accept json
// @Accept x-www-form-urlencoded
// @Produce json
// @Param request body LoginRequest true "Credentials"
// @Success 200 {object} AuthResponse
// @Failure 422 {object} HTTPValidationError
// @Router /api/login [post]
func apiLogin(c *gin.Context) {
	var creds LoginRequest
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

// apiRegister godoc
// @Summary Register a new user and set auth cookie
// @Tags Auth
// @Accept json
// @Accept x-www-form-urlencoded
// @Produce json
// @Param request body RegisterRequest true "Registration payload"
// @Success 200 {object} AuthResponse
// @Failure 422 {object} HTTPValidationError
// @Router /api/register [post]
func apiRegister(c *gin.Context) {
	var form RegisterRequest
	if err := c.ShouldBind(&form); err != nil {
		log.Printf("[REGISTER] Invalid form data: %v", err)
		userSignupCounter.WithLabelValues("failed").Inc()
		sendValidationError(c, "body", "invalid form data")
		return
	}

	if form.Username == "" {
		log.Printf("[REGISTER] Missing username")
		userSignupCounter.WithLabelValues("failed").Inc()
		sendValidationError(c, "username", "you have to enter a username")
		return
	}
	if form.Email == "" || !regexp.MustCompile(`.+@.+\..+`).MatchString(form.Email) {
		log.Printf("[REGISTER] Invalid email: %q", form.Email)
		userSignupCounter.WithLabelValues("failed").Inc()
		sendValidationError(c, "email", "you have to enter a valid email address")
		return
	}
	if form.Password == "" {
		log.Printf("[REGISTER] Missing password for username=%q", form.Username)
		userSignupCounter.WithLabelValues("failed").Inc()
		sendValidationError(c, "password", "you have to enter a password")
		return
	}
	if form.Password != form.Password2 {
		log.Printf("[REGISTER] Password mismatch for username=%q", form.Username)
		userSignupCounter.WithLabelValues("failed").Inc()
		sendValidationError(c, "password2", "the two passwords do not match")
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	userID, err := InsertUserQuery(db, form.Username, form.Email, string(hash))
	if err != nil {
		log.Printf("[REGISTER] Database error: %v", err)
		userSignupCounter.WithLabelValues("failed").Inc()
		c.JSON(422, HTTPValidationError{Detail: []ValidationError{{Loc: []any{"database", 0}, Msg: "username or email taken", Type: "db_error"}}})
		return
	}

	log.Printf("[REGISTER] User registered: %s", form.Username)
	userSignupCounter.WithLabelValues("success").Inc()

	util.SetAuthCookie(c, int(userID))
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

// apiLogout godoc
// @Summary Clear the auth cookie
// @Tags Auth
// @Produce json
// @Success 200 {object} AuthResponse
// @Router /api/logout [get]
func apiLogout(c *gin.Context) {
	util.RemoveAuthCookie(c)
	code := 200
	msg := "logged out"
	log.Printf("[LOGOUT] Logout request from IP=%s", c.ClientIP())
	c.JSON(http.StatusOK, AuthResponse{&code, &msg})
}

// apiSession godoc
// @Summary Report session state based on auth cookie
// @Tags Auth
// @Produce json
// @Success 200 {object} AuthResponse "statusCode 200 if cookie is present, 401 otherwise"
// @Router /api/session [get]
func apiSession(c *gin.Context) {
	_, err := c.Cookie("user_id")
	if err != nil {
		log.Printf("[SESSION] Missing user_id cookie")
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
