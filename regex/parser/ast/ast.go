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

func (*InvalidNode) primaryRegexNode()    {}
func (*CharNode) primaryRegexNode()       {}
func (*BellEscapeNode) primaryRegexNode() {}

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

// Represents a beginning of text anchor eg. `\A`
type BeginningOfTextAnchorNode struct {
	NodeBase
}

// Create a new beginning of text anchor node.
func NewBeginningOfTextAnchorNode(span *position.Span) *BeginningOfTextAnchorNode {
	return &BeginningOfTextAnchorNode{
		NodeBase: NodeBase{span: span},
	}
}

// Represents an end of text anchor eg. `\z`
type EndOfTextAnchorNode struct {
	NodeBase
}

// Create a new end of text anchor node.
func NewEndOfTextAnchorNode(span *position.Span) *EndOfTextAnchorNode {
	return &EndOfTextAnchorNode{
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
