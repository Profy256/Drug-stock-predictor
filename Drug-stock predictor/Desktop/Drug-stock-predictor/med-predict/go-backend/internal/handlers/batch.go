package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"med-predict/go-backend/internal/middleware"
	"med-predict/go-backend/internal/models"
)

type BatchHandler struct {
	db *sql.DB
}

// NewBatchHandler creates a new batch handler
func NewBatchHandler(db *sql.DB) *BatchHandler {
	return &BatchHandler{db: db}
}

// ListBatches lists all batches
func (h *BatchHandler) ListBatches(c *gin.Context) {
	_, pharmacyID, role, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "User not authenticated"})
		return
	}

	var rows *sql.Rows
	if role == models.RoleAdmin {
		// Admins see all batches
		rows, err = h.db.Query(`
			SELECT id, pharmacy_id, name, description, status, submitted_by_id, approved_by_id, rejection_reason, created_at, updated_at
			FROM batches
			ORDER BY created_at DESC
		`)
	} else {
		// Others see only their pharmacy's batches
		rows, err = h.db.Query(`
			SELECT id, pharmacy_id, name, description, status, submitted_by_id, approved_by_id, rejection_reason, created_at, updated_at
			FROM batches
			WHERE pharmacy_id = $1
			ORDER BY created_at DESC
		`, pharmacyID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to fetch batches"})
		return
	}
	defer rows.Close()

	batches := []models.BatchResponse{}
	for rows.Next() {
		var batch models.BatchResponse
		err := rows.Scan(
			&batch.ID, &batch.PharmacyID, &batch.Name, &batch.Description, &batch.Status,
			&batch.SubmittedByID, &batch.ApprovedByID, &batch.RejectionReason, &batch.CreatedAt, &batch.UpdatedAt,
		)
		if err != nil {
			continue
		}
		batches = append(batches, batch)
	}

	c.JSON(http.StatusOK, batches)
}

// CreateBatch creates a new batch
func (h *BatchHandler) CreateBatch(c *gin.Context) {
	userID, pharmacyID, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "User not authenticated"})
		return
	}

	var req models.BatchCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	batchID := generateID("batch")
	now := time.Now()

	query := `
		INSERT INTO batches (id, pharmacy_id, name, description, status, submitted_by_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	status := "pending"
	if req.Status != "" {
		status = string(req.Status)
	}

	_, err = h.db.Exec(query, batchID, pharmacyID, req.Name, req.Description, status, userID, now, now)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to create batch"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": batchID, "message": "Batch created successfully"})
}

// ApproveBatch approves a batch (admin only)
func (h *BatchHandler) ApproveBatch(c *gin.Context) {
	userID, _, role, err := middleware.GetUserFromContext(c)
	if err != nil || role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"detail": "Only admins can approve batches"})
		return
	}

	batchID := c.Param("id")

	query := `
		UPDATE batches
		SET status = $1, approved_by_id = $2, updated_at = $3
		WHERE id = $4
	`

	_, err = h.db.Exec(query, "approved", userID, time.Now(), batchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to approve batch"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Batch approved successfully"})
}

// RejectBatch rejects a batch (admin only)
func (h *BatchHandler) RejectBatch(c *gin.Context) {
	userID, _, role, err := middleware.GetUserFromContext(c)
	if err != nil || role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"detail": "Only admins can reject batches"})
		return
	}

	batchID := c.Param("id")

	var req struct {
		RejectionReason string `json:"rejection_reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	query := `
		UPDATE batches
		SET status = $1, approved_by_id = $2, rejection_reason = $3, updated_at = $4
		WHERE id = $5
	`

	_, err = h.db.Exec(query, "rejected", userID, req.RejectionReason, time.Now(), batchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to reject batch"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Batch rejected successfully"})
}
