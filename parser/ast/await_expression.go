package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents an `await` expression eg. `await foo()`
type AwaitExpressionNode struct {
	TypedNodeBase
	Value ExpressionNode
}

func (*AwaitExpressionNode) IsStatic() bool {
	return false
}

// Create a new `await` expression node eg. `await foo()`
func NewAwaitExpressionNode(span *position.Span, val ExpressionNode) *AwaitExpressionNode {
	return &AwaitExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

func (*AwaitExpressionNode) Class() *value.Class {
	return value.AwaitExpressionNodeClass
}

func (*AwaitExpressionNode) DirectClass() *value.Class {
	return value.AwaitExpressionNodeClass
}

func (n *AwaitExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::AwaitExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  value: ")
	indentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *AwaitExpressionNode) Error() string {
	return n.Inspect()
}
