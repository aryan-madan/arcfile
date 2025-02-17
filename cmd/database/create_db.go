package main

import (
	"github.com/joho/godotenv"
	"github.com/nxrmqlly/arcfile/storage"
	_ "modernc.org/sqlite"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	if _, err := storage.InitDatabase(); err != nil {
		panic(err)
	}
}
