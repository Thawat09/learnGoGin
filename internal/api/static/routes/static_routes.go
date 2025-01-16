package routes

import (
	"goGin/internal/api/static/handler"
	"goGin/internal/middleware"

	"github.com/gin-gonic/gin"
)

func StaticRoutes(r *gin.RouterGroup) {
	static := r.Group("/static")
	static.Use(middleware.AuthMiddleware())
	{
		static.GET("/data", handler.GetStatistics)
	}
}
