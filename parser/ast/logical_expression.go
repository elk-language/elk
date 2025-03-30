package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
)

// Expression of a logical operator with two operands eg. `foo && bar`
type LogicalExpressionNode struct {
	TypedNodeBase
	Op     *token.Token   // operator
	Left   ExpressionNode // left hand side
	Right  ExpressionNode // right hand side
	static bool
}

func (n *LogicalExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*LogicalExpressionNode)
	if !ok {
		return false
	}

	return n.Op.Equal(o.Op) &&
		n.Left.Equal(value.Ref(o.Left)) &&
		n.Right.Equal(value.Ref(o.Right)) &&
		n.span.Equal(o.span)
}

func (n *LogicalExpressionNode) String() string {
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

	buff.WriteString(" ")
	buff.WriteString(n.Op.String())
	buff.WriteString(" ")

	if rightParen {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Right.String())
	if rightParen {
		buff.WriteRune(')')
	}

	return buff.String()
}

func (l *LogicalExpressionNode) IsStatic() bool {
	return l.static
}

// Create a new logical expression node.
func NewLogicalExpressionNode(span *position.Span, op *token.Token, left, right ExpressionNode) *LogicalExpressionNode {
	return &LogicalExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		Left:          left,
		Right:         right,
		static:        areExpressionsStatic(left, right),
	}
}

// Same as [NewLogicalExpressionNode] but returns an interface
func NewLogicalExpressionNodeI(span *position.Span, op *token.Token, left, right ExpressionNode) ExpressionNode {
	return &LogicalExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		Left:          left,
		Right:         right,
	}
}

func (*LogicalExpressionNode) Class() *value.Class {
	return value.LogicalExpressionNodeClass
}

func (*LogicalExpressionNode) DirectClass() *value.Class {
	return value.LogicalExpressionNodeClass
}

func (n *LogicalExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::LogicalExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  op: ")
	indent.IndentStringFromSecondLine(&buff, n.Op.Inspect(), 1)

	buff.WriteString(",\n  left: ")
	indent.IndentStringFromSecondLine(&buff, n.Left.Inspect(), 1)

	buff.WriteString(",\n  right: ")
	indent.IndentStringFromSecondLine(&buff, n.Right.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *LogicalExpressionNode) Error() string {
	return n.Inspect()
}
