package handlers

import "github.com/nxrmqlly/arcfile/storage"

type Handlers struct {
	repo *storage.Repository
}

func New(repo *storage.Repository) *Handlers {
	return &Handlers{repo: repo}
}
