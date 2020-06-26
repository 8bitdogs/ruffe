package ruffe

type Middleware struct {
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

// Before create middleware which call Middleware handler before handler from argument
func (m *Middleware) Before(h Handler) *Middleware {
	return &Middleware{
		h: m.mwf(m.h, h),
	}
}

// BeforeFunc do same as Before
func (m *Middleware) BeforeFunc(f func(Context) error) *Middleware {
	return m.Before(HandlerFunc(f))
}

// After create middleware which call Middleware handler after handler from argument
func (m *Middleware) After(h Handler) *Middleware {
	return &Middleware{
		h: m.mwf(h, m.h),
	}
}

// AfterFunc do same as After
func (m *Middleware) AfterFunc(f func(Context) error) *Middleware {
	return m.After(HandlerFunc(f))
}

// Wrap crate middleware from f argument where pass Middleware handler as next argument
func (m *Middleware) Wrap(f func(ctx Context, next Handler) error) *Middleware {
	return &Middleware{
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
