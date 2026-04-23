package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"med-predict/go-backend/internal/db"
	"med-predict/go-backend/internal/handlers"
	"med-predict/go-backend/internal/middleware"
)

func main() {
	// Load environment variables
	err := godotenv.Load("../../.env")
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: Error loading .env file: %v\n", err)
	}

	// Initialize database
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v\n", err)
	}
	defer database.Close()

	// Set up Gin router
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	// Setup middleware
	middleware.SetupMiddleware(router)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "Med Predict Backend",
		})
	})

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Med Predict Backend",
			"version": "1.0.0",
			"docs":    "/docs",
		})
	})

	// Register all handlers
	handlers.RegisterHandlers(router, database)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Starting server on port %s...\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
