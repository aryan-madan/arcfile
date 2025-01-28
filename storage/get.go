package storage

import (
	"context"
	"fmt"

	"github.com/nxrmqlly/arcfile-backend/structures"
	_ "modernc.org/sqlite"
)

type FileNotFoundError struct {
	identifier string
}

func (e *FileNotFoundError) Error() string {
	return fmt.Sprintf("file not found with identifier: %s", e.identifier)
}

func (r *Repository) GetFile(ctx context.Context, identifier string) (structures.File, error) {
	query := "SELECT * FROM files WHERE identifier = ?"

	rows, err := r.db.QueryContext(ctx, query, identifier)

	if err != nil {
		return structures.File{}, &FileNotFoundError{identifier}
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
