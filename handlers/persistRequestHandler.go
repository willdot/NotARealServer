package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/willdot/NotARealServer/persistrequests"
)

var errNoRequestRouteFound = errors.New("no request route property found")
var errRequestRouteIsNotString = errors.New("the request route provided is not a string")
var errNoRequestMethodFound = errors.New("no request method property found")
var errRequestMethodIsNotString = errors.New("the request method provided is not a string")

// PersistServer allows the user to save or retrieve requests
type PersistServer struct {
	FileWriter persistrequests.Writer
	FileReader persistrequests.Reader
	LoadSaver  persistrequests.SaveLoadRequest
}

// NewPersistServer creates a new PersistServer and adds in dependencies
func NewPersistServer(requestDirectory string) PersistServer {
	return PersistServer{
		FileWriter: persistrequests.FileWriter{},
		FileReader: persistrequests.FileReader{},
		LoadSaver: persistrequests.JSONPersist{
			RequestDirectory: requestDirectory,
		},
	}
}

func verifyRequestRouteAndMethod(request map[string]interface{}) error {

	route, ok := request["RequestRoute"]
	if !ok {
		return errNoRequestRouteFound
	}

	if !isTypeString(route) {
		return errRequestRouteIsNotString
	}

	method, ok := request["RequestMethod"]
	if !ok {
		return errNoRequestMethodFound
	}

	if !isTypeString(method) {
		return errRequestMethodIsNotString
	}

	return nil
}

func isTypeString(item interface{}) bool {
	switch item.(type) {
	case string:
		return true
	default:
		return false
	}
}

// SaveRequestHandler takes the body of the request and saves it as a json file
func (p *PersistServer) SaveRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)

		var requestContent map[string]interface{}

		err := decoder.Decode(&requestContent)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = verifyRequestRouteAndMethod(requestContent)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = p.LoadSaver.Save(requestContent, p.FileWriter)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(requestContent)
	}
}

// RetreiveRequestHandler takes the first parameter of the url and tried to return a saved request with that name
func (p *PersistServer) RetreiveRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		requestRoute := params["RequestRoute"]
		requestMethod := r.Method

		result, err := p.LoadSaver.Load(requestRoute, requestMethod, p.FileReader)

		if err != nil {
			http.Error(w, fmt.Sprintf("Problem retreiving request '%v'", requestRoute), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
