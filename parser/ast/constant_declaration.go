package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a constant declaration eg. `const Foo: ArrayList[String] = ["foo", "bar"]`
type ConstantDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Constant    ExpressionNode // name of the constant
	TypeNode    TypeNode       // type of the constant
	Initialiser ExpressionNode // value assigned to the constant
}

func (*ConstantDeclarationNode) IsStatic() bool {
	return false
}

// Create a new constant declaration node eg. `const Foo: ArrayList[String] = ["foo", "bar"]`
func NewConstantDeclarationNode(span *position.Span, docComment string, constant ExpressionNode, typ TypeNode, init ExpressionNode) *ConstantDeclarationNode {
	return &ConstantDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Constant:    constant,
		TypeNode:    typ,
		Initialiser: init,
	}
}

func (*ConstantDeclarationNode) Class() *value.Class {
	return value.AwaitExpressionNodeClass
}

func (*ConstantDeclarationNode) DirectClass() *value.Class {
	return value.AwaitExpressionNodeClass
}

func (n *ConstantDeclarationNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::ConstantDeclarationNode{\n  &: %p", n)

	buff.WriteString(",\n  doc_comment: ")
	indentStringFromSecondLine(&buff, value.String(n.DocComment()).Inspect(), 1)

	buff.WriteString(",\n  constant: ")
	indentStringFromSecondLine(&buff, n.Constant.Inspect(), 1)

	buff.WriteString(",\n  type_node: ")
	indentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	buff.WriteString(",\n  initialiser: ")
	indentStringFromSecondLine(&buff, n.Initialiser.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ConstantDeclarationNode) Error() string {
	return n.Inspect()
}
