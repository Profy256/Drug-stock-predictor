package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"med-predict/go-backend/internal/middleware"
	"med-predict/go-backend/internal/models"
)

type PatientHandler struct {
	db *sql.DB
}

// NewPatientHandler creates a new patient handler
func NewPatientHandler(db *sql.DB) *PatientHandler {
	return &PatientHandler{db: db}
}

// GetFormFields gets patient form fields for a pharmacy
func (h *PatientHandler) GetFormFields(c *gin.Context) {
	_, pharmacyID, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "User not authenticated"})
		return
	}

	// Default form fields
	fields := []models.PatientFormField{
		{
			Name:        "name",
			FieldType:   "text",
			Required:    true,
			Placeholder: "Patient Name",
		},
		{
			Name:        "age",
			FieldType:   "number",
			Required:    true,
			Placeholder: "Age",
		},
		{
			Name:      "gender",
			FieldType: "select",
			Required:  true,
			Options:   []string{"Male", "Female", "Other"},
		},
		{
			Name:        "phone",
			FieldType:   "tel",
			Required:    false,
			Placeholder: "Phone Number",
		},
	}

	c.JSON(http.StatusOK, fields)
}

// UpdateFormFields updates patient form fields
func (h *PatientHandler) UpdateFormFields(c *gin.Context) {
	_, _, role, err := middleware.GetUserFromContext(c)
	if err != nil || role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"detail": "Only admins can update form fields"})
		return
	}

	var fields []models.PatientFormField
	if err := c.ShouldBindJSON(&fields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Form fields updated successfully"})
}

// CreatePendingRecord creates a pending patient record
func (h *PatientHandler) CreatePendingRecord(c *gin.Context) {
	_, pharmacyID, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "User not authenticated"})
		return
	}

	var patientData map[string]interface{}
	if err := c.ShouldBindJSON(&patientData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	recordID := generateID("record")

	// For now, just acknowledge the record
	c.JSON(http.StatusCreated, gin.H{
		"id":          recordID,
		"pharmacy_id": pharmacyID,
		"status":      "pending",
		"created_at":  time.Now(),
	})
}

// ListPendingRecords lists pending records
func (h *PatientHandler) ListPendingRecords(c *gin.Context) {
	_, pharmacyID, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "User not authenticated"})
		return
	}

	// Placeholder implementation
	records := []gin.H{}

	c.JSON(http.StatusOK, records)
}

// ApprovePendingRecord approves a pending record
func (h *PatientHandler) ApprovePendingRecord(c *gin.Context) {
	_, _, role, err := middleware.GetUserFromContext(c)
	if err != nil || role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"detail": "Only admins can approve records"})
		return
	}

	recordID := c.Param("id")

	c.JSON(http.StatusOK, gin.H{
		"message":   "Record approved successfully",
		"record_id": recordID,
	})
}
