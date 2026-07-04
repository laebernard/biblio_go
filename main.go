package main

import (
	"biblio_go/database"
	"biblio_go/security"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "DB connected ✅"})
	})

	r.POST("/user/register", security.Register)

	r.POST("/user/login", security.Login)

	auth := r.Group("/")
	auth.Use(security.AuthMiddleware())
	{
		auth.GET("/protected", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			c.JSON(200, gin.H{
				"message": "Access granted",
				"userID":  userID,
			})
		})
		auth.PUT("/user/me", security.UpdateMe)
		auth.PUT("/user/:id", security.UpdateUserByAdmin)
		auth.DELETE("/reset", security.ResetDatabase)
	}

	r.Run(":8080")
}
