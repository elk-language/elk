package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
)

// Expression of an operator with two operands eg. `2 + 5`, `foo > bar`
type BinaryExpressionNode struct {
	TypedNodeBase
	Op     *token.Token   // operator
	Left   ExpressionNode // left hand side
	Right  ExpressionNode // right hand side
	static bool
}

func (n *BinaryExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*BinaryExpressionNode)
	if !ok {
		return false
	}

	return n.Span().Equal(o.Span()) &&
		n.Op.Equal(o.Op) &&
		n.Left.Equal(value.Ref(o.Left)) &&
		n.Right.Equal(value.Ref(o.Right))
}

func (n *BinaryExpressionNode) String() string {
	var buff strings.Builder

	associativity := ExpressionAssociativity(n)

	var leftParen bool
	var rightParen bool
	if associativity == LEFT_ASSOCIATIVE {
		leftParen = ExpressionPrecedence(n) > ExpressionPrecedence(n.Left)
		rightParen = ExpressionPrecedence(n) >= ExpressionPrecedence(n.Right)
	} else {
		leftParen = ExpressionPrecedence(n) >= ExpressionPrecedence(n.Left)
		rightParen = ExpressionPrecedence(n) > ExpressionPrecedence(n.Right)
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

func (b *BinaryExpressionNode) IsStatic() bool {
	return b.static
}

// Create a new binary expression node.
func NewBinaryExpressionNode(span *position.Span, op *token.Token, left, right ExpressionNode) *BinaryExpressionNode {
	return &BinaryExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		Left:          left,
		Right:         right,
		static:        areExpressionsStatic(left, right),
	}
}

// Same as [NewBinaryExpressionNode] but returns an interface
func NewBinaryExpressionNodeI(span *position.Span, op *token.Token, left, right ExpressionNode) ExpressionNode {
	return NewBinaryExpressionNode(span, op, left, right)
}

func (*BinaryExpressionNode) Class() *value.Class {
	return value.BinaryExpressionNodeClass
}

func (*BinaryExpressionNode) DirectClass() *value.Class {
	return value.BinaryExpressionNodeClass
}

func (n *BinaryExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::BinaryExpressionNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  op: ")
	indent.IndentStringFromSecondLine(&buff, n.Op.Inspect(), 1)

	buff.WriteString(",\n  left: ")
	indent.IndentStringFromSecondLine(&buff, n.Left.Inspect(), 1)

	buff.WriteString(",\n  right: ")
	indent.IndentStringFromSecondLine(&buff, n.Right.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *BinaryExpressionNode) Error() string {
	return n.Inspect()
}
