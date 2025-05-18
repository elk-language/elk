package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
)

// Check whether the node can be used as a range pattern element.
func IsValidRangePatternElement(node Node) bool {
	switch node.(type) {
	case *TrueLiteralNode, *FalseLiteralNode, *NilLiteralNode, *CharLiteralNode,
		*RawCharLiteralNode, *RawStringLiteralNode, *DoubleQuotedStringLiteralNode,
		*InterpolatedStringLiteralNode, *SimpleSymbolLiteralNode, *InterpolatedSymbolLiteralNode,
		*FloatLiteralNode, *Float64LiteralNode, *Float32LiteralNode, *BigFloatLiteralNode,
		*IntLiteralNode, *Int64LiteralNode, *UInt64LiteralNode, *Int32LiteralNode, *UInt32LiteralNode,
		*Int16LiteralNode, *UInt16LiteralNode, *Int8LiteralNode, *UInt8LiteralNode,
		*PublicConstantNode, *PrivateConstantNode, *ConstantLookupNode, *UnaryExpressionNode:
		return true
	default:
		return false
	}
}

// Represents a Range literal eg. `1...5`
type RangeLiteralNode struct {
	TypedNodeBase
	Start  ExpressionNode
	End    ExpressionNode
	Op     *token.Token
	static bool
}

func (n *RangeLiteralNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	var start ExpressionNode
	if n.Start != nil {
		start = n.Start.Splice(loc, args, unquote).(ExpressionNode)
	}

	var end ExpressionNode
	if n.End != nil {
		end = n.End.Splice(loc, args, unquote).(ExpressionNode)
	}

	return &RangeLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Start:         start,
		End:           end,
		Op:            n.Op.Splice(loc, unquote),
		static:        areExpressionsStatic(start, end),
	}
}

func (n *RangeLiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Start != nil {
		if n.Start.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	if n.End != nil {
		if n.End.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

func (n *RangeLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*RangeLiteralNode)
	if !ok {
		return false
	}

	return n.Start.Equal(value.Ref(o.Start)) &&
		n.End.Equal(value.Ref(o.End)) &&
		n.Op.Equal(o.Op) &&
		n.loc.Equal(o.loc)
}

func (n *RangeLiteralNode) String() string {
	var buff strings.Builder

	leftParen := ExpressionPrecedence(n) > ExpressionPrecedence(n.Start)
	rightParen := ExpressionPrecedence(n) >= ExpressionPrecedence(n.End)

	if leftParen {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Start.String())
	if leftParen {
		buff.WriteRune(')')
	}

	buff.WriteString(n.Op.String())

	if rightParen {
		buff.WriteRune('(')
	}
	buff.WriteString(n.End.String())
	if rightParen {
		buff.WriteRune(')')
	}

	return buff.String()
}

func (r *RangeLiteralNode) IsStatic() bool {
	return r.static
}

// Create a Range literal node eg. `1...5`
func NewRangeLiteralNode(loc *position.Location, op *token.Token, start, end ExpressionNode) *RangeLiteralNode {
	return &RangeLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Op:            op,
		Start:         start,
		End:           end,
		static:        areExpressionsStatic(start, end),
	}
}

func (*RangeLiteralNode) Class() *value.Class {
	return value.RangeLiteralNodeClass
}

func (*RangeLiteralNode) DirectClass() *value.Class {
	return value.RangeLiteralNodeClass
}

func (n *RangeLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::RangeLiteralNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  start: ")
	if n.Start == nil {
		buff.WriteString("nil")
	} else {
		indent.IndentStringFromSecondLine(&buff, n.Start.Inspect(), 1)
	}

	buff.WriteString(",\n  end: ")
	if n.End == nil {
		buff.WriteString("nil")
	} else {
		indent.IndentStringFromSecondLine(&buff, n.End.Inspect(), 1)
	}

	buff.WriteString(",\n  op: ")
	indent.IndentStringFromSecondLine(&buff, n.Op.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *RangeLiteralNode) Error() string {
	return n.Inspect()
}
