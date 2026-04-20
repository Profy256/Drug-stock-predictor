package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"med-predict-backend/internal/config"
	"med-predict-backend/internal/models"
	"med-predict-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AuthorizationHeader = "Authorization"
	UserCtxKey          = "user"
	PharmacyCtxKey      = "pharmacy_id"
)

// Claims represents JWT payload
type Claims struct {
	UserID     string          `json:"user_id"`
	PharmacyID string          `json:"pharmacy_id"`
	Email      string          `json:"email"`
	Role       models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates JWT tokens
func AuthMiddleware(cfg *config.Config, logger *services.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			logger.Warn("missing authorization token", "ip", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization token"})
			c.Abort()
			return
		}

		claims, err := validateToken(token, cfg.JWTSecret)
		if err != nil {
			logger.Warn("invalid token", "error", err.Error(), "ip", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set(UserCtxKey, claims)
		c.Set(PharmacyCtxKey, claims.PharmacyID)
		c.Next()
	}
}

// GenerateToken creates a JWT token for a user
func GenerateToken(user *models.User, secret string) (string, error) {
	claims := &Claims{
		UserID:     user.ID,
		PharmacyID: user.PharmacyID,
		Email:      user.Email,
		Role:       user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GetUserFromContext extracts user claims from gin context
func GetUserFromContext(c *gin.Context) (*Claims, error) {
	user, exists := c.Get(UserCtxKey)
	if !exists {
		return nil, fmt.Errorf("user not found in context")
	}

	claims, ok := user.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid user claims type")
	}

	return claims, nil
}

func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader(AuthorizationHeader)
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func validateToken(tokenString, secret string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// RequireRole checks if user has one of the required roles
func RequireRole(roles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := GetUserFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		hasRole := false
		for _, r := range roles {
			if claims.Role == r {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
