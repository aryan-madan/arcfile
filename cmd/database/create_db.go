package main

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "./data/arcfile.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	createTable := `
CREATE TABLE IF NOT EXISTS files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    identifier TEXT NOT NULL,
    name TEXT NOT NULL,
    path TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL,
    email TEXT NOT NULL
);

`
	_, err = db.Exec(createTable)
	db.Close()

	if err != nil {
		panic(err)
	}
}
