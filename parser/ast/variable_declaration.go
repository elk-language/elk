package ast

import (
	"fmt"

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

func (v *VariableDeclarationNode) Inspect() string {
	return fmt.Sprintf(
		"Std::AST::VariableDeclarationNode{&: %p, name: %s, type_node: %s, initialiser: %s}",
		v,
		v.Name,
		v.TypeNode.Inspect(),
		v.Initialiser.Inspect(),
	)
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
