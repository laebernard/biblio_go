package main

import (
	"database/sql"
	"log"

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
