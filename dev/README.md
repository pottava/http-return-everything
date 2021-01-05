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

Reopen VSCode [in a container](https://code.visualstudio.com/docs/remote/containers).

```
$ cd app
$ air -c .air.toml
```

Now you can get a '200 OK' response from `http://localhost:8080`.
