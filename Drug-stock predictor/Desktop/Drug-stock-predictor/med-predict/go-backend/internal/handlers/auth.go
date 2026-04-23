package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"med-predict/go-backend/internal/middleware"
	"med-predict/go-backend/internal/models"
)

type AuthHandler struct {
	db *sql.DB
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to hash password"})
		return
	}

	// Create user in database
	userID := generateID("user")
	now := time.Now()

	query := `
		INSERT INTO users (id, pharmacy_id, name, email, password_hash, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = h.db.Exec(query, userID, req.PharmacyID, req.Name, req.Email, string(hashedPassword), "data_entrant", true, now, now)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user_id": userID})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	// Get user from database
	var userID, pharmacyID, passwordHash string
	var role string

	query := `
		SELECT id, pharmacy_id, password_hash, role
		FROM users
		WHERE email = $1 AND is_active = true
	`

	err := h.db.QueryRow(query, req.Email).Scan(&userID, &pharmacyID, &passwordHash, &role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "Invalid credentials"})
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "Invalid credentials"})
		return
	}

	// Generate token
	token, err := middleware.GenerateToken(userID, pharmacyID, models.UserRole(role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to generate token"})
		return
	}

	// Get user details
	var name string
	query = `SELECT name FROM users WHERE id = $1`
	h.db.QueryRow(query, userID).Scan(&name)

	c.JSON(http.StatusOK, models.LoginResponse{
		Token:      token,
		UserID:     userID,
		PharmacyID: pharmacyID,
		Name:       name,
		Email:      req.Email,
		Role:       models.UserRole(role),
	})
}

// GetCurrentUser returns the current authenticated user
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, pharmacyID, role, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "User not authenticated"})
		return
	}

	// Get user details
	var name, email string
	var isActive bool
	var createdAt, updatedAt time.Time

	query := `
		SELECT name, email, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err = h.db.QueryRow(query, userID).Scan(&name, &email, &isActive, &createdAt, &updatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, models.UserResponse{
		ID:         userID,
		PharmacyID: pharmacyID,
		Name:       name,
		Email:      email,
		Role:       role,
		IsActive:   isActive,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	})
}

// Helper function to generate IDs
func generateID(prefix string) string {
	hash := sha256.Sum256([]byte(time.Now().String()))
	return fmt.Sprintf("%s_%s", prefix, hex.EncodeToString(hash[:8]))
}
