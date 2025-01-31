package storage

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/nxrmqlly/arcfile-backend/structures"
)

// starts a background cleaner for files //
func (r *Repository) StartCleanupRoutine(interval time.Duration) {
	go func() {

		// startup
		r.CleanupExpiredEntries()

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			r.CleanupExpiredEntries()
		}
	}()
}

func (r *Repository) ExpiredFiles() ([]structures.File, error) {
	rows, err := r.db.Query(`
        SELECT identifier, uuid 
        FROM files 
        WHERE expires_at <= ?
        `,
		time.Now().Format(time.RFC3339))
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

func (r *Repository) CleanupExpiredEntries() {
	storagePath := path.Join("data", "uploads")

	files, err := r.ExpiredFiles()
	if err != nil {
		log.Printf("Failed to fetch expired files: %v", err)
		return
	}

	var deletedCount int
	for _, file := range files {
		filePath := filepath.Join(storagePath, file.UUID)
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			log.Printf("File deletion failed (%s): %v", file.UUID, err)
			continue
		}

		if err := r.deleteDatabaseEntry(file.Identifier); err != nil {
			log.Printf("Database cleanup failed (%s): %v", file.Identifier, err)
			continue
		}

		deletedCount++
		log.Printf("Cleaned up: %s (UUID: %s)", file.Identifier, file.UUID)
	}

	log.Printf("Cleanup completed. Removed %d expired entries", deletedCount)
}

func (r *Repository) deleteDatabaseEntry(identifier string) error {
	_, err := r.db.Exec(`
        DELETE FROM files 
        WHERE identifier = ?`,
		identifier)

	return err
}
