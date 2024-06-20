package types

type InterfaceProxy struct {
	parent ConstantContainer
	*Interface
}

func (i *InterfaceProxy) Parent() ConstantContainer {
	return i.parent
}

func (i *InterfaceProxy) SetParent(parent ConstantContainer) {
	i.parent = parent
}

func NewInterfaceProxy(iface *Interface, parent ConstantContainer) *InterfaceProxy {
	return &InterfaceProxy{
		parent:    parent,
		Interface: iface,
	}
}

func (i *InterfaceProxy) ToNonLiteral(env *GlobalEnvironment) Type {
	return i
}
