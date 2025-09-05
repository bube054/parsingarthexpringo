package parsingarthexpringo

import (
	"errors"
	"fmt"

	"github.com/bube054/go-js-array-methods/array"
)

var (
	ErrInvalidExpression       = errors.New("invalid expression")
	ErrCouldNotParseExpression = errors.New("could not parse expression")
	ErrInvalidReducingTokens   = errors.New("invalid reducing tokens type")
	ErrNoMatchingBracket       = errors.New("no matching right bracket")
)

// type OperatorPrecedence map[TokenKind]int
type OperatorPrecedence []TokenKind

var BIDMAS = OperatorPrecedence{EXP, MUL, DIV, SUB, ADD}

type Parser struct {
	originalTokens             Tokens
	reducingTokens             []any
	originalOperatorPrecedence OperatorPrecedence
	operatorPrecedence         OperatorPrecedence
}

func tokensToAny(tks Tokens) []any {
	reducing := make([]any, len(tks))
	for i, v := range tks {
		reducing[i] = v
	}
	return reducing
}

func anyToTokens(reducing []any) Tokens {
	tks := make(Tokens, 0, len(reducing))
	for _, v := range reducing {
		if tok, ok := asToken(v); ok {
			tks = append(tks, tok)
		}
	}
	return tks
}

func NewParser(tks Tokens, operatorPrecedence OperatorPrecedence) Parser {
	return Parser{
		originalTokens:             tks,
		reducingTokens:             tokensToAny(tks),
		operatorPrecedence:         operatorPrecedence,
		originalOperatorPrecedence: operatorPrecedence,
	}
}

func (p *Parser) Parse() (Node, error) {
	if len(p.originalTokens) > 0 && p.originalTokens[len(p.originalTokens)-1].Kind == ILLEGAL {
		return nil, ErrInvalidExpression
	}

	resolved, err := p.parse(p.reducingTokens)

	if err != nil || len(resolved) != 1 {
		return nil, ErrCouldNotParseExpression
	}

	raw := resolved[0]

	node, ok := raw.(Node)

	if !ok {
		return nil, ErrCouldNotParseExpression
	}

	return node, nil
}

func (p *Parser) parse(tokensAndNodes []any) ([]any, error) {
	if len(p.reducingTokens) == 1 {
		tokenOrNode := tokensAndNodes[0]

		if tk, ok := asToken(tokenOrNode); ok {
			node := NewOperand(tk.Kind, tk.Value)
			p.reducingTokens[0] = node
		}
		return p.reducingTokens, nil

	}

	resolvedGroupedTokensAndNodes, err := p.parseGroup()

	if err != nil {
		return p.reducingTokens, err
	}

	p.reducingTokens = resolvedGroupedTokensAndNodes

	p.parseExp()

	// fmt.Println(resolvedGroupedTokensAndNodes...)

	// return p.parse(p.reducingTokens)
	return resolvedGroupedTokensAndNodes, err
}

func (p *Parser) parseGroup() ([]any, error) {
	ind := array.FindIndex(p.reducingTokens, func(element any, index int, slice []any) bool {
		token, ok := asToken(element)
		return ok && token.Kind == LBRACKET
	})

	if ind == -1 {
		return p.reducingTokens, nil
	}

	leftCount := 1
	rightCount := 0
	start := ind + 1
	pos := start
	end := -1

	for pos < len(p.reducingTokens) {
		item := p.reducingTokens[pos]
		pos++

		token, ok := asToken(item)
		if !ok {
			continue
		}

		switch token.Kind {
		case LBRACKET:
			leftCount++
		case RBRACKET:
			rightCount++
		}

		if leftCount == rightCount {
			end = pos - 1
			break
		}
	}

	if end == -1 {
		return p.reducingTokens, fmt.Errorf("%w: left bracket at index %d", ErrNoMatchingBracket, ind)
	}

	group := p.reducingTokens[start:end]
	parser := NewParser(anyToTokens(group), p.originalOperatorPrecedence)
	parsedGroup, err := parser.Parse()

	if err != nil {
		return p.reducingTokens, nil
	}

	before := p.reducingTokens[:ind]
	// fmt.Println("before", before)
	after := p.reducingTokens[end+1:]
	// fmt.Println("after", after)

	result := make([]any, 0, len(before)+1+len(after))
	result = append(result, before...)
	result = append(result, parsedGroup)
	result = append(result, after...)
	// fmt.Println("result", result[0])

	return result, nil
}

func (p *Parser) parseExp() ([]any, error) {
	curOp := p.operatorPrecedence[0]

	fmt.Println("curOp", curOp)

	// remove if not
	if curOp != EXP {

	}

	return p.reducingTokens, nil
}

func (p *Parser) parseMul() ([]any, error) {
	return p.reducingTokens, nil
}

func (p *Parser) parseDiv() ([]any, error) {
	return p.reducingTokens, nil
}

func (p *Parser) parseAdd() ([]any, error) {
	return p.reducingTokens, nil
}

func (p *Parser) parseSub() ([]any, error) {
	return p.reducingTokens, nil
}

func asToken(v any) (Token, bool) {
	t, ok := v.(Token)
	return t, ok
}

func asNode(v any) (Node, bool) {
	n, ok := v.(Node)
	return n, ok
}
