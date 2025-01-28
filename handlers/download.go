package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nxrmqlly/arcfile-backend/storage"
)

// GET /api/file/:identifier/download
func (h *Handlers) FileDownload(c *gin.Context) {
	identifier := c.Param("identifier")

	file, err := h.repo.GetFile(c.Request.Context(), identifier)

	var ae *storage.FileNotFoundError
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
	}
	c.FileAttachment(file.UUID, file.Filename)
}
