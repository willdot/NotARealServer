package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
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

	if filename == "POST-badjson.json" {

		fakeJSON = strings.TrimRight(fakeJSON, `"}}`)

		return []byte(fakeJSON), nil
	}
	return []byte(fakeJSON), nil

}

var fakeJSON = `{"RequestMethod":"POST","RequestRoute":"Test","Response":{"something":"fake"}}`

var testThing = PersistServer{
	FileWriter: fakeFileReaderWriter{},
	FileReader: fakeFileReaderWriter{},
	LoadSaver:  persistrequests.JSONPersist{},
}

func TestNewPersistServer(t *testing.T) {

	got := NewPersistServer("")

	want := PersistServer{
		LoadSaver:  persistrequests.JSONPersist{},
		FileWriter: persistrequests.FileWriter{},
		FileReader: persistrequests.FileReader{},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got %v, wanted %v", got, want)
	}
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
