package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Union type eg. `String & Int & Float`
type IntersectionTypeNode struct {
	TypedNodeBase
	Elements []TypeNode
}

func (*IntersectionTypeNode) IsStatic() bool {
	return false
}

// Create a new binary type expression node eg. `String & Int`
func NewIntersectionTypeNode(span *position.Span, elements []TypeNode) *IntersectionTypeNode {
	return &IntersectionTypeNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

func (*IntersectionTypeNode) Class() *value.Class {
	return value.IntersectionTypeNodeClass
}

func (*IntersectionTypeNode) DirectClass() *value.Class {
	return value.IntersectionTypeNodeClass
}

func (n *IntersectionTypeNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::IntersectionTypeNode{\n  &: %p", n)

	buff.WriteString(",\n  elements: %%[\n")
	for i, stmt := range n.Elements {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *IntersectionTypeNode) Error() string {
	return n.Inspect()
}
