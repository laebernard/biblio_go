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

	r.Run(":8080")
}
