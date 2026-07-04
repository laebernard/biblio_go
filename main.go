package main

import (
	"biblio_go/database"
	"biblio_go/handlers"
	"biblio_go/security"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "biblio_go/docs"
)

func main() {
	database.InitDB()

	r := gin.Default()

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "DB connected ✅"})
	})

	r.POST("/user/register", security.Register)

	r.POST("/user/login", security.Login)

	// Routes protégées par authentification
	auth := r.Group("/")
	auth.Use(security.AuthMiddleware())
	{
		// Routes pour les films
		auth.GET("/movies", handlers.GetMovies)
		auth.GET("/movies/:id", handlers.GetMovie)
		auth.POST("/movies", handlers.CreateMovie)
		auth.PUT("/movies/:id", handlers.UpdateMovie)
		auth.DELETE("/movies/:id", handlers.DeleteMovie)

		// Routes utilisateur
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
