package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a Record pattern eg. `%{ foo: 5, bar: a, 5 => >= 10 }`
type RecordPatternNode struct {
	NodeBase
	Elements []PatternNode
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
		indentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *RecordPatternNode) Error() string {
	return n.Inspect()
}
