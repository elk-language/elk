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

func (m *MixinProxy) Copy() *MixinProxy {
	return &MixinProxy{
		parent: m.parent,
		Mixin:  m.Mixin,
	}
}

func (m *MixinProxy) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *MixinProxy {
	newMixin := &MixinProxy{}
	if m.parent != nil {
		newMixin.parent = DeepCopyEnv(m.parent, oldEnv, newEnv).(Namespace)
	}
	newMixin.Mixin = DeepCopyEnv(m.Mixin, oldEnv, newEnv).(*Mixin)
	return newMixin
}
