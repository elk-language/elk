package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a type expression `type String?`
type TypeExpressionNode struct {
	NodeBase
	TypeNode TypeNode
}

func (n *TypeExpressionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &TypeExpressionNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		TypeNode: n.TypeNode.splice(loc, args, unquote).(TypeNode),
	}
}

func (n *TypeExpressionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::TypeExpressionNode", env)
}

func (n *TypeExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.TypeNode.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (n *TypeExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*TypeExpressionNode)
	if !ok {
		return false
	}

	return n.TypeNode.Equal(value.Ref(o.TypeNode)) &&
		n.loc.Equal(o.loc)
}

func (n *TypeExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("type ")
	buff.WriteString(n.TypeNode.String())

	return buff.String()
}

func (*TypeExpressionNode) IsStatic() bool {
	return false
}

func (*TypeExpressionNode) Class() *value.Class {
	return value.TypeExpressionNodeClass
}

func (*TypeExpressionNode) DirectClass() *value.Class {
	return value.TypeExpressionNodeClass
}

func (n *TypeExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::TypeExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  type_node: ")
	indent.IndentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *TypeExpressionNode) Error() string {
	return n.Inspect()
}

// Create a new type expression `type String?`
func NewTypeExpressionNode(loc *position.Location, typeNode TypeNode) *TypeExpressionNode {
	return &TypeExpressionNode{
		NodeBase: NodeBase{loc: loc},
		TypeNode: typeNode,
	}
}
