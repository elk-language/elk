package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
)

// Expression of an operator with one operand eg. `!foo`, `-bar`
type UnaryExpressionNode struct {
	TypedNodeBase
	Op    *token.Token   // operator
	Right ExpressionNode // right hand side
}

func (n *UnaryExpressionNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &UnaryExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Op:            n.Op.Splice(loc, unquote),
		Right:         n.Right.Splice(loc, args, unquote).(ExpressionNode),
	}
}

func (n *UnaryExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UnaryExpressionNode)
	if !ok {
		return false
	}

	return n.Op.Equal(o.Op) &&
		n.Right.Equal(value.Ref(o.Right)) &&
		n.loc.Equal(o.loc)
}

func (n *UnaryExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.Op.FetchValue())

	parens := ExpressionPrecedence(n) > ExpressionPrecedence(n.Right)
	if parens {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Right.String())
	if parens {
		buff.WriteRune(')')
	}

	return buff.String()
}

func (u *UnaryExpressionNode) IsStatic() bool {
	return u.Right.IsStatic()
}

// Create a new unary expression node.
func NewUnaryExpressionNode(loc *position.Location, op *token.Token, right ExpressionNode) *UnaryExpressionNode {
	return &UnaryExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Op:            op,
		Right:         right,
	}
}

func (*UnaryExpressionNode) Class() *value.Class {
	return value.UnaryExpressionNodeClass
}

func (*UnaryExpressionNode) DirectClass() *value.Class {
	return value.UnaryExpressionNodeClass
}

func (n *UnaryExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::UnaryExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  op: ")
	indent.IndentStringFromSecondLine(&buff, n.Op.Inspect(), 1)

	buff.WriteString(",\n  right: ")
	indent.IndentStringFromSecondLine(&buff, n.Right.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *UnaryExpressionNode) Error() string {
	return n.Inspect()
}
