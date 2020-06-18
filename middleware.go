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
		parent: m,
		h:      h,
	}
}

func (m *Middleware) WrapFunc(f func(Context) error) *Middleware {
	return m.Wrap(HandlerFunc(f))
}

// WrapAfter invoke middleware handler after wrapped handler
func (m *Middleware) WrapAfter(h Handler) *Middleware {
	return &Middleware{
		parent: m.parent,
		h: HandlerFunc(func(ctx Context) error {
			err := m.err(ctx, h.Handle(ctx))
			if err != nil {
				return err
			}
			return m.Handle(ctx)
		}),
	}
}

func (m *Middleware) WrapAfterFunc(f func(Context) error) *Middleware {
	return m.WrapAfter(HandlerFunc(f))
}

func (m *Middleware) Handle(ctx Context) error {
	if m.parent != nil {
		if err := m.parent.Handle(ctx); err != nil {
			return m.err(ctx, err)
		}
	}
	if ctx.done() {
		return nil
	}
	return m.err(ctx, m.h.Handle(ctx))
}

func (m *Middleware) err(ctx Context, err error) error {
	if err == nil {
		return nil
	}
	if m.OnError == nil {
		return err
	}
	return m.OnError(ctx, err)
}
