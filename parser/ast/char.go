package ast

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Char literal eg. `a`
type CharLiteralNode struct {
	TypedNodeBase
	Value rune // value of the string literal
}

func (n *CharLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*CharLiteralNode)
	if !ok {
		return false
	}

	return n.Value == o.Value && n.Span().Equal(o.Span())
}

func (n *CharLiteralNode) String() string {
	return value.Char(n.Value).Inspect()
}

func (*CharLiteralNode) IsStatic() bool {
	return true
}

func (*CharLiteralNode) Class() *value.Class {
	return value.CharLiteralNodeClass
}

func (*CharLiteralNode) DirectClass() *value.Class {
	return value.CharLiteralNodeClass
}

func (n *CharLiteralNode) Inspect() string {
	return fmt.Sprintf(
		"Std::Elk::AST::CharLiteralNode{span: %s, value: %s}",
		(*value.Span)(n.span).Inspect(),
		value.Char(n.Value).Inspect(),
	)
}

func (n *CharLiteralNode) Error() string {
	return n.Inspect()
}

// Create a new char literal node eg. `c"a"`
func NewCharLiteralNode(span *position.Span, val rune) *CharLiteralNode {
	return &CharLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Raw Char literal eg. r`a`
type RawCharLiteralNode struct {
	TypedNodeBase
	Value rune // value of the char literal
}

func (n *RawCharLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*RawCharLiteralNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.span.Equal(o.span)
}

func (n *RawCharLiteralNode) String() string {
	return "r" + value.Char(n.Value).Inspect()
}

func (*RawCharLiteralNode) IsStatic() bool {
	return true
}

func (*RawCharLiteralNode) Class() *value.Class {
	return value.RawCharLiteralNodeClass
}

func (*RawCharLiteralNode) DirectClass() *value.Class {
	return value.RawCharLiteralNodeClass
}

func (n *RawCharLiteralNode) Inspect() string {
	return fmt.Sprintf(
		"Std::Elk::AST::RawCharLiteralNode{span: %s, value: %s}",
		(*value.Span)(n.span).Inspect(),
		value.Char(n.Value).Inspect(),
	)
}

func (n *RawCharLiteralNode) Error() string {
	return n.Inspect()
}

// Create a new raw char literal node eg. r`a`
func NewRawCharLiteralNode(span *position.Span, val rune) *RawCharLiteralNode {
	return &RawCharLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}
