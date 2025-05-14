package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a key value expression eg. `foo => bar`
type KeyValueExpressionNode struct {
	TypedNodeBase
	Key    ExpressionNode
	Value  ExpressionNode
	static bool
}

func (n *KeyValueExpressionNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	key := n.Key.Splice(loc, args, unquote).(ExpressionNode)
	val := n.Value.Splice(loc, args, unquote).(ExpressionNode)

	return &KeyValueExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Key:           key,
		Value:         val,
		static:        areExpressionsStatic(key, val),
	}
}

func (n *KeyValueExpressionNode) Traverse(yield func(Node) bool) bool {
	if !n.Key.Traverse(yield) {
		return false
	}
	if !n.Value.Traverse(yield) {
		return false
	}
	return yield(n)
}

func (n *KeyValueExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*KeyValueExpressionNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.Key.Equal(value.Ref(o.Key)) &&
		n.Value.Equal(value.Ref(o.Value))
}

func (n *KeyValueExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.Key.String())
	buff.WriteString(" => ")
	buff.WriteString(n.Value.String())

	return buff.String()
}

func (k *KeyValueExpressionNode) IsStatic() bool {
	return k.static
}

// Create a key value expression node eg. `foo => bar`
func NewKeyValueExpressionNode(loc *position.Location, key, val ExpressionNode) *KeyValueExpressionNode {
	return &KeyValueExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Key:           key,
		Value:         val,
		static:        areExpressionsStatic(key, val),
	}
}

func (*KeyValueExpressionNode) Class() *value.Class {
	return value.KeyValueExpressionNodeClass
}

func (*KeyValueExpressionNode) DirectClass() *value.Class {
	return value.KeyValueExpressionNodeClass
}

func (n *KeyValueExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::KeyValueExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  key: ")
	indent.IndentStringFromSecondLine(&buff, n.Key.Inspect(), 1)

	buff.WriteString(",\n  value: ")
	indent.IndentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *KeyValueExpressionNode) Error() string {
	return n.Inspect()
}

// Represents a symbol value expression eg. `foo: bar`
type SymbolKeyValueExpressionNode struct {
	NodeBase
	Key   string
	Value ExpressionNode
}

func (n *SymbolKeyValueExpressionNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &SymbolKeyValueExpressionNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Key:      n.Key,
		Value:    n.Value.Splice(loc, args, unquote).(ExpressionNode),
	}
}

func (n *SymbolKeyValueExpressionNode) Traverse(yield func(Node) bool) bool {
	if !n.Value.Traverse(yield) {
		return false
	}
	return yield(n)
}

func (n *SymbolKeyValueExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*SymbolKeyValueExpressionNode)
	if !ok {
		return false
	}

	return n.Key != o.Key &&
		n.Value.Equal(value.Ref(o.Value)) &&
		n.loc.Equal(o.loc)
}

func (n *SymbolKeyValueExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.Key)
	buff.WriteString(": ")
	buff.WriteString(n.Value.String())

	return buff.String()
}

func (s *SymbolKeyValueExpressionNode) IsStatic() bool {
	return s.Value.IsStatic()
}

// Create a symbol key value node eg. `foo: bar`
func NewSymbolKeyValueExpressionNode(loc *position.Location, key string, val ExpressionNode) *SymbolKeyValueExpressionNode {
	return &SymbolKeyValueExpressionNode{
		NodeBase: NodeBase{loc: loc},
		Key:      key,
		Value:    val,
	}
}

func (*SymbolKeyValueExpressionNode) Class() *value.Class {
	return value.SymbolKeyValueExpressionNodeClass
}

func (*SymbolKeyValueExpressionNode) DirectClass() *value.Class {
	return value.SymbolKeyValueExpressionNodeClass
}

func (n *SymbolKeyValueExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::SymbolKeyValueExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  key: ")
	indent.IndentStringFromSecondLine(&buff, value.String(n.Key).Inspect(), 1)

	buff.WriteString(",\n  value: ")
	indent.IndentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *SymbolKeyValueExpressionNode) Error() string {
	return n.Inspect()
}

// Represents a symbol value pattern eg. `foo: bar`
type SymbolKeyValuePatternNode struct {
	NodeBase
	Key   string
	Value PatternNode
}

func (n *SymbolKeyValuePatternNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &SymbolKeyValuePatternNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Key:      n.Key,
		Value:    n.Value.Splice(loc, args, unquote).(PatternNode),
	}
}

func (n *SymbolKeyValuePatternNode) Traverse(yield func(Node) bool) bool {
	if !n.Value.Traverse(yield) {
		return false
	}
	return yield(n)
}

func (n *SymbolKeyValuePatternNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*SymbolKeyValuePatternNode)
	if !ok {
		return false
	}

	return n.Key != o.Key &&
		n.Value.Equal(value.Ref(o.Value)) &&
		n.loc.Equal(o.loc)
}

func (n *SymbolKeyValuePatternNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.Key)
	buff.WriteString(": ")
	buff.WriteString(n.Value.String())

	return buff.String()
}

func (s *SymbolKeyValuePatternNode) IsStatic() bool {
	return false
}

func (*SymbolKeyValuePatternNode) Class() *value.Class {
	return value.SymbolKeyValuePatternNodeClass
}

func (*SymbolKeyValuePatternNode) DirectClass() *value.Class {
	return value.SymbolKeyValuePatternNodeClass
}

func (n *SymbolKeyValuePatternNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::SymbolKeyValuePatternNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  key: ")
	indent.IndentStringFromSecondLine(&buff, value.String(n.Key).Inspect(), 1)

	buff.WriteString(",\n  value: ")
	indent.IndentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *SymbolKeyValuePatternNode) Error() string {
	return n.Inspect()
}

// Create a symbol key value node eg. `foo: bar`
func NewSymbolKeyValuePatternNode(loc *position.Location, key string, val PatternNode) *SymbolKeyValuePatternNode {
	return &SymbolKeyValuePatternNode{
		NodeBase: NodeBase{loc: loc},
		Key:      key,
		Value:    val,
	}
}

// Represents a key value pattern eg. `foo => bar`
type KeyValuePatternNode struct {
	NodeBase
	Key   PatternExpressionNode
	Value PatternNode
}

func (n *KeyValuePatternNode) Splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &KeyValuePatternNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Key:      n.Key.Splice(loc, args, unquote).(PatternExpressionNode),
		Value:    n.Value.Splice(loc, args, unquote).(PatternNode),
	}
}

func (n *KeyValuePatternNode) Traverse(yield func(Node) bool) bool {
	if !n.Key.Traverse(yield) {
		return false
	}
	if !n.Value.Traverse(yield) {
		return false
	}
	return yield(n)
}

func (n *KeyValuePatternNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*KeyValuePatternNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.Key.Equal(value.Ref(o.Key)) &&
		n.Value.Equal(value.Ref(o.Value))
}

func (n *KeyValuePatternNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.Key.String())
	buff.WriteString(" => ")
	buff.WriteString(n.Value.String())

	return buff.String()
}

func (k *KeyValuePatternNode) IsStatic() bool {
	return false
}

func (*KeyValuePatternNode) Class() *value.Class {
	return value.KeyValuePatternNodeClass
}

func (*KeyValuePatternNode) DirectClass() *value.Class {
	return value.KeyValuePatternNodeClass
}

func (n *KeyValuePatternNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::SymbolKeyValuePatternNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  key: ")
	indent.IndentStringFromSecondLine(&buff, n.Key.Inspect(), 1)

	buff.WriteString(",\n  value: ")
	indent.IndentStringFromSecondLine(&buff, n.Value.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *KeyValuePatternNode) Error() string {
	return n.Inspect()
}

// Create a key value pattern node eg. `foo => bar`
func NewKeyValuePatternNode(loc *position.Location, key PatternExpressionNode, val PatternNode) *KeyValuePatternNode {
	return &KeyValuePatternNode{
		NodeBase: NodeBase{loc: loc},
		Key:      key,
		Value:    val,
	}
}
