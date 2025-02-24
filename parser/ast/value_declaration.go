package ast

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a value declaration eg. `val foo: String`
type ValueDeclarationNode struct {
	TypedNodeBase
	Name        string         // name of the value
	TypeNode    TypeNode       // type of the value
	Initialiser ExpressionNode // value assigned to the value
}

func (*ValueDeclarationNode) IsStatic() bool {
	return false
}

// Create a new value declaration node eg. `val foo: String`
func NewValueDeclarationNode(span *position.Span, name string, typ TypeNode, init ExpressionNode) *ValueDeclarationNode {
	return &ValueDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Name:          name,
		TypeNode:      typ,
		Initialiser:   init,
	}
}

func (*ValueDeclarationNode) Class() *value.Class {
	return value.ValueDeclarationNodeClass
}

func (*ValueDeclarationNode) DirectClass() *value.Class {
	return value.ValueDeclarationNodeClass
}

func (v *ValueDeclarationNode) Inspect() string {
	return fmt.Sprintf(
		"Std::AST::ValueDeclarationNode{&: %p, name: %s, type_node: %s, initialiser: %s}",
		v,
		v.Name,
		v.TypeNode.Inspect(),
		v.Initialiser.Inspect(),
	)
}

func (v *ValueDeclarationNode) Error() string {
	return v.Inspect()
}
