package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a `yield` expression eg. `yield`, `yield true`, `yield* foo()`
type YieldExpressionNode struct {
	NodeBase
	Value   ExpressionNode
	Forward bool
}

func (n *YieldExpressionNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	var val ExpressionNode
	if n.Value != nil {
		val = n.Value.Splice(loc, args, unquote).(ExpressionNode)
	}

	return &YieldExpressionNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Value:    val,
		Forward:  n.Forward,
	}
}

func (n *YieldExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*YieldExpressionNode)
	if !ok {
		return false
	}

	if !n.loc.Equal(o.loc) ||
		n.Forward != o.Forward {
		return false
	}

	if n.Value == o.Value {
	} else if n.Value == nil || o.Value == nil {
		return false
	} else if !n.Value.Equal(value.Ref(o.Value)) {
		return false
	}

	return true
}

func (n *YieldExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("yield")
	if n.Forward {
		buff.WriteRune('*')
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

func (*YieldExpressionNode) IsStatic() bool {
	return false
}

// Create a new `yield` expression node eg. `yield`, `yield true`, `yield* foo()`
func NewYieldExpressionNode(loc *position.Location, forward bool, val ExpressionNode) *YieldExpressionNode {
	return &YieldExpressionNode{
		NodeBase: NodeBase{loc: loc},
		Forward:  forward,
		Value:    val,
	}
}

func (*YieldExpressionNode) Class() *value.Class {
	return value.YieldExpressionNodeClass
}

func (*YieldExpressionNode) DirectClass() *value.Class {
	return value.YieldExpressionNodeClass
}

func (n *YieldExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::YieldExpressionNode{\n  loc: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  value: ")
	if n.Value == nil {
		buff.WriteString("nil")
	} else {
		indent.IndentStringFromSecondLine(&buff, n.Value.Inspect(), 1)
	}

	buff.WriteString("\n}")

	return buff.String()
}

func (n *YieldExpressionNode) Error() string {
	return n.Inspect()
}
