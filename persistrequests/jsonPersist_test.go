package persistrequests

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"testing"
)

var errFakeError = errors.New("Fake error")

type fakeFileReaderWriter struct {
}

func (f fakeFileReaderWriter) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return nil
}

func (f fakeFileReaderWriter) ReadFile(filename string) ([]byte, error) {

	// For when the file can't be read such as doesn't exist
	if filename == "POST-bad file.json" {
		return nil, errFakeError
	}

	// For when the data inside a json file is not correct such as missing a " somewhere
	if filename == "POST-bad.json" {
		result := []byte(testBadJSONFile)
		return result, nil
	}

	result := []byte(testGoodJSONFile)

	return result, nil
}

type TestStruct struct {
	RequestRoute  string
	RequestMethod string
	Request       interface{}
}

type TestSubStruct struct {
	Language string
	Count    float64
}

var testGoodJSONFile = `{
	"RequestRoute": "Hello",
	"RequestMethod": "POST",
	"Request": {
	 "Count": "1",
	 "Language": "Go"
	}
   }`

var testBadJSONFile = `{
	"RequestRoute": "Hello",
	"RequestMethod": "POST",
	"Request": {
	 "Count": "1,
	 "Language": "Go"
	}
   }`

func createData(good bool) interface{} {

	var result interface{}
	if good {
		thing := TestStruct{
			RequestRoute:  "Hello",
			RequestMethod: "POST",
			Request: TestSubStruct{
				Language: "Go",
				Count:    1.000,
			},
		}

		jsonByte, _ := json.Marshal(thing)
		json.Unmarshal(jsonByte, &result)

		return result
	}

	result = TestStruct{
		RequestRoute:  "Hello",
		RequestMethod: "POST",
		Request: TestSubStruct{
			Language: "Go",
			Count:    math.Inf(1),
		},
	}

	return result
}

/*func createData(good bool) map[string]interface{} {
	result := make(map[string]interface{})

	if good {
		thing := TestStruct{
			Thing: "Hello",
			AnotherThing: TestSubStruct{
				Language: "Go",
				Count:    1.000,
			},
		}

		jsonByte, _ := json.Marshal(thing)
		json.Unmarshal(jsonByte, &result)

		return result
	}

	result["test"] = TestStruct{
		Thing: "Hello",
		AnotherThing: TestSubStruct{
			Language: "Go",
			Count:    math.Inf(1),
		},
	}

	return result

}*/

// This creates the same error that will be returned from json.Marshal when sending in data that isn't valid
func createTestMarshalError() error {
	_, err := json.Marshal(math.Inf(1))
	return err.(*json.UnsupportedValueError)
}

// This creates the same error that will be returned from json.Unmarshal when sending in data that isn't valid
func createUnmarshalError() error {

	a := []byte(testBadJSONFile)
	var b map[string]interface{}
	err := json.Unmarshal(a, &b)

	return err.(*json.SyntaxError)
}

func TestSave(t *testing.T) {

	testCases := []struct {
		Name          string
		InputData     bool
		ExpectedError error
	}{
		{
			"Data input valid, no error returned",
			true,
			nil,
		},
		{
			"Data input invalid, error returned",
			false,
			createTestMarshalError(),
		},
	}

	testObj := JSONPersist{}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			testData := createData(test.InputData)

			got := testObj.Save("test", "POST", testData, fakeFileReaderWriter{})

			AssertErrors(got, test.ExpectedError, t)
		})
	}
}

func TestLoad(t *testing.T) {

	testCases := []struct {
		Name          string
		RequestRoute  string
		RequestMethod string
		OutputData    interface{}
		ExpectedError error
	}{
		{
			"File valid, data returned, no error returned",
			"good",
			"POST",
			TestSubStruct{
				Count:    1,
				Language: "Go",
			},
			nil,
		},
		{
			"File invalid, no data returned, error returned",
			"bad file",
			"POST",
			nil,
			errFakeError,
		},
		{
			"File valid, file data invalid, no data returned, error returned",
			"bad",
			"POST",
			nil,
			createUnmarshalError(),
		},
	}

	testObj := JSONPersist{}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			got, err := testObj.Load(test.RequestRoute, test.RequestMethod, fakeFileReaderWriter{})

			AssertErrors(err, test.ExpectedError, t)

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

func AssertErrors(got, want error, t *testing.T) {
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
