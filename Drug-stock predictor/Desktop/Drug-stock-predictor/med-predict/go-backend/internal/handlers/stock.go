package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"med-predict/go-backend/internal/middleware"
	"med-predict/go-backend/internal/models"
)

type StockHandler struct {
	db *sql.DB
}

// NewStockHandler creates a new stock handler
func NewStockHandler(db *sql.DB) *StockHandler {
	return &StockHandler{db: db}
}

// ListMedicines lists all medicines for a pharmacy
func (h *StockHandler) ListMedicines(c *gin.Context) {
	_, pharmacyID, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "User not authenticated"})
		return
	}

	rows, err := h.db.Query(`
		SELECT id, pharmacy_id, name, generic_name, dosage, unit, quantity_on_hand,
		       reorder_level, unit_cost, selling_price, expiry_date, manufacturer_id,
		       batch_number, is_active, created_at, updated_at
		FROM medicines
		WHERE pharmacy_id = $1 AND is_active = true
		ORDER BY name
	`, pharmacyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to fetch medicines"})
		return
	}
	defer rows.Close()

	medicines := []models.MedicineResponse{}
	for rows.Next() {
		var med models.MedicineResponse
		err := rows.Scan(
			&med.ID, &med.PharmacyID, &med.Name, &med.GenericName, &med.Dosage, &med.Unit,
			&med.QuantityOnHand, &med.ReorderLevel, &med.UnitCost, &med.SellingPrice,
			&med.ExpiryDate, &med.ManufacturerID, &med.BatchNumber, &med.IsActive,
			&med.CreatedAt, &med.UpdatedAt,
		)
		if err != nil {
			continue
		}
		medicines = append(medicines, med)
	}

	c.JSON(http.StatusOK, medicines)
}

// CreateMedicine creates a new medicine
func (h *StockHandler) CreateMedicine(c *gin.Context) {
	_, pharmacyID, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "User not authenticated"})
		return
	}

	var req models.MedicineCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	medicineID := generateID("med")
	now := time.Now()

	query := `
		INSERT INTO medicines (id, pharmacy_id, name, generic_name, dosage, unit, quantity_on_hand,
		                        reorder_level, unit_cost, selling_price, expiry_date, manufacturer_id,
		                        batch_number, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`

	_, err = h.db.Exec(query, medicineID, pharmacyID, req.Name, req.GenericName, req.Dosage, req.Unit,
		req.QuantityOnHand, req.ReorderLevel, req.UnitCost, req.SellingPrice, req.ExpiryDate,
		req.ManufacturerID, req.BatchNumber, true, now, now)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to create medicine"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": medicineID, "message": "Medicine created successfully"})
}

// UpdateMedicine updates a medicine
func (h *StockHandler) UpdateMedicine(c *gin.Context) {
	_, pharmacyID, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "User not authenticated"})
		return
	}

	medicineID := c.Param("id")

	// Check if medicine belongs to user's pharmacy
	var existingPharmacyID string
	err = h.db.QueryRow("SELECT pharmacy_id FROM medicines WHERE id = $1", medicineID).Scan(&existingPharmacyID)
	if err != nil || existingPharmacyID != pharmacyID {
		c.JSON(http.StatusForbidden, gin.H{"detail": "Access denied"})
		return
	}

	var req models.MedicineUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	// Build dynamic update query
	updates := "updated_at = $1"
	args := []interface{}{time.Now()}
	argIndex := 2

	if req.Name != nil {
		updates += ", name = $" + strconv.Itoa(argIndex)
		args = append(args, *req.Name)
		argIndex++
	}
	if req.QuantityOnHand != nil {
		updates += ", quantity_on_hand = $" + strconv.Itoa(argIndex)
		args = append(args, *req.QuantityOnHand)
		argIndex++
	}
	// Add other fields as needed...

	args = append(args, medicineID)

	query := `UPDATE medicines SET ` + updates + ` WHERE id = $` + strconv.Itoa(argIndex)
	_, err = h.db.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to update medicine"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Medicine updated successfully"})
}

// DeleteMedicine soft deletes a medicine
func (h *StockHandler) DeleteMedicine(c *gin.Context) {
	_, pharmacyID, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "User not authenticated"})
		return
	}

	medicineID := c.Param("id")

	// Check if medicine belongs to user's pharmacy
	var existingPharmacyID string
	err = h.db.QueryRow("SELECT pharmacy_id FROM medicines WHERE id = $1", medicineID).Scan(&existingPharmacyID)
	if err != nil || existingPharmacyID != pharmacyID {
		c.JSON(http.StatusForbidden, gin.H{"detail": "Access denied"})
		return
	}

	query := `UPDATE medicines SET is_active = false, updated_at = $1 WHERE id = $2`
	_, err = h.db.Exec(query, time.Now(), medicineID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to delete medicine"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Medicine deleted successfully"})
}

// GetMedicine retrieves a specific medicine
func (h *StockHandler) GetMedicine(c *gin.Context) {
	_, pharmacyID, _, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "User not authenticated"})
		return
	}

	medicineID := c.Param("id")

	var med models.MedicineResponse
	query := `
		SELECT id, pharmacy_id, name, generic_name, dosage, unit, quantity_on_hand,
		       reorder_level, unit_cost, selling_price, expiry_date, manufacturer_id,
		       batch_number, is_active, created_at, updated_at
		FROM medicines
		WHERE id = $1 AND pharmacy_id = $2
	`

	err = h.db.QueryRow(query, medicineID, pharmacyID).Scan(
		&med.ID, &med.PharmacyID, &med.Name, &med.GenericName, &med.Dosage, &med.Unit,
		&med.QuantityOnHand, &med.ReorderLevel, &med.UnitCost, &med.SellingPrice,
		&med.ExpiryDate, &med.ManufacturerID, &med.BatchNumber, &med.IsActive,
		&med.CreatedAt, &med.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Medicine not found"})
		return
	}

	c.JSON(http.StatusOK, med)
}
