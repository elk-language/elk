package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a box type eg. `^Int`, `*Int`
type BoxTypeNode struct {
	TypedNodeBase
	Immutable bool
	TypeNode  TypeNode // right hand side
}

func (n *BoxTypeNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &BoxTypeNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Immutable:     n.Immutable,
		TypeNode:      n.TypeNode.splice(loc, args, unquote).(ComplexConstantNode),
	}
}

func (n *BoxTypeNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::BoxTypeNode", env)
}

func (n *BoxTypeNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.TypeNode.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

// Equal checks if this node equals the other node.
func (n *BoxTypeNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*BoxTypeNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.Immutable == o.Immutable &&
		n.TypeNode.Equal(value.Ref(o.TypeNode))
}

// String returns the string representation of this node.
func (n *BoxTypeNode) String() string {
	var buff strings.Builder

	if n.Immutable {
		buff.WriteRune('*')
	} else {
		buff.WriteRune('^')
	}

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

func (*BoxTypeNode) IsStatic() bool {
	return false
}

// Create a new box type node eg. `^Int`
func NewBoxTypeNode(loc *position.Location, typ TypeNode, immutable bool) *BoxTypeNode {
	return &BoxTypeNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Immutable:     immutable,
		TypeNode:      typ,
	}
}

func (*BoxTypeNode) Class() *value.Class {
	return value.BoxTypeNodeClass
}

func (*BoxTypeNode) DirectClass() *value.Class {
	return value.BoxTypeNodeClass
}

func (n *BoxTypeNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::BoxTypeNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  type_node: ")
	indent.IndentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	fmt.Fprintf(&buff, ",\n  immutable: %t", n.Immutable)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *BoxTypeNode) Error() string {
	return n.Inspect()
}
