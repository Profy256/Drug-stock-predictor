package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"med-predict/go-backend/internal/middleware"
	"med-predict/go-backend/internal/models"
)

type AnalyticsHandler struct {
	db *sql.DB
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(db *sql.DB) *AnalyticsHandler {
	return &AnalyticsHandler{db: db}
}

// GetPredictions gets medicine predictions
func (h *AnalyticsHandler) GetPredictions(c *gin.Context) {
	_, pharmacyID, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "User not authenticated"})
		return
	}

	rows, err := h.db.Query(`
		SELECT m.id, m.name, m.quantity_on_hand
		FROM medicines m
		WHERE m.pharmacy_id = $1 AND m.is_active = true
	`, pharmacyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to fetch predictions"})
		return
	}
	defer rows.Close()

	predictions := []models.PredictionResponse{}
	for rows.Next() {
		var id, name string
		var current int
		if err := rows.Scan(&id, &name, &current); err != nil {
			continue
		}

		// Simple prediction logic (can be enhanced with ML)
		predicted := float64(current) * 1.1
		trend := "stable"
		if predicted > float64(current)*1.2 {
			trend = "increasing"
		} else if predicted < float64(current)*0.8 {
			trend = "decreasing"
		}

		predictions = append(predictions, models.PredictionResponse{
			MedicineID:   id,
			MedicineName: name,
			Current:      current,
			Predicted:    predicted,
			Trend:        trend,
			Confidence:   0.85,
		})
	}

	c.JSON(http.StatusOK, predictions)
}

// GetAlerts gets alerts for the pharmacy
func (h *AnalyticsHandler) GetAlerts(c *gin.Context) {
	_, pharmacyID, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "User not authenticated"})
		return
	}

	rows, err := h.db.Query(`
		SELECT m.id, m.pharmacy_id, m.name, m.quantity_on_hand, m.reorder_level
		FROM medicines m
		WHERE m.pharmacy_id = $1 AND m.is_active = true
		AND m.quantity_on_hand <= m.reorder_level
	`, pharmacyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to fetch alerts"})
		return
	}
	defer rows.Close()

	var alerts []gin.H
	for rows.Next() {
		var medicineID, medicineName string
		var pharmID string
		var quantity, reorderLevel int
		if err := rows.Scan(&medicineID, &pharmID, &medicineName, &quantity, &reorderLevel); err != nil {
			continue
		}

		alerts = append(alerts, gin.H{
			"id":      medicineID,
			"type":    "low_stock",
			"message": medicineName + " is running low on stock",
			"data": gin.H{
				"current":       quantity,
				"reorder_level": reorderLevel,
			},
		})
	}

	c.JSON(http.StatusOK, alerts)
}

// GetTrends gets analytics trends
func (h *AnalyticsHandler) GetTrends(c *gin.Context) {
	_, pharmacyID, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "User not authenticated"})
		return
	}

	// Placeholder for trends endpoint
	c.JSON(http.StatusOK, gin.H{
		"message":     "Trends data",
		"pharmacy_id": pharmacyID,
	})
}
