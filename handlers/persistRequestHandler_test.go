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

	return nil, nil
}

var testThing = PersistServer{
	FileWriter: fakeFileReaderWriter{},
	FileReader: fakeFileReaderWriter{},
	LoadSaver:  persistrequests.JSONPersist{},
}

type fake struct {
	requestName string
	something   string
}

func TestNewPersistServer(t *testing.T) {

	got := NewPersistServer()

	want := PersistServer{
		LoadSaver:  persistrequests.JSONPersist{},
		FileWriter: persistrequests.FileWriter{},
		FileReader: persistrequests.FileReader{},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}

var makeRequest = func(t *testing.T, url, body string, handler http.Handler, rr *httptest.ResponseRecorder) {

	t.Helper()

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	handler.ServeHTTP(rr, req)
}

func TestRetreiveRequestHandler(t *testing.T) {

	handler := mux.NewRouter()
	handler.HandleFunc("/{request}", testThing.RetreiveRequestHandler())

	t.Run("Param ok. Request Exists. Request returned", func(t *testing.T) {
		body := ""

		rr := httptest.NewRecorder()

		makeRequest(t, "/test", body, handler, rr)

	})
}

func TestSaveRequestHandler(t *testing.T) {
	handler := testThing.SaveRequestHandler()

	t.Run("Body ok. Returns 200", func(t *testing.T) {
		body := `{
			"requestName": "Test",
			"something" : "Hello"
		   }`

		rr := httptest.NewRecorder()

		makeRequest(t, "save", body, handler, rr)

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

		makeRequest(t, "save", body, handler, rr)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("Body doesn't have requestName property. Returns 400", func(t *testing.T) {
		body := `{
			"something" : "Hello"
		   }`

		rr := httptest.NewRecorder()

		makeRequest(t, "save", body, handler, rr)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})
}
