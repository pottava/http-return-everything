FROM golang:1.8.3-alpine3.6 AS build-env
RUN apk --no-cache add git
RUN go get github.com/pottava/http-return-everything

FROM alpine:3.6

COPY --from=build-env /go/bin/http-return-everything /app

EXPOSE 80
CMD ["/app"]
