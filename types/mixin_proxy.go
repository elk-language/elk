package types

type MixinProxy struct {
	parent Namespace
	*Mixin
}

func (m *MixinProxy) Parent() Namespace {
	return m.parent
}

func (m *MixinProxy) SetParent(parent Namespace) {
	m.parent = parent
}

func NewMixinProxy(mixin *Mixin, parent Namespace) *MixinProxy {
	return &MixinProxy{
		parent: parent,
		Mixin:  mixin,
	}
}

func (m *MixinProxy) ToNonLiteral(env *GlobalEnvironment) Type {
	return m
}
