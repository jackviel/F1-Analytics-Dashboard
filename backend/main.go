package main

import (
	"log"
	"os"

	"github.com/f1-analytics/config"
	"github.com/f1-analytics/handlers"
	"github.com/f1-analytics/middleware"
	"github.com/f1-analytics/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize database
	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Run database migrations
	if err := config.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Get database connection
	db := config.GetDB()

	// Initialize services
	openF1Service := services.NewOpenF1Service()

	// Initialize handlers with both OpenF1Service and database connection
	driverHandler := handlers.NewDriverHandler(openF1Service, db)
	teamHandler := handlers.NewTeamHandler(openF1Service, db)
	raceHandler := handlers.NewRaceHandler(openF1Service, db)

	// Initialize router
	router := gin.Default()

	// Add CORS middleware
	router.Use(middleware.CORSMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Driver routes
		api.GET("/drivers", driverHandler.GetDrivers)
		api.GET("/drivers/:id", driverHandler.GetDriver)
		api.GET("/drivers/:id/stats", driverHandler.GetDriverStats)

		// Team routes
		api.GET("/teams", teamHandler.GetTeams)
		api.GET("/teams/:id", teamHandler.GetTeam)
		api.GET("/teams/:id/stats", teamHandler.GetTeamStats)

		// Race routes
		api.GET("/races", raceHandler.GetRaces)
		api.GET("/races/:id", raceHandler.GetRace)
		api.GET("/races/:id/results", raceHandler.GetRaceResults)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
