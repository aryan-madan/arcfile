package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nxrmqlly/arcfile-backend/database"
)

func UploadHandler(c *gin.Context) {
	
	// Source
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": "get form error: " + err.Error(),
		})
		return
	}
	
	cwd, err := os.Getwd()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": "internal error: " + err.Error(),
		})
		return
	}
	
	email := c.PostForm("email")
	relPath := filepath.Join("data", "uploads", file.Filename)
	pathToSave := filepath.Join(cwd, relPath)
	currentTime := time.Now()
	expiresAt := currentTime.Add(15 * time.Minute)

	database.CreateFile("090000", file.Filename, relPath, currentTime, expiresAt, email)

	if err := c.SaveUploadedFile(file, pathToSave); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": "error uploading file: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"identifier": "090000",
		"created_at": currentTime.String(),
		"expires_at": expiresAt.String(),
		"filename":   file.Filename,
		"email":      email,
	})
}
