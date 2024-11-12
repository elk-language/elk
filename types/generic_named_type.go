package types

import (
	"strings"

	"github.com/elk-language/elk/value"
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
	if len(g.TypeParameters) < 1 {
		return g.Name
	}

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
		case BIVARIANT:
			buffer.WriteString("+-")
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

func (g *GenericNamedType) Copy() *GenericNamedType {
	return &GenericNamedType{
		Name:           g.Name,
		Type:           g.Type,
		TypeParameters: g.TypeParameters,
	}
}

func (g *GenericNamedType) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *GenericNamedType {
	if newType, ok := NameToTypeOk(g.Name, newEnv); ok {
		return newType.(*GenericNamedType)
	}

	newType := &GenericNamedType{
		Name: g.Name,
	}

	classConstantPath := GetConstantPath(g.Name)
	parentNamespace := DeepCopyNamespacePath(classConstantPath[:len(classConstantPath)-1], oldEnv, newEnv)
	classConstantName := classConstantPath[len(classConstantPath)-1]
	parentNamespace.DefineSubtype(value.ToSymbol(classConstantName), newType)

	newType.Type = DeepCopyEnv(g.Type, oldEnv, newEnv)
	newType.TypeParameters = TypeParametersDeepCopyEnv(g.TypeParameters, oldEnv, newEnv)

	return newType
}
