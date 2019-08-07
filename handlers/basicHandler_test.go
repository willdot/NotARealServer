package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBasic(t *testing.T) {
	var makeRequest = func(t *testing.T, url, body string, handler http.Handler, rr *httptest.ResponseRecorder) {

		t.Helper()

		req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))

		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")

		handler.ServeHTTP(rr, req)
	}

	testCases := []struct {
		Name               string
		Body               string
		Handler            http.HandlerFunc
		ExpectedStatusCode int
		ExpectedBody       string
	}{
		{
			"Basic handler, returns simple string",
			"",
			Basic(),
			http.StatusOK,
			"You hit basic",
		},
		{
			"BasicWithBody handler, request body ok, returns request",
			`{"Basic":"Request"}`,
			BasicWithBody(),
			http.StatusOK,
			`{"Basic":"Request"}`,
		},
		{
			"BasicWithBody handler, request body not ok, returns 400",
			`{"Basic`,
			BasicWithBody(),
			http.StatusBadRequest,
			"unexpected EOF",
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
