package ast

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
)

// Expression of an operator with one operand eg. `!foo`, `-bar`
type UnaryExpressionNode struct {
	TypedNodeBase
	Op    *token.Token   // operator
	Right ExpressionNode // right hand side
}

func (u *UnaryExpressionNode) IsStatic() bool {
	return u.Right.IsStatic()
}

// Create a new unary expression node.
func NewUnaryExpressionNode(span *position.Span, op *token.Token, right ExpressionNode) *UnaryExpressionNode {
	return &UnaryExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		Right:         right,
	}
}

func (*UnaryExpressionNode) Class() *value.Class {
	return value.UnaryExpressionNodeClass
}

func (*UnaryExpressionNode) DirectClass() *value.Class {
	return value.UnaryExpressionNodeClass
}

func (n *UnaryExpressionNode) Inspect() string {
	return fmt.Sprintf(
		"Std::AST::UnaryExpressionNode{&: %p, op: %s, right: %s}",
		n,
		n.Op.Inspect(),
		n.Right.Inspect(),
	)
}

func (n *UnaryExpressionNode) Error() string {
	return n.Inspect()
}
