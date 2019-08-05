package persistrequests

import (
	"encoding/json"
	"io/ioutil"
)

// JSONPersist will allow a request to be saved and loaded to/from a JSON file
type JSONPersist struct {
}

// Save will save a request to a json file
func (j JSONPersist) Save(filename string, requestData map[string]interface{}, w Writer) error {

	file, err := json.MarshalIndent(requestData, "", " ")

	if err != nil {
		return err
	}

	err = w.WriteFile(filename+".json", file, 0644)

	return err
}

// Load will load a json from a file
func (j JSONPersist) Load(filename string, r Reader) (map[string]interface{}, error) {

	byteValue, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	var result map[string]interface{}

	err = json.Unmarshal(byteValue, &result)

	if err != nil {
		return nil, err
	}

	return result, err
}
