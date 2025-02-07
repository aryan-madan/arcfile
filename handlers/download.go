package handlers

import (
	"fmt"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

// GET /api/download/:identifier
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

	filePath := path.Join("data", "uploads", file.UUID)

	fmt.Println("path is ", filePath)

	c.FileAttachment(filePath, file.Filename)

	err = h.repo.DelteFile(c, identifier, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		fmt.Println("ln46: ", err.Error())
		return
	}

}
