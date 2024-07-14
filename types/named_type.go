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

func (*NamedType) IsLiteral() bool {
	return false
}

func (n *NamedType) inspect() string {
	return n.Name
}
