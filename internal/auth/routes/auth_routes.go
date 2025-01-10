package auth

import (
	"goGin/internal/auth/handler"
	"goGin/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/register", handler.Register)

		auth.POST("/encrypt", middleware.TokenValidationMiddleware(), handler.EncryptMessage)
		auth.POST("/decrypt", middleware.TokenValidationMiddleware(), handler.DecryptMessage)
		auth.POST("/decryptToken", middleware.TokenValidationMiddleware(), handler.DecryptToken)
		auth.POST("/decryptRefreshToken", middleware.TokenValidationMiddleware(), handler.DecryptRefreshToken)
	}
}
