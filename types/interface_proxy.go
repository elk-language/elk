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
