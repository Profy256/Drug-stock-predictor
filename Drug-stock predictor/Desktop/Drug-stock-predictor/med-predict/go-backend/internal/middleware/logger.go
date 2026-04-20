package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggerMiddleware logs incoming requests
func LoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := "anonymous"
		if user, exists := c.Get(UserCtxKey); exists {
			if claims, ok := user.(*Claims); ok {
				userID = claims.UserID
			}
		}

		logger.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"ip":     c.ClientIP(),
			"user":   userID,
		}).Info("request received")

		c.Next()
	}
}
