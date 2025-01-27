package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nxrmqlly/arcfile-backend/handlers"
)

func main() {
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 10 << 20 // 10 MiB

	router.POST("/upload", handlers.UploadHandler)
	router.GET("/upload", handlers.UploadHandler)

	router.Run("localhost:8080")
}
