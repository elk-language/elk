package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents an as type downcast eg. `foo as String`
type AsExpressionNode struct {
	TypedNodeBase
	Value       ExpressionNode
	RuntimeType ComplexConstantNode
}

func (n *AsExpressionNode) Splice(loc *position.Location, args *[]Node) Node {
	return &AsExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: getLoc(loc, n.loc), typ: n.typ},
		Value:         n.Value.Splice(loc, args).(ExpressionNode),
		RuntimeType:   n.RuntimeType.Splice(loc, args).(ComplexConstantNode),
	}
}

func (*AsExpressionNode) IsStatic() bool {
	return false
}

func (*AsExpressionNode) Class() *value.Class {
	return value.PublicIdentifierNodeClass
}

func (*AsExpressionNode) DirectClass() *value.Class {
	return value.PublicIdentifierNodeClass
}

func (n *AsExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*AsExpressionNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.Value.Equal(value.Ref(o.Value)) &&
		n.RuntimeType.Equal(value.Ref(o.RuntimeType))
}

func (n *AsExpressionNode) String() string {
	var buff strings.Builder

	parens := ExpressionPrecedence(n) > ExpressionPrecedence(n.Value)
	if parens {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Value.String())
	if parens {
		buff.WriteRune(')')
	}

	buff.WriteString(" as ")

	buff.WriteString(n.RuntimeType.String())

	return buff.String()
}

func (n *AsExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::AsExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  value: ")
	indent.IndentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

	buff.WriteString(",\n  runtime_type: ")
	indent.IndentStringFromSecondLine(&buff, n.RuntimeType.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *AsExpressionNode) Error() string {
	return n.Inspect()
}

// Create a new private constant node eg. `_Foo`.
func NewAsExpressionNode(loc *position.Location, val ExpressionNode, runtimeType ComplexConstantNode) *AsExpressionNode {
	return &AsExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
		RuntimeType:   runtimeType,
	}
}
