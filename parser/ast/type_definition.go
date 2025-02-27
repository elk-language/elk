package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a new generic type definition eg. `typedef Nilable[T] = T | nil`
type GenericTypeDefinitionNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	TypeParameters []TypeParameterNode // Generic type variable definitions
	Constant       ComplexConstantNode // new name of the type
	TypeNode       TypeNode            // the type
}

func (*GenericTypeDefinitionNode) IsStatic() bool {
	return false
}

// Create a generic type definition node eg. `typedef Nilable[T] = T | nil`
func NewGenericTypeDefinitionNode(span *position.Span, docComment string, constant ComplexConstantNode, typeVars []TypeParameterNode, typ TypeNode) *GenericTypeDefinitionNode {
	return &GenericTypeDefinitionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Constant:       constant,
		TypeParameters: typeVars,
		TypeNode:       typ,
	}
}

func (*GenericTypeDefinitionNode) Class() *value.Class {
	return value.GenericTypeDefinitionNodeClass
}

func (*GenericTypeDefinitionNode) DirectClass() *value.Class {
	return value.GenericTypeDefinitionNodeClass
}

func (n *GenericTypeDefinitionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::GenericTypeDefinitionNode{\n  &: %p", n)

	buff.WriteString(",\n  doc_comment: ")
	buff.WriteString(n.DocComment())

	buff.WriteString(",\n  constant: ")
	indentStringFromSecondLine(&buff, n.Constant.Inspect(), 1)

	buff.WriteString(",\n  type_node: ")
	indentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	buff.WriteString(",\n  type_parameters: %%[\n")
	for i, element := range n.TypeParameters {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *GenericTypeDefinitionNode) Error() string {
	return n.Inspect()
}

// Represents a new type definition eg. `typedef StringList = ArrayList[String]`
type TypeDefinitionNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Constant ComplexConstantNode // new name of the type
	TypeNode TypeNode            // the type
}

func (*TypeDefinitionNode) IsStatic() bool {
	return false
}

func (*TypeDefinitionNode) Class() *value.Class {
	return value.TypeDefinitionNodeClass
}

func (*TypeDefinitionNode) DirectClass() *value.Class {
	return value.TypeDefinitionNodeClass
}

func (n *TypeDefinitionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::TypeDefinitionNode{\n  &: %p", n)

	buff.WriteString(",\n  doc_comment: ")
	buff.WriteString(n.DocComment())

	buff.WriteString(",\n  constant: ")
	indentStringFromSecondLine(&buff, n.Constant.Inspect(), 1)

	buff.WriteString(",\n  type_node: ")
	indentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *TypeDefinitionNode) Error() string {
	return n.Inspect()
}

// Create a type definition node eg. `typedef StringList = ArrayList[String]`
func NewTypeDefinitionNode(span *position.Span, docComment string, constant ComplexConstantNode, typ TypeNode) *TypeDefinitionNode {
	return &TypeDefinitionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Constant: constant,
		TypeNode: typ,
	}
}
