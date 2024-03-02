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

// Represents a concatenation element like a quantifier, a char, a char class etc
type ConcatenationElementNode interface {
	Node
	concatenationElementNode()
}

func (*InvalidNode) concatenationElementNode()              {}
func (*ZeroOrOneQuantifierNode) concatenationElementNode()  {}
func (*ZeroOrMoreQuantifierNode) concatenationElementNode() {}
func (*OneOrMoreQuantifierNode) concatenationElementNode()  {}
func (*NQuantifierNode) concatenationElementNode()          {}
func (*NMQuantifierNode) concatenationElementNode()         {}

func (*MetaCharEscapeNode) concatenationElementNode()              {}
func (*GroupNode) concatenationElementNode()                       {}
func (*CharNode) concatenationElementNode()                        {}
func (*HexEscapeNode) concatenationElementNode()                   {}
func (*UnicodeCharClassNode) concatenationElementNode()            {}
func (*NegatedUnicodeCharClassNode) concatenationElementNode()     {}
func (*BellEscapeNode) concatenationElementNode()                  {}
func (*FormFeedEscapeNode) concatenationElementNode()              {}
func (*TabEscapeNode) concatenationElementNode()                   {}
func (*NewlineEscapeNode) concatenationElementNode()               {}
func (*CarriageReturnEscapeNode) concatenationElementNode()        {}
func (*VerticalTabEscapeNode) concatenationElementNode()           {}
func (*StartOfStringAnchorNode) concatenationElementNode()         {}
func (*EndOfStringAnchorNode) concatenationElementNode()           {}
func (*AbsoluteStartOfStringAnchorNode) concatenationElementNode() {}
func (*AbsoluteEndOfStringAnchorNode) concatenationElementNode()   {}
func (*WordBoundaryAnchorNode) concatenationElementNode()          {}
func (*NotWordBoundaryAnchorNode) concatenationElementNode()       {}
func (*WordCharClassNode) concatenationElementNode()               {}
func (*NotWordCharClassNode) concatenationElementNode()            {}
func (*DigitCharClassNode) concatenationElementNode()              {}
func (*NotDigitCharClassNode) concatenationElementNode()           {}
func (*WhitespaceCharClassNode) concatenationElementNode()         {}
func (*NotWhitespaceCharClassNode) concatenationElementNode()      {}
func (*AnyCharClassNode) concatenationElementNode()                {}

// Represents a primary regex element like a char, an escape, a char class, a group etc
type PrimaryRegexNode interface {
	Node
	ConcatenationElementNode
	primaryRegexNode()
}

func (*InvalidNode) primaryRegexNode()                     {}
func (*MetaCharEscapeNode) primaryRegexNode()              {}
func (*GroupNode) primaryRegexNode()                       {}
func (*CharNode) primaryRegexNode()                        {}
func (*HexEscapeNode) primaryRegexNode()                   {}
func (*UnicodeCharClassNode) primaryRegexNode()            {}
func (*NegatedUnicodeCharClassNode) primaryRegexNode()     {}
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
	Elements []ConcatenationElementNode
}

// Create a new concatenation node.
func NewConcatenationNode(span *position.Span, elements []ConcatenationElementNode) *ConcatenationNode {
	return &ConcatenationNode{
		NodeBase: NodeBase{span: span},
		Elements: elements,
	}
}

// Represents groups eg. `(foo)`, `(\w|\d)`
type GroupNode struct {
	NodeBase
	Regex Node
}

// Create a new group node.
func NewGroupNode(span *position.Span, regex Node) *GroupNode {
	return &GroupNode{
		NodeBase: NodeBase{span: span},
		Regex:    regex,
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

// Represents a zero or one quantifier eg. `f?`, `f??`
type ZeroOrOneQuantifierNode struct {
	NodeBase
	Regex Node
	Alt   bool
}

// Create a new zero or one quantifier node.
func NewZeroOrOneQuantifierNode(span *position.Span, regex Node, alt bool) *ZeroOrOneQuantifierNode {
	return &ZeroOrOneQuantifierNode{
		NodeBase: NodeBase{span: span},
		Regex:    regex,
		Alt:      alt,
	}
}

// Represents a one or more quantifier eg. `f*`, `f*?`
type ZeroOrMoreQuantifierNode struct {
	NodeBase
	Regex Node
	Alt   bool
}

// Create a new zero or more quantifier node.
func NewZeroOrMoreQuantifierNode(span *position.Span, regex Node, alt bool) *ZeroOrMoreQuantifierNode {
	return &ZeroOrMoreQuantifierNode{
		NodeBase: NodeBase{span: span},
		Regex:    regex,
		Alt:      alt,
	}
}

// Represents a one or more quantifier eg. `f+`, `f+?`
type OneOrMoreQuantifierNode struct {
	NodeBase
	Regex Node
	Alt   bool
}

// Create a new one or more quantifier node.
func NewOneOrMoreQuantifierNode(span *position.Span, regex Node, alt bool) *OneOrMoreQuantifierNode {
	return &OneOrMoreQuantifierNode{
		NodeBase: NodeBase{span: span},
		Regex:    regex,
		Alt:      alt,
	}
}

// Represents an N quantifier eg. `f{5}`, `\w{15}?`
type NQuantifierNode struct {
	NodeBase
	Regex Node
	N     string
}

// Create a new N quantifier node.
func NewNQuantifierNode(span *position.Span, regex Node, n string) *NQuantifierNode {
	return &NQuantifierNode{
		NodeBase: NodeBase{span: span},
		Regex:    regex,
		N:        n,
	}
}

// Represents an NM quantifier eg. `f{5,}`, `\w{15, 6}?`
type NMQuantifierNode struct {
	NodeBase
	Regex Node
	N     string
	M     string
	Alt   bool
}

// Create a new NM quantifier node.
func NewNMQuantifierNode(span *position.Span, regex Node, n, m string, alt bool) *NMQuantifierNode {
	return &NMQuantifierNode{
		NodeBase: NodeBase{span: span},
		Regex:    regex,
		N:        n,
		M:        m,
		Alt:      alt,
	}
}

// Represents a meta-char escape eg. `\\`, `\.`, `\+`
type MetaCharEscapeNode struct {
	NodeBase
	Value rune
}

// Create a new meta-char escape node.
func NewMetaCharEscapeNode(span *position.Span, char rune) *MetaCharEscapeNode {
	return &MetaCharEscapeNode{
		NodeBase: NodeBase{span: span},
		Value:    char,
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

// Represents a unicode char class eg. `\pL`, `\p{Latin}`
type UnicodeCharClassNode struct {
	NodeBase
	Value string
}

// Create a new unicode char class node.
func NewUnicodeCharClassNode(span *position.Span, value string) *UnicodeCharClassNode {
	return &UnicodeCharClassNode{
		NodeBase: NodeBase{span: span},
		Value:    value,
	}
}

// Represents a negated unicode char class eg. `\PL`, `\P{Latin}`, `\p{^Latin}`
type NegatedUnicodeCharClassNode struct {
	NodeBase
	Value string
}

// Create a new negated unicode char class node.
func NewNegatedUnicodeCharClassNode(span *position.Span, value string) *NegatedUnicodeCharClassNode {
	return &NegatedUnicodeCharClassNode{
		NodeBase: NodeBase{span: span},
		Value:    value,
	}
}

// Represents a hex escape eg. `\x20`, `\x{357}`
type HexEscapeNode struct {
	NodeBase
	Value string
}

// Create a new hex escape node.
func NewHexEscapeNode(span *position.Span, value string) *HexEscapeNode {
	return &HexEscapeNode{
		NodeBase: NodeBase{span: span},
		Value:    value,
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
