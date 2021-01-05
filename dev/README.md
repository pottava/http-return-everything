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

## Run the application

```
$ docker-compose -f dev/runtime.yml up
```
