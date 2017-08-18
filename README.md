# HTTP-Return-Everything

Supported tags and respective `Dockerfile` links:  
・latest ([prod/Dockerfile](https://github.com/pottava/http-return-everything/blob/master/prod/Dockerfile))

## Usage

### 1. Set environment variables

Environment Variables     | Description                                       | Required
------------------------- | ------------------------------------------------- | ---------------------
ACCESS_LOG                | Send access logs to /dev/stdout. (default: true) | 
CONTENT_ENCODING          | Compress response data if the request allows. (default: true) |

### 2. Run the application

`$ docker run -d -p 8080:80 pottava/http-re`

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
