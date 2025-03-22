package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a numeric `for` expression eg. `fornum i := 0; i < 10; i += 1 then println(i)`
type NumericForExpressionNode struct {
	TypedNodeBase
	Initialiser ExpressionNode  // i := 0
	Condition   ExpressionNode  // i < 10
	Increment   ExpressionNode  // i += 1
	ThenBody    []StatementNode // then expression body
}

func (*NumericForExpressionNode) IsStatic() bool {
	return false
}

// Create a new numeric `fornum` expression eg. `for i := 0; i < 10; i += 1 then println(i)`
func NewNumericForExpressionNode(span *position.Span, init, cond, incr ExpressionNode, then []StatementNode) *NumericForExpressionNode {
	return &NumericForExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Initialiser:   init,
		Condition:     cond,
		Increment:     incr,
		ThenBody:      then,
	}
}

func (*NumericForExpressionNode) Class() *value.Class {
	return value.NumericForExpressionNodeClass
}

func (*NumericForExpressionNode) DirectClass() *value.Class {
	return value.NumericForExpressionNodeClass
}

func (n *NumericForExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::NumericForExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  initialiser: ")
	indent.IndentStringFromSecondLine(&buff, n.Initialiser.Inspect(), 1)

	buff.WriteString(",\n  condition: ")
	indent.IndentStringFromSecondLine(&buff, n.Condition.Inspect(), 1)

	buff.WriteString(",\n  increment: ")
	indent.IndentStringFromSecondLine(&buff, n.Increment.Inspect(), 1)

	buff.WriteString(",\n  then_body: %%[\n")
	for i, stmt := range n.ThenBody {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *NumericForExpressionNode) Error() string {
	return n.Inspect()
}
