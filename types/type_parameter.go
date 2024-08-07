package types

import "github.com/elk-language/elk/value"

type Variance uint8

const (
	INVARIANT Variance = iota
	COVARIANT
	CONTRAVARIANT
)

type TypeParameter struct {
	Name       value.Symbol
	LowerBound Type
	UpperBound Type
	Variance   Variance
}

func NewTypeParameter(name value.Symbol, lowerBound, upperBound Type, variance Variance) *TypeParameter {
	return &TypeParameter{
		Name:       name,
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
