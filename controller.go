package upload

import (
	"github.com/gin-gonic/gin"
	"hash"
	"net/http"
)

type FileController struct {
	svc Service
}

type dependency interface {
	Hash() hash.Hash
}

func NewFileController(d dependency) FileController {
	return FileController{
		svc: NewService(d),
	}
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
	req.Req = ctx.Request

	err := c.svc.Download(req, ctx.Writer.Header(), ctx.Writer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}
