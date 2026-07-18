// Package router assembles the application's HTTP handler by wiring up the
// routes exposed by each module.
package router

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/mzeahmed/noticoel/internal/modules/auth"
	"github.com/mzeahmed/noticoel/internal/modules/events"
	"github.com/mzeahmed/noticoel/internal/modules/health"
)

// New builds and returns the application's top-level http.Handler, with all
// module routes registered on a fresh http.ServeMux.
func New(db *sql.DB, appVersion, authToken string, log *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	health.New(appVersion).RegisterRoutes(mux)
	events.New(db, log).RegisterRoutes(mux, auth.Authenticate(authToken))

	return mux
}
