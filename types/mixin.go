package types

import (
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/value/symbol"
)

type Mixin struct {
	parent         Ref[Namespace]
	abstract       bool
	defined        bool
	native         bool
	Checked        bool
	singleton      Ref[*SingletonClass]
	typeParameters []Ref[*TypeParameter]
	NamespaceBase
}

func (m *Mixin) ToRef() Ref[*Mixin] {
	return ToRef(m)
}

func (m *Mixin) EqualAny(other any) bool {
	o, ok := other.(*Class)
	if !ok {
		return false
	}

	if m.id > 0 {
		return m.id == o.ID()
	}

	return m.name == o.name
}

func (m *Mixin) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(m, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(m, parent)
	}
}

func (m *Mixin) IsGeneric() bool {
	return len(m.typeParameters) > 0
}

func (m *Mixin) TypeParameters() []Ref[*TypeParameter] {
	return m.typeParameters
}

func (m *Mixin) SetTypeParameters(t []Ref[*TypeParameter]) {
	m.typeParameters = t
}

func IsMixin(typ Type) bool {
	_, ok := typ.(*Mixin)
	return ok
}

func (m *Mixin) Singleton() *SingletonClass {
	return m.singleton.Get()
}

func (m *Mixin) SetSingleton(singleton *SingletonClass) {
	m.singleton = singleton.ToRef()
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

func (m *Mixin) IsNative() bool {
	return m.native
}

func (m *Mixin) SetNative(native bool) {
	m.native = native
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

func (m *Mixin) IsImmutable() bool {
	return false
}

func (m *Mixin) Parent() Namespace {
	return m.parent.Get()
}

func (m *Mixin) SetParent(parent Namespace) {
	m.parent = ToRef(parent)
}

func (m *Mixin) RemoveTemporaryParents() {
	if _, ok := m.parent.Get().(*TemporaryParent); !ok {
		return
	}

	m.parent = ZERO_ID
	m.singleton.Get().parent = CastRef[Namespace](Env.StdSubtypeClass(symbol.C_Mixin))
}

func NewMixin(docComment string, abstract bool, name string) *Mixin {
	mixin := &Mixin{
		abstract:      abstract,
		defined:       Env.Init,
		NamespaceBase: MakeNamespaceBase(docComment, name),
	}
	mixin.singleton = NewSingletonClass(mixin, Env.StdSubtypeClass(symbol.C_Mixin)).ToRef()

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
) *Mixin {
	mixin := &Mixin{
		parent:   ToRef(parent),
		abstract: abstract,
		defined:  Env.Init,
		NamespaceBase: NamespaceBase{
			docComment: docComment,
			name:       name,
			constants:  consts,
			methods:    methods,
			subtypes:   subtypes,
		},
	}
	mixin.singleton = NewSingletonClass(mixin, Env.StdSubtypeClass(symbol.C_Mixin)).ToRef()

	return mixin
}

func (m *Mixin) DefineMethod(docComment string, flags bitfield.BitFlag16, name symbol.Symbol, typeParams []Ref[*TypeParameter], params []Ref[*Parameter], returnType, throwType Type) *Method {
	method := NewMethod(docComment, flags, name, typeParams, params, returnType, throwType, m)
	m.SetMethod(name, method)
	return method
}

func (m *Mixin) inspect() string {
	return m.name
}

func (m *Mixin) ToNonLiteral() Type {
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
			id:         m.id,
		},
	}
}
