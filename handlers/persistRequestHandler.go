package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/willdot/NotARealServer/persistrequests"
)

var errNoRequestRouteFound = errors.New("no request route property found")
var errRequestRouteIsNotString = errors.New("the request route provided is not a string")
var errNoRequestMethodFound = errors.New("no request method property found")
var errRequestMethodIsNotString = errors.New("the request method provided is not a string")

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
func (s *Server) SaveRequestHandler() http.HandlerFunc {
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

		err = s.HandleRequests.Save(requestContent, s.FileReadWriter)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(requestContent)
	}
}

// RetreiveRequestHandler takes the first parameter of the url and tried to return a saved request with that name
func (s *Server) RetreiveRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		requestRoute := params["RequestRoute"]
		requestMethod := r.Method

		result, err := s.HandleRequests.Load(requestRoute, requestMethod, s.FileReadWriter)

		if err != nil {
			http.Error(w, fmt.Sprintf("Problem retreiving request '%v'", requestRoute), http.StatusBadRequest)
			return
		}

		badResponse := checkHeaders(result.Headers, r.Header)

		if badResponse != nil {
			http.Error(w, badResponse.Message, badResponse.ErrorCode)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result.Response)
	}
}

func checkHeaders(savedHeaders []persistrequests.HeaderRequest, requestHeaders map[string][]string) *persistrequests.BadResponse {
	// Loop through all the saved headers so that we can see if the request contains that header
	for _, savedHeader := range savedHeaders {

		header := savedHeader.Header

		// As a header is a map, we need to loop through it to get the header key
		for k, v := range header {

			requestHeader := requestHeaders[k]
			// If the request doesn't contain a header with the key, return the bad response
			if requestHeader == nil {
				return &savedHeader.BadResponse
			}

			// Check the values of the saved header and the request header to make sure the request header values are valid
			if !reflect.DeepEqual(v, requestHeader) {
				return &savedHeader.BadResponse
			}
		}
	}

	return nil
}
