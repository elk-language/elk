package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a `go` expression eg. `go foo()`, `go; foo(); end`
type GoExpressionNode struct {
	TypedNodeBase
	Body []StatementNode
}

func (*GoExpressionNode) IsStatic() bool {
	return false
}

// Create a new `go` expression node eg. `go foo()`, `go; foo(); end`
func NewGoExpressionNode(span *position.Span, body []StatementNode) *GoExpressionNode {
	return &GoExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Body:          body,
	}
}

func (*GoExpressionNode) Class() *value.Class {
	return value.GoExpressionNodeClass
}

func (*GoExpressionNode) DirectClass() *value.Class {
	return value.GoExpressionNodeClass
}

func (n *GoExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::GoExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  body: %%[\n")
	for i, stmt := range n.Body {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *GoExpressionNode) Error() string {
	return n.Inspect()
}
