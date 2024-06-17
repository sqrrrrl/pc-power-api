package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pc-power-api/src/controller"
	"github.com/pc-power-api/src/controller/middleware"
	"log"
	"os"
)

func main() {
	r := gin.Default()
	r.Use(middleware.ExceptionHandler())
	_ = r.SetTrustedProxies(nil)
	port := os.Getenv("PORT")

	controller.NewDevicesHandler(r)

	log.Fatal(r.Run(":" + port))
}
