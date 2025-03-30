package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a not type eg. `~String`
type NotTypeNode struct {
	TypedNodeBase
	TypeNode TypeNode // right hand side
}

func (n *NotTypeNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*NotTypeNode)
	if !ok {
		return false
	}

	return n.span.Equal(o.span) &&
		n.TypeNode.Equal(value.Ref(o.TypeNode))
}

func (n *NotTypeNode) String() string {
	var buff strings.Builder

	buff.WriteRune('~')

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

func (*NotTypeNode) IsStatic() bool {
	return false
}

// Create a new not type node eg. `~String`
func NewNotTypeNode(span *position.Span, typ TypeNode) *NotTypeNode {
	return &NotTypeNode{
		TypedNodeBase: TypedNodeBase{span: span},
		TypeNode:      typ,
	}
}

func (*NotTypeNode) Class() *value.Class {
	return value.NotTypeNodeClass
}

func (*NotTypeNode) DirectClass() *value.Class {
	return value.NotTypeNodeClass
}

func (n *NotTypeNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::NotTypeNode{\n  &: %p", n)

	buff.WriteString(",\n  type_node: ")
	indent.IndentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *NotTypeNode) Error() string {
	return n.Inspect()
}
