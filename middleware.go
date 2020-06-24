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

func (m *Middleware) Before(h Handler) *Middleware {
	return &Middleware{
		parent: m,
		h:      m.mwf(m.h, h),
	}
}

func (m *Middleware) BeforeFunc(f func(Context) error) *Middleware {
	return m.Before(HandlerFunc(f))
}

// WrapAfter invoke middleware handler after wrapped handler
func (m *Middleware) After(h Handler) *Middleware {
	return &Middleware{
		parent: m,
		h:      m.mwf(h, m.h),
	}
}

func (m *Middleware) AfterFunc(f func(Context) error) *Middleware {
	return m.After(HandlerFunc(f))
}

func (m *Middleware) Wrap(f func(ctx Context, next Handler) error) *Middleware {
	return &Middleware{
		parent: m,
		h: HandlerFunc(func(ctx Context) error {
			if ctx.done() {
				return nil
			}
			return f(ctx, m.h)
		}),
	}
}

func (m *Middleware) mwf(first, last Handler) Handler {
	return HandlerFunc(func(ctx Context) error {
		if err := first.Handle(ctx); err != nil {
			return err
		}
		if ctx.done() {
			return nil
		}
		return last.Handle(ctx)
	})
}

func (m *Middleware) Handle(ctx Context) error {
	err := m.h.Handle(ctx)
	if err == nil {
		return nil
	}
	if m.OnError == nil {
		return err
	}
	return m.OnError(ctx, err)
}
