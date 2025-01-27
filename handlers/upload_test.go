package handlers_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nxrmqlly/arcfile-backend/handlers"
	"github.com/stretchr/testify/assert"
)

func TestUploadHandler(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/upload", handlers.UploadHandler)

	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Write some data to the file
	fileContent := []byte("Hello, World!")
	if _, err := tmpFile.Write(fileContent); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Create a new multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file to the request
	part, err := writer.CreateFormFile("file", filepath.Base(tmpFile.Name()))
	if err != nil {
		t.Fatal(err)
	}
	_, err = part.Write(fileContent)
	if err != nil {
		t.Fatal(err)
	}

	// Add other form fields
	_ = writer.WriteField("email", "test@example.com")

	// Close the writer
	err = writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Create a new request
	req, err := http.NewRequest("POST", "/upload", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)

	// Optionally, you can further assert the response body
	// For example, you can unmarshal the JSON response and check specific fields
	// Here's a simple check for the presence of the "identifier" field
	assert.Contains(t, w.Body.String(), "identifier")
}
