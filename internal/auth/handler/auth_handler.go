package handler

import (
	"goGin/internal/auth/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var loginReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	decryptedUsername, err := service.Decrypt(loginReq.Username)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt username"})
		return
	}

	decryptedPassword, err := service.Decrypt(loginReq.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt password"})
		return
	}

	err = service.Login(decryptedUsername, decryptedPassword)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	userClaims := &service.Claims{
		Username: decryptedUsername,
		Email:    loginReq.Username,
	}

	accessToken, err := service.CreateAccessToken(userClaims)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := service.CreateRefreshToken(userClaims)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Login successful",
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func Register(c *gin.Context) {
	var userReq struct {
		Username     string `json:"username"`
		Password     string `json:"password"`
		Email        string `json:"email"`
		DepartmentId int    `json:"departmentId"`
	}

	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if userReq.Username == "" || userReq.Email == "" || userReq.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username, email, and password must not be empty"})
		return
	}

	if err := service.Register(userReq.Username, userReq.Password, userReq.Email, userReq.DepartmentId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func EncryptMessage(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	encryptedUsername, err := service.Encrypt(req.Username)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt message"})
		return
	}

	encryptedPassword, err := service.Encrypt(req.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"encryptedUsername": encryptedUsername,
		"encryptedPassword": encryptedPassword,
	})
}
