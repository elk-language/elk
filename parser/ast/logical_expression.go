package ast

import (
	"fmt"
	"strings"

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
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::LogicalExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  op: ")
	indentStringFromSecondLine(&buff, n.Op.Inspect(), 1)

	buff.WriteString(",\n  left: ")
	indentStringFromSecondLine(&buff, n.Left.Inspect(), 1)

	buff.WriteString(",\n  right: ")
	indentStringFromSecondLine(&buff, n.Right.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *LogicalExpressionNode) Error() string {
	return n.Inspect()
}
