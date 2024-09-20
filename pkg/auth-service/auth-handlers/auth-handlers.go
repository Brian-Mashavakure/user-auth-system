package auth_handlers

import (
	"github.com/Brian-Mashavakure/user-auth-system/pkg/database"
	"github.com/Brian-Mashavakure/user-auth-system/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type User struct {
	gorm.Model
	ID          uint
	NAME        string `json:"name"`
	EMAIL       string `json:"email" gorm:"unique"`
	AGE         string `json:"age"`
	GENDER      string `json:"gender"`
	USERNAME    string `json:"username" gorm:"unique"`
	PASSWORD    string `json:"password"`
	USER_STATUS string `json:"user_status"`
}

func RegisterUser(c *gin.Context) {
	name := c.Request.FormValue("name")
	email := c.Request.FormValue("email")
	age := c.Request.FormValue("age")
	gender := c.Request.FormValue("gender")
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	user_status := "active"

	//hash user password
	hashedPassword := utils.HashPassword(password)

	//check for existing users
	var existingUser User

	oldEmail := database.DB.Where("email = ? ", email).First(&existingUser)
	if oldEmail.Error == nil {
		log.Println("Error occured trying to create user")
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"error": "User with that email address already exists"})
		return

	}

	oldUser := database.DB.Where("username = ? ", username).First(&existingUser)
	if oldUser.Error == nil {
		log.Println("Error occured trying to create user")
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"error": "User with that username already exists"})
		return

	}

	user := User{
		NAME:        name,
		EMAIL:       email,
		AGE:         age,
		GENDER:      gender,
		USERNAME:    username,
		PASSWORD:    hashedPassword,
		USER_STATUS: user_status,
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		log.Println("Error occured trying to create user")
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}
