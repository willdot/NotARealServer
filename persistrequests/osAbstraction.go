package persistrequests

import "os"

// Remover is an interface to use over os.Remove() so that it can be mocked or faked
type Remover interface {
	Remove(name string) error
}

// FileRemover implements an abstraction over os.Remove
type FileRemover struct {
}

// Remove implements the Remover interface that's been created so that os.Remove() can be mocked or faked
func (FileRemover) Remove(name string) error {
	return os.Remove(name)
}
