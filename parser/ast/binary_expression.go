package ast

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
)

// Expression of an operator with two operands eg. `2 + 5`, `foo > bar`
type BinaryExpressionNode struct {
	TypedNodeBase
	Op     *token.Token   // operator
	Left   ExpressionNode // left hand side
	Right  ExpressionNode // right hand side
	static bool
}

func (b *BinaryExpressionNode) IsStatic() bool {
	return b.static
}

// Create a new binary expression node.
func NewBinaryExpressionNode(span *position.Span, op *token.Token, left, right ExpressionNode) *BinaryExpressionNode {
	return &BinaryExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		Left:          left,
		Right:         right,
		static:        areExpressionsStatic(left, right),
	}
}

// Same as [NewBinaryExpressionNode] but returns an interface
func NewBinaryExpressionNodeI(span *position.Span, op *token.Token, left, right ExpressionNode) ExpressionNode {
	return NewBinaryExpressionNode(span, op, left, right)
}

func (*BinaryExpressionNode) Class() *value.Class {
	return value.BinaryExpressionNodeClass
}

func (*BinaryExpressionNode) DirectClass() *value.Class {
	return value.BinaryExpressionNodeClass
}

func (n *BinaryExpressionNode) Inspect() string {
	return fmt.Sprintf(
		"Std::AST::BinaryExpressionNode{&: %p, op: %s, left: %s, right: %s}",
		n,
		n.Op.Inspect(),
		n.Left.Inspect(),
		n.Right.Inspect(),
	)
}

func (n *BinaryExpressionNode) Error() string {
	return n.Inspect()
}
