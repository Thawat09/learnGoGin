package handler

import (
	// "encoding/json"
	"fmt"
	"goGin/internal/auth/service"
	"goGin/internal/database"
	"log"
	"net/http"
	"time"

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

	ipAddress := c.ClientIP()
	if ipAddress == "::1" || ipAddress == "" {
		ipAddress = "localhost"
	}

	user, err := service.Login(decryptedUsername, decryptedPassword, ipAddress)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	userClaims := &service.Claims{
		UserId:         user.UserId,
		Username:       decryptedUsername,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		RoleId:         user.RoleId,
		RoleName:       user.RoleName,
		DepartmentId:   user.DepartmentId,
		DepartmentName: user.DepartmentName,
	}

	// jsonData, err := json.MarshalIndent(userClaims, "", "  ")
	// if err != nil {
	// 	fmt.Println("Error converting to JSON:", err)
	// 	return
	// }

	// fmt.Println("JSON data:", string(jsonData))

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

	userId := userClaims.UserId
	redisClient := database.GetRedisClient()

	go func() {
		err := database.SetValue(redisClient, fmt.Sprintf("accessToken:%d", userId), accessToken, time.Hour)
		if err != nil {
			log.Println("Failed to save accessToken to Redis:", err)
		}
	}()

	go func() {
		err := database.SetValue(redisClient, fmt.Sprintf("refreshToken:%d", userId), refreshToken, 10*time.Hour)
		if err != nil {
			log.Println("Failed to save refreshToken to Redis:", err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"message":     "Login successful",
		"accessToken": accessToken,
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
