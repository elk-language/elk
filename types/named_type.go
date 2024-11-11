package types

import "github.com/elk-language/elk/value"

type NamedType struct {
	Name string
	Type Type
}

func NewNamedType(name string, typ Type) *NamedType {
	return &NamedType{
		Name: name,
		Type: typ,
	}
}

func (n *NamedType) ToNonLiteral(env *GlobalEnvironment) Type {
	return n
}

func (*NamedType) IsLiteral() bool {
	return false
}

func (n *NamedType) inspect() string {
	return n.Name
}

func (n *NamedType) Copy() *NamedType {
	return &NamedType{
		Name: n.Name,
		Type: n.Type,
	}
}

func (n *NamedType) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *NamedType {
	if newType, ok := NameToTypeOk(n.Name, newEnv); ok {
		return newType.(*NamedType)
	}

	newType := n.Copy()

	classConstantPath := GetConstantPath(n.Name)
	parentNamespace := DeepCopyNamespacePath(classConstantPath[:len(classConstantPath)-1], oldEnv, newEnv)
	classConstantName := classConstantPath[len(classConstantPath)-1]
	parentNamespace.DefineSubtype(value.ToSymbol(classConstantName), newType)

	newType.Type = DeepCopyEnv(n.Type, oldEnv, newEnv)

	return newType
}
