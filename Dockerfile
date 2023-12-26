FROM golang:1.21-alpine

RUN go install golang.org/x/tools/cmd/goimports@latest && \
    apk update && \
    apk add git vim curl

WORKDIR /workdir
