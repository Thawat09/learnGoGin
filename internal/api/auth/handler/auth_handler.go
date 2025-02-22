package handler

import (
	// "fmt"
	authService "goGin/internal/api/auth/service"
	tokenService "goGin/internal/api/token/service"
	"goGin/internal/config/database"
	"log"
	"net/http"
	"strconv"
	"time"

	// jsoniter "github.com/json-iterator/go"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var loginReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Invalid input",
			"status": http.StatusBadRequest,
		})
		return
	}

	decryptedUsername, err := tokenService.Decrypt(loginReq.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to decrypt username",
			"status": http.StatusInternalServerError,
		})
		return
	}

	decryptedPassword, err := tokenService.Decrypt(loginReq.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to decrypt password",
			"status": http.StatusInternalServerError,
		})
		return
	}

	ipAddress := c.ClientIP()
	if ipAddress == "::1" || ipAddress == "" {
		ipAddress = "localhost"
	}

	user, err := authService.Login(decryptedUsername, decryptedPassword, ipAddress)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  "Invalid credentials",
			"status": http.StatusUnauthorized,
		})
		return
	}

	userClaims := &authService.Claims{
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

	// var json = jsoniter.ConfigCompatibleWithStandardLibrary
	// jsonData, err := json.MarshalIndent(userClaims, "", "  ")
	// if err != nil {
	// 	fmt.Println("Error converting to JSON:", err)
	// 	return
	// }

	// fmt.Println("JSON data:", string(jsonData))

	accessToken, err := authService.CreateAccessToken(userClaims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to generate access token",
			"status": http.StatusInternalServerError,
		})
		return
	}

	refreshToken, err := authService.CreateRefreshToken(userClaims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to generate refresh token",
			"status": http.StatusInternalServerError,
		})
		return
	}

	userId := userClaims.UserId
	redisClient := database.GetRedisClient()

	go func() {
		key := "accessToken:" + strconv.Itoa(userId)
		err := database.SetValue(redisClient, key, accessToken, time.Hour)
		if err != nil {
			log.Println("Failed to save accessToken to Redis:", err)
		}
	}()

	go func() {
		key := "refreshToken:" + strconv.Itoa(userId)
		err := database.SetValue(redisClient, key, refreshToken, 10*time.Hour)
		if err != nil {
			log.Println("Failed to save refreshToken to Redis:", err)
		}
	}()

	encryptedAccessToken, err := tokenService.Encrypt(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to encrypt access token",
			"status": http.StatusInternalServerError,
		})
		return
	}

	encryptedRefreshToken, err := tokenService.Encrypt(refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to encrypt access token",
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Login successful",
		"accessToken":  encryptedAccessToken,
		"refreshToken": encryptedRefreshToken,
		"status":       http.StatusOK,
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Invalid input",
			"status": http.StatusBadRequest,
		})
		return
	}

	if userReq.Username == "" || userReq.Email == "" || userReq.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Username, email, and password must not be empty",
			"status": http.StatusBadRequest,
		})
		return
	}

	if err := authService.Register(userReq.Username, userReq.Password, userReq.Email, userReq.DepartmentId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to register",
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"status":  http.StatusOK,
	})
}
