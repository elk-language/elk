package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents an unquoted piecie of AST inside of a quote
type UnquoteExpressionNode struct {
	TypedNodeBase
	Expression ExpressionNode
}

func (n *UnquoteExpressionNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	if args == nil || len(*args) == 0 {
		panic("too few arguments for splicing AST nodes")
	}

	arg := (*args)[0]
	*args = (*args)[1:]

	targetLoc := loc.Copy()
	targetLoc.Parent = n.loc

	return arg.Splice(targetLoc, nil, true)
}

// Check if this node equals another node.
func (n *UnquoteExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UnquoteExpressionNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.Expression.Equal(value.Ref(o.Expression))
}

// Return a string representation of the node.
func (n *UnquoteExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("unquote(")

	exprStr := n.Expression.String()
	if strings.ContainsRune(exprStr, '\n') {
		buff.WriteRune('\n')
		indent.IndentString(&buff, exprStr, 1)
		buff.WriteRune('\n')
	} else {
		buff.WriteString(exprStr)
	}

	buff.WriteRune(')')

	return buff.String()
}

func (*UnquoteExpressionNode) IsStatic() bool {
	return false
}

func (*UnquoteExpressionNode) Class() *value.Class {
	return value.MacroBoundaryNodeClass
}

func (*UnquoteExpressionNode) DirectClass() *value.Class {
	return value.MacroBoundaryNodeClass
}

func (n *UnquoteExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::UnquoteExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  expression: ")
	indent.IndentStringFromSecondLine(&buff, n.Expression.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *UnquoteExpressionNode) Error() string {
	return n.Inspect()
}

// Create an unquote expression node eg.
//
//	unquote(x)
func NewUnquoteExpressionNode(loc *position.Location, expr ExpressionNode) *UnquoteExpressionNode {
	return &UnquoteExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Expression:    expr,
	}
}
