package services

import (
	"encoding/json"

	"med-predict-backend/internal/db"
	"med-predict-backend/internal/models"

	"github.com/google/uuid"
)

// AuditService handles audit logging
type AuditService struct {
	db     *db.Database
	logger *Logger
}

// NewAuditService creates an audit service instance
func NewAuditService(database *db.Database, logger *Logger) *AuditService {
	return &AuditService{
		db:     database,
		logger: logger,
	}
}

// LogAction logs a user action for compliance
func (as *AuditService) LogAction(userID, pharmacyID, action, entityType, entityID, ipAddress string, details interface{}) {
	log := &models.AuditLog{
		ID:         uuid.New().String(),
		UserID:     userID,
		PharmacyID: pharmacyID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		IPAddress:  ipAddress,
	}

	// Serialize details to JSON
	if details != nil {
		if data, err := json.Marshal(details); err == nil {
			log.Details = string(data)
		}
	}

	// Log to database - never fail the main operation
	if err := as.db.LogAuditEvent(log); err != nil {
		as.logger.Error("failed to write audit log", "error", err.Error(), "action", action)
	}
}

// Audit action constants
const (
	ActionLoginSuccess    = "LOGIN_SUCCESS"
	ActionLoginFailed     = "LOGIN_FAILED"
	ActionUserCreated     = "USER_CREATED"
	ActionUserUpdated     = "USER_UPDATED"
	ActionStockAdded      = "STOCK_ADDED"
	ActionStockAdjusted   = "STOCK_ADJUSTED"
	ActionBatchSubmitted  = "BATCH_SUBMITTED"
	ActionBatchApproved   = "BATCH_APPROVED"
	ActionBatchRejected   = "BATCH_REJECTED"
	ActionPharmacyCreated = "PHARMACY_CREATED"
)

// Common entity types
const (
	EntityTypeUser     = "USER"
	EntityTypePharmacy = "PHARMACY"
	EntityTypeMedicine = "MEDICINE"
	EntityTypeBatch    = "BATCH"
)
