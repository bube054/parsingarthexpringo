package ast

import (
	"fmt"

	"github.com/bube054/parsingarthexpringo/lexer"
)

type NodeKind int

const (
	NODE_NUMBER_ALPHABET NodeKind = iota
	NODE_BINARY_EXPR
)

type Node interface {
	NodeKind() NodeKind
	Equals(Node) bool
	fmt.Stringer
}

type Operand struct {
	Token lexer.Token
}

func NewOperand(kind lexer.TokenKind, value string) *Operand {
	return &Operand{
		Token: lexer.NewToken(kind, value),
	}
}

func (n *Operand) NodeKind() NodeKind {
	return NODE_NUMBER_ALPHABET
}

func (n *Operand) Equals(other Node) bool {
	o, ok := other.(*Operand)
	if !ok {
		return false
	}
	return n.Token.Equal(&o.Token)
}

func (o *Operand) String() string {
	return fmt.Sprintf("Operand{%s}", o.Token)
}

type BinaryExpr struct {
	Operator lexer.TokenKind
	Left     Node
	Right    Node
}

func NewBinaryExpr(operator lexer.TokenKind, left Node, right Node) *BinaryExpr {
	return &BinaryExpr{
		Operator: operator,
		Left:     left,
		Right:    right,
	}
}

func (b *BinaryExpr) NodeKind() NodeKind {
	return NODE_BINARY_EXPR
}

func (b *BinaryExpr) Equals(other Node) bool {
	o, ok := other.(*BinaryExpr)
	if !ok {
		return false
	}
	if b.Operator != o.Operator {
		return false
	}
	if b.Left == nil && o.Left != nil || b.Left != nil && o.Left == nil {
		return false
	}
	if b.Right == nil && o.Right != nil || b.Right != nil && o.Right == nil {
		return false
	}
	// recursive equality check
	leftEq := (b.Left == nil && o.Left == nil) || (b.Left != nil && b.Left.Equals(o.Left))
	rightEq := (b.Right == nil && o.Right == nil) || (b.Right != nil && b.Right.Equals(o.Right))
	return leftEq && rightEq
}

func (b *BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Left, b.Operator, b.Right)
}

func AsOperand(v any) (Operand, bool) {
	n, ok := v.(Operand)
	return n, ok
}

func AsOperator(v any) (Operand, bool) {
	n, ok := v.(Operand)
	return n, ok
}

func AsNode(v any) (Node, bool) {
	n, ok := v.(Node)
	return n, ok
}

func AsNodeV2(v any) (Node, bool) {
	switch x := v.(type) {
	case Node:
		return x, true
	case lexer.Token:
		return NewOperand(x.Kind, x.Value), true
	default:
		return nil, false
	}
}
