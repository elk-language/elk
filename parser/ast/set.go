package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a HashSet literal eg. `^[1, 5, -6]`
type HashSetLiteralNode struct {
	TypedNodeBase
	Elements []ExpressionNode
	Capacity ExpressionNode
	static   bool
}

// Check if this node equals another node.
func (n *HashSetLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*HashSetLiteralNode)
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
func (n *HashSetLiteralNode) String() string {
	var buff strings.Builder

	buff.WriteString("^[")

	var hasMultilineArgs bool
	elementStrings := make([]string, len(n.Elements))

	for i, element := range n.Elements {
		elementString := element.String()
		elementStrings[i] = elementString
		if strings.ContainsRune(elementString, '\n') {
			hasMultilineArgs = true
		}
	}

	if hasMultilineArgs || len(n.Elements) > 10 {
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

func (s *HashSetLiteralNode) IsStatic() bool {
	return s.static
}

// Create a HashSet literal node eg. `^[1, 5, -6]`
func NewHashSetLiteralNode(span *position.Span, elements []ExpressionNode, capacity ExpressionNode) *HashSetLiteralNode {
	var static bool
	if capacity != nil {
		static = isExpressionSliceStatic(elements) && capacity.IsStatic()
	} else {
		static = isExpressionSliceStatic(elements)
	}
	return &HashSetLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewHashSetLiteralNode] but returns an interface
func NewHashSetLiteralNodeI(span *position.Span, elements []ExpressionNode, capacity ExpressionNode) ExpressionNode {
	return NewHashSetLiteralNode(span, elements, capacity)
}

func (*HashSetLiteralNode) Class() *value.Class {
	return value.HashSetLiteralNodeClass
}

func (*HashSetLiteralNode) DirectClass() *value.Class {
	return value.HashSetLiteralNodeClass
}

func (n *HashSetLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::HashSetLiteralNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  elements: %[\n")
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

func (n *HashSetLiteralNode) Error() string {
	return n.Inspect()
}

// Represents a word HashSet literal eg. `^w[foo bar]`
type WordHashSetLiteralNode struct {
	TypedNodeBase
	Elements []WordCollectionContentNode
	Capacity ExpressionNode
	static   bool
}

func (n *WordHashSetLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*WordHashSetLiteralNode)
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

func (n *WordHashSetLiteralNode) String() string {
	var buff strings.Builder

	buff.WriteString("^w[")

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

func (w *WordHashSetLiteralNode) IsStatic() bool {
	return w.static
}

// Create a word HashSet literal node eg. `^w[foo bar]`
func NewWordHashSetLiteralNode(span *position.Span, elements []WordCollectionContentNode, capacity ExpressionNode) *WordHashSetLiteralNode {
	var static bool
	if capacity != nil {
		static = capacity.IsStatic()
	} else {
		static = true
	}
	return &WordHashSetLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewWordHashSetLiteralNode] but returns an interface.
func NewWordHashSetLiteralNodeI(span *position.Span, elements []WordCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewWordHashSetLiteralNode(span, elements, capacity)
}

// Same as [NewWordHashSetLiteralNode] but returns an interface.
func NewWordHashSetLiteralPatternExpressionNode(span *position.Span, elements []WordCollectionContentNode) PatternExpressionNode {
	return NewWordHashSetLiteralNode(span, elements, nil)
}

func (*WordHashSetLiteralNode) Class() *value.Class {
	return value.WordHashSetLiteralNodeClass
}

func (*WordHashSetLiteralNode) DirectClass() *value.Class {
	return value.WordHashSetLiteralNodeClass
}

func (n *WordHashSetLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::WordHashSetLiteralNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  elements: %[\n")
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

func (n *WordHashSetLiteralNode) Error() string {
	return n.Inspect()
}

// Represents a symbol HashSet literal eg. `^s[foo bar]`
type SymbolHashSetLiteralNode struct {
	TypedNodeBase
	Elements []SymbolCollectionContentNode
	Capacity ExpressionNode
	static   bool
}

func (n *SymbolHashSetLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*SymbolHashSetLiteralNode)
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

func (n *SymbolHashSetLiteralNode) String() string {
	var buff strings.Builder

	buff.WriteString("^s[")

	var i int
	for _, element := range n.Elements {
		element, ok := element.(*SimpleSymbolLiteralNode)
		if !ok {
			continue
		}
		if i != 0 {
			buff.WriteRune(' ')
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

func (s *SymbolHashSetLiteralNode) IsStatic() bool {
	return s.static
}

// Create a symbol HashSet literal node eg. `^s[foo bar]`
func NewSymbolHashSetLiteralNode(span *position.Span, elements []SymbolCollectionContentNode, capacity ExpressionNode) *SymbolHashSetLiteralNode {
	var static bool
	if capacity != nil {
		static = capacity.IsStatic()
	} else {
		static = true
	}
	return &SymbolHashSetLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewSymbolHashSetLiteralNode] but returns an interface.
func NewSymbolHashSetLiteralNodeI(span *position.Span, elements []SymbolCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewSymbolHashSetLiteralNode(span, elements, capacity)
}

// Same as [NewSymbolHashSetLiteralNode] but returns an interface.
func NewSymbolHashSetLiteralPatternExpressionNode(span *position.Span, elements []SymbolCollectionContentNode) PatternExpressionNode {
	return NewSymbolHashSetLiteralNode(span, elements, nil)
}

func (*SymbolHashSetLiteralNode) Class() *value.Class {
	return value.SymbolHashSetLiteralNodeClass
}

func (*SymbolHashSetLiteralNode) DirectClass() *value.Class {
	return value.SymbolHashSetLiteralNodeClass
}

func (n *SymbolHashSetLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::SymbolHashSetLiteralNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  elements: %[\n")
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

func (n *SymbolHashSetLiteralNode) Error() string {
	return n.Inspect()
}

// Represents a hex HashSet literal eg. `^x[ff ee]`
type HexHashSetLiteralNode struct {
	TypedNodeBase
	Elements []IntCollectionContentNode
	Capacity ExpressionNode
	static   bool
}

func (h *HexHashSetLiteralNode) IsStatic() bool {
	return h.static
}

// Check if this node equals another node.
func (n *HexHashSetLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*HexHashSetLiteralNode)
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
func (n *HexHashSetLiteralNode) String() string {
	var buff strings.Builder

	buff.WriteString("^x[")

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

// Create a hex HashSet literal node eg. `^x[ff ee]`
func NewHexHashSetLiteralNode(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) *HexHashSetLiteralNode {
	var static bool
	if capacity != nil {
		static = capacity.IsStatic()
	} else {
		static = true
	}
	return &HexHashSetLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewHexHashSetLiteralNode] but returns an interface.
func NewHexHashSetLiteralNodeI(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewHexHashSetLiteralNode(span, elements, capacity)
}

// Same as [NewHexHashSetLiteralNode] but returns an interface.
func NewHexHashSetLiteralPatternExpressionNode(span *position.Span, elements []IntCollectionContentNode) PatternExpressionNode {
	return NewHexHashSetLiteralNode(span, elements, nil)
}

func (*HexHashSetLiteralNode) Class() *value.Class {
	return value.HexHashSetLiteralNodeClass
}

func (*HexHashSetLiteralNode) DirectClass() *value.Class {
	return value.HexHashSetLiteralNodeClass
}

func (n *HexHashSetLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::HexHashSetLiteralNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  elements: %[\n")
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

func (n *HexHashSetLiteralNode) Error() string {
	return n.Inspect()
}

// Represents a bin HashSet literal eg. `^b[11 10]`
type BinHashSetLiteralNode struct {
	TypedNodeBase
	Elements []IntCollectionContentNode
	Capacity ExpressionNode
	static   bool
}

func (n *BinHashSetLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*BinHashSetLiteralNode)
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

func (n *BinHashSetLiteralNode) String() string {
	var buff strings.Builder

	buff.WriteString("^b[")
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

func (b *BinHashSetLiteralNode) IsStatic() bool {
	return b.static
}

// Create a bin HashSet literal node eg. `^b[11 10]`
func NewBinHashSetLiteralNode(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) *BinHashSetLiteralNode {
	var static bool
	if capacity != nil {
		static = capacity.IsStatic()
	} else {
		static = true
	}
	return &BinHashSetLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewBinHashSetLiteralNode] but returns an interface.
func NewBinHashSetLiteralNodeI(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewBinHashSetLiteralNode(span, elements, capacity)
}

// Same as [NewBinHashSetLiteralNode] but returns an interface.
func NewBinHashSetLiteralPatternExpressionNode(span *position.Span, elements []IntCollectionContentNode) PatternExpressionNode {
	return NewBinHashSetLiteralNode(span, elements, nil)
}

func (*BinHashSetLiteralNode) Class() *value.Class {
	return value.BinHashSetLiteralNodeClass
}

func (*BinHashSetLiteralNode) DirectClass() *value.Class {
	return value.BinHashSetLiteralNodeClass
}

func (n *BinHashSetLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::BinHashSetLiteralNode{\n  span: %s", (*value.Span)(n.span).Inspect())

	buff.WriteString(",\n  elements: %[\n")
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

func (n *BinHashSetLiteralNode) Error() string {
	return n.Inspect()
}

// Represents a Set pattern eg. `^[1, "foo"]`
type SetPatternNode struct {
	TypedNodeBase
	Elements []PatternNode
}

func (n *SetPatternNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*SetPatternNode)
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

func (n *SetPatternNode) String() string {
	var buff strings.Builder
	buff.WriteString("^[")
	for i, element := range n.Elements {
		if i != 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(element.String())
	}
	buff.WriteRune(']')
	return buff.String()
}

func (s *SetPatternNode) IsStatic() bool {
	return false
}

func (*SetPatternNode) Class() *value.Class {
	return value.SetPatternNodeClass
}

func (*SetPatternNode) DirectClass() *value.Class {
	return value.SetPatternNodeClass
}

func (n *SetPatternNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::SetPatternNode{\n  span: %s", (*value.Span)(n.span).Inspect())

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

func (n *SetPatternNode) Error() string {
	return n.Inspect()
}

// Create a Set pattern node eg. `^[1, "foo"]`
func NewSetPatternNode(span *position.Span, elements []PatternNode) *SetPatternNode {
	return &SetPatternNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

// Same as [NewSetPatternNode] but returns an interface
func NewSetPatternNodeI(span *position.Span, elements []PatternNode) PatternNode {
	return NewSetPatternNode(span, elements)
}
