# TodoService

Todo microservice written in [Go](https://golang.org/) and uses [Go-Kit](https://gokit.io/)

## Build

```
// Use Go 1.11 Modules
export GO111MODULE=on;

// Build TodoService binary
go build
```

## Run

```
./TodoService
```

TodoService will start a http server on the port specified in by Environment variable `PORT`, which defaults to `8000`.

## NB!

NOTE: Curently TodoService is an in memory service
