package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents an instance type eg. `^self`
type InstanceOfTypeNode struct {
	TypedNodeBase
	TypeNode TypeNode // right hand side
}

func (*InstanceOfTypeNode) IsStatic() bool {
	return false
}

// Create a new instance of type node eg. `^self`
func NewInstanceOfTypeNode(span *position.Span, typ TypeNode) *InstanceOfTypeNode {
	return &InstanceOfTypeNode{
		TypedNodeBase: TypedNodeBase{span: span},
		TypeNode:      typ,
	}
}

func (*InstanceOfTypeNode) Class() *value.Class {
	return value.InstanceOfTypeNodeClass
}

func (*InstanceOfTypeNode) DirectClass() *value.Class {
	return value.InstanceOfTypeNodeClass
}

func (n *InstanceOfTypeNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::InstanceOfTypeNode{\n  &: %p", n)

	buff.WriteString(",\n  type_node: ")
	indentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *InstanceOfTypeNode) Error() string {
	return n.Inspect()
}
