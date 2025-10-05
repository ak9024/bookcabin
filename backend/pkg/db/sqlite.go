package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func NewSQLiteConnection(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to SQLite database!")
	return db, nil
}
