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
	query := `DROP TABLE IF EXISTS users;`
	_, err := DB.Exec(query)
	if err != nil {
		return err
	}

	createTable()
	createAdmin()

	return nil
}
