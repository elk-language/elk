package ast

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Char literal eg. `a`
type CharLiteralNode struct {
	TypedNodeBase
	Value rune // value of the string literal
}

func (n *CharLiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &CharLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Value:         n.Value,
	}
}

func (n *CharLiteralNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::CharLiteralNode", env)
}

func (n *CharLiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
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

func (n *RawCharLiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &RawCharLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Value:         n.Value,
	}
}

func (n *RawCharLiteralNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::RawCharLiteralNode", env)
}

func (n *RawCharLiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
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
