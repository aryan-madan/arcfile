package database

import (
	"log"
	"math/rand"
	"time"

	"github.com/nxrmqlly/arcfile-backend/structures"
	_ "modernc.org/sqlite"
)

const alphanumericChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// generates a random identifier for the file
func generateIdentifier(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = alphanumericChars[rand.Intn(len(alphanumericChars))]
	}
	return string(b)
}

// creates a new file entry in database
func CreateFile(
	file *structures.File,
) error {

	var identifier string = generateIdentifier(6)

	log.Println("generated identifier:", identifier)

	query := `
        INSERT INTO files 
        (identifier, filename, uuid, created_at, expires_at, email) 
        VALUES (?, ?, ?, ?, ?, ?) 
        RETURNING identifier;`

	err := DB.QueryRow(
		query,
		identifier,
		file.Filename,
		file.UUID,
		file.CreatedAt,
		file.ExpiresAt,
		file.Email,
	).Scan(
		&file.Identifier,
	)
	if err != nil {
		return err
	}

	return nil

}
