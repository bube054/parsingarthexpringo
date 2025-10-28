package multipassparser

import (
	"testing"

	"github.com/bube054/parsingarthexpringo/ast"
	"github.com/bube054/parsingarthexpringo/lexer"
)

func TestParser(t *testing.T) {
	var tests = []struct {
		input    string
		expected ast.Node
		opPrec   OperatorPrecedence
	}{
		// Parse single arithmetic or algebraic expressions
		{
			input:    "3",
			expected: ast.NewOperand(lexer.NUM, "3"),
			opPrec:   BIDMAS,
		},
		{
			input:    "a",
			expected: ast.NewOperand(lexer.ALPHA, "a"),
			opPrec:   BIDMAS,
		},

		// Parse single arithmetic or algebraic expressions with parenthesis
		{
			input:    "( 3 )",
			expected: ast.NewOperand(lexer.NUM, "3"),
			opPrec:   BIDMAS,
		},
		{
			input:    "( a )",
			expected: ast.NewOperand(lexer.ALPHA, "a"),
			opPrec:   BIDMAS,
		},

		// Parse binary arithmetic or algebraic expressions
		{
			input:    "3 + 9",
			expected: ast.NewBinaryExpr(lexer.ADD, ast.NewOperand(lexer.NUM, "3"), ast.NewOperand(lexer.NUM, "9")),
			opPrec:   BIDMAS,
		},
		{
			input:    "4 - 5",
			expected: ast.NewBinaryExpr(lexer.SUB, ast.NewOperand(lexer.NUM, "4"), ast.NewOperand(lexer.NUM, "5")),
			opPrec:   BIDMAS,
		},
		{
			input:    "x / y",
			expected: ast.NewBinaryExpr(lexer.DIV, ast.NewOperand(lexer.ALPHA, "x"), ast.NewOperand(lexer.ALPHA, "y")),
			opPrec:   BIDMAS,
		},
		{
			input:    "a * b",
			expected: ast.NewBinaryExpr(lexer.MUL, ast.NewOperand(lexer.ALPHA, "a"), ast.NewOperand(lexer.ALPHA, "b")),
			opPrec:   BIDMAS,
		},
		{
			input:    "p ^ q",
			expected: ast.NewBinaryExpr(lexer.POW, ast.NewOperand(lexer.ALPHA, "p"), ast.NewOperand(lexer.ALPHA, "q")),
			opPrec:   BIDMAS,
		},

		// Parse binary arithmetic or algebraic expressions with parenthesis
		{
			input:    "( 15 + 1 )",
			expected: ast.NewBinaryExpr(lexer.ADD, ast.NewOperand(lexer.NUM, "15"), ast.NewOperand(lexer.NUM, "1")),
			opPrec:   BIDMAS,
		},
		{
			input:    "( 3 - 1 )",
			expected: ast.NewBinaryExpr(lexer.SUB, ast.NewOperand(lexer.NUM, "3"), ast.NewOperand(lexer.NUM, "1")),
			opPrec:   BIDMAS,
		},
		{
			input:    "( x / y )",
			expected: ast.NewBinaryExpr(lexer.DIV, ast.NewOperand(lexer.ALPHA, "x"), ast.NewOperand(lexer.ALPHA, "y")),
			opPrec:   BIDMAS,
		},
		{
			input:    "( a * b )",
			expected: ast.NewBinaryExpr(lexer.MUL, ast.NewOperand(lexer.ALPHA, "a"), ast.NewOperand(lexer.ALPHA, "b")),
			opPrec:   BIDMAS,
		},
		{
			input:    "( p ^ q )",
			expected: ast.NewBinaryExpr(lexer.POW, ast.NewOperand(lexer.ALPHA, "p"), ast.NewOperand(lexer.ALPHA, "q")),
			opPrec:   BIDMAS,
		},

		// Parse complex arithmetic or algebraic expressions
		{
			input:    "3 + 6 * 7",
			expected: ast.NewBinaryExpr(lexer.ADD, ast.NewOperand(lexer.NUM, "3"), ast.NewBinaryExpr(lexer.MUL, ast.NewOperand(lexer.NUM, "6"), ast.NewOperand(lexer.NUM, "7"))),
			opPrec:   BIDMAS,
		},
		{
			input:    "12 - 8 / 2",
			expected: ast.NewBinaryExpr(lexer.SUB, ast.NewOperand(lexer.NUM, "12"), ast.NewBinaryExpr(lexer.DIV, ast.NewOperand(lexer.NUM, "8"), ast.NewOperand(lexer.NUM, "2"))),
			opPrec:   BIDMAS,
		},
		{
			input:    "4 * 3 ^ 2",
			expected: ast.NewBinaryExpr(lexer.MUL, ast.NewOperand(lexer.NUM, "4"), ast.NewBinaryExpr(lexer.POW, ast.NewOperand(lexer.NUM, "3"), ast.NewOperand(lexer.NUM, "2"))),
			opPrec:   BIDMAS,
		},
		{
			input:    "2 + 4 * 6 - 3",
			expected: ast.NewBinaryExpr(lexer.SUB, ast.NewBinaryExpr(lexer.ADD, ast.NewOperand(lexer.NUM, "2"), ast.NewBinaryExpr(lexer.MUL, ast.NewOperand(lexer.NUM, "4"), ast.NewOperand(lexer.NUM, "6"))), ast.NewOperand(lexer.NUM, "3")),
			opPrec:   BIDMAS,
		},

		// Parse complex arithmetic or algebraic expressions with parenthesis
		{
			input:    "3 * ( 2 + 4 )",
			expected: ast.NewBinaryExpr(lexer.MUL, ast.NewOperand(lexer.NUM, "3"), ast.NewBinaryExpr(lexer.ADD, ast.NewOperand(lexer.NUM, "2"), ast.NewOperand(lexer.NUM, "4"))),
			opPrec:   BIDMAS,
		},
		{
			input:    "3 ( 2 + 4 )",
			expected: ast.NewBinaryExpr(lexer.MUL, ast.NewOperand(lexer.NUM, "3"), ast.NewBinaryExpr(lexer.ADD, ast.NewOperand(lexer.NUM, "2"), ast.NewOperand(lexer.NUM, "4"))),
			opPrec:   BIDMAS,
		},
		{
			input:    "( a - b ) * c",
			expected: ast.NewBinaryExpr(lexer.MUL, ast.NewBinaryExpr(lexer.SUB, ast.NewOperand(lexer.ALPHA, "a"), ast.NewOperand(lexer.ALPHA, "b")), ast.NewOperand(lexer.ALPHA, "c")),
			opPrec:   BIDMAS,
		},
		{
			input:    "( a - b ) c",
			expected: ast.NewBinaryExpr(lexer.MUL, ast.NewBinaryExpr(lexer.SUB, ast.NewOperand(lexer.ALPHA, "a"), ast.NewOperand(lexer.ALPHA, "b")), ast.NewOperand(lexer.ALPHA, "c")),
			opPrec:   BIDMAS,
		},
	}

	// nbe := ast.NewBinaryExpr()

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			tokens := lexer.Lexer(test.input)
			parser := NewParser(tokens, test.opPrec)
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
