package middleware

import (
	"fmt"
	"goGin/internal/auth/handler"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		status := c.Writer.Status()
		fmt.Printf("Response Status: %d\n", status)
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token != "valid-token" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "Unauthorized",
				"status": http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func TokenValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":  "Endpoint not found",
				"status": http.StatusNotFound,
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "Invalid authorization format",
				"status": http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		token := parts[1]

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "Unauthorized",
				"status": http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		claims, err := handler.DecryptTokenMiddleware(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  err.Error(),
				"status": http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		exp, ok := claims["exp"].(*jwt.NumericDate)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "Invalid expiration format",
				"status": http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		expirationTime := exp.Time
		fmt.Println("expirationTime:", expirationTime)

		if time.Now().After(expirationTime) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "Token expired",
				"status": http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		username, ok := claims["username"].(string)
		if !ok || username != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "Unauthorized - Not admin",
				"status": http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
