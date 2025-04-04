package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
)

// Type expression of an operator with two operands eg. `String | Int`
type BinaryTypeNode struct {
	TypedNodeBase
	Op    *token.Token // operator
	Left  TypeNode     // left hand side
	Right TypeNode     // right hand side
}

func (n *BinaryTypeNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*BinaryTypeNode)
	if !ok {
		return false
	}

	return n.Span().Equal(o.Span()) &&
		n.Op.Equal(o.Op) &&
		n.Left.Equal(value.Ref(o.Left)) &&
		n.Right.Equal(value.Ref(o.Right))
}

func (n *BinaryTypeNode) String() string {
	var buff strings.Builder

	associativity := TypeAssociativity(n)

	var leftParen bool
	var rightParen bool
	if associativity == LEFT_ASSOCIATIVE {
		leftParen = TypePrecedence(n) > TypePrecedence(n.Left)
		rightParen = TypePrecedence(n) >= TypePrecedence(n.Right)
	} else {
		leftParen = TypePrecedence(n) >= TypePrecedence(n.Left)
		rightParen = TypePrecedence(n) > TypePrecedence(n.Right)
	}

	if leftParen {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Left.String())
	if leftParen {
		buff.WriteRune(')')
	}

	buff.WriteRune(' ')
	buff.WriteString(n.Op.FetchValue())
	buff.WriteRune(' ')

	if rightParen {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Right.String())
	if rightParen {
		buff.WriteRune(')')
	}

	return buff.String()
}

func (*BinaryTypeNode) IsStatic() bool {
	return false
}

// Create a new binary type expression node eg. `String | Int`
func NewBinaryTypeNode(span *position.Span, op *token.Token, left, right TypeNode) *BinaryTypeNode {
	return &BinaryTypeNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		Left:          left,
		Right:         right,
	}
}

// Same as [NewBinaryTypeNode] but returns an interface
func NewBinaryTypeNodeI(span *position.Span, op *token.Token, left, right TypeNode) TypeNode {
	return &BinaryTypeNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		Left:          left,
		Right:         right,
	}
}

func (*BinaryTypeNode) Class() *value.Class {
	return value.BinaryTypeNodeClass
}

func (*BinaryTypeNode) DirectClass() *value.Class {
	return value.BinaryTypeNodeClass
}

func (n *BinaryTypeNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::BinaryTypeNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  op: ")
	indent.IndentStringFromSecondLine(&buff, n.Op.Inspect(), 1)

	buff.WriteString(",\n  left: ")
	indent.IndentStringFromSecondLine(&buff, n.Left.Inspect(), 1)

	buff.WriteString(",\n  right: ")
	indent.IndentStringFromSecondLine(&buff, n.Right.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *BinaryTypeNode) Error() string {
	return n.Inspect()
}
