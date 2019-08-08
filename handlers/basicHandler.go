package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

// BasicHandler is a simple handler that just returns a 200 status code and an ok message
func BasicHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		io.WriteString(w, "You hit basic")
	}
}

// BasicWithBodyHandler will return the body that was sent in with the request
func BasicWithBodyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)

		var request map[string]interface{}

		err := decoder.Decode(&request)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(request)
	}
}
