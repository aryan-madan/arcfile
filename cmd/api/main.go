package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nxrmqlly/arcfile-backend/handlers"
	"github.com/nxrmqlly/arcfile-backend/ratelimits"
	"github.com/nxrmqlly/arcfile-backend/storage"
)

func main() {

	// fmt.Println("Hello world!") // for the win!!!
	roDB, rwDB, err := storage.InitDatabase()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	defer func() {
		if err := roDB.Close(); err != nil {
			log.Fatalf("Error closing the database: %v", err)
		}
		if err := rwDB.Close(); err != nil {
			log.Fatalf("Error closing the database: %v", err)
		}
	}()

	// signal handling - graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")

		// close the database connection
		if err := roDB.Close(); err != nil {
			log.Printf("Error closing the database during shutdown: %v", err)
		}
		if err := rwDB.Close(); err != nil {
			log.Printf("Error closing the database during shutdown: %v", err)
		}

		os.Exit(0)
	}()

	repo := storage.NewRepository(roDB, rwDB)
	handler := handlers.New(repo)

	repo.StartCleanupRoutine(30 * time.Second)

	// Rate limiter setup
	limiters := ratelimits.SetupLimiters()

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router := gin.Default()
	router.MaxMultipartMemory = 10 << 20 // 10 MiB

	router.Static("/static", "./static") // Serves static CSS and JS
	router.LoadHTMLGlob("templates/*")   // Loads HTML templates

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	router.POST("/api/upload", limiters["postFile"], handler.Upload)
	router.GET("/api/file/:identifier", limiters["getFile"], handler.FileInfo)
	router.GET("/api/file/:identifier/download", handler.FileDownload)

	if err := router.Run("localhost:8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
