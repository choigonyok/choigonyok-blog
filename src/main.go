package main

import (
	"net/http"
	_ "net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	eg := gin.Default()
	eg.LoadHTMLGlob("./templates/**/*.html")
	eg.Static("/assets", "./assets")
	//왜 images 폴더는 캐싱이 안될까?

	// default homepage
	eg.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// redirect to homepage
	eg.GET("/index.html", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/")
	})

	// generic page
	eg.GET("/generic.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "generic.html", nil)
	})

	// elements page
	eg.GET("/elements.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "elements.html", nil)
	})

	eg.Run(":8080")
}
