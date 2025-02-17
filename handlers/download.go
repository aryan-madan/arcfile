package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/nxrmqlly/arcfile/storage"
)

// GET /api/download/:identifier
func (h *Handlers) FileDownload(c *gin.Context) {
	identifier := c.Param("identifier")
	fmt.Println("Downloading file with identifier:", identifier)

	file, err := h.repo.GetFile(c, identifier)
	if err != nil {
		var fileNotFoundErr *storage.FileNotFoundError
		if errors.As(err, &fileNotFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Internal server error",
			})
			fmt.Println("GetFile error:", err.Error())
		}
		return
	}

	filePath := path.Join("data", "uploads", file.UUID)
	fmt.Println("File path:", filePath)

	c.FileAttachment(filePath, file.Filename)

	if err := h.repo.DeleteFile(c, identifier, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to delete file",
		})
		fmt.Println("DeleteFile error:", err.Error())
		return
	}
}
