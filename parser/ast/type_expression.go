package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a type expression `type String?`
type TypeExpressionNode struct {
	NodeBase
	TypeNode TypeNode
}

func (n *TypeExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*TypeExpressionNode)
	if !ok {
		return false
	}

	return n.TypeNode.Equal(value.Ref(o.TypeNode)) &&
		n.span.Equal(o.span)
}

func (n *TypeExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("type ")
	buff.WriteString(n.TypeNode.String())

	return buff.String()
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
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::TypeExpressionNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  type_node: ")
	indent.IndentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
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
