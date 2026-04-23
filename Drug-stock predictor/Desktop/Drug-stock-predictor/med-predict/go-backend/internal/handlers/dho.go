package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"med-predict/go-backend/internal/middleware"
	"med-predict/go-backend/internal/models"
)

type DHOHandler struct {
	db *sql.DB
}

// NewDHOHandler creates a new DHO handler
func NewDHOHandler(db *sql.DB) *DHOHandler {
	return &DHOHandler{db: db}
}

// ReviewBatch reviews a batch (DHO only)
func (h *DHOHandler) ReviewBatch(c *gin.Context) {
	_, _, role, err := middleware.GetUserFromContext(c)
	if err != nil || role != models.RoleDHO {
		c.JSON(http.StatusForbidden, gin.H{"detail": "Only DHOs can review batches"})
		return
	}

	batchID := c.Param("id")

	var batch models.BatchResponse
	query := `
		SELECT id, pharmacy_id, name, description, status, submitted_by_id, approved_by_id, rejection_reason, created_at, updated_at
		FROM batches
		WHERE id = $1
	`

	err = h.db.QueryRow(query, batchID).Scan(
		&batch.ID, &batch.PharmacyID, &batch.Name, &batch.Description, &batch.Status,
		&batch.SubmittedByID, &batch.ApprovedByID, &batch.RejectionReason, &batch.CreatedAt, &batch.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Batch not found"})
		return
	}

	c.JSON(http.StatusOK, batch)
}

// ListBatchesForReview lists batches pending DHO review
func (h *DHOHandler) ListBatchesForReview(c *gin.Context) {
	_, _, role, err := middleware.GetUserFromContext(c)
	if err != nil || role != models.RoleDHO {
		c.JSON(http.StatusForbidden, gin.H{"detail": "Only DHOs can review batches"})
		return
	}

	rows, err := h.db.Query(`
		SELECT id, pharmacy_id, name, description, status, submitted_by_id, approved_by_id, rejection_reason, created_at, updated_at
		FROM batches
		WHERE status = 'pending'
		ORDER BY created_at DESC
	`)
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

// GetReviewStats gets review statistics
func (h *DHOHandler) GetReviewStats(c *gin.Context) {
	_, _, role, err := middleware.GetUserFromContext(c)
	if err != nil || role != models.RoleDHO {
		c.JSON(http.StatusForbidden, gin.H{"detail": "Only DHOs can view stats"})
		return
	}

	var pendingCount, approvedCount, rejectedCount int

	h.db.QueryRow("SELECT COUNT(*) FROM batches WHERE status = 'pending'").Scan(&pendingCount)
	h.db.QueryRow("SELECT COUNT(*) FROM batches WHERE status = 'approved'").Scan(&approvedCount)
	h.db.QueryRow("SELECT COUNT(*) FROM batches WHERE status = 'rejected'").Scan(&rejectedCount)

	c.JSON(http.StatusOK, gin.H{
		"pending":  pendingCount,
		"approved": approvedCount,
		"rejected": rejectedCount,
	})
}
