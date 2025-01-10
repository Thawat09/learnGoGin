package middleware

import (
	"fmt"
	authService "goGin/internal/auth/service"
	"goGin/internal/database"
	"goGin/internal/token/handler"

	tokenService "goGin/internal/token/service"
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

		parsedToken, err := tokenService.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error":  "Invalid Access Token",
				"status": http.StatusForbidden,
			})
			c.Abort()
			return
		}

		exp := parsedToken.ExpiresAt

		expirationTime := time.Unix(exp.Unix(), 0)
		if time.Now().After(expirationTime) {
			redisClient := database.GetRedisClient()
			accessTokenCacheKey := fmt.Sprintf("accessToken:%s", strconv.Itoa(int(parsedToken.UserId)))

			cachedRefreshToken, err := database.GetValue(redisClient, fmt.Sprintf("refreshToken:%s", strconv.Itoa(int(parsedToken.UserId))))
			if err != nil || cachedRefreshToken == "" {
				c.JSON(http.StatusForbidden, gin.H{
					"error":  "Refresh Token not found or expired",
					"status": http.StatusForbidden,
				})
				c.Abort()
				return
			}

			refreshParsedToken, err := jwt.Parse(cachedRefreshToken, nil)
			if err != nil || !refreshParsedToken.Valid {
				c.JSON(http.StatusForbidden, gin.H{
					"error":  "Invalid Refresh Token",
					"status": http.StatusForbidden,
				})
				c.Abort()
				return
			}

			refreshClaims, ok := refreshParsedToken.Claims.(jwt.MapClaims)
			if !ok || refreshClaims["exp"] == nil {
				c.JSON(http.StatusForbidden, gin.H{
					"error":  "Invalid Refresh Token Claims",
					"status": http.StatusForbidden,
				})
				c.Abort()
				return
			}

			refreshExp := int64(refreshClaims["exp"].(float64))
			refreshExpirationTime := time.Unix(refreshExp, 0)
			if time.Now().After(refreshExpirationTime) {
				c.JSON(http.StatusForbidden, gin.H{
					"error":  "Refresh Token has expired",
					"status": http.StatusForbidden,
				})
				c.Abort()
				return
			}

			serviceClaims := &authService.Claims{
				UserId:         int(refreshClaims["UserId"].(float64)),
				Username:       refreshClaims["Username"].(string),
				Email:          refreshClaims["Email"].(string),
				FirstName:      refreshClaims["FirstName"].(string),
				LastName:       refreshClaims["LastName"].(string),
				RoleId:         int(refreshClaims["RoleId"].(float64)),
				RoleName:       refreshClaims["RoleName"].(string),
				DepartmentId:   int(refreshClaims["DepartmentId"].(float64)),
				DepartmentName: refreshClaims["DepartmentName"].(string),
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
			database.SetValue(redisClient, accessTokenCacheKey, newAccessToken, 900)

			c.Set("user", serviceClaims)
			c.Set("newAccessToken", newAccessToken)

			c.Next()
			return
		}

		redisClient := database.GetRedisClient()
		accessTokenCacheKey := fmt.Sprintf("accessToken:%s", strconv.Itoa(int(parsedToken.UserId)))
		database.SetValue(redisClient, accessTokenCacheKey, token, 900)

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
