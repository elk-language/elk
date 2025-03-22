package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a double splat expression eg. `**foo`
type DoubleSplatExpressionNode struct {
	TypedNodeBase
	Value ExpressionNode
}

func (*DoubleSplatExpressionNode) IsStatic() bool {
	return false
}

func (*DoubleSplatExpressionNode) Class() *value.Class {
	return value.DoubleSplatExpressionNodeClass
}

func (*DoubleSplatExpressionNode) DirectClass() *value.Class {
	return value.DoubleSplatExpressionNodeClass
}

func (n *DoubleSplatExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::DoubleSplatExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  value: ")
	indent.IndentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *DoubleSplatExpressionNode) Error() string {
	return n.Inspect()
}

// Create a double splat expression node eg. `**foo`
func NewDoubleSplatExpressionNode(span *position.Span, val ExpressionNode) *DoubleSplatExpressionNode {
	return &DoubleSplatExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Represents a splat expression eg. `*foo`
type SplatExpressionNode struct {
	TypedNodeBase
	Value ExpressionNode
}

func (*SplatExpressionNode) IsStatic() bool {
	return false
}

func (*SplatExpressionNode) Class() *value.Class {
	return value.SplatExpressionNodeClass
}

func (*SplatExpressionNode) DirectClass() *value.Class {
	return value.SplatExpressionNodeClass
}

func (n *SplatExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::SplatExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  value: ")
	indent.IndentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *SplatExpressionNode) Error() string {
	return n.Inspect()
}

// Create a splat expression node eg. `*foo`
func NewSplatExpressionNode(span *position.Span, val ExpressionNode) *SplatExpressionNode {
	return &SplatExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}
