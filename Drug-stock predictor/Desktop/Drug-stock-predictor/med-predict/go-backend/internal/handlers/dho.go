package handlers

import (
	"net/http"

	"med-predict-backend/internal/db"
	"med-predict-backend/internal/middleware"
	"med-predict-backend/internal/models"
	"med-predict-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type DHOHandler struct {
	db     *db.Database
	audit  *services.AuditService
	log    *services.Logger
	secret string
}

func NewDHOHandler(database *db.Database, audit *services.AuditService, log *services.Logger, secret string) *DHOHandler {
	return &DHOHandler{
		db:     database,
		audit:  audit,
		log:    log,
		secret: secret,
	}
}

// GetPharmacies returns all pharmacies with stats (DHO only)
func (h *DHOHandler) GetPharmacies(c *gin.Context) {
	claims, _ := middleware.GetUserFromContext(c)

	if claims.Role != models.RoleDHO {
		c.JSON(http.StatusForbidden, gin.H{"error": "only DHO can access"})
		return
	}

	pharmacies, err := h.db.GetAllPharmacies()
	if err != nil {
		h.log.Error("failed to fetch pharmacies", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch pharmacies"})
		return
	}

	var results []gin.H
	for _, p := range pharmacies {
		users, _ := h.db.GetPharmacyUsers(p.ID)
		medicines, _ := h.db.GetPharmacyMedicines(p.ID)

		results = append(results, gin.H{
			"id":             p.ID,
			"name":           p.Name,
			"region":         p.Region,
			"district":       p.District,
			"lat":            p.Lat,
			"lng":            p.Lng,
			"contact_phone":  p.ContactPhone,
			"user_count":     len(users),
			"medicine_count": len(medicines),
			"is_active":      p.IsActive,
			"created_at":     p.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, results)
}

// RegisterPharmacy creates new pharmacy and admin user (DHO only)
func (h *DHOHandler) RegisterPharmacy(c *gin.Context) {
	claims, _ := middleware.GetUserFromContext(c)

	if claims.Role != models.RoleDHO {
		c.JSON(http.StatusForbidden, gin.H{"error": "only DHO can register pharmacies"})
		return
	}

	var req struct {
		PharmacyName  string  `json:"pharmacy_name" binding:"required"`
		Region        string  `json:"region" binding:"required"`
		District      string  `json:"district" binding:"required"`
		Lat           float64 `json:"lat"`
		Lng           float64 `json:"lng"`
		ContactPhone  string  `json:"contact_phone"`
		AdminName     string  `json:"admin_name" binding:"required"`
		AdminEmail    string  `json:"admin_email" binding:"required,email"`
		AdminPassword string  `json:"admin_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Create pharmacy
	pharmacyID := uuid.New().String()
	pharmacy := &models.Pharmacy{
		ID:           pharmacyID,
		Name:         req.PharmacyName,
		Region:       req.Region,
		District:     req.District,
		Lat:          req.Lat,
		Lng:          req.Lng,
		ContactPhone: req.ContactPhone,
		IsActive:     true,
	}

	if err := h.db.CreatePharmacy(pharmacy); err != nil {
		h.log.Error("failed to create pharmacy", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "pharmacy creation failed"})
		return
	}

	// Create admin user for pharmacy
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		h.log.Error("password hashing failed", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "admin creation failed"})
		return
	}

	adminID := uuid.New().String()
	admin := &models.User{
		ID:           adminID,
		PharmacyID:   pharmacyID,
		Name:         req.AdminName,
		Email:        req.AdminEmail,
		PasswordHash: string(hashedPassword),
		Role:         models.RoleAdmin,
		IsActive:     true,
	}

	if err := h.db.CreateUser(admin); err != nil {
		h.log.Error("failed to create admin user", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "admin creation failed"})
		return
	}

	h.audit.LogAction(claims.UserID, "", services.ActionPharmacyCreated, "PHARMACY", pharmacyID, c.ClientIP(), gin.H{
		"pharmacy_name": req.PharmacyName,
		"admin_email":   req.AdminEmail,
	})

	c.JSON(http.StatusCreated, gin.H{
		"pharmacy_id": pharmacyID,
		"admin_id":    adminID,
		"message":     "pharmacy registered successfully",
	})
}

// RegionalMap returns map data for DHO (DHO only)
func (h *DHOHandler) RegionalMap(c *gin.Context) {
	claims, _ := middleware.GetUserFromContext(c)

	if claims.Role != models.RoleDHO {
		c.JSON(http.StatusForbidden, gin.H{"error": "only DHO can access"})
		return
	}

	pharmacies, err := h.db.GetAllPharmacies()
	if err != nil {
		h.log.Error("failed to fetch pharmacies", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch map data"})
		return
	}

	var mapData []gin.H
	for _, p := range pharmacies {
		risks, _ := NewAnalyticsHandler(h.db, services.NewAnalyticsService(h.db, h.log), h.log).analytics.PredictStockoutRisk(p.ID)

		riskLevel := "ok"
		criticalCount := 0
		warningCount := 0
		for _, r := range risks {
			if r.RiskLevel == "critical" {
				criticalCount++
			} else if r.RiskLevel == "warning" {
				warningCount++
			}
		}

		if criticalCount > 0 {
			riskLevel = "critical"
		} else if warningCount > 0 {
			riskLevel = "warning"
		}

		mapData = append(mapData, gin.H{
			"id":         p.ID,
			"name":       p.Name,
			"lat":        p.Lat,
			"lng":        p.Lng,
			"risk_level": riskLevel,
			"region":     p.Region,
		})
	}

	c.JSON(http.StatusOK, mapData)
}
