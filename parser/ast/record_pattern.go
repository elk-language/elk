package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a Record pattern eg. `%{ foo: 5, bar: a, 5 => >= 10 }`
type RecordPatternNode struct {
	NodeBase
	Elements []PatternNode
}

func (n *RecordPatternNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*RecordPatternNode)
	if !ok {
		return false
	}

	if len(n.Elements) != len(o.Elements) ||
		!n.span.Equal(o.span) {
		return false
	}

	for i, element := range n.Elements {
		if !element.Equal(value.Ref(o.Elements[i])) {
			return false
		}
	}

	return true
}

func (n *RecordPatternNode) String() string {
	var buff strings.Builder

	buff.WriteString("%{")

	for i, element := range n.Elements {
		if i != 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(element.String())
	}

	buff.WriteString("}")

	return buff.String()
}

func (m *RecordPatternNode) IsStatic() bool {
	return false
}

// Create a Record pattern node eg. `%{ foo: 5, bar: a, 5 => >= 10 }`
func NewRecordPatternNode(span *position.Span, elements []PatternNode) *RecordPatternNode {
	return &RecordPatternNode{
		NodeBase: NodeBase{span: span},
		Elements: elements,
	}
}

// Same as [NewRecordPatternNode] but returns an interface
func NewRecordPatternNodeI(span *position.Span, elements []PatternNode) PatternNode {
	return NewRecordPatternNode(span, elements)
}

func (*RecordPatternNode) Class() *value.Class {
	return value.RecordPatternNodeClass
}

func (*RecordPatternNode) DirectClass() *value.Class {
	return value.RecordPatternNodeClass
}

func (n *RecordPatternNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::RecordPatternNode{\n  &: %p", n)

	buff.WriteString(",\n  elements: %%[\n")
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

func (n *RecordPatternNode) Error() string {
	return n.Inspect()
}
