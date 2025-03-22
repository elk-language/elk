package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

type ProgramState uint8

const (
	UNCHECKED ProgramState = iota
	CHECKING_NAMESPACES
	CHECKED_NAMESPACES

	CHECKING_METHODS
	CHECKED_METHODS

	CHECKING_EXPRESSIONS
	CHECKED_EXPRESSIONS
)

// Represents a single Elk program (usually a single file).
type ProgramNode struct {
	NodeBase
	Body        []StatementNode
	ImportPaths []string
	State       ProgramState
}

func (*ProgramNode) IsStatic() bool {
	return false
}

func (*ProgramNode) Class() *value.Class {
	return value.ProgramNodeClass
}

func (*ProgramNode) DirectClass() *value.Class {
	return value.ProgramNodeClass
}

func (n *ProgramNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ProgramNode{\n  &: %p", n)

	buff.WriteString(",\n  body: %%[\n")
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

func (p *ProgramNode) Error() string {
	return p.Inspect()
}

// Create a new program node.
func NewProgramNode(span *position.Span, body []StatementNode) *ProgramNode {
	return &ProgramNode{
		NodeBase: NodeBase{span: span},
		Body:     body,
	}
}
