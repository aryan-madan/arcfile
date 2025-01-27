package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)


func UploadHandler(c *gin.Context) {
	name := c.PostForm("name")

	email := c.PostForm("email")

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}

	filename := filepath.Base(file.Filename)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":  name,
		"email": email,
		"file":  filename,
	})
}
