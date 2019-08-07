package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var makeRequest = func(t *testing.T, url, body string, handler http.Handler, rr *httptest.ResponseRecorder) {

	t.Helper()

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	handler.ServeHTTP(rr, req)
}
