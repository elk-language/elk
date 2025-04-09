package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Nodes that implement this interface represent
// named arguments in method calls.
type NamedArgumentNode interface {
	Node
	namedArgumentNode()
}

func (*InvalidNode) namedArgumentNode()               {}
func (*NamedCallArgumentNode) namedArgumentNode()     {}
func (*DoubleSplatExpressionNode) namedArgumentNode() {}

// Represents a named argument in a function call eg. `foo: 123`
type NamedCallArgumentNode struct {
	NodeBase
	Name  string
	Value ExpressionNode
}

func (n *NamedCallArgumentNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*NamedCallArgumentNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.Name == o.Name &&
		n.Value.Equal(value.Ref(o.Value))
}

func (n *NamedCallArgumentNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.Name)
	buff.WriteString(": ")
	buff.WriteString(n.Value.String())

	return buff.String()
}

func (*NamedCallArgumentNode) IsStatic() bool {
	return false
}

func (*NamedCallArgumentNode) Class() *value.Class {
	return value.NamedCallArgumentNodeClass
}

func (*NamedCallArgumentNode) DirectClass() *value.Class {
	return value.NamedCallArgumentNodeClass
}

func (n *NamedCallArgumentNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::NamedCallArgumentNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  name: ")
	indent.IndentStringFromSecondLine(&buff, value.String(n.Name).Inspect(), 1)

	buff.WriteString(",\n  value: ")
	indent.IndentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *NamedCallArgumentNode) Error() string {
	return n.Inspect()
}

// Create a named argument node eg. `foo: 123`
func NewNamedCallArgumentNode(loc *position.Location, name string, val ExpressionNode) *NamedCallArgumentNode {
	return &NamedCallArgumentNode{
		NodeBase: NodeBase{loc: loc},
		Name:     name,
		Value:    val,
	}
}
