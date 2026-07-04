package repositories

import (
	"biblio_go/database"
	"biblio_go/models"
)

func GetAllMovies() ([]models.Movie, error) {
	var movies []models.Movie
	err := DB.Find(&movies).Error
	return movies, err
}

func GetMovieByID(id string) (models.Movie, error) {
	var movie models.Movie
	err := database.DB.First(&movie, id).Error
	return movie, err
}

func CreateMovie(movie models.Movie) error {
	return database.DB.Create(&movie).Error
}

func UpdateMovie(movie models.Movie) error {
	return database.DB.Save(&movie).Error
}

func DeleteMovie(id string) error {
	return database.DB.Delete(&models.Movie{}, id).Error
}
