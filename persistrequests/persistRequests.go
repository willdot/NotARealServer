package persistrequests

// SaveRequest is an interface to save a request
type SaveRequest interface {
	Save(requestData map[string]interface{}, w Writer) error
}

// LoadRequest is an interface to load a request
type LoadRequest interface {
	Load(requestRoute, requestMethod string, r Reader) (interface{}, error)
}

// RemoveRequest is an interface to Remove a single request
type RemoveRequest interface {
	Remove(requestsToRemove []DeleteRequest, r Remove) error
}

// RemoveAllRequests is an interface to remove all requests
type RemoveAllRequests interface {
	RemoveAll(r RemoveAll) error
}

// HandleRequests allows the saving and loading of a request
type HandleRequests interface {
	SaveRequest
	LoadRequest
	RemoveRequest
	RemoveAllRequests
}
