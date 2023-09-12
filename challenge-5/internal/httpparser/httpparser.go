package httpparser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"golang.org/x/exp/slog"
)

type HttpRequest struct {
	Method  string
	URI     string
	Version string
	Headers map[string]string
	Body    string
}

func Parse(reader io.Reader) (*HttpRequest, error) {
	var result HttpRequest
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	if scanner.Scan() {
		request := scanner.Text()
		fmt.Println(request)
		m, u, v, err := parseRequest(request)
		if err != nil {
			slog.Error(err.Error())
			return &result, err
		}
		result = HttpRequest{
			Method:  m,
			URI:     u,
			Version: v,
		}
	}

	headers := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if line == "" {
			break
		}
		tokens := strings.Split(line, ":")
		headers[tokens[0]] = strings.TrimSpace(tokens[1])
	}
	result.Headers = headers
	return &result, nil
}

func ReadBytes(reader io.Reader) ([]byte, error) {
	result := []byte{}
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	if scanner.Scan() {
		request := scanner.Bytes()
		result = append(result, request...)
		result = append(result, '\r', '\n')
		fmt.Println(string(request))
	}

	headers := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Bytes()
		result = append(result, line...)
		result = append(result, '\r', '\n')
		fmt.Println(string(line))
		if len(line) <= 0 {
			break
		}
		tokens := strings.Split(string(line), ":")
		headers[tokens[0]] = strings.TrimSpace(tokens[1])
	}

	body := []byte{}
	if headers["Content-Length"] != "" {
		length, err := strconv.Atoi(headers["Content-Length"])
		if err != nil {
			return nil, err
		}
		for scanner.Scan() {
			line := scanner.Bytes()
			body = append(body, line...)
			if len(body) >= length {
				break
			}
			body = append(body, '\r', '\n')
		}
	}

	fmt.Println(string(body))
	result = append(result, body...)
	return result, nil
}

func Read(reader io.Reader) ([]byte, error) {
	result := []byte{}

	r := bufio.NewReader(reader)
	for {
		bytes, _, err := r.ReadLine()
		if err != nil {
			slog.Error(err.Error())
			break
		}
		if len(bytes) <= 0 {
			result = append(result, '\r', '\n')
			break
		}
		result = append(result, bytes...)
		result = append(result, '\r', '\n')
	}

	fmt.Println(string(result))
	return result, nil
}

func parseRequest(req string) (method string, uri string, version string, err error) {
	tokens := strings.Split(req, " ")
	if len(tokens) < 3 {
		return "", "", "", fmt.Errorf("malformed request URI")
	}

	return tokens[0], tokens[1], tokens[2], nil
}
