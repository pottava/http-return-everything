# HTTP-Return-Everything

[![CircleCI](https://circleci.com/gh/pottava/http-return-everything.svg?style=svg)](https://circleci.com/gh/pottava/http-return-everything)

Supported tags and respective `Dockerfile` links:  
・latest ([prod/Dockerfile](https://github.com/pottava/http-return-everything/blob/master/prod/Dockerfile))

## Usage

### 1. Set environment variables

Environment Variables     | Description                                       |
------------------------- | ------------------------------------------------- |
ACCESS_LOG                | Send access logs to /dev/stdout. (default: true) | 
ACCESS_DETAIL_LOG         | Save HTTP request details (default: false) | 
CONTENT_ENCODING          | Compress response data if the request allows. (default: true) |

### 2. Run the application

`$ docker run -d -p 80:8080 pottava/http-re:1.0`

* with docker-compose.yml:  

```
check:
  image: pottava/http-re:1.0
  ports:
    - 80:8080
  environment:
    - ACCESS_LOG=false
    - CONTENT_ENCODING=false
  container_name: check
```

* with kubernetes-deployment.yaml

```
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: my-app
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
      - name: api
        image: pottava/http-re:1.0
        imagePullPolicy: Always
        ports:
        - protocol: TCP
          containerPort: 8080
        env:
        - name: APP_NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
        - name: APP_POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: APP_POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: APP_POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        - name: APP_POD_SERVICE_ACCOUNT
          valueFrom:
            fieldRef:
              fieldPath: spec.serviceAccountName
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
```

### 3. Make HTTP GET requests

- GET /

<img alt="" src="https://raw.githubusercontent.com/wiki/pottava/http-return-everything/images/everything.png" style="max-width: 100%;">

- GET /app/envs/

<img alt="" src="https://raw.githubusercontent.com/wiki/pottava/http-return-everything/images/app-envs.png" style="max-width: 100%;">

- GET /app/envs/{env}

<img alt="" src="https://raw.githubusercontent.com/wiki/pottava/http-return-everything/images/app-envs-key.png" style="max-width: 100%;">

- GET /req/

<img alt="" src="https://raw.githubusercontent.com/wiki/pottava/http-return-everything/images/req.png" style="max-width: 100%;">

- GET /req/headers/

<img alt="" src="https://raw.githubusercontent.com/wiki/pottava/http-return-everything/images/req-headers.png" style="max-width: 100%;">

- GET /req/headers/{header}

<img alt="" src="https://raw.githubusercontent.com/wiki/pottava/http-return-everything/images/req-headers-key.png" style="max-width: 100%;">
