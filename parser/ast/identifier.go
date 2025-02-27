package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a public identifier eg. `foo`.
type PublicIdentifierNode struct {
	TypedNodeBase
	Value string
}

func (*PublicIdentifierNode) IsStatic() bool {
	return false
}

func (*PublicIdentifierNode) Class() *value.Class {
	return value.PublicIdentifierNodeClass
}

func (*PublicIdentifierNode) DirectClass() *value.Class {
	return value.PublicIdentifierNodeClass
}

func (n *PublicIdentifierNode) Inspect() string {
	return fmt.Sprintf("Std::AST::PublicIdentifierNode{&: %p, value: %s}", n, n.Value)
}

func (n *PublicIdentifierNode) Error() string {
	return n.Inspect()
}

// Create a new public identifier node eg. `foo`.
func NewPublicIdentifierNode(span *position.Span, val string) *PublicIdentifierNode {
	return &PublicIdentifierNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Represents a private identifier eg. `_foo`
type PrivateIdentifierNode struct {
	TypedNodeBase
	Value string
}

func (*PrivateIdentifierNode) IsStatic() bool {
	return false
}

func (*PrivateIdentifierNode) Class() *value.Class {
	return value.PrivateIdentifierNodeClass
}

func (*PrivateIdentifierNode) DirectClass() *value.Class {
	return value.PrivateIdentifierNodeClass
}

func (n *PrivateIdentifierNode) Inspect() string {
	return fmt.Sprintf("Std::AST::PrivateIdentifierNode{&: %p, value: %s}", n, n.Value)
}

func (n *PrivateIdentifierNode) Error() string {
	return n.Inspect()
}

// Create a new private identifier node eg. `_foo`.
func NewPrivateIdentifierNode(span *position.Span, val string) *PrivateIdentifierNode {
	return &PrivateIdentifierNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Represents an identifier with as in using declarations
// eg. `foo as bar`.
type PublicIdentifierAsNode struct {
	NodeBase
	Target *PublicIdentifierNode
	AsName string
}

func (*PublicIdentifierAsNode) IsStatic() bool {
	return false
}

func (*PublicIdentifierAsNode) Class() *value.Class {
	return value.ConstantAsNodeClass
}

func (*PublicIdentifierAsNode) DirectClass() *value.Class {
	return value.ConstantAsNodeClass
}

func (n *PublicIdentifierAsNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::PublicIdentifierAsNode{\n  &: %p", n)

	buff.WriteString(",\n  target: ")
	indentStringFromSecondLine(&buff, n.Target.Inspect(), 1)

	buff.WriteString(",\n  as_name: ")
	buff.WriteString(n.AsName)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *PublicIdentifierAsNode) Error() string {
	return n.Inspect()
}

// Create a new identifier with as eg. `foo as bar`.
func NewPublicIdentifierAsNode(span *position.Span, target *PublicIdentifierNode, as string) *PublicIdentifierAsNode {
	return &PublicIdentifierAsNode{
		NodeBase: NodeBase{span: span},
		Target:   target,
		AsName:   as,
	}
}
