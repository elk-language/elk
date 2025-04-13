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

func (n *CharLiteralNode) Splice(loc *position.Location, args *[]Node) Node {
	return &CharLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: getLoc(loc, n.loc), typ: n.typ},
		Value:         n.Value,
	}
}

func (n *CharLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*CharLiteralNode)
	if !ok {
		return false
	}

	return n.Value == o.Value && n.loc.Equal(o.loc)
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
		"Std::Elk::AST::CharLiteralNode{location: %s, value: %s}",
		(*value.Location)(n.loc).Inspect(),
		value.Char(n.Value).Inspect(),
	)
}

func (n *CharLiteralNode) Error() string {
	return n.Inspect()
}

// Create a new char literal node eg. `c"a"`
func NewCharLiteralNode(loc *position.Location, val rune) *CharLiteralNode {
	return &CharLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
	}
}

// Raw Char literal eg. r`a`
type RawCharLiteralNode struct {
	TypedNodeBase
	Value rune // value of the char literal
}

func (n *RawCharLiteralNode) Splice(loc *position.Location, args *[]Node) Node {
	return &RawCharLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: getLoc(loc, n.loc), typ: n.typ},
		Value:         n.Value,
	}
}

func (n *RawCharLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*RawCharLiteralNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.loc.Equal(o.loc)
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
		"Std::Elk::AST::RawCharLiteralNode{location: %s, value: %s}",
		(*value.Location)(n.loc).Inspect(),
		value.Char(n.Value).Inspect(),
	)
}

func (n *RawCharLiteralNode) Error() string {
	return n.Inspect()
}

// Create a new raw char literal node eg. r`a`
func NewRawCharLiteralNode(loc *position.Location, val rune) *RawCharLiteralNode {
	return &RawCharLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
	}
}
