package ruffe

var emptyHandler = HandlerFunc(func(Context) error { return nil })

type HandlerFunc func(Context) error

func (h HandlerFunc) Handle(ctx Context) error {
	return h(ctx)
}

type Handler interface {
	Handle(h Context) error
}
