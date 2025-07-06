package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents an `await` expression eg. `await foo()`
type AwaitExpressionNode struct {
	TypedNodeBase
	Value ExpressionNode
	Sync  bool
}

func (n *AwaitExpressionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &AwaitExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Sync:          n.Sync,
		Value:         n.Value.splice(loc, args, unquote).(ExpressionNode),
	}
}

func (n *AwaitExpressionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::AwaitExpressionNode", env)
}

func (n *AwaitExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
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

func (*AwaitExpressionNode) IsStatic() bool {
	return false
}

// Create a new `await` expression node eg. `await foo()`
func NewAwaitExpressionNode(loc *position.Location, val ExpressionNode, sync bool) *AwaitExpressionNode {
	return &AwaitExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
		Sync:          sync,
	}
}

func (n *AwaitExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*AwaitExpressionNode)
	if !ok {
		return false
	}

	if n.Sync != o.Sync || !n.loc.Equal(o.loc) {
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

	if n.Sync {
		buff.WriteString("await_sync")
	} else {
		buff.WriteString("await ")
	}

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

	fmt.Fprintf(&buff, ",\n  sync: %t", n.Sync)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *AwaitExpressionNode) Error() string {
	return n.Inspect()
}
