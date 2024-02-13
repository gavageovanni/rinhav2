package main

import (
	"baitadores-rinhav2/routes"
	"github.com/labstack/echo/v4/middleware"
	"log"
)

func main() {
	//config.DatabaseInit()
	//gorm := config.DB()

	//dbGorm, err := gorm.DB()
	//if err != nil {
	//	panic(err)
	//}
	//
	//dbGorm.Ping()

	echo := routes.Routing.GetRoutes(routes.Routing{})

	echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowMethods:     []string{"GET", "HEAD", "PUT", "PATCH", "POST", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}))

	err := echo.Start(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
