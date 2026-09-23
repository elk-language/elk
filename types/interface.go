package types

import (
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/value/symbol"
)

type Interface struct {
	parent         Ref[Namespace]
	singleton      Ref[*SingletonClass]
	Checked        bool
	compiled       bool
	typeParameters []Ref[*TypeParameter]
	NamespaceBase
}

var _ Namespace = &Interface{}

func (i *Interface) ToRef() Ref[*Interface] {
	return ToRef(i)
}

func (i *Interface) EqualAny(other any) bool {
	o, ok := other.(*Interface)
	if !ok {
		return false
	}

	if i.id > 0 {
		return i.id == o.ID()
	}

	return i.name == o.name
}

func (i *Interface) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(i, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(i, parent)
	}
}

func (i *Interface) IsGeneric() bool {
	return len(i.typeParameters) > 0
}

func (i *Interface) TypeParameters() []Ref[*TypeParameter] {
	return i.typeParameters
}

func (i *Interface) SetTypeParameters(t []Ref[*TypeParameter]) {
	i.typeParameters = t
}

func IsInterface(typ Type) bool {
	_, ok := typ.(*Interface)
	return ok
}

func (i *Interface) Singleton() *SingletonClass {
	return i.singleton.Get()
}

func (i *Interface) SingletonRef() Ref[*SingletonClass] {
	return i.singleton
}

func (i *Interface) SetSingleton(singleton *SingletonClass) {
	i.singleton = ToRef(singleton)
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

func (*Interface) IsNative() bool {
	return false
}

func (*Interface) IsSealed() bool {
	return false
}

func (*Interface) IsPrimitive() bool {
	return false
}

func (*Interface) IsImmutable() bool {
	return false
}

func (i *Interface) Parent() Namespace {
	return i.parent.Get()
}

func (i *Interface) ParentRef() Ref[Namespace] {
	return i.parent
}

func (i *Interface) SetParent(parent Namespace) {
	i.parent = ToRef(parent)
}

func NewInterface(docComment string, name string) *Interface {
	iface := &Interface{
		NamespaceBase: MakeNamespaceBase(docComment, name),
		compiled:      Env.Init,
	}
	iface.singleton = NewSingletonClass(iface, Env.StdSubtypeClass(symbol.C_Interface)).ToRef()

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
		parent:   ToRef[Namespace](parent),
		compiled: env.Init,
		NamespaceBase: NamespaceBase{
			name:      name,
			constants: consts,
			methods:   methods,
			subtypes:  subtypes,
		},
	}
}

func (i *Interface) DefineMethod(docComment string, flags bitfield.BitFlag16, name symbol.Symbol, typeParams []Ref[*TypeParameter], params []Ref[*Parameter], returnType, throwType Type) *Method {
	method := NewMethod(docComment, flags, name, typeParams, params, returnType, throwType, i)
	i.SetMethod(name, method)
	return method
}

func (i *Interface) inspect() string {
	return i.name
}

func (i *Interface) ToNonLiteral() Type {
	return i
}

func (*Interface) IsLiteral() bool {
	return false
}

func (i *Interface) Copy() *Interface {
	return &Interface{
		parent:         i.parent,
		compiled:       i.compiled,
		Checked:        i.Checked,
		singleton:      i.singleton,
		typeParameters: i.typeParameters,
		NamespaceBase: NamespaceBase{
			name:      i.name,
			constants: i.constants,
			methods:   i.methods,
			subtypes:  i.subtypes,
		},
	}
}

func (i *Interface) CopyType() Type {
	return i.Copy()
}
