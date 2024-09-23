package mux

import "net/http"

func (m *mux) Handle(path string, handler http.Handler) {
	if m.emptyHandler {
		m.emptyHandler = false
	}

	wrapped := handler
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		wrapped = m.middlewares[i](wrapped)
	}
	m.ServeMux.Handle(path, wrapped)
}

func (m *mux) Get(path string, handler http.HandlerFunc) {
	m.Handle("GET "+path, http.HandlerFunc(handler))
}

func (m *mux) Post(path string, handler http.HandlerFunc) {
	m.Handle("POST "+path, http.HandlerFunc(handler))
}

func (m *mux) Put(path string, handler http.HandlerFunc) {
	m.Handle("PUT "+path, http.HandlerFunc(handler))
}

func (m *mux) Patch(path string, handler http.HandlerFunc) {
	m.Handle("PATCH "+path, http.HandlerFunc(handler))
}

func (m *mux) Delete(path string, handler http.HandlerFunc) {
	m.Handle("DELETE "+path, http.HandlerFunc(handler))
}

// Use adds middleware to the Mux.
func (m *mux) Use(mw func(http.Handler) http.Handler) {
	if m.emptyHandler {
		panic("mux: no handler registered")
	}
	m.middlewares = append(m.middlewares, mw)
}
