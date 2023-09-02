package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_cut(t *testing.T) {
	cases := []struct {
		input     string
		fields    []int
		output    string
		delimiter string
	}{
		{
			input:  `f0	f1	f2	f3	f4`,
			fields: []int{1},
			output: "f0\n",
		},
		{
			input:  `f0	f1	f2	f3	f4`,
			fields: []int{5},
			output: "f4\n",
		},
		{
			input:     `f0	f1:f2:f3,f4`,
			fields:    []int{2},
			output:    "f2\n",
			delimiter: ":",
		},
		{
			input:     `f1:f2:f3,f4`,
			fields:    []int{2},
			output:    "f2\n",
			delimiter: ":",
		},
		{
			input:     `f1:f2:f3`,
			fields:    []int{2},
			output:    "f1:f2:f3\n",
			delimiter: ",",
		},
		{
			input:     `f1:f2:f3`,
			fields:    []int{1, 2},
			output:    "f1:f2\n",
			delimiter: ":",
		},
	}

	for i, c := range cases {
		reader := strings.NewReader(c.input)
		var buffer bytes.Buffer
		cut(reader, &buffer, c.fields, c.delimiter)
		actual := buffer.String()

		if diff := cmp.Diff(c.output, actual); diff != "" {
			t.Error("Case no.", i+1, diff)
		}
	}
}
