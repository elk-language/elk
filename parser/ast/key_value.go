package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a key value expression eg. `foo => bar`
type KeyValueExpressionNode struct {
	TypedNodeBase
	Key    ExpressionNode
	Value  ExpressionNode
	static bool
}

func (n *KeyValueExpressionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	key := n.Key.splice(loc, args, unquote).(ExpressionNode)
	val := n.Value.splice(loc, args, unquote).(ExpressionNode)

	return &KeyValueExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Key:           key,
		Value:         val,
		static:        areExpressionsStatic(key, val),
	}
}

func (n *KeyValueExpressionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::KeyValueExpressionNode", env)
}

func (n *KeyValueExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Key.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	if n.Value.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
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
	Key   IdentifierNode
	Value ExpressionNode
}

func (n *SymbolKeyValueExpressionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &SymbolKeyValueExpressionNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Key:      n.Key.splice(loc, args, unquote).(IdentifierNode),
		Value:    n.Value.splice(loc, args, unquote).(ExpressionNode),
	}
}

func (n *SymbolKeyValueExpressionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::SymbolKeyValueExpressionNode", env)
}

func (n *SymbolKeyValueExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Value.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
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

	buff.WriteString(n.Key.String())
	buff.WriteString(": ")
	buff.WriteString(n.Value.String())

	return buff.String()
}

func (s *SymbolKeyValueExpressionNode) IsStatic() bool {
	return s.Value.IsStatic()
}

// Create a symbol key value node eg. `foo: bar`
func NewSymbolKeyValueExpressionNode(loc *position.Location, key IdentifierNode, val ExpressionNode) *SymbolKeyValueExpressionNode {
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
	indent.IndentStringFromSecondLine(&buff, n.Key.Inspect(), 1)

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

func (n *SymbolKeyValuePatternNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &SymbolKeyValuePatternNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Key:      n.Key,
		Value:    n.Value.splice(loc, args, unquote).(PatternNode),
	}
}

func (n *SymbolKeyValuePatternNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::SymbolKeyValuePatternNode", env)
}

func (n *SymbolKeyValuePatternNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Value.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
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

func (n *KeyValuePatternNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &KeyValuePatternNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Key:      n.Key.splice(loc, args, unquote).(PatternExpressionNode),
		Value:    n.Value.splice(loc, args, unquote).(PatternNode),
	}
}

func (n *KeyValuePatternNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::KeyValuePatternNode", env)
}

func (n *KeyValuePatternNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Key.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	if n.Value.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
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
