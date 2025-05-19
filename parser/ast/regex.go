package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/regex/flag"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

// Nodes that implement this interface can appear
// inside of a Regex literal.
type RegexLiteralContentNode interface {
	Node
	regexLiteralContentNode()
}

func (*InvalidNode) regexLiteralContentNode()                    {}
func (*RegexInterpolationNode) regexLiteralContentNode()         {}
func (*RegexLiteralContentSectionNode) regexLiteralContentNode() {}

// All nodes that represent regexes should
// implement this interface.
type RegexLiteralNode interface {
	Node
	PatternExpressionNode
	regexLiteralNode()
}

func (*InvalidNode) regexLiteralNode()                    {}
func (*UninterpolatedRegexLiteralNode) regexLiteralNode() {}
func (*InterpolatedRegexLiteralNode) regexLiteralNode()   {}

// Represents an uninterpolated regex literal eg. `%/foo/`
type UninterpolatedRegexLiteralNode struct {
	NodeBase
	Content string
	Flags   bitfield.BitField8
}

func (n *UninterpolatedRegexLiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &UninterpolatedRegexLiteralNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Content:  n.Content,
		Flags:    n.Flags,
	}
}

func (n *UninterpolatedRegexLiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	}

	return leave(n, parent)
}

func (n *UninterpolatedRegexLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UninterpolatedRegexLiteralNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.Content == o.Content &&
		n.Flags == o.Flags
}

func (n *UninterpolatedRegexLiteralNode) String() string {
	var buff strings.Builder

	buff.WriteString("%/")
	buff.WriteString(n.Content)
	buff.WriteRune('/')

	if n.IsCaseInsensitive() {
		buff.WriteString("i")
	}
	if n.IsMultiline() {
		buff.WriteString("m")
	}
	if n.IsDotAll() {
		buff.WriteString("s")
	}
	if n.IsUngreedy() {
		buff.WriteString("U")
	}
	if n.IsASCII() {
		buff.WriteString("a")
	}
	if n.IsExtended() {
		buff.WriteString("x")
	}

	return buff.String()
}

func (*UninterpolatedRegexLiteralNode) Type(env *types.GlobalEnvironment) types.Type {
	return env.StdSubtype(symbol.Regex)
}

func (*UninterpolatedRegexLiteralNode) IsStatic() bool {
	return true
}

func (r *UninterpolatedRegexLiteralNode) IsCaseInsensitive() bool {
	return r.Flags.HasFlag(flag.CaseInsensitiveFlag)
}

func (r *UninterpolatedRegexLiteralNode) SetCaseInsensitive() *UninterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.CaseInsensitiveFlag)
	return r
}

func (r *UninterpolatedRegexLiteralNode) IsMultiline() bool {
	return r.Flags.HasFlag(flag.MultilineFlag)
}

func (r *UninterpolatedRegexLiteralNode) SetMultiline() *UninterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.MultilineFlag)
	return r
}

func (r *UninterpolatedRegexLiteralNode) IsDotAll() bool {
	return r.Flags.HasFlag(flag.DotAllFlag)
}

func (r *UninterpolatedRegexLiteralNode) SetDotAll() *UninterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.DotAllFlag)
	return r
}

func (r *UninterpolatedRegexLiteralNode) IsUngreedy() bool {
	return r.Flags.HasFlag(flag.UngreedyFlag)
}

func (r *UninterpolatedRegexLiteralNode) SetUngreedy() *UninterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.UngreedyFlag)
	return r
}

func (r *UninterpolatedRegexLiteralNode) IsASCII() bool {
	return r.Flags.HasFlag(flag.ASCIIFlag)
}

func (r *UninterpolatedRegexLiteralNode) SetASCII() *UninterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.ASCIIFlag)
	return r
}

func (r *UninterpolatedRegexLiteralNode) IsExtended() bool {
	return r.Flags.HasFlag(flag.ExtendedFlag)
}

func (r *UninterpolatedRegexLiteralNode) SetExtended() *UninterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.ExtendedFlag)
	return r
}

func (*UninterpolatedRegexLiteralNode) Class() *value.Class {
	return value.UninterpolatedRegexLiteralNodeClass
}

func (*UninterpolatedRegexLiteralNode) DirectClass() *value.Class {
	return value.UninterpolatedRegexLiteralNodeClass
}

func (n *UninterpolatedRegexLiteralNode) Inspect() string {
	return fmt.Sprintf(
		"Std::Elk::AST::UninterpolatedRegexLiteralNode{location: %s, content: %s, flags: %d}",
		(*value.Location)(n.loc).Inspect(),
		value.String(n.Content).Inspect(),
		n.Flags.Byte(),
	)
}

func (n *UninterpolatedRegexLiteralNode) Error() string {
	return n.Inspect()
}

// Create a new uninterpolated regex literal node eg. `%/foo/`.
func NewUninterpolatedRegexLiteralNode(loc *position.Location, content string, flags bitfield.BitField8) *UninterpolatedRegexLiteralNode {
	return &UninterpolatedRegexLiteralNode{
		NodeBase: NodeBase{loc: loc},
		Content:  content,
		Flags:    flags,
	}
}

// Represents a single section of characters of a regex literal eg. `foo` in `%/foo${bar}/`.
type RegexLiteralContentSectionNode struct {
	NodeBase
	Value string
}

func (n *RegexLiteralContentSectionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &RegexLiteralContentSectionNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Value:    n.Value,
	}
}
func (n *RegexLiteralContentSectionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	}

	return leave(n, parent)
}

func (n *RegexLiteralContentSectionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*RegexLiteralContentSectionNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.loc.Equal(o.loc)
}

func (n *RegexLiteralContentSectionNode) String() string {
	return n.Value
}

func (*RegexLiteralContentSectionNode) Class() *value.Class {
	return value.RegexLiteralContentSectionNodeClass
}

func (*RegexLiteralContentSectionNode) DirectClass() *value.Class {
	return value.RegexLiteralContentSectionNodeClass
}

func (n *RegexLiteralContentSectionNode) Inspect() string {
	return fmt.Sprintf(
		"Std::Elk::AST::RegexLiteralContentSectionNode{location: %s, value: %s}",
		(*value.Location)(n.loc).Inspect(),
		value.String(n.Value).Inspect(),
	)
}

func (n *RegexLiteralContentSectionNode) Error() string {
	return n.Inspect()
}

func (*RegexLiteralContentSectionNode) IsStatic() bool {
	return true
}

// Create a new regex literal content section node eg. `foo` in `%/foo${bar}/`.
func NewRegexLiteralContentSectionNode(loc *position.Location, val string) *RegexLiteralContentSectionNode {
	return &RegexLiteralContentSectionNode{
		NodeBase: NodeBase{loc: loc},
		Value:    val,
	}
}

// Represents a single interpolated section of a regex literal eg. `bar + 2` in `%/foo${bar + 2}/`
type RegexInterpolationNode struct {
	NodeBase
	Expression ExpressionNode
}

func (n *RegexInterpolationNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &RegexInterpolationNode{
		NodeBase:   NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Expression: n.Expression.splice(loc, args, unquote).(ExpressionNode),
	}
}

func (n *RegexInterpolationNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Expression.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (n *RegexInterpolationNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*RegexInterpolationNode)
	if !ok {
		return false
	}

	return n.Expression.Equal(value.Ref(o.Expression)) &&
		n.loc.Equal(o.loc)
}

func (n *RegexInterpolationNode) String() string {
	var buff strings.Builder

	buff.WriteString("${")
	buff.WriteString(n.Expression.String())
	buff.WriteString("}")

	return buff.String()
}

func (*RegexInterpolationNode) IsStatic() bool {
	return false
}

func (*RegexInterpolationNode) Class() *value.Class {
	return value.RegexInterpolationNodeClass
}

func (*RegexInterpolationNode) DirectClass() *value.Class {
	return value.RegexInterpolationNodeClass
}

func (n *RegexInterpolationNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::RegexInterpolationNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  expression: ")
	indent.IndentStringFromSecondLine(&buff, n.Expression.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *RegexInterpolationNode) Error() string {
	return n.Inspect()
}

// Create a new regex interpolation node eg. `bar + 2` in `%/foo${bar + 2}/`
func NewRegexInterpolationNode(loc *position.Location, expr ExpressionNode) *RegexInterpolationNode {
	return &RegexInterpolationNode{
		NodeBase:   NodeBase{loc: loc},
		Expression: expr,
	}
}

// Represents an Interpolated regex literal eg. `%/foo${1 + 2}bar/`
type InterpolatedRegexLiteralNode struct {
	NodeBase
	Content []RegexLiteralContentNode
	Flags   bitfield.BitField8
}

func (n *InterpolatedRegexLiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &InterpolatedRegexLiteralNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Content:  SpliceSlice(n.Content, loc, args, unquote),
		Flags:    n.Flags,
	}
}

func (n *InterpolatedRegexLiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	for _, content := range n.Content {
		if content.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

func (n *InterpolatedRegexLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*InterpolatedRegexLiteralNode)
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

	return n.loc.Equal(o.loc) && n.Flags == o.Flags
}

func (n *InterpolatedRegexLiteralNode) String() string {
	var buff strings.Builder
	buff.WriteString("%/")

	for _, content := range n.Content {
		buff.WriteString(content.String())
	}

	buff.WriteString("/")

	if n.IsCaseInsensitive() {
		buff.WriteString("i")
	}
	if n.IsMultiline() {
		buff.WriteString("m")
	}
	if n.IsDotAll() {
		buff.WriteString("s")
	}
	if n.IsUngreedy() {
		buff.WriteString("U")
	}
	if n.IsASCII() {
		buff.WriteString("a")
	}
	if n.IsExtended() {
		buff.WriteString("x")
	}

	return buff.String()
}

func (*InterpolatedRegexLiteralNode) Type(env *types.GlobalEnvironment) types.Type {
	return env.StdSubtype(symbol.Regex)
}

func (*InterpolatedRegexLiteralNode) IsStatic() bool {
	return false
}

func (r *InterpolatedRegexLiteralNode) IsCaseInsensitive() bool {
	return r.Flags.HasFlag(flag.CaseInsensitiveFlag)
}

func (r *InterpolatedRegexLiteralNode) SetCaseInsensitive() *InterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.CaseInsensitiveFlag)
	return r
}

func (r *InterpolatedRegexLiteralNode) IsMultiline() bool {
	return r.Flags.HasFlag(flag.MultilineFlag)
}

func (r *InterpolatedRegexLiteralNode) SetMultiline() *InterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.MultilineFlag)
	return r
}

func (r *InterpolatedRegexLiteralNode) IsDotAll() bool {
	return r.Flags.HasFlag(flag.DotAllFlag)
}

func (r *InterpolatedRegexLiteralNode) SetDotAll() *InterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.DotAllFlag)
	return r
}

func (r *InterpolatedRegexLiteralNode) IsUngreedy() bool {
	return r.Flags.HasFlag(flag.UngreedyFlag)
}

func (r *InterpolatedRegexLiteralNode) SetUngreedy() *InterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.UngreedyFlag)
	return r
}

func (r *InterpolatedRegexLiteralNode) IsASCII() bool {
	return r.Flags.HasFlag(flag.ASCIIFlag)
}

func (r *InterpolatedRegexLiteralNode) SetASCII() *InterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.ASCIIFlag)
	return r
}

func (r *InterpolatedRegexLiteralNode) IsExtended() bool {
	return r.Flags.HasFlag(flag.ExtendedFlag)
}

func (r *InterpolatedRegexLiteralNode) SetExtended() *InterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.ExtendedFlag)
	return r
}

func (*InterpolatedRegexLiteralNode) Class() *value.Class {
	return value.InterpolatedRegexLiteralNodeClass
}

func (*InterpolatedRegexLiteralNode) DirectClass() *value.Class {
	return value.InterpolatedRegexLiteralNodeClass
}

func (n *InterpolatedRegexLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::InterpolatedRegexLiteralNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  content: %[\n")
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

func (n *InterpolatedRegexLiteralNode) Error() string {
	return n.Inspect()
}

// Create a new interpolated regex literal node eg. `%/foo${1 + 3}bar/`.
func NewInterpolatedRegexLiteralNode(loc *position.Location, content []RegexLiteralContentNode, flags bitfield.BitField8) *InterpolatedRegexLiteralNode {
	return &InterpolatedRegexLiteralNode{
		NodeBase: NodeBase{loc: loc},
		Content:  content,
		Flags:    flags,
	}
}
