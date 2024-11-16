package types

import (
	"strings"

	"github.com/elk-language/elk/lexer"
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
	case BIVARIANT:
		return "bivariant"
	case INVARIANT:
		return "invariant"
	default:
		return "<unknown>"
	}
}

const (
	INVARIANT Variance = iota
	COVARIANT
	CONTRAVARIANT
	BIVARIANT
)

var varianceString = []string{
	INVARIANT:     "INVARIANT",
	COVARIANT:     "COVARIANT",
	CONTRAVARIANT: "CONTRAVARIANT",
	BIVARIANT:     "BIVARIANT",
}

type TypeParameter struct {
	Name       value.Symbol
	Namespace  Namespace
	LowerBound Type
	UpperBound Type
	Variance   Variance
}

func (t *TypeParameter) Copy() *TypeParameter {
	return &TypeParameter{
		Name:       t.Name,
		Namespace:  t.Namespace,
		LowerBound: t.LowerBound,
		UpperBound: t.UpperBound,
		Variance:   t.Variance,
	}
}

func (t *TypeParameter) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *TypeParameter {
	newTypeParam := t.Copy()
	if t.Namespace != nil {
		newTypeParam.Namespace = DeepCopyEnv(t.Namespace, oldEnv, newEnv).(Namespace)
	}
	newTypeParam.LowerBound = DeepCopyEnv(t.LowerBound, oldEnv, newEnv)
	newTypeParam.UpperBound = DeepCopyEnv(t.UpperBound, oldEnv, newEnv)
	return newTypeParam
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
	case BIVARIANT:
		buffer.WriteString("+-")
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

func (t *TypeParameter) InspectSignatureWithColor() string {
	return lexer.Colorize(t.InspectSignature())
}
