package services

import (
	"biblio_go/models"
	"biblio_go/repositories"
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
