FROM alpine:3.6

ENV API_PORT=8080
ADD entrypoint.sh /

RUN apk --no-cache add tini ca-certificates
RUN apk --no-cache add --virtual build-dependencies bash gcc musl-dev openssl go git \

    # Install go 1.9
    && GOLANG_VERSION=1.9.2 \
    && GOLANG_SRC_URL=https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz \
    && GOLANG_SRC_SHA256=665f184bf8ac89986cfd5a4460736976f60b57df6b320ad71ad4cef53bb143dc \
    && export GOROOT_BOOTSTRAP="$(go env GOROOT)" \
    && wget -q "$GOLANG_SRC_URL" -O golang.tar.gz \
    && echo "$GOLANG_SRC_SHA256  golang.tar.gz" | sha256sum -c - \
    && tar -C /usr/local -xzf golang.tar.gz \
    && wget -q https://raw.githubusercontent.com/docker-library/golang/master/1.9/alpine3.6/no-pic.patch \
    && cd /usr/local/go/src \
    && patch -p2 -i /no-pic.patch \
    && ./make.bash \
    && export PATH=/go/bin:/usr/local/go/bin:$PATH \
    && export GOPATH=/go \
    && mkdir -p /go/src /go/bin \
    && chmod -R 777 /go \

    # Install go-swagger & golang/dep
    && cd /usr/local/go/bin \
    && wget -q https://github.com/go-swagger/go-swagger/releases/download/0.12.0/swagger_linux_amd64 \
    && chmod +x /usr/local/go/bin/swagger_linux_amd64 \
    && go get -u github.com/golang/dep/cmd/dep \

    # Compile the app
    && project_dir=/go/src/github.com/pottava/http-return-everything \
    && git clone https://github.com/pottava/http-return-everything.git ${project_dir} \
    && cd ${project_dir} \
    # && git checkout $(git describe --tags `git rev-list --tags --max-count=1`) \
    && /usr/local/go/bin/swagger_linux_amd64 generate server -f spec.yaml -t app/generated \
    && cd app && /go/bin/dep ensure \
    && go install github.com/pottava/http-return-everything/app/generated/cmd/return-everything-server \
    && mv /go/bin/return-everything-server /app \

    # Clean up
    && apk del --purge -r build-dependencies \
    && rm -rf /usr/local/go /usr/lib/go /go /golang.tar.gz /*.patch

ENTRYPOINT ["/sbin/tini", "--", "/entrypoint.sh"]
