package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"med-predict-backend/internal/config"
	"med-predict-backend/internal/db"
	"med-predict-backend/internal/handlers"
	"med-predict-backend/internal/middleware"
	"med-predict-backend/internal/models"
	"med-predict-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Initialize logger
	logger := services.NewLogger(cfg)

	// Connect to database
	database, err := db.Connect(cfg)
	if err != nil {
		logger.Error("database connection failed", "error", err.Error())
		os.Exit(1)
	}
	defer database.Close()

	logger.Info("database connected successfully")

	// Initialize services
	auditService := services.NewAuditService(database, logger)
	analyticsService := services.NewAnalyticsService(database, logger)

	// Initialize Gin
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(middleware.LoggerMiddleware(logger.Logger))
	router.Use(middleware.CORSMiddleware(cfg))
	router.Use(gin.Recovery())

	// Health check endpoint (no auth required)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
		})
	})

	// ============================================================
	// Initialize Handlers
	// ============================================================

	authHandler := handlers.NewAuthHandler(database, auditService, logger, cfg.JWTSecret)
	stockHandler := handlers.NewStockHandler(database, auditService, logger)
	batchHandler := handlers.NewBatchHandler(database, auditService, logger)
	analyticsHandler := handlers.NewAnalyticsHandler(database, analyticsService, logger)
	patientHandler := handlers.NewPatientHandler(database, logger)
	adminHandler := handlers.NewAdminHandler(database, auditService, logger)
	dhoHandler := handlers.NewDHOHandler(database, auditService, logger, cfg.JWTSecret)

	// ============================================================
	// API Routes
	// ============================================================

	// Authentication routes (rate limited, no auth required)
	authGroup := router.Group("/api/v1/auth")
	authGroup.Use(middleware.RateLimitAuth())
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/register", authHandler.Register)
	}

	// Authenticated routes
	apiGroup := router.Group("/api/v1")
	apiGroup.Use(middleware.AuthMiddleware(cfg, logger))
	apiGroup.Use(middleware.RateLimitAPI())
	{
		// Auth
		authGroup := apiGroup.Group("/auth")
		{
			authGroup.GET("/me", authHandler.GetMe)
		}

		// Stock Management
		stockGroup := apiGroup.Group("/stock")
		{
			stockGroup.GET("", stockHandler.List)
			stockGroup.POST("", stockHandler.Add)
			stockGroup.PUT("/:id", stockHandler.Update)
			stockGroup.GET("/expiring", stockHandler.Expiring)
			stockGroup.GET("/search", stockHandler.Search)
		}

		// Batch Processing
		batchGroup := apiGroup.Group("/batches")
		{
			batchGroup.POST("", batchHandler.Submit)
			batchGroup.GET("", batchHandler.List)
			batchGroup.GET("/:id", batchHandler.Get)
			batchGroup.POST("/:id/approve", middleware.RequireRole(models.RoleAdmin, models.RoleDHO)(batchHandler.Approve))
			batchGroup.POST("/:id/reject", middleware.RequireRole(models.RoleAdmin, models.RoleDHO)(batchHandler.Reject))
			batchGroup.DELETE("/:id/records/:rid", batchHandler.DeleteRecord)
		}

		// Analytics
		analyticsGroup := apiGroup.Group("/analytics")
		{
			analyticsGroup.GET("/trends", analyticsHandler.GetTrends)
			analyticsGroup.GET("/ai-summary", analyticsHandler.GetAISummary)
			analyticsGroup.GET("/stockout-risk", analyticsHandler.GetStockoutRisk)
			analyticsGroup.GET("/regional", analyticsHandler.GetRegionalOverview)
		}

		// Patient Data
		patientGroup := apiGroup.Group("/patients")
		{
			patientGroup.GET("/search", patientHandler.Search)
		}

		// Administration
		adminGroup := apiGroup.Group("/admin")
		adminGroup.Use(middleware.RequireRole(models.RoleAdmin, models.RoleDHO))
		{
			adminGroup.GET("/form-fields", adminHandler.GetFormFields)
			adminGroup.POST("/form-fields", adminHandler.CreateFormField)

			adminGroup.GET("/users", adminHandler.GetUsers)

			adminGroup.GET("/audit-logs", adminHandler.GetAuditLogs)
		}

		// DHO (District Health Officer)
		dhoGroup := apiGroup.Group("/dho")
		dhoGroup.Use(middleware.RequireRole(models.RoleDHO))
		{
			dhoGroup.GET("/pharmacies", dhoHandler.GetPharmacies)
			dhoGroup.POST("/pharmacies", dhoHandler.RegisterPharmacy)
			dhoGroup.GET("/regional-map", dhoHandler.RegionalMap)
		}
	}

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Port)
	logger.Info("server starting", "addr", addr, "env", cfg.Env)

	if err := router.Run(addr); err != nil {
		logger.Error("server error", "error", err.Error())
		os.Exit(1)
	}
}
