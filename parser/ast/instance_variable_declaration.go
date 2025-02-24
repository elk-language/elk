package ast

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents an instance variable declaration eg. `var @foo: String`
type InstanceVariableDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Name     string   // name of the variable
	TypeNode TypeNode // type of the variable
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

func (i *InstanceVariableDeclarationNode) Inspect() string {
	return fmt.Sprintf("Std::AST::InstanceVariableDeclarationNode{&: %p, name: %s, type_node: %s}", i, i.Name, i.TypeNode.Inspect())
}

func (p *InstanceVariableDeclarationNode) Error() string {
	return p.Inspect()
}

// Create a new instance variable declaration node eg. `var @foo: String`
func NewInstanceVariableDeclarationNode(span *position.Span, docComment string, name string, typ TypeNode) *InstanceVariableDeclarationNode {
	return &InstanceVariableDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Name:     name,
		TypeNode: typ,
	}
}
