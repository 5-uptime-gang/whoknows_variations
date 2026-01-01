package main

import "github.com/gin-gonic/gin"

func serveHTML(c *gin.Context, path string) {
	c.File(path)
}

// serveIndexFile godoc
// @Summary Serve index page
// @Tags Pages
// @Produce html
// @Success 200 {string} string "HTML page"
// @Router / [get]
func serveIndexFile(c *gin.Context) { serveHTML(c, "./public/index.html") }

// serveLoginFile godoc
// @Summary Serve login page
// @Tags Pages
// @Produce html
// @Success 200 {string} string "HTML page"
// @Router /login [get]
func serveLoginFile(c *gin.Context) { serveHTML(c, "./public/login.html") }

// serveRegisterFile godoc
// @Summary Serve registration page
// @Tags Pages
// @Produce html
// @Success 200 {string} string "HTML page"
// @Router /register [get]
func serveRegisterFile(c *gin.Context) { serveHTML(c, "./public/register.html") }

// serverWeatherFile godoc
// @Summary Serve weather page
// @Tags Pages
// @Produce html
// @Success 200 {string} string "HTML page"
// @Router /weather [get]
func serverWeatherFile(c *gin.Context) { serveHTML(c, "./public/weather.html") }

// serveAboutFile godoc
// @Summary Serve about page
// @Tags Pages
// @Produce html
// @Success 200 {string} string "HTML page"
// @Router /about [get]
func serveAboutFile(c *gin.Context) { serveHTML(c, "./public/about.html") }

// staticAssetsDoc is a documentation-only stub for the /public static file server.
// @Summary Serve static assets
// @Tags Assets
// @Produce octet-stream
// @Param filepath path string true "Asset path under /public (e.g. css/main.css)"
// @Success 200 {string} string "Static asset bytes"
// @Router /public/{filepath} [get]
func staticAssetsDoc() {}

// serveSwaggerSpecYaml godoc
// @Summary Serve generated Swagger spec (YAML)
// @Tags Docs
// @Produce plain
// @Success 200 {string} string "swagger.yaml"
// @Router /docs/swagger.yaml [get]
func serveSwaggerSpecYaml(c *gin.Context) { serveHTML(c, "./docs/api/swagger.yaml") }

// serveSwaggerSpecJSON godoc
// @Summary Serve generated Swagger spec (JSON)
// @Tags Docs
// @Produce json
// @Success 200 {string} string "swagger.json"
// @Router /docs/swagger.json [get]
func serveSwaggerSpecJSON(c *gin.Context) { serveHTML(c, "./docs/api/swagger.json") }

// serveSwaggerUI godoc
// @Summary Interactive Swagger UI
// @Tags Docs
// @Produce html
// @Success 200 {string} string "Swagger UI"
// @Router /docs [get]
func serveSwaggerUI(c *gin.Context) {
	const html = `<!doctype html>
<html>
  <head>
    <meta charset="utf-8">
    <title>API Docs</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.onload = () => {
        SwaggerUIBundle({
          url: '/docs/api/swagger.yaml',
          dom_id: '#swagger-ui'
        });
      };
    </script>
  </body>
</html>`
	c.Data(200, "text/html; charset=utf-8", []byte(html))
}
