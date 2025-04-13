package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a `return` expression eg. `return`, `return true`
type ReturnExpressionNode struct {
	NodeBase
	Value ExpressionNode
}

func (n *ReturnExpressionNode) Splice(loc *position.Location, args *[]Node) Node {
	var val ExpressionNode
	if n.Value != nil {
		val = n.Value.Splice(loc, args).(ExpressionNode)
	}

	return &ReturnExpressionNode{
		NodeBase: n.NodeBase,
		Value:    val,
	}
}

func (n *ReturnExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ReturnExpressionNode)
	if !ok {
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

func (n *ReturnExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("return")

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

func (*ReturnExpressionNode) IsStatic() bool {
	return false
}

func (*ReturnExpressionNode) Type(*types.GlobalEnvironment) types.Type {
	return types.Never{}
}

// Create a new `return` expression node eg. `return`, `return true`
func NewReturnExpressionNode(loc *position.Location, val ExpressionNode) *ReturnExpressionNode {
	return &ReturnExpressionNode{
		NodeBase: NodeBase{loc: loc},
		Value:    val,
	}
}

func (*ReturnExpressionNode) Class() *value.Class {
	return value.ReturnExpressionNodeClass
}

func (*ReturnExpressionNode) DirectClass() *value.Class {
	return value.ReturnExpressionNodeClass
}

func (n *ReturnExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ReturnExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  value: ")
	if n.Value == nil {
		buff.WriteString("nil")
	} else {
		indent.IndentStringFromSecondLine(&buff, n.Value.Inspect(), 1)
	}

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ReturnExpressionNode) Error() string {
	return n.Inspect()
}
