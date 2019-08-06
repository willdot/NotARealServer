package persistrequests

import (
	"encoding/json"
	"fmt"
)

// JSONPersist will allow a request to be saved and loaded to/from a JSON file
type JSONPersist struct {
}

// Save will save a request to a json file
func (j JSONPersist) Save(requestRoute, requestMethod string, requestData interface{}, w Writer) error {

	file, err := json.MarshalIndent(requestData, "", " ")

	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%v-%v.json", requestMethod, requestRoute)
	err = w.WriteFile(filename, file, 0644)

	return err
}

// Load will load a json from a file
func (j JSONPersist) Load(requestRoute, requestMethod string, r Reader) (interface{}, error) {

	filename := fmt.Sprintf("%v-%v.json", requestMethod, requestRoute)
	byteValue, err := r.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	var savedRequest SavedRequest

	err = json.Unmarshal(byteValue, &savedRequest)

	if err != nil {
		return nil, err
	}

	return savedRequest.Request, err
}

// SavedRequest is an entire saved request
type SavedRequest struct {
	RequestRoute  string
	RequestMethod string
	Request       interface{}
}
