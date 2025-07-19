package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Expression of an operator with two operands eg. `2 + 5`, `foo > bar`
type BinaryExpressionNode struct {
	TypedNodeBase
	Op     *token.Token   // operator
	Left   ExpressionNode // left hand side
	Right  ExpressionNode // right hand side
	static static
}

func (n *BinaryExpressionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	left := n.Left.splice(loc, args, unquote).(ExpressionNode)
	right := n.Right.splice(loc, args, unquote).(ExpressionNode)

	return &BinaryExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Op:            n.Op.Splice(loc, unquote),
		Left:          left,
		Right:         right,
	}
}

func (n *BinaryExpressionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::BinaryExpressionNode", env)
}

func (n *BinaryExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Left.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}
	if n.Right.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (n *BinaryExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*BinaryExpressionNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
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
	if b.static == staticUnset {
		if areExpressionsStatic(b.Left, b.Right) {
			b.static = staticTrue
		} else {
			b.static = staticFalse
		}
	}
	return b.static == staticTrue
}

// Create a new binary expression node.
func NewBinaryExpressionNode(loc *position.Location, op *token.Token, left, right ExpressionNode) *BinaryExpressionNode {
	return &BinaryExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Op:            op,
		Left:          left,
		Right:         right,
	}
}

// Same as [NewBinaryExpressionNode] but returns an interface
func NewBinaryExpressionNodeI(loc *position.Location, op *token.Token, left, right ExpressionNode) ExpressionNode {
	return NewBinaryExpressionNode(loc, op, left, right)
}

func (*BinaryExpressionNode) Class() *value.Class {
	return value.BinaryExpressionNodeClass
}

func (*BinaryExpressionNode) DirectClass() *value.Class {
	return value.BinaryExpressionNodeClass
}

func (n *BinaryExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::BinaryExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

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
