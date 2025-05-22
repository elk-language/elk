package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a singleton type eg. `&String`
type SingletonTypeNode struct {
	TypedNodeBase
	TypeNode TypeNode // right hand side
}

func (n *SingletonTypeNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &SingletonTypeNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		TypeNode:      n.TypeNode.splice(loc, args, unquote).(TypeNode),
	}
}

func (n *SingletonTypeNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::SingletonTypeNode", env)
}

func (n *SingletonTypeNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
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

func (n *SingletonTypeNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*SingletonTypeNode)
	if !ok {
		return false
	}

	return n.TypeNode.Equal(value.Ref(o.TypeNode)) &&
		n.loc.Equal(o.loc)
}

func (n *SingletonTypeNode) String() string {
	var buff strings.Builder

	buff.WriteRune('&')
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

func (*SingletonTypeNode) IsStatic() bool {
	return false
}

// Create a new singleton type node eg. `&String`
func NewSingletonTypeNode(loc *position.Location, typ TypeNode) *SingletonTypeNode {
	return &SingletonTypeNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		TypeNode:      typ,
	}
}

func (*SingletonTypeNode) Class() *value.Class {
	return value.SingletonTypeNodeClass
}

func (*SingletonTypeNode) DirectClass() *value.Class {
	return value.SingletonTypeNodeClass
}

func (n *SingletonTypeNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::SingletonTypeNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  type_node: ")
	indent.IndentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *SingletonTypeNode) Error() string {
	return n.Inspect()
}
