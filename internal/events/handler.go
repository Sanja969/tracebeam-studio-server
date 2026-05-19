package events

import (
	"encoding/json"
	"github.com/Sanja969/tracebeam-studio-server/internal/realtime"
	"net/http"
)

type Handler struct {
	store *Store
	hub   *realtime.Hub
}

func NewHandler(store *Store, hub *realtime.Hub) *Handler {
	return &Handler{
		store: store,
		hub:   hub,
	}
}

func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.store.GetAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load events",
		})
		return
	}

	writeJSON(w, http.StatusOK, events)
}

func (h *Handler) ClearEvents(w http.ResponseWriter, r *http.Request) {
	if err := h.store.Clear(); err != nil {
	writeJSON(w, http.StatusInternalServerError, map[string]string{
		"error": "failed to clear events",
	})
	return
}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "events cleared",
	})
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}
	createdEvent, err := h.store.Add(event)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to save event",
		})
		return
	}
	h.hub.Broadcast(createdEvent)

	writeJSON(w, http.StatusCreated, createdEvent)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
