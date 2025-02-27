package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents an `unless` expression eg. `unless foo then println("bar")`
type UnlessExpressionNode struct {
	TypedNodeBase
	Condition ExpressionNode  // unless condition
	ThenBody  []StatementNode // then expression body
	ElseBody  []StatementNode // else expression body
}

func (*UnlessExpressionNode) IsStatic() bool {
	return false
}

// Create a new `unless` expression node eg. `unless foo then println("bar")`
func NewUnlessExpressionNode(span *position.Span, cond ExpressionNode, then, els []StatementNode) *UnlessExpressionNode {
	return &UnlessExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		ThenBody:      then,
		Condition:     cond,
		ElseBody:      els,
	}
}

func (*UnlessExpressionNode) Class() *value.Class {
	return value.UnlessExpressionNodeClass
}

func (*UnlessExpressionNode) DirectClass() *value.Class {
	return value.UnlessExpressionNodeClass
}

func (n *UnlessExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::UnlessExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  condition: ")
	indentStringFromSecondLine(&buff, n.Condition.Inspect(), 1)

	buff.WriteString(",\n  then_body: %%[\n")
	for i, stmt := range n.ThenBody {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  else_body: %%[\n")
	for i, stmt := range n.ElseBody {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *UnlessExpressionNode) Error() string {
	return n.Inspect()
}
