package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// All nodes that should be valid identifiers
// should implement this interface.
type IdentifierNode interface {
	Node
	PatternExpressionNode
	identifierNode()
}

func (*InvalidNode) identifierNode()           {}
func (*PublicIdentifierNode) identifierNode()  {}
func (*PrivateIdentifierNode) identifierNode() {}

// Represents a public identifier eg. `foo`.
type PublicIdentifierNode struct {
	TypedNodeBase
	Value string
}

func (n *PublicIdentifierNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*PublicIdentifierNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.span.Equal(o.span)
}

func (n *PublicIdentifierNode) String() string {
	return n.Value
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
	return fmt.Sprintf("Std::Elk::AST::PublicIdentifierNode{span: %s, value: %s}", (*value.Span)(n.span).Inspect(), n.Value)
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

func (n *PrivateIdentifierNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*PrivateIdentifierNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.span.Equal(o.span)
}

func (n *PrivateIdentifierNode) String() string {
	return n.Value
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
	return fmt.Sprintf("Std::Elk::AST::PrivateIdentifierNode{span: %s, value: %s}", (*value.Span)(n.span).Inspect(), n.Value)
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

func (n *PublicIdentifierAsNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*PublicIdentifierAsNode)
	if !ok {
		return false
	}

	return n.Target.Equal(value.Ref(o.Target)) &&
		n.AsName == o.AsName &&
		n.span.Equal(o.span)
}

func (n *PublicIdentifierAsNode) String() string {
	return fmt.Sprintf("%s as %s", n.Target.String(), n.AsName)
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

	fmt.Fprintf(&buff, "Std::Elk::AST::PublicIdentifierAsNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  target: ")
	indent.IndentStringFromSecondLine(&buff, n.Target.Inspect(), 1)

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
