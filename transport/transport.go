package transport

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"log"
)

type Transport struct {
	FileHandler FileHandler
}

func NewTransport(fileHandler FileHandler) Transport {
	return Transport{
		FileHandler: fileHandler,
	}
}


func (t *Transport) InitFastHTTP() {

	r := router.New()
	r.GET("/{filepath:*}", t.FileHandler.HandleFile)
	log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
}