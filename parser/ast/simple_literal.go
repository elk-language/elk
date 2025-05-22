package ast

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// `true` literal.
type TrueLiteralNode struct {
	NodeBase
}

func (n *TrueLiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &TrueLiteralNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
	}
}

func (n *TrueLiteralNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::TrueLiteralNode", env)
}

func (n *TrueLiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n *TrueLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*TrueLiteralNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc)
}

func (n *TrueLiteralNode) String() string {
	return "true"
}

func (*TrueLiteralNode) IsStatic() bool {
	return true
}

func (*TrueLiteralNode) Type(globalEnv *types.GlobalEnvironment) types.Type {
	return types.True{}
}

func (*TrueLiteralNode) Class() *value.Class {
	return value.TrueLiteralNodeClass
}

func (*TrueLiteralNode) DirectClass() *value.Class {
	return value.TrueLiteralNodeClass
}

func (n *TrueLiteralNode) Inspect() string {
	return fmt.Sprintf("Std::Elk::AST::TrueLiteralNode{location: %s}", (*value.Location)(n.loc).Inspect())
}

func (n *TrueLiteralNode) Error() string {
	return n.Inspect()
}

// Create a new `true` literal node.
func NewTrueLiteralNode(loc *position.Location) *TrueLiteralNode {
	return &TrueLiteralNode{
		NodeBase: NodeBase{loc: loc},
	}
}

// `false` literal.
type FalseLiteralNode struct {
	NodeBase
}

func (n *FalseLiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &FalseLiteralNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
	}
}

func (n *FalseLiteralNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::FalseLiteralNode", env)
}

func (n *FalseLiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n *FalseLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*FalseLiteralNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc)
}

func (n *FalseLiteralNode) String() string {
	return "false"
}

func (*FalseLiteralNode) IsStatic() bool {
	return true
}

func (*FalseLiteralNode) Type(globalEnv *types.GlobalEnvironment) types.Type {
	return types.False{}
}

func (*FalseLiteralNode) Class() *value.Class {
	return value.FalseLiteralNodeClass
}

func (*FalseLiteralNode) DirectClass() *value.Class {
	return value.FalseLiteralNodeClass
}

func (n *FalseLiteralNode) Inspect() string {
	return fmt.Sprintf("Std::Elk::AST::FalseLiteralNode{location: %s}", (*value.Location)(n.loc).Inspect())
}

func (n *FalseLiteralNode) Error() string {
	return n.Inspect()
}

// Create a new `false` literal node.
func NewFalseLiteralNode(loc *position.Location) *FalseLiteralNode {
	return &FalseLiteralNode{
		NodeBase: NodeBase{loc: loc},
	}
}

// `self` literal.
type SelfLiteralNode struct {
	TypedNodeBase
}

func (n *SelfLiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &SelfLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
	}
}

func (n *SelfLiteralNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::SelfLiteralNode", env)
}

func (n *SelfLiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n *SelfLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*SelfLiteralNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc)
}

func (n *SelfLiteralNode) String() string {
	return "self"
}

func (*SelfLiteralNode) IsStatic() bool {
	return false
}

func (*SelfLiteralNode) Class() *value.Class {
	return value.SelfLiteralNodeClass
}

func (*SelfLiteralNode) DirectClass() *value.Class {
	return value.SelfLiteralNodeClass
}

func (n *SelfLiteralNode) Inspect() string {
	return fmt.Sprintf("Std::Elk::AST::SelfLiteralNode{location: %s}", (*value.Location)(n.loc).Inspect())
}

func (n *SelfLiteralNode) Error() string {
	return n.Inspect()
}

// Create a new `self` literal node.
func NewSelfLiteralNode(loc *position.Location) *SelfLiteralNode {
	return &SelfLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
	}
}

// `nil` literal.
type NilLiteralNode struct {
	NodeBase
}

func (n *NilLiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &NilLiteralNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
	}
}

func (n *NilLiteralNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::NilLiteralNode", env)
}

func (n *NilLiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n *NilLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*NilLiteralNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc)
}

func (n *NilLiteralNode) String() string {
	return "nil"
}

func (*NilLiteralNode) SetType(types.Type) {}

func (*NilLiteralNode) IsStatic() bool {
	return true
}

func (*NilLiteralNode) Type(globalEnv *types.GlobalEnvironment) types.Type {
	return types.Nil{}
}

func (*NilLiteralNode) Class() *value.Class {
	return value.NilLiteralNodeClass
}

func (*NilLiteralNode) DirectClass() *value.Class {
	return value.NilLiteralNodeClass
}

func (n *NilLiteralNode) Inspect() string {
	return fmt.Sprintf("Std::Elk::AST::NilLiteralNode{location: %s}", (*value.Location)(n.loc).Inspect())
}

func (n *NilLiteralNode) Error() string {
	return n.Inspect()
}

// Create a new `nil` literal node.
func NewNilLiteralNode(loc *position.Location) *NilLiteralNode {
	return &NilLiteralNode{
		NodeBase: NodeBase{loc: loc},
	}
}

// `undefined` literal.
type UndefinedLiteralNode struct {
	NodeBase
}

func (n *UndefinedLiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &UndefinedLiteralNode{
		NodeBase: NodeBase{loc: position.SpliceLocation(loc, n.loc, unquote)},
	}
}

func (n *UndefinedLiteralNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::UndefinedLiteralNode", env)
}

func (n *UndefinedLiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n *UndefinedLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UndefinedLiteralNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc)
}

func (n *UndefinedLiteralNode) String() string {
	return "undefined"
}

func (*UndefinedLiteralNode) IsStatic() bool {
	return true
}

func (*UndefinedLiteralNode) Class() *value.Class {
	return value.UndefinedLiteralNodeClass
}

func (*UndefinedLiteralNode) DirectClass() *value.Class {
	return value.UndefinedLiteralNodeClass
}

func (n *UndefinedLiteralNode) Inspect() string {
	return fmt.Sprintf("Std::Elk::AST::UndefinedLiteralNode{location: %s}", (*value.Location)(n.loc).Inspect())
}

func (n *UndefinedLiteralNode) Error() string {
	return n.Inspect()
}

// Create a new `undefined` literal node.
func NewUndefinedLiteralNode(loc *position.Location) *UndefinedLiteralNode {
	return &UndefinedLiteralNode{
		NodeBase: NodeBase{loc: loc},
	}
}
