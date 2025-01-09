package handler

import (
	"goGin/internal/static/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	userID := c.Param("id")
	user, err := service.GetUserByID(userID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":  "User not found",
			"status": http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":   user,
		"status": http.StatusOK,
	})
}
