package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a `while` expression eg. `while i < 5 then i += 5`
type WhileExpressionNode struct {
	TypedNodeBase
	Condition ExpressionNode  // while condition
	ThenBody  []StatementNode // then expression body
}

func (*WhileExpressionNode) IsStatic() bool {
	return false
}

// Create a new `while` expression node eg. `while i < 5 then i += 5`
func NewWhileExpressionNode(span *position.Span, cond ExpressionNode, then []StatementNode) *WhileExpressionNode {
	return &WhileExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Condition:     cond,
		ThenBody:      then,
	}
}

func (*WhileExpressionNode) Class() *value.Class {
	return value.WhileExpressionNodeClass
}

func (*WhileExpressionNode) DirectClass() *value.Class {
	return value.WhileExpressionNodeClass
}

func (n *WhileExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::WhileExpressionNode{\n  &: %p", n)

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

	buff.WriteString("\n}")

	return buff.String()
}

func (n *WhileExpressionNode) Error() string {
	return n.Inspect()
}
