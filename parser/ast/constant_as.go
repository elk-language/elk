package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a constant with as in using declarations
// eg. `Foo::Bar as Bar`.
type ConstantAsNode struct {
	NodeBase
	Constant ComplexConstantNode
	AsName   string
}

func (n *ConstantAsNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &ConstantAsNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Constant: n.Constant.splice(loc, args, unquote).(ComplexConstantNode),
		AsName:   n.AsName,
	}
}

func (n *ConstantAsNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::ConstantAsNode", env)
}

func (n *ConstantAsNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Constant.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

// Check if this node equals another node.
func (n *ConstantAsNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ConstantAsNode)
	if !ok {
		return false
	}

	return n.AsName == o.AsName &&
		n.Constant.Equal(value.Ref(o.Constant)) &&
		n.loc.Equal(o.loc)
}

// Return a string representation of the node.
func (n *ConstantAsNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.Constant.String())
	buff.WriteString(" as ")
	buff.WriteString(n.AsName)

	return buff.String()
}

func (*ConstantAsNode) IsStatic() bool {
	return false
}

func (*ConstantAsNode) Class() *value.Class {
	return value.ConstantAsNodeClass
}

func (*ConstantAsNode) DirectClass() *value.Class {
	return value.ConstantAsNodeClass
}

func (n *ConstantAsNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ConstantAsNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  constant: ")
	indent.IndentStringFromSecondLine(&buff, n.Constant.Inspect(), 1)

	buff.WriteString(",\n  as_name: ")
	buff.WriteString(n.AsName)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ConstantAsNode) Error() string {
	return n.Inspect()
}

// Create a new identifier with as eg. `Foo::Bar as Bar`.
func NewConstantAsNode(loc *position.Location, constant ComplexConstantNode, as string) *ConstantAsNode {
	return &ConstantAsNode{
		NodeBase: NodeBase{loc: loc},
		Constant: constant,
		AsName:   as,
	}
}
