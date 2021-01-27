package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/pantyukhov/imageresizeserver/services"
	"net/http"
)

type FileHandler struct {
	S3Service services.S3Service
}

func NewFileHandler(s3Service services.S3Service) FileHandler {
	return FileHandler{
		S3Service: s3Service,
	}
}

// Hello is the Hello handler
func (f *FileHandler) HandleFile(ctx *gin.Context) {

	file, err := f.S3Service.GetOrCreteFile("test.png")

	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	info, err := file.Stat()
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	extraHeaders := map[string]string{
		"Content-Disposition": "attachment; filename=" + info.Key,
	}
	ctx.DataFromReader(http.StatusOK, info.Size, "application/octet-stream", file, extraHeaders)
}
