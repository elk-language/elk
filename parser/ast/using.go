package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents all nodes that are valid in using declarations
type UsingEntryNode interface {
	Node
	ExpressionNode
	usingEntryNode()
}

func (*InvalidNode) usingEntryNode()                  {}
func (*PublicConstantNode) usingEntryNode()           {}
func (*PrivateConstantNode) usingEntryNode()          {}
func (*ConstantLookupNode) usingEntryNode()           {}
func (*MethodLookupNode) usingEntryNode()             {}
func (*UsingAllEntryNode) usingEntryNode()            {}
func (*UsingEntryWithSubentriesNode) usingEntryNode() {}
func (*ConstantAsNode) usingEntryNode()               {}
func (*MethodLookupAsNode) usingEntryNode()           {}
func (*GenericConstantNode) usingEntryNode()          {}
func (*NilLiteralNode) usingEntryNode()               {}

// Represents all nodes that are valid in using subentries
// in `UsingEntryWithSubentriesNode`
type UsingSubentryNode interface {
	Node
	ExpressionNode
	usingSubentryNode()
}

func (*InvalidNode) usingSubentryNode()            {}
func (*PublicConstantNode) usingSubentryNode()     {}
func (*PublicConstantAsNode) usingSubentryNode()   {}
func (*PublicIdentifierNode) usingSubentryNode()   {}
func (*PublicIdentifierAsNode) usingSubentryNode() {}

// Represents a using all entry node eg. `Foo::*`, `A::B::C::*`
type UsingAllEntryNode struct {
	TypedNodeBase
	Namespace UsingEntryNode
}

func (n *UsingAllEntryNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UsingAllEntryNode)
	if !ok {
		return false
	}

	return n.Namespace.Equal(value.Ref(o.Namespace)) &&
		n.span.Equal(o.span)
}

func (n *UsingAllEntryNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.Namespace.String())
	buff.WriteString("::*")

	return buff.String()
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

	fmt.Fprintf(&buff, "Std::Elk::AST::UsingAllEntryNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  namespace: ")
	indent.IndentStringFromSecondLine(&buff, n.Namespace.Inspect(), 1)

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
	Namespace  UsingEntryNode
	Subentries []UsingSubentryNode
}

func (n *UsingEntryWithSubentriesNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UsingEntryWithSubentriesNode)
	if !ok {
		return false
	}

	if len(n.Subentries) != len(o.Subentries) ||
		!n.Namespace.Equal(value.Ref(o.Namespace)) ||
		!n.span.Equal(o.span) {
		return false
	}

	for i, subentry := range n.Subentries {
		if !subentry.Equal(value.Ref(o.Subentries[i])) {
			return false
		}
	}

	return true
}

func (n *UsingEntryWithSubentriesNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.Namespace.String())
	buff.WriteString("::{")

	for i, subentry := range n.Subentries {
		if i != 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(subentry.String())
	}

	buff.WriteRune('}')

	return buff.String()
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

	fmt.Fprintf(&buff, "Std::Elk::AST::UsingEntryWithSubentriesNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  namespace: ")
	indent.IndentStringFromSecondLine(&buff, n.Namespace.Inspect(), 1)

	buff.WriteString(",\n  subentries: %[\n")
	for i, element := range n.Subentries {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *UsingEntryWithSubentriesNode) Error() string {
	return n.Inspect()
}

// Represents a using expression eg. `using Foo`
type UsingExpressionNode struct {
	TypedNodeBase
	Entries []UsingEntryNode
}

func (n *UsingExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UsingExpressionNode)
	if !ok {
		return false
	}

	if len(n.Entries) != len(o.Entries) ||
		!n.span.Equal(o.span) {
		return false
	}

	for i, entry := range n.Entries {
		if !entry.Equal(value.Ref(o.Entries[i])) {
			return false
		}
	}

	return true
}

func (n *UsingExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("using ")

	for i, entry := range n.Entries {
		if i != 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(entry.String())
	}

	return buff.String()
}

func (*UsingExpressionNode) SkipTypechecking() bool {
	return false
}

func (*UsingExpressionNode) IsStatic() bool {
	return false
}

func (*UsingExpressionNode) Class() *value.Class {
	return value.UsingExpressionNodeClass
}

func (*UsingExpressionNode) DirectClass() *value.Class {
	return value.UsingExpressionNodeClass
}

func (n *UsingExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::UsingExpressionNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  entries: %[\n")
	for i, element := range n.Entries {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *UsingExpressionNode) Error() string {
	return n.Inspect()
}

// Create a using expression node eg. `using Foo`
func NewUsingExpressionNode(span *position.Span, consts []UsingEntryNode) *UsingExpressionNode {
	return &UsingExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Entries:       consts,
	}
}
