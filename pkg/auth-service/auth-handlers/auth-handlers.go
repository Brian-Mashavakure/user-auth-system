package auth_handlers

import (
	"github.com/Brian-Mashavakure/user-auth-system/pkg/database"
	"github.com/Brian-Mashavakure/user-auth-system/pkg/utils"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
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

type Token struct {
	gorm.Model
	ID          uint
	USERNAME    string `json:"username" gorm:"unique"`
	TOKEN       string `json:"token" gorm:"unique"`
	START_DATE  string
	EXPIRY_DATE string
}

// TODO: Update function on registring a user that had been deleted
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

	oldEmail := database.DB.Where("email = ? AND user_status = ?", email, "active").First(&existingUser)
	if oldEmail.Error == nil {
		log.Println("Error occurred trying to create user")
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"error": "User with that email address already exists"})
		return

	}

	oldUser := database.DB.Where("username = ? AND user_status = ?", username, "active").First(&existingUser)
	if oldUser.Error == nil {
		log.Println("Error occurred trying to create user")
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
		log.Printf("Error occurred trying to create user\n: %n", result.Error)
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"error": "Failed to create user"})
		return
	}

	//generate token
	tokenString, tokenStartDate, tokenExpiryDate := utils.GenerateToken(username, email)

	//save token to db
	token := Token{
		USERNAME:    username,
		TOKEN:       tokenString,
		START_DATE:  tokenStartDate,
		EXPIRY_DATE: tokenExpiryDate,
	}

	tokenResult := database.DB.Create(&token)
	if tokenResult.Error != nil {
		log.Println("Error occcurred trying to save token")
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"error": "Error occurred trying to save token"})
		return
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary

	tokenJson, tokenErr := json.Marshal(tokenString)
	if tokenErr != nil {
		log.Println("Error trying to marshal token")
		c.Abort()
		return
	}

	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, string(tokenJson))
}

func LoginUser(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	//retrieve user
	user := User{USERNAME: username}
	result := database.DB.First(&user).Scan(&user)
	if result.Error != nil {
		log.Println("Error retrieving users from db")
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"error": "Incorrect email"})
		return
	}

	//compare passwords
	comparison := utils.ComparePasswordAndHash(user.PASSWORD, password)
	if comparison == false {
		log.Println("Passwords do not match")
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"error": "Incorrect password"})
		return
	}

	//check if user is active
	if user.USER_STATUS == "inactive" {
		log.Println("User was deleted")
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"error": "User was deleted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
}

func TokenStatus(c *gin.Context) {
	username := c.Request.FormValue("username")
	email := c.Request.FormValue("email")
	tokenString := c.Request.Header.Get("Authorization")

	token := Token{USERNAME: username}

	result := database.DB.Table("tokens").First(&token).Scan(&token)
	if result.Error != nil {
		log.Println("Error retrieving user from db")
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"error": "Failed to find username"})
		return
	}

	if tokenString == token.TOKEN {
		//check if token is still valid
		todayDate := time.Now().Format("02/01/2006")
		status := utils.CompareDates(todayDate, token.EXPIRY_DATE)
		if status == true {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "Token is still valid"})
		} else if status == false {
			newTokenString, tokenStartTime, tokenExpiryTime := utils.GenerateToken(username, email)
			newToken := Token{
				USERNAME:    username,
				TOKEN:       newTokenString,
				START_DATE:  tokenStartTime,
				EXPIRY_DATE: tokenExpiryTime,
			}

			newTokenResult := database.DB.Create(&newToken)
			if newTokenResult.Error != nil {
				log.Println("Error generating new token")
				c.AbortWithStatusJSON(http.StatusOK, gin.H{"error": "Failed to generate new token"})
				return
			}

			c.String(http.StatusOK, newTokenString)
		}

	} else {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "Token provided not in our database"})
	}

}

func DeleteUser(c *gin.Context) {
	username := c.Request.FormValue("username")

	//delete user
	result := database.DB.Table("users").Where("username = ?", username).Update("user_status", "inactive")
	if result.Error != nil {
		log.Println("Error deleting user")
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
