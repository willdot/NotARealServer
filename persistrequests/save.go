package persistrequests

// SaveRequest is an interface to save a request
type SaveRequest interface {
	Save(requestData map[string]interface{}) error
}

// JSONSaver will allow a request to be saved to a JSON file
type JSONSaver struct {

}

// Save will save a request to a json file
func (j JSONSaver) Save(requestData map[string]interface{}) error {

	return nil
}