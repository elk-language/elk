package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a ArrayTuple literal eg. `%[1, 5, -6]`
type ArrayTupleLiteralNode struct {
	TypedNodeBase
	Elements []ExpressionNode
	static   bool
}

func (n *ArrayTupleLiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	elements := SpliceSlice(n.Elements, loc, args, unquote)

	return &ArrayTupleLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Elements:      elements,
		static:        isExpressionSliceStatic(elements),
	}
}

func (n *ArrayTupleLiteralNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::ArrayTupleLiteralNode", env)
}

func (n *ArrayTupleLiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	for _, elem := range n.Elements {
		if elem.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

func (t *ArrayTupleLiteralNode) IsStatic() bool {
	return t.static
}

// Create a ArrayTuple literal node eg. `%[1, 5, -6]`
func NewArrayTupleLiteralNode(loc *position.Location, elements []ExpressionNode) *ArrayTupleLiteralNode {
	return &ArrayTupleLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Elements:      elements,
		static:        isExpressionSliceStatic(elements),
	}
}

// Same as [NewArrayTupleLiteralNode] but returns an interface
func NewArrayTupleLiteralNodeI(loc *position.Location, elements []ExpressionNode) ExpressionNode {
	return NewArrayTupleLiteralNode(loc, elements)
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

	fmt.Fprintf(&buff, "Std::Elk::AST::ArrayTupleLiteralNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

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

func (n *WordArrayTupleLiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &WordArrayTupleLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Elements:      SpliceSlice(n.Elements, loc, args, unquote),
	}
}

func (n *WordArrayTupleLiteralNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::WordArrayTupleLiteralNode", env)
}

func (n *WordArrayTupleLiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	for _, elem := range n.Elements {
		if elem.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

func (n *WordArrayTupleLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*WordArrayTupleLiteralNode)
	if !ok {
		return false
	}

	if !n.loc.Equal(o.loc) ||
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
func NewWordArrayTupleLiteralNode(loc *position.Location, elements []WordCollectionContentNode) *WordArrayTupleLiteralNode {
	return &WordArrayTupleLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Elements:      elements,
	}
}

// Same as [NewWordArrayTupleLiteralNode] but returns an interface.
func NewWordArrayTupleLiteralExpressionNode(loc *position.Location, elements []WordCollectionContentNode) ExpressionNode {
	return NewWordArrayTupleLiteralNode(loc, elements)
}

// Same as [NewWordArrayTupleLiteralNode] but returns an interface.
func NewWordArrayTupleLiteralPatternExpressionNode(loc *position.Location, elements []WordCollectionContentNode) PatternExpressionNode {
	return NewWordArrayTupleLiteralNode(loc, elements)
}

func (*WordArrayTupleLiteralNode) Class() *value.Class {
	return value.WordArrayTupleLiteralNodeClass
}

func (*WordArrayTupleLiteralNode) DirectClass() *value.Class {
	return value.WordArrayTupleLiteralNodeClass
}

func (n *WordArrayTupleLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::WordArrayTupleLiteralNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

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

func (n *SymbolArrayTupleLiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &SymbolArrayTupleLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Elements:      SpliceSlice(n.Elements, loc, args, unquote),
	}
}

func (n *SymbolArrayTupleLiteralNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::SymbolArrayTupleLiteralNode", env)
}

func (n *SymbolArrayTupleLiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	for _, elem := range n.Elements {
		if elem.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

// Check if this node equals another node.
func (n *SymbolArrayTupleLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*SymbolArrayTupleLiteralNode)
	if !ok {
		return false
	}

	if len(n.Elements) != len(o.Elements) ||
		!n.loc.Equal(o.loc) {
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
func NewSymbolArrayTupleLiteralNode(loc *position.Location, elements []SymbolCollectionContentNode) *SymbolArrayTupleLiteralNode {
	return &SymbolArrayTupleLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Elements:      elements,
	}
}

// Same as [NewSymbolArrayTupleLiteralNode] but returns an interface.
func NewSymbolArrayTupleLiteralExpressionNode(loc *position.Location, elements []SymbolCollectionContentNode) ExpressionNode {
	return NewSymbolArrayTupleLiteralNode(loc, elements)
}

// Same as [NewSymbolArrayTupleLiteralNode] but returns an interface.
func NewSymbolArrayTupleLiteralPatternExpressionNode(loc *position.Location, elements []SymbolCollectionContentNode) PatternExpressionNode {
	return NewSymbolArrayTupleLiteralNode(loc, elements)
}

func (*SymbolArrayTupleLiteralNode) Class() *value.Class {
	return value.SymbolArrayTupleLiteralNodeClass
}

func (*SymbolArrayTupleLiteralNode) DirectClass() *value.Class {
	return value.SymbolArrayTupleLiteralNodeClass
}

func (n *SymbolArrayTupleLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::SymbolArrayTupleLiteralNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

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

func (n *HexArrayTupleLiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &HexArrayTupleLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Elements:      SpliceSlice(n.Elements, loc, args, unquote),
	}
}

func (n *HexArrayTupleLiteralNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::HexArrayTupleLiteralNode", env)
}

func (n *HexArrayTupleLiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	for _, elem := range n.Elements {
		if elem.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
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

	return n.loc.Equal(o.loc)
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
func NewHexArrayTupleLiteralNode(loc *position.Location, elements []IntCollectionContentNode) *HexArrayTupleLiteralNode {
	return &HexArrayTupleLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Elements:      elements,
	}
}

// Same as [NewHexArrayTupleLiteralNode] but returns an interface.
func NewHexArrayTupleLiteralExpressionNode(loc *position.Location, elements []IntCollectionContentNode) ExpressionNode {
	return NewHexArrayTupleLiteralNode(loc, elements)
}

// Same as [NewHexArrayTupleLiteralNode] but returns an interface.
func NewHexArrayTupleLiteralPatternExpressionNode(loc *position.Location, elements []IntCollectionContentNode) PatternExpressionNode {
	return NewHexArrayTupleLiteralNode(loc, elements)
}

func (*HexArrayTupleLiteralNode) Class() *value.Class {
	return value.HexArrayTupleLiteralNodeClass
}

func (*HexArrayTupleLiteralNode) DirectClass() *value.Class {
	return value.HexArrayTupleLiteralNodeClass
}

func (n *HexArrayTupleLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::HexArrayTupleLiteralNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

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

func (n *BinArrayTupleLiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &BinArrayTupleLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Elements:      SpliceSlice(n.Elements, loc, args, unquote),
	}
}

func (n *BinArrayTupleLiteralNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::BinArrayTupleLiteralNode", env)
}

func (n *BinArrayTupleLiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	for _, elem := range n.Elements {
		if elem.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
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
func NewBinArrayTupleLiteralNode(loc *position.Location, elements []IntCollectionContentNode) *BinArrayTupleLiteralNode {
	return &BinArrayTupleLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Elements:      elements,
	}
}

// Same as [NewBinArrayTupleLiteralNode] but returns an interface.
func NewBinArrayTupleLiteralExpressionNode(loc *position.Location, elements []IntCollectionContentNode) ExpressionNode {
	return NewBinArrayTupleLiteralNode(loc, elements)
}

// Same as [NewBinArrayTupleLiteralNode] but returns an interface.
func NewBinArrayTupleLiteralPatternExpressionNode(loc *position.Location, elements []IntCollectionContentNode) PatternExpressionNode {
	return NewBinArrayTupleLiteralNode(loc, elements)
}

func (*BinArrayTupleLiteralNode) Class() *value.Class {
	return value.BinArrayTupleLiteralNodeClass
}

func (*BinArrayTupleLiteralNode) DirectClass() *value.Class {
	return value.BinArrayTupleLiteralNodeClass
}

func (n *BinArrayTupleLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::BinArrayTupleLiteralNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

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

func (n *TuplePatternNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &TuplePatternNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Elements:      SpliceSlice(n.Elements, loc, args, unquote),
	}
}

func (n *TuplePatternNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::TuplePatternNode", env)
}

func (n *TuplePatternNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	for _, elem := range n.Elements {
		if elem.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

func (n *TuplePatternNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*TuplePatternNode)
	if !ok {
		return false
	}

	if len(n.Elements) != len(o.Elements) ||
		!n.loc.Equal(o.loc) {
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
func NewTuplePatternNode(loc *position.Location, elements []PatternNode) *TuplePatternNode {
	return &TuplePatternNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Elements:      elements,
	}
}

// Same as [NewTuplePatternNode] but returns an interface
func NewTuplePatternNodeI(loc *position.Location, elements []PatternNode) PatternNode {
	return NewTuplePatternNode(loc, elements)
}

func (*TuplePatternNode) Class() *value.Class {
	return value.TuplePatternNodeClass
}

func (*TuplePatternNode) DirectClass() *value.Class {
	return value.TuplePatternNodeClass
}

func (n *TuplePatternNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::TuplePatternNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

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
