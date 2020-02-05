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

var basicAuthValidate = func(username, password string) bool {
	if passwd, exists := users[username]; exists && passwd == password {
		return true
	}
	return false
}

var basicAuthErrorHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "invalid credential", http.StatusUnauthorized)
})

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

	// basic auth
	basicAuth := middleware.BasicAuth(
		basicAuthValidate,
		middleware.BasicAuthRealm("Restricted"),                 // optional.
		middleware.BasicAuthErrorHandler(basicAuthErrorHandler), // optional.
	)
	auth := app.Group("/auth", clevergo.RouteGroupMiddleware(basicAuth))
	auth.Get("/", func(w http.ResponseWriter, r *http.Request) {
		user := middleware.GetBasicAuthUser(r)
		fmt.Fprintf(w, "hello %s", user)
	})

	log.Fatal(app.ListenAndServe())
}
```

```shell
$ curl -v http://localhost:1234/auth/
...
< HTTP/1.1 401 Unauthorized
< Www-Authenticate: Basic realm="Restricted"
invalid credential

$ curl -v -u foo:bar http://localhost:1234/auth/
...
> Authorization: Basic Zm9vOmJhcg==
...
< Www-Authenticate: Basic realm="Restricted"
hello foo
```
