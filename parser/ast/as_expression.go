package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents an as type downcast eg. `foo as String`
type AsExpressionNode struct {
	TypedNodeBase
	Value       ExpressionNode
	RuntimeType ComplexConstantNode
}

func (*AsExpressionNode) IsStatic() bool {
	return false
}

func (*AsExpressionNode) Class() *value.Class {
	return value.PublicIdentifierNodeClass
}

func (*AsExpressionNode) DirectClass() *value.Class {
	return value.PublicIdentifierNodeClass
}

func (n *AsExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::AsExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  value: ")
	indentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

	buff.WriteString(",\n  runtime_type: ")
	indentStringFromSecondLine(&buff, n.RuntimeType.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *AsExpressionNode) Error() string {
	return n.Inspect()
}

// Create a new private constant node eg. `_Foo`.
func NewAsExpressionNode(span *position.Span, val ExpressionNode, runtimeType ComplexConstantNode) *AsExpressionNode {
	return &AsExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
		RuntimeType:   runtimeType,
	}
}
