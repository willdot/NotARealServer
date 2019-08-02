package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// SaveRequest takes the body of the request and saves it as a json file
func SaveRequest() http.HandlerFunc {
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
			http.Error(w, "no requestName property found in body", http.StatusBadRequest)
			return
		}

		save(filename.(string)+".json", request)

		json.NewEncoder(w).Encode(request)
	}
}

// RetreiveRequest takes the first parameter of the url and tried to return a saved request with that name
func RetreiveRequest() http.HandlerFunc {
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

func save(filename string, data map[string]interface{}) {

	file, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile(filename, file, 0644)
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
