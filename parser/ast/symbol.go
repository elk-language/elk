package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
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

func (n *SimpleSymbolLiteralNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &SimpleSymbolLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Content:       n.Content,
	}
}

func (n *SimpleSymbolLiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	}

	return leave(n, parent)
}

func (n *SimpleSymbolLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*SimpleSymbolLiteralNode)
	if !ok {
		return false
	}

	return n.Content == o.Content &&
		n.loc.Equal(o.loc)
}

func (n *SimpleSymbolLiteralNode) String() string {
	return value.InspectSymbol(n.Content)
}

func (*SimpleSymbolLiteralNode) IsStatic() bool {
	return true
}

// Create a simple symbol literal node eg. `:foo`, `:'foo bar`, `:"lol"`
func NewSimpleSymbolLiteralNode(loc *position.Location, cont string) *SimpleSymbolLiteralNode {
	return &SimpleSymbolLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
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
		"Std::Elk::AST::SimpleSymbolLiteralNode{location: %s, content: %s}",
		(*value.Location)(n.loc).Inspect(),
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

func (n *InterpolatedSymbolLiteralNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &InterpolatedSymbolLiteralNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Content:  n.Content.Splice(loc, args, unquote).(*InterpolatedStringLiteralNode),
	}
}

func (n *InterpolatedSymbolLiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Content.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (n *InterpolatedSymbolLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*InterpolatedSymbolLiteralNode)
	if !ok {
		return false
	}

	return n.Content.Equal(value.Ref(o.Content)) &&
		n.loc.Equal(o.loc)
}

func (n *InterpolatedSymbolLiteralNode) String() string {
	var buff strings.Builder

	buff.WriteString(":")
	buff.WriteString(n.Content.String())

	return buff.String()
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

	fmt.Fprintf(&buff, "Std::Elk::AST::InterpolatedSymbolLiteralNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  content: ")
	indent.IndentStringFromSecondLine(&buff, n.Content.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *InterpolatedSymbolLiteralNode) Error() string {
	return n.Inspect()
}

// Create an interpolated symbol literal node eg. `:"foo ${bar + 2}"`
func NewInterpolatedSymbolLiteralNode(loc *position.Location, cont *InterpolatedStringLiteralNode) *InterpolatedSymbolLiteralNode {
	return &InterpolatedSymbolLiteralNode{
		NodeBase: NodeBase{loc: loc},
		Content:  cont,
	}
}
