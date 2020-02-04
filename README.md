# Middleware [![Build Status](https://travis-ci.org/clevergo/middleware.svg?branch=master)](https://travis-ci.org/clevergo/middleware) [![Coverage Status](https://coveralls.io/repos/github/clevergo/middleware/badge.svg?branch=master)](https://coveralls.io/github/clevergo/middleware?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/clevergo/middleware)](https://goreportcard.com/report/github.com/clevergo/middleware) [![GoDoc](https://img.shields.io/badge/godoc-reference-blue)](https://pkg.go.dev/github.com/clevergo/middleware) [![Release](https://img.shields.io/github/release/clevergo/middleware.svg?style=flat-square)](https://github.com/clevergo/middleware/releases)

Collections of HTTP middlewares.

## Middlewares

- [**Logging**](https://pkg.go.dev/github.com/clevergo/middleware#Logging)
- [**Compress**](https://pkg.go.dev/github.com/clevergo/middleware#Compress)
- [**Header**](https://pkg.go.dev/github.com/clevergo/middleware#Header)
- [**Basic Auth**](https://pkg.go.dev/github.com/clevergo/middleware#BasicAuth)


## Example

```go
package main

import (
	"compress/gzip"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/middleware"
	"github.com/gorilla/handlers"
)

var users = map[string]string{
	"foo": "bar",
}

func main() {
	app := clevergo.New(":1234")
	app.Use(
		handlers.RecoveryHandler(),
		middleware.Compress(gzip.DefaultCompression),
		middleware.Logging(os.Stdout),
	)
	app.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello world")
	})
	auth := app.Group("/auth", clevergo.RouteGroupMiddleware(
		middleware.BasicAuth(func(username, password string) bool {
			if passwd, exists := users[username]; exists && passwd == password {
				return true
			}
			return false
		}),
	))
	auth.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "authenticated")
	})
	log.Fatal(app.ListenAndServe())
}
```
