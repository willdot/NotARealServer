package persistrequests

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"testing"
)

func TestSave(t *testing.T) {

	testCases := []struct {
		Name          string
		InputData     bool
		ExpectedError error
	}{
		{
			Name:          "Data input valid, no error returned",
			InputData:     true,
			ExpectedError: nil,
		},
		{
			Name:          "Data input invalid, error returned",
			InputData:     false,
			ExpectedError: createTestMarshalError(),
		},
	}

	testObj := JSONPersist{}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			testData := createData(test.InputData)

			got := testObj.Save("test", "POST", testData, fakeFileReaderWriter{})

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

			gotString := fmt.Sprintf("%v", got)
			wantString := fmt.Sprintf("%v", want)

			if gotString != wantString {
				t.Errorf("Got %v, wanted %v", got, want)
			}

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

var errFakeError = errors.New("Fake error")

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

type Something struct {
	Language string
	Count    float64
}

func createData(goodData bool) interface{} {

	var result interface{}
	if goodData {
		thing := SavedRequest{
			RequestRoute:  "Hello",
			RequestMethod: "POST",
			Response: Something{
				Language: "Go",
				Count:    1.000,
			},
		}

		jsonByte, _ := json.Marshal(thing)
		json.Unmarshal(jsonByte, &result)

		return result
	}

	result = SavedRequest{
		RequestRoute:  "Hello",
		RequestMethod: "POST",
		Response: Something{
			Language: "Go",
			Count:    math.Inf(1),
		},
	}

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
