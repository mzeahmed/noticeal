// Package router assembles the application's HTTP handler by wiring up the
// routes exposed by each module.
package router

import (
	"log/slog"
	"net/http"

	"github.com/mzeahmed/noticeal/internal/modules/auth"
	"github.com/mzeahmed/noticeal/internal/modules/events"
	"github.com/mzeahmed/noticeal/internal/modules/health"
)

// New builds and returns the application's top-level http.Handler, with all
// module routes registered on a fresh http.ServeMux.
func New(appVersion, authToken string, log *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	health.New(appVersion).RegisterRoutes(mux)
	events.New(log).RegisterRoutes(mux, auth.Authenticate(authToken))

	return mux
}
