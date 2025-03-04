package ast

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Char literal eg. `c"a"`
type CharLiteralNode struct {
	TypedNodeBase
	Value rune // value of the string literal
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
		"Std::Elk::AST::CharLiteralNode{&: %p, value: %s}",
		n,
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

// Raw Char literal eg. `a`
type RawCharLiteralNode struct {
	TypedNodeBase
	Value rune // value of the char literal
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
		"Std::Elk::AST::RawCharLiteralNode{&: %p, value: %s}",
		n,
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
