package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
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

func (p *ProgramNode) Inspect() string {
	var buff strings.Builder
	fmt.Fprintf(&buff, "Std::AST::ProgramNode{&: %p, body: %%[", p)

	for i, stmt := range p.Body {
		if i != 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(stmt.Inspect())
	}

	buff.WriteString("]}")
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
