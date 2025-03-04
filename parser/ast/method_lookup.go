package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a method lookup expression eg. `Foo::bar`, `a::c`
type MethodLookupNode struct {
	TypedNodeBase
	Receiver ExpressionNode
	Name     string
}

func (*MethodLookupNode) IsStatic() bool {
	return false
}

// Create a new method lookup expression node eg. `Foo::bar`, `a::c`
func NewMethodLookupNode(span *position.Span, receiver ExpressionNode, name string) *MethodLookupNode {
	return &MethodLookupNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Receiver:      receiver,
		Name:          name,
	}
}

func (*MethodLookupNode) Class() *value.Class {
	return value.MethodLookupNodeClass
}

func (*MethodLookupNode) DirectClass() *value.Class {
	return value.MethodLookupNodeClass
}

func (n *MethodLookupNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::MethodLookupNode{\n  &: %p", n)

	buff.WriteString(",\n  receiver: ")
	indentStringFromSecondLine(&buff, n.Receiver.Inspect(), 1)

	buff.WriteString(",\n  name: ")
	buff.WriteString(n.Name)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *MethodLookupNode) Error() string {
	return n.Inspect()
}

// Represents a method lookup with as in using declarations
// eg. `Foo::bar as Bar`.
type MethodLookupAsNode struct {
	NodeBase
	MethodLookup *MethodLookupNode
	AsName       string
}

func (*MethodLookupAsNode) IsStatic() bool {
	return false
}

// Create a new identifier with as eg. `Foo::bar as Bar`.
func NewMethodLookupAsNode(span *position.Span, methodLookup *MethodLookupNode, as string) *MethodLookupAsNode {
	return &MethodLookupAsNode{
		NodeBase:     NodeBase{span: span},
		MethodLookup: methodLookup,
		AsName:       as,
	}
}

func (*MethodLookupAsNode) Class() *value.Class {
	return value.ConstantAsNodeClass
}

func (*MethodLookupAsNode) DirectClass() *value.Class {
	return value.ConstantAsNodeClass
}

func (n *MethodLookupAsNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::MethodLookupAsNode{\n  &: %p", n)

	buff.WriteString(",\n  method_lookup: ")
	indentStringFromSecondLine(&buff, n.MethodLookup.Inspect(), 1)

	buff.WriteString(",\n  as_name: ")
	buff.WriteString(n.AsName)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *MethodLookupAsNode) Error() string {
	return n.Inspect()
}
