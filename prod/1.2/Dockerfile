FROM alpine:3.8

RUN apk --no-cache add tini ca-certificates

ENV API_PORT=8080
ADD app/generated/cmd/return-everything-server/dist/linux_amd64 /app
ADD prod/1.2/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh /app

ENTRYPOINT ["tini", "--", "/entrypoint.sh"]
