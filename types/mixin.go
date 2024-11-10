package types

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type Mixin struct {
	parent         Namespace
	abstract       bool
	compiled       bool
	Checked        bool
	singleton      *SingletonClass
	typeParameters []*TypeParameter
	NamespaceBase
}

func (m *Mixin) IsGeneric() bool {
	return len(m.typeParameters) > 0
}

func (m *Mixin) TypeParameters() []*TypeParameter {
	return m.typeParameters
}

func (m *Mixin) SetTypeParameters(t []*TypeParameter) {
	m.typeParameters = t
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

func (m *Mixin) IsDefined() bool {
	return m.compiled
}

func (m *Mixin) SetDefined(compiled bool) {
	m.compiled = compiled
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
		compiled:      env.Init,
		NamespaceBase: MakeNamespaceBase(docComment, name),
	}
	mixin.singleton = NewSingletonClass(mixin, env.StdSubtypeClass(symbol.Mixin))

	return mixin
}

func NewMixinWithDetails(
	docComment string,
	abstract bool,
	name string,
	parent Namespace,
	consts ConstantMap,
	subtypes ConstantMap,
	methods MethodMap,
	env *GlobalEnvironment,
) *Mixin {
	mixin := &Mixin{
		parent:   parent,
		abstract: abstract,
		compiled: env.Init,
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
