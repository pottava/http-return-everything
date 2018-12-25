FROM mcr.microsoft.com/windows/nanoserver:10.0.14393.2485
ADD https://github.com/pottava/http-return-everything/releases/download/v1.2.1/windows_amd64.exe server.exe
ENTRYPOINT ["C:\\server.exe", "--host", "0.0.0.0", "--port", "8080"]
