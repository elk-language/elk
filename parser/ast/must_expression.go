package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a `must` expression eg. `must foo()`
type MustExpressionNode struct {
	TypedNodeBase
	Value ExpressionNode
}

func (n *MustExpressionNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &MustExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Value:         n.Value,
	}
}

func (n *MustExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Value.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (n *MustExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*MustExpressionNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.Value.Equal(value.Ref(o.Value))
}

func (n *MustExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("must ")
	buff.WriteString(n.Value.String())

	return buff.String()
}

func (*MustExpressionNode) IsStatic() bool {
	return false
}

// Create a new `must` expression node eg. `must foo()`
func NewMustExpressionNode(loc *position.Location, val ExpressionNode) *MustExpressionNode {
	return &MustExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
	}
}

func (*MustExpressionNode) Class() *value.Class {
	return value.MustExpressionNodeClass
}

func (*MustExpressionNode) DirectClass() *value.Class {
	return value.MustExpressionNodeClass
}

func (n *MustExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::MustExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  value: ")
	indent.IndentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *MustExpressionNode) Error() string {
	return n.Inspect()
}
