package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRemoveRequestHandler(t *testing.T) {

	handler := testThing.RemoveRequestHandler()
	testCases := []struct {
		Name               string
		Body               string
		ExpectedStatusCode int
		ExpectedBody       string
	}{
		{
			Name:               "Remove with valid request, no errors returned",
			Body:               `{"Requests": [{"RequestRoute": "CreateTest","RequestMethod": "POST"}]}`,
			ExpectedStatusCode: http.StatusOK,
			ExpectedBody:       "Successfully deleted requests",
		},
		{
			Name:               "Remove with invalid request, error returned",
			Body:               `{"Requests": [{"RequestRoute}`,
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody:       "unexpected EOF",
		},
		{
			Name:               "Remove with request, file doesn't exist error returned",
			Body:               `{"Requests": [{"RequestRoute": "exists","RequestMethod": "NOT"}]}`,
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody:       "file does not exist",
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			makeRequest(t, "remove", test.Body, handler, rr)

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

func TestRemoveAllRequestsHandler(t *testing.T) {

	handler := testThing.RemoveAllRequestsHandler()
	testCases := []struct {
		Name               string
		HasError           bool
		ExpectedStatusCode int
		ExpectedBody       string
	}{
		{
			Name:               "Remove all no errors returned",
			HasError:           false,
			ExpectedStatusCode: http.StatusOK,
			ExpectedBody:       "Successfully deleted all requests",
		},
		{
			Name:               "Remove with error, error returned",
			HasError:           true,
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody:       "fake error",
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			rr := httptest.NewRecorder()

			removeAllError = test.HasError
			makeRequest(t, "removeall", "", handler, rr)

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
