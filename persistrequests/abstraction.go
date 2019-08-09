package persistrequests

import (
	"io/ioutil"
	"os"
)

// Writer is an interface to use over the ioutil.WriteFile() function so that it can be mocked or faked
type Writer interface {
	WriteFile(filename string, data []byte, perm os.FileMode) error
}

// FileWriter implements an abstraction of ioutil.WriterFile
type FileWriter struct {
}

// WriteFile implements the Writer interface that's been created so that ioutil.WriteFile can be mocked or faked
func (w FileWriter) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}

// Reader is an interface to use over the ioutil.ReadFile() function so that it can be mocked or faked
type Reader interface {
	ReadFile(filename string) ([]byte, error)
}

//FileReader implements an abstraction of the ioutil.ReadFile
type FileReader struct {
}

// ReadFile implements the Reader interface that has been created so that ioutil.ReadFile can be mocked or faked
func (f FileReader) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

// Remove is an interface to use over os.Remove() so that it can be mocked or faked
type Remove interface {
	Remove(name string) error
}

// RemoveAll is an interface to use over os.RemoveAll() so that it can be mocked or faked
type RemoveAll interface {
	RemoveAll(path string) error
}

// Remover is an interface to allow the removal of files
type Remover interface {
	Remove
	RemoveAll
}

// FileRemover implements an abstraction over os.Remove and os.RemoveAll
type FileRemover struct {
}

// Remove implements the Remover interface that's been created so that os.Remove() can be mocked or faked
func (f FileRemover) Remove(name string) error {
	return os.Remove(name)
}

// RemoveAll implements the Remover interface thats been created so that os.RemoveAll() can be mocked or faked
func (f FileRemover) RemoveAll(path string) error {
	return os.RemoveAll(path)
}
