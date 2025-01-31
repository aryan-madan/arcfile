package main

import (
	"bytes"
	"log"
	"mime/multipart"
	"net/http"
)

func main() {
	for {
		buf := new(bytes.Buffer)
		w := multipart.NewWriter(buf)
		w.WriteField("email", "test")
		fw, err := w.CreateFormFile("file", "test.txt")
		if err != nil {
			panic(err)
		}
		fw.Write([]byte("hello, world!"))
		w.Close()
		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/upload", buf)
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", w.FormDataContentType())
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Printf("unexpected status code: %d", resp.StatusCode)
		}
	}
}
