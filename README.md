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
        version: v0.1.0
    spec:
      containers:
      - name: api
        image: pottava/http-re
        imagePullPolicy: IfNotPresent
        ports:
        - protocol: TCP
          containerPort: 80
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
            path: /version
            port: 80
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
