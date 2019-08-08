package persistrequests

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
)

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
			fr := fakeSingleRemover{err: test.ExpectedError}

			err := jp.RemoveRequest(test.Method, test.Route, fr)

			if err != test.ExpectedError {
				t.Errorf("Didn't want an error, but got %v", err)
			}
		})
	}

}

func TestRemoveMultipleRequests(t *testing.T) {
	testCases := []struct {
		Name             string
		RequestsToRemove []SavedRequest
		NumberOfErrors   int
		ExpectedError    error
	}{
		{
			Name:             "2 Requests provided and both removed successfully with no error",
			RequestsToRemove: createRequests(),
			ExpectedError:    createError(0),
		},
		{
			Name:             "2 Requests provided one doesn't exist and error returned stating the 1 file that doesn't exist",
			RequestsToRemove: createRequests(),
			NumberOfErrors:   1,
			ExpectedError:    createError(1),
		},
		{
			Name:             "2 Requests provided neither exist and error returned stating the 2 files don't exist",
			RequestsToRemove: createRequests(),
			NumberOfErrors:   2,
			ExpectedError:    createError(2),
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			jp := JSONPersist{directoryPath}

			fr := fakeMultipleRemover{errorRequests: test.RequestsToRemove[:test.NumberOfErrors]}

			err := jp.RemoveMultipleRequests(test.RequestsToRemove, &fr)

			assertErrors(err, test.ExpectedError, t)
		})
	}
}

const directoryPath = "path/"

// This is a fake remover that returns an error if required
type fakeSingleRemover struct {
	err error
}

// Remove mocks the os.Remove() but will just return and error if there is one
func (fsr fakeSingleRemover) Remove(name string) error {
	return fsr.err
}

// This is a fake Remover that has logic to send back an error message. The slice of SavedRequests are the files that don't exist.
type fakeMultipleRemover struct {
	errorRequests []SavedRequest
}

// Remove mocks the os.Remove() but in this case is will see if the incoming filename is in the slice of errorRequests on the struct, and if if is, then it'll return the error
func (fmr *fakeMultipleRemover) Remove(name string) error {

	for _, v := range fmr.errorRequests {
		filename := directoryPath + createFilename(v.RequestMethod, v.RequestRoute)

		if name == filename {
			return os.ErrNotExist
		}
	}
	return nil
}

// This creates an error message that contains an error message of a given number of file not exist errors
func createError(numberOfErrors int) error {

	if numberOfErrors == 0 {
		return nil
	}
	errorMessage := ""

	for i := 0; i < numberOfErrors; i++ {
		errorMessage += fmt.Sprintf("%v\n", os.ErrNotExist.Error())

	}
	errorMessage = strings.TrimRight(errorMessage, "\n")
	return errors.New(errorMessage)
}

func createRequests() []SavedRequest {
	result := make([]SavedRequest, 2)

	result[0] = SavedRequest{
		RequestRoute:  "Here",
		RequestMethod: "POST",
	}

	result[1] = SavedRequest{
		RequestRoute:  "Somewhere",
		RequestMethod: "POST",
	}

	return result
}
