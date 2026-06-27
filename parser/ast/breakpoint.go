package ast

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// `breakpoint` entry
type BreakpointNode struct {
	NodeBase
	TypecheckerContext value.Reference
}

func (n *BreakpointNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &BreakpointNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
	}
}

func (n *BreakpointNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::BreakpointNode", env)
}

func (n *BreakpointNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n *BreakpointNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*BreakpointNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc)
}

func (n *BreakpointNode) String() string {
	return "breakpoint"
}

func (*BreakpointNode) SetType(types.Type) {}

func (*BreakpointNode) IsStatic() bool {
	return true
}

func (*BreakpointNode) Type(globalEnv *types.GlobalEnvironment) types.Type {
	return types.Nil{}
}

func (*BreakpointNode) Class() *value.Class {
	return value.BreakpointNodeClass
}

func (*BreakpointNode) DirectClass() *value.Class {
	return value.BreakpointNodeClass
}

func (n *BreakpointNode) Inspect() string {
	return fmt.Sprintf("Std::Elk::AST::BreakpointNode{location: %s}", (*value.Location)(n.loc).Inspect())
}

func (n *BreakpointNode) ToValue() value.Value {
	return value.Ref(n)
}

func (n *BreakpointNode) Error() string {
	return n.Inspect()
}

// Create a new `breakpoint` entry node.
func NewBreakpointNode(loc *position.Location) *BreakpointNode {
	return &BreakpointNode{
		NodeBase: NodeBase{loc: loc},
	}
}
