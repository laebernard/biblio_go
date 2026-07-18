package handlers

import (
	"net/http"

	"biblio_go/logger"
	"biblio_go/models"
	"biblio_go/services"

	"strconv"

	"github.com/gin-gonic/gin"
)

func parseID(id string) uint {
	i, _ := strconv.Atoi(id)
	return uint(i)
}

// GetMovies godoc
// @Summary      Liste tous les films
// @Tags         movies
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header    string  true  "Bearer token (utilisateur ou admin): Bearer <token>"
// @Success      200  {array}   models.Movie
// @Failure      500  {object}  map[string]string
// @Router       /movies [get]
func GetMovies(c *gin.Context) {
	logger.Info("GET /movies - Fetching all movies")
	movies, err := services.GetMovies()
	if err != nil {
		logger.Error("GET /movies - Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot fetch movies"})
		return
	}

	logger.Info("GET /movies - Successfully fetched %d movies", len(movies))
	c.JSON(http.StatusOK, movies)
}

// GetMovie godoc
// @Summary      Récupère un film par ID
// @Tags         movies
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Movie ID"
// @Param        Authorization  header    string  true  "Bearer token (utilisateur ou admin): Bearer <token>"
// @Success      200  {object}  models.Movie
// @Failure      404  {object}  map[string]string
// @Router       /movies/{id} [get]
func GetMovie(c *gin.Context) {
	id := c.Param("id")

	movie, err := services.GetMovie(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}

	c.JSON(http.StatusOK, movie)
}

// CreateMovie godoc
// @Summary      Crée un film
// @Tags         movies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        movie  body      models.Movie  true  "Film à créer"
// @Param        Authorization  header    string  true  "Bearer token (utilisateur ou admin): Bearer <token>"
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movies [post]
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

// UpdateMovie godoc
// @Summary      Met à jour un film
// @Tags         movies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      string        true  "Movie ID"
// @Param        movie  body      models.Movie  true  "Film mis à jour"
// @Param        Authorization  header    string  true  "Bearer token (utilisateur ou admin): Bearer <token>"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movies/{id} [put]
func UpdateMovie(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Title       *string `json:"title"`
		Director    *string `json:"director"`
		Genre       *string `json:"genre"`
		ReleaseYear *int    `json:"release_year"`
		Description *string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	if input.Title == nil && input.Director == nil && input.Genre == nil && input.ReleaseYear == nil && input.Description == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	movie, err := services.GetMovie(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}
	if input.Director != nil {
		movie.Director = *input.Director
	}
	if input.Genre != nil {
		movie.Genre = *input.Genre
	}
	if input.ReleaseYear != nil {
		if *input.ReleaseYear <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid release_year"})
			return
		}
		movie.ReleaseYear = *input.ReleaseYear
	}
	if input.Description != nil {
		movie.Description = *input.Description
	}

	movie.ID = parseID(id)

	err = services.UpdateMovie(movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot update movie"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "movie updated"})
}

// DeleteMovie godoc
// @Summary      Supprime un film
// @Tags         movies
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Movie ID"
// @Param        Authorization  header    string  true  "Bearer token (utilisateur ou admin): Bearer <token>"
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movies/{id} [delete]
func DeleteMovie(c *gin.Context) {
	id := c.Param("id")

	err := services.DeleteMovie(id)
	if err != nil {
		logger.Error("DELETE /movies/:id - Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot delete movie"})
		return
	}

	logger.Info("DELETE /movies/:id - Movie deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "movie deleted"})
}

// SearchMovies godoc
// @Summary      Recherche des films sur tous les champs
// @Tags         movies
// @Produce      json
// @Param        query  query     string  true  "Terme de recherche"
// @Success      200    {array}   models.Movie
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /search [get]
func SearchMovies(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	logger.Info("GET /search - Searching movies with query: %s", query)
	movies, err := services.SearchMovies(query)
	if err != nil {
		logger.Error("GET /search - Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot search movies"})
		return
	}

	logger.Info("GET /search - Found %d movies", len(movies))
	c.JSON(http.StatusOK, movies)
}

// BulkDeleteMovies godoc
// @Summary      Supprime plusieurs films
// @Tags         movies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ids  body      []object  true  "Liste d'IDs ex: [{\"id\":1},{\"id\":2}]"
// @Param        Authorization  header    string  true  "Bearer token (utilisateur ou admin): Bearer <token>"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movies [delete]
func BulkDeleteMovies(c *gin.Context) {
	var input []struct {
		ID uint `json:"id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil || len(input) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body: expected [{\"id\":1},{\"id\":2}]"})
		return
	}

	ids := make([]uint, len(input))
	for i, item := range input {
		ids[i] = item.ID
	}

	if err := services.BulkDeleteMovies(ids); err != nil {
		logger.Error("DELETE /movies - BulkDelete error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot bulk delete movies"})
		return
	}

	logger.Info("DELETE /movies - Bulk deleted %d movies", len(ids))
	c.JSON(http.StatusOK, gin.H{"message": "movies deleted", "count": len(ids)})
}

// BulkUpdateMovies godoc
// @Summary      Met à jour plusieurs films
// @Tags         movies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        movies  body      []models.Movie  true  "Liste de films avec ID"
// @Param        Authorization  header    string  true  "Bearer token (utilisateur ou admin): Bearer <token>"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /movies/bulk [put]
func BulkUpdateMovies(c *gin.Context) {
	var movies []models.Movie

	if err := c.ShouldBindJSON(&movies); err != nil || len(movies) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body: expected array of movies with id"})
		return
	}

	if err := services.BulkUpdateMovies(movies); err != nil {
		logger.Error("PUT /movies/bulk - BulkUpdate error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot bulk update movies"})
		return
	}

	logger.Info("PUT /movies/bulk - Bulk updated %d movies", len(movies))
	c.JSON(http.StatusOK, gin.H{"message": "movies updated", "count": len(movies)})
}
