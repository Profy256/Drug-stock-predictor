package handlers

import (
	"net/http"

	"med-predict-backend/internal/db"
	"med-predict-backend/internal/middleware"
	"med-predict-backend/internal/models"
	"med-predict-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db     *db.Database
	audit  *services.AuditService
	log    *services.Logger
	secret string
}

func NewAuthHandler(database *db.Database, audit *services.AuditService, log *services.Logger, secret string) *AuthHandler {
	return &AuthHandler{
		db:     database,
		audit:  audit,
		log:    log,
		secret: secret,
	}
}

// Login authenticates a user and returns a JWT token
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Get user by email
	user, err := h.db.GetUserByEmail(req.Email)
	if err != nil {
		h.log.Warn("login failed: user not found", "email", req.Email, "ip", c.ClientIP())
		h.audit.LogAction("", "", services.ActionLoginFailed, services.EntityTypeUser, "", c.ClientIP(), map[string]string{"reason": "user not found"})
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		h.log.Warn("login failed: invalid password", "email", req.Email, "ip", c.ClientIP())
		h.audit.LogAction(user.ID, user.PharmacyID, services.ActionLoginFailed, services.EntityTypeUser, user.ID, c.ClientIP(), map[string]string{"reason": "invalid password"})
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate JWT
	token, err := middleware.GenerateToken(user, h.secret)
	if err != nil {
		h.log.Error("failed to generate token", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	// Log successful login
	h.audit.LogAction(user.ID, user.PharmacyID, services.ActionLoginSuccess, services.EntityTypeUser, user.ID, c.ClientIP(), nil)

	c.JSON(http.StatusOK, models.LoginResponse{
		Token:      token,
		UserID:     user.ID,
		PharmacyID: user.PharmacyID,
		Name:       user.Name,
		Email:      user.Email,
		Role:       user.Role,
	})
}

// Register creates a new user account
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.log.Error("password hashing failed", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}

	// Create user
	user := &models.User{
		ID:           uuid.New().String(),
		PharmacyID:   req.PharmacyID,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		IsActive:     true,
	}

	if err := h.db.CreateUser(user); err != nil {
		h.log.Error("failed to create user", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}

	h.audit.LogAction(user.ID, user.PharmacyID, services.ActionUserCreated, services.EntityTypeUser, user.ID, c.ClientIP(), map[string]string{"name": user.Name})
	h.log.Info("user registered successfully", "user_id", user.ID, "email", user.Email)

	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"role":  user.Role,
	})
}

// GetMe returns the current user's profile
func (h *AuthHandler) GetMe(c *gin.Context) {
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.db.GetUserByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          user.ID,
		"name":        user.Name,
		"email":       user.Email,
		"role":        user.Role,
		"pharmacy_id": user.PharmacyID,
	})
}
