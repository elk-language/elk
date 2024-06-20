package types

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

func (n *NamedType) inspect() string {
	return n.Name
}
