package events

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/mzeahmed/noticeal/internal/response"
)

// Handler handles all HTTP requests related to the events module.
type Handler struct {
	service *Service
	log     *slog.Logger
}

// NewHandler creates a new events handler.
func NewHandler(service *Service, log *slog.Logger) *Handler {
	return &Handler{service: service, log: log}
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

	created, err := h.service.Create(r.Context(), e)
	if err != nil {
		h.log.Error("failed to store event", "error", err)
		response.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})

		return
	}

	h.log.Info("event received",
		"id", created.ID,
		"source", created.Source,
		"type", created.Type,
		"status", created.Status,
	)

	response.JSON(w, http.StatusAccepted, created)
}
