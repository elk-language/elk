package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
)

// Type expression of an operator with two operands eg. `String | Int`
type BinaryTypeExpressionNode struct {
	TypedNodeBase
	Op    *token.Token // operator
	Left  TypeNode     // left hand side
	Right TypeNode     // right hand side
}

func (*BinaryTypeExpressionNode) IsStatic() bool {
	return false
}

// Create a new binary type expression node eg. `String | Int`
func NewBinaryTypeExpressionNode(span *position.Span, op *token.Token, left, right TypeNode) *BinaryTypeExpressionNode {
	return &BinaryTypeExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		Left:          left,
		Right:         right,
	}
}

// Same as [NewBinaryTypeExpressionNode] but returns an interface
func NewBinaryTypeExpressionNodeI(span *position.Span, op *token.Token, left, right TypeNode) TypeNode {
	return &BinaryTypeExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		Left:          left,
		Right:         right,
	}
}

func (*BinaryTypeExpressionNode) Class() *value.Class {
	return value.ReturnExpressionNodeClass
}

func (*BinaryTypeExpressionNode) DirectClass() *value.Class {
	return value.ReturnExpressionNodeClass
}

func (n *BinaryTypeExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::BinaryTypeExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  op: ")
	indentStringFromSecondLine(&buff, n.Op.Inspect(), 1)

	buff.WriteString(",\n  left: ")
	indentStringFromSecondLine(&buff, n.Left.Inspect(), 1)

	buff.WriteString(",\n  right: ")
	indentStringFromSecondLine(&buff, n.Right.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *BinaryTypeExpressionNode) Error() string {
	return n.Inspect()
}
