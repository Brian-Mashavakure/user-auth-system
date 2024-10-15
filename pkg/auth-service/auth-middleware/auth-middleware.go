package auth_middleware

import (
	auth_handlers "github.com/Brian-Mashavakure/user-auth-system/pkg/auth-service/auth-handlers"
	"github.com/Brian-Mashavakure/user-auth-system/pkg/database"
	"github.com/Brian-Mashavakure/user-auth-system/pkg/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func TokenCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Token")
		username := c.Request.FormValue("username")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Request not authorized"})
			c.Abort()
			return
		}

		//retrieve user token info
		token := auth_handlers.Token{USERNAME: username}
		result := database.DB.Table("tokens").Where("username = ?", username).Scan(&token)
		if result.Error != nil {
			log.Println("Error finding user or token in db")
			c.Abort()
			return
		}

		//check if token string is the same
		if tokenString != token.TOKEN {
			log.Println("Invalid token string")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token string provided"})
			return
		}

		//check if token is still valid
		todayDate := time.Now().Format("02-01-2006")
		status := utils.CompareDates(todayDate, token.EXPIRY_DATE)
		if status == false {
			log.Println("Token has expired")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Expired token string provided"})
			return
		}

		//continue
		c.Next()
	}

}
