package mux

import "net/http"

func (m *Mux) Handle(path string, handler http.Handler) {
	wrapped := handler
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		wrapped = m.middlewares[i](wrapped)
	}
	m.Handler = wrapped
}

func (m *Mux) Get(path string, handler http.HandlerFunc) {
	m.Handle("GET "+path, http.HandlerFunc(handler))
}

func (m *Mux) Post(path string, handler http.HandlerFunc) {
	m.Handle("POST "+path, http.HandlerFunc(handler))
}

func (m *Mux) Put(path string, handler http.HandlerFunc) {
	m.Handle("PUT "+path, http.HandlerFunc(handler))
}

func (m *Mux) Patch(path string, handler http.HandlerFunc) {
	m.Handle("PATCH "+path, http.HandlerFunc(handler))
}

func (m *Mux) Delete(path string, handler http.HandlerFunc) {
	m.Handle("DELETE "+path, http.HandlerFunc(handler))
}

// Use adds middleware to the Mux.
func (m *Mux) Use(mw func(http.Handler) http.Handler) {
	if m.Handler != nil {
		panic("mux: middleware must be added before any routes")
	}
	m.middlewares = append(m.middlewares, mw)
}
