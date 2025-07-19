package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents an instance method lookup expression eg. `Foo.:bar`
type InstanceMethodLookupNode struct {
	TypedNodeBase
	Receiver ExpressionNode
	Name     IdentifierNode
}

func (n *InstanceMethodLookupNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &InstanceMethodLookupNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Receiver:      n.Receiver.splice(loc, args, unquote).(ExpressionNode),
		Name:          n.Name.splice(loc, args, unquote).(IdentifierNode),
	}
}

func (n *InstanceMethodLookupNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::InstanceMethodLookupNode", env)
}

func (n *InstanceMethodLookupNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Receiver.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	if n.Name.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (n *InstanceMethodLookupNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*InstanceMethodLookupNode)
	if !ok {
		return false
	}

	return n.Span().Equal(o.Span()) &&
		n.Receiver.Equal(value.Ref(o.Receiver)) &&
		n.Name.Equal(value.Ref(o.Name))
}

func (n *InstanceMethodLookupNode) String() string {
	var buff strings.Builder

	parens := ExpressionPrecedence(n) > ExpressionPrecedence(n.Receiver)

	if parens {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Receiver.String())
	if parens {
		buff.WriteRune(')')
	}

	buff.WriteString("::")
	buff.WriteString(n.Name.String())

	return buff.String()
}

func (*InstanceMethodLookupNode) IsStatic() bool {
	return false
}

// Create a new method lookup expression node eg. `Foo::bar`, `a::c`
func NewInstanceMethodLookupNode(loc *position.Location, receiver ExpressionNode, name IdentifierNode) *InstanceMethodLookupNode {
	return &InstanceMethodLookupNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Receiver:      receiver,
		Name:          name,
	}
}

func (*InstanceMethodLookupNode) Class() *value.Class {
	return value.InstanceMethodLookupNodeClass
}

func (*InstanceMethodLookupNode) DirectClass() *value.Class {
	return value.InstanceMethodLookupNodeClass
}

func (n *InstanceMethodLookupNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::InstanceMethodLookupNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  receiver: ")
	indent.IndentStringFromSecondLine(&buff, n.Receiver.Inspect(), 1)

	buff.WriteString(",\n  name: ")
	indent.IndentStringFromSecondLine(&buff, n.Name.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *InstanceMethodLookupNode) Error() string {
	return n.Inspect()
}

// Represents a method lookup expression eg. `Foo::bar`, `a::c`
type MethodLookupNode struct {
	TypedNodeBase
	Receiver ExpressionNode
	Name     IdentifierNode
}

func (n *MethodLookupNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &MethodLookupNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Receiver:      n.Receiver.splice(loc, args, unquote).(ExpressionNode),
		Name:          n.Name.splice(loc, args, unquote).(IdentifierNode),
	}
}

func (n *MethodLookupNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::MethodLookupNode", env)
}

func (n *MethodLookupNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Receiver.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	if n.Name.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (n *MethodLookupNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*MethodLookupNode)
	if !ok {
		return false
	}

	return n.Span().Equal(o.Span()) &&
		n.Receiver.Equal(value.Ref(o.Receiver)) &&
		n.Name.Equal(value.Ref(o.Name))
}

func (n *MethodLookupNode) String() string {
	var buff strings.Builder

	parens := ExpressionPrecedence(n) > ExpressionPrecedence(n.Receiver)

	if parens {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Receiver.String())
	if parens {
		buff.WriteRune(')')
	}

	buff.WriteString("::")
	buff.WriteString(n.Name.String())

	return buff.String()
}

func (*MethodLookupNode) IsStatic() bool {
	return false
}

// Create a new method lookup expression node eg. `Foo::bar`, `a::c`
func NewMethodLookupNode(loc *position.Location, receiver ExpressionNode, name IdentifierNode) *MethodLookupNode {
	return &MethodLookupNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Receiver:      receiver,
		Name:          name,
	}
}

func (*MethodLookupNode) Class() *value.Class {
	return value.MethodLookupNodeClass
}

func (*MethodLookupNode) DirectClass() *value.Class {
	return value.MethodLookupNodeClass
}

func (n *MethodLookupNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::MethodLookupNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  receiver: ")
	indent.IndentStringFromSecondLine(&buff, n.Receiver.Inspect(), 1)

	buff.WriteString(",\n  name: ")
	indent.IndentStringFromSecondLine(&buff, n.Name.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *MethodLookupNode) Error() string {
	return n.Inspect()
}

// Represents a method lookup with as in using declarations
// eg. `Foo::bar as baz`.
type MethodLookupAsNode struct {
	NodeBase
	MethodLookup *MethodLookupNode
	AsName       IdentifierNode
}

func (n *MethodLookupAsNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &MethodLookupAsNode{
		NodeBase:     NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		MethodLookup: n.MethodLookup.splice(loc, args, unquote).(*MethodLookupNode),
		AsName:       n.AsName.splice(loc, args, unquote).(IdentifierNode),
	}
}

func (n *MethodLookupAsNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::MethodLookupAsNode", env)
}

func (n *MethodLookupAsNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.MethodLookup.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	if n.AsName.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

// Check if this method lookup as node is equal to another value.
func (n *MethodLookupAsNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*MethodLookupAsNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.MethodLookup.Equal(value.Ref(o.MethodLookup)) &&
		n.AsName.Equal(value.Ref(o.AsName))
}

// Return a string representation of this method lookup as node.
func (n *MethodLookupAsNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.MethodLookup.String())
	buff.WriteString(" as ")
	buff.WriteString(n.AsName.String())

	return buff.String()
}

func (*MethodLookupAsNode) IsStatic() bool {
	return false
}

// Create a new identifier with as eg. `Foo::bar as Bar`.
func NewMethodLookupAsNode(loc *position.Location, methodLookup *MethodLookupNode, as IdentifierNode) *MethodLookupAsNode {
	return &MethodLookupAsNode{
		NodeBase:     NodeBase{loc: loc},
		MethodLookup: methodLookup,
		AsName:       as,
	}
}

func (*MethodLookupAsNode) Class() *value.Class {
	return value.ConstantAsNodeClass
}

func (*MethodLookupAsNode) DirectClass() *value.Class {
	return value.ConstantAsNodeClass
}

func (n *MethodLookupAsNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::MethodLookupAsNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  method_lookup: ")
	indent.IndentStringFromSecondLine(&buff, n.MethodLookup.Inspect(), 1)

	buff.WriteString(",\n  as_name: ")
	indent.IndentStringFromSecondLine(&buff, n.AsName.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *MethodLookupAsNode) Error() string {
	return n.Inspect()
}
