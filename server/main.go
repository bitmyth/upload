package main

import (
	"github.com/bitmyth/upload"
	"github.com/bitmyth/upload/registry"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	upload.InitDir()

	r := gin.Default()

	var rg registry.Registry
	fileController := upload.NewFileController(rg)

	r.POST("/uploads", fileController.NewUpload)
	r.POST("/uploads/:id/:chunk", fileController.UploadChunk)
	r.POST("/uploads/:id", fileController.Reassemble)
	r.GET("/file/:id", fileController.Download)

	// Start the server on port 9090
	err := http.ListenAndServe(":9090", r)
	if err != nil {
		log.Fatal("Error starting server")
	}
}
