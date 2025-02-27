package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a `do` expression eg.
//
//	do
//		print("awesome!")
//	end
type DoExpressionNode struct {
	TypedNodeBase
	Body    []StatementNode // do expression body
	Catches []*CatchNode
	Finally []StatementNode
}

func (*DoExpressionNode) IsStatic() bool {
	return false
}

func (*DoExpressionNode) Class() *value.Class {
	return value.DoExpressionNodeClass
}

func (*DoExpressionNode) DirectClass() *value.Class {
	return value.DoExpressionNodeClass
}

func (n *DoExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::DoExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  body: %%[\n")
	for i, stmt := range n.Body {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  catches: %%[\n")
	for i, stmt := range n.Catches {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  finally: %%[\n")
	for i, stmt := range n.Finally {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *DoExpressionNode) Error() string {
	return n.Inspect()
}

// Create a new `do` expression node eg.
//
//	do
//		print("awesome!")
//	end
func NewDoExpressionNode(span *position.Span, body []StatementNode, catches []*CatchNode, finally []StatementNode) *DoExpressionNode {
	return &DoExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Body:          body,
		Catches:       catches,
		Finally:       finally,
	}
}
