FROM golang:1.15.6-alpine3.12 AS builder
RUN apk --no-cache add gcc musl-dev git
WORKDIR /go/src/github.com/pottava
ENV APP_PACKAGE="github.com/pottava/http-return-everything" \
    APP_VERSION=v2.0.0
RUN git clone --depth=1 -b "${APP_VERSION}" "https://${APP_PACKAGE}.git"
WORKDIR /go/src/github.com/pottava/http-return-everything
RUN git rev-parse --short HEAD 2>/dev/null > app/githash
COPY --from=ghcr.io/supinf/go-swagger:0.25 /usr/bin/swagger /usr/bin/
RUN swagger generate server -f spec.yaml -t app/generated
WORKDIR /go/src/github.com/pottava/http-return-everything/app
RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 GOOS="${GOOS:-linux}" GOARCH="${GOARCH:-amd64}"; \
    ldflags="-s -w \
      -X \"${APP_PACKAGE}/app/lib.ver=${APP_VERSION}\" \
      -X \"${APP_PACKAGE}/app/lib.commit=$( cat githash )\" \
      -X \"${APP_PACKAGE}/app/lib.date=$( date +%Y-%m-%dT%H:%M:%SZ --utc )\" \
    "; \
    go build -ldflags "${ldflags}" -o /app \
    generated/cmd/return-everything-server/main.go

FROM alpine:3.13.7
RUN apk --no-cache add tini ca-certificates
COPY --from=builder /app /app
ENV HOST=0.0.0.0 PORT=8080
USER nobody
ENTRYPOINT ["tini", "--", "/app"]
