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

func (n *BinaryTypeNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &BinaryTypeNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Op:            n.Op.Splice(loc, unquote),
		Left:          n.Left.splice(loc, args, unquote).(TypeNode),
		Right:         n.Right.splice(loc, args, unquote).(TypeNode),
	}
}

func (n *BinaryTypeNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
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

func (n *BinaryTypeNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*BinaryTypeNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
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
func NewBinaryTypeNode(loc *position.Location, op *token.Token, left, right TypeNode) *BinaryTypeNode {
	return &BinaryTypeNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Op:            op,
		Left:          left,
		Right:         right,
	}
}

// Same as [NewBinaryTypeNode] but returns an interface
func NewBinaryTypeNodeI(loc *position.Location, op *token.Token, left, right TypeNode) TypeNode {
	return &BinaryTypeNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
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

	fmt.Fprintf(&buff, "Std::Elk::AST::BinaryTypeNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

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
