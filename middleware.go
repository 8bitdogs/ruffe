package ruffe

type Middleware struct {
	parent  *Middleware
	h       Handler
	OnError func(Context, error) error
}

func NewMiddleware(h Handler) *Middleware {
	return &Middleware{
		h: h,
	}
}

func NewMiddlewareFunc(f func(Context) error) *Middleware {
	return NewMiddleware(HandlerFunc(f))
}

func (m *Middleware) Wrap(h Handler) *Middleware {
	return &Middleware{
		parent:  m,
		h:       h,
		OnError: m.OnError,
	}
}

func (m *Middleware) WrapFunc(f func(Context) error) *Middleware {
	return m.Wrap(HandlerFunc(f))
}

func (m *Middleware) Handle(ctx Context) error {
	if m.parent != nil {
		if err := m.parent.Handle(ctx); err != nil {
			return m.err(ctx, err)
		}
	}
	return m.err(ctx, m.h.Handle(ctx))
}

func (m *Middleware) err(ctx Context, err error) error {
	if m.OnError == nil {
		return err
	}
	return m.OnError(ctx, err)
}
