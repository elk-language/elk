package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents an `await` expression eg. `await foo()`
type AwaitExpressionNode struct {
	TypedNodeBase
	Value ExpressionNode
}

func (n *AwaitExpressionNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &AwaitExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Value:         n.Value.Splice(loc, args, unquote).(ExpressionNode),
	}
}

func (*AwaitExpressionNode) IsStatic() bool {
	return false
}

// Create a new `await` expression node eg. `await foo()`
func NewAwaitExpressionNode(loc *position.Location, val ExpressionNode) *AwaitExpressionNode {
	return &AwaitExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
	}
}
func (n *AwaitExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*AwaitExpressionNode)
	if !ok {
		return false
	}

	if !n.loc.Equal(o.loc) {
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

func (n *AwaitExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("await ")

	parens := ExpressionPrecedence(n) > ExpressionPrecedence(n.Value)
	if parens {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Value.String())
	if parens {
		buff.WriteRune(')')
	}

	return buff.String()
}

func (*AwaitExpressionNode) Class() *value.Class {
	return value.AwaitExpressionNodeClass
}

func (*AwaitExpressionNode) DirectClass() *value.Class {
	return value.AwaitExpressionNodeClass
}

func (n *AwaitExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::AwaitExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  value: ")
	indent.IndentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *AwaitExpressionNode) Error() string {
	return n.Inspect()
}
