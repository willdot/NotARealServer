package persistrequests

// SavedRequest is an entire saved request that requires a RequestRoute and RequestMethod. The Response is what the user wants to be returned when they make their fake API call.
type SavedRequest struct {
	RequestRoute  string
	RequestMethod string
	Response      interface{}
}

// DeleteRequest is a request to delete that is made up of the RequestRoute and RequestMethod which form the file names of the requests to delete
type DeleteRequest struct {
	RequestRoute  string
	RequestMethod string
}
