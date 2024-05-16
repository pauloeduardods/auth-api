package static

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupStaticFiles(r *gin.Engine) {
	r.Static("/static", "static")

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "static/home.html")
	})

	r.GET("/home", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "static/home.html")
	})

	r.GET("/login", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "static/login.html")
	})

	r.GET("/register", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "static/register.html")
	})

}
