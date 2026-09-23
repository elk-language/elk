package types

import (
	"fmt"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/value/ivar"
	"github.com/elk-language/elk/value/symbol"
)

type Module struct {
	defined     bool
	native      bool
	parent      Namespace
	ivarIndices *ivar.IvarIndices
	NamespaceBase
}

func (m *Module) ToRef() Ref[*Module] {
	return Ref[*Module](m.id)
}

func (m *Module) EqualAny(other any) bool {
	o, ok := other.(*Module)
	if !ok {
		return false
	}

	if m.id > 0 {
		return m.id == o.ID()
	}

	return m.name == o.name
}

func (m *Module) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(m, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(m, parent)
	}
}

func (m *Module) IvarIndices() *ivar.IvarIndices {
	return m.ivarIndices
}

func (m *Module) SetIvarIndices(in *ivar.IvarIndices) {
	m.ivarIndices = in
}

func (m *Module) IsGeneric() bool {
	return false
}

func (m *Module) TypeParameters() []Ref[*TypeParameter] {
	return nil
}

func (m *Module) SetTypeParameters(t []Ref[*TypeParameter]) {
	panic("cannot set type parameters on a module")
}

func (*Module) Singleton() *SingletonClass {
	return nil
}

func (*Module) SingletonRef() Ref[*SingletonClass] {
	return 0
}

func (m *Module) SetSingleton(*SingletonClass) {
	panic(fmt.Sprintf("cannot set singleton class of a module: %s", m.Name()))
}

func (m *Module) Parent() Namespace {
	return m.parent
}

func (m *Module) ParentRef() Ref[Namespace] {
	if m.parent == nil {
		return 0
	}
	return ToRef(m.parent)
}

func (m *Module) SetParent(parent Namespace) {
	m.parent = parent
}

func (m *Module) IsDefined() bool {
	return m.defined
}

func (m *Module) SetDefined(defined bool) {
	m.defined = defined
}

func (m *Module) IsNative() bool {
	return m.native
}

func (m *Module) SetNative(native bool) {
	m.native = native
}

func (m *Module) IsAbstract() bool {
	return false
}

func (m *Module) IsSealed() bool {
	return false
}

func (m *Module) IsPrimitive() bool {
	return false
}

func (m *Module) IsImmutable() bool {
	return false
}

func NewModule(docComment, name string) *Module {
	return &Module{
		native:        Env.Init,
		parent:        Env.StdSubtypeClass(symbol.C_Module),
		NamespaceBase: MakeNamespaceBase(docComment, name),
	}
}

func NewModuleWithDetails(
	docComment string,
	name string,
	consts ConstantMap,
	subtypes ConstantMap,
	methods MethodMap,
) *Module {
	return &Module{
		parent: Env.StdSubtypeClass(symbol.C_Module),
		native: Env.Init,
		NamespaceBase: NamespaceBase{
			docComment: docComment,
			name:       name,
			constants:  consts,
			subtypes:   subtypes,
			methods:    methods,
		},
	}
}

func (m *Module) ToNonLiteral() Type {
	return m
}

func (*Module) IsLiteral() bool {
	return false
}

func (m *Module) inspect() string {
	return m.Name()
}

func (m *Module) DefineMethod(docComment string, flags bitfield.BitFlag16, name symbol.Symbol, typeParams []Ref[*TypeParameter], params []Ref[*Parameter], returnType, throwType Type) *Method {
	method := NewMethod(docComment, flags, name, typeParams, params, returnType, throwType, m)
	m.SetMethod(name, method)
	return method
}

func (m *Module) Copy() *Module {
	return &Module{
		parent:  m.parent,
		defined: m.defined,
		native:  m.native,
		NamespaceBase: NamespaceBase{
			docComment: m.docComment,
			name:       m.name,
			constants:  m.constants,
			subtypes:   m.subtypes,
			methods:    m.methods,
		},
	}
}

func (m *Module) CopyType() Type {
	return m.Copy()
}

func (m *Module) RemoveTemporaryParents() {
	if _, ok := m.parent.(*TemporaryParent); !ok {
		return
	}
	m.parent = nil
}
