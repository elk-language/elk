package types

import "github.com/elk-language/elk/value/symbol"

type Mixin struct {
	parent    Namespace
	abstract  bool
	singleton *SingletonClass
	NamespaceBase
}

func (m *Mixin) IncludeMixin(includedMixin *Mixin) {
	headProxy, tailProxy := includedMixin.CreateProxy()
	tailProxy.SetParent(m.Parent())
	m.SetParent(headProxy)
}

func (m *Mixin) ImplementInterface(iface *Interface) {
	headProxy, tailProxy := iface.CreateProxy()
	tailProxy.SetParent(m.Parent())
	m.SetParent(headProxy)
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

func (m *Mixin) DeclaresInstanceVariables() bool {
	var currentNamespace Namespace = m

	for ; currentNamespace != nil; currentNamespace = currentNamespace.Parent() {
		if currentNamespace.InstanceVariables().Len() > 0 {
			return true
		}
	}

	return false
}

func NewMixin(docComment string, abstract bool, name string, env *GlobalEnvironment) *Mixin {
	mixin := &Mixin{
		abstract:      abstract,
		NamespaceBase: MakeNamespaceBase(docComment, name),
	}
	mixin.singleton = NewSingletonClass(mixin, env.StdSubtypeClass(symbol.Mixin))

	return mixin
}

func NewMixinWithDetails(docComment string, abstract bool, name string, parent *MixinProxy, consts *TypeMap, subtypes *TypeMap, methods *MethodMap, env *GlobalEnvironment) *Mixin {
	mixin := &Mixin{
		parent:   parent,
		abstract: abstract,
		NamespaceBase: NamespaceBase{
			docComment: docComment,
			name:       name,
			constants:  consts,
			methods:    methods,
			subtypes:   subtypes,
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

func (m *Mixin) DefineMethod(docComment string, abstract, sealed, native bool, name string, params []*Parameter, returnType, throwType Type) *Method {
	method := NewMethod(docComment, abstract, sealed, native, name, params, returnType, throwType, m)
	m.SetMethod(name, method)
	return method
}

func (m *Mixin) inspect() string {
	return m.name
}

func (m *Mixin) ToNonLiteral(env *GlobalEnvironment) Type {
	return m
}
