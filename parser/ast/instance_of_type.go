package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents an instance type eg. `%self`
type InstanceOfTypeNode struct {
	TypedNodeBase
	TypeNode TypeNode // right hand side
}

// Equal checks if this node equals the other node.
func (n *InstanceOfTypeNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*InstanceOfTypeNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.TypeNode.Equal(value.Ref(o.TypeNode))
}

// String returns the string representation of this node.
func (n *InstanceOfTypeNode) String() string {
	var buff strings.Builder

	buff.WriteRune('%')

	parens := TypePrecedence(n) > TypePrecedence(n.TypeNode)
	if parens {
		buff.WriteRune('(')
	}
	buff.WriteString(n.TypeNode.String())
	if parens {
		buff.WriteRune(')')
	}

	return buff.String()
}

func (*InstanceOfTypeNode) IsStatic() bool {
	return false
}

// Create a new instance of type node eg. `%self`
func NewInstanceOfTypeNode(loc *position.Location, typ TypeNode) *InstanceOfTypeNode {
	return &InstanceOfTypeNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		TypeNode:      typ,
	}
}

func (*InstanceOfTypeNode) Class() *value.Class {
	return value.InstanceOfTypeNodeClass
}

func (*InstanceOfTypeNode) DirectClass() *value.Class {
	return value.InstanceOfTypeNodeClass
}

func (n *InstanceOfTypeNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::InstanceOfTypeNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  type_node: ")
	indent.IndentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *InstanceOfTypeNode) Error() string {
	return n.Inspect()
}
