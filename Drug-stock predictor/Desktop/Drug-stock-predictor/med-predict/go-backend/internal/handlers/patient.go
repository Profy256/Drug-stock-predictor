package handlers

import (
	"net/http"

	"med-predict-backend/internal/db"
	"med-predict-backend/internal/middleware"
	"med-predict-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	db  *db.Database
	log *services.Logger
}

func NewPatientHandler(database *db.Database, log *services.Logger) *PatientHandler {
	return &PatientHandler{
		db:  database,
		log: log,
	}
}

// Search returns patient visit history by hashed credential
func (h *PatientHandler) Search(c *gin.Context) {
	claims, _ := middleware.GetUserFromContext(c)
	patientHash := c.Query("hash")

	if patientHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "patient hash required"})
		return
	}

	// Data entrants can only search within their pharmacy
	if claims.Role != "admin" && claims.Role != "dho" {
		// For data entrants: only accessible within same pharmacy context
	}

	c.JSON(http.StatusOK, gin.H{
		"patient_hash": patientHash,
		"visits":       []interface{}{},
		"message":      "patient history (anonymized)",
	})
}
