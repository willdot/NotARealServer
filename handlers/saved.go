package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/willdot/NotARealServer/persistrequests"
)

// PersistServer allows the user to save or retrieve requests
type PersistServer struct {
	Saver      persistrequests.SaveRequest
	FileWriter persistrequests.Writer
}

// NewPersistServer creates a new PersistServer and adds in dependencies
func NewPersistServer() PersistServer {
	return PersistServer{
		Saver:      persistrequests.JSONSaver{},
		FileWriter: persistrequests.FileWriter{},
	}
}

// SaveRequest takes the body of the request and saves it as a json file
func (p PersistServer) SaveRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)

		var request map[string]interface{}

		err := decoder.Decode(&request)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		filename, _ := request["requestName"]

		p.Saver.Save(filename.(string), request, persistrequests.FileWriter{})

		json.NewEncoder(w).Encode(request)
	}
}

// RetreiveRequest takes the first parameter of the url and tried to return a saved request with that name
func (p PersistServer) RetreiveRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		request, found := params["request"]

		if !found {
			http.Error(w, "Request asked for cannot be found", http.StatusBadRequest)
			return
		}

		decodedFile := load(request + ".json")

		// remove the filename
		delete(decodedFile, "requestName")

		json.NewEncoder(w).Encode(decodedFile)
	}
}

func load(filename string) map[string]interface{} {

	jsonFile, err := os.Open(filename)
	defer jsonFile.Close()

	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(jsonFile)

	var result map[string]interface{}
	decoder.Decode(&result)

	return result
}
