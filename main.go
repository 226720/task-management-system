package main

import (
	"task/config"
	"task/models"
	"task/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database
	config.ConnectDatabase()

	// Migrate the schema
	config.DB.AutoMigrate(&models.User{}, &models.Task{})

	// Initialize Gin router
	router := gin.Default()

	// Set up routes
	routes.SetupRoutes(router)

	// Run the server
	router.Run(":8080")
}
