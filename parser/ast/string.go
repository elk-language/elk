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

// Nodes that implement this interface can appear
// inside of a String literal.
type StringLiteralContentNode interface {
	Node
	stringLiteralContentNode()
}

func (*InvalidNode) stringLiteralContentNode()                     {}
func (*StringInspectInterpolationNode) stringLiteralContentNode()  {}
func (*StringInterpolationNode) stringLiteralContentNode()         {}
func (*StringLiteralContentSectionNode) stringLiteralContentNode() {}

// All nodes that represent strings should
// implement this interface.
type StringLiteralNode interface {
	Node
	PatternExpressionNode
	StringOrSymbolLiteralNode
	stringLiteralNode()
}

func (*InvalidNode) stringLiteralNode()                   {}
func (*DoubleQuotedStringLiteralNode) stringLiteralNode() {}
func (*RawStringLiteralNode) stringLiteralNode()          {}
func (*InterpolatedStringLiteralNode) stringLiteralNode() {}

// All nodes that represent simple strings (without interpolation)
// should implement this interface.
type SimpleStringLiteralNode interface {
	Node
	ExpressionNode
	StringLiteralNode
	StringOrSymbolLiteralNode
	simpleStringLiteralNode()
}

func (*InvalidNode) simpleStringLiteralNode()                   {}
func (*DoubleQuotedStringLiteralNode) simpleStringLiteralNode() {}
func (*RawStringLiteralNode) simpleStringLiteralNode()          {}

// Raw string literal enclosed with single quotes eg. `'foo'`.
type RawStringLiteralNode struct {
	TypedNodeBase
	Value string // value of the string literal
}

func (n *RawStringLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*RawStringLiteralNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.span.Equal(o.span)
}

func (n *RawStringLiteralNode) String() string {
	return fmt.Sprintf("'%s'", n.Value)
}

func (*RawStringLiteralNode) IsStatic() bool {
	return true
}

func (*RawStringLiteralNode) Class() *value.Class {
	return value.RawStringLiteralNodeClass
}

func (*RawStringLiteralNode) DirectClass() *value.Class {
	return value.RawStringLiteralNodeClass
}

func (n *RawStringLiteralNode) Inspect() string {
	return fmt.Sprintf(
		"Std::Elk::AST::RawStringLiteralNode{&: %p, value: %s}",
		n,
		value.String(n.Value).Inspect(),
	)
}

func (n *RawStringLiteralNode) Error() string {
	return n.Inspect()
}

// Create a new raw string literal node eg. `'foo'`.
func NewRawStringLiteralNode(span *position.Span, val string) *RawStringLiteralNode {
	return &RawStringLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Represents a single section of characters of a string literal eg. `foo` in `"foo${bar}"`.
type StringLiteralContentSectionNode struct {
	NodeBase
	Value string
}

func (*StringLiteralContentSectionNode) IsStatic() bool {
	return true
}

func (*StringLiteralContentSectionNode) Class() *value.Class {
	return value.StringLiteralContentSectionNodeClass
}

func (*StringLiteralContentSectionNode) DirectClass() *value.Class {
	return value.StringLiteralContentSectionNodeClass
}

func (n *StringLiteralContentSectionNode) Inspect() string {
	return fmt.Sprintf(
		"Std::Elk::AST::StringLiteralContentSectionNode{&: %p, value: %s}",
		n,
		value.String(n.Value).Inspect(),
	)
}

func (n *StringLiteralContentSectionNode) Error() string {
	return n.Inspect()
}

// Create a new string literal content section node eg. `foo` in `"foo${bar}"`.
func NewStringLiteralContentSectionNode(span *position.Span, val string) *StringLiteralContentSectionNode {
	return &StringLiteralContentSectionNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Represents a single inspect interpolated section of a string literal eg. `bar + 2` in `"foo#{bar + 2}"`
type StringInspectInterpolationNode struct {
	NodeBase
	Expression ExpressionNode
}

func (*StringInspectInterpolationNode) IsStatic() bool {
	return false
}

func (*StringInspectInterpolationNode) Class() *value.Class {
	return value.StringInspectInterpolationNodeClass
}

func (*StringInspectInterpolationNode) DirectClass() *value.Class {
	return value.StringInspectInterpolationNodeClass
}

func (n *StringInspectInterpolationNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::StringInspectInterpolationNode{\n  &: %p", n)

	buff.WriteString(",\n  expression: ")
	indent.IndentStringFromSecondLine(&buff, n.Expression.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *StringInspectInterpolationNode) Error() string {
	return n.Inspect()
}

// Create a new string inspect interpolation node eg. `bar + 2` in `"foo#{bar + 2}"`
func NewStringInspectInterpolationNode(span *position.Span, expr ExpressionNode) *StringInspectInterpolationNode {
	return &StringInspectInterpolationNode{
		NodeBase:   NodeBase{span: span},
		Expression: expr,
	}
}

// Represents a single interpolated section of a string literal eg. `bar + 2` in `"foo${bar + 2}"`
type StringInterpolationNode struct {
	NodeBase
	Expression ExpressionNode
}

func (*StringInterpolationNode) IsStatic() bool {
	return false
}

func (*StringInterpolationNode) Class() *value.Class {
	return value.StringInterpolationNodeClass
}

func (*StringInterpolationNode) DirectClass() *value.Class {
	return value.StringInterpolationNodeClass
}

func (n *StringInterpolationNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::StringInterpolationNode{\n  &: %p", n)

	buff.WriteString(",\n  expression: ")
	indent.IndentStringFromSecondLine(&buff, n.Expression.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *StringInterpolationNode) Error() string {
	return n.Inspect()
}

// Create a new string interpolation node eg. `bar + 2` in `"foo${bar + 2}"`
func NewStringInterpolationNode(span *position.Span, expr ExpressionNode) *StringInterpolationNode {
	return &StringInterpolationNode{
		NodeBase:   NodeBase{span: span},
		Expression: expr,
	}
}

// Represents an interpolated string literal eg. `"foo ${bar} baz"`
type InterpolatedStringLiteralNode struct {
	NodeBase
	Content []StringLiteralContentNode
}

func (n *InterpolatedStringLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*InterpolatedStringLiteralNode)
	if !ok {
		return false
	}

	if len(n.Content) != len(o.Content) {
		return false
	}

	for i, content := range n.Content {
		if !content.Equal(value.Ref(o.Content[i])) {
			return false
		}
	}

	return n.span.Equal(o.span)
}

func (n *InterpolatedStringLiteralNode) String() string {
	var buff strings.Builder
	buff.WriteString("\"")

	for _, content := range n.Content {
		buff.WriteString(content.String())
	}

	buff.WriteString("\"")
	return buff.String()
}

func (*InterpolatedStringLiteralNode) IsStatic() bool {
	return false
}

func (*InterpolatedStringLiteralNode) Type(globalEnv *types.GlobalEnvironment) types.Type {
	return globalEnv.StdSubtype(symbol.String)
}

func (*InterpolatedStringLiteralNode) Class() *value.Class {
	return value.InterpolatedStringLiteralNodeClass
}

func (*InterpolatedStringLiteralNode) DirectClass() *value.Class {
	return value.InterpolatedStringLiteralNodeClass
}

func (n *InterpolatedStringLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::InterpolatedStringLiteralNode{\n  &: %p", n)

	buff.WriteString(",\n  content: %%[\n")
	for i, stmt := range n.Content {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, stmt.Inspect(), 2)
	}

	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *InterpolatedStringLiteralNode) Error() string {
	return n.Inspect()
}

// Create a new interpolated string literal node eg. `"foo ${bar} baz"`
func NewInterpolatedStringLiteralNode(span *position.Span, cont []StringLiteralContentNode) *InterpolatedStringLiteralNode {
	return &InterpolatedStringLiteralNode{
		NodeBase: NodeBase{span: span},
		Content:  cont,
	}
}

// Represents a simple double quoted string literal eg. `"foo baz"`
type DoubleQuotedStringLiteralNode struct {
	TypedNodeBase
	Value string
}

// Check if this node equals another node.
func (n *DoubleQuotedStringLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*DoubleQuotedStringLiteralNode)
	if !ok {
		return false
	}

	return n.Value == o.Value && n.span.Equal(o.span)
}

// Return a string representation of the node.
func (n *DoubleQuotedStringLiteralNode) String() string {
	return value.String(n.Value).Inspect()
}

func (*DoubleQuotedStringLiteralNode) IsStatic() bool {
	return true
}

func (*DoubleQuotedStringLiteralNode) Class() *value.Class {
	return value.DoubleQuotedStringLiteralNodeClass
}

func (*DoubleQuotedStringLiteralNode) DirectClass() *value.Class {
	return value.DoubleQuotedStringLiteralNodeClass
}

func (n *DoubleQuotedStringLiteralNode) Inspect() string {
	return fmt.Sprintf(
		"Std::Elk::AST::DoubleQuotedStringLiteralNode{&: %p, value: %s}",
		n,
		value.String(n.Value).Inspect(),
	)
}

func (n *DoubleQuotedStringLiteralNode) Error() string {
	return n.Inspect()
}

// Create a new double quoted string literal node eg. `"foo baz"`
func NewDoubleQuotedStringLiteralNode(span *position.Span, val string) *DoubleQuotedStringLiteralNode {
	return &DoubleQuotedStringLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}
