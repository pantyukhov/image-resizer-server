package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
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

// Handle request to file from S3 storage
// If file not found, try select from url
func (f *FileHandler) HandleFile(ctx *gin.Context) {
	file, info, err := f.S3Service.GetOrCreteFile(ctx.Request.URL.Path, true)
	defer func(file *minio.Object) {
		_ = file.Close()
	}(file)

	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	extraHeaders := map[string]string{
		//"Content-Disposition": "attachment; filename=" + info.Key,
	}

	contentType := "application/octet-stream"

	if len(info.ContentType) > 0 {
		contentType = info.ContentType
	}

	ctx.DataFromReader(http.StatusOK, info.Size, contentType, file, extraHeaders)
}
