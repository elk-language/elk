package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents an include expression eg. `include Enumerable[V]`
type IncludeExpressionNode struct {
	TypedNodeBase
	Constants []ComplexConstantNode
}

func (n *IncludeExpressionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &IncludeExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Constants:     SpliceSlice(n.Constants, loc, args, unquote),
	}
}

func (n *IncludeExpressionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::IncludeExpressionNode", env)
}

func (n *IncludeExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	for _, constant := range n.Constants {
		if constant.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

// Check if this node equals another node.
func (n *IncludeExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*IncludeExpressionNode)
	if !ok {
		return false
	}

	if len(n.Constants) != len(o.Constants) {
		return false
	}

	for i, constant := range n.Constants {
		if !constant.Equal(value.Ref(o.Constants[i])) {
			return false
		}
	}

	return n.loc.Equal(o.loc)
}

// Return a string representation of the node.
func (n *IncludeExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("include ")

	for i, constant := range n.Constants {
		if i != 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(constant.String())
	}

	return buff.String()
}

func (*IncludeExpressionNode) SkipTypechecking() bool {
	return false
}

func (*IncludeExpressionNode) IsStatic() bool {
	return false
}

// Create an include expression node eg. `include Enumerable[V]`
func NewIncludeExpressionNode(loc *position.Location, consts []ComplexConstantNode) *IncludeExpressionNode {
	return &IncludeExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Constants:     consts,
	}
}

func (*IncludeExpressionNode) Class() *value.Class {
	return value.UsingAllEntryNodeClass
}

func (*IncludeExpressionNode) DirectClass() *value.Class {
	return value.UsingAllEntryNodeClass
}

func (n *IncludeExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::IncludeExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  constants: %[\n")
	for i, element := range n.Constants {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *IncludeExpressionNode) Error() string {
	return n.Inspect()
}
