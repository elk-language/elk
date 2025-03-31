package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
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

// Equal compares this node to another value for equality.
func (n *GenericTypeDefinitionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*GenericTypeDefinitionNode)
	if !ok {
		return false
	}

	if !n.Constant.Equal(value.Ref(o.Constant)) ||
		!n.TypeNode.Equal(value.Ref(o.TypeNode)) ||
		!n.span.Equal(o.span) ||
		n.DocComment() != o.DocComment() {
		return false
	}

	if len(n.TypeParameters) != len(o.TypeParameters) {
		return false
	}

	for i, param := range n.TypeParameters {
		if !param.Equal(value.Ref(o.TypeParameters[i])) {
			return false
		}
	}

	return true
}

// String returns a string representation of the GenericTypeDefinitionNode.
func (n *GenericTypeDefinitionNode) String() string {
	var buff strings.Builder

	doc := n.DocComment()
	if len(doc) > 0 {
		buff.WriteString("##[\n")
		indent.IndentString(&buff, doc, 1)
		buff.WriteString("\n]##\n")
	}

	buff.WriteString("typedef ")
	buff.WriteString(n.Constant.String())
	buff.WriteString("[")

	for i, param := range n.TypeParameters {
		if i > 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(param.String())
	}

	buff.WriteString("] = ")
	buff.WriteString(n.TypeNode.String())

	return buff.String()
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

	fmt.Fprintf(&buff, "Std::Elk::AST::GenericTypeDefinitionNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  doc_comment: ")
	buff.WriteString(n.DocComment())

	buff.WriteString(",\n  constant: ")
	indent.IndentStringFromSecondLine(&buff, n.Constant.Inspect(), 1)

	buff.WriteString(",\n  type_node: ")
	indent.IndentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	buff.WriteString(",\n  type_parameters: %%[\n")
	for i, element := range n.TypeParameters {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
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

func (n *TypeDefinitionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*TypeDefinitionNode)
	if !ok {
		return false
	}

	return n.Constant.Equal(value.Ref(o.Constant)) &&
		n.TypeNode.Equal(value.Ref(o.TypeNode)) &&
		n.comment == o.comment &&
		n.span.Equal(o.span)
}

func (n *TypeDefinitionNode) String() string {
	var buff strings.Builder

	doc := n.DocComment()
	if len(doc) > 0 {
		buff.WriteString("##[\n")
		indent.IndentString(&buff, doc, 1)
		buff.WriteString("\n]##\n")
	}

	buff.WriteString("typedef ")
	buff.WriteString(n.Constant.String())
	buff.WriteString(" = ")
	buff.WriteString(n.TypeNode.String())

	return buff.String()
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

	fmt.Fprintf(&buff, "Std::Elk::AST::TypeDefinitionNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  doc_comment: ")
	buff.WriteString(n.DocComment())

	buff.WriteString(",\n  constant: ")
	indent.IndentStringFromSecondLine(&buff, n.Constant.Inspect(), 1)

	buff.WriteString(",\n  type_node: ")
	indent.IndentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

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
