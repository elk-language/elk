package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents an optional or nilable type eg. `String?`
type NilableTypeNode struct {
	TypedNodeBase
	TypeNode TypeNode // right hand side
}

func (n *NilableTypeNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &NilableTypeNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		TypeNode:      n.TypeNode.splice(loc, args, unquote).(TypeNode),
	}
}

func (n *NilableTypeNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::NilableTypeNode", env)
}

func (n *NilableTypeNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
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

func (n *NilableTypeNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*NilableTypeNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.TypeNode.Equal(value.Ref(o.TypeNode))
}

func (n *NilableTypeNode) String() string {
	var buff strings.Builder

	parens := TypePrecedence(n) > TypePrecedence(n.TypeNode)
	if parens {
		buff.WriteRune('(')
	}
	buff.WriteString(n.TypeNode.String())
	if parens {
		buff.WriteRune(')')
	}

	buff.WriteString("?")

	return buff.String()
}

func (*NilableTypeNode) IsStatic() bool {
	return false
}

// Create a new nilable type node eg. `String?`
func NewNilableTypeNode(loc *position.Location, typ TypeNode) *NilableTypeNode {
	return &NilableTypeNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		TypeNode:      typ,
	}
}

func (*NilableTypeNode) Class() *value.Class {
	return value.NilableTypeNodeClass
}

func (*NilableTypeNode) DirectClass() *value.Class {
	return value.NilableTypeNodeClass
}

func (n *NilableTypeNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::NilableTypeNode{\n  loc: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  type_node: ")
	indent.IndentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *NilableTypeNode) Error() string {
	return n.Inspect()
}
