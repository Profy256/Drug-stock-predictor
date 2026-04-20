package handlers

import (
	"net/http"
	"strconv"

	"med-predict-backend/internal/db"
	"med-predict-backend/internal/middleware"
	"med-predict-backend/internal/models"
	"med-predict-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminHandler struct {
	db    *db.Database
	audit *services.AuditService
	log   *services.Logger
}

func NewAdminHandler(database *db.Database, audit *services.AuditService, log *services.Logger) *AdminHandler {
	return &AdminHandler{
		db:    database,
		audit: audit,
		log:   log,
	}
}

// ============================================================
// Form Fields Management
// ============================================================

// GetFormFields returns custom patient form fields
func (h *AdminHandler) GetFormFields(c *gin.Context) {
	claims, _ := middleware.GetUserFromContext(c)

	fields, err := h.db.GetPharmacyFormFields(claims.PharmacyID)
	if err != nil {
		h.log.Error("failed to fetch form fields", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch fields"})
		return
	}

	c.JSON(http.StatusOK, fields)
}

// CreateFormField creates a new custom form field
func (h *AdminHandler) CreateFormField(c *gin.Context) {
	claims, _ := middleware.GetUserFromContext(c)

	var req struct {
		FieldKey   string   `json:"field_key" binding:"required"`
		Label      string   `json:"label" binding:"required"`
		FieldType  string   `json:"field_type" binding:"required"`
		Options    []string `json:"options"`
		IsRequired bool     `json:"is_required"`
		SortOrder  int      `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	field := &models.PatientFormField{
		ID:         uuid.New().String(),
		PharmacyID: claims.PharmacyID,
		FieldKey:   req.FieldKey,
		Label:      req.Label,
		FieldType:  req.FieldType,
		Options:    req.Options,
		IsRequired: req.IsRequired,
		IsActive:   true,
		SortOrder:  req.SortOrder,
	}

	if err := h.db.CreateFormField(field); err != nil {
		h.log.Error("failed to create form field", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "creation failed"})
		return
	}

	h.audit.LogAction(claims.UserID, claims.PharmacyID, "FORM_FIELD_CREATED", "FORM_FIELD", field.ID, c.ClientIP(), nil)

	c.JSON(http.StatusCreated, field)
}

// ============================================================
// User Management
// ============================================================

// GetUsers returns all users in pharmacy
func (h *AdminHandler) GetUsers(c *gin.Context) {
	claims, _ := middleware.GetUserFromContext(c)

	users, err := h.db.GetPharmacyUsers(claims.PharmacyID)
	if err != nil {
		h.log.Error("failed to fetch users", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// ============================================================
// Audit Logs
// ============================================================

// GetAuditLogs returns audit trail with pagination
func (h *AdminHandler) GetAuditLogs(c *gin.Context) {
	claims, _ := middleware.GetUserFromContext(c)

	limit := 50
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	logs, err := h.db.GetAuditLogs(claims.PharmacyID, limit, offset)
	if err != nil {
		h.log.Error("failed to fetch audit logs", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}
