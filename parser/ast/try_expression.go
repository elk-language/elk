package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a `try` expression eg. `try foo()`
type TryExpressionNode struct {
	TypedNodeBase
	Value ExpressionNode
}

func (n *TryExpressionNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &TryExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Value:         n.Value.Splice(loc, args, unquote).(ExpressionNode),
	}
}

func (n *TryExpressionNode) Traverse(yield func(Node) bool) bool {
	if n.Value != nil {
		if n.Value.Traverse(yield) {
			return false
		}
	}
	return yield(n)
}

func (n *TryExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*TryExpressionNode)
	if !ok {
		return false
	}

	return n.Value.Equal(value.Ref(o.Value)) &&
		n.loc.Equal(o.loc)
}

func (n *TryExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("try ")

	if n.Value != nil {
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

func (*TryExpressionNode) IsStatic() bool {
	return false
}

// Create a new `try` expression node eg. `try foo()`
func NewTryExpressionNode(loc *position.Location, val ExpressionNode) *TryExpressionNode {
	return &TryExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
	}
}

func (*TryExpressionNode) Class() *value.Class {
	return value.TryExpressionNodeClass
}

func (*TryExpressionNode) DirectClass() *value.Class {
	return value.TryExpressionNodeClass
}

func (n *TryExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::TryExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  value: ")
	indent.IndentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *TryExpressionNode) Error() string {
	return n.Inspect()
}
