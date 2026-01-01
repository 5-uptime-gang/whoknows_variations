package main

import "github.com/gin-gonic/gin"

func newRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), loggingMiddleware(), BrowserMiddleware())

	router.GET("/metrics", metricsEndpoint)

	api := router.Group("/api")
	{
		api.GET("/weather", apiWeather)
		api.GET("/search", apiSearch)
		api.POST("/login", apiLogin)
		api.POST("/register", apiRegister)
		api.GET("/logout", apiLogout)
		api.GET("/session", apiSession)
	}

	router.GET("/docs", serveSwaggerUI)
	router.GET("/docs/swagger.yaml", serveSwaggerSpecYaml)
	router.GET("/docs/swagger.json", serveSwaggerSpecJSON)
	router.GET("/", serveIndexFile)
	router.GET("/login", serveLoginFile)
	router.GET("/register", serveRegisterFile)
	router.GET("/weather", serverWeatherFile)
	router.GET("/about", serveAboutFile)
	router.Static("/public", "./public")

	return router
}
