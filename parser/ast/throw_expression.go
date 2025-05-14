package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a `throw` expression eg. `throw ArgumentError("foo")`
type ThrowExpressionNode struct {
	NodeBase
	Unchecked bool
	Value     ExpressionNode
}

func (n *ThrowExpressionNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	var val ExpressionNode
	if n.Value != nil {
		val = n.Value.Splice(loc, args, unquote).(ExpressionNode)
	}

	return &ThrowExpressionNode{
		NodeBase:  NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Unchecked: n.Unchecked,
		Value:     val,
	}
}

func (n *ThrowExpressionNode) Traverse(yield func(Node) bool) bool {
	if n.Value != nil {
		if n.Value.Traverse(yield) {
			return false
		}
	}
	return yield(n)
}

func (n *ThrowExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ThrowExpressionNode)
	if !ok {
		return false
	}

	if n.Value == o.Value {
	} else if n.Value == nil || o.Value == nil {
		return false
	} else if !n.Value.Equal(value.Ref(o.Value)) {
		return false
	}

	return n.Unchecked == o.Unchecked &&
		n.loc.Equal(o.loc)
}

func (n *ThrowExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("throw ")

	if n.Unchecked {
		buff.WriteString("unchecked ")
	}

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

func (*ThrowExpressionNode) Type(*types.GlobalEnvironment) types.Type {
	return types.Never{}
}

func (*ThrowExpressionNode) IsStatic() bool {
	return false
}

// Create a new `throw` expression node eg. `throw ArgumentError("foo")`
func NewThrowExpressionNode(loc *position.Location, unchecked bool, val ExpressionNode) *ThrowExpressionNode {
	return &ThrowExpressionNode{
		NodeBase:  NodeBase{loc: loc},
		Unchecked: unchecked,
		Value:     val,
	}
}

func (*ThrowExpressionNode) Class() *value.Class {
	return value.ThrowExpressionNodeClass
}

func (*ThrowExpressionNode) DirectClass() *value.Class {
	return value.ThrowExpressionNodeClass
}

func (n *ThrowExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ThrowExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	fmt.Fprintf(&buff, ",\n  unchecked: %t", n.Unchecked)

	buff.WriteString(",\n  value: ")
	if n.Value == nil {
		buff.WriteString("nil")
	} else {
		indent.IndentStringFromSecondLine(&buff, n.Value.Inspect(), 1)
	}

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ThrowExpressionNode) Error() string {
	return n.Inspect()
}
