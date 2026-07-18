package health

import "net/http"

// Handler handles all HTTP requests related to the health module.
type Handler struct {
	version string
}

// NewHandler creates a new health handler.
func NewHandler(version string) *Handler {
	return &Handler{version: version}
}

// Health handles GET /health.
func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

// Version handles GET /version.
func (h *Handler) Version(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(h.version))
}
