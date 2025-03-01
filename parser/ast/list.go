package ast

import (
	"fmt"
	"strings"

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

func (n *ArrayListLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::ArrayListLiteralNode{\n  &: %p", n)

	buff.WriteString(",\n  elements: %%[\n")
	for i, element := range n.Elements {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  capacity: ")
	indentStringFromSecondLine(&buff, n.Capacity.Inspect(), 1)

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

	fmt.Fprintf(&buff, "Std::AST::WordArrayListLiteralNode{\n  &: %p", n)

	buff.WriteString(",\n  elements: %%[\n")
	for i, element := range n.Elements {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  capacity: ")
	indentStringFromSecondLine(&buff, n.Capacity.Inspect(), 1)

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

	fmt.Fprintf(&buff, "Std::AST::SymbolArrayListLiteralNode{\n  &: %p", n)

	buff.WriteString(",\n  elements: %%[\n")
	for i, element := range n.Elements {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  capacity: ")
	indentStringFromSecondLine(&buff, n.Capacity.Inspect(), 1)

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

	fmt.Fprintf(&buff, "Std::AST::HexArrayListLiteralNode{\n  &: %p", n)

	buff.WriteString(",\n  elements: %%[\n")
	for i, element := range n.Elements {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  capacity: ")
	indentStringFromSecondLine(&buff, n.Capacity.Inspect(), 1)

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

	fmt.Fprintf(&buff, "Std::AST::BinArrayListLiteralNode{\n  &: %p", n)

	buff.WriteString(",\n  elements: %%[\n")
	for i, element := range n.Elements {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  capacity: ")
	indentStringFromSecondLine(&buff, n.Capacity.Inspect(), 1)

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

	fmt.Fprintf(&buff, "Std::AST::ListPatternNode{\n  &: %p", n)

	buff.WriteString(",\n  elements: %%[\n")
	for i, element := range n.Elements {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indentString(&buff, element.Inspect(), 2)
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
