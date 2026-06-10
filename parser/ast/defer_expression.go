package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a `defer` expression eg. `defer foo()`
type DeferExpressionNode struct {
	NodeBase
	Expression ExpressionNode
}

func (n *DeferExpressionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &DeferExpressionNode{
		NodeBase:   NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Expression: n.Expression,
	}
}

func (*DeferExpressionNode) SetType(types.Type) {}

func (*DeferExpressionNode) Type(globalEnv *types.GlobalEnvironment) types.Type {
	return types.Nil{}
}

func (n *DeferExpressionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::DeferExpressionNode", env)
}

func (n *DeferExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Expression.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (n *DeferExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*DeferExpressionNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.Expression.Equal(value.Ref(o.Expression))
}

func (n *DeferExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("defer ")
	buff.WriteString(n.Expression.String())

	return buff.String()
}

func (*DeferExpressionNode) IsStatic() bool {
	return false
}

// Create a new `defer` expression node eg. `defer foo()`
func NewDeferExpressionNode(loc *position.Location, expr ExpressionNode) *DeferExpressionNode {
	return &DeferExpressionNode{
		NodeBase:   NodeBase{loc: loc},
		Expression: expr,
	}
}

func (*DeferExpressionNode) Class() *value.Class {
	return value.DeferExpressionNodeClass
}

func (*DeferExpressionNode) DirectClass() *value.Class {
	return value.DeferExpressionNodeClass
}

func (n *DeferExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::DeferExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  expression: ")
	indent.IndentStringFromSecondLine(&buff, n.Expression.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *DeferExpressionNode) ToValue() value.Value {
	return value.Ref(n)
}

func (n *DeferExpressionNode) Error() string {
	return n.Inspect()
}
