package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/nxrmqlly/arcfile-backend/handlers"
	"github.com/nxrmqlly/arcfile-backend/storage"
)

func main() {
	db, err := storage.InitDatabase()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}


	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Error closing the database: %v", err)
		}
	}()

	// Set up signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")

		// Close the database connection
		if err := db.Close(); err != nil {
			log.Printf("Error closing the database during shutdown: %v", err)
		}

		os.Exit(0) // Ensure the app exits after cleanup
	}()

	
	router := gin.Default()
	repo := storage.NewRepository(db)
	handler := handlers.New(repo)

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 10 << 20 // 10 MiB

	router.POST("/api/upload", handler.Upload)
	router.GET("/api/file/:identifier", handler.FileInfo)
	router.GET("/api/file/:identifier/download", handler.FileDownload)



	if err := router.Run("localhost:8080"); err != nil {
        log.Fatalf("Error starting server: %v", err)
    }
}
