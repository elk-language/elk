package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
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

func (n *PostfixExpressionNode) Splice(loc *position.Location, args *[]Node) Node {
	return &PostfixExpressionNode{
		TypedNodeBase: n.TypedNodeBase,
		Op:            n.Op,
		Expression:    n.Expression.Splice(loc, args).(ExpressionNode),
	}
}

func (n *PostfixExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*PostfixExpressionNode)
	if !ok {
		return false
	}

	return n.Op.Equal(o.Op) &&
		n.Expression.Equal(value.Ref(o.Expression)) &&
		n.loc.Equal(o.loc)
}

func (n *PostfixExpressionNode) String() string {
	var buff strings.Builder

	parens := ExpressionPrecedence(n) > ExpressionPrecedence(n.Expression)
	if parens {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Expression.String())
	if parens {
		buff.WriteRune(')')
	}

	buff.WriteString(n.Op.String())

	return buff.String()
}

func (i *PostfixExpressionNode) IsStatic() bool {
	return false
}

// Create a new postfix expression node.
func NewPostfixExpressionNode(loc *position.Location, op *token.Token, expr ExpressionNode) *PostfixExpressionNode {
	return &PostfixExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
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

	fmt.Fprintf(&buff, "Std::Elk::AST::PostfixExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  op: ")
	indent.IndentStringFromSecondLine(&buff, n.Op.Inspect(), 1)

	buff.WriteString(",\n  expression: ")
	indent.IndentStringFromSecondLine(&buff, n.Expression.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (p *PostfixExpressionNode) Error() string {
	return p.Inspect()
}
