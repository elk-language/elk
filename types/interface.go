package types

type Interface struct {
	parent *InterfaceProxy
	ConstantMap
}

func (*Interface) IsAbstract() bool {
	return true
}

func (*Interface) IsSealed() bool {
	return false
}

func (i *Interface) Parent() ConstantContainer {
	return i.parent
}

func (i *Interface) SetParent(parent ConstantContainer) {
	i.parent = parent.(*InterfaceProxy)
}

func NewInterface(name string) *Interface {
	return &Interface{
		ConstantMap: MakeConstantMap(name),
	}
}

func NewInterfaceWithDetails(name string, parent *InterfaceProxy, consts *TypeMap, subtypes *TypeMap, methods *MethodMap) *Interface {
	return &Interface{
		parent: parent,
		ConstantMap: ConstantMap{
			name:      name,
			constants: consts,
			methods:   methods,
			subtypes:  subtypes,
		},
	}
}

// Create a proxy that has a pointer to this interface.
//
// Returns two values, the head and tail proxies.
// This is because of the fact that it's possible to include
// one mixin in another, so there is an entire inheritance chain.
func (i *Interface) CreateProxy() (head, tail *InterfaceProxy) {
	var headParent ConstantContainer
	if i.parent != nil {
		headParent = i.parent
	}
	headProxy := NewInterfaceProxy(i, headParent)

	tailProxy := headProxy
	baseProxy := i.parent
	for baseProxy != nil {
		proxyCopy := NewInterfaceProxy(baseProxy.Interface, nil)
		tailProxy.parent = proxyCopy
		tailProxy = proxyCopy

		if baseProxy.parent == nil {
			break
		}
		baseProxy = baseProxy.parent.(*InterfaceProxy)
	}

	return headProxy, tailProxy
}

func (i *Interface) DefineMethod(name string, params []*Parameter, returnType, throwType Type) *Method {
	method := NewMethod(name, params, returnType, throwType, i)
	i.SetMethod(name, method)
	return method
}

func (i *Interface) inspect() string {
	return i.name
}

func (i *Interface) ToNonLiteral(env *GlobalEnvironment) Type {
	return i
}
