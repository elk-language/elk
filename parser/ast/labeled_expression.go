package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a labeled expression eg. `$foo: 1 + 2`
type LabeledExpressionNode struct {
	NodeBase
	Label      string
	Expression ExpressionNode
}

func (n *LabeledExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*LabeledExpressionNode)
	if !ok {
		return false
	}

	return n.span.Equal(o.span) &&
		n.Label == o.Label &&
		n.Expression.Equal(value.Ref(o.Expression))
}

func (n *LabeledExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteRune('$')
	buff.WriteString(n.Label)
	buff.WriteString(": ")

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

func (l *LabeledExpressionNode) Type(env *types.GlobalEnvironment) types.Type {
	return l.Expression.Type(env)
}

func (l *LabeledExpressionNode) IsStatic() bool {
	return l.Expression.IsStatic()
}

// Create a new labeled expression node eg. `$foo: 1 + 2`
func NewLabeledExpressionNode(span *position.Span, label string, expr ExpressionNode) *LabeledExpressionNode {
	return &LabeledExpressionNode{
		NodeBase:   NodeBase{span: span},
		Label:      label,
		Expression: expr,
	}
}

func (*LabeledExpressionNode) Class() *value.Class {
	return value.LabeledExpressionNodeClass
}

func (*LabeledExpressionNode) DirectClass() *value.Class {
	return value.LabeledExpressionNodeClass
}

func (n *LabeledExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::LabeledExpressionNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  label: ")
	buff.WriteString(n.Label)

	buff.WriteString(",\n  expression: ")
	indent.IndentStringFromSecondLine(&buff, n.Expression.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *LabeledExpressionNode) Error() string {
	return n.Inspect()
}
