package database

import (
	"database/sql"
	"time"

	"github.com/nxrmqlly/arcfile-backend/structures"
	_ "modernc.org/sqlite"
)

func CreateFile(
	identifier string,
	name string,
	path string,
	createdAt time.Time,
	expiresAt time.Time,
	email string,
) (structures.File, error) {
	db, err := sql.Open("sqlite", "./data/arcfile.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	query := "INSERT INTO files (identifier, name, path, created_at, expires_at, email) VALUES (?, ?, ?, ?, ?, ?)"

	rows, err := db.Query(query, identifier, name, path, createdAt, expiresAt, email)

	if err != nil {
		return structures.File{}, err
	}

	var File = structures.File{}

	rows.Next()

	if err := rows.Scan(
		&File.ID,
		&File.Identifier,
		&File.Filename,
		&File.Path,
		&File.CreatedAt,
		&File.ExpiresAt,
		&File.Email,
	); err != nil {
		return structures.File{}, err
	}

	return File, nil
}
