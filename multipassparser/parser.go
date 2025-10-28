package multipassparser

import (
	"errors"
	"fmt"

	"github.com/bube054/go-js-array-methods/array"
	"github.com/bube054/parsingarthexpringo/ast"
	"github.com/bube054/parsingarthexpringo/lexer"
)

var (
	ErrInvalidExpression       = errors.New("invalid expression")
	ErrCouldNotParseExpression = errors.New("could not parse expression")
	ErrInvalidReducingTokens   = errors.New("invalid reducing tokens type")
	ErrNoMatchingBracket       = errors.New("no matching right bracket")
	ErrInvalidInfixExpression  = errors.New("invalid infix expression")
)

// type OperatorPrecedence map[lexer.TokenKind]int
type OperatorPrecedence []lexer.TokenKind

var BIDMAS = OperatorPrecedence{lexer.POW, lexer.MUL, lexer.DIV, lexer.ADD, lexer.SUB}

type Parser struct {
	originalTokens             lexer.Tokens
	reducingTokens             []any
	originalOperatorPrecedence OperatorPrecedence
	operatorPrecedence         OperatorPrecedence
	hasResolvedGroups          bool
}

func tokensToAny(tks lexer.Tokens) []any {
	reducing := make([]any, len(tks))
	for i, v := range tks {
		reducing[i] = v
	}
	return reducing
}

func anyToTokens(reducing []any) lexer.Tokens {
	tks := make(lexer.Tokens, 0, len(reducing))
	for _, v := range reducing {
		if tok, ok := lexer.AsToken(v); ok {
			tks = append(tks, tok)
		}
	}
	return tks
}

func NewParser(tks lexer.Tokens, operatorPrecedence OperatorPrecedence) Parser {
	return Parser{
		originalTokens:             tks,
		reducingTokens:             tokensToAny(tks),
		operatorPrecedence:         operatorPrecedence,
		originalOperatorPrecedence: operatorPrecedence,
	}
}

func (p *Parser) Parse() (ast.Node, error) {
	if len(p.originalTokens) > 0 && p.originalTokens[len(p.originalTokens)-1].Kind == lexer.ILLEGAL {
		return nil, ErrInvalidExpression
	}

	resolved, err := p.parse(p.reducingTokens)

	if err != nil || len(resolved) != 1 {
		return nil, ErrCouldNotParseExpression
	}

	// fmt.Println("Resolved", resolved)

	raw := resolved[0]

	node, ok := raw.(ast.Node)

	if !ok {
		return nil, ErrCouldNotParseExpression
	}

	return node, nil
}

func (p *Parser) parse(tokensAndNodes []any) ([]any, error) {
	if len(p.reducingTokens) == 1 {
		tokenOrNode := tokensAndNodes[0]

		if tk, ok := lexer.AsToken(tokenOrNode); ok {
			node := ast.NewOperand(tk.Kind, tk.Value)
			p.reducingTokens[0] = node
		}

		return p.reducingTokens, nil
	}

	reducedTokens, err := p.parseGroup()
	if err != nil {
		return reducedTokens, nil
	}
	p.reducingTokens = reducedTokens

	reducedTokens, err = p.parsePow()
	if err != nil {
		return reducedTokens, nil
	}
	p.reducingTokens = reducedTokens

	reducedTokens, err = p.parseMul()
	if err != nil {
		return reducedTokens, nil
	}
	p.reducingTokens = reducedTokens

	reducedTokens, err = p.parseDiv()
	if err != nil {
		return reducedTokens, nil
	}
	p.reducingTokens = reducedTokens

	reducedTokens, err = p.parseAdd()
	if err != nil {
		return reducedTokens, nil
	}
	p.reducingTokens = reducedTokens

	reducedTokens, err = p.parseSub()
	if err != nil {
		return reducedTokens, nil
	}
	p.reducingTokens = reducedTokens

	// return p.parse(p.reducingTokens)
	return p.parse(p.reducingTokens)
}

func (p *Parser) parseGroup() ([]any, error) {
	ind := array.FindIndex(p.reducingTokens, func(element any, index int, slice []any) bool {
		token, ok := lexer.AsToken(element)
		return ok && token.Kind == lexer.LBRACKET
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

		token, ok := lexer.AsToken(item)
		if !ok {
			continue
		}

		switch token.Kind {
		case lexer.LBRACKET:
			leftCount++
		case lexer.RBRACKET:
			rightCount++
		}

		if leftCount == rightCount {
			end = pos - 1
			break
		}
	}

	if end == -1 {
		return nil, fmt.Errorf("%w: left bracket at index %d", ErrNoMatchingBracket, ind)
	}

	group := p.reducingTokens[start:end]
	parser := NewParser(anyToTokens(group), p.originalOperatorPrecedence)
	parsedGroup, err := parser.Parse()

	if err != nil {
		return p.reducingTokens, nil
	}

	before := p.reducingTokens[:ind]
	beforeLastItemRaw, _ := array.At(before, -1)
	beforeLastItem, ok := lexer.AsToken(beforeLastItemRaw)

	if len(before) > 0 && (!ok || (beforeLastItem.Category != lexer.CATEGORY_OPERATOR)) {
		// If the last item before this group is not a valid operator
		// (either it's not a Token, or it's a Token but not an operator, or array.At errored),
		// then we implicitly insert a multiplication (*) to keep the expression valid.
		before = append(before, lexer.NewToken(lexer.MUL, "*"))
	}

	// fmt.Println("before", before)

	after := p.reducingTokens[end+1:]

	afterLastItemRaw, _ := array.At(after, 0)
	afterLastItem, ok := lexer.AsToken(afterLastItemRaw)

	if len(after) > 0 && (!ok || (afterLastItem.Category != lexer.CATEGORY_OPERATOR)) {
		// If the first item after this group is not a valid operator
		// (either it's not a Token, or it's a Token but not an operator, or array.At errored),
		// then we implicitly insert a multiplication (*) before it.
		after = append([]any{lexer.NewToken(lexer.MUL, "*")}, after...)
	}

	// fmt.Println("after", after)

	result := make([]any, 0, len(before)+1+len(after))
	result = append(result, before...)
	result = append(result, parsedGroup)
	result = append(result, after...)
	// fmt.Println("result", result)

	return result, nil
}

func (p *Parser) parseInfixExpr(operatorKind lexer.TokenKind) ([]any, error) {
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
		token, ok := lexer.AsToken(element)
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
		return nil, ErrInvalidInfixExpression
	}

	leftOperandRaw := p.reducingTokens[start]
	leftOperand, ok := ast.AsNodeV2(leftOperandRaw)

	if !ok {
		return nil, ErrInvalidInfixExpression
	}

	operatorRaw := p.reducingTokens[ind]
	operator, ok := lexer.AsToken(operatorRaw)

	if !ok {
		return nil, ErrInvalidInfixExpression
	}

	rightOperandRaw := p.reducingTokens[end]
	rightOperand, ok := ast.AsNodeV2(rightOperandRaw)

	if !ok {
		return nil, ErrInvalidInfixExpression
	}

	resolvedBinaryExpression := ast.NewBinaryExpr(operator.Kind, leftOperand, rightOperand)
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
	return p.parseInfixExpr(lexer.POW)
}

func (p *Parser) parseMul() ([]any, error) {
	return p.parseInfixExpr(lexer.MUL)
}

func (p *Parser) parseDiv() ([]any, error) {
	return p.parseInfixExpr(lexer.DIV)
}

func (p *Parser) parseAdd() ([]any, error) {
	return p.parseInfixExpr(lexer.ADD)
}

func (p *Parser) parseSub() ([]any, error) {
	return p.parseInfixExpr(lexer.SUB)
}
