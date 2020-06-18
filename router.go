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
	middlewares  *Middleware
	mc           MuxCreator
	mux          map[string]Mux //*http.ServeMux
	interceptors []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	onError      func(Context, error) error
}

// New allocates and returns a new ruffe Router with http.ServeMux
func New() *Router {
	return NewMux(muxCreator{})
}

// NewMux allocates and returns a new ruffe Router with provided mux
func NewMux(mc MuxCreator) *Router {
	return &Router{
		mux:          make(map[string]Mux),
		mc:           mc,
		middlewares:  NewMiddleware(emptyHandler),
		interceptors: []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc){},
	}
}

// Use applies handler which invokes before executing each registered handler
func (s *Router) Use(h Handler) {
	s.middlewares = s.middlewares.Wrap(h)
}

// UseFunc applies handler which invokes before executing each registered handler
func (s *Router) UseFunc(f func(Context) error) {
	s.Use(HandlerFunc(f))
}

// UseAfter applies handler which invokes after executing each registered handler
func (s *Router) UseAfter(h Handler) {
	s.middlewares = s.middlewares.WrapAfter(h)
}

// UseAfterFunc applies handler which invokes after executing each registered handler
func (s *Router) UseAfterFunc(f func(Context) error) {
	s.UseAfter(HandlerFunc(f))
}

// AppendInterceptor adding http.Handler with reference on next interceptor which invokes before ruffe handler
// Warning: Don't forget to call next(w, r) inside interceptor, if it won't be called handler will stop on current executing interceptor
func (s *Router) AppendInterceptor(i func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)) {
	if i == nil {
		return
	}
	s.interceptors = append(s.interceptors, i)
}

// Handle registers the handler for the given pattern with method.
// If a handler already exists for pattern, Handle panics (Only for default Mux).
func (s *Router) Handle(pattern, method string, h Handler) {
	if h == nil {
		return
	}
	mux, ok := s.mux[method]
	if !ok {
		mux = s.mc.Create()
		s.mux[method] = mux
	}

	// apply middlewares
	handler := func(w http.ResponseWriter, r *http.Request) {
		ctx := ContextFromRequest(w, r)
		mw := s.middlewares.Wrap(h)
		mw.OnError = s.onError
		// TODO: this how to handle unhandled error
		// maybe make sense to store it into request context and pass it to interceptors?
		_ = mw.Handle(ctx)
	}

	// apply interceptors
	for i := len(s.interceptors) - 1; i >= 0; i-- {
		h := handler
		itc := s.interceptors[i]
		handler = func(w http.ResponseWriter, r *http.Request) {
			itc(w, r, h)
		}
	}

	mux.HandleFunc(pattern, handler)
}

// HandleFunc registers the handler for the given pattern with method.
// If a handler already exists for pattern, Handle panics (Only for default Mux).
func (s *Router) HandleFunc(pattern, method string, f func(Context) error) {
	s.Handle(pattern, method, HandlerFunc(f))
}

// OnError assign error handler for Route
func (s *Router) OnError(f func(Context, error) error) {
	s.onError = f
}

func (s *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux, ok := s.mux[r.Method]
	if !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	mux.ServeHTTP(w, r)
}
