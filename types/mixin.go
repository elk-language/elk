package types

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type Mixin struct {
	parent         Namespace
	abstract       bool
	defined        bool
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

func (m *Mixin) SetSingleton(singleton *SingletonClass) {
	m.singleton = singleton
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
	return m.defined
}

func (m *Mixin) SetDefined(compiled bool) {
	m.defined = compiled
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
		defined:       env.Init,
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
		defined:  env.Init,
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

func (m *Mixin) DefineMethod(docComment string, abstract, sealed, native, generator bool, name value.Symbol, typeParams []*TypeParameter, params []*Parameter, returnType, throwType Type) *Method {
	method := NewMethod(docComment, abstract, sealed, native, generator, name, typeParams, params, returnType, throwType, m)
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

func (m *Mixin) Copy() *Mixin {
	return &Mixin{
		parent:         m.parent,
		abstract:       m.abstract,
		defined:        m.defined,
		Checked:        m.Checked,
		typeParameters: m.typeParameters,
		NamespaceBase: NamespaceBase{
			docComment: m.docComment,
			name:       m.name,
			constants:  m.constants,
			methods:    m.methods,
			subtypes:   m.subtypes,
		},
	}
}

func (m *Mixin) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *Mixin {
	mixinConstantPath := GetConstantPath(m.name)
	parentNamespace := DeepCopyNamespacePath(mixinConstantPath[:len(mixinConstantPath)-1], oldEnv, newEnv)

	if newType, ok := NameToTypeOk(m.name, newEnv); ok {
		return newType.(*Mixin)
	}

	newMixin := &Mixin{
		abstract:      m.abstract,
		defined:       m.defined,
		Checked:       m.Checked,
		NamespaceBase: MakeNamespaceBase(m.docComment, m.name),
	}
	parentNamespace.DefineSubtype(value.ToSymbol(mixinConstantPath[len(mixinConstantPath)-1]), newMixin)

	newMixin.methods = MethodsDeepCopyEnv(m.methods, oldEnv, newEnv)
	newMixin.subtypes = ConstantsDeepCopyEnv(m.subtypes, oldEnv, newEnv)
	newMixin.instanceVariables = TypesDeepCopyEnv(m.instanceVariables, oldEnv, newEnv)
	newMixin.constants = ConstantsDeepCopyEnv(m.constants, oldEnv, newEnv)
	newMixin.typeParameters = TypeParametersDeepCopyEnv(m.typeParameters, oldEnv, newEnv)

	if m.parent != nil {
		newMixin.parent = DeepCopyEnv(m.parent, oldEnv, newEnv).(Namespace)
	}
	var singletonParent Namespace
	if m.singleton.parent != nil {
		singletonParent = DeepCopyEnv(m.singleton.parent, oldEnv, newEnv).(Namespace)
	}
	newMixin.singleton = NewSingletonClass(
		newMixin,
		singletonParent,
	)
	return newMixin
}
