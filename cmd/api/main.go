package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/khai1301/moneywise-backend/internal/config"
	"github.com/khai1301/moneywise-backend/internal/models"
	"github.com/khai1301/moneywise-backend/internal/routes"
)

func main() {
	// 1. Load configuration variables
	config.LoadConfig()

	// 2. Connect to Database (Must have PostgreSQL running)
	config.ConnectDatabase()

	// 3. Migrate database tables based on Models struct
	err := config.DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Transaction{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	fmt.Println("Database Migration completed successfully")

	// 4. Initialize Gin router
	router := gin.Default()

	// Configure routing
	routes.SetupRoutes(router)

	// 5. Start Server
	router.Run(":8080")
}
