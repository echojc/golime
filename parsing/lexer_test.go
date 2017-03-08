package parsing

import (
	"reflect"
	"strings"
	"testing"

	"github.com/kr/pretty"
)

func TestLexer(t *testing.T) {
	t.Run("doesn't care about balanced brackets", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected []Token
		}{
			{`(`, []Token{
				{0, TokenOpenParen, nil},
			}},
			{`)`, []Token{
				{0, TokenCloseParen, nil},
			}},
			{`)(`, []Token{
				{0, TokenCloseParen, nil},
				{1, TokenOpenParen, nil},
			}},
			{`(()`, []Token{
				{0, TokenOpenParen, nil},
				{1, TokenOpenParen, nil},
				{2, TokenCloseParen, nil},
			}},
		}
		for _, testCase := range testCases {
			l := NewLexer(strings.NewReader(testCase.input))

			tokens, err := l.Tokens()
			if err != nil {
				t.Errorf("couldn't parse input [err: %+v]", err)
				continue
			}

			if !reflect.DeepEqual(testCase.expected, tokens) {
				t.Errorf(
					"expected and actual differ [input: %s] [diff: %s]",
					testCase.input,
					pretty.Diff(testCase.expected, tokens),
				)
			}
		}
	})

	t.Run("parses strings", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected []Token
		}{
			{`"a"`, []Token{
				{0, TokenString, "a"},
			}},
			{`"abc"`, []Token{
				{0, TokenString, "abc"},
			}},
			{`""`, []Token{
				{0, TokenString, ""},
			}},
			{`" "`, []Token{
				{0, TokenString, " "},
			}},
			{`"   "`, []Token{
				{0, TokenString, "   "},
			}},
			{`"( )"`, []Token{
				{0, TokenString, "( )"},
			}},
			{`"ab '!@#$%^&*()"`, []Token{
				{0, TokenString, "ab '!@#$%^&*()"},
			}},
			{`"日本語"`, []Token{
				{0, TokenString, "日本語"},
			}},
		}
		for _, testCase := range testCases {
			l := NewLexer(strings.NewReader(testCase.input))

			tokens, err := l.Tokens()
			if err != nil {
				t.Errorf("couldn't parse input [err: %+v]", err)
				continue
			}

			if !reflect.DeepEqual(testCase.expected, tokens) {
				t.Errorf(
					"expected and actual differ [input: %s] [diff: %s]",
					testCase.input,
					pretty.Diff(testCase.expected, tokens),
				)
			}
		}
	})
}
