package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents an `extend where` block expression eg.
//
//	extend where T < Foo
//		def hello then println("awesome!")
//	end
type ExtendWhereBlockExpressionNode struct {
	TypedNodeBase
	Body  []StatementNode
	Where []TypeParameterNode
}

func (*ExtendWhereBlockExpressionNode) SkipTypechecking() bool {
	return false
}

func (*ExtendWhereBlockExpressionNode) IsStatic() bool {
	return false
}

// Create a new `singleton` block expression node eg.
//
//	singleton
//		def hello then println("awesome!")
//	end
func NewExtendWhereBlockExpressionNode(span *position.Span, body []StatementNode, where []TypeParameterNode) *ExtendWhereBlockExpressionNode {
	return &ExtendWhereBlockExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Body:          body,
		Where:         where,
	}
}

func (*ExtendWhereBlockExpressionNode) Class() *value.Class {
	return value.ExtendWhereBlockExpressionNodeClass
}

func (*ExtendWhereBlockExpressionNode) DirectClass() *value.Class {
	return value.ExtendWhereBlockExpressionNodeClass
}

func (n *ExtendWhereBlockExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ExtendWhereBlockExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  body: %%[\n")
	for i, element := range n.Body {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  where: %%[\n")
	for i, element := range n.Where {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ExtendWhereBlockExpressionNode) Error() string {
	return n.Inspect()
}
