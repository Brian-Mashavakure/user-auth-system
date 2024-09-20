package auth_routes

import (
	auth_handlers "github.com/Brian-Mashavakure/user-auth-system/pkg/auth-service/auth-handlers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	api := router.Group("api/auth")

	api.POST("/registeruser", auth_handlers.RegisterUser)

}
