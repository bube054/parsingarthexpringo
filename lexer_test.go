package parsingarthexpringo

import (
	"reflect"
	"testing"
)

func TestLexer(t *testing.T) {
	var tests = []struct {
		input    string
		expected Tokens
	}{
		// {input: "1 + 2 - 3 * 4 / 5", expected: Tokens{{NUM, "1"}, {ADD, "+"}, {NUM, "2"}, {SUB, "-"}, {NUM, "3"}, {MUL, "*"}, {NUM, "4"}, {DIV, "/"}, {NUM, "5"}}},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			actual := Lexer(test.input)

			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf("got %v, want %v", actual, test.expected)
			}
		})
	}
}
