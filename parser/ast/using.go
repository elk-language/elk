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

func (n *UsingAllEntryNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &UsingAllEntryNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Namespace:     n.Namespace.Splice(loc, args, unquote).(UsingEntryNode),
	}
}

func (n *UsingAllEntryNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Namespace.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (n *UsingAllEntryNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UsingAllEntryNode)
	if !ok {
		return false
	}

	return n.Namespace.Equal(value.Ref(o.Namespace)) &&
		n.loc.Equal(o.loc)
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

	fmt.Fprintf(&buff, "Std::Elk::AST::UsingAllEntryNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  namespace: ")
	indent.IndentStringFromSecondLine(&buff, n.Namespace.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *UsingAllEntryNode) Error() string {
	return n.Inspect()
}

// Create a new using all entry node eg. `Foo::*`, `A::B::C::*`
func NewUsingAllEntryNode(loc *position.Location, namespace UsingEntryNode) *UsingAllEntryNode {
	return &UsingAllEntryNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Namespace:     namespace,
	}
}

// Represents a using entry node with subentries eg. `Foo::{Bar, baz}`, `A::B::C::{lol, foo as epic, Gro as Moe}`
type UsingEntryWithSubentriesNode struct {
	NodeBase
	Namespace  UsingEntryNode
	Subentries []UsingSubentryNode
}

func (n *UsingEntryWithSubentriesNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &UsingEntryWithSubentriesNode{
		NodeBase:   NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Namespace:  n.Namespace.Splice(loc, args, unquote).(UsingEntryNode),
		Subentries: SpliceSlice(n.Subentries, loc, args, unquote),
	}
}

func (n *UsingEntryWithSubentriesNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Namespace.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	for _, entry := range n.Subentries {
		if entry.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

func (n *UsingEntryWithSubentriesNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UsingEntryWithSubentriesNode)
	if !ok {
		return false
	}

	if len(n.Subentries) != len(o.Subentries) ||
		!n.Namespace.Equal(value.Ref(o.Namespace)) ||
		!n.loc.Equal(o.loc) {
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
func NewUsingEntryWithSubentriesNode(loc *position.Location, namespace UsingEntryNode, subentries []UsingSubentryNode) *UsingEntryWithSubentriesNode {
	return &UsingEntryWithSubentriesNode{
		NodeBase:   NodeBase{loc: loc},
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

	fmt.Fprintf(&buff, "Std::Elk::AST::UsingEntryWithSubentriesNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

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

func (n *UsingExpressionNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &UsingExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Entries:       SpliceSlice(n.Entries, loc, args, unquote),
	}
}

func (n *UsingExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	for _, entry := range n.Entries {
		if entry.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

func (n *UsingExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UsingExpressionNode)
	if !ok {
		return false
	}

	if len(n.Entries) != len(o.Entries) ||
		!n.loc.Equal(o.loc) {
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

	fmt.Fprintf(&buff, "Std::Elk::AST::UsingExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

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
func NewUsingExpressionNode(loc *position.Location, consts []UsingEntryNode) *UsingExpressionNode {
	return &UsingExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Entries:       consts,
	}
}
