package services

import (
	"movie-catalog-api/models"
	"movie-catalog-api/repositories"
)

func GetMovies() ([]models.Movie, error) {
	return repositories.GetAllMovies()
}

func GetMovie(id string) (models.Movie, error) {
	return repositories.GetMovieByID(id)
}

func CreateMovie(movie models.Movie) error {
	return repositories.CreateMovie(movie)
}

func UpdateMovie(movie models.Movie) error {
	return repositories.UpdateMovie(movie)
}

func DeleteMovie(id string) error {
	return repositories.DeleteMovie(id)
}
