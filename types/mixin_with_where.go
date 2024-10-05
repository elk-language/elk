package types

// Represents a mixin parent with a where clause
// `extend where Bar < Baz`
type MixinWithWhere struct {
	*MixinProxy
	Namespace Namespace
	Where     []*TypeParameter
}

func IsMixinWithWhere(typ Type) bool {
	_, ok := typ.(*MixinWithWhere)
	return ok
}

func NewMixinWithWhere(mixin *MixinProxy, namespace Namespace, where []*TypeParameter) *MixinWithWhere {
	return &MixinWithWhere{
		MixinProxy: mixin,
		Namespace:  namespace,
		Where:      where,
	}
}

func (m *MixinWithWhere) ToNonLiteral(env *GlobalEnvironment) Type {
	return m
}

func (m *MixinWithWhere) inspect() string {
	return m.Namespace.inspect()
}
