package middleware

import (
	"fmt"
	authService "goGin/internal/api/auth/service"
	"goGin/internal/api/token/handler"
	"goGin/internal/config/database"

	tokenService "goGin/internal/api/token/service"
	"net/http"

	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		ip := c.ClientIP()
		status := c.Writer.Status()

		if ip == "::1" || ip == "" {
			ip = "localhost"
		}

		fmt.Printf("Response Status: %d IP: %s\n", status, ip)
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "Token is required",
				"status": http.StatusUnauthorized,
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

		deToken, err := tokenService.Decrypt(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "Failed to decrypt data",
				"status": http.StatusInternalServerError,
			})
			c.Abort()
			return
		}

		parsedToken, err := tokenService.ParseTokenForExp(deToken)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error":  "Invalid Access Token",
				"status": http.StatusForbidden,
			})
			c.Abort()
			return
		}

		redisClient := database.GetRedisClient()

		userId, ok := parsedToken["UserId"].(float64)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"error":  "Invalid UserId in token",
				"status": http.StatusForbidden,
			})
			c.Abort()
			return
		}

		cachedAccessToken, err := database.GetValue(redisClient, fmt.Sprintf("accessToken:%s", strconv.Itoa(int(userId))))
		if err != nil || cachedAccessToken == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":  "Access Token not found or expired",
				"status": http.StatusForbidden,
			})
			c.Abort()
			return
		}

		cachedRefreshToken, err := database.GetValue(redisClient, fmt.Sprintf("refreshToken:%s", strconv.Itoa(int(userId))))
		if err != nil || cachedRefreshToken == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":  "Refresh Token not found or expired",
				"status": http.StatusForbidden,
			})
			c.Abort()
			return
		}

		refreshParsedToken, err := tokenService.ParseRefeshToken(cachedRefreshToken)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error":  "Invalid Refresh Token",
				"status": http.StatusForbidden,
			})
			c.Abort()
			return
		}

		refreshExp := refreshParsedToken.ExpiresAt
		refreshExpirationTime := time.Unix(refreshExp.Unix(), 0)
		if time.Now().After(refreshExpirationTime) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":  "Refresh Token has expired",
				"status": http.StatusForbidden,
			})
			c.Abort()
			return
		} else {
			exp, _ := parsedToken["Exp"].(float64)
			expirationTime := time.Unix(int64(exp), 0)
			if time.Now().After(expirationTime) {
				serviceClaims := &authService.Claims{
					UserId:         refreshParsedToken.UserId,
					Username:       refreshParsedToken.Username,
					Email:          refreshParsedToken.Email,
					FirstName:      refreshParsedToken.FirstName,
					LastName:       refreshParsedToken.LastName,
					RoleId:         refreshParsedToken.RoleId,
					RoleName:       refreshParsedToken.RoleName,
					DepartmentId:   refreshParsedToken.DepartmentId,
					DepartmentName: refreshParsedToken.DepartmentName,
				}

				newAccessToken, err := authService.CreateAccessToken(serviceClaims)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":  "Failed to create new access token",
						"status": http.StatusInternalServerError,
					})
					c.Abort()
					return
				}

				database.SetValue(redisClient, fmt.Sprintf("accessToken:%s", newAccessToken), newAccessToken, 3600)

				encryptedAccessToken, err := tokenService.Encrypt(newAccessToken)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":  "Failed to encrypt access token",
						"status": http.StatusInternalServerError,
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"accessToken": encryptedAccessToken,
					"status":      http.StatusOK,
				})

				c.Set("user", serviceClaims)
				c.Set("newAccessToken", newAccessToken)

				c.Abort()
				return
			}
		}

		c.Set("user", parsedToken)
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
