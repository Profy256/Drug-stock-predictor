package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"med-predict/go-backend/internal/middleware"
)

// RegisterHandlers registers all route handlers
func RegisterHandlers(router *gin.Engine, db *sql.DB) {
	// Create handlers
	authHandler := NewAuthHandler(db)
	stockHandler := NewStockHandler(db)
	batchHandler := NewBatchHandler(db)
	analyticsHandler := NewAnalyticsHandler(db)
	patientHandler := NewPatientHandler(db)
	recordsHandler := NewRecordsHandler(db)
	adminHandler := NewAdminHandler(db)
	dhoHandler := NewDHOHandler(db)

	// ============================================================
	// Auth Routes (Public)
	// ============================================================
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/me", middleware.AuthMiddleware(), authHandler.GetCurrentUser)
	}

	// ============================================================
	// Stock Routes (Authenticated)
	// ============================================================
	stock := router.Group("/api/v1/stock")
	stock.Use(middleware.AuthMiddleware())
	{
		stock.GET("/medicines", stockHandler.ListMedicines)
		stock.POST("/medicines", stockHandler.CreateMedicine)
		stock.GET("/medicines/:id", stockHandler.GetMedicine)
		stock.PUT("/medicines/:id", stockHandler.UpdateMedicine)
		stock.DELETE("/medicines/:id", stockHandler.DeleteMedicine)
	}

	// ============================================================
	// Batch Routes (Authenticated)
	// ============================================================
	batch := router.Group("/api/v1/batch")
	batch.Use(middleware.AuthMiddleware())
	{
		batch.GET("", batchHandler.ListBatches)
		batch.POST("", batchHandler.CreateBatch)
		batch.POST("/:id/approve", middleware.RequireRole("admin"), batchHandler.ApproveBatch)
		batch.POST("/:id/reject", middleware.RequireRole("admin"), batchHandler.RejectBatch)
	}

	// ============================================================
	// Analytics Routes (Authenticated)
	// ============================================================
	analytics := router.Group("/api/v1/analytics")
	analytics.Use(middleware.AuthMiddleware())
	{
		analytics.GET("/predictions", analyticsHandler.GetPredictions)
		analytics.GET("/alerts", analyticsHandler.GetAlerts)
		analytics.GET("/trends", analyticsHandler.GetTrends)
	}

	// ============================================================
	// Patient Routes (Authenticated)
	// ============================================================
	patient := router.Group("/api/v1/patient")
	patient.Use(middleware.AuthMiddleware())
	{
		patient.GET("/form-fields", patientHandler.GetFormFields)
		patient.PUT("/form-fields", middleware.RequireRole("admin"), patientHandler.UpdateFormFields)
		patient.POST("/pending-records", patientHandler.CreatePendingRecord)
		patient.GET("/pending-records", patientHandler.ListPendingRecords)
		patient.POST("/pending-records/:id/approve", middleware.RequireRole("admin"), patientHandler.ApprovePendingRecord)
	}

	// ============================================================
	// Records Routes (Authenticated)
	// ============================================================
	records := router.Group("/api/v1/records")
	records.Use(middleware.AuthMiddleware())
	{
		records.GET("/pending", recordsHandler.ListPendingRecords)
		records.GET("/approved", recordsHandler.ListApprovedVisits)
		records.GET("/:id", recordsHandler.GetRecordDetails)
	}

	// ============================================================
	// Admin Routes (Admin Only)
	// ============================================================
	admin := router.Group("/api/v1/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.RequireRole("admin"))
	{
		admin.GET("/users", adminHandler.ListUsers)
		admin.DELETE("/users/:id", adminHandler.DeactivateUser)
		admin.GET("/audit-logs", adminHandler.GetAuditLogs)
		admin.GET("/pharmacies", adminHandler.ListPharmacies)
	}

	// ============================================================
	// DHO Routes (DHO Only)
	// ============================================================
	dho := router.Group("/api/v1/dho")
	dho.Use(middleware.AuthMiddleware(), middleware.RequireRole("dho"))
	{
		dho.GET("/batches/:id/review", dhoHandler.ReviewBatch)
		dho.GET("/batches-for-review", dhoHandler.ListBatchesForReview)
		dho.GET("/stats", dhoHandler.GetReviewStats)
	}
}
