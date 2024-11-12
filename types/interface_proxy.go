package types

type InterfaceProxy struct {
	parent Namespace
	*Interface
}

func (i *InterfaceProxy) Parent() Namespace {
	return i.parent
}

func (i *InterfaceProxy) SetParent(parent Namespace) {
	i.parent = parent
}

func NewInterfaceProxy(iface *Interface, parent Namespace) *InterfaceProxy {
	return &InterfaceProxy{
		parent:    parent,
		Interface: iface,
	}
}

func (i *InterfaceProxy) ToNonLiteral(env *GlobalEnvironment) Type {
	return i
}

func (m *InterfaceProxy) Copy() *InterfaceProxy {
	return &InterfaceProxy{
		parent:    m.parent,
		Interface: m.Interface,
	}
}

func (i *InterfaceProxy) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *InterfaceProxy {
	newIface := &InterfaceProxy{}
	if i.parent != nil {
		newIface.parent = DeepCopyEnv(i.parent, oldEnv, newEnv).(Namespace)
	}
	newIface.Interface = DeepCopyEnv(i.Interface, oldEnv, newEnv).(*Interface)
	return newIface
}
