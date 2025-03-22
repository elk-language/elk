package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

// Represents a class declaration eg. `class Foo; end`
type ClassDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Abstract       bool
	Sealed         bool
	Primitive      bool
	NoInit         bool
	Constant       ExpressionNode      // The constant that will hold the class value
	TypeParameters []TypeParameterNode // Generic type variable definitions
	Superclass     ExpressionNode      // the super/parent class of this class
	Body           []StatementNode     // body of the class
	Bytecode       *vm.BytecodeFunction
}

func (*ClassDeclarationNode) SkipTypechecking() bool {
	return false
}

func (*ClassDeclarationNode) IsStatic() bool {
	return false
}

// Create a new class declaration node eg. `class Foo; end`
func NewClassDeclarationNode(
	span *position.Span,
	docComment string,
	abstract bool,
	sealed bool,
	primitive bool,
	noinit bool,
	constant ExpressionNode,
	typeParams []TypeParameterNode,
	superclass ExpressionNode,
	body []StatementNode,
) *ClassDeclarationNode {

	return &ClassDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Abstract:       abstract,
		Sealed:         sealed,
		Primitive:      primitive,
		NoInit:         noinit,
		Constant:       constant,
		TypeParameters: typeParams,
		Superclass:     superclass,
		Body:           body,
	}
}

func (*ClassDeclarationNode) Class() *value.Class {
	return value.ClassDeclarationNodeClass
}

func (*ClassDeclarationNode) DirectClass() *value.Class {
	return value.ClassDeclarationNodeClass
}

func (n *ClassDeclarationNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ClassDeclarationNode{\n  &: %p", n)

	fmt.Fprintf(&buff, ",\n  abstract: %t", n.Abstract)
	fmt.Fprintf(&buff, ",\n  sealed: %t", n.Sealed)
	fmt.Fprintf(&buff, ",\n  primitive: %t", n.Primitive)
	fmt.Fprintf(&buff, ",\n  noinit: %t", n.NoInit)

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

	buff.WriteString(",\n  superclass: ")
	indent.IndentStringFromSecondLine(&buff, n.Superclass.Inspect(), 1)

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

func (n *ClassDeclarationNode) Error() string {
	return n.Inspect()
}
