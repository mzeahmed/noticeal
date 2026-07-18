// Package events receives notification events over HTTP.
package events

import (
	"log/slog"
	"net/http"
)

// Module wires together the events module's dependencies and exposes its
// HTTP routes.
type Module struct {
	handler *Handler
}

// New builds an events Module with its handler dependencies initialized.
func New(log *slog.Logger) *Module {
	return &Module{
		handler: NewHandler(log),
	}
}

// RegisterRoutes registers the events module's routes on the given mux.
//
// authenticate guards the route, requiring a valid bearer token; the
// caller (see router.New) is expected to pass auth.Authenticate(token).
func (m *Module) RegisterRoutes(mux *http.ServeMux, authenticate func(http.Handler) http.Handler) {
	mux.Handle("POST /api/v1/events", authenticate(http.HandlerFunc(m.handler.Create)))
}
