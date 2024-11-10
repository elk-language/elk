package types

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type Interface struct {
	parent         Namespace
	singleton      *SingletonClass
	Checked        bool
	compiled       bool
	typeParameters []*TypeParameter
	NamespaceBase
}

func (i *Interface) IsGeneric() bool {
	return len(i.typeParameters) > 0
}

func (i *Interface) TypeParameters() []*TypeParameter {
	return i.typeParameters
}

func (i *Interface) SetTypeParameters(t []*TypeParameter) {
	i.typeParameters = t
}

func IsInterface(typ Type) bool {
	_, ok := typ.(*Interface)
	return ok
}

func (i *Interface) Singleton() *SingletonClass {
	return i.singleton
}

func (i *Interface) IsDefined() bool {
	return i.compiled
}

func (i *Interface) SetDefined(compiled bool) {
	i.compiled = compiled
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
		compiled:      env.Init,
	}
	iface.singleton = NewSingletonClass(iface, env.StdSubtypeClass(symbol.Interface))

	return iface
}

func NewInterfaceWithDetails(
	name string,
	parent *InterfaceProxy,
	consts ConstantMap,
	subtypes ConstantMap,
	methods MethodMap,
	env *GlobalEnvironment,
) *Interface {
	return &Interface{
		parent:   parent,
		compiled: env.Init,
		NamespaceBase: NamespaceBase{
			name:      name,
			constants: consts,
			methods:   methods,
			subtypes:  subtypes,
		},
	}
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
