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
		h:      m.mwf(m.h, h),
	}
}

func (m *Middleware) WrapFunc(f func(Context) error) *Middleware {
	return m.Wrap(HandlerFunc(f))
}

// WrapAfter invoke middleware handler after wrapped handler
func (m *Middleware) WrapAfter(h Handler) *Middleware {
	return &Middleware{
		parent: m,
		h:      m.mwf(h, m.h),
	}
}

func (m *Middleware) mwf(first, last Handler) Handler {
	return HandlerFunc(func(ctx Context) error {
		if err := m.err(ctx, first.Handle(ctx)); err != nil {
			return err
		}
		if ctx.done() {
			return nil
		}
		return m.err(ctx, last.Handle(ctx))
	})
}

func (m *Middleware) WrapAfterFunc(f func(Context) error) *Middleware {
	return m.WrapAfter(HandlerFunc(f))
}

func (m *Middleware) Handle(ctx Context) error {
	return m.h.Handle(ctx)
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
