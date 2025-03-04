package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

// Nodes that implement this interface represent
// symbol literals.
type SymbolLiteralNode interface {
	Node
	ExpressionNode
	symbolLiteralNode()
}

func (*InvalidNode) symbolLiteralNode()                   {}
func (*SimpleSymbolLiteralNode) symbolLiteralNode()       {}
func (*InterpolatedSymbolLiteralNode) symbolLiteralNode() {}

// Represents a symbol literal with simple content eg. `:foo`, `:'foo bar`, `:"lol"`
type SimpleSymbolLiteralNode struct {
	TypedNodeBase
	Content string
}

func (*SimpleSymbolLiteralNode) IsStatic() bool {
	return true
}

// Create a simple symbol literal node eg. `:foo`, `:'foo bar`, `:"lol"`
func NewSimpleSymbolLiteralNode(span *position.Span, cont string) *SimpleSymbolLiteralNode {
	return &SimpleSymbolLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Content:       cont,
	}
}

func (*SimpleSymbolLiteralNode) Class() *value.Class {
	return value.SimpleSymbolLiteralNodeClass
}

func (*SimpleSymbolLiteralNode) DirectClass() *value.Class {
	return value.SimpleSymbolLiteralNodeClass
}

func (n *SimpleSymbolLiteralNode) Inspect() string {
	return fmt.Sprintf(
		"Std::Elk::AST::SimpleSymbolLiteralNode{&: %p, content: %s}",
		n,
		value.String(n.Content).Inspect(),
	)
}

func (n *SimpleSymbolLiteralNode) Error() string {
	return n.Inspect()
}

// Represents an interpolated symbol eg. `:"foo ${bar + 2}"`
type InterpolatedSymbolLiteralNode struct {
	NodeBase
	Content *InterpolatedStringLiteralNode
}

func (*InterpolatedSymbolLiteralNode) IsStatic() bool {
	return false
}

func (*InterpolatedSymbolLiteralNode) Type(globalEnv *types.GlobalEnvironment) types.Type {
	return globalEnv.StdSubtype(symbol.Symbol)
}

func (*InterpolatedSymbolLiteralNode) Class() *value.Class {
	return value.InterpolatedSymbolLiteralNodeClass
}

func (*InterpolatedSymbolLiteralNode) DirectClass() *value.Class {
	return value.InterpolatedSymbolLiteralNodeClass
}

func (n *InterpolatedSymbolLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::InterpolatedSymbolLiteralNode{\n  &: %p", n)

	buff.WriteString(",\n  content: ")
	indentStringFromSecondLine(&buff, n.Content.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *InterpolatedSymbolLiteralNode) Error() string {
	return n.Inspect()
}

// Create an interpolated symbol literal node eg. `:"foo ${bar + 2}"`
func NewInterpolatedSymbolLiteralNode(span *position.Span, cont *InterpolatedStringLiteralNode) *InterpolatedSymbolLiteralNode {
	return &InterpolatedSymbolLiteralNode{
		NodeBase: NodeBase{span: span},
		Content:  cont,
	}
}
