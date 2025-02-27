package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

// Represents a module declaration eg. `module Foo; end`
type ModuleDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Constant ExpressionNode  // The constant that will hold the module value
	Body     []StatementNode // body of the module
	Bytecode *vm.BytecodeFunction
}

func (*ModuleDeclarationNode) SkipTypechecking() bool {
	return false
}

func (*ModuleDeclarationNode) IsStatic() bool {
	return false
}

// Create a new module declaration node eg. `module Foo; end`
func NewModuleDeclarationNode(
	span *position.Span,
	docComment string,
	constant ExpressionNode,
	body []StatementNode,
) *ModuleDeclarationNode {

	return &ModuleDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Constant: constant,
		Body:     body,
	}
}

func (*ModuleDeclarationNode) Class() *value.Class {
	return value.ClassDeclarationNodeClass
}

func (*ModuleDeclarationNode) DirectClass() *value.Class {
	return value.ClassDeclarationNodeClass
}

func (n *ModuleDeclarationNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::ModuleDeclarationNode{\n  &: %p", n)

	buff.WriteString(",\n  doc_comment: ")
	indentStringFromSecondLine(&buff, value.String(n.DocComment()).Inspect(), 1)

	buff.WriteString(",\n  constant: ")
	indentStringFromSecondLine(&buff, n.Constant.Inspect(), 1)

	buff.WriteString(",\n  body: %%[\n")
	for i, element := range n.Body {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ModuleDeclarationNode) Error() string {
	return n.Inspect()
}
