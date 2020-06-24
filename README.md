
# Ruffe
[![Go Report](https://goreportcard.com/badge/github.com/8bitdogs/ruffe)](https://goreportcard.com/report/github.com/8bitdogs/ruffe)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/8bitdogs/ruffe)
![License](https://img.shields.io/github/license/8bitdogs/ruffe)
![Tag](https://img.shields.io/github/v/tag/8bitdogs/ruffe)

![Ruffe preview](ruffe.png?raw=true "Ruffe")

Golang HTTP handler
				
- [Installing](#installing)	
- [Router](#router)
	- [Customization](#customization)
- [Interceptors](#interceptors)
- [Middlewares](#middlewares)
	- [Handler middleware](#handler-middleware)
	- [Router middleware](#router-middleware)
- [Error Handling](#error-handling)

## Guide
### Installing
```
go get -u github.com/8bitdogs/ruffe
```

### Router
```go
package main

import (
	"net/http"

	"github.com/8bitdogs/ruffe"
)

func main() {
	// Ruffe instance
	rr := ruffe.New()

	// add handler
	rr.HandleFunc("/", http.MethodGet, hello)

	// Start server
	http.ListenAndServe(":3030", rs)
}

// hello handler
func hello(ctx ruffe.Context) error {
	return ctx.Result(http.StatusOK, "hello world")
}
```
#### Customization
Using custom request router
```go
package main

import (
	"net/http"

	"github.com/8bitdogs/ruffe"
	"github.com/gorilla/mux"
)

type Router struct {
	*mux.Router
}

// We have to override gorilla Handle and HandleFunc, because those two functions are returning gorilla Router instance

func (r *Router) Handle(pattern string, handler http.Handler) {
	r.Router.Handle(pattern, handler)
}

func (r *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	r.Router.HandleFunc(pattern, handler)
}

type muxCreator struct{}

func (muxCreator) Create() ruffe.Mux {
	return &Router{
		Router: mux.NewRouter(),
	}
}

func main() {
	r := ruffe.NewMux(muxCreator{})
	r.HandleFunc("/foo/{id}", "GET", func(ctx ruffe.Context) error {
		// as you can see, gorilla mux features are available :) 
		return ctx.Result(http.StatusOK, "bar"+mux.Vars(ctx.Request())["id"])
	})
	http.ListenAndServe(":3030", r)
}
```

### Interceptors
example with logging interceptor
```go
rr := ruffe.New()

// adding interceptor
rr.AppendInterceptor(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Println(r.URL, r.Method, r.Header)
	next(w, r)
	log.Println("done")
})

// ... handlers registration

http.ListenAndServe(":3030", rr)
```

### Middlewares
#### Handler middleware
```go
// Initializing Ruffe Middleware
// Middleware implements ruffe.Handler interface
mw := ruffe.NewMiddlewareFunc(func(_ ruffe.Context) error {
	// Middleware logic
	return nil
})

// Add middleware handler before calling <ruffe handler> 
mwh := mw.Wrap(<ruffe handler>) // WrapFunc returns middleware

// Add middleware handler after calling <ruffe handler> 
mwh := mw.WrapAfter(<ruffe handler>) // WrapAfterFunc returns middleware
```
#### Router middleware
```go
rr := ruffe.New()

// applies handler which invokes before executing each registered handler
rr.Use(<ruffe handler>)

// applies handler which invokes after executing each registered handler
rr.UseAfter(<ruffe handler>)
```
#### Error handling
```go
package main

import (
	"errors"
	"net/http"

	"github.com/8bitdogs/ruffe"
)

var Err = errors.New("error")

func main() {
	// Ruffe instance
	rr := ruffe.New()

	// Define error handler
	rr.OnError = func(_ ruffe.Context, err error) error {
		if err == Err {
			// Caught!
			return nil
		}
		return nil
	}

	// add handler
	rr.HandleFunc("/", http.MethodGet, hello)

	// Start server
	http.ListenAndServe(":3030", rs)
}

// hello handler
func hello(_ ruffe.Context) error {
	return Err
}
```
