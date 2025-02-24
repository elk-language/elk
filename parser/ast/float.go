package ast

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Float literal eg. `5.2`, `.5`, `45e20`
type FloatLiteralNode struct {
	TypedNodeBase
	Value string
}

func (*FloatLiteralNode) IsStatic() bool {
	return true
}

func (*FloatLiteralNode) Class() *value.Class {
	return value.FloatLiteralNodeClass
}

func (*FloatLiteralNode) DirectClass() *value.Class {
	return value.FloatLiteralNodeClass
}

func (n *FloatLiteralNode) Inspect() string {
	return fmt.Sprintf("Std::AST::FloatLiteralNode{&: %p, value: %s}", n, n.Value)
}

func (n *FloatLiteralNode) Error() string {
	return n.Inspect()
}

// Create a new float literal node eg. `5.2`, `.5`, `45e20`
func NewFloatLiteralNode(span *position.Span, val string) *FloatLiteralNode {
	return &FloatLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// BigFloat literal eg. `5.2bf`, `.5bf`, `45e20bf`
type BigFloatLiteralNode struct {
	TypedNodeBase
	Value string
}

func (*BigFloatLiteralNode) IsStatic() bool {
	return true
}

func (*BigFloatLiteralNode) Class() *value.Class {
	return value.BigFloatLiteralNodeClass
}

func (*BigFloatLiteralNode) DirectClass() *value.Class {
	return value.BigFloatLiteralNodeClass
}

func (n *BigFloatLiteralNode) Inspect() string {
	return fmt.Sprintf("Std::AST::BigFloatLiteralNode{&: %p, value: %s}", n, n.Value)
}

func (n *BigFloatLiteralNode) Error() string {
	return n.Inspect()
}

// Create a new BigFloat literal node eg. `5.2bf`, `.5bf`, `45e20bf`
func NewBigFloatLiteralNode(span *position.Span, val string) *BigFloatLiteralNode {
	return &BigFloatLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Float64 literal eg. `5.2f64`, `.5f64`, `45e20f64`
type Float64LiteralNode struct {
	TypedNodeBase
	Value string
}

func (*Float64LiteralNode) IsStatic() bool {
	return true
}

func (*Float64LiteralNode) Class() *value.Class {
	return value.Float64LiteralNodeClass
}

func (*Float64LiteralNode) DirectClass() *value.Class {
	return value.Float64LiteralNodeClass
}

func (n *Float64LiteralNode) Inspect() string {
	return fmt.Sprintf("Std::AST::Float64LiteralNode{&: %p, value: %s}", n, n.Value)
}

func (n *Float64LiteralNode) Error() string {
	return n.Inspect()
}

// Create a new Float64 literal node eg. `5.2f64`, `.5f64`, `45e20f64`
func NewFloat64LiteralNode(span *position.Span, val string) *Float64LiteralNode {
	return &Float64LiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Float32 literal eg. `5.2f32`, `.5f32`, `45e20f32`
type Float32LiteralNode struct {
	TypedNodeBase
	Value string
}

func (*Float32LiteralNode) IsStatic() bool {
	return true
}

func (*Float32LiteralNode) Class() *value.Class {
	return value.Float32LiteralNodeClass
}

func (*Float32LiteralNode) DirectClass() *value.Class {
	return value.Float32LiteralNodeClass
}

func (n *Float32LiteralNode) Inspect() string {
	return fmt.Sprintf("Std::AST::Float32LiteralNode{&: %p, value: %s}", n, n.Value)
}

func (n *Float32LiteralNode) Error() string {
	return n.Inspect()
}

// Create a new Float32 literal node eg. `5.2f32`, `.5f32`, `45e20f32`
func NewFloat32LiteralNode(span *position.Span, val string) *Float32LiteralNode {
	return &Float32LiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}
