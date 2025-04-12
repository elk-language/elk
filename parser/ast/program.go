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

func (n *ProgramNode) Splice(loc *position.Location, args *[]Node) Node {
	return &ProgramNode{
		NodeBase:    n.NodeBase,
		Body:        SpliceSlice(n.Body, loc, args),
		ImportPaths: n.ImportPaths,
		State:       n.State,
	}
}

func (n *ProgramNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ProgramNode)
	if !ok {
		return false
	}

	if len(n.Body) != len(o.Body) {
		return false
	}

	for i, stmt := range n.Body {
		if !value.Equal(value.Ref(stmt), value.Ref(o.Body[i])) {
			return false
		}
	}

	return n.loc.Equal(o.loc)
}

func (n *ProgramNode) String() string {
	var buff strings.Builder
	for _, stmt := range n.Body {
		buff.WriteString(stmt.String())
		buff.WriteRune('\n')
	}
	return buff.String()
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

	fmt.Fprintf(&buff, "Std::Elk::AST::ProgramNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

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

func (p *ProgramNode) Error() string {
	return p.Inspect()
}

// Create a new program node.
func NewProgramNode(loc *position.Location, body []StatementNode) *ProgramNode {
	return &ProgramNode{
		NodeBase: NodeBase{loc: loc},
		Body:     body,
	}
}
