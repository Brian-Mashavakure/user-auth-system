package main

import (
	"fmt"
	auth_handlers "github.com/Brian-Mashavakure/user-auth-system/pkg/auth-service/auth-handlers"
	auth_routes "github.com/Brian-Mashavakure/user-auth-system/pkg/auth-service/auth-routes"
	"github.com/Brian-Mashavakure/user-auth-system/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error trying to load .env file, please check")
	}

	database.DatabaseConnector()

	port := os.Getenv("PORT")

	//auto-migrate user table
	database.DB.AutoMigrate(&auth_handlers.User{}, &auth_handlers.Token{})

	//router
	router := gin.Default()

	auth_routes.AuthRoutes(router)

	address := fmt.Sprintf("localhost:%s", port)

	router.Run(address)

}
