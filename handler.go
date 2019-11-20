package ruffe

import "net/http"

var emptyHandler = HandlerFunc(func(Context) error { return nil })

type HandlerFunc func(Context) error

func (h HandlerFunc) Handle(ctx Context) error {
	return h(ctx)
}

type Handler interface {
	Handle(h Context) error
}

type HTTPHandlerFunc func(http.ResponseWriter, *http.Request)

func (h HTTPHandlerFunc) Handle(ctx Context) error {
	h(ctx, ctx.Request())
	return nil
}
