VERSION 0.7
FROM golang:1.20-alpine3.17
WORKDIR /go-workdir

deps:
    COPY go.mod go.sum ./
    COPY ./*.go ./
    COPY ./cmd/instances/main.go ./cmd/instances

build:
    FROM +deps
    RUN go build -v

test:
    FROM +build
    RUN go test -v
