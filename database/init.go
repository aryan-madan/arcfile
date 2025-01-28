package database

import (
	"database/sql"
	"os"
)

var DB *sql.DB

func InitDatabase() {
	//ensure dir exisits
	if err := os.MkdirAll("./data", os.ModePerm); err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite", "./data/arcfile.db")
	if err != nil {
		panic(err)
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
		panic(err)
	}
}
