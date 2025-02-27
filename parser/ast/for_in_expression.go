package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a `for in` expression eg. `for i in 5..15 then println(i)`
type ForInExpressionNode struct {
	TypedNodeBase
	Pattern      PatternNode
	InExpression ExpressionNode  // expression that will be iterated through
	ThenBody     []StatementNode // then expression body
}

func (*ForInExpressionNode) IsStatic() bool {
	return false
}

// Create a new `for in` expression node eg. `for i in 5..15 then println(i)`
func NewForInExpressionNode(span *position.Span, pattern PatternNode, inExpr ExpressionNode, then []StatementNode) *ForInExpressionNode {
	return &ForInExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Pattern:       pattern,
		InExpression:  inExpr,
		ThenBody:      then,
	}
}

func (*ForInExpressionNode) Class() *value.Class {
	return value.ForInExpressionNodeClass
}

func (*ForInExpressionNode) DirectClass() *value.Class {
	return value.ForInExpressionNodeClass
}

func (n *ForInExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::ForInExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  pattern: ")
	indentStringFromSecondLine(&buff, n.Pattern.Inspect(), 1)

	buff.WriteString(",\n  in_expression: ")
	indentStringFromSecondLine(&buff, n.InExpression.Inspect(), 1)

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

func (n *ForInExpressionNode) Error() string {
	return n.Inspect()
}
