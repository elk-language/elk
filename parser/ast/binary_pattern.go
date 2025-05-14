package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
)

// Pattern with two operands eg. `> 10 && < 50`
type BinaryPatternNode struct {
	TypedNodeBase
	Op    *token.Token // operator
	Left  PatternNode  // left hand side
	Right PatternNode  // right hand side
}

func (n *BinaryPatternNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &BinaryPatternNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Op:            n.Op.Splice(loc, unquote),
		Left:          n.Left.Splice(loc, args, unquote).(PatternNode),
		Right:         n.Right.Splice(loc, args, unquote).(PatternNode),
	}
}

func (n *BinaryPatternNode) Traverse(yield func(Node) bool) bool {
	if !n.Left.Traverse(yield) {
		return false
	}
	if !n.Right.Traverse(yield) {
		return false
	}
	return yield(n)
}

func (n *BinaryPatternNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*BinaryPatternNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.Op.Equal(o.Op) &&
		n.Left.Equal(value.Ref(o.Left)) &&
		n.Right.Equal(value.Ref(o.Right))
}

func (n *BinaryPatternNode) String() string {
	var buff strings.Builder

	associativity := PatternAssociativity(n)

	var leftParen bool
	var rightParen bool
	if associativity == LEFT_ASSOCIATIVE {
		leftParen = PatternPrecedence(n) > PatternPrecedence(n.Left)
		rightParen = PatternPrecedence(n) >= PatternPrecedence(n.Right)
	} else {
		leftParen = PatternPrecedence(n) >= PatternPrecedence(n.Left)
		rightParen = PatternPrecedence(n) > PatternPrecedence(n.Right)
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

func (*BinaryPatternNode) IsStatic() bool {
	return false
}

// Create a new binary pattern node eg. `> 10 && < 50`
func NewBinaryPatternNode(loc *position.Location, op *token.Token, left, right PatternNode) *BinaryPatternNode {
	return &BinaryPatternNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Op:            op,
		Left:          left,
		Right:         right,
	}
}

// Same as [NewBinaryPatternNode] but returns an interface
func NewBinaryPatternNodeI(loc *position.Location, op *token.Token, left, right PatternNode) PatternNode {
	return NewBinaryPatternNode(loc, op, left, right)
}

func (*BinaryPatternNode) Class() *value.Class {
	return value.BinaryPatternNodeClass
}

func (*BinaryPatternNode) DirectClass() *value.Class {
	return value.BinaryPatternNodeClass
}

func (n *BinaryPatternNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::BinaryPatternNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  op: ")
	indent.IndentStringFromSecondLine(&buff, n.Op.Inspect(), 1)

	buff.WriteString(",\n  left: ")
	indent.IndentStringFromSecondLine(&buff, n.Left.Inspect(), 1)

	buff.WriteString(",\n  right: ")
	indent.IndentStringFromSecondLine(&buff, n.Right.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *BinaryPatternNode) Error() string {
	return n.Inspect()
}
