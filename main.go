package main

import "github.com/gin-gonic/gin"

func main() {
	InitDB()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "DB connected ✅"})
	})

	r.POST("/user/register", Register)

	r.POST("/user/login", Login)

	auth := r.Group("/")
	auth.Use(AuthMiddleware())
	{
		auth.GET("/protected", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			c.JSON(200, gin.H{
				"message": "Access granted",
				"userID":  userID,
			})
		})
	}

	r.Run(":8080")
}
