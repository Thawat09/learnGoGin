package handler

import (
	"goGin/internal/api/token/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func EncryptMessage(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Invalid input",
			"status": http.StatusBadRequest,
		})
		return
	}

	encryptedUsername, err := service.Encrypt(req.Username)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to encrypt username",
			"status": http.StatusInternalServerError,
		})
		return
	}

	encryptedPassword, err := service.Encrypt(req.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to encrypt password",
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"encryptedUsername": encryptedUsername,
		"encryptedPassword": encryptedPassword,
		"status":            http.StatusOK,
	})
}

func DecryptMessage(c *gin.Context) {
	var req struct {
		EncryptedData string `json:"data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Invalid input",
			"status": http.StatusBadRequest,
		})
		return
	}

	decryptedData, err := service.Decrypt(req.EncryptedData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to decrypt data",
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"decryptedData": decryptedData,
		"status":        http.StatusOK,
	})
}

func DecryptToken(c *gin.Context) {
	var req struct {
		EncryptedToken string `json:"encryptedToken"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Invalid input",
			"status": http.StatusBadRequest,
		})
		return
	}

	decryptedToken, err := service.Decrypt(req.EncryptedToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to decrypt token",
			"status": http.StatusInternalServerError,
		})
		return
	}

	claims, err := service.ParseToken(decryptedToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  "Invalid or expired token",
			"status": http.StatusUnauthorized,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"claims": claims,
		"status": http.StatusOK,
	})
}

func DecryptRefreshToken(c *gin.Context) {
	var req struct {
		EncryptedRefreshToken string `json:"encryptedRefreshToken"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Invalid input",
			"status": http.StatusBadRequest,
		})
		return
	}

	decryptedRefreshToken, err := service.Decrypt(req.EncryptedRefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to decrypt refresh token",
			"status": http.StatusInternalServerError,
		})
		return
	}

	claims, err := service.ParseRefeshToken(decryptedRefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  "Invalid or expired refresh token",
			"status": http.StatusUnauthorized,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"claims": claims,
		"status": http.StatusOK,
	})
}

func DecryptTokenMiddleware(token string) (map[string]interface{}, error) {
	decryptedToken, err := service.Decrypt(token)
	if err != nil {
		return nil, err
	}

	claims, err := service.ParseToken(decryptedToken)
	if err != nil {
		return nil, err
	}

	claimsMap := map[string]interface{}{
		"userId":         claims.UserId,
		"username":       claims.Username,
		"email":          claims.Email,
		"firstName":      claims.FirstName,
		"lastName":       claims.LastName,
		"roleId":         claims.RoleId,
		"roleName":       claims.RoleName,
		"departmentId":   claims.DepartmentId,
		"departmentName": claims.DepartmentName,
		"exp":            claims.ExpiresAt,
	}

	return claimsMap, nil
}
