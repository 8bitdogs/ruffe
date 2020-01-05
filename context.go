package ruffe

import (
	"errors"
	"io"
	"net/http"
)

var (
	ErrResponseWasAlreadySent = errors.New("result was sent")
)

type Context interface {
	done() bool
	http.ResponseWriter
	Request() *http.Request
	Bind(interface{}) error
	Result(int, interface{}) error
}

type requestUnmarshaler interface {
	Unmarshal(r io.Reader, v interface{}) error
}

type responseMarshaler interface {
	ContentType() string
	Marshal(w io.Writer, v interface{}) error
}

type ctx struct {
	isSent bool
	http.ResponseWriter
	r  *http.Request
	ru requestUnmarshaler
	rm responseMarshaler
}

func ContextFromRequest(w http.ResponseWriter, r *http.Request) Context {
	jc := &jsonContent{}
	result := &ctx{
		// TODO: parse Accept header to define correct marshaler
		// Github issue: https://github.com/8bitdogs/ruffe/issues/1
		rm:             jc,
		ru:             emptyUnmarshaler{},
		r:              r,
		ResponseWriter: w,
	}
	if r.ContentLength > 0 {
		// TODO: parse Content-type to define correct unmarshaler
		// currently ruffe supports only json
		// Github issue: https://github.com/8bitdogs/ruffe/issues/2
		result.ru = jc
	}
	return result
}

func (c *ctx) done() bool {
	return c.isSent
}

func (c *ctx) Request() *http.Request {
	return c.r
}

func (c *ctx) Bind(v interface{}) error {
	return c.ru.Unmarshal(c.r.Body, v)
}

func (c *ctx) Result(code int, v interface{}) error {
	if c.isSent {
		return ErrResponseWasAlreadySent
	}
	if v != nil {
		c.Header().Add(ContentTypeHeader, c.rm.ContentType())
	}
	c.WriteHeader(code)
	c.isSent = true
	if v != nil {
		return c.rm.Marshal(c.ResponseWriter, v)
	}
	return nil
}

type emptyUnmarshaler struct{}

func (emptyUnmarshaler) Unmarshal(io.Reader, interface{}) error { return nil }

type emptyMarshaler struct{}

func (emptyMarshaler) ContentType() string { return "" }

func (emptyMarshaler) Marshal(io.Writer, interface{}) error { return nil }
