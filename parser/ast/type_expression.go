package ast

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a type expression `type String?`
type TypeExpressionNode struct {
	NodeBase
	TypeNode TypeNode
}

func (*TypeExpressionNode) IsStatic() bool {
	return false
}

func (*TypeExpressionNode) Class() *value.Class {
	return value.TypeExpressionNodeClass
}

func (*TypeExpressionNode) DirectClass() *value.Class {
	return value.TypeExpressionNodeClass
}

func (n *TypeExpressionNode) Inspect() string {
	return fmt.Sprintf("Std::AST::TypeExpressionNode{&: %p, type_node: %s}", n, n.TypeNode.Inspect())
}

func (n *TypeExpressionNode) Error() string {
	return n.Inspect()
}

// Create a new type expression `type String?`
func NewTypeExpressionNode(span *position.Span, typeNode TypeNode) *TypeExpressionNode {
	return &TypeExpressionNode{
		NodeBase: NodeBase{span: span},
		TypeNode: typeNode,
	}
}
