package database

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func DatabaseConnector() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error trying to load .env file, please check")
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("USER")
	name := os.Getenv("DB_NAME")
	password := os.Getenv("PASSWORD")

	postgressetup := fmt.Sprintf("host=%s port=%s user=%s dbname=%s  password=%s sslmode=require", host, port, user, name, password)

	db, sqlErr := gorm.Open(postgres.Open(postgressetup), &gorm.Config{})
	if sqlErr != nil {
		panic(sqlErr)
	}

	DB = db
}
