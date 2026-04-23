package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"med-predict/go-backend/internal/middleware"
	"med-predict/go-backend/internal/models"
)

type AdminHandler struct {
	db *sql.DB
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(db *sql.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

// ListUsers lists all users
func (h *AdminHandler) ListUsers(c *gin.Context) {
	_, _, role, err := middleware.GetUserFromContext(c)
	if err != nil || role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"detail": "Only admins can list users"})
		return
	}

	rows, err := h.db.Query(`
		SELECT id, pharmacy_id, name, email, role, is_active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	users := []models.UserResponse{}
	for rows.Next() {
		var user models.UserResponse
		var role string
		err := rows.Scan(
			&user.ID, &user.PharmacyID, &user.Name, &user.Email,
			&role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			continue
		}
		user.Role = models.UserRole(role)
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

// DeactivateUser deactivates a user
func (h *AdminHandler) DeactivateUser(c *gin.Context) {
	_, _, role, err := middleware.GetUserFromContext(c)
	if err != nil || role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"detail": "Only admins can deactivate users"})
		return
	}

	userID := c.Param("id")

	query := `UPDATE users SET is_active = false WHERE id = $1`
	_, err = h.db.Exec(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to deactivate user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deactivated successfully"})
}

// GetAuditLogs gets audit logs
func (h *AdminHandler) GetAuditLogs(c *gin.Context) {
	_, _, role, err := middleware.GetUserFromContext(c)
	if err != nil || role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"detail": "Only admins can view audit logs"})
		return
	}

	// Placeholder implementation
	logs := []gin.H{}

	c.JSON(http.StatusOK, logs)
}

// ListPharmacies lists all pharmacies
func (h *AdminHandler) ListPharmacies(c *gin.Context) {
	_, _, role, err := middleware.GetUserFromContext(c)
	if err != nil || role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"detail": "Only admins can list pharmacies"})
		return
	}

	rows, err := h.db.Query(`
		SELECT id, name, region, district, lat, lng, contact_phone, whatsapp_number, is_active, created_at, updated_at
		FROM pharmacies
		ORDER BY name
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to fetch pharmacies"})
		return
	}
	defer rows.Close()

	pharmacies := []models.PharmacyResponse{}
	for rows.Next() {
		var pharm models.PharmacyResponse
		err := rows.Scan(
			&pharm.ID, &pharm.Name, &pharm.Region, &pharm.District,
			&pharm.Lat, &pharm.Lng, &pharm.ContactPhone, &pharm.WhatsAppNumber,
			&pharm.IsActive, &pharm.CreatedAt, &pharm.UpdatedAt,
		)
		if err != nil {
			continue
		}
		pharmacies = append(pharmacies, pharm)
	}

	c.JSON(http.StatusOK, pharmacies)
}
