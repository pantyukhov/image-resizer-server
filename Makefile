# Common
NAME ?= resizer
VERSION ?= $(shell git tag --points-at HEAD --sort -version:refname | head -1)
THIS_FILE := $(lastword $(MAKEFILE_LIST))

# Build
BUILD_CMD ?= CGO_ENABLED=0 go build -o /bin/${NAME} -ldflags '-v -w -s' ./cmd/${NAME}
DEBUG_CMD ?= CGO_ENABLED=0 go build -o /bin/${NAME} -gcflags "all=-N -l" ./cmd/${NAME}

.DEFAULT_GOAL := build

.PHONY: golangci
golangci:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run

.PHONY: tests
tests: golangci
	@go test -v ./...

.PHONY: build
build:
	${BUILD_CMD}

.PHONY: build_debug
build_debug:
	${DEBUG_CMD}

