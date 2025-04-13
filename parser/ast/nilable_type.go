package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents an optional or nilable type eg. `String?`
type NilableTypeNode struct {
	TypedNodeBase
	TypeNode TypeNode // right hand side
}

func (n *NilableTypeNode) Splice(loc *position.Location, args *[]Node) Node {
	return &NilableTypeNode{
		TypedNodeBase: TypedNodeBase{loc: getLoc(loc, n.loc), typ: n.typ},
		TypeNode:      n.TypeNode.Splice(loc, args).(TypeNode),
	}
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
