package ast

import (
	"fmt"

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

func (p *PostfixExpressionNode) Inspect() string {
	return fmt.Sprintf("Std::AST::PostfixExpressionNode{&: %p, op: %s, expression: %s}", p, p.Op.Inspect(), p.Expression.Inspect())
}

func (p *PostfixExpressionNode) Error() string {
	return p.Inspect()
}
