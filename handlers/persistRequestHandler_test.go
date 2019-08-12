package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/willdot/NotARealServer/persistrequests"
)

type FakeFileWriter struct {
}

type fakeFileReaderWriter struct {
}

func (f fakeFileReaderWriter) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return nil
}

func (f fakeFileReaderWriter) ReadFile(filename string) ([]byte, error) {

	var json string

	// If the test is to do with bad header, take the fake json and add in "nothing as the content-type. Otherwise, put in "application/json" which is valid.
	if filename == "POST-badheader.json" {
		json = fmt.Sprintf(fakeJSON, "nothing")
	} else {
		json = fmt.Sprintf(fakeJSON, "application/json")
	}

	if filename == "POST-badjson.json" {

		json = strings.TrimRight(json, `"}}`)

		return []byte(json), nil
	}

	fmt.Println(json)

	return []byte(json), nil

}

func (f fakeFileReaderWriter) CreateDirIfNotFound(path string) error {
	return nil
}

type fakeFileRemover struct {
}

// Remove implements the Remover interface that's been created so that os.Remove() can be mocked or faked
func (f fakeFileRemover) Remove(name string) error {

	if name == "NOT-exists.json" {
		return os.ErrNotExist
	}

	return nil
}

// This is a flag to test if the RemoveAll will return an error or not
var removeAllError = false

// RemoveAll implements the Remover interface thats been created so that os.RemoveAll() can be mocked or faked
func (f fakeFileRemover) RemoveAll(path string) error {

	if removeAllError {
		return errors.New("fake error")
	}
	return nil
}

var fakeJSON = `{"RequestMethod":"POST","Headers":[{"Header": {"Content-Type": ["%v"]},"BadResponse": "Content type not valid"}],"RequestRoute":"Test","Response":{"something":"fake"}}`

var directoryPath = ""

var testThing = Server{
	FileReadWriter: fakeFileReaderWriter{},
	FileRemover:    fakeFileRemover{},
	HandleRequests: persistrequests.JSONPersist{RequestDirectory: directoryPath},
}

func TestRetreiveRequestHandler(t *testing.T) {

	testCases := []struct {
		Name               string
		Route              string
		ExpectedStatusCode int
		ExpectedBody       string
	}{
		{
			Name:               "Param ok. Request Exists. Request returned. requestName removed from data",
			Route:              "/test",
			ExpectedStatusCode: http.StatusOK,
			ExpectedBody:       `{"something":"fake"}`,
		},
		{
			Name:               "Param ok. Request file doesn't exist. Returns 400 bad request",
			Route:              "/badjson",
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody:       "Problem retreiving request 'badjson'",
		},
		{
			Name:               "Param ok. Headers not ok. Returns 400 bad request",
			Route:              "/badheader",
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody:       "Content type not valid",
		},
	}

	handler := mux.NewRouter()
	handler.HandleFunc("/{RequestRoute}", testThing.RetreiveRequestHandler())

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			makeRequest(t, test.Route, "", handler, rr)

			if status := rr.Code; status != test.ExpectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, test.ExpectedStatusCode)
			}

			got := strings.TrimSuffix(rr.Body.String(), "\n")

			if got != test.ExpectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), test.ExpectedBody)
			}
		})
	}
}

func TestSaveRequestHandler(t *testing.T) {
	handler := testThing.SaveRequestHandler()

	testCases := []struct {
		Name               string
		Body               string
		ExpectedStatusCode int
		ExpectedBody       string
	}{
		{
			Name:               "Body ok. Returns 200",
			Body:               `{"RequestRoute" : "Test","RequestMethod" : "POST","Request" : {"Something" : "Fake"}}`,
			ExpectedStatusCode: http.StatusOK,
			ExpectedBody:       `{"Request":{"Something":"Fake"},"RequestMethod":"POST","RequestRoute":"Test"}`,
		},
		{
			Name:               "Body not ok. Returns 400",
			Body:               `{"RequestRoute" : "Test}`,
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody:       "unexpected EOF",
		},
		{
			Name:               "Body doesn't have RequestRoute property. Returns 400",
			Body:               `{"RequestMethod" : "POST","Request" : {"Something" : "Fake"}}`,
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody:       "no request route property found",
		},
		{
			Name:               "Body doesn't have RequestMethod property. Returns 400",
			Body:               `{"RequestRoute" : "Test","Request" : {"Something" : "Fake"}}`,
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody:       "no request method property found",
		},
		{
			Name:               "Bodies RequestRoute is empty. Returns 400",
			Body:               `{"RequestRoute" : "","RequestMethod" : "POST","Request" : {"Something" : "Fake"}}`,
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody:       "no request route property found",
		},
		{
			Name:               "Bodies RequestMethod is empty. Returns 400",
			Body:               `{"RequestRoute" : "Test","RequestMethod" : "","Request" : {"Something" : "Fake"}}`,
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody:       "no request method property found",
		},
		{
			Name:               "Bodies RequestRoute is not a string. Returns 400",
			Body:               `{"RequestRoute" : {"Route" : "WRONG"},"RequestMethod" : "POST","Request" : {"Something" : "Fake"}}`,
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody:       "the request route provided is not a string",
		},
		{
			Name:               "Bodies RequestMethod is not a string. Returns 400",
			Body:               `{"RequestRoute" : "Test","RequestMethod" : {"Method" : "WRONG"},"Request" : {"Something" : "Fake"}}`,
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody:       "the request method provided is not a string",
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			makeRequest(t, "save", test.Body, handler, rr)

			if status := rr.Code; status != test.ExpectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, test.ExpectedStatusCode)
			}

			got := strings.TrimSuffix(rr.Body.String(), "\n")

			if got != test.ExpectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), test.ExpectedBody)
			}
		})
	}
}

func TestIsTypeString(t *testing.T) {
	testCases := []struct {
		Name           string
		Input          interface{}
		ExpectedResult bool
	}{
		{
			Name:           "Input is string. Returns true",
			Input:          "Hello",
			ExpectedResult: true,
		},
		{
			Name:           "Input is int. Returns false",
			Input:          1,
			ExpectedResult: false,
		},
		{
			Name:           "Input is struct. Returns false",
			Input:          FakeFileWriter{},
			ExpectedResult: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			got := isTypeString(test.Input)

			if got != test.ExpectedResult {
				t.Errorf("got %v, want %v", got, test.ExpectedResult)
			}
		})
	}
}

func TestCheckHeaders(t *testing.T) {
	testCases := []struct {
		Name                 string
		SavedHeaders         []persistrequests.HeaderRequest
		RequestHeaders       map[string][]string
		ExpectedResult       bool
		ExpectedErrorMessage string
	}{
		{
			Name:                 "1 header saved, 1 header supplied. Matching. Success",
			SavedHeaders:         createSavedHeaders(1),
			RequestHeaders:       createRequestHeaders(1, true),
			ExpectedResult:       true,
			ExpectedErrorMessage: "",
		},
		{
			Name:                 "2 headers saved, 1 header supplied. Not enough supplied. Failure",
			SavedHeaders:         createSavedHeaders(2),
			RequestHeaders:       createRequestHeaders(1, true),
			ExpectedResult:       false,
			ExpectedErrorMessage: "Error",
		},
		{
			Name:                 "1 headers saved, 1 header supplied. Header values don't match. Failure",
			SavedHeaders:         createSavedHeaders(1),
			RequestHeaders:       createRequestHeaders(1, false),
			ExpectedResult:       false,
			ExpectedErrorMessage: "Error",
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			got, errorMessage := checkHeaders(test.SavedHeaders, test.RequestHeaders)

			if got != test.ExpectedResult {
				t.Errorf("got %v, want %v", got, test.ExpectedResult)
			}

			if errorMessage != test.ExpectedErrorMessage {
				t.Errorf("got %v, want %v", errorMessage, test.ExpectedErrorMessage)
			}
		})
	}

}

func createSavedHeaders(headerCount int) []persistrequests.HeaderRequest {

	result := make([]persistrequests.HeaderRequest, headerCount)

	for i := 0; i < headerCount; i++ {
		result[i] = persistrequests.HeaderRequest{
			BadResponse: "Error",
			Header: map[string][]string{
				fmt.Sprintf("header %v", i): []string{"value"}},
		}
	}

	return result
}

func createRequestHeaders(headerCount int, matchSavedHeader bool) map[string][]string {

	headerValue := "value"

	if !matchSavedHeader {
		headerValue = "something"
	}

	result := make(map[string][]string, headerCount)

	for i := 0; i < headerCount; i++ {
		result[fmt.Sprintf("header %v", i)] = []string{headerValue}
	}

	return result
}
