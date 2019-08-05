package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/willdot/NotARealServer/persistrequests"
)

var errNoRequestNameFound = errors.New("no request property found")

// PersistServer allows the user to save or retrieve requests
type PersistServer struct {
	LoadSaver  persistrequests.JSONPersist
	FileWriter persistrequests.Writer
	FileReader persistrequests.Reader
}

// NewPersistServer creates a new PersistServer and adds in dependencies
func NewPersistServer() PersistServer {
	return PersistServer{
		FileWriter: persistrequests.FileWriter{},
		FileReader: persistrequests.FileReader{},
		LoadSaver:  persistrequests.JSONPersist{},
	}
}

// SaveRequestHandler takes the body of the request and saves it as a json file
func (p PersistServer) SaveRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)

		var request map[string]interface{}

		err := decoder.Decode(&request)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		filename, found := request["requestName"]

		if !found {
			http.Error(w, errNoRequestNameFound.Error(), http.StatusBadRequest)
			return
		}

		p.LoadSaver.Save(filename.(string), request, p.FileWriter)

		json.NewEncoder(w).Encode(request)
	}
}

// RetreiveRequestHandler takes the first parameter of the url and tried to return a saved request with that name
func (p PersistServer) RetreiveRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		request, found := params["request"]

		if !found {
			http.Error(w, "Request asked for cannot be found", http.StatusBadRequest)
			return
		}

		decodedFile, err := p.LoadSaver.Load(request+".json", p.FileReader)

		if err != nil {
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// remove the filename
		delete(decodedFile, "requestName")

		json.NewEncoder(w).Encode(decodedFile)
	}
}
