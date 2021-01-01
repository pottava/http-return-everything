# Retrieving server context - A REST API server

[![GitHub Actions](https://github.com/pottava/http-return-everything/workflows/Publish%20artifacts/badge.svg?branch=master)](https://github.com/pottava/http-return-everything/actions)

[gcr.io/pottava/re](https://gcr.io/pottava/re/)

Supported tags and respective `Dockerfile` links:  
・v2.0 ([prod/2.0/Dockerfile](https://github.com/pottava/http-return-everything/blob/master/prod/2.0/Dockerfile))  
・v1.3 ([prod/1.3/Dockerfile](https://github.com/pottava/http-return-everything/blob/master/prod/1.3/Dockerfile))  

## Usage

### 1. Set environment variables

Environment Variables     | Description                                       |
------------------------- | ------------------------------------------------- |
PORT                      | Listening port. (default: 8080) | 
ENABLE_AWS                | Enable the AWS metadata endpoint. (default: true) | 
ENABLE_GCP                | Enable the Google Cloud metadata endpoint. (default: true) | 
ACCESS_LOG                | Send access logs to /dev/stdout. (default: true) | 
ACCESS_DETAIL_LOG         | Save HTTP request details (default: false) | 
CONTENT_ENCODING          | Compress response data if the request allows. (default: true) |
CORS_ORIGIN               | Allowed CORS origin (default: *) |

### 2. Run the application

`$ docker run -d --rm -p 80:8080 gcr.io/pottava/re:v2.0`

* with Google [Cloud Run](https://cloud.google.com/run):  

```bash
$ gcloud run deploy re --allow-unauthenticated \
    --image gcr.io/pottava/re:v2.0 \
    --set-env-vars ENABLE_GCP=1,ENABLE_AWS=0
```

* with docker-compose.yml:  

```yaml
check:
  image: gcr.io/pottava/re:v2.0
  ports:
    - 80:8080
  environment:
    - ENABLE_AWS=false
    - ENABLE_GCP=false
    - ACCESS_LOG=false
    - CONTENT_ENCODING=false
  container_name: check
```

* with Kubernetes deployment.yaml

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  replicas: 1
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
      - name: api
        image: gcr.io/pottava/re:v2.0
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

[API spec](https://github.com/pottava/http-return-everything/blob/master/spec.yaml)

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
