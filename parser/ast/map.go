package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a HashMap literal eg. `{ foo: 1, 'bar' => 5, baz }`
type HashMapLiteralNode struct {
	TypedNodeBase
	Elements []ExpressionNode
	Capacity ExpressionNode
	static   bool
}

func (n *HashMapLiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	elements := SpliceSlice(n.Elements, loc, args, unquote)
	var capacity ExpressionNode
	var static bool

	if n.Capacity != nil {
		capacity = n.Capacity.splice(loc, args, unquote).(ExpressionNode)
		static = isExpressionSliceStatic(elements) && capacity.IsStatic()
	} else {
		static = isExpressionSliceStatic(elements)
	}

	return &HashMapLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

func (n *HashMapLiteralNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::HashMapLiteralNode", env)
}

func (n *HashMapLiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
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

	if n.Capacity != nil {
		if n.Capacity.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

// Check if this node equals another node.
func (n *HashMapLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*HashMapLiteralNode)
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

	return n.loc.Equal(o.loc)
}

// Return a string representation of the node.
func (n *HashMapLiteralNode) String() string {
	var buff strings.Builder

	buff.WriteRune('{')

	var hasMultilineArgs bool
	elementStrings := make([]string, len(n.Elements))

	for i, element := range n.Elements {
		elementString := element.String()
		elementStrings[i] = elementString
		if strings.ContainsRune(elementString, '\n') {
			hasMultilineArgs = true
		}
	}

	if hasMultilineArgs || len(n.Elements) > 6 {
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

	buff.WriteRune('}')

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

func (m *HashMapLiteralNode) IsStatic() bool {
	return m.static
}

// Create a HashMap literal node eg. `{ foo: 1, 'bar' => 5, baz }`
func NewHashMapLiteralNode(loc *position.Location, elements []ExpressionNode, capacity ExpressionNode) *HashMapLiteralNode {
	var static bool
	if capacity != nil {
		static = isExpressionSliceStatic(elements) && capacity.IsStatic()
	} else {
		static = isExpressionSliceStatic(elements)
	}
	return &HashMapLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewHashMapLiteralNode] but returns an interface
func NewHashMapLiteralNodeI(loc *position.Location, elements []ExpressionNode, capacity ExpressionNode) ExpressionNode {
	return NewHashMapLiteralNode(loc, elements, capacity)
}

func (*HashMapLiteralNode) Class() *value.Class {
	return value.HashMapLiteralNodeClass
}

func (*HashMapLiteralNode) DirectClass() *value.Class {
	return value.HashMapLiteralNodeClass
}

func (n *HashMapLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::HashMapLiteralNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  elements: %[\n")
	for i, element := range n.Elements {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  capacity: ")
	if n.Capacity == nil {
		buff.WriteString("nil")
	} else {
		indent.IndentStringFromSecondLine(&buff, n.Capacity.Inspect(), 1)
	}

	buff.WriteString("\n}")

	return buff.String()
}

func (n *HashMapLiteralNode) Error() string {
	return n.Inspect()
}

// Represents a Record literal eg. `%{ foo: 1, 'bar' => 5, baz }`
type HashRecordLiteralNode struct {
	TypedNodeBase
	Elements []ExpressionNode
	static   bool
}

func (n *HashRecordLiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	elements := SpliceSlice(n.Elements, loc, args, unquote)
	static := isExpressionSliceStatic(elements)

	return &HashRecordLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Elements:      elements,
		static:        static,
	}
}

func (n *HashRecordLiteralNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::HashRecordLiteralNode", env)
}

func (n *HashRecordLiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
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
func (n *HashRecordLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*HashRecordLiteralNode)
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
func (n *HashRecordLiteralNode) String() string {
	var buff strings.Builder

	buff.WriteString("%{")

	var hasMultilineArgs bool
	elementStrings := make([]string, len(n.Elements))

	for i, element := range n.Elements {
		elementString := element.String()
		elementStrings[i] = elementString
		if strings.ContainsRune(elementString, '\n') {
			hasMultilineArgs = true
		}
	}

	if hasMultilineArgs || len(n.Elements) > 6 {
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

	buff.WriteRune('}')

	return buff.String()
}

func (r *HashRecordLiteralNode) IsStatic() bool {
	return r.static
}

// Create a Record literal node eg. `%{ foo: 1, 'bar' => 5, baz }`
func NewHashRecordLiteralNode(loc *position.Location, elements []ExpressionNode) *HashRecordLiteralNode {
	return &HashRecordLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Elements:      elements,
		static:        isExpressionSliceStatic(elements),
	}
}

// Same as [NewHashRecordLiteralNode] but returns an interface
func NewHashRecordLiteralNodeI(loc *position.Location, elements []ExpressionNode) ExpressionNode {
	return NewHashRecordLiteralNode(loc, elements)
}

func (*HashRecordLiteralNode) Class() *value.Class {
	return value.HashRecordLiteralNodeClass
}

func (*HashRecordLiteralNode) DirectClass() *value.Class {
	return value.HashRecordLiteralNodeClass
}

func (n *HashRecordLiteralNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::HashRecordLiteralNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

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

func (n *HashRecordLiteralNode) Error() string {
	return n.Inspect()
}
