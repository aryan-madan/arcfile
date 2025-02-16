package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nxrmqlly/arcfile/handlers"
	"github.com/nxrmqlly/arcfile/public"
	"github.com/nxrmqlly/arcfile/ratelimits"
	"github.com/nxrmqlly/arcfile/storage"
)

var tag = `
    ___              _____ __   
   /   |  __________/ __(_) /__ 
  / /| | / ___/ ___/ /_/ / / _ \
 / ___ |/ /  / /__/ __/ / /  __/
/_/  |_/_/   \___/_/ /_/_/\___/ 

Repository: https://github.com/nxrmqlly/arcfile
Author: nxrmqlly (Ritam Das)
`

func main() {
	godotenv.Load()
	fmt.Fprintf(os.Stdout, "\033[0;31m%s\033[0m\n\n", tag)

	pool, err := storage.InitDatabase()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	defer pool.Close()

	// signal handling - graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")

		// close the database connection
		pool.Close()

		os.Exit(0)
	}()

	repo := storage.NewRepository(pool)
	handler := handlers.New(repo)

	repo.StartCleanupRoutine(30 * time.Second)

	// Rate limiter setup
	limiters := ratelimits.SetupLimiters()

	mode := os.Getenv("GIN_MODE")
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 10 << 20 // 10 MiB

	router.StaticFS("/public", http.FS(public.Templates))
	router.StaticFileFS("/about", "about.html", http.FS(public.Templates))
	router.StaticFileFS("/", "app.html", http.FS(public.Templates))

	router.POST("/api/upload", limiters["postFile"], handler.Upload)
	router.GET("/api/file/:identifier", limiters["getFile"], handler.FileInfo)
	router.GET("/api/download/:identifier", handler.FileDownload)

	var addr string
	if mode == "release" {
		addr = "0.0.0.0:8080"
	} else {
		addr = "localhost:8080"
	}
	if err := router.Run(addr); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
