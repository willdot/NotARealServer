package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

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

	return nil, nil
}

var testThing = PersistServer{
	FileWriter: fakeFileReaderWriter{},
	LoadSaver:  persistrequests.JSONPersist{},
}

type fake struct {
	requestName string
	something   string
}

func TestSaveRequestHandler(t *testing.T) {

	makeRequest := func(t *testing.T, body string, rr *httptest.ResponseRecorder) {

		t.Helper()

		req, err := http.NewRequest(http.MethodPost, "/save", strings.NewReader(body))

		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")
		handler := http.HandlerFunc(testThing.SaveRequestHandler())

		handler.ServeHTTP(rr, req)
	}

	t.Run("Body ok. Returns 200", func(t *testing.T) {
		body := `{
			"requestName": "Test",
			"something" : "Hello"
		   }`

		rr := httptest.NewRecorder()

		makeRequest(t, body, rr)

		got := strings.TrimSuffix(rr.Body.String(), "\n")

		want := `{"requestName":"Test","something":"Hello"}`
		if got != want {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), want)
		}
	})

	t.Run("Body not ok. Returns 400", func(t *testing.T) {
		body := `{
			"requestName": "Test",
			"something" : "Hello
		   }`

		rr := httptest.NewRecorder()

		makeRequest(t, body, rr)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("Body doesn't have requestName property. Returns 400", func(t *testing.T) {
		body := `{
			"something" : "Hello"
		   }`

		rr := httptest.NewRecorder()

		makeRequest(t, body, rr)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})
}
