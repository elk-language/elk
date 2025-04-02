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

func (n *FloatLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*FloatLiteralNode)
	if !ok {
		return false
	}

	return n.span.Equal(o.span) &&
		n.Value == o.Value
}

func (n *FloatLiteralNode) String() string {
	return n.Value
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
	return fmt.Sprintf("Std::Elk::AST::FloatLiteralNode{span: %s, value: %s}", (*value.Span)(n.span).Inspect(), n.Value)
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

func (n *BigFloatLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*BigFloatLiteralNode)
	if !ok {
		return false
	}

	return n.span.Equal(o.span) &&
		n.Value == o.Value
}

func (n *BigFloatLiteralNode) String() string {
	return fmt.Sprintf("%sbf", n.Value)
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
	return fmt.Sprintf("Std::Elk::AST::BigFloatLiteralNode{span: %s, value: %s}", (*value.Span)(n.span).Inspect(), n.Value)
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

func (n *Float64LiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*Float64LiteralNode)
	if !ok {
		return false
	}

	return n.span.Equal(o.span) &&
		n.Value == o.Value
}

func (n *Float64LiteralNode) String() string {
	return fmt.Sprintf("%sf64", n.Value)
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
	return fmt.Sprintf("Std::Elk::AST::Float64LiteralNode{span: %s, value: %s}", (*value.Span)(n.span).Inspect(), n.Value)
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

func (n *Float32LiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*Float32LiteralNode)
	if !ok {
		return false
	}

	return n.span.Equal(o.span) &&
		n.Value == o.Value
}

func (n *Float32LiteralNode) String() string {
	return fmt.Sprintf("%sf32", n.Value)
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
	return fmt.Sprintf("Std::Elk::AST::Float32LiteralNode{span: %s, value: %s}", (*value.Span)(n.span).Inspect(), n.Value)
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
