package main

import (
	"github.com/bitmyth/upload"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func init() {
	err := os.MkdirAll("./uploads/chunks", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := gin.Default()

	var fileController upload.FileController

	r.POST("/uploads", fileController.NewUpload)
	r.POST("/uploads/:id/:chunk", fileController.UploadChunk)
	r.POST("/uploads/:id", fileController.Reassemble)
	r.GET("/file/:id", fileController.Download)

	// Start the server on port 8080
	err := http.ListenAndServe(":9090", r)
	if err != nil {
		log.Fatal("Error starting server")
	}
}
