package parsingarthexpringo

import (
	"testing"
)

func TestParser(t *testing.T) {
	var tests = []struct {
		input    string
		expected Node
		opPrec   OperatorPrecedence
	}{
		// // Parse single arithmetic or algebraic expressions
		// {
		// 	input:    "3",
		// 	expected: NewOperand(NUM, "3"),
		// },
		// {
		// 	input:    "a",
		// 	expected: NewOperand(ALPHA, "a"),
		// },

		// Parse single arithmetic or algebraic expressions with parenthesis
		// {
		// 	input:    "( 3 )",
		// 	expected: NewOperand(NUM, "3"),
		// },
		// {
		// 	input:    "( a )",
		// 	expected: NewOperand(ALPHA, "a"),
		// },

		// Parse binary arithmetic or algebraic expressions
		// {
		// 	input:    "3 + 9",
		// 	expected: NewBinaryExpr(ADD, NewOperand(NUM, "3"), NewOperand(NUM, "9")),
		// },
		// {
		// 	input:    "4 - 5",
		// 	expected: NewBinaryExpr(SUB, NewOperand(NUM, "4"), NewOperand(NUM, "5")),
		// },
		// {
		// 	input:    "x / y",
		// 	expected: NewBinaryExpr(DIV, NewOperand(ALPHA, "x"), NewOperand(ALPHA, "y")),
		// },
		// {
		// 	input:    "a * b",
		// 	expected: NewBinaryExpr(MUL, NewOperand(ALPHA, "a"), NewOperand(ALPHA, "b")),
		// },
		{
			input:    "p ^ q",
			expected: NewBinaryExpr(EXP, NewOperand(ALPHA, "p"), NewOperand(ALPHA, "q")),
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			tokens := Lexer(test.input)
			parser := NewParser(tokens, BIDMAS)
			actual, _ := parser.Parse()

			if actual == nil && test.expected != nil {
				t.Errorf("got nil, want %v", test.expected)
				return
			}
			if actual != nil && test.expected == nil {
				t.Errorf("got %v, want nil", actual)
				return
			}
			if actual != nil && !actual.Equals(test.expected) {
				t.Errorf("got %v, want %v", actual, test.expected)
			}

		})
	}
}
