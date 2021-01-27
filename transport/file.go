package transport

import (
	"fmt"
	"github.com/pantyukhov/imageresizeserver/services"
	"github.com/valyala/fasthttp"
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
func (f *FileHandler) HandleFile(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "hello, %s!\n", ctx.UserValue("filepath"))
}


