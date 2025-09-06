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
	ErrInvalidInfixExpression  = errors.New("invalid infix expression")
)

// type OperatorPrecedence map[TokenKind]int
type OperatorPrecedence []TokenKind

var BIDMAS = OperatorPrecedence{POW, MUL, DIV, SUB, ADD}

type Parser struct {
	originalTokens             Tokens
	reducingTokens             []any
	originalOperatorPrecedence OperatorPrecedence
	operatorPrecedence         OperatorPrecedence
	hasResolvedGroups          bool
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

	reducedTokens, err := p.parseGroup()
	_ = err
	p.reducingTokens = reducedTokens
	
	reducedTokens, err = p.parsePow()
	p.reducingTokens = reducedTokens
	
	reducedTokens, err = p.parseMul()
	p.reducingTokens = reducedTokens
	
	reducedTokens, err = p.parseDiv()
	p.reducingTokens = reducedTokens
	
	reducedTokens, err = p.parseAdd()
	p.reducingTokens = reducedTokens
	
	reducedTokens, err = p.parseSub()
	p.reducingTokens = reducedTokens

	// return p.parse(p.reducingTokens)
	return p.parse(p.reducingTokens)
}

func (p *Parser) parseGroup() ([]any, error) {
	ind := array.FindIndex(p.reducingTokens, func(element any, index int, slice []any) bool {
		token, ok := asToken(element)
		return ok && token.Kind == LBRACKET
	})

	if ind == -1 {
		p.hasResolvedGroups = true
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

func (p *Parser) parseInfixExpr(operatorKind TokenKind) ([]any, error) {
	if !(len(p.operatorPrecedence) > 0) {
		return p.reducingTokens, nil
	}

	curOp := p.operatorPrecedence[0]

	if !p.hasResolvedGroups || curOp != operatorKind {
		return p.reducingTokens, nil
	}

	// fmt.Println("reducingTokens", p.reducingTokens)
	// fmt.Println("curOp", curOp)
	// fmt.Println("hasResolvedGroups", p.hasResolvedGroups)

	ind := array.FindIndex(p.reducingTokens, func(element any, index int, slice []any) bool {
		token, ok := asToken(element)
		return ok && token.Kind == operatorKind
	})

	if ind == -1 {
		p.operatorPrecedence = p.operatorPrecedence[1:]
		return p.reducingTokens, nil
	}

	start := ind - 1
	end := ind + 1

	// fmt.Println("start", start)
	// fmt.Println("end", end)

	if len(p.reducingTokens) < 3 {
		return p.reducingTokens, ErrInvalidInfixExpression
	}

	leftOperandRaw := p.reducingTokens[start]
	leftOperand, ok := asNodeV2(leftOperandRaw)

	if !ok {
		return p.reducingTokens, ErrInvalidInfixExpression
	}

	operatorRaw := p.reducingTokens[ind]
	operator, ok := asToken(operatorRaw)

	if !ok {
		return p.reducingTokens, ErrInvalidInfixExpression
	}

	rightOperandRaw := p.reducingTokens[end]
	rightOperand, ok := asNodeV2(rightOperandRaw)

	if !ok {
		return p.reducingTokens, ErrInvalidInfixExpression
	}

	resolvedBinaryExpression := NewBinaryExpr(operator.Kind, leftOperand, rightOperand)
	// fmt.Println("resolvedBinaryExpression", resolvedBinaryExpression)

	before := p.reducingTokens[:start]
	// fmt.Println("before", before)
	after := p.reducingTokens[end+1:]
	// fmt.Println("after", after)

	result := make([]any, 0, len(before)+1+len(after))
	result = append(result, before...)
	result = append(result, resolvedBinaryExpression)
	result = append(result, after...)
	// fmt.Println("result", result[0])

	return result, nil
}

func (p *Parser) parsePow() ([]any, error) {
	return p.parseInfixExpr(POW)
}

func (p *Parser) parseMul() ([]any, error) {
	return p.parseInfixExpr(MUL)
}

func (p *Parser) parseDiv() ([]any, error) {
	return p.parseInfixExpr(DIV)
}

func (p *Parser) parseAdd() ([]any, error) {
	return p.parseInfixExpr(ADD)
}

func (p *Parser) parseSub() ([]any, error) {
	return p.parseInfixExpr(SUB)
}
