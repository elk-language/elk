package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Exact type of a class instance eg. `exact Foo`, `exact Bar[String]`
type ExactTypeNode struct {
	TypedNodeBase
	TypeNode TypeNode // right hand side
}

func (n *ExactTypeNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &ExactTypeNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		TypeNode:      n.TypeNode.splice(loc, args, unquote).(TypeNode),
	}
}

func (n *ExactTypeNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::ExactTypeNode", env)
}

func (n *ExactTypeNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
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

func (n *ExactTypeNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ExactTypeNode)
	if !ok {
		return false
	}

	return n.TypeNode.Equal(value.Ref(o.TypeNode)) &&
		n.loc.Equal(o.loc)
}

func (n *ExactTypeNode) String() string {
	var buff strings.Builder

	buff.WriteString("exact ")

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

func (u *ExactTypeNode) IsStatic() bool {
	return false
}

// Create a new exact type node.
func NewExactTypeNode(loc *position.Location, typeNode TypeNode) *ExactTypeNode {
	return &ExactTypeNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		TypeNode:      typeNode,
	}
}

func (*ExactTypeNode) Class() *value.Class {
	return value.ExactTypeNodeClass
}

func (*ExactTypeNode) DirectClass() *value.Class {
	return value.ExactTypeNodeClass
}

func (n *ExactTypeNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ExactTypeNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  type_node: ")
	indent.IndentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ExactTypeNode) ToValue() value.Value {
	return value.Ref(n)
}

func (n *ExactTypeNode) Error() string {
	return n.Inspect()
}
