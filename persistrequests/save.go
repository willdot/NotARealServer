package persistrequests

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// FileWriter implements is an abstraction of ioutil.WriterFile
type FileWriter struct {
}

// WriteFile implements the Writer interface that's been created so that ioutil.WriteFile can be mocked or faked
func (w FileWriter) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}

// Writer is an interface to use over the ioutil.WriteFile() function so that it can be mocked or faked
type Writer interface {
	WriteFile(filename string, data []byte, perm os.FileMode) error
}

// SaveRequest is an interface to save a request
type SaveRequest interface {
	Save(filename string, requestData map[string]interface{}, w Writer) error
}

// JSONPersist will allow a request to be saved and loaded to/from a JSON file
type JSONPersist struct {
}

// Save will save a request to a json file
func (j JSONPersist) Save(filename string, requestData map[string]interface{}, w Writer) error {

	file, err := json.MarshalIndent(requestData, "", " ")

	if err != nil {
		return err
	}

	//err = ioutil.WriteFile(filename+".json", file, 0644)
	err = w.WriteFile(filename+".json", file, 0644)

	return err
}
