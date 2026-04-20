package handlers

import (
	"net/http"
	"time"

	"med-predict-backend/internal/db"
	"med-predict-backend/internal/middleware"
	"med-predict-backend/internal/models"
	"med-predict-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	db        *db.Database
	analytics *services.AnalyticsService
	log       *services.Logger
}

func NewAnalyticsHandler(database *db.Database, analytics *services.AnalyticsService, log *services.Logger) *AnalyticsHandler {
	return &AnalyticsHandler{
		db:        database,
		analytics: analytics,
		log:       log,
	}
}

// GetTrends returns analytics trends
func (h *AnalyticsHandler) GetTrends(c *gin.Context) {
	claims, _ := middleware.GetUserFromContext(c)

	// Default to last 30 days
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	if start := c.Query("start_date"); start != "" {
		if t, err := time.Parse("2006-01-02", start); err == nil {
			startDate = t
		}
	}

	if end := c.Query("end_date"); end != "" {
		if t, err := time.Parse("2006-01-02", end); err == nil {
			endDate = t
		}
	}

	trends, err := h.analytics.GetTrends(claims.PharmacyID, startDate, endDate)
	if err != nil {
		h.log.Error("failed to compute trends", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to compute trends"})
		return
	}

	c.JSON(http.StatusOK, trends)
}

// GetStockoutRisk returns predicted stockout risks
func (h *AnalyticsHandler) GetStockoutRisk(c *gin.Context) {
	claims, _ := middleware.GetUserFromContext(c)

	risks, err := h.analytics.PredictStockoutRisk(claims.PharmacyID)
	if err != nil {
		h.log.Error("failed to predict stockout risk", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to compute risks"})
		return
	}

	c.JSON(http.StatusOK, risks)
}

// GetAISummary returns AI-generated briefing
func (h *AnalyticsHandler) GetAISummary(c *gin.Context) {
	claims, _ := middleware.GetUserFromContext(c)

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	trends, err := h.analytics.GetTrends(claims.PharmacyID, startDate, endDate)
	if err != nil {
		h.log.Error("failed to compute trends", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate summary"})
		return
	}

	pharmacy, err := h.db.GetPharmacy(claims.PharmacyID)
	if err != nil {
		pharmacy = &models.Pharmacy{Name: "Your Pharmacy"}
	}

	summary := h.analytics.GenerateAISummary(trends, pharmacy.Name, "Last 30 Days")

	c.JSON(http.StatusOK, summary)
}

// GetRegionalOverview returns DHO cross-pharmacy overview
func (h *AnalyticsHandler) GetRegionalOverview(c *gin.Context) {
	// This is DHO-only endpoint
	claims, _ := middleware.GetUserFromContext(c)

	if claims.Role != models.RoleDHO {
		c.JSON(http.StatusForbidden, gin.H{"error": "only DHO can access regional data"})
		return
	}

	// Get all pharmacies and their stats
	pharmacies, err := h.db.GetAllPharmacies()
	if err != nil {
		h.log.Error("failed to fetch pharmacies", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch data"})
		return
	}

	var overview []gin.H
	for _, p := range pharmacies {
		// Get trends for each pharmacy
		endDate := time.Now()
		startDate := endDate.AddDate(0, 0, -30)

		trends, _ := h.analytics.GetTrends(p.ID, startDate, endDate)
		risks, _ := h.analytics.PredictStockoutRisk(p.ID)

		overview = append(overview, gin.H{
			"pharmacy_id":    p.ID,
			"pharmacy_name":  p.Name,
			"region":         p.Region,
			"district":       p.District,
			"lat":            p.Lat,
			"lng":            p.Lng,
			"trends":         trends,
			"stockout_risks": risks,
		})
	}

	c.JSON(http.StatusOK, overview)
}
