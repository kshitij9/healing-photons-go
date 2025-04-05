package main

import (
	"healing_photons/internal/config"
	"healing_photons/internal/database"
	"healing_photons/internal/handlers"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := database.InitializeDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Setup Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Initialize routes
	handlers.SetupRoutes(router, db)
	handlers.SetupWeightTypeRoutes(router, db)
	handlers.SetupPeelingMachineRoutes(router, db)
	handlers.SetupHumidifierRoutes(router, db)

	// Start server
	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
