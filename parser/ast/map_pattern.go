package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a Map pattern eg. `{ foo: 5, bar: a, 5 => >= 10 }`
type MapPatternNode struct {
	TypedNodeBase
	Elements []PatternNode
}

func (n *MapPatternNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*MapPatternNode)
	if !ok {
		return false
	}

	if len(n.Elements) != len(o.Elements) {
		return false
	}

	for i, elem := range n.Elements {
		if !elem.Equal(value.Ref(o.Elements[i])) {
			return false
		}
	}

	return n.span.Equal(o.span)
}

func (n *MapPatternNode) String() string {
	var buff strings.Builder

	buff.WriteString("{")

	for i, elem := range n.Elements {
		if i != 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(elem.String())
	}

	buff.WriteString("}")

	return buff.String()
}

func (m *MapPatternNode) IsStatic() bool {
	return false
}

// Create a Map pattern node eg. `{ foo: 5, bar: a, 5 => >= 10 }`
func NewMapPatternNode(span *position.Span, elements []PatternNode) *MapPatternNode {
	return &MapPatternNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

// Same as [NewMapPatternNode] but returns an interface
func NewMapPatternNodeI(span *position.Span, elements []PatternNode) PatternNode {
	return NewMapPatternNode(span, elements)
}

func (*MapPatternNode) Class() *value.Class {
	return value.MapPatternNodeClass
}

func (*MapPatternNode) DirectClass() *value.Class {
	return value.MapPatternNodeClass
}

func (n *MapPatternNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::MapPatternNode{\n  span: %s", (*value.Span)(n.span).Inspect())

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

func (n *MapPatternNode) Error() string {
	return n.Inspect()
}
