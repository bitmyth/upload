package upload

import (
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

type FileController struct {
	svc UploadService
}

func (c *FileController) NewUpload(ctx *gin.Context) {
	upload := c.svc.CreateUpload()
	ctx.JSON(http.StatusOK, upload)
}

func (c *FileController) UploadChunk(ctx *gin.Context) {
	chunk, err := c.svc.UploadChunk(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, chunk)
}

func (c *FileController) Reassemble(ctx *gin.Context) {
	var req ReassembleChunksRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.UploadId = ctx.Param("id")

	file, err := c.svc.ReassembleChunk(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Return a success message and the file metadata
	ctx.JSON(http.StatusOK, file)
}

func (c *FileController) Download(ctx *gin.Context) {
	var req DownloadRequest
	req.UploadId = ctx.Param("id")

	err := c.svc.Download(req, ctx.Writer.Header(), ctx.Writer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}
