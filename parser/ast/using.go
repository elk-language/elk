package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a using all entry node eg. `Foo::*`, `A::B::C::*`
type UsingAllEntryNode struct {
	TypedNodeBase
	Namespace ExpressionNode
}

func (*UsingAllEntryNode) IsStatic() bool {
	return false
}

func (*UsingAllEntryNode) Class() *value.Class {
	return value.UsingAllEntryNodeClass
}

func (*UsingAllEntryNode) DirectClass() *value.Class {
	return value.UsingAllEntryNodeClass
}

func (n *UsingAllEntryNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::UsingAllEntryNode{\n  &: %p", n)

	buff.WriteString(",\n  namespace: ")
	indentStringFromSecondLine(&buff, n.Namespace.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *UsingAllEntryNode) Error() string {
	return n.Inspect()
}

// Create a new using all entry node eg. `Foo::*`, `A::B::C::*`
func NewUsingAllEntryNode(span *position.Span, namespace UsingEntryNode) *UsingAllEntryNode {
	return &UsingAllEntryNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Namespace:     namespace,
	}
}

// Represents a using entry node with subentries eg. `Foo::{Bar, baz}`, `A::B::C::{lol, foo as epic, Gro as Moe}`
type UsingEntryWithSubentriesNode struct {
	NodeBase
	Namespace  ExpressionNode
	Subentries []UsingSubentryNode
}

func (*UsingEntryWithSubentriesNode) IsStatic() bool {
	return false
}

// Create a new using all entry node eg. `Foo::*`, `A::B::C::*`
func NewUsingEntryWithSubentriesNode(span *position.Span, namespace UsingEntryNode, subentries []UsingSubentryNode) *UsingEntryWithSubentriesNode {
	return &UsingEntryWithSubentriesNode{
		NodeBase:   NodeBase{span: span},
		Namespace:  namespace,
		Subentries: subentries,
	}
}

func (*UsingEntryWithSubentriesNode) Class() *value.Class {
	return value.UsingEntryWithSubentriesNodeClass
}

func (*UsingEntryWithSubentriesNode) DirectClass() *value.Class {
	return value.UsingEntryWithSubentriesNodeClass
}

func (n *UsingEntryWithSubentriesNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::UsingEntryWithSubentriesNode{\n  &: %p", n)

	buff.WriteString(",\n  namespace: ")
	indentStringFromSecondLine(&buff, n.Namespace.Inspect(), 1)

	buff.WriteString(",\n  subentries: %%[\n")
	for i, element := range n.Subentries {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *UsingEntryWithSubentriesNode) Error() string {
	return n.Inspect()
}
