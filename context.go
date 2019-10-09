package ruffe

import (
	"encoding/json"
	"net/http"
)

type Context interface {
	Request() *http.Request
	Bind(interface{}) error
	Result(int, interface{}) error
}

type jsonCtx struct {
	r *http.Request
	w http.ResponseWriter
}

func (c *jsonCtx) Request() *http.Request {
	return c.r
}

func (c *jsonCtx) Bind(v interface{}) error {
	return json.NewDecoder(c.r.Body).Decode(v)
}

func (c *jsonCtx) Result(code int, v interface{}) error {
	c.w.WriteHeader(code)
	return json.NewEncoder(c.w).Encode(v)
}
