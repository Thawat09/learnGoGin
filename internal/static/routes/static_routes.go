package routes

import (
	"goGin/internal/middleware"
	"goGin/internal/static/handler"

	"github.com/gin-gonic/gin"
)

func RegisterStaticRoutes(r *gin.RouterGroup) {
	static := r.Group("/static")
	static.Use(middleware.AuthMiddleware())
	{
		static.GET("/:id", handler.GetUser)
	}
}
