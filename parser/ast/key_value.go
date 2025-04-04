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

func (n *KeyValueExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*KeyValueExpressionNode)
	if !ok {
		return false
	}

	return n.span.Equal(o.span) &&
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
func NewKeyValueExpressionNode(span *position.Span, key, val ExpressionNode) *KeyValueExpressionNode {
	return &KeyValueExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
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

	fmt.Fprintf(&buff, "Std::Elk::AST::KeyValueExpressionNode{\n  span: %s", (*value.Span)(n.span).Inspect())

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

func (n *SymbolKeyValueExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*SymbolKeyValueExpressionNode)
	if !ok {
		return false
	}

	return n.Key != o.Key &&
		n.Value.Equal(value.Ref(o.Value)) &&
		n.span.Equal(o.span)
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
func NewSymbolKeyValueExpressionNode(span *position.Span, key string, val ExpressionNode) *SymbolKeyValueExpressionNode {
	return &SymbolKeyValueExpressionNode{
		NodeBase: NodeBase{span: span},
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

	fmt.Fprintf(&buff, "Std::Elk::AST::SymbolKeyValueExpressionNode{\n  span: %s", (*value.Span)(n.span).Inspect())

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

func (n *SymbolKeyValuePatternNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*SymbolKeyValuePatternNode)
	if !ok {
		return false
	}

	return n.Key != o.Key &&
		n.Value.Equal(value.Ref(o.Value)) &&
		n.span.Equal(o.span)
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

	fmt.Fprintf(&buff, "Std::Elk::AST::SymbolKeyValuePatternNode{\n  span: %s", (*value.Span)(n.span).Inspect())

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
func NewSymbolKeyValuePatternNode(span *position.Span, key string, val PatternNode) *SymbolKeyValuePatternNode {
	return &SymbolKeyValuePatternNode{
		NodeBase: NodeBase{span: span},
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

func (n *KeyValuePatternNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*KeyValuePatternNode)
	if !ok {
		return false
	}

	return n.span.Equal(o.span) &&
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

	fmt.Fprintf(&buff, "Std::Elk::AST::SymbolKeyValuePatternNode{\n  span: %s", (*value.Span)(n.span).Inspect())

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
func NewKeyValuePatternNode(span *position.Span, key PatternExpressionNode, val PatternNode) *KeyValuePatternNode {
	return &KeyValuePatternNode{
		NodeBase: NodeBase{span: span},
		Key:      key,
		Value:    val,
	}
}
