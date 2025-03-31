package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
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

func (n *UnaryExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UnaryExpressionNode)
	if !ok {
		return false
	}

	return n.Op.Equal(o.Op) &&
		n.Right.Equal(value.Ref(o.Right)) &&
		n.span.Equal(o.span)
}

func (n *UnaryExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.Op.FetchValue())

	parens := ExpressionPrecedence(n) > ExpressionPrecedence(n.Right)
	if parens {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Right.String())
	if parens {
		buff.WriteRune(')')
	}

	return buff.String()
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
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::UnaryExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  op: ")
	indent.IndentStringFromSecondLine(&buff, n.Op.Inspect(), 1)

	buff.WriteString(",\n  right: ")
	indent.IndentStringFromSecondLine(&buff, n.Right.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *UnaryExpressionNode) Error() string {
	return n.Inspect()
}
