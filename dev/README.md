# How to setup local environment

## Generate source code from the OpenAPI spec file

```
$ docker-compose -f dev/tools.yml run --rm codegen
```

## Vendoring go dependencies

```
$ docker-compose -f dev/tools.yml run --rm vendor
```

## To edit OpenAPI spec

```
$ docker-compose -f dev/openapi.yml up
```

## Run the application as a Cloud Run service

According to [Debugging your Cloud Run service](https://cloud.google.com/code/docs/vscode/debugging-a-cloud-run-app), you can get a '200 OK' response from `http://localhost:8080`.
