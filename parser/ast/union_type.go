package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Union type eg. `String | Int | Float`
type UnionTypeNode struct {
	TypedNodeBase
	Elements []TypeNode
}

func (n *UnionTypeNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UnionTypeNode)
	if !ok {
		return false
	}

	if len(n.Elements) != len(o.Elements) {
		return false
	}

	for i, element := range n.Elements {
		if !element.Equal(value.Ref(o.Elements[i])) {
			return false
		}
	}

	return n.span.Equal(o.span)
}

func (n *UnionTypeNode) String() string {
	var buff strings.Builder

	for i, element := range n.Elements {
		if i > 0 {
			buff.WriteString(" | ")
		}
		buff.WriteString(element.String())
	}

	return buff.String()
}

func (*UnionTypeNode) IsStatic() bool {
	return false
}

// Create a new binary type expression node eg. `String | Int`
func NewUnionTypeNode(span *position.Span, elements []TypeNode) *UnionTypeNode {
	return &UnionTypeNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

func (*UnionTypeNode) Class() *value.Class {
	return value.UnionTypeNodeClass
}

func (*UnionTypeNode) DirectClass() *value.Class {
	return value.UnionTypeNodeClass
}

func (n *UnionTypeNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::UnionTypeNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  elements: %[\n")
	for i, stmt := range n.Elements {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *UnionTypeNode) Error() string {
	return n.Inspect()
}
