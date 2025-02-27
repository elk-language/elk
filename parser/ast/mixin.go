package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

// Represents a mixin declaration eg. `mixin Foo; end`
type MixinDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Abstract              bool
	Constant              ExpressionNode      // The constant that will hold the mixin value
	TypeParameters        []TypeParameterNode // Generic type variable definitions
	Body                  []StatementNode     // body of the mixin
	IncludesAndImplements []ExpressionNode
	Bytecode              *vm.BytecodeFunction
}

func (*MixinDeclarationNode) SkipTypechecking() bool {
	return false
}

func (*MixinDeclarationNode) IsStatic() bool {
	return false
}

// Create a new mixin declaration node eg. `mixin Foo; end`
func NewMixinDeclarationNode(
	span *position.Span,
	docComment string,
	abstract bool,
	constant ExpressionNode,
	typeParams []TypeParameterNode,
	body []StatementNode,
) *MixinDeclarationNode {

	return &MixinDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Abstract:       abstract,
		Constant:       constant,
		TypeParameters: typeParams,
		Body:           body,
	}
}

func (*MixinDeclarationNode) Class() *value.Class {
	return value.MixinDeclarationNodeClass
}

func (*MixinDeclarationNode) DirectClass() *value.Class {
	return value.MixinDeclarationNodeClass
}

func (n *MixinDeclarationNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::MixinDeclarationNode{\n  &: %p", n)

	fmt.Fprintf(&buff, ",\n  abstract: %t", n.Abstract)

	buff.WriteString(",\n  doc_comment: ")
	indentStringFromSecondLine(&buff, value.String(n.DocComment()).Inspect(), 1)

	buff.WriteString(",\n  constant: ")
	indentStringFromSecondLine(&buff, n.Constant.Inspect(), 1)

	buff.WriteString(",\n  type_parameters: %%[\n")
	for i, element := range n.TypeParameters {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

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

func (n *MixinDeclarationNode) Error() string {
	return n.Inspect()
}
