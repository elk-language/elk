package types

import "github.com/elk-language/elk/value/symbol"

type Mixin struct {
	parent    Namespace
	abstract  bool
	singleton *SingletonClass
	NamespaceBase
}

func (m *Mixin) Singleton() *SingletonClass {
	return m.singleton
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

func (m *Mixin) IsPrimitive() bool {
	return false
}

func (m *Mixin) Parent() Namespace {
	return m.parent
}

func (m *Mixin) SetParent(parent Namespace) {
	m.parent = parent
}

func NewMixin(name string, env *GlobalEnvironment) *Mixin {
	mixin := &Mixin{
		NamespaceBase: MakeNamespaceBase(name),
	}
	mixin.singleton = NewSingletonClass(mixin, env.StdSubtypeClass(symbol.Mixin))

	return mixin
}

func NewMixinWithDetails(name string, parent *MixinProxy, consts *TypeMap, subtypes *TypeMap, methods *MethodMap, env *GlobalEnvironment) *Mixin {
	mixin := &Mixin{
		parent: parent,
		NamespaceBase: NamespaceBase{
			name:      name,
			constants: consts,
			methods:   methods,
			subtypes:  subtypes,
		},
	}
	mixin.singleton = NewSingletonClass(mixin, env.StdSubtypeClass(symbol.Mixin))

	return mixin
}

// Create a proxy that has a pointer to this mixin.
//
// Returns two values, the head and tail proxies.
// This is because of the fact that it's possible to include
// one mixin in another, so there is an entire inheritance chain.
func (m *Mixin) CreateProxy() (head *MixinProxy, tail Namespace) {
	var headParent Namespace
	if m.parent != nil {
		headParent = m.parent
	}
	headProxy := NewMixinProxy(m, headParent)

	var tailProxy Namespace = headProxy
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
