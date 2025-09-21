package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a box of expression eg. `&a`
type BoxOfExpressionNode struct {
	TypedNodeBase
	Expression ExpressionNode // right hand side
}

func (n *BoxOfExpressionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &BoxOfExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Expression:    n.Expression.splice(loc, args, unquote).(ExpressionNode),
	}
}

func (n *BoxOfExpressionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::BoxOfExpressionNode", env)
}

func (n *BoxOfExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Expression.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

// Equal checks if this node equals the other node.
func (n *BoxOfExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*BoxOfExpressionNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.Expression.Equal(value.Ref(o.Expression))
}

// String returns the string representation of this node.
func (n *BoxOfExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteRune('&')

	parens := ExpressionPrecedence(n) > ExpressionPrecedence(n.Expression)
	if parens {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Expression.String())
	if parens {
		buff.WriteRune(')')
	}

	return buff.String()
}

func (*BoxOfExpressionNode) IsStatic() bool {
	return false
}

// Create a new box of expression node eg. `&a`
func NewBoxOfExpressionNode(loc *position.Location, expr ExpressionNode) *BoxOfExpressionNode {
	return &BoxOfExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Expression:    expr,
	}
}

func (*BoxOfExpressionNode) Class() *value.Class {
	return value.BoxOfExpressionNodeClass
}

func (*BoxOfExpressionNode) DirectClass() *value.Class {
	return value.BoxOfExpressionNodeClass
}

func (n *BoxOfExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::BoxOfExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  expression: ")
	indent.IndentStringFromSecondLine(&buff, n.Expression.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *BoxOfExpressionNode) Error() string {
	return n.Inspect()
}
