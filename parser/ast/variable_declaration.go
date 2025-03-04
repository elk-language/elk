package ast

import (
	"fmt"
	"strings"

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

	fmt.Fprintf(&buff, "Std::Elk::AST::VariableDeclarationNode{\n  &: %p", n)

	buff.WriteString(",\n  name: ")
	buff.WriteString(n.Name)

	buff.WriteString(",\n  type_node: ")
	indentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	buff.WriteString(",\n  initialiser: ")
	indentStringFromSecondLine(&buff, n.Initialiser.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (v *VariableDeclarationNode) Error() string {
	return v.Inspect()
}

// Create a new variable declaration node eg. `var foo: String`
func NewVariableDeclarationNode(span *position.Span, docComment string, name string, typ TypeNode, init ExpressionNode) *VariableDeclarationNode {
	return &VariableDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Name:        name,
		TypeNode:    typ,
		Initialiser: init,
	}
}
