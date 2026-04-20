package handlers

import (
	"net/http"
	"strconv"
	"time"

	"med-predict-backend/internal/db"
	"med-predict-backend/internal/middleware"
	"med-predict-backend/internal/models"
	"med-predict-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BatchHandler struct {
	db    *db.Database
	audit *services.AuditService
	log   *services.Logger
}

func NewBatchHandler(database *db.Database, audit *services.AuditService, log *services.Logger) *BatchHandler {
	return &BatchHandler{
		db:    database,
		audit: audit,
		log:   log,
	}
}

// Submit creates a new batch with pending records
func (h *BatchHandler) Submit(c *gin.Context) {
	claims, _ := middleware.GetUserFromContext(c)

	var req models.SubmitBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	batchID := uuid.New().String()
	batch := &models.Batch{
		ID:          batchID,
		PharmacyID:  claims.PharmacyID,
		SubmittedBy: claims.UserID,
		Status:      models.BatchStatusPending,
		RecordCount: len(req.Records),
	}

	if err := h.db.CreateBatch(batch); err != nil {
		h.log.Error("failed to create batch", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create batch"})
		return
	}

	// Create pending records
	for _, record := range req.Records {
		pr := &models.PendingRecord{
			ID:                uuid.New().String(),
			BatchID:           batchID,
			PatientHash:       record.PatientHash,
			MedicineID:        record.MedicineID,
			QuantityDispensed: record.QuantityDispensed,
			Diagnosis:         record.Diagnosis,
			PatientData:       record.PatientData,
		}

		if err := h.db.CreatePendingRecord(pr); err != nil {
			h.log.Error("failed to create pending record", "error", err.Error())
		}
	}

	h.audit.LogAction(claims.UserID, claims.PharmacyID, services.ActionBatchSubmitted, services.EntityTypeBatch, batchID, c.ClientIP(), map[string]interface{}{"record_count": len(req.Records)})

	c.JSON(http.StatusCreated, gin.H{"batch_id": batchID})
}

// List returns batches for a pharmacy with pagination
func (h *BatchHandler) List(c *gin.Context) {
	claims, _ := middleware.GetUserFromContext(c)

	limit := 20
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

	batches, err := h.db.GetPharmacyBatches(claims.PharmacyID, limit, offset)
	if err != nil {
		h.log.Error("failed to fetch batches", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch batches"})
		return
	}

	c.JSON(http.StatusOK, batches)
}

// Get returns batch details with pending records
func (h *BatchHandler) Get(c *gin.Context) {
	claims, _ := middleware.GetUserFromContext(c)
	batchID := c.Param("id")

	batch, err := h.db.GetBatch(batchID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "batch not found"})
		return
	}

	if batch.PharmacyID != claims.PharmacyID {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	records, err := h.db.GetBatchRecords(batchID)
	if err != nil {
		h.log.Error("failed to fetch records", "error", err.Error())
		records = []models.PendingRecord{}
	}

	c.JSON(http.StatusOK, gin.H{
		"batch":   batch,
		"records": records,
	})
}

// Approve approves a batch, anonymizes data, and deducts stock
func (h *BatchHandler) Approve(c *gin.Context) {
	claims, _ := middleware.GetUserFromContext(c)
	batchID := c.Param("id")

	batch, err := h.db.GetBatch(batchID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "batch not found"})
		return
	}

	if batch.PharmacyID != claims.PharmacyID {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	// Fetch pending records
	records, err := h.db.GetBatchRecords(batchID)
	if err != nil {
		h.log.Error("failed to fetch records", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to approve batch"})
		return
	}

	// Move records to approved_visits and deduct stock
	for _, record := range records {
		visit := &models.ApprovedVisit{
			ID:                uuid.New().String(),
			PharmacyID:        batch.PharmacyID,
			MedicineID:        record.MedicineID,
			QuantityDispensed: record.QuantityDispensed,
			Diagnosis:         record.Diagnosis,
			PatientData:       record.PatientData,
			VisitDate:         time.Now(),
		}

		if err := h.db.CreateApprovedVisit(visit); err != nil {
			h.log.Error("failed to create approved visit", "error", err.Error())
		}

		// Deduct stock
		if err := h.db.UpdateMedicineQuantity(record.MedicineID, -record.QuantityDispensed); err != nil {
			h.log.Error("failed to deduct stock", "error", err.Error())
		}

		// Delete pending record
		h.db.DeletePendingRecord(record.ID)
	}

	// Update batch status
	approvedBy := claims.UserID
	if err := h.db.UpdateBatchStatus(batchID, models.BatchStatusApproved, "", approvedBy); err != nil {
		h.log.Error("failed to update batch", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "approval failed"})
		return
	}

	h.audit.LogAction(claims.UserID, claims.PharmacyID, services.ActionBatchApproved, services.EntityTypeBatch, batchID, c.ClientIP(), map[string]interface{}{"record_count": len(records)})

	c.JSON(http.StatusOK, gin.H{"status": "approved"})
}

// Reject rejects a batch
func (h *BatchHandler) Reject(c *gin.Context) {
	claims, _ := middleware.GetUserFromContext(c)
	batchID := c.Param("id")

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reason required"})
		return
	}

	batch, err := h.db.GetBatch(batchID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "batch not found"})
		return
	}

	if batch.PharmacyID != claims.PharmacyID {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	// Delete all pending records
	records, _ := h.db.GetBatchRecords(batchID)
	for _, record := range records {
		h.db.DeletePendingRecord(record.ID)
	}

	// Update batch status
	if err := h.db.UpdateBatchStatus(batchID, models.BatchStatusRejected, req.Reason, claims.UserID); err != nil {
		h.log.Error("failed to reject batch", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "rejection failed"})
		return
	}

	h.audit.LogAction(claims.UserID, claims.PharmacyID, services.ActionBatchRejected, services.EntityTypeBatch, batchID, c.ClientIP(), map[string]interface{}{"reason": req.Reason})

	c.JSON(http.StatusOK, gin.H{"status": "rejected"})
}

// DeleteRecord removes a single pending record from a batch
func (h *BatchHandler) DeleteRecord(c *gin.Context) {
	claims, _ := middleware.GetUserFromContext(c)
	batchID := c.Param("id")
	recordID := c.Param("rid")

	batch, err := h.db.GetBatch(batchID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "batch not found"})
		return
	}

	if batch.PharmacyID != claims.PharmacyID {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.db.DeletePendingRecord(recordID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
