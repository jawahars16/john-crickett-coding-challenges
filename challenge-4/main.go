package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	invalidFieldInput = errors.New("Invalid field")
)

func cut(reader io.Reader, writer io.Writer, fields []int, delimiter string) {
	var result string
	if delimiter == "" {
		delimiter = "\t"
	}

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, delimiter)
		switch {
		case len(tokens) == 1:
			result = tokens[0]
		case len(fields) > 1:
			str := []string{}
			for _, field := range fields {
				str = append(str, tokens[field-1])
			}
			result = strings.Join(str, delimiter)
		default:
			result = tokens[fields[0]-1]
		}
		writer.Write([]byte(result + "\n"))
	}
}

func main() {
	fieldsStr := flag.String("f", "", "-f '1,2'")
	delimiter := flag.String("d", "\t", "-d ','")
	flag.Parse()
	fileName := flag.Arg(0)

	var (
		reader *os.File
		err    error
	)

	if fileName != "" && fileName != "-" {
		reader, err = os.Open(fileName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer reader.Close()
	} else {
		reader = os.Stdin
	}
	fieldsArr := strings.FieldsFunc(*fieldsStr, func(r rune) bool {
		return r == ',' || r == ' '
	})

	var fields []int
	for _, f := range fieldsArr {
		value, err := strconv.Atoi(f)
		if err != nil {
			fmt.Println(invalidFieldInput)
			os.Exit(1)
		}
		fields = append(fields, value)
	}

	cut(reader, os.Stdout, fields, *delimiter)
}
