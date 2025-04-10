package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a variable declaration eg. `var foo: String`
type VariableDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Name        string         // name of the variable
	TypeNode    TypeNode       // type of the variable
	Initialiser ExpressionNode // value assigned to the variable
}

func (n *VariableDeclarationNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*VariableDeclarationNode)
	if !ok {
		return false
	}

	if n.Name != o.Name ||
		n.comment != o.comment ||
		!n.loc.Equal(o.loc) {
		return false
	}

	if n.TypeNode == o.TypeNode {
	} else if n.TypeNode == nil || o.TypeNode == nil {
		return false
	} else if !n.TypeNode.Equal(value.Ref(o.TypeNode)) {
		return false
	}

	if n.Initialiser == o.Initialiser {
	} else if n.Initialiser == nil || o.Initialiser == nil {
		return false
	} else if !n.Initialiser.Equal(value.Ref(o.Initialiser)) {
		return false
	}

	return true
}

func (n *VariableDeclarationNode) String() string {
	var buff strings.Builder

	buff.WriteString("var ")
	buff.WriteString(n.Name)

	if n.TypeNode != nil {
		buff.WriteString(": ")
		buff.WriteString(n.TypeNode.String())
	}

	if n.Initialiser != nil {
		buff.WriteString(" = ")
		buff.WriteString(n.Initialiser.String())
	}

	return buff.String()
}

func (*VariableDeclarationNode) IsStatic() bool {
	return false
}

func (*VariableDeclarationNode) Class() *value.Class {
	return value.VariableDeclarationNodeClass
}

func (*VariableDeclarationNode) DirectClass() *value.Class {
	return value.VariableDeclarationNodeClass
}

func (n *VariableDeclarationNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::VariableDeclarationNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  name: ")
	buff.WriteString(n.Name)

	buff.WriteString(",\n  type_node: ")
	indent.IndentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	buff.WriteString(",\n  initialiser: ")
	indent.IndentStringFromSecondLine(&buff, n.Initialiser.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (v *VariableDeclarationNode) Error() string {
	return v.Inspect()
}

// Create a new variable declaration node eg. `var foo: String`
func NewVariableDeclarationNode(loc *position.Location, docComment string, name string, typ TypeNode, init ExpressionNode) *VariableDeclarationNode {
	return &VariableDeclarationNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Name:        name,
		TypeNode:    typ,
		Initialiser: init,
	}
}
