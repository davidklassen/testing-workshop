package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// NewRouter creates an RPC-style http server for the greeter service.
func NewRouter(greeter *Greeter) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/greetings", func(w http.ResponseWriter, r *http.Request) {
		res, err := greeter.GetGreetings()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("failed to write http response: %v", err)
		}
	})

	router.HandleFunc("/say-hello", func(w http.ResponseWriter, r *http.Request) {
		var req = struct {
			Name string `json:"name"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := greeter.SayHello(req.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("failed to write http response: %v", err)
		}
	})

	router.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	return router
}
