package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a double splat expression eg. `**foo`
type DoubleSplatExpressionNode struct {
	TypedNodeBase
	Value ExpressionNode
}

func (n *DoubleSplatExpressionNode) Splice(loc *position.Location, args *[]Node) Node {
	return &DoubleSplatExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: getLoc(loc, n.loc), typ: n.typ},
		Value:         n.Value.Splice(loc, args).(ExpressionNode),
	}
}

// Check if this node equals another node.
func (n *DoubleSplatExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*DoubleSplatExpressionNode)
	if !ok {
		return false
	}

	return n.Value.Equal(value.Ref(o.Value)) &&
		n.loc.Equal(o.loc)
}

// Return a string representation of the node.
func (n *DoubleSplatExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("**")

	parens := ExpressionPrecedence(n) > ExpressionPrecedence(n.Value)
	if parens {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Value.String())
	if parens {
		buff.WriteRune(')')
	}

	return buff.String()
}

func (*DoubleSplatExpressionNode) IsStatic() bool {
	return false
}

func (*DoubleSplatExpressionNode) Class() *value.Class {
	return value.DoubleSplatExpressionNodeClass
}

func (*DoubleSplatExpressionNode) DirectClass() *value.Class {
	return value.DoubleSplatExpressionNodeClass
}

func (n *DoubleSplatExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::DoubleSplatExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  value: ")
	indent.IndentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *DoubleSplatExpressionNode) Error() string {
	return n.Inspect()
}

// Create a double splat expression node eg. `**foo`
func NewDoubleSplatExpressionNode(loc *position.Location, val ExpressionNode) *DoubleSplatExpressionNode {
	return &DoubleSplatExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
	}
}

// Represents a splat expression eg. `*foo`
type SplatExpressionNode struct {
	TypedNodeBase
	Value ExpressionNode
}

func (n *SplatExpressionNode) Splice(loc *position.Location, args *[]Node) Node {
	return &SplatExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: getLoc(loc, n.loc), typ: n.typ},
		Value:         n.Value.Splice(loc, args).(ExpressionNode),
	}
}

func (n *SplatExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*SplatExpressionNode)
	if !ok {
		return false
	}

	return n.Value.Equal(value.Ref(o.Value)) &&
		n.loc.Equal(o.loc)
}

func (n *SplatExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteRune('*')
	parens := ExpressionPrecedence(n) > ExpressionPrecedence(n.Value)

	if parens {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Value.String())
	if parens {
		buff.WriteRune(')')
	}

	return buff.String()
}

func (*SplatExpressionNode) IsStatic() bool {
	return false
}

func (*SplatExpressionNode) Class() *value.Class {
	return value.SplatExpressionNodeClass
}

func (*SplatExpressionNode) DirectClass() *value.Class {
	return value.SplatExpressionNodeClass
}

func (n *SplatExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::SplatExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  value: ")
	indent.IndentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *SplatExpressionNode) Error() string {
	return n.Inspect()
}

// Create a splat expression node eg. `*foo`
func NewSplatExpressionNode(loc *position.Location, val ExpressionNode) *SplatExpressionNode {
	return &SplatExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
	}
}
