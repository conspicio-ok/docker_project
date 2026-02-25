package main

import (
	"log"
	"net/http"
	"store-jv/handlers"
	"github.com/gorilla/mux"
)

func main() {
	InitDB()
	handlers.SetDB(DB)
	r := mux.NewRouter()
	r.Use(corsMiddleware)
	r.HandleFunc("/games", handlers.GetGames).Methods("GET")
	r.HandleFunc("/games/{id}", handlers.GetGame).Methods("GET")
	r.HandleFunc("/cart/{session_id}", handlers.GetCart).Methods("GET")
	r.HandleFunc("/cart/{session_id}/add", handlers.AddToCart).Methods("POST")
	r.HandleFunc("/cart/{session_id}/remove/{game_id}", handlers.RemoveFromCart).Methods("DELETE")

	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Println("API démarrée sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		next.ServeHTTP(w, r)
	})
}









