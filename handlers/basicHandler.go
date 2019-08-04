package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

// Basic is a simple handler that just returns a 200 status code and an ok message
func Basic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		io.WriteString(w, "You hit basic")
	}
}

// BasicWithBody will return the body that was sent in with the request
func BasicWithBody() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)

		var request map[string]interface{}

		err := decoder.Decode(&request)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(request)
	}
}
