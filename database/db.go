package database

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "database.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable()
	createMovieTable()
	createAdmin()
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT UNIQUE,
		password TEXT,
		isAdmin BOOLEAN
	);
	`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func createMovieTable() {
	query := `
	CREATE TABLE IF NOT EXISTS movies (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		director TEXT NOT NULL,
		genre TEXT NOT NULL,
		release_year INTEGER NOT NULL,
		description TEXT
	);
	`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func createAdmin() {
	password, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

	query := `
	INSERT OR IGNORE INTO users (name, email, password, isAdmin)
	VALUES (?, ?, ?, ?)
	`

	_, err := DB.Exec(query, "Admin", "admin@mail.com", string(password), true)
	if err != nil {
		log.Fatal(err)
	}
}

func ResetDB() error {
	query := `DROP TABLE IF EXISTS users; DROP TABLE IF EXISTS movies;`
	_, err := DB.Exec(query)
	if err != nil {
		return err
	}

	createTable()
	createMovieTable()
	createAdmin()

	return nil
}
