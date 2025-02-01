package storage

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

func InitDatabase() (readonly, readwrite *sql.DB, err error) {
	// ensure dir exisits
	if err := os.MkdirAll("./data", os.ModePerm); err != nil {
		return nil, nil, fmt.Errorf("create data dir: %w", err)
	}
	readwrite, err = sql.Open("sqlite", "./data/arcfile.db?mode=readwrite&_journal_mode=wal&_txlock=immediate")
	if err != nil {
		return nil, nil, fmt.Errorf("open db: %w", err)
	}
	readwrite.SetMaxOpenConns(1)

	readonly, err = sql.Open("sqlite", "./data/arcfile.db?mode=readonly&_journal_mode=wal")
	if err != nil {
		return nil, nil, fmt.Errorf("open db: %w", err)
	}

	createTable := `
    CREATE TABLE IF NOT EXISTS files (
        identifier TEXT PRIMARY KEY,
        filename TEXT NOT NULL,
        uuid TEXT NOT NULL,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        expires_at DATETIME NOT NULL,
        email TEXT NOT NULL
    );`

	_, err = readwrite.Exec(createTable)
	if err != nil {
		return nil, nil, fmt.Errorf("create table: %w", err)
	}

	return readonly, readwrite, nil
}
