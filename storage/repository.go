package storage

import "database/sql"

type Repository struct {
	rdb, wdb *sql.DB
}

func NewRepository(rdb *sql.DB, wdb *sql.DB) *Repository {
	return &Repository{rdb: rdb, wdb: wdb}
}
