package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"strconv"
	"time"
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

// hash strings for custom token
func HashString(input string) string {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	salt := os.Getenv("HASH_SALT")
	byteInput := []byte(input + salt)

	md5hash := md5.Sum(byteInput)

	return hex.EncodeToString(md5hash[:])
}

// compare two date strings
func CompareDates(todayDate, expiryDate string) bool {
	//parse dates with the correct format (DD-MM-YYYY)
	parseToday, err1 := time.Parse("02-01-2006", todayDate)
	if err1 != nil {
		log.Println("Failed to parse today's date")
		panic(err1)
	}

	parseExpiry, err2 := time.Parse("02-01-2006", expiryDate)
	if err2 != nil {
		log.Println("Failed to parse expiry date")
		panic(err2)
	}

	if parseExpiry.After(parseToday) {
		return true
	} else {
		return false
	}
}

func GenerateToken(username, email string) (string, string, string) {
	startTime := time.Now()
	expiryTime := startTime.Add(10 * 24 * time.Hour)
	fmtStartDate := startTime.Format("02-01-2006")
	fmtExpiryDate := expiryTime.Format("02-01-2006")
	stringToHash := username + email + fmtStartDate + fmtExpiryDate
	token := HashString(stringToHash)
	return token, fmtStartDate, fmtExpiryDate
}
