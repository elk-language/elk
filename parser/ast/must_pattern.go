package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a `must` pattern eg. `must as foo`
type MustPatternNode struct {
	TypedNodeBase
}

func (n *MustPatternNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &MustPatternNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
	}
}

func (n *MustPatternNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::MustPatternNode", env)
}

func (n *MustPatternNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n *MustPatternNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*MustPatternNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc)
}

func (n *MustPatternNode) String() string {
	return "must"
}

func (*MustPatternNode) IsStatic() bool {
	return false
}

// Create a new `must` pattern node eg. `must`
func NewMustPatternNode(loc *position.Location) *MustPatternNode {
	return &MustPatternNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
	}
}

func (*MustPatternNode) Class() *value.Class {
	return value.MustPatternNodeClass
}

func (*MustPatternNode) DirectClass() *value.Class {
	return value.MustPatternNodeClass
}

func (n *MustPatternNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::MustPatternNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString("\n}")

	return buff.String()
}

func (n *MustPatternNode) ToValue() value.Value {
	return value.Ref(n)
}

func (n *MustPatternNode) Error() string {
	return n.Inspect()
}
