package main

import (
	"os"
	"strings"
	"testing"
)

func Test_count(t *testing.T) {
	testcases := map[string]struct {
		expected int
		argument string
		input    string
	}{
		"Should return number of bytes": {
			argument: "c",
			input:    "one two three ท",
			expected: 17,
		},
		"Should return number of words": {
			argument: "w",
			input:    "one two three",
			expected: 3,
		},
		"Should return number of lines": {
			argument: "l",
			input:    "one two three\nfour five six",
			expected: 2,
		},
		"Should return number of characters": {
			argument: "m",
			input:    "one two three ท",
			expected: 15,
		},
	}

	for k, v := range testcases {
		reader := strings.NewReader(v.input)
		var actual int
		err := count(v.argument, reader, &actual)
		if err != nil {
			t.Error(k, err)
		}
		if actual != v.expected {
			t.Error(k, "Expected:", v.expected, "Actual", actual)
		}
	}
}

func Test_CountFrom(t *testing.T) {
	reader, err := os.Open("file.txt")
	if err != nil {
		t.Error(err)
	}

	actualLines, actualWords, actualChars, actualBytes, err := countFrom(reader)
	if err != nil {
		t.Error(err)
	}

	if 23 != actualLines {
		t.Error("Expected lines", 22, "Actual lines", actualLines)
	}

	if 549 != actualWords {
		t.Error("Expected words", 549, "Actual words", actualWords)
	}

	if 3568 != actualBytes {
		t.Error("Expected bytes", 3568, "Actual bytes", actualBytes)
	}

	if 3566 != actualChars {
		t.Error("Expected chars", 3566, "Actual chars", actualChars)
	}
}
