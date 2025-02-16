package storage

import (
	"context"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/nxrmqlly/arcfile/structures"
)

// starts a background cleaner for files
func (r *Repository) StartCleanupRoutine(interval time.Duration) {
	ctx := context.TODO()
	defer ctx.Done()
	go func() {

		// startup
		expired, err := r.ExpiredFiles(context.TODO())
		if err != nil {
			log.Printf("Failed to fetch expired files: %v", err)
			return
		}

		r.CleanupExpired(ctx, expired)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			expired, err := r.ExpiredFiles(context.TODO())
			if err != nil {
				log.Printf("Failed to fetch expired files: %v", err)
				return
			}

			r.CleanupExpired(ctx, expired)
		}
	}()
}

// delete all expired files
func (r *Repository) CleanupExpired(ctx context.Context, files []structures.File) {
	storagePath := path.Join("data", "uploads")

	var deletedCount int
	for _, file := range files {
		filePath := filepath.Join(storagePath, file.UUID)
		err := r.DelteFile(ctx, file.Identifier, filePath)
		if err != nil {
			log.Printf("Failed to delete file %s: %v", filePath, err)
			continue
		}

		deletedCount++
		log.Printf("Cleaned up: %s (UUID: %s)", file.Identifier, file.UUID)
	}

	if deletedCount != 0 {
		log.Printf("Cleanup completed. Removed %d expired entries", deletedCount)
	}

}

// get expired files
func (r *Repository) ExpiredFiles(ctx context.Context) ([]structures.File, error) {
	timeNow := time.Now().UTC()
	rows, err := r.pool.Query(ctx, `
        SELECT identifier, uuid 
        FROM arcfile_files 
        WHERE expires_at <= $1
        `,
		timeNow)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []structures.File

	for rows.Next() {
		var file structures.File
		if err := rows.Scan(&file.Identifier, &file.UUID); err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

// deletes file from disk
func (r *Repository) DelteFile(ctx context.Context, identifier string, filePath string) error {
	if err := r.DeleteDatabaseEntry(ctx, identifier); err != nil {
		return err
	}
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// deletes from database
func (r *Repository) DeleteDatabaseEntry(ctx context.Context, identifier string) error {
	_, err := r.pool.Exec(ctx, `
        DELETE FROM arcfile_files 
        WHERE identifier = $1`,
		identifier)

	return err
}
