package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"med-predict/go-backend/internal/models"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func init() {
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("your-secret-key-change-in-production")
	}
}

// GenerateToken generates a JWT token for a user
func GenerateToken(userID, pharmacyID string, role models.UserRole) (string, error) {
	claims := jwt.MapClaims{
		"user_id":     userID,
		"pharmacy_id": pharmacyID,
		"role":        role,
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
		"iat":         time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// AuthMiddleware verifies JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"detail": "Missing authorization header"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix if present
		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"detail": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"detail": "Invalid token claims"})
			c.Abort()
			return
		}

		// Store claims in context
		c.Set("user_id", claims["user_id"])
		c.Set("pharmacy_id", claims["pharmacy_id"])
		c.Set("role", models.UserRole(claims["role"].(string)))

		c.Next()
	}
}

// RequireRole checks if user has required role
func RequireRole(roles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleInterface, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"detail": "User role not found"})
			c.Abort()
			return
		}

		userRole := roleInterface.(models.UserRole)
		hasRole := false

		for _, role := range roles {
			if userRole == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"detail": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserFromContext extracts user info from context
func GetUserFromContext(c *gin.Context) (userID, pharmacyID string, role models.UserRole, err error) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return "", "", "", errors.New("user_id not found in context")
	}

	pharmacyIDInterface, exists := c.Get("pharmacy_id")
	if !exists {
		return "", "", "", errors.New("pharmacy_id not found in context")
	}

	roleInterface, exists := c.Get("role")
	if !exists {
		return "", "", "", errors.New("role not found in context")
	}

	return userIDInterface.(string), pharmacyIDInterface.(string), roleInterface.(models.UserRole), nil
}
