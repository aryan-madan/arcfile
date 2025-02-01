package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/nxrmqlly/arcfile-backend/structures"
)

// POST /api/upload
func (h *Handlers) Upload(c *gin.Context) {
	formFile, err := c.FormFile("file")
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

	filename := filepath.Base(formFile.Filename)
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

	currentTime := time.Now().UTC()
	expiresAt := currentTime.Add(10 * time.Minute)

	file := structures.File{
		Identifier: "",
		Filename:   filename,
		UUID:       fileUUID,
		CreatedAt:  currentTime,
		ExpiresAt:  expiresAt,
		Email:      email,
	}

	if err := h.repo.CreateFile(c, &file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "error saving file metadata: " + err.Error(),
		})
		return
	}

	if err := c.SaveUploadedFile(formFile, pathToSave); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "error uploading file: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "file uploaded successfully",
		"data": gin.H{
			"identifier": file.Identifier,
			"filename":   file.Filename,
			"uuid":       file.UUID,
			"created_at": file.CreatedAt,
			"expires_at": file.ExpiresAt,
			"email":      file.Email,
		},
	})
}
