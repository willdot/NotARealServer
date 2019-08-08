package persistrequests

import (
	"errors"
	"os"
	"testing"
)

type fakeRemover struct{}

const directoryPath = "path/"

var errFileNotExists = errors.New("file does not exist")

func (fakeRemover) Remove(name string) error {

	if name == directoryPath+"WRONG-hello.json" {
		return os.ErrNotExist
	}

	return nil
}

func TestRemoveRequest(t *testing.T) {

	fs := fakeRemover{}

	testCases := []struct {
		Name          string
		Route         string
		Method        string
		ExpectedError error
	}{
		{
			Name:          "File exists. Deleted. No Error",
			Route:         "hello",
			Method:        "POST",
			ExpectedError: nil,
		},
		{
			Name:          "File does not exist. Error returned",
			Route:         "hello",
			Method:        "WRONG",
			ExpectedError: os.ErrNotExist,
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			err := JSONPersist{directoryPath}.RemoveRequest(test.Method, test.Route, fs)

			if err != test.ExpectedError {
				t.Errorf("Didn't want an error, but got %v", err)
			}
		})
	}

}
