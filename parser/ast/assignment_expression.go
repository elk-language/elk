package ast

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
)

// Assignment with the specified operator.
type AssignmentExpressionNode struct {
	TypedNodeBase
	Op    *token.Token   // operator
	Left  ExpressionNode // left hand side
	Right ExpressionNode // right hand side
}

func (*AssignmentExpressionNode) IsStatic() bool {
	return false
}

// Create a new assignment expression node eg. `foo = 3`
func NewAssignmentExpressionNode(span *position.Span, op *token.Token, left, right ExpressionNode) *AssignmentExpressionNode {
	return &AssignmentExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		Left:          left,
		Right:         right,
	}
}

func (*AssignmentExpressionNode) Class() *value.Class {
	return value.AssignmentExpressionNodeClass
}

func (*AssignmentExpressionNode) DirectClass() *value.Class {
	return value.AssignmentExpressionNodeClass
}

func (n *AssignmentExpressionNode) Inspect() string {
	return fmt.Sprintf(
		"Std::AST::AssignmentExpressionNode{&: %p, op: %s, left: %s, right: %s}",
		n,
		n.Op.Inspect(),
		n.Left.Inspect(),
		n.Right.Inspect(),
	)
}

func (p *AssignmentExpressionNode) Error() string {
	return p.Inspect()
}
