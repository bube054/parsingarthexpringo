package parsingarthexpringo

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type TokenKind int

const (
	NUM TokenKind = iota
	ALPHA
	ADD
	SUB
	DIV
	MUL
	EXP
	LBRACKET
	RBRACKET
	ILLEGAL
)

func (tk TokenKind) String() string {
	switch tk {
	case NUM:
		return "NUM"
	case ALPHA:
		return "ALPHA"
	case ADD:
		return "ADD"
	case SUB:
		return "SUB"
	case DIV:
		return "DIV"
	case MUL:
		return "MUL"
	case EXP:
		return "EXP"
	case LBRACKET:
		return "LBRACKET"
	case RBRACKET:
		return "RBRACKET"
	case ILLEGAL:
		return "ILLEGAL"
	default:
		return fmt.Sprintf("Unknown TokenKind(%d)", int(tk))
	}
}

type TokenCategory int

const (
	CATEGORY_VALUE TokenCategory = iota
	CATEGORY_OPERATOR
	CATEGORY_PAREN
	CATEGORY_ILLEGAL
)

func categorize(kind TokenKind) TokenCategory {
	switch kind {
	case NUM, ALPHA:
		return CATEGORY_VALUE
	case ADD, SUB, MUL, DIV, EXP:
		return CATEGORY_OPERATOR
	case LBRACKET, RBRACKET:
		return CATEGORY_PAREN
	default:
		return CATEGORY_ILLEGAL
	}
}

type Token struct {
	Kind     TokenKind
	Value    string
	Category TokenCategory
}

func NewToken(kind TokenKind, value string) Token {
	return Token{
		Kind:     kind,
		Value:    value,
		Category: categorize(kind),
	}
}

type Tokens []Token

func (tks *Tokens) Push(tk Token) {
	*tks = append(*tks, tk)
}

func (tks *Tokens) Equal(other *Tokens) bool {
	if tks == nil && other == nil {
		return true
	}
	if tks == nil || other == nil {
		return false
	}
	if len(*tks) != len(*other) {
		return false
	}
	for i := range *tks {
		if !(*tks)[i].Equal(&(*other)[i]) {
			return false
		}
	}
	return true
}

func (t *Token) Equal(other *Token) bool {
	if t == nil && other == nil {
		return true
	}
	if t == nil || other == nil {
		return false
	}
	return t.Kind == other.Kind &&
		t.Value == other.Value &&
		t.Category == other.Category
}

func Lexer(input string) Tokens {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)

	tokens := Tokens{}

	for scanner.Scan() {
		text := scanner.Text()
		switch text {
		case "+":
			tokens.Push(NewToken(ADD, text))
		case "-":
			tokens.Push(NewToken(SUB, text))
		case "*":
			tokens.Push(NewToken(MUL, text))
		case "/":
			tokens.Push(NewToken(DIV, text))
		case "^":
			tokens.Push(NewToken(EXP, text))
		case "(":
			tokens.Push(NewToken(LBRACKET, text))
		case ")":
			tokens.Push(NewToken(RBRACKET, text))
		default:
			if IsLetter(text) {
				tokens.Push(NewToken(ALPHA, text))
				continue
			}

			if IsNumber(text) {
				tokens.Push(NewToken(NUM, text))
				continue
			}

			tokens.Push(NewToken(ILLEGAL, text))
			return tokens
		}
	}

	return tokens
}

func IsLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func IsNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

