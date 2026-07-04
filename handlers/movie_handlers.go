package handlers

import (
	"net/http"

	"biblio_go/models"
	"biblio_go/services"

	"strconv"

	"github.com/gin-gonic/gin"
)

func parseID(id string) uint {
	i, _ := strconv.Atoi(id)
	return uint(i)
}

func GetMovies(c *gin.Context) {
	movies, err := services.GetMovies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot fetch movies"})
		return
	}

	c.JSON(http.StatusOK, movies)
}

func GetMovie(c *gin.Context) {
	id := c.Param("id")

	movie, err := services.GetMovie(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}

	c.JSON(http.StatusOK, movie)
}

func CreateMovie(c *gin.Context) {
	var movie models.Movie

	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	err := services.CreateMovie(movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create movie"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "movie created"})
}

func UpdateMovie(c *gin.Context) {
	id := c.Param("id")

	var movie models.Movie

	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	// 👉 on applique l'id de l'URL au movie
	movie.ID = parseID(id)

	err := services.UpdateMovie(movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot update movie"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "movie updated"})
}

func DeleteMovie(c *gin.Context) {
	id := c.Param("id")

	err := services.DeleteMovie(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot delete movie"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "movie deleted"})
}
