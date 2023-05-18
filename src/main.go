package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("hi")
	http.ListenAndServe(":8080", nil)
	eg := gin.Default()
	eg.GET("/")
}
