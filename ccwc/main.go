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

var stream io.ReadWriteCloser

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

	// input reader -> pipe 1 writer -> pipe 1 reader -> pipe 2 writer -> pipe 2 reader
	// 			   \-> read and count lines       \-> read and count words
	cpr, cpw := io.Pipe()
	ltee := io.TeeReader(reader, cpw)

	wpr, wpw := io.Pipe()
	ctee := io.TeeReader(cpr, wpw)

	mpr, mpw := io.Pipe()
	wtee := io.TeeReader(wpr, mpw)

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		defer cpw.Close()
		count("l", ltee, &lines)
	}()

	go func() {
		defer wg.Done()
		defer wpw.Close()
		count("w", ctee, &words)
	}()

	go func() {
		defer wg.Done()
		defer mpw.Close()
		count("c", wtee, &bytes)
	}()

	go func() {
		defer wg.Done()
		count("m", mpr, &chars)
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
