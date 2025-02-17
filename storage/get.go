package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/nxrmqlly/arcfile/structures"
	_ "modernc.org/sqlite"
)

var ErrFileNotFound = errors.New("file not found")

type FileNotFoundError struct {
	Identifier string
}

func (e *FileNotFoundError) Error() string {
	return fmt.Sprintf("file not found with identifier: %s", e.Identifier)
}

// Fetch a file from the database
func (r *Repository) GetFile(ctx context.Context, identifier string) (structures.File, error) {
	query := "SELECT * FROM arcfile_files WHERE identifier = $1"

	row := r.pool.QueryRow(ctx, query, identifier)

	var file structures.File

	if err := row.Scan(
		&file.Identifier,
		&file.Filename,
		&file.UUID,
		&file.CreatedAt,
		&file.ExpiresAt,
		&file.Email,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return structures.File{}, &FileNotFoundError{Identifier: identifier}
		}
		return structures.File{}, fmt.Errorf("failed to get file: %w", err)
	}

	return file, nil
}
