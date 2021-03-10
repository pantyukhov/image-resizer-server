package transport

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pantyukhov/imageresizeserver/pkg/setting"
	"github.com/zsais/go-gin-prometheus"
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

	p := ginprometheus.NewPrometheus("gin")
	p.Use(baseRoute)

	corsConfig := cors.DefaultConfig()
	if len(setting.Settings.CorsConfig.AllowOrigins) > 0 {
		corsConfig.AllowOrigins = setting.Settings.CorsConfig.AllowOrigins
	}

	baseRoute.Use(cors.New(corsConfig))

	baseRoute.GET("/health", HandleHealthCheck)
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
