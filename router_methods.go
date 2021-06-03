package ruffe

import "net/http"

const (
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH" // RFC 5789
	MethodDelete  = "DELETE"
	MethodConnect = "CONNECT"
	MethodOptions = "OPTIONS"
	MethodTrace   = "TRACE"
)

func (r *Router) Get(pattern string, h Handler) {
	r.Handle(pattern, http.MethodGet, h)
}

func (r *Router) GetFunc(pattern string, f func(Context) error) {
	r.HandleFunc(pattern, http.MethodGet, f)
}

func (r *Router) Head(pattern string, h Handler) {
	r.Handle(pattern, http.MethodHead, h)
}

func (r *Router) HeadFunc(pattern string, f func(Context) error) {
	r.HandleFunc(pattern, http.MethodHead, f)
}

func (r *Router) Post(pattern string, h Handler) {
	r.Handle(pattern, http.MethodPost, h)
}

func (r *Router) PostFunc(pattern string, f func(Context) error) {
	r.HandleFunc(pattern, http.MethodPost, f)
}

func (r *Router) Put(pattern string, h Handler) {
	r.Handle(pattern, http.MethodPut, h)
}

func (r *Router) PutFunc(pattern string, f func(Context) error) {
	r.HandleFunc(pattern, http.MethodPut, f)
}

func (r *Router) Patch(pattern string, h Handler) {
	r.Handle(pattern, http.MethodPatch, h)
}

func (r *Router) PatchFunc(pattern string, f func(Context) error) {
	r.HandleFunc(pattern, http.MethodPatch, f)
}

func (r *Router) Delete(pattern string, h Handler) {
	r.Handle(pattern, http.MethodDelete, h)
}

func (r *Router) DeleteFunc(pattern string, f func(Context) error) {
	r.HandleFunc(pattern, http.MethodDelete, f)
}
