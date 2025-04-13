package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a `break` expression eg. `break`, `break false`
type BreakExpressionNode struct {
	NodeBase
	Label string
	Value ExpressionNode
}

func (n *BreakExpressionNode) Splice(loc *position.Location, args *[]Node) Node {
	var val ExpressionNode
	if n.Value != nil {
		val = n.Value.Splice(loc, args).(ExpressionNode)
	}

	return &BreakExpressionNode{
		NodeBase: n.NodeBase,
		Label:    n.Label,
		Value:    val,
	}
}

func (n *BreakExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*BreakExpressionNode)
	if !ok {
		return false
	}

	if n.Label != o.Label {
		return false
	}

	if n.Value == o.Value {
	} else if n.Value == nil || o.Value == nil {
		return false
	} else if !n.Value.Equal(value.Ref(o.Value)) {
		return false
	}

	return n.loc.Equal(o.loc)
}

func (n *BreakExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("break")

	if n.Label != "" {
		buff.WriteRune('$')
		buff.WriteString(n.Label)
	}

	if n.Value != nil {
		buff.WriteRune(' ')

		parens := ExpressionPrecedence(n) > ExpressionPrecedence(n.Value)
		if parens {
			buff.WriteRune('(')
		}
		buff.WriteString(n.Value.String())
		if parens {
			buff.WriteRune(')')
		}
	}

	return buff.String()
}

func (*BreakExpressionNode) IsStatic() bool {
	return false
}

func (*BreakExpressionNode) Type(*types.GlobalEnvironment) types.Type {
	return types.Never{}
}

// Create a new `break` expression node eg. `break`
func NewBreakExpressionNode(loc *position.Location, label string, val ExpressionNode) *BreakExpressionNode {
	return &BreakExpressionNode{
		NodeBase: NodeBase{loc: loc},
		Label:    label,
		Value:    val,
	}
}

func (*BreakExpressionNode) Class() *value.Class {
	return value.BreakExpressionNodeClass
}

func (*BreakExpressionNode) DirectClass() *value.Class {
	return value.BreakExpressionNodeClass
}

func (n *BreakExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::BreakExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  label: ")
	buff.WriteString(n.Label)

	buff.WriteString(",\n  value: ")
	if n.Value == nil {
		buff.WriteString("nil")
	} else {
		indent.IndentStringFromSecondLine(&buff, n.Value.Inspect(), 1)
	}

	buff.WriteString("\n}")

	return buff.String()
}

func (n *BreakExpressionNode) Error() string {
	return n.Inspect()
}
