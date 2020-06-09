package ruffe

import "net/http"

type muxCreator struct{}

func (m muxCreator) Create() Mux {
	return http.NewServeMux()
}
