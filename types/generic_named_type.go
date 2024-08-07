package types

import (
	"strings"
)

type GenericNamedType struct {
	Name           string
	Type           Type
	TypeParameters []*TypeParameter
}

func NewGenericNamedType(name string, typ Type, typeVars []*TypeParameter) *GenericNamedType {
	return &GenericNamedType{
		Name:           name,
		Type:           typ,
		TypeParameters: typeVars,
	}
}

func (g *GenericNamedType) ToNonLiteral(env *GlobalEnvironment) Type {
	return g
}

func (*GenericNamedType) IsLiteral() bool {
	return false
}

func (g *GenericNamedType) inspect() string {
	buffer := new(strings.Builder)

	buffer.WriteString(g.Name)
	buffer.WriteRune('[')
	for i, typeVar := range g.TypeParameters {
		if i > 0 {
			buffer.WriteString(", ")
		}
		switch typeVar.Variance {
		case COVARIANT:
			buffer.WriteRune('+')
		case CONTRAVARIANT:
			buffer.WriteRune('-')
		}

		buffer.WriteString(typeVar.Name.String())
		if !IsNever(typeVar.LowerBound) {
			buffer.WriteString(" > ")
			buffer.WriteString(Inspect(typeVar.LowerBound))
		}
		if !IsAny(typeVar.UpperBound) {
			buffer.WriteString(" < ")
			buffer.WriteString(Inspect(typeVar.UpperBound))
		}
	}
	buffer.WriteRune(']')
	return buffer.String()
}
