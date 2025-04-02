package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a singleton type eg. `&String`
type SingletonTypeNode struct {
	TypedNodeBase
	TypeNode TypeNode // right hand side
}

func (n *SingletonTypeNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*SingletonTypeNode)
	if !ok {
		return false
	}

	return n.TypeNode.Equal(value.Ref(o.TypeNode)) &&
		n.span.Equal(o.span)
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
func NewSingletonTypeNode(span *position.Span, typ TypeNode) *SingletonTypeNode {
	return &SingletonTypeNode{
		TypedNodeBase: TypedNodeBase{span: span},
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

	fmt.Fprintf(&buff, "Std::Elk::AST::SingletonTypeNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  type_node: ")
	indent.IndentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *SingletonTypeNode) Error() string {
	return n.Inspect()
}
