package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Represents a type variable in generics like `class Foo[+V]; end`
type TypeParameterNode interface {
	Node
	typeVariableNode()
}

func (*InvalidNode) typeVariableNode()              {}
func (*VariantTypeParameterNode) typeVariableNode() {}

// Represents the variance of a type parameter.
type Variance uint8

const (
	INVARIANT Variance = iota
	COVARIANT
	CONTRAVARIANT
)

// Represents a type parameter eg. `+V`
type VariantTypeParameterNode struct {
	TypedNodeBase
	Variance   Variance // Variance level of this type parameter
	Name       string   // Name of the type parameter eg. `T`
	LowerBound TypeNode
	UpperBound TypeNode
	Default    TypeNode
}

func (*VariantTypeParameterNode) IsStatic() bool {
	return false
}

// Create a new type variable node eg. `+V`
func NewVariantTypeParameterNode(span *position.Span, variance Variance, name string, lower, upper, def TypeNode) *VariantTypeParameterNode {
	return &VariantTypeParameterNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Variance:      variance,
		Name:          name,
		LowerBound:    lower,
		UpperBound:    upper,
		Default:       def,
	}
}

func (*VariantTypeParameterNode) Class() *value.Class {
	return value.VariantTypeParameterNodeClass
}

func (*VariantTypeParameterNode) DirectClass() *value.Class {
	return value.VariantTypeParameterNodeClass
}

func (n *VariantTypeParameterNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::AST::VariantTypeParameterNode{\n  &: %p", n)

	buff.WriteString(",\n  variance: ")
	indentStringFromSecondLine(&buff, value.UInt8(n.Variance).Inspect(), 1)

	buff.WriteString(",\n  name: ")
	indentStringFromSecondLine(&buff, value.String(n.Name).Inspect(), 1)

	buff.WriteString(",\n  lower_bound: ")
	indentStringFromSecondLine(&buff, n.LowerBound.Inspect(), 1)

	buff.WriteString(",\n  upper_bound: ")
	indentStringFromSecondLine(&buff, n.UpperBound.Inspect(), 1)

	buff.WriteString(",\n  default: ")
	indentStringFromSecondLine(&buff, n.Default.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *VariantTypeParameterNode) Error() string {
	return n.Inspect()
}
