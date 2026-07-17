package main

// @title Movie Catalog API
// @version 1.0
// @description API REST pour gerer un catalogue de films avec authentification JWT
// @termsOfService http://swagger.io/terms/
//
// @contact.name API Support
// @contact.email support@example.com
//
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host
// @BasePath /
// @schemes
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Saisir le token JWT au format: ******
//
// @externalDocs.description OpenAPI
// @externalDocs.url https://swagger.io
import (
	"os"

	"biblio_go/database"
	"biblio_go/docs"
	"biblio_go/handlers"
	"biblio_go/security"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	database.InitDB()

	// Keep Swagger host/scheme dynamic so production (Render HTTPS domain) works.
	docs.SwaggerInfo.Host = ""
	docs.SwaggerInfo.Schemes = []string{}

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "DB connected ✅"})
	})

	r.POST("/user/register", security.Register)
	r.POST("/user/login", security.Login)
	r.GET("/search", handlers.SearchMovies)

	auth := r.Group("/")
	auth.Use(security.AuthMiddleware())
	{
		auth.GET("/movies", handlers.GetMovies)
		auth.GET("/movies/:id", handlers.GetMovie)
		auth.POST("/movies", handlers.CreateMovie)
		auth.PUT("/movies/bulk", handlers.BulkUpdateMovies)
		auth.PUT("/movies/:id", handlers.UpdateMovie)
		auth.DELETE("/movies", handlers.BulkDeleteMovies)
		auth.DELETE("/movies/:id", handlers.DeleteMovie)

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

	backup := r.Group("/backup")
	backup.Use(security.AuthMiddleware())
	{
		backup.GET("/config", handlers.BackupConfig)
		backup.POST("/config", handlers.RestoreConfig)
		backup.GET("/state", handlers.BackupState)
		backup.POST("/state", handlers.RestoreState)
	}

	// Render impose d’utiliser le port fourni par l’environnement
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback local
	}

	r.Run(":" + port)
}
