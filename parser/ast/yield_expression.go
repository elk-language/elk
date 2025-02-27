package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a `yield` expression eg. `yield`, `yield true`, `yield* foo()`
type YieldExpressionNode struct {
	NodeBase
	Value   ExpressionNode
	Forward bool
}

func (*YieldExpressionNode) IsStatic() bool {
	return false
}

// Create a new `yield` expression node eg. `yield`, `yield true`, `yield* foo()`
func NewYieldExpressionNode(span *position.Span, forward bool, val ExpressionNode) *YieldExpressionNode {
	return &YieldExpressionNode{
		NodeBase: NodeBase{span: span},
		Forward:  forward,
		Value:    val,
	}
}

func (*YieldExpressionNode) Class() *value.Class {
	return value.YieldExpressionNodeClass
}

func (*YieldExpressionNode) DirectClass() *value.Class {
	return value.YieldExpressionNodeClass
}

func (n *YieldExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::YieldExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  value: ")
	indentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *YieldExpressionNode) Error() string {
	return n.Inspect()
}
