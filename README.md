
# Ruffe 
![Alt text](ruffe.png?raw=true "Ruffe")

Golang HTTP handler

## Guide

### Examples
#### 1. Create ruffe instance, add http handler and start it with http.ListenAndServer
```go
package main

import (
	"net/http"

	"github.com/8bitdogs/ruffe"
)

func main() {
	// Ruffe instance
	rs := ruffe.New()

	// add handler
	rs.HandleFunc("/", http.MethodGet, hello)

	// Start server
	http.ListenAndServe(":8080", rs)
}

// hello handler
func hello(ctx ruffe.Context) error {
	return ctx.Result(http.StatusOK, "hello world")
}
```
#### Middleware
__handler middleware__
```go
// Ruffe Middleware
mw := ruffe.NewMiddlewareFunc(func(_ ruffe.Context) error {
    // this handler will occurs before `hello handler`
    return nil
})

// Wrap hello handler with mw middleware
mwh := mw.WrapFunc(hello) // WrapFunc returns middleware

// add handler
_ = rs.Handle("/", http.MethodGet, mwh) // Handle returns middleware
```
__server middleware__
```go
// Applying middleware for all handler
rs.UseFunc(func(_ ruffe.Context) error {
    // server middleware calling before handler
    return nil
})

// add handler
rs.HandleFunc("/", http.MethodGet, hello)
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
	rs := ruffe.New()

	// Define error handler
	rs.OnError = func(_ ruffe.Context, err error) error {
		if err == Err {
			// Caught!
			return nil
		}
		return nil
	}

	// add handler
	rs.HandleFunc("/", http.MethodGet, hello)

	// Start server
	http.ListenAndServe(":8080", rs)
}

// hello handler
func hello(_ ruffe.Context) error {
	return Err
}
```