package storage

import "database/sql"

type Repository struct {
	ro, rw *sql.DB
}

func NewRepository(ro, rw *sql.DB) *Repository {
	return &Repository{ro: ro, rw: rw}
}
