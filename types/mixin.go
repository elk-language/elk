package types

type Mixin struct {
	parent *MixinProxy
	ConstantMap
}

func (m *Mixin) Parent() ConstantContainer {
	if m.parent == nil {
		return nil
	}
	return m.parent
}

func (m *Mixin) SetParent(parent *MixinProxy) {
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
func (m *Mixin) CreateProxy() (head, tail *MixinProxy) {
	headProxy := NewMixinProxy(m, m.parent)

	tailProxy := headProxy
	baseProxy := m.parent
	for baseProxy != nil {
		proxyCopy := NewMixinProxy(baseProxy.Mixin, nil)
		tailProxy.parent = proxyCopy

		if baseProxy.parent == nil {
			break
		}
		baseProxy = baseProxy.parent.(*MixinProxy)
		tailProxy = proxyCopy
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
