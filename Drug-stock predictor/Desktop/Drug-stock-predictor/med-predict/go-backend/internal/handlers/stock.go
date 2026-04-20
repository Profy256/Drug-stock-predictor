package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"med-predict-backend/internal/db"
	"med-predict-backend/internal/middleware"
	"med-predict-backend/internal/models"
	"med-predict-backend/internal/services"
)

type StockHandler struct {
	db    *db.Database
	audit *services.AuditService
	log   *services.Logger
}

func NewStockHandler(database *db.Database, audit *services.AuditService, log *services.Logger) *StockHandler {
	return &StockHandler{
		db:    database,
		audit: audit,
		log:   log,
	}
}

// List returns all medicines for a pharmacy
func (h *StockHandler) List(c *gin.Context) {
	claims, _ := middleware.GetUserFromContext(c)

	medicines, err := h.db.GetPharmacyMedicines(claims.PharmacyID)
	if err != nil {
		h.log.Error("failed to fetch medicines", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch medicines"})
		return
	}

	// Add status to each medicine
	for i := range medicines {
		medicines[i].Status = calculateMedicineStatus(&medicines[i])
	}

	c.JSON(http.StatusOK, medicines)
}

// Add creates new stock entry
func (h *StockHandler) Add(c *gin.Context) {
	claims, _ := middleware.GetUserFromContext(c)

	var req models.AddStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	med := &models.Medicine{
		ID:               uuid.New().String(),
		PharmacyID:       claims.PharmacyID,
		Name:             req.Name,
		GenericName:      req.GenericName,
		Category:         req.Category,
		Unit:             req.Unit,
		QuantityTotal:    req.QuantityTotal,
		QuantityRemaining: req.QuantityTotal,
		ExpiryDate:       req.ExpiryDate,
		BatchNumber:      req.BatchNumber,
		Supplier:         req.Supplier,
		UnitCost:         req.UnitCost,
		ReorderLevel:     req.ReorderLevel,
		NotificationDays: req.NotificationDays,
		CreatedBy:        claims.UserID,
	}

	if err := h.db.CreateMedicine(med); err != nil {
		h.log.Error("failed to create medicine", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create medicine"})
		return
	}

	h.audit.LogAction(claims.UserID, claims.PharmacyID, services.ActionStockAdded, services.EntityTypeMedicine, med.ID, c.ClientIP(), map[string]interface{}{"name": med.Name, "quantity": req.QuantityTotal})

	c.JSON(http.StatusCreated, med)
}

// Update adjusts medicine quantity
func (h *StockHandler) Update(c *gin.Context) {
	claims, _ := middleware.GetUserFromContext(c)
	medicineID := c.Param("id")

	var req models.AdjustStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	med, err := h.db.GetMedicine(medicineID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "medicine not found"})
		return
	}

	// Verify pharmacy ownership
	if med.PharmacyID != claims.PharmacyID {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.db.UpdateMedicineQuantity(medicineID, req.QuantityAdjustment); err != nil {
		h.log.Error("failed to update medicine", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}

	h.audit.LogAction(claims.UserID, claims.PharmacyID, services.ActionStockAdjusted, services.EntityTypeMedicine, medicineID, c.ClientIP(), map[string]interface{}{"adjustment": req.QuantityAdjustment, "reason": req.Reason})

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// Search returns typeahead results for medicine search
func (h *StockHandler) Search(c *gin.Context) {
	claims, _ := middleware.GetUserFromContext(c)
	query := c.Query("q")

	if len(query) < 2 {
		c.JSON(http.StatusOK, []models.Medicine{})
		return
	}

	medicines, err := h.db.SearchMedicines(claims.PharmacyID, query)
	if err != nil {
		h.log.Error("search failed", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}

	c.JSON(http.StatusOK, medicines)
}

// Expiring returns medicines expiring soon
func (h *StockHandler) Expiring(c *gin.Context) {
	claims, _ := middleware.GetUserFromContext(c)

	medicines, err := h.db.GetExpiringMedicines(claims.PharmacyID, 14)
	if err != nil {
		h.log.Error("failed to fetch expiring medicines", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch medicines"})
		return
	}

	c.JSON(http.StatusOK, medicines)
}

// Helper functions
func calculateMedicineStatus(med *models.Medicine) string {
	now := time.Now()
	
	if med.ExpiryDate.Before(now) {
		return models.StatusExpired
	}
	if med.QuantityRemaining <= med.ReorderLevel {
		return models.StatusLowStock
	}
	
	daysRemaining := med.ExpiryDate.Sub(now).Hours() / 24
	if daysRemaining <= 7 {
		return models.StatusExpiring
	}
	
	return models.StatusOK
}
