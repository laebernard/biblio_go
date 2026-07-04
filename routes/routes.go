package routes

import (
	"biblio_go/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/movies", handlers.GetMovies)
	r.GET("/movies/:id", handlers.GetMovie)
	r.POST("/movies", handlers.CreateMovie)
	r.PUT("/movies/:id", handlers.UpdateMovie)
	r.DELETE("/movies/:id", handlers.DeleteMovie)
}
