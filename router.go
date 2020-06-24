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
	pathPrefix   string
	head         *Middleware
	tail         *Middleware
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
		head:         NewMiddleware(emptyHandler),
		tail:         NewMiddleware(emptyHandler),
		interceptors: []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc){},
	}
}

// Use applies handler which invokes before executing each registered handler
func (r *Router) Use(h Handler) {
	r.head = r.head.Wrap(h)
}

// UseFunc applies handler which invokes before executing each registered handler
func (r *Router) UseFunc(f func(Context) error) {
	r.Use(HandlerFunc(f))
}

// UseAfter applies handler which invokes after executing each registered handler
func (r *Router) UseAfter(h Handler) {
	r.tail = r.tail.Wrap(h)
}

// UseAfterFunc applies handler which invokes after executing each registered handler
func (r *Router) UseAfterFunc(f func(Context) error) {
	r.UseAfter(HandlerFunc(f))
}

// AppendInterceptor adding http.Handler with reference on next interceptor which invokes before ruffe handler
// Warning: Don't forget to call next(w, r) inside interceptor, if it won't be called handler will stop on current executing interceptor
func (r *Router) AppendInterceptor(i func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)) {
	if i == nil {
		return
	}
	r.interceptors = append(r.interceptors, i)
}

// Handle registers the handler for the given pattern with method.
// If a handler already exists for pattern, Handle panics (Only for default Mux).
func (r *Router) Handle(pattern, method string, h Handler) {
	if h == nil {
		panic("handler cannot be nil")
	}

	mux, ok := r.mux[method]
	if !ok {
		mux = r.mc.Create()
		r.mux[method] = mux
	}

	// apply middlewares
	handler := func(w http.ResponseWriter, rq *http.Request) {
		ctx := ContextFromRequest(w, rq)
		mw := r.tail.WrapAfter(r.head.Wrap(h))
		mw.OnError = r.onError
		// TODO: this how to handle unhandled error
		// maybe make sense to store it into request context and pass it to interceptors?
		_ = mw.Handle(ctx)
	}

	// apply interceptors
	for i := len(r.interceptors) - 1; i >= 0; i-- {
		h := handler
		itc := r.interceptors[i]
		handler = func(w http.ResponseWriter, rq *http.Request) {
			itc(w, rq, h)
		}
	}

	mux.HandleFunc(r.pathPrefix+pattern, handler)
}

// HandleFunc registers the handler for the given pattern with method.
// If a handler already exists for pattern, Handle panics (Only for default Mux).
func (r *Router) HandleFunc(pattern, method string, f func(Context) error) {
	r.Handle(pattern, method, HandlerFunc(f))
}

// OnError assign error handler for Route
func (r *Router) OnError(f func(Context, error) error) {
	r.onError = f
}

// Subrouter creates copy of router
// Handle and HandleFunc will register handlers to parent router mux
func (r *Router) Subrouter(pathPrefix string) *Router {
	interceptors := make([]func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc), len(r.interceptors))
	for i := range r.interceptors {
		interceptors[i] = r.interceptors[i]
	}
	return &Router{
		pathPrefix:   pathPrefix,
		head:         r.head,
		tail:         r.tail,
		mc:           r.mc,
		mux:          r.mux,
		interceptors: interceptors,
		onError:      r.onError,
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	mux, ok := r.mux[rq.Method]
	if !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	mux.ServeHTTP(w, rq)
}
