package requestextractor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func getSavedRequest() SavedRequest {
	byteValue, _ := ioutil.ReadFile("ExampleSavedRequest.json")

	var extracted SavedRequest

	json.Unmarshal(byteValue, &extracted)

	return extracted
}

func TestExtract(t *testing.T) {

	t.Run("POST method. Returns the post request part of saved file", func(t *testing.T) {

		method := "POST"
		savedRequest := getSavedRequest()

		result, _ := Extract(method, savedRequest)

		fmt.Println(result)
		fmt.Print("s")

		if fmt.Sprintf("%v", result) != "{POST map[Id:1111]}" {
			t.Errorf("nope")
		}
	})
}
