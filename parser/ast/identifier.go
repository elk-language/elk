package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
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
func (*UnquoteNode) identifierNode()           {}

// Represents a public identifier eg. `foo`.
type PublicIdentifierNode struct {
	TypedNodeBase
	Value string
}

func (n *PublicIdentifierNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &PublicIdentifierNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Value:         n.Value,
	}
}

func (n *PublicIdentifierNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::PublicIdentifierNode", env)
}

func (n *PublicIdentifierNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n *PublicIdentifierNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*PublicIdentifierNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.loc.Equal(o.loc)
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
	return fmt.Sprintf(
		"Std::Elk::AST::PublicIdentifierNode{location: %s, value: %s}",
		(*value.Location)(n.loc).Inspect(),
		n.Value,
	)
}

func (n *PublicIdentifierNode) Error() string {
	return n.Inspect()
}

// Create a new public identifier node eg. `foo`.
func NewPublicIdentifierNode(loc *position.Location, val string) *PublicIdentifierNode {
	return &PublicIdentifierNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
	}
}

// Represents a private identifier eg. `_foo`
type PrivateIdentifierNode struct {
	TypedNodeBase
	Value string
}

func (n *PrivateIdentifierNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &PrivateIdentifierNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Value:         n.Value,
	}
}

func (n *PrivateIdentifierNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::PrivateIdentifierNode", env)
}

func (n *PrivateIdentifierNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n *PrivateIdentifierNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*PrivateIdentifierNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.loc.Equal(o.loc)
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
	return fmt.Sprintf(
		"Std::Elk::AST::PrivateIdentifierNode{location: %s, value: %s}",
		(*value.Location)(n.loc).Inspect(),
		n.Value,
	)
}

func (n *PrivateIdentifierNode) Error() string {
	return n.Inspect()
}

// Create a new private identifier node eg. `_foo`.
func NewPrivateIdentifierNode(loc *position.Location, val string) *PrivateIdentifierNode {
	return &PrivateIdentifierNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
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

func (n *PublicIdentifierAsNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &PublicIdentifierAsNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Target:   n.Target.splice(loc, args, unquote).(*PublicIdentifierNode),
		AsName:   n.AsName,
	}
}

func (n *PublicIdentifierAsNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::PublicIdentifierAsNode", env)
}

func (n *PublicIdentifierAsNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Target.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (n *PublicIdentifierAsNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*PublicIdentifierAsNode)
	if !ok {
		return false
	}

	return n.Target.Equal(value.Ref(o.Target)) &&
		n.AsName == o.AsName &&
		n.loc.Equal(o.loc)
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

	fmt.Fprintf(&buff, "Std::Elk::AST::PublicIdentifierAsNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

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
func NewPublicIdentifierAsNode(loc *position.Location, target *PublicIdentifierNode, as string) *PublicIdentifierAsNode {
	return &PublicIdentifierAsNode{
		NodeBase: NodeBase{loc: loc},
		Target:   target,
		AsName:   as,
	}
}
