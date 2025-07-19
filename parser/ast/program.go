package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

type ProgramState uint8

const (
	UNCHECKED ProgramState = iota
	CHECKING_NAMESPACES
	CHECKED_NAMESPACES

	EXPANDING_TOP_LEVEL_MACROS
	EXPANDED_TOP_LEVEL_MACROS

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

func (n *ProgramNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &ProgramNode{
		NodeBase:    NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Body:        SpliceSlice(n.Body, loc, args, unquote),
		ImportPaths: n.ImportPaths,
		State:       n.State,
	}
}

func (n *ProgramNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::ProgramNode", env)
}

func (n *ProgramNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	for _, stmt := range n.Body {
		if stmt.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
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
