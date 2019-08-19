package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/willdot/NotARealServer/pkg/persistrequests"
)

// RemoveRequestHandler removes the requests provided by the user in the body
func (s Server) RemoveRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)

		var x DeleteJson

		err := decoder.Decode(&x)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = s.HandleRequests.Remove(x.Requests, s.FileRemover)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		io.WriteString(w, "Successfully deleted requests")
	}
}

// RemoveAllRequestsHandler removes the requests provided by the user in the body
func (s Server) RemoveAllRequestsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := s.HandleRequests.RemoveAll(s.FileRemover)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		io.WriteString(w, "Successfully deleted all requests")
	}
}

type DeleteJson struct {
	Requests []persistrequests.DeleteRequest
}
