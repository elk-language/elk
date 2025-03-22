package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents an enhance expression eg. `implement Enumerable[V]`
type ImplementExpressionNode struct {
	TypedNodeBase
	Constants []ComplexConstantNode
}

func (*ImplementExpressionNode) SkipTypechecking() bool {
	return false
}

func (*ImplementExpressionNode) IsStatic() bool {
	return false
}

// Create an enhance expression node eg. `implement Enumerable[V]`
func NewImplementExpressionNode(span *position.Span, consts []ComplexConstantNode) *ImplementExpressionNode {
	return &ImplementExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Constants:     consts,
	}
}

func (*ImplementExpressionNode) Class() *value.Class {
	return value.ImplementExpressionNodeClass
}

func (*ImplementExpressionNode) DirectClass() *value.Class {
	return value.ImplementExpressionNodeClass
}

func (n *ImplementExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ImplementExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  constants: %%[\n")
	for i, element := range n.Constants {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ImplementExpressionNode) Error() string {
	return n.Inspect()
}
