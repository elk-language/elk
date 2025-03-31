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

// Check if this node equals another node.
func (n *DoubleSplatExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*DoubleSplatExpressionNode)
	if !ok {
		return false
	}

	return n.Value.Equal(value.Ref(o.Value)) &&
		n.span.Equal(o.span)
}

// Return a string representation of the node.
func (n *DoubleSplatExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("**")

	parens := ExpressionPrecedence(n) > ExpressionPrecedence(n.Value)
	if parens {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Value.String())
	if parens {
		buff.WriteRune(')')
	}

	return buff.String()
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

func (n *SplatExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*SplatExpressionNode)
	if !ok {
		return false
	}

	return n.Value.Equal(value.Ref(o.Value)) &&
		n.span.Equal(o.span)
}

func (n *SplatExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteRune('*')
	parens := ExpressionPrecedence(n) > ExpressionPrecedence(n.Value)

	if parens {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Value.String())
	if parens {
		buff.WriteRune(')')
	}

	return buff.String()
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
