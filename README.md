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

This service requires

- PostgreSQL to be running and have a `public.todos` which matches that of the Todo struct in `model.go`.
- POSTGRES_URL environment varible set to the postgres url of the PostrgeSQL server above. e.g. postgres://etc...

```
./TodoService
```

TodoService will start a http server on the port specified in by Environment variable `PORT`, which defaults to `8000`.
