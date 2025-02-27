package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a `throw` expression eg. `throw ArgumentError("foo")`
type ThrowExpressionNode struct {
	NodeBase
	Unchecked bool
	Value     ExpressionNode
}

func (*ThrowExpressionNode) Type(*types.GlobalEnvironment) types.Type {
	return types.Never{}
}

func (*ThrowExpressionNode) IsStatic() bool {
	return false
}

// Create a new `throw` expression node eg. `throw ArgumentError("foo")`
func NewThrowExpressionNode(span *position.Span, unchecked bool, val ExpressionNode) *ThrowExpressionNode {
	return &ThrowExpressionNode{
		NodeBase:  NodeBase{span: span},
		Unchecked: unchecked,
		Value:     val,
	}
}

func (*ThrowExpressionNode) Class() *value.Class {
	return value.ThrowExpressionNodeClass
}

func (*ThrowExpressionNode) DirectClass() *value.Class {
	return value.ThrowExpressionNodeClass
}

func (n *ThrowExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::ThrowExpressionNode{\n  &: %p", n)

	fmt.Fprintf(&buff, ",\n  unchecked: %t", n.Unchecked)

	buff.WriteString(",\n  value: ")
	indentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ThrowExpressionNode) Error() string {
	return n.Inspect()
}
