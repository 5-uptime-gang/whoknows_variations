package main

import "github.com/gin-gonic/gin"

func serveHTML(c *gin.Context, path string) {
	c.File(path)
}

func serveIndexFile(c *gin.Context)    { serveHTML(c, "./public/index.html") }
func serveLoginFile(c *gin.Context)    { serveHTML(c, "./public/login.html") }
func serveRegisterFile(c *gin.Context) { serveHTML(c, "./public/register.html") }
func serverWeatherFile(c *gin.Context) { serveHTML(c, "./public/weather.html") }
func serveAboutFile(c *gin.Context)    { serveHTML(c, "./public/about.html") }
