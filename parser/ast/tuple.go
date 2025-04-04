package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a ArrayTuple literal eg. `%[1, 5, -6]`
type ArrayTupleLiteralNode struct {
	TypedNodeBase
	Elements []ExpressionNode
	static   bool
}

func (t *ArrayTupleLiteralNode) IsStatic() bool {
	return t.static
}

// Create a ArrayTuple literal node eg. `%[1, 5, -6]`
func NewArrayTupleLiteralNode(span *position.Span, elements []ExpressionNode) *ArrayTupleLiteralNode {
	return &ArrayTupleLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		static:        isExpressionSliceStatic(elements),
	}
}

// Same as [NewArrayTupleLiteralNode] but returns an interface
func NewArrayTupleLiteralNodeI(span *position.Span, elements []ExpressionNode) ExpressionNode {
	return NewArrayTupleLiteralNode(span, elements)
}

func (*ArrayTupleLiteralNode) Class() *value.Class {
	return value.ArrayTupleLiteralNodeClass
}

func (*ArrayTupleLiteralNode) DirectClass() *value.Class {
	return value.ArrayTupleLiteralNodeClass
}

func (n *ArrayTupleLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ArrayTupleLiteralNode)
	if !ok {
		return false
	}

	if !n.Span().Equal(o.Span()) {
		return false
	}

	if len(n.Elements) != len(o.Elements) {
		return false
	}

	for i, element := range n.Elements {
		if !element.Equal(value.Ref(o.Elements[i])) {
			return false
		}
	}

	return true
}

func (n *ArrayTupleLiteralNode) String() string {
	var buff strings.Builder

	buff.WriteString("%[")

	var hasMultilineElements bool
	elementStrings := make([]string, len(n.Elements))

	for i, element := range n.Elements {
		elementString := element.String()
		elementStrings[i] = elementString
		if strings.ContainsRune(elementString, '\n') {
			hasMultilineElements = true
		}
	}

	if hasMultilineElements || len(n.Elements) > 10 {
		buff.WriteRune('\n')
		for i, elementStr := range elementStrings {
			if i != 0 {
				buff.WriteString(",\n")
			}
			indent.IndentString(&buff, elementStr, 1)
		}
		buff.WriteRune('\n')
	} else {
		for i, elementStr := range elementStrings {
			if i != 0 {
				buff.WriteString(", ")
			}
			buff.WriteString(elementStr)
		}
	}
	buff.WriteRune(']')

	return buff.String()
}

func (n *ArrayTupleLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ArrayTupleLiteralNode{\n  span: %s", (*value.Span)(n.Span()).Inspect())

	buff.WriteString(",\n  elements: %[\n")
	for i, element := range n.Elements {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ArrayTupleLiteralNode) Error() string {
	return n.Inspect()
}

// Represents a word ArrayTuple literal eg. `%w[foo bar]`
type WordArrayTupleLiteralNode struct {
	TypedNodeBase
	Elements []WordCollectionContentNode
}

func (n *WordArrayTupleLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*WordArrayTupleLiteralNode)
	if !ok {
		return false
	}

	if !n.span.Equal(o.span) ||
		len(n.Elements) != len(o.Elements) {
		return false
	}

	for i, element := range n.Elements {
		if !element.Equal(value.Ref(o.Elements[i])) {
			return false
		}
	}

	return true
}

func (n *WordArrayTupleLiteralNode) String() string {
	var buff strings.Builder

	buff.WriteString("%w[")

	var i int
	for _, element := range n.Elements {
		element, ok := element.(*RawStringLiteralNode)
		if !ok {
			continue
		}
		if i != 0 {
			buff.WriteRune(' ')
		}
		buff.WriteString(element.Value)
		i++
	}

	buff.WriteRune(']')

	return buff.String()
}

func (*WordArrayTupleLiteralNode) IsStatic() bool {
	return true
}

// Create a word ArrayTuple literal node eg. `%w[foo bar]`
func NewWordArrayTupleLiteralNode(span *position.Span, elements []WordCollectionContentNode) *WordArrayTupleLiteralNode {
	return &WordArrayTupleLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

// Same as [NewWordArrayTupleLiteralNode] but returns an interface.
func NewWordArrayTupleLiteralExpressionNode(span *position.Span, elements []WordCollectionContentNode) ExpressionNode {
	return NewWordArrayTupleLiteralNode(span, elements)
}

// Same as [NewWordArrayTupleLiteralNode] but returns an interface.
func NewWordArrayTupleLiteralPatternExpressionNode(span *position.Span, elements []WordCollectionContentNode) PatternExpressionNode {
	return NewWordArrayTupleLiteralNode(span, elements)
}

func (*WordArrayTupleLiteralNode) Class() *value.Class {
	return value.WordArrayTupleLiteralNodeClass
}

func (*WordArrayTupleLiteralNode) DirectClass() *value.Class {
	return value.WordArrayTupleLiteralNodeClass
}

func (n *WordArrayTupleLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::WordArrayTupleLiteralNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  elements: %[\n")
	for i, element := range n.Elements {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *WordArrayTupleLiteralNode) Error() string {
	return n.Inspect()
}

// Represents a symbol ArrayTuple literal eg. `%s[foo bar]`
type SymbolArrayTupleLiteralNode struct {
	TypedNodeBase
	Elements []SymbolCollectionContentNode
}

// Check if this node equals another node.
func (n *SymbolArrayTupleLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*SymbolArrayTupleLiteralNode)
	if !ok {
		return false
	}

	if len(n.Elements) != len(o.Elements) ||
		!n.span.Equal(o.span) {
		return false
	}

	for i, element := range n.Elements {
		if !element.Equal(value.Ref(o.Elements[i])) {
			return false
		}
	}

	return true
}

// Return a string representation of the node.
func (n *SymbolArrayTupleLiteralNode) String() string {
	var buff strings.Builder

	buff.WriteString("%s[")

	var i int
	for _, element := range n.Elements {
		element, ok := element.(*SimpleSymbolLiteralNode)
		if !ok {
			continue
		}
		if i != 0 {
			buff.WriteString(" ")
		}
		buff.WriteString(element.Content)
		i++
	}

	buff.WriteRune(']')

	return buff.String()
}

func (*SymbolArrayTupleLiteralNode) IsStatic() bool {
	return true
}

// Create a symbol arrayTuple literal node eg. `%s[foo bar]`
func NewSymbolArrayTupleLiteralNode(span *position.Span, elements []SymbolCollectionContentNode) *SymbolArrayTupleLiteralNode {
	return &SymbolArrayTupleLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

// Same as [NewSymbolArrayTupleLiteralNode] but returns an interface.
func NewSymbolArrayTupleLiteralExpressionNode(span *position.Span, elements []SymbolCollectionContentNode) ExpressionNode {
	return NewSymbolArrayTupleLiteralNode(span, elements)
}

// Same as [NewSymbolArrayTupleLiteralNode] but returns an interface.
func NewSymbolArrayTupleLiteralPatternExpressionNode(span *position.Span, elements []SymbolCollectionContentNode) PatternExpressionNode {
	return NewSymbolArrayTupleLiteralNode(span, elements)
}

func (*SymbolArrayTupleLiteralNode) Class() *value.Class {
	return value.SymbolArrayTupleLiteralNodeClass
}

func (*SymbolArrayTupleLiteralNode) DirectClass() *value.Class {
	return value.SymbolArrayTupleLiteralNodeClass
}

func (n *SymbolArrayTupleLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::SymbolArrayTupleLiteralNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  elements: %[\n")
	for i, element := range n.Elements {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *SymbolArrayTupleLiteralNode) Error() string {
	return n.Inspect()
}

// Represents a hex ArrayTuple literal eg. `%x[ff ee]`
type HexArrayTupleLiteralNode struct {
	TypedNodeBase
	Elements []IntCollectionContentNode
}

// Check if this node equals another node.
func (n *HexArrayTupleLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*HexArrayTupleLiteralNode)
	if !ok {
		return false
	}

	if len(n.Elements) != len(o.Elements) {
		return false
	}

	for i, element := range n.Elements {
		if !element.Equal(value.Ref(o.Elements[i])) {
			return false
		}
	}

	return n.span.Equal(o.span)
}

// Return a string representation of the node.
func (n *HexArrayTupleLiteralNode) String() string {
	var buff strings.Builder

	buff.WriteString("%x[")

	var i int
	for _, element := range n.Elements {
		element, ok := element.(*IntLiteralNode)
		if !ok {
			continue
		}
		if i != 0 {
			buff.WriteRune(' ')
		}

		buff.WriteString(element.Value[2:]) // skip "0x"
		i++
	}
	buff.WriteRune(']')

	return buff.String()
}

func (*HexArrayTupleLiteralNode) IsStatic() bool {
	return true
}

// Create a hex ArrayTuple literal node eg. `%x[ff ee]`
func NewHexArrayTupleLiteralNode(span *position.Span, elements []IntCollectionContentNode) *HexArrayTupleLiteralNode {
	return &HexArrayTupleLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

// Same as [NewHexArrayTupleLiteralNode] but returns an interface.
func NewHexArrayTupleLiteralExpressionNode(span *position.Span, elements []IntCollectionContentNode) ExpressionNode {
	return NewHexArrayTupleLiteralNode(span, elements)
}

// Same as [NewHexArrayTupleLiteralNode] but returns an interface.
func NewHexArrayTupleLiteralPatternExpressionNode(span *position.Span, elements []IntCollectionContentNode) PatternExpressionNode {
	return NewHexArrayTupleLiteralNode(span, elements)
}

func (*HexArrayTupleLiteralNode) Class() *value.Class {
	return value.HexArrayTupleLiteralNodeClass
}

func (*HexArrayTupleLiteralNode) DirectClass() *value.Class {
	return value.HexArrayTupleLiteralNodeClass
}

func (n *HexArrayTupleLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::HexArrayTupleLiteralNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  elements: %[\n")
	for i, element := range n.Elements {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *HexArrayTupleLiteralNode) Error() string {
	return n.Inspect()
}

// Represents a bin ArrayTuple literal eg. `%b[11 10]`
type BinArrayTupleLiteralNode struct {
	TypedNodeBase
	Elements []IntCollectionContentNode
}

func (n *BinArrayTupleLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*BinArrayTupleLiteralNode)
	if !ok {
		return false
	}

	if !n.Span().Equal(o.Span()) {
		return false
	}

	if len(n.Elements) != len(o.Elements) {
		return false
	}

	for i, element := range n.Elements {
		if element == o.Elements[i] {
			continue
		}
		if !element.Equal(value.Ref(o.Elements[i])) {
			return false
		}
	}

	return true
}

func (n *BinArrayTupleLiteralNode) String() string {
	var buff strings.Builder

	buff.WriteString("%b[")
	var i int
	for _, element := range n.Elements {
		element, ok := element.(*IntLiteralNode)
		if !ok {
			continue
		}
		if i != 0 {
			buff.WriteRune(' ')
		}

		buff.WriteString(element.Value[2:]) // skip "0b"
		i++
	}
	buff.WriteRune(']')

	return buff.String()
}

func (*BinArrayTupleLiteralNode) IsStatic() bool {
	return true
}

// Create a bin ArrayList literal node eg. `%b[11 10]`
func NewBinArrayTupleLiteralNode(span *position.Span, elements []IntCollectionContentNode) *BinArrayTupleLiteralNode {
	return &BinArrayTupleLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

// Same as [NewBinArrayTupleLiteralNode] but returns an interface.
func NewBinArrayTupleLiteralExpressionNode(span *position.Span, elements []IntCollectionContentNode) ExpressionNode {
	return NewBinArrayTupleLiteralNode(span, elements)
}

// Same as [NewBinArrayTupleLiteralNode] but returns an interface.
func NewBinArrayTupleLiteralPatternExpressionNode(span *position.Span, elements []IntCollectionContentNode) PatternExpressionNode {
	return NewBinArrayTupleLiteralNode(span, elements)
}

func (*BinArrayTupleLiteralNode) Class() *value.Class {
	return value.BinArrayTupleLiteralNodeClass
}

func (*BinArrayTupleLiteralNode) DirectClass() *value.Class {
	return value.BinArrayTupleLiteralNodeClass
}

func (n *BinArrayTupleLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::BinArrayTupleLiteralNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  elements: %[\n")
	for i, element := range n.Elements {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *BinArrayTupleLiteralNode) Error() string {
	return n.Inspect()
}

// Represents a Tuple pattern eg. `%[1, a, >= 10]`
type TuplePatternNode struct {
	TypedNodeBase
	Elements []PatternNode
}

func (n *TuplePatternNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*TuplePatternNode)
	if !ok {
		return false
	}

	if len(n.Elements) != len(o.Elements) ||
		!n.span.Equal(o.span) {
		return false
	}

	for i, element := range n.Elements {
		if !element.Equal(value.Ref(o.Elements[i])) {
			return false
		}
	}

	return true
}

func (n *TuplePatternNode) String() string {
	var buff strings.Builder

	buff.WriteString("%[")
	for i, element := range n.Elements {
		if i != 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(element.String())
	}
	buff.WriteRune(']')

	return buff.String()
}

func (l *TuplePatternNode) IsStatic() bool {
	return false
}

// Create a Tuple pattern node eg. `%[1, a, >= 10]`
func NewTuplePatternNode(span *position.Span, elements []PatternNode) *TuplePatternNode {
	return &TuplePatternNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

// Same as [NewTuplePatternNode] but returns an interface
func NewTuplePatternNodeI(span *position.Span, elements []PatternNode) PatternNode {
	return NewTuplePatternNode(span, elements)
}

func (*TuplePatternNode) Class() *value.Class {
	return value.TuplePatternNodeClass
}

func (*TuplePatternNode) DirectClass() *value.Class {
	return value.TuplePatternNodeClass
}

func (n *TuplePatternNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::TuplePatternNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  elements: %[\n")
	for i, element := range n.Elements {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *TuplePatternNode) Error() string {
	return n.Inspect()
}
