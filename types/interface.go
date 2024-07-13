package types

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type Interface struct {
	parent    *InterfaceProxy
	singleton *SingletonClass
	NamespaceBase
}

func (i *Interface) ImplementInterface(implementedInterface *Interface) {
	headProxy, tailProxy := implementedInterface.CreateProxy()
	tailProxy.SetParent(i.Parent())
	i.SetParent(headProxy)
}

func (i *Interface) Singleton() *SingletonClass {
	return i.singleton
}

func (*Interface) IsAbstract() bool {
	return true
}

func (*Interface) IsSealed() bool {
	return false
}

func (*Interface) IsPrimitive() bool {
	return false
}

func (i *Interface) Parent() Namespace {
	if i.parent == nil {
		return nil
	}
	return i.parent
}

func (i *Interface) SetParent(parent Namespace) {
	i.parent = parent.(*InterfaceProxy)
}

func NewInterface(docComment string, name string, env *GlobalEnvironment) *Interface {
	iface := &Interface{
		NamespaceBase: MakeNamespaceBase(docComment, name),
	}
	iface.singleton = NewSingletonClass(iface, env.StdSubtypeClass(symbol.Interface))

	return iface
}

func NewInterfaceWithDetails(name string, parent *InterfaceProxy, consts *TypeMap, subtypes *TypeMap, methods *MethodMap) *Interface {
	return &Interface{
		parent: parent,
		NamespaceBase: NamespaceBase{
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
	var headParent Namespace
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

func (i *Interface) DefineMethod(docComment string, abstract, sealed, native bool, name value.Symbol, params []*Parameter, returnType, throwType Type) *Method {
	method := NewMethod(docComment, abstract, sealed, native, name, params, returnType, throwType, i)
	i.SetMethod(name, method)
	return method
}

func (i *Interface) inspect() string {
	return i.name
}

func (i *Interface) ToNonLiteral(env *GlobalEnvironment) Type {
	return i
}
