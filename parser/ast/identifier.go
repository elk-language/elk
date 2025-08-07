package ast

import (
	"fmt"
	"regexp"
	"strings"

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

func IdentifierToString(ident IdentifierNode) string {
	switch ident := ident.(type) {
	case *PublicIdentifierNode:
		return ident.Value
	case *PrivateIdentifierNode:
		return ident.Value
	}

	return ""
}

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

var IdentifierRegexp = regexp.MustCompile(`^\p{Ll}[\p{L}\p{N}_]*$`)
var PrefixedIdentifierRegexp = regexp.MustCompile(`^[\p{L}\p{N}_]+$`)

func (n *PublicIdentifierNode) String() string {
	if IdentifierRegexp.MatchString(n.Value) {
		return n.Value
	}

	var buff strings.Builder
	buff.WriteByte('$')

	if PrefixedIdentifierRegexp.MatchString(n.Value) {
		buff.WriteString(n.Value)
		return buff.String()
	}

	buff.WriteString(value.String(n.Value).Inspect())
	return buff.String()
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

var PrivateIdentifierRegexp = regexp.MustCompile(`^_[\p{L}\p{N}_]*$`)

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
