package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a boundary for expanded macros
type MacroBoundaryNode struct {
	TypedNodeBase
	Name string
	Body []StatementNode
}

func (n *MacroBoundaryNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &MacroBoundaryNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Name:          n.Name,
		Body:          SpliceSlice(n.Body, loc, args, unquote),
	}
}

func (n *MacroBoundaryNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::MacroBoundaryNode", env)
}

func (n *MacroBoundaryNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
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

	return n.loc.Equal(o.loc) && n.Name == o.Name
}

// Return a string representation of the node.
func (n *MacroBoundaryNode) String() string {
	var buff strings.Builder

	buff.WriteString("do macro ")
	if n.Name != "" {
		fmt.Fprintf(&buff, "'%s' ", n.Name)
	}
	buff.WriteString("\n")

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

	fmt.Fprintf(&buff, "Std::Elk::AST::MacroBoundaryNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	if n.Name != "" {
		fmt.Fprintf(&buff, ",\n  name: %s", value.String(n.Name).Inspect())
	}
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
//	do macro
//		print("awesome!")
//	end
func NewMacroBoundaryNode(loc *position.Location, body []StatementNode, name string) *MacroBoundaryNode {
	return &MacroBoundaryNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Name:          name,
		Body:          body,
	}
}
