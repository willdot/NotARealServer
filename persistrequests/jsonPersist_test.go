package persistrequests

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
	"testing"
)

const directoryPath = "path/"

var errFakeError = errors.New("Fake error")

func TestSave(t *testing.T) {

	testCases := []struct {
		Name             string
		ValidRequestData bool
		RequestRoute     string
		RequestMethod    string
		ExpectedError    error
	}{
		{
			Name:             "Data input valid, no error returned",
			ValidRequestData: true,
			RequestRoute:     "POST",
			RequestMethod:    "Test",
			ExpectedError:    nil,
		},
		{
			Name:             "Data input invalid, error returned",
			ValidRequestData: false,
			RequestRoute:     "POST",
			RequestMethod:    "Test",
			ExpectedError:    createTestMarshalError(),
		},
		{
			Name:             "Request route not found, error returned",
			ValidRequestData: true,
			RequestRoute:     "",
			RequestMethod:    "POST",
			ExpectedError:    errNoRequestRouteFound,
		},
		{
			Name:             "Request method not found, error returned",
			ValidRequestData: true,
			RequestRoute:     "Test",
			RequestMethod:    "",
			ExpectedError:    errNoRequestMethodFound,
		},
	}

	testObj := JSONPersist{}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			testData := createData(test.ValidRequestData, test.RequestRoute, test.RequestMethod)

			got := testObj.Save(testData, fakeFileReaderWriter{})

			assertErrors(got, test.ExpectedError, t)
		})
	}
}

func TestLoad(t *testing.T) {

	testCases := []struct {
		Name          string
		RequestRoute  string
		OutputData    interface{}
		ExpectedError error
	}{
		{
			Name:         "File valid, data returned, no error returned",
			RequestRoute: "good",
			OutputData: Something{
				Count:    1,
				Language: "Go",
			},
			ExpectedError: nil,
		},
		{
			Name:          "File invalid, no data returned, error returned",
			RequestRoute:  "badfile",
			OutputData:    nil,
			ExpectedError: errFakeError,
		},
		{
			Name:          "File valid, file data invalid, no data returned, error returned",
			RequestRoute:  "badjsonformat",
			OutputData:    nil,
			ExpectedError: createUnmarshalError(),
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			got, err := JSONPersist{}.Load(test.RequestRoute, "POST", fakeFileReaderWriter{})

			assertErrors(err, test.ExpectedError, t)

			var want interface{}
			jsonByte, _ := json.Marshal(test.OutputData)
			json.Unmarshal(jsonByte, &want)

			gotString := fmt.Sprintf("%v", got.Response)
			wantString := fmt.Sprintf("%v", want)

			if gotString != wantString {
				t.Errorf("Got %v, wanted %v", got, want)
			}

		})
	}
}

func TestRemoveRequests(t *testing.T) {
	testCases := []struct {
		Name             string
		RequestsToRemove []DeleteRequest
		NumberOfErrors   int
		ExpectedError    error
	}{
		{
			Name:             "1 Request provided and it's removed successfully with no error",
			RequestsToRemove: createDeleteRequests(1),
			ExpectedError:    createError(0),
		},
		{
			Name:             "2 Requests provided and both removed successfully with no error",
			RequestsToRemove: createDeleteRequests(2),
			ExpectedError:    createError(0),
		},
		{
			Name:             "2 Requests provided one doesn't exist and error returned stating the 1 file that doesn't exist",
			RequestsToRemove: createDeleteRequests(2),
			NumberOfErrors:   1,
			ExpectedError:    createError(1),
		},
		{
			Name:             "2 Requests provided neither exist and error returned stating the 2 files don't exist",
			RequestsToRemove: createDeleteRequests(2),
			NumberOfErrors:   2,
			ExpectedError:    createError(2),
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			jp := JSONPersist{directoryPath}

			fr := fakeMultipleRemover{errorRequests: test.RequestsToRemove[:test.NumberOfErrors]}

			err := jp.Remove(test.RequestsToRemove, &fr)

			assertErrors(err, test.ExpectedError, t)
		})
	}
}

func TestRemoveAll(t *testing.T) {

	testCases := []struct {
		Name          string
		Path          string
		ExpectedError error
	}{
		{
			Name:          "Path exists, deletes, no error",
			Path:          "this/path/exists",
			ExpectedError: nil,
		},
		{
			Name:          "Path does not exist, doesn't delete, no error",
			Path:          "does/not/exist",
			ExpectedError: nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			jp := JSONPersist{test.Path}
			fr := fakeRemover{err: test.ExpectedError}

			err := jp.RemoveAll(fr)

			if err != test.ExpectedError {
				t.Errorf("Didn't want an error, but got %v", err)
			}
		})
	}
}

func TestMergeErrors(t *testing.T) {

	testCases := []struct {
		Name           string
		NumberOfErrors int
		ExpectedError  error
	}{
		{
			Name:           "One error, error returned",
			NumberOfErrors: 1,
			ExpectedError:  errors.New("error 1"),
		},
		{
			Name:           "Two errors, error returned",
			NumberOfErrors: 2,
			ExpectedError:  errors.New("error 1\nerror 2"),
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			errors := make([]error, test.NumberOfErrors)

			for i := 0; i < test.NumberOfErrors; i++ {

				err := fmt.Errorf("error %v", i+1)
				errors[i] = err
			}

			got := mergeErrors(errors)

			assertErrors(got, test.ExpectedError, t)
		})
	}
}

func assertErrors(got, want error, t *testing.T) {
	// If both are nil, then all is fine
	if got == nil && want == nil {
		return
	}

	// if got or want is nil, then return the comparision
	if got == nil || want == nil {

		if got != want {
			t.Errorf("Got %v, wanted %v", got, want)
		}
		return
	}

	// Neither error will be nil so now an actual check for the errors can be done
	if got.Error() != want.Error() {
		t.Errorf("Got %v, wanted %v", got, want)
		return
	}

	return
}

// This is a fake remover that returns an error if required
type fakeRemover struct {
	err error
}

// Remove mocks the os.Remove() but will just return and error if there is one
func (fr fakeRemover) Remove(name string) error {
	return fr.err
}

// Remove mocks the os.RemoveAll() but will just return and error if there is one
func (fr fakeRemover) RemoveAll(path string) error {
	return fr.err
}

// This is a fake Remover that has logic to send back an error message. The slice of SavedRequests are the files that don't exist.
type fakeMultipleRemover struct {
	errorRequests []DeleteRequest
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

func createRequests(numberOfRequests int) []SavedRequest {
	result := make([]SavedRequest, numberOfRequests)

	for i := 0; i < numberOfRequests; i++ {
		result[i] = SavedRequest{
			RequestRoute:  fmt.Sprintf("HERE %v", i+1),
			RequestMethod: "POST",
		}
	}

	return result
}

func createDeleteRequests(numberOfRequests int) []DeleteRequest {
	result := make([]DeleteRequest, numberOfRequests)

	for i := 0; i < numberOfRequests; i++ {
		result[i] = DeleteRequest{
			RequestRoute:  fmt.Sprintf("HERE %v", i+1),
			RequestMethod: "POST",
		}
	}

	return result
}

type fakeFileReaderWriter struct {
}

func (f fakeFileReaderWriter) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return nil
}

func (f fakeFileReaderWriter) ReadFile(filename string) ([]byte, error) {

	// For when the file can't be read such as doesn't exist
	if filename == "POST-badfile.json" {
		return nil, errFakeError
	}

	// For when the data inside a json file is not correct such as missing a " somewhere
	if filename == "POST-badjsonformat.json" {

		result := []byte(`{"Something" : "Wrong}`)
		return result, nil
	}

	JSONFile := `{
		"RequestRoute": "Hello",
		"RequestMethod": "POST",
		"Response": {
		 "Count": "1",
		 "Language": "Go"
		}
	   }`

	result := []byte(JSONFile)

	return result, nil
}

func (f fakeFileReaderWriter) CreateDirIfNotFound(path string) error {
	return nil
}

type Something struct {
	Language string
	Count    float64
}

func createData(goodData bool, requestRoute, requestMethod string) map[string]interface{} {

	var result map[string]interface{}
	result = make(map[string]interface{})
	if goodData {
		thing := SavedRequest{
			RequestRoute:  requestRoute,
			RequestMethod: requestMethod,
			Response: Something{
				Language: "Go",
				Count:    1.000,
			},
		}

		jsonByte, _ := json.Marshal(thing)
		json.Unmarshal(jsonByte, &result)

		return result
	}

	request := SavedRequest{
		RequestRoute:  requestRoute,
		RequestMethod: "requestMethod",
		Response: Something{
			Language: "Go",
			Count:    math.Inf(1),
		},
	}

	result["test"] = request

	return result
}

// This creates the same error that will be returned from json.Marshal when sending in data that isn't valid
func createTestMarshalError() error {
	_, err := json.Marshal(math.Inf(1))
	return err.(*json.UnsupportedValueError)
}

// This creates the same error that will be returned from json.Unmarshal when sending in data that isn't valid
func createUnmarshalError() error {

	a := []byte(`{"Something" : "Wrong}`)
	var b map[string]interface{}
	err := json.Unmarshal(a, &b)

	return err.(*json.SyntaxError)
}
