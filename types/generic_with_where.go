package types

// Represents a generic namespace parent with a where clause
// `include Foo[Bar] where Bar < Baz`
type GenericWithWere struct {
	*Generic
	Where []*TypeParameter
}

func IsGenericWithWere(typ Type) bool {
	_, ok := typ.(*GenericWithWere)
	return ok
}

func NewGenericWithWhere(generic *Generic, where []*TypeParameter) *GenericWithWere {
	return &GenericWithWere{
		Generic: generic,
		Where:   where,
	}
}

func (g *GenericWithWere) ToNonLiteral(env *GlobalEnvironment) Type {
	return g
}
