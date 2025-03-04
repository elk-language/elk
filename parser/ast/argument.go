package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Nodes that implement this interface represent
// named arguments in method calls.
type NamedArgumentNode interface {
	Node
	namedArgumentNode()
}

func (*InvalidNode) namedArgumentNode()               {}
func (*NamedCallArgumentNode) namedArgumentNode()     {}
func (*DoubleSplatExpressionNode) namedArgumentNode() {}

// Represents a named argument in a function call eg. `foo: 123`
type NamedCallArgumentNode struct {
	NodeBase
	Name  string
	Value ExpressionNode
}

func (*NamedCallArgumentNode) IsStatic() bool {
	return false
}

func (*NamedCallArgumentNode) Class() *value.Class {
	return value.NamedCallArgumentNodeClass
}

func (*NamedCallArgumentNode) DirectClass() *value.Class {
	return value.NamedCallArgumentNodeClass
}

func (n *NamedCallArgumentNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::NamedCallArgumentNode{\n  &: %p", n)

	buff.WriteString(",\n  name: ")
	indentStringFromSecondLine(&buff, value.String(n.Name).Inspect(), 1)

	buff.WriteString(",\n  value: ")
	indentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *NamedCallArgumentNode) Error() string {
	return n.Inspect()
}

// Create a named argument node eg. `foo: 123`
func NewNamedCallArgumentNode(span *position.Span, name string, val ExpressionNode) *NamedCallArgumentNode {
	return &NamedCallArgumentNode{
		NodeBase: NodeBase{span: span},
		Name:     name,
		Value:    val,
	}
}
