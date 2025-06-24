package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// All nodes that should be valid instance variables
// should implement this interface.
type InstanceVariableNode interface {
	Node
	ExpressionNode
	ivarNode()
}

func (*InvalidNode) ivarNode()                {}
func (*PublicInstanceVariableNode) ivarNode() {}
func (*UnquoteNode) ivarNode()                {}

// Represents an instance variable eg. `@foo`
type PublicInstanceVariableNode struct {
	TypedNodeBase
	Value string
}

func (n *PublicInstanceVariableNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &PublicInstanceVariableNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Value:         n.Value,
	}
}

func (n *PublicInstanceVariableNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::PublicInstanceVariableNode", env)
}

func (n *PublicInstanceVariableNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n *PublicInstanceVariableNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*PublicInstanceVariableNode)
	if !ok {
		return false
	}

	return n.Value != o.Value ||
		!n.loc.Equal(o.loc)
}

func (n *PublicInstanceVariableNode) String() string {
	var buff strings.Builder
	buff.WriteByte('@')

	if PrefixedIdentifierRegexp.MatchString(n.Value) {
		buff.WriteString(n.Value)
		return buff.String()
	}

	buff.WriteString(value.String(n.Value).Inspect())
	return buff.String()
}

func (*PublicInstanceVariableNode) IsStatic() bool {
	return false
}

// Create an instance variable node eg. `@foo`.
func NewPublicInstanceVariableNode(loc *position.Location, val string) *PublicInstanceVariableNode {
	return &PublicInstanceVariableNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
	}
}

func (*PublicInstanceVariableNode) Class() *value.Class {
	return value.PublicInstanceVariableNodeClass
}

func (*PublicInstanceVariableNode) DirectClass() *value.Class {
	return value.PublicInstanceVariableNodeClass
}

func (i *PublicInstanceVariableNode) Inspect() string {
	return fmt.Sprintf(
		"Std::Elk::AST::PublicInstanceVariableNode{location: %s, value: %s}",
		(*value.Location)(i.loc).Inspect(),
		value.String(i.Value).Inspect(),
	)
}

func (p *PublicInstanceVariableNode) Error() string {
	return p.Inspect()
}
