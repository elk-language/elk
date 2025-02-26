package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
)

// Postfix expression eg. `foo++`, `bar--`
type PostfixExpressionNode struct {
	TypedNodeBase
	Op         *token.Token // operator
	Expression ExpressionNode
}

func (i *PostfixExpressionNode) IsStatic() bool {
	return false
}

// Create a new postfix expression node.
func NewPostfixExpressionNode(span *position.Span, op *token.Token, expr ExpressionNode) *PostfixExpressionNode {
	return &PostfixExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		Expression:    expr,
	}
}

func (*PostfixExpressionNode) Class() *value.Class {
	return value.PostfixExpressionNodeClass
}

func (*PostfixExpressionNode) DirectClass() *value.Class {
	return value.PostfixExpressionNodeClass
}

func (n *PostfixExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::PostfixExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  op: ")
	indentStringFromSecondLine(&buff, n.Op.Inspect(), 1)

	buff.WriteString(",\n  expression: ")
	indentStringFromSecondLine(&buff, n.Expression.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (p *PostfixExpressionNode) Error() string {
	return p.Inspect()
}
