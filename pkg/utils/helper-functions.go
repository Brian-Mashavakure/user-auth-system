package utils

import (
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"strconv"
)

// hash user password
func HashPassword(pwd string) string {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	hashingCost := os.Getenv("HashCost")
	intCost, intErr := strconv.Atoi(hashingCost)
	if intErr != nil {
		log.Println("Failed to convert string to int")
	}

	pwdHash, hashErr := bcrypt.GenerateFromPassword([]byte(pwd), intCost)
	if hashErr != nil {
		fmt.Printf("Error trying to hash password: %s\n", hashErr)
	}
	return string(pwdHash)
}

// compare user password and stored hash
func ComparePasswordAndHash(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePwd := []byte(plainPwd)

	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
