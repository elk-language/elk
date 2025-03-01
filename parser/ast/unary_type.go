package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
)

// Type of an operator with one operand eg. `-2`, `+3`
type UnaryTypeNode struct {
	TypedNodeBase
	Op       *token.Token // operator
	TypeNode TypeNode     // right hand side
}

func (u *UnaryTypeNode) IsStatic() bool {
	return false
}

// Create a new unary expression node.
func NewUnaryTypeNode(span *position.Span, op *token.Token, typeNode TypeNode) *UnaryTypeNode {
	return &UnaryTypeNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		TypeNode:      typeNode,
	}
}

func (*UnaryTypeNode) Class() *value.Class {
	return value.UnaryTypeNodeClass
}

func (*UnaryTypeNode) DirectClass() *value.Class {
	return value.UnaryTypeNodeClass
}

func (n *UnaryTypeNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::UnaryTypeNode{\n  &: %p", n)

	buff.WriteString(",\n  op: ")
	indentStringFromSecondLine(&buff, n.Op.Inspect(), 1)

	buff.WriteString(",\n  type_node: ")
	indentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *UnaryTypeNode) Error() string {
	return n.Inspect()
}
