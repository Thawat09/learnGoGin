package handler

import (
	"fmt"
	"goGin/internal/api/static/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetStatistics(c *gin.Context) {
	userID := c.Param("id")
	user, _ := service.GetUserByID(userID)

	fmt.Println("User ID: ", userID)

	c.JSON(http.StatusOK, gin.H{
		"user":   user,
		"status": http.StatusOK,
	})
}
