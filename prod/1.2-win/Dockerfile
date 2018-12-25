FROM golang:1.11.4-windowsservercore as build
RUN go get github.com/go-swagger/go-swagger/cmd/swagger
WORKDIR C:\\gopath\\src\\github.com\\go-swagger\\go-swagger
RUN git checkout v0.17.2
RUN go build cmd/swagger/swagger.go
RUN go get github.com/golang/dep/...
WORKDIR C:\\gopath\\src\\github.com\\pottava
RUN git clone https://github.com/pottava/http-return-everything.git
WORKDIR C:\\gopath\\src\\github.com\\pottava\\http-return-everything
RUN C:\\gopath\\src\\github.com\\go-swagger\\go-swagger\\swagger.exe generate server -f spec.yaml -t app/generated
WORKDIR C:\\gopath\\src\\github.com\\pottava\\http-return-everything\\app
RUN dep ensure
RUN go build generated/cmd/return-everything-server/main.go

FROM mcr.microsoft.com/windows/nanoserver:10.0.14393.2485
COPY --from=build C:\\gopath\\src\\github.com\\pottava\\http-return-everything\\app\\main.exe server.exe
ENTRYPOINT ["C:\\server.exe", "--host", "0.0.0.0", "--port", "8080"]
