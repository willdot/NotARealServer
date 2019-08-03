package persistrequests

import (
	"encoding/json"
	"os"
)

type Writer interface {
	WriteFile(filename string, data []byte, perm os.FileMode) error
}

// SaveRequest is an interface to save a request
type SaveRequest interface {
	Save(filename string, requestData map[string]interface{}, w Writer) error
}

// JSONSaver will allow a request to be saved to a JSON file
type JSONSaver struct {
}

// Save will save a request to a json file
func (j JSONSaver) Save(filename string, requestData map[string]interface{}, w Writer) error {

	file, err := json.MarshalIndent(requestData, "", " ")

	if err != nil {
		return err
	}

	//err = ioutil.WriteFile(filename+".json", file, 0644)
	err = w.WriteFile(filename+".json", file, 0644)

	return err
}
