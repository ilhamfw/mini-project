package main

import (
	"log"
	"rental-games/config"
	"rental-games/docs"
	"rental-games/handler"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger" // This is important
	
)

// @title Rental Games API
// @version 1.0
// @description This is the API documentation for the Rental Games project.
// @host localhost:8080
// @BasePath /api
func main() {
	// Inisialisasi koneksi database
	db, err := config.GetGormDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	e := echo.New()
	// Inisialisasi dokumen Swagger
	docs.SwaggerInfo.Title = "Rental Playstatioms"
	docs.SwaggerInfo.Description = "API for Rent Playstations "
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"

	// Endpoint untuk mendaftarkan pengguna
	e.POST("/users/register", handler.RegisterUser)

	// Endpoint untuk Login pengguna
	e.POST("/users/login", handler.LoginUser)

	// Routing
	e.GET("/console", handler.GetAvailableConsoles(db))
	e.POST("/rent", handler.RentConsole(db), handler.AuthMiddleware)

	// Handle untuk menampilkan dokumentasi Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Mulai server
	e.Start(":8080")
}
