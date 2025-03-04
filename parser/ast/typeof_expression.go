package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a `typeof` expression eg. `typeof foo()`
type TypeofExpressionNode struct {
	TypedNodeBase
	Value ExpressionNode
}

func (*TypeofExpressionNode) IsStatic() bool {
	return false
}

// Create a new `typeof` expression node eg. `typeof foo()`
func NewTypeofExpressionNode(span *position.Span, val ExpressionNode) *TypeofExpressionNode {
	return &TypeofExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

func (*TypeofExpressionNode) Class() *value.Class {
	return value.TypeofExpressionNodeClass
}

func (*TypeofExpressionNode) DirectClass() *value.Class {
	return value.TypeofExpressionNodeClass
}

func (n *TypeofExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::TypeofExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  value: ")
	indentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *TypeofExpressionNode) Error() string {
	return n.Inspect()
}
