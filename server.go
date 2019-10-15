package ruffe

import "net/http"

type Server struct {
	handlers []Handler
	mux      map[string]*http.ServeMux
	OnError  func(ctx Context, err error) error
}

func New() *Server {
	return &Server{
		mux:      make(map[string]*http.ServeMux),
		handlers: make([]Handler, 0),
	}
}

func (s *Server) Use(h Handler) {
	s.handlers = append(s.handlers, h)
}

func (s *Server) UseFunc(f func(Context) error) {
	s.Use(HandlerFunc(f))
}

func (s *Server) Handle(pattern, method string, h Handler) *Middleware {
	mux, ok := s.mux[method]
	if !ok {
		mux = http.NewServeMux()
		s.mux[method] = mux
	}
	mw := NewMiddleware(h)
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		// TODO: context should be created depends on accept and content-type header
		ctx := &jsonCtx{r: r, w: w}
		for _, h := range s.handlers {
			if err := h.Handle(ctx); err != nil {
				s.err(ctx, err)
				return
			}
		}
		err := mw.Handle(ctx)
		if err != nil {
			s.err(ctx, err)
		}
	})
	return mw
}

func (s *Server) HandleFunc(pattern, method string, f func(Context) error) *Middleware {
	return s.Handle(pattern, method, HandlerFunc(f))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux, ok := s.mux[r.Method]
	if !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	mux.ServeHTTP(w, r)
}

func (s *Server) err(ctx Context, err error) error {
	if s.OnError == nil {
		return err
	}
	return s.OnError(ctx, err)
}
