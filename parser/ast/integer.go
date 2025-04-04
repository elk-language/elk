package ast

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Int literal eg. `5`, `125_355`, `0xff`
type IntLiteralNode struct {
	TypedNodeBase
	Value string
}

func (n *IntLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*IntLiteralNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.span.Equal(o.span)
}

func (n *IntLiteralNode) String() string {
	return n.Value
}

func (*IntLiteralNode) IsStatic() bool {
	return true
}

func (*IntLiteralNode) Class() *value.Class {
	return value.IntLiteralNodeClass
}

func (*IntLiteralNode) DirectClass() *value.Class {
	return value.IntLiteralNodeClass
}

func (n *IntLiteralNode) Inspect() string {
	return fmt.Sprintf("Std::Elk::AST::IntLiteralNode{span: %s, value: %s}", (*value.Span)(n.span).Inspect(), n.Value)
}

func (n *IntLiteralNode) Error() string {
	return n.Inspect()
}

// Create a new int literal node eg. `5`, `125_355`, `0xff`
func NewIntLiteralNode(span *position.Span, val string) *IntLiteralNode {
	return &IntLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Int64 literal eg. `5i64`, `125_355i64`, `0xffi64`
type Int64LiteralNode struct {
	TypedNodeBase
	Value string
}

func (n *Int64LiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*Int64LiteralNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.span.Equal(o.span)
}

func (n *Int64LiteralNode) String() string {
	return fmt.Sprintf("%si64", n.Value)
}

func (*Int64LiteralNode) IsStatic() bool {
	return true
}

func (*Int64LiteralNode) Class() *value.Class {
	return value.Int64LiteralNodeClass
}

func (*Int64LiteralNode) DirectClass() *value.Class {
	return value.Int64LiteralNodeClass
}

func (n *Int64LiteralNode) Inspect() string {
	return fmt.Sprintf("Std::Elk::AST::Int64LiteralNode{span: %s, value: %s}", (*value.Span)(n.span).Inspect(), n.Value)
}

func (n *Int64LiteralNode) Error() string {
	return n.Inspect()
}

// Create a new Int64 literal node eg. `5i64`, `125_355i64`, `0xffi64`
func NewInt64LiteralNode(span *position.Span, val string) *Int64LiteralNode {
	return &Int64LiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// UInt64 literal eg. `5u64`, `125_355u64`, `0xffu64`
type UInt64LiteralNode struct {
	TypedNodeBase
	Value string
}

func (n *UInt64LiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UInt64LiteralNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.span.Equal(o.span)
}

func (n *UInt64LiteralNode) String() string {
	return fmt.Sprintf("%su64", n.Value)
}

func (*UInt64LiteralNode) IsStatic() bool {
	return true
}

func (*UInt64LiteralNode) Class() *value.Class {
	return value.UInt64LiteralNodeClass
}

func (*UInt64LiteralNode) DirectClass() *value.Class {
	return value.UInt64LiteralNodeClass
}

func (n *UInt64LiteralNode) Inspect() string {
	return fmt.Sprintf("Std::Elk::AST::UInt64LiteralNode{span: %s, value: %s}", (*value.Span)(n.span).Inspect(), n.Value)
}

func (n *UInt64LiteralNode) Error() string {
	return n.Inspect()
}

// Create a new UInt64 literal node eg. `5u64`, `125_355u64`, `0xffu64`
func NewUInt64LiteralNode(span *position.Span, val string) *UInt64LiteralNode {
	return &UInt64LiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Int32 literal eg. `5i32`, `1_20i32`, `0xffi32`
type Int32LiteralNode struct {
	TypedNodeBase
	Value string
}

func (n *Int32LiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*Int32LiteralNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.span.Equal(o.span)
}

func (n *Int32LiteralNode) String() string {
	return fmt.Sprintf("%si32", n.Value)
}

func (*Int32LiteralNode) IsStatic() bool {
	return true
}

func (*Int32LiteralNode) Class() *value.Class {
	return value.Int32LiteralNodeClass
}

func (*Int32LiteralNode) DirectClass() *value.Class {
	return value.Int32LiteralNodeClass
}

func (n *Int32LiteralNode) Inspect() string {
	return fmt.Sprintf("Std::Elk::AST::Int32LiteralNode{span: %s, value: %s}", (*value.Span)(n.span).Inspect(), n.Value)
}

func (n *Int32LiteralNode) Error() string {
	return n.Inspect()
}

// Create a new Int32 literal node eg. `5i32`, `1_20i32`, `0xffi32`
func NewInt32LiteralNode(span *position.Span, val string) *Int32LiteralNode {
	return &Int32LiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// UInt32 literal eg. `5u32`, `1_20u32`, `0xffu32`
type UInt32LiteralNode struct {
	TypedNodeBase
	Value string
}

func (n *UInt32LiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UInt32LiteralNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.span.Equal(o.span)
}

func (n *UInt32LiteralNode) String() string {
	return fmt.Sprintf("%su32", n.Value)
}

func (*UInt32LiteralNode) IsStatic() bool {
	return true
}

func (*UInt32LiteralNode) Class() *value.Class {
	return value.UInt32LiteralNodeClass
}

func (*UInt32LiteralNode) DirectClass() *value.Class {
	return value.UInt32LiteralNodeClass
}

func (n *UInt32LiteralNode) Inspect() string {
	return fmt.Sprintf("Std::Elk::AST::UInt32LiteralNode{span: %s, value: %s}", (*value.Span)(n.span).Inspect(), n.Value)
}

func (n *UInt32LiteralNode) Error() string {
	return n.Inspect()
}

// Create a new UInt32 literal node eg. `5u32`, `1_20u32`, `0xffu32`
func NewUInt32LiteralNode(span *position.Span, val string) *UInt32LiteralNode {
	return &UInt32LiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Int16 literal eg. `5i16`, `1_20i16`, `0xffi16`
type Int16LiteralNode struct {
	TypedNodeBase
	Value string
}

func (n *Int16LiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*Int16LiteralNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.span.Equal(o.span)
}

func (n *Int16LiteralNode) String() string {
	return fmt.Sprintf("%si16", n.Value)
}

func (*Int16LiteralNode) IsStatic() bool {
	return true
}

func (*Int16LiteralNode) Class() *value.Class {
	return value.Int16LiteralNodeClass
}

func (*Int16LiteralNode) DirectClass() *value.Class {
	return value.Int16LiteralNodeClass
}

func (n *Int16LiteralNode) Inspect() string {
	return fmt.Sprintf("Std::Elk::AST::Int16LiteralNode{span: %s, value: %s}", (*value.Span)(n.span).Inspect(), n.Value)
}

func (n *Int16LiteralNode) Error() string {
	return n.Inspect()
}

// Create a new Int16 literal node eg. `5i16`, `1_20i16`, `0xffi16`
func NewInt16LiteralNode(span *position.Span, val string) *Int16LiteralNode {
	return &Int16LiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// UInt16 literal eg. `5u16`, `1_20u16`, `0xffu16`
type UInt16LiteralNode struct {
	TypedNodeBase
	Value string
}

func (n *UInt16LiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UInt16LiteralNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.span.Equal(o.span)
}

func (n *UInt16LiteralNode) String() string {
	return fmt.Sprintf("%su16", n.Value)
}

func (*UInt16LiteralNode) IsStatic() bool {
	return true
}

func (*UInt16LiteralNode) Class() *value.Class {
	return value.UInt16LiteralNodeClass
}

func (*UInt16LiteralNode) DirectClass() *value.Class {
	return value.UInt16LiteralNodeClass
}

func (n *UInt16LiteralNode) Inspect() string {
	return fmt.Sprintf("Std::Elk::AST::UInt16LiteralNode{span: %s, value: %s}", (*value.Span)(n.span).Inspect(), n.Value)
}

func (n *UInt16LiteralNode) Error() string {
	return n.Inspect()
}

// Create a new UInt16 literal node eg. `5u16`, `1_20u16`, `0xffu16`
func NewUInt16LiteralNode(span *position.Span, val string) *UInt16LiteralNode {
	return &UInt16LiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Int8 literal eg. `5i8`, `1_20i8`, `0xffi8`
type Int8LiteralNode struct {
	TypedNodeBase
	Value string
}

func (n *Int8LiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*Int8LiteralNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.span.Equal(o.span)
}

func (n *Int8LiteralNode) String() string {
	return fmt.Sprintf("%si8", n.Value)
}

func (*Int8LiteralNode) IsStatic() bool {
	return true
}

func (*Int8LiteralNode) Class() *value.Class {
	return value.Int8LiteralNodeClass
}

func (*Int8LiteralNode) DirectClass() *value.Class {
	return value.Int8LiteralNodeClass
}

func (n *Int8LiteralNode) Inspect() string {
	return fmt.Sprintf("Std::Elk::AST::Int8LiteralNode{span: %s, value: %s}", (*value.Span)(n.span).Inspect(), n.Value)
}

func (n *Int8LiteralNode) Error() string {
	return n.Inspect()
}

// Create a new Int8 literal node eg. `5i8`, `1_20i8`, `0xffi8`
func NewInt8LiteralNode(span *position.Span, val string) *Int8LiteralNode {
	return &Int8LiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// UInt8 literal eg. `5u8`, `1_20u8`, `0xffu8`
type UInt8LiteralNode struct {
	TypedNodeBase
	Value string
}

func (n *UInt8LiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UInt8LiteralNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.span.Equal(o.span)
}

func (n *UInt8LiteralNode) String() string {
	return fmt.Sprintf("%su8", n.Value)
}

func (*UInt8LiteralNode) IsStatic() bool {
	return true
}

func (*UInt8LiteralNode) Class() *value.Class {
	return value.UInt8LiteralNodeClass
}

func (*UInt8LiteralNode) DirectClass() *value.Class {
	return value.UInt8LiteralNodeClass
}

func (n *UInt8LiteralNode) Inspect() string {
	return fmt.Sprintf("Std::Elk::AST::UInt8LiteralNode{span: %s, value: %s}", (*value.Span)(n.span).Inspect(), n.Value)
}

func (n *UInt8LiteralNode) Error() string {
	return n.Inspect()
}

// Create a new UInt8 literal node eg. `5u8`, `1_20u8`, `0xffu8`
func NewUInt8LiteralNode(span *position.Span, val string) *UInt8LiteralNode {
	return &UInt8LiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}
