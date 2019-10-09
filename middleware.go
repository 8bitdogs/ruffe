package ruffe

type Middleware struct {
	mw      *Middleware
	h       Handler
	OnError func(ctx Context, err error) error
}

func newMiddleware(h Handler) *Middleware {
	return &Middleware{
		h: h,
	}
}

func (m *Middleware) Add(h Handler) *Middleware {
	m.mw = &Middleware{
		h: m.h,
	}
	m.h = h
	return m.mw
}

func (m *Middleware) Handle(ctx Context) error {
	if err := m.h.Handle(ctx); err != nil {
		return m.err(ctx, err)
	}
	if m.mw == nil {
		return nil
	}
	return m.err(ctx, m.mw.Handle(ctx))
}

func (m *Middleware) err(ctx Context, err error) error {
	if m.OnError == nil {
		return err
	}
	return m.OnError(ctx, err)
}
