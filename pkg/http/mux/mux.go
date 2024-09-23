package mux

import "net/http"

type Mux struct {
	http.Handler
	middlewares []func(http.Handler) http.Handler
}

func New() *Mux {
	return &Mux{
		Handler:     nil,
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
}
