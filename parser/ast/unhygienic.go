package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents an unhygienic node
type UnhygienicNode struct {
	TypedNodeBase
	Node Node
}

func (n *UnhygienicNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &UnhygienicNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Node:          n.Node.splice(loc, args, unquote),
	}
}

func (n *UnhygienicNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::UnhygienicNode", env)
}

func (n *UnhygienicNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.traverse(n.Node, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

// Check if this node equals another node.
func (n *UnhygienicNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UnhygienicNode)
	if !ok {
		return false
	}

	if !n.Node.Equal(value.Ref(o.Node)) {
		return false
	}

	return n.loc.Equal(o.loc) && n.Node.Equal(value.Ref(o.Node))
}

// Return a string representation of the node.
func (n *UnhygienicNode) String() string {
	var buff strings.Builder

	buff.WriteString("unhygienic(")
	buff.WriteString(n.Node.String())
	buff.WriteRune(')')

	return buff.String()
}

func (n *UnhygienicNode) IsStatic() bool {
	return n.IsStatic()
}

func (*UnhygienicNode) Class() *value.Class {
	return value.UnhygienicNodeClass
}

func (*UnhygienicNode) DirectClass() *value.Class {
	return value.UnhygienicNodeClass
}

func (n *UnhygienicNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::UnhygienicNode{\n  location: %s", (*value.Location)(n.loc).Inspect())
	fmt.Fprintf(&buff, ",\n  node: %s", n.Node.Inspect())
	buff.WriteString("\n}")

	return buff.String()
}

func (n *UnhygienicNode) Error() string {
	return n.Inspect()
}

// Create a new unhygienic node.
// It enables macro-generated code to access
// local variables from outer scopes.
func NewUnhygienicNode(loc *position.Location, node Node) *UnhygienicNode {
	return &UnhygienicNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Node:          node,
	}
}
