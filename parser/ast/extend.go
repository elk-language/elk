package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
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

// Check if this node equals another node.
func (n *ExtendWhereBlockExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ExtendWhereBlockExpressionNode)
	if !ok {
		return false
	}

	if len(n.Body) != len(o.Body) ||
		len(n.Where) != len(o.Where) {
		return false
	}

	for i, stmt := range n.Body {
		if !stmt.Equal(value.Ref(o.Body[i])) {
			return false
		}
	}

	for i, param := range n.Where {
		if !param.Equal(value.Ref(o.Where[i])) {
			return false
		}
	}

	return n.span.Equal(o.span)
}

// Return a string representation of the node.
func (n *ExtendWhereBlockExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("extend where ")

	for i, param := range n.Where {
		if i != 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(param.String())
	}

	buff.WriteString("\n")

	for _, stmt := range n.Body {
		indent.IndentString(&buff, stmt.String(), 1)
		buff.WriteString("\n")
	}

	buff.WriteString("end")

	return buff.String()
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

	fmt.Fprintf(&buff, "Std::Elk::AST::ExtendWhereBlockExpressionNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  body: %[\n")
	for i, element := range n.Body {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  where: %[\n")
	for i, element := range n.Where {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ExtendWhereBlockExpressionNode) Error() string {
	return n.Inspect()
}
