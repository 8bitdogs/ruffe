package ruffe

import "net/http"

type Server struct {
	mux     map[string]*http.ServeMux
	OnError func(ctx Context, err error) error
}

func (s *Server) Handle(pattern, method string, h Handler) *Middleware {
	mux, ok := s.mux[method]
	if !ok {
		mux = http.NewServeMux()
		s.mux[method] = mux
	}
	rt := newMiddleware(h)
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		// TODO: handle content-type
		ctx := &jsonCtx{r: r, w: w}
		err := rt.Handle(ctx)
		if err == nil {
			return
		}
		s.OnError(ctx, err)
	})
	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux, ok := s.mux[r.Method]
	if !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	mux.ServeHTTP(w, r)
}
