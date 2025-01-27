package database

import (
	"database/sql"
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
	name string,
	uuid string,
	createdAt time.Time,
	expiresAt time.Time,
	email string,
) (structures.File, error) {
	db, err := sql.Open("sqlite", "./data/arcfile.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var identifier string

	for {
		identifier = generateIdentifier(6)

		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM files WHERE identifier = ?", identifier).Scan(&count)
		if err != nil {
			log.Fatalln("error getting identifier count:", err)
			return structures.File{}, err
		}
		if count == 0 {
			break // unique
		}
	}

	query := "INSERT INTO files (identifier, filename, uuid, created_at, expires_at, email) VALUES (?, ?, ?, ?, ?, ?)"

	result, err := db.Exec(query, identifier, name, uuid, createdAt, expiresAt, email)
	if err != nil {
		log.Fatalln("error inserting file:", err)
		return structures.File{}, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Fatalln("error getting last insert ID:", err)
		return structures.File{}, err
	}

	var file structures.File
	err = db.QueryRow("SELECT id, identifier, filename, uuid, created_at, expires_at, email FROM files WHERE id = ?", lastInsertID).Scan(
		&file.ID,
		&file.Identifier,
		&file.Filename,
		&file.UUID,
		&file.CreatedAt,
		&file.ExpiresAt,
		&file.Email,
	)
	if err != nil {
		log.Fatalln("error fetching inserted file:", err)
		return structures.File{}, err
	}

	return file, nil
}
