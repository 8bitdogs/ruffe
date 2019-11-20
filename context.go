package ruffe

import (
	"encoding/json"
	"net/http"
)

type Context interface {
	http.ResponseWriter
	Request() *http.Request
	Bind(interface{}) error
	Result(int, interface{}) error
}

type jsonCtx struct {
	r *http.Request
	http.ResponseWriter
}

func (c *jsonCtx) Request() *http.Request {
	return c.r
}

func (c *jsonCtx) Bind(v interface{}) error {
	return json.NewDecoder(c.r.Body).Decode(v)
}

func (c *jsonCtx) Result(code int, v interface{}) error {
	c.WriteHeader(code)
	return json.NewEncoder(c).Encode(v)
}
