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
	query := "SELECT * FROM arcfile_files WHERE identifier = $1"

	row := r.pool.QueryRow(ctx, query, identifier)

	var File = structures.File{}

	if err := row.Scan(
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
