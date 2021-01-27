// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/pantyukhov/imageresizeserver/services"
	"github.com/pantyukhov/imageresizeserver/transport"
)

func InitializeTransport() transport.Transport {
	wire.Build(transport.NewTransport, transport.NewFileHandler, services.NewS3Service)
	return transport.Transport{}
}