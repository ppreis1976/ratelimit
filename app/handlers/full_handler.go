package handlers

import (
	"log"
	"net/http"
)

func HandleFull(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s", r.Method, r.URL.Path)
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"name": "Full Cycle"}`))
		log.Printf("Responded with: %s", `{"name": "Full Cycle"}`)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("Method not allowed: %s", r.Method)
	}
}
