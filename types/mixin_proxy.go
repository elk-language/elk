package types

type MixinProxy struct {
	parent ConstantContainer
	*Mixin
}

func (m *MixinProxy) Parent() ConstantContainer {
	return m.parent
}

func (m *MixinProxy) SetParent(parent ConstantContainer) {
	m.parent = parent
}

func NewMixinProxy(mixin *Mixin, parent ConstantContainer) *MixinProxy {
	return &MixinProxy{
		parent: parent,
		Mixin:  mixin,
	}
}

func (m *MixinProxy) ToNonLiteral(env *GlobalEnvironment) Type {
	return m
}
