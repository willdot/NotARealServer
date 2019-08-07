package main

import "testing"

func TestRequestPathTrailingSlash(t *testing.T) {

	testCases := []struct {
		Name           string
		Input          string
		ExpectedOutput string
	}{
		{
			Name:           "Has trailing slash, returned unchanged",
			Input:          "request/",
			ExpectedOutput: "request/",
		},
		{
			Name:           "Does not have trailing slash, returned with slash",
			Input:          "request",
			ExpectedOutput: "request/",
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			validateRequestDirectory(&test.Input)

			if test.Input != test.ExpectedOutput {
				t.Errorf("got %v want %v", test.Input, test.ExpectedOutput)
			}
		})
	}

}
