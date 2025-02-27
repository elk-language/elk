package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a `until` expression eg. `until i >= 5 then i += 5`
type UntilExpressionNode struct {
	TypedNodeBase
	Condition ExpressionNode  // until condition
	ThenBody  []StatementNode // then expression body
}

func (*UntilExpressionNode) IsStatic() bool {
	return false
}

// Create a new `until` expression node eg. `until i >= 5 then i += 5`
func NewUntilExpressionNode(span *position.Span, cond ExpressionNode, then []StatementNode) *UntilExpressionNode {
	return &UntilExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Condition:     cond,
		ThenBody:      then,
	}
}

func (*UntilExpressionNode) Class() *value.Class {
	return value.UntilExpressionNodeClass
}

func (*UntilExpressionNode) DirectClass() *value.Class {
	return value.UntilExpressionNodeClass
}

func (n *UntilExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::UntilExpressionNode{\n  &: %p", n)

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

func (n *UntilExpressionNode) Error() string {
	return n.Inspect()
}
