package events

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/mzeahmed/noticeal/internal/response"
)

// Handler handles all HTTP requests related to the events module.
type Handler struct {
	log *slog.Logger
}

// NewHandler creates a new events handler.
func NewHandler(log *slog.Logger) *Handler {
	return &Handler{log: log}
}

// Create handles POST /api/v1/events.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var e Event
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	if err := e.Validate(); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	h.log.Info("event received",
		"source", e.Source,
		"type", e.Type,
		"status", e.Status,
	)

	response.JSON(w, http.StatusAccepted, map[string]string{"status": "accepted"})
}
