package main

import (
	_ "go.uber.org/automaxprocs" // optimization for k8s
)

func main() {
	transport := InitializeTransport()
	transport.InitHttp()
}
