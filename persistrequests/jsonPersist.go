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

	err = w.CreateDirIfNotFound(j.RequestDirectory)

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

// Remove will remove all the requests that the user has requested to be removed. An error will be returned with any files that don't exist
func (j JSONPersist) Remove(requestsToRemove []DeleteRequest, r Remove) error {

	var errors []error
	for _, req := range requestsToRemove {
		filename := createFilename(req.RequestMethod, req.RequestRoute)

		err := r.Remove(j.RequestDirectory + filename)

		if err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) > 0 {
		return mergeErrors(errors)
	}

	return nil
}

// RemoveAll will delete the contents of the requests directory
func (j JSONPersist) RemoveAll(r RemoveAll) error {

	return r.RemoveAll(j.RequestDirectory)
}

func mergeErrors(errs []error) error {

	var errorMessage string
	for _, e := range errs {
		errorMessage += fmt.Sprintf("%v\n", e.Error())
	}

	errorMessage = strings.TrimRight(errorMessage, "\n")

	return errors.New(errorMessage)
}

func createFilename(requestMethod, requestRoute string) string {
	return fmt.Sprintf("%v-%v.json", strings.ToUpper(requestMethod), strings.ToLower(requestRoute))
}
