package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a public constant eg. `Foo`.
type PublicConstantNode struct {
	TypedNodeBase
	Value string
}

func (*PublicConstantNode) IsStatic() bool {
	return false
}

func (*PublicConstantNode) Class() *value.Class {
	return value.PublicIdentifierNodeClass
}

func (*PublicConstantNode) DirectClass() *value.Class {
	return value.PublicIdentifierNodeClass
}

func (n *PublicConstantNode) Inspect() string {
	return fmt.Sprintf("Std::AST::PublicConstantNode{&: %p, value: %s}", n, n.Value)
}

func (n *PublicConstantNode) Error() string {
	return n.Inspect()
}

// Create a new public constant node eg. `Foo`.
func NewPublicConstantNode(span *position.Span, val string) *PublicConstantNode {
	return &PublicConstantNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Represents a private constant eg. `_Foo`
type PrivateConstantNode struct {
	TypedNodeBase
	Value string
}

func (*PrivateConstantNode) IsStatic() bool {
	return false
}

func (*PrivateConstantNode) Class() *value.Class {
	return value.PublicIdentifierNodeClass
}

func (*PrivateConstantNode) DirectClass() *value.Class {
	return value.PublicIdentifierNodeClass
}

func (n *PrivateConstantNode) Inspect() string {
	return fmt.Sprintf("Std::AST::PrivateConstantNode{&: %p, value: %s}", n, n.Value)
}

func (n *PrivateConstantNode) Error() string {
	return n.Inspect()
}

// Create a new private constant node eg. `_Foo`.
func NewPrivateConstantNode(span *position.Span, val string) *PrivateConstantNode {
	return &PrivateConstantNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Represents a constant with as in using declarations
// eg. `Foo as Bar`.
type PublicConstantAsNode struct {
	NodeBase
	Target *PublicConstantNode
	AsName string
}

func (*PublicConstantAsNode) IsStatic() bool {
	return false
}

func (*PublicConstantAsNode) Class() *value.Class {
	return value.PublicIdentifierNodeClass
}

func (*PublicConstantAsNode) DirectClass() *value.Class {
	return value.PublicIdentifierNodeClass
}

func (n *PublicConstantAsNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::PublicConstantAsNode{\n  &: %p", n)

	buff.WriteString(",\n  target: ")
	indentStringFromSecondLine(&buff, n.Target.Inspect(), 1)

	buff.WriteString(",\n  as_name: ")
	buff.WriteString(n.AsName)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *PublicConstantAsNode) Error() string {
	return n.Inspect()
}

// Create a new identifier with as eg. `Foo as Bar`.
func NewPublicConstantAsNode(span *position.Span, target *PublicConstantNode, as string) *PublicConstantAsNode {
	return &PublicConstantAsNode{
		NodeBase: NodeBase{span: span},
		Target:   target,
		AsName:   as,
	}
}

// Represents a constant lookup expressions eg. `Foo::Bar`
type ConstantLookupNode struct {
	TypedNodeBase
	Left  ExpressionNode      // left hand side
	Right ComplexConstantNode // right hand side
}

func (*ConstantLookupNode) IsStatic() bool {
	return false
}

func (*ConstantLookupNode) Class() *value.Class {
	return value.PublicIdentifierNodeClass
}

func (*ConstantLookupNode) DirectClass() *value.Class {
	return value.PublicIdentifierNodeClass
}

func (n *ConstantLookupNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::ConstantLookupNode{\n  &: %p", n)

	buff.WriteString(",\n  left: ")
	indentStringFromSecondLine(&buff, n.Left.Inspect(), 1)

	buff.WriteString(",\n  right: ")
	indentStringFromSecondLine(&buff, n.Right.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ConstantLookupNode) Error() string {
	return n.Inspect()
}

// Create a new constant lookup expression node eg. `Foo::Bar`
func NewConstantLookupNode(span *position.Span, left ExpressionNode, right ComplexConstantNode) *ConstantLookupNode {
	return &ConstantLookupNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Left:          left,
		Right:         right,
	}
}
