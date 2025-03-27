package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents an include expression eg. `include Enumerable[V]`
type IncludeExpressionNode struct {
	TypedNodeBase
	Constants []ComplexConstantNode
}

// Check if this node equals another node.
func (n *IncludeExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*IncludeExpressionNode)
	if !ok {
		return false
	}

	if len(n.Constants) != len(o.Constants) {
		return false
	}

	for i, constant := range n.Constants {
		if !constant.Equal(value.Ref(o.Constants[i])) {
			return false
		}
	}

	return n.span.Equal(o.span)
}

// Return a string representation of the node.
func (n *IncludeExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("include ")

	for i, constant := range n.Constants {
		if i != 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(constant.String())
	}

	return buff.String()
}

func (*IncludeExpressionNode) SkipTypechecking() bool {
	return false
}

func (*IncludeExpressionNode) IsStatic() bool {
	return false
}

// Create an include expression node eg. `include Enumerable[V]`
func NewIncludeExpressionNode(span *position.Span, consts []ComplexConstantNode) *IncludeExpressionNode {
	return &IncludeExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Constants:     consts,
	}
}

func (*IncludeExpressionNode) Class() *value.Class {
	return value.UsingAllEntryNodeClass
}

func (*IncludeExpressionNode) DirectClass() *value.Class {
	return value.UsingAllEntryNodeClass
}

func (n *IncludeExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::IncludeExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  constants: %%[\n")
	for i, element := range n.Constants {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *IncludeExpressionNode) Error() string {
	return n.Inspect()
}
