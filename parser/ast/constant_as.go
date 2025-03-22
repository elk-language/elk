package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a constant with as in using declarations
// eg. `Foo::Bar as Bar`.
type ConstantAsNode struct {
	NodeBase
	Constant ComplexConstantNode
	AsName   string
}

func (*ConstantAsNode) IsStatic() bool {
	return false
}

func (*ConstantAsNode) Class() *value.Class {
	return value.ConstantAsNodeClass
}

func (*ConstantAsNode) DirectClass() *value.Class {
	return value.ConstantAsNodeClass
}

func (n *ConstantAsNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ConstantAsNode{\n  &: %p", n)

	buff.WriteString(",\n  constant: ")
	indent.IndentStringFromSecondLine(&buff, n.Constant.Inspect(), 1)

	buff.WriteString(",\n  as_name: ")
	buff.WriteString(n.AsName)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ConstantAsNode) Error() string {
	return n.Inspect()
}

// Create a new identifier with as eg. `Foo::Bar as Bar`.
func NewConstantAsNode(span *position.Span, constant ComplexConstantNode, as string) *ConstantAsNode {
	return &ConstantAsNode{
		NodeBase: NodeBase{span: span},
		Constant: constant,
		AsName:   as,
	}
}
