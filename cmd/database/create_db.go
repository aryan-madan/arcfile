package main

import (
	"github.com/joho/godotenv"
	"github.com/nxrmqlly/arcfile/storage"
	_ "modernc.org/sqlite"
)

func main() {
	godotenv.Load()

	if _, err := storage.InitDatabase(); err != nil {
		panic(err)
	}
}
