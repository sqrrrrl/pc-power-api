package main

import (
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/pc-power-api/src/controller"
	"github.com/pc-power-api/src/controller/middleware"
	"github.com/pc-power-api/src/infra/entity"
	"github.com/pc-power-api/src/infra/repo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	r := gin.Default()
	_ = r.SetTrustedProxies(nil)

	db := connectDatabase()
	err := db.AutoMigrate(&entity.User{}, &entity.Device{})
	if err != nil {
		log.Fatal(err)
	}

	deviceRepository := repo.NewDeviceRepository(db)
	userRepository := repo.NewUserRepository(db)

	authenticationMiddleWare := middleware.NewAuthenticationMiddleware(userRepository)
	authMiddlewareHandlerFunction, authMiddlewareHandler := authenticationMiddleWare.AuthMiddleware()

	r.Use(authMiddlewareHandlerFunction)
	r.Use(middleware.ExceptionHandler())

	controller.NewAuthHandler(r, authMiddlewareHandler, userRepository)
	controller.NewUsersHandler(r, authMiddlewareHandler, deviceRepository, userRepository)
	controller.NewDevicesHandler(r, authMiddlewareHandler, deviceRepository, userRepository)

	port := os.Getenv("PORT")
	log.Fatal(r.Run(":" + port))
}

func connectDatabase() *gorm.DB {
	if os.Getenv("DBTYPE") == "mysql" {
		username := os.Getenv("DBUSER")
		password := os.Getenv("DBPASS")
		host := os.Getenv("DBHOST")
		port := os.Getenv("DBPORT")
		dbname := os.Getenv("DBNAME")
		dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}
		return db
	} else if os.Getenv("DBTYPE") == "sqlite" {
		db, err := gorm.Open(sqlite.Open("db/gorm.db"), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}
		return db
	} else {
		log.Fatal("Database type not supported")
		return nil
	}
}
