FROM golang:1.15-alpine AS builder

WORKDIR /opt/resizer

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -ldflags '-v -w -s'

RUN go build -o /bin/resizer

FROM alpine:latest

WORKDIR /opt/resizer

COPY --from=builder /bin/resizer ./resizer

CMD ./resizer
