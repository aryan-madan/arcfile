package main

import (
	"github.com/nxrmqlly/arcfile-backend/database"
	_ "modernc.org/sqlite"
)

func main() {
	database.InitDatabase()
}
