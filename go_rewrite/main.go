package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	// swagger (docs will be generated later; the blank import is fine even before docs exist)
	_ "go_rewrite/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// ---------- API metadata ----------
// @title           WhoKnows
// @version         0.1.0
// @description     OpenAPI generated from Go code to match the provided specification.
// @BasePath        /
// @schemes         http
// @accept          json
// @produce         json

// ---------- Schemas (components) ----------
type AuthResponse struct {
	StatusCode *int    `json:"statusCode,omitempty"` // nullable int
	Message    *string `json:"message,omitempty"`    // nullable string
}

type SearchResponse struct {
	Data []map[string]any `json:"data" binding:"required"`
}

type StandardResponse struct {
	Data map[string]any `json:"data" binding:"required"`
}

type ValidationError struct {
	Loc  []any  `json:"loc" binding:"required"` // (string|int)[]
	Msg  string `json:"msg" binding:"required"`
	Type string `json:"type" binding:"required"`
}

type HTTPValidationError struct {
	Detail []ValidationError `json:"detail"`
}

type RequestValidationError struct {
	StatusCode int     `json:"statusCode" example:"422"`
	Message    *string `json:"message,omitempty"`
}

// ---------- Form bodies ----------
type LoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}
type RegisterForm struct {
	Username  string `form:"username" binding:"required"`
	Email     string `form:"email"    binding:"required"`
	Password  string `form:"password" binding:"required"`
	Password2 string `form:"password2"` // optional
}

func main() {
	r := gin.Default()

	// ---------- HTML pages (return text/html as string) ----------

	// @Summary      Serve Root Page
	// @ID           serve_root_page__get
	// @Produce      html
	// @Success      200 {string} string "HTML"
	// @Router       / [get]
	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html", []byte("<h1>WhoKnows</h1>"))
	})

	// @Summary      Serve Weather Page
	// @ID           serve_weather_page_weather_get
	// @Produce      html
	// @Success      200 {string} string "HTML"
	// @Router       /weather [get]
	r.GET("/weather", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html", []byte("<h1>Weather</h1>"))
	})

	// @Summary      Serve Register Page
	// @ID           serve_register_page_register_get
	// @Produce      html
	// @Success      200 {string} string "HTML"
	// @Router       /register [get]
	r.GET("/register", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html", []byte("<h1>Register</h1>"))
	})

	// @Summary      Serve Login Page
	// @ID           serve_login_page_login_get
	// @Produce      html
	// @Success      200 {string} string "HTML"
	// @Router       /login [get]
	r.GET("/login", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html", []byte("<h1>Login</h1>"))
	})

	// ---------- JSON API ----------

	// @Summary      Search
	// @ID           search_api_search_get
	// @Produce      json
	// @Param        q         query string true  "Q"
	// @Param        language  query string false "Language code (e.g., 'en')"
	// @Success      200 {object} SearchResponse
	// @Failure      422 {object} RequestValidationError
	// @Router       /api/search [get]
	r.GET("/api/search", func(c *gin.Context) {
		q := c.Query("q")
		if q == "" {
			msg := "Missing required query parameter: q"
			c.JSON(422, RequestValidationError{StatusCode: 422, Message: &msg})
			return
		}
		lang := c.Query("language")
		resp := map[string]any{"q": q}
		if lang != "" {
			resp["language"] = lang
		}
		c.JSON(http.StatusOK, SearchResponse{Data: []map[string]any{resp}})
	})

	// @Summary      Weather
	// @ID           weather_api_weather_get
	// @Produce      json
	// @Success      200 {object} StandardResponse
	// @Router       /api/weather [get]
	r.GET("/api/weather", func(c *gin.Context) {
		c.JSON(http.StatusOK, StandardResponse{Data: map[string]any{
			"temp": 20, "unit": "C",
		}})
	})

	// @Summary      Register
	// @ID           register_api_register_post
	// @Accept       x-www-form-urlencoded
	// @Produce      json
	// @Param        username formData string true  "Username"
	// @Param        email    formData string true  "Email"
	// @Param        password formData string true  "Password"
	// @Param        password2 formData string false "Password2"
	// @Success      200 {object} AuthResponse
	// @Failure      422 {object} HTTPValidationError
	// @Router       /api/register [post]
	r.POST("/api/register", func(c *gin.Context) {
		var f RegisterForm
		if err := c.ShouldBind(&f); err != nil {
			c.JSON(422, HTTPValidationError{Detail: []ValidationError{
				{Loc: []any{"body"}, Msg: "invalid form", Type: "value_error"},
			}})
			return
		}
		code := 200
		msg := "registered"
		c.JSON(http.StatusOK, AuthResponse{StatusCode: &code, Message: &msg})
	})

	// @Summary      Login
	// @ID           login_api_login_post
	// @Accept       x-www-form-urlencoded
	// @Produce      json
	// @Param        username formData string true "Username"
	// @Param        password formData string true "Password"
	// @Success      200 {object} AuthResponse
	// @Failure      422 {object} HTTPValidationError
	// @Router       /api/login [post]
	r.POST("/api/login", func(c *gin.Context) {
		var f LoginForm
		if err := c.ShouldBind(&f); err != nil {
			c.JSON(422, HTTPValidationError{Detail: []ValidationError{
				{Loc: []any{"body"}, Msg: "invalid form", Type: "value_error"},
			}})
			return
		}
		code := http.StatusOK
		msg := "ok"
		c.JSON(http.StatusOK, AuthResponse{StatusCode: &code, Message: &msg})
	})

	// @Summary      Logout
	// @ID           logout_api_logout_get
	// @Produce      json
	// @Success      200 {object} AuthResponse
	// @Router       /api/logout [get]
	r.GET("/api/logout", func(c *gin.Context) {
		code := http.StatusOK
		msg := "ok"
		c.JSON(http.StatusOK, AuthResponse{StatusCode: intPtr(code), Message: &msg})
	})

	// Swagger UI endpoint (will work after we generate docs)
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	r.Run(":8080")
}

func intPtr(i int) *int { return &i }