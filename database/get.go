package database

import (
	"database/sql"

	"github.com/nxrmqlly/arcfile-backend/structures"
	_ "modernc.org/sqlite"
)

func GetFile(identifier string) (structures.File, error) {
	db, err := sql.Open("sqlite", "./data/arcfile.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	query := "SELECT * FROM files WHERE identifier = ?"

	rows, err := db.Query(query, identifier)

	if err != nil {
		return structures.File{}, err
	}

	var File = structures.File{}

	rows.Next()

	if err := rows.Scan(
		&File.Identifier,
		&File.Filename,
		&File.UUID,
		&File.CreatedAt,
		&File.ExpiresAt,
		&File.Email,
	); err != nil {
		return structures.File{}, err
	}

	return File, nil
}
