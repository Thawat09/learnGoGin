package health

import (
	"goGin/internal/check/handler"

	"github.com/gin-gonic/gin"
)

func CheckRoutes(r *gin.RouterGroup) {
	check := r.Group("/check")
	{
		check.GET("/health", handler.Health)
	}
}
