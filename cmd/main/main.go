package main

import (
	auth_handlers "github.com/Brian-Mashavakure/user-auth-system/pkg/auth-service/auth-handlers"
	auth_routes "github.com/Brian-Mashavakure/user-auth-system/pkg/auth-service/auth-routes"
	"github.com/Brian-Mashavakure/user-auth-system/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	database.DatabaseConnector()

	//auto-migrate user table
	database.DB.AutoMigrate(&auth_handlers.User{}, &auth_handlers.Token{})

	//router
	router := gin.Default()

	auth_routes.AuthRoutes(router)

	router.Run("192.168.0.7:8080")

}
