package persistrequests

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var errNoRequestRouteFound = errors.New("no request route property found")
var errNoRequestMethodFound = errors.New("no request method property found")

// JSONPersist will allow a request to be saved and loaded to/from a JSON file
type JSONPersist struct {
	RequestDirectory string
}

// Save will save a request to a json file
func (j JSONPersist) Save(requestData map[string]interface{}, w Writer) error {

	requestMethod := requestData["RequestMethod"]
	if requestMethod == "" {
		return errNoRequestMethodFound
	}

	requestRoute := requestData["RequestRoute"]
	if requestRoute == "" {
		return errNoRequestRouteFound
	}
	file, err := json.MarshalIndent(requestData, "", " ")

	if err != nil {
		return err
	}

	filename := createFilename(requestMethod.(string), requestRoute.(string))
	err = w.WriteFile(j.RequestDirectory+filename, file, 0644)

	return err
}

// Load will load a json from a file
func (j JSONPersist) Load(requestRoute, requestMethod string, r Reader) (interface{}, error) {

	filename := createFilename(requestMethod, requestRoute)
	byteValue, err := r.ReadFile(j.RequestDirectory + filename)

	if err != nil {
		return nil, err
	}

	var savedRequest SavedRequest

	err = json.Unmarshal(byteValue, &savedRequest)

	if err != nil {
		return nil, err
	}

	return savedRequest.Response, nil
}

func createFilename(requestMethod, requestRoute string) string {
	return fmt.Sprintf("%v-%v.json", strings.ToUpper(requestMethod), strings.ToLower(requestRoute))
}

// SavedRequest is an entire saved request that requires a RequestRoute and RequestMethod. The Response is what the user wants to be returned when they make their fake API call.
type SavedRequest struct {
	RequestRoute  string
	RequestMethod string
	Response      interface{}
}
