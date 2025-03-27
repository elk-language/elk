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
	return fmt.Sprintf("Std::Elk::AST::TrueLiteralNode{&: %p}", n)
}

func (n *TrueLiteralNode) Error() string {
	return n.Inspect()
}

// Create a new `true` literal node.
func NewTrueLiteralNode(span *position.Span) *TrueLiteralNode {
	return &TrueLiteralNode{
		NodeBase: NodeBase{span: span},
	}
}

// `false` literal.
type FalseLiteralNode struct {
	NodeBase
}

func (n *FalseLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*FalseLiteralNode)
	if !ok {
		return false
	}

	return n.Span().Equal(o.Span())
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
	return fmt.Sprintf("Std::Elk::AST::FalseLiteralNode{&: %p}", n)
}

func (n *FalseLiteralNode) Error() string {
	return n.Inspect()
}

// Create a new `false` literal node.
func NewFalseLiteralNode(span *position.Span) *FalseLiteralNode {
	return &FalseLiteralNode{
		NodeBase: NodeBase{span: span},
	}
}

// `self` literal.
type SelfLiteralNode struct {
	TypedNodeBase
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
	return fmt.Sprintf("Std::Elk::AST::SelfLiteralNode{&: %p}", n)
}

func (n *SelfLiteralNode) Error() string {
	return n.Inspect()
}

// Create a new `self` literal node.
func NewSelfLiteralNode(span *position.Span) *SelfLiteralNode {
	return &SelfLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
	}
}

// `nil` literal.
type NilLiteralNode struct {
	NodeBase
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
	return fmt.Sprintf("Std::Elk::AST::NilLiteralNode{&: %p}", n)
}

func (n *NilLiteralNode) Error() string {
	return n.Inspect()
}

// Create a new `nil` literal node.
func NewNilLiteralNode(span *position.Span) *NilLiteralNode {
	return &NilLiteralNode{
		NodeBase: NodeBase{span: span},
	}
}

// `undefined` literal.
type UndefinedLiteralNode struct {
	NodeBase
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
	return fmt.Sprintf("Std::Elk::AST::UndefinedLiteralNode{&: %p}", n)
}

func (n *UndefinedLiteralNode) Error() string {
	return n.Inspect()
}

// Create a new `undefined` literal node.
func NewUndefinedLiteralNode(span *position.Span) *UndefinedLiteralNode {
	return &UndefinedLiteralNode{
		NodeBase: NodeBase{span: span},
	}
}
