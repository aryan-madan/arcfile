package handlers

import (
	"errors"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/nxrmqlly/arcfile-backend/storage"
)

// GET /api/file/:identifier/download
func (h *Handlers) FileDownload(c *gin.Context) {
	identifier := c.Param("identifier")

	file, err := h.repo.GetFile(c, identifier)

	var ae *storage.FileNotFoundError
	if err != nil {
		if errors.As(err, &ae) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": err.Error(),
			})
			return
		} else {
			// some other internal error
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error(),
			})
			return
		}
	}
	path := path.Join("data", "uploads", file.UUID)
	c.FileAttachment(path, file.Filename)
}
