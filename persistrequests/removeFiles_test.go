package persistrequests

import (
	"errors"
	"os"
	"testing"
)

type fakeRemover struct {
	err error
}

const directoryPath = "path/"

var errFileNotExists = errors.New("file does not exist")

func (f fakeRemover) Remove(name string) error {

	return f.err
}

func TestRemoveRequest(t *testing.T) {

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

			jp := JSONPersist{directoryPath}
			fr := fakeRemover{err: test.ExpectedError}

			err := jp.RemoveRequest(test.Method, test.Route, fr)

			if err != test.ExpectedError {
				t.Errorf("Didn't want an error, but got %v", err)
			}
		})
	}

}
