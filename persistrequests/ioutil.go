package persistrequests

import (
	"io/ioutil"
	"os"
)

// FileWriter implements an abstraction of ioutil.WriterFile
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

//FileReader implements an abstraction of the ioutil.ReadFile
type FileReader struct {
}

// ReadFile implements the Reader interface that has been created so that ioutil.ReadFile can be mocked or faked
func (f FileReader) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

// Reader is an interface to use over the ioutil.ReadFile() function so that it can be mocked or faked
type Reader interface {
	ReadFile(filename string) ([]byte, error)
}

// LoadRequest is an interface to load a request
type LoadRequest interface {
	Load(filename string, r Reader) (map[string]interface{}, error)
}
