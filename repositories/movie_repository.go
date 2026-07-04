package repositories

import (
	"fmt"
	"strings"

	"biblio_go/database"
	"biblio_go/models"
)

func GetAllMovies() ([]models.Movie, error) {
	var movies []models.Movie

	rows, err := database.DB.Query("SELECT id, title, director, genre, release_year, description FROM movies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(&movie.ID, &movie.Title, &movie.Director, &movie.Genre, &movie.ReleaseYear, &movie.Description)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	return movies, rows.Err()
}

func GetMovieByID(id string) (models.Movie, error) {
	var movie models.Movie

	err := database.DB.QueryRow("SELECT id, title, director, genre, release_year, description FROM movies WHERE id = ?", id).
		Scan(&movie.ID, &movie.Title, &movie.Director, &movie.Genre, &movie.ReleaseYear, &movie.Description)

	return movie, err
}

func CreateMovie(movie models.Movie) error {
	_, err := database.DB.Exec(
		"INSERT INTO movies (title, director, genre, release_year, description) VALUES (?, ?, ?, ?, ?)",
		movie.Title,
		movie.Director,
		movie.Genre,
		movie.ReleaseYear,
		movie.Description,
	)
	return err
}

func UpdateMovie(movie models.Movie) error {
	_, err := database.DB.Exec(
		"UPDATE movies SET title = ?, director = ?, genre = ?, release_year = ?, description = ? WHERE id = ?",
		movie.Title,
		movie.Director,
		movie.Genre,
		movie.ReleaseYear,
		movie.Description,
		movie.ID,
	)
	return err
}

func DeleteMovie(id string) error {
	_, err := database.DB.Exec("DELETE FROM movies WHERE id = ?", id)
	return err
}

func SearchMovies(query string) ([]models.Movie, error) {
	var movies []models.Movie
	searchPattern := fmt.Sprintf("%%%s%%", query)

	rows, err := database.DB.Query(`
		SELECT id, title, director, genre, release_year, description FROM movies
		WHERE title LIKE ? OR director LIKE ? OR genre LIKE ? OR description LIKE ?
	`, searchPattern, searchPattern, searchPattern, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(&movie.ID, &movie.Title, &movie.Director, &movie.Genre, &movie.ReleaseYear, &movie.Description)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	return movies, rows.Err()
}

func BulkDeleteMovies(ids []uint) error {
	if len(ids) == 0 {
		return fmt.Errorf("ids list cannot be empty")
	}

	placeholders := strings.Repeat("?,", len(ids)-1) + "?"
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	query := fmt.Sprintf("DELETE FROM movies WHERE id IN (%s)", placeholders)
	_, err := database.DB.Exec(query, args...)
	return err
}

func BulkUpdateMovies(movies []models.Movie) error {
	for _, movie := range movies {
		err := UpdateMovie(movie)
		if err != nil {
			return err
		}
	}
	return nil
}
