package types

import (
	"fmt"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type Mixin struct {
	parent         Namespace
	abstract       bool
	singleton      *SingletonClass
	TypeParameters []*TypeParameter
	NamespaceBase
}

func (m *Mixin) IsGeneric() bool {
	return len(m.TypeParameters) > 0
}

func IsMixin(typ Type) bool {
	_, ok := typ.(*Mixin)
	return ok
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
		case *Generic:
			switch n := base.Namespace.(type) {
			case *MixinProxy:
				var proxyCopy Namespace
				proxyCopy = NewMixinProxy(n.Mixin, nil)
				proxyCopy = NewGeneric(proxyCopy, base.TypeArguments)
				tailProxy.SetParent(proxyCopy)
				tailProxy = proxyCopy

				if n.parent == nil {
					break loop
				}
				baseProxy = n.parent
			case *InterfaceProxy:
				var proxyCopy Namespace
				proxyCopy = NewInterfaceProxy(n.Interface, nil)
				proxyCopy = NewGeneric(proxyCopy, base.TypeArguments)
				tailProxy.SetParent(proxyCopy)
				tailProxy = proxyCopy

				if n.parent == nil {
					break loop
				}
				baseProxy = n.parent
			default:
				panic(fmt.Sprintf("invalid mixin ancestor: %T", base.Namespace))
			}
		default:
			panic(fmt.Sprintf("invalid mixin ancestor: %T", baseProxy))
		}
	}

	return headProxy, tailProxy
}

func (m *Mixin) DefineMethod(docComment string, abstract, sealed, native bool, name value.Symbol, typeParams []*TypeParameter, params []*Parameter, returnType, throwType Type) *Method {
	method := NewMethod(docComment, abstract, sealed, native, name, typeParams, params, returnType, throwType, m)
	m.SetMethod(name, method)
	return method
}

func (m *Mixin) inspect() string {
	return m.name
}

func (m *Mixin) ToNonLiteral(env *GlobalEnvironment) Type {
	return m
}

func (*Mixin) IsLiteral() bool {
	return false
}
