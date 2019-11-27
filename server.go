package ruffe

import "net/http"

type Server struct {
	middlewares *Middleware
	mux         map[string]*http.ServeMux
}

func New() *Server {
	return &Server{
		mux:         make(map[string]*http.ServeMux),
		middlewares: NewMiddlewareFunc(func(Context) error { return nil }),
	}
}

func (s *Server) Use(h Handler) {
	s.middlewares = s.middlewares.Wrap(h)
}

func (s *Server) UseFunc(f func(Context) error) {
	s.Use(HandlerFunc(f))
}

func (s *Server) Handle(pattern, method string, h Handler) {
	if h == nil {
		return
	}
	mux, ok := s.mux[method]
	if !ok {
		mux = http.NewServeMux()
		s.mux[method] = mux
	}
	h = s.middlewares.Wrap(h)
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx := ContextFromRequest(w, r)
		h.Handle(ctx)
	})
}

func (s *Server) HandleFunc(pattern, method string, f func(Context) error) {
	s.Handle(pattern, method, HandlerFunc(f))
}

func (s *Server) OnError(f func(Context, error) error) {
	s.middlewares.OnError = f
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux, ok := s.mux[r.Method]
	if !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	mux.ServeHTTP(w, r)
}
