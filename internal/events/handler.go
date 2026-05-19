package events

import (
	"encoding/json"
	"net/http"
	"github.com/Sanja969/tracebeam-studio-server/internal/realtime"
)

type Handler struct {
	store *Store
	hub *realtime.Hub
}

func NewHandler(store *Store, hub *realtime.Hub) *Handler {
	return &Handler{
		store: store,
		hub: hub,
	}
}

func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	events := h.store.GetAll()

	writeJSON(w, http.StatusOK, events)
}

func (h *Handler) ClearEvents(w http.ResponseWriter, r *http.Request) {
	h.store.Clear()
	
	writeJSON(w, http.StatusOK, map[string]string{
		"message": "events cleared",
	})
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event Event
	if err :=json.NewDecoder(r.Body).Decode(&event); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}
	createdEvent := h.store.Add(event)
	h.hub.Broadcast(createdEvent)
	
	writeJSON(w, http.StatusCreated, createdEvent)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}