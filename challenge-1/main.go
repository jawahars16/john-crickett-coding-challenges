package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
)

func main() {
	var (
		reader io.Reader
		err    error
	)

	filePath, arg := parseInput()
	if filePath != "" {
		reader, err = os.Open(filePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		meta, err := os.Stdin.Stat()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if meta.Size() > 0 {
			reader = os.Stdin
		}
	}

	if reader == nil {
		fmt.Println("No source found")
		os.Exit(1)
	}

	lines, words, chars, bytes, err := countFrom(reader)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch arg {
	case "l":
		fmt.Println(lines, filePath)
	case "w":
		fmt.Println(words, filePath)
	case "c":
		fmt.Println(bytes, filePath)
	case "m":
		fmt.Println(chars, filePath)
	default:
		fmt.Println(lines, words, chars, bytes, filePath)
	}
}

func parseInput() (string, string) {
	var (
		filePath string
		argument string
	)

	initializeFlags()
	filePath = flag.Arg(0)

	if flag.NFlag() > 0 {
		flag.Visit(func(f *flag.Flag) {
			argument = f.Name
			filePath = f.Value.String()
		})
	}

	return filePath, argument
}

func countFrom(reader io.Reader) (int, int, int, int, error) {
	var (
		lines int
		words int
		chars int
		bytes int
	)

	pr1, pw1 := io.Pipe()
	pr2, pw2 := io.Pipe()
	pr3, pw3 := io.Pipe()
	pr4, pw4 := io.Pipe()

	var wg sync.WaitGroup
	wg.Add(5)

	writer := io.MultiWriter(pw1, pw2, pw3, pw4)
	go func() {
		defer wg.Done()
		defer pw1.Close()
		defer pw2.Close()
		defer pw3.Close()
		defer pw4.Close()
		io.Copy(writer, reader)
	}()

	go func() {
		defer wg.Done()
		defer pw1.Close()
		count("l", pr1, &lines)
	}()

	go func() {
		defer wg.Done()
		defer pw2.Close()
		count("w", pr2, &words)
	}()

	go func() {
		defer wg.Done()
		defer pw3.Close()
		count("c", pr3, &bytes)
	}()

	go func() {
		defer wg.Done()
		defer pw4.Close()
		count("m", pr4, &chars)
	}()

	wg.Wait()
	return lines, words, chars, bytes, nil
}

func count(argument string, reader io.Reader, result *int) error {
	var splitFunc bufio.SplitFunc

	switch argument {
	case "c":
		splitFunc = bufio.ScanBytes
	case "w":
		splitFunc = bufio.ScanWords
	case "l":
		splitFunc = bufio.ScanLines
	case "m":
		splitFunc = bufio.ScanRunes
	default:
		return errors.New(fmt.Sprintln("Invalid argument:", argument))
	}

	scanner := bufio.NewScanner(reader)
	scanner.Split(splitFunc)
	count := 0
	for scanner.Scan() {
		scanner.Text()
		count += 1
	}
	*result = count
	return nil
}

func availableFlags() map[string]string {
	return map[string]string{
		"c": "Count bytes in the given file",
		"w": "Count words in the given file",
		"l": "Count lines in the given file",
		"m": "Count characters in the given file",
	}
}

func initializeFlags() {
	for k, v := range availableFlags() {
		flag.String(k, "", v)
	}
	flag.Parse()
}
