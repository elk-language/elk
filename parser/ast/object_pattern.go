package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents an Object pattern eg. `Foo(foo: 5, bar: a, c)`
type ObjectPatternNode struct {
	TypedNodeBase
	ObjectType ComplexConstantNode
	Attributes []PatternNode
}

func (m *ObjectPatternNode) IsStatic() bool {
	return false
}

// Create an Object pattern node eg. `Foo(foo: 5, bar: a, c)`
func NewObjectPatternNode(span *position.Span, objectType ComplexConstantNode, attrs []PatternNode) *ObjectPatternNode {
	return &ObjectPatternNode{
		TypedNodeBase: TypedNodeBase{span: span},
		ObjectType:    objectType,
		Attributes:    attrs,
	}
}

func (*ObjectPatternNode) Class() *value.Class {
	return value.ObjectPatternNodeClass
}

func (*ObjectPatternNode) DirectClass() *value.Class {
	return value.ObjectPatternNodeClass
}

func (n *ObjectPatternNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ObjectPatternNode{\n  &: %p", n)

	buff.WriteString(",\n  object_type: ")
	indent.IndentStringFromSecondLine(&buff, n.ObjectType.Inspect(), 1)

	buff.WriteString(",\n  attributes: %%[\n")
	for i, element := range n.Attributes {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ObjectPatternNode) Error() string {
	return n.Inspect()
}
