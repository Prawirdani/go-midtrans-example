package mux

import "net/http"

type mux struct {
	*http.ServeMux
	middlewares  []func(http.Handler) http.Handler
	emptyHandler bool
}

func New() *mux {
	return &mux{
		ServeMux:     http.NewServeMux(),
		emptyHandler: true,
		middlewares:  make([]func(http.Handler) http.Handler, 0),
	}
}
