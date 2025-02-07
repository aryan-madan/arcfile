package handlers

import (
	"fmt"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

// GET /api/file/:identifier/download
func (h *Handlers) FileDownload(c *gin.Context) {
	identifier := c.Param("identifier")
	fmt.Println("Downloading file with identifier:", identifier)

	file, err := h.repo.GetFile(c, identifier)
	fmt.Println("File:", file)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	path := path.Join("data", "uploads", file.UUID)
	fmt.Println("path is ", path)
	c.FileAttachment(path, file.Filename)
}
