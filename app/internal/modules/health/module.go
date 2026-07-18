// Package health exposes the liveness and version endpoints used to check
// that the application is up and to identify which build is running.
package health

import "net/http"

// Module wires together the health module's dependencies and exposes its
// HTTP routes.
type Module struct {
	handler *Handler
}

// New builds a health Module with its handler dependencies initialized.
func New(version string) *Module {
	return &Module{
		handler: NewHandler(version),
	}
}

// RegisterRoutes registers the health module's routes on the given mux.
func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", m.handler.Health)
	mux.HandleFunc("GET /version", m.handler.Version)
}
