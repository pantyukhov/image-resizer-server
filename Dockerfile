FROM golang:1.15-alpine AS builder

ENV NAME "resizer"
WORKDIR /opt/${NAME}

RUN apk add --no-cache make

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o /bin/${NAME} -gcflags "all=-N -l" ./cmd/${NAME}

FROM alpine:latest
ENV NAME "resizer"

WORKDIR /opt/${NAME}

COPY --from=builder /bin/${NAME} ./${NAME}

CMD ./${NAME}
