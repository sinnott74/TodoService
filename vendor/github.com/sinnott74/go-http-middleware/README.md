# go-http-middleware

[![Build Status](https://travis-ci.org/sinnott74/go-http-middleware.svg?branch=master)](https://travis-ci.org/sinnott74/go-http-middleware)

Collection of Golang HTTP middlewares fo use with Go's `net/http` package

- [**etag**](https://github.com/sinnott74/go-http-middleware/blob/master/etag.go) Adds [ETag](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/ETag) support for each resource.

- [**https**](https://github.com/sinnott74/go-http-middleware/blob/master/https.go) Forces [X-Forwarded-Proto](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-Proto) header to be set to HTTPS. Useful when behind a load balancer i.e. aws, cloudfoundry, etc.

- [**auth**](https://github.com/sinnott74/go-http-middleware/blob/master/auth.go) handles authenication using a user supplied authenication function.

- [**JWT**](https://github.com/sinnott74/go-http-middleware/blob/master/JWT.go) handles JWT authenication which allows a user supplied validation function.

## Installation

`go get https://github.com/sinnott74/go-http-middleware`

## Example Usage

```go
package main

import (
	"net/http"

	"github.com/sinnott74/go-http-middleware"
)

func main() {
	http.Handle("/", middleware.DefaultEtag(helloWorldHandler()))
	http.ListenAndServe(":8080", nil)
}

// Simple hello world handler
func helloWorldHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})
}

// Visiting localhost:8080/ returns a response body of "Hello, world" and an ETag header 'W/"d-ZajifYh5KDgxtmS9i38K1A=="'
```
