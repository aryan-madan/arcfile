package storage

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func InitDatabase() (ro, rw *sql.DB, err error) {
	// ensure dir exisits
	if err := os.MkdirAll("./data", os.ModePerm); err != nil {
		return nil, nil, fmt.Errorf("create data dir: %w", err)
	}

	ro, err = sql.Open("sqlite3", "file:./data/arcfile.db?mode=ro&_journal_mode=wal")
	if err != nil {
		return nil, nil, fmt.Errorf("open db: %w", err)
	}

	rw, err = sql.Open("sqlite3", "file:./data/arcfile.db?mode=rw&_journal_mode=wal&_txlock=immediate")
	if err != nil {
		return nil, nil, fmt.Errorf("open db: %w", err)
	}
	rw.SetMaxOpenConns(1)

	createTable := `
    CREATE TABLE IF NOT EXISTS files (
        identifier TEXT PRIMARY KEY,
        filename TEXT NOT NULL,
        uuid TEXT NOT NULL,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        expires_at DATETIME NOT NULL,
        email TEXT NOT NULL
    );`

	_, err = rw.Exec(createTable)
	if err != nil {
		return nil, nil, fmt.Errorf("create table: %w", err)
	}

	return ro, rw, nil
}
