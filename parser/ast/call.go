package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a new expression eg. `new(123)`
type NewExpressionNode struct {
	TypedNodeBase
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (*NewExpressionNode) IsStatic() bool {
	return false
}

// Create a new expression node eg. `new(123)`
func NewNewExpressionNode(span *position.Span, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *NewExpressionNode {
	return &NewExpressionNode{
		TypedNodeBase:       TypedNodeBase{span: span},
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

func (*NewExpressionNode) Class() *value.Class {
	return value.NewExpressionNodeClass
}

func (*NewExpressionNode) DirectClass() *value.Class {
	return value.NewExpressionNodeClass
}

func (n *NewExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::NewExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  positional_arguments: %%[\n")
	for i, element := range n.PositionalArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  named_arguments: %%[\n")
	for i, element := range n.NamedArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *NewExpressionNode) Error() string {
	return n.Inspect()
}
