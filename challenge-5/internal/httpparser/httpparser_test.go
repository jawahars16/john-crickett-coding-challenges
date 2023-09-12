package httpparser_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jawahars16/john-crickett-coding-challenges/challenge-5/internal/httpparser"
)

func Test_read(t *testing.T) {
	request := "GET /test HTTP/1.1\r\nContent-Length: 5\r\n\r\nHello"
	expected := []byte("GET /test HTTP/1.1\r\nContent-Length: 5\r\n\r\nHello")

	reader := strings.NewReader(request)
	actual, err := httpparser.ReadBytes(reader)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}
