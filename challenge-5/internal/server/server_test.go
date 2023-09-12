package server

import (
	"bytes"
	"testing"
)

func Test_handle(t *testing.T) {
	reader := bytes.NewBuffer([]byte("GET /test HTTP/1.1\r\nContent-Length: 5\r\n\r\nHello"))
	writer := bytes.NewBuffer([]byte{})
	handle(reader, writer)
	expected := []byte("GET /test HTTP/1.1\r\nContent-Length: 5\r\n\r\nHello")
	actual := reader.Bytes()
	if !bytes.Equal(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
