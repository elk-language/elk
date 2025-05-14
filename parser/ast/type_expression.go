package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a type expression `type String?`
type TypeExpressionNode struct {
	NodeBase
	TypeNode TypeNode
}

func (n *TypeExpressionNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &TypeExpressionNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		TypeNode: n.TypeNode.Splice(loc, args, unquote).(TypeNode),
	}
}

func (n *TypeExpressionNode) Traverse(yield func(Node) bool) bool {
	if n.TypeNode.Traverse(yield) {
		return false
	}
	return yield(n)
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
