package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents an instance variable declaration eg. `var @foo: String`
type InstanceVariableDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Name     string   // name of the variable
	TypeNode TypeNode // type of the variable
}

func (n *InstanceVariableDeclarationNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.TypeNode != nil {
		if n.TypeNode.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

func (n *InstanceVariableDeclarationNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	var typeNode ComplexConstantNode
	if n.TypeNode != nil {
		typeNode = n.TypeNode.splice(loc, args, unquote).(ComplexConstantNode)
	}

	return &InstanceVariableDeclarationNode{
		TypedNodeBase:          TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		DocCommentableNodeBase: n.DocCommentableNodeBase,
		Name:                   n.Name,
		TypeNode:               typeNode,
	}
}

func (n *InstanceVariableDeclarationNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::InstanceVariableDeclarationNode", env)
}

func (n *InstanceVariableDeclarationNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*InstanceVariableDeclarationNode)
	if !ok {
		return false
	}

	if n.comment != o.comment ||
		n.Name != o.Name ||
		!n.loc.Equal(o.loc) {
		return false
	}

	if n.TypeNode == o.TypeNode {
	} else if n.TypeNode == nil || o.TypeNode == nil {
		return false
	} else if !n.TypeNode.Equal(value.Ref(o.TypeNode)) {
		return false
	}

	return true
}

func (n *InstanceVariableDeclarationNode) String() string {
	var buff strings.Builder

	doc := n.DocComment()
	if len(doc) > 0 {
		buff.WriteString("##[\n")
		indent.IndentString(&buff, doc, 1)
		buff.WriteString("\n]##\n")
	}

	buff.WriteString("var @")
	buff.WriteString(n.Name)

	if n.TypeNode != nil {
		buff.WriteString(": ")
		buff.WriteString(n.TypeNode.String())
	}

	return buff.String()
}

func (*InstanceVariableDeclarationNode) IsStatic() bool {
	return false
}

func (*InstanceVariableDeclarationNode) Class() *value.Class {
	return value.InstanceVariableDeclarationNodeClass
}

func (*InstanceVariableDeclarationNode) DirectClass() *value.Class {
	return value.InstanceVariableDeclarationNodeClass
}

func (n *InstanceVariableDeclarationNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::InstanceVariableDeclarationNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  name: ")
	buff.WriteString(n.Name)

	buff.WriteString(",\n  type_node: ")
	if n.TypeNode == nil {
		buff.WriteString("nil")
	} else {
		indent.IndentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)
	}

	buff.WriteString("\n}")

	return buff.String()
}

func (p *InstanceVariableDeclarationNode) Error() string {
	return p.Inspect()
}

// Create a new instance variable declaration node eg. `var @foo: String`
func NewInstanceVariableDeclarationNode(loc *position.Location, docComment string, name string, typ TypeNode) *InstanceVariableDeclarationNode {
	return &InstanceVariableDeclarationNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Name:     name,
		TypeNode: typ,
	}
}
