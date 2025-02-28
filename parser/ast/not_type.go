package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a not type eg. `~String`
type NotTypeNode struct {
	TypedNodeBase
	TypeNode TypeNode // right hand side
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

	fmt.Fprintf(&buff, "Std::AST::NotTypeNode{\n  &: %p", n)

	buff.WriteString(",\n  type_node: ")
	indentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *NotTypeNode) Error() string {
	return n.Inspect()
}
