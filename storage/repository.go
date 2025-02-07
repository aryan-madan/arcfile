package storage

import (
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	pool *pgx.Conn
}

func NewRepository(pool *pgx.Conn) *Repository {
	return &Repository{pool: pool}
}
