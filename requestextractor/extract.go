package requestextractor

import "errors"

// SavedRequest is an entire saved request
type SavedRequest struct {
	RequestRoute string
	Requests     []Request
}

// Request is an individual request that can be sent through a route
type Request struct {
	MethodType string
	Result     interface{}
}

var errRequestNotFound = errors.New("Request can't be found")

// Extract will return part of the request that has been saved based on the method type they ask for
func Extract(methodtype string, savedRequest SavedRequest) (Request, error) {

	var request Request
	for _, x := range savedRequest.Requests {
		if x.MethodType == methodtype {
			request = x
			return request, nil
		}
	}
	return request, errRequestNotFound
}
