package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Check whether the node is a constant.
func IsConstant(node Node) bool {
	switch node.(type) {
	case *PrivateConstantNode, *PublicConstantNode:
		return true
	default:
		return false
	}
}

// Check whether the node is a complex constant.
func IsComplexConstant(node Node) bool {
	switch node.(type) {
	case *PrivateConstantNode, *PublicConstantNode, *ConstantLookupNode:
		return true
	default:
		return false
	}
}

// All nodes that should be valid in constant lookups
// should implement this interface.
type ComplexConstantNode interface {
	Node
	TypeNode
	ExpressionNode
	PatternNode
	PatternExpressionNode
	UsingEntryNode
	complexConstantNode()
}

func (*InvalidNode) complexConstantNode()         {}
func (*PublicConstantNode) complexConstantNode()  {}
func (*PrivateConstantNode) complexConstantNode() {}
func (*ConstantLookupNode) complexConstantNode()  {}
func (*GenericConstantNode) complexConstantNode() {}
func (*NilLiteralNode) complexConstantNode()      {}
func (*UnquoteNode) complexConstantNode()         {}

// All nodes that should be valid constants
// should implement this interface.
type ConstantNode interface {
	Node
	TypeNode
	ExpressionNode
	UsingEntryNode
	ComplexConstantNode
	constantNode()
}

func (*InvalidNode) constantNode()         {}
func (*PublicConstantNode) constantNode()  {}
func (*PrivateConstantNode) constantNode() {}
func (*UnquoteNode) constantNode()         {}

// Represents a public constant eg. `Foo`.
type PublicConstantNode struct {
	TypedNodeBase
	Value string
}

func (n *PublicConstantNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &PublicConstantNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Value:         n.Value,
	}
}

func (n *PublicConstantNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::PublicConstantNode", env)
}

func (n *PublicConstantNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n *PublicConstantNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*PublicConstantNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.loc.Equal(o.loc)
}

func (n *PublicConstantNode) String() string {
	return n.Value
}

func (*PublicConstantNode) IsStatic() bool {
	return false
}

func (*PublicConstantNode) Class() *value.Class {
	return value.PublicIdentifierNodeClass
}

func (*PublicConstantNode) DirectClass() *value.Class {
	return value.PublicIdentifierNodeClass
}

func (n *PublicConstantNode) Inspect() string {
	return fmt.Sprintf(
		"Std::Elk::AST::PublicConstantNode{location: %s, value: %s}",
		(*value.Location)(n.loc).Inspect(),
		n.Value,
	)
}

func (n *PublicConstantNode) Error() string {
	return n.Inspect()
}

// Create a new public constant node eg. `Foo`.
func NewPublicConstantNode(loc *position.Location, val string) *PublicConstantNode {
	return &PublicConstantNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
	}
}

// Represents a private constant eg. `_Foo`
type PrivateConstantNode struct {
	TypedNodeBase
	Value string
}

func (n *PrivateConstantNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &PrivateConstantNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Value:         n.Value,
	}
}

func (n *PrivateConstantNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::PrivateConstantNode", env)
}

func (n *PrivateConstantNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n *PrivateConstantNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*PrivateConstantNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.loc.Equal(o.loc)
}

func (n *PrivateConstantNode) String() string {
	return n.Value
}

func (*PrivateConstantNode) IsStatic() bool {
	return false
}

func (*PrivateConstantNode) Class() *value.Class {
	return value.PublicIdentifierNodeClass
}

func (*PrivateConstantNode) DirectClass() *value.Class {
	return value.PublicIdentifierNodeClass
}

func (n *PrivateConstantNode) Inspect() string {
	return fmt.Sprintf(
		"Std::Elk::AST::PrivateConstantNode{location: %s, value: %s}",
		(*value.Location)(n.loc).Inspect(),
		n.Value,
	)
}

func (n *PrivateConstantNode) Error() string {
	return n.Inspect()
}

// Create a new private constant node eg. `_Foo`.
func NewPrivateConstantNode(loc *position.Location, val string) *PrivateConstantNode {
	return &PrivateConstantNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
	}
}

// Represents a constant with as in using declarations
// eg. `Foo as Bar`.
type PublicConstantAsNode struct {
	NodeBase
	Target *PublicConstantNode
	AsName string
}

func (n *PublicConstantAsNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &PublicConstantAsNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Target:   n.Target.splice(loc, args, unquote).(*PublicConstantNode),
		AsName:   n.AsName,
	}
}

func (n *PublicConstantAsNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::PublicConstantAsNode", env)
}

func (n *PublicConstantAsNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Target.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (n *PublicConstantAsNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*PublicConstantAsNode)
	if !ok {
		return false
	}

	return n.Target.Equal(value.Ref(o.Target)) &&
		n.AsName == o.AsName &&
		n.loc.Equal(o.loc)
}

func (n *PublicConstantAsNode) String() string {
	return fmt.Sprintf("%s as %s", n.Target.String(), n.AsName)
}

func (*PublicConstantAsNode) IsStatic() bool {
	return false
}

func (*PublicConstantAsNode) Class() *value.Class {
	return value.PublicConstantAsNodeClass
}

func (*PublicConstantAsNode) DirectClass() *value.Class {
	return value.PublicConstantAsNodeClass
}

func (n *PublicConstantAsNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::PublicConstantAsNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  target: ")
	indent.IndentStringFromSecondLine(&buff, n.Target.Inspect(), 1)

	buff.WriteString(",\n  as_name: ")
	buff.WriteString(n.AsName)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *PublicConstantAsNode) Error() string {
	return n.Inspect()
}

// Create a new identifier with as eg. `Foo as Bar`.
func NewPublicConstantAsNode(loc *position.Location, target *PublicConstantNode, as string) *PublicConstantAsNode {
	return &PublicConstantAsNode{
		NodeBase: NodeBase{loc: loc},
		Target:   target,
		AsName:   as,
	}
}

// Represents a constant lookup expressions eg. `Foo::Bar`
type ConstantLookupNode struct {
	TypedNodeBase
	Left  ExpressionNode      // left hand side
	Right ComplexConstantNode // right hand side
}

func (n *ConstantLookupNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	var left ExpressionNode
	if n.Left != nil {
		left = n.Left.splice(loc, args, unquote).(ExpressionNode)
	}
	right := n.Right.splice(loc, args, unquote).(ComplexConstantNode)

	return &ConstantLookupNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Left:          left,
		Right:         right,
	}
}

func (n *ConstantLookupNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::ConstantLookupNode", env)
}

func (n *ConstantLookupNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Left != nil {
		if n.Left.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	if n.Right.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

// Check if this node equals another node.
func (n *ConstantLookupNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ConstantLookupNode)
	if !ok {
		return false
	}

	if n.Left == o.Left {
	} else if n.Left == nil || o.Left == nil {
		return false
	} else if !n.Left.Equal(value.Ref(o.Left)) {
		return false
	}

	return n.Right.Equal(value.Ref(o.Right)) &&
		n.loc.Equal(o.loc)
}

// Return a string representation of the node.
func (n *ConstantLookupNode) String() string {
	var buff strings.Builder

	if n.Left != nil {
		buff.WriteString(n.Left.String())
	}
	buff.WriteString("::")
	buff.WriteString(n.Right.String())

	return buff.String()
}

func (*ConstantLookupNode) IsStatic() bool {
	return false
}

func (*ConstantLookupNode) Class() *value.Class {
	return value.ConstantLookupNodeClass
}

func (*ConstantLookupNode) DirectClass() *value.Class {
	return value.ConstantLookupNodeClass
}

func (n *ConstantLookupNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ConstantLookupNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  left: ")
	if n.Left == nil {
		buff.WriteString("nil")
	} else {
		indent.IndentStringFromSecondLine(&buff, n.Left.Inspect(), 1)
	}

	buff.WriteString(",\n  right: ")
	indent.IndentStringFromSecondLine(&buff, n.Right.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ConstantLookupNode) Error() string {
	return n.Inspect()
}

// Create a new constant lookup expression node eg. `Foo::Bar`
func NewConstantLookupNode(loc *position.Location, left ExpressionNode, right ComplexConstantNode) *ConstantLookupNode {
	return &ConstantLookupNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Left:          left,
		Right:         right,
	}
}

// Represents a generic constant in type annotations eg. `ArrayList[String]`
type GenericConstantNode struct {
	TypedNodeBase
	Constant      ComplexConstantNode
	TypeArguments []TypeNode
}

func (n *GenericConstantNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &GenericConstantNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Constant:      n.Constant.splice(loc, args, unquote).(ComplexConstantNode),
		TypeArguments: SpliceSlice(n.TypeArguments, loc, args, unquote),
	}
}

func (n *GenericConstantNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::GenericConstantNode", env)
}

func (n *GenericConstantNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Constant.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	for _, arg := range n.TypeArguments {
		if arg.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

// Equal checks if the given GenericConstantNode is equal to another value.
func (n *GenericConstantNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*GenericConstantNode)
	if !ok {
		return false
	}

	if !n.Constant.Equal(value.Ref(o.Constant)) {
		return false
	}

	if len(n.TypeArguments) != len(o.TypeArguments) {
		return false
	}

	for i, arg := range n.TypeArguments {
		if !arg.Equal(value.Ref(o.TypeArguments[i])) {
			return false
		}
	}

	return n.loc.Equal(o.loc)
}

// String returns a string representation of the GenericConstantNode.
func (n *GenericConstantNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.Constant.String())
	buff.WriteString("[")

	for i, arg := range n.TypeArguments {
		if i > 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(arg.String())
	}

	buff.WriteString("]")

	return buff.String()
}

func (*GenericConstantNode) IsStatic() bool {
	return true
}

func (*GenericConstantNode) Class() *value.Class {
	return value.GenericConstantNodeClass
}

func (*GenericConstantNode) DirectClass() *value.Class {
	return value.GenericConstantNodeClass
}

func (n *GenericConstantNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::GenericConstantNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  constant: ")
	indent.IndentStringFromSecondLine(&buff, n.Constant.Inspect(), 1)

	buff.WriteString(",\n  type_arguments: %[\n")
	for i, element := range n.TypeArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *GenericConstantNode) Error() string {
	return n.Inspect()
}

// Create a generic constant node eg. `ArrayList[String]`
func NewGenericConstantNode(loc *position.Location, constant ComplexConstantNode, args []TypeNode) *GenericConstantNode {
	return &GenericConstantNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Constant:      constant,
		TypeArguments: args,
	}
}
