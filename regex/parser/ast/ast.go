// Package ast defines types
// used by the regex parser.
package ast

import (
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/regex/token"
)

// Every node type implements this interface.
type Node interface {
	position.SpanInterface
}

// Base struct of every AST node.
type NodeBase struct {
	span *position.Span
}

func (n *NodeBase) Span() *position.Span {
	return n.span
}

func (n *NodeBase) SetSpan(span *position.Span) {
	n.span = span
}

// Represents a primary regex element like a char, a char class, a group etc
type PrimaryRegexNode interface {
	Node
	primaryRegexNode()
}

func (*InvalidNode) primaryRegexNode()                     {}
func (*CharNode) primaryRegexNode()                        {}
func (*BellEscapeNode) primaryRegexNode()                  {}
func (*FormFeedEscapeNode) primaryRegexNode()              {}
func (*TabEscapeNode) primaryRegexNode()                   {}
func (*NewlineEscapeNode) primaryRegexNode()               {}
func (*CarriageReturnEscapeNode) primaryRegexNode()        {}
func (*VerticalTabEscapeNode) primaryRegexNode()           {}
func (*StartOfStringAnchorNode) primaryRegexNode()         {}
func (*EndOfStringAnchorNode) primaryRegexNode()           {}
func (*AbsoluteStartOfStringAnchorNode) primaryRegexNode() {}
func (*AbsoluteEndOfStringAnchorNode) primaryRegexNode()   {}
func (*WordBoundaryAnchorNode) primaryRegexNode()          {}
func (*NotWordBoundaryAnchorNode) primaryRegexNode()       {}
func (*WordCharClassNode) primaryRegexNode()               {}
func (*NotWordCharClassNode) primaryRegexNode()            {}
func (*DigitCharClassNode) primaryRegexNode()              {}
func (*NotDigitCharClassNode) primaryRegexNode()           {}
func (*WhitespaceCharClassNode) primaryRegexNode()         {}
func (*NotWhitespaceCharClassNode) primaryRegexNode()      {}
func (*AnyCharClassNode) primaryRegexNode()                {}

// Represents a syntax error.
type InvalidNode struct {
	NodeBase
	Token *token.Token
}

// Create a new invalid node.
func NewInvalidNode(span *position.Span, tok *token.Token) *InvalidNode {
	return &InvalidNode{
		NodeBase: NodeBase{span: span},
		Token:    tok,
	}
}

// Represents concatenated elements eg. `foo`, `\w-\d`
type ConcatenationNode struct {
	NodeBase
	Elements []PrimaryRegexNode
}

// Create a new concatenation node.
func NewConcatenationNode(span *position.Span, elements []PrimaryRegexNode) *ConcatenationNode {
	return &ConcatenationNode{
		NodeBase: NodeBase{span: span},
		Elements: elements,
	}
}

// Represents union eg. `foo|bar`, `\w|\d`
type UnionNode struct {
	NodeBase
	Left  Node
	Right Node
}

// Create a new union node.
func NewUnionNode(span *position.Span, left, right Node) *UnionNode {
	return &UnionNode{
		NodeBase: NodeBase{span: span},
		Left:     left,
		Right:    right,
	}
}

// Represents a char eg. `f`, `ę`
type CharNode struct {
	NodeBase
	Value rune
}

// Create a new char node.
func NewCharNode(span *position.Span, char rune) *CharNode {
	return &CharNode{
		NodeBase: NodeBase{span: span},
		Value:    char,
	}
}

// Represents a bell escape eg. `\a`
type BellEscapeNode struct {
	NodeBase
}

// Create a new bell escape node.
func NewBellEscapeNode(span *position.Span) *BellEscapeNode {
	return &BellEscapeNode{
		NodeBase: NodeBase{span: span},
	}
}

// Represents a form feed escape eg. `\f`
type FormFeedEscapeNode struct {
	NodeBase
}

// Create a new form feed escape node.
func NewFormFeedEscapeNode(span *position.Span) *FormFeedEscapeNode {
	return &FormFeedEscapeNode{
		NodeBase: NodeBase{span: span},
	}
}

// Represents a tab escape eg. `\t`
type TabEscapeNode struct {
	NodeBase
}

// Create a new tab escape node.
func NewTabEscapeNode(span *position.Span) *TabEscapeNode {
	return &TabEscapeNode{
		NodeBase: NodeBase{span: span},
	}
}

// Represents a newline escape eg. `\n`
type NewlineEscapeNode struct {
	NodeBase
}

// Create a new tab escape node.
func NewNewlineEscapeNode(span *position.Span) *NewlineEscapeNode {
	return &NewlineEscapeNode{
		NodeBase: NodeBase{span: span},
	}
}

// Represents a carriage return escape eg. `\r`
type CarriageReturnEscapeNode struct {
	NodeBase
}

// Create a new carriage return escape node.
func NewCarriageReturnEscapeNode(span *position.Span) *CarriageReturnEscapeNode {
	return &CarriageReturnEscapeNode{
		NodeBase: NodeBase{span: span},
	}
}

// Represents a vertical tab escape eg. `\v`
type VerticalTabEscapeNode struct {
	NodeBase
}

// Create a new vertical tab escape node.
func NewVerticalTabEscapeNode(span *position.Span) *VerticalTabEscapeNode {
	return &VerticalTabEscapeNode{
		NodeBase: NodeBase{span: span},
	}
}

// Represents the start of string anchor eg. `^`
type StartOfStringAnchorNode struct {
	NodeBase
}

// Create a new start of string anchor node.
func NewStartOfStringAnchorNode(span *position.Span) *StartOfStringAnchorNode {
	return &StartOfStringAnchorNode{
		NodeBase: NodeBase{span: span},
	}
}

// Represents the end of string anchor eg. `$`
type EndOfStringAnchorNode struct {
	NodeBase
}

// Create a new end of string anchor node.
func NewEndOfStringAnchorNode(span *position.Span) *EndOfStringAnchorNode {
	return &EndOfStringAnchorNode{
		NodeBase: NodeBase{span: span},
	}
}

// Represents the absolute start of text anchor eg. `\A`
type AbsoluteStartOfStringAnchorNode struct {
	NodeBase
}

// Create a new absolute start of text anchor node.
func NewAbsoluteStartOfStringAnchorNode(span *position.Span) *AbsoluteStartOfStringAnchorNode {
	return &AbsoluteStartOfStringAnchorNode{
		NodeBase: NodeBase{span: span},
	}
}

// Represents the absolute end of text anchor eg. `\z`
type AbsoluteEndOfStringAnchorNode struct {
	NodeBase
}

// Create a new absolute end of text anchor node.
func NewAbsoluteEndOfStringAnchorNode(span *position.Span) *AbsoluteEndOfStringAnchorNode {
	return &AbsoluteEndOfStringAnchorNode{
		NodeBase: NodeBase{span: span},
	}
}

// Represents a word boundary anchor eg. `\b`
type WordBoundaryAnchorNode struct {
	NodeBase
}

// Create a new word boundary anchor node.
func NewWordBoundaryAnchorNode(span *position.Span) *WordBoundaryAnchorNode {
	return &WordBoundaryAnchorNode{
		NodeBase: NodeBase{span: span},
	}
}

// Represents a not word boundary anchor eg. `\B`
type NotWordBoundaryAnchorNode struct {
	NodeBase
}

// Create a new not word boundary anchor node.
func NewNotWordBoundaryAnchorNode(span *position.Span) *NotWordBoundaryAnchorNode {
	return &NotWordBoundaryAnchorNode{
		NodeBase: NodeBase{span: span},
	}
}

// Represents a word char class eg. `\w`
type WordCharClassNode struct {
	NodeBase
}

// Create a new word char class node.
func NewWordCharClassNode(span *position.Span) *WordCharClassNode {
	return &WordCharClassNode{
		NodeBase: NodeBase{span: span},
	}
}

// Represents a not word char class eg. `\W`
type NotWordCharClassNode struct {
	NodeBase
}

// Create a new not word char class node.
func NewNotWordCharClassNode(span *position.Span) *NotWordCharClassNode {
	return &NotWordCharClassNode{
		NodeBase: NodeBase{span: span},
	}
}

// Represents a digit char class eg. `\d`
type DigitCharClassNode struct {
	NodeBase
}

// Create a new digit char class node.
func NewDigitCharClassNode(span *position.Span) *DigitCharClassNode {
	return &DigitCharClassNode{
		NodeBase: NodeBase{span: span},
	}
}

// Represents a not digit char class eg. `\D`
type NotDigitCharClassNode struct {
	NodeBase
}

// Create a new not digit char class node.
func NewNotDigitCharClassNode(span *position.Span) *NotDigitCharClassNode {
	return &NotDigitCharClassNode{
		NodeBase: NodeBase{span: span},
	}
}

// Represents a whitespace char class eg. `\s`
type WhitespaceCharClassNode struct {
	NodeBase
}

// Create a new whitespace char class node.
func NewWhitespaceCharClassNode(span *position.Span) *WhitespaceCharClassNode {
	return &WhitespaceCharClassNode{
		NodeBase: NodeBase{span: span},
	}
}

// Represents a not whitespace char class eg. `\S`
type NotWhitespaceCharClassNode struct {
	NodeBase
}

// Create a new not whitespace char class node.
func NewNotWhitespaceCharClassNode(span *position.Span) *NotWhitespaceCharClassNode {
	return &NotWhitespaceCharClassNode{
		NodeBase: NodeBase{span: span},
	}
}

// Represents the any char class eg. `.`
type AnyCharClassNode struct {
	NodeBase
}

// Create a new any char class node.
func NewAnyCharClassNode(span *position.Span) *AnyCharClassNode {
	return &AnyCharClassNode{
		NodeBase: NodeBase{span: span},
	}
}
