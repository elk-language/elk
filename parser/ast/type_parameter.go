package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
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

func (n *VariantTypeParameterNode) Splice(loc *position.Location, args *[]Node) Node {
	var lowerBound TypeNode
	if n.LowerBound != nil {
		lowerBound = n.LowerBound.Splice(loc, args).(TypeNode)
	}

	var upperBound TypeNode
	if n.UpperBound != nil {
		upperBound = n.UpperBound.Splice(loc, args).(TypeNode)
	}

	var def TypeNode
	if n.Default != nil {
		def = n.Default.Splice(loc, args).(TypeNode)
	}

	return &VariantTypeParameterNode{
		TypedNodeBase: TypedNodeBase{loc: getLoc(loc, n.loc), typ: n.typ},
		Variance:      n.Variance,
		Name:          n.Name,
		LowerBound:    lowerBound,
		UpperBound:    upperBound,
		Default:       def,
	}
}

func (n *VariantTypeParameterNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*VariantTypeParameterNode)
	if !ok {
		return false
	}

	if n.LowerBound == o.LowerBound {
	} else if n.LowerBound == nil || o.LowerBound == nil {
		return false
	} else if !n.LowerBound.Equal(value.Ref(o.LowerBound)) {
		return false
	}

	if n.UpperBound == o.UpperBound {
	} else if n.UpperBound == nil || o.UpperBound == nil {
		return false
	} else if !n.UpperBound.Equal(value.Ref(o.UpperBound)) {
		return false
	}

	if n.Default == o.Default {
	} else if n.Default == nil || o.Default == nil {
		return false
	} else if !n.Default.Equal(value.Ref(o.Default)) {
		return false
	}

	return n.loc.Equal(o.loc) &&
		n.Variance == o.Variance &&
		n.Name == o.Name
}

func (n *VariantTypeParameterNode) String() string {
	var buff strings.Builder

	switch n.Variance {
	case COVARIANT:
		buff.WriteRune('+')
	case CONTRAVARIANT:
		buff.WriteRune('-')
	}

	buff.WriteString(n.Name)

	if n.LowerBound != nil {
		buff.WriteString(" > ")
		buff.WriteString(n.LowerBound.String())
	}

	if n.UpperBound != nil {
		buff.WriteString(" < ")
		buff.WriteString(n.UpperBound.String())
	}

	if n.Default != nil {
		buff.WriteString(" = ")
		buff.WriteString(n.Default.String())
	}

	return buff.String()
}

func (*VariantTypeParameterNode) IsStatic() bool {
	return false
}

// Create a new type variable node eg. `+V`
func NewVariantTypeParameterNode(loc *position.Location, variance Variance, name string, lower, upper, def TypeNode) *VariantTypeParameterNode {
	return &VariantTypeParameterNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
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

	fmt.Fprintf(&buff, "Std::Elk::AST::VariantTypeParameterNode{\n  location: %s", (*value.Location)(n.loc).Inspect())

	buff.WriteString(",\n  variance: ")
	indent.IndentStringFromSecondLine(&buff, value.UInt8(n.Variance).Inspect(), 1)

	buff.WriteString(",\n  name: ")
	indent.IndentStringFromSecondLine(&buff, value.String(n.Name).Inspect(), 1)

	buff.WriteString(",\n  lower_bound: ")
	if n.LowerBound == nil {
		buff.WriteString("nil")
	} else {
		indent.IndentStringFromSecondLine(&buff, n.LowerBound.Inspect(), 1)
	}

	buff.WriteString(",\n  upper_bound: ")
	if n.UpperBound == nil {
		buff.WriteString("nil")
	} else {
		indent.IndentStringFromSecondLine(&buff, n.UpperBound.Inspect(), 1)
	}

	buff.WriteString(",\n  default: ")
	if n.Default == nil {
		buff.WriteString("nil")
	} else {
		indent.IndentStringFromSecondLine(&buff, n.Default.Inspect(), 1)
	}

	buff.WriteString("\n}")

	return buff.String()
}

func (n *VariantTypeParameterNode) Error() string {
	return n.Inspect()
}
