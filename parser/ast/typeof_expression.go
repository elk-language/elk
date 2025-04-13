package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a `typeof` expression eg. `typeof foo()`
type TypeofExpressionNode struct {
	TypedNodeBase
	Value ExpressionNode
}

func (n *TypeofExpressionNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &TypeofExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Value:         n.Value.Splice(loc, args, unquote).(ExpressionNode),
	}
}

func (n *TypeofExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*TypeofExpressionNode)
	if !ok {
		return false
	}

	return n.Value.Equal(value.Ref(o.Value)) &&
		n.loc.Equal(o.loc)
}

func (n *TypeofExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("typeof ")
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

func (*TypeofExpressionNode) IsStatic() bool {
	return false
}

// Create a new `typeof` expression node eg. `typeof foo()`
func NewTypeofExpressionNode(loc *position.Location, val ExpressionNode) *TypeofExpressionNode {
	return &TypeofExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
	}
}

func (*TypeofExpressionNode) Class() *value.Class {
	return value.TypeofExpressionNodeClass
}

func (*TypeofExpressionNode) DirectClass() *value.Class {
	return value.TypeofExpressionNodeClass
}

func (n *TypeofExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::TypeofExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  value: ")
	indent.IndentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *TypeofExpressionNode) Error() string {
	return n.Inspect()
}
