package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents all nodes that are valid in using declarations
type UsingEntryNode interface {
	Node
	ExpressionNode
	usingEntryNode()
}

func (*InvalidNode) usingEntryNode()                  {}
func (*PublicConstantNode) usingEntryNode()           {}
func (*PrivateConstantNode) usingEntryNode()          {}
func (*ConstantLookupNode) usingEntryNode()           {}
func (*MethodLookupNode) usingEntryNode()             {}
func (*UsingAllEntryNode) usingEntryNode()            {}
func (*UsingEntryWithSubentriesNode) usingEntryNode() {}
func (*ConstantAsNode) usingEntryNode()               {}
func (*MethodLookupAsNode) usingEntryNode()           {}
func (*GenericConstantNode) usingEntryNode()          {}
func (*NilLiteralNode) usingEntryNode()               {}
func (*UnquoteNode) usingEntryNode()                  {}
func (*MacroCallNode) usingEntryNode()                {}
func (*ReceiverlessMacroCallNode) usingEntryNode()    {}
func (*ScopedMacroCallNode) usingEntryNode()          {}

// Represents all nodes that are valid in using subentries
// in `UsingEntryWithSubentriesNode`
type UsingSubentryNode interface {
	Node
	ExpressionNode
	usingSubentryNode()
}

func (*InvalidNode) usingSubentryNode()          {}
func (*PublicConstantNode) usingSubentryNode()   {}
func (*PublicConstantAsNode) usingSubentryNode() {}
func (*MacroNameNode) usingSubentryNode()        {}
func (*PublicIdentifierNode) usingSubentryNode() {}
func (*UsingSubentryAsNode) usingSubentryNode()  {}

// Represents a using all entry node eg. `Foo::*`, `A::B::C::*`
type UsingAllEntryNode struct {
	TypedNodeBase
	Namespace UsingEntryNode
}

func (n *UsingAllEntryNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &UsingAllEntryNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Namespace:     n.Namespace.splice(loc, args, unquote).(UsingEntryNode),
	}
}

func (n *UsingAllEntryNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::UsingAllEntryNode", env)
}

func (n *UsingAllEntryNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Namespace.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (n *UsingAllEntryNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UsingAllEntryNode)
	if !ok {
		return false
	}

	return n.Namespace.Equal(value.Ref(o.Namespace)) &&
		n.loc.Equal(o.loc)
}

func (n *UsingAllEntryNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.Namespace.String())
	buff.WriteString("::*")

	return buff.String()
}

func (*UsingAllEntryNode) IsStatic() bool {
	return false
}

func (*UsingAllEntryNode) Class() *value.Class {
	return value.UsingAllEntryNodeClass
}

func (*UsingAllEntryNode) DirectClass() *value.Class {
	return value.UsingAllEntryNodeClass
}

func (n *UsingAllEntryNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::UsingAllEntryNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  namespace: ")
	indent.IndentStringFromSecondLine(&buff, n.Namespace.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *UsingAllEntryNode) Error() string {
	return n.Inspect()
}

// Create a new using all entry node eg. `Foo::*`, `A::B::C::*`
func NewUsingAllEntryNode(loc *position.Location, namespace UsingEntryNode) *UsingAllEntryNode {
	return &UsingAllEntryNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Namespace:     namespace,
	}
}

// Represents a using entry node with subentries eg. `Foo::{Bar, baz}`, `A::B::C::{lol, foo as epic, Gro as Moe}`
type UsingEntryWithSubentriesNode struct {
	NodeBase
	Namespace  UsingEntryNode
	Subentries []UsingSubentryNode
}

func (n *UsingEntryWithSubentriesNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &UsingEntryWithSubentriesNode{
		NodeBase:   NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Namespace:  n.Namespace.splice(loc, args, unquote).(UsingEntryNode),
		Subentries: SpliceSlice(n.Subentries, loc, args, unquote),
	}
}

func (n *UsingEntryWithSubentriesNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::UsingEntryWithSubentriesNode", env)
}

func (n *UsingEntryWithSubentriesNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Namespace.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	for _, entry := range n.Subentries {
		if entry.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

func (n *UsingEntryWithSubentriesNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UsingEntryWithSubentriesNode)
	if !ok {
		return false
	}

	if len(n.Subentries) != len(o.Subentries) ||
		!n.Namespace.Equal(value.Ref(o.Namespace)) ||
		!n.loc.Equal(o.loc) {
		return false
	}

	for i, subentry := range n.Subentries {
		if !subentry.Equal(value.Ref(o.Subentries[i])) {
			return false
		}
	}

	return true
}

func (n *UsingEntryWithSubentriesNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.Namespace.String())
	buff.WriteString("::{")

	for i, subentry := range n.Subentries {
		if i != 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(subentry.String())
	}

	buff.WriteRune('}')

	return buff.String()
}

func (*UsingEntryWithSubentriesNode) IsStatic() bool {
	return false
}

// Create a new using all entry node eg. `Foo::*`, `A::B::C::*`
func NewUsingEntryWithSubentriesNode(loc *position.Location, namespace UsingEntryNode, subentries []UsingSubentryNode) *UsingEntryWithSubentriesNode {
	return &UsingEntryWithSubentriesNode{
		NodeBase:   NodeBase{loc: loc},
		Namespace:  namespace,
		Subentries: subentries,
	}
}

func (*UsingEntryWithSubentriesNode) Class() *value.Class {
	return value.UsingEntryWithSubentriesNodeClass
}

func (*UsingEntryWithSubentriesNode) DirectClass() *value.Class {
	return value.UsingEntryWithSubentriesNodeClass
}

func (n *UsingEntryWithSubentriesNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::UsingEntryWithSubentriesNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  namespace: ")
	indent.IndentStringFromSecondLine(&buff, n.Namespace.Inspect(), 1)

	buff.WriteString(",\n  subentries: %[\n")
	for i, element := range n.Subentries {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *UsingEntryWithSubentriesNode) Error() string {
	return n.Inspect()
}

// Represents a using expression eg. `using Foo`
type UsingExpressionNode struct {
	TypedNodeBase
	Entries []UsingEntryNode
}

func (n *UsingExpressionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &UsingExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Entries:       SpliceSlice(n.Entries, loc, args, unquote),
	}
}

func (n *UsingExpressionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::UsingExpressionNode", env)
}

func (n *UsingExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	for _, entry := range n.Entries {
		if entry.traverse(n, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(n, parent)
}

func (n *UsingExpressionNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UsingExpressionNode)
	if !ok {
		return false
	}

	if len(n.Entries) != len(o.Entries) ||
		!n.loc.Equal(o.loc) {
		return false
	}

	for i, entry := range n.Entries {
		if !entry.Equal(value.Ref(o.Entries[i])) {
			return false
		}
	}

	return true
}

func (n *UsingExpressionNode) String() string {
	var buff strings.Builder

	buff.WriteString("using ")

	for i, entry := range n.Entries {
		if i != 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(entry.String())
	}

	return buff.String()
}

func (*UsingExpressionNode) SkipTypechecking() bool {
	return false
}

func (*UsingExpressionNode) IsStatic() bool {
	return false
}

func (*UsingExpressionNode) Class() *value.Class {
	return value.UsingExpressionNodeClass
}

func (*UsingExpressionNode) DirectClass() *value.Class {
	return value.UsingExpressionNodeClass
}

func (n *UsingExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::UsingExpressionNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  entries: %[\n")
	for i, element := range n.Entries {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *UsingExpressionNode) Error() string {
	return n.Inspect()
}

// Create a using expression node eg. `using Foo`
func NewUsingExpressionNode(loc *position.Location, consts []UsingEntryNode) *UsingExpressionNode {
	return &UsingExpressionNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Entries:       consts,
	}
}

// Represents an identifier with as in using declarations
// eg. `foo as bar`.
type UsingSubentryAsNode struct {
	NodeBase
	Target IdentifierNode
	AsName IdentifierNode
}

func (n *UsingSubentryAsNode) IsMacro() bool {
	_, ok := n.Target.(*MacroNameNode)
	return ok
}

func (n *UsingSubentryAsNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &UsingSubentryAsNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
		Target:   n.Target.splice(loc, args, unquote).(IdentifierNode),
		AsName:   n.AsName.splice(loc, args, unquote).(IdentifierNode),
	}
}

func (n *UsingSubentryAsNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::UsingSubentryAsNode", env)
}

func (n *UsingSubentryAsNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	if n.Target.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	if n.AsName.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (n *UsingSubentryAsNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UsingSubentryAsNode)
	if !ok {
		return false
	}

	return n.Target.Equal(value.Ref(o.Target)) &&
		n.AsName.Equal(value.Ref(o.AsName)) &&
		n.loc.Equal(o.loc)
}

func (n *UsingSubentryAsNode) String() string {
	return fmt.Sprintf("%s as %s", n.Target.String(), n.AsName.String())
}

func (*UsingSubentryAsNode) IsStatic() bool {
	return false
}

func (*UsingSubentryAsNode) Class() *value.Class {
	return value.UsingSubentryAsNodeClass
}

func (*UsingSubentryAsNode) DirectClass() *value.Class {
	return value.UsingSubentryAsNodeClass
}

func (n *UsingSubentryAsNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::UsingSubentryAsNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  target: ")
	indent.IndentStringFromSecondLine(&buff, n.Target.Inspect(), 1)

	buff.WriteString(",\n  as_name: ")
	indent.IndentStringFromSecondLine(&buff, n.AsName.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *UsingSubentryAsNode) Error() string {
	return n.Inspect()
}

// Create a new identifier with as eg. `foo as bar`.
func NewUsingSubentryAsNode(loc *position.Location, target IdentifierNode, as IdentifierNode) *UsingSubentryAsNode {
	return &UsingSubentryAsNode{
		NodeBase: NodeBase{loc: loc},
		Target:   target,
		AsName:   as,
	}
}

// Represents a macro name eg. `foo!`.
type MacroNameNode struct {
	TypedNodeBase
	Value string
}

func (n *MacroNameNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &MacroNameNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Value:         n.Value,
	}
}

func (n *MacroNameNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::MacroNameNode", env)
}

func (n *MacroNameNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n *MacroNameNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*MacroNameNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.loc.Equal(o.loc)
}

func (n *MacroNameNode) String() string {
	if IdentifierRegexp.MatchString(n.Value) {
		return n.Value
	}

	var buff strings.Builder
	buff.WriteByte('$')

	if PrefixedIdentifierRegexp.MatchString(n.Value) {
		buff.WriteString(n.Value)
		return buff.String()
	}

	buff.WriteString(value.String(n.Value).Inspect())
	buff.WriteByte('!')
	return buff.String()
}

func (*MacroNameNode) IsStatic() bool {
	return false
}

func (*MacroNameNode) Class() *value.Class {
	return value.MacroNameNodeClass
}

func (*MacroNameNode) DirectClass() *value.Class {
	return value.MacroNameNodeClass
}

func (n *MacroNameNode) Inspect() string {
	return fmt.Sprintf(
		"Std::Elk::AST::MacroNameNode{location: %s, value: %s}",
		(*value.Location)(n.loc).Inspect(),
		n.Value,
	)
}

func (n *MacroNameNode) Error() string {
	return n.Inspect()
}

// Create a new macro name node eg. `foo!`.
func NewMacroNameNode(loc *position.Location, val string) *MacroNameNode {
	return &MacroNameNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
	}
}
