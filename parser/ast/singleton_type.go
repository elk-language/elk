package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a singleton type eg. `&String`
type SingletonTypeNode struct {
	TypedNodeBase
	TypeNode TypeNode // right hand side
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

	fmt.Fprintf(&buff, "Std::AST::SingletonTypeNode{\n  &: %p", n)

	buff.WriteString(",\n  type_node: ")
	indentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *SingletonTypeNode) Error() string {
	return n.Inspect()
}
