package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/nxrmqlly/arcfile-backend/database"
)

func UploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "get form error: " + err.Error(),
		})
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "internal error: " + err.Error(),
		})
		return
	}

	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "email is required",
		})
		return
	}

	filename := filepath.Base(file.Filename)
	fileUUID := uuid.NewString()
	pathToSave := filepath.Join(cwd, "data", "uploads", fileUUID)

	uploadDir := filepath.Join(cwd, "data", "uploads")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "error creating upload directory: " + err.Error(),
		})
		return
	}

	currentTime := time.Now()
	expiresAt := currentTime.Add(15 * time.Minute)

	retFile, err := database.CreateFile(filename, fileUUID, currentTime, expiresAt, email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "error saving file metadata: " + err.Error(),
		})
		return
	}

	if err := c.SaveUploadedFile(file, pathToSave); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "error uploading file: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file_id":    retFile.ID,
		"identifier": retFile.Identifier,
		"created_at": retFile.CreatedAt.String(),
		"expires_at": retFile.ExpiresAt.String(),
		"filename":   retFile.Filename,
		"email":      retFile.Email,
	})
}
