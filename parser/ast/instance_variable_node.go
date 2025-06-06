package ast

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents an instance variable eg. `@foo`
type InstanceVariableNode struct {
	TypedNodeBase
	Value string
}

func (n *InstanceVariableNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &InstanceVariableNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Value:         n.Value,
	}
}

func (n *InstanceVariableNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::InstanceVariableNode", env)
}

func (n *InstanceVariableNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n *InstanceVariableNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*InstanceVariableNode)
	if !ok {
		return false
	}

	return n.Value != o.Value ||
		!n.loc.Equal(o.loc)
}

func (n *InstanceVariableNode) String() string {
	return fmt.Sprintf("@%s", n.Value)
}

func (*InstanceVariableNode) IsStatic() bool {
	return false
}

// Create an instance variable node eg. `@foo`.
func NewInstanceVariableNode(loc *position.Location, val string) *InstanceVariableNode {
	return &InstanceVariableNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
	}
}

func (*InstanceVariableNode) Class() *value.Class {
	return value.InstanceVariableNodeClass
}

func (*InstanceVariableNode) DirectClass() *value.Class {
	return value.InstanceVariableNodeClass
}

func (i *InstanceVariableNode) Inspect() string {
	return fmt.Sprintf(
		"Std::Elk::AST::InstanceVariableNode{location: %s, value: %s}",
		(*value.Location)(i.loc).Inspect(),
		value.String(i.Value).Inspect(),
	)
}

func (p *InstanceVariableNode) Error() string {
	return p.Inspect()
}
