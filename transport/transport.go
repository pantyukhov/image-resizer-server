package transport

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Transport struct {
	FileHandler FileHandler
}

func NewTransport(fileHandler FileHandler) Transport {
	return Transport{
		FileHandler: fileHandler,
	}
}

func (t *Transport) InitHttp() {
	baseRoute := gin.New()
	baseRoute.Use(gin.Logger())
	baseRoute.Use(gin.Recovery())
	baseRoute.NoRoute(t.FileHandler.HandleFile)

	baseEndPoint := ":8080"
	maxHeaderBytes := 1 << 20

	baseServer := &http.Server{
		Addr:           baseEndPoint,
		Handler:        baseRoute,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info] start http server listening %s", baseEndPoint)
	err := baseServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
