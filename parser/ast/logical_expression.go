package ast

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
)

// Expression of a logical operator with two operands eg. `foo && bar`
type LogicalExpressionNode struct {
	TypedNodeBase
	Op     *token.Token   // operator
	Left   ExpressionNode // left hand side
	Right  ExpressionNode // right hand side
	static bool
}

func (l *LogicalExpressionNode) IsStatic() bool {
	return l.static
}

// Create a new logical expression node.
func NewLogicalExpressionNode(span *position.Span, op *token.Token, left, right ExpressionNode) *LogicalExpressionNode {
	return &LogicalExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		Left:          left,
		Right:         right,
		static:        areExpressionsStatic(left, right),
	}
}

// Same as [NewLogicalExpressionNode] but returns an interface
func NewLogicalExpressionNodeI(span *position.Span, op *token.Token, left, right ExpressionNode) ExpressionNode {
	return &LogicalExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		Left:          left,
		Right:         right,
	}
}

func (*LogicalExpressionNode) Class() *value.Class {
	return value.LogicalExpressionNodeClass
}

func (*LogicalExpressionNode) DirectClass() *value.Class {
	return value.LogicalExpressionNodeClass
}

func (n *LogicalExpressionNode) Inspect() string {
	return fmt.Sprintf(
		"Std::AST::LogicalExpressionNode{&: %p, op: %s, left: %s, right: %s}",
		n,
		n.Op.Inspect(),
		n.Left.Inspect(),
		n.Right.Inspect(),
	)
}

func (n *LogicalExpressionNode) Error() string {
	return n.Inspect()
}
