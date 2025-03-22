package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

// Represents an interface declaration eg. `interface Foo; end`
type InterfaceDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Constant       ExpressionNode      // The constant that will hold the interface value
	TypeParameters []TypeParameterNode // Generic type variable definitions
	Body           []StatementNode     // body of the interface
	Implements     []*ImplementExpressionNode
	Bytecode       *vm.BytecodeFunction
}

func (*InterfaceDeclarationNode) SkipTypechecking() bool {
	return false
}

func (*InterfaceDeclarationNode) IsStatic() bool {
	return false
}

// Create a new interface declaration node eg. `interface Foo; end`
func NewInterfaceDeclarationNode(
	span *position.Span,
	docComment string,
	constant ExpressionNode,
	typeParams []TypeParameterNode,
	body []StatementNode,
) *InterfaceDeclarationNode {

	return &InterfaceDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Constant:       constant,
		TypeParameters: typeParams,
		Body:           body,
	}
}

func (*InterfaceDeclarationNode) Class() *value.Class {
	return value.InterfaceDeclarationNodeClass
}

func (*InterfaceDeclarationNode) DirectClass() *value.Class {
	return value.InterfaceDeclarationNodeClass
}

func (n *InterfaceDeclarationNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::InterfaceDeclarationNode{\n  &: %p", n)

	buff.WriteString(",\n  doc_comment: ")
	indent.IndentStringFromSecondLine(&buff, value.String(n.DocComment()).Inspect(), 1)

	buff.WriteString(",\n  constant: ")
	indent.IndentStringFromSecondLine(&buff, n.Constant.Inspect(), 1)

	buff.WriteString(",\n  type_parameters: %%[\n")
	for i, element := range n.TypeParameters {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  body: %%[\n")
	for i, element := range n.Body {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *InterfaceDeclarationNode) Error() string {
	return n.Inspect()
}
