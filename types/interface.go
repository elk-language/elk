package types

import (
	"fmt"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type Interface struct {
	parent         Namespace
	singleton      *SingletonClass
	Checked        bool
	TypeParameters []*TypeParameter
	NamespaceBase
}

func (i *Interface) IsGeneric() bool {
	return len(i.TypeParameters) > 0
}

func IsInterface(typ Type) bool {
	_, ok := typ.(*Interface)
	return ok
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
	return i.parent
}

func (i *Interface) SetParent(parent Namespace) {
	i.parent = parent
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
func (i *Interface) CreateProxy() (head *InterfaceProxy, tail Namespace) {
	var headParent Namespace
	if i.parent != nil {
		headParent = i.parent
	}
	headProxy := NewInterfaceProxy(i, headParent)

	var tailProxy Namespace = headProxy
	baseProxy := i.parent

loop:
	for baseProxy != nil {
		switch base := baseProxy.(type) {
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
				panic(fmt.Sprintf("invalid interface ancestor: %T", base.Namespace))
			}
		default:
			panic(fmt.Sprintf("invalid interface ancestor: %T", baseProxy))
		}
	}

	return headProxy, tailProxy
}

func (i *Interface) DefineMethod(docComment string, abstract, sealed, native bool, name value.Symbol, typeParams []*TypeParameter, params []*Parameter, returnType, throwType Type) *Method {
	method := NewMethod(docComment, abstract, sealed, native, name, typeParams, params, returnType, throwType, i)
	i.SetMethod(name, method)
	return method
}

func (i *Interface) inspect() string {
	return i.name
}

func (i *Interface) ToNonLiteral(env *GlobalEnvironment) Type {
	return i
}

func (*Interface) IsLiteral() bool {
	return false
}
