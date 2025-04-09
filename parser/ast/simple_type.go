package ast

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// `bool` literal.
type BoolLiteralNode struct {
	NodeBase
}

func (n *BoolLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*BoolLiteralNode)
	if !ok {
		return false
	}

	return n.loc.Equal(o.loc)
}

func (n *BoolLiteralNode) String() string {
	return "bool"
}

func (*BoolLiteralNode) IsStatic() bool {
	return true
}

func (*BoolLiteralNode) Type(globalEnv *types.GlobalEnvironment) types.Type {
	return types.Bool{}
}

func (*BoolLiteralNode) Class() *value.Class {
	return value.BoolLiteralNodeClass
}

func (*BoolLiteralNode) DirectClass() *value.Class {
	return value.BoolLiteralNodeClass
}

func (n *BoolLiteralNode) Inspect() string {
	return fmt.Sprintf("Std::Elk::AST::BoolLiteralNode{location: %s}", (*value.Location)(n.loc).Inspect())
}

func (n *BoolLiteralNode) Error() string {
	return n.Inspect()
}

// Create a new `bool` literal node.
func NewBoolLiteralNode(loc *position.Location) *BoolLiteralNode {
	return &BoolLiteralNode{
		NodeBase: NodeBase{loc: loc},
	}
}

// `void` type.
type VoidTypeNode struct {
	NodeBase
}

func (n *VoidTypeNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*VoidTypeNode)
	if !ok {
		return false
	}
	return n.loc.Equal(o.loc)
}

func (n *VoidTypeNode) String() string {
	return "void"
}

func (*VoidTypeNode) IsStatic() bool {
	return true
}

func (*VoidTypeNode) Class() *value.Class {
	return value.VoidTypeNodeClass
}

func (*VoidTypeNode) DirectClass() *value.Class {
	return value.VoidTypeNodeClass
}

func (n *VoidTypeNode) Inspect() string {
	return fmt.Sprintf("Std::Elk::AST::VoidTypeNode{location: %s}", (*value.Location)(n.loc).Inspect())
}

func (n *VoidTypeNode) Error() string {
	return n.Inspect()
}

// Create a new `void` type node.
func NewVoidTypeNode(loc *position.Location) *VoidTypeNode {
	return &VoidTypeNode{
		NodeBase: NodeBase{loc: loc},
	}
}

// `never` type.
type NeverTypeNode struct {
	NodeBase
}

func (n *NeverTypeNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*NeverTypeNode)
	if !ok {
		return false
	}
	return n.loc.Equal(o.loc)
}

func (n *NeverTypeNode) String() string {
	return "never"
}

func (*NeverTypeNode) IsStatic() bool {
	return true
}

func (*NeverTypeNode) Type(globalEnv *types.GlobalEnvironment) types.Type {
	return types.Never{}
}

func (*NeverTypeNode) Class() *value.Class {
	return value.VoidTypeNodeClass
}

func (*NeverTypeNode) DirectClass() *value.Class {
	return value.VoidTypeNodeClass
}

func (n *NeverTypeNode) Inspect() string {
	return fmt.Sprintf("Std::Elk::AST::NeverTypeNode{location: %s}", (*value.Location)(n.loc).Inspect())
}

func (n *NeverTypeNode) Error() string {
	return n.Inspect()
}

// Create a new `never` type node.
func NewNeverTypeNode(loc *position.Location) *NeverTypeNode {
	return &NeverTypeNode{
		NodeBase: NodeBase{loc: loc},
	}
}

// `any` type.
type AnyTypeNode struct {
	NodeBase
}

func (n *AnyTypeNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*AnyTypeNode)
	if !ok {
		return false
	}
	return n.loc.Equal(o.loc)
}

func (n *AnyTypeNode) String() string {
	return "any"
}

func (*AnyTypeNode) IsStatic() bool {
	return true
}

func (*AnyTypeNode) Type(globalEnv *types.GlobalEnvironment) types.Type {
	return types.Any{}
}

func (*AnyTypeNode) Class() *value.Class {
	return value.AnyTypeNodeClass
}

func (*AnyTypeNode) DirectClass() *value.Class {
	return value.AnyTypeNodeClass
}

func (n *AnyTypeNode) Inspect() string {
	return fmt.Sprintf("Std::Elk::AST::AnyTypeNode{location: %s}", (*value.Location)(n.loc).Inspect())
}

func (n *AnyTypeNode) Error() string {
	return n.Inspect()
}

// Create a new `any` type node.
func NewAnyTypeNode(loc *position.Location) *AnyTypeNode {
	return &AnyTypeNode{
		NodeBase: NodeBase{loc: loc},
	}
}
