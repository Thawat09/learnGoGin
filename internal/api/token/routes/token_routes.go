package token

import (
	"goGin/internal/api/token/handler"
	"goGin/internal/middleware"

	"github.com/gin-gonic/gin"
)

func TokenRoutes(r *gin.RouterGroup) {
	token := r.Group("/token")
	{
		token.POST("/encrypt", middleware.TokenValidationMiddleware(), handler.EncryptMessage)
		token.POST("/decrypt", middleware.TokenValidationMiddleware(), handler.DecryptMessage)
		token.POST("/decryptToken", middleware.TokenValidationMiddleware(), handler.DecryptToken)
		token.POST("/decryptRefreshToken", middleware.TokenValidationMiddleware(), handler.DecryptRefreshToken)
	}
}
