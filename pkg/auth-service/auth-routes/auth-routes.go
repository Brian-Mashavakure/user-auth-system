package auth_routes

import (
	auth_handlers "github.com/Brian-Mashavakure/user-auth-system/pkg/auth-service/auth-handlers"
	auth_middleware "github.com/Brian-Mashavakure/user-auth-system/pkg/auth-service/auth-middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	api := router.Group("api/auth")

	api.POST("/register", auth_handlers.RegisterUser)
	api.POST("/login", auth_middleware.TokenCheckMiddleware(), auth_handlers.LoginUser)
	api.POST("/deleteuser", auth_middleware.TokenCheckMiddleware(), auth_handlers.DeleteUser)
	api.POST("/tokenstatus", auth_middleware.TokenCheckMiddleware(), auth_handlers.TokenStatus)

}
