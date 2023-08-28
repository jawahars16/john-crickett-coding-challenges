package parser

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func Test_simpleParse(t *testing.T) {
	testcases := []struct {
		input string
		valid bool
	}{
		{
			input: "{}",
			valid: true,
		},
		{
			input: "{",
			valid: false,
		},
		{
			input: `{"key":"value"}`,
			valid: true,
		},
		{
			input: `{"key"}`,
			valid: false,
		},
		{
			input: `{"key":12}`,
			valid: true,
		},
		{
			input: `{"key":312.212}`,
			valid: true,
		},
		{
			input: `{33:312.212}`,
			valid: false,
		},
		{
			input: `
			{
				"string": "value",
				"numeric":100.1,
				"boolean": true,
				"object":{},
				"array":[[]]
			}
			`,
			valid: true,
		},
		{
			input: `{"Numbers cannot be hex": 0x14}`,
			valid: false,
		},
		{
			input: `["123",]`,
			valid: false,
		},
		{
			input: `{"key":null}`,
			valid: true,
		},
	}

	for k, testCase := range testcases {
		err := Parse(strings.NewReader(testCase.input))
		actual := err == nil
		if actual != testCase.valid {
			t.Error(k, testCase.input, "Expected:", testCase.valid, "Actual:", actual, err)
		}
	}
}

func Test_parseFunctionOnTestData(t *testing.T) {
	// iterate through files in testdata
	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Error(err)
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "s") {
			// skip tests
			continue
		}
		func() {
			reader, err := os.Open("./testdata/" + file.Name())
			defer reader.Close()
			if err != nil {
				t.Error(err)
			}
			err = Parse(reader)
			if err != nil && strings.Contains(file.Name(), "pass") {
				t.Error("Expecting no error", file.Name(), err)
			}
			if err == nil && strings.Contains(file.Name(), "fail") {
				t.Error("Expecting error", file.Name(), err)
			}
		}()
	}
}
