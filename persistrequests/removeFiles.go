package persistrequests

import (
	"errors"
	"fmt"
	"strings"
)

// RemoveRequest generates a filename from the given parameters the user has provided and will try to delete that request if it exists. An error will be returned if it doesn't exist
func (j JSONPersist) RemoveRequest(method, route string, r Remover) error {

	filename := createFilename(method, route)

	return r.Remove(j.RequestDirectory + filename)
}

// RemoveMultipleRequests will remove all the requests that the user has requested to be removed. An error will be returned with any files that don't exist
func (j JSONPersist) RemoveMultipleRequests(requestsToRemove []SavedRequest, r Remover) error {

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

func mergeErrors(errs []error) error {

	var errorMessage string
	for _, e := range errs {
		errorMessage += fmt.Sprintf("%v\n", e.Error())
	}

	errorMessage = strings.TrimRight(errorMessage, "\n")

	return errors.New(errorMessage)
}
