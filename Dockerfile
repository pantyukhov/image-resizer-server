FROM golang:1.15-alpine AS builder

ENV NAME "resizer"
WORKDIR /opt/${NAME}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -ldflags '-v -w -s'

RUN go build -o ./bin/${NAME} -ldflags '-v -w -s' ./

FROM alpine:latest

ENV NAME "resizer"
WORKDIR /opt/${NAME}

COPY --from=builder ./bin/${NAME} ./${NAME}

CMD ./${NAME}
