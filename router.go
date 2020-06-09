package ruffe

import (
	"net/http"
)

type MuxCreator interface {
	Create() Mux
}

type Mux interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Router struct {
	middlewares *Middleware
	mc          MuxCreator
	mux         map[string]Mux //*http.ServeMux
}

func New() *Router {
	return NewMux(muxCreator{})
}

func NewMux(mc MuxCreator) *Router {
	return &Router{
		mux:         make(map[string]Mux),
		mc:          mc,
		middlewares: NewMiddleware(emptyHandler),
	}
}

func (s *Router) Use(h Handler) {
	s.middlewares = s.middlewares.Wrap(h)
}

func (s *Router) UseFunc(f func(Context) error) {
	s.Use(HandlerFunc(f))
}

func (s *Router) Handle(pattern, method string, h Handler) {
	if h == nil {
		return
	}
	mux, ok := s.mux[method]
	if !ok {
		mux = s.mc.Create()
		s.mux[method] = mux
	}
	h = s.middlewares.Wrap(h)
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx := ContextFromRequest(w, r)
		h.Handle(ctx)
	})
}

func (s *Router) HandleFunc(pattern, method string, f func(Context) error) {
	s.Handle(pattern, method, HandlerFunc(f))
}

func (s *Router) OnError(f func(Context, error) error) {
	s.middlewares.OnError = f
}

func (s *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux, ok := s.mux[r.Method]
	if !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	mux.ServeHTTP(w, r)
}
