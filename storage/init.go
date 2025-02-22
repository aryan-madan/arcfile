package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDatabase() (*pgxpool.Pool, error) {
	// ensure dir exisits
	if err := os.MkdirAll("./data", os.ModePerm); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	createTable := `
    CREATE TABLE IF NOT EXISTS arcfile_files (
        identifier VARCHAR(6) PRIMARY KEY,
        filename TEXT NOT NULL,
        uuid TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        expires_at TIMESTAMP NOT NULL,
        email TEXT NOT NULL
    );`

	_, err = pool.Exec(context.Background(), createTable)
	if err != nil {
		return nil, fmt.Errorf("create table: %w", err)
	}

	return pool, nil
}
