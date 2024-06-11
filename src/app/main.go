package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pc-power-api/src/controller/middleware"
	"log"
	"net/http"
	"os"
)

func main() {
	r := gin.Default()
	r.Use(middleware.ExceptionHandler())
	_ = r.SetTrustedProxies(nil)
	port := os.Getenv("PORT")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	log.Fatal(r.Run(":" + port))
}
