package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a `continue` expression eg. `continue`, `continue "foo"`
type ContinueExpressionNode struct {
	NodeBase
	Label string
	Value ExpressionNode
}

// Check if this node equals another node.
func (n *ContinueExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ContinueExpressionNode)
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

	return n.span.Equal(o.span)
}

// Return a string representation of the node.
func (n *ContinueExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("continue")

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

func (*ContinueExpressionNode) Type(*types.GlobalEnvironment) types.Type {
	return types.Never{}
}

func (*ContinueExpressionNode) IsStatic() bool {
	return false
}

// Create a new `continue` expression node eg. `continue`, `continue "foo"`
func NewContinueExpressionNode(span *position.Span, label string, val ExpressionNode) *ContinueExpressionNode {
	return &ContinueExpressionNode{
		NodeBase: NodeBase{span: span},
		Label:    label,
		Value:    val,
	}
}

func (*ContinueExpressionNode) Class() *value.Class {
	return value.ContinueExpressionNodeClass
}

func (*ContinueExpressionNode) DirectClass() *value.Class {
	return value.ContinueExpressionNodeClass
}

func (n *ContinueExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ContinueExpressionNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  value: ")
	indent.IndentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ContinueExpressionNode) Error() string {
	return n.Inspect()
}
