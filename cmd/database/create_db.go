package main

import (
	"github.com/nxrmqlly/arcfile-backend/storage"
	_ "modernc.org/sqlite"
)

func main() {
	storage.InitDatabase()
}
