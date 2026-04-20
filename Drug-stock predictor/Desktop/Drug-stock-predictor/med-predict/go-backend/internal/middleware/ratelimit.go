package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	// authLimiter: 10 requests per 15 minutes for auth endpoints
	authLimiter = rate.NewLimiter(rate.Every(90*60/10), 10) // 15min / 10 requests

	// apiLimiter: 200 requests per minute for general API
	apiLimiter = rate.NewLimiter(rate.Every(60/200), 200) // 1min / 200 requests
)

// RateLimitAuth applies strict rate limiting to auth endpoints
func RateLimitAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !authLimiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many login attempts. try again in 15 minutes",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RateLimitAPI applies general rate limiting to API endpoints
func RateLimitAPI() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !apiLimiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
