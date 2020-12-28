# How to setup local environment

## Generate source code from the OpenAPI spec file

```
$ docker-compose -f dev/tools.yml run --rm codegen
```

## Run the application

```
$ docker-compose -f dev/runtime.yml up
```

Now you can get a '200 OK' response from `http://localhost:9000` .
