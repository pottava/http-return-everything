FROM golang:1.15.6-alpine3.12

RUN apk --no-cache add git
RUN adduser -D -u 1000 vscode

ENV GO111MODULE=on \
    ENABLE_AWS=false \
    ENABLE_GCP=false

RUN go get -u "github.com/cosmtrek/air@v1.15.1"
RUN go get -u "github.com/go-delve/delve/cmd/dlv@v1.5.1"
