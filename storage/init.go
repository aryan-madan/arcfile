package storage

import (
	"database/sql"
	"fmt"
	"os"
)

var DB *sql.DB

func InitDatabase() (*sql.DB, error) {
	// ensure dir exisits
	if err := os.MkdirAll("./data", os.ModePerm); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}

	db, err := sql.Open("sqlite", "./data/arcfile.db")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	createTable := `
    CREATE TABLE IF NOT EXISTS files (
        identifier TEXT PRIMARY KEY,
        filename TEXT NOT NULL,
        uuid TEXT NOT NULL,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        expires_at DATETIME NOT NULL,
        email TEXT NOT NULL
    );`

	_, err = db.Exec(createTable)
	if err != nil {
		return nil, fmt.Errorf("create table: %w", err)
	}

	return db, nil
}
