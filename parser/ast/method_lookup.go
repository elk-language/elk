package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a method lookup expression eg. `Foo::bar`, `a::c`
type MethodLookupNode struct {
	TypedNodeBase
	Receiver ExpressionNode
	Name     string
}

func (n *MethodLookupNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*MethodLookupNode)
	if !ok {
		return false
	}

	return n.Span().Equal(o.Span()) &&
		n.Receiver.Equal(value.Ref(o.Receiver)) &&
		n.Name == o.Name
}

func (n *MethodLookupNode) String() string {
	var buff strings.Builder

	parens := ExpressionPrecedence(n) > ExpressionPrecedence(n.Receiver)

	if parens {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Receiver.String())
	if parens {
		buff.WriteRune(')')
	}

	buff.WriteString("::")
	buff.WriteString(n.Name)

	return buff.String()
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

	fmt.Fprintf(&buff, "Std::Elk::AST::MethodLookupNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  receiver: ")
	indent.IndentStringFromSecondLine(&buff, n.Receiver.Inspect(), 1)

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

// Check if this method lookup as node is equal to another value.
func (n *MethodLookupAsNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*MethodLookupAsNode)
	if !ok {
		return false
	}

	return n.span.Equal(o.span) &&
		n.MethodLookup.Equal(value.Ref(o.MethodLookup)) &&
		n.AsName == o.AsName
}

// Return a string representation of this method lookup as node.
func (n *MethodLookupAsNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.MethodLookup.String())
	buff.WriteString(" as ")
	buff.WriteString(n.AsName)

	return buff.String()
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

	fmt.Fprintf(&buff, "Std::Elk::AST::MethodLookupAsNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  method_lookup: ")
	indent.IndentStringFromSecondLine(&buff, n.MethodLookup.Inspect(), 1)

	buff.WriteString(",\n  as_name: ")
	buff.WriteString(n.AsName)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *MethodLookupAsNode) Error() string {
	return n.Inspect()
}
