package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// An expression that matches a value against a pattern eg. `a match Foo(a: "lol", b: > 2)`
type MatchExpressionNode struct {
	TypedNodeBase
	Expression ExpressionNode // left hand side
	Pattern    PatternNode    // right hand side
}

func (n *MatchExpressionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	expr := n.Expression.splice(loc, args, unquote).(ExpressionNode)
	pattern := n.Pattern.splice(loc, args, unquote).(PatternNode)

	return &MatchExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Expression:    expr,
		Pattern:       pattern,
	}
}

func (n *MatchExpressionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::MatchExpressionNode", env)
}

func (n *MatchExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Expression.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}
	if n.Pattern.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (n *MatchExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*MatchExpressionNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.Expression.Equal(value.Ref(o.Expression)) &&
		n.Pattern.Equal(value.Ref(o.Pattern))
}

func (n *MatchExpressionNode) String() string {
	var buff strings.Builder

	leftParen := ExpressionPrecedence(n) > ExpressionPrecedence(n.Expression)
	if leftParen {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Expression.String())
	if leftParen {
		buff.WriteRune(')')
	}

	buff.WriteString(" match ")
	buff.WriteString(n.Pattern.String())

	return buff.String()
}

func (b *MatchExpressionNode) IsStatic() bool {
	return false
}

// Create a new match expression node.
func NewMatchExpressionNode(loc *position.Location, expression ExpressionNode, pattern PatternNode) *MatchExpressionNode {
	return &MatchExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Expression:    expression,
		Pattern:       pattern,
	}
}

func (*MatchExpressionNode) Class() *value.Class {
	return value.MatchExpressionNodeClass
}

func (*MatchExpressionNode) DirectClass() *value.Class {
	return value.MatchExpressionNodeClass
}

func (n *MatchExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::MatchExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  expression: ")
	indent.IndentStringFromSecondLine(&buff, n.Expression.Inspect(), 1)

	buff.WriteString(",\n  pattern: ")
	indent.IndentStringFromSecondLine(&buff, n.Pattern.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *MatchExpressionNode) Error() string {
	return n.Inspect()
}
