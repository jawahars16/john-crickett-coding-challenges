package parser

import (
	"container/list"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

type TokenType int

const (
	ObjectOpener      TokenType = 1
	ObjectCloser      TokenType = 2
	StringLiteral     TokenType = 3
	NumericLiteral    TokenType = 4
	KeyValueSeparator TokenType = 5
	BooleanLiteral    TokenType = 6
	ItemSepartor      TokenType = 7
	ArrayOpener       TokenType = 8
	ArrayCloser       TokenType = 9
	NullLiteral       TokenType = 10
)

const (
	ObjectContext int = 0
	ArrayContext  int = 1
)

const (
	trueLiteral  string = "true"
	falseLiteral string = "false"
	nullLiteral  string = "null"
)

type Token struct {
	Type  TokenType
	Value string
}

func Parse(reader io.Reader) error {
	tokens, err := tokenize(reader)
	if err != nil {
		return err
	}
	err = checkGrammar(tokens)
	if err != nil {
		return err
	}
	return nil
}

func tokenize(reader io.Reader) ([]Token, error) {
	buffer := make([]byte, 1024)
	var tokens []Token
	for {
		n, err := reader.Read(buffer)
		start := 0
		for start < n {
			advance, token, err := scan(buffer[start:], err == io.EOF)
			if err != nil {
				start = 0
				return tokens, err
			}
			if token != nil {
				tokens = append(tokens, *token)
			}
			start = start + advance
		}

		if err == io.EOF {
			break
		}
	}
	return tokens, nil
}

func scan(data []byte, atEOF bool) (advance int, token *Token, err error) {
	var current, width int
	// Skip all the spaces
	for current = 0; current < len(data); current += width {
		var r rune
		r, width = utf8.DecodeRune(data[current:])
		if !unicode.IsSpace(r) {
			break
		}
	}

	// Check if the current pointer is at expected token
	for current < len(data) {
		var r rune
		r, width = utf8.DecodeRune(data[current:])
		current += width

		if isObjectOpener(r) {
			return current, &Token{Type: ObjectOpener, Value: string(r)}, nil
		}

		if isObjectCloser(r) {
			return current, &Token{Type: ObjectCloser, Value: string(r)}, nil
		}

		if isArrayOpener(r) {
			return current, &Token{Type: ArrayOpener, Value: string(r)}, nil
		}

		if isArrayCloser(r) {
			return current, &Token{Type: ArrayCloser, Value: string(r)}, nil
		}

		if isQuote(r) {
			// we found a quote, now grab the literal
			literal, width, err := grabStringLiteral(data[current:])
			if err != nil {
				return current, nil, err
			}
			return current + width, &Token{Type: StringLiteral, Value: literal}, nil
		}

		if unicode.IsDigit(r) {
			literal, width, err := grabNumericLiteral(data[current-width:])
			if err != nil {
				return current, nil, err
			}
			return current + width, &Token{Type: NumericLiteral, Value: literal}, nil
		}

		if ok, bLiteral := getBooleanLiteral(data[current-1:]); ok {
			return current - 1 + len(bLiteral), &Token{Type: BooleanLiteral, Value: bLiteral}, nil
		}

		if ok, nLiteral := getNullLiteral(data[current-1:]); ok {
			return current - 1 + len(nLiteral), &Token{Type: NullLiteral, Value: nullLiteral}, nil
		}

		if isKeyValueSeprator(r) {
			return current, &Token{Type: KeyValueSeparator, Value: string(r)}, nil
		}

		if isItemSeparator(r) {
			return current, &Token{Type: ItemSepartor, Value: string(r)}, nil
		}

		if unicode.IsLetter(r) {
			return current, nil, fmt.Errorf("Unexpected token %s", string(r))
		}
	}

	return current, nil, nil
}

func grabStringLiteral(data []byte) (string, int, error) {
	current := 0
	var literal []rune
	for current < len(data) {
		// iterate until we find the closing quote
		r, width := utf8.DecodeRune(data[current:])
		current += width

		if unicode.IsControl(r) {
			return "", current, fmt.Errorf("Control character '%s' not allowed.", string(r))
		}

		if isQuote(r) { // if closing quote
			return string(literal), current, nil
		}

		literal = append(literal, r) // build the string until we find closing quote
	}
	return "", current, errors.New(`Invalid string literal. Expecting '"'`)
}

func grabNumericLiteral(data []byte) (string, int, error) {
	current := 0
	var literal []rune
	for current < len(data) {
		// iterate until we find the closing quote
		r, width := utf8.DecodeRune(data[current:])

		if !unicode.IsDigit(r) && r != '.' { // Let's stop here
			return string(literal), current - 1, nil
		}

		current += width
		literal = append(literal, r) // build the string until we find closing quote

		if unicode.IsSpace(r) {
			// 💣 boom !!! that's a space.
			return "", current, errors.New(fmt.Sprintf("Unexpected %s in literal", string(r)))
		}
	}
	return "", current, errors.New(`Invalid string literal. Expecting '"'`)
}

func linkedList(tokens []Token) *list.List {
	l := list.New()
	for _, token := range tokens {
		l.PushBack(token)
	}
	return l
}

func checkGrammar(tokens []Token) error {
	linkedTokens := linkedList(tokens)
	context := NewStack()

	for node := linkedTokens.Front(); node != nil; node = node.Next() {
		currentToken := node.Value.(Token)
		var (
			nextToken Token
			prevToken Token
		)

		if next := node.Next(); next != nil {
			nextToken = next.Value.(Token)
		}
		if prev := node.Prev(); prev != nil {
			prevToken = prev.Value.(Token)
		}

		switch currentToken.Type {
		case ObjectOpener:
			if nextToken.Type == 0 {
				return fmt.Errorf("Expected end of file.")
			}
			context.Push(ObjectContext)
			if node.Next() == nil {
				return fmt.Errorf("Invalid object expression.")
			}
			if nextToken.Type != ObjectCloser && nextToken.Type != StringLiteral {
				return fmt.Errorf("Unexpected token %s", nextToken.Value)
			}
		case ArrayOpener:
			if nextToken.Type == 0 {
				return fmt.Errorf("Expected end of file.")
			}
			context.Push(ArrayContext)
			if nextToken.Type != ArrayCloser &&
				!isValidValue(nextToken) {
				return fmt.Errorf("Invalid array.")
			}
		case KeyValueSeparator:
			if nextToken.Type == 0 {
				return fmt.Errorf("Expected end of file.")
			}
			if node.Prev() == nil || node.Next() == nil {
				return fmt.Errorf("Missing key or value in key value pair.")
			}
			if prevToken.Type != StringLiteral {
				return fmt.Errorf("Unexpected token. Expecting key for the key value pair.")
			}
			if !isValidValue(nextToken) {
				return fmt.Errorf("Unexpected token. Expecting value for the key value pair.")
			}
		case StringLiteral:
			currentContext := context.Peek()
			switch currentContext {
			case ObjectContext:
				if nextToken.Type != KeyValueSeparator && prevToken.Type != KeyValueSeparator {
					return fmt.Errorf("Invalid string literal in object expression.")
				}
			case ArrayContext:
				if prevToken.Type != ItemSepartor &&
					nextToken.Type != ItemSepartor &&
					nextToken.Type != ArrayCloser {
					return fmt.Errorf("Orphan string literal in array")
				}
			default:
				return fmt.Errorf("Orphan string literal")
			}
		case NumericLiteral:
			if nextToken.Type == 0 {
				return fmt.Errorf("Expected end of file.")
			}
			if strings.HasPrefix(currentToken.Value, "0") {
				return fmt.Errorf("Numeric literal cannot begin with 0.")
			}
			currentContext := context.Peek()
			switch currentContext {
			case ObjectContext:
				if prevToken.Type != KeyValueSeparator {
					return fmt.Errorf("Invalid numeric literal in object expression.")
				}
			case ArrayContext:
				if prevToken.Type != ItemSepartor &&
					nextToken.Type != ItemSepartor &&
					nextToken.Type != ArrayCloser {
					return fmt.Errorf("Orphan numeric literal in array")
				}
			default:
				return fmt.Errorf("Orphan numeric literal")
			}
		case ObjectCloser:
			if prevToken.Type == ItemSepartor || prevToken.Type == KeyValueSeparator {
				return fmt.Errorf("Unexpected end of file")
			}
			currentContext := context.Peek()
			if currentContext == ObjectContext {
				context.Pop()
			} else {
				return fmt.Errorf("Unexpected token %s.", currentToken.Value)
			}
		case ArrayCloser:
			if prevToken.Type == ItemSepartor || prevToken.Type == KeyValueSeparator {
				return fmt.Errorf("Unexpected end of file")
			}
			currentContext := context.Peek()
			if currentContext == ArrayContext {
				context.Pop()
			} else {
				return fmt.Errorf("Unexpected token %s.", currentToken.Value)
			}

		case ItemSepartor:
			if nextToken.Type == 0 {
				return fmt.Errorf("Expected end of file.")
			}
			currentContext := context.Peek()
			if currentContext == ObjectContext {
				if nextToken.Type != StringLiteral {
					return fmt.Errorf("Not expecting ',' here")
				}
			}
		}
	}

	return nil
}

func isValidValue(nextToken Token) bool {
	return nextToken.Type == StringLiteral ||
		nextToken.Type == NumericLiteral ||
		nextToken.Type == BooleanLiteral ||
		nextToken.Type == ObjectOpener ||
		nextToken.Type == ArrayOpener ||
		nextToken.Type == NullLiteral
}

func isQuote(r rune) bool {
	return r == '"'
}

func isObjectOpener(r rune) bool {
	return r == '{'
}

func isObjectCloser(r rune) bool {
	return r == '}'
}

func isArrayOpener(r rune) bool {
	return r == '['
}

func isArrayCloser(r rune) bool {
	return r == ']'
}

func isKeyValueSeprator(r rune) bool {
	return r == ':'
}

func isItemSeparator(r rune) bool {
	return r == ','
}

func isEscapeSequence(r rune) bool {
	return r == '\\'
}

func getBooleanLiteral(data []byte) (bool, string) {
	r, _ := utf8.DecodeRune(data)
	if r != 't' && r != 'f' {
		return false, ""
	}

	if len(data) > 4 {
		literal := string(data[:len(falseLiteral)])
		if literal == falseLiteral {
			return true, literal
		}
	}

	if len(data) > 3 {
		literal := string(data[:len(trueLiteral)])
		if literal == trueLiteral {
			return true, literal
		}
	}

	return false, ""
}

func getNullLiteral(data []byte) (bool, string) {
	if data[0] != 'n' {
		return false, ""
	}

	literal := data[:4]
	if string(literal) == "null" {
		return true, "null"
	}
	return false, ""
}
