package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a `must` expression eg. `must foo()`
type MustExpressionNode struct {
	TypedNodeBase
	Value ExpressionNode
}

func (*MustExpressionNode) IsStatic() bool {
	return false
}

// Create a new `must` expression node eg. `must foo()`
func NewMustExpressionNode(span *position.Span, val ExpressionNode) *MustExpressionNode {
	return &MustExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

func (*MustExpressionNode) Class() *value.Class {
	return value.MustExpressionNodeClass
}

func (*MustExpressionNode) DirectClass() *value.Class {
	return value.MustExpressionNodeClass
}

func (n *MustExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::MustExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  value: ")
	indent.IndentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *MustExpressionNode) Error() string {
	return n.Inspect()
}
