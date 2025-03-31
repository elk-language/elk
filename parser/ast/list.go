package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a ArrayList literal eg. `[1, 5, -6]`
type ArrayListLiteralNode struct {
	TypedNodeBase
	Elements []ExpressionNode
	Capacity ExpressionNode
	static   bool
}

func (l *ArrayListLiteralNode) IsStatic() bool {
	return l.static
}

// Create a ArrayList literal node eg. `[1, 5, -6]`
func NewArrayListLiteralNode(span *position.Span, elements []ExpressionNode, capacity ExpressionNode) *ArrayListLiteralNode {
	var static bool
	if capacity != nil {
		static = isExpressionSliceStatic(elements) && capacity.IsStatic()
	} else {
		static = isExpressionSliceStatic(elements)
	}
	return &ArrayListLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewArrayListLiteralNode] but returns an interface
func NewArrayListLiteralNodeI(span *position.Span, elements []ExpressionNode, capacity ExpressionNode) ExpressionNode {
	return NewArrayListLiteralNode(span, elements, capacity)
}

func (*ArrayListLiteralNode) Class() *value.Class {
	return value.ArrayListLiteralNodeClass
}

func (*ArrayListLiteralNode) DirectClass() *value.Class {
	return value.ArrayListLiteralNodeClass
}

func (n *ArrayListLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ArrayListLiteralNode)
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

	if n.Capacity == o.Capacity {
		return true
	}

	if n.Capacity == nil || o.Capacity == nil {
		return false
	}

	return n.Capacity.Equal(value.Ref(o.Capacity))
}

func (n *ArrayListLiteralNode) String() string {
	var buff strings.Builder

	buff.WriteRune('[')

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

	if n.Capacity != nil {
		buff.WriteRune(':')

		parens := ExpressionPrecedence(n) > ExpressionPrecedence(n.Capacity)
		if parens {
			buff.WriteRune('(')
		}
		buff.WriteString(n.Capacity.String())
		if parens {
			buff.WriteRune(')')
		}
	}

	return buff.String()
}

func (n *ArrayListLiteralNode) Inspect() string {
	var buff strings.Builder

	buff.WriteString("Std::Elk::AST::ArrayListLiteralNode{\n")

	fmt.Fprintf(&buff, "span: %s", (*value.Span)(n.span).Inspect())
	buff.WriteString(",\n  elements: %%[\n")
	for i, element := range n.Elements {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  capacity: ")
	indent.IndentStringFromSecondLine(&buff, n.Capacity.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ArrayListLiteralNode) Error() string {
	return n.Inspect()
}

// Represents a word ArrayList literal eg. `\w[foo bar]`
type WordArrayListLiteralNode struct {
	TypedNodeBase
	Elements []WordCollectionContentNode
	Capacity ExpressionNode
	static   bool
}

func (n *WordArrayListLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*WordArrayListLiteralNode)
	if !ok {
		return false
	}

	if !n.span.Equal(o.span) ||
		len(n.Elements) != len(o.Elements) {
		return false
	}

	if n.Capacity == o.Capacity {
	} else if n.Capacity == nil || o.Capacity == nil {
		return false
	} else if !n.Capacity.Equal(value.Ref(o.Capacity)) {
		return false
	}

	for i, element := range n.Elements {
		if !element.Equal(value.Ref(o.Elements[i])) {
			return false
		}
	}

	return true
}

func (n *WordArrayListLiteralNode) String() string {
	var buff strings.Builder

	buff.WriteString("\\w[")

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

	if n.Capacity != nil {
		buff.WriteRune(':')

		parens := ExpressionPrecedence(n) > ExpressionPrecedence(n.Capacity)
		if parens {
			buff.WriteRune('(')
		}
		buff.WriteString(n.Capacity.String())
		if parens {
			buff.WriteRune(')')
		}
	}

	return buff.String()
}

func (w *WordArrayListLiteralNode) IsStatic() bool {
	return w.static
}

// Create a word ArrayList literal node eg. `\w[foo bar]`
func NewWordArrayListLiteralNode(span *position.Span, elements []WordCollectionContentNode, capacity ExpressionNode) *WordArrayListLiteralNode {
	var static bool
	if capacity != nil {
		static = capacity.IsStatic()
	} else {
		static = true
	}
	return &WordArrayListLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewWordArrayListLiteralNode] but returns an interface.
func NewWordArrayListLiteralExpressionNode(span *position.Span, elements []WordCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewWordArrayListLiteralNode(span, elements, capacity)
}

// Same as [NewWordArrayListLiteralNode] but returns an interface.
func NewWordArrayListLiteralPatternExpressionNode(span *position.Span, elements []WordCollectionContentNode) PatternExpressionNode {
	return NewWordArrayListLiteralNode(span, elements, nil)
}

func (*WordArrayListLiteralNode) Class() *value.Class {
	return value.WordArrayListLiteralNodeClass
}

func (*WordArrayListLiteralNode) DirectClass() *value.Class {
	return value.WordArrayListLiteralNodeClass
}

func (n *WordArrayListLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::WordArrayListLiteralNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  elements: %%[\n")
	for i, element := range n.Elements {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  capacity: ")
	indent.IndentStringFromSecondLine(&buff, n.Capacity.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *WordArrayListLiteralNode) Error() string {
	return n.Inspect()
}

// Represents a symbol ArrayList literal eg. `\s[foo bar]`
type SymbolArrayListLiteralNode struct {
	TypedNodeBase
	Elements []SymbolCollectionContentNode
	Capacity ExpressionNode
	static   bool
}

func (n *SymbolArrayListLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*SymbolArrayListLiteralNode)
	if !ok {
		return false
	}

	if len(n.Elements) != len(o.Elements) ||
		!n.span.Equal(o.span) {
		return false
	}

	if n.Capacity == o.Capacity {
	} else if n.Capacity == nil || o.Capacity == nil {
		return false
	} else if !n.Capacity.Equal(value.Ref(o.Capacity)) {
		return false
	}

	for i, element := range n.Elements {
		if !element.Equal(value.Ref(o.Elements[i])) {
			return false
		}
	}

	return true
}

func (n *SymbolArrayListLiteralNode) String() string {
	var buff strings.Builder

	buff.WriteString("\\s[")

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

	if n.Capacity != nil {
		buff.WriteRune(':')

		parens := ExpressionPrecedence(n) > ExpressionPrecedence(n.Capacity)
		if parens {
			buff.WriteRune('(')
		}
		buff.WriteString(n.Capacity.String())
		if parens {
			buff.WriteRune(')')
		}
	}

	return buff.String()
}

func (s *SymbolArrayListLiteralNode) IsStatic() bool {
	return s.static
}

// Create a symbol ArrayList literal node eg. `\s[foo bar]`
func NewSymbolArrayListLiteralNode(span *position.Span, elements []SymbolCollectionContentNode, capacity ExpressionNode) *SymbolArrayListLiteralNode {
	var static bool
	if capacity != nil {
		static = capacity.IsStatic()
	} else {
		static = true
	}
	return &SymbolArrayListLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewSymbolArrayListLiteralNode] but returns an interface.
func NewSymbolArrayListLiteralExpressionNode(span *position.Span, elements []SymbolCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewSymbolArrayListLiteralNode(span, elements, capacity)
}

// Same as [NewSymbolArrayListLiteralNode] but returns an interface.
func NewSymbolArrayListLiteralPatternExpressionNode(span *position.Span, elements []SymbolCollectionContentNode) PatternExpressionNode {
	return NewSymbolArrayListLiteralNode(span, elements, nil)
}

func (*SymbolArrayListLiteralNode) Class() *value.Class {
	return value.SymbolArrayListLiteralNodeClass
}

func (*SymbolArrayListLiteralNode) DirectClass() *value.Class {
	return value.SymbolArrayListLiteralNodeClass
}

func (n *SymbolArrayListLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::SymbolArrayListLiteralNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  elements: %%[\n")
	for i, element := range n.Elements {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  capacity: ")
	indent.IndentStringFromSecondLine(&buff, n.Capacity.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *SymbolArrayListLiteralNode) Error() string {
	return n.Inspect()
}

// Represents a hex ArrayList literal eg. `\x[ff ee]`
type HexArrayListLiteralNode struct {
	TypedNodeBase
	Elements []IntCollectionContentNode
	Capacity ExpressionNode
	static   bool
}

// Check if this node equals another node.
func (n *HexArrayListLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*HexArrayListLiteralNode)
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

	if n.Capacity == o.Capacity {
	} else if n.Capacity == nil || o.Capacity == nil {
		return false
	} else if !n.Capacity.Equal(value.Ref(o.Capacity)) {
		return false
	}

	return n.span.Equal(o.span)
}

// Return a string representation of the node.
func (n *HexArrayListLiteralNode) String() string {
	var buff strings.Builder

	buff.WriteString("\\x[")

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

	if n.Capacity != nil {
		buff.WriteRune(':')

		parens := ExpressionPrecedence(n) > ExpressionPrecedence(n.Capacity)
		if parens {
			buff.WriteRune('(')
		}
		buff.WriteString(n.Capacity.String())
		if parens {
			buff.WriteRune(')')
		}
	}

	return buff.String()
}

func (h *HexArrayListLiteralNode) IsStatic() bool {
	return h.static
}

// Create a hex ArrayList literal node eg. `\x[ff ee]`
func NewHexArrayListLiteralNode(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) *HexArrayListLiteralNode {
	var static bool
	if capacity != nil {
		static = capacity.IsStatic()
	} else {
		static = true
	}
	return &HexArrayListLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewHexArrayListLiteralNode] but returns an interface.
func NewHexArrayListLiteralExpressionNode(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewHexArrayListLiteralNode(span, elements, capacity)
}

// Same as [NewHexArrayListLiteralNode] but returns an interface.
func NewHexArrayListLiteralPatternExpressionNode(span *position.Span, elements []IntCollectionContentNode) PatternExpressionNode {
	return NewHexArrayListLiteralNode(span, elements, nil)
}

func (*HexArrayListLiteralNode) Class() *value.Class {
	return value.HexArrayListLiteralNodeClass
}

func (*HexArrayListLiteralNode) DirectClass() *value.Class {
	return value.HexArrayListLiteralNodeClass
}

func (n *HexArrayListLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::HexArrayListLiteralNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  elements: %%[\n")
	for i, element := range n.Elements {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  capacity: ")
	indent.IndentStringFromSecondLine(&buff, n.Capacity.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *HexArrayListLiteralNode) Error() string {
	return n.Inspect()
}

// Represents a bin ArrayList literal eg. `\b[11 10]`
type BinArrayListLiteralNode struct {
	TypedNodeBase
	Elements []IntCollectionContentNode
	Capacity ExpressionNode
	static   bool
}

func (b *BinArrayListLiteralNode) IsStatic() bool {
	return b.static
}
func (n *BinArrayListLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*BinArrayListLiteralNode)
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

	if n.Capacity == o.Capacity {
	} else if n.Capacity == nil || o.Capacity == nil {
		return false
	} else if !n.Capacity.Equal(value.Ref(o.Capacity)) {
		return false
	}

	return true
}

func (n *BinArrayListLiteralNode) String() string {
	var buff strings.Builder

	buff.WriteString("\\b[")
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

	if n.Capacity != nil {
		buff.WriteRune(':')

		parens := ExpressionPrecedence(n) > ExpressionPrecedence(n.Capacity)
		if parens {
			buff.WriteRune('(')
		}
		buff.WriteString(n.Capacity.String())
		if parens {
			buff.WriteRune(')')
		}
	}

	return buff.String()
}

// Create a bin ArrayList literal node eg. `\b[11 10]`
func NewBinArrayListLiteralNode(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) *BinArrayListLiteralNode {
	var static bool
	if capacity != nil {
		static = capacity.IsStatic()
	} else {
		static = true
	}
	return &BinArrayListLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewBinArrayListLiteralNode] but returns an interface.
func NewBinArrayListLiteralExpressionNode(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewBinArrayListLiteralNode(span, elements, capacity)
}

// Same as [NewBinArrayListLiteralNode] but returns an interface.
func NewBinArrayListLiteralPatternExpressionNode(span *position.Span, elements []IntCollectionContentNode) PatternExpressionNode {
	return NewBinArrayListLiteralNode(span, elements, nil)
}

func (*BinArrayListLiteralNode) Class() *value.Class {
	return value.BinArrayListLiteralNodeClass
}

func (*BinArrayListLiteralNode) DirectClass() *value.Class {
	return value.BinArrayListLiteralNodeClass
}

func (n *BinArrayListLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::BinArrayListLiteralNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  elements: %%[\n")
	for i, element := range n.Elements {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  capacity: ")
	indent.IndentStringFromSecondLine(&buff, n.Capacity.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *BinArrayListLiteralNode) Error() string {
	return n.Inspect()
}

// Represents a List pattern eg. `[1, a, >= 10]`
type ListPatternNode struct {
	TypedNodeBase
	Elements []PatternNode
}

func (n *ListPatternNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ListPatternNode)
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

func (n *ListPatternNode) String() string {
	var buff strings.Builder

	buff.WriteRune('[')
	for i, element := range n.Elements {
		if i != 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(element.String())
	}
	buff.WriteRune(']')

	return buff.String()
}

func (l *ListPatternNode) IsStatic() bool {
	return false
}

func (*ListPatternNode) Class() *value.Class {
	return value.ListPatternNodeClass
}

func (*ListPatternNode) DirectClass() *value.Class {
	return value.ListPatternNodeClass
}

func (n *ListPatternNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ListPatternNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  elements: %%[\n")
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

func (n *ListPatternNode) Error() string {
	return n.Inspect()
}

// Create a List pattern node eg. `[1, a, >= 10]`
func NewListPatternNode(span *position.Span, elements []PatternNode) *ListPatternNode {
	return &ListPatternNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

// Same as [NewListPatternNode] but returns an interface
func NewListPatternNodeI(span *position.Span, elements []PatternNode) PatternNode {
	return NewListPatternNode(span, elements)
}
