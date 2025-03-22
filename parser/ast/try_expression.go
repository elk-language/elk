package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a `try` expression eg. `try foo()`
type TryExpressionNode struct {
	TypedNodeBase
	Value ExpressionNode
}

func (*TryExpressionNode) IsStatic() bool {
	return false
}

// Create a new `try` expression node eg. `try foo()`
func NewTryExpressionNode(span *position.Span, val ExpressionNode) *TryExpressionNode {
	return &TryExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

func (*TryExpressionNode) Class() *value.Class {
	return value.TryExpressionNodeClass
}

func (*TryExpressionNode) DirectClass() *value.Class {
	return value.TryExpressionNodeClass
}

func (n *TryExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::TryExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  value: ")
	indent.IndentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *TryExpressionNode) Error() string {
	return n.Inspect()
}
