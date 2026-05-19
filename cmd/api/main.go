package main

import (
	"log"
	"net/http"
	"github.com/Sanja969/tracebeam-studio-server/internal/database"
	"github.com/Sanja969/tracebeam-studio-server/internal/events"
	"github.com/Sanja969/tracebeam-studio-server/internal/realtime"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	db, err := database.NewSQLiteConnection("tracebeam_studio.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := events.NewStore(db)
	hub := realtime.NewHub()
	handler := events.NewHandler(store, hub)

	router := mux.NewRouter()

	const eventsEndPoint = "/events"

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173", "http://localhost:5174"}),
		handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)(router)

	router.HandleFunc(eventsEndPoint, handler.GetEvents).Methods(http.MethodGet)
	router.HandleFunc(eventsEndPoint, handler.CreateEvent).Methods(http.MethodPost)
	router.HandleFunc(eventsEndPoint, handler.ClearEvents).Methods(http.MethodDelete)

	router.HandleFunc("/ws", hub.HandleWebSocket)

	log.Println("Tracebeam Studio running on port http://localhost:8080")

	if err := http.ListenAndServe(":8080", corsHandler); err != nil {
		log.Fatal(err)
	}

}
