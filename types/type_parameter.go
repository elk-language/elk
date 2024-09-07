package types

import (
	"strings"

	"github.com/elk-language/elk/value"
)

type Variance uint8

func (v Variance) String() string {
	return varianceString[v]
}

func (v Variance) Name() string {
	switch v {
	case COVARIANT:
		return "covariant"
	case CONTRAVARIANT:
		return "contravariant"
	default:
		return "invariant"
	}
}

const (
	INVARIANT Variance = iota
	COVARIANT
	CONTRAVARIANT
)

var varianceString = []string{
	INVARIANT:     "INVARIANT",
	COVARIANT:     "COVARIANT",
	CONTRAVARIANT: "CONTRAVARIANT",
}

type TypeParameter struct {
	Name       value.Symbol
	Namespace  Namespace
	LowerBound Type
	UpperBound Type
	Variance   Variance
}

func (m *TypeParameter) Copy() *TypeParameter {
	return &TypeParameter{
		Name:       m.Name,
		Namespace:  m.Namespace,
		LowerBound: m.LowerBound,
		UpperBound: m.UpperBound,
		Variance:   m.Variance,
	}
}

func NewTypeParameter(name value.Symbol, namespace Namespace, lowerBound, upperBound Type, variance Variance) *TypeParameter {
	return &TypeParameter{
		Name:       name,
		Namespace:  namespace,
		LowerBound: lowerBound,
		UpperBound: upperBound,
		Variance:   variance,
	}
}

func (t *TypeParameter) ToNonLiteral(env *GlobalEnvironment) Type {
	return t
}

func (*TypeParameter) IsLiteral() bool {
	return false
}

func (t *TypeParameter) inspect() string {
	return t.Name.String()
}

func (t *TypeParameter) InspectSignature() string {
	buffer := new(strings.Builder)
	switch t.Variance {
	case COVARIANT:
		buffer.WriteRune('+')
	case CONTRAVARIANT:
		buffer.WriteRune('-')
	}
	buffer.WriteString(t.Name.String())

	if !IsNever(t.LowerBound) {
		buffer.WriteString(" > ")
		buffer.WriteString(Inspect(t.LowerBound))
	}

	if !IsAny(t.UpperBound) {
		buffer.WriteString(" < ")
		buffer.WriteString(Inspect(t.UpperBound))
	}

	return buffer.String()
}
