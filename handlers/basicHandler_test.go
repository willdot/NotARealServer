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

	t.Run("Basic handler, returns simple string", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := Basic()
		body := ""

		makeRequest(t, "", body, handler, rr)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		got := strings.TrimSuffix(rr.Body.String(), "\n")

		want := "You hit basic"

		if got != want {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), want)
		}
	})

	t.Run("BasicWithBody handler, request body ok, returns request", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := BasicWithBody()
		body := `{"Basic":"Request"}`

		makeRequest(t, "", body, handler, rr)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		got := strings.TrimSuffix(rr.Body.String(), "\n")

		want := body

		if got != want {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), want)
		}

	})

	t.Run("BasicWithBody handler, request body not ok, returns 400", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := BasicWithBody()
		body := `{"Basic`

		makeRequest(t, "", body, handler, rr)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})
}
