package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBasic(t *testing.T) {

	testCases := []struct {
		Name               string
		Body               string
		Handler            http.HandlerFunc
		ExpectedStatusCode int
		ExpectedBody       string
	}{
		{
			Name:               "Basic handler, returns simple string",
			Body:               "",
			Handler:            Basic(),
			ExpectedStatusCode: http.StatusOK,
			ExpectedBody:       "You hit basic",
		},
		{
			Name:               "BasicWithBody handler, request body ok, returns request",
			Body:               `{"Basic":"Request"}`,
			Handler:            BasicWithBody(),
			ExpectedStatusCode: http.StatusOK,
			ExpectedBody:       `{"Basic":"Request"}`,
		},
		{
			Name:               "BasicWithBody handler, request body not ok, returns 400",
			Body:               `{"Basic`,
			Handler:            BasicWithBody(),
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody:       "unexpected EOF",
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			makeRequest(t, "", test.Body, test.Handler, rr)

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
