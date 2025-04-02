package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a boundary for expanded macros
type MacroBoundaryNode struct {
	TypedNodeBase
	Body []StatementNode
}

// Check if this node equals another node.
func (n *MacroBoundaryNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*MacroBoundaryNode)
	if !ok {
		return false
	}

	if len(n.Body) != len(o.Body) {
		return false
	}

	for i, stmt := range n.Body {
		if !stmt.Equal(value.Ref(o.Body[i])) {
			return false
		}
	}

	return n.span.Equal(o.span)
}

// Return a string representation of the node.
func (n *MacroBoundaryNode) String() string {
	var buff strings.Builder

	buff.WriteString("macro do\n")

	for _, stmt := range n.Body {
		indent.IndentString(&buff, stmt.String(), 1)
		buff.WriteString("\n")
	}

	buff.WriteString("end")

	return buff.String()
}

func (*MacroBoundaryNode) IsStatic() bool {
	return false
}

func (*MacroBoundaryNode) Class() *value.Class {
	return value.MacroBoundaryNodeClass
}

func (*MacroBoundaryNode) DirectClass() *value.Class {
	return value.MacroBoundaryNodeClass
}

func (n *MacroBoundaryNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::MacroBoundaryNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  body: %[\n")
	for i, stmt := range n.Body {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, stmt.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *MacroBoundaryNode) Error() string {
	return n.Inspect()
}

// Create a new macro boundary node eg.
//
//	macro do
//		print("awesome!")
//	end
func NewMacroBoundaryNode(span *position.Span, body []StatementNode) *MacroBoundaryNode {
	return &MacroBoundaryNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Body:          body,
	}
}
