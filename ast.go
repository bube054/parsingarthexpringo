package parsingarthexpringo

type NodeKind int

const (
	NODE_NUMBER_ALPHABET NodeKind = iota
	NODE_BINARY_EXPR
)

type Node interface {
	NodeKind() NodeKind
	Equals(Node) bool
}

type Operand struct {
	Token Token
}

func NewOperand(kind TokenKind, value string) *Operand {
	return &Operand{
		Token: NewToken(kind, value),
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

type BinaryExpr struct {
	Operator TokenKind
	Left     Node
	Right    Node
}

func NewBinaryExpr(operator TokenKind, left Node, right Node) *BinaryExpr {
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

func asOperand(v any) (Operand, bool) {
	n, ok := v.(Operand)
	return n, ok
}
