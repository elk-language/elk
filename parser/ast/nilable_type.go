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

func (n *NilableTypeNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*NilableTypeNode)
	if !ok {
		return false
	}

	return n.span.Equal(o.span) &&
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
func NewNilableTypeNode(span *position.Span, typ TypeNode) *NilableTypeNode {
	return &NilableTypeNode{
		TypedNodeBase: TypedNodeBase{span: span},
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

	fmt.Fprintf(&buff, "Std::Elk::AST::NilableTypeNode{\n  &: %p", n)

	buff.WriteString(",\n  type_node: ")
	indent.IndentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *NilableTypeNode) Error() string {
	return n.Inspect()
}
