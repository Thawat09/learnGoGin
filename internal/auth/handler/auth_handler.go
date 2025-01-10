package handler

import (
	// "fmt"
	"goGin/internal/auth/service"
	"goGin/internal/database"
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

	decryptedUsername, err := service.Decrypt(loginReq.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to decrypt username",
			"status": http.StatusInternalServerError,
		})
		return
	}

	decryptedPassword, err := service.Decrypt(loginReq.Password)
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

	user, err := service.Login(decryptedUsername, decryptedPassword, ipAddress)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  "Invalid credentials",
			"status": http.StatusUnauthorized,
		})
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

	// var json = jsoniter.ConfigCompatibleWithStandardLibrary
	// jsonData, err := json.MarshalIndent(userClaims, "", "  ")
	// if err != nil {
	// 	fmt.Println("Error converting to JSON:", err)
	// 	return
	// }

	// fmt.Println("JSON data:", string(jsonData))

	accessToken, err := service.CreateAccessToken(userClaims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to generate access token",
			"status": http.StatusInternalServerError,
		})
		return
	}

	refreshToken, err := service.CreateRefreshToken(userClaims)
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

	encryptedAccessToken, err := service.Encrypt(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to encrypt access token",
			"status": http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Login successful",
		"accessToken": encryptedAccessToken,
		"status":      http.StatusOK,
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

	if err := service.Register(userReq.Username, userReq.Password, userReq.Email, userReq.DepartmentId); err != nil {
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
