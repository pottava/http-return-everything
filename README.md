# HTTP-Return-Everything

Supported tags and respective `Dockerfile` links:  
・latest ([prod/Dockerfile](https://github.com/pottava/http-return-everything/blob/master/prod/Dockerfile))

## Usage

### 1. Set environment variables

Environment Variables     | Description                                       |
------------------------- | ------------------------------------------------- |
ACCESS_LOG                | Send access logs to /dev/stdout. (default: true) | 
CONTENT_ENCODING          | Compress response data if the request allows. (default: true) |

### 2. Run the application

`$ docker run -d -p 80:80 pottava/http-re`

* with docker-compose.yml:  

```
check:
  image: pottava/http-re
  ports:
    - 80:80
  environment:
    - ACCESS_LOG=false
  container_name: check
```

### 3. Make HTTP requests

- GET /

<img alt="" src="https://raw.githubusercontent.com/wiki/pottava/http-return-everything/images/result.png" style="max-width: 100%;">

- GET /envs/

<img alt="" src="https://raw.githubusercontent.com/wiki/pottava/http-return-everything/images/envs.png" style="max-width: 100%;">

- GET /envs/key

<img alt="" src="https://raw.githubusercontent.com/wiki/pottava/http-return-everything/images/envs-key.png" style="max-width: 100%;">

- GET /request/

<img alt="" src="https://raw.githubusercontent.com/wiki/pottava/http-return-everything/images/request.png" style="max-width: 100%;">

- GET /request/headers/

<img alt="" src="https://raw.githubusercontent.com/wiki/pottava/http-return-everything/images/request-headers.png" style="max-width: 100%;">

- GET /request/headers/key

<img alt="" src="https://raw.githubusercontent.com/wiki/pottava/http-return-everything/images/request-headers-key.png" style="max-width: 100%;">
