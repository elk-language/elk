package types

import "strings"

type GenericNamedType struct {
	Name          string
	Type          Type
	TypeVariables []*TypeVariable
}

func NewGenericNamedType(name string, typ Type, typeVars []*TypeVariable) *GenericNamedType {
	return &GenericNamedType{
		Name:          name,
		Type:          typ,
		TypeVariables: typeVars,
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
	for i, typeVar := range g.TypeVariables {
		if i > 0 {
			buffer.WriteString(", ")
		}
		switch typeVar.Variance {
		case COVARIANT:
			buffer.WriteRune('+')
		case CONTRAVARIANT:
			buffer.WriteRune('-')
		}

		buffer.WriteString(typeVar.Name)
		if typeVar.LowerBound != nil {
			buffer.WriteString(" > ")
			buffer.WriteString(Inspect(typeVar.LowerBound))
		}
		if typeVar.UpperBound != nil {
			buffer.WriteString(" < ")
			buffer.WriteString(Inspect(typeVar.UpperBound))
		}
	}
	buffer.WriteRune(']')
	return buffer.String()
}
