package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type TokenType string

const (
	TokenOpenParen  TokenType = "openparen"
	TokenCloseParen           = "closeparen"
	TokenIdent                = "ident"
	TokenString               = "string"
	TokenNumber               = "number"
	TokenSpace                = "space"
)

var (
	ErrInvalidNumber   = errors.New("invalid number")
	ErrUnmatchedQuotes = errors.New("unmatched quotes for string")
)

type Lexer struct {
	scanner *bufio.Scanner
}

type Token struct {
	Pos   int
	Type  TokenType
	Value interface{}
}

type LexError struct {
	Pos int
	Err error
}

func (e LexError) Error() string {
	return fmt.Sprintf("%s at pos %d", e.Err.Error(), e.Pos)
}

func NewLexer(r io.Reader) Lexer {
	scanner := bufio.NewScanner(r)
	scanner.Split(tokenizerFunc)
	return Lexer{
		scanner: scanner,
	}
}

func (l *Lexer) Tokens() (tokens []Token, err error) {
	pos := 0

	for l.scanner.Scan() {
		var (
			token *Token
			r     rune
		)

		tokens = append(tokens, Token{})
		token = &tokens[len(tokens)-1]

		bytes := l.scanner.Bytes()
		if r, _ = utf8.DecodeRune(bytes); err != nil {
			return
		}

		token.Pos = pos
		pos += len(bytes)

		if r == '(' {
			token.Type = TokenOpenParen
		} else if r == ')' {
			token.Type = TokenCloseParen
		} else if unicode.IsSpace(r) {
			token.Type = TokenSpace
		} else {
			s := l.scanner.Text()
			if s[0] == '"' {
				token.Type = TokenString
				if len(s) == 1 || s[len(s)-1] != '"' {
					token.Value = s
					err = LexError{Pos: pos, Err: ErrUnmatchedQuotes}
					return
				}
				token.Value = s[1 : len(s)-1]
			} else if unicode.IsDigit(r) || (len(bytes) > 1 && (r == '+' || r == '-')) {
				var n float64
				token.Type = TokenNumber
				if n, err = strconv.ParseFloat(s, 64); err != nil {
					token.Value = s
					err = LexError{Pos: pos, Err: ErrInvalidNumber}
					return
				} else {
					token.Value = n
				}
			} else {
				token.Type = TokenIdent
				token.Value = s
			}
		}
	}

	return
}

func isTokenSeparator(r rune) bool {
	return unicode.IsSpace(r) || r == '(' || r == ')'
}

func tokenizerFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	r, width := utf8.DecodeRune(data)
	isExpectingQuote := r == '"'
	isExpectingSpace := unicode.IsSpace(r)

	if r == '(' || r == ')' {
		return width, data[:width], nil
	}

	for start := width; start < len(data); start += width {
		r, width = utf8.DecodeRune(data[start:])
		if isExpectingQuote {
			if r == '"' {
				return start + width, data[:start+width], nil
			}
		} else if isExpectingSpace {
			if !unicode.IsSpace(r) {
				return start, data[:start], nil
			}
		} else if isTokenSeparator(r) {
			return start, data[:start], nil
		}
	}

	// return the last token
	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}
