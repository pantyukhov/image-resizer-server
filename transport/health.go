package transport

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleHealthCheck(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Ok")
}