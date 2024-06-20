package types

type Mixin struct {
	parent   ConstantContainer
	abstract bool
	ConstantMap
}

func (m *Mixin) SetAbstract(abstract bool) *Mixin {
	m.abstract = abstract
	return m
}

func (m *Mixin) IsAbstract() bool {
	return m.abstract
}

func (m *Mixin) IsSealed() bool {
	return false
}

func (m *Mixin) Parent() ConstantContainer {
	return m.parent
}

func (m *Mixin) SetParent(parent ConstantContainer) {
	m.parent = parent
}

func NewMixin(name string) *Mixin {
	return &Mixin{
		ConstantMap: MakeConstantMap(name),
	}
}

func NewMixinWithDetails(name string, parent *MixinProxy, consts *TypeMap, subtypes *TypeMap, methods *MethodMap) *Mixin {
	return &Mixin{
		parent: parent,
		ConstantMap: ConstantMap{
			name:      name,
			constants: consts,
			methods:   methods,
			subtypes:  subtypes,
		},
	}
}

// Create a proxy that has a pointer to this mixin.
//
// Returns two values, the head and tail proxies.
// This is because of the fact that it's possible to include
// one mixin in another, so there is an entire inheritance chain.
func (m *Mixin) CreateProxy() (head *MixinProxy, tail ConstantContainer) {
	var headParent ConstantContainer
	if m.parent != nil {
		headParent = m.parent
	}
	headProxy := NewMixinProxy(m, headParent)

	var tailProxy ConstantContainer = headProxy
	baseProxy := m.parent
loop:
	for baseProxy != nil {
		switch base := baseProxy.(type) {
		case *MixinProxy:
			proxyCopy := NewMixinProxy(base.Mixin, nil)
			tailProxy.SetParent(proxyCopy)
			tailProxy = proxyCopy

			if base.parent == nil {
				break loop
			}
			baseProxy = base.parent
		case *InterfaceProxy:
			proxyCopy := NewInterfaceProxy(base.Interface, nil)
			tailProxy.SetParent(proxyCopy)
			tailProxy = proxyCopy

			if base.parent == nil {
				break loop
			}
			baseProxy = base.parent
		}
	}

	return headProxy, tailProxy
}

func (m *Mixin) DefineMethod(name string, params []*Parameter, returnType, throwType Type) *Method {
	method := NewMethod(name, params, returnType, throwType, m)
	m.SetMethod(name, method)
	return method
}

func (m *Mixin) inspect() string {
	return m.name
}

func (m *Mixin) ToNonLiteral(env *GlobalEnvironment) Type {
	return m
}
