package main

import (
	"log"
	"net/http"

	"github.com/Sanja969/tracebeam-studio-server/internal/events"
	"github.com/Sanja969/tracebeam-studio-server/internal/realtime"
	"github.com/gorilla/mux"
)


func main() {
	store := events.NewStore()
	hub := realtime.NewHub()
	handler := events.NewHandler(store, hub)

	router := mux.NewRouter()

	const eventsEndPoint = "/events"

	router.HandleFunc(eventsEndPoint, handler.GetEvents).Methods(http.MethodGet)
	router.HandleFunc(eventsEndPoint, handler.CreateEvent).Methods(http.MethodPost)
	router.HandleFunc(eventsEndPoint, handler.ClearEvents).Methods(http.MethodDelete)

	router.HandleFunc("/ws", hub.HandleWebSocket)

	log.Println("Tracebeam Studio running on port http://localhost:8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}

}